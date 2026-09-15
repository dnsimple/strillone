package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/dnsimple/strillone/internal/config"
	xhttp "github.com/dnsimple/strillone/internal/http"
	"github.com/dnsimple/strillone/internal/logging"
)

func main() {
	logging.SetupDefault()
	cfg := config.LoadConfiguration()

	slog.Info("Starting", "program", config.Program, "version", config.Version)
	if err := run(cfg); err != nil {
		slog.Error("strillone failed", logging.Err(err))
		os.Exit(1)
	}
}

func run(cfg *config.Configuration) error {
	server := xhttp.NewServer(cfg.DNSimpleURL)

	addr := cfg.WebServerHost + ":" + cfg.WebServerPort
	slog.Info("WebServer listening", "addr", addr)
	return http.ListenAndServe(addr, server)
}
