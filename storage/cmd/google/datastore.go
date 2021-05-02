package google

import (
	"context"
	"cloud.google.com/go/datastore"

	log "github.com/sirupsen/logrus"
)

type PageVisits struct {
	Visits 		int32
}

type DatastoreClient struct {
	dsClient 			*datastore.Client
	context				context.Context
	entity				string
}

func NewDatastoreClient(projectName string, entity string) *DatastoreClient {
	context := context.Background()
	dsClient, err := datastore.NewClient(context, datastore.DetectProjectID)
	if err != nil {
		log.Fatalf("Error connecting to Datastore. Err: %s", err)
	}

	log.Infof("New client for Datastore created.")
	return &DatastoreClient { dsClient, context, entity }
}

func (client *DatastoreClient) UpdateVisits(page string) error {
	log.Infof("Updating Datastore page %s visits.", page)

	var pageVisits PageVisits
	key := datastore.NameKey(client.entity, page, nil)

	_, err := client.dsClient.RunInTransaction(client.context, func(transaction *datastore.Transaction) error {

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
		return err
	}

	return nil
}
