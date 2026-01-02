package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"api.etin.dev/internal/data"
)

func (app *application) getCreateContentNoteHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	contentTypeStr := r.PathValue("contentType")
	itemIDStr := r.PathValue("id")

	if contentTypeStr == "" || itemIDStr == "" {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemType := data.ItemType(strings.ToLower(contentTypeStr))

	var input struct {
		Title       string     `json:"title"`
		Subtitle    string     `json:"subtitle"`
		Body        string     `json:"body"`
		PublishedAt *string    `json:"publishedAt"` // Allow optional publishedAt
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.logger.Printf("Could not parse content note payload: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	note := &data.Note{
		Title:    input.Title,
		Subtitle: input.Subtitle,
		Body:     input.Body,
	}

	if input.PublishedAt != nil {
		t, err := time.Parse(time.RFC3339, *input.PublishedAt)
		if err != nil {
			app.writeError(w, http.StatusBadRequest)
			return
		}
		note.PublishedAt = &t
	}

	// Create the note
	if err := app.models.Notes.Insert(note); err != nil {
		app.logger.Printf("Could not create note: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	// Link the note to the item
	itemNote := &data.ItemNote{
		NoteID:   note.ID,
		ItemID:   itemID,
		ItemType: itemType,
	}

	if err := app.models.ItemNotes.Insert(itemNote); err != nil {
		// Rollback note creation? Or just log?
		// Since we don't have transaction helper readily available across models without passing DB,
		// we should probably try to delete the note.
		_ = app.models.Notes.Delete(note.ID)

		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.logger.Printf("Could not create item note association: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"note": note})
}

func (app *application) getContentNotesHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	contentTypeStr := r.PathValue("contentType")
	itemIDStr := r.PathValue("id")

	if contentTypeStr == "" || itemIDStr == "" {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemType := data.ItemType(strings.ToLower(contentTypeStr))

	filters := data.CursorFilters{
		Limit:  20, // Default limit
	}

	qs := r.URL.Query()
	if cursor := qs.Get("cursor"); cursor != "" {
		filters.Cursor = cursor
	}
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err == nil && l > 0 {
			filters.Limit = l
		}
	}

	notes, metadata, err := app.models.ItemNotes.GetNotesForItem(itemType, itemID, filters)
	if err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}
		app.logger.Printf("Error fetching notes for item: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"notes": notes, "metadata": metadata})
}

func (app *application) getAllContentNotesHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	contentTypeStr := r.PathValue("contentType")
	if contentTypeStr == "" {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemType := data.ItemType(strings.ToLower(contentTypeStr))

	filters := data.CursorFilters{
		Limit:  20,
	}

	qs := r.URL.Query()
	if cursor := qs.Get("cursor"); cursor != "" {
		filters.Cursor = cursor
	}
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err == nil && l > 0 {
			filters.Limit = l
		}
	}

	notes, metadata, err := app.models.ItemNotes.GetNotesForContentType(itemType, filters)
	if err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}
		app.logger.Printf("Error fetching notes for content type: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"notes": notes, "metadata": metadata})
}
