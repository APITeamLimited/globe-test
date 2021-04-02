// +build appengine

package colorable

import (
	"io"
	"os"

	_ "github.com/mattn/go-isatty"
)

// NewColorable returns new instance of Writer which handles escape sequence.
func NewColorable(file *os.File) io.Writer ***REMOVED***
	if file == nil ***REMOVED***
		panic("nil passed instead of *os.File to NewColorable()")
	***REMOVED***

	return file
***REMOVED***

// NewColorableStdout returns new instance of Writer which handles escape sequence for stdout.
func NewColorableStdout() io.Writer ***REMOVED***
	return os.Stdout
***REMOVED***

// NewColorableStderr returns new instance of Writer which handles escape sequence for stderr.
func NewColorableStderr() io.Writer ***REMOVED***
	return os.Stderr
***REMOVED***

// EnableColorsStdout enable colors if possible.
func EnableColorsStdout(enabled *bool) func() ***REMOVED***
	if enabled != nil ***REMOVED***
		*enabled = true
	***REMOVED***
	return func() ***REMOVED******REMOVED***
***REMOVED***
