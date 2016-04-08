package worker

import (
	"github.com/loadimpact/speedboat/comm"
	"testing"
)

func TestRegisterProcessor(t *testing.T) ***REMOVED***
	oldGlobalProcessors := GlobalProcessors
	GlobalProcessors = nil
	defer func() ***REMOVED*** GlobalProcessors = oldGlobalProcessors ***REMOVED***()

	RegisterProcessor(func(w *Worker) comm.Processor ***REMOVED*** return nil ***REMOVED***)
	if len(GlobalProcessors) != 1 ***REMOVED***
		t.Error("Processor not registered")
	***REMOVED***
***REMOVED***
