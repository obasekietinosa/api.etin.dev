package querybuilder

import (
	"reflect"
	"testing"
)

func TestUpdateQueryBuilder_Fails_Without_Where(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{},
		Clause{ColumnName: "name", Value: "John"},
		Clause{ColumnName: "age", Value: 31},
	)

	updateQB := qb.Update(values)

	_, err := updateQB.buildQuery()

	if err == nil {
		t.Fatalf("Expected error for missing conditions, got nil")
	}
}

func TestUpdateQueryBuilder_WhereEqual(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{},
		Clause{ColumnName: "name", Value: "John"},
		Clause{ColumnName: "age", Value: 31},
	)

	updateQB := qb.SetBaseTable("users").Update(values).WhereEqual("id", 1)

	expectedConditions := append(Clauses{}, Clause{ColumnName: "id:=", Value: 1})
	if !reflect.DeepEqual(updateQB.conditions, expectedConditions) {
		t.Errorf("Expected conditions to be %v, got %v", expectedConditions, updateQB.conditions)
	}

	query, err := updateQB.buildQuery()

	expectedQuery := "UPDATE users SET name = $1, age = $2 WHERE id = $3"
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if *query != expectedQuery {

		t.Errorf("Expected query to be '%s', got '%s'", expectedQuery, *query)

	}

}

func TestUpdateQueryBuilder_NoTable(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{}, Clause{ColumnName: "name", Value: "John"})

	updateQB := qb.Update(values).WhereEqual("id", 1)

	_, err := updateQB.buildQuery()

	if err == nil {
		t.Error("Expected error for missing table, got nil")
	}
}

func TestUpdateQueryBuilder_Returning(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{},
		Clause{ColumnName: "name", Value: "John"},
		Clause{ColumnName: "age", Value: 31},
	)

	updateQB := qb.SetBaseTable("users").Update(values).Returning("id", "updated_at")

	query, err := updateQB.buildQuery()

	expectedQuery := "UPDATE users SET name = $1, age = $2 RETURNING id, updated_at"
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if *query != expectedQuery {
		t.Errorf("Expected query to be '%s', got '%s'", expectedQuery, *query)
	}
}
