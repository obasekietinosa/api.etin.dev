package querybuilder

import (
	"reflect"
	"testing"
)

func TestSetBaseTable(t *testing.T) {
	qb := QueryBuilder{}
	qb.SetBaseTable("users")

	if qb.table != "users" {
		t.Errorf("Expected table to be 'users', got %s", qb.table)
	}
}

func TestInsertQueryBuilder(t *testing.T) {
	qb := QueryBuilder{}
	insertQB := qb.SetBaseTable("users").Insert(ClauseMap{"name": "John Doe", "age": 30}).Returning("id", "created_at")

	if insertQB.table != "users" {
		t.Errorf("Expected table to be 'users', got %s", insertQB.table)
	}

	expectedValues := ClauseMap{"name": "John Doe", "age": 30}
	if !reflect.DeepEqual(insertQB.values, expectedValues) {
		t.Errorf("Expected values to be %v, got %v", expectedValues, insertQB.values)
	}

	query, err := insertQB.buildQuery()
	if err != nil {
		t.Errorf("Unexpected error when building insert query, got %s", err)
		return
	}

	if *query != "INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, created_at" {
		t.Errorf("Expected query was not generated got: %s", *query)
	}
}

func TestDeleteQueryBuilder(t *testing.T) {
	qb := QueryBuilder{}
	deleteQB := qb.SetBaseTable("users").Delete()

	if deleteQB.table != "users" {
		t.Errorf("Expected table to be 'users', got %s", deleteQB.table)
	}
}

func TestAddCondition(t *testing.T) {
	qb := QueryBuilder{}
	conditions := make(ClauseMap)
	qb.addCondition("age", 30, "=", &conditions)

	expectedConditions := ClauseMap{"age:=": 30}
	if !reflect.DeepEqual(conditions, expectedConditions) {
		t.Errorf("Expected conditions to be %v, got %v", expectedConditions, conditions)
	}
}

func TestBuildValuesStatement(t *testing.T) {
	qb := QueryBuilder{}
	qb.preparedVariableOffset = 0
	values := ClauseMap{"name": "John Doe", "age": 30}

	stmt := qb.buildValuesStatement(values)

	expectedStmt := "$1, $2"
	if stmt != expectedStmt {
		t.Errorf("Expected statement to be '%s', got '%s'", expectedStmt, stmt)
	}
}

func TestBuildConditionalStatement(t *testing.T) {
	qb := QueryBuilder{}
	conditions := ClauseMap{"age:>": 30, "name:=": "John Doe"}

	stmt := qb.buildConditionalStatement(conditions)

	expectedStmt := " WHERE age > $1 AND name = $2"
	if stmt != expectedStmt {
		t.Errorf("Expected statement to be '%s', got '%s'", expectedStmt, stmt)
	}
}

func TestBuildParameters(t *testing.T) {
	qb := QueryBuilder{}
	parameters := ClauseMap{"age:>": 30, "name:=": "John Doe"}

	params := qb.buildParameters(parameters)

	expectedParams := []interface{}{30, "John Doe"}
	if !reflect.DeepEqual(params, expectedParams) {
		t.Errorf("Expected parameters to be %v, got %v", expectedParams, params)
	}
}

func TestBuildReturnedColumns(t *testing.T) {
	qb := QueryBuilder{}
	fields := []string{"id", "name", "email"}

	stmt := qb.buildReturnedColumns(fields)

	expectedStmt := "id, name, email"
	if stmt != expectedStmt {
		t.Errorf("Expected statement to be '%s', got '%s'", expectedStmt, stmt)
	}
}
