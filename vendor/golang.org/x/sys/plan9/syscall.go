// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9
// +build plan9

// Package plan9 contains an interface to the low-level operating system
// primitives. OS details vary depending on the underlying system, and
// by default, godoc will display the OS-specific documentation for the current
// system. If you want godoc to display documentation for another
// system, set $GOOS and $GOARCH to the desired system. For example, if
// you want to view documentation for freebsd/arm on linux/amd64, set $GOOS
// to freebsd and $GOARCH to arm.
//
// The primary use of this package is inside other packages that provide a more
// portable interface to the system, such as "os", "time" and "net".  Use
// those packages rather than this one if you can.
//
// For details of the functions and data types in this package consult
// the manuals for the appropriate operating system.
//
// These calls return err == nil to indicate success; otherwise
// err represents an operating system error describing the failure and
// holds a value of type syscall.ErrorString.
package plan9 // import "golang.org/x/sys/plan9"

import (
	"bytes"
	"strings"
	"unsafe"

	"golang.org/x/sys/internal/unsafeheader"
)

// ByteSliceFromString returns a NUL-terminated slice of bytes
// containing the text of s. If s contains a NUL byte at any
// location, it returns (nil, EINVAL).
func ByteSliceFromString(s string) ([]byte, error) ***REMOVED***
	if strings.IndexByte(s, 0) != -1 ***REMOVED***
		return nil, EINVAL
	***REMOVED***
	a := make([]byte, len(s)+1)
	copy(a, s)
	return a, nil
***REMOVED***

// BytePtrFromString returns a pointer to a NUL-terminated array of
// bytes containing the text of s. If s contains a NUL byte at any
// location, it returns (nil, EINVAL).
func BytePtrFromString(s string) (*byte, error) ***REMOVED***
	a, err := ByteSliceFromString(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &a[0], nil
***REMOVED***

// ByteSliceToString returns a string form of the text represented by the slice s, with a terminating NUL and any
// bytes after the NUL removed.
func ByteSliceToString(s []byte) string ***REMOVED***
	if i := bytes.IndexByte(s, 0); i != -1 ***REMOVED***
		s = s[:i]
	***REMOVED***
	return string(s)
***REMOVED***

// BytePtrToString takes a pointer to a sequence of text and returns the corresponding string.
// If the pointer is nil, it returns the empty string. It assumes that the text sequence is terminated
// at a zero byte; if the zero byte is not present, the program may crash.
func BytePtrToString(p *byte) string ***REMOVED***
	if p == nil ***REMOVED***
		return ""
	***REMOVED***
	if *p == 0 ***REMOVED***
		return ""
	***REMOVED***

	// Find NUL terminator.
	n := 0
	for ptr := unsafe.Pointer(p); *(*byte)(ptr) != 0; n++ ***REMOVED***
		ptr = unsafe.Pointer(uintptr(ptr) + 1)
	***REMOVED***

	var s []byte
	h := (*unsafeheader.Slice)(unsafe.Pointer(&s))
	h.Data = unsafe.Pointer(p)
	h.Len = n
	h.Cap = n

	return string(s)
***REMOVED***

// Single-word zero for use when we need a valid pointer to 0 bytes.
// See mksyscall.pl.
var _zero uintptr

func (ts *Timespec) Unix() (sec int64, nsec int64) ***REMOVED***
	return int64(ts.Sec), int64(ts.Nsec)
***REMOVED***

func (tv *Timeval) Unix() (sec int64, nsec int64) ***REMOVED***
	return int64(tv.Sec), int64(tv.Usec) * 1000
***REMOVED***

func (ts *Timespec) Nano() int64 ***REMOVED***
	return int64(ts.Sec)*1e9 + int64(ts.Nsec)
***REMOVED***

func (tv *Timeval) Nano() int64 ***REMOVED***
	return int64(tv.Sec)*1e9 + int64(tv.Usec)*1000
***REMOVED***

// use is a no-op, but the compiler cannot see that it is.
// Calling use(p) ensures that p is kept live until that point.
//
//go:noescape
func use(p unsafe.Pointer)
