package integration

import (
	"context"
	"encoding/json"
	"fmt"
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
			name:           "missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"cookie token is empty"}`,
		},
		{
			name:           "invalid token",
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

func TestGetStudentIntegration(t *testing.T) {
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

	// First create a student to retrieve
	created, err := studentService.CreateStudent(model.Student{
		Name:  "Jane Doe",
		Email: "jane.doe@example.com",
		Age:   22,
	}, context.Background())
	assert.NoError(t, err)

	tests := []struct {
		name           string
		id             string
		token          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "get student with valid admin token",
			id:             created.ID,
			token:          generateToken("admin"),
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"name": "Jane Doe",
				"email": "jane.doe@example.com",
				"age": 22
			}`,
		},
		{
			name:           "student not found",
			id:             "non-existent-id",
			token:          generateToken("admin"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"student not found"}`,
		},
		{
			name:           "invalid token",
			id:             created.ID,
			token:          "invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid character '\\u008a' looking for beginning of value"}`,
		},
		{
			name:           "missing token",
			id:             created.ID,
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"cookie token is empty"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/students/%s", tc.id)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tc.token)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.NoError(t, err)

			if tc.expectedStatus == http.StatusOK {
				delete(actualBody, "id")
			}

			cleaned, err := json.Marshal(actualBody)
			assert.NoError(t, err)

			assert.JSONEq(t, tc.expectedBody, string(cleaned))
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

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		requestBody    string
		expectedBody   string
	}{
		{
			name:           "create student with admin role",
			token:          generateToken("admin"),
			expectedStatus: http.StatusCreated,
			requestBody: `{
				"name": "John Doe",
				"email": "john.doe@example.com",
				"age": 21
			}`,
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
			requestBody: `{
				"name": "John Doe",
				"email": "john.doe@example.com",
				"age": 21
			}`,
			expectedBody: `{
				"name": "John Doe",
				"email": "john.doe@example.com",
				"age": 21
			}`,
		},
		{
			name:           "missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			requestBody: `{
				"name": "John Doe",
				"email": "john.doe@example.com",
				"age": 21
			}`,
			expectedBody: `{"error":"cookie token is empty"}`,
		},
		{
			name:           "invalid token",
			token:          "invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			requestBody: `{
				"name": "John Doe",
				"email": "john.doe@example.com",
				"age": 21
			}`,
			expectedBody: `{"error":"invalid character '\\u008a' looking for beginning of value"}`,
		},
		{
			name:           "invalid age type",
			token:          generateToken("regular"),
			expectedStatus: http.StatusBadRequest,
			requestBody: `{
				"name": "Invalid User",
				"email": "invalid@example.com",
				"age": "twenty"
			}`,
			expectedBody: `{"error":"json: cannot unmarshal string into Go struct field Student.age of type int"}`,
		},
		{
			name:           "missing required field",
			token:          generateToken("regular"),
			expectedStatus: http.StatusBadRequest,
			requestBody: `{
				"email": "missing.name@example.com",
				"age": 18
			}`,
			expectedBody: `{"error":"Key: 'Student.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/students", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tc.token)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.NoError(t, err)

			if tc.expectedStatus == http.StatusCreated {
				delete(actualBody, "id")
			}

			cleaned, err := json.Marshal(actualBody)
			assert.NoError(t, err)

			assert.JSONEq(t, tc.expectedBody, string(cleaned))
		})
	}
}

func TestUpdateStudentIntegration(t *testing.T) {
	setupLoggerForIntegrationTests()
	db := SetupInMemoryDB()

	studentService := &service.StudentServiceImpl{DB: db}
	teacherService := &service.TeacherServiceImpl{DB: db}

	authMiddleware, err := middleware.AuthMiddleware(teacherService)
	assert.NoError(t, err)

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

	// create a student to update
	initialStudent := model.Student{Name: "Jane", Email: "jane@example.com", Age: 22}
	created, err := studentService.CreateStudent(initialStudent, context.Background())
	assert.NoError(t, err)

	tests := []struct {
		name           string
		id             string
		token          string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:  "update student with valid data and admin role",
			id:    created.ID,
			token: generateToken("admin"),
			requestBody: `{
				"name": "Jane Updated",
				"email": "updated@example.com",
				"age": 23
			}`,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"name": "Jane Updated",
				"email": "updated@example.com",
				"age": 23
			}`,
		},
		{
			name:           "update student with missing token",
			id:             created.ID,
			token:          "",
			requestBody:    `{"name":"X","email":"x@example.com","age":10}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"cookie token is empty"}`,
		},
		{
			name:           "update student with invalid ID",
			id:             "non-existent-id",
			token:          generateToken("regular"),
			requestBody:    `{"name":"New","email":"new@example.com","age":24}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Student not found"}`,
		},
		{
			name:           "update student with invalid body",
			id:             created.ID,
			token:          generateToken("regular"),
			requestBody:    `{"name": 123, "email": "bad", "age": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"json: cannot unmarshal number into Go struct field Student.name of type string"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/students/%s", tc.id)
			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tc.token)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			var actualBody map[string]interface{}
			err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
			assert.NoError(t, err)

			if tc.expectedStatus == http.StatusOK {
				delete(actualBody, "id")
			}

			cleaned, err := json.Marshal(actualBody)
			assert.NoError(t, err)

			assert.JSONEq(t, tc.expectedBody, string(cleaned))
		})
	}
}

func TestDeleteStudentIntegration(t *testing.T) {
	setupLoggerForIntegrationTests()
	db := SetupInMemoryDB()

	studentService := &service.StudentServiceImpl{DB: db}
	teacherService := &service.TeacherServiceImpl{DB: db}

	authMiddleware, err := middleware.AuthMiddleware(teacherService)
	assert.NoError(t, err)

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

	// Create a student to delete
	initialStudent := model.Student{Name: "Delete Me", Email: "deleteme@example.com", Age: 25}
	created, err := studentService.CreateStudent(initialStudent, context.Background())
	assert.NoError(t, err)

	tests := []struct {
		name           string
		id             string
		token          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "delete student with admin role",
			id:             created.ID,
			token:          generateToken("admin"),
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "delete student with missing token",
			id:             created.ID,
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"cookie token is empty"}`,
		},
		{
			name:           "delete student with invalid token",
			id:             created.ID,
			token:          "invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid character '\\u008a' looking for beginning of value"}`,
		},
		{
			name:           "delete student that does not exist",
			id:             "non-existent-id",
			token:          generateToken("admin"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Student not found"}`,
		},
		{
			name:           "delete student with regular role",
			id:             created.ID,
			token:          generateToken("regular"),
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"you don't have permission to access this resource"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/students/%s", tc.id)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tc.token)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedStatus != http.StatusNoContent {
				assert.JSONEq(t, tc.expectedBody, resp.Body.String())
			}

			if tc.expectedStatus == http.StatusNoContent {
				assert.Equal(t, tc.expectedBody, resp.Body.String())
			}
		})
	}
}
