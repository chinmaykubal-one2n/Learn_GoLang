// Package model contains the data structures for the student API
package model

// Student represents a student entity in the system
// @Description Student information including personal details
type Student struct {
	// Unique identifier of the student
	// @example "550e8400-e29b-41d4-a716-446655440000"
	ID string `gorm:"primaryKey" json:"id"`

	// Full name of the student
	// @example "John Smith"
	// @required true
	Name string `json:"name" binding:"required"`

	// Age of the student (must be between 1 and 120)
	// @minimum 1
	// @maximum 120
	// @example 18
	// @required true
	Age int `json:"age" binding:"required,gte=1,lte=120"`

	// Email address of the student
	// @example "john.smith@example.com"
	// @pattern ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
	// @required true
	Email string `json:"email" binding:"required,email"`
}
