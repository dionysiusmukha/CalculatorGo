package net

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurl_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok-test"))
	}))
	defer srv.Close()

	body, err := Curl(srv.URL)
	if err != nil {
		t.Fatalf("Curl error: %v", err)
	}
	if body != "ok-test" {
		t.Fatalf("expected ok-test, got %q", body)
	}
}

func TestCurl_BadURL(t *testing.T) {
	_, err := Curl("http://127.0.0.1:0") // заведомо плохой порт
	if err == nil {
		t.Fatalf("expected error for bad url, got nil")
	}
}
