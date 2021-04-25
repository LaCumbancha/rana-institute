package services

import(
	"net/http"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type LegalService struct {
	template			*template.Template
	visitorService		*VisitorService
}

type legalRenderizationData struct {
	Visits 				int32
}

const LEGAL_HTML_URL = "./html/legal.html"

func NewLegalService(visitorService *VisitorService) *LegalService {
	templ, err := template.ParseFiles(LEGAL_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load legal HTML. Err: %s", err)
	}

	return &LegalService { template: templ, visitorService: visitorService }
}

func (service *LegalService) LegalHandler(writer http.ResponseWriter, _ *http.Request) {
	visitorNumber := service.visitorService.NewVisitor(LEGAL)

	renderData := legalRenderizationData { Visits: visitorNumber }

	if err := service.template.Execute(writer, renderData); err != nil {
		log.Errorf("Error rendering Legal HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}
