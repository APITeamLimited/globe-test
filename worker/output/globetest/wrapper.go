package globetest

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

type SampleData struct ***REMOVED***
	Time  time.Time                 `json:"time"`
	Value float64                   `json:"value"`
	Tags  *workerMetrics.SampleTags `json:"tags"`
***REMOVED***

// Stores a metric and data point
type SampleEnvelope struct ***REMOVED***
	Type   string                `json:"type"`
	Data   SampleData            `json:"data"`
	Metric *workerMetrics.Metric `json:"metric"`
***REMOVED***

// wrapSample is used to package a metric sample in a way that's nice to export
// to JSON.
func wrapSample(sample workerMetrics.Sample) SampleEnvelope ***REMOVED***
	s := SampleEnvelope***REMOVED***
		Type:   "Point",
		Metric: sample.Metric,
	***REMOVED***
	s.Data.Time = sample.Time
	s.Data.Value = sample.Value
	s.Data.Tags = sample.Tags
	return s
***REMOVED***
