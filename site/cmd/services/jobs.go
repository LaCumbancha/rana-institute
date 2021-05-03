package services

import(
	"net/http"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type JobsService struct {
	template			*template.Template
	visitorService		*VisitorService
}

type jobsRenderizationData struct {
	Visits 				int
}

const JOBS_PAGE = "JOBS"
const JOBS_HTML_URL = "./html/jobs.html"

func NewJobsService(visitorService *VisitorService) *JobsService {
	templ, err := template.ParseFiles(JOBS_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load Jobs HTML. Err: %s", err)
	}

	return &JobsService { template: templ, visitorService: visitorService }
}

func (service *JobsService) JobsHandler(writer http.ResponseWriter, _ *http.Request) {
	visitorNumber := service.visitorService.UpdateVisitorCount(JOBS_PAGE)

	renderData := jobsRenderizationData { Visits: visitorNumber }

	if err := service.template.Execute(writer, renderData); err != nil {
		log.Errorf("Error rendering Jobs HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}
