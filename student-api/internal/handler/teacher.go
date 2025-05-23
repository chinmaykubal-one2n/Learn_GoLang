// Package handler provides HTTP request handlers for the teacher API
package handler

import (
	"net/http"
	logging "student-api/internal/logger"
	"student-api/internal/model"
	"student-api/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	logging.Logger.Info("[register-teacher-handler]: Received request to register new teacher")

	var input model.Teacher
	if err := c.ShouldBindJSON(&input); err != nil {
		logging.Logger.Error("[register-teacher-handler]: Invalid input data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher, err := ts.Service.CreateTeacher(input, c.Request.Context())
	if err != nil {
		logging.Logger.Error("[register-teacher-handler]: Failed to create teacher",
			zap.String("username", input.Username),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logging.Logger.Info("[register-teacher-handler]: Successfully registered teacher",
		zap.String("teacher_id", teacher.ID),
		zap.String("username", teacher.Username))
	c.JSON(http.StatusCreated, teacher)
}
