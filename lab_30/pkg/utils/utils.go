package utils

import (
	u "lab_30/pkg/user"
	//sr "lab_30/pkg/service"
	"io"
	"log"
	"encoding/json"
	"net/http"
	"strings"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func LogRequest(name string, r *http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n",
		name, r.Method, r.Body, r.Header.Get("Content-Type"))
}

func IsPostAndCtJson(method string, header string) bool {
	return method == http.MethodPost && strings.ContainsAny(header, "application/json")
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

func MarshalData(store map[string]*u.User, w http.ResponseWriter) (string, bool) {
	response := ""
	for _, user := range store {
		log.Printf("user = %s", user.ToString())
		encUser, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Printf("Error: %s", err.Error())
			return "", true
		}
		response += string(encUser) + "\n"
	}
	return response, false
}