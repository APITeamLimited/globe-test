// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package term

import (
	"golang.org/x/sys/unix"
)

type state struct ***REMOVED***
	termios unix.Termios
***REMOVED***

func isTerminal(fd int) bool ***REMOVED***
	_, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	return err == nil
***REMOVED***

func makeRaw(fd int) (*State, error) ***REMOVED***
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	oldState := State***REMOVED***state***REMOVED***termios: *termios***REMOVED******REMOVED***

	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, termios); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &oldState, nil
***REMOVED***

func getState(fd int) (*State, error) ***REMOVED***
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &State***REMOVED***state***REMOVED***termios: *termios***REMOVED******REMOVED***, nil
***REMOVED***

func restore(fd int, state *State) error ***REMOVED***
	return unix.IoctlSetTermios(fd, ioctlWriteTermios, &state.termios)
***REMOVED***

func getSize(fd int) (width, height int, err error) ***REMOVED***
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil ***REMOVED***
		return -1, -1, err
	***REMOVED***
	return int(ws.Col), int(ws.Row), nil
***REMOVED***

// passwordReader is an io.Reader that reads from a specific file descriptor.
type passwordReader int

func (r passwordReader) Read(buf []byte) (int, error) ***REMOVED***
	return unix.Read(int(r), buf)
***REMOVED***

func readPassword(fd int) ([]byte, error) ***REMOVED***
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	newState := *termios
	newState.Lflag &^= unix.ECHO
	newState.Lflag |= unix.ICANON | unix.ISIG
	newState.Iflag |= unix.ICRNL
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)

	return readPasswordLine(passwordReader(fd))
***REMOVED***
