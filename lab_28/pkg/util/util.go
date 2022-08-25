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
		e = fmt.Errorf("неправильный формат данных!%v", "")
		return
	}
	name = studentArr[0]
	age, e = strconv.Atoi(studentArr[1])
	if e != nil {
		return
	}
	grade, e = strconv.Atoi(studentArr[2])
	if e != nil {
		return
	}
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
