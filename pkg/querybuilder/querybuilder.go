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

type QueryBuilder struct {
	DB                     *sql.DB
	table                  string
	preparedVariableOffset int
	commonTableExpressions []CommonQueryBuilder
}

type CommonQueryBuilder interface {
	buildQuery() (*string, error)
	buildPreparedStatementValues() []interface{}
}

func (q *QueryBuilder) SetBaseTable(table string) *QueryBuilder {
	q.table = table
	return q
}

func (q QueryBuilder) Select(fields ...string) SelectQueryBuilder {
	return SelectQueryBuilder{queryBuilder: q, table: q.table, fields: fields}
}

func (q QueryBuilder) Update(values Clauses) UpdateQueryBuilder {
	return UpdateQueryBuilder{queryBuilder: q, table: q.table, values: values}
}

func (q QueryBuilder) Insert(values Clauses) InsertQueryBuilder {
	return InsertQueryBuilder{queryBuilder: q, table: q.table, values: values}
}

func (q QueryBuilder) Delete() DeleteQueryBuilder {
	return DeleteQueryBuilder{queryBuilder: q, table: q.table}
}

func (q QueryBuilder) With(query CommonQueryBuilder, name string) QueryBuilder {
	if q.commonTableExpressions == nil {
		q.commonTableExpressions = make([]CommonQueryBuilder, 0)
	}

	q.commonTableExpressions = append(q.commonTableExpressions, query)
	return q
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
