package repository

import (
	"context"
	"demo/internal/core/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return gorm.G[domain.Product](r.db).Create(ctx, product)
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	p, err := gorm.G[domain.Product](r.db).Where("id = ?", id).First(ctx)
	return &p, err
}

func (r *productRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	return gorm.G[domain.Product](r.db).Find(ctx)
}
