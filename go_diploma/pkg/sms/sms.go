package sms

import (
	"fmt"
	"go_diploma/pkg/utils"
	"os"
	"strconv"
	"strings"
)

type SmsData struct {
	Сountry      string
	Bandwidth    string
	ResponseTime string
	Provider     string
}

func buildStruct(fields []string) *SmsData {
	return &SmsData{fields[0], fields[1], fields[2], fields[3]}
}

func (sms *SmsData) ToString() string {
	return fmt.Sprintf("Country = %s, Bandwidth = %s, ResponseTime = %s, Provider = %s",
		sms.Сountry, sms.Bandwidth, sms.ResponseTime, sms.Provider)
}

func LoadData(conf utils.Config) []*SmsData {
	smsDataBytes, err := os.ReadFile(conf.EmuPath + conf.SmsDataFile)
	if err != nil {
		panic(err)
	}
	if len(smsDataBytes) == 0 {
		fmt.Println("SMS data file is empty!")
	}
	smsDataStr := string(smsDataBytes)
	return parceSmsData(smsDataStr, conf)
}

func parceSmsData(smsDataStr string, conf utils.Config) (result []*SmsData) {
	smsDataRows := strings.Split(smsDataStr, "\n")
	for i, row := range smsDataRows {
		fmt.Println(i, row)
		fields := strings.Split(row, ";")
		isValid := validateRecords(fields, conf)
		if !isValid {
			continue
		}
		result = append(result, buildStruct(fields))
	}
	return
}

func validateRecords(fields []string, conf utils.Config) bool {
	if len(fields) < 4 {
		fmt.Printf("sms.validateRecords: len < 4, fields = %v\n", fields)
		return false
	}
	for i, field := range fields {
		field := strings.Trim(field, " ")
		if len(field) == 0 {
			fmt.Printf("sms.validateRecords: len == 0\n")
			return false
		}
		if i == 0 {
			if !utils.SliceContains(conf.CountryCodes, field) {
				fmt.Printf("sms.validateRecords: not in Country Codes, fields = %v\n", fields)
				return false
			}
		} else if i == 1 || i == 2 {
			_, err := strconv.Atoi(field)
			if err != nil {
				fmt.Printf("sms.validateRecords: not int, fields = %v\n", fields)
				return false
			}
		} else if i == 3 {
			if !utils.SliceContains(conf.SmsProviders, field) {
				fmt.Printf("sms.validateRecords: not in Providers, fields = %v\n", fields)
				return false
			}
		}
	}
	return true
}
