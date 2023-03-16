package aggregator

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/worker/metrics"
	"google.golang.org/protobuf/proto"
)

func (aggregator *aggregator) evaluateThresholds() {
	aggregator.lockAllMutexes()
	defer aggregator.unlockAllMutexes()

	if len(aggregator.previousIntervals) == 0 {
		return
	}

	currentDuration := time.Since(*aggregator.thresholdStartTime)
	standalone := aggregator.gs.Standalone()
	mostRecentInterval := aggregator.previousIntervals[len(aggregator.previousIntervals)-1]

	// Need to populate sinks from aggregator.intervals

	abort := false

	// Need to copy exact values of thresholds for hash
	newThresholds := make(map[string]*metrics.Thresholds, len(aggregator.thresholds))

	for metricName, thresholds := range aggregator.thresholds {
		newThresholds[metricName] = &thresholds

		sinkName := getSinkName(metricName, standalone)

		sink, ok := mostRecentInterval.Sinks[sinkName]
		if !ok {
			// Haven't found the sink, so we can't evaluate thresholds for it.
			continue
		}

		// If either the metric has no thresholds defined, or its sinks
		// are empty, let's ignore its thresholds execution at this point.
		if len(thresholds.Thresholds) == 0 || len(sink.Labels) == 0 {
			continue
		}

		_, err := runThresholdEvaluation(newThresholds[metricName], sink, currentDuration)
		if err != nil {
			libOrch.HandleError(aggregator.gs, fmt.Errorf("error evaluating thresholds for metric %s: %v", metricName, err))
			return
		}

		if thresholds.Abort {
			abort = true
		}

	}

	hashedThresholds, err := hashThresholds(newThresholds)
	if err != nil {
		libOrch.HandleError(aggregator.gs, fmt.Errorf("error hashing thresholds: %v", err))
		return
	}

	if hashedThresholds != aggregator.previousThresholdsHash {
		aggregator.previousThresholdsHash = hashedThresholds
		err = aggregator.sendThresholds(newThresholds)
		if err != nil {
			libOrch.HandleError(aggregator.gs, fmt.Errorf("error sending thresholds: %v", err))
			return
		}
	}

	if abort {
		libOrch.HandleStringError(aggregator.gs, "thresholds failed and abort was set to true, aborting test")
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

func (aggregator *aggregator) sendThresholds(thresholds map[string]*metrics.Thresholds) error {
	// Mutexes are already locked

	thresholdCount := 0
	for _, threshold := range thresholds {
		thresholdCount += len(threshold.Thresholds)
	}

	streamedData := &StreamedData{
		DataPoints: make([]*DataPoint, thresholdCount),
	}

	i := 0
	for metricName, thresholdGroup := range thresholds {
		for _, threshold := range thresholdGroup.Thresholds {
			streamedData.DataPoints[i] = &DataPoint{
				Data: &DataPoint_Threshold{
					Threshold: &Threshold{
						Source:         threshold.Source,
						Metric:         metricName,
						AbortOnFail:    &threshold.AbortOnFail,
						DelayAbortEval: &threshold.Parsed.AbortGracePeriodSource,
					},
				},
			}
			i++
		}
	}

	encodedBytes, err := proto.Marshal(streamedData)
	if err != nil {
		return err
	}

	libOrch.DispatchMessage(aggregator.gs, base64.StdEncoding.EncodeToString(encodedBytes), "THRESHOLD")

	return nil
}

func (aggregator *aggregator) sendInitialThresholds(thresholds map[string]metrics.Thresholds) error {
	mappedThresholds := make(map[string]*metrics.Thresholds, len(thresholds))
	for metricName, threshold := range thresholds {
		mappedThresholds[metricName] = &threshold
	}

	return aggregator.sendThresholds(mappedThresholds)
}
