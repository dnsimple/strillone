package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple/webhook"
	"github.com/dnsimple/strillone/internal/logging"
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
	slog.Info(text, "event_id", eventID)

	// Don't send to Slack
	if s.Token[0] == '-' {
		return "", nil
	}

	slackWebhookURL := fmt.Sprintf("https://hooks.slack.com/services/%s", s.Token)
	slog.Info("Sending event to slack", "event_id", eventID, "webhook_url", slackWebhookURL)

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
		slog.Error("Error sending to slack", "event_id", eventID, logging.Err(err))
		return "", err
	}

	return text, nil
}
