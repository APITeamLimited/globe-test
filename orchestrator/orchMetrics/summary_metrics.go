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
	location    string
	subFraction float64
}

type summaryBank struct {
	gs libOrch.BaseGlobalState
	// Locations and list of summary summaryMetrics for each location
	collectedSummaryMetrics map[string][]map[string]wrappedSummaryMetric
	mu                      *sync.Mutex
}

func NewSummaryBank(gs libOrch.BaseGlobalState, options *libWorker.Options) *summaryBank {
	sb := &summaryBank{
		gs:                      gs,
		collectedSummaryMetrics: make(map[string][]map[string]wrappedSummaryMetric),
		mu:                      &sync.Mutex{},
	}

	// Add load distribution locations
	for _, location := range options.LoadDistribution.Value {
		sb.collectedSummaryMetrics[location.Location] = make([]map[string]wrappedSummaryMetric, 0)
	}

	// Add global location
	sb.collectedSummaryMetrics[libOrch.GlobalName] = make([]map[string]wrappedSummaryMetric, 0)

	return sb
}

func (sb *summaryBank) Cleanup() {
	sb = nil
}

func (sb *summaryBank) AddMessage(message libOrch.WorkerMessage, workerLocation string, subFraction float64) error {
	var summaryMetrics map[string]summaryMetric

	err := json.Unmarshal([]byte(message.Message), &summaryMetrics)
	if err != nil {
		fmt.Println("Error unmarshalling summaryMetrics", err)
		return err
	}

	wrappedSummaryMetrics := make(map[string]wrappedSummaryMetric)
	for metricName, metric := range summaryMetrics {
		wrappedSummaryMetrics[metricName] = wrappedSummaryMetric{
			summaryMetric: metric,
			location:      workerLocation,
			subFraction:   subFraction,
		}
	}

	sb.mu.Lock()
	defer sb.mu.Unlock()

	// Add the summaryMetrics to correct location array
	sb.collectedSummaryMetrics[workerLocation] = append(sb.collectedSummaryMetrics[workerLocation], wrappedSummaryMetrics)

	// And to the global array
	sb.collectedSummaryMetrics[libOrch.GlobalName] = append(sb.collectedSummaryMetrics[libOrch.GlobalName], wrappedSummaryMetrics)

	return nil
}

func (sb *summaryBank) CalculateAndDispatchSummaryMetrics() error {
	if len(sb.collectedSummaryMetrics) == 0 {
		return errors.New("need at least one group of summaryMetrics to calculate summary")
	}

	firstKey := ""

	// Ensure at least one metric in each location
	for key, zoneSummaryMetrics := range sb.collectedSummaryMetrics {
		if len(zoneSummaryMetrics) == 0 {
			return errors.New("need at least one metric in each location to calculate summary")
		}

		if firstKey == "" {
			firstKey = key
		}
	}

	metricKeys := make([]string, 0)
	metricTypes := make([]workerMetrics.MetricType, 0)
	metricContains := make([]string, 0)

	// Determine the metric keys from first metric
	for metricKey, zoneSummaryMetric := range sb.collectedSummaryMetrics[firstKey][0] {
		metricKeys = append(metricKeys, metricKey)
		metricTypes = append(metricTypes, zoneSummaryMetric.Type)
		metricContains = append(metricContains, zoneSummaryMetric.Contains)
	}

	outputSummaryMetrics := make(map[string]map[string]summaryMetric)

	// Add all locations to output summaryMetrics
	for location, zoneSummaryMetrics := range sb.collectedSummaryMetrics {
		outputSummaryMetrics[location] = make(map[string]summaryMetric)

		for i, metricKey := range metricKeys {
			// Find all summaryMetrics is this zone that match the metric key
			matchingKeyMetrics := make([]map[string]wrappedSummaryMetric, 0)
			for _, zoneMetric := range zoneSummaryMetrics {
				if zoneMetric[metricKey].Contains != "" {
					matchingKeyMetrics = append(matchingKeyMetrics, zoneMetric)
				}
			}

			// Combine the summaryMetrics
			if metricTypes[i] == workerMetrics.Counter {
				// Calculate the total for this metric
				outputSummaryMetrics[location][metricKey] = determineSummaryCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			} else if metricTypes[i] == workerMetrics.Gauge {
				// Gauges are summed
				outputSummaryMetrics[location][metricKey] = determineSummaryCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Gauge)
			} else if metricTypes[i] == workerMetrics.Rate {
				// Rates are summed
				outputSummaryMetrics[location][metricKey] = determineSummaryCounter(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Counter)
			} else if metricTypes[i] == workerMetrics.Trend {
				// Calculate the trend for this metric
				outputSummaryMetrics[location][metricKey] = determineSummaryTrend(matchingKeyMetrics, metricKey, metricContains[i], workerMetrics.Trend)
			} else {
				return fmt.Errorf("unknown metric type %s", metricTypes[i])
			}
		}
	}

	// Dispatch the summary summaryMetrics
	marshalledOutputSummary, err := json.Marshal(outputSummaryMetrics)
	if err != nil {
		return err
	}

	// Send the summary summaryMetrics to the orchestrator
	libOrch.DispatchMessage(sb.gs, string(marshalledOutputSummary), "SUMMARY_METRICS")

	return nil
}

// Calculates an aggregated counter summaryMetric for a zone
func determineSummaryCounter(matchingKeyMetrics []map[string]wrappedSummaryMetric, metricName string,
	metricContains string, metricType workerMetrics.MetricType) summaryMetric {
	aggregatedMetric := summaryMetric{
		Contains: metricContains,
		Type:     metricType,
		Values:   make(map[string]float64),
	}

	// If no value keys, return an empty metric
	if len(matchingKeyMetrics) == 0 || len(matchingKeyMetrics[0][metricName].Values) == 0 {
		return aggregatedMetric
	}

	// Determine the value keys from first metric
	valueKeys := make([]string, 0)
	for valueKey := range matchingKeyMetrics[0][metricName].Values {
		valueKeys = append(valueKeys, valueKey)
	}

	for _, valueKey := range valueKeys {
		for _, zoneMetric := range matchingKeyMetrics {
			aggregatedMetric.Values[valueKey] += zoneMetric[metricName].Values[valueKey]
		}
	}

	return aggregatedMetric
}

// Calculates a weighted mean value summaryMetric for a zone
func determineSummaryTrend(matchingKeyMetrics []map[string]wrappedSummaryMetric, metricName string,
	metricContains string, metricType workerMetrics.MetricType) summaryMetric {
	aggregatedMetric := summaryMetric{
		Contains: metricContains,
		Type:     metricType,
		Values:   make(map[string]float64),
	}

	// If no value keys, return an empty metric
	if len(matchingKeyMetrics) == 0 || len(matchingKeyMetrics[0][metricName].Values) == 0 {
		return aggregatedMetric
	}

	// Determine the value keys from first metric
	valueKeys := make([]string, 0)
	for valueKey := range matchingKeyMetrics[0][metricName].Values {
		valueKeys = append(valueKeys, valueKey)
	}

	// Determine the weighted average of each value key from each metric
	subFractionTotal := 0.0

	for _, valueKey := range valueKeys {
		if valueKey == "max" {
			// Find biggest value
			for _, zoneMetric := range matchingKeyMetrics {
				if zoneMetric[metricName].Values[valueKey] > aggregatedMetric.Values[valueKey] {
					aggregatedMetric.Values[valueKey] = zoneMetric[metricName].Values[valueKey]
				}
			}
		} else if valueKey == "min" {
			// Find smallest value
			aggregatedMetric.Values[valueKey] = math.MaxFloat64

			for _, zoneMetric := range matchingKeyMetrics {
				if zoneMetric[metricName].Values[valueKey] < aggregatedMetric.Values[valueKey] {
					aggregatedMetric.Values[valueKey] = zoneMetric[metricName].Values[valueKey]
				}
			}
		} else {
			// This isn't ideal for all all remaining value keys but better than nothing

			for _, zoneMetric := range matchingKeyMetrics {
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

	if sb == nil || sb.collectedSummaryMetrics == nil {
		return count
	}

	for k, v := range sb.collectedSummaryMetrics {
		if v == nil {
			continue
		}

		if k != libOrch.GlobalName {
			count += len(v)
		}
	}

	return count
}
