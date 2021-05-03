package google

import (
	"fmt"
	"context"
	"math/rand"
	"cloud.google.com/go/datastore"

	log "github.com/sirupsen/logrus"
)

type StoredVisits struct {
	Page				string
	Visits 				int
}

type DatastoreClient struct {
	dsClient 			*datastore.Client
	context				context.Context
	entity				string
	partitions			int
}

const DATASTORE_FILTER_KEY = "Page"

func NewDatastoreClient(projectId string, entity string, partitions int) *DatastoreClient {
	context := context.Background()
	dsClient, err := datastore.NewClient(context, projectId)
	if err != nil {
		log.Fatalf("Error connecting to Datastore. Err: %s", err)
	}

	log.Infof("New client for Datastore created.")
	return &DatastoreClient { dsClient, context, entity, partitions }
}

func (client *DatastoreClient) RegisterNewVisitor(page string) error {
	log.Infof("Updating Datastore page %s visits.", page)

	_, err := client.dsClient.RunInTransaction(client.context, func(transaction *datastore.Transaction) error {
		var storedVisits StoredVisits
		entityId := client.randomKeyName(page)
		key := datastore.NameKey(client.entity, entityId, nil)

		if err := transaction.Get(key, &storedVisits); err == datastore.ErrNoSuchEntity {
			log.Infof("Entity %s not found. Initializing with visitors count at 1", entityId)
			storedVisits = StoredVisits { Page: page, Visits: 1 }
		} else if err != nil {
			log.Errorf("Error retrieving page %s visits counter. Err: %s", page, err)
			return err
		} else {
			log.Infof("Entity %s retrieved with visitors count at %d. Incrementiny by 1.", entityId, storedVisits.Visits)
			storedVisits.Visits++
		}

		if _, err := transaction.Put(key, &storedVisits); err != nil {
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

func (client *DatastoreClient) RetrieveVisitorCount(page string) (int, error) {
	log.Infof("Retrieving Datastore page %s visits.", page)

	var storedVisits []StoredVisits
	query := datastore.NewQuery(client.entity).Filter(fmt.Sprintf("%s = ", DATASTORE_FILTER_KEY), page)
	if _, err := client.dsClient.GetAll(client.context, query, &storedVisits); err != nil {
		log.Errorf("Couldn't retrieve the page %s visitor count.", page)
		return 0, err
	}

	maxVisits := 0
	for _, storedVisit := range storedVisits {
		maxVisits += storedVisit.Visits
    }

    return maxVisits, nil
}

func (client *DatastoreClient) randomKeyName(page string) string {
	shard := rand.Intn(client.partitions)
	return fmt.Sprintf("%s-%d", page, shard)
}
