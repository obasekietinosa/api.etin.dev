package data

import (
	"database/sql"
	"errors"

	"api.etin.dev/pkg/querybuilder"
)

var ErrInvalidItemType = errors.New("invalid item type")

type ItemType string

const (
	ItemTypeNotes    ItemType = "notes"
	ItemTypeRoles    ItemType = "roles"
	ItemTypeProjects ItemType = "projects"
)

type TagItem struct {
	ID       int64    `json:"id"`
	TagID    int64    `json:"tagId"`
	ItemID   int64    `json:"itemId"`
	ItemType ItemType `json:"itemType"`
}

type TagItemModel struct {
	DB    *sql.DB
	Query *querybuilder.QueryBuilder
}

func (t TagItemModel) Insert(tagItem *TagItem) error {
	if err := validateItemType(tagItem.ItemType); err != nil {
		return err
	}

	values := querybuilder.Clauses{
		{ColumnName: "tagId", Value: tagItem.TagID},
		{ColumnName: "itemId", Value: tagItem.ItemID},
		{ColumnName: "itemType", Value: string(tagItem.ItemType)},
	}

	row, err := t.Query.SetBaseTable("tagged_items").Insert(values).Returning("id").QueryRow()
	if err != nil {
		return err
	}

	return row.Scan(&tagItem.ID)
}

func (t TagItemModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}

	results, err := t.Query.SetBaseTable("tagged_items").Delete().WhereEqual("id", id).Exec()
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

func (t TagItemModel) RemoveTagFromItem(tagID, itemID int64, itemType ItemType) error {
	if err := validateItemType(itemType); err != nil {
		return err
	}

	results, err := t.Query.SetBaseTable("tagged_items").Delete().
		WhereEqual("tagId", tagID).
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

func (t TagItemModel) GetTagsForItem(itemType ItemType, itemID int64) ([]*Tag, error) {
	if err := validateItemType(itemType); err != nil {
		return nil, err
	}

	rows, err := t.Query.SetBaseTable("tagged_items").Select(
		"tags.id AS id",
		"tags.createdAt AS createdAt",
		"tags.updatedAt AS updatedAt",
		"tags.deletedAt AS deletedAt",
		"tags.name AS name",
		"tags.slug AS slug",
		"tags.icon AS icon",
		"tags.theme AS theme",
	).LeftJoin("tags", "tagId", "id").
		WhereEqual("tagged_items.itemId", itemID).
		WhereEqual("tagged_items.itemType", string(itemType)).
		WhereEqual("tags.deletedAt", nil).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*Tag, 0)

	for rows.Next() {
		tag := &Tag{}
		var deletedAt sql.NullTime
		var iconValue sql.NullString
		var themeValue sql.NullString

		if err := rows.Scan(
			&tag.ID,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&deletedAt,
			&tag.Name,
			&tag.Slug,
			&iconValue,
			&themeValue,
		); err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			tag.DeletedAt = &deletedAt.Time
		}

		if iconValue.Valid {
			value := iconValue.String
			tag.Icon = &value
		}

		if themeValue.Valid {
			value := themeValue.String
			tag.Theme = &value
		}

		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func validateItemType(itemType ItemType) error {
	switch itemType {
	case ItemTypeNotes, ItemTypeRoles, ItemTypeProjects:
		return nil
	default:
		return ErrInvalidItemType
	}
}
