package services

import(
	"net/http"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type HomeService struct {
	template			*template.Template
	visitorService		*VisitorService
}

type homeRenderizationData struct {
	Visits 				int32
}

const HOME_HTML_URL = "./html/home.html"

func NewHomeService(visitorService *VisitorService) *HomeService {
	templ, err := template.ParseFiles(HOME_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load Home HTML. Err: %s", err)
	}

	return &HomeService { template: templ, visitorService: visitorService }
}

func (service *HomeService) HomeHandler(writer http.ResponseWriter, _ *http.Request) {
	visitorNumber := service.visitorService.NewVisitor(HOME)

	renderData := homeRenderizationData { Visits: visitorNumber }

	if err := service.template.Execute(writer, renderData); err != nil {
		log.Errorf("Error rendering home HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}
