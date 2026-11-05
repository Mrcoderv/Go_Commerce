package incoming

import (
	"Go_Commerce/internal/domain/auth"
	"context"
)

type RegisterAuthInput struct {
	Email    string
	Password string
	Role     auth.Role
	Name     string
	Phone    string
	Address  string
}

type LoginInput struct {
	Email    string
	Password string
	Role     auth.Role
}

type AuthResult struct {
	Token string    `json:"token"`
	User  auth.User `json:"user"`
}

type AuthService interface {
	Register(ctx context.Context, in RegisterAuthInput) (AuthResult, error)
	Login(ctx context.Context, in LoginInput) (AuthResult, error)
}
