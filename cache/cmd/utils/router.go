package utils

import (
	"sync"
	"strings"
	"net/http"
)

type RouterHelperT struct {
	updateEndpoint		string
}

var helper *RouterHelperT
var once sync.Once

func RouterHelper() *RouterHelperT {
	once.Do(func() { helper = &RouterHelperT{ updateEndpoint: "/get-visits/" } })
	return helper
}

func (helper *RouterHelperT) GetPage(request *http.Request) string {
	return strings.TrimPrefix(request.URL.Path, helper.updateEndpoint)
}
