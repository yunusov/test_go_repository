package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

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

func GetRequestParam(r *http.Request, paramName string) string {
	params := mux.Vars(r)
	userId := params[paramName]
	return userId
}

func DecodeData(r *http.Request, w http.ResponseWriter) (map[string]interface{}, bool) {
	content, shouldReturn := GetContent(r, w)
	if shouldReturn {
		return nil, true
	}

	log.Printf("content = %s", string(content))
	var dat map[string]interface{}
	if shouldReturn1 := UnMarshalData(content, &dat, w); shouldReturn1 {
		return nil, true
	}
	return dat, false
}