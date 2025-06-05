package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"student-api/internal/handler"
	"student-api/internal/model"
	"testing"

	logging "student-api/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// setupLoggerForHandlerTests initializes a zap logger for use in tests
func setupLoggerForTeacherHandlerTests() {
	zapLogger := zap.NewExample()
	logging.Logger = otelzap.New(zapLogger)
}

type MockTeacherService struct {
	mock.Mock
}

func (m *MockTeacherService) CreateTeacher(t model.Teacher, ctx context.Context) (model.Teacher, error) {
	args := m.Called(t, ctx)
	return args.Get(0).(model.Teacher), args.Error(1)
}

func (m *MockTeacherService) GetTeacher(username string, ctx context.Context) (model.Teacher, error) {
	args := m.Called(username, ctx)
	return args.Get(0).(model.Teacher), args.Error(1)
}

func TestRegisterTeacher(t *testing.T) {
	setupLoggerForTeacherHandlerTests()

	mockService := new(MockTeacherService)
	h := handler.TeacherHandler{Service: mockService}
	router := gin.Default()
	router.POST("/register", h.RegisterTeacher)

	t.Run("success", func(t *testing.T) {
		input := model.Teacher{Username: "teacher1", Password: "password123", Email: "teacher1@gmail.com", Role: "admin"}
		output := input
		output.ID = "new-teacher-id"
		body := `{"username":"teacher1","password":"password123", "email":"teacher1@gmail.com", "role":"admin"}`

		mockService.On("CreateTeacher", input, mock.Anything).Return(output, nil)

		req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertCalled(t, "CreateTeacher", input, mock.Anything)
	})

	t.Run("bad request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/register", strings.NewReader(`bad-json`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		input := model.Teacher{Username: "teacher2", Password: "password1234", Email: "teacher2@gmail.com", Role: "admin"}
		body := `{"username":"teacher2","password":"password1234", "email":"teacher2@gmail.com", "role":"admin"}`

		mockService.On("CreateTeacher", input, mock.Anything).Return(model.Teacher{}, errors.New("something went wrong"))

		req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertCalled(t, "CreateTeacher", input, mock.Anything)
	})
}
