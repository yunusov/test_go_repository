package main

import (
	s "lab_30/pkg/service"
	u "lab_30/pkg/utils"
	"log"
	"net/http"
	"github.com/gorilla/mux" //go get github.com/gorilla/mux
)

func main() {
	srv := s.NewService(0)
	router := mux.NewRouter()
	router.HandleFunc("/get", srv.GetAll).Methods(http.MethodGet) 
	router.HandleFunc("/create", srv.Create).Methods(http.MethodPost)
	router.HandleFunc("/", u.Hello)
	router.HandleFunc("/make_friends", srv.MakeFriends).Methods(http.MethodPost)
	router.HandleFunc("/user", srv.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/friends/{userid:[\\d]+}", srv.GetFriendsById).Methods(http.MethodGet)
	router.HandleFunc("/{userid:[\\d]+}", srv.UpdateAgeById).Methods(http.MethodPut)

	http.Handle("/", router)
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
