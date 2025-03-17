package http_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dnsimple/strillone/internal/config"
	appServer "github.com/dnsimple/strillone/internal/http"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

var server *appServer.Server

func init() {
	server = appServer.NewServer()
}

func TestMain(m *testing.M) {
	// Load configuration here
	// This ensures configuration is available before any tests run
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	config.Config = cfg

	// Run the tests
	exitCode := m.Run()

	// Exit with the same code
	os.Exit(exitCode)
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

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "[<https://dnsimple.com/a/1010/account|User>] example@example.com created the domain <https://dnsimple.com/a/1010/domains/example.com|example.com>\n", response.Body.String())
}

func TestSlackTwice(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-0000-f03b1f8521f1"}`
	request, _ := http.NewRequest("POST", "/slack/-/-/-", strings.NewReader(payload))
	response := httptest.NewRecorder()

	server.Slack(response, request, httprouter.Params{})

	if want := http.StatusOK; want != response.Code {
		t.Errorf("POST /slack expected HTTP %v, got %v", want, response.Code)
	}
	assert.Empty(t, response.Header().Get(appServer.HeaderProcessingStatus))

	requestDuplicate, _ := http.NewRequest("POST", "/slack/-/-/-", strings.NewReader(payload))
	responseDuplicate := httptest.NewRecorder()

	server.Slack(responseDuplicate, requestDuplicate, httprouter.Params{})

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "skipped;already-processed", responseDuplicate.Header().Get(appServer.HeaderProcessingStatus))
}
