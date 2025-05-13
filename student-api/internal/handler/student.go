package handler

import (
	"net/http"
	"student-api/internal/model"
	"student-api/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	student := r.Group("/students")
	{
		student.GET("", list)
		student.GET("/:id", get)
		student.POST("", create)
		student.PUT("/:id", update)
		student.DELETE("/:id", deleteStudent)
	}
}

func list(c *gin.Context) {
	students, err := service.ListStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

func get(c *gin.Context) {
	s, err := service.GetStudent(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func create(c *gin.Context) {
	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := service.CreateStudent(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, student)
}

func update(c *gin.Context) {
	var input model.Student
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := service.UpdateStudent(c.Param("id"), input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, student)
}

func deleteStudent(c *gin.Context) {
	err := service.DeleteStudent(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
