package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const what = "dnsimple-slackhooks"

var (
	httpPort        string
	slackWebhookURL string
)

func init() {
	// for now read the URL from the ENV.
	// in the future we may probably want to be able to provide a flexible configuration.
	slackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhookURL == "" {
		log.Fatalln("Slack Webhook URL is missing")
	}

	httpPort = os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "5000"
	}
}

func main() {
	log.Printf("Starting %s...\n", what)

	server := NewServer()

	log.Printf("%s listening on %s...\n", what, httpPort)
	if err := http.ListenAndServe(":"+httpPort, server); err != nil {
		log.Panic(err)
	}
}

// Server represents a front-end web server.
type Server struct {
	// Router which handles incoming requests
	mux *http.ServeMux
}

// NewServer returns a new front-end web server that handles HTTP requests for the app.
func NewServer() *Server {
	router := http.NewServeMux()
	server := &Server{mux: router}
	router.HandleFunc("/", server.Root)
	return server
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Root is the handler for the HTTP requests to /.
// It returns a simple uptime message useful for monitoring.
func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL.RequestURI())
	w.Header().Set("Content-type", "application/json")

	fmt.Fprintln(w, fmt.Sprintf(`{"ping":"%v","what":%s}`, time.Now().Unix(), what))
}
