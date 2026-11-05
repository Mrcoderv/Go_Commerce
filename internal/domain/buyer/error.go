package buyer

import "errors"

var (
	ErrNotFound           = errors.New("buyer not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidInput       = errors.New("invalid buyer input")
	ErrConflict           = errors.New("buyer cannot be deleted: has existing orders")
)
