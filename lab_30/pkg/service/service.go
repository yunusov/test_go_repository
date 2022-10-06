package service

import (
	u "lab_30/pkg/user"
	ut "lab_30/pkg/utils"
	"fmt"
	"log"
	"net/http"
)

type Service struct {
	idGen int
	store map[string]*u.User
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
		user.Name += userId
		s.store[userId] = user
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(userId))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	ut.LogRequest("GetAll", r)
	if r.Method == http.MethodGet {
		response := ""
		for _, user := range s.store {
			log.Printf("user = %s", user.ToString())
			resp, shouldReturn := ut.MarshalData(user, w)
			if shouldReturn {
				return
			}
			response += resp
		}
		defer r.Body.Close()

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
		dat, shouldReturn1 := s.unMarshalData(r, w)
		if shouldReturn1 {
			return
		}
		defer r.Body.Close()
		sourceUser, targetUser, err := s.makeFriend(dat)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name)))
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (*Service) unMarshalData(r *http.Request, w http.ResponseWriter) (map[string]interface{}, bool) {
	content, shouldReturn := ut.GetContent(r, w)
	if shouldReturn {
		return nil, true
	}
	defer r.Body.Close()

	log.Printf("content = %s", string(content))
	var dat map[string]interface{}
	if shouldReturn1 := ut.UnMarshalData(content, &dat, w); shouldReturn1 {
		return nil, true
	}
	return dat, false
}

func (s *Service) makeFriend(dat map[string]interface{}) (*u.User, *u.User, error) {
	sourceId := dat["source_id"].(string)
	targetId := dat["target_id"].(string)
	sourceUser := s.getUser(sourceId)
	targetUser := s.getUser(targetId)
	if sourceUser != nil && targetUser != nil && sourceId != targetId {
		if err := sourceUser.AddFriend(targetId); err != nil {
			return sourceUser, targetUser, err
		}
		if err := targetUser.AddFriend(sourceId); err != nil {
			return sourceUser, targetUser, err
		}
	} else {
		return sourceUser, targetUser, fmt.Errorf("unexisted users %d", 1)
	}
	return sourceUser, targetUser, nil
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	/*
		3. Сделайте обработчик, который удаляет пользователя.
	*/
	ut.LogRequest("Delete", r)
	if ut.IsDeleteAndCtJson(r.Method, r.Header.Get("Content-Type")) {
		dat, shouldReturn := s.unMarshalData(r, w)
		if shouldReturn {
			return
		}
		defer r.Body.Close()
		userName := s.deleteUser(dat)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userName))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) deleteUser(dat map[string]interface{}) (userName string) {
	targetId := dat["target_id"].(string)
	if s.getUser(targetId) != nil {
		userName = s.unFriend(targetId)
		delete(s.store, targetId)
	}
	return
}

func (s *Service) unFriend(userId string) string {
	log.Printf("unFriend userId = %s", userId)
	user := s.getUser(userId)
	friendIds := user.GetFriendIds()
	for _, friendId := range friendIds {
		friend := s.getUser(friendId)
		friend.UnFriend(userId)
	}
	return user.GetName()
}

func (s *Service) GetFriendsById(w http.ResponseWriter, r *http.Request) {
	ut.LogRequest("GetFriendsById", r)
	/*r.get
	if r.Method == http.MethodGet {
		response := ""
		for _, user := range s.store {
			log.Printf("user = %s", user.ToString())
			resp, shouldReturn := ut.MarshalData(user, w)
			if shouldReturn {
				return
			}
			response += resp
		}

		log.Println(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)*/
}