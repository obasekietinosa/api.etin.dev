package data

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"api.etin.dev/pkg/querybuilder"
)

type Company struct {
	ID          int64      `json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
	Name        string     `json:"name"`
	Icon        *string    `json:"icon,omitempty"`
	Description *string    `json:"description,omitempty"`
}

type CompanyModel struct {
	DB     *sql.DB
	Query  *querybuilder.QueryBuilder
	Logger *log.Logger
}

func (c CompanyModel) Insert(company *Company) error {
	var icon interface{}
	if company.Icon != nil {
		icon = *company.Icon
	}
	var description interface{}
	if company.Description != nil {
		description = *company.Description
	}

	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "name", Value: company.Name},
		querybuilder.Clause{ColumnName: "icon", Value: icon},
		querybuilder.Clause{ColumnName: "description", Value: description},
	}

	row, err := c.Query.SetBaseTable("companies").Insert(values).Returning(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
	).QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime

	err = row.Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt, &deletedAt)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		company.DeletedAt = &deletedAt.Time
	}

	return nil
}

func (c CompanyModel) Get(id int64) (*Company, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	row, err := c.Query.SetBaseTable("companies").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"name",
		"icon",
		"description",
	).WhereEqual("deletedAt", nil).WhereEqual("id", id).QueryRow()
	if err != nil {
		return nil, err
	}

	var company Company
	var deletedAt sql.NullTime
	var icon sql.NullString
	var description sql.NullString

	err = row.Scan(
		&company.ID,
		&company.CreatedAt,
		&company.UpdatedAt,
		&deletedAt,
		&company.Name,
		&icon,
		&description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	if deletedAt.Valid {
		company.DeletedAt = &deletedAt.Time
	}

	if icon.Valid {
		company.Icon = &icon.String
	}

	if description.Valid {
		company.Description = &description.String
	}

	return &company, nil
}

func (c CompanyModel) Update(company *Company) error {
	var icon interface{}
	if company.Icon != nil {
		icon = *company.Icon
	}
	var description interface{}
	if company.Description != nil {
		description = *company.Description
	}

	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "name", Value: company.Name},
		querybuilder.Clause{ColumnName: "icon", Value: icon},
		querybuilder.Clause{ColumnName: "description", Value: description},
		querybuilder.Clause{ColumnName: "updatedAt", Value: time.Now()},
	}

	row, err := c.Query.With(
		c.Query.SetBaseTable("companies").Update(values).WhereEqual("id", company.ID).WhereEqual("deletedAt", nil).Returning(
			"id",
			"createdAt",
			"updatedAt",
			"deletedAt",
			"name",
			"icon",
			"description",
		),
		"updated_company",
	).Select(
		"updated_company.id",
		"updated_company.createdAt",
		"updated_company.updatedAt",
		"updated_company.deletedAt",
		"updated_company.name",
		"updated_company.icon",
		"updated_company.description",
	).From("updated_company").QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var iconVal sql.NullString
	var descriptionVal sql.NullString

	err = row.Scan(
		&company.ID,
		&company.CreatedAt,
		&company.UpdatedAt,
		&deletedAt,
		&company.Name,
		&iconVal,
		&descriptionVal,
	)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		company.DeletedAt = &deletedAt.Time
	} else {
		company.DeletedAt = nil
	}

	if iconVal.Valid {
		company.Icon = &iconVal.String
	} else {
		company.Icon = nil
	}

	if descriptionVal.Valid {
		company.Description = &descriptionVal.String
	} else {
		company.Description = nil
	}

	return nil
}

func (c CompanyModel) Delete(company *Company) error {
	if company.ID < 1 {
		return errors.New("record not found")
	}

	values := querybuilder.Clauses{
		querybuilder.Clause{ColumnName: "updatedAt", Value: time.Now()},
		querybuilder.Clause{ColumnName: "deletedAt", Value: time.Now()},
	}

	results, err := c.Query.SetBaseTable("companies").Update(values).WhereEqual("id", company.ID).WhereEqual("deletedAt", nil).Exec()
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

func (c CompanyModel) GetAll() ([]*Company, error) {
	rows, err := c.Query.SetBaseTable("companies").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"name",
		"icon",
		"description",
	).WhereEqual("deletedAt", nil).OrderBy("createdAt", "desc").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*Company{}

	for rows.Next() {
		var company Company
		var deletedAt sql.NullTime
		var icon sql.NullString
		var description sql.NullString

		err := rows.Scan(
			&company.ID,
			&company.CreatedAt,
			&company.UpdatedAt,
			&deletedAt,
			&company.Name,
			&icon,
			&description,
		)
		if err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			company.DeletedAt = &deletedAt.Time
		}

		if icon.Valid {
			company.Icon = &icon.String
		}

		if description.Valid {
			company.Description = &description.String
		}

		companies = append(companies, &company)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}
