package config

import "errors"

var (
	// ErrInvalidCommandTimeout error.
	ErrInvalidCommandTimeout = errors.New("config: invalid command timeout")

	// ErrInvalidCommandWorkerCount error.
	ErrInvalidCommandWorkerCount = errors.New("config: invalid command worker count")

	// ErrInvalidCommandHistoryLimit error.
	ErrInvalidCommandHistoryLimit = errors.New("config: invalid command history limit")
)
