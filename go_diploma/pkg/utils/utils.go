package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	contentType     = "Content-Type"
	contentTypeJson = "application/json"
)

type Config struct {
	CountryCodes []string
	EmuPath      string
	MmsPage      string
	SmsProviders []string
	SmsDataFile  string
	SrvAddress   string
}

type RequestStruct struct {
	UrlRequest string
	Content    string
	HttpMethod string
}

func (conf *Config) ToString() string {
	return fmt.Sprintf("CountryCodes = %v, EmuPath = %v, Providers = %v", conf.CountryCodes, conf.EmuPath, conf.SmsProviders)
}

func LoadSettings() (conf Config) {
	if _, err := toml.DecodeFile("../settings.toml", &conf); err != nil {
		panic(err)
	}
	conf.SmsProviders = sorting(conf.SmsProviders)
	conf.CountryCodes = sorting(conf.CountryCodes)
	if len(strings.Trim(conf.EmuPath, " ")) == 0 {
		panic("Please fulfill field EmuPath in settings.toml!")
	}
	return
}

func SliceContains(s []string, searchterm string) bool {
	/*s - sorted slice*/
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func sorting(arr []string) []string {
	bsSize := len(arr)
	for i := 1; i < bsSize; i++ {
		x := arr[i]
		j := i
		for ; j > 0 && arr[j-1] > x; j-- {
			arr[j] = arr[j-1]
		}
		arr[j] = x
	}
	return arr
}

func SendRequest(rs *RequestStruct) ([]byte, error) {
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
