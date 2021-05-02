package main

import (
	"os"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	google "github.com/LaCumbancha/rana-institute/infra/cmd/google"
	services "github.com/LaCumbancha/rana-institute/infra/cmd/services"
)

func main() {
	projectId := os.Getenv("project_id")
	visitsEndpoint := os.Getenv("visits_endpoint")
	datastoreEntity := os.Getenv("datastore_entity")

	datastoreClient := google.NewDatastoreClient(projectId, datastoreEntity)
	visitorService := services.NewVisitorService(datastoreClient)
	http.HandleFunc(visitsEndpoint, visitorService.VisitHandler)



	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, "Hello, World!")
	})

	

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
