package service

import (
	"encoding/json"
	"io"
	u "lab_30/pkg/user"
	ut "lab_30/pkg/utils"
	"log"
	"net/http"
	"strings"
)

type service struct {
	id_gen int
	store  map[string]*u.User
}

func NewService(id int) *service {
	return &service{id, make(map[string]*u.User)}
}

func (s *service) getId() int {
	s.id_gen++
	return s.id_gen
}

func (s *service) Create(w http.ResponseWriter, r *http.Request) {
	/*
		1. Сделайте обработчик создания пользователя.
	*/
	ut.LogRequest("Create", r)
	if r.Method == http.MethodPost &&
		strings.ContainsAny(r.Header.Get("Content-Type"), "application/json") {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Printf("Error: %s", err.Error())
			return
		}
		defer r.Body.Close()

		log.Printf("content = %s", string(content))
		user := u.NewUser(s.getId())
		if err := json.Unmarshal(content, &user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Printf("Error: %s", err.Error())
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

func (s *service) GetAll(w http.ResponseWriter, r *http.Request) {
	ut.LogRequest("GetAll", r)
	if r.Method == "GET" {
		response := ""
		for _, user := range s.store {
			log.Printf("user = %s", user.ToString())
			encUser, err := json.Marshal(user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				log.Printf("Error: %s", err.Error())
				return
			}
			response += string(encUser) + "\n"
		}
		log.Println(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	/*
		2. Сделайте обработчик, который делает друзей из двух пользователей.
	*/
	ut.LogRequest("MakeFriends", r)

}
