package services

import (
	"context"
	"cloud.google.com/go/datastore"

	log "github.com/sirupsen/logrus"
)

const HOME = "HOME"
const JOBS = "JOBS"
const ABOUT = "ABOUT"
const LEGAL = "LEGAL"

type PageVisits struct {
	Visits 		int32
}

type VisitorService struct {
	client 				*datastore.Client
	context				context.Context
	entity				string
}

func NewVisitorService(projectName string, entity string) *VisitorService {
	context := context.Background()
	client, err := datastore.NewClient(context, datastore.DetectProjectID)
	if err != nil {
		log.Fatalf("Error connecting to Datastore. Err: %s", err)
	}

	log.Infof("New client for Datastore created.")
	return &VisitorService { client, context, entity }
}

func (service *VisitorService) VisitHandler(writer http.ResponseWriter, request *http.Request) {
	taskName := request.Header.Get("X-Appengine-Taskname")
	if taskName == "" {
		// You may use the presence of the X-Appengine-Taskname header to validate
		// the request comes from Cloud Tasks.
		log.Warnf("Invalid Task: Task received with no X-Appengine-Taskname request header.")
		http.Error(writer, "Bad Request - Invalid Task", http.StatusBadRequest)
		return
	}

	// Pull useful headers from Task request.
	queueName := request.Header.Get("X-Appengine-Queuename")

	// Extract the request body for further task details.
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("ReadAll: %v", err)
		http.Error(writer, "Internal Error", http.StatusInternalServerError)
		return
	}

	// Log & output details of the task.
	output := fmt.Sprintf("Completed task: task queue(%s), task name(%s), payload(%s)",
		queueName,
		taskName,
		string(body),
	)
	log.Println(output)

	// Set a non-2xx status code to indicate a failure in task processing that should be retried.
	// For example, http.Error(writer, "Internal Server Error: Task Processing", http.StatusInternalServerError)
	fmt.Fprintln(writer, output)
}
