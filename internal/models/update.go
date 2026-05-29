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

// TODO: implement

func (m *UpdateModel) Insert(author *User, title, body string) error {
	return ErrNotImplemented
}

func (m *UpdateModel) GetByID(id int) (*Update, error) {
	return nil, ErrNotImplemented
}

func (m *UpdateModel) GetAll(title string, filters Filters) ([]*Update, error) {
	return nil, ErrNotImplemented
}

func (m *UpdateModel) Update(id int, author *User, title, body string) error {
	return ErrNotImplemented
}

func (m *UpdateModel) Delete(id int) error {
	return ErrNotImplemented
}

func (m *UpdateModel) GetLatest(count int) ([]*Update, error) {
	return nil, ErrNotImplemented
}
