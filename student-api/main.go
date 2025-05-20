package main

import (
	"log"
	"student-api/internal/db"
	"student-api/internal/handler"
	"student-api/internal/middleware"
	"student-api/internal/service"

	_ "student-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbInstance := db.Connect()

	studentService := &service.StudentServiceImpl{
		DB: dbInstance,
	}

	teacherService := &service.TeacherServiceImpl{
		DB: dbInstance,
	}

	teacherHandler := &handler.TeacherHandler{
		Service: &service.TeacherServiceImpl{
			DB: dbInstance,
		},
	}

	h := handler.NewHandler(studentService)

	routerEngine := gin.Default()
	routerEngine.GET("/healthz", h.HealthCheck)
	routerEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	routerEngine.POST("/register", teacherHandler.RegisterTeacher)

	authMiddleware, err := middleware.AuthMiddleware(teacherService)
	if err != nil {
		log.Fatalf("JWT Error: %s", err.Error())
	}

	routerEngine.POST("/login", authMiddleware.LoginHandler)
	routerEngine.GET("/refresh_token", authMiddleware.RefreshHandler)

	api := routerEngine.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		h.RegisterRoutes(api)
	}

	routerEngine.Run(":8080")
}
