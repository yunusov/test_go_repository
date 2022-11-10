package main

import (
	"context"
	"flag"
	s "lab_30/pkg/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
)

// go run main.go -port "<port_number>"
// go run main.go

func main() {
	port := getParamsPort()
	service := s.NewService(0)
	router := mux.NewRouter()
	makeHandleFuncs(router, service)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Запуск веб-сервера на %v\n", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	go initGracefulShutdown(stop, &wg, srv, ctx)
	wg.Wait()
	log.Println("Server Exited Properly")
}

func initGracefulShutdown(cancelFunc context.CancelFunc, wg *sync.WaitGroup, srv *http.Server, ctx context.Context) {
	defer wg.Done()
	log.Println("Graceful Shutdown")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	cancelFunc()
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
		//log.Fatal("Param 'port' is absent! Service cannot start!")
	}
	return ":" + port
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n",
		"hello", r.Method, r.Body, r.Header.Get("Content-Type"))
	w.Write([]byte("Hello"))
}
