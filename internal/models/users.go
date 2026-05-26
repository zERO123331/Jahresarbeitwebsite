package models

import (
	"database/sql"
	"fmt"
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}
	user.PasswordHash = string(passwordHash)

	stmt := "INSERT INTO users (name, email, password_hash, activated) VALUES ($1, $2, $3, false) RETURNING id, created_at, version"

	err = m.DB.QueryRow(stmt, user.Name, user.Email, passwordHash).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var user User
	query := `
SELECT id, created_at, name, email, password_hash, activated, version
FROM users
WHERE email = $1`

	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.PasswordHash, &user.Activated, &user.Version)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return 0, err
	}
	if !user.Activated {
		return 0, fmt.Errorf("user not activated")
	}

	return user.ID, nil
}

func (m *UserModel) Activate(id int) error {
	return nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func ValidateEmail(v *validator.Validator)
