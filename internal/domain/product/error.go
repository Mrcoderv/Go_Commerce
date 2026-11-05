package product

import "errors"

var (
	ErrNotFound          = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrConflict          = errors.New("product cannot be deleted: has existing orders")
)
