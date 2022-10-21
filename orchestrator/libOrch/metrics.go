package libOrch

import "github.com/APITeamLimited/globe-test/worker/libWorker"

// Cached metrics are stored before being collated and sent
type BaseMetricsStore interface ***REMOVED***
	InitMetricsStore(options *libWorker.Options)
	AddMessage(message WorkerMessage, workerLocation string, subFraction float64) error
	Stop()
	FlushMetrics()
***REMOVED***
