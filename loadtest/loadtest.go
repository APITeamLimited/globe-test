package loadtest

import (
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
***REMOVED***
