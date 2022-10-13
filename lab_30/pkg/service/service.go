package service

import (
	"fmt"
	u "lab_30/pkg/user"
	ut "lab_30/pkg/utils"
	"log"
	"net/http"
	"strconv"
)

type Service struct {
	idGen int
	store map[string]*u.User
}

func NewService(id int) *Service {
	return &Service{id, make(map[string]*u.User)}
}

func (s *Service) getId() string {
	s.idGen++
	return strconv.Itoa(s.idGen)
}

func (s *Service) getUser(id string) (*u.User, error) {
	user := s.store[id]
	if user == nil {
		return user, fmt.Errorf("user is nil")
	}
	return user, nil
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	/*
		1. Сделайте обработчик создания пользователя.
	*/
	ut.LogRequest("Create", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
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
		user.Name += userIdStr
		s.store[userIdStr] = user
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(userIdStr))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	ut.LogRequest("GetAll", r)
	response := ""
	for _, user := range s.store {
		user.RefreshFriends()
		log.Printf("user = %s", user.ToString())
		resp, shouldReturn := ut.MarshalData(user, w)
		if shouldReturn {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response += resp
	}
	defer r.Body.Close()

	log.Println(response)
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
		return sourceUser, targetUser, fmt.Errorf("incorrect friends %d", 1)
	} else {
		if err := sourceUser.AddFriend(targetUser); err != nil {
			return sourceUser, targetUser, err
		}
		if err := targetUser.AddFriend(sourceUser); err != nil {
			return sourceUser, targetUser, err
		}
	}
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
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) deleteUser(targetId string) (string, error) {
	user, err := s.getUser(targetId)
	if err != nil {
		return "", err
	}
	s.unFriend(targetId)
	delete(s.store, targetId)
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
		friend.RefreshFriends()
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
}

func (s *Service) UpdateAgeById(w http.ResponseWriter, r *http.Request) {
	/*
		 5. Сделайте обработчик, который обновляет возраст пользователя.
	*/
	ut.LogRequest("UpdateAgeById", r)
	if ut.IsCtJson(r.Header.Get("Content-Type")) {
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
		response := "возраст пользователя успешно обновлён"
		log.Println(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
	w.WriteHeader(http.StatusBadRequest)
}
