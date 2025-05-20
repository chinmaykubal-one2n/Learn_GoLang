package handler

import (
	"net/http"
	"student-api/internal/model"
	"student-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	Service service.TeacherService
}

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
