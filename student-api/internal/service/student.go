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
	logging.Logger.Info("[list-service]: Fetching students from database")

	var students []model.Student
	result := s.DB.WithContext(ctx).Find(&students)
	// result := s.DB.Find(&students)
	if result.Error != nil {
		logging.Logger.Error("[list-service]: Failed to fetch students")
		return []model.Student{}, errors.New("Student not found")
	}
	logging.Logger.Info("[list-service]: Successfully fetched students", zap.Int("count", len(students)))

	return students, nil
}

func (s *StudentServiceImpl) GetStudent(id string, ctx context.Context) (model.Student, error) {
	logging.Logger.Info("[get-service]: Fetching student from database")

	var student model.Student
	result := s.DB.WithContext(ctx).First(&student, "id = ?", id)
	if result.Error != nil {
		logging.Logger.Error("[get-service]: Failed to get student", zap.String("id", id))
		return model.Student{}, errors.New("Student not found")
	}
	logging.Logger.Info("[get-service]: Successfully fetched student", zap.String("id", id))

	return student, nil
}

func (s *StudentServiceImpl) CreateStudent(st model.Student, ctx context.Context) (model.Student, error) {
	logging.Logger.Info("[create-service]: Creating student in database")

	st.ID = uuid.New().String()
	result := s.DB.WithContext(ctx).Create(&st)
	if result.Error != nil {
		logging.Logger.Error("[create-service]: Failed to create student")
		return model.Student{}, result.Error
	}
	logging.Logger.Info("[create-service]: Successfully created student", zap.String("id", st.ID))

	return st, nil
}

func (s *StudentServiceImpl) UpdateStudent(id string, updated model.Student, ctx context.Context) (model.Student, error) {
	logging.Logger.Info("[update-service]: Updating student in database")

	var student model.Student
	if err := s.DB.WithContext(ctx).First(&student, "id = ?", id).Error; err != nil {
		logging.Logger.Error("[update-service]: Failed to update student", zap.String("id", id))
		return model.Student{}, errors.New("Student not found")
	}
	student.Name = updated.Name
	student.Age = updated.Age
	student.Email = updated.Email
	s.DB.WithContext(ctx).Save(&student)

	logging.Logger.Info("[update-service]: Successfully updated student", zap.String("id", id))
	return student, nil
}

func (s *StudentServiceImpl) DeleteStudent(id string, ctx context.Context) error {
	logging.Logger.Info("[delete-service]: Deleting student from database")

	result := s.DB.WithContext(ctx).Delete(&model.Student{}, "id = ?", id)
	if result.RowsAffected == 0 {
		logging.Logger.Error("[delete-service]: Failed to delete student", zap.String("id", id))
		return errors.New("Student not found")
	}

	logging.Logger.Info("[delete-service]: Successfully deleted student", zap.String("id", id))
	return nil
}
