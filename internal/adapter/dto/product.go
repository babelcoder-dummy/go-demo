package dto

import "mime/multipart"

type CreateProductRequest struct {
	Name  string                `form:"name" binding:"required"`
	Price int                   `form:"price" binding:"required,min=1"`
	Image *multipart.FileHeader `form:"image"`
}

type ProductResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price int     `json:"price"`
	Image *string `json:"image"`
}
