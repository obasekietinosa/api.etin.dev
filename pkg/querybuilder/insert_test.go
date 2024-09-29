package querybuilder

import (
	"testing"
)

func TestInsertQueryBuilder_buildColumnNameStatement(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{}, Clause{ColumnName: "name", Value: "John"}, Clause{ColumnName: "age", Value: 30})

	insertQB := InsertQueryBuilder{queryBuilder: &qb, values: values}

	columnStmt := insertQB.buildColumnNameStatement()

	expectedStmt := "name, age"
	if columnStmt != expectedStmt {
		t.Errorf("Expected column statement to be '%s', got '%s'", expectedStmt, columnStmt)
	}
}

func TestInsertQueryBuilder_buildQuery(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{}, Clause{ColumnName: "name", Value: "John"}, Clause{ColumnName: "age", Value: 30})

	insertQB := qb.SetBaseTable("users").Insert(values)

	query, err := insertQB.buildQuery()

	expectedQuery := "INSERT INTO users (name, age) VALUES ($1, $2)"
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if *query != expectedQuery {
		t.Errorf("Expected query to be '%s', got '%s'", expectedQuery, *query)
	}
}

func TestInsertQueryBuilder_NoTable(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{}, Clause{ColumnName: "name", Value: "John"})

	insertQB := qb.Insert(values)

	_, err := insertQB.buildQuery()

	if err == nil {
		t.Error("Expected error for missing table, got nil")
	}
}

func TestInsertQueryBuilder_Returning(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{}, Clause{ColumnName: "name", Value: "John"}, Clause{ColumnName: "age", Value: 30})

	insertQB := qb.SetBaseTable("users").Insert(values).Returning("id", "created_at")

	query, err := insertQB.buildQuery()

	expectedQuery := "INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, created_at"
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if *query != expectedQuery {
		t.Errorf("Expected query to be '%s', got '%s'", expectedQuery, *query)
	}
}
