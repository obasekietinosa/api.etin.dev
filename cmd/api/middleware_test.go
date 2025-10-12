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
