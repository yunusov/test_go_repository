package main

import (
	"context"
	"flag"
	s "go_attest/pkg/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
)

// go run main.go -port 8080

func main() {
	port := getParamsPort()
	service := s.NewService()
	service.Init()
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
		log.Printf("Server started on %v\n", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()
	go initGracefulShutdown(stop, &wg, srv, ctx, service)
	wg.Wait()
	log.Println("Server exited properly")
}

func initGracefulShutdown(cancelFunc context.CancelFunc, wg *sync.WaitGroup, srv *http.Server, ctx context.Context, service *s.Service) {
	defer wg.Done()
	service.SaveStore()
	log.Println("Graceful Shutdown")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	cancelFunc()
}

func makeHandleFuncs(router *mux.Router, service *s.Service) {
	router.HandleFunc("/city/{cityid:[\\d]+}", service.GetCityById).Methods(http.MethodGet)
	router.HandleFunc("/create", service.Create).Methods(http.MethodPost)
	router.HandleFunc("/", hello)
	router.HandleFunc("/region", service.GetByRegion).Methods(http.MethodPost)
	router.HandleFunc("/district", service.GetByDistrict).Methods(http.MethodPost)
	router.HandleFunc("/population", service.GetByPopulationRange).Methods(http.MethodPost)
	router.HandleFunc("/foundation", service.GetByFoundationRange).Methods(http.MethodPost)
	router.HandleFunc("/city", service.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/city/{cityid:[\\d]+}", service.UpdatePopulationById).Methods(http.MethodPut)
	http.Handle("/", router)
}

func getParamsPort() (port string) {
	flag.StringVar(&port, "port", "", "set port")
	flag.Parse()
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n",
		"hello", r.Method, r.Body, r.Header.Get("Content-Type"))
	w.Write([]byte("Hello"))
}
