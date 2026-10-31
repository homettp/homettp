package web

import "errors"

var (
	// ErrInvalidCall error.
	ErrInvalidCall = errors.New("web: invalid call")

	// ErrMissingParameter error.
	ErrMissingParameter = errors.New("web: parameter not found")

	// ErrQueueClosed error.
	ErrQueueClosed = errors.New("web: queue is closed")

	// ErrQueueFull error.
	ErrQueueFull = errors.New("web: queue is full")
)
