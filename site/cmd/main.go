package main

import (
	"os"
	"net/http"

	log "github.com/sirupsen/logrus"
	services "github.com/LaCumbancha/rana-institute/app/cmd/services"
)

func main() {
	taskQueue := os.Getenv("task_queue")
	visitorCounter := services.NewVisitorService(taskQueue)
	
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
