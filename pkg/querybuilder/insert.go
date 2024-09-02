package querybuilder

import (
	"database/sql"
	"errors"
	"fmt"
)

type InsertQueryBuilder struct {
	queryBuilder QueryBuilder
	table        string
	values       ClauseMap
	fields       []string
}

func (q InsertQueryBuilder) buildValuesStatement() string {
	keys := make([]string, 0, len(q.values))
	for k := range q.values {
		keys = append(keys, k)
	}

	stmt := ""

	if len(keys) != 0 {
		for i := range keys {
			if i > 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("$%d", i+1)
		}
	}

	return stmt
}

func (q InsertQueryBuilder) buildColumnNameStatement() string {
	keys := make([]string, 0, len(q.values))
	for k := range q.values {
		keys = append(keys, k)
	}

	stmt := ""

	if len(keys) != 0 {
		for i, column := range keys {
			if i > 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("%s", column)
		}
	}

	return stmt
}

func (q InsertQueryBuilder) buildReturnedColumns() string {
	stmt := ""

	if len(q.fields) != 0 {
		for i, field := range q.fields {
			if i > 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("%s", field)
		}
	}

	return stmt
}

func (q InsertQueryBuilder) buildQuery() (*string, error) {
	if len(q.values) == 0 {
		err := errors.New("Incorrectly formatted query. Ensure fields are set")
		return nil, err
	}
	if q.table == "" {
		err := errors.New("Incorrectly formatted query. Ensure base table is set.")
		return nil, err
	}

	columnNames := q.buildColumnNameStatement()
	columnValues := q.buildValuesStatement()
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.table, columnNames, columnValues)
	if len(q.fields) > 0 {
		returnedColumns := q.buildReturnedColumns()
		query += fmt.Sprintf(" RETURNING %s", returnedColumns)
	}

	return &query, nil
}

func (q InsertQueryBuilder) Returning(fields ...string) InsertQueryBuilder {
	q.fields = fields
	return q
}

func (q InsertQueryBuilder) Query() (*sql.Rows, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}
	values := q.queryBuilder.buildParameters(q.values)
	return q.queryBuilder.DB.Query(*query, values...)
}

func (q InsertQueryBuilder) QueryRow() (*sql.Row, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}

	values := q.queryBuilder.buildParameters(q.values)
	return q.queryBuilder.DB.QueryRow(*query, values...), nil
}

func (q InsertQueryBuilder) Exec() (sql.Result, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}

	values := q.queryBuilder.buildParameters(q.values)
	return q.queryBuilder.DB.Exec(*query, values...)
}
