package data

import (
	"database/sql"

	"api.etin.dev/pkg/querybuilder"
)

type Models struct {
	Roles     RoleModel
	Companies CompanyModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Roles:     RoleModel{DB: db},
		Companies: CompanyModel{DB: db, QB: &querybuilder.QueryBuilder{DB: db}},
	}

}
