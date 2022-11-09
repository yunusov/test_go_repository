package service

import (
	"bytes"
	"fmt"
	"io"
	u "lab_30/pkg/user"
	ut "lab_30/pkg/utils"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

const (
	contentType     = "Content-Type"
	contentTypeJson = "application/json"

	urlCreate      = "/create"
	urlDelete      = "/delete"
	urlMakeFriends = "/make_friends"
	urlFriends     = "/friends/?"
	urlGetAll      = "/get"
	urlUpdateAge   = "/?"
)

type RequestStruct struct {
	t          *testing.T
	urlRequest string
	content    string
	httpMethod string
}

var (
	contentCreate = "" +
		"{" +
		"	\"name\": \"User9\"," +
		"	\"age\": \"241\"," +
		"	\"friends\": []" +
		"}"
	contentDelete = "" +
		"{" +
		"	\"target_id\": \"?\"" +
		"}"
	contentMakeFriends = "" +
		"{" +
		"	\"source_id\": \"?\"," +
		"	\"target_id\": \"?\"" +
		"}"
	contentUpdateAge = "" +
		"{" +
		"	\"new age\": \"?\"" +
		"}"
)

func TestCreate(t *testing.T) {
	s := NewService(0)
	srv := httptest.NewServer(http.HandlerFunc(s.Create))
	rs := &RequestStruct{t, srv.URL + urlCreate, contentCreate, http.MethodPost}
	response := sendRequest(rs)
	user, err := s.getUser(response)
	errProcessor(err, t)
	_, err = strconv.Atoi(response)
	errProcessor(err, t)
	t.Logf("User %v was created!", user.ToString())
}

func TestDelete(t *testing.T) {
	s := NewService(0)
	srv := httptest.NewServer(http.HandlerFunc(s.Delete))
	user := s.createUser()
	userName := user.Name
	userId := user.GetId()
	t.Logf("User %v was created!", userName)
	contentDeleteBytes := strings.Replace(contentDelete, "?", userId, 1)
	rs := &RequestStruct{t, srv.URL + urlDelete, contentDeleteBytes, http.MethodPost}
	response := sendRequest(rs)

	text := string(response)
	if text != userName {
		errProcessor(fmt.Errorf("Actual response (%v) differs from expected (%v)", text, userName), t)
	}
	delUser, err := s.getUser(userId)
	if err != nil && strings.Contains(err.Error(), "user is nil with ID=") && delUser == nil {
		t.Logf("User %v was deleted!", userName)
	} else {
		errProcessor(err, t)
		errProcessor(fmt.Errorf("Deleted user (%v) still alive!", userId), t)
	}
}

func TestMakeFriends(t *testing.T) {
	s := NewService(0)
	srv := httptest.NewServer(http.HandlerFunc(s.MakeFriends))
	user1 := s.createUser()
	user2 := s.createUser()
	userName1 := user1.Name
	userName2 := user2.Name
	userId1 := user1.GetId()
	userId2 := user2.GetId()
	t.Logf("User %v and %v were created!", userName1, userName2)
	contentMakeFriends = strings.Replace(contentMakeFriends, "?", userId1, 1)
	contentMakeFriends = strings.Replace(contentMakeFriends, "?", userId2, 1)

	rs := &RequestStruct{t, srv.URL + urlMakeFriends, contentMakeFriends, http.MethodPost}
	response := sendRequest(rs)
	trueResponse := "? и ? теперь друзья"
	trueResponse = strings.Replace(trueResponse, "?", userName1, 1)
	trueResponse = strings.Replace(trueResponse, "?", userName2, 1)

	if response != trueResponse {
		errProcessor(fmt.Errorf("Actual response ('%s') does not matched with expected ('%s')!\n", response, trueResponse), t)
	}
	user1, err := s.getUser(userId1)
	errProcessor(err, t)
	user2, err = s.getUser(userId2)
	errProcessor(err, t)

	friendIdsUser1 := user1.GetFriendIds()
	friendIdsUser2 := user2.GetFriendIds()
	if !(len(friendIdsUser1) == 1 && ut.SliceContains(friendIdsUser1, userId2) &&
		len(friendIdsUser2) == 1 && ut.SliceContains(friendIdsUser2, userId1)) {
		t.Logf("friendIdsUser1 = %v, friendIdsUser2 = %v", friendIdsUser1, friendIdsUser2)
		errProcessor(fmt.Errorf("MakeFriends function failed. Users aren't friends!%v", " "), t)
	}
}

func TestUpdateAge(t *testing.T) {
	s := NewService(0)
	router := mux.NewRouter()
	router.HandleFunc("/{userid:[\\d]+}", s.UpdateAgeById).Methods(http.MethodPut)
	srv := httptest.NewServer(router)
	rand.Seed(time.Now().UnixNano())
	newAgeStr := strconv.Itoa(rand.Intn(101))
	user := s.createUser()
	userId := user.GetId()
	urlUpdAge := strings.Replace(urlUpdateAge, "?", userId, 1)

	contentUpdateAge = strings.Replace(contentUpdateAge, "?", newAgeStr, 1)
	rs := &RequestStruct{t, srv.URL + urlUpdAge, contentUpdateAge, http.MethodPut}
	response := sendRequest(rs)
	expectedResponse := "возраст пользователя успешно обновлён"
	if response != expectedResponse {
		errProcessor(fmt.Errorf("Expected response differs from actual (%v)!", response), t)
	}
	user, err := s.getUser(userId)
	errProcessor(err, t)
	actualAge := user.GetAge()
	if newAgeStr != actualAge {
		errProcessor(fmt.Errorf("Actual age (%v) differs from expected (%v)!", actualAge, newAgeStr), t)
	}
}

func TestGetFriendsById(t *testing.T) {
	s := NewService(0)
	router := mux.NewRouter()
	router.HandleFunc("/friends/{userid:[\\d]+}", s.GetFriendsById).Methods(http.MethodGet)
	srv := httptest.NewServer(router)
	user1 := s.createUser()
	user2 := s.createUser()
	user3 := s.createUser()
	friends := make(map[string]*u.User)
	friends[user2.GetId()] = user2
	friends[user3.GetId()] = user3

	s.MakeFriend(user1.GetId(), user2.GetId(), user3.GetId())
	urlFrnds := strings.Replace(urlFriends, "?", user1.GetId(), 1)

	rs := &RequestStruct{t, srv.URL + urlFrnds, "", http.MethodGet}
	response := sendRequest(rs)
	if len(response) == 0 {
		errProcessor(fmt.Errorf("Response empty!"), t)
	}
	t.Logf("response = %v", response)
}

func errProcessor(err error, t *testing.T) {
	if err != nil {
		t.Log("errProcessor fail: ", err)
		t.Fail()
	}
}

func sendRequest(rs *RequestStruct) string {
	var resp *http.Response
	var err error
	if rs.httpMethod == http.MethodPost {
		resp, err = http.Post(rs.urlRequest, contentTypeJson, bytes.NewReader([]byte(rs.content)))
		errProcessor(err, rs.t)
	} else if rs.httpMethod == http.MethodGet {
		resp, err = http.Get(rs.urlRequest)
		errProcessor(err, rs.t)
	} else {
		client := &http.Client{}
		req, err := http.NewRequest(rs.httpMethod, rs.urlRequest, bytes.NewReader([]byte(rs.content)))
		errProcessor(err, rs.t)
		req.Header.Add(contentType, contentTypeJson)
		resp, err = client.Do(req)
		errProcessor(err, rs.t)
	}
	textBytes, err := io.ReadAll(resp.Body)
	errProcessor(err, rs.t)
	defer resp.Body.Close()

	text := string(textBytes)
	return text
}
