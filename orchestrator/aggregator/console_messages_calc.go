package aggregator

// Combines console messages of same message and level into single messages, to be used
// for aggregating console messages during cleanup
func AggregateConsoleMessages(consoleMessages []*ConsoleMessage) []*ConsoleMessage {
	aggregateMessages := make([]*ConsoleMessage, 0)

	for _, consoleMessage := range consoleMessages {
		// Check if already in aggregateMessages
		found := false

		for _, aggregateMessage := range aggregateMessages {
			if aggregateMessage.Message == consoleMessage.Message && aggregateMessage.Level == consoleMessage.Level {
				found = true

				// Add counts
				for location, count := range consoleMessage.Count {
					aggregateMessage.Count[location] += count
				}

				// Update first occurred
				if consoleMessage.FirstOccurred.Seconds < aggregateMessage.FirstOccurred.Seconds {
					aggregateMessage.FirstOccurred = consoleMessage.FirstOccurred
				}

				// Update last occurred
				if consoleMessage.LastOccurred.Seconds > aggregateMessage.LastOccurred.Seconds {
					aggregateMessage.LastOccurred = consoleMessage.LastOccurred
				}
			}
		}

		if !found {
			// If no console message count, add it
			if consoleMessage.Count == nil {
				consoleMessage.Count = make(map[string]int32)
			}

			aggregateMessages = append(aggregateMessages, consoleMessage)
		}
	}

	return aggregateMessages
}
