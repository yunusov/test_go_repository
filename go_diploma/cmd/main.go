package main

import (
	"fmt"
	"go_diploma/pkg/mms"
	"go_diploma/pkg/sms"
	"go_diploma/pkg/utils"
)

func main() {
	config := utils.LoadSettings()
	smsData := sms.LoadData(config)
	fmt.Println("smsData = ", len(smsData))
	for _, v := range smsData {
		fmt.Println(v.ToString())
	}
	mmsData := mms.LoadData(config)
	fmt.Println("mmsData = ", len(mmsData))
	for _, v := range mmsData {
		fmt.Println(v.ToString())
	}
}
