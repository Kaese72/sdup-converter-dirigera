package dirigeraclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Device represents a Dirigera device
type Device struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes DeviceAttributes `json:"attributes"`
	Room       *Room            `json:"room,omitempty"`
}

// DeviceAttributes holds the state attributes of a device
type DeviceAttributes struct {
	IsOn       *bool  `json:"isOn,omitempty"`
	CustomName string `json:"customName,omitempty"`
}

// Room represents the room a device belongs to
type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DevicePatchItem is a single patch operation for a device
type DevicePatchItem struct {
	Attributes DeviceAttributes `json:"attributes"`
}

// Client is a HTTP client for the Dirigera local API
type Client struct {
	host   string
	port   int
	token  string
	client *http.Client
}

// New creates a new Dirigera client
func New(host, token string, port int, ignoreTLSErrors bool) *Client {
	transport := &http.Transport{}
	if ignoreTLSErrors {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &Client{
		host:   host,
		port:   port,
		token:  token,
		client: &http.Client{Transport: transport},
	}
}

func (c *Client) baseURL() string {
	return fmt.Sprintf("https://%s:%d/v1", c.host, c.port)
}

func (c *Client) newRequest(method, path string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, c.baseURL()+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("dirigera API error %d: %s", resp.StatusCode, string(data))
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// GetDevices returns all devices from the Dirigera hub
func (c *Client) GetDevices() ([]Device, error) {
	req, err := c.newRequest(http.MethodGet, "/devices", nil)
	if err != nil {
		return nil, err
	}
	var devices []Device
	if err := c.do(req, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

// UpdateDevice patches the attributes of a device
func (c *Client) UpdateDevice(id string, attributes DeviceAttributes) error {
	body := []DevicePatchItem{{Attributes: attributes}}
	req, err := c.newRequest(http.MethodPatch, "/devices/"+id, body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
