package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"api.etin.dev/internal/data"
)

func (app *application) getTagItemsHandler(w http.ResponseWriter, r *http.Request) {
	tagItems, err := app.models.TagItems.GetAll()
	if err != nil {
		app.logger.Printf("Error retrieving tag associations: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"taggedItems": tagItems})
}

func (app *application) createTagItemHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
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
}

func (app *application) getTagItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
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

func (app *application) updateTagItemHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
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

func (app *application) deleteTagItemHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
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
	// Replaced manual path parsing with r.PathValue
	// Pattern: GET /v1/tagged-items/items/{itemType}/{itemId}

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

	tags, err := app.models.TagItems.GetTagsForItem(itemType, itemID)
	if err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.logger.Printf("Could not retrieve tags for %s item %d: %s", itemTypeStr, itemID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"tags": tags})
}
