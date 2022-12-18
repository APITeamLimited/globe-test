package orchMetrics

import (
	"encoding/json"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/output"
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

// Cached metrics are stored before being collated and sent
type cachedMetricsStore struct {
	// Each envelope in the map is a certain metric
	collectedMetrics map[string][]wrappedMetric
	mu               sync.RWMutex
	flusher          *output.PeriodicFlusher
	gs               libOrch.BaseGlobalState
	options          *libWorker.Options
}

var (
	_ libOrch.BaseMetricsStore = &cachedMetricsStore{}
)

func NewCachedMetricsStore(gs libOrch.BaseGlobalState) *cachedMetricsStore {
	return &cachedMetricsStore{
		gs:               gs,
		collectedMetrics: make(map[string][]wrappedMetric),
		mu:               sync.RWMutex{},
	}
}

func (store *cachedMetricsStore) InitMetricsStore(options *libWorker.Options) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.options = options

	// Add load distribution locations
	for _, location := range options.LoadDistribution.Value {
		store.collectedMetrics[location.Location] = make([]wrappedMetric, 0)
	}

	// Add global location
	store.collectedMetrics[libOrch.GlobalName] = make([]wrappedMetric, 0)

	// This will never return an error
	pf, _ := output.NewPeriodicFlusher(1000*time.Millisecond, store.FlushMetrics)

	store.flusher = pf
}

func (store *cachedMetricsStore) AddMessage(message libOrch.WorkerMessage, workerLocation string, subFraction float64) error {
	if store.options == nil {
		return errors.New("metrics store not initialised")
	}

	// Ensure childJobId is not already in the store, prevents duplicate metrics
	for _, existingWrappedMetric := range store.collectedMetrics[workerLocation] {
		if existingWrappedMetric.childJobId == message.ChildJobId {
			return nil
		}
	}

	var sampleEnvelopes []globetest.SampleEnvelope

	err := json.Unmarshal([]byte(message.Message), &sampleEnvelopes)
	if err != nil {
		return err
	}

	wrappedMetrics := make([]wrappedMetric, 0)
	for _, sampleEnvelope := range sampleEnvelopes {
		wrappedMetrics = append(wrappedMetrics, wrappedMetric{
			metric: metric{
				Contains: sampleEnvelope.Metric.Contains.String(),
				Type:     sampleEnvelope.Metric.Type,
				Value:    sampleEnvelope.Data.Value,
			},
			name:        sampleEnvelope.Metric.Name,
			location:    workerLocation,
			subFraction: subFraction,
			childJobId:  message.ChildJobId,
		})
	}

	// Add the metrics to the store
	store.mu.Lock()
	defer store.mu.Unlock()

	for _, wrappedMetric := range wrappedMetrics {
		// Add the metrics to the correct location array
		store.collectedMetrics[workerLocation] = append(store.collectedMetrics[workerLocation], wrappedMetric)

		// Add the metrics to the global array
		store.collectedMetrics[libOrch.GlobalName] = append(store.collectedMetrics[libOrch.GlobalName], wrappedMetric)
	}

	return nil
}

// Empty the store and returns its contents
func (store *cachedMetricsStore) emptyStore() (map[string][]wrappedMetric, error) {
	if store.options == nil {
		return nil, errors.New("metrics store not initialised")
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	// Copy the map
	result := make(map[string][]wrappedMetric, len(store.collectedMetrics))
	for metricName, sampleEnvelopes := range store.collectedMetrics {
		result[metricName] = make([]wrappedMetric, len(sampleEnvelopes))
		copy(result[metricName], sampleEnvelopes)
	}

	// Empty the map
	store.collectedMetrics = make(map[string][]wrappedMetric)

	// Add load distribution locations
	for _, location := range store.options.LoadDistribution.Value {
		store.collectedMetrics[location.Location] = make([]wrappedMetric, 0)
	}

	// Add global location
	store.collectedMetrics[libOrch.GlobalName] = make([]wrappedMetric, 0)

	return result, nil
}

func (store *cachedMetricsStore) FlushMetrics() {
	collectedMetrics, err := store.getMetrics()
	if err != nil {
		// Sometimes there are not enough metrics, this will throw an expected error
		// that can be ignored
		return
	}

	if len(collectedMetrics) == 0 {
		return
	}

	// Marshall the envelopes
	marshalledCollectedMetrics, err := json.Marshal(collectedMetrics)
	if err != nil {
		libOrch.HandleError(store.gs, err)
		return
	}

	libOrch.DispatchMessage(store.gs, string(marshalledCollectedMetrics), "METRICS")
}

func (store *cachedMetricsStore) getMetrics() (map[string]map[string]metric, error) {
	collectedMetrics, err := store.emptyStore()
	if err != nil {
		return nil, err
	}

	if len(collectedMetrics) == 0 {
		return nil, errors.New("need at least one group of metrics")
	}

	firstKey := ""

	// Ensure at least one metric in each location
	for key, zoneMetrics := range collectedMetrics {
		if len(zoneMetrics) == 0 {
			return nil, errors.New("need at least one metric in each location")
		}

		if firstKey == "" {
			firstKey = key
		}
	}

	metricKeys := make([]string, 0)
	metricTypes := make([]workerMetrics.MetricType, 0)
	metricContains := make([]string, 0)

	// Find first metric in first location
	for _, metric := range collectedMetrics[firstKey] {
		alreadExists := false
		for i := 0; i < len(metricKeys); i++ {
			if metricKeys[i] == metric.name {
				alreadExists = true
				break
			}
		}

		if !alreadExists {
			metricKeys = append(metricKeys, metric.name)
			metricTypes = append(metricTypes, metric.metric.Type)
			metricContains = append(metricContains, metric.metric.Contains)
		}
	}

	// Combined metrics is the collated metrics
	combinedMetrics := make(map[string]map[string]metric)

	for location, collectedMetrics := range collectedMetrics {
		combinedMetrics[location] = make(map[string]metric, 0)

		for i, metricKey := range metricKeys {
			// Find all metrics in this zone that match the key
			matchingKeyMetrics := make([]wrappedMetric, 0)
			for _, metric := range collectedMetrics {
				if metric.name == metricKey {
					matchingKeyMetrics = append(matchingKeyMetrics, metric)
				}
			}

			// Combine the metrics
			if metricTypes[i] == workerMetrics.Counter {
				combinedMetrics[location][metricKey] = determineCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			} else if metricTypes[i] == workerMetrics.Gauge {
				// Gauges are summed
				combinedMetrics[location][metricKey] = determineCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Gauge)
			} else if metricTypes[i] == workerMetrics.Rate {
				// Rates are summed
				combinedMetrics[location][metricKey] = determineCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Rate)
			} else if metricTypes[i] == workerMetrics.Trend {
				// Trends are summed
				combinedMetrics[location][metricKey] = determineTrend(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Trend, "mean")
			}
		}

	}

	// If any keys are empty, remove them
	for key, value := range combinedMetrics {
		if len(value) == 0 {
			delete(combinedMetrics, key)
		}
	}

	return combinedMetrics, nil
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

		aggregatedMetric.Value /= subFractionTotal
	}

	return aggregatedMetric
}

func (store *cachedMetricsStore) Stop() {
	if store.flusher != nil {
		store.flusher.Stop()
	}
}
