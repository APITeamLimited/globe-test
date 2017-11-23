// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,solaris

package unix

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: sec, Nsec: nsec***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: sec, Usec: usec***REMOVED***
***REMOVED***

func (iov *Iovec) SetLen(length int) ***REMOVED***
	iov.Len = uint64(length)
***REMOVED***

func (cmsg *Cmsghdr) SetLen(length int) ***REMOVED***
	cmsg.Len = uint32(length)
***REMOVED***

func sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) ***REMOVED***
	// TODO(aram): implement this, see issue 5847.
	panic("unimplemented")
***REMOVED***
