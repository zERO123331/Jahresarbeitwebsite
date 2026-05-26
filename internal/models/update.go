package models

import (
	"database/sql"
	"time"
)

type Update struct {
	Title       string
	Author      User
	Body        string
	ID          int
	Created     time.Time
	LastUpdated time.Time
}

type UpdateModel struct {
	DB *sql.DB
}

func (m *UpdateModel) Insert(author User, title, body string) error {
	return nil
}

func (m *UpdateModel) GetByID(id int) (*Update, error) {
	return nil, nil
}

func (m *UpdateModel) GetAll() ([]*Update, error) {
	return nil, nil
}

func (m *UpdateModel) Update(id int, author User, title, body string) error {
	return nil
}

func (m *UpdateModel) Delete(id int) error {
	return nil
}

func (m *UpdateModel) GetLatest(count int) (*[]Update, error) {
	return nil, nil
}
