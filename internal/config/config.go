package config

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/dnsimple/strillone/internal/logging"
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
	WebServerHost string `env:"WEB_SERVER_HOST"` // Defaults to http.Server default.
	WebServerPort string `env:"WEB_SERVER_PORT" envDefault:"4000"`
	DNSimpleURL   string `env:"DNSIMPLE_URL" envDefault:"https://dnsimple.com"`
}

// LoadConfiguration loads environment variables into a Configuration struct.
// If parsing fails, it logs a fatal error and exits the program.
func LoadConfiguration() *Configuration {
	cfg := &Configuration{}
	if err := env.Parse(cfg); err != nil {
		slog.Error("Cannot parse environment configuration", logging.Err(err))
		os.Exit(1)
	}

	return cfg
}
