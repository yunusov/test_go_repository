package service

import (
	"fmt"
	u "lab_30/pkg/user"
	ut "lab_30/pkg/utils"
	"log"
	"net/http"
)

type Service struct {
	IdGen int                `json:"id"`
	Store map[string]*u.User `json:"store"`
}

func NewService(id int) *Service {
	return &Service{id, make(map[string]*u.User)}
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	/*
		1. Сделайте обработчик создания пользователя.
	*/
	ut.LogRequest("Create", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		content, err := ut.GetContent(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		user := s.createUser()
		if err := ut.UnMarshalData(content, &user, w); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		s.saveUser(user)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(user.GetId()))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	ut.LogRequest("GetAll", r)
	response := ""
	s.loadStore()
	for _, user := range s.GetAllUsers() {
		user.RefreshStrFriends()
		log.Printf("user = %s", user.ToString())
		resp, err := ut.MarshalData(user, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
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
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		sourceUser, targetUser, err := s.makeFriend(dat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	/*
		3. Сделайте обработчик, который удаляет пользователя.
	*/
	ut.LogRequest("Delete", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		userId := dat["target_id"].(string)
		userName, err := s.deleteUser(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userName))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
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
		w.Write([]byte(err.Error()))
		return
	}
	friendIds := user.GetFriendIds()

	response := ""
	for _, friendId := range friendIds {
		friend, err := s.getUser(friendId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		friend.RefreshStrFriends()
		resp, err := ut.MarshalData(friend, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
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
		dat, err := ut.DecodeData(r, w)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		userId := ut.GetRequestParam(r, "userid")
		log.Printf("userId = %v", userId)
		user, err := s.getUser(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		s.updateUserAge(user, dat["new age"].(string))
		response := "возраст пользователя успешно обновлён"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
