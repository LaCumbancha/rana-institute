package services

import(
	"net/http"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type HomeService struct {
	host 				string
	template			*template.Template
	visitorService		*VisitorService
}

type homeRenderizationData struct {
	URL 				string
	Visits 				int32
}

const HOME_HTML_URL = "./html/home.html"

func NewHomeService(host string, visitorService *VisitorService) *HomeService {
	templ, err := template.ParseFiles(HOME_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load Home HTML. Err: %s", err)
	}

	return &HomeService { template: templ, host: host, visitorService: visitorService }
}

func (service *HomeService) HomeHandler(writer http.ResponseWriter, _ *http.Request) {
	visitorNumber := service.visitorService.NewVisitor(HOME)

	renderData := homeRenderizationData { URL: service.host, Visits: visitorNumber }

	if err := service.template.Execute(writer, renderData); err != nil {
		log.Errorf("Error rendering home HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}
