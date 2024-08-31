package data

import (
	"database/sql"
	"errors"

	"api.etin.dev/pkg/querybuilder"
	_ "github.com/lib/pq"
)

type Company struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type CompanyModel struct {
	DB    *sql.DB
	Query *querybuilder.QueryBuilder
}

func (c CompanyModel) GetAll() ([]*Company, error) {
	rows, err := c.Query.Select(
		"id",
		"name",
		"icon",
		"description").From("companies").OrderBy("id", "asc").Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	companies := []*Company{}

	for rows.Next() {
		company := &Company{}

		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.Icon,
			&company.Description,
		)

		if err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func (c CompanyModel) Insert(company *Company) error {
	query := `
		INSERT INTO companies (
			name,
			icon,
			description
		) VALUES ($1, $2, $3)
		RETURNING id;
	`

	args := []interface{}{company.Name, company.Icon, company.Description}
	return c.DB.QueryRow(query, args...).Scan(&company.ID)
}

func (c CompanyModel) Get(companyId int64) (*Company, error) {
	if companyId < 1 {
		return nil, errors.New("No record found")
	}

	row, err := c.Query.Select("id", "name", "icon", "description").From("companies").WhereEqual("id", companyId).QueryRow()
	if err != nil {
		return nil, err
	}

	var company Company

	err = row.Scan(&company.ID, &company.Name, &company.Icon, &company.Description)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (c CompanyModel) Update(company *Company) error {
	values := querybuilder.ClauseMap{
		"name":        company.Name,
		"icon":        company.Icon,
		"description": company.Description,
	}

	results, err := c.Query.SetBaseTable("companies").Update(values).WhereEqual("id", company.ID).Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No record found")
	}

	return nil
}

func (c CompanyModel) Delete(company *Company) error {
	query := `
		DELETE from companies
		WHERE id = $1;
	`

	results, err := c.DB.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No record found")
	}

	return nil
}
