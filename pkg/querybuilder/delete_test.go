package querybuilder

import (
	"reflect"
	"testing"
)

func TestDeleteQueryBuilder_Fails_Without_Where(t *testing.T) {
	qb := QueryBuilder{}

	DeleteQB := qb.Delete()

	_, err := DeleteQB.buildQuery()

	if err == nil {
		t.Fatalf("Expected error for missing conditions, got nil")
	}
}

func TestDeleteQueryBuilder_WhereEqual(t *testing.T) {
	qb := QueryBuilder{}

	DeleteQB := qb.SetBaseTable("users").Delete().WhereEqual("id", 1)

	expectedConditions := ClauseMap{"id:=": 1}
	if !reflect.DeepEqual(DeleteQB.conditions, expectedConditions) {
		t.Errorf("Expected conditions to be %v, got %v", expectedConditions, DeleteQB.conditions)
	}

	query, err := DeleteQB.buildQuery()

	expectedQuery := "DELETE FROM users WHERE id = $1"
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if *query != expectedQuery {
		t.Errorf("Expected query to be '%s', got '%s'", expectedQuery, *query)
	}
}

func TestDeleteQueryBuilder_NoTable(t *testing.T) {
	qb := QueryBuilder{}

	DeleteQB := qb.Delete().WhereEqual("id", 1)

	_, err := DeleteQB.buildQuery()

	if err == nil {
		t.Error("Expected error for missing table, got nil")
	}
}
