package services

import(
	"net/http"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type LegalService struct {
	host 				string
	template			*template.Template
	visitorService		*VisitorService
}

type legalRenderizationData struct {
	URL 				string
	Visits 				int32
}

const LEGAL_HTML_URL = "./html/legal.html"

func NewLegalService(host string, visitorService *VisitorService) *LegalService {
	templ, err := template.ParseFiles(LEGAL_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load legal HTML. Err: %s", err)
	}

	return &LegalService { template: templ, host: host, visitorService: visitorService }
}

func (service *LegalService) LegalHandler(writer http.ResponseWriter, _ *http.Request) {
	visitorNumber := service.visitorService.NewVisitor(LEGAL)

	renderData := legalRenderizationData { URL: service.host, Visits: visitorNumber }

	if err := service.template.Execute(writer, renderData); err != nil {
		log.Errorf("Error rendering Legal HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}
