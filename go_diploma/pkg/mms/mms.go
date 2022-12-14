package mms

import (
	"encoding/json"
	"fmt"
	"go_diploma/pkg/config"
	"go_diploma/pkg/request"
	"go_diploma/pkg/utils"
	"net/http"
	"strconv"
	"strings"
)

type MmsData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func BuildStruct(fields []string) *MmsData {
	return &MmsData{fields[0], fields[1], fields[2], fields[3]}
}

func (mms *MmsData) ToString() string {
	return fmt.Sprintf("Country = %s, Bandwidth = %s, ResponseTime = %s, Provider = %s",
		mms.Country, mms.Bandwidth, mms.ResponseTime, mms.Provider)
}

func LoadData(conf *config.Config) ([]*MmsData, error) {
	response, err := sendMmsRequest(conf)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(response))
	mmsData, err := parceResponse([]byte(response))
	if err != nil {
		return nil, err
	}
	mmsData = validateRecords(mmsData, conf)
	return mmsData, nil
}

func sendMmsRequest(conf *config.Config) ([]byte, error) {
	rs := &request.RequestStruct{UrlRequest: conf.GetMmsServerAddress(),
		Content: "", HttpMethod: http.MethodGet}
	result, err := rs.Send()
	if err != nil {
		fmt.Println("sendSupportRequest error:", err.Error())
		return nil, err
	}
	return result, nil
}

func parceResponse(response []byte) ([]*MmsData, error) {
	if len(response) == 0 {
		return nil, fmt.Errorf("Response zero length")
	}
	mmsData := []*MmsData{}
	if err := json.Unmarshal(response, &mmsData); err != nil {
		fmt.Printf("service.loadStore: Error: %s", err.Error())
		return nil, err
	}
	return mmsData, nil
}

func validateRecords(mmsData []*MmsData, conf *config.Config) (result []*MmsData) {
	for _, record := range mmsData {
		bandwidth := strings.Trim(record.Bandwidth, " ")
		provider := strings.Trim(record.Provider, " ")
		responseTime := strings.Trim(record.ResponseTime, " ")
		country := strings.Trim(record.Country, " ")
		if len(bandwidth) == 0 || len(provider) == 0 || len(responseTime) == 0 || len(country) == 0 {
			fmt.Printf("sms.validateRecords: len == 0, fields = %v\n", record.ToString())
			continue
		}
		if !utils.SliceContains(conf.GetCoutryCodes(), country) {
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
		if !utils.SliceContains(conf.GetSmsProviders(), provider) {
			fmt.Printf("mms.validateRecords: not in Providers, fields = %v\n", record.ToString())
			continue
		}
		result = append(result, record)
	}
	return
}
