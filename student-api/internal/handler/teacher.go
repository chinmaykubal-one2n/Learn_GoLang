// Package handler provides HTTP request handlers for the teacher API
package handler

import (
	"net/http"
	"student-api/internal/model"
	"student-api/internal/service"

	"github.com/gin-gonic/gin"
)

// TeacherHandler manages HTTP requests for teacher operations
type TeacherHandler struct {
	Service service.TeacherService
}

// @Summary Register a new teacher
// @Description Register a new teacher with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param teacher body model.Teacher true "Teacher registration details"
// @Success 201 {object} model.Teacher "Successfully registered teacher"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /register [post]
func (ts *TeacherHandler) RegisterTeacher(c *gin.Context) {
	var input model.Teacher
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	teacher, err := ts.Service.CreateTeacher(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, teacher)
}
