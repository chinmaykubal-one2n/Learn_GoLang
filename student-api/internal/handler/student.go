// Package handler provides HTTP request handlers for the student API
package handler

import (
	"net/http"
	"strconv"
	logging "student-api/internal/logger"
	"student-api/internal/model"
	"student-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	_, span := otel.Tracer("student-handler").Start(c.Request.Context(), "health-check")
	defer span.End()

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
	ctx, span := otel.Tracer("student-handler").Start(c.Request.Context(), "list-students")
	defer span.End()

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	if err != nil || pageSize < 1 || pageSize > 500 {
		pageSize = 100
	}

	logging.Logger.Info("[list-handler]: Received request to list students", zap.Int("page", page), zap.Int("page_size", pageSize))

	students, err := h.service.ListStudents(ctx, page, pageSize)
	if err != nil {
		span.SetStatus(codes.Error, "failed to list students")
		span.SetAttributes(attribute.String("error", err.Error()))
		logging.Logger.Error("[list-handler]: Failed to list students")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.Int("student_count", len(students)))
	span.SetStatus(codes.Ok, "successfully listed students")
	logging.Logger.Info("[list-handler]: Successfully listed students", zap.Int("count", len(students)))

	c.JSON(http.StatusOK, gin.H{
		"students":  students,
		"page":      page,
		"page_size": pageSize,
	})
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
	ctx, span := otel.Tracer("student-handler").Start(c.Request.Context(), "get-student")
	defer span.End()

	studentID := c.Param("id")
	span.SetAttributes(attribute.String("student_id", studentID))

	logging.Logger.Info("[get-handler]: Received request to get a student")

	s, err := h.service.GetStudent(studentID, ctx)
	if err != nil {
		span.SetStatus(codes.Error, "failed to get student")
		span.SetAttributes(attribute.String("error", err.Error()))
		logging.Logger.Error("[get-handler]: Failed to get student", zap.String("student_id", studentID))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, "successfully retrieved student")
	logging.Logger.Info("[get-handler]: Successfully retrieved student", zap.String("student_id", s.ID))

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
	ctx, span := otel.Tracer("student-handler").Start(c.Request.Context(), "create-student")
	defer span.End()

	logging.Logger.Info("[create-handler]: Received request to create a student")

	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		span.SetStatus(codes.Error, "invalid input data")
		span.SetAttributes(attribute.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(
		attribute.String("student_name", input.Name),
		attribute.String("student_email", input.Email),
		attribute.Int("student_age", input.Age),
	)

	student, err := h.service.CreateStudent(input, ctx)
	if err != nil {
		span.SetStatus(codes.Error, "failed to create student")
		span.SetAttributes(attribute.String("error", err.Error()))
		logging.Logger.Error("[create-handler]: Failed to create student")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(attribute.String("created_student_id", student.ID))
	span.SetStatus(codes.Ok, "successfully created student")
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
	ctx, span := otel.Tracer("student-handler").Start(c.Request.Context(), "update-student")
	defer span.End()

	studentID := c.Param("id")
	span.SetAttributes(attribute.String("student_id", studentID))

	logging.Logger.Info("[update-handler]: Received request to update a student")

	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		span.SetStatus(codes.Error, "invalid input data")
		span.SetAttributes(attribute.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	span.SetAttributes(
		attribute.String("updated_name", input.Name),
		attribute.String("updated_email", input.Email),
		attribute.Int("updated_age", input.Age),
	)

	student, err := h.service.UpdateStudent(studentID, input, ctx)
	if err != nil {
		span.SetStatus(codes.Error, "failed to update student")
		span.SetAttributes(attribute.String("error", err.Error()))
		logging.Logger.Error("[update-handler]: Failed to update student", zap.String("student_id", studentID))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, "successfully updated student")
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
	ctx, span := otel.Tracer("student-handler").Start(c.Request.Context(), "delete-student")
	defer span.End()

	studentID := c.Param("id")
	span.SetAttributes(attribute.String("student_id", studentID))

	logging.Logger.Info("[delete-handler]: Received request to delete a student")

	err := h.service.DeleteStudent(studentID, ctx)
	if err != nil {
		span.SetStatus(codes.Error, "failed to delete student")
		span.SetAttributes(attribute.String("error", err.Error()))
		logging.Logger.Error("[delete-handler]: Failed to delete student", zap.String("student_id", studentID))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	span.SetStatus(codes.Ok, "successfully deleted student")
	logging.Logger.Info("[delete-handler]: Successfully deleted student", zap.String("student_id", studentID))

	c.Status(http.StatusNoContent)
}
