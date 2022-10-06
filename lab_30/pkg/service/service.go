package service

import (
	"fmt"
	u "lab_30/pkg/user"
	ut "lab_30/pkg/utils"

	//"encoding/json"
	"log"
	"net/http"
)

type Service struct {
	idGen int
	store  map[string]*u.User
}

func NewService(id int) *Service {
	return &Service{id, make(map[string]*u.User)}
}

func (s *Service) getId() int {
	s.idGen++
	return s.idGen
}

func (s *Service) getUser(id string) *u.User {
	return s.store[id]
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	/*
		1. Сделайте обработчик создания пользователя.
	*/
	ut.LogRequest("Create", r)
	if ut.IsPostAndCtJson(r.Method, r.Header.Get("Content-Type")) {
		content, shouldReturn := ut.GetContent(r, w)
		if shouldReturn {
			return
		}
		defer r.Body.Close()

		log.Printf("content = %s", string(content))
		user := u.NewUser(s.getId())
		if shouldReturn1 := ut.UnMarshalData(content, user, w); shouldReturn1 {
			return
		}
		userId := user.GetId()
		s.store[userId] = user
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(userId))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	ut.LogRequest("GetAll", r)
	if r.Method == "GET" {
		response, shouldReturn := ut.MarshalData(s.store, w)
		if shouldReturn {
			return
		}
		log.Println(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	/*
		2. Сделайте обработчик, который делает друзей из двух пользователей.
	*/
	ut.LogRequest("MakeFriends", r)
	if ut.IsPostAndCtJson(r.Method, r.Header.Get("Content-Type")) {
		content, shouldReturn := ut.GetContent(r, w)
		if shouldReturn {
			return
		}
		defer r.Body.Close()

		log.Printf("content = %s", string(content))
		var dat map[string]interface{}
		if shouldReturn1 := ut.UnMarshalData(content, &dat, w); shouldReturn1 {
			return
		}
		sourceUser, targetUser := s.makeFriend(dat)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name)))
	}
}

func (s *Service) makeFriend(dat map[string]interface{}) (sourceUser *u.User, targetUser *u.User) {
	sourceId := dat["source_id"].(string)
	targetId := dat["target_id"].(string)
	sourceUser = s.getUser(sourceId)
	targetUser = s.getUser(targetId)
	sourceUser.MakeFriend(targetId)
	targetUser.MakeFriend(sourceId)
	return
}
