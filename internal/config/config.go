// Package config reads the configuration of the application from the environment.
package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/caarlos0/env/v11"
)

// Configuration holds the settings of the application. README.md documents each variable.
type Configuration struct {
	BindAddress string `env:"BIND_ADDRESS"` // Empty listens on all interfaces.
	Port        int    `env:"PORT" envDefault:"4000"`
	DNSimpleURL string `env:"DNSIMPLE_URL" envDefault:"https://app.dnsimple.com"`
	SlackURL    string `env:"SLACK_URL" envDefault:"https://hooks.slack.com/services"`
}

// Load reads the configuration from the process environment.
func Load() (*Configuration, error) {
	return Parse(env.ToMap(os.Environ()))
}

// Parse reads the configuration from environ in place of the process environment.
func Parse(environ map[string]string) (*Configuration, error) {
	cfg := &Configuration{}
	if err := env.ParseWithOptions(cfg, env.Options{Environment: environ}); err != nil {
		return nil, fmt.Errorf("parsing the environment: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Addr returns the listen address of the HTTP server.
func (c *Configuration) Addr() string {
	return net.JoinHostPort(c.BindAddress, strconv.Itoa(c.Port))
}

func (c *Configuration) validate() error {
	var errs []error

	if c.Port < 1 || c.Port > 65535 {
		errs = append(errs, errors.New("PORT must be between 1 and 65535"))
	}

	return errors.Join(errs...)
}
