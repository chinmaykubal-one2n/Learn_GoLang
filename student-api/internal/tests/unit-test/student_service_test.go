package unit_test

import (
	"context"
	"database/sql"
	"errors"
	logging "student-api/internal/logger"
	"student-api/internal/model"
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

		assert.NoError(t, err)
		assert.Len(t, students, 2)
		assert.Equal(t, "Luffy", students[0].Name)
		assert.Equal(t, "Zoro", students[1].Name)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "students"`).
			WillReturnError(errors.New("student not found"))

		ctx := context.Background()
		students, err := svc.ListStudents(ctx, 1, 2)

		assert.Error(t, err)
		assert.Empty(t, students)
	})
}

func TestDeleteStudentService(t *testing.T) {
	setupLoggerForServiceTests()

	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	svc := &service.StudentServiceImpl{DB: db}

	t.Run("deletes student successfully", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "students" WHERE id = \$1`).
			WithArgs("123").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		ctx := context.Background()

		// Perform the delete operation
		err := svc.DeleteStudent("123", ctx)

		assert.NoError(t, err)
	})

	t.Run("returns error when student not found", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "students" WHERE id = \$1`).
			WithArgs("123").
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(errors.New("student not found"))
		mock.ExpectRollback()

		ctx := context.Background()
		err := svc.DeleteStudent("123", ctx)

		assert.Error(t, err)
	})
}

func TestCreateStudentService(t *testing.T) {
	setupLoggerForServiceTests()

	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	svc := &service.StudentServiceImpl{DB: db}

	t.Run("creates student successfully", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "students"`).
			WithArgs(sqlmock.AnyArg(), "Luffy", 19, "luffy@onepiece.com").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		ctx := context.Background()
		student, err := svc.CreateStudent(model.Student{Name: "Luffy", Age: 19, Email: "luffy@onepiece.com"}, ctx)

		assert.NoError(t, err)
		assert.Equal(t, "Luffy", student.Name)
		assert.Equal(t, 19, student.Age)
		assert.Equal(t, "luffy@onepiece.com", student.Email)
		assert.NotEmpty(t, student.ID)
	})

	t.Run("returns error when creation fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "students"`).
			WillReturnError(sql.ErrConnDone).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectRollback()

		ctx := context.Background()
		student, err := svc.CreateStudent(model.Student{Name: "Luffy", Age: 19, Email: "luffy@onepiece.com"}, ctx)

		assert.Error(t, err)
		assert.Empty(t, student)
	})
}

func TestGetStudentService(t *testing.T) {
	setupLoggerForServiceTests()

	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	svc := &service.StudentServiceImpl{DB: db}

	t.Run("returns student on success", func(t *testing.T) {
		studentID := "123"
		rows := sqlmock.NewRows([]string{"id", "name", "age", "email"}).
			AddRow(studentID, "Luffy", 19, "luffy@onepiece.com")

		mock.ExpectQuery(`(?i)^SELECT \* FROM "students" WHERE id = \$1 ORDER BY "students"\."id" LIMIT \$2$`).
			WithArgs(studentID, 1).
			WillReturnRows(rows)

		ctx := context.Background()
		student, err := svc.GetStudent(studentID, ctx)

		assert.NoError(t, err)
		assert.Equal(t, studentID, student.ID)
		assert.Equal(t, "Luffy", student.Name)
		assert.Equal(t, 19, student.Age)
		assert.Equal(t, "luffy@onepiece.com", student.Email)
	})

	t.Run("returns error when student not found", func(t *testing.T) {
		studentID := "999"

		mock.ExpectQuery(`(?i)^SELECT \* FROM "students" WHERE id = \$1 ORDER BY "students"\."id" LIMIT \$2$`).
			WithArgs(studentID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		ctx := context.Background()
		student, err := svc.GetStudent(studentID, ctx)

		assert.Error(t, err)
		assert.EqualError(t, err, "Student not found")
		assert.Empty(t, student.ID)
	})
}

func TestUpdateStudentService(t *testing.T) {
	setupLoggerForServiceTests()

	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	svc := &service.StudentServiceImpl{DB: db}

	t.Run("updates student successfully", func(t *testing.T) {
		studentID := "123"

		rows := sqlmock.NewRows([]string{"id", "name", "age", "email"}).
			AddRow(studentID, "Luffy", 19, "luffy@onepiece.com")

		mock.ExpectQuery(`(?i)^SELECT \* FROM "students" WHERE id = \$1 ORDER BY "students"\."id" LIMIT \$2$`).
			WithArgs(studentID, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "students" SET`).
			WithArgs("Zoro", 21, "zoro@onepiece.com", studentID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		ctx := context.Background()
		updatedStudent, err := svc.UpdateStudent(studentID, model.Student{Name: "Zoro", Age: 21, Email: "zoro@onepiece.com"}, ctx)

		assert.NoError(t, err)
		assert.Equal(t, "Zoro", updatedStudent.Name)
		assert.Equal(t, 21, updatedStudent.Age)
		assert.Equal(t, "zoro@onepiece.com", updatedStudent.ID)
	})

	t.Run("returns error when student not found", func(t *testing.T) {
		studentID := "123"
		mock.ExpectQuery(`(?i)^SELECT \* FROM "students" WHERE id = \$1 ORDER BY "students"\."id" LIMIT \$2$`).
			WithArgs(studentID).
			WillReturnError(sql.ErrNoRows)

		ctx := context.Background()
		updatedStudent, err := svc.UpdateStudent(studentID, model.Student{Name: "Zoro", Age: 21, Email: "zoro@onepiece.com"}, ctx)

		assert.Error(t, err)
		assert.Empty(t, updatedStudent)
	})
}
