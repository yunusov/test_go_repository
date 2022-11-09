package user

import (
	"fmt"
	ut "lab_30/pkg/utils"
	"strconv"
)

type User struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
	friends map[string]*User
}

func (u *User) ToString() string {
	keys := make([]string, 0, len(u.friends))
	for k := range u.friends {
		keys = append(keys, k)
	}
	return fmt.Sprintf("Id = %d, Name = %s, Age = %s, Friends = %s, friends = %s\n", u.Id,
		u.Name, u.Age, u.Friends, keys)
}

func (u *User) GetId() string {
	return strconv.Itoa(u.Id)
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) GetStrFriends() []string {
	if u.Friends == nil {
		u.Friends = []string{}
	}
	return u.Friends
}

func (u *User) GetFriendIds() []string {
	keys := make([]string, 0, len(u.friends))
	for k := range u.friends {
		keys = append(keys, k)
	}
	return keys
}

func (u *User) RefreshStrFriends() {
	if u.friends == nil {
		u.friends = make(map[string]*User)
	}
	keys := make([]string, 0, len(u.friends))
	for k := range u.friends {
		keys = append(keys, k)
	}
	u.Friends = keys
}

func NewUser(id int) *User {
	return &User{id, "", "", []string{}, make(map[string]*User)}
}

func (u *User) AddFriend(user *User) (*User, error) {
	if user == nil {
		return nil, fmt.Errorf("User cannot be null %d", 1)
	}
	userId := user.GetId()
	friend := u.friends[user.GetId()]
	if friend == nil && u.GetId() != userId {
		if !ut.SliceContains(u.Friends, userId) {
			u.Friends = append(u.Friends, userId)
		}
		u.friends[userId] = user
		return u, nil
	}
	return nil, fmt.Errorf("already friends %d", 1)
}

func (u *User) InitFriends() {
	if u.friends == nil {
		u.friends = make(map[string]*User)
	}
	if u.Friends == nil {
		u.Friends = []string{}
	}
}

func (u *User) UnFriend(userId string) {
	delete(u.friends, userId)
	u.RefreshStrFriends()
}

func (u *User) UpdateAge(newAge string) {
	u.Age = newAge
}

func (u *User) GetAge() string {
	return u.Age
}
