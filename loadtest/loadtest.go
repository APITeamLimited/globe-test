package loadtest

import (
	"github.com/loadimpact/speedboat/message"
	"io/ioutil"
	"time"
)

// Specification for a VU curve.
type VUSpec struct ***REMOVED***
	Start int // Start at this many
	End   int // Interpolate to this many
***REMOVED***

// A load test is composed of at least one "stage", which controls VU distribution.
type Stage struct ***REMOVED***
	Duration time.Duration // Duration of this stage.
	VUs      VUSpec        // VU specification
***REMOVED***

// A load test definition.
type LoadTest struct ***REMOVED***
	Script string  // Script filename.
	URL    string  // URL for simple tests.
	Stages []Stage // Test stages.

	scriptSource string
***REMOVED***

func (t *LoadTest) Load() error ***REMOVED***
	srcb, err := ioutil.ReadFile(t.Script)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	t.scriptSource = string(srcb)
	return nil
***REMOVED***

func (t *LoadTest) Run(in <-chan message.Message, out chan message.Message, errors <-chan error) (<-chan message.Message, chan message.Message, <-chan error) ***REMOVED***
	oin := make(chan message.Message)

	go func() ***REMOVED***
		out <- message.NewToWorker("run.run", message.Fields***REMOVED***
			"filename": t.Script,
			"src":      t.scriptSource,
			"vus":      t.Stages[0].VUs.Start,
		***REMOVED***)

		duration := time.Duration(0)
		for i := range t.Stages ***REMOVED***
			duration += t.Stages[i].Duration
		***REMOVED***

		timeout := time.After(duration)
	runLoop:
		for ***REMOVED***
			select ***REMOVED***
			case msg := <-in:
				oin <- msg
			case <-timeout:
				out <- message.NewToWorker("run.stop", message.Fields***REMOVED******REMOVED***)
				close(oin)
				break runLoop
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return oin, out, errors
***REMOVED***
