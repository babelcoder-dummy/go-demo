package main

import (
	"demo/internal/adapter/handler"
	"demo/internal/adapter/repository"
	"demo/internal/core/service"
	"demo/internal/infrastructure/auth"
	"demo/internal/infrastructure/config"
	"demo/internal/infrastructure/database"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = false
	corsCfg.AllowOrigins = []string{"http://127.0.0.1:5500"}
	corsCfg.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsCfg.AllowCredentials = true

	r.Use(cors.New(corsCfg))
	r.Static("/assets", "./uploads") // /assets/image.png => ./uploads/image.png

	enforcer, err := config.NewCasbinEnforcer()
	if err != nil {
		log.Fatalf("Failed to create casbin enforcer: %v", err)
	}

	// GET /api/v1/products
	rGroup := r.Group("/api/v1")

	userRepository := repository.NewUserRepository(db)
	productRepository := repository.NewProductRepository(db)

	hasher := auth.NewBcryptHasher()
	jwtManager := auth.NewTokenManager(cfg)
	authService := service.NewAuthService(userRepository, hasher, jwtManager)
	productService := service.NewProductService(productRepository)

	authHandler := handler.NewAuthHandler(authService, jwtManager)
	productHandler := handler.NewProductHandler(productService, jwtManager, enforcer)

	productHandler.Register(rGroup)
	authHandler.Register(rGroup)

	if err := r.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
