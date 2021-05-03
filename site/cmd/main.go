package main

import (
	"os"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	google "github.com/LaCumbancha/rana-institute/app/cmd/google"
	clients "github.com/LaCumbancha/rana-institute/app/cmd/clients"
	services "github.com/LaCumbancha/rana-institute/app/cmd/services"
)

const DEFAULT_PORT = "8080"

func main() {
	queueId := os.Getenv("queue_id")
	projectId := os.Getenv("project_id")
	locationId := os.Getenv("location_id")
	cacheId := os.Getenv("cache_service")
	storageId := os.Getenv("storage_service")

	cacheClient := clients.NewCacheClient(projectId, cacheId)
	storageClient := clients.NewStorageClient(projectId, storageId)
	tasksProducer := google.NewTaskProducer(projectId, locationId, queueId)
	visitorService := services.NewVisitorService(tasksProducer, cacheClient, storageClient)
	
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
		port = DEFAULT_PORT
		log.Infof("Defaulting to port %s.", port)
	}

	log.Printf("Listening on port %s.", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal("Couldn't listen on port %d. Err: %s", port, err)
	}
}
