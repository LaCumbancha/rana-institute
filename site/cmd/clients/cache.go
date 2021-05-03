package clients

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type CacheTransferData struct {
	Page 				string
	Visits 				int
}

type CacheClient struct {
	projectId			string
	serviceId 			string
}

func NewCacheClient(projectId string, serviceId string) *CacheClient {
	return &CacheClient { projectId, serviceId }
}

func (client *CacheClient) GetCachedVisitorCount(page string) int {
	if response, err := http.Get(client.getCacheEndpoint(page)); err != nil {
		log.Errorf("Error executing GET to the cache service. Err: %s", err)
		return -1
	} else {
		if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
			log.Infof("The information retrieved from the cache service was not valid. Status: %s.", response.Status)
			return -1
		}

		decoder := json.NewDecoder(response.Body)
		var cachedData CacheTransferData
		err := decoder.Decode(&cachedData)
		if err != nil {
			log.Warnf("There was an error mapping cache response body. Err: %s", err)
			return -1
		}

		return cachedData.Visits
	}
}

func (client *CacheClient) SetCachedVisitorCount(page string, visits int) {
	transferData := CacheTransferData { Page: page, Visits: visits }

	output, err := json.Marshal(transferData)
	if err != nil {
		log.Errorf("Error serializing caché updating data. Err: %s", err)
	} else {
		if _, err := http.Post(client.setCacheEndpoint(), "application/json", bytes.NewBuffer(output)); err != nil {
			log.Errorf("Error executin POST to the cache service. Err: %s", err)
		}
	}
}

func (client *CacheClient) getCacheEndpoint(page string) string {
	return fmt.Sprintf("https://%s-dot-%s.uc.r.appspot.com/get-visits/%s", client.serviceId, client.projectId, page)
}

func (client *CacheClient) setCacheEndpoint() string {
	return fmt.Sprintf("https://%s-dot-%s.uc.r.appspot.com/set-visits", client.serviceId, client.projectId)
}
