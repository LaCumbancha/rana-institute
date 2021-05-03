package main

import (
	"os"
	"fmt"
	"strconv"
	"net/http"

	log "github.com/sirupsen/logrus"
	google "github.com/LaCumbancha/rana-institute/storage/cmd/google"
	services "github.com/LaCumbancha/rana-institute/storage/cmd/services"
)

const DEFAULT_PORT = "8080"
const DEFAULT_PARTITIONS = 10

func main() {
	projectId := os.Getenv("project_id")
	entity := os.Getenv("entity")

	partitions, err := strconv.Atoi(os.Getenv("partitions"))
	if err != nil {
		log.Errorf("Error retrieving the number of Datastore partitions. Defaulting at 10.")
		partitions = DEFAULT_PARTITIONS
	}

	datastoreClient := google.NewDatastoreClient(projectId, entity, partitions)
	visitorService := services.NewVisitorService(datastoreClient)
	http.HandleFunc("/register-visits", visitorService.RegisterVisitHandler)
	http.HandleFunc("/retrieve-visits/", visitorService.RetrieveVisitsHandler)

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
