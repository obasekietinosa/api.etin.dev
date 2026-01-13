package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"api.etin.dev/internal/data"
)

func (app *application) getRolesHandler(w http.ResponseWriter, r *http.Request) {
	roles, err := app.getModels(r).Roles.GetAll()
	if err != nil {
		app.logger.Printf("Error: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"roles": roles})
}

func (app *application) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}
	var input struct {
		StartDate   time.Time `json:"startDate"`
		EndDate     time.Time `json:"endDate"`
		Title       string    `json:"title"`
		Subtitle    string    `json:"subtitle"`
		CompanyId   int64     `json:"companyId"`
		Description string    `json:"description"`
		Skills      []string  `json:"skills"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Print(err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	role := &data.Role{
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Title:       input.Title,
		Subtitle:    input.Subtitle,
		CompanyId:   input.CompanyId,
		Description: input.Description,
		Skills:      input.Skills,
	}
	err = app.getModels(r).Roles.Insert(role)
	if err != nil {
		app.logger.Printf("Error: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{"role": role})
}

func (app *application) getRoleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	role, err := app.getModels(r).Roles.Get(id)
	if err != nil {
		app.logger.Printf("A problem fetching roleid: %d Error: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"role": role})
}

func (app *application) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}
	role, err := app.getModels(r).Roles.Get(id)
	if err != nil {
		app.logger.Printf("Could not retrieve model. Error: %s", err)
		app.writeError(w, http.StatusNotFound)
		return
	}
	var input struct {
		StartDate   *time.Time `json:"startDate"`
		EndDate     *time.Time `json:"endDate"`
		Title       *string    `json:"title"`
		Subtitle    *string    `json:"subtitle"`
		CompanyId   *int64     `json:"companyId"`
		Description *string    `json:"description"`
		Skills      []string   `json:"skills"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Printf("Could not parse input. Error: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	if !input.StartDate.IsZero() {
		role.StartDate = *input.StartDate
	}
	if !input.EndDate.IsZero() {
		role.EndDate = *input.EndDate
	}
	if input.Title != nil {
		role.Title = *input.Title
	}
	if input.Subtitle != nil {
		role.Subtitle = *input.Subtitle
	}
	if input.CompanyId != nil {
		role.CompanyId = *input.CompanyId
	}
	if input.Description != nil {
		role.Description = *input.Description
	}
	if len(input.Skills) > 0 {
		role.Skills = input.Skills
	}
	err = app.getModels(r).Roles.Update(role)
	if err != nil {
		app.logPostgresError(fmt.Sprintf("Could not update role %d", id), err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"role": role})
}

func (app *application) deleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}
	err = app.getModels(r).Roles.Delete(id)
	if err != nil {
		app.writeError(w, http.StatusNotFound)
		return
	}
	app.writeJSON(w, http.StatusNoContent, envelope{"role": nil})
}
