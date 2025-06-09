package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handler "student-api/internal/handler"
	logging "student-api/internal/logger"
	"student-api/internal/middleware"
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

	studentService := &service.StudentServiceImpl{DB: db}
	teacherService := &service.TeacherServiceImpl{DB: db}

	authMiddleware, authMiddlewareErr := middleware.AuthMiddleware(teacherService)
	assert.NoError(t, authMiddlewareErr)

	h := handler.NewHandler(studentService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	api := router.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		h.RegisterRoutes(api)
	}

	generateToken := func(role string) string {
		user := &middleware.User{
			UserName: "Noel Johnson",
			Role:     role,
		}
		token, _, err := authMiddleware.TokenGenerator(user)
		assert.NoError(t, err)
		return token
	}

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "admin role - success",
			token:          generateToken("admin"),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"students":[],"page":1,"page_size":100}`,
		},
		{
			name:           "regular role - success",
			token:          generateToken("regular"),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"students":[],"page":1,"page_size":100}`,
		},
		{
			name:           "missing token - unauthorized",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"cookie token is empty"}`,
		},
		{
			name:           "invalid token - unauthorized",
			token:          "invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid character '\\u008a' looking for beginning of value"}`, // just matching the error message for invalid token, from the terminal output
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/students", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)
			assert.JSONEq(t, tc.expectedBody, resp.Body.String())
		})
	}
}

func TestCreateStudentIntegration(t *testing.T) {
	setupLoggerForIntegrationTests()
	db := SetupInMemoryDB()

	studentService := &service.StudentServiceImpl{DB: db}
	teacherService := &service.TeacherServiceImpl{DB: db}

	authMiddleware, authMiddlewareErr := middleware.AuthMiddleware(teacherService)
	assert.NoError(t, authMiddlewareErr)

	h := handler.NewHandler(studentService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	api := router.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		h.RegisterRoutes(api)
	}

	generateToken := func(role string) string {
		user := &middleware.User{
			UserName: "Noel Johnson",
			Role:     role,
		}
		token, _, err := authMiddleware.TokenGenerator(user)
		assert.NoError(t, err)
		return token
	}

	requestBody := `{
		"name": "John Doe",
		"email": "john.doe@example.com",
		"age": 21
	}`

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "create student with admin role",
			token:          generateToken("admin"),
			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"name": "John Doe",
				"email": "john.doe@example.com",
				"age": 21
			}`,
		},
		{
			name:           "create student with regular role",
			token:          generateToken("regular"),
			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"name": "John Doe",
				"email": "john.doe@example.com",
				"age": 21
			}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/students", strings.NewReader(requestBody))
			req.Header.Set("Authorization", "Bearer "+tc.token)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			var actualBody map[string]interface{} // striping field ID from actual body for proper comparison
			err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.NoError(t, err)

			delete(actualBody, "id")
			actualBodyCleaned, err := json.Marshal(actualBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, resp.Code)
			assert.JSONEq(t, tc.expectedBody, string(actualBodyCleaned))
		})
	}
}
