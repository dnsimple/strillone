package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple/webhook"
	"github.com/dnsimple/strillone/internal/config"
	"github.com/dnsimple/strillone/internal/service"
	"github.com/wunderlist/ttlcache"
)

const (
	cacheTTL               = 300
	HeaderProcessingStatus = "X-Processing-Status"
)

// Server represents a front-end web server.
type Server struct {
	mux          *http.ServeMux
	webhookCache *ttlcache.Cache
}

// NewServer returns a new front-end web server that handles HTTP requests for the app.
func NewServer() *Server {
	cache := ttlcache.NewCache(cacheTTL * time.Second)

	mux := http.NewServeMux()
	server := &Server{
		mux:          mux,
		webhookCache: cache,
	}

	mux.Handle("GET /", http.HandlerFunc(server.Root))
	mux.Handle("POST /slack/{slackAlpha}/{slackBeta}/{slackGamma}", http.HandlerFunc(server.Slack))
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

	fmt.Fprintf(w, `{"ping":"%v","what":"%s"}`, time.Now().Unix(), config.Program)
}

// Slack handles a request to publish a webhook to a Slack channel.
func (s *Server) Slack(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL.RequestURI())

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error parsing body: %v\n", err)
		return
	}

	event, err := webhook.ParseEvent(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error parsing event: %v\n", err)
		return
	}

	// Check if the event was already processed
	_, cacheExists := s.webhookCache.Get(event.RequestID)
	if cacheExists {
		log.Printf("Skipping event %v as already processed\n", event.RequestID)
		w.Header().Set(HeaderProcessingStatus, "skipped;already-processed")
		w.WriteHeader(http.StatusOK)
		return
	}

	slackAlpha := r.PathValue("slackAlpha")
	slackBeta := r.PathValue("slackBeta")
	slackGamma := r.PathValue("slackGamma")
	slackToken := fmt.Sprintf("%s/%s/%s", slackAlpha, slackBeta, slackGamma)

	service := &service.SlackService{Token: slackToken}
	text, err := service.PostEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Internal Error: %v\n", err)
		return
	}

	s.webhookCache.Set(event.RequestID, "1")

	fmt.Fprintln(w, text)
}
