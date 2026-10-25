package service

import (
	"context"
	"demo/internal/core/domain"
	"demo/internal/core/port"
)

type productService struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) *productService {
	return &productService{repo: repo}
}

func (s *productService) GetAll(ctx context.Context) ([]domain.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) Create(ctx context.Context, product *domain.Product) error {
	return s.repo.Create(ctx, product)
}
