package service_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	xservice "github.com/dnsimple/strillone/internal/service"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestIsClientError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"not found", slack.StatusCodeError{Code: http.StatusNotFound}, true},
		{"forbidden wrapped", fmt.Errorf("post: %w", slack.StatusCodeError{Code: http.StatusForbidden}), true},
		{"rate limited", &slack.RateLimitedError{RetryAfter: time.Second}, true},
		{"server error", slack.StatusCodeError{Code: http.StatusInternalServerError}, false},
		{"network error", errors.New("failed to post webhook"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, xservice.IsClientError(tt.err))
		})
	}
}
