package handler

import (
	"demo/internal/adapter/dto"
	"demo/internal/adapter/middleware"
	"demo/internal/core/domain"
	"demo/internal/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type authHandler struct {
	service      port.AuthService
	tokenManager port.TokenManager
}

func NewAuthHandler(service port.AuthService, tokenManager port.TokenManager) *authHandler {
	return &authHandler{service: service, tokenManager: tokenManager}
}

func (h *authHandler) Register(r *gin.RouterGroup) {
	authMiddleware := middleware.Authenticate(h.tokenManager)
	authGroup := r.Group("/auth")

	authGroup.POST("/register", h.RegisterUser)
	authGroup.POST("/login", h.Login)
	authGroup.GET("/me", authMiddleware, h.GetProfile)
	authGroup.POST("/refresh", h.RefreshToken)
	authGroup.DELETE("/logout", authMiddleware, h.Logout)
}

func (h *authHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user domain.User
	copier.Copy(&user, &req)

	if err := h.service.Register(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res dto.UserResponse
	copier.Copy(&res, &user)
	c.JSON(http.StatusCreated, res)
}

func (h *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var userRes dto.UserResponse
	copier.Copy(&userRes, user)
	c.JSON(http.StatusCreated, dto.LoginResponse{
		User:         userRes,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *authHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res dto.UserResponse
	copier.Copy(&res, &user)
	c.JSON(http.StatusOK, res)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *authHandler) Logout(c *gin.Context) {
	userID := c.GetUint("userID")
	if err := h.service.Logout(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
