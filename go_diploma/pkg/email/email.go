package email

import (
	"fmt"
	"go_diploma/pkg/config"
	"go_diploma/pkg/utils"
	"os"
	"strconv"
	"strings"
)

type EmailData struct {
	Country      string
	Provider     string
	DeliveryTime int
}

func (email *EmailData) ToString() string {
	return fmt.Sprintf("Country = %s, Provider = %s, DeliveryTime = %d",
		email.Country, email.Provider, email.DeliveryTime)
}

func LoadData(conf *config.Config) []*EmailData {
	emailDataBytes, err := os.ReadFile(conf.GetEmuPath() + conf.GetEmailDataFile())
	if err != nil {
		panic(err)
	}
	if len(emailDataBytes) == 0 {
		fmt.Println("Email data file is empty!")
	}
	return parceEmailData(string(emailDataBytes), conf)
}

func parceEmailData(dataStr string, conf *config.Config) (result []*EmailData) {
	vcDataRows := strings.Split(dataStr, "\n")
	for i, row := range vcDataRows {
		fmt.Println(i, row)
		fields := strings.Split(row, ";")
		isValid, record := validateRecords(fields, conf)
		if !isValid {
			continue
		}
		result = append(result, record)
	}
	return
}

func validateRecords(fields []string, conf *config.Config) (bool, *EmailData) {
	if len(fields) < 3 {
		fmt.Printf("voicecall.validateRecords: len < 3, fields = %v\n", fields)
		return false, nil
	}
	for i, field := range fields {
		field := strings.Trim(field, " ")
		if len(field) == 0 {
			fmt.Printf("voicecall.validateRecords: len == 0\n")
			return false, nil
		}
		if i == 0 {
			if !utils.SliceContains(conf.GetCoutryCodes(), field) {
				fmt.Printf("voicecall.validateRecords: not in Country Codes, fields = %v\n", fields)
				return false, nil
			}
		} else if i == 1 {
			if !utils.SliceContains(conf.GetEmailProviders(), field) {
				fmt.Printf("voicecall.validateRecords: not in Providers, fields = %v\n", fields)
				return false, nil
			}
		}
	}
	dtValue, err := strconv.Atoi(fields[2])
	if err != nil {
		fmt.Printf("voicecall.validateRecords: not int, fields = %v\n", fields)
		return false, nil
	}
	return true, &EmailData{fields[0], fields[1], dtValue}
}
