package main

import (
	"net/http"
	"strings"
)

func (app *application) adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Printf("admin login: could not parse credentials: %v", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	expectedEmail := strings.ToLower(app.config.adminEmail)

	if email == "" || !secureCompare(email, expectedEmail) || !secureCompare(input.Password, app.config.adminPassword) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	token, expiresAt, err := app.sessions.create()
	if err != nil {
		app.logger.Printf("admin login: could not create session: %v", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{
		"token":     token,
		"expiresAt": expiresAt,
	})
}

func (app *application) adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	token, err := parseBearerToken(r.Header.Get("Authorization"))
	if err != nil || !app.sessions.validate(token) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	app.sessions.revoke(token)
	app.writeJSON(w, http.StatusOK, envelope{"message": "logged out"})
}
