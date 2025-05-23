package service

import (
	"context"
	logging "student-api/internal/logger"
	"student-api/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type TeacherService interface {
	CreateTeacher(t model.Teacher, ctx context.Context) (model.Teacher, error)
	GetTeacher(username string, ctx context.Context) (model.Teacher, error)
}

type TeacherServiceImpl struct {
	DB *gorm.DB
}

func (ts *TeacherServiceImpl) CreateTeacher(t model.Teacher, ctx context.Context) (model.Teacher, error) {
	logging.Logger.Info("[create-teacher-service]: Creating new teacher")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.Logger.Error("[create-teacher-service]: Failed to hash password", zap.Error(err))
		return model.Teacher{}, err
	}
	t.Password = string(hashedPassword)

	t.ID = uuid.New().String()

	result := ts.DB.WithContext(ctx).Create(&t)
	if result.Error != nil {
		logging.Logger.Error("[create-teacher-service]: Failed to create teacher",
			zap.String("username", t.Username),
		)
		return model.Teacher{}, result.Error
	}

	logging.Logger.Info("[create-teacher-service]: Successfully created teacher",
		zap.String("teacher_id", t.ID),
		zap.String("username", t.Username))
	return t, nil
}

func (ts *TeacherServiceImpl) GetTeacher(username string, ctx context.Context) (model.Teacher, error) {
	logging.Logger.Info("[get-teacher-service]: Fetching teacher", zap.String("username", username))

	var teacher model.Teacher
	result := ts.DB.WithContext(ctx).Where("username = ?", username).First(&teacher)
	if result.Error != nil {
		logging.Logger.Error("[get-teacher-service]: Failed to fetch teacher",
			zap.String("username", username),
		)
		return model.Teacher{}, result.Error
	}

	logging.Logger.Info("[get-teacher-service]: Successfully fetched teacher",
		zap.String("teacher_id", teacher.ID),
		zap.String("username", teacher.Username))
	return teacher, nil
}
