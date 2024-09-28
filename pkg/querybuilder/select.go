package querybuilder

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type SelectQueryBuilder struct {
	queryBuilder QueryBuilder

	fields     []string
	table      string
	conditions ClauseMap

	sortDirection string
	sortColumn    string

	leftJoinTable      string
	leftJoinOwnKey     string
	leftJoinForeignKey string
}

func (q SelectQueryBuilder) From(table string) SelectQueryBuilder {
	q.table = table
	return q
}

func (q SelectQueryBuilder) LeftJoin(table string, ownKey string, foreignKey string) SelectQueryBuilder {
	q.leftJoinTable = table
	q.leftJoinOwnKey = ownKey
	q.leftJoinForeignKey = foreignKey

	return q
}

func (q SelectQueryBuilder) OrderBy(column string, sortDirection string) SelectQueryBuilder {
	q.sortColumn = column
	q.sortDirection = sortDirection

	return q
}

func (q SelectQueryBuilder) WhereEqual(column string, value interface{}) SelectQueryBuilder {
	if value == nil {
		q.queryBuilder.addCondition(column, nil, "IS NULL", &q.conditions)
	} else {
		q.queryBuilder.addCondition(column, value, "=", &q.conditions)
	}
	return q
}

func (q SelectQueryBuilder) buildPreparedStatementValues() []interface{} {
	values := make([]interface{}, 0)

	if len(q.queryBuilder.commonTableExpressions) > 0 {
		for _, cte := range q.queryBuilder.commonTableExpressions {
			values = append(values, cte.buildPreparedStatementValues()...)
		}
	}

	values = append(values, q.queryBuilder.buildParameters(q.conditions)...)

	return values
}

func (q SelectQueryBuilder) buildQuery() (*string, error) {
	if len(q.fields) == 0 || q.table == "" {
		err := errors.New("Incorrectly formatted query. Ensure fields and base tables are set")
		return nil, err
	}

	fields := strings.Join(q.fields, ", ")

	query := fmt.Sprintf("SELECT %s FROM %s", fields, q.table)

	if q.leftJoinTable != "" {
		query += fmt.Sprintf(" LEFT JOIN %s ON %s.%s = %s.%s", q.leftJoinTable, q.table, q.leftJoinOwnKey, q.leftJoinTable, q.leftJoinForeignKey)
	}

	query += q.queryBuilder.buildConditionalStatement(q.conditions)

	if q.sortColumn != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", q.sortColumn, q.sortDirection)
	}

	return &query, nil
}

func (q SelectQueryBuilder) Query() (*sql.Rows, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}
	if len(q.conditions) > 0 {
		values := q.queryBuilder.buildParameters(q.conditions)
		return q.queryBuilder.DB.Query(*query, values...)
	}
	return q.queryBuilder.DB.Query(*query)
}

func (q SelectQueryBuilder) QueryRow() (*sql.Row, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}
	if len(q.conditions) > 0 || len(q.queryBuilder.commonTableExpressions) > 0 {
		values := q.buildPreparedStatementValues()
		return q.queryBuilder.DB.QueryRow(*query, values...), nil
	}
	return q.queryBuilder.DB.QueryRow(*query), nil
}
