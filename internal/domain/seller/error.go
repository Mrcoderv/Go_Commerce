package seller

import "errors"

var (
	ErrNotFound           = errors.New("seller not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)