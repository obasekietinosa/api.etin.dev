package main

import (
	"net/http"
	"strconv"
	"time"

	"api.etin.dev/internal/data"
)

func (app *application) getCreateNotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		notes, err := app.models.Notes.GetAll()
		if err != nil {
			app.logger.Printf("Error retrieving notes: %s", err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}

		app.writeJSON(w, http.StatusOK, envelope{"notes": notes})
	case http.MethodPost:
		if !app.isRequestAuthenticated(r) {
			app.writeError(w, http.StatusForbidden)
			return
		}

		var input struct {
			Title       string     `json:"title"`
			Subtitle    string     `json:"subtitle"`
			Body        string     `json:"body"`
			PublishedAt *time.Time `json:"publishedAt"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			app.logger.Printf("Could not parse request body: %s", err)
			app.writeError(w, http.StatusBadRequest)
			return
		}

		if input.Title == "" {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		var publishedAt *time.Time
		if input.PublishedAt != nil {
			t := *input.PublishedAt
			publishedAt = &t
		}

		note := &data.Note{
			Title:       input.Title,
			Subtitle:    input.Subtitle,
			Body:        input.Body,
			PublishedAt: publishedAt,
		}

		err = app.models.Notes.Insert(note)
		if err != nil {
			app.logger.Printf("Could not create note: %s", err)
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.writeJSON(w, http.StatusCreated, envelope{"note": note})
		app.triggerDeployWebhook()
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getUpdateDeleteNotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getNote(w, r)
	case http.MethodPut:
		app.updateNote(w, r)
	case http.MethodDelete:
		app.deleteNote(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/notes/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	note, err := app.models.Notes.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve note %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"note": note})
}

func (app *application) updateNote(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/v1/notes/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	note, err := app.models.Notes.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve note %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	var input struct {
		Title       *string    `json:"title"`
		Subtitle    *string    `json:"subtitle"`
		Body        *string    `json:"body"`
		PublishedAt *time.Time `json:"publishedAt"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Printf("Could not parse request body: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	if input.Title != nil {
		note.Title = *input.Title
	}

	if input.Subtitle != nil {
		note.Subtitle = *input.Subtitle
	}

	if input.Body != nil {
		note.Body = *input.Body
	}

	if input.PublishedAt != nil {
		t := *input.PublishedAt
		note.PublishedAt = &t
	}

	err = app.models.Notes.Update(note)
	if err != nil {
		app.logger.Printf("Could not update note %d: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"note": note})
	app.triggerDeployWebhook()
}

func (app *application) deleteNote(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/v1/notes/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	err = app.models.Notes.Delete(id)
	if err != nil {
		app.logger.Printf("Could not delete note %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusNoContent, envelope{"note": nil})
	app.triggerDeployWebhook()
}
