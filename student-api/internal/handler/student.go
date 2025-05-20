// Package handler provides HTTP request handlers for the student API
package handler

import (
	"net/http"
	"student-api/internal/model"
	"student-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
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
	students, err := h.service.ListStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	s, err := h.service.GetStudent(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
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
	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := h.service.CreateStudent(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := h.service.UpdateStudent(c.Param("id"), input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
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
	err := h.service.DeleteStudent(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
