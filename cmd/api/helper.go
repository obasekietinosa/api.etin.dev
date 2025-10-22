package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope) {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		app.writeError(w, http.StatusInternalServerError)
	}

	jsonData = append(jsonData, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
	return
}

func (app *application) writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain a single JSON object")
	}

	return nil
}

func (app *application) isRequestAuthenticated(r *http.Request) bool {
	token, err := parseBearerToken(r.Header.Get("Authorization"))
	if err != nil {
		return false
	}

	return app.sessions.validate(token)
}

func (app *application) logPostgresError(context string, err error) {
	if app.logger == nil || err == nil {
		return
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		app.logger.Printf("%s: %v", context, err)
		return
	}

	details := []string{string(pqErr.Code)}
	if pqErr.Table != "" {
		details = append(details, "table="+pqErr.Table)
	}
	if pqErr.Column != "" {
		details = append(details, "column="+pqErr.Column)
	}
	if pqErr.Constraint != "" {
		details = append(details, "constraint="+pqErr.Constraint)
	}
	if pqErr.Detail != "" {
		details = append(details, "detail="+pqErr.Detail)
	}
	if pqErr.Where != "" {
		details = append(details, "where="+pqErr.Where)
	}

	app.logger.Printf("%s: %s (%s)", context, pqErr.Message, strings.Join(details, ", "))

	if pqErr.InternalQuery != "" {
		app.logger.Printf("%s: internal query=%s", context, pqErr.InternalQuery)
	}
}
