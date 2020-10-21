// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build aix
// +build ppc64

package unix

//sysnb	Getrlimit(resource int, rlim *Rlimit) (err error)
//sysnb	Setrlimit(resource int, rlim *Rlimit) (err error)
//sys	Seek(fd int, offset int64, whence int) (off int64, err error) = lseek

//sys	mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) = mmap64

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: sec, Nsec: nsec***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: int64(sec), Usec: int32(usec)***REMOVED***
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

// In order to only have Timespec structure, type of Stat_t's fields
// Atim, Mtim and Ctim is changed from StTimespec to Timespec during
// ztypes generation.
// On ppc64, Timespec.Nsec is an int64 while StTimespec.Nsec is an
// int32, so the fields' value must be modified.
func fixStatTimFields(stat *Stat_t) ***REMOVED***
	stat.Atim.Nsec >>= 32
	stat.Mtim.Nsec >>= 32
	stat.Ctim.Nsec >>= 32
***REMOVED***

func Fstat(fd int, stat *Stat_t) error ***REMOVED***
	err := fstat(fd, stat)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	fixStatTimFields(stat)
	return nil
***REMOVED***

func Fstatat(dirfd int, path string, stat *Stat_t, flags int) error ***REMOVED***
	err := fstatat(dirfd, path, stat, flags)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	fixStatTimFields(stat)
	return nil
***REMOVED***

func Lstat(path string, stat *Stat_t) error ***REMOVED***
	err := lstat(path, stat)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	fixStatTimFields(stat)
	return nil
***REMOVED***

func Stat(path string, statptr *Stat_t) error ***REMOVED***
	err := stat(path, statptr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	fixStatTimFields(statptr)
	return nil
***REMOVED***
