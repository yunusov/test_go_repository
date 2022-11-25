package service

import (
	"go_diploma/pkg/sms"
)

func loadSmsData(resultSetT *ResultSetT) error {
	smsData, err := sms.LoadData(srvConfig)
	if err != nil {
		return err
	}
	smsData1 := make([]sms.SmsData, 0)
	smsData2 := make([]sms.SmsData, 0)
	countryCodesMap := srvConfig.GetCoutryCodesMap()
	for _, data := range smsData {
		data1 := *data
		data2 := *data
		data2.Country = countryCodesMap[data2.Country]
		smsData2 = append(smsData2, data2)
		smsData1 = append(smsData1, data1)
	}
	sortingSmsDataByCountry(&smsData1)
	sortingSmsDataByProvider(&smsData2)
	resultSetT.SMS = [][]sms.SmsData{smsData1, smsData2}
	return nil
}

func sortingSmsDataByCountry(arr *[]sms.SmsData) {
	arrSmsData := *arr
	bsSize := len(arrSmsData)
	for i := 1; i < bsSize; i++ {
		itemX := arrSmsData[i]
		j := i
		for ; j > 0 && arrSmsData[j-1].Country > itemX.Country; j-- {
			arrSmsData[j] = arrSmsData[j-1]
		}
		arrSmsData[j] = itemX
	}
}

func sortingSmsDataByProvider(arr *[]sms.SmsData) {
	arrSmsData := *arr
	bsSize := len(arrSmsData)
	for i := 1; i < bsSize; i++ {
		itemX := arrSmsData[i]
		j := i
		for ; j > 0 && arrSmsData[j-1].Provider > itemX.Provider; j-- {
			arrSmsData[j] = arrSmsData[j-1]
		}
		arrSmsData[j] = itemX
	}
}
