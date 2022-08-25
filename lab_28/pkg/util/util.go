package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func GetStudentAttributes(studentStr string) (name string, age, grade int, e error) {
	studentArr := strings.Split(studentStr, " ")
	if len(studentArr) != 3 {
		e = fmt.Errorf("Неправильный формат данных!")
		return
	}
	name = studentArr[0]
	age, e = getIntValue(studentArr[1])
	grade, e = getIntValue(studentArr[2])
	return
}

func EnterStringValue(str string) (result string, isEof bool) {
	fmt.Print(str)
	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadString('\r')
	result = strings.TrimSuffix(result, "\r")
	if err != nil {
		isEof = err == io.EOF
		if !isEof {
			panic(err)
		}
	}
	return
}

func getIntValue(str string) (result int, e error) {
	result, err := strconv.Atoi(str)
	e = err
	return
}