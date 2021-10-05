// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term

import (
	"os"

	"golang.org/x/sys/windows"
)

type state struct ***REMOVED***
	mode uint32
***REMOVED***

func isTerminal(fd int) bool ***REMOVED***
	var st uint32
	err := windows.GetConsoleMode(windows.Handle(fd), &st)
	return err == nil
***REMOVED***

func makeRaw(fd int) (*State, error) ***REMOVED***
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	raw := st &^ (windows.ENABLE_ECHO_INPUT | windows.ENABLE_PROCESSED_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_PROCESSED_OUTPUT)
	if err := windows.SetConsoleMode(windows.Handle(fd), raw); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &State***REMOVED***state***REMOVED***st***REMOVED******REMOVED***, nil
***REMOVED***

func getState(fd int) (*State, error) ***REMOVED***
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &State***REMOVED***state***REMOVED***st***REMOVED******REMOVED***, nil
***REMOVED***

func restore(fd int, state *State) error ***REMOVED***
	return windows.SetConsoleMode(windows.Handle(fd), state.mode)
***REMOVED***

func getSize(fd int) (width, height int, err error) ***REMOVED***
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(fd), &info); err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	return int(info.Window.Right - info.Window.Left + 1), int(info.Window.Bottom - info.Window.Top + 1), nil
***REMOVED***

func readPassword(fd int) ([]byte, error) ***REMOVED***
	var st uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	old := st

	st &^= (windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT)
	st |= (windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_PROCESSED_INPUT)
	if err := windows.SetConsoleMode(windows.Handle(fd), st); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer windows.SetConsoleMode(windows.Handle(fd), old)

	var h windows.Handle
	p, _ := windows.GetCurrentProcess()
	if err := windows.DuplicateHandle(p, windows.Handle(fd), p, &h, 0, false, windows.DUPLICATE_SAME_ACCESS); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	f := os.NewFile(uintptr(h), "stdin")
	defer f.Close()
	return readPasswordLine(f)
***REMOVED***
