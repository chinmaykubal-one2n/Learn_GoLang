package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"student-api/internal/handler"
	"student-api/internal/model"
	"student-api/internal/pkg/tests"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStudentService struct {
	mock.Mock
}

func (m *MockStudentService) ListStudents(ctx context.Context) ([]model.Student, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Student), args.Error(1)
}

func (m *MockStudentService) GetStudent(id string, ctx context.Context) (model.Student, error) {
	args := m.Called(id, ctx)
	return args.Get(0).(model.Student), args.Error(1)
}

func (m *MockStudentService) CreateStudent(s model.Student, ctx context.Context) (model.Student, error) {
	args := m.Called(s, ctx)
	return args.Get(0).(model.Student), args.Error(1)
}

func (m *MockStudentService) UpdateStudent(id string, s model.Student, ctx context.Context) (model.Student, error) {
	args := m.Called(id, s, ctx)
	return args.Get(0).(model.Student), args.Error(1)
}

func (m *MockStudentService) DeleteStudent(id string, ctx context.Context) error {
	args := m.Called(id, ctx)
	return args.Error(0)
}

func TestGetStudent(t *testing.T) {
	tests.SetupLoggerForTests()
	mockService := new(MockStudentService)
	h := handler.NewHandler(mockService)
	router := gin.Default()
	h.RegisterRoutes(router.Group("/"))

	t.Run("success", func(t *testing.T) {
		mockStudent := model.Student{ID: "1", Name: "Alice", Age: 21, Email: "alice@example.com"}
		mockService.On("GetStudent", "1", mock.Anything).Return(mockStudent, nil)

		req, _ := http.NewRequest("GET", "/students/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "GetStudent", "1", mock.Anything)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetStudent", "999", mock.Anything).Return(model.Student{}, errors.New("Student not found"))

		req, _ := http.NewRequest("GET", "/students/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertCalled(t, "GetStudent", "999", mock.Anything)

		var respBody map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &respBody)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Student not found", respBody["error"])
	})
}

func TestCreateStudent(t *testing.T) {
	tests.SetupLoggerForTests()
	mockService := new(MockStudentService)
	h := handler.NewHandler(mockService)
	router := gin.Default()
	h.RegisterRoutes(router.Group("/"))

	t.Run("success", func(t *testing.T) {
		input := model.Student{Name: "Bob", Age: 22, Email: "bob@example.com"}
		output := input
		output.ID = "new-id"
		body := `{"name":"Bob","age":22,"email":"bob@example.com"}`

		mockService.On("CreateStudent", input, mock.Anything).Return(output, nil)

		req, _ := http.NewRequest("POST", "/students", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertCalled(t, "CreateStudent", input, mock.Anything)
	})

	t.Run("bad request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/students", strings.NewReader(`bad-json`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		input := model.Student{Name: "Charlie", Age: 23, Email: "charlie@example.com"}
		body := `{"name":"Charlie","age":23,"email":"charlie@example.com"}`

		mockService.On("CreateStudent", input, mock.Anything).Return(model.Student{}, errors.New("something went wrong"))

		req, _ := http.NewRequest("POST", "/students", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertCalled(t, "CreateStudent", input, mock.Anything)
	})
}

func TestUpdateStudent(t *testing.T) {
	tests.SetupLoggerForTests()
	mockService := new(MockStudentService)
	h := handler.NewHandler(mockService)
	router := gin.Default()
	h.RegisterRoutes(router.Group("/"))

	t.Run("success", func(t *testing.T) {
		id := "1"
		input := model.Student{Name: "Updated Name", Age: 23, Email: "updated@example.com"}
		output := input
		output.ID = id

		body := `{"name":"Updated Name","age":23,"email":"updated@example.com"}`

		mockService.On("UpdateStudent", id, input, mock.Anything).Return(output, nil)

		req, _ := http.NewRequest("PUT", "/students/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "UpdateStudent", id, input, mock.Anything)
	})

	t.Run("not found", func(t *testing.T) {
		id := "999"
		input := model.Student{Name: "Ghost", Age: 100, Email: "ghost@example.com"}
		body := `{"name":"Ghost","age":100,"email":"ghost@example.com"}`

		mockService.On("UpdateStudent", id, input, mock.Anything).Return(model.Student{}, errors.New("Student not found"))

		req, _ := http.NewRequest("PUT", "/students/999", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertCalled(t, "UpdateStudent", id, input, mock.Anything)

		var respBody map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &respBody)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Student not found", respBody["error"])
	})

	t.Run("bad request", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/students/1", strings.NewReader(`invalid-json`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteStudent(t *testing.T) {
	tests.SetupLoggerForTests()
	mockService := new(MockStudentService)
	h := handler.NewHandler(mockService)
	router := gin.Default()
	h.RegisterRoutes(router.Group("/"))

	t.Run("success", func(t *testing.T) {
		id := "1"
		mockService.On("DeleteStudent", id, mock.Anything).Return(nil)

		req, _ := http.NewRequest("DELETE", "/students/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertCalled(t, "DeleteStudent", id, mock.Anything)
	})

	t.Run("not found", func(t *testing.T) {
		id := "999"
		mockService.On("DeleteStudent", id, mock.Anything).Return(errors.New("Student not found"))

		req, _ := http.NewRequest("DELETE", "/students/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertCalled(t, "DeleteStudent", id, mock.Anything)

		var respBody map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &respBody)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Student not found", respBody["error"])
	})
}

func TestListStudents(t *testing.T) {
	tests.SetupLoggerForTests()
	t.Run("success", func(t *testing.T) {
		mockService := new(MockStudentService)
		h := handler.NewHandler(mockService)
		router := gin.Default()
		h.RegisterRoutes(router.Group("/"))

		students := []model.Student{
			{ID: "1", Name: "Alice", Age: 21, Email: "alice@example.com"},
			{ID: "2", Name: "Bob", Age: 22, Email: "bob@example.com"},
		}

		mockService.On("ListStudents", mock.Anything).Return(students, nil)

		req, _ := http.NewRequest("GET", "/students", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertCalled(t, "ListStudents", mock.Anything)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockService := new(MockStudentService)
		h := handler.NewHandler(mockService)
		router := gin.Default()
		h.RegisterRoutes(router.Group("/"))

		mockService.On("ListStudents", mock.Anything).Return([]model.Student{}, errors.New("Student not found"))

		req, _ := http.NewRequest("GET", "/students", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertCalled(t, "ListStudents", mock.Anything)

		var respBody map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &respBody)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Student not found", respBody["error"])
	})
}
