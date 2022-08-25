package storage

import (
	"fmt"
	s "lab_28/pkg/student"
)

type Storage map[string]*s.Student

func (storage *Storage) Put(student *s.Student) {
	temp := *storage
	temp[student.GetName()] = student
}

func (storage *Storage) Get(name string) (result *s.Student) {
	temp := *storage
	result = temp[name]
	return
}

func (storage *Storage) PrintStudents() {
	fmt.Println("\n\nСтуденты из хранилища:")
	for _, val := range *storage {
		fmt.Println(val.ToString())
	}
}

func New() *Storage {
	return &Storage{}
}
