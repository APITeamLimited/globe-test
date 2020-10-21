// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build aix
// +build ppc

package unix

//sysnb	Getrlimit(resource int, rlim *Rlimit) (err error) = getrlimit64
//sysnb	Setrlimit(resource int, rlim *Rlimit) (err error) = setrlimit64
//sys	Seek(fd int, offset int64, whence int) (off int64, err error) = lseek64

//sys	mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error)

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: int32(sec), Nsec: int32(nsec)***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: int32(sec), Usec: int32(usec)***REMOVED***
***REMOVED***

func (iov *Iovec) SetLen(length int) ***REMOVED***
	iov.Len = uint32(length)
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

func Fstat(fd int, stat *Stat_t) error ***REMOVED***
	return fstat(fd, stat)
***REMOVED***

func Fstatat(dirfd int, path string, stat *Stat_t, flags int) error ***REMOVED***
	return fstatat(dirfd, path, stat, flags)
***REMOVED***

func Lstat(path string, stat *Stat_t) error ***REMOVED***
	return lstat(path, stat)
***REMOVED***

func Stat(path string, statptr *Stat_t) error ***REMOVED***
	return stat(path, statptr)
***REMOVED***
