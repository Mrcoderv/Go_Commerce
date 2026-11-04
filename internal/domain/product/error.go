package product

import "errors"

var (
	ErrNotFound          = errors.New(" not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
)
