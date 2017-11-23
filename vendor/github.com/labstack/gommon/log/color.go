// +build !appengine

package log

import (
	"io"

	"github.com/mattn/go-colorable"
)

func output() io.Writer ***REMOVED***
	return colorable.NewColorableStdout()
***REMOVED***
