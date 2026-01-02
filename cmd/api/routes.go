package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger", app.swaggerHandler)
	mux.HandleFunc("GET /public/v1/notes", app.getPublicNotesHandler)
	mux.HandleFunc("GET /public/v1/{contentType}/{id}/notes", app.getPublicNotesForContentHandler)
	mux.HandleFunc("GET /public/v1/{contentType}/notes", app.getPublicAllNotesForContentHandler)
	mux.HandleFunc("GET /public/v1/projects", app.getPublicProjectsHandler)
	mux.HandleFunc("GET /public/v1/roles", app.getPublicRolesHandler)
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheck)
	mux.HandleFunc("POST /v1/admin/login", app.adminLoginHandler)
	mux.HandleFunc("POST /v1/admin/logout", app.adminLogoutHandler)

	mux.HandleFunc("POST /v1/assets", app.getCreateAssetsHandler)

	mux.Handle("GET /v1/roles", app.deployWebhook(http.HandlerFunc(app.getRolesHandler)))
	mux.Handle("POST /v1/roles", app.deployWebhook(http.HandlerFunc(app.createRoleHandler)))
	mux.Handle("GET /v1/roles/{id}", app.deployWebhook(http.HandlerFunc(app.getRoleHandler)))
	mux.Handle("PUT /v1/roles/{id}", app.deployWebhook(http.HandlerFunc(app.updateRoleHandler)))
	mux.Handle("DELETE /v1/roles/{id}", app.deployWebhook(http.HandlerFunc(app.deleteRoleHandler)))

	mux.Handle("GET /v1/companies", app.deployWebhook(http.HandlerFunc(app.getCompaniesHandler)))
	mux.Handle("POST /v1/companies", app.deployWebhook(http.HandlerFunc(app.createCompanyHandler)))
	mux.Handle("GET /v1/companies/{id}", app.deployWebhook(http.HandlerFunc(app.getCompanyHandler)))
	mux.Handle("PUT /v1/companies/{id}", app.deployWebhook(http.HandlerFunc(app.updateCompanyHandler)))
	mux.Handle("DELETE /v1/companies/{id}", app.deployWebhook(http.HandlerFunc(app.deleteCompanyHandler)))

	mux.Handle("GET /v1/notes", app.deployWebhook(http.HandlerFunc(app.getNotesHandler)))
	mux.Handle("POST /v1/notes", app.deployWebhook(http.HandlerFunc(app.createNoteHandler)))
	mux.Handle("GET /v1/notes/{id}", app.deployWebhook(http.HandlerFunc(app.getNoteHandler)))
	mux.Handle("PUT /v1/notes/{id}", app.deployWebhook(http.HandlerFunc(app.updateNoteHandler)))
	mux.Handle("DELETE /v1/notes/{id}", app.deployWebhook(http.HandlerFunc(app.deleteNoteHandler)))

	mux.Handle("GET /v1/item-notes", app.deployWebhook(http.HandlerFunc(app.getItemNotesHandler)))
	mux.Handle("POST /v1/item-notes", app.deployWebhook(http.HandlerFunc(app.createItemNoteHandler)))
	mux.Handle("GET /v1/item-notes/{id}", app.deployWebhook(http.HandlerFunc(app.getItemNoteHandler)))
	mux.Handle("PUT /v1/item-notes/{id}", app.deployWebhook(http.HandlerFunc(app.updateItemNoteHandler)))
	mux.Handle("DELETE /v1/item-notes/{id}", app.deployWebhook(http.HandlerFunc(app.deleteItemNoteHandler)))
	mux.HandleFunc("GET /v1/item-notes/items/{itemType}/{itemId}", app.getNotesForItemHandler)

	mux.HandleFunc("POST /v1/{contentType}/{id}/notes", app.getCreateContentNoteHandler)
	mux.HandleFunc("GET /v1/{contentType}/{id}/notes", app.getContentNotesHandler)
	mux.HandleFunc("GET /v1/{contentType}/notes", app.getAllContentNotesHandler)

	mux.Handle("GET /v1/projects", app.deployWebhook(http.HandlerFunc(app.getProjectsHandler)))
	mux.Handle("POST /v1/projects", app.deployWebhook(http.HandlerFunc(app.createProjectHandler)))
	mux.Handle("GET /v1/projects/{id}", app.deployWebhook(http.HandlerFunc(app.getProjectHandler)))
	mux.Handle("PUT /v1/projects/{id}", app.deployWebhook(http.HandlerFunc(app.updateProjectHandler)))
	mux.Handle("DELETE /v1/projects/{id}", app.deployWebhook(http.HandlerFunc(app.deleteProjectHandler)))

	mux.Handle("GET /v1/tagged-items", app.deployWebhook(http.HandlerFunc(app.getTagItemsHandler)))
	mux.Handle("POST /v1/tagged-items", app.deployWebhook(http.HandlerFunc(app.createTagItemHandler)))
	mux.Handle("GET /v1/tagged-items/{id}", app.deployWebhook(http.HandlerFunc(app.getTagItemHandler)))
	mux.Handle("PUT /v1/tagged-items/{id}", app.deployWebhook(http.HandlerFunc(app.updateTagItemHandler)))
	mux.Handle("DELETE /v1/tagged-items/{id}", app.deployWebhook(http.HandlerFunc(app.deleteTagItemHandler)))
	mux.HandleFunc("GET /v1/tagged-items/items/{itemType}/{itemId}", app.getTagsForItemHandler)

	return app.enableCORS(mux)
}
