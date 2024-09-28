package querybuilder

import (
	"database/sql"
	"fmt"
	"strings"
)

type ClauseMap map[string]interface{}

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

func (q QueryBuilder) Update(values ClauseMap) UpdateQueryBuilder {
	return UpdateQueryBuilder{queryBuilder: q, table: q.table, values: values}
}

func (q QueryBuilder) Insert(values ClauseMap) InsertQueryBuilder {
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

func (q *QueryBuilder) addCondition(column string, value interface{}, comparer string, conditions *ClauseMap) {
	if *conditions == nil {
		*conditions = (make(map[string]interface{}))
	}
	key := fmt.Sprintf("%s:%s", column, comparer)
	(*conditions)[key] = value
}

func (q *QueryBuilder) buildColumnUpdateStatement(values ClauseMap) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}

	stmt := ""
	i, column := 0, ""

	if len(keys) != 0 {
		stmt += " SET"

		for i, column = range keys {
			if i > 0 {
				stmt += ","
			}
			stmt += fmt.Sprintf(" %s = $%d", column, i+q.preparedVariableOffset+1)
		}
	}
	q.preparedVariableOffset = q.preparedVariableOffset + i + 1
	// return fmt.Sprintf("number: %d", q.preparedVariableOffset)
	return stmt
}

func (q *QueryBuilder) buildValuesStatement(values ClauseMap) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}

	stmt := ""
	i := 0

	if len(keys) != 0 {
		for i = range keys {
			if i > 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("$%d", i+q.preparedVariableOffset+1)
		}
	}

	q.preparedVariableOffset += i + 1
	return stmt
}

func (q *QueryBuilder) buildConditionalStatement(conditions ClauseMap) string {
	keys := make([]string, 0, len(conditions))
	for k := range conditions {
		keys = append(keys, k)
	}

	stmt := ""

	if len(keys) != 0 {
		stmt += " WHERE"
		preparedStatementCount := 0
		for conditionIndex, v := range keys {
			if conditionIndex > 0 {
				stmt += " AND"
			}
			columnAndComparator := strings.Split(v, ":")
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

func (q *QueryBuilder) buildParameters(parameters ClauseMap) []interface{} {
	values := make([]interface{}, 0, len(parameters))
	for _, v := range parameters {
		if v != nil {
			values = append(values, v)
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
