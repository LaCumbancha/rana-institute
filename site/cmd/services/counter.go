package services

import (
	log "github.com/sirupsen/logrus"
)

const HOME = "HOME"
const JOBS = "JOBS"
const ABOUT = "ABOUT"
const LEGAL = "LEGAL"

type VisitorService struct {
	visits 				map[string]int64
	taskQueue 			string
}

func NewVisitorService(taskQueue string) *VisitorService {
	return &VisitorService { visits: make(map[string]int64), taskQueue: taskQueue }
}

func (service *VisitorService) HandleNewVisitor(page string) int64 {
	
	// TODO: Insert TASK into TASK QUEUE
	log.Infof("Sending new visit for page %s to the task queue %s.", page, service.taskQueue)
	
	// TODO: Replace by CACHE
	if previousVisits, found := service.visits[page]; found {
		log.Infof("Visit counter for page %s found at %d. Increasing by 1.", page, previousVisits)
		service.visits[page]++
	} else {
		log.Infof("Visit counter for page %s not found. Defaulting at 1.", page, previousVisits)
		service.visits[page] = 1
	}

	return service.visits[page]
}
