package querybuilder

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type QueryBuilder struct {
	DB         *sql.DB
	fields     []string
	conditions []string
	table      string
}

func (q QueryBuilder) Select(fields ...string) QueryBuilder {
	q.fields = fields
	return q
}

func (q QueryBuilder) From(table string) QueryBuilder {
	q.table = table
	return q
}

func (q QueryBuilder) Where(conditions []string) QueryBuilder {

	return q
}

func (q QueryBuilder) BuildQuery() (*string, error) {
	if len(q.fields) == 0 || q.table == "" {
		err := errors.New("Incorrectly formatted query. Ensure fields and base tables are set")
		return nil, err
	}

	fields := strings.Join(q.fields, ", ")

	query := fmt.Sprintf("SELECT %s FROM %s", fields, q.table)

	if len(q.conditions) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(q.conditions, "AND "))
	}

	return &query, nil
}

func (q QueryBuilder) Query() (*sql.Rows, error) {
	query, err := q.BuildQuery()
	if err != nil {
		return nil, err
	}

	return q.DB.Query(*query)
}
