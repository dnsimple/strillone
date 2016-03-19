package main

import (
	"fmt"
)

type MessagingService interface {
	FormatLink(name, url string) string
	FormatMessage(message string) string
}

// SlackService represents the Slack message service.
type SlackService struct {
}

// Implements MessagingService
func (s *SlackService) FormatLink(name, url string) string {
	return fmt.Sprintf("<%s|%s>", url, name)
}

// Implements MessagingService
func (s *SlackService) FormatMessage(message string) string {
	return message
}
