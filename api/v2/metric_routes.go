package v2

import (
	"github.com/julienschmidt/httprouter"
	"github.com/loadimpact/k6/api/common"
	"github.com/manyminds/api2go/jsonapi"
	"net/http"
)

func HandleGetMetrics(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	engine := common.GetEngine(r.Context())

	metrics := make([]Metric, 0)
	for m, _ := range engine.Metrics ***REMOVED***
		metrics = append(metrics, NewMetric(*m))
	***REMOVED***

	data, err := jsonapi.Marshal(metrics)
	if err != nil ***REMOVED***
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***

func HandleGetMetric(rw http.ResponseWriter, r *http.Request, p httprouter.Params) ***REMOVED***
	id := p.ByName("id")
	engine := common.GetEngine(r.Context())

	var metric Metric
	var found bool
	for m, _ := range engine.Metrics ***REMOVED***
		if m.Name == id ***REMOVED***
			metric = NewMetric(*m)
			found = true
			break
		***REMOVED***
	***REMOVED***

	if !found ***REMOVED***
		http.Error(rw, "No such metric", http.StatusNotFound)
		return
	***REMOVED***

	data, err := jsonapi.Marshal(metric)
	if err != nil ***REMOVED***
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***
