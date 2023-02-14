package globetest

// import (
// 	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
// )

// // Takes sample envelopes and aggregates them into a single envelope
// func aggregateSamples(samples []workerMetrics.Sample) []workerMetrics.Sample {
// 	newSamples := make([]workerMetrics.Sample, 0)

// 	for index, sample := range samples {
// 		if index == 0 {
// 			newSamples = append(newSamples, sample)
// 			continue
// 		}

// 		matched := false
// 		for i, newSample := range newSamples {
// 			if sample.Metric.Name == newSample.Metric.Name {
// 				if sample.Data.Tags != nil && newSample.Data.Tags != nil {
// 					name, _ := sample.Data.Tags.Get("name")
// 					newName, _ := newSample.Data.Tags.Get("name")

// 					if name != newName {
// 						continue
// 					}
// 				}

// 				if newSample.Metric.Type == workerMetrics.Counter || newSample.Metric.Type == workerMetrics.Rate {
// 					newSamples[i].Data.Value += sample.Data.Value
// 				} else if newSample.Metric.Type == workerMetrics.Gauge {
// 					if newSample.Data.Value > sample.Data.Value {
// 						newSamples[i].Data.Value = sample.Data.Value
// 					}
// 				} else if newSample.Metric.Type == workerMetrics.Trend {
// 					newSamples[i].Data.Value += sample.Data.Value
// 				}

// 				newSamples[i].samplesAdded++
// 				matched = true
// 				break
// 			}
// 		}

// 		if !matched {
// 			newSamples = append(newSamples, sample)
// 		}
// 	}

// 	for index, sample := range newSamples {
// 		if sample.Metric.Type == workerMetrics.Trend && sample.samplesAdded > 0 {
// 			newSamples[index].Data.Value /= float64(sample.samplesAdded)
// 		}
// 	}

// 	return newSamples
// }
