package services

import(
	"net/http"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type JobsService struct {
	host 				string
	template			*template.Template
	visitorService		*VisitorService
}

type jobsRenderizationData struct {
	URL 				string
	Visits 				int32
}

const JOBS_HTML_URL = "./html/jobs.html"

func NewJobsService(host string, visitorService *VisitorService) *JobsService {
	templ, err := template.ParseFiles(JOBS_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load jobs HTML. Err: %s", err)
	}

	return &JobsService { template: templ, host: host, visitorService: visitorService }
}

func (service *JobsService) JobsHandler(writer http.ResponseWriter, _ *http.Request) {
	visitorNumber := service.visitorService.NewVisitor(JOBS)

	renderData := jobsRenderizationData { URL: service.host, Visits: visitorNumber }

	if err := service.template.Execute(writer, renderData); err != nil {
		log.Errorf("Error rendering jobs HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}
