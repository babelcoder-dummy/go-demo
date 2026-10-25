package port

import (
	"context"
	"demo/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, id uint, user *domain.User) error
	ClearRefreshToken(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}
