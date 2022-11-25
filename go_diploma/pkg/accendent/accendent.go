package accendent

import (
	"encoding/json"
	"fmt"
	"go_diploma/pkg/config"
	"go_diploma/pkg/request"
	"net/http"
)

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` // возможные статусы active и closed
}

func (inc *IncidentData) ToString() string {
	return fmt.Sprintf("Topic = %s, Status = %s", inc.Topic, inc.Status)
}

func LoadData(conf *config.Config) ([]*IncidentData, error) {
	response, err := sendAccedentRequest(conf)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(response))
	return parceResponse([]byte(response))
}

func sendAccedentRequest(conf *config.Config) ([]byte, error) {
	rs := &request.RequestStruct{UrlRequest: conf.GetIncindentServerAddress(),
		Content: "", HttpMethod: http.MethodGet}
	result, err := rs.Send()
	if err != nil {
		fmt.Println("sendAccedentRequest error:", err.Error())
		return nil, err
	}
	return result, nil
}

func parceResponse(response []byte) ([]*IncidentData, error) {
	if len(response) == 0 {
		return nil, fmt.Errorf("Response zero length")
	}
	mmsData := []*IncidentData{}
	if err := json.Unmarshal(response, &mmsData); err != nil {
		fmt.Printf("service.loadStore: Error: %s", err.Error())
		return nil, err
	}
	return mmsData, nil
}
