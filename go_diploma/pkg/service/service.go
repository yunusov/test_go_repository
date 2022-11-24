package service

import (
	"go_diploma/pkg/accendent"
	"go_diploma/pkg/billing"
	"go_diploma/pkg/config"
	"go_diploma/pkg/email"
	"go_diploma/pkg/mms"
	"go_diploma/pkg/sms"
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

func Start(conf *config.Config) {
	router := mux.NewRouter()
	router.HandleFunc("/", handleConnection)

	http.Handle("/", router)
	log.Printf("Запуск веб-сервера на http://%s", conf.GetServiceServerAddress())
	err := http.ListenAndServe(conf.GetServiceServerAddress(), nil)
	log.Fatal(err)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s method = %v, body = %v, ct = %s\n",
		"hello", r.Method, r.Body, r.Header.Get("Content-Type"))
	w.Write([]byte("OK"))
	//w.WriteHeader(http.StatusOK)
}
