// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm64 && netbsd
// +build arm64,netbsd

package unix

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: sec, Nsec: nsec***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: sec, Usec: int32(usec)***REMOVED***
***REMOVED***

func SetKevent(k *Kevent_t, fd, mode, flags int) ***REMOVED***
	k.Ident = uint64(fd)
	k.Filter = uint32(mode)
	k.Flags = uint32(flags)
***REMOVED***

func (iov *Iovec) SetLen(length int) ***REMOVED***
	iov.Len = uint64(length)
***REMOVED***

func (msghdr *Msghdr) SetControllen(length int) ***REMOVED***
	msghdr.Controllen = uint32(length)
***REMOVED***

func (msghdr *Msghdr) SetIovlen(length int) ***REMOVED***
	msghdr.Iovlen = int32(length)
***REMOVED***

func (cmsg *Cmsghdr) SetLen(length int) ***REMOVED***
	cmsg.Len = uint32(length)
***REMOVED***
