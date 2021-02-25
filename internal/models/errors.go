package models

import "errors"

var (
	// ErrDuplicateEmail error.
	ErrDuplicateEmail = errors.New("models: duplicate email")

	// ErrDuplicateName error.
	ErrDuplicateName = errors.New("models: duplicate name")

	// ErrDuplicateUsername error.
	ErrDuplicateUsername = errors.New("models: duplicate username")

	// ErrInvalidCredentials error.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrInvalidValue error.
	ErrInvalidValue = errors.New("models: invalid value")

	// ErrNoRecord error.
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrTimestamp error.
	ErrTimestamp = errors.New("models: invalid timestamp type")
)
