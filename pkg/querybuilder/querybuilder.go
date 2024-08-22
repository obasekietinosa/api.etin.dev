package querybuilder

import (
	"database/sql"
	"fmt"
	"strings"
)

type QueryBuilder struct {
	DB     *sql.DB
	fields []string
	table  string
}

func (q QueryBuilder) Select(fields ...string) SelectQueryBuilder {
	return SelectQueryBuilder{queryBuilder: q, table: q.table, fields: fields}
}

func (q *QueryBuilder) addCondition(column string, value interface{}, comparer string, conditions *map[string]interface{}) {
	if *conditions == nil {
		*conditions = (make(map[string]interface{}))
	}
	key := fmt.Sprintf("%s:%s", column, comparer)
	(*conditions)[key] = value
}

func (q QueryBuilder) buildConditionalStatement(conditions map[string]interface{}) string {
	keys := make([]string, 0, len(conditions))
	for k := range conditions {
		keys = append(keys, k)
	}

	stmt := ""

	if len(keys) != 0 {
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

func (q QueryBuilder) buildParameters(parameters map[string]interface{}) []interface{} {
	values := make([]interface{}, 0, len(parameters))
	for _, v := range parameters {
		values = append(values, v)
	}
	return values
}
