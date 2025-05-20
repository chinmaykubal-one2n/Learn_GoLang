package main

import (
	"log"
	"student-api/internal/db"
	"student-api/internal/handler"
	"student-api/internal/middleware"
	"student-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbInstance := db.Connect()

	studentService := &service.StudentServiceImpl{
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

	routerEngine.POST("/register", teacherHandler.RegisterTeacher)

	authMiddleware, err := middleware.AuthMiddleware()
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
