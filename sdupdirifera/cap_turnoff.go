package sdupdirifera

import (
	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/Kaese72/sdup-converter-dirigera/dirigeraclient"
)

func init() {
	capRegistry[CapabilityDeactivate] = TriggerTurnOff
}

func TriggerTurnOff(target SDUPDirigeraTarget, id string, _ ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	off := false
	err := target.client.UpdateDevice(id, dirigeraclient.DeviceAttributes{IsOn: &off})
	return adapterErrorFromErr(err)
}
