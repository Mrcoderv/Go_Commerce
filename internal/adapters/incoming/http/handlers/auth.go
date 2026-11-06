package handlers

import (
	"Go_Commerce/internal/adapters/incoming/http/dto"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/ports/incoming"
	"context"
)

type AuthHandler struct {
	svc incoming.AuthService
	log logger.Logger
}

func NewAuthHandler(svc incoming.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{svc: svc, log: log}
}

func (h *AuthHandler) Register(ctx context.Context, req *dto.RegisterAuthRequest) (*dto.AuthResponse, error) {
	res, err := h.svc.Register(ctx, incoming.RegisterAuthInput{
		Email:    req.Body.Email,
		Password: req.Body.Password,
		Role:     req.Body.Role,
		Name:     req.Body.Name,
		Phone:    req.Body.Phone,
		Address:  req.Body.Address,
	})
	if err != nil {
		return nil, MapError(h.log, err)
	}
	var out dto.AuthResponse
	out.Body.Token = res.Token
	out.Body.User = res.User
	return &out, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *dto.LoginAuthRequest) (*dto.AuthResponse, error) {
	res, err := h.svc.Login(ctx, incoming.LoginInput{
		Email:    req.Body.Email,
		Password: req.Body.Password,
		Role:     req.Body.Role,
	})
	if err != nil {
		return nil, MapError(h.log, err)
	}
	var out dto.AuthResponse
	out.Body.Token = res.Token
	out.Body.User = res.User
	return &out, nil
}
