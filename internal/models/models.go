package models

import (
	"database/sql"
	"errors"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type Models struct {
	Shop ShopModel
	User UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Shop: ShopModel{DB: db},
		User: UserModel{DB: db},
	}
}
