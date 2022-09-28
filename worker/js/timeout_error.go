package js

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/k6-worker/errext"
	"github.com/APITeamLimited/k6-worker/errext/exitcodes"
	"github.com/APITeamLimited/k6-worker/lib/consts"
)

// timeoutError is used when some operation times out.
type timeoutError struct ***REMOVED***
	place string
	d     time.Duration
***REMOVED***

var (
	_ errext.HasExitCode = timeoutError***REMOVED******REMOVED***
	_ errext.HasHint     = timeoutError***REMOVED******REMOVED***
)

// newTimeoutError returns a new timeout error, reporting that a timeout has
// happened at the given place and given duration.
func newTimeoutError(place string, d time.Duration) timeoutError ***REMOVED***
	return timeoutError***REMOVED***place: place, d: d***REMOVED***
***REMOVED***

// String returns the timeout error in human readable format.
func (t timeoutError) Error() string ***REMOVED***
	return fmt.Sprintf("%s() execution timed out after %.f seconds", t.place, t.d.Seconds())
***REMOVED***

// Hint potentially returns a hint message for fixing the error.
func (t timeoutError) Hint() string ***REMOVED***
	hint := ""

	switch t.place ***REMOVED***
	case consts.SetupFn:
		hint = "You can increase the time limit via the setupTimeout option"
	case consts.TeardownFn:
		hint = "You can increase the time limit via the teardownTimeout option"
	***REMOVED***
	return hint
***REMOVED***

// ExitCode returns the coresponding exit code value to the place.
func (t timeoutError) ExitCode() exitcodes.ExitCode ***REMOVED***
	// TODO: add handleSummary()
	switch t.place ***REMOVED***
	case consts.SetupFn:
		return exitcodes.SetupTimeout
	case consts.TeardownFn:
		return exitcodes.TeardownTimeout
	default:
		return exitcodes.GenericTimeout
	***REMOVED***
***REMOVED***
