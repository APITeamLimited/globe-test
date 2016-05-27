package speedboat

import (
	"time"
)

// A load test is composed of at least one "stage", which controls VU distribution.
type TestStage struct ***REMOVED***
	Duration time.Duration // Duration of this stage.
	StartVUs int           // VUs at the start of this stage.
	EndVUs   int           // VUs at the end of this stage.
***REMOVED***

// A load test definition.
type Test struct ***REMOVED***
	Script string      // Script filename.
	URL    string      // URL for simple tests.
	Stages []TestStage // Test stages.

	Source string // Script source
***REMOVED***

/*func (t *LoadTest) StageAt(d time.Duration) (start time.Duration, stage TestStage, stop bool) ***REMOVED***
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
	vus = stage.VUs.Start + int(float64(stage.EndVUs-stage.StartVUs)*percentage)

	return vus, false
***REMOVED****/
