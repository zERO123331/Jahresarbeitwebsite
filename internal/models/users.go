package models

import (
	"Jahresarbeitwebsite/internal/validator"
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	CreatedAt    time.Time
	Name         string
	Email        string
	Password     password
	Password2    string
	PasswordHash string
	Activated    bool
	Version      int
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.hash = hash
	p.plaintext = &plaintextPassword
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(user *User) (int, error) {
	stmt := `INSERT INTO users (name, email, password_hash, activated) VALUES ($1, $2, $3, false) RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password.hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return 0, ErrDuplicateEmail
		default:
			return 0, err
		}
	}

	return user.ID, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
SELECT id, created_at, name, email, password_hash, activated, version
FROM user
WHERE email = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var user User
	query := `
SELECT id, created_at, name, email, password_hash, activated, version
FROM users
WHERE email = $1`

	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.PasswordHash, &user.Activated, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrNoRecord
		default:
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return 0, ErrInvalidCredentials
	}
	if !user.Activated {
		return 0, ErrUserNotActivated
	}

	return user.ID, nil
}

func (m UserModel) Update(user *User) error {
	query := `
UPDATE users
SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m *UserModel) Activate(id int) error {
	return nil
}

func (m *UserModel) Exists(id int) (bool, error) {

	return false, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "Email must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "Email must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "Password must be provided")
	v.Check(len(password) >= 8, "password", "Password must be at least 8 characters long")
	v.Check(len(password) <= 72, "password", "Password must be at most 255 characters long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "Name must be provided")
	v.Check(len(user.Name) <= 255, "name", "Name must be at most 255 characters long")

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
