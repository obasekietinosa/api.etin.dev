package data

import (
	"database/sql"

	"api.etin.dev/pkg/querybuilder"
)

type Models struct {
	Roles     RoleModel
	Companies CompanyModel
	Notes     NoteModel
	Projects  ProjectModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Roles:     RoleModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		Companies: CompanyModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		Notes:     NoteModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		Projects:  ProjectModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
	}

}
