package libOrch

import "github.com/APITeamLimited/globe-test/worker/metrics"

// Cached metrics are stored before being collated and sent
type BaseAggregator interface {
	InitAggregator(childJobs map[string]ChildJobDistribution, thresholds map[string]metrics.Thresholds)
	StartConsoleLogging()
	StartThresholdEvaluation()
	AddInterval(message WorkerMessage, workerLocation string, subFraction float64) error
	AddConsoleMessages(message WorkerMessage, workerLocation string) error
	StopAndCleanup()
	GetThresholds() map[string]metrics.Thresholds
}
