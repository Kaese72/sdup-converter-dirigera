package main

import (
	"os"

	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-converter-dirigera/config"
	"github.com/Kaese72/sdup-converter-dirigera/sdupdirifera"
	"github.com/Kaese72/sdup-lib/adapter"
)

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		log.Error("Failed to read config", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
	dirigeraTarget, err := sdupdirifera.InitSDUPDirigeraTarget(
		conf.Adapter.Dirigera.Host,
		conf.Adapter.Dirigera.Token,
		conf.Adapter.Dirigera.Port,
		conf.Adapter.Dirigera.IgnoreTLSErrors,
	)
	if err != nil {
		log.Error("Failed to initialize Dirigera target", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
	// Make sure we fulfil the trigger interface
	var _ adapter.DeviceTriggerCapability = dirigeraTarget
	err = adapter.StartAdapter(dirigeraTarget)
	if err != nil {
		log.Error("Failed to start adapter", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
}
