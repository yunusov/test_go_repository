package billing

import (
	"fmt"
	"go_diploma/pkg/config"
	"math"
	"os"
	"strconv"
)

type BillingData struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

func (bd *BillingData) ToString() string {
	return fmt.Sprintf("CreateCustomer = %t, Purchase = %t, Payout = %t, Recurring = %t, "+
		"FraudControl = %t, CheckoutPage = %t",
		bd.CreateCustomer, bd.Purchase, bd.Payout, bd.Recurring, bd.FraudControl, bd.CheckoutPage)
}

func LoadData(conf *config.Config) *BillingData {
	emailDataBytes, err := os.ReadFile(conf.GetEmuPath() + conf.GetBillingDataFile())
	if err != nil {
		panic(err)
	}
	if len(emailDataBytes) == 0 {
		fmt.Println("Billing data file is empty!")
	}
	return parceBillingData(emailDataBytes)
}

func parceBillingData(dataBytes []byte) *BillingData {
	num := convBytesToInt(dataBytes)
	return loadBillingData(num)
}

func convBytesToInt(dataBytes []byte) (result int) {
	x := 0.0
	for i := len(dataBytes) - 1; i >= 0; i-- {
		b, _ := strconv.Atoi(string(dataBytes[i]))
		result += int(math.Pow(2, x)) * b
		x++
	}
	return
}

func loadBillingData(data int) *BillingData {
	result := &BillingData{(data>>0)&1 == 1,
		(data>>1)&1 == 1,
		(data>>2)&1 == 1,
		(data>>3)&1 == 1,
		(data>>4)&1 == 1,
		(data>>5)&1 == 1}
	return result
}
