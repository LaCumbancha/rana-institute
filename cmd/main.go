package main

import (
	"os"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	services "github.com/LaCumbancha/institutional-site/cmd/services"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	host := os.Getenv("HOST")
	visitorCounter := services.NewVisitorService()
	homeService := services.NewHomeService(host, visitorCounter)
	jobsService := services.NewJobsService(host, visitorCounter)
	aboutService := services.NewAboutService(host, visitorCounter)
	legalService := services.NewLegalService(host, visitorCounter)

	http.HandleFunc("/", homeService.HomeHandler)
	http.HandleFunc("/jobs", jobsService.JobsHandler)
	http.HandleFunc("/about", aboutService.AboutHandler)
	http.HandleFunc("/about/legal", legalService.LegalHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s.", port)
	}

	log.Printf("Listening on port %s.", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Couldn't listen on port %d. Err: %s", port, err)
	}
}
