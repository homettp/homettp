package models

import "errors"

var (
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateName      = errors.New("models: duplicate name")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrInvalidValue       = errors.New("models: invalid value")
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrTimestamp          = errors.New("models: invalid timestamp type")
)
