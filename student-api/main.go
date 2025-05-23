package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"student-api/internal/db"
	"student-api/internal/handler"
	logging "student-api/internal/logger"
	"student-api/internal/middleware"
	"student-api/internal/otel"
	"student-api/internal/service"
	"syscall"
	"time"

	_ "student-api/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

//	@title			Swagger Student API
//	@version		1.0
//	@description	This is a sample server Student API server.

func main() {
	// https://github.dev/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	// Initialize the logger
	err := logging.InitLogger(ctx, os.Getenv("SERVICE_NAME"), os.Getenv("OTLP_ENDPOINT"))
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logging.Logger.Sync()
	logging.Logger.Info("Starting Student API")

	// Initialize the OpenTelemetry tracer
	cleanup := otel.InitTracer()
	defer cleanup(context.Background())

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

	// OpenTelemetry otelgin middleware
	routerEngine.Use(otelgin.Middleware(os.Getenv("SERVICE_NAME")))

	routerEngine.GET("/healthz", h.HealthCheck)
	routerEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // http://localhost:8090/swagger/index.html

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

	srv := &http.Server{
		Addr:    ":8090",
		Handler: routerEngine,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
