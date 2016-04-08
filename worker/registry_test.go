package worker

import (
	"github.com/loadimpact/speedboat/comm"
	"testing"
)

func TestRegisterProcessor(t *testing.T) ***REMOVED***
	GlobalProcessors = nil
	RegisterProcessor(func(w *Worker) comm.Processor ***REMOVED*** return nil ***REMOVED***)
	if len(GlobalProcessors) != 1 ***REMOVED***
		t.Error("Processor not registered")
	***REMOVED***
***REMOVED***
