package service

import (
	"encoding/json"
	"fmt"
	u "lab_30/pkg/user"
	ut "lab_30/pkg/utils"
	"log"
	"net/http"
	"strconv"
)

type Service struct {
	IdGen int                `json:"id"`
	Store map[string]*u.User `json:"store"`
}

func NewService(id int) *Service {
	return &Service{id, make(map[string]*u.User)}
}

func (s *Service) getId() string {
	s.IdGen++
	return strconv.Itoa(s.IdGen)
}

func (s *Service) loadStore() {
	content := ut.ReadFromFile()
	log.Println("service.loadStore content =", string(content))
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

func (s *Service) saveStore() {
	serv, err := json.Marshal(&s)
	if err != nil {
		log.Printf("service.saveStore: Error: %s", err.Error())
	}
	log.Println("service.saveStore: ", string(serv))
	ut.WriteToFile("[" + string(serv) + "]")
}

func (s *Service) saveUser(user *u.User) {
	log.Println("service.saveUser user =", user.ToString())
	s.Store[user.GetId()] = user
	s.saveStore()
}

func (s *Service) getUser(id string) (*u.User, error) {
	user := s.Store[id]
	if user == nil {
		return user, fmt.Errorf("user is nil")
	}
	return user, nil
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

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	/*
		1. Сделайте обработчик создания пользователя.
	*/
	ut.LogRequest("Create", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		s.loadStore()
		content, shouldReturn := ut.GetContent(r, w)
		if shouldReturn {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("content = %s", string(content))
		userIdStr := s.getId()
		userId, _ := strconv.Atoi(userIdStr)
		user := u.NewUser(userId)
		if shouldReturn1 := ut.UnMarshalData(content, user, w); shouldReturn1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.saveUser(user)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(userIdStr))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	ut.LogRequest("GetAll", r)
	response := ""
	s.loadStore()
	for _, user := range s.Store {
		user.RefreshStrFriends()
		log.Printf("user = %s", user.ToString())
		resp, shouldReturn := ut.MarshalData(user, w)
		if shouldReturn {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response += resp
	}
	defer r.Body.Close()

	log.Println("response =", response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	/*
		2. Сделайте обработчик, который делает друзей из двух пользователей.
	*/
	ut.LogRequest("MakeFriends", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, shouldReturn1 := ut.DecodeData(r, w)
		if shouldReturn1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		sourceUser, targetUser, err := s.makeFriend(dat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) makeFriend(dat map[string]interface{}) (*u.User, *u.User, error) {
	sourceId := dat["source_id"].(string)
	targetId := dat["target_id"].(string)
	sourceUser, err := s.getUser(sourceId)
	if err != nil {
		return sourceUser, nil, err
	}
	targetUser, err := s.getUser(targetId)
	if err != nil {
		return sourceUser, targetUser, err
	}
	if sourceId == targetId {
		return sourceUser, targetUser, fmt.Errorf("service.makeFriend: incorrect friends %d", 1)
	} else {
		if targetUser, err := sourceUser.AddFriend(targetUser); err != nil {
			return sourceUser, targetUser, err
		}
		if sourceUser, err := targetUser.AddFriend(sourceUser); err != nil {
			return sourceUser, targetUser, err
		}
	}
	s.saveUser(sourceUser)
	s.saveUser(targetUser)
	return sourceUser, targetUser, nil
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	/*
		3. Сделайте обработчик, который удаляет пользователя.
	*/
	ut.LogRequest("Delete", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, shouldReturn := ut.DecodeData(r, w)
		if shouldReturn {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		userId := dat["target_id"].(string)
		userName, err := s.deleteUser(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userName))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
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

func (s *Service) GetFriendsById(w http.ResponseWriter, r *http.Request) {
	/*
	 4. Сделайте обработчик, который возвращает всех друзей пользователя.
	*/
	ut.LogRequest("GetFriendsById", r)
	s.loadStore()
	defer r.Body.Close()

	userId := ut.GetRequestParam(r, "userid")
	user, err := s.getUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	friendIds := user.GetFriendIds()

	response := ""
	for _, friendId := range friendIds {
		friend, err := s.getUser(friendId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		friend.RefreshStrFriends()
		resp, shouldReturn := ut.MarshalData(friend, w)
		if shouldReturn {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response += resp
	}

	log.Println(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	s.saveStore()
}

func (s *Service) UpdateAgeById(w http.ResponseWriter, r *http.Request) {
	/*
	 5. Сделайте обработчик, который обновляет возраст пользователя.
	*/
	ut.LogRequest("UpdateAgeById", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		s.loadStore()
		dat, shouldReturn := ut.DecodeData(r, w)
		if shouldReturn {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		userId := ut.GetRequestParam(r, "userid")
		user, err := s.getUser(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user.UpdateAge(dat["new age"].(string))
		s.saveStore()
		response := "возраст пользователя успешно обновлён"
		log.Println(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
