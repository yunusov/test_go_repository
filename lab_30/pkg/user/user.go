package user

import (
	"fmt"
	"strconv"
)

type User struct {
	id 			int
	Name		string	`json:"name"`
	Age  		int		`json:"age"`
	Friends []	string 	`json:"friends"`
}

func(u*User) ToString()string {
	return fmt.Sprintf("ID = %d, Name = %s, Age = %d, friends = %s\n", u.id, u.Name, u.Age, u.Friends)
}

func(u*User) GetId()string {
	return strconv.Itoa(u.id)
}

func NewUser(id int) *User {
	return &User{id, "", 0, []string{}}
}