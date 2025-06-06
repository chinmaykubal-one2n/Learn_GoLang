package integration

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "student-api/internal/handler"
	logging "student-api/internal/logger"
	"student-api/internal/model"
	"student-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupLoggerForIntegrationTests() {
	zapLogger := zap.NewExample()
	logging.Logger = otelzap.New(zapLogger)
}

func SetupInMemoryDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("[test-db]: Failed to connect to in-memory database: %v", err)
	}

	err = db.AutoMigrate(
		&model.Student{},
		&model.Teacher{},
	)
	if err != nil {
		log.Fatalf("[test-db]: Failed to run migrations on in-memory DB: %v", err)
	}

	return db
}

func TestHealthCheckIntegration(t *testing.T) {
	setupLoggerForIntegrationTests()

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "health check",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"ok","timestamp":"`,
		},
	}

	gin.SetMode(gin.TestMode)

	h := handler.NewHandler(nil)

	router := gin.Default()
	router.GET("/healthz", h.HealthCheck)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			assert.Contains(t, resp.Body.String(), tc.expectedBody)
		})
	}
}

func TestListStudentsIntegration(t *testing.T) {
	setupLoggerForIntegrationTests()

	db := SetupInMemoryDB()
	studentService := &service.StudentServiceImpl{
		DB: db,
	}

	h := handler.NewHandler(studentService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api := router.Group("/api")
	h.RegisterRoutes(api)

	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "list students",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"students":[],"page":1,"page_size":100}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/students", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)
			assert.JSONEq(t, tc.expectedBody, resp.Body.String())
		})
	}
}
