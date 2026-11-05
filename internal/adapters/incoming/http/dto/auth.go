package dto

import "Go_Commerce/internal/domain/auth"

type RegisterAuthRequest struct {
	Body struct {
		Email    string    `json:"email" format:"email"`
		Password string    `json:"password" minLength:"6"`
		Role     auth.Role `json:"role" enum:"seller,buyer"`
		Name     string    `json:"name" minLength:"1"`
		Phone    string    `json:"phone,omitempty"`
		Address  string    `json:"address,omitempty"`
	}
}

type LoginAuthRequest struct {
	Body struct {
		Email    string    `json:"email" format:"email"`
		Password string    `json:"password"`
		Role     auth.Role `json:"role" enum:"seller,buyer"`
	}
}

type AuthResponse struct {
	Body struct {
		Token string    `json:"token"`
		User  auth.User `json:"user"`
	}
}
