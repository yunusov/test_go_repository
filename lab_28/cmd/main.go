package main

import (
	"fmt"
	"lab_28/pkg/storage"
	"lab_28/pkg/student"
	"lab_28/pkg/util"
)

func main() {
	m := storage.New()
	for {
		studentStr, isEof := util.EnterStringValue("Введите студента: ")
		if isEof {
			break
		}
		name, age, grade, e := util.GetStudentAttributes(studentStr)
		if e != nil {
			fmt.Printf("Ошибка ввода студента! Попробуйте ввести данные снова.\r e = %v\r\r", e)
			continue
		}
		m.Put(student.New(name, age, grade))
	}
	m.PrintStudents()
}
