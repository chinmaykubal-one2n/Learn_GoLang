package service

import (
	"student-api/internal/db"
	"student-api/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateTeacher(t model.Teacher) (model.Teacher, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.Teacher{}, err
	}
	t.Password = string(hashedPassword)

	t.ID = uuid.New().String()

	result := db.DB.Create(&t)
	if result.Error != nil {
		return model.Teacher{}, result.Error
	}
	return t, nil
}
