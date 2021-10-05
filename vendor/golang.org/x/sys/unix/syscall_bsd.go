// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || dragonfly || freebsd || netbsd || openbsd
// +build darwin dragonfly freebsd netbsd openbsd

// BSD system call wrappers shared by *BSD based systems
// including OS X (Darwin) and FreeBSD.  Like the other
// syscall_*.go files it is compiled as Go code but also
// used as input to mksyscall which parses the //sys
// lines and generates system call stubs.

package unix

import (
	"runtime"
	"syscall"
	"unsafe"
)

const ImplementsGetwd = true

func Getwd() (string, error) ***REMOVED***
	var buf [PathMax]byte
	_, err := Getcwd(buf[0:])
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	n := clen(buf[:])
	if n < 1 ***REMOVED***
		return "", EINVAL
	***REMOVED***
	return string(buf[:n]), nil
***REMOVED***

/*
 * Wrapped
 */

//sysnb	getgroups(ngid int, gid *_Gid_t) (n int, err error)
//sysnb	setgroups(ngid int, gid *_Gid_t) (err error)

func Getgroups() (gids []int, err error) ***REMOVED***
	n, err := getgroups(0, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if n == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	// Sanity check group count. Max is 16 on BSD.
	if n < 0 || n > 1000 ***REMOVED***
		return nil, EINVAL
	***REMOVED***

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	gids = make([]int, n)
	for i, v := range a[0:n] ***REMOVED***
		gids[i] = int(v)
	***REMOVED***
	return
***REMOVED***

func Setgroups(gids []int) (err error) ***REMOVED***
	if len(gids) == 0 ***REMOVED***
		return setgroups(0, nil)
	***REMOVED***

	a := make([]_Gid_t, len(gids))
	for i, v := range gids ***REMOVED***
		a[i] = _Gid_t(v)
	***REMOVED***
	return setgroups(len(a), &a[0])
***REMOVED***

// Wait status is 7 bits at bottom, either 0 (exited),
// 0x7F (stopped), or a signal number that caused an exit.
// The 0x80 bit is whether there was a core dump.
// An extra number (exit code, signal causing a stop)
// is in the high bits.

type WaitStatus uint32

const (
	mask  = 0x7F
	core  = 0x80
	shift = 8

	exited  = 0
	killed  = 9
	stopped = 0x7F
)

func (w WaitStatus) Exited() bool ***REMOVED*** return w&mask == exited ***REMOVED***

func (w WaitStatus) ExitStatus() int ***REMOVED***
	if w&mask != exited ***REMOVED***
		return -1
	***REMOVED***
	return int(w >> shift)
***REMOVED***

func (w WaitStatus) Signaled() bool ***REMOVED*** return w&mask != stopped && w&mask != 0 ***REMOVED***

func (w WaitStatus) Signal() syscall.Signal ***REMOVED***
	sig := syscall.Signal(w & mask)
	if sig == stopped || sig == 0 ***REMOVED***
		return -1
	***REMOVED***
	return sig
***REMOVED***

func (w WaitStatus) CoreDump() bool ***REMOVED*** return w.Signaled() && w&core != 0 ***REMOVED***

func (w WaitStatus) Stopped() bool ***REMOVED*** return w&mask == stopped && syscall.Signal(w>>shift) != SIGSTOP ***REMOVED***

func (w WaitStatus) Killed() bool ***REMOVED*** return w&mask == killed && syscall.Signal(w>>shift) != SIGKILL ***REMOVED***

func (w WaitStatus) Continued() bool ***REMOVED*** return w&mask == stopped && syscall.Signal(w>>shift) == SIGSTOP ***REMOVED***

func (w WaitStatus) StopSignal() syscall.Signal ***REMOVED***
	if !w.Stopped() ***REMOVED***
		return -1
	***REMOVED***
	return syscall.Signal(w>>shift) & 0xFF
***REMOVED***

func (w WaitStatus) TrapCause() int ***REMOVED*** return -1 ***REMOVED***

//sys	wait4(pid int, wstatus *_C_int, options int, rusage *Rusage) (wpid int, err error)

func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) ***REMOVED***
	var status _C_int
	wpid, err = wait4(pid, &status, options, rusage)
	if wstatus != nil ***REMOVED***
		*wstatus = WaitStatus(status)
	***REMOVED***
	return
***REMOVED***

//sys	accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error)
//sys	bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error)
//sys	connect(s int, addr unsafe.Pointer, addrlen _Socklen) (err error)
//sysnb	socket(domain int, typ int, proto int) (fd int, err error)
//sys	getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error)
//sys	setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error)
//sysnb	getpeername(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error)
//sysnb	getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error)
//sys	Shutdown(s int, how int) (err error)

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Port < 0 || sa.Port > 0xFFFF ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Len = SizeofSockaddrInet4
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
***REMOVED***

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Port < 0 || sa.Port > 0xFFFF ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Len = SizeofSockaddrInet6
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
***REMOVED***

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) || n == 0 ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Len = byte(3 + n) // 2 for Family, Len; 1 for NUL
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ ***REMOVED***
		sa.raw.Path[i] = int8(name[i])
	***REMOVED***
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
***REMOVED***

func (sa *SockaddrDatalink) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Index == 0 ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Len = sa.Len
	sa.raw.Family = AF_LINK
	sa.raw.Index = sa.Index
	sa.raw.Type = sa.Type
	sa.raw.Nlen = sa.Nlen
	sa.raw.Alen = sa.Alen
	sa.raw.Slen = sa.Slen
	for i := 0; i < len(sa.raw.Data); i++ ***REMOVED***
		sa.raw.Data[i] = sa.Data[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), SizeofSockaddrDatalink, nil
***REMOVED***

func anyToSockaddr(fd int, rsa *RawSockaddrAny) (Sockaddr, error) ***REMOVED***
	switch rsa.Addr.Family ***REMOVED***
	case AF_LINK:
		pp := (*RawSockaddrDatalink)(unsafe.Pointer(rsa))
		sa := new(SockaddrDatalink)
		sa.Len = pp.Len
		sa.Family = pp.Family
		sa.Index = pp.Index
		sa.Type = pp.Type
		sa.Nlen = pp.Nlen
		sa.Alen = pp.Alen
		sa.Slen = pp.Slen
		for i := 0; i < len(sa.Data); i++ ***REMOVED***
			sa.Data[i] = pp.Data[i]
		***REMOVED***
		return sa, nil

	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		if pp.Len < 2 || pp.Len > SizeofSockaddrUnix ***REMOVED***
			return nil, EINVAL
		***REMOVED***
		sa := new(SockaddrUnix)

		// Some BSDs include the trailing NUL in the length, whereas
		// others do not. Work around this by subtracting the leading
		// family and len. The path is then scanned to see if a NUL
		// terminator still exists within the length.
		n := int(pp.Len) - 2 // subtract leading Family, Len
		for i := 0; i < n; i++ ***REMOVED***
			if pp.Path[i] == 0 ***REMOVED***
				// found early NUL; assume Len included the NUL
				// or was overestimating.
				n = i
				break
			***REMOVED***
		***REMOVED***
		bytes := (*[len(pp.Path)]byte)(unsafe.Pointer(&pp.Path[0]))[0:n]
		sa.Name = string(bytes)
		return sa, nil

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		for i := 0; i < len(sa.Addr); i++ ***REMOVED***
			sa.Addr[i] = pp.Addr[i]
		***REMOVED***
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		for i := 0; i < len(sa.Addr); i++ ***REMOVED***
			sa.Addr[i] = pp.Addr[i]
		***REMOVED***
		return sa, nil
	***REMOVED***
	return anyToSockaddrGOOS(fd, rsa)
***REMOVED***

func Accept(fd int) (nfd int, sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept(fd, &rsa, &len)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if (runtime.GOOS == "darwin" || runtime.GOOS == "ios") && len == 0 ***REMOVED***
		// Accepted socket has no address.
		// This is likely due to a bug in xnu kernels,
		// where instead of ECONNABORTED error socket
		// is accepted, but has no address.
		Close(nfd)
		return 0, nil, ECONNABORTED
	***REMOVED***
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil ***REMOVED***
		Close(nfd)
		nfd = 0
	***REMOVED***
	return
***REMOVED***

func Getsockname(fd int) (sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil ***REMOVED***
		return
	***REMOVED***
	// TODO(jsing): DragonFly has a "bug" (see issue 3349), which should be
	// reported upstream.
	if runtime.GOOS == "dragonfly" && rsa.Addr.Family == AF_UNSPEC && rsa.Addr.Len == 0 ***REMOVED***
		rsa.Addr.Family = AF_UNIX
		rsa.Addr.Len = SizeofSockaddrUnix
	***REMOVED***
	return anyToSockaddr(fd, &rsa)
***REMOVED***

//sysnb	socketpair(domain int, typ int, proto int, fd *[2]int32) (err error)

// GetsockoptString returns the string value of the socket option opt for the
// socket associated with fd at the given socket level.
func GetsockoptString(fd, level, opt int) (string, error) ***REMOVED***
	buf := make([]byte, 256)
	vallen := _Socklen(len(buf))
	err := getsockopt(fd, level, opt, unsafe.Pointer(&buf[0]), &vallen)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(buf[:vallen-1]), nil
***REMOVED***

//sys	recvfrom(fd int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, err error)
//sys	sendto(s int, buf []byte, flags int, to unsafe.Pointer, addrlen _Socklen) (err error)
//sys	recvmsg(s int, msg *Msghdr, flags int) (n int, err error)

func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error) ***REMOVED***
	var msg Msghdr
	var rsa RawSockaddrAny
	msg.Name = (*byte)(unsafe.Pointer(&rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var iov Iovec
	if len(p) > 0 ***REMOVED***
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	***REMOVED***
	var dummy byte
	if len(oob) > 0 ***REMOVED***
		// receive at least one normal byte
		if len(p) == 0 ***REMOVED***
			iov.Base = &dummy
			iov.SetLen(1)
		***REMOVED***
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
	***REMOVED***
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); err != nil ***REMOVED***
		return
	***REMOVED***
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
	// source address is only specified if the socket is unconnected
	if rsa.Addr.Family != AF_UNSPEC ***REMOVED***
		from, err = anyToSockaddr(fd, &rsa)
	***REMOVED***
	return
***REMOVED***

//sys	sendmsg(s int, msg *Msghdr, flags int) (n int, err error)

func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error) ***REMOVED***
	_, err = SendmsgN(fd, p, oob, to, flags)
	return
***REMOVED***

func SendmsgN(fd int, p, oob []byte, to Sockaddr, flags int) (n int, err error) ***REMOVED***
	var ptr unsafe.Pointer
	var salen _Socklen
	if to != nil ***REMOVED***
		ptr, salen, err = to.sockaddr()
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(ptr))
	msg.Namelen = uint32(salen)
	var iov Iovec
	if len(p) > 0 ***REMOVED***
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	***REMOVED***
	var dummy byte
	if len(oob) > 0 ***REMOVED***
		// send at least one normal byte
		if len(p) == 0 ***REMOVED***
			iov.Base = &dummy
			iov.SetLen(1)
		***REMOVED***
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
	***REMOVED***
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = sendmsg(fd, &msg, flags); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if len(oob) > 0 && len(p) == 0 ***REMOVED***
		n = 0
	***REMOVED***
	return n, nil
***REMOVED***

//sys	kevent(kq int, change unsafe.Pointer, nchange int, event unsafe.Pointer, nevent int, timeout *Timespec) (n int, err error)

func Kevent(kq int, changes, events []Kevent_t, timeout *Timespec) (n int, err error) ***REMOVED***
	var change, event unsafe.Pointer
	if len(changes) > 0 ***REMOVED***
		change = unsafe.Pointer(&changes[0])
	***REMOVED***
	if len(events) > 0 ***REMOVED***
		event = unsafe.Pointer(&events[0])
	***REMOVED***
	return kevent(kq, change, len(changes), event, len(events), timeout)
***REMOVED***

// sysctlmib translates name to mib number and appends any additional args.
func sysctlmib(name string, args ...int) ([]_C_int, error) ***REMOVED***
	// Translate name to mib number.
	mib, err := nametomib(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for _, a := range args ***REMOVED***
		mib = append(mib, _C_int(a))
	***REMOVED***

	return mib, nil
***REMOVED***

func Sysctl(name string) (string, error) ***REMOVED***
	return SysctlArgs(name)
***REMOVED***

func SysctlArgs(name string, args ...int) (string, error) ***REMOVED***
	buf, err := SysctlRaw(name, args...)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	n := len(buf)

	// Throw away terminating NUL.
	if n > 0 && buf[n-1] == '\x00' ***REMOVED***
		n--
	***REMOVED***
	return string(buf[0:n]), nil
***REMOVED***

func SysctlUint32(name string) (uint32, error) ***REMOVED***
	return SysctlUint32Args(name)
***REMOVED***

func SysctlUint32Args(name string, args ...int) (uint32, error) ***REMOVED***
	mib, err := sysctlmib(name, args...)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	n := uintptr(4)
	buf := make([]byte, 4)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if n != 4 ***REMOVED***
		return 0, EIO
	***REMOVED***
	return *(*uint32)(unsafe.Pointer(&buf[0])), nil
***REMOVED***

func SysctlUint64(name string, args ...int) (uint64, error) ***REMOVED***
	mib, err := sysctlmib(name, args...)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	n := uintptr(8)
	buf := make([]byte, 8)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if n != 8 ***REMOVED***
		return 0, EIO
	***REMOVED***
	return *(*uint64)(unsafe.Pointer(&buf[0])), nil
***REMOVED***

func SysctlRaw(name string, args ...int) ([]byte, error) ***REMOVED***
	mib, err := sysctlmib(name, args...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Find size.
	n := uintptr(0)
	if err := sysctl(mib, nil, &n, nil, 0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if n == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	// Read into buffer of that size.
	buf := make([]byte, n)
	if err := sysctl(mib, &buf[0], &n, nil, 0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// The actual call may return less than the original reported required
	// size so ensure we deal with that.
	return buf[:n], nil
***REMOVED***

func SysctlClockinfo(name string) (*Clockinfo, error) ***REMOVED***
	mib, err := sysctlmib(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	n := uintptr(SizeofClockinfo)
	var ci Clockinfo
	if err := sysctl(mib, (*byte)(unsafe.Pointer(&ci)), &n, nil, 0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if n != SizeofClockinfo ***REMOVED***
		return nil, EIO
	***REMOVED***
	return &ci, nil
***REMOVED***

func SysctlTimeval(name string) (*Timeval, error) ***REMOVED***
	mib, err := sysctlmib(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var tv Timeval
	n := uintptr(unsafe.Sizeof(tv))
	if err := sysctl(mib, (*byte)(unsafe.Pointer(&tv)), &n, nil, 0); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if n != unsafe.Sizeof(tv) ***REMOVED***
		return nil, EIO
	***REMOVED***
	return &tv, nil
***REMOVED***

//sys	utimes(path string, timeval *[2]Timeval) (err error)

func Utimes(path string, tv []Timeval) error ***REMOVED***
	if tv == nil ***REMOVED***
		return utimes(path, nil)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

func UtimesNano(path string, ts []Timespec) error ***REMOVED***
	if ts == nil ***REMOVED***
		err := utimensat(AT_FDCWD, path, nil, 0)
		if err != ENOSYS ***REMOVED***
			return err
		***REMOVED***
		return utimes(path, nil)
	***REMOVED***
	if len(ts) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	// Darwin setattrlist can set nanosecond timestamps
	err := setattrlistTimes(path, ts, 0)
	if err != ENOSYS ***REMOVED***
		return err
	***REMOVED***
	err = utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
	if err != ENOSYS ***REMOVED***
		return err
	***REMOVED***
	// Not as efficient as it could be because Timespec and
	// Timeval have different types in the different OSes
	tv := [2]Timeval***REMOVED***
		NsecToTimeval(TimespecToNsec(ts[0])),
		NsecToTimeval(TimespecToNsec(ts[1])),
	***REMOVED***
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

func UtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) error ***REMOVED***
	if ts == nil ***REMOVED***
		return utimensat(dirfd, path, nil, flags)
	***REMOVED***
	if len(ts) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	err := setattrlistTimes(path, ts, flags)
	if err != ENOSYS ***REMOVED***
		return err
	***REMOVED***
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), flags)
***REMOVED***

//sys	futimes(fd int, timeval *[2]Timeval) (err error)

func Futimes(fd int, tv []Timeval) error ***REMOVED***
	if tv == nil ***REMOVED***
		return futimes(fd, nil)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return futimes(fd, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

//sys	poll(fds *PollFd, nfds int, timeout int) (n int, err error)

func Poll(fds []PollFd, timeout int) (n int, err error) ***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return poll(nil, 0, timeout)
	***REMOVED***
	return poll(&fds[0], len(fds), timeout)
***REMOVED***

// TODO: wrap
//	Acct(name nil-string) (err error)
//	Gethostuuid(uuid *byte, timeout *Timespec) (err error)
//	Ptrace(req int, pid int, addr uintptr, data int) (ret uintptr, err error)

var mapper = &mmapper***REMOVED***
	active: make(map[*byte][]byte),
	mmap:   mmap,
	munmap: munmap,
***REMOVED***

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) ***REMOVED***
	return mapper.Mmap(fd, offset, length, prot, flags)
***REMOVED***

func Munmap(b []byte) (err error) ***REMOVED***
	return mapper.Munmap(b)
***REMOVED***

//sys	Madvise(b []byte, behav int) (err error)
//sys	Mlock(b []byte) (err error)
//sys	Mlockall(flags int) (err error)
//sys	Mprotect(b []byte, prot int) (err error)
//sys	Msync(b []byte, flags int) (err error)
//sys	Munlock(b []byte) (err error)
//sys	Munlockall() (err error)
