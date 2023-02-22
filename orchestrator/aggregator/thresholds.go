package aggregator

import (
	"fmt"
	"strings"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/metrics"
)

func (aggregator *aggregator) evaluateThresholds() {
	aggregator.lockAllMutexes()
	defer aggregator.unlockAllMutexes()

	if len(aggregator.previousIntervals) == 0 {
		return
	}

	currentDuration := time.Since(*aggregator.thresholdStartTime)

	mostRecentInterval := aggregator.previousIntervals[len(aggregator.previousIntervals)-1]

	//Need to populate sinks from aggregator.intervals

	for metricName, thresholds := range aggregator.thresholds {
		sinkName := getSinkName(metricName, aggregator.gs.Standalone())

		sink, ok := mostRecentInterval.Sinks[sinkName]
		if !ok {
			libOrch.HandleError(aggregator.gs, fmt.Errorf("could not find sink %s, make sure it is specified correctly", sinkName))
		}

		// If either the metric has no thresholds defined, or its sinks
		// are empty, let's ignore its thresholds execution at this point.
		if len(thresholds.Thresholds) == 0 || len(sink.Labels) == 0 {
			continue
		}

		_, err := run_theshold_evaluation(&thresholds, sink, currentDuration)

		if err != nil {
			libOrch.HandleError(aggregator.gs, fmt.Errorf("error evaluating thresholds for metric %s: %v", metricName, err))
			return
		}

		if thresholds.Abort {
			libOrch.HandleStringError(aggregator.gs, fmt.Sprintf("Thresholds failed for metric %s, aborting", metricName))
		}

	}

}

func getSinkName(thresholdDescriptor string, standalone bool) string {
	parts := strings.Split(thresholdDescriptor, "::")
	partsLen := len(parts)

	if partsLen == 1 {
		if standalone {
			return fmt.Sprintf("%s::%s::group::default", libOrch.GlobalName, parts[0])
		}

		return fmt.Sprintf("localhost::%s::group::default", parts[0])
	}

	if partsLen == 2 {
		if standalone {
			return fmt.Sprintf("%s::%s::group::default", parts[0], parts[1])
		}

		return fmt.Sprintf("%s::%s::group::default", parts[0], parts[1])
	}

	// Assume user has specified the sink name correctly as they are using a long path

	return thresholdDescriptor
}

func (aggregator *aggregator) GetThresholds() map[string]metrics.Thresholds {
	aggregator.thresholdsMutex.Lock()
	defer aggregator.thresholdsMutex.Unlock()

	return aggregator.thresholds
}

// func GetOutputThresholds(thresholds map[string]metrics.Thresholds) []Thresholds {
// 	outputThresholds := make(map[string]metrics.Thresholds)

// 	for metricName, threshold := range thresholds {
// 		if threshold.Abort {
// 			outputThresholds[metricName] = threshold
// 		}
// 	}

// 	return outputThresholds
// }
