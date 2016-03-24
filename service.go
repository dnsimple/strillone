package main

import (
	"fmt"
	"log"

	"github.com/aetrion/dnsimple-go/dnsimple/webhook"
	"github.com/bluele/slack"
)

// MessagingService represents a service where the event is published.
// Some examples are Slack, HipChat, and Campfire.
type MessagingService interface {
	FormatLink(name, url string) string
	PostEvent(event webhook.Event) error
}

// SlackService represents the Slack message service.
type SlackService struct {
	Token string
}

// Implements MessagingService
func (s *SlackService) FormatLink(name, url string) string {
	return fmt.Sprintf("<%s|%s>", url, name)
}

// Implements MessagingService
func (s *SlackService) FormatMessage(message string) string {
	return message
}

// Implements MessagingService
func (s *SlackService) PostEvent(event webhook.Event) error {
	eventID := eventRequestID(event)
	text := Message(s, event)

	// Send the webhook to Logs
	log.Printf("[event:%v] %s", eventID, text)

	// Don't send to Slack
	if s.Token[0] == '-' {
		return nil
	}

	slackWebhookURL := fmt.Sprintf("https://hooks.slack.com/services/%s", s.Token)
	log.Printf("[event:%v] Sending event to slack %v\n", eventID, slackWebhookURL)

	webhook := slack.NewWebHook(slackWebhookURL)
	webhookErr := webhook.PostMessage(&slack.WebHookPostPayload{
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
	if webhookErr != nil {
		log.Printf("[event:%v] Error sending to slack: %v\n", eventID, webhookErr)
	}

	return webhookErr
}
