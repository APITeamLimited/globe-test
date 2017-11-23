// +build !windows
// +build !appengine

package colorable

import (
	"io"
	"os"

	_ "github.com/mattn/go-isatty"
)

// NewColorable return new instance of Writer which handle escape sequence.
func NewColorable(file *os.File) io.Writer ***REMOVED***
	if file == nil ***REMOVED***
		panic("nil passed instead of *os.File to NewColorable()")
	***REMOVED***

	return file
***REMOVED***

// NewColorableStdout return new instance of Writer which handle escape sequence for stdout.
func NewColorableStdout() io.Writer ***REMOVED***
	return os.Stdout
***REMOVED***

// NewColorableStderr return new instance of Writer which handle escape sequence for stderr.
func NewColorableStderr() io.Writer ***REMOVED***
	return os.Stderr
***REMOVED***
