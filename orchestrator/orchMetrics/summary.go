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
	childJobId  string
	location    string
	subFraction float64
***REMOVED***

type zoneSummaryMetrics struct ***REMOVED***
	// Fraction is the total weighting of a zone
	ZoneMetrics []map[string]wrappedSummaryMetric `json:"zoneMetrics"`
***REMOVED***

type summaryBank struct ***REMOVED***
	gs libOrch.BaseGlobalState
	// Locations and list of summary metrics for each location
	unknownMetrics map[string]zoneSummaryMetrics
	mu             *sync.Mutex
***REMOVED***

func NewSummaryBank(gs libOrch.BaseGlobalState, options *libWorker.Options) *summaryBank ***REMOVED***
	sb := &summaryBank***REMOVED***
		gs:             gs,
		unknownMetrics: make(map[string]zoneSummaryMetrics),
		mu:             &sync.Mutex***REMOVED******REMOVED***,
	***REMOVED***

	// Add load distribution locations
	for _, location := range options.LoadDistribution.Value ***REMOVED***
		sb.unknownMetrics[location.Location] = zoneSummaryMetrics***REMOVED***
			ZoneMetrics: make([]map[string]wrappedSummaryMetric, 0),
		***REMOVED***
	***REMOVED***

	// Add global location
	sb.unknownMetrics[libOrch.GlobalName] = zoneSummaryMetrics***REMOVED***
		ZoneMetrics: make([]map[string]wrappedSummaryMetric, 0),
	***REMOVED***

	return sb
***REMOVED***

func (sb *summaryBank) AddMessage(message libOrch.WorkerMessage, workerLocation string, subFraction float64) error ***REMOVED***
	var metrics map[string]summaryMetric

	err := json.Unmarshal([]byte(message.Message), &metrics)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	wrappedMetrics := make(map[string]wrappedSummaryMetric)
	for metricName, metric := range metrics ***REMOVED***
		wrappedMetrics[metricName] = wrappedSummaryMetric***REMOVED***
			summaryMetric: metric,
			childJobId:    message.ChildJobId,
			location:      workerLocation,
			subFraction:   subFraction,
		***REMOVED***
	***REMOVED***

	sb.mu.Lock()
	defer sb.mu.Unlock()

	// Add the metrics to correct location array
	sb.unknownMetrics[workerLocation] = zoneSummaryMetrics***REMOVED***
		ZoneMetrics: append(sb.unknownMetrics[workerLocation].ZoneMetrics, wrappedMetrics),
	***REMOVED***

	// And to the global array
	sb.unknownMetrics[libOrch.GlobalName] = zoneSummaryMetrics***REMOVED***
		ZoneMetrics: append(sb.unknownMetrics[libOrch.GlobalName].ZoneMetrics, wrappedMetrics),
	***REMOVED***

	return nil
***REMOVED***

func (sb *summaryBank) CalculateAndDispatchSummaryMetrics() error ***REMOVED***
	if len(sb.unknownMetrics) == 0 ***REMOVED***
		return errors.New("need at least one group of metrics to calculate summary")
	***REMOVED***

	fistKey := ""

	// Ensure at least one metric in each location
	for key, zoneSummaryMetrics := range sb.unknownMetrics ***REMOVED***
		if fistKey == "" ***REMOVED***
			fistKey = key
		***REMOVED***

		if len(zoneSummaryMetrics.ZoneMetrics) == 0 ***REMOVED***
			return errors.New("need at least one metric in each location to calculate summary")
		***REMOVED***
	***REMOVED***

	metricKeys := make([]string, 0)
	metricTypes := make([]workerMetrics.MetricType, 0)
	metricContains := make([]string, 0)

	// Determine the metric keys from first metric
	for metricKey, zoneSummaryMetric := range sb.unknownMetrics[fistKey].ZoneMetrics[0] ***REMOVED***
		metricKeys = append(metricKeys, metricKey)
		metricTypes = append(metricTypes, zoneSummaryMetric.Type)
		metricContains = append(metricContains, zoneSummaryMetric.Contains)
	***REMOVED***

	outputSummaryMetrics := make(map[string]map[string]summaryMetric)

	// Add all locations to output metrics
	for location, zoneSummaryMetrics := range sb.unknownMetrics ***REMOVED***
		outputSummaryMetrics[location] = make(map[string]summaryMetric)

		for i, metricKey := range metricKeys ***REMOVED***
			// Find all metrics is this zone that match the metric key
			zoneMetrics := make([]map[string]wrappedSummaryMetric, 0)

			for _, zoneMetric := range zoneSummaryMetrics.ZoneMetrics ***REMOVED***
				if zoneMetric[metricKey].Contains != "" ***REMOVED***
					zoneMetrics = append(zoneMetrics, zoneMetric)
				***REMOVED***
			***REMOVED***

			if metricTypes[i] == workerMetrics.Counter ***REMOVED***
				// Calculate the total for this metric
				outputSummaryMetrics[location][metricKey] = determineCounter(zoneMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Gauge ***REMOVED***
				// Gauges are summed
				outputSummaryMetrics[location][metricKey] = determineCounter(zoneMetrics, metricKey, metricContains[i], workerMetrics.Gauge)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Rate ***REMOVED***
				// Rates are summed
				outputSummaryMetrics[location][metricKey] = determineCounter(zoneMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			***REMOVED*** else if metricTypes[i] == workerMetrics.Trend ***REMOVED***
				// Calculate the trend for this metric
				outputSummaryMetrics[location][metricKey] = determineTrend(zoneMetrics, metricKey, metricContains[i], workerMetrics.Trend)
			***REMOVED*** else ***REMOVED***
				return fmt.Errorf("unknown metric type %s", metricTypes[i])
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Dispatch the summary metrics
	marshalledOutputSummary, err := json.Marshal(outputSummaryMetrics)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Send the summary metrics to the orchestrator
	libOrch.DispatchMessage(sb.gs, string(marshalledOutputSummary), "SUMMARY_METRICS")

	return nil
***REMOVED***

// Calculates an aggregated counter metric for a zone
func determineCounter(zoneMetrics []map[string]wrappedSummaryMetric, metricName string,
	metricContains string, metricType workerMetrics.MetricType) summaryMetric ***REMOVED***
	aggregatedMetric := summaryMetric***REMOVED***
		Contains: metricContains,
		Type:     metricType,
		Values:   make(map[string]float64),
	***REMOVED***

	// If no value keys, return an empty metric
	if len(zoneMetrics[0][metricName].Values) == 0 ***REMOVED***
		return aggregatedMetric
	***REMOVED***

	// Determine the value keys from first metric
	valueKeys := make([]string, 0)
	for valueKey := range zoneMetrics[0][metricName].Values ***REMOVED***
		valueKeys = append(valueKeys, valueKey)
	***REMOVED***

	for _, valueKey := range valueKeys ***REMOVED***
		for _, zoneMetric := range zoneMetrics ***REMOVED***
			aggregatedMetric.Values[valueKey] += zoneMetric[metricName].Values[valueKey]
		***REMOVED***
	***REMOVED***

	return aggregatedMetric
***REMOVED***

// Calculates a weighted mean value metric for a zone
func determineTrend(zoneMetrics []map[string]wrappedSummaryMetric, metricName string,
	metricContains string, metricType workerMetrics.MetricType) summaryMetric ***REMOVED***
	aggregatedMetric := summaryMetric***REMOVED***
		Contains: metricContains,
		Type:     metricType,
		Values:   make(map[string]float64),
	***REMOVED***

	// If no value keys, return an empty metric
	if len(zoneMetrics[0][metricName].Values) == 0 ***REMOVED***
		return aggregatedMetric
	***REMOVED***

	// Determine the value keys from first metric
	valueKeys := make([]string, 0)
	for valueKey := range zoneMetrics[0][metricName].Values ***REMOVED***
		valueKeys = append(valueKeys, valueKey)
	***REMOVED***

	// Determine the weighted average of each value key from each metric
	subFractionTotal := 0.0

	for _, valueKey := range valueKeys ***REMOVED***
		if valueKey == "max" ***REMOVED***
			// Find biggest value
			for _, zoneMetric := range zoneMetrics ***REMOVED***
				if zoneMetric[metricName].Values[valueKey] > aggregatedMetric.Values[valueKey] ***REMOVED***
					aggregatedMetric.Values[valueKey] = zoneMetric[metricName].Values[valueKey]
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if valueKey == "min" ***REMOVED***
			// Find smallest value
			aggregatedMetric.Values[valueKey] = math.MaxFloat64

			for _, zoneMetric := range zoneMetrics ***REMOVED***
				if zoneMetric[metricName].Values[valueKey] < aggregatedMetric.Values[valueKey] ***REMOVED***
					aggregatedMetric.Values[valueKey] = zoneMetric[metricName].Values[valueKey]
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// This isn't ideal for all all remaining value keys but better than nothing

			for _, zoneMetric := range zoneMetrics ***REMOVED***
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

	for k, v := range sb.unknownMetrics ***REMOVED***
		if k != libOrch.GlobalName ***REMOVED***
			count += len(v.ZoneMetrics)
		***REMOVED***
	***REMOVED***

	return count
***REMOVED***
