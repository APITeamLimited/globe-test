package globetest

import (
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

type SampleData struct {
	// Type is point
	Value float64                   `json:"value"`
	Tags  *workerMetrics.SampleTags `json:"tags"`
}

// Stores a metric and data point
type SampleEnvelope struct {
	Data         SampleData            `json:"data"`
	Metric       *workerMetrics.Metric `json:"metric"`
	samplesAdded int
}

// wrapSample is used to package a metric sample in a way that's nice to export
// to JSON and be easily aggregated in the orchestrator, unnecessary fields are
// removed and the metric name is split into its parts.
func wrapSample(sample workerMetrics.Sample) SampleEnvelope {
	s := SampleEnvelope{
		Metric: sample.Metric,
	}
	s.Data.Value = sample.Value
	s.Data.Tags = sample.Tags
	return s
}
