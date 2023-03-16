package globetest

import (
	"encoding/base64"
	"time"

	"github.com/APITeamLimited/globe-test/orchestrator/aggregator"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (o *Output) addConsoleMessage(consoleMessage *aggregator.ConsoleMessage) error {
	o.consoleMutex.Lock()
	defer o.consoleMutex.Unlock()

	o.addedSinceLastFlush = true

	// Check if message and level
	alreadySeen := false
	for _, seenMessage := range o.consoleMessages {
		if seenMessage.Message == consoleMessage.Message && seenMessage.Level == consoleMessage.Level {
			alreadySeen = true

			totalCountDelta := int32(0)
			for location, count := range consoleMessage.Count {
				totalCountDelta += count

				seenMessage.Count[location] += count
			}

			if consoleMessage.FirstOccurred.Seconds < seenMessage.FirstOccurred.Seconds {
				seenMessage.FirstOccurred = consoleMessage.FirstOccurred
			}

			if consoleMessage.LastOccurred.Seconds > seenMessage.LastOccurred.Seconds {
				seenMessage.LastOccurred = consoleMessage.LastOccurred
			}

			break
		}
	}

	if !alreadySeen {
		// The max console logs message is sent frrom the orchestrator, so we don't need to send it here
		// Just make sure limit is breached slightly to ensure message is sent
		if o.consoleMessageCount > aggregator.MaxUniqueConsoleLogs+1 {
			return nil
		}

		o.consoleMessages = append(o.consoleMessages, consoleMessage)
		o.consoleMessageCount++
	}

	return nil
}

func (o *Output) listenOnLoggerChannel() {
	for log := range o.gs.GetLoggerChannel() {
		level := log["level"].(string)

		var msg string

		if log["error"] != nil {
			msg = log["error"].(string)
		} else {
			msg = log["msg"].(string)
		}

		consoleMessage := &aggregator.ConsoleMessage{
			Message: msg,
			Level:   level,
			Count: map[string]int32{
				o.location: 1,
			},
			FirstOccurred: timestamppb.New(time.Now()),
			LastOccurred:  timestamppb.New(time.Now()),
		}

		o.addConsoleMessage(consoleMessage)
	}
}

func (o *Output) flushConsoleMessages() error {
	o.consoleMutex.Lock()
	defer o.consoleMutex.Unlock()

	if !o.addedSinceLastFlush {
		return nil
	}

	streamedData := &aggregator.StreamedData{
		DataPoints: make([]*aggregator.DataPoint, len(o.consoleMessages)),
	}

	for i, message := range o.consoleMessages {
		streamedData.DataPoints[i] = &aggregator.DataPoint{
			Data: &aggregator.DataPoint_ConsoleMessage{
				ConsoleMessage: message,
			},
		}
	}

	encodedBytes, err := proto.Marshal(streamedData)
	if err != nil {
		return err
	}

	libWorker.DispatchMessage(o.gs, base64.StdEncoding.EncodeToString(encodedBytes), "CONSOLE")

	for _, cm := range o.consoleMessages {
		cm.Count = make(map[string]int32)
	}

	return nil
}
