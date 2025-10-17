package querybuilder

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSelectQueryBuilder_QueryWhereNullDoesNotSendArgs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating sqlmock: %s", err)
	}
	defer db.Close()

	qb := QueryBuilder{DB: db}

	mock.ExpectQuery("SELECT id FROM projects WHERE deletedAt IS NULL").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	rows, err := qb.SetBaseTable("projects").Select("id").WhereEqual("deletedAt", nil).Query()
	if err != nil {
		t.Fatalf("unexpected error querying: %s", err)
	}
	rows.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unmet expectations: %s", err)
	}
}
