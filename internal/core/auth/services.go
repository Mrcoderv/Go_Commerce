package auth

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"Go_Commerce/internal/domain/auth"
	buyerdomain "Go_Commerce/internal/domain/buyer"
	sellerdomain "Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID   int64     `json:"uid"`
	Role     auth.Role `json:"role"`
	SellerID *int64    `json:"seller_id,omitempty"`
	BuyerID  *int64    `json:"buyer_id,omitempty"`
	jwt.RegisteredClaims
}

type Service struct {
	users   outgoing.UserRepository
	sellers outgoing.SellerRepository
	buyers  outgoing.BuyerRepository
	tx      outgoing.TxManager
	cfg     *config.Config
	log     logger.Logger
}

func NewService(
	users outgoing.UserRepository,
	sellers outgoing.SellerRepository,
	buyers outgoing.BuyerRepository,
	tx outgoing.TxManager,
	cfg *config.Config,
	log logger.Logger,
) incoming.AuthService {
	return &Service{users: users, sellers: sellers, buyers: buyers, tx: tx, cfg: cfg, log: log}
}

func (s *Service) Register(ctx context.Context, in incoming.RegisterAuthInput) (incoming.AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	name := strings.TrimSpace(in.Name)
	if email == "" || in.Password == "" || name == "" {
		return incoming.AuthResult{}, auth.ErrInvalidInput
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return incoming.AuthResult{}, auth.ErrInvalidInput
	}
	if in.Role != auth.RoleSeller && in.Role != auth.RoleBuyer {
		return incoming.AuthResult{}, auth.ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("hash password failed", "error", err)
		return incoming.AuthResult{}, fmt.Errorf("hash password: %w", err)
	}

	var user auth.User
	err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		u := auth.User{Email: email, PasswordHash: string(hash), Role: in.Role}
		switch in.Role {
		case auth.RoleSeller:
			sel, err := s.sellers.Create(ctx, sellerdomain.Seller{
				Name: name, Email: email, Phone: strings.TrimSpace(in.Phone), Address: strings.TrimSpace(in.Address),
			})
			if err != nil {
				return err
			}
			id := sel.ID
			u.SellerID = &id
		case auth.RoleBuyer:
			buy, err := s.buyers.Create(ctx, buyerdomain.Buyer{
				Name: name, Email: email, Phone: strings.TrimSpace(in.Phone), Address: strings.TrimSpace(in.Address),
			})
			if err != nil {
				return err
			}
			id := buy.ID
			u.BuyerID = &id
		}
		created, err := s.users.Create(ctx, u)
		if err != nil {
			return err
		}
		user = created
		return nil
	})
	if err != nil {
		s.log.Warn("registration failed", "email", email, "role", in.Role, "error", err)
		return incoming.AuthResult{}, err
	}

	s.log.Info("user registered", "user_id", user.ID, "role", user.Role)

	token, err := s.sign(user)
	if err != nil {
		s.log.Error("sign token failed", "user_id", user.ID, "error", err)
		return incoming.AuthResult{}, err
	}
	return incoming.AuthResult{Token: token, User: user}, nil
}

func (s *Service) Login(ctx context.Context, in incoming.LoginInput) (incoming.AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if email == "" || in.Password == "" {
		return incoming.AuthResult{}, auth.ErrInvalidInput
	}
	if in.Role != auth.RoleSeller && in.Role != auth.RoleBuyer {
		return incoming.AuthResult{}, auth.ErrInvalidInput
	}
	u, err := s.users.GetByEmailAndRole(ctx, email, in.Role)
	if err != nil {
		if err == auth.ErrNotFound {
			s.log.Warn("login failed", "email", email, "role", in.Role, "reason", "user not found")
			return incoming.AuthResult{}, auth.ErrInvalidCreds
		}
		s.log.Error("login lookup failed", "email", email, "role", in.Role, "error", err)
		return incoming.AuthResult{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		s.log.Warn("login failed", "email", email, "role", in.Role, "reason", "invalid password")
		return incoming.AuthResult{}, auth.ErrInvalidCreds
	}

	s.log.Info("user logged in", "user_id", u.ID, "role", u.Role)

	token, err := s.sign(u)
	if err != nil {
		s.log.Error("sign token failed", "user_id", u.ID, "error", err)
		return incoming.AuthResult{}, err
	}
	return incoming.AuthResult{Token: token, User: u}, nil
}

func (s *Service) sign(u auth.User) (string, error) {
	claims := Claims{
		UserID:   u.ID,
		Role:     u.Role,
		SellerID: u.SellerID,
		BuyerID:  u.BuyerID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", u.ID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func ParseToken(secret, tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrUnauthorized
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, auth.ErrUnauthorized
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, auth.ErrUnauthorized
	}
	return claims, nil
}
