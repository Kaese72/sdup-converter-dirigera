package sdupdirifera

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-converter-dirigera/dirigeraclient"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/mitchellh/mapstructure"
)

type setExtendedArgs struct {
	Value *float32 `mapstructure:"value"`
}

func init() {
	capRegistry[CapabilitySetExtended] = TriggerSetExtended
}

func TriggerSetExtended(target SDUPDirigeraTarget, id string, args ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	var pArgs setExtendedArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Value == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "value must be set"}
	}
	err := target.client.UpdateDevice(id, dirigeraclient.DeviceAttributes{BlindsTargetLevel: pArgs.Value})
	return adapterErrorFromErr(err)
}
