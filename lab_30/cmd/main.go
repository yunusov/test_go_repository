package main

import (
	"flag"
	"github.com/gorilla/mux"
	s "lab_30/pkg/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	onlyOnePortPermittedErrMsg = "Only one usage of each socket address (protocol/network address/port) is normally permitted."
)

// go run main.go -port "<port_number>"
// go run main.go

func main() {
	port := getParamsPort()
	service := s.NewService(0)
	router := mux.NewRouter()

	makeHandleFuncs(router, service)
	startWebService(port)
}

func startWebService(port string) {
	for {
		log.Printf("Запуск веб-сервера на http://127.0.0.1:%v\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if strings.Contains(err.Error(), onlyOnePortPermittedErrMsg) {
			log.Printf("Port %v is busy. Trying to use next port.", port)
			port = getNextPort(port)
			continue
		}
		log.Fatal(err)
	}
}

func getNextPort(port string) string {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}
	portInt++
	return strconv.Itoa(portInt)
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
	if port == "" {
		port = "8080"
		//panic("Param 'port' is absent! Service cannot start!")
	}
	return port
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n",
		"hello", r.Method, r.Body, r.Header.Get("Content-Type"))
	w.Write([]byte("Hello"))
}
