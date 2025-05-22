package config

import (
	"github.com/caarlos0/env/v11"
)

var (
	// Config is the global configuration.
	Config *Configuration

	// Program name.
	Program = "dnsimple-strillone"

	// Version is replaced at compilation time.
	Version string
)

// Configuration holds all the environment-based configuration settings for the application.
type Configuration struct {
	// Port is the HTTP port the server listens on.
	Port string `env:"PORT" envDefault:"4000"`
	// DNSimpleURL is the DNSimple app URL.
	DNSimpleURL string `env:"DNSIMPLE_URL" envDefault:"https://dnsimple.com"`
}

// NewConfig returns a new Configuration instance.
func NewConfig() (*Configuration, error) {
	cfg := &Configuration{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
