package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"student-api/internal/db"
	"student-api/internal/handler"
	"student-api/internal/middleware"
	"student-api/internal/service"
	"syscall"
	"time"

	_ "student-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Swagger Student API
//	@version		1.0
//	@description	This is a sample server Student API server.

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

	// https://gin-gonic.com/en/docs/examples/graceful-restart-or-stop/
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routerEngine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")
}
