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

type summaryMetric struct {
	Contains string                   `json:"contains"`
	Type     workerMetrics.MetricType `json:"type"`
	Values   map[string]float64       `json:"values"`
}

type wrappedSummaryMetric struct {
	summaryMetric
	childJobId  string
	location    string
	subFraction float64
}

type zoneSummaryMetrics struct {
	// Fraction is the total weighting of a zone
	ZoneMetrics []map[string]wrappedSummaryMetric `json:"zoneMetrics"`
}

type summaryBank struct {
	gs libOrch.BaseGlobalState
	// Locations and list of summary metrics for each location
	unknownMetrics map[string]zoneSummaryMetrics
	mu             *sync.Mutex
}

func NewSummaryBank(gs libOrch.BaseGlobalState, options *libWorker.Options) *summaryBank {
	sb := &summaryBank{
		gs:             gs,
		unknownMetrics: make(map[string]zoneSummaryMetrics),
		mu:             &sync.Mutex{},
	}

	// Add load distribution locations
	for _, location := range options.LoadDistribution.Value {
		sb.unknownMetrics[location.Location] = zoneSummaryMetrics{
			ZoneMetrics: make([]map[string]wrappedSummaryMetric, 0),
		}
	}

	// Add global location
	sb.unknownMetrics[libOrch.GlobalName] = zoneSummaryMetrics{
		ZoneMetrics: make([]map[string]wrappedSummaryMetric, 0),
	}

	return sb
}

func (sb *summaryBank) AddMessage(message libOrch.WorkerMessage, workerLocation string, subFraction float64) error {
	var metrics map[string]summaryMetric

	err := json.Unmarshal([]byte(message.Message), &metrics)
	if err != nil {
		return err
	}

	wrappedMetrics := make(map[string]wrappedSummaryMetric)
	for metricName, metric := range metrics {
		wrappedMetrics[metricName] = wrappedSummaryMetric{
			summaryMetric: metric,
			childJobId:    message.ChildJobId,
			location:      workerLocation,
			subFraction:   subFraction,
		}
	}

	sb.mu.Lock()
	defer sb.mu.Unlock()

	// Add the metrics to correct location array
	sb.unknownMetrics[workerLocation] = zoneSummaryMetrics{
		ZoneMetrics: append(sb.unknownMetrics[workerLocation].ZoneMetrics, wrappedMetrics),
	}

	// And to the global array
	sb.unknownMetrics[libOrch.GlobalName] = zoneSummaryMetrics{
		ZoneMetrics: append(sb.unknownMetrics[libOrch.GlobalName].ZoneMetrics, wrappedMetrics),
	}

	return nil
}

func (sb *summaryBank) CalculateAndDispatchSummaryMetrics() error {
	if len(sb.unknownMetrics) == 0 {
		return errors.New("need at least one group of metrics to calculate summary")
	}

	fistKey := ""

	// Ensure at least one metric in each location
	for key, zoneSummaryMetrics := range sb.unknownMetrics {
		if fistKey == "" {
			fistKey = key
		}

		if len(zoneSummaryMetrics.ZoneMetrics) == 0 {
			return errors.New("need at least one metric in each location to calculate summary")
		}
	}

	metricKeys := make([]string, 0)
	metricTypes := make([]workerMetrics.MetricType, 0)
	metricContains := make([]string, 0)

	// Determine the metric keys from first metric
	for metricKey, zoneSummaryMetric := range sb.unknownMetrics[fistKey].ZoneMetrics[0] {
		metricKeys = append(metricKeys, metricKey)
		metricTypes = append(metricTypes, zoneSummaryMetric.Type)
		metricContains = append(metricContains, zoneSummaryMetric.Contains)
	}

	outputSummaryMetrics := make(map[string]map[string]summaryMetric)

	// Add all locations to output metrics
	for location, zoneSummaryMetrics := range sb.unknownMetrics {
		outputSummaryMetrics[location] = make(map[string]summaryMetric)

		for i, metricKey := range metricKeys {
			// Find all metrics is this zone that match the metric key
			zoneMetrics := make([]map[string]wrappedSummaryMetric, 0)

			for _, zoneMetric := range zoneSummaryMetrics.ZoneMetrics {
				if zoneMetric[metricKey].Contains != "" {
					zoneMetrics = append(zoneMetrics, zoneMetric)
				}
			}

			if metricTypes[i] == workerMetrics.Counter {
				// Calculate the total for this metric
				outputSummaryMetrics[location][metricKey] = determineCounter(zoneMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			} else if metricTypes[i] == workerMetrics.Gauge {
				// Gauges are summed
				outputSummaryMetrics[location][metricKey] = determineCounter(zoneMetrics, metricKey, metricContains[i], workerMetrics.Gauge)
			} else if metricTypes[i] == workerMetrics.Rate {
				// Rates are summed
				outputSummaryMetrics[location][metricKey] = determineCounter(zoneMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			} else if metricTypes[i] == workerMetrics.Trend {
				// Calculate the trend for this metric
				outputSummaryMetrics[location][metricKey] = determineTrend(zoneMetrics, metricKey, metricContains[i], workerMetrics.Trend)
			} else {
				return fmt.Errorf("unknown metric type %s", metricTypes[i])
			}
		}
	}

	// Dispatch the summary metrics
	marshalledOutputSummary, err := json.Marshal(outputSummaryMetrics)
	if err != nil {
		return err
	}

	// Send the summary metrics to the orchestrator
	libOrch.DispatchMessage(sb.gs, string(marshalledOutputSummary), "SUMMARY_METRICS")

	return nil
}

// Calculates an aggregated counter metric for a zone
func determineCounter(zoneMetrics []map[string]wrappedSummaryMetric, metricName string,
	metricContains string, metricType workerMetrics.MetricType) summaryMetric {
	aggregatedMetric := summaryMetric{
		Contains: metricContains,
		Type:     metricType,
		Values:   make(map[string]float64),
	}

	// If no value keys, return an empty metric
	if len(zoneMetrics[0][metricName].Values) == 0 {
		return aggregatedMetric
	}

	// Determine the value keys from first metric
	valueKeys := make([]string, 0)
	for valueKey := range zoneMetrics[0][metricName].Values {
		valueKeys = append(valueKeys, valueKey)
	}

	for _, valueKey := range valueKeys {
		for _, zoneMetric := range zoneMetrics {
			aggregatedMetric.Values[valueKey] += zoneMetric[metricName].Values[valueKey]
		}
	}

	return aggregatedMetric
}

// Calculates a weighted mean value metric for a zone
func determineTrend(zoneMetrics []map[string]wrappedSummaryMetric, metricName string,
	metricContains string, metricType workerMetrics.MetricType) summaryMetric {
	aggregatedMetric := summaryMetric{
		Contains: metricContains,
		Type:     metricType,
		Values:   make(map[string]float64),
	}

	// If no value keys, return an empty metric
	if len(zoneMetrics[0][metricName].Values) == 0 {
		return aggregatedMetric
	}

	// Determine the value keys from first metric
	valueKeys := make([]string, 0)
	for valueKey := range zoneMetrics[0][metricName].Values {
		valueKeys = append(valueKeys, valueKey)
	}

	// Determine the weighted average of each value key from each metric
	subFractionTotal := 0.0

	for _, valueKey := range valueKeys {
		if valueKey == "max" {
			// Find biggest value
			for _, zoneMetric := range zoneMetrics {
				if zoneMetric[metricName].Values[valueKey] > aggregatedMetric.Values[valueKey] {
					aggregatedMetric.Values[valueKey] = zoneMetric[metricName].Values[valueKey]
				}
			}
		} else if valueKey == "min" {
			// Find smallest value
			aggregatedMetric.Values[valueKey] = math.MaxFloat64

			for _, zoneMetric := range zoneMetrics {
				if zoneMetric[metricName].Values[valueKey] < aggregatedMetric.Values[valueKey] {
					aggregatedMetric.Values[valueKey] = zoneMetric[metricName].Values[valueKey]
				}
			}
		} else {
			// This isn't ideal for all all remaining value keys but better than nothing

			for _, zoneMetric := range zoneMetrics {
				subFractionTotal += zoneMetric[metricName].subFraction
				aggregatedMetric.Values[valueKey] += zoneMetric[metricName].Values[valueKey] * zoneMetric[metricName].subFraction
			}

			aggregatedMetric.Values[valueKey] /= subFractionTotal

		}
	}

	return aggregatedMetric
}

func (sb *summaryBank) Size() int {
	count := 0

	for k, v := range sb.unknownMetrics {
		if k != libOrch.GlobalName {
			count += len(v.ZoneMetrics)
		}
	}

	return count
}
