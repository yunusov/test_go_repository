package main

import (
	s "lab_30/pkg/service"
	u "lab_30/pkg/utils"
	"log"
	"net/http"
	//"go-chi/chi"
)

func main() {
	mux := http.NewServeMux()
	srv := s.NewService(0)

	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetAll)
	mux.HandleFunc("/", u.Hello)
	mux.HandleFunc("/make_friends", srv.MakeFriends)
	mux.HandleFunc("/user", srv.Delete)
	mux.HandleFunc("/friends/{user_id}", srv.GetFriendsById)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
