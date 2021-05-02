package google

import (
	"fmt"
	"context"

	log "github.com/sirupsen/logrus"
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
   	publisher "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

type TaskProducer struct {
	client 				*cloudtasks.Client
	context				context.Context
	queueId				string
	queuePath			string
	endpoint			string
}

func NewTaskProducer(projectId string, locationId string, queueId string, endpoint string) *TaskProducer {
	context := context.Background()
	client, err := cloudtasks.NewClient(context)
	if err != nil {
		log.Fatalf("Error connecting to CloudTasks. Err: %s", err)
	}

	log.Infof("New client for CloudTasks created.")

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectId, locationId, queueId)
	log.Debugf("Tasks queue path defined as '%s'", queuePath)

	return &TaskProducer { client, context, queueId, queuePath, endpoint }
}

func (producer *TaskProducer) RegisterNewVisit(page string) {
	request := &publisher.CreateTaskRequest {
		Parent: producer.queuePath,
		Task: &publisher.Task {
			MessageType: &publisher.Task_AppEngineHttpRequest {
				AppEngineHttpRequest: &publisher.AppEngineHttpRequest {
					HttpMethod:  publisher.HttpMethod_POST,
					RelativeUri: producer.endpoint,
				},
			},
		},
	}
	log.Debugf("New task request generated for page %s and for endpoint %s.", page, producer.endpoint)

	// Add the page as a a payload message.
	request.Task.GetAppEngineHttpRequest().Body = []byte(page)

	if createdTask, err := producer.client.CreateTask(producer.context, request); err != nil {
		log.Errorf("Error publishing page %s task to queue %s. Err: %s", page, producer.queueId, err)
	} else {
		log.Infof("Page %s task published to queue %s with id '%s'.", page, producer.queueId, createdTask.Name)
	}
}
