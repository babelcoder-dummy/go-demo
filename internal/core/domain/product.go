package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name  string  `json:"name" gorm:"not null;default:null"`
	Price int     `json:"price"  gorm:"not null;default:null"`
	Image *string `json:"image"`
}
