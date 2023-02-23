package aggregator

import (
	"encoding/json"
	"hash/fnv"
	"time"

	"github.com/APITeamLimited/globe-test/worker/metrics"
)

// Run processes all the thresholds with the provided Sink at the provided time and returns if any
// of them fails
func runThresholdEvaluation(ts *metrics.Thresholds, sink *Sink, duration time.Duration) (succeeded bool, err error) {
	// Initialize the sinks store
	ts.Sinked = sink.Labels

	return ts.RunAll(duration)
}

func hashThresholds(thresholds map[string]*metrics.Thresholds) (uint32, error) {
	marshalled, err := json.Marshal(thresholds)
	if err != nil {
		return 0, err
	}

	h := fnv.New32a()
	h.Write(marshalled)

	return h.Sum32(), nil
}
