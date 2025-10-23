package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple/webhook"
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

	t.Run("account.sso_user_add", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"user": {"id": 1111, "email": "xxxxx@xxxxxx.xxx", "created_at": "2025-09-16T22:12:34Z", "updated_at": "2025-09-18T10:46:19Z"}, "account": {"id": 4, "email": "yyyyy@yyyyyy.yyy", "created_at": "2025-08-13T23:09:47Z", "updated_at": "2025-08-13T23:10:05Z", "plan_identifier": "teams-v2-monthly"}, "account_identity_provider": {"organization_identifier": "51fae1e9-ce56-4df2-8364-cdab573027aa"}}, "name": "account.sso_user_add", "actor": {"id": "", "entity": "dnsimple", "pretty": "support@dnsimple.com"}, "account": {"id": 4, "display": "Personal", "identifier": "xxxxxx"}, "api_version": "v2", "request_identifier": "4aedf8d3-f93d-4a42-99d9-ec20c9349358"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "support@dnsimple.com added xxxxx@xxxxxx.xxx to account <4|https://dnsimple.com/a/4/account/members> via SSO", result)
	})
}

func Test_Message_CertificateEventData(t *testing.T) {
	t.Run("certificate.issue", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"certificate": {"id": 101967, "csr": "-----BEGIN CERTIFICATE REQUEST-----\nMIICmTCCAYECAQAwGjEYMBYGA1UEAwwPd3d3LmJpbmdvLnBpenphMIIBIjANBgkq\nhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw4+KoZ9IDCK2o5qAQpi+Icu5kksmjQzx\n5o5g4B6XhRxhsfHlK/i3iU5hc8CONjyVv8j82835RNsiKrflnxGa9SH68vbQfcn4\nIpbMz9c+Eqv5h0Euqlc3A4DBzp0unEu5QAUhR6Xu1TZIWDPjhrBOGiszRlLQcp4F\nzy6fD6j5/d/ylpzTp5v54j+Ey31Bz86IaBPtSpHI+Qk87Hs8DVoWxZk/6RlAkyur\nXDGWnPu9n3RMfs9ag5anFhggLIhCNtVN4+0vpgPQ59pqwYo8TfdYzK7WSKeL7geu\nCqVE3bHAqU6dLtgHOZfTkLwGycUh4p9aawuc6fsXHHYDpIL8s3vAvwIDAQABoDow\nOAYJKoZIhvcNAQkOMSswKTAnBgNVHREEIDAeggtiaW5nby5waXp6YYIPd3d3LmJp\nbmdvLnBpenphMA0GCSqGSIb3DQEBCwUAA4IBAQBwOLKv+PO5hSJkgqS6wL/wRqLh\nQ1zbcHRHAjRjnpRz06cDvN3X3aPI+lpKSNFCI0A1oKJG7JNtgxX3Est66cuO8ESQ\nPIb6WWN7/xlVlBCe7ZkjAFgN6JurFdclwCp/NI5wBCwj1yb3Ar5QQMFIZOezIgTI\nAWkQSfCmgkB96d6QlDWgidYDDjcsXugQveOQRPlHr0TsElu47GakxZdJCFZU+WPM\nodQQf5SaqiIK2YaH1dWO//4KpTS9QoTy1+mmAa27apHcmz6X6+G5dvpHZ1qH14V0\nJoMWIK+39HRPq6mDo1UMVet/xFUUrG/H7/tFlYIDVbSpVlpVAFITd/eQkaW/\n-----END CERTIFICATE REQUEST-----\n", "name": "www", "state": "issued", "years": 1, "domain_id": 289333, "auto_renew": false, "contact_id": 2511, "created_at": "2020-06-18T18:54:17Z", "expires_at": "2020-09-16T18:10:13Z", "expires_on": "2020-09-16", "updated_at": "2020-06-18T19:10:14Z", "common_name": "www.bingo.pizza", "alternate_names": [], "authority_identifier": "letsencrypt"}}, "name": "certificate.issue", "actor": {"id": "system", "entity": "dnsimple", "pretty": "support@dnsimple.com"}, "account": {"id": 5623, "display": "DNSimple", "identifier": "dnsimple"}, "api_version": "v2", "request_identifier": "2092d585-bd6c-40ba-b0ac-bd7ebab8b4b1"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<DNSimple|https://dnsimple.com/a/5623/account>] support@dnsimple.com issued the certificate <www.bingo.pizza|https://dnsimple.com/a/5623/domains/289333/certificates/101967>", result)
	})

	t.Run("certificate.remove_private_key", func(t *testing.T) {
		service := NewTestMessagingService("dummyMessagingService")
		payload := `{"data": {"certificate": {"id": 101972, "csr": "-----BEGIN CERTIFICATE REQUEST-----\nMIICmTCCAYECAQAwGjEYMBYGA1UEAwwPd3d3LmJpbmdvLnBpenphMIIBIjANBgkq\nhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuWImvp8rGULVobq+6dTpMPbxdePpIjki\nUJ5fTcpNfuRw/8EPU9UcQ5QaHToVyUhPtdv1QWtZbTVJpl0u9uvUviZEL0+NxBcR\nj4ymiPtyH6ZErYisUgIULaSEmkCz6YFf1X+fSDBcASvibHqkwhzN3ugVcKgqTHPm\nbbCpNv2EPOqkvwiqVPrjRmgqtjQmfO60K7+8aMZqyIjMslsDufGJ4sRaiiusRJUG\n1QzZqnlGp5Vrz5XdowHAQfLcUd+lPevk/lDkfmV6bxuoZyEKFAHVRFCM8Aw3nks4\nONrtWTdnOd5QoxwcOnbtl1S0bydQpulJjefy8sXQq/XUwsyAP6uLHwIDAQABoDow\nOAYJKoZIhvcNAQkOMSswKTAnBgNVHREEIDAeggtiaW5nby5waXp6YYIPd3d3LmJp\nbmdvLnBpenphMA0GCSqGSIb3DQEBCwUAA4IBAQBUPXgSgTO8hDjIvz8SKCE0EDbF\nEmgVqjynO2wvc8Dn9E3xp9GsYLNChUItUSzh0dnxn+XYBtia356bw5EaA3ZCbZIJ\nA/JyGavwNqIeBVSMsMCVXiM9NFunkWchid7bh1mS+W4/8gqEElIYRRIIOP7LEHq+\nxE7ZUH9qjKpiHKL/YTf2zVo4y6opjY4WnDxonQ2nMeJxfj8GdVskXYQoMxyVneRI\n0Of1gTZWvP1f1F9ddcjZDnb9VLdKcqrY395Zvy+FkNetd0xHRu2VBJDFbMnH8Gsr\nTd6BwijqU3kNM1j2zWvOhfO9tPbcl4BbVRg+/V2bq3jLCld3Bj38d1CJLz21\n-----END CERTIFICATE REQUEST-----\n", "name": "www", "state": "issued", "years": 1, "domain_id": 289333, "auto_renew": false, "contact_id": 2511, "created_at": "2020-06-18T19:56:20Z", "expires_at": "2020-09-16T19:10:05Z", "expires_on": "2020-09-16", "updated_at": "2020-06-18T20:41:05Z", "common_name": "www.bingo.pizza", "alternate_names": [], "authority_identifier": "letsencrypt"}}, "name": "certificate.remove_private_key", "actor": {"id": "109237", "entity": "user", "pretty": "xxxxxx.xxxxxxx@dnsimple.com"}, "account": {"id": 5623, "display": "DNSimple", "identifier": "dnsimple"}, "api_version": "v2", "request_identifier": "bdee91db-7339-49d1-a93c-104b46d72235"}`
		event, err := webhook.ParseEvent([]byte(payload))
		assert.NoError(t, err)

		result := xservice.Message(service, event)
		assert.Equal(t, "[<DNSimple|https://dnsimple.com/a/5623/account>] xxxxxx.xxxxxxx@dnsimple.com deleted the private key for the certificate <www.bingo.pizza|https://dnsimple.com/a/5623/domains/289333/certificates/101972>", result)
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
