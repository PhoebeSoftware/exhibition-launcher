package proxy_client

import (
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"net/http"
	"net/url"
	"strconv"
)

func (proxyClient *ProxyClient) GetMetadataByName(name string) ([]igdb.Metadata, error) {
	params := url.Values{}
	params.Add("name", name)

	req, err := proxyClient.newRequest(http.MethodGet, "/game", nil, params)
	if err != nil {
		return nil, err
	}

	var metadataList []igdb.Metadata

	err = proxyClient.do(req, &metadataList)
	if err != nil {
		return nil, err
	}

	return metadataList, nil
}

func (proxyClient *ProxyClient) GetMetadataByIGDBID(igdbId int) (*igdb.Metadata,error){

	req, err := proxyClient.newRequest(http.MethodGet, "/game/" + strconv.Itoa(igdbId), nil, nil)
	if err != nil {
		return nil, err
	}

	var metadata igdb.Metadata

	err = proxyClient.do(req, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}





