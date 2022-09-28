// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gccgo && linux && amd64
// +build gccgo,linux,amd64

package unix

import "syscall"

//extern gettimeofday
func realGettimeofday(*Timeval, *byte) int32

func gettimeofday(tv *Timeval) (err syscall.Errno) ***REMOVED***
	r := realGettimeofday(tv, nil)
	if r < 0 ***REMOVED***
		return syscall.GetErrno()
	***REMOVED***
	return 0
***REMOVED***
