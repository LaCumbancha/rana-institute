package services

import (
	"fmt"
	"net/http"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	google "github.com/LaCumbancha/rana-institute/storage/cmd/google"
)

const HOME = "HOME"
const JOBS = "JOBS"
const ABOUT = "ABOUT"
const LEGAL = "LEGAL"

type VisitorService struct {
	datastoreClient 	*google.DatastoreClient
}

func NewVisitorService(datastoreClient *google.DatastoreClient) *VisitorService {
	return &VisitorService { datastoreClient }
}

func (service *VisitorService) VisitHandler(writer http.ResponseWriter, request *http.Request) {
	log.Infof("New message received at visits endpoint.")

	taskName := request.Header.Get("X-Appengine-Taskname")
	if taskName == "" {
		// You may use the presence of the X-Appengine-Taskname header to validate the request comes from Cloud Tasks.
		log.Warnf("Invalid Task: Task received with no X-Appengine-Taskname request header.")
		http.Error(writer, "Bad Request - Invalid Task", http.StatusBadRequest)
		return
	}

	// Pull useful headers from Task request.
	queueName := request.Header.Get("X-Appengine-Queuename")
	log.Debug("Task %s received from queue %s.", taskName, queueName)

	// Extract the request body for further task details.
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Warnf("Couldn't read task body. Err: %s", err)
		http.Error(writer, "Internal Error", http.StatusInternalServerError)
		return
	}

	page := string(body)
	log.Infof("New visit received from page %s.", page)

	if service.validatePage(page) {
		service.datastoreClient.UpdateVisits(page)
	} else {
		log.Warnf("Unknown page visit received: %s.", page)
	}

	// Log & output details of the task.
	output := fmt.Sprintf("Task %s from queue %s completed.", taskName, queueName)
	log.Infof(output)

	// Set a non-2xx status code to indicate a failure in task processing that should be retried.
	// For example, http.Error(writer, "Internal Server Error: Task Processing", http.StatusInternalServerError)
	fmt.Fprintln(writer, output)
}

func (service *VisitorService) validatePage(page string) bool {
	return page == HOME || page == JOBS || page == ABOUT || page == LEGAL
}
