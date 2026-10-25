package service

import (
	"context"
	"demo/internal/core/domain"
	"demo/internal/core/port"
	"errors"

	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type authService struct {
	repo         port.UserRepository
	hasher       port.PasswordHasher
	tokenManager port.TokenManager
}

func NewAuthService(repo port.UserRepository, hasher port.PasswordHasher, tokenManager port.TokenManager) *authService {
	return &authService{repo: repo, hasher: hasher, tokenManager: tokenManager}
}

func (s *authService) Register(ctx context.Context, user *domain.User) error {
	if _, err := s.repo.GetByEmail(ctx, user.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("email already exists")
	}

	hashedPassword, err := s.hasher.Hash(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Role = "customer"

	return s.repo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, string, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", "", ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.Password, password); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	user.RefreshToken = &refreshToken
	if err := s.repo.Update(ctx, user.ID, user); err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, token string) (string, string, error) {
	// Validate Token
	userID, _, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		return "", "", err
	}

	// Find User
	user, err := s.repo.GetByID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", errors.New("user not found")
	}

	if *user.RefreshToken != token {
		return "", "", errors.New("invalid refresh token")
	}

	// Generate Tokens
	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(user)

	user.RefreshToken = &refreshToken
	if err := s.repo.Update(ctx, user.ID, user); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) GetProfile(ctx context.Context, userID uint) (*domain.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	return s.repo.ClearRefreshToken(ctx, userID)
}
