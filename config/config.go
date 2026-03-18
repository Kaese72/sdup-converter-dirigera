package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// Config is the top level config structure
type Config struct {
	Adapter struct {
		Dirigera struct {
			Host            string `mapstructure:"host"`
			Port            int    `mapstructure:"port"`
			Token           string `mapstructure:"token"`
			IgnoreTLSErrors bool   `mapstructure:"ignore-tls-errors"`
		} `mapstructure:"dirigera"`
	} `mapstructure:"adapter"`
}

func ReadConfig() (Config, error) {
	myVip := viper.New()
	// Set replaces to allow keys like "adapter.dirigera.token"
	// WARNING. Overriding any of these may have unintended consequences.
	myVip.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	// Hostname or IP to Dirigera hub (no scheme), eg. 192.168.1.10
	// ADAPTER_DIRIGERA_HOST
	myVip.BindEnv("adapter.dirigera.host")
	// Port of the Dirigera hub API (default 8443)
	// ADAPTER_DIRIGERA_PORT
	myVip.SetDefault("adapter.dirigera.port", 8443)
	myVip.BindEnv("adapter.dirigera.port")
	// Whether to ignore TLS errors when connecting to Dirigera hub
	// ADAPTER_DIRIGERA_IGNORE_TLS_ERRORS
	myVip.SetDefault("adapter.dirigera.ignore-tls-errors", true)
	myVip.BindEnv("adapter.dirigera.ignore-tls-errors")
	// Bearer token for Dirigera API
	// ADAPTER_DIRIGERA_TOKEN
	myVip.BindEnv("adapter.dirigera.token")

	var conf Config
	err := myVip.Unmarshal(&conf)
	if err != nil {
		return Config{}, err
	}
	if conf.Adapter.Dirigera.Token == "" {
		return Config{}, errors.New("must provide adapter.dirigera.token")
	}
	if conf.Adapter.Dirigera.Host == "" {
		return Config{}, errors.New("must provide adapter.dirigera.host")
	}
	if strings.Contains(conf.Adapter.Dirigera.Host, "://") {
		return Config{}, errors.New("adapter.dirigera.host must be a hostname or IP, not a URL")
	}
	return conf, nil
}
