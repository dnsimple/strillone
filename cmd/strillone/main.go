package main

import (
	"log"
	"net/http"

	"github.com/dnsimple/strillone/internal/config"
	xhttp "github.com/dnsimple/strillone/internal/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	config.Config = cfg

	server := xhttp.NewServer()

	log.Printf("%s listening on %s...\n", config.Program, cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, server); err != nil {
		log.Fatal(err.Error())
	}
}
