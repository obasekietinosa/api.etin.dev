package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"api.etin.dev/internal/data"
)

func (app *application) getPublicProjectHandler(w http.ResponseWriter, r *http.Request) {
	idOrSlug := r.PathValue("idOrSlug")

	var project *data.Project
	var err error

	id, errInt := strconv.ParseInt(idOrSlug, 10, 64)
	if errInt == nil {
		project, err = app.getModels(r).Projects.Get(id)
	} else {
		project, err = app.getModels(r).Projects.GetBySlug(idOrSlug)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "record not found" {
			app.writeError(w, http.StatusNotFound)
			return
		}
		app.logger.Printf("Error retrieving project: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	tags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeProjects, project.ID)
	if err != nil {
		app.logger.Printf("Error retrieving tags for project %d: %s", project.ID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	notes, _, err := app.getModels(r).ItemNotes.GetNotesForItem(string(data.ItemTypeProjects), project.ID, data.CursorFilters{OnlyPublished: true})
	if err != nil {
		app.logger.Printf("Error retrieving notes for project %d: %s", project.ID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	relatedItem := &publicRelatedItem{
		ID:    project.ID,
		Title: project.Title,
		Type:  "project",
		Slug:  project.Slug,
	}

	publicNotes := make([]publicNote, 0, len(notes))
	for _, note := range notes {
		noteTags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeNotes, note.ID)
		if err != nil {
			app.logger.Printf("Error retrieving tags for note %d: %s", note.ID, err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}
		publicNotes = append(publicNotes, buildPublicNote(note, noteTags, relatedItem))
	}

	app.writeJSON(w, http.StatusOK, envelope{"project": buildPublicProject(project, tags, publicNotes)})
}

func (app *application) getPublicNoteHandler(w http.ResponseWriter, r *http.Request) {
	idOrSlug := r.PathValue("idOrSlug")

	var note *data.Note
	var err error

	id, errInt := strconv.ParseInt(idOrSlug, 10, 64)
	if errInt == nil {
		note, err = app.getModels(r).Notes.Get(id)
	} else {
		note, err = app.getModels(r).Notes.GetBySlug(idOrSlug)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "record not found" {
			app.writeError(w, http.StatusNotFound)
			return
		}
		app.logger.Printf("Error retrieving note: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	if note.PublishedAt == nil || note.PublishedAt.After(time.Now()) {
		app.writeError(w, http.StatusNotFound)
		return
	}

	tags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeNotes, note.ID)
	if err != nil {
		app.logger.Printf("Error retrieving tags for note %d: %s", note.ID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	var relatedItem *publicRelatedItem
	itemNotes, err := app.getModels(r).ItemNotes.GetByNoteIDs([]int64{note.ID})
	if err != nil {
		app.logger.Printf("Error retrieving item notes: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	if len(itemNotes) > 0 {
		in := itemNotes[0]
		if in.ItemType == string(data.ItemTypeProjects) {
			project, err := app.getModels(r).Projects.Get(in.ItemID)
			if err == nil {
				relatedItem = &publicRelatedItem{
					ID:    project.ID,
					Title: project.Title,
					Type:  "project",
					Slug:  project.Slug,
				}
			}
		} else if in.ItemType == string(data.ItemTypeRoles) {
			role, err := app.getModels(r).Roles.Get(in.ItemID)
			if err == nil {
				relatedItem = &publicRelatedItem{
					ID:    role.ID,
					Title: role.Title,
					Type:  "role",
					Slug:  role.Slug,
				}
			}
		}
	}

	app.writeJSON(w, http.StatusOK, envelope{"note": buildPublicNote(note, tags, relatedItem)})
}

func (app *application) getPublicRoleHandler(w http.ResponseWriter, r *http.Request) {
	idOrSlug := r.PathValue("idOrSlug")

	var role *data.Role
	var err error

	id, errInt := strconv.ParseInt(idOrSlug, 10, 64)
	if errInt == nil {
		role, err = app.getModels(r).Roles.Get(id)
	} else {
		role, err = app.getModels(r).Roles.GetBySlug(idOrSlug)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "record not found" {
			app.writeError(w, http.StatusNotFound)
			return
		}
		app.logger.Printf("Error retrieving role: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	notes, _, err := app.getModels(r).ItemNotes.GetNotesForItem(string(data.ItemTypeRoles), role.ID, data.CursorFilters{OnlyPublished: true})
	if err != nil {
		app.logger.Printf("Error retrieving notes for role %d: %s", role.ID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	relatedItem := &publicRelatedItem{
		ID:    role.ID,
		Title: role.Title,
		Type:  "role",
		Slug:  role.Slug,
	}

	publicNotes := make([]publicNote, 0, len(notes))
	for _, note := range notes {
		noteTags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeNotes, note.ID)
		if err != nil {
			app.logger.Printf("Error retrieving tags for note %d: %s", note.ID, err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}
		publicNotes = append(publicNotes, buildPublicNote(note, noteTags, relatedItem))
	}

	app.writeJSON(w, http.StatusOK, envelope{"role": buildPublicRole(role, publicNotes)})
}
