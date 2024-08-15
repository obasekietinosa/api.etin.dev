package main

import (
	"encoding/json"
	"fmt"
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
			roles := []data.Role{
				{
					ID:          1,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					DeletedAt:   time.Time{},
					StartDate:   time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Time{},
					Title:       "Technical Lead",
					Subtitle:    "Content Management",
					Company:     "Accurx",
					CompanyIcon: "/accurx.png",
					Slug:        "tech-lead-accurx-w87xbv9402",
					Description: "I'm currently technical lead for the Content Management team at Accurx. Our remit involves maintaining a broad spectrum of existing functionality as well as evolving...",
					Skills:      []string{"JavaScript", "React", "Typescript", "C#"},
				},
				{
					ID:          2,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					DeletedAt:   time.Time{},
					StartDate:   time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Time{},
					Title:       "Software Engineer",
					Subtitle:    "",
					Company:     "Proper",
					CompanyIcon: "/proper.png",
					Slug:        "software-engineer-proper-k921b2j0",
					Description: "At Proper, I was one of 3 engineers who worked to launch our digital sleep improvement tools as well as maintain our retail website. On a small team, I worked as a full-stack software engineer, delivering on the frontend in React and on the backend in Nest.js",
					Skills:      []string{"React", "Typescript", "Node.js", "Nest.js", "Headless Content Management"},
				},
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
				Company     string    `json:"company"`
				CompanyIcon string    `json:"companyIcon"`
				Description string    `json:"description"`
				Skills      []string  `json:"skills"`
			}

			err := app.readJSON(w, r, &input)
			if err != nil {
				app.writeError(w, http.StatusBadRequest)
				return
			}

			fmt.Fprintf(w, "%v\n", input)
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

	role := data.Role{
		ID:          id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   time.Time{},
		StartDate:   time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Time{},
		Title:       "Technical Lead",
		Subtitle:    "Content Management",
		Company:     "Accurx",
		CompanyIcon: "/accurx.png",
		Slug:        "tech-lead-accurx-w87xbv9402",
		Description: "I'm currently technical lead for the Content Management team at Accurx. Our remit involves maintaining a broad spectrum of existing functionality as well as evolving...",
		Skills:      []string{"JavaScript", "React", "Typescript", "C#"},
	}

	app.writeJSON(w, http.StatusOK, envelope{"role": role})
}

func (app *application) updateRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/roles/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	var input struct {
		StartDate   *time.Time `json:"startDate"`
		EndDate     *time.Time `json:"endDate"`
		Title       *string    `json:"title"`
		Subtitle    *string    `json:"subtitle"`
		Company     *string    `json:"company"`
		CompanyIcon *string    `json:"companyIcon"`
		Description *string    `json:"description"`
		Skills      []string   `json:"skills"`
	}

	role := data.Role{
		ID:          id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   time.Time{},
		StartDate:   time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Time{},
		Title:       "Technical Lead",
		Subtitle:    "Content Management",
		Company:     "Accurx",
		CompanyIcon: "/accurx.png",
		Slug:        "tech-lead-accurx-w87xbv9402",
		Description: "I'm currently technical lead for the Content Management team at Accurx. Our remit involves maintaining a broad spectrum of existing functionality as well as evolving...",
		Skills:      []string{"JavaScript", "React", "Typescript", "C#"},
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
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

	if input.Company != nil {
		role.Company = *input.Company
	}

	if input.CompanyIcon != nil {
		role.CompanyIcon = *input.CompanyIcon
	}

	if input.Description != nil {
		role.Description = *input.Description
	}

	if len(input.Skills) > 0 {
		role.Skills = input.Skills
	}

	app.writeJSON(w, http.StatusOK, envelope{"role": role})
}

func (app *application) deleteRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Path[len("/v1/roles/"):], 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Deleting role with ID: %d", id)
}
