package repository

import (
	"context"
	"demo/internal/core/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return gorm.G[domain.User](r.db).Create(ctx, user)
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := gorm.G[domain.User](r.db).Where("id = ?", id).First(ctx)
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := gorm.G[domain.User](r.db).Where("email = ?", email).First(ctx)
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, id uint, user *domain.User) error {
	_, err := gorm.G[domain.User](r.db).Where("id = ?", id).Updates(ctx, *user)
	return err
}

func (r *userRepository) ClearRefreshToken(ctx context.Context, id uint) error {
	_, err := gorm.G[domain.User](r.db).Where("id = ?", id).Update(ctx, "refresh_token", nil)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	_, err := gorm.G[domain.User](r.db).Where("id = ?", id).Delete(ctx)
	return err
}
