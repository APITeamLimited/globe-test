// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build arm,freebsd

package unix

import (
	"syscall"
	"unsafe"
)

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: sec, Nsec: int32(nsec)***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: sec, Usec: int32(usec)***REMOVED***
***REMOVED***

func SetKevent(k *Kevent_t, fd, mode, flags int) ***REMOVED***
	k.Ident = uint32(fd)
	k.Filter = int16(mode)
	k.Flags = uint16(flags)
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

func sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) ***REMOVED***
	var writtenOut uint64 = 0
	_, _, e1 := Syscall9(SYS_SENDFILE, uintptr(infd), uintptr(outfd), uintptr(*offset), uintptr((*offset)>>32), uintptr(count), 0, uintptr(unsafe.Pointer(&writtenOut)), 0, 0)

	written = int(writtenOut)

	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***

func Syscall9(num, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r1, r2 uintptr, err syscall.Errno)

func PtraceIO(req int, pid int, addr uintptr, out []byte, countin int) (count int, err error) ***REMOVED***
	ioDesc := PtraceIoDesc***REMOVED***Op: int32(req), Offs: (*byte)(unsafe.Pointer(addr)), Addr: (*byte)(unsafe.Pointer(&out[0])), Len: uint32(countin)***REMOVED***
	err = ptrace(PTRACE_IO, pid, uintptr(unsafe.Pointer(&ioDesc)), 0)
	return int(ioDesc.Len), err
***REMOVED***
