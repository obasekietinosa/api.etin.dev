package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"api.etin.dev/internal/data"
)

func (app *application) getCreateTagItemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tagItems, err := app.models.TagItems.GetAll()
		if err != nil {
			app.logger.Printf("Error retrieving tag associations: %s", err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}

		app.writeJSON(w, http.StatusOK, envelope{"taggedItems": tagItems})
	case http.MethodPost:
		if !app.isRequestAuthenticated(r) {
			app.writeError(w, http.StatusForbidden)
			return
		}

		var input struct {
			TagID    int64  `json:"tagId"`
			ItemID   int64  `json:"itemId"`
			ItemType string `json:"itemType"`
		}

		if err := app.readJSON(w, r, &input); err != nil {
			app.logger.Printf("Could not parse tag association payload: %s", err)
			app.writeError(w, http.StatusBadRequest)
			return
		}

		if input.TagID < 1 || input.ItemID < 1 || strings.TrimSpace(input.ItemType) == "" {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		tagItem := &data.TagItem{
			TagID:    input.TagID,
			ItemID:   input.ItemID,
			ItemType: data.ItemType(strings.ToLower(input.ItemType)),
		}

		if err := app.models.TagItems.Insert(tagItem); err != nil {
			if errors.Is(err, data.ErrInvalidItemType) {
				app.writeError(w, http.StatusBadRequest)
				return
			}

			app.logger.Printf("Could not create tag association: %s", err)
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.writeJSON(w, http.StatusCreated, envelope{"taggedItem": tagItem})
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getUpdateDeleteTagItemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getTagItem(w, r)
	case http.MethodPut:
		app.updateTagItem(w, r)
	case http.MethodDelete:
		app.deleteTagItem(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getTagItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/tagged-items/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	tagItem, err := app.models.TagItems.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve tag association %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"taggedItem": tagItem})
}

func (app *application) updateTagItem(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/v1/tagged-items/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	tagItem, err := app.models.TagItems.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve tag association %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	var input struct {
		TagID    *int64  `json:"tagId"`
		ItemID   *int64  `json:"itemId"`
		ItemType *string `json:"itemType"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.logger.Printf("Could not parse tag association update payload: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	if input.TagID != nil {
		tagItem.TagID = *input.TagID
	}

	if input.ItemID != nil {
		tagItem.ItemID = *input.ItemID
	}

	if input.ItemType != nil {
		tagItem.ItemType = data.ItemType(strings.ToLower(*input.ItemType))
	}

	if err := app.models.TagItems.Update(tagItem); err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.logger.Printf("Could not update tag association %d: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"taggedItem": tagItem})
}

func (app *application) deleteTagItem(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/v1/tagged-items/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	if err := app.models.TagItems.Delete(id); err != nil {
		app.logger.Printf("Could not delete tag association %d: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusNoContent, envelope{"taggedItem": nil})
}

func (app *application) getTagsForItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/v1/tagged-items/items/")
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

	tags, err := app.models.TagItems.GetTagsForItem(itemType, itemID)
	if err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.logger.Printf("Could not retrieve tags for %s item %d: %s", parts[0], itemID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"tags": tags})
}
