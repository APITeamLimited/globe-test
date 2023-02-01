package orchMetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/output/globetest"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

type metric struct {
	Contains string                   `json:"contains"`
	Type     workerMetrics.MetricType `json:"type"`
	Value    float64                  `json:"value"`
}

type wrappedMetric struct {
	metric
	name        string
	location    string
	subFraction float64
	childJobId  string
}

type locationSubJobs struct {
	childJobId     string
	wrappedMetrics []wrappedMetric
	collected      bool
}

type collectedMetricLocations struct {
	locations map[string][]locationSubJobs
}

type collectedInterval struct {
	flushCount       int
	collectedMetrics collectedMetricLocations
}

// Cached metrics are stored before being collated and sent
type cachedMetricsStore struct {
	// Each envelope in the map is a certain metric
	collectedIntervals []collectedInterval
	mu                 sync.RWMutex
	gs                 libOrch.BaseGlobalState

	childJobs     map[string]libOrch.ChildJobDistribution
	childJobCount int
}

var (
	_ libOrch.BaseMetricsStore = &cachedMetricsStore{}
)

func NewCachedMetricsStore(gs libOrch.BaseGlobalState) *cachedMetricsStore {
	return &cachedMetricsStore{
		gs:                 gs,
		collectedIntervals: make([]collectedInterval, 0),
		mu:                 sync.RWMutex{},
	}
}

func (store *cachedMetricsStore) InitMetricsStore(childJobs map[string]libOrch.ChildJobDistribution) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.childJobs = childJobs

	// Determine total number of child jobs
	for _, childJob := range childJobs {
		store.childJobCount += len(childJob.ChildJobs)
	}

	// This will never return an error
	// pf, _ := output.NewPeriodicFlusher(1000*time.Millisecond, store.FlushMetrics)

	// store.flusher = pf
}

func (store *cachedMetricsStore) AddMessage(message libOrch.WorkerMessage, workerLocation string, subFraction float64) error {
	if store.childJobs == nil {
		return errors.New("metrics store not initialised")
	}

	// TODO: Implement gzip decompression

	// // Gzip decompress

	// buf := bytes.NewBuffer([]byte(message.Message))

	// // Create a new flate reader
	// fr := flate.NewReader(buf)

	// // Read the decompressed bytes from the flate reader
	// decompressed, err := ioutil.ReadAll(fr)
	// if err != nil {
	// 	fmt.Printf("Error decompressing: %v", err)
	// 	return err
	// }

	// fmt.Printf("Time to decompress %v\n", time.Since(startTime))

	var wrappedFormattedSamples globetest.WrappedFormattedSamples
	err := json.Unmarshal([]byte(message.Message), &wrappedFormattedSamples)
	if err != nil {
		fmt.Printf("Error unmarshalling while adding message: %v", err)
		return err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	store.extendStoreIfRequired(wrappedFormattedSamples)

	store.addSamplesToStore(wrappedFormattedSamples, workerLocation, message.ChildJobId, subFraction)

	store.determineIfCanSendMetrics(wrappedFormattedSamples.FlushCount)

	store.cleanupIrretrievableMetrics()

	return nil
}

// Each flush count is a new set of metrics, and the array is extended to accommodate
func (store *cachedMetricsStore) extendStoreIfRequired(wrappedFormattedSamples globetest.WrappedFormattedSamples) {
	// Check last item in array
	if !(len(store.collectedIntervals) == 0 || store.collectedIntervals[len(store.collectedIntervals)-1].flushCount != wrappedFormattedSamples.FlushCount) {
		return
	}

	store.collectedIntervals = append(store.collectedIntervals, collectedInterval{
		flushCount: wrappedFormattedSamples.FlushCount,
		collectedMetrics: collectedMetricLocations{
			locations: make(map[string][]locationSubJobs),
		},
	})

	newEndIndex := len(store.collectedIntervals) - 1

	// Add load distribution locations
	for location, childJobDistribution := range store.childJobs {
		// Get number of child jobs for this location
		numChildJobs := len(childJobDistribution.ChildJobs)

		store.collectedIntervals[newEndIndex].collectedMetrics.locations[location] = make([]locationSubJobs, 0, numChildJobs)

		// Add sub jobs
		for _, childJob := range childJobDistribution.ChildJobs {
			store.collectedIntervals[newEndIndex].collectedMetrics.locations[location] = append(store.collectedIntervals[newEndIndex].collectedMetrics.locations[location], locationSubJobs{
				childJobId:     childJob.ChildJobId,
				wrappedMetrics: make([]wrappedMetric, 0),
				collected:      false,
			})
		}
	}
}

func (store *cachedMetricsStore) addSamplesToStore(wrappedFormattedSamples globetest.WrappedFormattedSamples, workerLocation string, childJobId string, subfraction float64) {
	// Find the correct interval
	for i := len(store.collectedIntervals) - 1; i >= 0; i-- {
		if store.collectedIntervals[i].flushCount == wrappedFormattedSamples.FlushCount {
			// Find the correct location
			for j := range store.collectedIntervals[i].collectedMetrics.locations[workerLocation] {
				if store.collectedIntervals[i].collectedMetrics.locations[workerLocation][j].childJobId == childJobId {
					// Add the samples
					locationSubJob := &store.collectedIntervals[i].collectedMetrics.locations[workerLocation][j]

					for _, sampleEnvelope := range wrappedFormattedSamples.SampleEnvelopes {
						locationSubJob.wrappedMetrics = append(locationSubJob.wrappedMetrics, wrappedMetric{
							metric: metric{
								Contains: sampleEnvelope.Metric.Contains.String(),
								Type:     sampleEnvelope.Metric.Type,
								Value:    sampleEnvelope.Data.Value,
							},
							name:        sampleEnvelope.Metric.Name,
							location:    workerLocation,
							subFraction: subfraction,
							childJobId:  childJobId,
						})
					}

					locationSubJob.collected = true

					return
				}
			}
		}
	}
}

// Checks to see if all locationSubJobs have been collected, if so, they can be agregated,
// sent, and the flush count can be removed
func (store *cachedMetricsStore) determineIfCanSendMetrics(flushCount int) {
	// Find the correct interval
	for i := len(store.collectedIntervals) - 1; i >= 0; i-- {
		if store.collectedIntervals[i].flushCount == flushCount {
			// Check if all locationSubJobs have been collected
			for _, locationSubJobs := range store.collectedIntervals[i].collectedMetrics.locations {
				for _, locationSubJob := range locationSubJobs {
					if !locationSubJob.collected {
						return
					}
				}
			}

			// All locationSubJobs have been collected, so we can send the metrics
			// TODO: Send metrics
			store.aggreagateAndSendMetrics(i)
			return
		}
	}
}

type possibleMetric struct {
	metricKey      string
	metricType     workerMetrics.MetricType
	metricContains string
}

func determinePossibleMetrics(interval collectedInterval) []possibleMetric {
	possibleMetrics := make([]possibleMetric, 0)

	for _, locationSubJobs := range interval.collectedMetrics.locations {
		for _, locationSubJob := range locationSubJobs {
			for _, wrappedMetric := range locationSubJob.wrappedMetrics {
				// Check if this metric has already been added
				metricAlreadyAdded := false

				for _, possibleMetric := range possibleMetrics {
					if possibleMetric.metricKey == wrappedMetric.name && possibleMetric.metricType == wrappedMetric.metric.Type && possibleMetric.metricContains == wrappedMetric.metric.Contains {
						metricAlreadyAdded = true
						break
					}
				}

				if !metricAlreadyAdded {
					possibleMetrics = append(possibleMetrics, possibleMetric{
						metricKey:      wrappedMetric.name,
						metricType:     wrappedMetric.metric.Type,
						metricContains: wrappedMetric.metric.Contains,
					})
				}
			}
		}
	}

	return possibleMetrics
}

func (store *cachedMetricsStore) aggreagateAndSendMetrics(intervalIndex int) {
	if store.gs.Standalone() {
		store.addGlobalLocation(intervalIndex)
	}

	interval := store.collectedIntervals[intervalIndex]

	// Combined metrics is map[location]map[metricKey]metric
	combinedMetrics := calculateCombinedMetrics(interval)

	// Send metrics
	go sendMetrics(store.gs, combinedMetrics)

	// Remove the interval
	store.collectedIntervals = append(store.collectedIntervals[:intervalIndex], store.collectedIntervals[intervalIndex+1:]...)
}

func sendMetrics(gs libOrch.BaseGlobalState, combinedMetrics map[string]map[string]metric) {
	// Marshall the envelopes
	marshalledCollectedMetrics, err := json.Marshal(combinedMetrics)
	if err != nil {
		libOrch.HandleError(gs, err)
		return
	}

	libOrch.DispatchMessage(gs, string(marshalledCollectedMetrics), "METRICS")
}

func calculateCombinedMetrics(interval collectedInterval) map[string]map[string]metric {
	// Determine metric types in this interval
	possibleMetrics := determinePossibleMetrics(interval)

	combinedMetrics := make(map[string]map[string]metric)
	for location := range interval.collectedMetrics.locations {
		combinedMetrics[location] = make(map[string]metric)
	}

	for location, locationSubJobs := range interval.collectedMetrics.locations {
		for _, possibleMetric := range possibleMetrics {
			// Find all metrics in this zone that match the key
			matchingKeyMetrics := make([]wrappedMetric, 0)
			for _, locationSubJob := range locationSubJobs {
				for _, wrappedMetric := range locationSubJob.wrappedMetrics {
					if wrappedMetric.name == possibleMetric.metricKey && wrappedMetric.metric.Type == possibleMetric.metricType && wrappedMetric.metric.Contains == possibleMetric.metricContains {
						matchingKeyMetrics = append(matchingKeyMetrics, wrappedMetric)
					}
				}
			}

			// Combine the metrics
			if possibleMetric.metricType == workerMetrics.Counter {
				combinedMetrics[location][possibleMetric.metricKey] = determineCounter(matchingKeyMetrics, possibleMetric.metricKey, possibleMetric.metricContains, workerMetrics.Counter)
			} else if possibleMetric.metricType == workerMetrics.Gauge {
				// Gauges are summed
				combinedMetrics[location][possibleMetric.metricKey] = determineCounter(matchingKeyMetrics, possibleMetric.metricKey, possibleMetric.metricContains, workerMetrics.Gauge)
			} else if possibleMetric.metricType == workerMetrics.Rate {
				// Rates are summed
				combinedMetrics[location][possibleMetric.metricKey] = determineCounter(matchingKeyMetrics, possibleMetric.metricKey, possibleMetric.metricContains, workerMetrics.Rate)
			} else if possibleMetric.metricType == workerMetrics.Trend {
				// Trends are summed
				combinedMetrics[location][possibleMetric.metricKey] = determineTrend(matchingKeyMetrics, possibleMetric.metricKey, possibleMetric.metricContains, workerMetrics.Trend, "mean")
			}
		}
	}

	return combinedMetrics
}

func (store *cachedMetricsStore) addGlobalLocation(intervalIndex int) {
	// Make a global location
	store.collectedIntervals[intervalIndex].collectedMetrics.locations["global"] = make([]locationSubJobs, 0, store.childJobCount)

	// Add all sub jobs to global location
	for location, locationSubJobs := range store.collectedIntervals[intervalIndex].collectedMetrics.locations {
		if location == "global" {
			continue
		}

		store.collectedIntervals[intervalIndex].collectedMetrics.locations["global"] = append(store.collectedIntervals[intervalIndex].collectedMetrics.locations["global"], locationSubJobs...)
	}
}

// Calculates an aggregated counter metric for a zone
func determineCounter(matchingKeyMetrics []wrappedMetric, metricName string, metricContains string,
	metricType workerMetrics.MetricType) metric {
	aggregatedMetric := metric{
		Contains: metricContains,
		Type:     metricType,
		Value:    0.0,
	}

	for _, zoneMetric := range matchingKeyMetrics {
		aggregatedMetric.Value += zoneMetric.Value
	}

	return aggregatedMetric
}

// Calculates a weighted mean value metric for a zone
func determineTrend(matchingKeyMetrics []wrappedMetric, metricName string, metricContains string,
	metricType workerMetrics.MetricType, valueKey string) metric {
	aggregatedMetric := metric{
		Contains: metricContains,
		Type:     metricType,
		Value:    0.0,
	}

	// Determine the weighted average of each value key from each metric
	subFractionTotal := 0.0

	if valueKey == "max" {
		// Find biggest value
		for _, zoneMetric := range matchingKeyMetrics {
			if zoneMetric.Value > aggregatedMetric.Value {
				aggregatedMetric.Value = zoneMetric.Value
			}
		}
	} else if valueKey == "min" {
		// Find smallest value
		aggregatedMetric.Value = math.MaxFloat64

		for _, zoneMetric := range matchingKeyMetrics {
			if zoneMetric.Value < aggregatedMetric.Value {
				aggregatedMetric.Value = zoneMetric.Value
			}
		}
	} else {
		// This isn't ideal for all all remaining value keys but better than nothing

		for _, zoneMetric := range matchingKeyMetrics {
			subFractionTotal += zoneMetric.subFraction
			aggregatedMetric.Value += zoneMetric.Value * zoneMetric.subFraction
		}

		if subFractionTotal > 0 {
			aggregatedMetric.Value /= subFractionTotal
		}
	}

	return aggregatedMetric
}

func (store *cachedMetricsStore) Cleanup() {
	store = nil
}

// If flushCounts lags considerably behind the leading flushCount, then we can assume that the flushCount is no longer retrievable
// and can be cleaned up, this is to prevent the memory from growing too large in the case that there is a lagging worker
func (store *cachedMetricsStore) cleanupIrretrievableMetrics() {
	// Find the leading flushCount
	leadingFlushCount := 0
	for _, interval := range store.collectedIntervals {
		if interval.flushCount > leadingFlushCount {
			leadingFlushCount = interval.flushCount
		}
	}

	indexesToRemove := make([]int, 0)

	// Remove all metrics that are more than 5 behind the leading flushCount
	for index, interval := range store.collectedIntervals {
		if interval.flushCount < leadingFlushCount-20 {
			indexesToRemove = append(indexesToRemove, index)
		}
	}

	for _, index := range indexesToRemove {
		if index < len(store.collectedIntervals) {
			store.collectedIntervals = append(store.collectedIntervals[:index], store.collectedIntervals[index+1:]...)
		} else {
			store.collectedIntervals = store.collectedIntervals[:index]
		}
	}
}
