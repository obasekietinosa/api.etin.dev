package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger", app.swaggerHandler)
	mux.HandleFunc("/public/v1/notes", app.getPublicNotesHandler)
	mux.HandleFunc("/public/v1/projects", app.getPublicProjectsHandler)
	mux.HandleFunc("/public/v1/roles", app.getPublicRolesHandler)
	mux.HandleFunc("/v1/healthcheck", app.healthcheck)
	mux.HandleFunc("/v1/admin/login", app.adminLoginHandler)
	mux.HandleFunc("/v1/admin/logout", app.adminLogoutHandler)

	mux.HandleFunc("/v1/assets", app.getCreateAssetsHandler)

	mux.Handle("/v1/roles", app.deployWebhook(http.HandlerFunc(app.getCreateRolesHandler), http.MethodPost))
	mux.Handle("/v1/roles/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteRolesHandler), http.MethodPut, http.MethodDelete))

	mux.Handle("/v1/companies", app.deployWebhook(http.HandlerFunc(app.getCreateCompaniesHandler), http.MethodPost))
	mux.Handle("/v1/companies/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteCompaniesHandler), http.MethodPut, http.MethodDelete))

	mux.Handle("/v1/notes", app.deployWebhook(http.HandlerFunc(app.getCreateNotesHandler), http.MethodPost))
	mux.Handle("/v1/notes/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteNotesHandler), http.MethodPut, http.MethodDelete))
	mux.Handle("/v1/item-notes", app.deployWebhook(http.HandlerFunc(app.getCreateItemNotesHandler), http.MethodPost))
	mux.Handle("/v1/item-notes/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteItemNotesHandler), http.MethodPut, http.MethodDelete))
	mux.HandleFunc("/v1/item-notes/items/", app.getNotesForItemHandler)

	mux.Handle("/v1/projects", app.deployWebhook(http.HandlerFunc(app.getCreateProjectsHandler), http.MethodPost))
	mux.Handle("/v1/projects/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteProjectsHandler), http.MethodPut, http.MethodDelete))

	mux.Handle("/v1/tagged-items", app.deployWebhook(http.HandlerFunc(app.getCreateTagItemsHandler), http.MethodPost))
	mux.Handle("/v1/tagged-items/", app.deployWebhook(http.HandlerFunc(app.getUpdateDeleteTagItemsHandler), http.MethodPut, http.MethodDelete))
	mux.HandleFunc("/v1/tagged-items/items/", app.getTagsForItemHandler)

	return app.enableCORS(mux)
}
