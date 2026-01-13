package main

import (
	"api.etin.dev/internal/data"
	"net/http"
	"strconv"
)

func (app *application) getCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	companies, err := app.getModels(r).Companies.GetAll()
	if err != nil {
		app.logger.Printf("Error getting all companies, error: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"companies": companies})
	return
}

func (app *application) createCompanyHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	var input struct {
		Name        string `json:"name"`
		Icon        string `json:"icon"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Print(err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	company := &data.Company{
		Name:        input.Name,
		Icon:        &input.Icon,
		Description: &input.Description,
	}

	err = app.getModels(r).Companies.Insert(company)
	if err != nil {
		app.logger.Print(err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"company": company})
	return
}

func (app *application) getCompanyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	company, err := app.getModels(r).Companies.Get(id)
	if err != nil {
		app.logger.Printf("Error getting company with ID: %d, error: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"company": company})
	return
}

func (app *application) updateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	company, err := app.getModels(r).Companies.Get(id)
	if err != nil {
		app.logger.Printf("Error getting company with ID: %d, error: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Icon        *string `json:"icon"`
		Description *string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Printf("Could not parse input. Error: %s", err)
		app.writeError(w, http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		company.Name = *input.Name
	}

	if input.Icon != nil {
		company.Icon = input.Icon
	}

	if input.Description != nil {
		company.Description = input.Description
	}

	err = app.getModels(r).Companies.Update(company)
	if err != nil {
		app.logger.Printf("Could not update company with ID: %d. Error: %s", id, err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusAccepted, envelope{"company": company})
	return
}

func (app *application) deleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	company, err := app.getModels(r).Companies.Get(id)
	if err != nil {
		app.logger.Printf("Error getting company with ID: %d, error: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}
	err = app.getModels(r).Companies.Delete(company)
	if err != nil {
		app.logger.Printf("Could not delete company with ID: %d. Error: %s", id, err)
		app.writeError(w, http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusNoContent, nil)
	return
}
