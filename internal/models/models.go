package models

import (
	"database/sql"
)

type Models struct {
	Shop        ShopModel
	User        UserModel
	Update      UpdateModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Shop:        ShopModel{DB: db},
		User:        UserModel{DB: db},
		Update:      UpdateModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}
