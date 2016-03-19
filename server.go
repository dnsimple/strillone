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

	text := MexText(event)
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
							Value: text,
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

func MexText(e webhook.Event) (text string) {
	header := e.EventHeader()
	account := header.Account
	prefix := fmt.Sprintf("[%v] %v", MexDURL(account.Display, fmt.Sprintf("/a/%d/account", account.ID)), header.Actor.Pretty)

	switch event := e.(type) {
	case *webhook.ContactEvent:
		contactLink := MexDURL(fmt.Sprintf("%s %s", event.Contact.FirstName, event.Contact.LastName), fmt.Sprintf("/a/%d/contacts/%d", account.ID, event.Contact.ID))
		switch event.Name {
		case "contact.create":
			text = fmt.Sprintf("%s created the contact %s", prefix, contactLink)
		case "contact.update":
			text = fmt.Sprintf("%s update the contact %s", prefix, contactLink)
		case "contact.delete":
			text = fmt.Sprintf("%s deleted the contact %s", prefix, contactLink)
		default:
			text = fmt.Sprintf("%s performed %s", prefix, event.EventName())
		}
	case *webhook.DomainEvent:
		domainLink := MexDURL(event.Domain.Name, fmt.Sprintf("/a/%d/domains/%s", account.ID, event.Domain.Name))
		switch event.Name {
		case "domain.auto_renewal_enable":
			text = fmt.Sprintf("%s enabled auto-renewal for the domain %s", prefix, domainLink)
		case "domain.auto_renewal_disable":
			text = fmt.Sprintf("%s disabled auto-renewal for the domain %s", prefix, domainLink)
		case "domain.create":
			text = fmt.Sprintf("%s created the domain %s", prefix, domainLink)
		case "domain.delete":
			text = fmt.Sprintf("%s deleted the domain %s", prefix, domainLink)
		case "domain.renew":
			text = fmt.Sprintf("%s renewed the domain %s", prefix, domainLink)
		case "domain.resolution_enable":
			text = fmt.Sprintf("%s enabled resolution for the domain %s", prefix, domainLink)
		case "domain.resolution_disable":
			text = fmt.Sprintf("%s disabled resolution for the domain %s", prefix, domainLink)
		case "domain.token_reset":
			text = fmt.Sprintf("%s reset the token for the domain %s", prefix, domainLink)
		default:
			text = fmt.Sprintf("%s performed %s on domain %s", prefix, event.Name, domainLink)
		}
	case *webhook.ZoneRecordEvent:
		zoneRecordDisplay := fmt.Sprintf("%s %s.%s %s", event.ZoneRecord.Type, event.ZoneRecord.Name, event.ZoneRecord.ZoneID, event.ZoneRecord.Content)
		zoneRecordLink := MexDURL(zoneRecordDisplay, fmt.Sprintf("/a/%d/domains/%s/records/%d", account.ID, event.ZoneRecord.ZoneID, event.ZoneRecord.ID))
		switch event.Name {
		case "record.create":
			text = fmt.Sprintf("%s created the record %s", prefix, zoneRecordLink)
		case "record.update":
			text = fmt.Sprintf("%s updated the record %s", prefix, zoneRecordLink)
		case "record.delete":
			text = fmt.Sprintf("%s deleted the record %s", prefix, zoneRecordLink)
		}
	case *webhook.WebhookEvent:
		webhookLink := MexDURL(event.Webhook.URL, fmt.Sprintf("/a/%d/webhooks/%d", account.ID, event.Webhook.ID))
		switch event.Name {
		case "webhook.create":
			text = fmt.Sprintf("%s created the webhook %s", prefix, webhookLink)
		case "webhook.delete":
			text = fmt.Sprintf("%s deleted the webhook %s", prefix, webhookLink)
		}
	default:
		text = fmt.Sprintf("%s performed %s", prefix, event.EventName())
	}

	return
}

func MexDURL(name, url string) string {
	return MexURL(name, dnsimpleURL+url)
}

func MexURL(name, url string) string {
	return fmt.Sprintf("<%s|%s>", url, name)
}
