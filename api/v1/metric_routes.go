package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"go.k6.io/k6/api/common"
)

func handleGetMetrics(rw http.ResponseWriter, r *http.Request) ***REMOVED***
	engine := common.GetEngine(r.Context())

	var t time.Duration
	if engine.ExecutionScheduler != nil ***REMOVED***
		t = engine.ExecutionScheduler.GetState().GetCurrentTestRunDuration()
	***REMOVED***

	engine.MetricsEngine.MetricsLock.Lock()
	metrics := newMetricsJSONAPI(engine.MetricsEngine.ObservedMetrics, t)
	engine.MetricsEngine.MetricsLock.Unlock()

	data, err := json.Marshal(metrics)
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***

func handleGetMetric(rw http.ResponseWriter, r *http.Request, id string) ***REMOVED***
	engine := common.GetEngine(r.Context())

	var t time.Duration
	if engine.ExecutionScheduler != nil ***REMOVED***
		t = engine.ExecutionScheduler.GetState().GetCurrentTestRunDuration()
	***REMOVED***

	engine.MetricsEngine.MetricsLock.Lock()
	metric, ok := engine.MetricsEngine.ObservedMetrics[id]
	if !ok ***REMOVED***
		engine.MetricsEngine.MetricsLock.Unlock()
		apiError(rw, "Not Found", "No metric with that ID was found", http.StatusNotFound)
		return
	***REMOVED***
	wrappedMetric := newMetricEnvelope(metric, t)
	engine.MetricsEngine.MetricsLock.Unlock()

	data, err := json.Marshal(wrappedMetric)
	if err != nil ***REMOVED***
		apiError(rw, "Encoding error", err.Error(), http.StatusInternalServerError)
		return
	***REMOVED***
	_, _ = rw.Write(data)
***REMOVED***
