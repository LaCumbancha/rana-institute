package services

import (
	log "github.com/sirupsen/logrus"
	google "github.com/LaCumbancha/rana-institute/app/cmd/google"
)

const HOME = "HOME"
const JOBS = "JOBS"
const ABOUT = "ABOUT"
const LEGAL = "LEGAL"

type VisitorService struct {
	tasksProducer 		*google.TaskProducer
	visitorsCache		*google.VisitorsCache
}

func NewVisitorService(tasksProducer *google.TaskProducer, visitorsCache *google.VisitorsCache) *VisitorService {
	return &VisitorService { tasksProducer, visitorsCache }
}

func (service *VisitorService) HandleNewVisitor(page string) int {
	log.Infof("Generating new task for page %s new visitor.", page)
	service.tasksProducer.RegisterNewVisit(page)
	
	log.Debugf("Retrieving visits count for page %s from caché.", page)
	visits := service.visitorsCache.UpdateVisits(page)
	log.Infof("Visits count for page %s at %d retrieved from caché.", page, visits)

	return visits
}
