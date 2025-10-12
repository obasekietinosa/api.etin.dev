package main

import (
	"net/http"
	"net/url"
	"strings"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

		app.logger.Printf("CORS origin: %s", r.Header.Get("Origin"))

		if origin := r.Header.Get("Origin"); origin != "" {
			if allowedOrigin, ok := app.getAllowedOrigin(origin); ok {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				if allowedOrigin != "*" {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) getAllowedOrigin(origin string) (string, bool) {
	normalizedOrigin := normalizeOrigin(origin)

	for _, trustedOrigin := range app.config.cors.trustedOrigins {
		if trustedOrigin == "*" {
			return "*", true
		}
		if trustedOrigin == normalizedOrigin {
			return origin, true
		}
	}

	return "", false
}

func normalizeOrigin(origin string) string {
	trimmed := strings.TrimSpace(origin)
	trimmed = strings.TrimRight(trimmed, "/")
	if trimmed == "" {
		return ""
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return trimmed
	}

	if parsed.Scheme != "" {
		parsed.Scheme = strings.ToLower(parsed.Scheme)
	}

	if parsed.Host != "" {
		parsed.Host = strings.ToLower(parsed.Host)
	}

	if parsed.Path != "" {
		parsed.Path = strings.TrimRight(parsed.Path, "/")
	}

	return strings.TrimRight(parsed.String(), "/")
}
