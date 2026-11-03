package product

import "errors"

var (
	ErrNotFound           = errors.New(" not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInsufficientStock     = errors.New("insufficient stock")
	ErrInvalidProductID     = errors.New("invalid product id")
)
