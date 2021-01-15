package models

import "errors"

var (
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrNoRecord           = errors.New("models: no matching record found")
)
