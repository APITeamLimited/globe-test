// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package term provides support functions for dealing with terminals, as
// commonly found on UNIX systems.
//
// Putting a terminal into raw mode is the most common requirement:
//
// 	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
// 	if err != nil ***REMOVED***
// 	        panic(err)
// 	***REMOVED***
// 	defer term.Restore(int(os.Stdin.Fd()), oldState)
package term

// State contains the state of a terminal.
type State struct ***REMOVED***
	state
***REMOVED***

// IsTerminal returns whether the given file descriptor is a terminal.
func IsTerminal(fd int) bool ***REMOVED***
	return isTerminal(fd)
***REMOVED***

// MakeRaw puts the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*State, error) ***REMOVED***
	return makeRaw(fd)
***REMOVED***

// GetState returns the current state of a terminal which may be useful to
// restore the terminal after a signal.
func GetState(fd int) (*State, error) ***REMOVED***
	return getState(fd)
***REMOVED***

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func Restore(fd int, oldState *State) error ***REMOVED***
	return restore(fd, oldState)
***REMOVED***

// GetSize returns the visible dimensions of the given terminal.
//
// These dimensions don't include any scrollback buffer height.
func GetSize(fd int) (width, height int, err error) ***REMOVED***
	return getSize(fd)
***REMOVED***

// ReadPassword reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func ReadPassword(fd int) ([]byte, error) ***REMOVED***
	return readPassword(fd)
***REMOVED***
