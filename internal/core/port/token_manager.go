package port

import "demo/internal/core/domain"

type TokenManager interface {
	GenerateTokens(user *domain.User) (string, string, error)
	ValidateToken(tokenString string) (uint, string, error)
}
