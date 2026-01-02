package main

import "net/http"

func (app *application) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(app.swagger)
}
