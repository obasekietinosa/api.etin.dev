package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger", app.swaggerHandler)
	mux.HandleFunc("/v1/healthcheck", app.healthcheck)
	mux.HandleFunc("/v1/admin/login", app.adminLoginHandler)
	mux.HandleFunc("/v1/admin/logout", app.adminLogoutHandler)

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

	mux.HandleFunc("/v1/tagged-items", app.getCreateTagItemsHandler)
	mux.HandleFunc("/v1/tagged-items/", app.getUpdateDeleteTagItemsHandler)
	mux.HandleFunc("/v1/tagged-items/items/", app.getTagsForItemHandler)

	return app.enableCORS(mux)
}
