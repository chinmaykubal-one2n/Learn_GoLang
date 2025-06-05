package service_test

import (
	"context"
	"database/sql"
	logging "student-api/internal/logger"
	"student-api/internal/model"
	"student-api/internal/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupLoggerForTeacherServiceTests() {
	zapLogger := zap.NewExample()
	logging.Logger = otelzap.New(zapLogger)
}

func setupTeacherTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
	}
	return gormDB, mock, cleanup
}

func TestCreateTeacherService(t *testing.T) {
	setupLoggerForTeacherServiceTests()
	const (
		username = "Luffy"
		password = "randomPassword"
		email    = "luffy@onepiece.com"
		role     = "admin"
	)

	db, mock, cleanup := setupTeacherTestDB(t)
	defer cleanup()

	svc := &service.TeacherServiceImpl{DB: db}

	t.Run("creates teacher successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "teachers"`).
			WithArgs(sqlmock.AnyArg(), username, sqlmock.AnyArg(), email, role).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		ctx := context.Background()
		teacher, err := svc.CreateTeacher(model.Teacher{
			Username: username,
			Password: password,
			Email:    email,
			Role:     role,
		}, ctx)

		assert.NoError(t, err)
		assert.Equal(t, username, teacher.Username)
		assert.Equal(t, email, teacher.Email)
		assert.Equal(t, role, teacher.Role)
		assert.NotEmpty(t, teacher.ID)

		bcryptErr := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(password))
		assert.NoError(t, bcryptErr)
	})

	t.Run("returns error when creation fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "teachers"`).
			WillReturnError(sql.ErrConnDone).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectRollback()

		ctx := context.Background()
		teacher, err := svc.CreateTeacher(model.Teacher{
			Username: username,
			Password: password,
			Email:    email,
			Role:     role,
		}, ctx)

		assert.Error(t, err)
		assert.Empty(t, teacher)
	})
}
