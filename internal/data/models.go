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
	Tags      TagModel
	TagItems  TagItemModel
	ItemNotes ItemNoteModel
	Assets    AssetModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Roles:     RoleModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		Companies: CompanyModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		Notes:     NoteModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		Projects:  ProjectModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		Tags:      TagModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		TagItems:  TagItemModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		ItemNotes: ItemNoteModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
		Assets:    AssetModel{DB: db, Query: &querybuilder.QueryBuilder{DB: db}},
	}

}
