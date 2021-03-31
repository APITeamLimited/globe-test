// +build !appengine,!js,windows

package logrus

import (
	"io"
	"os"

	"golang.org/x/sys/windows"
)

func checkIfTerminal(w io.Writer) bool ***REMOVED***
	switch v := w.(type) ***REMOVED***
	case *os.File:
		handle := windows.Handle(v.Fd())
		var mode uint32
		if err := windows.GetConsoleMode(handle, &mode); err != nil ***REMOVED***
			return false
		***REMOVED***
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		if err := windows.SetConsoleMode(handle, mode); err != nil ***REMOVED***
			return false
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***
