package handler

import (
	"net/http"
	"student-api/internal/model"
	"student-api/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.StudentService
}

func NewHandler(s service.StudentService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/students", h.list)
	r.GET("/students/:id", h.get)
	r.POST("/students", h.create)
	r.PUT("/students/:id", h.update)
	r.DELETE("/students/:id", h.deleteStudent)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *Handler) list(c *gin.Context) {
	students, err := h.service.ListStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *Handler) get(c *gin.Context) {
	s, err := h.service.GetStudent(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

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

func (h *Handler) deleteStudent(c *gin.Context) {
	err := h.service.DeleteStudent(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
