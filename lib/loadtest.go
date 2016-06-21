package lib

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
***REMOVED***

func (t *Test) TotalDuration() time.Duration ***REMOVED***
	var total time.Duration
	for _, stage := range t.Stages ***REMOVED***
		total += stage.Duration
	***REMOVED***
	return total
***REMOVED***

func (t *Test) VUsAt(at time.Duration) int ***REMOVED***
	stageStart := time.Duration(0)
	for _, stage := range t.Stages ***REMOVED***
		if stageStart+stage.Duration < at ***REMOVED***
			stageStart += stage.Duration
			continue
		***REMOVED***
		progress := float64(at-stageStart) / float64(stage.Duration)
		return stage.StartVUs + int(float64(stage.EndVUs-stage.StartVUs)*progress)
	***REMOVED***
	return 0
***REMOVED***
