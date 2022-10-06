package utils

import (
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func LogRequest(name string, r *http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n",
		name, r.Method, r.Body, r.Header.Get("Content-Type"))
}
