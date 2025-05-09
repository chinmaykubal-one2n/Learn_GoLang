package model

type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"required,gte=1,lte=120"`
	Email string `json:"email" binding:"required,email"`
}
