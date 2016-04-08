package worker

import (
	"github.com/loadimpact/speedboat/comm"
)

// All registered worker processors.
var GlobalProcessors []func(*Worker) comm.Processor

// Register a worker processor.
func RegisterProcessor(factory func(*Worker) comm.Processor) ***REMOVED***
	GlobalProcessors = append(GlobalProcessors, factory)
***REMOVED***
