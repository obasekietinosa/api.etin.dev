package data

import (
	"database/sql"
	"errors"
	"time"

	"api.etin.dev/pkg/querybuilder"
)

type Note struct {
	ID          int64      `json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	Title       string     `json:"title"`
	Subtitle    string     `json:"subtitle"`
	Body        string     `json:"body"`
}

type NoteModel struct {
	DB    *sql.DB
	Query *querybuilder.QueryBuilder
}

func (n NoteModel) Insert(note *Note) error {
	var publishedAt interface{}
	if note.PublishedAt != nil {
		publishedAt = *note.PublishedAt
	}

	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "publishedAt", Value: publishedAt},
		querybuilder.Clause{ColumnName: "title", Value: note.Title},
		querybuilder.Clause{ColumnName: "subtitle", Value: note.Subtitle},
		querybuilder.Clause{ColumnName: "body", Value: note.Body},
	}

	row, err := n.Query.SetBaseTable("notes").Insert(values).Returning("id", "createdAt", "updatedAt", "deletedAt", "publishedAt").QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var published sql.NullTime

	err = row.Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt, &deletedAt, &published)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		note.DeletedAt = &deletedAt.Time
	}

	if published.Valid {
		note.PublishedAt = &published.Time
	} else {
		note.PublishedAt = nil
	}

	return nil
}

func (n NoteModel) Get(id int64) (*Note, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	row, err := n.Query.SetBaseTable("notes").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"publishedAt",
		"title",
		"subtitle",
		"body",
	).WhereEqual("deletedAt", nil).WhereEqual("id", id).QueryRow()
	if err != nil {
		return nil, err
	}

	var note Note
	var deletedAt sql.NullTime
	var publishedAt sql.NullTime

	err = row.Scan(
		&note.ID,
		&note.CreatedAt,
		&note.UpdatedAt,
		&deletedAt,
		&publishedAt,
		&note.Title,
		&note.Subtitle,
		&note.Body,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	if deletedAt.Valid {
		note.DeletedAt = &deletedAt.Time
	}

	if publishedAt.Valid {
		note.PublishedAt = &publishedAt.Time
	}

	return &note, nil
}

func (n NoteModel) Update(note *Note) error {
	var publishedAt interface{}
	if note.PublishedAt != nil {
		publishedAt = *note.PublishedAt
	}

	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "publishedAt", Value: publishedAt},
		querybuilder.Clause{ColumnName: "title", Value: note.Title},
		querybuilder.Clause{ColumnName: "subtitle", Value: note.Subtitle},
		querybuilder.Clause{ColumnName: "body", Value: note.Body},
		querybuilder.Clause{ColumnName: "updatedAt", Value: time.Now()},
	}

	row, err := n.Query.With(
		n.Query.SetBaseTable("notes").Update(values).WhereEqual("id", note.ID).WhereEqual("deletedAt", nil).Returning("id", "createdAt", "updatedAt", "deletedAt", "publishedAt", "title", "subtitle", "body"),
		"updated_note",
	).Select(
		"updated_note.id",
		"updated_note.createdAt",
		"updated_note.updatedAt",
		"updated_note.deletedAt",
		"updated_note.publishedAt",
		"updated_note.title",
		"updated_note.subtitle",
		"updated_note.body",
	).From("updated_note").QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var published sql.NullTime

	err = row.Scan(
		&note.ID,
		&note.CreatedAt,
		&note.UpdatedAt,
		&deletedAt,
		&published,
		&note.Title,
		&note.Subtitle,
		&note.Body,
	)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		note.DeletedAt = &deletedAt.Time
	} else {
		note.DeletedAt = nil
	}

	if published.Valid {
		note.PublishedAt = &published.Time
	} else {
		note.PublishedAt = nil
	}

	return nil
}

func (n NoteModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}

	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "updatedAt", Value: time.Now()},
		querybuilder.Clause{ColumnName: "deletedAt", Value: time.Now()},
	}

	results, err := n.Query.SetBaseTable("notes").Update(values).WhereEqual("id", id).WhereEqual("deletedAt", nil).Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (n NoteModel) GetAll() ([]*Note, error) {
	rows, err := n.Query.SetBaseTable("notes").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"publishedAt",
		"title",
		"subtitle",
		"body",
	).WhereEqual("deletedAt", nil).OrderBy("COALESCE(publishedAt, createdAt)", "desc").Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notes := []*Note{}

	for rows.Next() {
		var note Note
		var deletedAt sql.NullTime
		var publishedAt sql.NullTime

		err := rows.Scan(
			&note.ID,
			&note.CreatedAt,
			&note.UpdatedAt,
			&deletedAt,
			&publishedAt,
			&note.Title,
			&note.Subtitle,
			&note.Body,
		)
		if err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			note.DeletedAt = &deletedAt.Time
		}

		if publishedAt.Valid {
			note.PublishedAt = &publishedAt.Time
		}

		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
