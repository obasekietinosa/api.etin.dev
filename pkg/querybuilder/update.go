package querybuilder

import (
	"database/sql"
	"errors"
	"fmt"
)

type UpdateQueryBuilder struct {
	queryBuilder *QueryBuilder
	table        string
	values       Clauses
	conditions   Clauses
	fields       []string
}

func (q *UpdateQueryBuilder) buildPreparedStatementValues() []interface{} {
	values := make([]interface{}, 0)

	if len(q.queryBuilder.commonTableExpressions) > 0 {
		for _, cte := range q.queryBuilder.commonTableExpressions {
			values = append(values, cte.buildPreparedStatementValues()...)
		}
	}

	updateValues := q.queryBuilder.buildParameters(q.values)
	conditionValues := q.queryBuilder.buildParameters(q.conditions)

	values = append(values, updateValues...)
	values = append(values, conditionValues...)
	return values
}

func (q *UpdateQueryBuilder) WhereEqual(column string, value interface{}) *UpdateQueryBuilder {
	q.queryBuilder.addCondition(column, value, "=", &q.conditions)
	return q
}

func (q *UpdateQueryBuilder) Returning(fields ...string) *UpdateQueryBuilder {
	q.fields = fields
	return q
}

func (q *UpdateQueryBuilder) buildQuery() (*string, error) {
	if len(q.values) == 0 || q.table == "" {
		err := errors.New("Incorrectly formatted query. Ensure fields and base tables are set")
		return nil, err
	}

	query := fmt.Sprintf("UPDATE %s", q.table)
	query += q.queryBuilder.buildColumnUpdateStatement(q.values)
	query += q.queryBuilder.buildConditionalStatement(q.conditions)

	if len(q.fields) > 0 {
		returnedColumns := q.queryBuilder.buildReturnedColumns(q.fields)
		query += fmt.Sprintf(" RETURNING %s", returnedColumns)
	}

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
