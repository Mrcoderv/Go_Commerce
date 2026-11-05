package order

import "errors"

// Business errors. Every layer can use them, and the HTTP adapter
// turns them into proper status codes (404, 409).
var (
	ErrNotFound          = errors.New("order not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidTransition = errors.New("invalid status transition")
)
