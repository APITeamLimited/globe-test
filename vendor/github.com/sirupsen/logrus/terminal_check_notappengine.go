// +build !appengine

package logrus

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func checkIfTerminal(w io.Writer) bool ***REMOVED***
	switch v := w.(type) ***REMOVED***
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	***REMOVED***
***REMOVED***
