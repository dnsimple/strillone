package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aetrion/dnsimple-go/dnsimple/webhook"
	"github.com/bluele/slack"
	"github.com/julienschmidt/httprouter"
)

const what = "dnsimple-slackhooks"
const dnsimpleURL = "https://dnsimple.com"

var (
	httpPort    string
	slackDryRun bool
)

func init() {
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

	fmt.Fprintln(w, fmt.Sprintf(`{"ping":"%v","what":"%s"}`, time.Now().Unix(), what))
}

func (s *Server) Slack(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	slackAlpha, slackBeta, slackGamma := params.ByName("slackAlpha"), params.ByName("slackBeta"), params.ByName("slackGamma")
	log.Printf("%s %s\n", r.Method, strings.Replace(r.URL.RequestURI(), slackGamma, slackGamma[0:6]+"...", 1))

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var err error

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error parsing body: %v\n", err)
	}

	event, err := webhook.Parse(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error parsing event: %v\n", err)
	}

	service := &SlackService{}
	text := Message(service, event)
	eventHeader := event.EventHeader()

	// Send the webhook to Logs
	log.Printf("[event:%v] %s", eventHeader.RequestID, text)

	// Send the webhook to Slack
	slackWebhookURL := fmt.Sprintf("https://hooks.slack.com/services/%s/%s/%s", slackAlpha, slackBeta, slackGamma)
	if slackAlpha != "-" {
		log.Printf("[event:%v] Sending event to slack %v\n", eventHeader.RequestID, slackAlpha+"/"+slackBeta)

		webhook := slack.NewWebHook(slackWebhookURL)
		slackErr := webhook.PostMessage(&slack.WebHookPostPayload{
			Username: "DNSimple",
			IconUrl:  "http://cl.ly/2t0u2Q380N3y/trusty.png",
			Attachments: []*slack.Attachment{
				&slack.Attachment{
					Fallback: text,
					Color:    "good",
					Fields: []*slack.AttachmentField{
						&slack.AttachmentField{
							Title: event.EventName(),
							Value: service.FormatMessage(text),
						},
					},
				},
			},
		})
		if slackErr != nil {
			log.Printf("[event:%v] Error sending to slack: %v\n", eventHeader.RequestID, err)
		}
	}
}
