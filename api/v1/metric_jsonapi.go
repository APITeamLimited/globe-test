package v1

import (
	"time"

	"go.k6.io/k6/metrics"
)

// MetricsJSONAPI is JSON API envelop for metrics
type MetricsJSONAPI struct ***REMOVED***
	Data []metricData `json:"data"`
***REMOVED***

type metricJSONAPI struct ***REMOVED***
	Data metricData `json:"data"`
***REMOVED***

type metricData struct ***REMOVED***
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes Metric `json:"attributes"`
***REMOVED***

func newMetricEnvelope(m *metrics.Metric, t time.Duration) metricJSONAPI ***REMOVED***
	return metricJSONAPI***REMOVED***
		Data: newMetricData(m, t),
	***REMOVED***
***REMOVED***

func newMetricsJSONAPI(list map[string]*metrics.Metric, t time.Duration) MetricsJSONAPI ***REMOVED***
	metrics := make([]metricData, 0, len(list))

	for _, m := range list ***REMOVED***
		metrics = append(metrics, newMetricData(m, t))
	***REMOVED***

	return MetricsJSONAPI***REMOVED***
		Data: metrics,
	***REMOVED***
***REMOVED***

func newMetricData(m *metrics.Metric, t time.Duration) metricData ***REMOVED***
	metric := NewMetric(m, t)

	return metricData***REMOVED***
		Type:       "metrics",
		ID:         metric.Name,
		Attributes: metric,
	***REMOVED***
***REMOVED***

// Metrics extract the []v1.Metric from the JSON API envelop
func (m MetricsJSONAPI) Metrics() []Metric ***REMOVED***
	list := make([]Metric, 0, len(m.Data))

	for _, metric := range m.Data ***REMOVED***
		m := metric.Attributes
		m.Name = metric.ID
		list = append(list, m)
	***REMOVED***

	return list
***REMOVED***
