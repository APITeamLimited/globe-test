package orchMetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/APITeamLimited/globe-test/worker/workerMetrics"
)

type summaryMetric struct ***REMOVED***
	Contains string                   `json:"contains"`
	Type     workerMetrics.MetricType `json:"type"`
	Values   map[string]float64       `json:"values"`
***REMOVED***

type wrappedSummaryMetric struct ***REMOVED***
	summaryMetric
	location    string
	subFraction float64
***REMOVED***

type summaryBank struct ***REMOVED***
	gs libOrch.BaseGlobalState
	// Locations and list of summary summaryMetrics for each location
	collectedSummaryMetrics map[string][]map[string]wrappedSummaryMetric
	mu                      *sync.Mutex
***REMOVED***

func NewSummaryBank(gs libOrch.BaseGlobalState, options *libWorker.Options) *summaryBank ***REMOVED***
	sb := &summaryBank***REMOVED***
		gs:                      gs,
		collectedSummaryMetrics: make(map[string][]map[string]wrappedSummaryMetric),
		mu:                      &sync.Mutex***REMOVED******REMOVED***,
	***REMOVED***

	// Add load distribution locations
	for _, location := range options.LoadDistribution.Value ***REMOVED***
		sb.collectedSummaryMetrics[location.Location] = make([]map[string]wrappedSummaryMetric, 0)
	***REMOVED***

	// Add global location
	sb.collectedSummaryMetrics[libOrch.GlobalName] = make([]map[string]wrappedSummaryMetric, 0)

	return sb
***REMOVED***

func (sb *summaryBank) AddMessage(message libOrch.WorkerMessage, workerLocation string, subFraction float64) error ***REMOVED***
	var summaryMetrics map[string]summaryMetric

	err := json.Unmarshal([]byte(message.Message), &summaryMetrics)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	wrappedSummaryMetrics := make(map[string]wrappedSummaryMetric)
	for metricName, metric := range summaryMetrics ***REMOVED***
		wrappedSummaryMetrics[metricName] = wrappedSummaryMetric***REMOVED***
			summaryMetric: metric,
			location:      workerLocation,
			subFraction:   subFraction,
		***REMOVED***
	***REMOVED***

	sb.mu.Lock()
	defer sb.mu.Unlock()

	// Add the summaryMetrics to correct location array
	sb.collectedSummaryMetrics[workerLocation] = append(sb.collectedSummaryMetrics[workerLocation], wrappedSummaryMetrics)

	// And to the global array
	sb.collectedSummaryMetrics[libOrch.GlobalName] = append(sb.collectedSummaryMetrics[libOrch.GlobalName], wrappedSummaryMetrics)

	return nil
***REMOVED***

func (sb *summaryBank) CalculateAndDispatchSummaryMetrics() error ***REMOVED***
	if len(sb.collectedSummaryMetrics) == 0 ***REMOVED***
		return errors.New("need at least one group of summaryMetrics to calculate summary")
	***REMOVED***

	firstKey := ""

	// Ensure at least one metric in each location
	for key, zoneSummaryMetrics := range sb.collectedSummaryMetrics ***REMOVED***
		if len(zoneSummaryMetrics) == 0 ***REMOVED***
			return errors.New("need at least one metric in each location to calculate summary")
		***REMOVED***

		if firstKey == "" ***REMOVED***
			firstKey = key
		***REMOVED***
	***REMOVED***

	metricKeys := make([]string, 0)
	metricTypes := make([]workerMetrics.MetricType, 0)
	metricContains := make([]string, 0)

	// Determine the metric keys from first metric
	for metricKey, zoneSummaryMetric := range sb.collectedSummaryMetrics[firstKey][0] ***REMOVED***
		metricKeys = append(metricKeys, metricKey)
		metricTypes = append(metricTypes, zoneSummaryMetric.Type)
		metricContains = append(metricContains, zoneSummaryMetric.Contains)
	***REMOVED***

	outputSummaryMetrics := make(map[string]map[string]summaryMetric)

	// Add all locations to output summaryMetrics
	for location, zoneSummaryMetrics := range sb.collectedSummaryMetrics ***REMOVED***
		outputSummaryMetrics[location] = make(map[string]summaryMetric)

		for i, metricKey := range metricKeys ***REMOVED***
			// Find all summaryMetrics is this zone that match the metric key
			matchingKeyMetrics := make([]map[string]wrappedSummaryMetric, 0)
			for _, zoneMetric := range zoneSummaryMetrics ***REMOVED***
				if zoneMetric[metricKey].Contains != "" ***REMOVED***
					matchingKeyMetrics = append(matchingKeyMetrics, zoneMetric)
				***REMOVED***
			***REMOVED***

			// Combine the summaryMetrics
			if metricTypes[i] == workerMetrics.Counter ***REMOVED***
				// Calculate the total for this metric
				outputSummaryMetrics[location][metricKey] = determineSummaryCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Gauge ***REMOVED***
				// Gauges are summed
				outputSummaryMetrics[location][metricKey] = determineSummaryCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Gauge)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Rate ***REMOVED***
				// Rates are summed
				outputSummaryMetrics[location][metricKey] = determineSummaryCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Trend ***REMOVED***
				// Calculate the trend for this metric
				outputSummaryMetrics[location][metricKey] = determineSummaryTrend(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Trend)
			***REMOVED*** else ***REMOVED***
				return fmt.Errorf("unknown metric type %s", metricTypes[i])
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Dispatch the summary summaryMetrics
	marshalledOutputSummary, err := json.Marshal(outputSummaryMetrics)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Send the summary summaryMetrics to the orchestrator
	libOrch.DispatchMessage(sb.gs, string(marshalledOutputSummary), "SUMMARY_METRICS")

	return nil
***REMOVED***

// Calculates an aggregated counter summaryMetric for a zone
func determineSummaryCounter(matchingKeyMetrics []map[string]wrappedSummaryMetric, metricName string,
	metricContains string, metricType workerMetrics.MetricType) summaryMetric ***REMOVED***
	aggregatedMetric := summaryMetric***REMOVED***
		Contains: metricContains,
		Type:     metricType,
		Values:   make(map[string]float64),
	***REMOVED***

	// If no value keys, return an empty metric
	if len(matchingKeyMetrics[0][metricName].Values) == 0 ***REMOVED***
		return aggregatedMetric
	***REMOVED***

	// Determine the value keys from first metric
	valueKeys := make([]string, 0)
	for valueKey := range matchingKeyMetrics[0][metricName].Values ***REMOVED***
		valueKeys = append(valueKeys, valueKey)
	***REMOVED***

	for _, valueKey := range valueKeys ***REMOVED***
		for _, zoneMetric := range matchingKeyMetrics ***REMOVED***
			aggregatedMetric.Values[valueKey] += zoneMetric[metricName].Values[valueKey]
		***REMOVED***
	***REMOVED***

	return aggregatedMetric
***REMOVED***

// Calculates a weighted mean value summaryMetric for a zone
func determineSummaryTrend(matchingKeyMetrics []map[string]wrappedSummaryMetric, metricName string,
	metricContains string, metricType workerMetrics.MetricType) summaryMetric ***REMOVED***
	aggregatedMetric := summaryMetric***REMOVED***
		Contains: metricContains,
		Type:     metricType,
		Values:   make(map[string]float64),
	***REMOVED***

	// If no value keys, return an empty metric
	if len(matchingKeyMetrics[0][metricName].Values) == 0 ***REMOVED***
		return aggregatedMetric
	***REMOVED***

	// Determine the value keys from first metric
	valueKeys := make([]string, 0)
	for valueKey := range matchingKeyMetrics[0][metricName].Values ***REMOVED***
		valueKeys = append(valueKeys, valueKey)
	***REMOVED***

	// Determine the weighted average of each value key from each metric
	subFractionTotal := 0.0

	for _, valueKey := range valueKeys ***REMOVED***
		if valueKey == "max" ***REMOVED***
			// Find biggest value
			for _, zoneMetric := range matchingKeyMetrics ***REMOVED***
				if zoneMetric[metricName].Values[valueKey] > aggregatedMetric.Values[valueKey] ***REMOVED***
					aggregatedMetric.Values[valueKey] = zoneMetric[metricName].Values[valueKey]
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if valueKey == "min" ***REMOVED***
			// Find smallest value
			aggregatedMetric.Values[valueKey] = math.MaxFloat64

			for _, zoneMetric := range matchingKeyMetrics ***REMOVED***
				if zoneMetric[metricName].Values[valueKey] < aggregatedMetric.Values[valueKey] ***REMOVED***
					aggregatedMetric.Values[valueKey] = zoneMetric[metricName].Values[valueKey]
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// This isn't ideal for all all remaining value keys but better than nothing

			for _, zoneMetric := range matchingKeyMetrics ***REMOVED***
				subFractionTotal += zoneMetric[metricName].subFraction
				aggregatedMetric.Values[valueKey] += zoneMetric[metricName].Values[valueKey] * zoneMetric[metricName].subFraction
			***REMOVED***

			aggregatedMetric.Values[valueKey] /= subFractionTotal

		***REMOVED***
	***REMOVED***

	return aggregatedMetric
***REMOVED***

func (sb *summaryBank) Size() int ***REMOVED***
	count := 0

	for k, v := range sb.collectedSummaryMetrics ***REMOVED***
		if k != libOrch.GlobalName ***REMOVED***
			count += len(v)
		***REMOVED***
	***REMOVED***

	return count
***REMOVED***
