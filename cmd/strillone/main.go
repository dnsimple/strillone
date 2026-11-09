package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/dnsimple/strillone/internal/config"
	xhttp "github.com/dnsimple/strillone/internal/http"
	"github.com/dnsimple/strillone/internal/logging"
	"github.com/dnsimple/strillone/internal/version"
)

func main() {
	logging.SetupDefault()

	slog.Info("Starting", "program", version.Program, "version", version.Version)
	if err := run(); err != nil {
		slog.Error("strillone failed", logging.Err(err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	server := xhttp.NewServer(cfg.DNSimpleURL, cfg.SlackURL)

	addr := cfg.Addr()
	slog.Info("WebServer listening", "addr", addr)
	return http.ListenAndServe(addr, server)
}
