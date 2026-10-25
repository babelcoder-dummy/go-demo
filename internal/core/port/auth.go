package port

import (
	"context"
	"demo/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, password string) (*domain.User, string, string, error)
	GetProfile(ctx context.Context, userID uint) (*domain.User, error)
	RefreshToken(ctx context.Context, token string) (string, string, error)
	Logout(ctx context.Context, userID uint) error
}
