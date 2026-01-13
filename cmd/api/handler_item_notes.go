package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"api.etin.dev/internal/data"
)

func (app *application) getItemNotesHandler(w http.ResponseWriter, r *http.Request) {
	itemNotes, err := app.getModels(r).ItemNotes.GetAll()
	if err != nil {
		app.logger.Printf("Error retrieving item note associations: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"itemNotes": itemNotes})
}

func (app *application) createItemNoteHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
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
		ItemType: strings.ToLower(input.ItemType),
	}

	if err := app.getModels(r).ItemNotes.Insert(itemNote); err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.logger.Printf("Could not create item note association: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"itemNote": itemNote})
}

func (app *application) getItemNoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemNote, err := app.getModels(r).ItemNotes.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve item note association %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"itemNote": itemNote})
}

func (app *application) updateItemNoteHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemNote, err := app.getModels(r).ItemNotes.Get(id)
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
		itemNote.ItemType = strings.ToLower(*input.ItemType)
	}

	if err := app.getModels(r).ItemNotes.Update(itemNote); err != nil {
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

func (app *application) deleteItemNoteHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	if err := app.getModels(r).ItemNotes.Delete(id); err != nil {
		app.logger.Printf("Could not delete item note association %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusNoContent, envelope{"itemNote": nil})
}

func (app *application) getNotesForItemHandler(w http.ResponseWriter, r *http.Request) {
	// r.PathValue("itemType") and r.PathValue("itemId") should be used now.
	// But let's check how the route was defined:
	// mux.HandleFunc("GET /v1/item-notes/items/{itemType}/{itemId}", app.getNotesForItemHandler)

	itemTypeStr := r.PathValue("itemType")
	itemIdStr := r.PathValue("itemId")

	if itemTypeStr == "" || itemIdStr == "" {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemType := data.ItemType(strings.ToLower(itemTypeStr))

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

	notes, metadata, err := app.getModels(r).ItemNotes.GetNotesForItem(string(itemType), itemID, filters)
	if err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.logger.Printf("Could not retrieve notes for %s item %d: %s", itemTypeStr, itemID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"notes": notes, "metadata": metadata})
}
