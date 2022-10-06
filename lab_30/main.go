package main

import (
	//"fmt"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"strings"
)

type User struct {
	id 			int
	Name		string	`json:"name"`
	Age  		int		`json:"age"`
	Friends []	string 	`json:"friends"`
}

func(u*User) toString()string {
	return fmt.Sprintf("ID = %d, Name = %s, Age = %d, friends = %s\n", u.id, u.Name, u.Age, u.Friends)
}

func newUser(id int) *User {
	return &User{id, "", 0, []string{}}
}

type service struct {
	id_gen	int
	store	map[string]*User
}

func newService(id int) *service {
	return &service{id, make(map[string]*User)}
}

func (s *service) getId() int{
	s.id_gen++
	return s.id_gen
}

func main() {
	mux := http.NewServeMux()
	srv := newService(0)

	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetAll)
	mux.HandleFunc("/", Hello)
	//mux.HandleFunc("/make_friens", srv.MakeFriends)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func logRequest(name string, r * http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n", 
		name, r.Method, r.Body, r.Header.Get("Content-Type"))
}

func (s *service) Create(w http.ResponseWriter, r * http.Request) {
	/*
		1. Сделайте обработчик создания пользователя.
	*/
	logRequest("Create", r)
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
		u := newUser(s.getId())
		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Printf("Error: %s", err.Error())
			return
		}
		userId := strconv.Itoa(u.id)
		s.store[userId] = u	
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(userId))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *service) GetAll(w http.ResponseWriter, r * http.Request) {
	logRequest("GetAll", r)
	if r.Method == "GET" {
		response := ""
		for _, user := range s.store {
			log.Printf("user = %s", user.toString())
			json.Marshal(user)
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