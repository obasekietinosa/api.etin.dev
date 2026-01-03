package data

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"api.etin.dev/pkg/querybuilder"
)

type ItemNote struct {
	ID       int64    `json:"id"`
	NoteID   int64    `json:"noteId"`
	ItemID   int64    `json:"itemId"`
	ItemType ItemType `json:"itemType"`
}

type ItemNoteModel struct {
	DB    *sql.DB
	Query *querybuilder.QueryBuilder
}

func (i ItemNoteModel) Insert(itemNote *ItemNote) error {
	if err := validateItemType(itemNote.ItemType); err != nil {
		return err
	}

	values := querybuilder.Clauses{
		{ColumnName: "noteId", Value: itemNote.NoteID},
		{ColumnName: "itemId", Value: itemNote.ItemID},
		{ColumnName: "itemType", Value: string(itemNote.ItemType)},
	}

	row, err := i.Query.SetBaseTable("item_notes").Insert(values).Returning("id").QueryRow()
	if err != nil {
		return err
	}

	return row.Scan(&itemNote.ID)
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

func (i ItemNoteModel) RemoveNoteFromItem(noteID, itemID int64, itemType ItemType) error {
	if err := validateItemType(itemType); err != nil {
		return err
	}

	results, err := i.Query.SetBaseTable("item_notes").Delete().
		WhereEqual("noteId", noteID).
		WhereEqual("itemId", itemID).
		WhereEqual("itemType", string(itemType)).
		Exec()
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

func (i ItemNoteModel) GetNotesForItem(itemType ItemType, itemID int64, filters CursorFilters) ([]*Note, Metadata, error) {
	if err := validateItemType(itemType); err != nil {
		return nil, Metadata{}, err
	}

	query := i.Query.SetBaseTable("item_notes").Select(
		"notes.id AS id",
		"notes.createdAt AS createdAt",
		"notes.updatedAt AS updatedAt",
		"notes.deletedAt AS deletedAt",
		"notes.publishedAt AS publishedAt",
		"notes.title AS title",
		"notes.subtitle AS subtitle",
		"notes.slug AS slug",
		"notes.body AS body",
	).LeftJoin("notes", "noteId", "id").
		WhereEqual("item_notes.itemId", itemID).
		WhereEqual("item_notes.itemType", string(itemType)).
		WhereEqual("notes.deletedAt", nil)

	if filters.Cursor != "" {
		decodedCursor, err := base64.URLEncoding.DecodeString(filters.Cursor)
		if err == nil {
			id, _ := strconv.ParseInt(string(decodedCursor), 10, 64)
			if id > 0 {
				query.WhereLessThan("notes.id", id)
			}
		}
	}

	query.OrderBy("notes.id", "DESC")

	if filters.Limit > 0 {
		query.Limit(filters.Limit)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	notes := make([]*Note, 0)

	for rows.Next() {
		note := &Note{}
		var deletedAt sql.NullTime
		var publishedAt sql.NullTime

		if err := rows.Scan(
			&note.ID,
			&note.CreatedAt,
			&note.UpdatedAt,
			&deletedAt,
			&publishedAt,
			&note.Title,
			&note.Subtitle,
			&note.Slug,
			&note.Body,
		); err != nil {
			return nil, Metadata{}, err
		}

		if deletedAt.Valid {
			note.DeletedAt = &deletedAt.Time
		}

		if publishedAt.Valid {
			note.PublishedAt = &publishedAt.Time
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := Metadata{}
	if len(notes) > 0 {
		lastNote := notes[len(notes)-1]
		metadata.NextCursor = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", lastNote.ID)))
	}

	return notes, metadata, nil
}

func (i ItemNoteModel) GetNotesForContentType(itemType ItemType, filters CursorFilters) ([]*Note, Metadata, error) {
	if err := validateItemType(itemType); err != nil {
		return nil, Metadata{}, err
	}

	query := i.Query.SetBaseTable("item_notes").Select(
		"notes.id AS id",
		"notes.createdAt AS createdAt",
		"notes.updatedAt AS updatedAt",
		"notes.deletedAt AS deletedAt",
		"notes.publishedAt AS publishedAt",
		"notes.title AS title",
		"notes.subtitle AS subtitle",
		"notes.slug AS slug",
		"notes.body AS body",
	).LeftJoin("notes", "noteId", "id").
		WhereEqual("item_notes.itemType", string(itemType)).
		WhereEqual("notes.deletedAt", nil)

	if filters.Cursor != "" {
		decodedCursor, err := base64.URLEncoding.DecodeString(filters.Cursor)
		if err == nil {
			id, _ := strconv.ParseInt(string(decodedCursor), 10, 64)
			if id > 0 {
				query.WhereLessThan("notes.id", id)
			}
		}
	}

	query.OrderBy("notes.id", "DESC")

	if filters.Limit > 0 {
		query.Limit(filters.Limit)
	}

	rows, err := query.Query()
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	notes := make([]*Note, 0)

	for rows.Next() {
		note := &Note{}
		var deletedAt sql.NullTime
		var publishedAt sql.NullTime

		if err := rows.Scan(
			&note.ID,
			&note.CreatedAt,
			&note.UpdatedAt,
			&deletedAt,
			&publishedAt,
			&note.Title,
			&note.Subtitle,
			&note.Slug,
			&note.Body,
		); err != nil {
			return nil, Metadata{}, err
		}

		if deletedAt.Valid {
			note.DeletedAt = &deletedAt.Time
		}

		if publishedAt.Valid {
			note.PublishedAt = &publishedAt.Time
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := Metadata{}
	if len(notes) > 0 {
		lastNote := notes[len(notes)-1]
		metadata.NextCursor = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", lastNote.ID)))
	}

	return notes, metadata, nil
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
	var itemType string

	if err := row.Scan(
		&itemNote.ID,
		&itemNote.NoteID,
		&itemNote.ItemID,
		&itemType,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	itemNote.ItemType = ItemType(itemType)

	return &itemNote, nil
}

func (i ItemNoteModel) Update(itemNote *ItemNote) error {
	if err := validateItemType(itemNote.ItemType); err != nil {
		return err
	}

	values := querybuilder.Clauses{
		{ColumnName: "noteId", Value: itemNote.NoteID},
		{ColumnName: "itemId", Value: itemNote.ItemID},
		{ColumnName: "itemType", Value: string(itemNote.ItemType)},
	}

	row, err := i.Query.SetBaseTable("item_notes").Update(values).
		WhereEqual("id", itemNote.ID).
		Returning("id", "noteId", "itemId", "itemType").
		QueryRow()
	if err != nil {
		return err
	}

	var itemType string

	if err := row.Scan(
		&itemNote.ID,
		&itemNote.NoteID,
		&itemNote.ItemID,
		&itemType,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("record not found")
		}
		return err
	}

	itemNote.ItemType = ItemType(itemType)

	return nil
}

func (i ItemNoteModel) GetAll() ([]*ItemNote, error) {
	rows, err := i.Query.SetBaseTable("item_notes").Select(
		"id",
		"noteId",
		"itemId",
		"itemType",
	).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemNotes := make([]*ItemNote, 0)

	for rows.Next() {
		itemNote := &ItemNote{}
		var itemType string

		if err := rows.Scan(
			&itemNote.ID,
			&itemNote.NoteID,
			&itemNote.ItemID,
			&itemType,
		); err != nil {
			return nil, err
		}

		itemNote.ItemType = ItemType(itemType)

		itemNotes = append(itemNotes, itemNote)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return itemNotes, nil
}
