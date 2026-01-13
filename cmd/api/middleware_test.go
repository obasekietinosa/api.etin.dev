package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAllowedOrigin_NormalizedMatch(t *testing.T) {
	app := &application{}
	app.config.cors.trustedOrigins = []string{
		" https://ADMIN.ETIN.dev/ ",
	}

	for i, origin := range app.config.cors.trustedOrigins {
		app.config.cors.trustedOrigins[i] = normalizeOrigin(origin)
	}

	allowedOrigin, ok := app.getAllowedOrigin("https://admin.etin.dev")
	if !ok {
		t.Fatalf("expected origin to be allowed")
	}

	if allowedOrigin != "https://admin.etin.dev" {
		t.Fatalf("expected allowed origin to equal request origin; got %q", allowedOrigin)
	}
}

func TestLogRequest(t *testing.T) {
	var buf bytes.Buffer
	app := &application{
		logger: log.New(&buf, "", 0),
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := app.logRequest(next)

	req := httptest.NewRequest(http.MethodGet, "/test/url", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	logOutput := buf.String()
	expectedLogPart := "1.2.3.4:1234 - HTTP/1.1 GET /test/url"

	if !strings.Contains(logOutput, expectedLogPart) {
		t.Errorf("expected log output to contain %q, got %q", expectedLogPart, logOutput)
	}
}
