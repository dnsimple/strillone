package main

import (
	"log"
	"net/http"
	"os"

	xhttp "github.com/dnsimple/strillone/internal/http"
)

var (
	// Program name
	Program = "dnsimple-strillone"

	// Version is replaced at compilation time
	Version string
)

func main() {
	log.Printf("Starting %s/%s\n", Program, Version)

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "4000"
	}

	server := xhttp.NewServer()

	log.Printf("%s listening on %s...\n", Program, httpPort)
	if err := http.ListenAndServe(":"+httpPort, server); err != nil {
		log.Fatal(err.Error())
	}
}
