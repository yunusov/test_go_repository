package support

import (
	"encoding/json"
	"fmt"
	"go_diploma/pkg/config"
	"go_diploma/pkg/request"
	"net/http"
)

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

func (sup *SupportData) ToString() string {
	return fmt.Sprintf("Topic = %s, ActiveTickets = %d",
		sup.Topic, sup.ActiveTickets)
}

func LoadData(conf *config.Config) ([]*SupportData, error) {
	response, err := sendSupportRequest(conf)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(response))
	return parceResponse([]byte(response))
}

func sendSupportRequest(conf *config.Config) ([]byte, error) {
	rs := &request.RequestStruct{UrlRequest: conf.GetSupportServerAddress(),
		Content: "", HttpMethod: http.MethodGet}
	result, err := rs.Send()
	if err != nil {
		fmt.Println("sendSupportRequest error:", err.Error())
		return nil, err
	}
	return result, nil
}

func parceResponse(response []byte) ([]*SupportData, error) {
	if len(response) == 0 {
		return nil, fmt.Errorf("Zero length response")
	}
	mmsData := []*SupportData{}
	if err := json.Unmarshal(response, &mmsData); err != nil {
		fmt.Printf("service.loadStore: Error: %s", err.Error())
		return nil, err
	}
	return mmsData, nil
}
