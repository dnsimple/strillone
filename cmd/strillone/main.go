package main

import (
	"log"
	"net/http"

	"github.com/dnsimple/strillone/internal/config"
	xhttp "github.com/dnsimple/strillone/internal/http"
)

func main() {
	server := xhttp.NewServer()

	log.Printf("Starting %s/%s", config.Program, config.Version)

	addr := config.Config.WebServerHost + ":" + config.Config.WebServerPort
	log.Printf("WebServer listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal(err)
	}
}
