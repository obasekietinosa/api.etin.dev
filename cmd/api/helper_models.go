package main

import (
	"fmt"
	"log"
	"net/http"

	"api.etin.dev/internal/data"
)

func (app *application) getModels(r *http.Request) data.Models {
	id, ok := r.Context().Value(requestIdKey).(string)
	if !ok {
		id = "unknown"
	}

	newLogger := log.New(app.logger.Writer(), fmt.Sprintf("[%s] ", id), app.logger.Flags())

	models := app.models
	models.Roles.Logger = newLogger
	models.Companies.Logger = newLogger
	models.Notes.Logger = newLogger
	models.Projects.Logger = newLogger
	models.Tags.Logger = newLogger
	models.TagItems.Logger = newLogger
	models.ItemNotes.Logger = newLogger
	models.Assets.Logger = newLogger

	return models
}
