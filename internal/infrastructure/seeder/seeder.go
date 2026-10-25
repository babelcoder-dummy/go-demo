package seeder

import (
	"context"
	"demo/internal/core/domain"
	"fmt"

	"gorm.io/gorm"
)

func SeedAll(ctx context.Context, db *gorm.DB) {
	seedProducts(ctx, db)
}

func seedProducts(ctx context.Context, db *gorm.DB) {
	var products []domain.Product

	for i := range 10 {
		products = append(products, domain.Product{
			Name:  fmt.Sprintf("Name: %d", i+1),
			Price: (i + 1) * 100,
		})
	}

	gorm.G[domain.Product](db).CreateInBatches(ctx, &products, 10)
}
