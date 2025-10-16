package main

import (
	"net/http"
	"strconv"
	"time"

	"api.etin.dev/internal/data"
)

func (app *application) getCreateProjectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		projects, err := app.models.Projects.GetAll()
		if err != nil {
			app.logger.Printf("Error retrieving projects. Error: %s", err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}

		app.writeJSON(w, http.StatusOK, envelope{"projects": projects})
	case http.MethodPost:
		if !app.isRequestAuthenticated(r) {
			app.writeError(w, http.StatusForbidden)
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

		err = app.models.Projects.Insert(project)
		if err != nil {
			app.logger.Printf("Error creating project. Error: %s", err)
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.writeJSON(w, http.StatusCreated, envelope{"project": project})
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getUpdateDeleteProjectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getProject(w, r)
	case http.MethodPut:
		app.updateProject(w, r)
	case http.MethodDelete:
		app.deleteProject(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/projects/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	project, err := app.models.Projects.Get(id)
	if err != nil {
		app.logger.Printf("Error retrieving project with ID %d. Error: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"project": project})
}

func (app *application) updateProject(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/v1/projects/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	project, err := app.models.Projects.Get(id)
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

	err = app.models.Projects.Update(project)
	if err != nil {
		app.logger.Printf("Error updating project with ID %d. Error: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"project": project})
}

func (app *application) deleteProject(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusForbidden)
		return
	}

	id, err := strconv.ParseInt(r.URL.Path[len("/v1/projects/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	err = app.models.Projects.Delete(id)
	if err != nil {
		app.logger.Printf("Error deleting project with ID %d. Error: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusNoContent, envelope{"project": nil})
}
