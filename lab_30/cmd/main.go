package main

import (
	"flag"
	s "lab_30/pkg/service"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

// go run main.go -port "<port_number>"

func main() {
	port := getParamsPort()
	service := s.NewService(0)
	router := mux.NewRouter()

	makeHandleFuncs(router, service)

	log.Printf("Запуск веб-сервера на http://127.0.0.1:%v\n", port)
	err := http.ListenAndServe(":" + port, nil)
	log.Fatal(err)
}

func makeHandleFuncs(router *mux.Router, service *s.Service) {
	router.HandleFunc("/get", service.GetAll).Methods(http.MethodGet)
	router.HandleFunc("/create", service.Create).Methods(http.MethodPost)
	router.HandleFunc("/", hello)
	router.HandleFunc("/make_friends", service.MakeFriends).Methods(http.MethodPost)
	router.HandleFunc("/user", service.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/friends/{userid:[\\d]+}", service.GetFriendsById).Methods(http.MethodGet)
	router.HandleFunc("/{userid:[\\d]+}", service.UpdateAgeById).Methods(http.MethodPut)
	http.Handle("/", router)
}

func getParamsPort() (port string) {
	flag.StringVar(&port, "port", "", "set port")
	flag.Parse()
	return port
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n",
		"hello", r.Method, r.Body, r.Header.Get("Content-Type"))
	w.Write([]byte("Hello"))
}
