package aggregator

import (
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"google.golang.org/protobuf/proto"
)

var TEST_INFO_KEYS = []string{"INTERVAL", "CONSOLE", "MESSAGE"}

func DeterminePostTestInfo(gs libOrch.BaseGlobalState, messages *[]libOrch.OrchestratorOrWorkerMessage) (*TestInfo, error) {
	testInfo := TestInfo{
		Intervals:       make([]*Interval, 0),
		ConsoleMessages: make([]*ConsoleMessage, 0),
		Thresholds:      getFinalThresolds(gs),
		Messages:        make([]string, 0),
	}

	for _, message := range *messages {
		if message.MessageType == "INTERVAL" || message.MessageType == "CONSOLE" {
			streamedData := StreamedData{}
			messageBytes, err := base64.StdEncoding.DecodeString(message.Message)
			if err != nil {
				return nil, fmt.Errorf("error decoding message bytes: %s", err.Error())
			}

			err = proto.Unmarshal(messageBytes, &streamedData)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling message bytes: %s", err.Error())
			}

			for _, dataPoint := range streamedData.DataPoints {
				if consoleMessage, ok := dataPoint.Data.(*DataPoint_ConsoleMessage); ok {
					testInfo.ConsoleMessages = AggregateConsoleMessages(append(testInfo.ConsoleMessages, consoleMessage.ConsoleMessage))
				} else if interval, ok := dataPoint.Data.(*DataPoint_Interval); ok {
					testInfo.Intervals = append(testInfo.Intervals, interval.Interval)
				}
			}
		} else if message.MessageType == "MESSAGE" {
			testInfo.Messages = append(testInfo.Messages, message.Message)
		}
	}

	sortIntervalsByPeriod(testInfo.Intervals)

	return &testInfo, nil
}

func sortIntervalsByPeriod(intervals []*Interval) {
	sort.SliceStable(intervals, func(i, j int) bool {
		return intervals[i].Period < intervals[j].Period
	})
}

func getFinalThresolds(gs libOrch.BaseGlobalState) []*Threshold {
	thresholds := gs.MetricsStore().GetThresholds()
	finalThresholds := make([]*Threshold, 0, len(thresholds))

	for metricName, metricThresholds := range thresholds {
		for _, threshold := range metricThresholds.Thresholds {
			finalThresholds = append(finalThresholds, &Threshold{
				Metric:         metricName,
				Source:         threshold.Source,
				AbortOnFail:    &threshold.AbortOnFail,
				DelayAbortEval: &threshold.Parsed.AbortGracePeriodSource,
			})
		}
	}

	return finalThresholds
}
