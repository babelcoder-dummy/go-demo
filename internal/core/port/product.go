package port

import (
	"context"
	"demo/internal/core/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uint) (*domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
}

type ProductService interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id uint) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
}
