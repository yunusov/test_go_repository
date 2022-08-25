package student

import (
	"strconv"
)

type Student struct {
	name  string
	age   int
	grade int
}

func (student *Student) ToString() (result string) {
	result = student.name + " " + strconv.Itoa(student.age) + " " + strconv.Itoa(student.grade)
	return
}

func (student *Student) SetName(name string) {
	student.name = name
}

func (student *Student) GetName() (name string) {
	return student.name
}

func New(name string, age, grade int) *Student {
	return &Student{name, age, grade}
}