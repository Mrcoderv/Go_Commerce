package outgoing

import (
	"Go_Commerce/internal/domain/auth"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, u auth.User) (auth.User, error)
	GetByEmailAndRole(ctx context.Context, email string, role auth.Role) (auth.User, error)
	GetByID(ctx context.Context, id int64) (auth.User, error)
}
