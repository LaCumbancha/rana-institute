package main

import (
	"os"
	"net/http"

	log "github.com/sirupsen/logrus"
	services "github.com/LaCumbancha/rana-institute/infra/cmd/services"
)

func main() {
	projectName := os.Getenv("project")
	datastoreEntity := os.Getenv("datastore_entity")
	visitorCounter := services.NewVisitorService(projectName, datastoreEntity)
	
	indexService := services.NewIndexService()
	homeService := services.NewHomeService(visitorCounter)
	jobsService := services.NewJobsService(visitorCounter)
	aboutService := services.NewAboutService(visitorCounter)
	legalService := services.NewLegalService(visitorCounter)

	http.HandleFunc("/", indexService.IndexHandler)
	http.HandleFunc("/home", homeService.HomeHandler)
	http.HandleFunc("/jobs", jobsService.JobsHandler)
	http.HandleFunc("/about", aboutService.AboutHandler)
	http.HandleFunc("/about/legal", legalService.LegalHandler)

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
