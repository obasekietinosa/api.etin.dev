package data

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type Company struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type CompanyModel struct {
	DB *sql.DB
}

func (c CompanyModel) GetAll() ([]*Company, error) {
	query := `
			SELECT id, name, icon, description FROM companies;
	`

	rows, err := c.DB.Query(query)
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

	query := `
		SELECT * from companies
		WHERE id = $1;
	`

	var company Company

	err := c.DB.QueryRow(query, companyId).Scan(&company.ID, &company.Name, &company.Icon, &company.Description)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (c CompanyModel) Update(company *Company) error {
	query := `
		UPDATE companies
		SET
			name = $1
			icon = $2
			description = $3
		WHERE id = $4;
	`

	args := []interface{}{company.Name, company.Icon, company.Description, company.ID}
	return c.DB.QueryRow(query, args...).Scan()
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
