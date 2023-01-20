package libOrch

// Cached metrics are stored before being collated and sent
type BaseMetricsStore interface {
	InitMetricsStore(childJobs map[string]ChildJobDistribution)
	AddMessage(message WorkerMessage, workerLocation string, subFraction float64) error
	Cleanup()
}
