package services

import(
	"net/http"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type AboutService struct {
	template			*template.Template
	visitorService		*VisitorService
}

type aboutRenderizationData struct {
	Visits 				int
}

const ABOUT_HTML_URL = "./html/about.html"

func NewAboutService(visitorService *VisitorService) *AboutService {
	templ, err := template.ParseFiles(ABOUT_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load About HTML. Err: %s", err)
	}

	return &AboutService { template: templ, visitorService: visitorService }
}

func (service *AboutService) AboutHandler(writer http.ResponseWriter, _ *http.Request) {
	visitorNumber := service.visitorService.HandleNewVisitor(ABOUT)

	renderData := aboutRenderizationData { Visits: visitorNumber }

	if err := service.template.Execute(writer, renderData); err != nil {
		log.Errorf("Error rendering About HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}
