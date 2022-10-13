package user

import (
	"fmt"
	"strconv"
)

type User struct {
	id      int				 
	Name    string   		 `json:"name"`
	Age     string   		 `json:"age"`
	Friends []string 		 `json:"friends"`
	friends map[string]*User 
}

func (u *User) ToString() string {
	keys := make([]string, 0, len(u.friends))
	for k := range u.friends {
		keys = append(keys, k)
	}	
	return fmt.Sprintf("ID = %d, Name = %s, Age = %s, Friends = %s\n", u.id,
		u.Name, u.Age, keys)
}

func (u *User) GetId() string {
	return strconv.Itoa(u.id)
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetFriendIds() []string {
	keys := make([]string, 0, len(u.friends))
	for k := range u.friends {
		keys = append(keys, k)
	}
	return keys
}

func (u *User) RefreshFriends() {
	keys := make([]string, 0, len(u.friends))
	for k := range u.friends {
		keys = append(keys, k)
	}
	u.Friends = keys
}

func NewUser(id int) *User {
	return &User{id, "", "", []string{}, make(map[string]*User)}
}

func (u *User) AddFriend(user *User) error {
	if user == nil {
		return fmt.Errorf("User cannot be null %d", 1)
	}
	userId := user.GetId()
	friend := u.friends[user.GetId()]
	if friend == nil && u.GetId() != userId {
		u.friends[userId] = user
		return nil
	} 
	return fmt.Errorf("already friends %d", 1)	
}

func (u *User) UnFriend(userId string) {
	delete(u.friends, userId)
}

func (u *User) UpdateAge(newAge string) {
	u.Age = newAge
}