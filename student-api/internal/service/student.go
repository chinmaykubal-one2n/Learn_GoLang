package service

import (
	"context"
	"errors"
	logging "student-api/internal/logger"
	"student-api/internal/model"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	_, span := otel.Tracer("student-service").Start(ctx, "list-students")
	defer span.End()

	logging.Logger.Ctx(ctx).Info("[list-service]: Fetching students from database")

	var students []model.Student
	result := s.DB.WithContext(ctx).Find(&students)
	if result.Error != nil {
		span.SetStatus(codes.Error, "failed to fetch students")
		span.SetAttributes(attribute.String("error", result.Error.Error()))
		logging.Logger.Error("[list-service]: Failed to fetch students")
		return []model.Student{}, errors.New("Student not found")
	}

	span.SetAttributes(attribute.Int("student_count", len(students)))
	span.SetStatus(codes.Ok, "successfully fetched students")
	logging.Logger.Ctx(ctx).Info("[list-service]: Successfully fetched students", zap.Int("count", len(students)))

	return students, nil
}

func (s *StudentServiceImpl) GetStudent(id string, ctx context.Context) (model.Student, error) {
	_, span := otel.Tracer("student-service").Start(ctx, "get-student")
	defer span.End()
	span.SetAttributes(attribute.String("student_id", id))

	logging.Logger.Ctx(ctx).Info("[get-service]: Fetching student from database")

	var student model.Student
	result := s.DB.WithContext(ctx).First(&student, "id = ?", id)
	if result.Error != nil {
		span.SetStatus(codes.Error, "student not found")
		span.SetAttributes(attribute.String("error", result.Error.Error()))
		logging.Logger.Error("[get-service]: Failed to get student", zap.String("id", id))
		return model.Student{}, errors.New("Student not found")
	}

	span.SetStatus(codes.Ok, "successfully fetched student")
	logging.Logger.Ctx(ctx).Info("[get-service]: Successfully fetched student", zap.String("id", id))

	return student, nil
}

func (s *StudentServiceImpl) CreateStudent(st model.Student, ctx context.Context) (model.Student, error) {
	_, span := otel.Tracer("student-service").Start(ctx, "create-student")
	defer span.End()

	logging.Logger.Ctx(ctx).Info("[create-service]: Creating student in database")

	st.ID = uuid.New().String()
	span.SetAttributes(
		attribute.String("student_id", st.ID),
		attribute.String("student_name", st.Name),
		attribute.String("student_email", st.Email),
	)

	result := s.DB.WithContext(ctx).Create(&st)
	if result.Error != nil {
		span.SetStatus(codes.Error, "failed to create student")
		span.SetAttributes(attribute.String("error", result.Error.Error()))
		logging.Logger.Error("[create-service]: Failed to create student")
		return model.Student{}, result.Error
	}

	span.SetStatus(codes.Ok, "successfully created student")
	logging.Logger.Ctx(ctx).Info("[create-service]: Successfully created student", zap.String("id", st.ID))

	return st, nil
}

func (s *StudentServiceImpl) UpdateStudent(id string, updated model.Student, ctx context.Context) (model.Student, error) {
	_, span := otel.Tracer("student-service").Start(ctx, "update-student")
	defer span.End()
	span.SetAttributes(attribute.String("student_id", id))

	logging.Logger.Ctx(ctx).Info("[update-service]: Updating student in database")

	var student model.Student
	if err := s.DB.WithContext(ctx).First(&student, "id = ?", id).Error; err != nil {
		span.SetStatus(codes.Error, "student not found")
		span.SetAttributes(attribute.String("error", err.Error()))
		logging.Logger.Error("[update-service]: Failed to update student", zap.String("id", id))
		return model.Student{}, errors.New("Student not found")
	}

	student.Name = updated.Name
	student.Age = updated.Age
	student.Email = updated.Email

	span.SetAttributes(
		attribute.String("updated_name", updated.Name),
		attribute.Int("updated_age", updated.Age),
		attribute.String("updated_email", updated.Email),
	)

	if err := s.DB.WithContext(ctx).Save(&student).Error; err != nil {
		span.SetStatus(codes.Error, "failed to save student updates")
		span.SetAttributes(attribute.String("error", err.Error()))
		return model.Student{}, err
	}

	span.SetStatus(codes.Ok, "successfully updated student")
	logging.Logger.Ctx(ctx).Info("[update-service]: Successfully updated student", zap.String("id", id))
	return student, nil
}

func (s *StudentServiceImpl) DeleteStudent(id string, ctx context.Context) error {
	_, span := otel.Tracer("student-service").Start(ctx, "delete-student")
	defer span.End()
	span.SetAttributes(attribute.String("student_id", id))

	logging.Logger.Ctx(ctx).Info("[delete-service]: Deleting student from database")

	result := s.DB.WithContext(ctx).Delete(&model.Student{}, "id = ?", id)
	if result.RowsAffected == 0 {
		span.SetStatus(codes.Error, "student not found")
		logging.Logger.Error("[delete-service]: Failed to delete student", zap.String("id", id))
		return errors.New("Student not found")
	}

	span.SetStatus(codes.Ok, "successfully deleted student")
	span.SetAttributes(attribute.Int64("rows_affected", result.RowsAffected))
	logging.Logger.Ctx(ctx).Info("[delete-service]: Successfully deleted student", zap.String("id", id))
	return nil
}
