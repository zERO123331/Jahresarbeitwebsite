package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrUserNotActivated   = errors.New("models: user not activated")
	ErrEditConflict       = errors.New("models: edit conflict")
	ErrProductNotFound    = errors.New("product not found")
	ErrNotImplemented     = errors.New("not implemented")
)
