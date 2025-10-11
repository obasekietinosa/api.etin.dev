package data

import (
	"database/sql"
	"errors"

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

func (i ItemNoteModel) GetNotesForItem(itemType ItemType, itemID int64) ([]*Note, error) {
	if err := validateItemType(itemType); err != nil {
		return nil, err
	}

	rows, err := i.Query.SetBaseTable("item_notes").Select(
		"notes.id AS id",
		"notes.createdAt AS createdAt",
		"notes.updatedAt AS updatedAt",
		"notes.deletedAt AS deletedAt",
		"notes.publishedAt AS publishedAt",
		"notes.title AS title",
		"notes.subtitle AS subtitle",
		"notes.body AS body",
	).LeftJoin("notes", "noteId", "id").
		WhereEqual("item_notes.itemId", itemID).
		WhereEqual("item_notes.itemType", string(itemType)).
		WhereEqual("notes.deletedAt", nil).
		Query()
	if err != nil {
		return nil, err
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
			&note.Body,
		); err != nil {
			return nil, err
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
		return nil, err
	}

	return notes, nil
}
