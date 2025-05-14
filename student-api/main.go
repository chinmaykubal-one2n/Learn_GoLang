package main

import (
	"log"
	"student-api/internal/db"
	"student-api/internal/handler"
	"student-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db.Connect()

	r := gin.Default()

	authMiddleware, err := middleware.AuthMiddleware()
	if err != nil {
		log.Fatalf("JWT Error: %s", err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		handler.RegisterRoutes(api)
	}

	r.Run(":8080")
}
