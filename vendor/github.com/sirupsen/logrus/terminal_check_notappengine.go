// +build !appengine,!js,!windows,!nacl,!plan9

package logrus

import (
	"io"
	"os"
)

func checkIfTerminal(w io.Writer) bool ***REMOVED***
	switch v := w.(type) ***REMOVED***
	case *os.File:
		return isTerminal(int(v.Fd()))
	default:
		return false
	***REMOVED***
***REMOVED***
