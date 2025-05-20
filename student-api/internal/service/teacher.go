package service

import (
	"student-api/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TeacherService interface {
	CreateTeacher(t model.Teacher) (model.Teacher, error)
	GetTeacher(username string) (model.Teacher, error)
}

type TeacherServiceImpl struct {
	DB *gorm.DB
}

func (ts *TeacherServiceImpl) CreateTeacher(t model.Teacher) (model.Teacher, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.Teacher{}, err
	}
	t.Password = string(hashedPassword)

	t.ID = uuid.New().String()

	result := ts.DB.Create(&t)
	if result.Error != nil {
		return model.Teacher{}, result.Error
	}
	return t, nil
}

func (ts *TeacherServiceImpl) GetTeacher(username string) (model.Teacher, error) {
	var teacher model.Teacher
	result := ts.DB.Where("username = ?", username).First(&teacher)
	if result.Error != nil {
		return model.Teacher{}, result.Error
	}
	return teacher, nil
}
