package service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"
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

func Test_Message_AccountUserInviteEvent(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	payload := `{
		"name": "account.user_invite",
		"actor": {"pretty": "john.doe@email.com"},
		"account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"},
		"data": {
			"account": {
				"id": 12345,
				"email": "john.doe@email.com"
			},
			"account_invitation": {
				"email": "jane.doe@email.com",
				"account_id": 12345
			}
		}
	}`
	event, err := webhook.ParseEvent([]byte(payload))
	assert.NoError(t, err)

	result := xservice.Message(service, event)
	assert.Equal(t, "john.doe@email.com invited jane.doe@email.com to account <12345|https://dnsimple.com/a/12345/account/members>", result)
}

func Test_Message_AccountUserInvitationAcceptEvent(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	payload := `{
		"name":"account.user_invitation_accept",
		"actor": {"pretty": "jane.doe@email.com"},
		"account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"},
		"data":{
			"account":{
				"id": 12345,
				"email":"john.doe@email.com"
			},
			"account_invitation":{
				"email":"jane.doe@email.com",
				"account_id":12345,
				"invitation_sent_at":"2020-05-12T18:42:44Z",
				"invitation_accepted_at":"2020-05-12T18:43:44Z"
			}
		}
	}`
	event, err := webhook.ParseEvent([]byte(payload))
	assert.NoError(t, err)

	result := xservice.Message(service, event)
	assert.Equal(t, "jane.doe@email.com accepted invitation to account <12345|https://dnsimple.com/a/12345/account/members>", result)
}

func Test_Message_AccountUserInvitationRevokeEvent(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	payload := `{
		"name":"account.user_invitation_revoke",
		"actor": {"pretty": "jane.doe@email.com"},
		"account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"},
		"data":{
			"account":{
				"id": 12345,
				"email":"john.doe@email.com"
			},
			"account_invitation":{
				"email":"jane.doe@email.com",
				"account_id":12345,
				"invitation_sent_at":"2020-05-12T18:42:44Z",
				"invitation_accepted_at":null
			}
		}
	}`
	event, err := webhook.ParseEvent([]byte(payload))
	assert.NoError(t, err)

	result := xservice.Message(service, event)
	assert.Equal(t, "jane.doe@email.com rejected invitation to account <12345|https://dnsimple.com/a/12345/account/members>", result)
}

func Test_Message_AccountUserRemoveEvent(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	payload := `{
		"name":"account.user_remove",
		"actor": {"pretty": "john.doe@email.com"},
		"account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"},
		"data":{
			"user":{
				"id":1120,
				"email":"jane.doe@email.com"
			},
			"account":{
				"id":12345,
				"email":"john.doe@email.com"
			}
		}
	}`
	event, err := webhook.ParseEvent([]byte(payload))
	assert.NoError(t, err)

	result := xservice.Message(service, event)
	assert.Equal(t, "john.doe@email.com removed jane.doe@email.com from account <12345|https://dnsimple.com/a/12345/account/members>", result)
}

func Test_Message_DomainTransferLockDisableEvent(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "account_id": 1010, "auto_renew": false, "created_at": "2023-03-02T02:39:18Z", "expires_at": "2024-03-02T02:39:22Z", "expires_on": "2024-03-02", "updated_at": "2023-08-31T06:46:48Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 101}}, "name": "domain.transfer_lock_disable", "actor": {"id": "1010", "entity": "account", "pretty": "xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com"}, "account": {"id": 1010, "display": "xxxxxxx-xxxxxxx-xxxxxxx", "identifier": "xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com"}, "api_version": "v2", "request_identifier": "0f31483c-c303-497b-8a88-2edb48aa111e"}`
	event, err := webhook.ParseEvent([]byte(payload))
	assert.NoError(t, err)

	result := xservice.Message(service, event)
	assert.Equal(t, "[<xxxxxxx-xxxxxxx-xxxxxxx|https://dnsimple.com/a/1010/account>] xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com disabled transfer lock for the domain <example.com|https://dnsimple.com/a/1010/domains/example.com>", result)
}

func Test_Message_DomainTransferLockEnableEvent(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "registered", "account_id": 1010, "auto_renew": false, "created_at": "2023-03-02T02:39:18Z", "expires_at": "2024-03-02T02:39:22Z", "expires_on": "2024-03-02", "updated_at": "2023-08-31T06:46:48Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": 101}}, "name": "domain.transfer_lock_enable", "actor": {"id": "1010", "entity": "account", "pretty": "xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com"}, "account": {"id": 1010, "display": "xxxxxxx-xxxxxxx-xxxxxxx", "identifier": "xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com"}, "api_version": "v2", "request_identifier": "0f31483c-c303-497b-8a88-2edb48aa111e"}`
	event, err := webhook.ParseEvent([]byte(payload))
	assert.NoError(t, err)

	result := xservice.Message(service, event)
	assert.Equal(t, "[<xxxxxxx-xxxxxxx-xxxxxxx|https://dnsimple.com/a/1010/account>] xxxxxxx-xxxxxxx-xxxxxxx@xxxxx.com enabled transfer lock for the domain <example.com|https://dnsimple.com/a/1010/domains/example.com>", result)
}

func Test_Message_DefaultMessage(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	account := webhook.Account{Identifier: "ID", Display: "john.doe@gmail.com"}
	actor := webhook.Actor{Pretty: "john.doe@email.com"}
	event := webhook.Event{Actor: &actor, Account: &account, Name: "event.name"}

	result := xservice.Message(service, &event)
	assert.Equal(t, "[<john.doe@gmail.com|https://dnsimple.com/a/0/account>] john.doe@email.com performed event.name", result)
}

func Test_fmtURL(t *testing.T) {
	assert.Equal(t, "https://dnsimple.com/a/1010/domains/1", xservice.FmtURL("/a/%v/domains/%v", "1010", 1))
}
