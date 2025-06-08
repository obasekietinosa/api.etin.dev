package querybuilder

import (
	"database/sql"
	"fmt"
	"strings"
)

type Clause struct {
	ColumnName string
	Value      interface{}
}
type Clauses []Clause

type CommonQuery struct {
	Builder CommonQueryBuilder
	Table   string
}

type QueryBuilder struct {
	DB                     *sql.DB
	table                  string
	preparedVariableOffset int
	commonTableExpressions []CommonQuery
}

type CommonQueryBuilder interface {
	buildQuery() (*string, error)
	buildPreparedStatementValues() []interface{}
}

func (q *QueryBuilder) SetBaseTable(table string) *QueryBuilder {
	q.table = table
	return q
}

func (q *QueryBuilder) clone() *QueryBuilder {
	return &QueryBuilder{
		DB:                     q.DB,
		table:                  q.table,
		commonTableExpressions: q.commonTableExpressions,
		preparedVariableOffset: 0,
	}
}

func (q *QueryBuilder) Select(fields ...string) *SelectQueryBuilder {
	cloned := q.clone()
	return &SelectQueryBuilder{queryBuilder: cloned, table: cloned.table, fields: fields}
}

func (q *QueryBuilder) Update(values Clauses) *UpdateQueryBuilder {
	cloned := q.clone()
	return &UpdateQueryBuilder{queryBuilder: cloned, table: cloned.table, values: values}
}

func (q *QueryBuilder) Insert(values Clauses) *InsertQueryBuilder {
	cloned := q.clone()
	return &InsertQueryBuilder{queryBuilder: cloned, table: cloned.table, values: values}
}

func (q *QueryBuilder) Delete() *DeleteQueryBuilder {
	cloned := q.clone()
	return &DeleteQueryBuilder{queryBuilder: cloned, table: cloned.table}
}

func (q *QueryBuilder) With(query CommonQueryBuilder, name string) *QueryBuilder {
	clonedQueryBuilder := *q
	if clonedQueryBuilder.commonTableExpressions == nil {
		clonedQueryBuilder.commonTableExpressions = make([]CommonQuery, 0)
	}

	clonedQueryBuilder.commonTableExpressions = append(q.commonTableExpressions, CommonQuery{Builder: query, Table: name})
	return &clonedQueryBuilder
}

func (q *QueryBuilder) addCondition(column string, value interface{}, comparer string, conditions *Clauses) {
	if *conditions == nil {
		*conditions = make(Clauses, 0)
	}
	key := fmt.Sprintf("%s:%s", column, comparer)
	*conditions = append(*conditions, Clause{ColumnName: key, Value: value})
}

func (q *QueryBuilder) buildColumnUpdateStatement(values Clauses) string {
	stmt := ""
	i := 0

	if len(values) != 0 {
		stmt += " SET"

		for _, column := range values {
			if i > 0 {
				stmt += ","
			}
			stmt += fmt.Sprintf(" %s = $%d", column.ColumnName, i+q.preparedVariableOffset+1)
			i++
		}
	}
	q.preparedVariableOffset = q.preparedVariableOffset + i
	return stmt
}

func (q *QueryBuilder) buildValuesStatement(values Clauses) string {
	stmt := ""
	i := 0

	if len(values) != 0 {
		for i = range values {
			if i > 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("$%d", i+q.preparedVariableOffset+1)
		}
	}

	q.preparedVariableOffset += i + 1
	return stmt
}

func (q *QueryBuilder) buildConditionalStatement(conditions Clauses) string {
	stmt := ""

	if len(conditions) != 0 {
		stmt += " WHERE"
		preparedStatementCount := 0
		for conditionIndex, clause := range conditions {
			if conditionIndex > 0 {
				stmt += " AND"
			}
			columnAndComparator := strings.Split(clause.ColumnName, ":")
			column, comparer := columnAndComparator[0], columnAndComparator[1]
			if comparer == "IS NULL" || comparer == "IS NOT NULL" {
				stmt += fmt.Sprintf(" %s %s", column, comparer)
			} else {
				stmt += fmt.Sprintf(" %s %s $%d", column, comparer, preparedStatementCount+q.preparedVariableOffset+1)
				preparedStatementCount++
			}
		}

		q.preparedVariableOffset += preparedStatementCount
	}

	return stmt
}

func (q *QueryBuilder) buildParameters(parameters Clauses) []interface{} {
	values := make([]interface{}, 0, len(parameters))
	for _, clause := range parameters {
		if clause.Value != "" {
			values = append(values, clause.Value)
		}
	}
	return values
}

func (q *QueryBuilder) buildReturnedColumns(fields []string) string {
	stmt := ""

	if len(fields) != 0 {
		for i, field := range fields {
			if i > 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("%s", field)
		}
	}

	return stmt
}

func (q *QueryBuilder) buildCommonTableExpressionParameters() []interface{} {
	values := make([]interface{}, 0)

	if len(q.commonTableExpressions) > 0 {
		for _, cte := range q.commonTableExpressions {
			values = append(values, cte.Builder.buildPreparedStatementValues()...)
		}
	}

	return values
}

func (q *QueryBuilder) buildCommonTableExpressions() (string, error) {
	stmt := ""

	if len(q.commonTableExpressions) > 0 {
		for _, cte := range q.commonTableExpressions {
			query, err := cte.Builder.buildQuery()
			if err != nil {
				return "", err
			}
			stmt += fmt.Sprintf("WITH %s AS (%s) ", cte.Table, *query)
		}
	}

	return stmt, nil
}
