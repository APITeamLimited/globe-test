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
