package main

import "net/http"

const adminOrigin = "https://admin.etin.dev"

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

               w.Header().Set("Access-Control-Allow-Origin", adminOrigin)
               w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
               w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

               if r.Method == http.MethodOptions {
                       w.WriteHeader(http.StatusNoContent)
                       return
               }

		next.ServeHTTP(w, r)
	})
}
