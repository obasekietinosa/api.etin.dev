package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger", app.swaggerHandler)
	mux.HandleFunc("/public/v1/notes", app.getPublicNotesHandler)
	mux.HandleFunc("GET /public/v1/{contentType}/{id}/notes", app.getPublicNotesForContentHandler)
	mux.HandleFunc("GET /public/v1/{contentType}/notes", app.getPublicAllNotesForContentHandler)
	mux.HandleFunc("/public/v1/projects", app.getPublicProjectsHandler)
	mux.HandleFunc("/public/v1/roles", app.getPublicRolesHandler)
	mux.HandleFunc("/v1/healthcheck", app.healthcheck)
	mux.HandleFunc("/v1/admin/login", app.adminLoginHandler)
	mux.HandleFunc("/v1/admin/logout", app.adminLogoutHandler)

	mux.HandleFunc("/v1/assets", app.getCreateAssetsHandler)

	mux.Handle("/v1/roles", app.deployWebhook(http.HandlerFunc(app.getCreateRolesHandler)))
	mux.Handle("/v1/roles/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteRolesHandler)))

	mux.Handle("/v1/companies", app.deployWebhook(http.HandlerFunc(app.getCreateCompaniesHandler)))
	mux.Handle("/v1/companies/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteCompaniesHandler)))

	mux.Handle("/v1/notes", app.deployWebhook(http.HandlerFunc(app.getCreateNotesHandler)))
	mux.Handle("/v1/notes/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteNotesHandler)))
	mux.Handle("/v1/item-notes", app.deployWebhook(http.HandlerFunc(app.getCreateItemNotesHandler)))
	mux.Handle("/v1/item-notes/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteItemNotesHandler)))
	mux.HandleFunc("/v1/item-notes/items/", app.getNotesForItemHandler)

	mux.HandleFunc("POST /v1/{contentType}/{id}/notes", app.getCreateContentNoteHandler)
	mux.HandleFunc("GET /v1/{contentType}/{id}/notes", app.getContentNotesHandler)
	mux.HandleFunc("GET /v1/{contentType}/notes", app.getAllContentNotesHandler)

	mux.Handle("/v1/projects", app.deployWebhook(http.HandlerFunc(app.getCreateProjectsHandler)))
	mux.Handle("/v1/projects/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteProjectsHandler)))

	mux.Handle("/v1/tagged-items", app.deployWebhook(http.HandlerFunc(app.getCreateTagItemsHandler)))
	mux.Handle("/v1/tagged-items/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteTagItemsHandler)))
	mux.HandleFunc("/v1/tagged-items/items/", app.getTagsForItemHandler)

	return app.enableCORS(mux)
}
