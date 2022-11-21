package mms

import (
	"encoding/json"
	"fmt"
	"go_diploma/pkg/utils"
	"net/http"
	"strconv"
	"strings"
)

type MmsData struct {
	Сountry      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func BuildStruct(fields []string) *MmsData {
	return &MmsData{fields[0], fields[1], fields[2], fields[3]}
}

func (mms *MmsData) ToString() string {
	return fmt.Sprintf("Country = %s, Bandwidth = %s, ResponseTime = %s, Provider = %s",
		mms.Сountry, mms.Bandwidth, mms.ResponseTime, mms.Provider)
}

func LoadData(conf utils.Config) (result []*MmsData) {
	response := sendMmsRequest(conf)
	fmt.Println(string(response))
	mmsData := parceResponse([]byte(response))
	mmsData = validateRecords(mmsData, conf)
	return mmsData
}

func sendMmsRequest(conf utils.Config) (result []byte) {
	rs := &utils.RequestStruct{UrlRequest: conf.SrvAddress + conf.MmsPage, Content: "", HttpMethod: http.MethodGet}
	result, err := utils.SendRequest(rs)
	if err != nil {
		panic(err)
	}
	return
}

func parceResponse(response []byte) (result []*MmsData) {
	if len(response) == 0 {
		return
	}
	mmsData := []*MmsData{}
	if err := json.Unmarshal(response, &mmsData); err != nil {
		fmt.Printf("service.loadStore: Error: %s", err.Error())
		return
	}
	return mmsData
}

func validateRecords(mmsData []*MmsData, conf utils.Config) (result []*MmsData) {
	for _, record := range mmsData {
		bandwidth := strings.Trim(record.Bandwidth, " ")
		provider := strings.Trim(record.Provider, " ")
		responseTime := strings.Trim(record.ResponseTime, " ")
		country := strings.Trim(record.Сountry, " ")
		if len(bandwidth) == 0 || len(provider) == 0 || len(responseTime) == 0 || len(country) == 0 {
			fmt.Printf("sms.validateRecords: len == 0, fields = %v\n", record.ToString())
			continue
		}
		if !utils.SliceContains(conf.CountryCodes, country) {
			fmt.Printf("mms.validateRecords: not in Country Codes, fields = %v\n", record.ToString())
			continue
		}
		_, err := strconv.Atoi(bandwidth)
		if err != nil {
			fmt.Printf("sms.validateRecords: not int, fields = %v\n", record.ToString())
			continue
		}
		_, err = strconv.Atoi(responseTime)
		if err != nil {
			fmt.Printf("sms.validateRecords: not int, fields = %v\n", record.ToString())
			continue
		}
		if !utils.SliceContains(conf.SmsProviders, provider) {
			fmt.Printf("mms.validateRecords: not in Providers, fields = %v\n", record.ToString())
			continue
		}
		result = append(result, record)
	}
	return
}
