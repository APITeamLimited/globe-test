package libOrch

// Cached metrics are stored before being collated and sent
type BaseMetricsStore interface {
	AddMessage(message WorkerMessage, workerLocation string) error
	Stop()
	FlushMetrics()
}
