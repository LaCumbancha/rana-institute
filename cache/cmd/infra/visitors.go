package infra

import (
	"fmt"
	"net/http"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	utils "github.com/LaCumbancha/rana-institute/cache/cmd/utils"
)

type VisitorData struct {
	Page 				string
	Visits 				int
}

type VisitorCache struct {
	visits				map[string]int
}

func NewVisitorCache() *VisitorCache {
	return &VisitorCache { visits: make(map[string]int) }
}

func (cache *VisitorCache) NewVisitHandler(writer http.ResponseWriter, request *http.Request) {
	page := utils.RouterHelper().GetPage(request)
	log.Infof("New visitor received at page %s.", page)

	visitorData := VisitorData { Page: page, Visits: cache.updateVisits(page) }
	output, err := json.Marshal(visitorData)
	if err != nil {
		log.Errorf("Error serializing visits data for page %s. Visitor counter at %d.", visitorData.Page, visitorData.Visits)
		http.Error(writer, "Internal Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(writer, string(output))
}

func (cache *VisitorCache) updateVisits(page string) int {
	// TODO: Replace by proper CACHE
	if previousVisits, found := cache.visits[page]; found {
		log.Infof("Visit counter for page %s found at %d. Increasing by 1.", page, previousVisits)
		cache.visits[page]++
	} else {
		log.Infof("Visit counter for page %s not found. Defaulting at 1.", page)
		cache.visits[page] = 1
	}

	return cache.visits[page]
}
