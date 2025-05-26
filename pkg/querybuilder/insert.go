package querybuilder

import (
	"database/sql"
	"errors"
	"fmt"
)

type InsertQueryBuilder struct {
	queryBuilder *QueryBuilder
	table        string
	values       Clauses
	fields       []string
}

func (q *InsertQueryBuilder) buildPreparedStatementValues() []interface{} {
	values := q.queryBuilder.buildCommonTableExpressionParameters()
	values = append(values, q.queryBuilder.buildParameters(q.values)...)

	return values
}

func (q *InsertQueryBuilder) buildColumnNameStatement() string {
	stmt := ""

	if len(q.values) != 0 {
		for i, column := range q.values {
			if i > 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("%s", column.ColumnName)
		}
	}

	return stmt
}

func (q *InsertQueryBuilder) buildQuery() (*string, error) {
	if len(q.values) == 0 {
		err := errors.New("Incorrectly formatted query. Ensure fields are set")
		return nil, err
	}
	if q.table == "" {
		err := errors.New("Incorrectly formatted query. Ensure base table is set.")
		return nil, err
	}

	columnNames := q.buildColumnNameStatement()
	columnValues := q.queryBuilder.buildValuesStatement(q.values)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.table, columnNames, columnValues)
	if len(q.fields) > 0 {
		returnedColumns := q.queryBuilder.buildReturnedColumns(q.fields)
		query += fmt.Sprintf(" RETURNING %s", returnedColumns)
	}

	return &query, nil
}

func (q *InsertQueryBuilder) Returning(fields ...string) *InsertQueryBuilder {
	q.fields = fields
	return q
}

func (q *InsertQueryBuilder) Query() (*sql.Rows, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}
	values := q.buildPreparedStatementValues()
	return q.queryBuilder.DB.Query(*query, values...)
}

func (q *InsertQueryBuilder) QueryRow() (*sql.Row, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}

	values := q.buildPreparedStatementValues()
	return q.queryBuilder.DB.QueryRow(*query, values...), nil
}

func (q *InsertQueryBuilder) Exec() (sql.Result, error) {
	query, err := q.buildQuery()
	if err != nil {
		return nil, err
	}

	values := q.buildPreparedStatementValues()
	return q.queryBuilder.DB.Exec(*query, values...)
}
