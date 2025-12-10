package http_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	appServer "github.com/dnsimple/strillone/internal/http"
	"github.com/stretchr/testify/assert"
)

var server *appServer.Server

func init() {
	server = appServer.NewServer()
}

func TestMain(m *testing.M) {
	// cfg := config.LoadConfiguration()

	// Run the tests
	exitCode := m.Run()

	// Exit with the same code
	os.Exit(exitCode)
}

func TestRoot(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	server.Root(response, request)

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
	// Set path values for Go 1.22+ routing
	request = request.WithContext(http.ContextWithPathValue(request.Context(), "slackAlpha", "-"))
	request = request.WithContext(http.ContextWithPathValue(request.Context(), "slackBeta", "-"))
	request = request.WithContext(http.ContextWithPathValue(request.Context(), "slackGamma", "-"))
	response := httptest.NewRecorder()

	server.Slack(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "[<https://dnsimple.com/a/1010/account|User>] example@example.com created the domain <https://dnsimple.com/a/1010/domains/example.com|example.com>\n", response.Body.String())
}

func TestSlackTwice(t *testing.T) {
	payload := `{"data": {"domain": {"id": 1, "name": "example.com", "state": "hosted", "token": "domain-token", "account_id": 1010, "auto_renew": false, "created_at": "2016-02-07T14:46:29.142Z", "expires_on": null, "updated_at": "2016-02-07T14:46:29.142Z", "unicode_name": "example.com", "private_whois": false, "registrant_id": null}}, "actor": {"id": "1", "entity": "user", "pretty": "example@example.com"}, "account": {"id": 1010, "display": "User", "identifier": "user"}, "name": "domain.create", "api_version": "v2", "request_identifier": "096bfc29-2bf0-40c6-0000-f03b1f8521f1"}`
	request, _ := http.NewRequest("POST", "/slack/-/-/-", strings.NewReader(payload))
	// Set path values for Go 1.22+ routing
	request = request.WithContext(http.ContextWithPathValue(request.Context(), "slackAlpha", "-"))
	request = request.WithContext(http.ContextWithPathValue(request.Context(), "slackBeta", "-"))
	request = request.WithContext(http.ContextWithPathValue(request.Context(), "slackGamma", "-"))
	response := httptest.NewRecorder()

	server.Slack(response, request)

	if want := http.StatusOK; want != response.Code {
		t.Errorf("POST /slack expected HTTP %v, got %v", want, response.Code)
	}
	assert.Empty(t, response.Header().Get(appServer.HeaderProcessingStatus))

	requestDuplicate, _ := http.NewRequest("POST", "/slack/-/-/-", strings.NewReader(payload))
	// Set path values for Go 1.22+ routing
	requestDuplicate = requestDuplicate.WithContext(http.ContextWithPathValue(requestDuplicate.Context(), "slackAlpha", "-"))
	requestDuplicate = requestDuplicate.WithContext(http.ContextWithPathValue(requestDuplicate.Context(), "slackBeta", "-"))
	requestDuplicate = requestDuplicate.WithContext(http.ContextWithPathValue(requestDuplicate.Context(), "slackGamma", "-"))
	responseDuplicate := httptest.NewRecorder()

	server.Slack(responseDuplicate, requestDuplicate)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "skipped;already-processed", responseDuplicate.Header().Get(appServer.HeaderProcessingStatus))
}
