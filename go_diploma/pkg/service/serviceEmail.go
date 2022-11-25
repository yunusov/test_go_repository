package service

import (
	"go_diploma/pkg/email"
)

func loadEmailData(resultSetT *ResultSetT) error {
	emailData, err := email.LoadData(srvConfig)
	if err != nil {
		return err
	}
	emailData1 := make(map[string][][]email.EmailData, 0)
	for _, data := range emailData {
		data1 := *data
		if emailData1[data1.Country] == nil {
			emailData1[data1.Country] = getSortingEmailDataByDTCountry(emailData, data1.Country)
		}
	}
	resultSetT.Email = emailData1
	return nil
}

func getSortingEmailDataByDTCountry(arr []*email.EmailData, country string) [][]email.EmailData {
	filteredArr := filterByCountry(arr, country)
	bsSize := len(filteredArr)
	for i := 1; i < bsSize; i++ {
		itemX := filteredArr[i]
		j := i
		for ; j > 0 && filteredArr[j-1].DeliveryTime > itemX.DeliveryTime; j-- {
			filteredArr[j] = filteredArr[j-1]
		}
		filteredArr[j] = itemX
	}
	arrSize := 0
	if bsSize < 3 {
		arrSize = bsSize
	} else {
		arrSize = 3
	}
	maxEmail := []email.EmailData{}
	minEmail := []email.EmailData{}
	for i := 0; i < arrSize; i++ {
		maxEmail = append(maxEmail, filteredArr[i])
		minEmail = append(minEmail, filteredArr[bsSize-i-1])
	}
	return [][]email.EmailData{maxEmail, minEmail}
}

func filterByCountry(arr []*email.EmailData, country string) (result []email.EmailData) {
	for _, v := range arr {
		if v.Country == country {
			result = append(result, *v)
		}
	}
	return
}
