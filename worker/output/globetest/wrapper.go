package globetest

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

type SampleData struct {
	Time  time.Time                 `json:"time"`
	Value float64                   `json:"value"`
	Tags  *workerMetrics.SampleTags `json:"tags"`
}

// Stores a metric and data point
type SampleEnvelope struct {
	Type   string                `json:"type"`
	Data   SampleData            `json:"data"`
	Metric *workerMetrics.Metric `json:"metric"`
}

// wrapSample is used to package a metric sample in a way that's nice to export
// to JSON.
func wrapSample(sample workerMetrics.Sample) SampleEnvelope {
	s := SampleEnvelope{
		Type:   "Point",
		Metric: sample.Metric,
	}
	s.Data.Time = sample.Time
	s.Data.Value = sample.Value
	s.Data.Tags = sample.Tags
	return s
}
