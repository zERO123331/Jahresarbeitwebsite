package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
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
	stmt := `INSERT INTO shopentry (title, description, price, quantity, image_urls, categories, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	args := []any{title, description, price, quantity, pq.Array(imageURLs), pq.Array(categories), userID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ShopModel) Buy(id, quantity int) (int, error) {
	stmt := `UPDATE shopentry SET quantity = quantity - $1 WHERE id = $2 RETURNING quantity`
	args := []any{quantity, id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var quantityLeft int
	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&quantityLeft)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrNoRecord
		default:
			return 0, err
		}
	}
	return quantityLeft, nil

}

func (m *ShopModel) Restock(id, quantity int) (int, error) {
	stmt := `UPDATE shopentry SET quantity = quantity + $1 WHERE id = $2 RETURNING quantity`
	args := []any{quantity, id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var quantityLeft int
	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&quantityLeft)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			quantityLeft = 0
		default:
			return 0, err
		}
	}
	return quantityLeft, nil
}

func (m *ShopModel) GetByID(id int) (*ShopEntry, error) {
	stmt := `SELECT id, created_at, title, description, price, quantity, image_urls, categories, user_id FROM shopentry WHERE id = $1`
	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	entry := &ShopEntry{}
	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&entry.ID, &entry.CreatedAt, &entry.Title, &entry.Description, &entry.Price, &entry.Quantity, &entry.ImageURLS, &entry.Categories, &entry.UserID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}
	return entry, nil
}

func (m *ShopModel) GetAll(title string, categories []string, filters Filters) ([]*ShopEntry, error) {
	// TODO: add pagination
	query := fmt.Sprintf(`
SELECT id, created_at, title, description, price, quantity, image_urls, categories, user_id 
FROM shopentry 
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
  AND (categories @> $2 OR $2 = '{}') 
ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(categories))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []*ShopEntry
	for rows.Next() {
		var entry ShopEntry
		err := rows.Scan(&entry.ID, &entry.CreatedAt, &entry.Title, &entry.Description, &entry.Price, &entry.Quantity, pq.Array(&entry.ImageURLS), pq.Array(&entry.Categories), &entry.UserID)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}
