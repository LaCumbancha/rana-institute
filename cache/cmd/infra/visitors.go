package infra

import (
	"fmt"
	"time"
	"sync"
	"net/http"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	utils "github.com/LaCumbancha/rana-institute/cache/cmd/utils"
)

type TransferData struct {
	Page 				string
	Visits 				int
}

type VisitorData struct {
	visits				int
	expiresAt			time.Time
}

type VisitorCache struct {
	TTL					utils.TTL
	data				map[string]VisitorData
	mutex				*sync.Mutex
}

func NewVisitorCache(TTL string) *VisitorCache {
	return &VisitorCache { TTL: utils.TimeToLive(TTL), data: make(map[string]VisitorData), mutex: &sync.Mutex{} }
}

func (cache *VisitorCache) GetVisitsHandler(writer http.ResponseWriter, request *http.Request) {
	page := utils.RouterHelper().GetPage(request)
	log.Infof("New GET request received for page %s.", page)

	now := time.Now()
	visits, cErr := cache.retrieveVisits(page, now)
	if cErr == EXPIRED {
		log.Infof("There was a value found for page %s visits (and it was %d) but it was expired.", page, visits)
		http.Error(writer, fmt.Sprintf("Visitors count for page %s expired", page), http.StatusGone)
		return
	} else if cErr == NOT_FOUND {
		log.Infof("The visitors value for page %s was not found.", page)
		http.Error(writer, fmt.Sprintf("Visitors count for page %s not found", page), http.StatusNotFound)
		return
	}
	
	visitorData := TransferData { Page: page, Visits: visits }
	output, err := json.Marshal(visitorData)
	if err != nil {
		log.Errorf("Error serializing visits data for page %s. Visitor count at %d.", visitorData.Page, visitorData.Visits)
		http.Error(writer, "Internal Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(writer, string(output))
}

func (cache *VisitorCache) SetVisitsHandler(writer http.ResponseWriter, request *http.Request) {
	log.Infof("New SET request received.")

	decoder := json.NewDecoder(request.Body)
	var inputData TransferData
	err := decoder.Decode(&inputData)
	if err != nil {
		log.Warnf("There was an error mapping request body. Err: %s", err)
		http.Error(writer, "Internal Error", http.StatusInternalServerError)
		return
	}

	if inputData.Page == "" {
		log.Warnf("There was no page received for the set visitors.")
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	cache.setVisits(inputData.Page, inputData.Visits)
	log.Infof("Cached visits value updated for page %s (setted at %d).", inputData.Page, inputData.Visits)

	fmt.Fprintln(writer, "OK")
}

func (cache *VisitorCache) retrieveVisits(page string, now time.Time) (int, CacheErrors) {
	if cachedVisits, found := cache.data[page]; found {
		if now.Before(cachedVisits.expiresAt) {
			log.Infof("Visitors count for page %s found at %d.", page, cachedVisits.visits)
			return cachedVisits.visits, OK
		} else {
			log.Infof("Cached visitor count for page %s expired %s ago.", page, cachedVisits.expiresAt.Sub(now).String)
			return cachedVisits.visits, EXPIRED
		}
	} else {
		log.Infof("Cached visitor count for page %s not found.", page)
		return 0, NOT_FOUND
	}
}

func (cache *VisitorCache) setVisits(page string, visits int) {
	newExpiration := time.Now().Add(time.Hour * cache.TTL.Hours + time.Minute * cache.TTL.Minutes + time.Second * cache.TTL.Seconds)
	visitorData := VisitorData { visits, newExpiration }
	cache.data[page] = visitorData
	log.Infof("Cached visitor count for page %s setted at %d with expiration at %s", page, visits, newExpiration.Format(time.RFC3339))
}
