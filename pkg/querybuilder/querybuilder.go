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
	conditions map[string]interface{}
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

func (q QueryBuilder) WhereEqual(column string, value interface{}) QueryBuilder {
	q.initialiseConditions()
	key := fmt.Sprintf("%s:%s", column, "=")
	q.conditions[key] = value
	return q
}

func (q *QueryBuilder) initialiseConditions() {
	if q.conditions == nil {
		q.conditions = make(map[string]interface{})
	}
}

func (q QueryBuilder) buildConditionalStatement() string {
	keys := make([]string, 0, len(q.conditions))
	for k := range q.conditions {
		keys = append(keys, k)
	}

	stmt := ""

	if len(keys) != 0 || q.conditions != nil {
		stmt += " WHERE"

		for i, v := range keys {
			if i > 0 {
				stmt += " AND"
			}
			columnAndComparator := strings.Split(v, ":")
			stmt += fmt.Sprintf(" %s %s $%d", columnAndComparator[0], columnAndComparator[1], i+1)
		}
	}

	return stmt
}

func (q QueryBuilder) BuildQuery() (*string, error) {
	if len(q.fields) == 0 || q.table == "" {
		err := errors.New("Incorrectly formatted query. Ensure fields and base tables are set")
		return nil, err
	}

	fields := strings.Join(q.fields, ", ")

	query := fmt.Sprintf("SELECT %s FROM %s", fields, q.table)
	query += q.buildConditionalStatement()

	return &query, nil
}

func (q QueryBuilder) Query() (*sql.Rows, error) {
	query, err := q.BuildQuery()
	if err != nil {
		return nil, err
	}
	if len(q.conditions) > 0 {
		values := make([]interface{}, 0, len(q.conditions))
		for _, v := range q.conditions {
			values = append(values, v)
		}

		return q.DB.Query(*query, values...)
	}
	return q.DB.Query(*query)
}

func (q QueryBuilder) QueryRow() (*sql.Row, error) {
	query, err := q.BuildQuery()
	if err != nil {
		return nil, err
	}
	if len(q.conditions) > 0 {
		values := make([]interface{}, 0, len(q.conditions))
		for _, v := range q.conditions {
			values = append(values, v)
		}
		return q.DB.QueryRow(*query, values...), nil
	}
	return q.DB.QueryRow(*query), nil
}
