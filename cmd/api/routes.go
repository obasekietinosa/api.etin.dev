package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthcheck", app.healthcheck)

	mux.HandleFunc("/v1/roles", app.getCreateRolesHandler)
	mux.HandleFunc("/v1/roles/", app.getUpdateDeleteRolesHandler)

	mux.HandleFunc("/v1/companies", app.getCreateCompaniesHandler)
	mux.HandleFunc("/v1/companies/", app.getUpdateDeleteCompaniesHandler)

	mux.HandleFunc("/v1/notes", app.getCreateNotesHandler)
	mux.HandleFunc("/v1/notes/", app.getUpdateDeleteNotesHandler)
	mux.HandleFunc("/v1/item-notes", app.getCreateItemNotesHandler)
	mux.HandleFunc("/v1/item-notes/", app.getUpdateDeleteItemNotesHandler)
	mux.HandleFunc("/v1/item-notes/items/", app.getNotesForItemHandler)

	mux.HandleFunc("/v1/projects", app.getCreateProjectsHandler)
	mux.HandleFunc("/v1/projects/", app.getUpdateDeleteProjectsHandler)

	return mux
}
