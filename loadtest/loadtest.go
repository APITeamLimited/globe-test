package loadtest

import (
	"github.com/loadimpact/speedboat/message"
	"io/ioutil"
	"path"
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
	currentVUs   int
***REMOVED***

func (t *LoadTest) Load(base string) error ***REMOVED***
	srcb, err := ioutil.ReadFile(path.Join(base, t.Script))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	t.scriptSource = string(srcb)
	return nil
***REMOVED***

func (t *LoadTest) StageAt(d time.Duration) (start time.Duration, stage Stage, stop bool) ***REMOVED***
	at := time.Duration(0)
	for i := range t.Stages ***REMOVED***
		stage = t.Stages[i]
		if d > at+stage.Duration ***REMOVED***
			at += stage.Duration
		***REMOVED*** else if d < at+stage.Duration ***REMOVED***
			return at, stage, false
		***REMOVED***
	***REMOVED***
	return at, stage, true
***REMOVED***

func (t *LoadTest) VUsAt(at time.Duration) (vus int, stop bool) ***REMOVED***
	start, stage, stop := t.StageAt(at)
	if stop ***REMOVED***
		return 0, true
	***REMOVED***

	stageElapsed := at - start
	percentage := (stageElapsed.Seconds() / stage.Duration.Seconds())
	vus = stage.VUs.Start + int(float64(stage.VUs.End-stage.VUs.Start)*percentage)

	return vus, false
***REMOVED***

func (t *LoadTest) Run(in <-chan message.Message, out chan message.Message) (<-chan message.Message, chan message.Message) ***REMOVED***
	oin := make(chan message.Message)

	go func() ***REMOVED***
		out <- message.ToWorker("test.run").With(MessageTestRun***REMOVED***
			Filename: t.Script,
			Source:   t.scriptSource,
			VUs:      t.Stages[0].VUs.Start,
		***REMOVED***)

		startTime := time.Now()
		intervene := time.Tick(time.Duration(1) * time.Second)
	runLoop:
		for ***REMOVED***
			select ***REMOVED***
			case msg := <-in:
				oin <- msg
			case <-intervene:
				vus, stop := t.VUsAt(time.Since(startTime))
				if stop ***REMOVED***
					out <- message.ToWorker("test.stop")
					close(oin)
					break runLoop
				***REMOVED***
				if vus != t.currentVUs ***REMOVED***
					out <- message.ToWorker("test.scale").With(MessageTestScale***REMOVED***VUs: vus***REMOVED***)
					t.currentVUs = vus
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	return oin, out
***REMOVED***
