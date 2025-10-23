package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple/webhook"
	"github.com/slack-go/slack"
)

// MessagingService represents a service where the event is published.
// Some examples are Slack, HipChat, and Campfire.
type MessagingService interface {
	FormatLink(name, url string) string
	PostEvent(event *webhook.Event) (string, error)
}

// SlackService represents the Slack message service.
type SlackService struct {
	Token string
}

// FormatLink implements MessagingService
func (s *SlackService) FormatLink(name, url string) string {
	return fmt.Sprintf("<%s|%s>", url, name)
}

// FormatMessage implements MessagingService
func (s *SlackService) FormatMessage(message string) string {
	return message
}

// PostEvent implements MessagingService
func (s *SlackService) PostEvent(event *webhook.Event) (string, error) {
	eventID := eventRequestID(event)
	text := Message(s, event)

	// Send the webhook to Logs
	log.Printf("[event:%v] %s", eventID, text)

	// Don't send to Slack
	if s.Token[0] == '-' {
		return "", nil
	}

	slackWebhookURL := fmt.Sprintf("https://hooks.slack.com/services/%s", s.Token)
	log.Printf("[event:%v] Sending event to slack %v\n", eventID, slackWebhookURL)

	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      text,
		AuthorName:    "DNSimple",
		AuthorSubname: "Strillone",
		AuthorLink:    "https://github.com/dnsimple/strillone",
		AuthorIcon:    "https://cdn.dnsimple.com/assets/strillone/icon128.png",
		Title:         event.Name,
		Text:          text,
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(slackWebhookURL, &msg)
	if err != nil {
		log.Printf("[event:%v] Error sending to slack: %v\n", eventID, err)
		return "", err
	}

	return text, nil
}
