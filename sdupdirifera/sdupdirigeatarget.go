package sdupdirifera

import (
	"strings"
	"time"

	"github.com/Kaese72/device-store/ingestmodels"
	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-converter-dirigera/dirigeraclient"
	"github.com/Kaese72/sdup-lib/adapter"
)

// SDUPDirigeraTarget is the adapter implementation for IKEA Dirigera
type SDUPDirigeraTarget struct {
	client          *dirigeraclient.Client
	host            string
	port            int
	token           string
	ignoreTLSErrors bool
}

// Initialize starts the updater. On each (re)connect it first pushes the full
// device state, then opens a WebSocket event stream to the Dirigera hub and
// listens for changes. On disconnect it waits 5 seconds before reconnecting.
// The returned channel is never closed — the adapter is expected to run indefinitely.
func (target SDUPDirigeraTarget) Initialize() (chan adapter.Update, error) {
	channel := make(chan adapter.Update)
	go func() {
		for {
			if err := target.sendInitialState(channel); err != nil {
				log.Error("Failed to send full device state", map[string]interface{}{"error": err.Error()})
			}
			if err := target.listenEventStream(channel); err != nil {
				log.Error("Event stream disconnected", map[string]interface{}{"error": err.Error()})
			}
			log.Info("Reconnecting to event stream in 5 seconds", nil)
			time.Sleep(5 * time.Second)
		}
	}()
	return channel, nil
}

func (target SDUPDirigeraTarget) sendInitialState(updates chan adapter.Update) error {
	devices, err := target.getAllDevices()
	if err != nil {
		return err
	}
	for _, d := range devices {
		updates <- adapter.Update{Device: &d}
	}
	return nil
}

// DeviceTriggerCapability routes a capability trigger to the appropriate handler
func (target SDUPDirigeraTarget) DeviceTriggerCapability(deviceID string, capabilityKey string, argument ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	capability, ok := capRegistry[capabilityKey]
	if !ok {
		log.Debug("Could not find capability", map[string]interface{}{"device": deviceID, "capability": capabilityKey})
		return &adapter.AdapterError{Code: 404, Message: "No such capability"}
	}
	if strings.TrimSpace(deviceID) == "" {
		return &adapter.AdapterError{Code: 404, Message: "No such device"}
	}
	return capability(target, deviceID, argument)
}

// InitSDUPDirigeraTarget creates and returns a new SDUPDirigeraTarget
func InitSDUPDirigeraTarget(host, token string, port int, ignoreTLSErrors bool) (SDUPDirigeraTarget, error) {
	client := dirigeraclient.New(host, token, port, ignoreTLSErrors)
	return SDUPDirigeraTarget{
		client:          client,
		host:            host,
		port:            port,
		token:           token,
		ignoreTLSErrors: ignoreTLSErrors,
	}, nil
}
