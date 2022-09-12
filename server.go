package strillone

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"
	"github.com/julienschmidt/httprouter"
	"github.com/wunderlist/ttlcache"
)

const (
	dnsimpleURL            = "https://dnsimple.com"
	cacheTTL               = 300
	HeaderProcessingStatus = "X-Processing-Status"
)

var (
	// Program name
	Program = "dnsimple-strillone"

	// Version is replaced at compilation time
	Version string
)

// Server represents a front-end web server.
type Server struct {
	mux          *httprouter.Router
	webhookCache *ttlcache.Cache
}

// NewServer returns a new front-end web server that handles HTTP requests for the app.
func NewServer() *Server {
	cache := ttlcache.NewCache(cacheTTL * time.Second)

	router := httprouter.New()
	server := &Server{
		mux:          router,
		webhookCache: cache,
	}

	router.GET("/", server.Root)
	router.POST("/slack/:slackAlpha/:slackBeta/:slackGamma", server.Slack)
	return server
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Root is the handler for the HTTP requests to /.
// It returns a simple uptime message useful for monitoring.
func (s *Server) Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("%s %s\n", r.Method, r.URL.RequestURI())
	w.Header().Set("Content-type", "application/json")

	fmt.Fprintf(w, `{"ping":"%v","what":"%s"}`, time.Now().Unix(), Program)
}

// Slack handles a request to publish a webhook to a Slack channel.
func (s *Server) Slack(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

	slackAlpha, slackBeta, slackGamma := params.ByName("slackAlpha"), params.ByName("slackBeta"), params.ByName("slackGamma")
	slackToken := fmt.Sprintf("%s/%s/%s", slackAlpha, slackBeta, slackGamma)

	service := &SlackService{Token: slackToken}
	text, err := service.PostEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Internal Error: %v\n", err)
		return
	}

	s.webhookCache.Set(event.RequestID, "1")

	fmt.Fprintln(w, text)
}
