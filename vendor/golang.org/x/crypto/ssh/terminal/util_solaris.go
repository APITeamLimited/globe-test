// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package terminal // import "golang.org/x/crypto/ssh/terminal"

import (
	"golang.org/x/sys/unix"
	"io"
	"syscall"
)

// State contains the state of a terminal.
type State struct ***REMOVED***
	termios unix.Termios
***REMOVED***

// IsTerminal returns whether the given file descriptor is a terminal.
func IsTerminal(fd int) bool ***REMOVED***
	_, err := unix.IoctlGetTermio(fd, unix.TCGETA)
	return err == nil
***REMOVED***

// ReadPassword reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func ReadPassword(fd int) ([]byte, error) ***REMOVED***
	// see also: http://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libast/common/uwin/getpass.c
	val, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	oldState := *val

	newState := oldState
	newState.Lflag &^= syscall.ECHO
	newState.Lflag |= syscall.ICANON | syscall.ISIG
	newState.Iflag |= syscall.ICRNL
	err = unix.IoctlSetTermios(fd, unix.TCSETS, &newState)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	defer unix.IoctlSetTermios(fd, unix.TCSETS, &oldState)

	var buf [16]byte
	var ret []byte
	for ***REMOVED***
		n, err := syscall.Read(fd, buf[:])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if n == 0 ***REMOVED***
			if len(ret) == 0 ***REMOVED***
				return nil, io.EOF
			***REMOVED***
			break
		***REMOVED***
		if buf[n-1] == '\n' ***REMOVED***
			n--
		***REMOVED***
		ret = append(ret, buf[:n]...)
		if n < len(buf) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return ret, nil
***REMOVED***

// MakeRaw puts the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
// see http://cr.illumos.org/~webrev/andy_js/1060/
func MakeRaw(fd int) (*State, error) ***REMOVED***
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	oldState := State***REMOVED***termios: *termios***REMOVED***

	termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &oldState, nil
***REMOVED***

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func Restore(fd int, oldState *State) error ***REMOVED***
	return unix.IoctlSetTermios(fd, unix.TCSETS, &oldState.termios)
***REMOVED***

// GetState returns the current state of a terminal which may be useful to
// restore the terminal after a signal.
func GetState(fd int) (*State, error) ***REMOVED***
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &State***REMOVED***termios: *termios***REMOVED***, nil
***REMOVED***

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (width, height int, err error) ***REMOVED***
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***
	return int(ws.Col), int(ws.Row), nil
***REMOVED***
