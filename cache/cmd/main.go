package main

import (
	"os"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	infra "github.com/LaCumbancha/rana-institute/cache/cmd/infra"
)

const DEFAULT_PORT = "8080"

func main() {
	visitorsCache := infra.NewVisitorCache()
	http.HandleFunc("/update-resource/", visitorsCache.NewVisitHandler)
	
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
