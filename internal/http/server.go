package http

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple/webhook"
	"github.com/dnsimple/strillone/internal/config"
	"github.com/dnsimple/strillone/internal/logging"
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
	dnsimpleURL  string
}

// NewServer returns a new front-end web server that handles HTTP requests for the app.
func NewServer(dnsimpleURL string) *Server {
	cache := ttlcache.NewCache(cacheTTL * time.Second)

	mux := http.NewServeMux()
	server := &Server{
		mux:          mux,
		webhookCache: cache,
		dnsimpleURL:  dnsimpleURL,
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
	slog.Info("Received request", "http_method", r.Method, "http_url", r.URL.RequestURI())
	w.Header().Set("Content-type", "application/json")

	fmt.Fprintf(w, `{"ping":"%v","what":"%s"}`, time.Now().Unix(), config.Program)
}

// Slack handles a request to publish a webhook to a Slack channel.
func (s *Server) Slack(w http.ResponseWriter, r *http.Request) {
	// The URL path contains the Slack token.
	slog.Info("Received request", "http_method", r.Method)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Warn("Error parsing body", logging.Err(err))
		return
	}

	event, err := webhook.ParseEvent(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Warn("Error parsing event", logging.Err(err))
		return
	}

	logger := slog.With("request_id", event.RequestID)

	// Check if the event was already processed
	_, cacheExists := s.webhookCache.Get(event.RequestID)
	if cacheExists {
		logger.Info("Skipping event, already processed")
		w.Header().Set(HeaderProcessingStatus, "skipped;already-processed")
		w.WriteHeader(http.StatusOK)
		return
	}

	slackAlpha := r.PathValue("slackAlpha")
	slackBeta := r.PathValue("slackBeta")
	slackGamma := r.PathValue("slackGamma")
	slackToken := fmt.Sprintf("%s/%s/%s", slackAlpha, slackBeta, slackGamma)

	service := &service.SlackService{Token: slackToken, DNSimpleURL: s.dnsimpleURL}
	text, err := service.PostEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("Error sending to slack", logging.Err(err))
		return
	}

	s.webhookCache.Set(event.RequestID, "1")

	fmt.Fprintln(w, text)
}
