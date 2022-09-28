// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unix

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

// Pledge implements the pledge syscall.
//
// The pledge syscall does not accept execpromises on OpenBSD releases
// before 6.3.
//
// execpromises must be empty when Pledge is called on OpenBSD
// releases predating 6.3, otherwise an error will be returned.
//
// For more information see pledge(2).
func Pledge(promises, execpromises string) error ***REMOVED***
	maj, min, err := majmin()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = pledgeAvailable(maj, min, execpromises)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	pptr, err := syscall.BytePtrFromString(promises)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// This variable will hold either a nil unsafe.Pointer or
	// an unsafe.Pointer to a string (execpromises).
	var expr unsafe.Pointer

	// If we're running on OpenBSD > 6.2, pass execpromises to the syscall.
	if maj > 6 || (maj == 6 && min > 2) ***REMOVED***
		exptr, err := syscall.BytePtrFromString(execpromises)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		expr = unsafe.Pointer(exptr)
	***REMOVED***

	_, _, e := syscall.Syscall(SYS_PLEDGE, uintptr(unsafe.Pointer(pptr)), uintptr(expr), 0)
	if e != 0 ***REMOVED***
		return e
	***REMOVED***

	return nil
***REMOVED***

// PledgePromises implements the pledge syscall.
//
// This changes the promises and leaves the execpromises untouched.
//
// For more information see pledge(2).
func PledgePromises(promises string) error ***REMOVED***
	maj, min, err := majmin()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = pledgeAvailable(maj, min, "")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// This variable holds the execpromises and is always nil.
	var expr unsafe.Pointer

	pptr, err := syscall.BytePtrFromString(promises)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	_, _, e := syscall.Syscall(SYS_PLEDGE, uintptr(unsafe.Pointer(pptr)), uintptr(expr), 0)
	if e != 0 ***REMOVED***
		return e
	***REMOVED***

	return nil
***REMOVED***

// PledgeExecpromises implements the pledge syscall.
//
// This changes the execpromises and leaves the promises untouched.
//
// For more information see pledge(2).
func PledgeExecpromises(execpromises string) error ***REMOVED***
	maj, min, err := majmin()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = pledgeAvailable(maj, min, execpromises)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// This variable holds the promises and is always nil.
	var pptr unsafe.Pointer

	exptr, err := syscall.BytePtrFromString(execpromises)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	_, _, e := syscall.Syscall(SYS_PLEDGE, uintptr(pptr), uintptr(unsafe.Pointer(exptr)), 0)
	if e != 0 ***REMOVED***
		return e
	***REMOVED***

	return nil
***REMOVED***

// majmin returns major and minor version number for an OpenBSD system.
func majmin() (major int, minor int, err error) ***REMOVED***
	var v Utsname
	err = Uname(&v)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	major, err = strconv.Atoi(string(v.Release[0]))
	if err != nil ***REMOVED***
		err = errors.New("cannot parse major version number returned by uname")
		return
	***REMOVED***

	minor, err = strconv.Atoi(string(v.Release[2]))
	if err != nil ***REMOVED***
		err = errors.New("cannot parse minor version number returned by uname")
		return
	***REMOVED***

	return
***REMOVED***

// pledgeAvailable checks for availability of the pledge(2) syscall
// based on the running OpenBSD version.
func pledgeAvailable(maj, min int, execpromises string) error ***REMOVED***
	// If OpenBSD <= 5.9, pledge is not available.
	if (maj == 5 && min != 9) || maj < 5 ***REMOVED***
		return fmt.Errorf("pledge syscall is not available on OpenBSD %d.%d", maj, min)
	***REMOVED***

	// If OpenBSD <= 6.2 and execpromises is not empty,
	// return an error - execpromises is not available before 6.3
	if (maj < 6 || (maj == 6 && min <= 2)) && execpromises != "" ***REMOVED***
		return fmt.Errorf("cannot use execpromises on OpenBSD %d.%d", maj, min)
	***REMOVED***

	return nil
***REMOVED***
