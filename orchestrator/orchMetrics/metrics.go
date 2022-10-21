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

type metric struct ***REMOVED***
	Contains string                   `json:"contains"`
	Type     workerMetrics.MetricType `json:"type"`
	Value    float64                  `json:"value"`
***REMOVED***

type wrappedMetric struct ***REMOVED***
	metric
	name        string
	location    string
	subFraction float64
	childJobId  string
***REMOVED***

// Cached metrics are stored before being collated and sent
type cachedMetricsStore struct ***REMOVED***
	// Each envelope in the map is a certain metric
	collectedMetrics map[string][]wrappedMetric
	mu               sync.RWMutex
	flusher          *output.PeriodicFlusher
	gs               libOrch.BaseGlobalState
	options          *libWorker.Options
***REMOVED***

var (
	_ libOrch.BaseMetricsStore = &cachedMetricsStore***REMOVED******REMOVED***
)

func NewCachedMetricsStore(gs libOrch.BaseGlobalState) *cachedMetricsStore ***REMOVED***
	return &cachedMetricsStore***REMOVED***
		gs:               gs,
		collectedMetrics: make(map[string][]wrappedMetric),
		mu:               sync.RWMutex***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

func (store *cachedMetricsStore) InitMetricsStore(options *libWorker.Options) ***REMOVED***
	store.mu.Lock()
	defer store.mu.Unlock()
	store.options = options

	// Add load distribution locations
	for _, location := range options.LoadDistribution.Value ***REMOVED***
		store.collectedMetrics[location.Location] = make([]wrappedMetric, 0)
	***REMOVED***

	// Add global location
	store.collectedMetrics[libOrch.GlobalName] = make([]wrappedMetric, 0)

	// This will never return an error
	pf, _ := output.NewPeriodicFlusher(1000*time.Millisecond, store.FlushMetrics)

	store.flusher = pf
***REMOVED***

func (store *cachedMetricsStore) AddMessage(message libOrch.WorkerMessage, workerLocation string, subFraction float64) error ***REMOVED***
	if store.options == nil ***REMOVED***
		return errors.New("metrics store not initialised")
	***REMOVED***

	// Ensure childJobId is not already in the store, prevents duplicate metrics
	for _, existingWrappedMetric := range store.collectedMetrics[workerLocation] ***REMOVED***
		if existingWrappedMetric.childJobId == message.ChildJobId ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	var sampleEnvelopes []globetest.SampleEnvelope

	err := json.Unmarshal([]byte(message.Message), &sampleEnvelopes)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	wrappedMetrics := make([]wrappedMetric, 0)
	for _, sampleEnvelope := range sampleEnvelopes ***REMOVED***
		wrappedMetrics = append(wrappedMetrics, wrappedMetric***REMOVED***
			metric: metric***REMOVED***
				Contains: sampleEnvelope.Metric.Contains.String(),
				Type:     sampleEnvelope.Metric.Type,
				Value:    sampleEnvelope.Data.Value,
			***REMOVED***,
			name:        sampleEnvelope.Metric.Name,
			location:    workerLocation,
			subFraction: subFraction,
			childJobId:  message.ChildJobId,
		***REMOVED***)
	***REMOVED***

	// Add the metrics to the store
	store.mu.Lock()
	defer store.mu.Unlock()

	for _, wrappedMetric := range wrappedMetrics ***REMOVED***
		// Add the metrics to the correct location array
		store.collectedMetrics[workerLocation] = append(store.collectedMetrics[workerLocation], wrappedMetric)

		// Add the metrics to the global array
		store.collectedMetrics[libOrch.GlobalName] = append(store.collectedMetrics[libOrch.GlobalName], wrappedMetric)
	***REMOVED***

	return nil
***REMOVED***

// Empty the store and returns its contents
func (store *cachedMetricsStore) emptyStore() (map[string][]wrappedMetric, error) ***REMOVED***
	if store.options == nil ***REMOVED***
		return nil, errors.New("metrics store not initialised")
	***REMOVED***

	store.mu.Lock()
	defer store.mu.Unlock()

	// Copy the map
	result := make(map[string][]wrappedMetric, len(store.collectedMetrics))
	for metricName, sampleEnvelopes := range store.collectedMetrics ***REMOVED***
		result[metricName] = make([]wrappedMetric, len(sampleEnvelopes))
		copy(result[metricName], sampleEnvelopes)
	***REMOVED***

	// Empty the map
	store.collectedMetrics = make(map[string][]wrappedMetric)

	// Add load distribution locations
	for _, location := range store.options.LoadDistribution.Value ***REMOVED***
		store.collectedMetrics[location.Location] = make([]wrappedMetric, 0)
	***REMOVED***

	// Add global location
	store.collectedMetrics[libOrch.GlobalName] = make([]wrappedMetric, 0)

	return result, nil
***REMOVED***

func (store *cachedMetricsStore) FlushMetrics() ***REMOVED***
	collectedMetrics, err := store.getMetrics()
	if err != nil ***REMOVED***
		// Sometimes there are not enough metrics, this will throw an error
		return
	***REMOVED***

	if len(collectedMetrics) == 0 ***REMOVED***
		return
	***REMOVED***

	// Marshall the envelopes
	marshalledCollectedMetrics, err := json.Marshal(collectedMetrics)
	if err != nil ***REMOVED***
		libOrch.HandleError(store.gs, err)
		return
	***REMOVED***

	libOrch.DispatchMessage(store.gs, string(marshalledCollectedMetrics), "METRICS")
***REMOVED***

func (store *cachedMetricsStore) getMetrics() (map[string]map[string]metric, error) ***REMOVED***
	collectedMetrics, err := store.emptyStore()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if len(collectedMetrics) == 0 ***REMOVED***
		return nil, errors.New("need at least one group of metrics")
	***REMOVED***

	firstKey := ""

	// Ensure at least one metric in each location
	for key, zoneMetrics := range collectedMetrics ***REMOVED***
		if len(zoneMetrics) == 0 ***REMOVED***
			return nil, errors.New("need at least one metric in each location")
		***REMOVED***

		if firstKey == "" ***REMOVED***
			firstKey = key
		***REMOVED***
	***REMOVED***

	metricKeys := make([]string, 0)
	metricTypes := make([]workerMetrics.MetricType, 0)
	metricContains := make([]string, 0)

	// Find first metric in first location
	for _, metric := range collectedMetrics[firstKey] ***REMOVED***
		alreadExists := false
		for i := 0; i < len(metricKeys); i++ ***REMOVED***
			if metricKeys[i] == metric.name ***REMOVED***
				alreadExists = true
				break
			***REMOVED***
		***REMOVED***

		if !alreadExists ***REMOVED***
			metricKeys = append(metricKeys, metric.name)
			metricTypes = append(metricTypes, metric.metric.Type)
			metricContains = append(metricContains, metric.metric.Contains)
		***REMOVED***
	***REMOVED***

	// Combined metrics is the collated metrics
	combinedMetrics := make(map[string]map[string]metric)

	for location, collectedMetrics := range collectedMetrics ***REMOVED***
		combinedMetrics[location] = make(map[string]metric, 0)

		for i, metricKey := range metricKeys ***REMOVED***
			// Find all metrics in this zone that match the key
			matchingKeyMetrics := make([]wrappedMetric, 0)
			for _, metric := range collectedMetrics ***REMOVED***
				if metric.name == metricKey ***REMOVED***
					matchingKeyMetrics = append(matchingKeyMetrics, metric)
				***REMOVED***
			***REMOVED***

			// Combine the metrics
			if metricTypes[i] == workerMetrics.Counter ***REMOVED***
				combinedMetrics[location][metricKey] = determineCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Gauge ***REMOVED***
				// Gauges are summed
				combinedMetrics[location][metricKey] = determineCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Gauge)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Rate ***REMOVED***
				// Rates are summed
				combinedMetrics[location][metricKey] = determineCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Rate)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Trend ***REMOVED***
				// Trends are summed
				combinedMetrics[location][metricKey] = determineTrend(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Trend, "mean")
			***REMOVED***
		***REMOVED***

	***REMOVED***

	// If any keys are empty, remove them
	for key, value := range combinedMetrics ***REMOVED***
		if len(value) == 0 ***REMOVED***
			delete(combinedMetrics, key)
		***REMOVED***
	***REMOVED***

	return combinedMetrics, nil
***REMOVED***

// Calculates an aggregated counter metric for a zone
func determineCounter(matchingKeyMetrics []wrappedMetric, metricName string, metricContains string,
	metricType workerMetrics.MetricType) metric ***REMOVED***
	aggregatedMetric := metric***REMOVED***
		Contains: metricContains,
		Type:     metricType,
		Value:    0.0,
	***REMOVED***

	for _, zoneMetric := range matchingKeyMetrics ***REMOVED***
		aggregatedMetric.Value += zoneMetric.Value
	***REMOVED***

	return aggregatedMetric
***REMOVED***

// Calculates a weighted mean value metric for a zone
func determineTrend(matchingKeyMetrics []wrappedMetric, metricName string, metricContains string,
	metricType workerMetrics.MetricType, valueKey string) metric ***REMOVED***
	aggregatedMetric := metric***REMOVED***
		Contains: metricContains,
		Type:     metricType,
		Value:    0.0,
	***REMOVED***

	// Determine the weighted average of each value key from each metric
	subFractionTotal := 0.0

	if valueKey == "max" ***REMOVED***
		// Find biggest value
		for _, zoneMetric := range matchingKeyMetrics ***REMOVED***
			if zoneMetric.Value > aggregatedMetric.Value ***REMOVED***
				aggregatedMetric.Value = zoneMetric.Value
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if valueKey == "min" ***REMOVED***
		// Find smallest value
		aggregatedMetric.Value = math.MaxFloat64

		for _, zoneMetric := range matchingKeyMetrics ***REMOVED***
			if zoneMetric.Value < aggregatedMetric.Value ***REMOVED***
				aggregatedMetric.Value = zoneMetric.Value
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// This isn't ideal for all all remaining value keys but better than nothing

		for _, zoneMetric := range matchingKeyMetrics ***REMOVED***
			subFractionTotal += zoneMetric.subFraction
			aggregatedMetric.Value += zoneMetric.Value * zoneMetric.subFraction
		***REMOVED***

		aggregatedMetric.Value /= subFractionTotal
	***REMOVED***

	return aggregatedMetric
***REMOVED***

func (store *cachedMetricsStore) Stop() ***REMOVED***
	if store.flusher != nil ***REMOVED***
		store.flusher.Stop()
	***REMOVED***
***REMOVED***
