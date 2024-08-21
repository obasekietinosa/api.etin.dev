package data

import "database/sql"

type Models struct {
	Roles     RoleModel
	Companies CompanyModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Roles:     RoleModel{DB: db},
		Companies: CompanyModel{DB: db},
	}

}
