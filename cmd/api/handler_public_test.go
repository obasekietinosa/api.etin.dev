package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"log"
	"os"

	"api.etin.dev/internal/data"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetPublicNotesForContentHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating sqlmock: %s", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "", 0)
	models := data.NewModels(db, logger)

	app := &application{
		logger: logger,
		models: models,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /public/v1/{contentType}/{idOrSlug}/notes", app.getPublicNotesForContentHandler)

	t.Run("Valid Project Slug", func(t *testing.T) {
		slug := "my-project"
		projectID := int64(10)

		// Expect GetBySlug for Project
		// The query builder generates SELECT ... FROM projects WHERE deletedAt IS NULL AND slug = $1
		mock.ExpectQuery(`SELECT id, createdAt, updatedAt, deletedAt, startDate, endDate, title, slug, description, imageUrl FROM projects WHERE deletedAt IS NULL AND slug = \$1`).
			WithArgs(slug).
			WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "startDate", "endDate", "title", "slug", "description", "imageUrl"}).
				AddRow(projectID, time.Now(), time.Now(), nil, time.Now(), nil, "Title", slug, "Desc", "img.jpg"))

		// Expect GetNotesForItem
		mock.ExpectQuery(`SELECT notes.id, .* FROM item_notes .*`).
			WithArgs("projects", projectID, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "slug", "body"}))

		req := httptest.NewRequest(http.MethodGet, "/public/v1/projects/"+slug+"/notes", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Valid Project ID", func(t *testing.T) {
		projectID := int64(20)

		// Expect Get for Project (Verification step)
		mock.ExpectQuery(`SELECT id, createdAt, updatedAt, deletedAt, startDate, endDate, title, slug, description, imageUrl FROM projects WHERE deletedAt IS NULL AND id = \$1`).
			WithArgs(projectID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "startDate", "endDate", "title", "slug", "description", "imageUrl"}).
				AddRow(projectID, time.Now(), time.Now(), nil, time.Now(), nil, "Title", "slug", "Desc", "img.jpg"))

		// Expect GetNotesForItem
		mock.ExpectQuery(`SELECT notes.id, .* FROM item_notes .*`).
			WithArgs("projects", projectID, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "slug", "body"}))

		req := httptest.NewRequest(http.MethodGet, "/public/v1/projects/20/notes", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("Invalid Slug", func(t *testing.T) {
		slug := "invalid-slug"

		// Expect GetBySlug for Project to fail
		mock.ExpectQuery(`SELECT id, createdAt, updatedAt, deletedAt, startDate, endDate, title, slug, description, imageUrl FROM projects WHERE deletedAt IS NULL AND slug = \$1`).
			WithArgs(slug).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, "/public/v1/projects/"+slug+"/notes", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rr.Code)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		id := int64(999)

		// Expect Get for Project to fail
		mock.ExpectQuery(`SELECT id, createdAt, updatedAt, deletedAt, startDate, endDate, title, slug, description, imageUrl FROM projects WHERE deletedAt IS NULL AND id = \$1`).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, "/public/v1/projects/999/notes", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rr.Code)
		}
	})
}
