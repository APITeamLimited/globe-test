package errext

import (
	"errors"

	"go.k6.io/k6/errext/exitcodes"
)

// InterruptError is an error that halts engine execution
type InterruptError struct ***REMOVED***
	Reason string
***REMOVED***

var _ HasExitCode = &InterruptError***REMOVED******REMOVED***

// Error returns the reason of the interruption.
func (i *InterruptError) Error() string ***REMOVED***
	return i.Reason
***REMOVED***

// ExitCode returns the status code used when the k6 process exits.
func (i *InterruptError) ExitCode() exitcodes.ExitCode ***REMOVED***
	return exitcodes.ScriptAborted
***REMOVED***

// AbortTest is the reason emitted when a test script calls test.abort()
const AbortTest = "test aborted"

// IsInterruptError returns true if err is *InterruptError.
func IsInterruptError(err error) bool ***REMOVED***
	if err == nil ***REMOVED***
		return false
	***REMOVED***
	var intErr *InterruptError
	return errors.As(err, &intErr)
***REMOVED***
