package sdupdirifera

import (
	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/Kaese72/sdup-converter-dirigera/dirigeraclient"
)

func init() {
	capRegistry[CapabilityActivate] = TriggerTurnOn
}

func TriggerTurnOn(target SDUPDirigeraTarget, id string, _ ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	on := true
	err := target.client.UpdateDevice(id, dirigeraclient.DeviceAttributes{IsOn: &on})
	return adapterErrorFromErr(err)
}
