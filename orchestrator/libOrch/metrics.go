package libOrch

// Cached metrics are stored before being collated and sent
type BaseMetricsStore interface ***REMOVED***
	AddMessage(message WorkerMessage, workerLocation string) error
	Stop()
	FlushMetrics()
***REMOVED***
