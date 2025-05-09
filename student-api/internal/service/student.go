package service

import (
	"errors"
	"student-api/internal/model"

	"github.com/google/uuid"
)

var students = map[string]model.Student{}

func ListStudents() []model.Student {
	result := []model.Student{}
	for _, s := range students {
		result = append(result, s)
	}
	return result
}

func GetStudent(id string) (model.Student, error) {
	s, ok := students[id]
	if !ok {
		return model.Student{}, errors.New("student not found")
	}
	return s, nil
}

func CreateStudent(s model.Student) model.Student {
	s.ID = uuid.New().String()
	students[s.ID] = s
	return s
}

func UpdateStudent(id string, updated model.Student) (model.Student, error) {
	_, exists := students[id]
	if !exists {
		return model.Student{}, errors.New("student not found")
	}
	updated.ID = id
	students[id] = updated
	return updated, nil
}

func DeleteStudent(id string) error {
	_, exists := students[id]
	if !exists {
		return errors.New("student not found")
	}
	delete(students, id)
	return nil
}
