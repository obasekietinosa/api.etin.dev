package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

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
	for _, trustedOrigin := range app.config.cors.trustedOrigins {
		if trustedOrigin == "*" {
			return "*", true
		}
		if trustedOrigin == origin {
			return origin, true
		}
	}

	return "", false
}
