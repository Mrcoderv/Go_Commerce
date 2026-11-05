package seller

import "errors"

var (
	ErrNotFound           = errors.New("seller not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrConflict           = errors.New("seller conflict")
)
