package aggregator

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"google.golang.org/protobuf/proto"
)

const intervalMaxLagPeriods = 20

func (aggregator *aggregator) AddInterval(message libOrch.WorkerMessage, location string, subFraction float64) error {
	if message.Message == "" {
		return nil
	}

	if aggregator.childJobs == nil {
		return errors.New("metrics aggregator not initialised")
	}

	messageBytes, err := base64.StdEncoding.DecodeString(message.Message)
	if err != nil {
		fmt.Printf("Error decoding interval message: %v", err)
		return err
	}

	streamedData := StreamedData{}
	err = proto.Unmarshal(messageBytes, &streamedData)
	if err != nil {
		fmt.Printf("Error unmarshalling while adding message: %v", err)
		return err
	}

	// Switch on the type of data point
	interval, ok := streamedData.DataPoints[0].Data.(*DataPoint_Interval)
	if !ok || interval == nil {
		return errors.New("failed to add interval to aggregator, data point not an interval")
	}

	aggregator.intervalsMutex.Lock()
	defer aggregator.intervalsMutex.Unlock()

	err = aggregator.addIntervalToStore(interval.Interval, location, subFraction)
	if err != nil {
		return err
	}

	aggregator.cleanupIrretrievableIntervals()

	return nil
}

// Each flush count is a new set of metrics, and the array is extended to accommodate
func (aggregator *aggregator) addIntervalToStore(interval *Interval, location string, subFraction float64) error {
	// Extend the intervals array if necessary
	storeLength := len(aggregator.intervals)
	if storeLength == 0 || int32(aggregator.intervals[storeLength-1].period) != interval.Period {
		aggregator.intervals = append(aggregator.intervals, &periodIntervals{
			period:    interval.Period,
			intervals: make([]*intervalWithSubfraction, 0),
		})
	}

	// Add the interval to the aggregator and determine if data is ready to be sent
	for _, periodIntervals := range aggregator.intervals {
		if periodIntervals.period == interval.Period {
			// Add the interval to the correct location
			periodIntervals.intervals = append(periodIntervals.intervals,
				&intervalWithSubfraction{
					interval:    interval,
					subFraction: subFraction,
					location:    location,
				})

			return aggregator.sendIntervalsIfAggregated(periodIntervals)
		}
	}

	return errors.New("failed to add interval to aggregator, period not found")
}

func (aggregator *aggregator) sendIntervalsIfAggregated(periodIntervals *periodIntervals) error {
	// Check if subfractions add up to 1 (or within epsilon)
	subFraction := 0.0
	for _, iwsf := range periodIntervals.intervals {
		subFraction += iwsf.subFraction
	}

	if subFraction-1 > 0.0000001 {
		return nil
	}

	aggregatedIntervals, err := aggregator.aggregateIntervals(periodIntervals.intervals, aggregator.gs.Standalone())
	if err != nil {
		fmt.Printf("Error aggregating intervals: %v", err)
		return err
	}

	streamedData := &StreamedData{
		DataPoints: []*DataPoint{
			{
				Data: &DataPoint_Interval{
					Interval: aggregatedIntervals,
				},
			},
		},
	}
	encodedBytes, err := proto.Marshal(streamedData)
	if err != nil {
		return err
	}

	libOrch.DispatchMessage(aggregator.gs, base64.StdEncoding.EncodeToString(encodedBytes), "INTERVAL")

	// Find the index of the periodIntervals and remove it
	for index, interval := range aggregator.intervals {
		if interval.period == periodIntervals.period {
			aggregator.intervals = append(aggregator.intervals[:index], aggregator.intervals[index+1:]...)
			break
		}
	}

	return nil
}

// If period lags considerably behind the leading period, then we can assume that the period is no longer retrievable
// and can be cleaned up, this is to prevent the memory from growing too large in the case that there is a lagging worker
func (aggregator *aggregator) cleanupIrretrievableIntervals() {
	if time.Since(aggregator.lastIntervalCleanupAt) < time.Second {
		return
	}

	aggregator.lastIntervalCleanupAt = time.Now()

	// Find the leading period
	leadingFlushCount := int32(0)
	for _, interval := range aggregator.intervals {
		if interval.period > leadingFlushCount {
			leadingFlushCount = interval.period
		}
	}

	indexesToRemove := make([]int, 0)

	// Remove all metrics that are more than intervalMaxLagPeriods behind the leading period
	for index, interval := range aggregator.intervals {
		if interval.period < leadingFlushCount-intervalMaxLagPeriods {
			indexesToRemove = append(indexesToRemove, index)
		}
	}

	for _, index := range indexesToRemove {
		if index < len(aggregator.intervals) {
			aggregator.intervals = append(aggregator.intervals[:index], aggregator.intervals[index+1:]...)
		} else {
			aggregator.intervals = aggregator.intervals[:index]
		}
	}
}
