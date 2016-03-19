package main

import (
	"testing"
)

func Test_fmtURL(t *testing.T) {
	if want, got := "https://dnsimple.com/a/1010/domains/1", fmtURL("/a/%v/domains/%v", "1010", 1); want != got {
		t.Fatalf("Expected %v, got %v", want, got)
	}
}
