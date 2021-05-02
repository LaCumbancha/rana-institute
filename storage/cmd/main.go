package main

import (
	"os"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	google "github.com/LaCumbancha/rana-institute/storage/cmd/google"
	services "github.com/LaCumbancha/rana-institute/storage/cmd/services"
)

const DEFAULT_PORT = "8080"

func main() {
	projectId := os.Getenv("project_id")
	visitsEndpoint := os.Getenv("visits_endpoint")
	datastoreEntity := os.Getenv("datastore_entity")

	datastoreClient := google.NewDatastoreClient(projectId, datastoreEntity)
	visitorService := services.NewVisitorService(datastoreClient)
	http.HandleFunc(visitsEndpoint, visitorService.VisitHandler)

	port := os.Getenv("port")
	if port == "" {
		port = DEFAULT_PORT
		log.Infof("Defaulting to port %s.", port)
	}

	log.Printf("Listening on port %s.", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Couldn't listen on port %d. Err: %s", port, err)
	}
}
