package service

import (
	"errors"
	"student-api/internal/db"
	"student-api/internal/model"

	"github.com/google/uuid"
)

func ListStudents() ([]model.Student, error) {
	var students []model.Student
	result := db.DB.Find(&students)
	if result.Error != nil {
		return []model.Student{}, errors.New("student not found")
	}
	return students, nil
}

func GetStudent(id string) (model.Student, error) {
	var student model.Student
	result := db.DB.First(&student, "id = ?", id)
	if result.Error != nil {
		return model.Student{}, errors.New("student not found")
	}
	return student, nil
}

func CreateStudent(s model.Student) (model.Student, error) {
	s.ID = uuid.New().String()
	result := db.DB.Create(&s)
	if result.Error != nil {
		return model.Student{}, result.Error
	}
	return s, nil
}

func UpdateStudent(id string, updated model.Student) (model.Student, error) {
	var student model.Student
	if err := db.DB.First(&student, "id = ?", id).Error; err != nil {
		return model.Student{}, errors.New("student not found")
	}

	student.Name = updated.Name
	student.Age = updated.Age
	student.Email = updated.Email
	db.DB.Save(&student)
	return student, nil
}

func DeleteStudent(id string) error {
	result := db.DB.Delete(&model.Student{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return errors.New("student not found")
	}
	return nil
}
