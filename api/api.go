package api

import (
	"github.com/loadimpact/speedboat/api/console"
	"github.com/loadimpact/speedboat/api/global"
	"github.com/loadimpact/speedboat/api/http"
)

type RegisterFunc func() map[string]interface***REMOVED******REMOVED***

var API = map[string]RegisterFunc***REMOVED***
	"global":  global.New,
	"console": console.New,
	"http":    http.New,
***REMOVED***

func New() map[string]map[string]interface***REMOVED******REMOVED*** ***REMOVED***
	res := make(map[string]map[string]interface***REMOVED******REMOVED***)
	for name, factory := range API ***REMOVED***
		res[name] = factory()
	***REMOVED***
	return res
***REMOVED***
