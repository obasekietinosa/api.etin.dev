package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"api.etin.dev/internal/data"
)

func (app *application) getCreateItemNotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		itemNotes, err := app.models.ItemNotes.GetAll()
		if err != nil {
			app.logger.Printf("Error retrieving item note associations: %s", err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}

		app.writeJSON(w, http.StatusOK, envelope{"itemNotes": itemNotes})
	case http.MethodPost:
		if !app.isRequestAuthenticated(r) {
			app.writeError(w, http.StatusForbidden)
			return
		}

		var input struct {
			NoteID   int64  `json:"noteId"`
			ItemID   int64  `json:"itemId"`
			ItemType string `json:"itemType"`
		}

		if err := app.readJSON(w, r, &input); err != nil {
			app.logger.Printf("Could not parse item note association payload: %s", err)
			app.writeError(w, http.StatusBadRequest)
			return
		}

		if input.NoteID < 1 || input.ItemID < 1 || strings.TrimSpace(input.ItemType) == "" {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		itemNote := &data.ItemNote{
			NoteID:   input.NoteID,
			ItemID:   input.ItemID,
			ItemType: data.ItemType(strings.ToLower(input.ItemType)),
		}

		if err := app.models.ItemNotes.Insert(itemNote); err != nil {
			if errors.Is(err, data.ErrInvalidItemType) {
				app.writeError(w, http.StatusBadRequest)
				return
			}

			app.logger.Printf("Could not create item note association: %s", err)
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.writeJSON(w, http.StatusCreated, envelope{"itemNote": itemNote})
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getUpdateDeleteItemNotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getItemNote(w, r)
	case http.MethodPut:
		app.updateItemNote(w, r)
	case http.MethodDelete:
		app.deleteItemNote(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getItemNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/item-notes/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemNote, err := app.models.ItemNotes.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve item note association %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"itemNote": itemNote})
}

func (app *application) updateItemNote(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/v1/item-notes/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemNote, err := app.models.ItemNotes.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve item note association %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	var input struct {
		NoteID   *int64  `json:"noteId"`
		ItemID   *int64  `json:"itemId"`
		ItemType *string `json:"itemType"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.logger.Printf("Could not parse item note association update payload: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	if input.NoteID != nil {
		itemNote.NoteID = *input.NoteID
	}

	if input.ItemID != nil {
		itemNote.ItemID = *input.ItemID
	}

	if input.ItemType != nil {
		itemNote.ItemType = data.ItemType(strings.ToLower(*input.ItemType))
	}

	if err := app.models.ItemNotes.Update(itemNote); err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.logger.Printf("Could not update item note association %d: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"itemNote": itemNote})
}

func (app *application) deleteItemNote(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/v1/item-notes/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	if err := app.models.ItemNotes.Delete(id); err != nil {
		app.logger.Printf("Could not delete item note association %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusNoContent, envelope{"itemNote": nil})
}

func (app *application) getNotesForItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/v1/item-notes/items/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemType := data.ItemType(strings.ToLower(parts[0]))

	filters := data.CursorFilters{
		Limit: 20,
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

		app.logger.Printf("Could not retrieve notes for %s item %d: %s", parts[0], itemID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"notes": notes, "metadata": metadata})
}
