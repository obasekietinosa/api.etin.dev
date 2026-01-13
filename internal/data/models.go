package data

import (
	"database/sql"
	"log"

	"api.etin.dev/pkg/querybuilder"
)

type Models struct {
	Roles     RoleModel
	Companies CompanyModel
	Notes     NoteModel
	Projects  ProjectModel
	Tags      TagModel
	TagItems  TagItemModel
	ItemNotes ItemNoteModel
	Assets    AssetModel
}

func NewModels(db *sql.DB, logger *log.Logger) Models {
	return Models{
		Roles:     RoleModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}, Logger: logger},
		Companies: CompanyModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}, Logger: logger},
		Notes:     NoteModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}, Logger: logger},
		Projects:  ProjectModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}, Logger: logger},
		Tags:      TagModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}, Logger: logger},
		TagItems:  TagItemModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}, Logger: logger},
		ItemNotes: ItemNoteModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}, Logger: logger},
		Assets:    AssetModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}, Logger: logger},
	}
}
