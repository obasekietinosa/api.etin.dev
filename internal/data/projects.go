package data

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"api.etin.dev/pkg/querybuilder"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

type Project struct {
	ID          int64      `json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	ImageURL    *string    `json:"imageUrl,omitempty"`
}

type ProjectModel struct {
	DB     *sql.DB
	Query  *querybuilder.QueryBuilder
	Logger *log.Logger
}

func (p ProjectModel) Insert(project *Project) error {
	var endDate interface{}
	if project.EndDate != nil {
		endDate = *project.EndDate
	}

	var imageURL interface{}
	if project.ImageURL != nil {
		imageURL = *project.ImageURL
	}

	if project.Slug == "" {
		project.Slug = slug.Make(project.Title)
	}

	if err := p.ensureUniqueSlug(project); err != nil {
		return err
	}

	values := querybuilder.Clauses{
		{ColumnName: "startDate", Value: project.StartDate},
		{ColumnName: "endDate", Value: endDate},
		{ColumnName: "title", Value: project.Title},
		{ColumnName: "slug", Value: project.Slug},
		{ColumnName: "description", Value: project.Description},
		{ColumnName: "imageUrl", Value: imageURL},
	}

	row, err := p.Query.SetBaseTable("projects").Insert(values).Returning(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"startDate",
		"endDate",
		"title",
		"slug",
		"description",
		"imageUrl",
	).QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var savedEndDate sql.NullTime
	var savedImageURL sql.NullString
	var slug sql.NullString

	err = row.Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&deletedAt,
		&project.StartDate,
		&savedEndDate,
		&project.Title,
		&slug,
		&project.Description,
		&savedImageURL,
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

	if savedImageURL.Valid {
		project.ImageURL = &savedImageURL.String
	} else {
		project.ImageURL = nil
	}

	if slug.Valid {
		project.Slug = slug.String
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
		"slug",
		"description",
		"imageUrl",
	).WhereEqual("deletedAt", nil).WhereEqual("id", projectID).QueryRow()
	if err != nil {
		return nil, err
	}

	var project Project
	var deletedAt sql.NullTime
	var endDate sql.NullTime
	var imageURL sql.NullString
	var slug sql.NullString

	err = row.Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&deletedAt,
		&project.StartDate,
		&endDate,
		&project.Title,
		&slug,
		&project.Description,
		&imageURL,
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

	if imageURL.Valid {
		project.ImageURL = &imageURL.String
	}

	if slug.Valid {
		project.Slug = slug.String
	}

	return &project, nil
}

func (p ProjectModel) Update(project *Project) error {
	var endDate interface{}
	if project.EndDate != nil {
		endDate = *project.EndDate
	}

	var imageURL interface{}
	if project.ImageURL != nil {
		imageURL = *project.ImageURL
	}

	if project.Slug == "" {
		project.Slug = slug.Make(project.Title)
	}

	if err := p.ensureUniqueSlug(project); err != nil {
		return err
	}

	values := querybuilder.Clauses{
		{ColumnName: "startDate", Value: project.StartDate},
		{ColumnName: "endDate", Value: endDate},
		{ColumnName: "title", Value: project.Title},
		{ColumnName: "slug", Value: project.Slug},
		{ColumnName: "description", Value: project.Description},
		{ColumnName: "imageUrl", Value: imageURL},
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
		"slug",
		"description",
		"imageUrl",
	).QueryRow()
	if err != nil {
		return err
	}

	var deletedAt sql.NullTime
	var savedEndDate sql.NullTime
	var savedImageURL sql.NullString
	var slug sql.NullString

	err = row.Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&deletedAt,
		&project.StartDate,
		&savedEndDate,
		&project.Title,
		&slug,
		&project.Description,
		&savedImageURL,
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

	if savedImageURL.Valid {
		project.ImageURL = &savedImageURL.String
	} else {
		project.ImageURL = nil
	}

	if slug.Valid {
		project.Slug = slug.String
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
		"slug",
		"description",
		"imageUrl",
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
		var imageURL sql.NullString
		var slug sql.NullString

		err := rows.Scan(
			&project.ID,
			&project.CreatedAt,
			&project.UpdatedAt,
			&deletedAt,
			&project.StartDate,
			&endDate,
			&project.Title,
			&slug,
			&project.Description,
			&imageURL,
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

		if imageURL.Valid {
			project.ImageURL = &imageURL.String
		}

		if slug.Valid {
			project.Slug = slug.String
		}

		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p ProjectModel) GetByIDs(ids []int64) ([]*Project, error) {
	if len(ids) == 0 {
		return []*Project{}, nil
	}

	query := `
        SELECT id, createdAt, updatedAt, deletedAt, startDate, endDate, title, slug, description, imageUrl
        FROM projects
        WHERE deletedAt IS NULL AND id = ANY($1)
    `

	rows, err := p.DB.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []*Project{}

	for rows.Next() {
		var project Project
		var deletedAt sql.NullTime
		var endDate sql.NullTime
		var imageURL sql.NullString
		var slug sql.NullString

		err := rows.Scan(
			&project.ID,
			&project.CreatedAt,
			&project.UpdatedAt,
			&deletedAt,
			&project.StartDate,
			&endDate,
			&project.Title,
			&slug,
			&project.Description,
			&imageURL,
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

		if imageURL.Valid {
			project.ImageURL = &imageURL.String
		}

		if slug.Valid {
			project.Slug = slug.String
		}

		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p ProjectModel) GetBySlug(slug string) (*Project, error) {
	row, err := p.Query.SetBaseTable("projects").Select(
		"id",
		"createdAt",
		"updatedAt",
		"deletedAt",
		"startDate",
		"endDate",
		"title",
		"slug",
		"description",
		"imageUrl",
	).WhereEqual("deletedAt", nil).WhereEqual("slug", slug).QueryRow()
	if err != nil {
		return nil, err
	}

	var project Project
	var deletedAt sql.NullTime
	var endDate sql.NullTime
	var imageURL sql.NullString
	var slugVal sql.NullString

	err = row.Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
		&deletedAt,
		&project.StartDate,
		&endDate,
		&project.Title,
		&slugVal,
		&project.Description,
		&imageURL,
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

	if imageURL.Valid {
		project.ImageURL = &imageURL.String
	}

	if slugVal.Valid {
		project.Slug = slugVal.String
	}

	return &project, nil
}

func (p ProjectModel) ensureUniqueSlug(project *Project) error {
	originalSlug := project.Slug
	counter := 1

	for {
		var count int
		row, err := p.Query.SetBaseTable("projects").Select("COUNT(*)").
			WhereEqual("slug", project.Slug).
			WhereNotEqual("id", project.ID).
			WhereEqual("deletedAt", nil).
			QueryRow()
		if err != nil {
			return err
		}

		if err := row.Scan(&count); err != nil {
			return err
		}

		if count == 0 {
			break
		}

		project.Slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}
	return nil
}
