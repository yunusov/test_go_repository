package main

import (
	"fmt"
	"go_diploma/pkg/mms"
	"go_diploma/pkg/sms"
	"go_diploma/pkg/billing"
	"go_diploma/pkg/accendent"
	"go_diploma/pkg/support"
	"go_diploma/pkg/email"
	"go_diploma/pkg/service"
	"go_diploma/pkg/utils"
	"go_diploma/pkg/voicecall"
)

func main() {
	fmt.Println("Start main")
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
	vcData := voicecall.LoadData(config)
	fmt.Println("voicecall = ", len(vcData))
	for _, v := range vcData {
		fmt.Println(v.ToString())
	}
	emailData := email.LoadData(config)
	fmt.Println("email = ", len(emailData))
	for _, v := range emailData {
		fmt.Println(v.ToString())
	}
	billData := billing.LoadData(config)
	fmt.Println("billing = ", billData.ToString())
	supportData := support.LoadData(config)
	fmt.Println("supportData = ", len(supportData))
	for _, v := range supportData {
		fmt.Println(v.ToString())
	}
	accendData := accendent.LoadData(config)
	fmt.Println("accendData = ", len(accendData))
	for _, v := range accendData {
		fmt.Println(v.ToString())
	}
	service.Start(config)
}
