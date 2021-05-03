package clients

import (
	"fmt"
	"net/http"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type StorageTransferData struct {
	Page 				string
	Visits 				int
}

type StorageClient struct {
	projectId			string
	serviceId 			string
}

func NewStorageClient(projectId string, serviceId string) *StorageClient {
	return &StorageClient { projectId, serviceId }
}

func (client *StorageClient) GetVisitorCount(page string) int {
	if response, err := http.Get(client.getEndpoint(page)); err != nil {
		log.Errorf("Error executing GET to the storage service. Err: %s", err)
		return -1
	} else {
		if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
			log.Infof("The information retrieved from the storage service was not valid. Status: %s.", response.Status)
			return -1
		}

		decoder := json.NewDecoder(response.Body)
		var storedData StorageTransferData
		err := decoder.Decode(&storedData)
		if err != nil {
			log.Warnf("There was an error mapping storage response body. Err: %s", err)
			return -1
		}

		return storedData.Visits
	}
}

func (client *StorageClient) getEndpoint(page string) string {
	return fmt.Sprintf("https://%s-dot-%s.uc.r.appspot.com/retrieve-visits/%s", client.serviceId, client.projectId, page)
}
