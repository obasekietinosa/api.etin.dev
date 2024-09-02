package data

import (
	"database/sql"
	"errors"
	"time"

	"api.etin.dev/pkg/querybuilder"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Role struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	CompanyId   int64     `json:"companyId"`
	Company     string    `json:"company"`
	CompanyIcon string    `json:"companyIcon"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Skills      []string  `json:"skills"`
}

type RoleModel struct {
	DB    *sql.DB
	Query *querybuilder.QueryBuilder
}

func (r RoleModel) Insert(role *Role) error {
	query := `
		WITH inserted_role as (
			INSERT INTO roles (
				startDate, 
				endDate,
				title,
				subtitle,
				companyId,
				slug,
				description,
				skills
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING *
		) SELECT 
				inserted_role.id as id, 
				inserted_role.createdAt as createdAt, 
				inserted_role.updatedAt as updatedAt, 
				companies.name as company,
				companies.icon as companyIcon 
			FROM inserted_role
			LEFT JOIN companies ON inserted_role.companyId = companies.id;
	`

	args := []interface{}{role.StartDate, role.EndDate, role.Title, role.Subtitle, role.CompanyId, role.Slug, role.Description, pq.Array(role.Skills)}
	return r.DB.QueryRow(query, args...).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Company, &role.CompanyIcon)
}

func (r RoleModel) Get(roleId int64) (*Role, error) {
	if roleId < 1 {
		return nil, errors.New("record not found")
	}

	row, err := r.Query.SetBaseTable("roles").
		Select(
			"roles.id AS id", "roles.createdAt AS createdAt", "roles.updatedAt AS updatedAt", "roles.startDate AS startDate",
			"roles.endDate AS endDate", "roles.title AS title", "roles.subtitle AS subtitle", "roles.slug AS slug",
			"roles.description AS description", "roles.skills AS skills",
			"companies.id as companyId", "companies.name as company", "companies.icon as companyIcon").
		LeftJoin("companies", "id", "companyId").
		WhereEqual("roles.deletedAt", nil).
		WhereEqual("roles.id", roleId).
		QueryRow()

	var role Role

	err = row.Scan(
		&role.ID,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.StartDate,
		&role.EndDate,
		&role.Title,
		&role.Subtitle,
		&role.Slug,
		&role.Description,
		pq.Array(&role.Skills),
		&role.CompanyId,
		&role.Company,
		&role.CompanyIcon,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}
	}

	return &role, nil
}

func (r RoleModel) Update(role *Role) error {
	query := `
		WITH updated_role as (
			UPDATE roles
			SET 
				startDate = $1,
				endDate = $2,
				title = $3,
				subtitle = $4,
				slug = $5,
				description = $6,
				skills = $7,
				companyId = $8,
				updatedAt = NOW()
			WHERE id = $9
			RETURNING *
		) 
		SELECT 
			updated_role.updatedAt AS updatedAt,
			companies.id as companyId,
			companies.name as company,
			companies.icon as companyIcon
		FROM updated_role
		LEFT JOIN companies ON updated_role.companyId = companies.id;
	`

	args := []interface{}{role.StartDate, role.EndDate, role.Title, role.Subtitle, role.Slug, role.Description, pq.Array(role.Skills), role.CompanyId, role.ID}
	return r.DB.QueryRow(query, args...).Scan(&role.UpdatedAt, &role.CompanyId, &role.Company, &role.CompanyIcon)
}

func (r RoleModel) Delete(roleId int64) error {
	if roleId < 1 {
		return errors.New("No record found")
	}

	query := `
		UPDATE roles
		SET updatedAt = NOW(), deletedAt = NOW()
		WHERE id = $1;
	`

	results, err := r.DB.Exec(query, roleId)
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

func (r RoleModel) GetAll() ([]*Role, error) {
	rows, err := r.Query.SetBaseTable("roles").Select(
		"roles.id AS id", "roles.createdAt AS createdAt", "roles.updatedAt AS updatedAt", "roles.startDate AS startDate",
		"roles.endDate AS endDate", "roles.title AS title", "roles.subtitle AS subtitle", "roles.slug AS slug",
		"roles.description AS description", "roles.skills AS skills", "roles.companyId AS companyId",
		"companies.name AS company", "companies.icon AS companyIcon",
	).
		LeftJoin("companies", "id", "companyId").
		OrderBy("startDate", "desc").
		Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := []*Role{}

	for rows.Next() {
		var role Role
		err := rows.Scan(
			&role.ID,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.StartDate,
			&role.EndDate,
			&role.Title,
			&role.Subtitle,
			&role.Slug,
			&role.Description,
			pq.Array(&role.Skills),
			&role.CompanyId,
			&role.Company,
			&role.CompanyIcon,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, &role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
