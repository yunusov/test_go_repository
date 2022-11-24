package voicecall

import (
	"fmt"
	"go_diploma/pkg/config"
	"go_diploma/pkg/utils"
	"os"
	"strconv"
	"strings"
)

type VoiceCallData struct {
	Country       string
	Load          int
	ResponseTime  int
	Provider      string
	Stability     float32
	TtfbClearence int
	CallDuration  int
	UnknowValue   int
}

func (vc *VoiceCallData) ToString() string {
	return fmt.Sprintf("Country = %s, Load = %d, ResponseTime = %d, Provider = %s, "+
		"Stability = %f, TtfbClearence = %d, CallDuration = %d, UnknowValue = %d",
		vc.Country, vc.Load, vc.ResponseTime, vc.Provider, vc.Stability,
		vc.TtfbClearence, vc.CallDuration, vc.UnknowValue)
}

func LoadData(conf *config.Config) []*VoiceCallData {
	vcBytes, err := os.ReadFile(conf.GetEmuPath() + conf.GetVoiceDataFile())
	if err != nil {
		panic(err)
	}
	if len(vcBytes) == 0 {
		fmt.Println("VoiceCall data file is empty!")
	}
	return parceVcData(string(vcBytes), conf)
}

func parceVcData(dataStr string, conf *config.Config) (result []*VoiceCallData) {
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

func validateRecords(fields []string, conf *config.Config) (bool, *VoiceCallData) {
	if len(fields) < 8 {
		fmt.Printf("voicecall.validateRecords: len < 8, fields = %v\n", fields)
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
		} else if i == 3 {
			if !utils.SliceContains(conf.GetVoiceProviders(), field) {
				fmt.Printf("voicecall.validateRecords: not in Providers, fields = %v\n", fields)
				return false, nil
			}
		}
	}
	loadValue, err := strconv.Atoi(fields[1])
	if err != nil {
		fmt.Printf("voicecall.validateRecords: not int, fields = %v\n", fields)
		return false, nil
	}
	responseTimeValue, err := strconv.Atoi(fields[2])
	if err != nil {
		fmt.Printf("voicecall.validateRecords: not int, fields = %v\n", fields)
		return false, nil
	}
	stabilityValue, err := strconv.ParseFloat(fields[4], 32)
	if err != nil {
		fmt.Printf("voicecall.validateRecords: not float32, fields = %v\n", fields)
		return false, nil
	}
	clearenceValue, err := strconv.Atoi(fields[5])
	if err != nil {
		fmt.Printf("voicecall.validateRecords: not int, fields = %v\n", fields)
		return false, nil
	}
	callDurationValue, err := strconv.Atoi(fields[6])
	if err != nil {
		fmt.Printf("voicecall.validateRecords: not int, fields = %v\n", fields)
		return false, nil
	}
	unknownValue, err := strconv.Atoi(fields[7])
	if err != nil {
		fmt.Printf("voicecall.validateRecords: not int, fields = %v\n", fields)
		return false, nil
	}
	return true, &VoiceCallData{fields[0], loadValue, responseTimeValue, fields[3],
		float32(stabilityValue), clearenceValue, callDurationValue, unknownValue}
}
