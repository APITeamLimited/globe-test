package aggregator

import (
	"errors"
	"time"

	"github.com/APITeamLimited/globe-test/worker/metrics"
)

// Run processes all the thresholds with the provided Sink at the provided time and returns if any
// of them fails
func run_theshold_evaluation(ts *metrics.Thresholds, sink *Sink, duration time.Duration) (bool, error) {
	return false, errors.New("Not implemented")

	/*

		// Initialize the sinks store
		parsedSink := make(map[string]float64)

		for label, value := range sink.Labels {
			parsedSink[label] = value
		}

		switch sinkImpl := sink.(type) {
		case *metrics.CounterSink:
			parsedSink["count"] = sinkImpl.Value
			parsedSink["rate"] = sinkImpl.Value
		case *metrics.GaugeSink:
			parsedSink["value"] = sinkImpl.Value
		case *metrics.TrendSink:
			parsedSink["min"] = sinkImpl.Min
			parsedSink["max"] = sinkImpl.Max
			parsedSink["avg"] = sinkImpl.Avg
			parsedSink["med"] = sinkImpl.Med

			// Parse the percentile thresholds and insert them in
			// the sinks mapping.
			for _, threshold := range ts.Thresholds {
				if threshold.Parsed.AggregationMethod != metrics.TokenPercentile {
					continue
				}

				key := fmt.Sprintf("p(%g)", threshold.Parsed.AggregationValue.Float64)
				parsedSink[key] = sinkImpl.P(threshold.Parsed.AggregationValue.Float64 / 100)
			}
		case *metrics.RateSink:
			// We want to avoid division by zero, which
			// would lead to [#2520](https://github.com/grafana/k6/issues/2520)
			if sinkImpl.Total > 0 {
				parsedSink["rate"] = float64(sinkImpl.Trues) / float64(sinkImpl.Total)
			}
		case metrics.DummySink:
			for k, v := range sinkImpl {
				parsedSink[k] = v
			}
		default:
			return false, fmt.Errorf("unable to run Thresholds; reason: unknown sink type")
		}

		ts.Sinked = parsedSink

		return ts.RunAll(duration)

	*/
}
