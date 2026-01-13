package data

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"api.etin.dev/pkg/querybuilder"
)

type Asset struct {
	ID           int64      `json:"id"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `json:"-"`
	URL          string     `json:"url"`
	SecureURL    string     `json:"secureUrl"`
	PublicID     string     `json:"publicId"`
	Format       string     `json:"format"`
	ResourceType string     `json:"resourceType"`
	Bytes        int        `json:"bytes"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
}

type AssetModel struct {
	DB     *sql.DB
	Query  *querybuilder.QueryBuilder
	Logger *log.Logger
}

func (m AssetModel) Insert(asset *Asset) error {
	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "url", Value: asset.URL},
		querybuilder.Clause{ColumnName: "secureUrl", Value: asset.SecureURL},
		querybuilder.Clause{ColumnName: "publicId", Value: asset.PublicID},
		querybuilder.Clause{ColumnName: "format", Value: asset.Format},
		querybuilder.Clause{ColumnName: "resourceType", Value: asset.ResourceType},
		querybuilder.Clause{ColumnName: "bytes", Value: asset.Bytes},
		querybuilder.Clause{ColumnName: "width", Value: asset.Width},
		querybuilder.Clause{ColumnName: "height", Value: asset.Height},
	}

	row, err := m.Query.SetBaseTable("assets").Insert(values).Returning(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
	).QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime

	err = row.Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt, &deletedAt)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		asset.DeletedAt = &deletedAt.Time
	}

	return nil
}

func (m AssetModel) Get(id int64) (*Asset, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	row, err := m.Query.SetBaseTable("assets").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"url",
		"secureUrl",
		"publicId",
		"format",
		"resourceType",
		"bytes",
		"width",
		"height",
	).WhereEqual("deletedAt", nil).WhereEqual("id", id).QueryRow()
	if err != nil {
		return nil, err
	}

	var asset Asset
	var deletedAt sql.NullTime

	err = row.Scan(
		&asset.ID,
		&asset.CreatedAt,
		&asset.UpdatedAt,
		&deletedAt,
		&asset.URL,
		&asset.SecureURL,
		&asset.PublicID,
		&asset.Format,
		&asset.ResourceType,
		&asset.Bytes,
		&asset.Width,
		&asset.Height,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	if deletedAt.Valid {
		asset.DeletedAt = &deletedAt.Time
	}

	return &asset, nil
}
