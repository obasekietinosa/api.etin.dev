package querybuilder

import "testing"

func TestSelectQueryBuilder_PreparedValuesForNullCondition(t *testing.T) {
	qb := QueryBuilder{}
	selectQB := qb.SetBaseTable("projects").Select("id", "title").WhereEqual("deletedAt", nil).OrderBy("startDate", "desc")

	query, err := selectQB.buildQuery()
	if err != nil {
		t.Fatalf("unexpected error building query: %s", err)
	}

	values := selectQB.buildPreparedStatementValues()
	if len(values) != 0 {
		t.Fatalf("expected no prepared values, got %d for query %s", len(values), *query)
	}
}
