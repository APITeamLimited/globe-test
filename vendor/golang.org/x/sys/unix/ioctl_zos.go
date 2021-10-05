// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build zos && s390x
// +build zos,s390x

package unix

import (
	"runtime"
	"unsafe"
)

// ioctl itself should not be exposed directly, but additional get/set
// functions for specific types are permissible.

// IoctlSetInt performs an ioctl operation which sets an integer value
// on fd, using the specified request number.
func IoctlSetInt(fd int, req uint, value int) error ***REMOVED***
	return ioctl(fd, req, uintptr(value))
***REMOVED***

// IoctlSetWinsize performs an ioctl on fd with a *Winsize argument.
//
// To change fd's window size, the req argument should be TIOCSWINSZ.
func IoctlSetWinsize(fd int, req uint, value *Winsize) error ***REMOVED***
	// TODO: if we get the chance, remove the req parameter and
	// hardcode TIOCSWINSZ.
	err := ioctl(fd, req, uintptr(unsafe.Pointer(value)))
	runtime.KeepAlive(value)
	return err
***REMOVED***

// IoctlSetTermios performs an ioctl on fd with a *Termios.
//
// The req value is expected to be TCSETS, TCSETSW, or TCSETSF
func IoctlSetTermios(fd int, req uint, value *Termios) error ***REMOVED***
	if (req != TCSETS) && (req != TCSETSW) && (req != TCSETSF) ***REMOVED***
		return ENOSYS
	***REMOVED***
	err := Tcsetattr(fd, int(req), value)
	runtime.KeepAlive(value)
	return err
***REMOVED***

// IoctlGetInt performs an ioctl operation which gets an integer value
// from fd, using the specified request number.
//
// A few ioctl requests use the return value as an output parameter;
// for those, IoctlRetInt should be used instead of this function.
func IoctlGetInt(fd int, req uint) (int, error) ***REMOVED***
	var value int
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return value, err
***REMOVED***

func IoctlGetWinsize(fd int, req uint) (*Winsize, error) ***REMOVED***
	var value Winsize
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
***REMOVED***

// IoctlGetTermios performs an ioctl on fd with a *Termios.
//
// The req value is expected to be TCGETS
func IoctlGetTermios(fd int, req uint) (*Termios, error) ***REMOVED***
	var value Termios
	if req != TCGETS ***REMOVED***
		return &value, ENOSYS
	***REMOVED***
	err := Tcgetattr(fd, &value)
	return &value, err
***REMOVED***
