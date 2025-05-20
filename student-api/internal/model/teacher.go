// Package model contains the data structures used in the student API
package model

// Teacher represents a teacher in the system
// @Description Teacher account information for authentication and authorization
type Teacher struct {
	// Unique identifier for the teacher
	// @example "123e4567-e89b-12d3-a456-426614174000"
	ID string `gorm:"primaryKey" json:"id"`

	// Username for login
	// @example "john.doe"
	// @required
	Username string `gorm:"uniqueIndex" json:"username" binding:"required"`

	// Password for authentication (will be hashed)
	// @example "securePassword123"
	// @required
	Password string `json:"password" binding:"required"`

	// Email address of the teacher
	// @example "john.doe@school.com"
	// @required
	// @pattern ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	Email string `json:"email" binding:"required,email"`

	// Role determines the teacher's permissions (admin or regular)
	// @example "admin"
	// @enum ["admin","regular"]
	// @required
	Role string `json:"role" binding:"required"` // e.g., "admin", "regular", etc.
}
