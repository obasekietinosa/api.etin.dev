package main

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"api.etin.dev/internal/data"
)

type publicTag struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Slug  string  `json:"slug"`
	Icon  *string `json:"icon"`
	Theme *string `json:"theme"`
}

type publicNote struct {
	ID          int64       `json:"id"`
	PublishedAt string      `json:"publishedAt"`
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Preview     string      `json:"preview"`
	Body        string      `json:"body"`
	IsFeatured  bool        `json:"isFeatured"`
	Tags        []publicTag `json:"tags"`
}

type publicProject struct {
	ID           int64        `json:"id"`
	StartDate    string       `json:"startDate"`
	EndDate      *string      `json:"endDate"`
	Title        string       `json:"title"`
	Image        string       `json:"image"`
	Slug         string       `json:"slug"`
	Status       *publicTag   `json:"status"`
	Description  string       `json:"description"`
	Technologies []string     `json:"technologies"`
	Notes        []publicNote `json:"notes"`
}

type publicRole struct {
	RoleID      int64        `json:"roleId"`
	StartDate   string       `json:"startDate"`
	EndDate     *string      `json:"endDate"`
	Title       string       `json:"title"`
	Subtitle    *string      `json:"subtitle"`
	Company     string       `json:"company"`
	CompanyIcon string       `json:"companyIcon"`
	Slug        string       `json:"slug"`
	Description string       `json:"description"`
	Skills      []string     `json:"skills"`
	Notes       []publicNote `json:"notes"`
}

var slugPattern = regexp.MustCompile(`[^a-z0-9]+`)

func (app *application) getPublicNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes, err := app.getModels(r).Notes.GetAllPublished()
	if err != nil {
		app.logger.Printf("Error retrieving notes: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	response := make([]publicNote, 0, len(notes))

	for _, note := range notes {
		tags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeNotes, note.ID)
		if err != nil {
			app.logger.Printf("Error retrieving tags for note %d: %s", note.ID, err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}

		response = append(response, buildPublicNote(note, tags))
	}

	app.writeJSON(w, http.StatusOK, envelope{"notes": response})
}

func (app *application) getPublicProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := app.getModels(r).Projects.GetAll()
	if err != nil {
		app.logger.Printf("Error retrieving projects: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	response := make([]publicProject, 0, len(projects))

	for _, project := range projects {
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

		publicNotes := make([]publicNote, 0, len(notes))
		for _, note := range notes {
			noteTags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeNotes, note.ID)
			if err != nil {
				app.logger.Printf("Error retrieving tags for note %d: %s", note.ID, err)
				app.writeError(w, http.StatusInternalServerError)
				return
			}
			publicNotes = append(publicNotes, buildPublicNote(note, noteTags))
		}

		response = append(response, buildPublicProject(project, tags, publicNotes))
	}

	app.writeJSON(w, http.StatusOK, envelope{"projects": response})
}

func (app *application) getPublicRolesHandler(w http.ResponseWriter, r *http.Request) {
	roles, err := app.getModels(r).Roles.GetAll()
	if err != nil {
		app.logger.Printf("Error retrieving roles: %s", err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	response := make([]publicRole, 0, len(roles))

	for _, role := range roles {
		notes, _, err := app.getModels(r).ItemNotes.GetNotesForItem(string(data.ItemTypeRoles), role.ID, data.CursorFilters{OnlyPublished: true})
		if err != nil {
			app.logger.Printf("Error retrieving notes for role %d: %s", role.ID, err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}

		publicNotes := make([]publicNote, 0, len(notes))
		for _, note := range notes {
			noteTags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeNotes, note.ID)
			if err != nil {
				app.logger.Printf("Error retrieving tags for note %d: %s", note.ID, err)
				app.writeError(w, http.StatusInternalServerError)
				return
			}
			publicNotes = append(publicNotes, buildPublicNote(note, noteTags))
		}

		response = append(response, buildPublicRole(role, publicNotes))
	}

	app.writeJSON(w, http.StatusOK, envelope{"roles": response})
}

func (app *application) getPublicNotesForContentHandler(w http.ResponseWriter, r *http.Request) {
	// Replaced manual path parsing with r.PathValue
	// Pattern: GET /public/v1/{contentType}/{id}/notes

	contentTypeStr := r.PathValue("contentType")
	itemIDStr := r.PathValue("id")

	if contentTypeStr == "" || itemIDStr == "" {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemType := data.ItemType(strings.ToLower(contentTypeStr))

	filters := data.CursorFilters{
		Limit:         20,
		OnlyPublished: true,
	}

	qs := r.URL.Query()
	if cursor := qs.Get("cursor"); cursor != "" {
		filters.Cursor = cursor
	}
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err == nil && l > 0 {
			filters.Limit = l
		}
	}

	notes, metadata, err := app.getModels(r).ItemNotes.GetNotesForItem(string(itemType), itemID, filters)
	if err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}
		app.logger.Printf("Error retrieving notes for %s item %d: %s", itemType, itemID, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	publicNotes := make([]publicNote, 0, len(notes))
	for _, note := range notes {
		noteTags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeNotes, note.ID)
		if err != nil {
			app.logger.Printf("Error retrieving tags for note %d: %s", note.ID, err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}
		publicNotes = append(publicNotes, buildPublicNote(note, noteTags))
	}

	app.writeJSON(w, http.StatusOK, envelope{"notes": publicNotes, "metadata": metadata})
}

func (app *application) getPublicAllNotesForContentHandler(w http.ResponseWriter, r *http.Request) {
	// Replaced manual path parsing with r.PathValue
	// Pattern: GET /public/v1/{contentType}/notes

	contentTypeStr := r.PathValue("contentType")
	if contentTypeStr == "" {
		// Fallback for explicit routes (projects/roles/notes) where {contentType} wildcard is missing
		// Expected path format: /public/v1/{contentType}/notes
		segments := strings.Split(r.URL.Path, "/")
		if len(segments) >= 4 {
			contentTypeStr = segments[3] // /public/v1/projects/notes -> [ "", "public", "v1", "projects", "notes" ]
		}
	}

	if contentTypeStr == "" {
		app.writeError(w, http.StatusBadRequest)
		return
	}

	itemType := data.ItemType(strings.ToLower(contentTypeStr))

	filters := data.CursorFilters{
		Limit:         20,
		OnlyPublished: true,
	}

	qs := r.URL.Query()
	if cursor := qs.Get("cursor"); cursor != "" {
		filters.Cursor = cursor
	}
	if limit := qs.Get("limit"); limit != "" {
		l, err := strconv.Atoi(limit)
		if err == nil && l > 0 {
			filters.Limit = l
		}
	}

	notes, metadata, err := app.getModels(r).ItemNotes.GetNotesForContentType(string(itemType), filters)
	if err != nil {
		if errors.Is(err, data.ErrInvalidItemType) {
			app.writeError(w, http.StatusBadRequest)
			return
		}
		app.logger.Printf("Error retrieving notes for content type %s: %s", itemType, err)
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	publicNotes := make([]publicNote, 0, len(notes))
	for _, note := range notes {
		noteTags, err := app.getModels(r).TagItems.GetTagsForItem(data.ItemTypeNotes, note.ID)
		if err != nil {
			app.logger.Printf("Error retrieving tags for note %d: %s", note.ID, err)
			app.writeError(w, http.StatusInternalServerError)
			return
		}
		publicNotes = append(publicNotes, buildPublicNote(note, noteTags))
	}

	app.writeJSON(w, http.StatusOK, envelope{"notes": publicNotes, "metadata": metadata})
}

func buildPublicNote(note *data.Note, tags []*data.Tag) publicNote {
	publishedAt := ""
	if note.PublishedAt != nil {
		publishedAt = formatTime(*note.PublishedAt)
	}

	publicTags := convertTags(tags)

	return publicNote{
		ID:          note.ID,
		PublishedAt: publishedAt,
		Title:       note.Title,
		Slug:        note.Slug,
		Preview:     buildPreview(note.Subtitle, note.Body),
		Body:        note.Body,
		IsFeatured:  hasFeaturedTag(publicTags),
		Tags:        publicTags,
	}
}

func buildPublicProject(project *data.Project, tags []*data.Tag, notes []publicNote) publicProject {
	startDate := formatTime(project.StartDate)

	var endDate *string
	if project.EndDate != nil {
		formatted := formatTime(*project.EndDate)
		endDate = &formatted
	}

	publicTags := convertTags(tags)

	var status *publicTag
	technologies := make([]string, 0)

	for i := range publicTags {
		tag := &publicTags[i]

		if tag.Theme != nil && strings.EqualFold(*tag.Theme, "status") {
			status = tag
			continue
		}

		if tag.Theme != nil && strings.EqualFold(*tag.Theme, "technology") {
			technologies = append(technologies, tag.Name)
		}
	}

	image := ""
	if project.ImageURL != nil {
		image = *project.ImageURL
	}

	return publicProject{
		ID:           project.ID,
		StartDate:    startDate,
		EndDate:      endDate,
		Title:        project.Title,
		Image:        image,
		Slug:         project.Slug,
		Status:       status,
		Description:  project.Description,
		Technologies: technologies,
		Notes:        notes,
	}
}

func buildPublicRole(role *data.Role, notes []publicNote) publicRole {
	startDate := formatTime(role.StartDate)

	var endDate *string
	if !role.EndDate.IsZero() {
		formatted := formatTime(role.EndDate)
		endDate = &formatted
	}

	var subtitle *string
	if trimmed := strings.TrimSpace(role.Subtitle); trimmed != "" {
		subtitle = &trimmed
	}

	return publicRole{
		RoleID:      role.ID,
		StartDate:   startDate,
		EndDate:     endDate,
		Title:       role.Title,
		Subtitle:    subtitle,
		Company:     role.Company,
		CompanyIcon: role.CompanyIcon,
		Slug:        role.Slug,
		Description: role.Description,
		Skills:      append([]string{}, role.Skills...),
		Notes:       notes,
	}
}

func convertTags(tags []*data.Tag) []publicTag {
	publicTags := make([]publicTag, 0, len(tags))

	for _, tag := range tags {
		publicTag := publicTag{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		}

		if tag.Icon != nil {
			icon := *tag.Icon
			publicTag.Icon = &icon
		}

		if tag.Theme != nil {
			theme := *tag.Theme
			publicTag.Theme = &theme
		}

		publicTags = append(publicTags, publicTag)
	}

	return publicTags
}

func hasFeaturedTag(tags []publicTag) bool {
	for _, tag := range tags {
		if strings.EqualFold(tag.Slug, "featured") {
			return true
		}
	}
	return false
}

func buildPreview(subtitle, body string) string {
	trimmedSubtitle := strings.TrimSpace(subtitle)
	if trimmedSubtitle != "" {
		return trimmedSubtitle
	}

	trimmedBody := strings.TrimSpace(body)
	if trimmedBody == "" {
		return ""
	}

	const maxPreviewRunes = 200
	runes := []rune(trimmedBody)
	if len(runes) <= maxPreviewRunes {
		return trimmedBody
	}

	trimmed := strings.TrimSpace(string(runes[:maxPreviewRunes]))
	if trimmed == "" {
		return ""
	}

	return trimmed + "…"
}

func slugify(input string) string {
	lowered := strings.ToLower(strings.TrimSpace(input))
	if lowered == "" {
		return ""
	}

	slug := slugPattern.ReplaceAllString(lowered, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
