package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	contentType     = "Content-Type"
	contentTypeJson = "application/json"
)

type RequestStruct struct {
	UrlRequest string
	Content    string
	HttpMethod string
}

func (rs *RequestStruct) Send() ([]byte, error) {
	var resp *http.Response
	var err error
	if rs.HttpMethod == http.MethodPost {
		resp, err = http.Post(rs.UrlRequest, contentTypeJson, bytes.NewReader([]byte(rs.Content)))
		if err != nil {
			return nil, err
		}
	} else if rs.HttpMethod == http.MethodGet {
		resp, err = http.Get(rs.UrlRequest)
		if err != nil {
			return nil, err
		}
	} else {
		client := &http.Client{}
		req, err := http.NewRequest(rs.HttpMethod, rs.UrlRequest, bytes.NewReader([]byte(rs.Content)))
		if err != nil {
			return nil, err
		}
		req.Header.Add(contentType, contentTypeJson)
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}
	if resp.StatusCode == 500 {
		return nil, fmt.Errorf("code 500 has been returned%s", "")
	}
	textBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return textBytes, nil
}
