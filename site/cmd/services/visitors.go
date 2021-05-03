package services

import (
	log "github.com/sirupsen/logrus"
	google "github.com/LaCumbancha/rana-institute/app/cmd/google"
	clients "github.com/LaCumbancha/rana-institute/app/cmd/clients"
)

type VisitorService struct {
	tasksProducer 		*google.TaskProducer
	cacheClient			*clients.CacheClient
	storageClient		*clients.StorageClient
}

func NewVisitorService(tasksProducer *google.TaskProducer, cacheClient *clients.CacheClient, storageClient *clients.StorageClient) *VisitorService {
	return &VisitorService { tasksProducer, cacheClient, storageClient }
}

func (service *VisitorService) UpdateVisitorCount(page string) int {
	log.Infof("Generating new task for page %s new visitor.", page)
	service.tasksProducer.RegisterNewVisit(page)
	
	log.Debugf("Retrieving visitors count for page %s from cache.", page)
	visits := service.cacheClient.GetCachedVisitorCount(page)

	if visits < 0 {
		log.Infof("Retrieving visitors count from permanent storage.")
		visits = service.storageClient.GetVisitorCount(page)

		if visits < 0 {
			log.Errorf("Couldn't retrieve appropiate visitors count for page %s from permanente storage.", page)
			return 0
		} else {
			log.Infof("Updating cache with new visitor counter (%d) for page %s.", visits, page)
			service.cacheClient.SetCachedVisitorCount(page, visits)
		}
	}

	return visits
}
