package main

import (
	"fmt"
	"net/http"
	"strconv"

	"api.etin.dev/internal/data"
)

func (app *application) getTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := app.models.Tags.GetAll()
	if err != nil {
		app.logger.Printf("Error: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"tags": tags})
}

func (app *application) createTagHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}
	var input struct {
		Name  string  `json:"name"`
		Slug  string  `json:"slug"`
		Icon  *string `json:"icon"`
		Theme *string `json:"theme"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Print(err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	tag := &data.Tag{
		Name:  input.Name,
		Slug:  input.Slug,
		Icon:  input.Icon,
		Theme: input.Theme,
	}
	err = app.models.Tags.Insert(tag)
	if err != nil {
		app.logger.Printf("Error: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{"tag": tag})
}

func (app *application) getTagHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}
	tag, err := app.models.Tags.Get(id)
	if err != nil {
		app.logger.Printf("A problem fetching tag id: %d Error: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"tag": tag})
}

func (app *application) updateTagHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}
	tag, err := app.models.Tags.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve model. Error: %s", err)
		app.writeError(w, http.StatusNotFound)
		return
	}
	var input struct {
		Name  *string `json:"name"`
		Slug  *string `json:"slug"`
		Icon  *string `json:"icon"`
		Theme *string `json:"theme"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Printf("Could not parse input. Error: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	if input.Name != nil {
		tag.Name = *input.Name
	}
	if input.Slug != nil {
		tag.Slug = *input.Slug
	}
	if input.Icon != nil {
		tag.Icon = input.Icon
	}
	if input.Theme != nil {
		tag.Theme = input.Theme
	}
	err = app.models.Tags.Update(tag)
	if err != nil {
		app.logPostgresError(fmt.Sprintf("Could not update tag %d", id), err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"tag": tag})
}

func (app *application) deleteTagHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}
	err = app.models.Tags.Delete(id)
	if err != nil {
		app.writeError(w, http.StatusNotFound)
		return
	}
	app.writeJSON(w, http.StatusNoContent, envelope{"tag": nil})
}
