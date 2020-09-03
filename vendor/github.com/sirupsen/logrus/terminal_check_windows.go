// +build !appengine,!js,windows

package logrus

import (
	"io"
	"os"
	"syscall"

	sequences "github.com/konsorten/go-windows-terminal-sequences"
)

func initTerminal(w io.Writer) ***REMOVED***
	switch v := w.(type) ***REMOVED***
	case *os.File:
		sequences.EnableVirtualTerminalProcessing(syscall.Handle(v.Fd()), true)
	***REMOVED***
***REMOVED***

func checkIfTerminal(w io.Writer) bool ***REMOVED***
	var ret bool
	switch v := w.(type) ***REMOVED***
	case *os.File:
		var mode uint32
		err := syscall.GetConsoleMode(syscall.Handle(v.Fd()), &mode)
		ret = (err == nil)
	default:
		ret = false
	***REMOVED***
	if ret ***REMOVED***
		initTerminal(w)
	***REMOVED***
	return ret
***REMOVED***
