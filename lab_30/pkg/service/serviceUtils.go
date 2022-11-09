package service

import (
	"encoding/json"
	"fmt"
	u "lab_30/pkg/user"
	ut "lab_30/pkg/utils"
	"log"
	"strconv"
)

const (
	storeJson = "storeJson.txt"
)

func (s *Service) loadStore() {
	s.LoadStoreFromFile(storeJson)
}

func (s *Service) GetAllUsers() map[string]*u.User {
	return s.Store
}

func (s *Service) LoadStoreFromFile(fileName string) {
	content := ut.ReadFromFile(fileName)
	if len(content) == 0 {
		return
	}
	serv := []Service{}
	if err := json.Unmarshal(content, &serv); err != nil {
		log.Printf("service.loadStore: Error: %s", err.Error())
		return
	}
	for _, v := range serv {
		s.IdGen = v.IdGen
		s.Store = v.Store
		s.refreshFriends()
	}
}

func (s *Service) getId() string {
	s.IdGen++
	return strconv.Itoa(s.IdGen)
}

func (s *Service) saveStore() {
	s.saveStoreToFile(storeJson)
}

func (s *Service) saveStoreToFile(fileName string) {
	serv, err := json.Marshal(&s)
	if err != nil {
		log.Printf("service.saveStore: Error: %s", err.Error())
	}
	ut.WriteToFile("["+string(serv)+"]", fileName)
}

func (s *Service) saveUser(user *u.User) {
	log.Println("service.saveUser user:", user.ToString())
	s.Store[user.GetId()] = user
	s.saveStore()
}

func (s *Service) getUser(id string) (*u.User, error) {
	user := s.Store[id]
	if user == nil {
		return nil, fmt.Errorf("user is nil with ID=%v", id)
	}
	return user, nil
}

func (s *Service) createUser() *u.User {
	s.loadStore()
	userIdStr := s.getId()
	userId, _ := strconv.Atoi(userIdStr)
	user := u.NewUser(userId)
	user.SetName("User" + user.GetId())
	s.saveUser(user)
	return user
}

func (s *Service) refreshFriends() {
	for _, user := range s.Store {
		user.InitFriends()
		for _, friendId := range user.GetStrFriends() {
			friend, err := s.getUser(friendId)
			if err != nil {
				log.Println("service refreshFriends: Error ", err.Error())
			}
			user.AddFriend(friend)
		}
	}
}

func (s *Service) makeFriend(dat map[string]interface{}) (*u.User, *u.User, error) {
	sourceId := dat["source_id"].(string)
	targetId := dat["target_id"].(string)

	err := s.MakeFriend(sourceId, targetId)
	if err != nil {
		return nil, nil, err
	}
	sourceUser, err := s.getUser(sourceId)
	if err != nil {
		return nil, nil, err
	}
	targetUser, err := s.getUser(targetId)
	if err != nil {
		return nil, nil, err
	}
	return sourceUser, targetUser, nil
}

func (s *Service) MakeFriend(userId1 string, friendIds ...string) error {
	sourceUser, err := s.getUser(userId1)
	if err != nil {
		return err
	}
	for _, targetId := range friendIds {
		targetUser, err := s.getUser(targetId)
		if err != nil {
			return err
		}
		if userId1 == targetId {
			return fmt.Errorf("service.makeFriend: incorrect friends %d", 1)
		} else {
			if _, err := sourceUser.AddFriend(targetUser); err != nil {
				return err
			}
			if _, err := targetUser.AddFriend(sourceUser); err != nil {
				return err
			}
		}
		s.saveUser(targetUser)
	}
	s.saveUser(sourceUser)
	return nil
}

func (s *Service) deleteUser(targetId string) (string, error) {
	s.loadStore()
	user, err := s.getUser(targetId)
	if err != nil {
		return "", err
	}
	s.unFriend(targetId)
	delete(s.Store, targetId)
	s.saveStore()
	return user.GetName(), nil
}

func (s *Service) unFriend(userId string) error {
	log.Printf("unFriend userId = %s", userId)
	user, err := s.getUser(userId)
	if err != nil {
		return err
	}
	friendIds := user.GetFriendIds()
	for _, friendId := range friendIds {
		friend, err := s.getUser(friendId)
		if err != nil {
			return err
		}
		friend.UnFriend(userId)
	}
	return nil
}

func (s *Service) updateUserAge(user *u.User, age string) {
	user.UpdateAge(age)
	s.saveUser(user)
	s.saveStore()
}
