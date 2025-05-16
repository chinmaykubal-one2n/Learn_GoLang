package service

import (
	"student-api/internal/db"
	"student-api/internal/model"

	"github.com/google/uuid"
)

func CreateTeacher(t model.Teacher) (model.Teacher, error) {
	t.ID = uuid.New().String()
	result := db.DB.Create(&t)
	if result.Error != nil {
		return model.Teacher{}, result.Error
	}
	return t, nil
}
