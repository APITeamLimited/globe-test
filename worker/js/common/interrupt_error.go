package common

import (
	"errors"

	"github.com/dop251/goja"
)

// UnwrapGojaInterruptedError returns the internal error handled by goja.
func UnwrapGojaInterruptedError(err error) error ***REMOVED***
	var gojaErr *goja.InterruptedError
	if errors.As(err, &gojaErr) ***REMOVED***
		if e, ok := gojaErr.Value().(error); ok ***REMOVED***
			return e
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***
