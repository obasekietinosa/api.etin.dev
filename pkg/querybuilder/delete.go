package querybuilder

import (
	"database/sql"
	"errors"
	"fmt"
)

type DeleteQueryBuilder struct {
	queryBuilder QueryBuilder
	table        string
	conditions   ClauseMap
}

func (q DeleteQueryBuilder) buildQuery() (*string, error) {
	if len(q.conditions) == 0 {
		err := errors.New("Incorrectly formatted query. Ensure fields are set")
		return nil, err
	}

	if q.table == "" {
		err := errors.New("Incorrectly formatted query. Ensure base tables is set")
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %s", q.table)
	query += q.queryBuilder.buildConditionalStatement(q.conditions)

	return &query, nil
}

func (q DeleteQueryBuilder) WhereEqual(column string, value interface{}) DeleteQueryBuilder {
	q.queryBuilder.addCondition(column, value, "=", &q.conditions)
	return q
}

func (q DeleteQueryBuilder) Exec() (sql.Result, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}

	values := q.queryBuilder.buildParameters(q.conditions)
	return q.queryBuilder.DB.Exec(*query, values...)
}
