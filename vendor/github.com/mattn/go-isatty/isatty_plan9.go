//go:build plan9
// +build plan9

package isatty

import (
	"syscall"
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd uintptr) bool ***REMOVED***
	path, err := syscall.Fd2path(int(fd))
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return path == "/dev/cons" || path == "/mnt/term/dev/cons"
***REMOVED***

// IsCygwinTerminal return true if the file descriptor is a cygwin or msys2
// terminal. This is also always false on this environment.
func IsCygwinTerminal(fd uintptr) bool ***REMOVED***
	return false
***REMOVED***
