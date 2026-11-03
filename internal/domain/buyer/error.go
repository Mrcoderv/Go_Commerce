package buyer

import "errors"

var (
	ErrNotFound           = errors.New(" not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
