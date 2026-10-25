package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"not null;default:null"`
	Email        string `gorm:"unique;not null;default:null"`
	Password     string `gorm:"not null;default:null"`
	Role         string `gorm:"not null;default:null;default:'customer'"`
	RefreshToken *string
}
