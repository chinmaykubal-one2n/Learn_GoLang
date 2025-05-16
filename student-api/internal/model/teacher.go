package model

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"` // e.g., "admin", "regular", etc.
}
