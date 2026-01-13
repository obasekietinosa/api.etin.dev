package main

import (
	"net/http"
	"strconv"
	"time"

	"api.etin.dev/internal/data"
)

func (app *application) getProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := app.getModels(r).Projects.GetAll()
	if err != nil {
		app.logger.Printf("Error retrieving projects. Error: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"projects": projects})
}

func (app *application) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	var input struct {
		StartDate   time.Time  `json:"startDate"`
		EndDate     *time.Time `json:"endDate"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		ImageURL    *string    `json:"imageUrl"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Printf("Error parsing project payload. Error: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	project := &data.Project{
		StartDate:   input.StartDate,
		Title:       input.Title,
		Description: input.Description,
		EndDate:     input.EndDate,
		ImageURL:    input.ImageURL,
	}

	err = app.getModels(r).Projects.Insert(project)
	if err != nil {
		app.logger.Printf("Error creating project. Error: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"project": project})
}

func (app *application) getProjectHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	project, err := app.getModels(r).Projects.Get(id)
	if err != nil {
		app.logger.Printf("Error retrieving project with ID %d. Error: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"project": project})
}

func (app *application) updateProjectHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	project, err := app.getModels(r).Projects.Get(id)
	if err != nil {
		app.logger.Printf("Error retrieving project with ID %d for update. Error: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	var input struct {
		StartDate   *time.Time `json:"startDate"`
		EndDate     *time.Time `json:"endDate"`
		Title       *string    `json:"title"`
		Description *string    `json:"description"`
		ImageURL    *string    `json:"imageUrl"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Printf("Error parsing project update payload. Error: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	if input.StartDate != nil {
		project.StartDate = *input.StartDate
	}

	if input.EndDate != nil {
		project.EndDate = input.EndDate
	}

	if input.Title != nil {
		project.Title = *input.Title
	}

	if input.Description != nil {
		project.Description = *input.Description
	}

	if input.ImageURL != nil {
		project.ImageURL = input.ImageURL
	}

	err = app.getModels(r).Projects.Update(project)
	if err != nil {
		app.logger.Printf("Error updating project with ID %d. Error: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"project": project})
}

func (app *application) deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	err = app.getModels(r).Projects.Delete(id)
	if err != nil {
		app.logger.Printf("Error deleting project with ID %d. Error: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusNoContent, envelope{"project": nil})
}
