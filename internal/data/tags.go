package data

import (
	"database/sql"
	"errors"
	"time"

	"api.etin.dev/pkg/querybuilder"
)

type Tag struct {
	ID        int64      `json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Icon      *string    `json:"icon,omitempty"`
	Theme     *string    `json:"theme,omitempty"`
}

type TagModel struct {
	DB    *sql.DB
	Query *querybuilder.QueryBuilder
}

func (t TagModel) Insert(tag *Tag) error {
	var icon interface{}
	if tag.Icon != nil {
		icon = *tag.Icon
	}

	var theme interface{}
	if tag.Theme != nil {
		theme = *tag.Theme
	}

	values := querybuilder.Clauses{
		{ColumnName: "name", Value: tag.Name},
		{ColumnName: "slug", Value: tag.Slug},
		{ColumnName: "icon", Value: icon},
		{ColumnName: "theme", Value: theme},
	}

	row, err := t.Query.SetBaseTable("tags").Insert(values).Returning(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"icon",
		"theme",
	).QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var iconValue sql.NullString
	var themeValue sql.NullString

	if err := row.Scan(
		&tag.ID,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&deletedAt,
		&iconValue,
		&themeValue,
	); err != nil {
		return err
	}

	if deletedAt.Valid {
		tag.DeletedAt = &deletedAt.Time
	} else {
		tag.DeletedAt = nil
	}

	if iconValue.Valid {
		value := iconValue.String
		tag.Icon = &value
	} else {
		tag.Icon = nil
	}

	if themeValue.Valid {
		value := themeValue.String
		tag.Theme = &value
	} else {
		tag.Theme = nil
	}

	return nil
}

func (t TagModel) Get(id int64) (*Tag, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	row, err := t.Query.SetBaseTable("tags").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"name",
		"slug",
		"icon",
		"theme",
	).WhereEqual("deletedAt", nil).WhereEqual("id", id).QueryRow()
	if err != nil {
		return nil, err
	}

	var tag Tag
	var deletedAt sql.NullTime
	var iconValue sql.NullString
	var themeValue sql.NullString

	if err := row.Scan(
		&tag.ID,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&deletedAt,
		&tag.Name,
		&tag.Slug,
		&iconValue,
		&themeValue,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
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

	return &tag, nil
}

func (t TagModel) Update(tag *Tag) error {
	var icon interface{}
	if tag.Icon != nil {
		icon = *tag.Icon
	}

	var theme interface{}
	if tag.Theme != nil {
		theme = *tag.Theme
	}

	values := querybuilder.Clauses{
		{ColumnName: "name", Value: tag.Name},
		{ColumnName: "slug", Value: tag.Slug},
		{ColumnName: "icon", Value: icon},
		{ColumnName: "theme", Value: theme},
		{ColumnName: "updatedAt", Value: time.Now()},
	}

	row, err := t.Query.SetBaseTable("tags").Update(values).WhereEqual("id", tag.ID).WhereEqual("deletedAt", nil).Returning(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"name",
		"slug",
		"icon",
		"theme",
	).QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var iconValue sql.NullString
	var themeValue sql.NullString

	if err := row.Scan(
		&tag.ID,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&deletedAt,
		&tag.Name,
		&tag.Slug,
		&iconValue,
		&themeValue,
	); err != nil {
		return err
	}

	if deletedAt.Valid {
		tag.DeletedAt = &deletedAt.Time
	} else {
		tag.DeletedAt = nil
	}

	if iconValue.Valid {
		value := iconValue.String
		tag.Icon = &value
	} else {
		tag.Icon = nil
	}

	if themeValue.Valid {
		value := themeValue.String
		tag.Theme = &value
	} else {
		tag.Theme = nil
	}

	return nil
}

func (t TagModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}

	values := querybuilder.Clauses{
		{ColumnName: "updatedAt", Value: time.Now()},
		{ColumnName: "deletedAt", Value: time.Now()},
	}

	results, err := t.Query.SetBaseTable("tags").Update(values).WhereEqual("id", id).WhereEqual("deletedAt", nil).Exec()
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

func (t TagModel) GetAll() ([]*Tag, error) {
	rows, err := t.Query.SetBaseTable("tags").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"name",
		"slug",
		"icon",
		"theme",
	).WhereEqual("deletedAt", nil).OrderBy("name", "asc").Query()
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
