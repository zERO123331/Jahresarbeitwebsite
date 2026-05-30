package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Update struct {
	Title       string
	UserID      int
	Body        string
	ID          int
	Created     time.Time
	LastUpdated time.Time
	Version     int
}

type UpdateModel struct {
	DB *sql.DB
}

func (m *UpdateModel) Insert(userID int, title, body string) (int, error) {
	stmt := `INSERT INTO updates (user_id, title, body) VALUES ($1, $2, $3) RETURNING id`
	args := []any{userID, title, body}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var ID int
	err := m.DB.QueryRowContext(ctx, stmt, args).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (m *UpdateModel) GetByID(id int) (*Update, error) {
	update := &Update{}
	stmt := `SELECT id, user_id, title, body, created_at, updated_at, version FROM updates WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(update.ID, update.UserID, update.Title, update.Body, update.Created, update.LastUpdated, update.Version)
	if err != nil {
		return nil, err
	}
	return update, nil
}

func (m *UpdateModel) GetAll(title string, filters Filters) ([]*Update, error) {
	stmt := fmt.Sprintf(`
SELECT id, user_id, title, body, created_at, updated_at, version 
FROM updates 
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, stmt, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	updates := []*Update{}
	for rows.Next() {
		update := &Update{}
		err := rows.Scan(&update.ID, &update.UserID, &update.Title, &update.Body, &update.Created, &update.LastUpdated, &update.Version)
		if err != nil {
			return nil, err
		}
		updates = append(updates, update)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return updates, nil
}

func (m *UpdateModel) Update(id int, title, body string, version int) error {
	stmt := `UPDATE updates SET title = $1, body = $2, version = version + 1, updated_at = NOW() WHERE id = $3 AND version = $4 RETURNING version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var updatedVersion int

	err := m.DB.QueryRowContext(ctx, stmt, title, body, id, version).Scan(&updatedVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}

	return nil
}

func (m *UpdateModel) Delete(id int) error {
	stmt := `DELETE FROM updates WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}
