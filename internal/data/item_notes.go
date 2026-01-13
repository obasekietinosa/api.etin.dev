package data

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"api.etin.dev/pkg/querybuilder"
)

type ItemNote struct {
	ID        int64      `json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	NoteID    int64      `json:"noteId"`
	ItemID    int64      `json:"itemId"`
	ItemType  string     `json:"itemType"`
}

type ItemNoteModel struct {
	DB     *sql.DB
	Query  *querybuilder.QueryBuilder
	Logger *log.Logger
}

func (i ItemNoteModel) Insert(itemNote *ItemNote) error {
	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "noteId", Value: itemNote.NoteID},
		querybuilder.Clause{ColumnName: "itemId", Value: itemNote.ItemID},
		querybuilder.Clause{ColumnName: "itemType", Value: itemNote.ItemType},
	}

	row, err := i.Query.SetBaseTable("item_notes").Insert(values).Returning("id", "createdAt", "updatedAt", "deletedAt").QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime

	err = row.Scan(&itemNote.ID, &itemNote.CreatedAt, &itemNote.UpdatedAt, &deletedAt)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		itemNote.DeletedAt = &deletedAt.Time
	}

	return nil
}

func (i ItemNoteModel) Get(id int64) (*ItemNote, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	row, err := i.Query.SetBaseTable("item_notes").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"noteId",
		"itemId",
		"itemType",
	).WhereEqual("deletedAt", nil).WhereEqual("id", id).QueryRow()
	if err != nil {
		return nil, err
	}

	var itemNote ItemNote
	var deletedAt sql.NullTime

	err = row.Scan(
		&itemNote.ID,
		&itemNote.CreatedAt,
		&itemNote.UpdatedAt,
		&deletedAt,
		&itemNote.NoteID,
		&itemNote.ItemID,
		&itemNote.ItemType,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	if deletedAt.Valid {
		itemNote.DeletedAt = &deletedAt.Time
	}

	return &itemNote, nil
}

func (i ItemNoteModel) Update(itemNote *ItemNote) error {
	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "noteId", Value: itemNote.NoteID},
		querybuilder.Clause{ColumnName: "itemId", Value: itemNote.ItemID},
		querybuilder.Clause{ColumnName: "itemType", Value: itemNote.ItemType},
		querybuilder.Clause{ColumnName: "updatedAt", Value: time.Now()},
	}

	row, err := i.Query.With(
		i.Query.SetBaseTable("item_notes").Update(values).WhereEqual("id", itemNote.ID).WhereEqual("deletedAt", nil).Returning("id", "createdAt", "updatedAt", "deletedAt", "noteId", "itemId", "itemType"),
		"updated_item_note",
	).Select(
		"updated_item_note.id",
		"updated_item_note.createdAt",
		"updated_item_note.updatedAt",
		"updated_item_note.deletedAt",
		"updated_item_note.noteId",
		"updated_item_note.itemId",
		"updated_item_note.itemType",
	).From("updated_item_note").QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime

	err = row.Scan(
		&itemNote.ID,
		&itemNote.CreatedAt,
		&itemNote.UpdatedAt,
		&deletedAt,
		&itemNote.NoteID,
		&itemNote.ItemID,
		&itemNote.ItemType,
	)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		itemNote.DeletedAt = &deletedAt.Time
	} else {
		itemNote.DeletedAt = nil
	}

	return nil
}

func (i ItemNoteModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}

	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "updatedAt", Value: time.Now()},
		querybuilder.Clause{ColumnName: "deletedAt", Value: time.Now()},
	}

	results, err := i.Query.SetBaseTable("item_notes").Update(values).WhereEqual("id", id).WhereEqual("deletedAt", nil).Exec()
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

func (i ItemNoteModel) GetAll() ([]*ItemNote, error) {
	rows, err := i.Query.SetBaseTable("item_notes").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"noteId",
		"itemId",
		"itemType",
	).WhereEqual("deletedAt", nil).OrderBy("createdAt", "desc").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemNotes := []*ItemNote{}

	for rows.Next() {
		var itemNote ItemNote
		var deletedAt sql.NullTime

		err := rows.Scan(
			&itemNote.ID,
			&itemNote.CreatedAt,
			&itemNote.UpdatedAt,
			&deletedAt,
			&itemNote.NoteID,
			&itemNote.ItemID,
			&itemNote.ItemType,
		)
		if err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			itemNote.DeletedAt = &deletedAt.Time
		}

		itemNotes = append(itemNotes, &itemNote)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return itemNotes, nil
}

func (i ItemNoteModel) GetNotesForItem(itemType string, itemID int64, filters CursorFilters) ([]*Note, Metadata, error) {
	query := i.Query.SetBaseTable("item_notes").Select(
		"notes.id",
		"notes.createdAt",
		"notes.updatedAt",
		"notes.deletedAt",
		"notes.publishedAt",
		"notes.title",
		"notes.subtitle",
		"notes.slug",
		"notes.body",
	).LeftJoin("notes", "noteId", "id").WhereEqual("itemType", itemType).WhereEqual("itemId", itemID).WhereEqual("deletedAt", nil)

	if filters.Cursor != "" {
		query.WhereLessThan("notes.id", filters.Cursor)
	}

	rows, err := query.Limit(filters.Limit).OrderBy("notes.id", "desc").Query()

	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	notes := []*Note{}

	for rows.Next() {
		var note Note
		var deletedAt sql.NullTime
		var publishedAt sql.NullTime
		var slug sql.NullString

		err := rows.Scan(
			&note.ID,
			&note.CreatedAt,
			&note.UpdatedAt,
			&deletedAt,
			&publishedAt,
			&note.Title,
			&note.Subtitle,
			&slug,
			&note.Body,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		if deletedAt.Valid {
			note.DeletedAt = &deletedAt.Time
		}

		if publishedAt.Valid {
			note.PublishedAt = &publishedAt.Time
		}

		if slug.Valid {
			note.Slug = slug.String
		}

		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(len(notes), filters.Limit, notes)

	return notes, metadata, nil
}

func (i ItemNoteModel) GetNotesForContentType(contentType string, filters CursorFilters) ([]*Note, Metadata, error) {
	query := i.Query.SetBaseTable("item_notes").Select(
		"notes.id",
		"notes.createdAt",
		"notes.updatedAt",
		"notes.deletedAt",
		"notes.publishedAt",
		"notes.title",
		"notes.subtitle",
		"notes.slug",
		"notes.body",
	).LeftJoin("notes", "noteId", "id").WhereEqual("itemType", contentType).WhereEqual("deletedAt", nil)

	if filters.Cursor != "" {
		query.WhereLessThan("notes.id", filters.Cursor)
	}

	rows, err := query.Limit(filters.Limit).OrderBy("notes.id", "desc").Query()

	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	notes := []*Note{}

	for rows.Next() {
		var note Note
		var deletedAt sql.NullTime
		var publishedAt sql.NullTime
		var slug sql.NullString

		err := rows.Scan(
			&note.ID,
			&note.CreatedAt,
			&note.UpdatedAt,
			&deletedAt,
			&publishedAt,
			&note.Title,
			&note.Subtitle,
			&slug,
			&note.Body,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		if deletedAt.Valid {
			note.DeletedAt = &deletedAt.Time
		}

		if publishedAt.Valid {
			note.PublishedAt = &publishedAt.Time
		}

		if slug.Valid {
			note.Slug = slug.String
		}

		notes = append(notes, &note)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(len(notes), filters.Limit, notes)

	return notes, metadata, nil
}
