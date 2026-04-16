package sdupdirifera

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Kaese72/device-store/ingestmodels"
	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/Kaese72/sdup-converter-dirigera/dirigeraclient"
	"github.com/gorilla/websocket"
)

const (
	// pingInterval is how often we send a WebSocket ping to keep the connection alive.
	pingInterval = 30 * time.Second
	// pongWait is how long we wait for any traffic (pong or message) before
	// treating the connection as dead.
	pongWait = 60 * time.Second
)

// DirigeraEvent represents an event from the Dirigera WebSocket stream
type DirigeraEvent struct {
	Type string            `json:"type"`
	Data DirigeraEventData `json:"data"`
}

// DirigeraEventData holds the device state from an event
type DirigeraEventData struct {
	ID         string                          `json:"id"`
	Type       string                          `json:"type"`
	Attributes dirigeraclient.DeviceAttributes `json:"attributes"`
}

func (target SDUPDirigeraTarget) listenEventStream(updates chan adapter.Update) error {
	wsURL := fmt.Sprintf("wss://%s:%d/v1", target.host, target.port)
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: target.ignoreTLSErrors},
	}
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+target.token)

	conn, _, err := dialer.Dial(wsURL, headers)
	if err != nil {
		return fmt.Errorf("failed to connect to event stream: %w", err)
	}
	defer conn.Close()

	log.Info("Connected to Dirigera event stream", nil)

	// Any traffic (pong or message) resets the read deadline.
	resetDeadline := func() {
		conn.SetReadDeadline(time.Now().Add(pongWait))
	}
	resetDeadline()
	conn.SetPongHandler(func(string) error {
		resetDeadline()
		return nil
	})

	// Ping goroutine — sends a ping every pingInterval.
	// Stops when the connection is closed (WriteMessage returns an error).
	done := make(chan struct{})
	defer close(done)
	go func() {
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			case <-done:
				return
			}
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("event stream read error: %w", err)
		}
		// Any received message is also proof the connection is alive.
		resetDeadline()
		if err := target.handleEvent(msg, updates); err != nil {
			log.Error("Failed to handle event", map[string]interface{}{"error": err.Error()})
		}
	}
}

func (target SDUPDirigeraTarget) handleEvent(msg []byte, updates chan adapter.Update) error {
	var event DirigeraEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	if !strings.EqualFold(event.Type, "deviceStateChanged") {
		log.Debug("Skipping non-device-state event", map[string]interface{}{"type": event.Type})
		return nil
	}

	id := strings.TrimSpace(event.Data.ID)
	if id == "" {
		return nil
	}

	var device *ingestmodels.IngestDevice
	switch {
	case strings.EqualFold(event.Data.Type, DeviceTypeLightbulb):
		device = handleLightbulbEvent(id, event.Data.Attributes)
	case strings.EqualFold(event.Data.Type, DeviceTypeOutlet):
		device = handleOutletEvent(id, event.Data.Attributes)
	case strings.EqualFold(event.Data.Type, DeviceTypeBlinds):
		device = handleBlindsEvent(id, event.Data.Attributes)
	default:
		log.Info("Ignoring event for unrecognized device type", map[string]interface{}{"type": event.Data.Type, "id": id})
		return nil
	}

	if device == nil {
		return nil
	}

	updates <- adapter.Update{Device: device}
	return nil
}

func handleLightbulbEvent(id string, attrs dirigeraclient.DeviceAttributes) *ingestmodels.IngestDevice {
	if attrs.IsOn == nil {
		log.Debug("Ignoring lightbulb event with no isOn attribute", map[string]interface{}{"id": id})
		return nil
	}
	on := *attrs.IsOn
	return &ingestmodels.IngestDevice{
		BridgeIdentifier: id,
		Attributes:       []ingestmodels.IngestAttribute{{Name: AttributeActive, Boolean: &on}},
	}
}

func handleBlindsEvent(id string, attrs dirigeraclient.DeviceAttributes) *ingestmodels.IngestDevice {
	if attrs.BlindsCurrentLevel == nil {
		log.Debug("Ignoring blinds event with no blindsCurrentLevel attribute", map[string]interface{}{"id": id})
		return nil
	}
	level := *attrs.BlindsCurrentLevel
	return &ingestmodels.IngestDevice{
		BridgeIdentifier: id,
		Attributes:       []ingestmodels.IngestAttribute{{Name: AttributeExtended, Numeric: &level}},
	}
}

func handleOutletEvent(id string, attrs dirigeraclient.DeviceAttributes) *ingestmodels.IngestDevice {
	if attrs.IsOn == nil {
		log.Debug("Ignoring outlet event with no isOn attribute", map[string]interface{}{"id": id})
		return nil
	}
	on := *attrs.IsOn
	return &ingestmodels.IngestDevice{
		BridgeIdentifier: id,
		Attributes:       []ingestmodels.IngestAttribute{{Name: AttributeActive, Boolean: &on}},
	}
}
