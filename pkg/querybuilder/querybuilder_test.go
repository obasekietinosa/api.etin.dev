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
	values := append(Clauses{}, Clause{ColumnName: "name", Value: "John Doe"}, Clause{ColumnName: "age", Value: 30})
	qb := QueryBuilder{}
	insertQB := qb.SetBaseTable("users").Insert(values).Returning("id", "created_at")

	if insertQB.table != "users" {
		t.Errorf("Expected table to be 'users', got %s", insertQB.table)
	}

	expectedValues := append(Clauses{}, Clause{ColumnName: "name", Value: "John Doe"}, Clause{ColumnName: "age", Value: 30})
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
	conditions := make(Clauses, 0)
	qb.addCondition("age", 30, "=", &conditions)

	expectedConditions := append(Clauses{}, Clause{ColumnName: "age:=", Value: 30})
	if !reflect.DeepEqual(conditions, expectedConditions) {
		t.Errorf("Expected conditions to be %v, got %v", expectedConditions, conditions)
	}
}

func TestBuildValuesStatement(t *testing.T) {
	qb := QueryBuilder{}
	qb.preparedVariableOffset = 0
	values := make(Clauses, 0)
	values = append(values,
		Clause{ColumnName: "age:>", Value: 30},
		Clause{ColumnName: "name:=", Value: "John Doe"})

	stmt := qb.buildValuesStatement(values)

	expectedStmt := "$1, $2"
	if stmt != expectedStmt {
		t.Errorf("Expected statement to be '%s', got '%s'", expectedStmt, stmt)
	}
}

func TestBuildConditionalStatement(t *testing.T) {
	qb := QueryBuilder{}
	conditions := make(Clauses, 0)
	conditions = append(conditions,
		Clause{ColumnName: "age:>", Value: 30},
		Clause{ColumnName: "name:=", Value: "John Doe"})

	stmt := qb.buildConditionalStatement(conditions)

	expectedStmt := " WHERE age > $1 AND name = $2"
	if stmt != expectedStmt {
		t.Errorf("Expected statement to be '%s', got '%s'", expectedStmt, stmt)
	}
}

func TestBuildParameters(t *testing.T) {
	qb := QueryBuilder{}
	parameters := make(Clauses, 0)
	parameters = append(parameters,
		Clause{ColumnName: "age:>", Value: 30},
		Clause{ColumnName: "name:=", Value: "John Doe"})

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

func Test_CommonTableExpression_WithUpdate_AndSelect(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{}, Clause{ColumnName: "name", Value: "John"})

	selectWithCte := qb.With(
		qb.SetBaseTable("users").Update(values).WhereEqual("id", 3).Returning("*"), "updated_user").
		Select(
			"updated_user.updatedAt AS updatedAt",
			"companies.name AS companyName",
			"companies.icon AS companyIcon").From("updated_user").LeftJoin("companies", "company_id", "id")

	query, err := selectWithCte.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	expected := "WITH updated_user AS (UPDATE users SET name = $1 WHERE id = $2 RETURNING *) SELECT updated_user.updatedAt AS updatedAt, companies.name AS companyName, companies.icon AS companyIcon FROM updated_user LEFT JOIN companies ON updated_user.company_id = companies.id"
	if *query != expected {
		t.Fatalf("Expected query was not generated\nexpected: %s \ngot: %s", expected, *query)
	}
}

func Test_CommonTableExpression_WithInsert_AndSelect(t *testing.T) {
	qb := QueryBuilder{}
	values := append(Clauses{}, Clause{ColumnName: "name", Value: "John"})

	selectWithCte := qb.With(
		qb.SetBaseTable("users").Insert(values).Returning("*"), "inserted_user").
		Select(
			"inserted_user.updatedAt AS updatedAt",
			"companies.name AS companyName",
			"companies.icon AS companyIcon").From("updated_user").LeftJoin("companies", "company_id", "id")

	query, err := selectWithCte.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	expected := "WITH inserted_user AS (INSERT INTO users (name) VALUES ($1) RETURNING *) SELECT inserted_user.updatedAt AS updatedAt, companies.name AS companyName, companies.icon AS companyIcon FROM updated_user LEFT JOIN companies ON updated_user.company_id = companies.id"
	if *query != expected {
		t.Fatalf("Expected query was not generated\nexpected: %s \ngot: %s", expected, *query)
	}
}

func Test_CommonTableExpression_WithSelect_AndSelect(t *testing.T) {
	qb := QueryBuilder{}

	selectWithCte := qb.With(
		qb.SetBaseTable("users").Select("id", "name").WhereEqual("age", 40), "selected_user").
		Select(
			"selected_user.updatedAt AS updatedAt",
			"companies.name AS companyName",
			"companies.icon AS companyIcon").From("selected_user").LeftJoin("companies", "company_id", "id")

	query, err := selectWithCte.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error when building select query, got %s", err)
	}

	expected := "WITH selected_user AS (SELECT id, name FROM users WHERE age = $1) SELECT selected_user.updatedAt AS updatedAt, companies.name AS companyName, companies.icon AS companyIcon FROM selected_user LEFT JOIN companies ON selected_user.company_id = companies.id"
	if *query != expected {
		t.Fatalf("Expected query was not generated\nexpected: \n%s \ngot: \n%s", expected, *query)
	}
}

func TestPreparedVariableOffset_IsResetPerQuery(t *testing.T) {
	qb := QueryBuilder{}
	values := Clauses{
		{ColumnName: "name", Value: "Alice"},
	}

	// First query
	insertQuery1 := qb.SetBaseTable("users").Insert(values).Returning("id")
	query1, err := insertQuery1.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error building first insert query: %s", err)
	}
	expected1 := "INSERT INTO users (name) VALUES ($1) RETURNING id"
	if *query1 != expected1 {
		t.Errorf("First query not as expected. Expected: %s, Got: %s", expected1, *query1)
	}

	// Second query
	updateQuery := qb.SetBaseTable("users").Update(values).WhereEqual("id", 3).Returning("id")
	query2, err := updateQuery.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error building update query: %s", err)
	}
	expected2 := "UPDATE users SET name = $1 WHERE id = $2 RETURNING id"
	if *query2 != expected2 {
		t.Errorf("Second query not as expected. Expected: %s, Got: %s", expected2, *query2)
	}

	// Third query
	insertQuery2 := qb.SetBaseTable("users").Insert(values).Returning("id")
	query3, err := insertQuery2.buildQuery()
	if err != nil {
		t.Fatalf("Unexpected error building second insert query: %s", err)
	}
	expected3 := "INSERT INTO users (name) VALUES ($1) RETURNING id"
	if *query3 != expected3 {
		t.Errorf("Third query not as expected. Expected: %s, Got: %s", expected3, *query3)
	}
}
