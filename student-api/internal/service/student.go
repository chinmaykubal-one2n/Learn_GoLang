package service

import (
	"context"
	"errors"
	"student-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentService interface {
	ListStudents(ctx context.Context) ([]model.Student, error)
	GetStudent(id string, ctx context.Context) (model.Student, error)
	CreateStudent(s model.Student, ctx context.Context) (model.Student, error)
	UpdateStudent(id string, updated model.Student, ctx context.Context) (model.Student, error)
	DeleteStudent(id string, ctx context.Context) error
}

type StudentServiceImpl struct {
	DB *gorm.DB
}

func (s *StudentServiceImpl) ListStudents(ctx context.Context) ([]model.Student, error) {
	var students []model.Student
	result := s.DB.WithContext(ctx).Find(&students)
	// result := s.DB.Find(&students)
	if result.Error != nil {
		return []model.Student{}, errors.New("Student not found")
	}
	return students, nil
}

func (s *StudentServiceImpl) GetStudent(id string, ctx context.Context) (model.Student, error) {
	var student model.Student
	result := s.DB.WithContext(ctx).First(&student, "id = ?", id)
	if result.Error != nil {
		return model.Student{}, errors.New("Student not found")
	}
	return student, nil
}

func (s *StudentServiceImpl) CreateStudent(st model.Student, ctx context.Context) (model.Student, error) {
	st.ID = uuid.New().String()
	result := s.DB.WithContext(ctx).Create(&st)
	if result.Error != nil {
		return model.Student{}, result.Error
	}
	return st, nil
}

func (s *StudentServiceImpl) UpdateStudent(id string, updated model.Student, ctx context.Context) (model.Student, error) {
	var student model.Student
	if err := s.DB.WithContext(ctx).First(&student, "id = ?", id).Error; err != nil {
		return model.Student{}, errors.New("Student not found")
	}
	student.Name = updated.Name
	student.Age = updated.Age
	student.Email = updated.Email
	s.DB.WithContext(ctx).Save(&student)
	return student, nil
}

func (s *StudentServiceImpl) DeleteStudent(id string, ctx context.Context) error {
	result := s.DB.WithContext(ctx).Delete(&model.Student{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return errors.New("Student not found")
	}
	return nil
}
