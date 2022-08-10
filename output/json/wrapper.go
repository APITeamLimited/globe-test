package json

import (
	"time"

	"go.k6.io/k6/metrics"
)

//go:generate easyjson -pkg -no_std_marshalers -gen_build_flags -mod=mod .

//easyjson:json
type sampleEnvelope struct ***REMOVED***
	Type string `json:"type"`
	Data struct ***REMOVED***
		Time  time.Time           `json:"time"`
		Value float64             `json:"value"`
		Tags  *metrics.SampleTags `json:"tags"`
	***REMOVED*** `json:"data"`
	Metric string `json:"metric"`
***REMOVED***

// wrapSample is used to package a metric sample in a way that's nice to export
// to JSON.
func wrapSample(sample metrics.Sample) sampleEnvelope ***REMOVED***
	s := sampleEnvelope***REMOVED***
		Type:   "Point",
		Metric: sample.Metric.Name,
	***REMOVED***
	s.Data.Time = sample.Time
	s.Data.Value = sample.Value
	s.Data.Tags = sample.Tags
	return s
***REMOVED***

//easyjson:json
type metricEnvelope struct ***REMOVED***
	Type string `json:"type"`
	Data struct ***REMOVED***
		Name       string               `json:"name"`
		Type       metrics.MetricType   `json:"type"`
		Contains   metrics.ValueType    `json:"contains"`
		Thresholds metrics.Thresholds   `json:"thresholds"`
		Submetrics []*metrics.Submetric `json:"submetrics"`
	***REMOVED*** `json:"data"`
	Metric string `json:"metric"`
***REMOVED***
