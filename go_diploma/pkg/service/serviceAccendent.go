package service

import (
	"go_diploma/pkg/accendent"
)

func loadIncindentData(resultSetT *ResultSetT) error {
	incendData, err := accendent.LoadData(srvConfig)
	if err != nil {
		return err
	}
	incendData1 := make([]accendent.IncidentData, 0)
	for _, data := range incendData {
		incendData1 = append(incendData1, *data)
	}
	sortingIncidentDataByStatus(&incendData1)
	resultSetT.Incidents = incendData1
	return nil
}

func sortingIncidentDataByStatus(arr *[]accendent.IncidentData) {
	arrSmsData := *arr
	bsSize := len(arrSmsData)
	for i := 1; i < bsSize; i++ {
		itemX := arrSmsData[i]
		j := i
		for ; j > 0 && arrSmsData[j-1].Status > itemX.Status; j-- {
			arrSmsData[j] = arrSmsData[j-1]
		}
		arrSmsData[j] = itemX
	}
}
