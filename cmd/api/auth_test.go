package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{name: "standard", header: "Bearer token", want: "token"},
		{name: "lowercase", header: "bearer token", want: "token"},
		{name: "extra spaces", header: "  BEARER   token   ", want: "token"},
		{name: "missing token", header: "Bearer", wantErr: true},
		{name: "wrong scheme", header: "Basic token", wantErr: true},
		{name: "empty", header: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBearerToken(tt.header)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %q got %q", tt.want, got)
			}
		})
	}
}

func TestAdminVerifyHandler(t *testing.T) {
	app := &application{
		sessions: newSessionManager(time.Hour),
	}

	token, _, err := app.sessions.create()
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	reqValid := httptest.NewRequest(http.MethodGet, "/v1/admin/verify", nil)
	reqValid.Header.Set("Authorization", "Bearer "+token)
	recValid := httptest.NewRecorder()
	app.adminVerifyHandler(recValid, reqValid)

	if recValid.Code != http.StatusOK {
		t.Fatalf("expected status 200 for valid token, got %d", recValid.Code)
	}

	reqInvalid := httptest.NewRequest(http.MethodGet, "/v1/admin/verify", nil)
	reqInvalid.Header.Set("Authorization", "Bearer invalid-token")
	recInvalid := httptest.NewRecorder()
	app.adminVerifyHandler(recInvalid, reqInvalid)

	if recInvalid.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for invalid token, got %d", recInvalid.Code)
	}
}
