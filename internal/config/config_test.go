package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigurationDNSimpleURL(t *testing.T) {
	t.Run("defaults to the application host", func(t *testing.T) {
		t.Setenv("DNSIMPLE_URL", "")

		cfg := LoadConfiguration()

		assert.Equal(t, "https://app.dnsimple.com", cfg.DNSimpleURL)
	})

	t.Run("uses an explicit override", func(t *testing.T) {
		t.Setenv("DNSIMPLE_URL", "https://app.example.test")

		cfg := LoadConfiguration()

		assert.Equal(t, "https://app.example.test", cfg.DNSimpleURL)
	})
}
