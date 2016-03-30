package loadtest

import (
	"time"
)

// A load test is composed of at least one "stage", which controls VU distribution.
type Stage struct ***REMOVED***
	Duration time.Duration // Duration of this stage.
	StartVUs int           // Set this many VUs at the start of the stage.
	EndVUs   int           // Ramp until there are this many VUs.
***REMOVED***

type LoadTest struct ***REMOVED***
	Stages []Stage // Test stages.
***REMOVED***
