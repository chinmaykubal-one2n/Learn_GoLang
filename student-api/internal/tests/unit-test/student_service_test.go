package unit_test

import (
	"context"
	"database/sql"
	logging "student-api/internal/logger"
	"student-api/internal/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupLoggerForServiceTests() {
	zapLogger := zap.NewExample()
	logging.Logger = otelzap.New(zapLogger)
}

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

func TestListStudentsService(t *testing.T) {
	setupLoggerForServiceTests()

	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	svc := &service.StudentServiceImpl{DB: db}

	t.Run("returns students on success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "age", "email"}).
			AddRow("123", "Luffy", 19, "luffy@onepiece.com").
			AddRow("456", "Zoro", 21, "zoro@onepiece.com")

		mock.ExpectQuery(`SELECT \* FROM "students"`).
			WillReturnRows(rows)

		ctx := context.Background()
		students, err := svc.ListStudents(ctx, 1, 2)

		require.NoError(t, err)
		assert.Len(t, students, 2)
		assert.Equal(t, "Luffy", students[0].Name)
		assert.Equal(t, "Zoro", students[1].Name)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "students"`).
			WillReturnError(sql.ErrConnDone)

		ctx := context.Background()
		students, err := svc.ListStudents(ctx, 1, 2)

		assert.Error(t, err)
		assert.Empty(t, students)
	})
}
