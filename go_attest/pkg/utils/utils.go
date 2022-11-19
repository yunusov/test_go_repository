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

func GetContent(r *http.Request, w http.ResponseWriter) ([]byte, error) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("utils.GetContent Error: %s", err.Error())
		return nil, err
	}
	return content, nil
}

func UnMarshalData(content []byte, dat any, w http.ResponseWriter) error {
	if err := json.Unmarshal(content, dat); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("utils.UnMarshalData Error: %s", err.Error())
		return err
	}
	return nil
}

func MarshalData(user any) (string, error) {
	encUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("utils.MarshalData Error: %s", err.Error())
		return "", err
	}
	return string(encUser) + "\n", nil
}

func GetRequestParam(r *http.Request, paramName string) string {
	params := mux.Vars(r)
	userId := params[paramName]
	return userId
}

func DecodeData(r *http.Request, w http.ResponseWriter) (map[string]interface{}, error) {
	content, err := GetContent(r, w)
	if err != nil {
		return nil, err
	}
	log.Printf("content = %s", string(content))
	dat := make(map[string]interface{})
	if err := UnMarshalData(content, &dat, w); err != nil {
		return nil, err
	}
	return dat, nil
}
