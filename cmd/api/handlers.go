package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"api.etin.dev/internal/data"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	j, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	j = append(j, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (app *application) getCreateRolesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			roles, err := app.models.Roles.GetAll()

			if err != nil {
				app.logger.Printf("Error: %s", err)
				app.writeError(w, http.StatusInternalServerError)
				return
			}

			app.writeJSON(w, http.StatusOK, envelope{"roles": roles})
			return
		}
	case http.MethodPost:
		{
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

			err = app.models.Roles.Insert(role)
			if err != nil {
				app.logger.Printf("Error: %s", err)
				app.writeError(w, http.StatusBadRequest)
				return
			}

			app.writeJSON(w, http.StatusCreated, envelope{"role": role})
		}
	}
}

func (app *application) getUpdateDeleteRolesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			app.getRole(w, r)
		}

	case http.MethodPut:
		{
			app.updateRole(w, r)
		}

	case http.MethodDelete:
		{
			app.deleteRole(w, r)
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/roles/"):], 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	role, err := app.models.Roles.Get(id)
	if err != nil {
		app.logger.Printf("A problem fetching roleid: %d Error: %s", id, err)
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"role": role})
}

func (app *application) updateRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/roles/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	role, err := app.models.Roles.Get(id)
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

	err = app.models.Roles.Update(role)
	if err != nil {
		app.logger.Printf("Could not update role: %d. Error: %s", id, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"role": role})
}

func (app *application) deleteRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/roles/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	err = app.models.Roles.Delete(id)
	if err != nil {
		app.writeError(w, http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusNoContent, envelope{"role": nil})
}

func (app *application) getCreateCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			companies, err := app.models.Companies.GetAll()
			if err != nil {
				app.logger.Printf("Error getting all companies, error: %s", err)
				app.writeError(w, http.StatusInternalServerError)
				return
			}

			app.writeJSON(w, http.StatusOK, envelope{"companies": companies})
			return
		}
	case http.MethodPost:
		{
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
				Icon:        input.Icon,
				Description: input.Description,
			}

			err = app.models.Companies.Insert(company)
			if err != nil {
				app.logger.Print(err)
				app.writeError(w, http.StatusBadRequest)
				return
			}
			app.writeJSON(w, http.StatusOK, envelope{"company": company})
			return
		}
	}
}

func (app *application) getUpdateDeleteCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/companies/"):], 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		{
			company, err := app.models.Companies.Get(id)
			if err != nil {
				app.logger.Printf("Error getting company with ID: %d, error: %s", id, err)
				app.writeError(w, http.StatusInternalServerError)
				return
			}

			app.writeJSON(w, http.StatusOK, envelope{"company": company})
			return
		}
	case http.MethodPut:
		{
			company, err := app.models.Companies.Get(id)
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
				company.Icon = *input.Icon
			}

			if input.Description != nil {
				company.Description = *input.Description
			}

			err = app.models.Companies.Update(company)
			if err != nil {
				app.logger.Printf("Could not update company with ID: %d. Error: %s", id, err)
				app.writeError(w, http.StatusBadRequest)
				return
			}
			app.writeJSON(w, http.StatusAccepted, envelope{"company": company})
			return
		}
	case http.MethodDelete:
		{
			company, err := app.models.Companies.Get(id)
			if err != nil {
				app.logger.Printf("Error getting company with ID: %d, error: %s", id, err)
				app.writeError(w, http.StatusInternalServerError)
				return
			}
			err = app.models.Companies.Delete(company)
			if err != nil {
				app.logger.Printf("Could not update company with ID: %d. Error: %s", id, err)
				app.writeError(w, http.StatusBadRequest)
				return
			}
			app.writeJSON(w, http.StatusNoContent, nil)
			return

		}
	}
}
