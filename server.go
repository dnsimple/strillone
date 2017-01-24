package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aetrion/dnsimple-go/dnsimple/webhook"
	"github.com/julienschmidt/httprouter"
	"github.com/wunderlist/ttlcache"
)

// What represents the program name
var What = "dnsimple-strillone"

// Version is replaced at compilation time
var Version string

const (
	dnsimpleURL = "https://dnsimple.com"
	cacheTTL    = 300
)

var (
	httpPort        string
	processedEvents *ttlcache.Cache
)

func init() {
	httpPort = os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "5000"
	}

	processedEvents = ttlcache.NewCache(time.Second * cacheTTL)
}

func main() {
	log.Printf("Starting %s %s\n", What, Version)

	server := NewServer()

	log.Printf("%s listening on %s...\n", What, httpPort)
	if err := http.ListenAndServe(":"+httpPort, server); err != nil {
		log.Panic(err)
	}
}

// Server represents a front-end web server.
type Server struct {
	// Router which handles incoming requests
	mux *httprouter.Router
}

// NewServer returns a new front-end web server that handles HTTP requests for the app.
func NewServer() *Server {
	router := httprouter.New()
	server := &Server{mux: router}
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

	fmt.Fprintln(w, fmt.Sprintf(`{"ping":"%v","what":"%s"}`, time.Now().Unix(), What))
}

// Slack handles a request to publish a webhook to a Slack channel.
func (s *Server) Slack(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Printf("%s %s\n", r.Method, r.URL.RequestURI())

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error parsing body: %v\n", err)
		return
	}

	event, err := webhook.Parse(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error parsing event: %v\n", err)
		return
	}

	// Check if the event was already processed
	_, cacheExists := processedEvents.Get(event.EventHeader().RequestID)
	if cacheExists {
		log.Printf("Event %v already processed, skipping\n", event.EventHeader().RequestID)
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

	processedEvents.Set(event.EventHeader().RequestID, event.EventHeader().RequestID)

	fmt.Fprintln(w, text)
}
