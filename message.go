package main

import (
	"fmt"
	"strings"

	"github.com/aetrion/dnsimple-go/dnsimple/webhook"
)

// Message formats the event into a text message suitable for being sent to a messaging service.
func Message(s MessagingService, e webhook.Event) (text string) {
	header := e.EventHeader()
	account := header.Account
	prefix := fmt.Sprintf("[%v] %v", s.FormatLink(account.Display, fmtURL("/a/%d/account", account.ID)), header.Actor.Pretty)

	switch event := e.(type) {
	case *webhook.ContactEvent:
		contactLink := s.FormatLink(fmt.Sprintf("%s %s", event.Contact.FirstName, event.Contact.LastName), fmtURL("/a/%d/contacts/%d", account.ID, event.Contact.ID))
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
		domainLink := s.FormatLink(event.Domain.Name, fmtURL("/a/%d/domains/%s", account.ID, event.Domain.Name))
		switch event.Name {
		case "domain.auto_renewal_enable":
			text = fmt.Sprintf("%s enabled auto-renewal for the domain %s", prefix, domainLink)
		case "domain.auto_renewal_disable":
			text = fmt.Sprintf("%s disabled auto-renewal for the domain %s", prefix, domainLink)
		case "domain.create":
			text = fmt.Sprintf("%s created the domain %s", prefix, domainLink)
		case "domain.delete":
			text = fmt.Sprintf("%s deleted the domain %s", prefix, domainLink)
		case "domain.register":
			text = fmt.Sprintf("%s registered the domain %s", prefix, domainLink)
		case "domain.renew":
			text = fmt.Sprintf("%s renewed the domain %s", prefix, domainLink)
		case "domain.delegation_change":
			servers := strings.Join(*event.Delegation, ", ")
			text = fmt.Sprintf("%s changed the delegation for the domain %s to %s", prefix, domainLink, servers)
		case "domain.registrant_change":
			registrant := event.Registrant.Label
			text = fmt.Sprintf("%s changed the registrant for the domain %s to %s", prefix, domainLink, registrant)
		case "domain.resolution_enable":
			text = fmt.Sprintf("%s enabled resolution for the domain %s", prefix, domainLink)
		case "domain.resolution_disable":
			text = fmt.Sprintf("%s disabled resolution for the domain %s", prefix, domainLink)
		case "domain.token_reset":
			text = fmt.Sprintf("%s reset the token for the domain %s", prefix, domainLink)
		case "domain.transfer":
			text = fmt.Sprintf("%s transferred the domain %s", prefix, domainLink)
		default:
			text = fmt.Sprintf("%s performed %s on domain %s", prefix, event.Name, domainLink)
		}
	case *webhook.ZoneRecordEvent:
		zoneRecordDisplay := fmt.Sprintf("%s %s.%s %s", event.ZoneRecord.Type, event.ZoneRecord.Name, event.ZoneRecord.ZoneID, event.ZoneRecord.Content)
		zoneRecordLink := s.FormatLink(zoneRecordDisplay, fmtURL("/a/%d/domains/%s/records/%d", account.ID, event.ZoneRecord.ZoneID, event.ZoneRecord.ID))
		switch event.Name {
		case "record.create":
			text = fmt.Sprintf("%s created the record %s", prefix, zoneRecordLink)
		case "record.update":
			text = fmt.Sprintf("%s updated the record %s", prefix, zoneRecordLink)
		case "record.delete":
			text = fmt.Sprintf("%s deleted the record %s", prefix, zoneRecordLink)
		}
	case *webhook.WebhookEvent:
		webhookLink := s.FormatLink(event.Webhook.URL, fmtURL("/a/%d/webhooks/%d", account.ID, event.Webhook.ID))
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

func eventRequestID(event webhook.Event) string {
	return event.EventHeader().RequestID
}

func fmtURL(path string, a ...interface{}) string {
	return fmt.Sprintf(dnsimpleURL+path, a...)
}
