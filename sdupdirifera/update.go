package sdupdirifera

import (
	"strings"

	"github.com/Kaese72/device-store/ingestmodels"
	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-converter-dirigera/dirigeraclient"
)

const (
	// AttributeActive represents whether the device is currently on or off
	AttributeActive = "active"
)

const (
	// CapabilityActivate turns the device on
	CapabilityActivate = "activate"
	// CapabilityDeactivate turns the device off
	CapabilityDeactivate = "deactivate"
)

// DeviceTypeLightbulb is the Dirigera device type identifier for lights
const DeviceTypeLightbulb = "lightbulb"

// DeviceTypeOutlet is the Dirigera device type identifier for outlets
const DeviceTypeOutlet = "outlet"

func (target SDUPDirigeraTarget) getAllDevices() ([]ingestmodels.IngestDevice, error) {
	allDevices, err := target.client.GetDevices()
	if err != nil {
		return nil, err
	}
	var specs []ingestmodels.IngestDevice
	for _, d := range allDevices {
		switch {
		case strings.EqualFold(d.Type, DeviceTypeLightbulb):
			specs = append(specs, createLightbulbDevice(d))
		case strings.EqualFold(d.Type, DeviceTypeOutlet):
			specs = append(specs, createOutletDevice(d))
		default:
			log.Info("Ignoring device with unrecognized type", map[string]interface{}{"type": d.Type, "id": d.ID})
		}
	}
	return specs, nil
}

func createLightbulbDevice(d dirigeraclient.Device) ingestmodels.IngestDevice {
	device := ingestmodels.IngestDevice{
		BridgeIdentifier: d.ID,
		Attributes:       []ingestmodels.IngestAttribute{},
		Capabilities: []ingestmodels.IngestDeviceCapability{
			{
				Name:          CapabilityActivate,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{},
			},
			{
				Name:          CapabilityDeactivate,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{},
			},
		},
	}

	if d.Attributes.IsOn != nil {
		on := *d.Attributes.IsOn
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name:    AttributeActive,
			Boolean: &on,
		})
	}

	return device
}

func createOutletDevice(d dirigeraclient.Device) ingestmodels.IngestDevice {
	device := ingestmodels.IngestDevice{
		BridgeIdentifier: d.ID,
		Attributes:       []ingestmodels.IngestAttribute{},
		Capabilities: []ingestmodels.IngestDeviceCapability{
			{
				Name:          CapabilityActivate,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{},
			},
			{
				Name:          CapabilityDeactivate,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{},
			},
		},
	}

	if d.Attributes.IsOn != nil {
		on := *d.Attributes.IsOn
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name:    AttributeActive,
			Boolean: &on,
		})
	}

	return device
}
