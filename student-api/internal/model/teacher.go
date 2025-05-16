package model

type Teacher struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"` // e.g., "admin", "regular", etc.
}
