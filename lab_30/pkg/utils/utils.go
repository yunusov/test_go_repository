package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func LogRequest(name string, r *http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n",
		name, r.Method, r.Body, r.Header.Get("Content-Type"))
}

func IsCtJson(header string) bool {
	return strings.ContainsAny(header, "application/json")
}

func GetContent(r *http.Request, w http.ResponseWriter) ([]byte, bool) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("Error: %s", err.Error())
		return nil, true
	}
	return content, false
}

func UnMarshalData(content []byte, dat any, w http.ResponseWriter) bool {
	if err := json.Unmarshal(content, dat); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("Error: %s", err.Error())
		return true
	}
	return false
}

func MarshalData(user any, w http.ResponseWriter) (string, bool) {
	encUser, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("Error: %s", err.Error())
		return "", true
	}
	return string(encUser) + "\n", false
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func GetRequestParam(r *http.Request, paramName string) string {
	params := mux.Vars(r)
	userId := params[paramName]
	return userId
}