package realdebrid

import (
	"fmt"
	"net/http"
)

type TrafficInfo struct {
	Left  int
	Bytes int
	Links int
	Limit int
	Type  string
	Extra int
	Reset string
}
func (client *RealDebridClient) GetTraffic() (map[string]TrafficInfo, error) {

	type TrafficResponse map[string]TrafficInfo

	var trafficResponse TrafficResponse

	req, err := client.newRequest(http.MethodGet, "/traffic", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get request failed while requesting trafic: %w", err)
	}

	err = client.do(req, &trafficResponse)
	if err != nil {
		return nil, err
	}

	return trafficResponse, nil
}

