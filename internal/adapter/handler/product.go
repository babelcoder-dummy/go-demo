package handler

import (
	"demo/internal/adapter/dto"
	"demo/internal/adapter/middleware"
	"demo/internal/core/domain"
	"demo/internal/core/port"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type productHandler struct {
	service      port.ProductService
	tokenManager port.TokenManager
	enforcer     *casbin.Enforcer
}

func NewProductHandler(service port.ProductService, tokenManager port.TokenManager, enforcer *casbin.Enforcer) *productHandler {
	return &productHandler{service: service, tokenManager: tokenManager, enforcer: enforcer}
}

func (h *productHandler) Register(r *gin.RouterGroup) {
	authnMiddleware := middleware.Authenticate(h.tokenManager)
	authzMiddleware := middleware.Authorize(h.enforcer)
	// /api/v1/products
	productGroup := r.Group("/products")

	productGroup.GET("", h.GetAll)
	productGroup.GET("/:id", h.GetByID)
	productGroup.POST("", authnMiddleware, authzMiddleware, h.Create)
}

func (h *productHandler) Create(c *gin.Context) {
	var form dto.CreateProductRequest // DTO
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadedPath := filepath.Join("uploads", form.Image.Filename)
	if err := c.SaveUploadedFile(form.Image, uploadedPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product := domain.Product{
		Name:  form.Name,
		Price: form.Price,
		Image: &form.Image.Filename,
	}

	ctx := c.Request.Context()
	err := h.service.Create(ctx, &product)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var resp dto.ProductResponse
	copier.Copy(&resp, &product)

	c.JSON(http.StatusCreated, resp)
}

func (h *productHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	ctx := c.Request.Context()
	product, err := h.service.GetByID(ctx, uint(productID))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var resp dto.ProductResponse
	copier.Copy(&resp, &product)
	c.JSON(http.StatusOK, resp)
}

func (h *productHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	products, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var resp []dto.ProductResponse
	copier.Copy(&resp, &products)
	c.JSON(http.StatusOK, resp)
}
