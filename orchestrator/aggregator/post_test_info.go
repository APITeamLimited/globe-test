package aggregator

import (
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"google.golang.org/protobuf/proto"
)

func DeterminePostTestInfo(gs libOrch.BaseGlobalState, messages *[]libOrch.OrchestratorOrWorkerMessage) (*TestInfo, error) {
	testInfo := TestInfo{
		Intervals:       make([]*Interval, 0),
		ConsoleMessages: make([]*ConsoleMessage, 0),
		Thresholds:      make([]*Threshold, 0),
	}

	for _, message := range *messages {
		if message.MessageType == "INTERVAL" {
			interval := Interval{}

			messageBytes, err := base64.StdEncoding.DecodeString(message.Message)
			if err != nil {
				fmt.Printf("Error decoding interval message during post test cleanup: %v", err)
				return nil, err
			}

			err = proto.Unmarshal(messageBytes, &interval)
			if err != nil {
				return nil, err
			}

			testInfo.Intervals = append(testInfo.Intervals, &interval)
		} else if message.MessageType == "CONSOLE" {
			consoleMessage := ConsoleMessage{}

			messageBytes, err := base64.StdEncoding.DecodeString(message.Message)
			if err != nil {
				fmt.Printf("Error decoding message: %v", err)
				return nil, err
			}

			err = proto.Unmarshal(messageBytes, &consoleMessage)
			if err != nil {
				return nil, err
			}

			testInfo.ConsoleMessages = AggregateConsoleMessages(append(testInfo.ConsoleMessages, &consoleMessage))
		}
	}

	// Todo get thresholds
	//testInfo.Thresholds = gs.MetricsStore().GetThresholds()

	// Sort intervals by Period

	sortIntervalsByPeriod(testInfo.Intervals)

	return &testInfo, nil
}

func sortIntervalsByPeriod(intervals []*Interval) {
	sort.SliceStable(intervals, func(i, j int) bool {
		return intervals[i].Period < intervals[j].Period
	})
}
