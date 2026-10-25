package auth

import (
	"demo/internal/core/domain"
	"demo/internal/infrastructure/config"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtManager struct {
	cfg             *config.Config
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewTokenManager(cfg *config.Config) *jwtManager {
	accessDuration, err := time.ParseDuration(cfg.JWTAccessExpiry)
	if err != nil {
		accessDuration = 5 * time.Minute
	}

	refreshDuration, err := time.ParseDuration(cfg.JWTRefreshExpiry)
	if err != nil {
		refreshDuration = 2_160 * time.Hour
	}

	return &jwtManager{
		cfg:             cfg,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (t *jwtManager) GenerateTokens(user *domain.User) (string, string, error) {

	accessToken, err := t.sign("access", user, t.accessDuration)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.sign("refresh", user, t.refreshDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (t *jwtManager) ValidateToken(tokenString string) (uint, string, error) {
	keyFunc := func(*jwt.Token) (any, error) { return []byte(t.cfg.JWTSecret), nil }
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !token.Valid {
		return 0, "", errors.New("Invalid or expired token")
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)

	return uint(userID), claims.Role, nil
}

func (t *jwtManager) sign(typ string, user *domain.User, d time.Duration) (string, error) {
	now := time.Now()
	role := ""
	if typ == "access" {
		role = user.Role
	}

	claims := Claims{
		Role:      role,
		Subject:   strconv.FormatUint(uint64(user.ID), 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(d)),
	}

	JWTToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return JWTToken.SignedString([]byte(t.cfg.JWTSecret))
}
