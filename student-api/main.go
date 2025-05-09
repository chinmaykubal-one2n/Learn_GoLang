package main

import (
	"student-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	handler.RegisterRoutes(r)
	r.Run(":8080")
}
