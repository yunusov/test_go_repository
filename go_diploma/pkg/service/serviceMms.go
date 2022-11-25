package service

import "go_diploma/pkg/mms"

func loadMmsData(resultSetT *ResultSetT) error {
	mmsData, err := mms.LoadData(srvConfig)
	if err != nil {
		return err
	}
	mmsData1 := make([]mms.MmsData, 0)
	mmsData2 := make([]mms.MmsData, 0)
	countryCodesMap := srvConfig.GetCoutryCodesMap()
	for _, data := range mmsData {
		data1 := *data
		data2 := *data
		data2.Country = countryCodesMap[data2.Country]
		mmsData2 = append(mmsData2, data2)
		mmsData1 = append(mmsData1, data1)
	}
	sortingMmsDataByCountry(&mmsData1)
	sortingMmsDataByProvider(&mmsData2)
	resultSetT.MMS = [][]mms.MmsData{mmsData1, mmsData2}
	return nil
}

func sortingMmsDataByCountry(arr *[]mms.MmsData) {
	arrMmsData := *arr
	bsSize := len(arrMmsData)
	for i := 1; i < bsSize; i++ {
		itemX := arrMmsData[i]
		j := i
		for ; j > 0 && arrMmsData[j-1].Country > itemX.Country; j-- {
			arrMmsData[j] = arrMmsData[j-1]
		}
		arrMmsData[j] = itemX
	}
}

func sortingMmsDataByProvider(arr *[]mms.MmsData) {
	arrMmsData := *arr
	bsSize := len(arrMmsData)
	for i := 1; i < bsSize; i++ {
		itemX := arrMmsData[i]
		j := i
		for ; j > 0 && arrMmsData[j-1].Provider > itemX.Provider; j-- {
			arrMmsData[j] = arrMmsData[j-1]
		}
		arrMmsData[j] = itemX
	}
}
