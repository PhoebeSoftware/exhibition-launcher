package realdebrid

import (
	"fmt"
	"net/http"
)

type TraficInfo struct {
	Left  int
	Bytes int
	Links int
	Limit int
	Type  string
	Extra int
	Reset string
}
func (client *RealDebridClient) GetTrafic() (map[string]TraficInfo, error) {

	type TraficResponse map[string]TraficInfo

	var traficResponse TraficResponse

	req, err := client.newRequest(http.MethodGet, "/traffic", nil, "", nil)
	if err != nil {
		return nil, fmt.Errorf("get request failed while requesting trafic: %w", err)
	}

	err = client.do(req, &traficResponse)
	if err != nil {
		return nil, err
	}

	return traficResponse, nil
}

