package querybuilder

import (
	"reflect"
	"testing"
)

func TestSelectQueryBuilder_SetBaseTable(t *testing.T) {
	qb := QueryBuilder{}
	query, err := qb.SetBaseTable("users").Select("id", "name").buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	if *query != "SELECT id, name FROM users" {
		t.Fatalf("Expected query was not generated got: %s", *query)
	}
}

func TestSelectQueryBuilder_From(t *testing.T) {
	qb := QueryBuilder{}
	selectQB := qb.Select("id", "name").From("users")

	if selectQB.table != "users" {
		t.Errorf("Expected table to be 'users', got %s", selectQB.table)
	}
	query, err := selectQB.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	if *query != "SELECT id, name FROM users" {
		t.Fatalf("Expected query was not generated got: %s", *query)
	}
}

func TestSelectQueryBuilder_NoTableSpecified(t *testing.T) {
	qb := QueryBuilder{}
	selectQB := qb.Select("id", "name")

	if selectQB.table != "" {
		t.Errorf("Expected table to be empty, got %s", selectQB.table)
	}
	query, err := selectQB.buildQuery()
	if query != nil {
		t.Fatalf("Expected nil query, got %s", *query)
	}

	if err == nil {
		t.Fatalf("Expected table error, got nil")
	}
}

func TestSelectQueryBuilder_NoFieldsSpecified(t *testing.T) {
	qb := QueryBuilder{}
	selectQB := qb.Select().From("users")

	if selectQB.table != "users" {
		t.Errorf("Expected table to be 'users', got %s", selectQB.table)
	}
	query, err := selectQB.buildQuery()
	if query != nil {
		t.Fatalf("Expected nil query, got %s", *query)
	}

	if err == nil {
		t.Fatalf("Expected fields error, got nil")
	}
}

func TestSelectQueryBuilder_LeftJoin(t *testing.T) {
	qb := QueryBuilder{}
	selectQB := qb.Select("id", "name").From("users").LeftJoin("orders", "id", "user_id")

	if selectQB.leftJoinTable != "orders" {
		t.Errorf("Expected left join table to be 'orders', got %s", selectQB.leftJoinTable)
	}

	if selectQB.leftJoinOwnKey != "id" || selectQB.leftJoinForeignKey != "user_id" {
		t.Errorf("Expected join keys to be 'id' and 'user_id', got %s and %s", selectQB.leftJoinOwnKey, selectQB.leftJoinForeignKey)
	}

	query, err := selectQB.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	if *query != "SELECT id, name FROM users LEFT JOIN orders ON users.id = orders.user_id" {
		t.Fatalf("Expected query was not generated got: %s", *query)
	}
}

func TestSelectQueryBuilder_OrderBy(t *testing.T) {
	qb := QueryBuilder{}
	selectQB := qb.Select("id", "name").From("users").OrderBy("id", "DESC")

	if selectQB.sortColumn != "id" || selectQB.sortDirection != "DESC" {
		t.Errorf("Expected order by 'id DESC', got %s %s", selectQB.sortColumn, selectQB.sortDirection)
	}

	query, err := selectQB.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	if *query != "SELECT id, name FROM users ORDER BY id DESC" {
		t.Fatalf("Expected query was not generated got: %s", *query)
	}
}

func TestSelectQueryBuilder_WhereEqual(t *testing.T) {
	qb := QueryBuilder{}
	selectQB := qb.Select("id", "name").From("users").WhereEqual("age", 30)

	expectedConditions := ClauseMap{"age:=": 30}
	if !reflect.DeepEqual(selectQB.conditions, expectedConditions) {
		t.Errorf("Expected conditions to be %v, got %v", expectedConditions, selectQB.conditions)
	}

	query, err := selectQB.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	if *query != "SELECT id, name FROM users WHERE age = $1" {
		t.Fatalf("Expected query was not generated got: %s", *query)
	}
}

func TestSelectQueryBuilder_WhereEqual_Null(t *testing.T) {
	qb := QueryBuilder{}
	selectQB := qb.Select("id", "name").From("users").WhereEqual("age", 30).WhereEqual("deleted_at", nil).WhereEqual("category", "admin")

	expectedConditions := ClauseMap{"age:=": 30, "deleted_at:IS NULL": nil, "category:=": "admin"}
	if !reflect.DeepEqual(selectQB.conditions, expectedConditions) {
		t.Errorf("Expected conditions to be %v, got %v", expectedConditions, selectQB.conditions)
	}

	query, err := selectQB.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	if *query != "SELECT id, name FROM users WHERE age = $1 AND deleted_at IS NULL AND category = $2" {
		t.Fatalf("Expected query was not generated got: %s", *query)
	}
}
