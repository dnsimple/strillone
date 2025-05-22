package config

import (
	"log"

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

func init() {
	Config = LoadConfiguration()
}

// Configuration holds all the environment-based configuration settings for the application.
type Configuration struct {
	// Port is the HTTP port the server listens on.
	Port string `env:"PORT" envDefault:"4000"`
	// DNSimpleURL is the DNSimple app URL.
	DNSimpleURL string `env:"DNSIMPLE_URL" envDefault:"https://dnsimple.com"`
}

// LoadConfiguration loads environment variables into a Configuration struct.
// If parsing fails, it logs a fatal error and exits the program.
func LoadConfiguration() *Configuration {
	cfg := &Configuration{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Cannot parse environment configuration")
	}

	return cfg
}
