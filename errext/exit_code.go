package errext

import (
	"errors"

	"github.com/APITeamLimited/k6-worker/errext/exitcodes"
)

// ExitCode is the code with which the application should exit if this error
// bubbles up to the top of the scope. Values should be between 0 and 125:
// https://unix.stackexchange.com/questions/418784/what-is-the-min-and-max-values-of-exit-codes-in-linux

// HasExitCode is a wrapper around an error with an attached exit code.
type HasExitCode interface ***REMOVED***
	error
	ExitCode() exitcodes.ExitCode
***REMOVED***

// WithExitCodeIfNone can attach an exit code to the given error, if it doesn't
// have one already. It won't do anything if the error already had an exit code
// attached. Similarly, if there is no error (i.e. the given error is nil), it
// also won't do anything.
func WithExitCodeIfNone(err error, exitCode exitcodes.ExitCode) error ***REMOVED***
	if err == nil ***REMOVED***
		// No error, do nothing
		return nil
	***REMOVED***
	var ecerr HasExitCode
	if errors.As(err, &ecerr) ***REMOVED***
		// The given error already has an exit code, do nothing
		return err
	***REMOVED***
	return withExitCode***REMOVED***err, exitCode***REMOVED***
***REMOVED***

type withExitCode struct ***REMOVED***
	error
	exitCode exitcodes.ExitCode
***REMOVED***

func (wh withExitCode) Unwrap() error ***REMOVED***
	return wh.error
***REMOVED***

func (wh withExitCode) ExitCode() exitcodes.ExitCode ***REMOVED***
	return wh.exitCode
***REMOVED***

var _ HasExitCode = withExitCode***REMOVED******REMOVED***
