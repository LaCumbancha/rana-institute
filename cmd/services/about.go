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
	Visits 				int32
}

const ABOUT_HTML_URL = "./html/about.html"

func NewAboutService(visitorService *VisitorService) *AboutService {
	templ, err := template.ParseFiles(ABOUT_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load about HTML. Err: %s", err)
	}

	return &AboutService { template: templ, visitorService: visitorService }
}

func (service *AboutService) AboutHandler(writer http.ResponseWriter, _ *http.Request) {
	visitorNumber := service.visitorService.NewVisitor(ABOUT)

	renderData := aboutRenderizationData { Visits: visitorNumber }

	if err := service.template.Execute(writer, renderData); err != nil {
		log.Errorf("Error rendering about HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}