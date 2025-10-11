package data

import (
	"database/sql"
	"errors"
	"time"

	"api.etin.dev/pkg/querybuilder"
)

type Project struct {
	ID          int64      `json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

type ProjectModel struct {
	DB    *sql.DB
	Query *querybuilder.QueryBuilder
}

func (p ProjectModel) Insert(project *Project) error {
	var endDate interface{}
	if project.EndDate != nil {
		endDate = *project.EndDate
	}

	values := querybuilder.Clauses{
		{ColumnName: "startDate", Value: project.StartDate},
		{ColumnName: "endDate", Value: endDate},
		{ColumnName: "title", Value: project.Title},
		{ColumnName: "description", Value: project.Description},
	}

	row, err := p.Query.SetBaseTable("projects").Insert(values).Returning(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"startDate",
		"endDate",
		"title",
		"description",
	).QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var savedEndDate sql.NullTime

	err = row.Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&deletedAt,
		&project.StartDate,
		&savedEndDate,
		&project.Title,
		&project.Description,
	)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		project.DeletedAt = &deletedAt.Time
	} else {
		project.DeletedAt = nil
	}

	if savedEndDate.Valid {
		project.EndDate = &savedEndDate.Time
	} else {
		project.EndDate = nil
	}

	return nil
}

func (p ProjectModel) Get(projectID int64) (*Project, error) {
	if projectID < 1 {
		return nil, errors.New("record not found")
	}

	row, err := p.Query.SetBaseTable("projects").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"startDate",
		"endDate",
		"title",
		"description",
	).WhereEqual("deletedAt", nil).WhereEqual("id", projectID).QueryRow()
	if err != nil {
		return nil, err
	}

	var project Project
	var deletedAt sql.NullTime
	var endDate sql.NullTime

	err = row.Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&deletedAt,
		&project.StartDate,
		&endDate,
		&project.Title,
		&project.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	if deletedAt.Valid {
		project.DeletedAt = &deletedAt.Time
	}

	if endDate.Valid {
		project.EndDate = &endDate.Time
	}

	return &project, nil
}

func (p ProjectModel) Update(project *Project) error {
	var endDate interface{}
	if project.EndDate != nil {
		endDate = *project.EndDate
	}

	values := querybuilder.Clauses{
		{ColumnName: "startDate", Value: project.StartDate},
		{ColumnName: "endDate", Value: endDate},
		{ColumnName: "title", Value: project.Title},
		{ColumnName: "description", Value: project.Description},
		{ColumnName: "updatedAt", Value: time.Now()},
	}

	row, err := p.Query.SetBaseTable("projects").Update(values).WhereEqual("id", project.ID).WhereEqual("deletedAt", nil).Returning(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"startDate",
		"endDate",
		"title",
		"description",
	).QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var savedEndDate sql.NullTime

	err = row.Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&deletedAt,
		&project.StartDate,
		&savedEndDate,
		&project.Title,
		&project.Description,
	)
	if err != nil {
		return err
	}

	if deletedAt.Valid {
		project.DeletedAt = &deletedAt.Time
	} else {
		project.DeletedAt = nil
	}

	if savedEndDate.Valid {
		project.EndDate = &savedEndDate.Time
	} else {
		project.EndDate = nil
	}

	return nil
}

func (p ProjectModel) Delete(projectID int64) error {
	if projectID < 1 {
		return errors.New("record not found")
	}

	values := querybuilder.Clauses{
		{ColumnName: "updatedAt", Value: time.Now()},
		{ColumnName: "deletedAt", Value: time.Now()},
	}

	results, err := p.Query.SetBaseTable("projects").Update(values).WhereEqual("id", projectID).WhereEqual("deletedAt", nil).Exec()
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

func (p ProjectModel) GetAll() ([]*Project, error) {
	rows, err := p.Query.SetBaseTable("projects").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"startDate",
		"endDate",
		"title",
		"description",
	).WhereEqual("deletedAt", nil).OrderBy("startDate", "desc").Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projects := []*Project{}

	for rows.Next() {
		var project Project
		var deletedAt sql.NullTime
		var endDate sql.NullTime

		err := rows.Scan(
			&project.ID,
			&project.CreatedAt,
			&project.UpdatedAt,
			&deletedAt,
			&project.StartDate,
			&endDate,
			&project.Title,
			&project.Description,
		)
		if err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			project.DeletedAt = &deletedAt.Time
		}

		if endDate.Valid {
			project.EndDate = &endDate.Time
		}

		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
