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

	routerEngine := gin.Default()
	routerEngine.GET("/healthz", handler.HealthCheck)

	authMiddleware, err := middleware.AuthMiddleware()
	if err != nil {
		log.Fatalf("JWT Error: %s", err.Error())
	}

	routerEngine.POST("/login", authMiddleware.LoginHandler)
	routerEngine.GET("/refresh_token", authMiddleware.RefreshHandler)

	api := routerEngine.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		handler.RegisterRoutes(api)
	}

	routerEngine.Run(":8080")
}
