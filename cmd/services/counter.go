package services

type Page string
const HOME Page  = "HOME"
const JOBS Page = "JOBS"
const ABOUT Page = "ABOUT"
const LEGAL Page = "LEGAL"

type VisitorService struct {
	counters			map[Page]int32
}

func NewVisitorService() *VisitorService {
	return &VisitorService { counters: map[Page]int32{ HOME: 0, JOBS: 0, ABOUT: 0, LEGAL: 0 } }
}

func (service *VisitorService) NewVisitor(page Page) int32 {
	service.counters[page]++
	return service.counters[page]
}
