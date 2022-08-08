package strillone_test

import (
	"fmt"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"

	"github.com/dnsimple/strillone"
)

type TestMessagingService struct {
	Name string
}

func NewTestMessagingService(name string) *TestMessagingService {
	return &TestMessagingService{Name: name}
}

func (*TestMessagingService) FormatLink(url, name string) string {
	return fmt.Sprintf("<%s|%s>", url, name)
}

func (*TestMessagingService) PostEvent(event *webhook.Event) (string, error) {
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
	if err != nil {
		t.Fatalf("Error parsing: %v.\n%v", err, payload)
	}

	result := strillone.Message(service, event)

	if want, got := "john.doe@email.com invited jane.doe@email.com to account <12345|https://dnsimple.com/a/12345/account/members>", result; want != got {
		t.Fatalf("Expected '%v', got '%v'", want, got)
	}
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
	if err != nil {
		t.Fatalf("Error parsing: %v.\n%v", err, payload)
	}

	result := strillone.Message(service, event)

	if want, got := "jane.doe@email.com accepted invitation to account <12345|https://dnsimple.com/a/12345/account/members>", result; want != got {
		t.Fatalf("Expected '%v', got '%v'", want, got)
	}
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
	if err != nil {
		t.Fatalf("Error parsing: %v.\n%v", err, payload)
	}

	result := strillone.Message(service, event)

	if want, got := "jane.doe@email.com rejected invitation to account <12345|https://dnsimple.com/a/12345/account/members>", result; want != got {
		t.Fatalf("Expected '%v', got '%v'", want, got)
	}
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
	if err != nil {
		t.Fatalf("Error parsing: %v.\n%v", err, payload)
	}

	result := strillone.Message(service, event)

	if want, got := "john.doe@email.com removed jane.doe@email.com from account <12345|https://dnsimple.com/a/12345/account/members>", result; want != got {
		t.Fatalf("Expected '%v', got '%v'", want, got)
	}
}

func Test_Message_DefaultMessage(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	account := webhook.Account{Identifier: "ID", Display: "john.doe@gmail.com"}
	actor := webhook.Actor{Pretty: "john.doe@email.com"}
	event := webhook.Event{Actor: &actor, Account: &account, Name: "event.name"}

	result := strillone.Message(service, &event)

	if want, got := "[<john.doe@gmail.com|https://dnsimple.com/a/0/account>] john.doe@email.com performed event.name", result; want != got {
		t.Fatalf("Expected %v, got %v", want, got)
	}
}

func Test_fmtURL(t *testing.T) {
	if want, got := "https://dnsimple.com/a/1010/domains/1", strillone.FmtURL("/a/%v/domains/%v", "1010", 1); want != got {
		t.Fatalf("Expected %v, got %v", want, got)
	}
}
