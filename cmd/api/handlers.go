package main

import (
	"encoding/json"
	"net/http"

	"api.etin.dev/internal/version"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version.Number,
	}
	j, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	j = append(j, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
