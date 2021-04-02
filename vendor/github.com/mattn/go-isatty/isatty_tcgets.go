// +build linux aix
// +build !appengine

package isatty

import "golang.org/x/sys/unix"

// IsTerminal return true if the file descriptor is terminal.
func IsTerminal(fd uintptr) bool ***REMOVED***
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	return err == nil
***REMOVED***

// IsCygwinTerminal return true if the file descriptor is a cygwin or msys2
// terminal. This is also always false on this environment.
func IsCygwinTerminal(fd uintptr) bool ***REMOVED***
	return false
***REMOVED***
