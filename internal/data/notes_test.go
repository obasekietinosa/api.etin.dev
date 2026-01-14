package data

import (
	"log"
	"os"
	"testing"

	"api.etin.dev/pkg/querybuilder"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNoteModel_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating sqlmock: %s", err)
	}
	defer db.Close()

	qb := &querybuilder.QueryBuilder{DB: db}
	m := NoteModel{
		DB:     db,
		Query:  qb,
		Logger: log.New(os.Stdout, "", 0),
	}

	// Expectation for GetAll (no filtering by publishedAt)
	mock.ExpectQuery(`SELECT id, createdAt, updatedAt, deletedAt, publishedAt, title, subtitle, slug, body FROM notes WHERE deletedAt IS NULL ORDER BY COALESCE\(publishedAt, createdAt\) desc`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "slug", "body"}))

	_, err = m.GetAll()
	if err != nil {
		t.Fatalf("unexpected error calling GetAll: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unmet expectations: %s", err)
	}
}

func TestNoteModel_GetAllPublished(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating sqlmock: %s", err)
	}
	defer db.Close()

	qb := &querybuilder.QueryBuilder{DB: db}
	m := NoteModel{
		DB:     db,
		Query:  qb,
		Logger: log.New(os.Stdout, "", 0),
	}

	// Expectation for GetAllPublished (filtering by publishedAt <= NOW)
	mock.ExpectQuery(`SELECT id, createdAt, updatedAt, deletedAt, publishedAt, title, subtitle, slug, body FROM notes WHERE deletedAt IS NULL AND publishedAt <= \$1 ORDER BY COALESCE\(publishedAt, createdAt\) desc`).
		WithArgs(sqlmock.AnyArg()). // Time argument
		WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "slug", "body"}))

	_, err = m.GetAllPublished()
	if err != nil {
		t.Fatalf("unexpected error calling GetAllPublished: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unmet expectations: %s", err)
	}
}
