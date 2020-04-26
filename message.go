package strillone

import (
	"fmt"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"
)

// Message formats the event into a text message suitable for being sent to a messaging service.
func Message(s MessagingService, e *webhook.Event) (text string) {
	account := e.Account
	prefix := fmt.Sprintf("[%v] %v", s.FormatLink(account.Display, fmtURL("/a/%d/account", account.ID)), e.Actor.Pretty)

	switch data := e.GetData().(type) {
	case *webhook.CertificateEventData:
		certificate := data.Certificate
		certificateDisplay := certificate.CommonName
		certificateLink := s.FormatLink(certificateDisplay, fmtURL("/a/%d/domains/%d/certificates/%d", account.ID, certificate.DomainID, certificate.ID))
		switch e.Name {
		case "certificate.remove_private_key":
			text = fmt.Sprintf("%s deleted the private key for the certificate %s", prefix, certificateLink)
		default:
			text = fmt.Sprintf("%s performed %s", prefix, e.Name)
		}
	case *webhook.ContactEventData:
		contactDisplay := fmt.Sprintf("%s %s", data.Contact.FirstName, data.Contact.LastName)
		contactLink := s.FormatLink(contactDisplay, fmtURL("/a/%d/contacts/%d", account.ID, data.Contact.ID))
		switch e.Name {
		case "contact.create":
			text = fmt.Sprintf("%s created the contact %s", prefix, contactLink)
		case "contact.update":
			text = fmt.Sprintf("%s update the contact %s", prefix, contactLink)
		case "contact.delete":
			text = fmt.Sprintf("%s deleted the contact %s", prefix, contactLink)
		default:
			text = fmt.Sprintf("%s performed %s", prefix, e.Name)
		}
	case *webhook.DomainEventData:
		domainDisplay := data.Domain.Name
		domainLink := s.FormatLink(domainDisplay, fmtURL("/a/%d/domains/%s", account.ID, data.Domain.Name))
		switch e.Name {
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
			servers := strings.Join(*data.Delegation, ", ")
			text = fmt.Sprintf("%s changed the delegation for the domain %s to %s", prefix, domainLink, servers)
		case "domain.registrant_change":
			registrant := data.Registrant.Label
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
			text = fmt.Sprintf("%s performed %s on domain %s", prefix, e.Name, domainLink)
		}
	case *webhook.WhoisPrivacyEventData:
		domainDisplay := data.Domain.Name
		domainLink := s.FormatLink(domainDisplay, fmtURL("/a/%d/domains/%s", account.ID, data.Domain.Name))
		switch e.Name {
		case "whois_privacy.disable":
			text = fmt.Sprintf("%s disabled whois privacy for the domain %s", prefix, domainLink)
		case "whois_privacy.enable":
			text = fmt.Sprintf("%s enabled whois privacy for the domain %s", prefix, domainLink)
		case "whois_privacy.purchase":
			text = fmt.Sprintf("%s purchased whois privacy for the domain %s", prefix, domainLink)
		case "whois_privacy.renew":
			text = fmt.Sprintf("%s renewed whois privacy for the domain %s", prefix, domainLink)
		}

	case *webhook.ZoneRecordEventData:
		zoneRecordDisplay := fmt.Sprintf("%s %s.%s %s", data.ZoneRecord.Type, data.ZoneRecord.Name, data.ZoneRecord.ZoneID, data.ZoneRecord.Content)
		zoneRecordLink := s.FormatLink(zoneRecordDisplay, fmtURL("/a/%d/domains/%s/records/%d", account.ID, data.ZoneRecord.ZoneID, data.ZoneRecord.ID))
		switch e.Name {
		case "zone_record.create":
			text = fmt.Sprintf("%s created the record %s", prefix, zoneRecordLink)
		case "zone_record.update":
			text = fmt.Sprintf("%s updated the record %s", prefix, zoneRecordLink)
		case "zone_record.delete":
			text = fmt.Sprintf("%s deleted the record %s", prefix, zoneRecordLink)
		}
	case *webhook.WebhookEventData:
		webhookDisplay := data.Webhook.URL
		webhookLink := s.FormatLink(webhookDisplay, fmtURL("/a/%d/webhooks/%d", account.ID, data.Webhook.ID))
		switch e.Name {
		case "webhook.create":
			text = fmt.Sprintf("%s created the webhook %s", prefix, webhookLink)
		case "webhook.delete":
			text = fmt.Sprintf("%s deleted the webhook %s", prefix, webhookLink)
		}
	default:
		text = fmt.Sprintf("%s performed %s", prefix, e.Name)
	}

	return
}

func eventRequestID(e *webhook.Event) string {
	return e.RequestID
}

func fmtURL(path string, a ...interface{}) string {
	return fmt.Sprintf(dnsimpleURL+path, a...)
}
