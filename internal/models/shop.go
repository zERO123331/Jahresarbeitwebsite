package models

import (
	"database/sql"
	"time"
)

type ShopEntry struct {
	ID          int
	CreatedAt   time.Time
	Title       string
	Description string
	Price       int
	Quantity    int
	ImageURLS   []string
	Categories  []string
	UserID      int
}

type ShopModel struct {
	DB *sql.DB
}

func (m *ShopModel) Insert(title, description string, price, quantity int, imageURLs, categories []string, userID int) error {
	return nil
}

func (m *ShopModel) Buy(id, quantity int) error {
	return nil
}

func (m *ShopModel) Restock(id, quantity int) error {
	return nil
}

func (m *ShopModel) GetByID(id int) (*ShopEntry, error) {
	return nil, nil
}

func (m *ShopModel) GetByTitle(title string, itemsPerPage, Page int) ([]*ShopEntry, error) {
	return nil, nil
}

func (m *ShopModel) SearchByCategories(categories []string, itemsPerPage, Page int) ([]*ShopEntry, error) {
	return nil, nil
}

func (m *ShopModel) GetAll() ([]*ShopEntry, error) {
	return nil, nil
}

func (m *ShopModel) Update(id int, title, description string, price, quantity int, imageURLs, categories []string) error {
	return nil
}

func (m *ShopModel) Delete(id int) error {
	return nil
}

func (m *ShopModel) GetByUserID(userID int) ([]*ShopEntry, error) {
	return nil, nil
}

func (m *ShopModel) GetLowStock(threshold int) ([]*ShopEntry, error) {
	return nil, nil
}
