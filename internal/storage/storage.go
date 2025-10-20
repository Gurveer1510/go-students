package storage

import "github.com/Gurveer1510/student-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	DeleteById(id int64) error
	UpdateStudent(id int64, name , email string, age int64) error
}