package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

const (
	cannotFindFileErrMsg = "The system cannot find the file specified."
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

func MarshalData(user any, w http.ResponseWriter) (string, error) {
	encUser, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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

func ReadFromFile(fileName string) (result []byte) {
	result, err := os.ReadFile(fileName)
	if len(result) == 0 {
		log.Printf("File %v is empty!", fileName)
		return nil
	} else if err != nil {
		if strings.Contains(err.Error(), cannotFindFileErrMsg) {
			log.Printf("utils.ReadFromFile: File %v is empty!", fileName)
			return nil
		}
		panic(err)
	}
	return result
}

func WriteToFile(str string, fileName string) {
	var b bytes.Buffer
	b.WriteString(str)
	if err := os.WriteFile(fileName, b.Bytes(), 0666); err != nil {
		log.Println("utils.WriteToFile: Error in file write: ", err)
	}
}

func SliceContains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
