// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package terminal provides support functions for dealing with terminals, as
// commonly found on UNIX systems.
//
// Putting a terminal into raw mode is the most common requirement:
//
// 	oldState, err := terminal.MakeRaw(0)
// 	if err != nil ***REMOVED***
// 	        panic(err)
// 	***REMOVED***
// 	defer terminal.Restore(0, oldState)
package terminal

import (
	"golang.org/x/sys/windows"
)

type State struct ***REMOVED***
	mode uint32
***REMOVED***

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool ***REMOVED***
	var st uint32
	err := windows.GetConsoleMode(windows.Handle(fd), &st)
	return err == nil
***REMOVED***

// MakeRaw put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*State, error) ***REMOVED***
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	raw := st &^ (windows.ENABLE_ECHO_INPUT | windows.ENABLE_PROCESSED_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_PROCESSED_OUTPUT)
	if err := windows.SetConsoleMode(windows.Handle(fd), raw); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &State***REMOVED***st***REMOVED***, nil
***REMOVED***

// GetState returns the current state of a terminal which may be useful to
// restore the terminal after a signal.
func GetState(fd int) (*State, error) ***REMOVED***
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &State***REMOVED***st***REMOVED***, nil
***REMOVED***

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func Restore(fd int, state *State) error ***REMOVED***
	return windows.SetConsoleMode(windows.Handle(fd), state.mode)
***REMOVED***

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (width, height int, err error) ***REMOVED***
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(fd), &info); err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	return int(info.Size.X), int(info.Size.Y), nil
***REMOVED***

// passwordReader is an io.Reader that reads from a specific Windows HANDLE.
type passwordReader int

func (r passwordReader) Read(buf []byte) (int, error) ***REMOVED***
	return windows.Read(windows.Handle(r), buf)
***REMOVED***

// ReadPassword reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func ReadPassword(fd int) ([]byte, error) ***REMOVED***
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	old := st

	st &^= (windows.ENABLE_ECHO_INPUT)
	st |= (windows.ENABLE_PROCESSED_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_PROCESSED_OUTPUT)
	if err := windows.SetConsoleMode(windows.Handle(fd), st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer func() ***REMOVED***
		windows.SetConsoleMode(windows.Handle(fd), old)
	***REMOVED***()

	return readPasswordLine(passwordReader(fd))
***REMOVED***
