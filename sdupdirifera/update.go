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
	// AttributeName represents the human-readable name of the device
	AttributeName = "description"
	// AttributeExtended represents how far a blind is extended (0-100)
	AttributeExtended = "extended"
)

const (
	// CapabilityActivate turns the device on
	CapabilityActivate = "activate"
	// CapabilityDeactivate turns the device off
	CapabilityDeactivate = "deactivate"
	// CapabilitySetExtended sets how far a blind should be extended (0-100)
	CapabilitySetExtended = "extend"
)

// DeviceTypeLightbulb is the Dirigera device type identifier for lights
const DeviceTypeLightbulb = "lightbulb"

// DeviceTypeOutlet is the Dirigera device type identifier for outlets
const DeviceTypeOutlet = "outlet"

// DeviceTypeBlinds is the Dirigera device type identifier for blinds
const DeviceTypeBlinds = "blinds"

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
		case strings.EqualFold(d.Type, DeviceTypeBlinds):
			specs = append(specs, createBlindsDevice(d))
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

	if d.Attributes.CustomName != "" {
		name := d.Attributes.CustomName
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name: AttributeName,
			Text: &name,
		})
	}

	if d.Attributes.CustomName != "" {
		name := d.Attributes.CustomName
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name: AttributeName,
			Text: &name,
		})
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

func createBlindsDevice(d dirigeraclient.Device) ingestmodels.IngestDevice {
	device := ingestmodels.IngestDevice{
		BridgeIdentifier: d.ID,
		Attributes:       []ingestmodels.IngestAttribute{},
		Capabilities: []ingestmodels.IngestDeviceCapability{
			{
				Name: CapabilitySetExtended,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{
					{
						Name:    "value",
						Numeric: &ingestmodels.IngestNumericArgumentSpec{Min: 0, Max: 100},
					},
				},
			},
		},
	}

	if d.Attributes.CustomName != "" {
		name := d.Attributes.CustomName
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name: AttributeName,
			Text: &name,
		})
	}

	if d.Attributes.BlindsCurrentLevel != nil {
		level := *d.Attributes.BlindsCurrentLevel
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name:    AttributeExtended,
			Numeric: &level,
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

	if d.Attributes.CustomName != "" {
		name := d.Attributes.CustomName
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name: AttributeName,
			Text: &name,
		})
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
