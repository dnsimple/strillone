package strillone_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dnsimple/strillone"
	"github.com/julienschmidt/httprouter"
)

var server *strillone.Server

func init() {
	server = strillone.NewServer()
}

func TestRoot(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	server.Root(response, request, httprouter.Params{})

	if want := http.StatusOK; want != response.Code {
		t.Errorf("GET / expected HTTP %v, got %v", want, response.Code)
	}
	if want, got := "application/json", response.Header().Get("Content-Type"); want != got {
		t.Errorf("GET / expected Content-Type %v, got %v", want, got)
	}
}

func TestSlack(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-991b-f03b1f8521f1"}`

	request, _ := http.NewRequest("POST", "/slack/-/-/-", strings.NewReader(payload))
	response := httptest.NewRecorder()

	server.Slack(response, request, httprouter.Params{})

	if want := http.StatusOK; want != response.Code {
		t.Errorf("POST /slack expected HTTP %v, got %v", want, response.Code)
	}

	want := "[<https://dnsimple.com/a/1010/account|User>] example@example.com created the domain <https://dnsimple.com/a/1010/domains/example.com|example.com>\n"
	if got := response.Body.String(); want != got {
		t.Errorf("POST /slack expected response\n\t%v\ngot\n\t%v", want, got)
	}
}

func TestSlackTwice(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-0000-f03b1f8521f1"}`

	request, _ := http.NewRequest("POST", "/slack/-/-/-", strings.NewReader(payload))
	response := httptest.NewRecorder()

	server.Slack(response, request, httprouter.Params{})

	if want := http.StatusOK; want != response.Code {
		t.Errorf("POST /slack expected HTTP %v, got %v", want, response.Code)
	}
	if want, got := "", response.Header().Get(strillone.HeaderProcessingStatus); want != got {
		t.Errorf("POST /slack X-Processing-Status expected empty, got %v", got)
	}

	requestDuplicate, _ := http.NewRequest("POST", "/slack/-/-/-", strings.NewReader(payload))
	responseDuplicate := httptest.NewRecorder()

	server.Slack(responseDuplicate, requestDuplicate, httprouter.Params{})

	if want := http.StatusOK; want != responseDuplicate.Code {
		t.Errorf("POST /slack (duplicate) expected HTTP %v, got %v", want, responseDuplicate.Code)
	}
	if want, got := "skipped;already-processed", responseDuplicate.Header().Get(strillone.HeaderProcessingStatus); want != got {
		t.Errorf("POST /slack (duplicate) X-Processing-Status expected %v, got %v", want, got)
	}
}
