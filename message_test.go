package strillone

import (
	"fmt"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"
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
    "actor": {"pretty": "xxxxxxxxxxxxxxxxx@xxxxxx.xxx"},
    "account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"},
    "data": {
      "account": {
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

	result := Message(service, event)

	if want, got := "john.doe@email.com invited jane.doe@email.com to account <12345|https://dnsimple.com/a/12345/account/members>", result; want != got {
		t.Fatalf("Expected '%v', got '%v'", want, got)
	}
}

func Test_Message_AccountUserInvitationAcceptEvent(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	payload := `{
    "name":"account.user_invitation_accept",
    "actor": {"pretty": "xxxxxxxxxxxxxxxxx@xxxxxx.xxx"},
    "account": {"display": "xxxxxxxx", "identifier": "xxxxxxxx"},
    "data":{
      "account":{
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

	result := Message(service, event)

	if want, got := "jane.doe@email.com accepted invitation to account <12345|https://dnsimple.com/a/12345/account/members>", result; want != got {
		t.Fatalf("Expected '%v', got '%v'", want, got)
	}
}

func Test_Message_DefaultMessage(t *testing.T) {
	service := NewTestMessagingService("dummyMessagingService")
	account := webhook.Account{Identifier: "ID", Display: "john.doe@gmail.com"}
	actor := webhook.Actor{Pretty: "john.doe@email.com"}
	event := webhook.Event{Actor: &actor, Account: &account, Name: "event.name"}

	result := Message(service, &event)

	if want, got := "[<john.doe@gmail.com|https://dnsimple.com/a/0/account>] john.doe@email.com performed event.name", result; want != got {
		t.Fatalf("Expected %v, got %v", want, got)
	}
}

func Test_fmtURL(t *testing.T) {
	if want, got := "https://dnsimple.com/a/1010/domains/1", fmtURL("/a/%v/domains/%v", "1010", 1); want != got {
		t.Fatalf("Expected %v, got %v", want, got)
	}
}
