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
	slog.SetDefault(logging.New(logging.ParseLevel(os.Getenv("LOG_LEVEL"))))
	config.Config = config.LoadConfiguration()

	if err := run(); err != nil {
		slog.Error("strillone failed", logging.Err(err))
		os.Exit(1)
	}
}

func run() error {
	server := xhttp.NewServer()

	slog.Info("Starting", "program", config.Program, "version", config.Version)

	addr := config.Config.WebServerHost + ":" + config.Config.WebServerPort
	slog.Info("WebServer listening", "address", addr)
	return http.ListenAndServe(addr, server)
}
