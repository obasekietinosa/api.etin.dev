package data

import (
	"log"
	"os"
	"testing"

	"api.etin.dev/pkg/querybuilder"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestItemNoteModel_GetNotesForItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating sqlmock: %s", err)
	}
	defer db.Close()

	qb := &querybuilder.QueryBuilder{DB: db}
	m := ItemNoteModel{
		DB:     db,
		Query:  qb,
		Logger: log.New(os.Stdout, "", 0),
	}

	// Case 1: No published filtering (default)
	mock.ExpectQuery(`SELECT notes.id, notes.createdAt, notes.updatedAt, notes.deletedAt, notes.publishedAt, notes.title, notes.subtitle, notes.slug, notes.body FROM item_notes LEFT JOIN notes ON item_notes.noteId = notes.id WHERE itemType = \$1 AND itemId = \$2 AND notes.deletedAt IS NULL ORDER BY notes.id desc LIMIT 20`).
		WithArgs("projects", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "slug", "body"}))

	filters := CursorFilters{Limit: 20}
	_, _, err = m.GetNotesForItem("projects", 1, filters)
	if err != nil {
		t.Fatalf("unexpected error calling GetNotesForItem: %s", err)
	}

	// Case 2: With published filtering
	mock.ExpectQuery(`SELECT notes.id, notes.createdAt, notes.updatedAt, notes.deletedAt, notes.publishedAt, notes.title, notes.subtitle, notes.slug, notes.body FROM item_notes LEFT JOIN notes ON item_notes.noteId = notes.id WHERE itemType = \$1 AND itemId = \$2 AND notes.deletedAt IS NULL AND notes.publishedAt <= \$3 ORDER BY notes.id desc LIMIT 20`).
		WithArgs("projects", 1, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "slug", "body"}))

	filtersPublished := CursorFilters{Limit: 20, OnlyPublished: true}
	_, _, err = m.GetNotesForItem("projects", 1, filtersPublished)
	if err != nil {
		t.Fatalf("unexpected error calling GetNotesForItem with published filter: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unmet expectations: %s", err)
	}
}

func TestItemNoteModel_GetNotesForContentType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating sqlmock: %s", err)
	}
	defer db.Close()

	qb := &querybuilder.QueryBuilder{DB: db}
	m := ItemNoteModel{
		DB:     db,
		Query:  qb,
		Logger: log.New(os.Stdout, "", 0),
	}

	// Case 1: No published filtering (default)
	mock.ExpectQuery(`SELECT notes.id, notes.createdAt, notes.updatedAt, notes.deletedAt, notes.publishedAt, notes.title, notes.subtitle, notes.slug, notes.body FROM item_notes LEFT JOIN notes ON item_notes.noteId = notes.id WHERE itemType = \$1 AND notes.deletedAt IS NULL ORDER BY notes.id desc LIMIT 20`).
		WithArgs("projects").
		WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "slug", "body"}))

	filters := CursorFilters{Limit: 20}
	_, _, err = m.GetNotesForContentType("projects", filters)
	if err != nil {
		t.Fatalf("unexpected error calling GetNotesForContentType: %s", err)
	}

	// Case 2: With published filtering
	mock.ExpectQuery(`SELECT notes.id, notes.createdAt, notes.updatedAt, notes.deletedAt, notes.publishedAt, notes.title, notes.subtitle, notes.slug, notes.body FROM item_notes LEFT JOIN notes ON item_notes.noteId = notes.id WHERE itemType = \$1 AND notes.deletedAt IS NULL AND notes.publishedAt <= \$2 ORDER BY notes.id desc LIMIT 20`).
		WithArgs("projects", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "slug", "body"}))

	filtersPublished := CursorFilters{Limit: 20, OnlyPublished: true}
	_, _, err = m.GetNotesForContentType("projects", filtersPublished)
	if err != nil {
		t.Fatalf("unexpected error calling GetNotesForContentType with published filter: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unmet expectations: %s", err)
	}
}
