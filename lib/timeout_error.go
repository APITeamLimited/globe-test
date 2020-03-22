package lib

import (
	"fmt"
	"time"
)

//nolint:gochecknoglobals
// Keep stages in sync with js/runner.go
// We set it here to prevent import cycle.
var (
	stageSetup    = "setup"
	stageTeardown = "teardown"
)

// TimeoutError is used when somethings timeouts
type TimeoutError struct ***REMOVED***
	place string
	d     time.Duration
***REMOVED***

// NewTimeoutError returns a new TimeoutError reporting that timeout has happened
// at the given place and given duration.
func NewTimeoutError(place string, d time.Duration) TimeoutError ***REMOVED***
	return TimeoutError***REMOVED***place: place, d: d***REMOVED***
***REMOVED***

// String returns timeout error in human readable format.
func (t TimeoutError) String() string ***REMOVED***
	return fmt.Sprintf("%s() execution timed out after %.f seconds", t.place, t.d.Seconds())
***REMOVED***

// Error implements error interface.
func (t TimeoutError) Error() string ***REMOVED***
	return t.String()
***REMOVED***

// Place returns the place where timeout occurred.
func (t TimeoutError) Place() string ***REMOVED***
	return t.place
***REMOVED***

// Hint returns a hint message for logging with given stage.
func (t TimeoutError) Hint() string ***REMOVED***
	hint := ""

	switch t.place ***REMOVED***
	case stageSetup:
		hint = "You can increase the time limit via the setupTimeout option"
	case stageTeardown:
		hint = "You can increase the time limit via the teardownTimeout option"
	***REMOVED***
	return hint
***REMOVED***
