package querybuilder

import (
	"database/sql"
	"errors"
	"fmt"
)

type UpdateQueryBuilder struct {
	queryBuilder QueryBuilder
	table        string
	values       ClauseMap
	conditions   ClauseMap
}

func (q UpdateQueryBuilder) buildColumnUpdateStatement() string {
	keys := make([]string, 0, len(q.values))
	for k := range q.values {
		keys = append(keys, k)
	}

	stmt := ""

	if len(keys) != 0 {
		stmt += " SET"

		for i, column := range keys {
			if i > 0 {
				stmt += ","
			}
			stmt += fmt.Sprintf(" %s = $%d", column, i+1)
		}
	}

	return stmt
}

func (q UpdateQueryBuilder) buildPreparedStatementValues() []interface{} {
	updateValues := q.queryBuilder.buildParameters(q.values)
	conditionValues := q.queryBuilder.buildParameters(q.conditions)

	return append(updateValues, conditionValues...)
}

func (q UpdateQueryBuilder) WhereEqual(column string, value interface{}) UpdateQueryBuilder {
	q.queryBuilder.addCondition(column, value, "=", &q.conditions)
	return q
}

func (q UpdateQueryBuilder) buildQuery() (*string, error) {
	if len(q.values) == 0 || q.table == "" {
		err := errors.New("Incorrectly formatted query. Ensure fields and base tables are set")
		return nil, err
	}

	query := fmt.Sprintf("UPDATE %s", q.table)
	query += q.buildColumnUpdateStatement()
	query += q.queryBuilder.buildConditionalStatement(q.conditions, len(q.values))

	return &query, nil
}

func (q UpdateQueryBuilder) Query() (*sql.Rows, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}
	values := q.buildPreparedStatementValues()
	return q.queryBuilder.DB.Query(*query, values...)
}

func (q UpdateQueryBuilder) QueryRow() (*sql.Row, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}

	values := q.buildPreparedStatementValues()
	return q.queryBuilder.DB.QueryRow(*query, values...), nil
}

func (q UpdateQueryBuilder) Exec() (sql.Result, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}

	values := q.buildPreparedStatementValues()
	return q.queryBuilder.DB.Exec(*query, values...)
}
