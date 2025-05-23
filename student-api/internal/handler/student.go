// Package handler provides HTTP request handlers for the student API
package handler

import (
	"net/http"
	logging "student-api/internal/logger"
	"student-api/internal/model"
	"student-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler manages HTTP requests for student operations
type Handler struct {
	service service.StudentService
}

// NewHandler creates a new handler with the given student service
func NewHandler(s service.StudentService) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up the routing for the student API
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/students", h.list)
	r.GET("/students/:id", h.get)
	r.POST("/students", h.create)
	r.PUT("/students/:id", h.update)
	r.DELETE("/students/:id", h.deleteStudent)
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Get the health status of the API
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	logging.Logger.Info("[health-check-handler]: Received health check request")
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// @Summary List all students
// @Description Get a list of all students
// @Tags students
// @Produce json
// @Success 200 {array} model.Student
// @Failure 500 {object} map[string]string
// @Router /api/students [get]
func (h *Handler) list(c *gin.Context) {
	logging.Logger.Info("[list-handler]: Received request to list students")

	students, err := h.service.ListStudents(c.Request.Context())
	if err != nil {
		logging.Logger.Error("[list-handler]: Failed to list students")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logging.Logger.Info("[list-handler]: Successfully listed students", zap.Int("count", len(students)))

	c.JSON(http.StatusOK, students)
}

// get godoc
// @Summary Get a student by ID
// @Description Get a student's details by their ID
// @Tags students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} model.Student
// @Failure 404 {object} map[string]string
// @Router /api/students/{id} [get]
func (h *Handler) get(c *gin.Context) {
	logging.Logger.Info("[get-handler]: Received request to get a student")

	s, err := h.service.GetStudent(c.Param("id"), c.Request.Context())
	if err != nil {
		logging.Logger.Error("[get-handler]: Failed to get student", zap.String("student_id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	logging.Logger.Info("[get-handler]: Successfully retrived student", zap.String("student_id", s.ID))

	c.JSON(http.StatusOK, s)
}

// create godoc
// @Summary Create a new student
// @Description Create a new student with the provided details
// @Tags students
// @Accept json
// @Produce json
// @Param student body model.Student true "Student details"
// @Success 201 {object} model.Student
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/students [post]
func (h *Handler) create(c *gin.Context) {
	logging.Logger.Info("[create-handler]: Received request to create a student")

	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := h.service.CreateStudent(input, c.Request.Context())
	if err != nil {
		logging.Logger.Error("[create-handler]: Failed to create student", zap.String("student_id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logging.Logger.Info("[create-handler]: Successfully created student", zap.String("student_id", student.ID))

	c.JSON(http.StatusCreated, student)
}

// update godoc
// @Summary Update a student
// @Description Update a student's details by their ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param student body model.Student true "Updated student details"
// @Success 200 {object} model.Student
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/students/{id} [put]
func (h *Handler) update(c *gin.Context) {
	logging.Logger.Info("[update-handler]: Received request to update a student")
	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := h.service.UpdateStudent(c.Param("id"), input, c.Request.Context())
	if err != nil {
		logging.Logger.Error("[update-handler]: Failed to update student", zap.String("student_id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	logging.Logger.Info("[update-handler]: Successfully updated student", zap.String("student_id", student.ID))

	c.JSON(http.StatusOK, student)
}

// deleteStudent godoc
// @Summary Delete a student
// @Description Delete a student by their ID
// @Tags students
// @Produce json
// @Param id path string true "Student ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /api/students/{id} [delete]
func (h *Handler) deleteStudent(c *gin.Context) {
	logging.Logger.Info("[delete-handler]:  Received request to delete a student")

	err := h.service.DeleteStudent(c.Param("id"), c.Request.Context())
	if err != nil {
		logging.Logger.Error("[delete-handler]: Failed to delete student", zap.String("student_id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	logging.Logger.Info("[update-handler]: Successfully deleted student", zap.String("student_id", c.Param("id")))

	c.Status(http.StatusNoContent)
}
