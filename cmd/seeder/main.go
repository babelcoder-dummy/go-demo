package main

import (
	"context"
	"demo/internal/infrastructure/config"
	"demo/internal/infrastructure/database"
	"demo/internal/infrastructure/seeder"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

	seeder.SeedAll(ctx, db)
}
