package main

import (
	"os"
	"net/http"

	log "github.com/sirupsen/logrus"
	services "github.com/LaCumbancha/rana-institute/infra/cmd/services"
)

func main() {
	projectName := os.Getenv("project_id")
	visitsEndpoint := os.Getenv("visits_endpoint")
	datastoreEntity := os.Getenv("datastore_entity")

	visitorService := services.NewVisitorService(projectName, datastoreEntity)
	http.HandleFunc(visitsEndpoint, visitorService.VisitHandler)

	port := os.Getenv("port")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s.", port)
	}

	log.Printf("Listening on port %s.", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Couldn't listen on port %d. Err: %s", port, err)
	}
}
