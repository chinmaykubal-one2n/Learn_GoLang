package handler

import (
	"net/http"
	"student-api/internal/model"
	"student-api/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterTeacher(c *gin.Context) {
	var input model.Teacher
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	teacher, err := service.CreateTeacher(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, teacher)
}
