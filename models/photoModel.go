package models

import (
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	Title     string    `gorm:"not null"`
	Caption   string
	PhotoURL  string    `gorm:"not null"`
	UserID    uint      `gorm:"not null;"`
	User   User `gorm:"constraint:OnDelete:CASCADE;"`
}