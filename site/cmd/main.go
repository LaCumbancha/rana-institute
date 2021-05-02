package main

import (
	"os"
	"net/http"

	log "github.com/sirupsen/logrus"
	google "github.com/LaCumbancha/rana-institute/app/cmd/google"
	services "github.com/LaCumbancha/rana-institute/app/cmd/services"
)

func main() {
	queueId := os.Getenv("queue_id")
	projectId := os.Getenv("project_id")
	locationId := os.Getenv("location_id")
	visitsEndpoint := os.Getenv("visits_endpoint")

	visitorsCache := google.NewVisitorsCache()
	tasksProducer := google.NewTaskProducer(projectId, locationId, queueId, visitsEndpoint)
	visitorService := services.NewVisitorService(tasksProducer, visitorsCache)
	
	indexService := services.NewIndexService()
	homeService := services.NewHomeService(visitorService)
	jobsService := services.NewJobsService(visitorService)
	aboutService := services.NewAboutService(visitorService)
	legalService := services.NewLegalService(visitorService)

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
