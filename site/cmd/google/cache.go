package google

import (
	log "github.com/sirupsen/logrus"
)

type VisitorsCache struct {
	visits				map[string]int
}

func NewVisitorsCache() *VisitorsCache {
	return &VisitorsCache { visits: make(map[string]int) }
}

func (cache *VisitorsCache) UpdateVisits(page string) int {
	// TODO: Replace by proper CACHE
	if previousVisits, found := cache.visits[page]; found {
		log.Infof("Visit counter for page %s found at %d. Increasing by 1.", page, previousVisits)
		cache.visits[page]++
	} else {
		log.Infof("Visit counter for page %s not found. Defaulting at 1.", page, previousVisits)
		cache.visits[page] = 1
	}

	return cache.visits[page]
}
