package auth

import "errors"

var (
	ErrNotFound     = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrEmailExists  = errors.New("email already exists")
)
