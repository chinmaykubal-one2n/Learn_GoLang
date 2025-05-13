package main

import (
	"log"
	"student-api/internal/db"
	"student-api/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db.Connect()

	r := gin.Default()
	handler.RegisterRoutes(r)

	r.Run(":8080")
}
