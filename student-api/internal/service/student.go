package service

import (
	"context"
	"errors"
	logging "student-api/internal/logger"
	"student-api/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
	logging.Logger.Info("Fetching students from database")

	var students []model.Student
	result := s.DB.WithContext(ctx).Find(&students)
	// result := s.DB.Find(&students)
	if result.Error != nil {
		return []model.Student{}, errors.New("Student not found")
	}
	logging.Logger.Info("Successfully fetched students", zap.Int("count", len(students)))

	return students, nil
}

func (s *StudentServiceImpl) GetStudent(id string, ctx context.Context) (model.Student, error) {
	logging.Logger.Info("Fetching student from database")

	var student model.Student
	result := s.DB.WithContext(ctx).First(&student, "id = ?", id)
	if result.Error != nil {
		return model.Student{}, errors.New("Student not found")
	}
	logging.Logger.Info("Successfully fetched student", zap.String("id", id))

	return student, nil
}

func (s *StudentServiceImpl) CreateStudent(st model.Student, ctx context.Context) (model.Student, error) {
	logging.Logger.Info("Creating student in database")

	st.ID = uuid.New().String()
	result := s.DB.WithContext(ctx).Create(&st)
	if result.Error != nil {
		return model.Student{}, result.Error
	}
	logging.Logger.Info("Successfully created student", zap.String("id", st.ID))

	return st, nil
}

func (s *StudentServiceImpl) UpdateStudent(id string, updated model.Student, ctx context.Context) (model.Student, error) {
	logging.Logger.Info("Updating student in database")

	var student model.Student
	if err := s.DB.WithContext(ctx).First(&student, "id = ?", id).Error; err != nil {
		return model.Student{}, errors.New("Student not found")
	}
	student.Name = updated.Name
	student.Age = updated.Age
	student.Email = updated.Email
	s.DB.WithContext(ctx).Save(&student)

	logging.Logger.Info("Successfully updated student", zap.String("id", id))
	return student, nil
}

func (s *StudentServiceImpl) DeleteStudent(id string, ctx context.Context) error {
	logging.Logger.Info("Deleting student from database")

	result := s.DB.WithContext(ctx).Delete(&model.Student{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return errors.New("Student not found")
	}

	logging.Logger.Info("Successfully deleted student", zap.String("id", id))
	return nil
}
