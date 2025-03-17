package config

import (
	"github.com/caarlos0/env/v11"
)

var (
	// Config is the global configuration
	Config *AppConfig

	// Program name
	Program = "dnsimple-strillone"
)

// AppConfig represents the configuration of the application.
type AppConfig struct {
	// Port is the HTTP port the server listens on.
	Port string `env:"PORT" envDefault:"4000"`
	// DNSimpleUrl is the DNSimple app URL.
	DNSimpleUrl string `env:"DNSIMPLE_URL" envDefault:"https://dnsimple.com"`
}

// NewConfig returns a new AppConfig instance.
func NewConfig() (*AppConfig, error) {
	cfg := &AppConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
