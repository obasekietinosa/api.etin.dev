package data

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"api.etin.dev/pkg/querybuilder"
)

type ItemNote struct {
	ID       int64  `json:"id"`
	NoteID   int64  `json:"noteId"`
	ItemID   int64  `json:"itemId"`
	ItemType string `json:"itemType"`
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

	row, err := i.Query.SetBaseTable("item_notes").Insert(values).Returning("id").QueryRow()
	if err != nil {
		return err
	}

	err = row.Scan(&itemNote.ID)
	if err != nil {
		return err
	}

	return nil
}

func (i ItemNoteModel) Get(id int64) (*ItemNote, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	row, err := i.Query.SetBaseTable("item_notes").Select(
		"id",
		"noteId",
		"itemId",
		"itemType",
	).WhereEqual("id", id).QueryRow()
	if err != nil {
		return nil, err
	}

	var itemNote ItemNote

	err = row.Scan(
		&itemNote.ID,
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

	return &itemNote, nil
}

func (i ItemNoteModel) Update(itemNote *ItemNote) error {
	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "noteId", Value: itemNote.NoteID},
		querybuilder.Clause{ColumnName: "itemId", Value: itemNote.ItemID},
		querybuilder.Clause{ColumnName: "itemType", Value: itemNote.ItemType},
	}

	row, err := i.Query.With(
		i.Query.SetBaseTable("item_notes").Update(values).WhereEqual("id", itemNote.ID).Returning("id", "noteId", "itemId", "itemType"),
		"updated_item_note",
	).Select(
		"updated_item_note.id",
		"updated_item_note.noteId",
		"updated_item_note.itemId",
		"updated_item_note.itemType",
	).From("updated_item_note").QueryRow()
	if err != nil {
		return err
	}

	err = row.Scan(
		&itemNote.ID,
		&itemNote.NoteID,
		&itemNote.ItemID,
		&itemNote.ItemType,
	)
	if err != nil {
		return err
	}

	return nil
}

func (i ItemNoteModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}

	results, err := i.Query.SetBaseTable("item_notes").Delete().WhereEqual("id", id).Exec()
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
		"noteId",
		"itemId",
		"itemType",
	).OrderBy("id", "desc").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemNotes := []*ItemNote{}

	for rows.Next() {
		var itemNote ItemNote

		err := rows.Scan(
			&itemNote.ID,
			&itemNote.NoteID,
			&itemNote.ItemID,
			&itemNote.ItemType,
		)
		if err != nil {
			return nil, err
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
	).LeftJoin("notes", "noteId", "id").WhereEqual("itemType", itemType).WhereEqual("itemId", itemID).WhereEqual("notes.deletedAt", nil)

	if filters.Cursor != "" {
		query.WhereLessThan("notes.id", filters.Cursor)
	}

	if filters.OnlyPublished {
		query.WhereLessThanEqual("notes.publishedAt", time.Now())
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
	).LeftJoin("notes", "noteId", "id").WhereEqual("itemType", contentType).WhereEqual("notes.deletedAt", nil)

	if filters.Cursor != "" {
		query.WhereLessThan("notes.id", filters.Cursor)
	}

	if filters.OnlyPublished {
		query.WhereLessThanEqual("notes.publishedAt", time.Now())
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
