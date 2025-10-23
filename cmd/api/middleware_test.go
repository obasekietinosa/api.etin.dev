package main

import "testing"

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

func TestGetAllowedOrigin_HostOnlyConfig(t *testing.T) {
	app := &application{}
	app.config.cors.trustedOrigins = []string{
		"admin.etin.dev",
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

func TestNormalizeOriginHostDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "https host", input: "admin.etin.dev", expected: "https://admin.etin.dev"},
		{name: "localhost", input: "localhost:3000", expected: "http://localhost:3000"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeOrigin(tc.input); got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}
