// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package unix

import "time"

// TimespecToNSec returns the time stored in ts as nanoseconds.
func TimespecToNsec(ts Timespec) int64 ***REMOVED*** return ts.Nano() ***REMOVED***

// NsecToTimespec converts a number of nanoseconds into a Timespec.
func NsecToTimespec(nsec int64) Timespec ***REMOVED***
	sec := nsec / 1e9
	nsec = nsec % 1e9
	if nsec < 0 ***REMOVED***
		nsec += 1e9
		sec--
	***REMOVED***
	return setTimespec(sec, nsec)
***REMOVED***

// TimeToTimespec converts t into a Timespec.
// On some 32-bit systems the range of valid Timespec values are smaller
// than that of time.Time values.  So if t is out of the valid range of
// Timespec, it returns a zero Timespec and ERANGE.
func TimeToTimespec(t time.Time) (Timespec, error) ***REMOVED***
	sec := t.Unix()
	nsec := int64(t.Nanosecond())
	ts := setTimespec(sec, nsec)

	// Currently all targets have either int32 or int64 for Timespec.Sec.
	// If there were a new target with floating point type for it, we have
	// to consider the rounding error.
	if int64(ts.Sec) != sec ***REMOVED***
		return Timespec***REMOVED******REMOVED***, ERANGE
	***REMOVED***
	return ts, nil
***REMOVED***

// TimevalToNsec returns the time stored in tv as nanoseconds.
func TimevalToNsec(tv Timeval) int64 ***REMOVED*** return tv.Nano() ***REMOVED***

// NsecToTimeval converts a number of nanoseconds into a Timeval.
func NsecToTimeval(nsec int64) Timeval ***REMOVED***
	nsec += 999 // round up to microsecond
	usec := nsec % 1e9 / 1e3
	sec := nsec / 1e9
	if usec < 0 ***REMOVED***
		usec += 1e6
		sec--
	***REMOVED***
	return setTimeval(sec, usec)
***REMOVED***

// Unix returns the time stored in ts as seconds plus nanoseconds.
func (ts *Timespec) Unix() (sec int64, nsec int64) ***REMOVED***
	return int64(ts.Sec), int64(ts.Nsec)
***REMOVED***

// Unix returns the time stored in tv as seconds plus nanoseconds.
func (tv *Timeval) Unix() (sec int64, nsec int64) ***REMOVED***
	return int64(tv.Sec), int64(tv.Usec) * 1000
***REMOVED***

// Nano returns the time stored in ts as nanoseconds.
func (ts *Timespec) Nano() int64 ***REMOVED***
	return int64(ts.Sec)*1e9 + int64(ts.Nsec)
***REMOVED***

// Nano returns the time stored in tv as nanoseconds.
func (tv *Timeval) Nano() int64 ***REMOVED***
	return int64(tv.Sec)*1e9 + int64(tv.Usec)*1000
***REMOVED***
