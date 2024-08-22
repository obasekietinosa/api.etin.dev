package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthcheck", app.healthcheck)

	mux.HandleFunc("/v1/roles", app.getCreateRolesHandler)
	mux.HandleFunc("/v1/roles/", app.getUpdateDeleteRolesHandler)

	mux.HandleFunc("/v1/companies", app.getCreateCompaniesHandler)
	mux.HandleFunc("/v1/companies/", app.getUpdateDeleteCompaniesHandler)

	return mux
}
