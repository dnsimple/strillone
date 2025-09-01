package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnsimple/dnsimple-go/v6/dnsimple/webhook"
	xservice "github.com/dnsimple/strillone/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// cfg := config.LoadConfiguration()

	// Run the tests
	exitCode := m.Run()

	// Exit with the same code
	os.Exit(exitCode)
}

type TestMessagingService struct {
	Name string
}

func NewTestMessagingService(name string) *TestMessagingService {
	return &TestMessagingService{Name: name}
}

func (*TestMessagingService) FormatLink(url, name string) string {
	return fmt.Sprintf("<%s|%s>", url, name)
}

func (*TestMessagingService) PostEvent(_ *webhook.Event) (string, error) {
	return "ok", nil
}

func Test_Message(t *testing.T) {
	t.Run("unknown event", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		account := webhook.Account{Identifier: "ID", Display: "john.doe@gmail.com"}
		actor := webhook.Actor{Pretty: "john.doe@email.com"}
		event := webhook.Event{Actor: &actor, Account: &account, Name: "event.name"}

		result := xservice.Message(service, &event)
		assert.Equal(t, "[<john.doe@gmail.com|https://dnsimple.com/a/0/account>] john.doe@email.com performed event.name", result)
	})
}

func Test_Message_AccountEventData(t *testing.T) {
	t.Run("account.user_invite", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"name": "account.user_invite", "actor": {"pretty": "john.doe@email.com"}, "account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"}, "data": {"account": {"id": 12345, "email": "john.doe@email.com"}, "account_invitation": {"email": "jane.doe@email.com", "account_id": 12345}}}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "john.doe@email.com invited jane.doe@email.com to account <12345|https://dnsimple.com/a/12345/account/members>", result)
	})

	t.Run("account.user_invitation_accept", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"name":"account.user_invitation_accept", "actor": {"pretty": "jane.doe@email.com"}, "account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"}, "data":{"account":{"id": 12345, "email":"john.doe@email.com"}, "account_invitation":{"email":"jane.doe@email.com", "account_id":12345, "invitation_sent_at":"2020-05-12T18:42:44Z", "invitation_accepted_at":"2020-05-12T18:43:44Z"}}}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "jane.doe@email.com accepted invitation to account <12345|https://dnsimple.com/a/12345/account/members>", result)
	})

	t.Run("account.user_invitation_revoke", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"name":"account.user_invitation_revoke", "actor": {"pretty": "jane.doe@email.com"}, "account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"}, "data":{"account":{"id": 12345, "email":"john.doe@email.com"}, "account_invitation":{"email":"jane.doe@email.com", "account_id":12345, "invitation_sent_at":"2020-05-12T18:42:44Z", "invitation_accepted_at":null}}}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "jane.doe@email.com rejected invitation to account <12345|https://dnsimple.com/a/12345/account/members>", result)
	})

	t.Run("account.user_remove", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"name":"account.user_remove", "actor": {"pretty": "john.doe@email.com"}, "account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"}, "data":{"user":{"id":1120, "email":"jane.doe@email.com"}, "account":{"id":12345, "email":"john.doe@email.com"}}}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "john.doe@email.com removed jane.doe@email.com from account <12345|https://dnsimple.com/a/12345/account/members>", result)
	})
}

func Test_Message_DomainEventData(t *testing.T) {
	t.Run("domain.transfer_lock_disable", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "account_id": 1010, "auto_renew": false, "created_at": "2023-03-02T02:39:18Z", "expires_at": "2024-03-02T02:39:22Z", "expires_on": "2024-03-02", "updated_at": "2023-08-31T06:46:48Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 101}}, "name": "domain.transfer_lock_disable", "actor": {"id": "1010", "entity": "account", "pretty": "xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com"}, "account": {"id": 1010, "display": "xxxxxxx-xxxxxxx-xxxxxxx", "identifier": "xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com"}, "api_version": "v2", "request_identifier": "0f31483c-c303-497b-8a88-2edb48aa111e"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<xxxxxxx-xxxxxxx-xxxxxxx|https://dnsimple.com/a/1010/account>] xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com disabled transfer lock for the domain <example.com|https://dnsimple.com/a/1010/domains/example.com>", result)
	})

	t.Run("domain.transfer_lock_enable", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "account_id": 1010, "auto_renew": false, "created_at": "2023-03-02T02:39:18Z", "expires_at": "2024-03-02T02:39:22Z", "expires_on": "2024-03-02", "updated_at": "2023-08-31T06:46:48Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 101}}, "name": "domain.transfer_lock_enable", "actor": {"id": "1010", "entity": "account", "pretty": "xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com"}, "account": {"id": 1010, "display": "xxxxxxx-xxxxxxx-xxxxxxx", "identifier": "xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com"}, "api_version": "v2", "request_identifier": "0f31483c-c303-497b-8a88-2edb48aa111e"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<xxxxxxx-xxxxxxx-xxxxxxx|https://dnsimple.com/a/1010/account>] xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com enabled transfer lock for the domain <example.com|https://dnsimple.com/a/1010/domains/example.com>", result)
	})
}

func Test_Message_DNSSECEventData(t *testing.T) {
	t.Run("dnssec.create", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"zone": {"id": 315333, "name": "example-20230920163010.com", "active": true, "reverse": false, "secondary": false, "account_id": 625, "created_at": "2023-09-20T14:30:19Z", "updated_at": "2025-06-13T13:11:52Z", "last_transferred_at": null}, "dnssec": {"enabled": true, "created_at": "2025-06-13T13:11:52Z", "updated_at": "2025-06-13T13:11:52Z"}}, "name": "dnssec.create", "actor": {"id": "2", "entity": "user", "pretty": "simone.carletti@dnsimple.com"}, "account": {"id": 625, "display": "Webhook Tests", "identifier": "webhooks@example.com"}, "api_version": "v2", "request_identifier": "e3e3d10b-b0cd-498c-9fb6-9afdb44fd19a"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<Webhook Tests|https://dnsimple.com/a/625/account>] simone.carletti@dnsimple.com enabled DNSSEC for the zone <example-20230920163010.com|https://dnsimple.com/a/625/domains/example-20230920163010.com>", result)
	})

	t.Run("dnssec.delete", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"zone": {"id": 315333, "name": "example-20230920163010.com", "active": true, "reverse": false, "secondary": false, "account_id": 625, "created_at": "2023-09-20T14:30:19Z", "updated_at": "2025-06-13T13:11:58Z", "last_transferred_at": null}, "dnssec": {"enabled": true, "created_at": "2025-06-13T13:11:52Z", "updated_at": "2025-06-13T13:11:52Z"}}, "name": "dnssec.delete", "actor": {"id": "2", "entity": "user", "pretty": "simone.carletti@dnsimple.com"}, "account": {"id": 625, "display": "Webhook Tests", "identifier": "webhooks@example.com"}, "api_version": "v2", "request_identifier": "1096e2f6-71d2-4d8f-a7c0-05858eb68454"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<Webhook Tests|https://dnsimple.com/a/625/account>] simone.carletti@dnsimple.com disabled DNSSEC for the zone <example-20230920163010.com|https://dnsimple.com/a/625/domains/example-20230920163010.com>", result)
	})

	t.Run("dnssec.rotation_start", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"zone": {"id": 315333, "name": "example-20230920163010.com", "active": true, "reverse": false, "secondary": false, "account_id": 625, "created_at": "2023-09-20T14:30:19Z", "updated_at": "2025-06-13T13:11:52Z", "last_transferred_at": null}, "dnssec": {"enabled": true, "created_at": "2025-06-13T13:11:52Z", "updated_at": "2025-06-13T13:11:52Z"}, "delegation_signer_record": {"id": null, "digest": "42AEE231E98FECE484E9FA983CEF28AFFA56E99AD26347806BC6AF291F67DE83", "keytag": "60812", "algorithm": "8", "domain_id": 30943, "created_at": null, "public_key": "AwEAAefkQW+2ZO79nSaQ2eUVGhdmcapkGvZcmc5Xd9ig50k76eldueP198qtMsCV+27KZLqphTbYb4zOh1cF432TyKluZu89VzeVNC7Lq4kxDN1ahJfOCmBXg+/JAbb+NtKzH751CP/cWbwAShCwRb10TipwmTdZRYdOs3y9tKQq7BIE4YnEGrGb4lfCXrKK15Jn2im2f/MtVEuSF+eDB3X/XPU=", "updated_at": null, "digest_type": "2"}}, "name": "dnssec.rotation_start", "actor": {"id": "", "entity": "dnsimple", "pretty": "support@dnsimple.com"}, "account": {"id": 625, "display": "Webhook Tests", "identifier": "webhooks@example.com"}, "api_version": "v2", "request_identifier": "cf0252ab-80ea-40ab-8590-7b145d28dd61"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<Webhook Tests|https://dnsimple.com/a/625/account>] support@dnsimple.com started DNSSEC key rotation for the zone <example-20230920163010.com|https://dnsimple.com/a/625/domains/example-20230920163010.com>", result)
	})

	t.Run("dnssec.rotation_complete", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"zone": {"id": 315333, "name": "example-20230920163010.com", "active": true, "reverse": false, "secondary": false, "account_id": 625, "created_at": "2023-09-20T14:30:19Z", "updated_at": "2025-06-13T13:11:52Z", "last_transferred_at": null}, "dnssec": {"enabled": true, "created_at": "2025-06-13T13:11:52Z", "updated_at": "2025-06-13T13:11:52Z"}, "delegation_signer_record": {"id": null, "digest": "992059C73169F2D049377884F210F893CF19CB56A4F8198B6424FF3D9BA1B4AA", "keytag": "25337", "algorithm": "8", "domain_id": 41557, "created_at": null, "public_key": "AwEAAb/7eqMeecFAp+KygQzEBGR204F35ATbY000IpetSFWJMlBHPnh6tKP9gJFQyI3hfPTO3fAcr5G7dRrWE6zECHicMm062h8IcHukIpTTOF1PnHhVI47Mk1ZHwAgmx12e8MhoTUPNqutpZG/AsLfN9T3vMQWKDDGwk5sM2JgWXcBa4ys2f0AnYOQain0LGdTUxl76B16MK1BqaNwZrTNo71U=", "updated_at": null, "digest_type": "2"}}, "name": "dnssec.rotation_complete", "actor": {"id": "", "entity": "dnsimple", "pretty": "support@dnsimple.com"}, "account": {"id": 625, "display": "Webhook Tests", "identifier": "webhooks@example.com"}, "api_version": "v2", "request_identifier": "f26050c7-cb4e-49be-bab1-e42f27f5b60a"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<Webhook Tests|https://dnsimple.com/a/625/account>] support@dnsimple.com completed DNSSEC key rotation for the zone <example-20230920163010.com|https://dnsimple.com/a/625/domains/example-20230920163010.com>", result)
	})
}

func Test_Message_ZoneEventData(t *testing.T) {
	t.Run("zone.delete", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"name": "zone.delete", "actor": {"pretty": "john.doe@email.com"}, "account": {"id": 1010, "display": "example-account", "identifier": "example-account@email.com"}, "data": {"zone": {"id": 12345, "name": "example.com", "account_id": 1010, "created_at": "2023-03-02T02:39:18Z", "updated_at": "2023-08-31T06:46:48Z"}}}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<example-account|https://dnsimple.com/a/1010/account>] john.doe@email.com deleted the zone <example.com|https://dnsimple.com/a/1010/domains/example.com>", result)
	})
}

func Test_fmtURL(t *testing.T) {
	assert.Equal(t, "https://dnsimple.com/a/1010/domains/1", xservice.FmtURL("/a/%v/domains/%v", "1010", 1))
}
