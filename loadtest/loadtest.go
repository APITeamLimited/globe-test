package loadtest

import (
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

	Source string // Script source
***REMOVED***

func (t *LoadTest) Load(base string) error ***REMOVED***
	srcb, err := ioutil.ReadFile(path.Join(base, t.Script))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	t.Source = string(srcb)
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
