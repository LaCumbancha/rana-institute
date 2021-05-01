package services

import(
	"net/http"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type IndexService struct {
	template			*template.Template
}

const INDEX_HTML_URL = "./html/index.html"

func NewIndexService() *IndexService {
	templ, err := template.ParseFiles(INDEX_HTML_URL)
	if err != nil {
		log.Fatalf("Coudn't load Index HTML. Err: %s", err)
	}

	return &IndexService { template: templ }
}

func (service *IndexService) IndexHandler(writer http.ResponseWriter, _ *http.Request) {
	if err := service.template.Execute(writer, nil); err != nil {
		log.Errorf("Error rendering Index HTML. Err: %s", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}
