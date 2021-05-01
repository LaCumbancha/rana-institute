package services

import (
	"context"
	"cloud.google.com/go/datastore"

	log "github.com/sirupsen/logrus"
)

const HOME = "HOME"
const JOBS = "JOBS"
const ABOUT = "ABOUT"
const LEGAL = "LEGAL"

type PageVisits struct {
	Visits 		int32
}

type VisitorService struct {
	client 				*datastore.Client
	context				context.Context
	entity				string
}

func NewVisitorService(projectName string, entity string) *VisitorService {
	context := context.Background()
	client, err := datastore.NewClient(context, datastore.DetectProjectID)
	if err != nil {
		log.Fatalf("Error connecting to Datastore. Err: %s", err)
	}

	log.Infof("New client for Datastore created.")
	return &VisitorService { client, context, entity }
}

func (service *VisitorService) HandleNewVisitor(page string) int32 {
	log.Infof("Requesting new visit for page %s.", page)

	pageVisits := PageVisits { Visits: -1 }
	key := datastore.NameKey(service.entity, page, nil)

	_, err := service.client.RunInTransaction(service.context, func(transaction *datastore.Transaction) error {

		if err := transaction.Get(key, &pageVisits); err != nil {
			log.Errorf("Error retrieving page %s visits counter. Err: %s", page, err)
			return err
		}

		pageVisits.Visits++
		if _, err := transaction.Put(key, &pageVisits); err != nil {
			log.Errorf("Error updating page %s visits counter. Err: %s", page, err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Errorf("There was an error running the visit counter transaction. Err: %s", err)
	}

	return pageVisits.Visits
}
