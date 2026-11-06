package dto

import "Go_Commerce/internal/domain/auth"

type RegisterAuthRequest struct {
	Body struct {
		Email    string    `json:"email" format:"email" doc:"Account email (must be unique per role)"`
		Password string    `json:"password" minLength:"6" doc:"Password, at least 6 characters"`
		Role     auth.Role `json:"role" enum:"seller,buyer" doc:"seller or buyer"`
		Name     string    `json:"name" minLength:"1" doc:"Display name"`
		Phone    string    `json:"phone,omitempty" doc:"Optional phone number"`
		Address  string    `json:"address,omitempty" doc:"Optional address"`
	}
}

type LoginAuthRequest struct {
	Body struct {
		Email    string    `json:"email" format:"email" doc:"Account email"`
		Password string    `json:"password" doc:"Account password"`
		Role     auth.Role `json:"role" enum:"seller,buyer" doc:"seller or buyer"`
	}
}

type AuthResponse struct {
	Body struct {
		Token string    `json:"token"`
		User  auth.User `json:"user"`
	}
}
