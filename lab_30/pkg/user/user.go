package user

import (
	ut "lab_30/pkg/utils"
	"fmt"
	"strconv"
)

type User struct {
	id      int
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

func (u *User) ToString() string {
	return fmt.Sprintf("ID = %d, Name = %s, Age = %d, friends = %s\n", u.id,
		u.Name, u.Age, u.Friends)
}

func (u *User) GetId() string {
	return strconv.Itoa(u.id)
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetFriendIds() []string {
	return u.Friends
}

func NewUser(id int) *User {
	return &User{id, "", 0, []string{}}
}

func (u *User) AddFriend(userId string) error {
	if !ut.Contains(u.GetFriendIds(), userId) && u.GetId() != userId {
		u.Friends = append(u.Friends, userId)
		return nil
	}
	return fmt.Errorf("already friends %d", 1)
}

func (u *User) UnFriend(userId string) {
	i := ut.Find(u.Friends, userId)
	u.Friends[i] = u.Friends[len(u.Friends)-1]
	u.Friends[len(u.Friends)-1] = ""
	u.Friends = u.Friends[:len(u.Friends)-1]
}
