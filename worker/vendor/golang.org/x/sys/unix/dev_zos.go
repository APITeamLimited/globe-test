// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build zos && s390x
// +build zos,s390x

// Functions to access/create device major and minor numbers matching the
// encoding used by z/OS.
//
// The information below is extracted and adapted from <sys/stat.h> macros.

package unix

// Major returns the major component of a z/OS device number.
func Major(dev uint64) uint32 ***REMOVED***
	return uint32((dev >> 16) & 0x0000FFFF)
***REMOVED***

// Minor returns the minor component of a z/OS device number.
func Minor(dev uint64) uint32 ***REMOVED***
	return uint32(dev & 0x0000FFFF)
***REMOVED***

// Mkdev returns a z/OS device number generated from the given major and minor
// components.
func Mkdev(major, minor uint32) uint64 ***REMOVED***
	return (uint64(major) << 16) | uint64(minor)
***REMOVED***
