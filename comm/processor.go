package comm

import (
	"sync"
)

type Processor interface ***REMOVED***
	Process(msg Message) <-chan Message
***REMOVED***

func Process(processors []Processor, msg Message) <-chan Message ***REMOVED***
	ch := make(chan Message)
	wg := sync.WaitGroup***REMOVED******REMOVED***

	// Dispatch processing across a number of processors, using a WaitGroup to record the
	// completion of each one
	for _, processor := range processors ***REMOVED***
		processor := processor
		wg.Add(1)
		go func() ***REMOVED***
			// No matter what happens, mark this processor as done once this goroutine returns
			defer wg.Done()

			// Forward resulting messages from the processor
			for m := range processor.Process(msg) ***REMOVED***
				ch <- m
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// Wait on the WaitGroup before closing the channel, signalling that we're done here
	go func() ***REMOVED***
		wg.Wait()
		close(ch)
	***REMOVED***()

	return ch
***REMOVED***
