package models

import (
	"context"
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

func (m *ShopModel) Insert(title, description string, price, quantity int, imageURLs, categories []string, userID int) (int, error) {
	stmt := `INSERT INTO shop (title, description, price, quantity, image_urls, categories, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	args := []any{title, description, price, quantity, imageURLs, categories, userID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	err := m.DB.QueryRowContext(ctx, stmt, args).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ShopModel) Buy(id, quantity int, user *User) error {
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
