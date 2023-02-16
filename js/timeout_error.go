package js

import (
	"fmt"
	"time"

	"github.com/APITeamLimited/globe-test/errext"
	"github.com/APITeamLimited/globe-test/errext/exitcodes"
	"github.com/APITeamLimited/globe-test/worker/libWorker/consts"
)

// timeoutError is used when some operation times out.
type timeoutError struct {
	place string
	d     time.Duration
}

var (
	_ errext.HasExitCode = timeoutError{}
)

// newTimeoutError returns a new timeout error, reporting that a timeout has
// happened at the given place and given duration.
func newTimeoutError(place string, d time.Duration) timeoutError {
	return timeoutError{place: place, d: d}
}

// String returns the timeout error in human readable format.
func (t timeoutError) Error() string {
	return fmt.Sprintf("%s() execution timed out after %.f seconds", t.place, t.d.Seconds())
}

// ExitCode returns the coresponding exit code value to the place.
func (t timeoutError) ExitCode() exitcodes.ExitCode {
	switch t.place {
	case consts.SetupFn:
		return exitcodes.SetupTimeout
	case consts.TeardownFn:
		return exitcodes.TeardownTimeout
	default:
		return exitcodes.GenericTimeout
	}
}
