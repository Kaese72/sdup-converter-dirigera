package sdupdirifera

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
)

// capRegistry contains device capability functions keyed by capability name
var capRegistry = map[string]func(target SDUPDirigeraTarget, id string, _ ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError{}

func adapterErrorFromErr(err error) *adapter.AdapterError {
	if err == nil {
		return nil
	}
	return &adapter.AdapterError{Code: http.StatusInternalServerError, Message: err.Error()}
}
