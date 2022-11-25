package service

import (
	"encoding/json"
	"fmt"
	"go_diploma/pkg/accendent"
	"go_diploma/pkg/billing"
	"go_diploma/pkg/config"
	"go_diploma/pkg/email"
	"go_diploma/pkg/mms"
	"go_diploma/pkg/sms"
	"go_diploma/pkg/support"
	"go_diploma/pkg/voicecall"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ResultT struct {
	Status bool       `json:"status"` // True, если все этапы сбора данных прошли успешно, False во всех остальных случаях
	Data   ResultSetT `json:"data"`   // Заполнен, если все этапы сбора  данных прошли успешно, nil во всех остальных случаях
	Error  string     `json:"error"`  // Пустая строка, если все этапы сбора данных прошли успешно, в случае ошибки заполнено текстом ошибки
}

type ResultSetT struct {
	SMS       [][]sms.SmsData                `json:"sms"`
	MMS       [][]mms.MmsData                `json:"mms"`
	VoiceCall []voicecall.VoiceCallData      `json:"voice_call"`
	Email     map[string][][]email.EmailData `json:"email"`
	Billing   billing.BillingData            `json:"billing"`
	Support   []int                          `json:"support"`
	Incidents []accendent.IncidentData       `json:"incident"`
}

var srvConfig *config.Config

func Start(conf *config.Config) {
	srvConfig = conf
	router := mux.NewRouter()
	router.HandleFunc("/", handleConnection)
	router.HandleFunc("/test", gatheringData)
	http.Handle("/", router)
	log.Printf("Запуск веб-сервера на http://%s", conf.GetServiceServerAddress())
	err := http.ListenAndServe(conf.GetServiceServerAddress(), nil)
	log.Fatal(err)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s method = %v\n", "hello", r.Method)
	w.Write([]byte("OK"))
}

func gatheringData(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s method = %v\n", "gatheringData", r.Method)
	status := true
	errorStr := ""
	resultSetT, err := getResultData()
	if err != nil {
		status = false
		errorStr = err.Error()
	}
	resultT := ResultT{status, *resultSetT, errorStr}
	response, err := marshalData(resultT)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func getResultData() (*ResultSetT, error) {
	result := &ResultSetT{}
	emptyResult := &ResultSetT{}
	err := loadSmsData(result)
	if err != nil {
		return emptyResult, err
	}
	err = loadMmsData(result)
	if err != nil {
		return emptyResult, err
	}
	err = loadVoiceCallData(result)
	if err != nil {
		return emptyResult, err
	}
	err = loadEmailData(result)
	if err != nil {
		return emptyResult, err
	}
	err = loadBillingData(result)
	if err != nil {
		return emptyResult, err
	}
	err = loadIncindentData(result)
	if err != nil {
		return emptyResult, err
	}
	err = loadSupportData(result)
	if err != nil {
		return emptyResult, err
	}
	return result, nil
}

func marshalData(data any) (string, error) {
	encData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("service.MarshalData Error: %s", err.Error())
		return "", err
	}
	return string(encData) + "\n", nil
}

func loadSupportData(result *ResultSetT) error {
	supportData, err := support.LoadData(srvConfig)
	if err != nil {
		return err
	}
	cntTickets := 0
	for _, data := range supportData {
		cntTickets += data.ActiveTickets
	}
	wlState := 3
	avgTimePerTicket := 60 / 18
	workload := avgTimePerTicket * cntTickets / 7
	if workload < 9 {
		wlState = 1
	} else if workload >= 9 && workload <= 16 {
		wlState = 2
	}
	result.Support = append(result.Support, wlState)
	result.Support = append(result.Support, workload)
	return nil
}

func loadVoiceCallData(result *ResultSetT) error {
	voiceData, err := voicecall.LoadData(srvConfig)
	if err != nil {
		return err
	}
	for _, data := range voiceData {
		data1 := *data
		result.VoiceCall = append(result.VoiceCall, data1)
	}
	return nil
}

func loadBillingData(result *ResultSetT) error {
	billData, err := billing.LoadData(srvConfig)
	if err != nil {
		return err
	}
	result.Billing = *billData
	return nil
}
