package config_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/dnsimple/strillone/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDefaults(t *testing.T) {
	cfg, err := config.Parse(map[string]string{})

	require.NoError(t, err)
	assert.Equal(t, ":4000", cfg.Addr())
	assert.Equal(t, "https://app.dnsimple.com", cfg.DNSimpleURL)
	assert.Equal(t, "https://hooks.slack.com/services", cfg.SlackURL)
}

func TestParseDNSimpleURL(t *testing.T) {
	t.Run("defaults to the application host when empty", func(t *testing.T) {
		cfg, err := config.Parse(map[string]string{"DNSIMPLE_URL": ""})

		require.NoError(t, err)
		assert.Equal(t, "https://app.dnsimple.com", cfg.DNSimpleURL)
	})

	t.Run("uses an explicit override", func(t *testing.T) {
		cfg, err := config.Parse(map[string]string{"DNSIMPLE_URL": "https://app.example.test"})

		require.NoError(t, err)
		assert.Equal(t, "https://app.example.test", cfg.DNSimpleURL)
	})
}

func TestParseListenAddress(t *testing.T) {
	cfg, err := config.Parse(map[string]string{"BIND_ADDRESS": "127.0.0.1", "PORT": "5000"})

	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1:5000", cfg.Addr())
}

func TestParseRefusesInvalidPort(t *testing.T) {
	for name, port := range map[string]string{
		"zero":      "0",
		"too large": "65536",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := config.Parse(map[string]string{"PORT": port})

			require.ErrorContains(t, err, "PORT")
		})
	}

	t.Run("not a number", func(t *testing.T) {
		_, err := config.Parse(map[string]string{"PORT": "http"})

		require.Error(t, err)
	})
}

func TestEveryVariableIsDocumented(t *testing.T) {
	params, err := env.GetFieldParams(&config.Configuration{})
	require.NoError(t, err)

	readme := readFile(t, "README.md")
	example := readFile(t, ".env.example")

	for _, param := range params {
		key := regexp.QuoteMeta(param.Key)
		assert.Regexp(t, `(?m)^\|\s*`+key+`\s*\|`, readme, "README.md does not document %s", param.Key)
		assert.Regexp(t, `(?m)^#?\s*`+key+`=`, example, ".env.example does not list %s", param.Key)
	}
}

func readFile(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("..", "..", name))
	require.NoError(t, err)

	return string(content)
}
