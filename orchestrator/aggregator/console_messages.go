package aggregator

import (
	"encoding/base64"
	"fmt"
	"hash/fnv"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"google.golang.org/protobuf/proto"
)

const MaxUniqueConsoleLogs = 100

func (aggregator *aggregator) AddConsoleMessages(message libOrch.WorkerMessage, workerLocation string) error {
	aggregator.consoleMutex.Lock()
	defer aggregator.consoleMutex.Unlock()

	messageBytes, err := base64.StdEncoding.DecodeString(message.Message)
	if err != nil {
		fmt.Printf("Error decoding console message: %v", err)
		return err
	}

	streamedData := StreamedData{}
	err = proto.Unmarshal(messageBytes, &streamedData)
	if err != nil {
		return err
	}

	for _, dataPoint := range streamedData.DataPoints {
		// Check if console message
		if consoleMessage, ok := dataPoint.Data.(*DataPoint_ConsoleMessage); ok {
			err = aggregator.addConsoleMessage(consoleMessage.ConsoleMessage, workerLocation)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (aggregator *aggregator) addConsoleMessage(consoleMessage *ConsoleMessage, workerLocation string) error {
	includeGlobalLocation := aggregator.gs.Standalone()

	// Check if message and level
	alreadySeen := false
	for _, seenMessage := range aggregator.consoleMessages {
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

			if includeGlobalLocation {
				seenMessage.Count[libOrch.GlobalName] += totalCountDelta
			}

			break
		}
	}

	if !alreadySeen {
		if aggregator.consoleMessageHashCount >= MaxUniqueConsoleLogs {
			aggregator.sendMaxConsoleLogsReached()
			return nil
		}

		aggregator.consoleMessageHashes = append(aggregator.consoleMessageHashes, consoleMessageHash(consoleMessage))
		aggregator.consoleMessageHashCount++

		targetMessage := consoleMessage

		if includeGlobalLocation {
			totalCount := int32(0)
			for _, count := range targetMessage.Count {
				totalCount += count
			}

			targetMessage.Count[libOrch.GlobalName] = totalCount
		}

		aggregator.consoleMessages = append(aggregator.consoleMessages, targetMessage)
	}

	return nil
}

func (aggregator *aggregator) flushConsoleMessages() error {
	aggregator.consoleMutex.Lock()
	defer aggregator.consoleMutex.Unlock()

	consoleMessageCount := len(aggregator.consoleMessages)

	if consoleMessageCount == 0 {
		return nil
	}

	streamedData := &StreamedData{
		DataPoints: make([]*DataPoint, consoleMessageCount),
	}

	for i, message := range aggregator.consoleMessages {
		streamedData.DataPoints[i] = &DataPoint{
			Data: &DataPoint_ConsoleMessage{
				ConsoleMessage: message,
			},
		}
	}

	encodedBytes, err := proto.Marshal(streamedData)
	if err != nil {
		return err
	}

	libOrch.DispatchMessage(aggregator.gs, base64.StdEncoding.EncodeToString(encodedBytes), "CONSOLE")

	for _, cm := range aggregator.consoleMessages {
		cm.Count = make(map[string]int32)
	}

	return nil
}

func consoleMessageHash(message *ConsoleMessage) string {
	// Create hash for message.msg + message.level
	h := fnv.New32a()
	h.Write([]byte(message.Message))

	return fmt.Sprintf("%d%s", h.Sum32(), message.Level)
}

func (aggregator *aggregator) sendMaxConsoleLogsReached() {
	if aggregator.sentMaxLogsMessage {
		return
	}

	libOrch.DispatchMessage(aggregator.gs, "MAX_CONSOLE_LOGS_REACHED", "MESSAGE")

	aggregator.sentMaxLogsMessage = true
}
