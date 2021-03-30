// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Solaris system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and wrap
// it in our own nicer implementation, either here or in
// syscall_solaris.go or syscall_unix.go.

package unix

import (
	"runtime"
	"syscall"
	"unsafe"
)

// Implemented in runtime/syscall_solaris.go.
type syscallFunc uintptr

func rawSysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)
func sysvicall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

// SockaddrDatalink implements the Sockaddr interface for AF_LINK type sockets.
type SockaddrDatalink struct ***REMOVED***
	Family uint16
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [244]int8
	raw    RawSockaddrDatalink
***REMOVED***

func direntIno(buf []byte) (uint64, bool) ***REMOVED***
	return readInt(buf, unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Ino), unsafe.Sizeof(Dirent***REMOVED******REMOVED***.Ino))
***REMOVED***

func direntReclen(buf []byte) (uint64, bool) ***REMOVED***
	return readInt(buf, unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Reclen), unsafe.Sizeof(Dirent***REMOVED******REMOVED***.Reclen))
***REMOVED***

func direntNamlen(buf []byte) (uint64, bool) ***REMOVED***
	reclen, ok := direntReclen(buf)
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	return reclen - uint64(unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Name)), true
***REMOVED***

//sysnb	pipe(p *[2]_C_int) (n int, err error)

func Pipe(p []int) (err error) ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	var pp [2]_C_int
	n, err := pipe(&pp)
	if n != 0 ***REMOVED***
		return err
	***REMOVED***
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return nil
***REMOVED***

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Port < 0 || sa.Port > 0xFFFF ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet4, nil
***REMOVED***

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Port < 0 || sa.Port > 0xFFFF ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet6, nil
***REMOVED***

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ ***REMOVED***
		sa.raw.Path[i] = int8(name[i])
	***REMOVED***
	// length is family (uint16), name, NUL.
	sl := _Socklen(2)
	if n > 0 ***REMOVED***
		sl += _Socklen(n) + 1
	***REMOVED***
	if sa.raw.Path[0] == '@' ***REMOVED***
		sa.raw.Path[0] = 0
		// Don't count trailing NUL for abstract address.
		sl--
	***REMOVED***

	return unsafe.Pointer(&sa.raw), sl, nil
***REMOVED***

//sys	getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) = libsocket.getsockname

func Getsockname(fd int) (sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil ***REMOVED***
		return
	***REMOVED***
	return anyToSockaddr(fd, &rsa)
***REMOVED***

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

const ImplementsGetwd = true

//sys	Getcwd(buf []byte) (n int, err error)

func Getwd() (wd string, err error) ***REMOVED***
	var buf [PathMax]byte
	// Getcwd will return an error if it failed for any reason.
	_, err = Getcwd(buf[0:])
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
	// Check for error and sanity check group count. Newer versions of
	// Solaris allow up to 1024 (NGROUPS_MAX).
	if n < 0 || n > 1024 ***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return nil, EINVAL
	***REMOVED*** else if n == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if n == -1 ***REMOVED***
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

// ReadDirent reads directory entries from fd and writes them into buf.
func ReadDirent(fd int, buf []byte) (n int, err error) ***REMOVED***
	// Final argument is (basep *uintptr) and the syscall doesn't take nil.
	// TODO(rsc): Can we use a single global basep for all calls?
	return Getdents(fd, buf, new(uintptr))
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

func (w WaitStatus) Continued() bool ***REMOVED*** return w&mask == stopped && syscall.Signal(w>>shift) == SIGSTOP ***REMOVED***

func (w WaitStatus) StopSignal() syscall.Signal ***REMOVED***
	if !w.Stopped() ***REMOVED***
		return -1
	***REMOVED***
	return syscall.Signal(w>>shift) & 0xFF
***REMOVED***

func (w WaitStatus) TrapCause() int ***REMOVED*** return -1 ***REMOVED***

//sys	wait4(pid int32, statusp *_C_int, options int, rusage *Rusage) (wpid int32, err error)

func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (int, error) ***REMOVED***
	var status _C_int
	rpid, err := wait4(int32(pid), &status, options, rusage)
	wpid := int(rpid)
	if wpid == -1 ***REMOVED***
		return wpid, err
	***REMOVED***
	if wstatus != nil ***REMOVED***
		*wstatus = WaitStatus(status)
	***REMOVED***
	return wpid, nil
***REMOVED***

//sys	gethostname(buf []byte) (n int, err error)

func Gethostname() (name string, err error) ***REMOVED***
	var buf [MaxHostNameLen]byte
	n, err := gethostname(buf[:])
	if n != 0 ***REMOVED***
		return "", err
	***REMOVED***
	n = clen(buf[:])
	if n < 1 ***REMOVED***
		return "", EFAULT
	***REMOVED***
	return string(buf[:n]), nil
***REMOVED***

//sys	utimes(path string, times *[2]Timeval) (err error)

func Utimes(path string, tv []Timeval) (err error) ***REMOVED***
	if tv == nil ***REMOVED***
		return utimes(path, nil)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

//sys	utimensat(fd int, path string, times *[2]Timespec, flag int) (err error)

func UtimesNano(path string, ts []Timespec) error ***REMOVED***
	if ts == nil ***REMOVED***
		return utimensat(AT_FDCWD, path, nil, 0)
	***REMOVED***
	if len(ts) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
***REMOVED***

func UtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) error ***REMOVED***
	if ts == nil ***REMOVED***
		return utimensat(dirfd, path, nil, flags)
	***REMOVED***
	if len(ts) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), flags)
***REMOVED***

//sys	fcntl(fd int, cmd int, arg int) (val int, err error)

// FcntlInt performs a fcntl syscall on fd with the provided command and argument.
func FcntlInt(fd uintptr, cmd, arg int) (int, error) ***REMOVED***
	valptr, _, errno := sysvicall6(uintptr(unsafe.Pointer(&procfcntl)), 3, uintptr(fd), uintptr(cmd), uintptr(arg), 0, 0, 0)
	var err error
	if errno != 0 ***REMOVED***
		err = errno
	***REMOVED***
	return int(valptr), err
***REMOVED***

// FcntlFlock performs a fcntl syscall for the F_GETLK, F_SETLK or F_SETLKW command.
func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error ***REMOVED***
	_, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procfcntl)), 3, uintptr(fd), uintptr(cmd), uintptr(unsafe.Pointer(lk)), 0, 0, 0)
	if e1 != 0 ***REMOVED***
		return e1
	***REMOVED***
	return nil
***REMOVED***

//sys	futimesat(fildes int, path *byte, times *[2]Timeval) (err error)

func Futimesat(dirfd int, path string, tv []Timeval) error ***REMOVED***
	pathp, err := BytePtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if tv == nil ***REMOVED***
		return futimesat(dirfd, pathp, nil)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return futimesat(dirfd, pathp, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

// Solaris doesn't have an futimes function because it allows NULL to be
// specified as the path for futimesat. However, Go doesn't like
// NULL-style string interfaces, so this simple wrapper is provided.
func Futimes(fd int, tv []Timeval) error ***REMOVED***
	if tv == nil ***REMOVED***
		return futimesat(fd, nil, nil)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return futimesat(fd, nil, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

func anyToSockaddr(fd int, rsa *RawSockaddrAny) (Sockaddr, error) ***REMOVED***
	switch rsa.Addr.Family ***REMOVED***
	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		// Assume path ends at NUL.
		// This is not technically the Solaris semantics for
		// abstract Unix domain sockets -- they are supposed
		// to be uninterpreted fixed-size binary blobs -- but
		// everyone uses this convention.
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 ***REMOVED***
			n++
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
	return nil, EAFNOSUPPORT
***REMOVED***

//sys	accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error) = libsocket.accept

func Accept(fd int) (nfd int, sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept(fd, &rsa, &len)
	if nfd == -1 ***REMOVED***
		return
	***REMOVED***
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil ***REMOVED***
		Close(nfd)
		nfd = 0
	***REMOVED***
	return
***REMOVED***

//sys	recvmsg(s int, msg *Msghdr, flags int) (n int, err error) = libsocket.__xnet_recvmsg

func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error) ***REMOVED***
	var msg Msghdr
	var rsa RawSockaddrAny
	msg.Name = (*byte)(unsafe.Pointer(&rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var iov Iovec
	if len(p) > 0 ***REMOVED***
		iov.Base = (*int8)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	***REMOVED***
	var dummy int8
	if len(oob) > 0 ***REMOVED***
		// receive at least one normal byte
		if len(p) == 0 ***REMOVED***
			iov.Base = &dummy
			iov.SetLen(1)
		***REMOVED***
		msg.Accrightslen = int32(len(oob))
	***REMOVED***
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); n == -1 ***REMOVED***
		return
	***REMOVED***
	oobn = int(msg.Accrightslen)
	// source address is only specified if the socket is unconnected
	if rsa.Addr.Family != AF_UNSPEC ***REMOVED***
		from, err = anyToSockaddr(fd, &rsa)
	***REMOVED***
	return
***REMOVED***

func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error) ***REMOVED***
	_, err = SendmsgN(fd, p, oob, to, flags)
	return
***REMOVED***

//sys	sendmsg(s int, msg *Msghdr, flags int) (n int, err error) = libsocket.__xnet_sendmsg

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
		iov.Base = (*int8)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	***REMOVED***
	var dummy int8
	if len(oob) > 0 ***REMOVED***
		// send at least one normal byte
		if len(p) == 0 ***REMOVED***
			iov.Base = &dummy
			iov.SetLen(1)
		***REMOVED***
		msg.Accrightslen = int32(len(oob))
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

//sys	acct(path *byte) (err error)

func Acct(path string) (err error) ***REMOVED***
	if len(path) == 0 ***REMOVED***
		// Assume caller wants to disable accounting.
		return acct(nil)
	***REMOVED***

	pathp, err := BytePtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return acct(pathp)
***REMOVED***

//sys	__makedev(version int, major uint, minor uint) (val uint64)

func Mkdev(major, minor uint32) uint64 ***REMOVED***
	return __makedev(NEWDEV, uint(major), uint(minor))
***REMOVED***

//sys	__major(version int, dev uint64) (val uint)

func Major(dev uint64) uint32 ***REMOVED***
	return uint32(__major(NEWDEV, dev))
***REMOVED***

//sys	__minor(version int, dev uint64) (val uint)

func Minor(dev uint64) uint32 ***REMOVED***
	return uint32(__minor(NEWDEV, dev))
***REMOVED***

/*
 * Expose the ioctl function
 */

//sys	ioctl(fd int, req uint, arg uintptr) (err error)

func IoctlSetTermio(fd int, req uint, value *Termio) error ***REMOVED***
	err := ioctl(fd, req, uintptr(unsafe.Pointer(value)))
	runtime.KeepAlive(value)
	return err
***REMOVED***

func IoctlGetTermio(fd int, req uint) (*Termio, error) ***REMOVED***
	var value Termio
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
***REMOVED***

//sys   poll(fds *PollFd, nfds int, timeout int) (n int, err error)

func Poll(fds []PollFd, timeout int) (n int, err error) ***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return poll(nil, 0, timeout)
	***REMOVED***
	return poll(&fds[0], len(fds), timeout)
***REMOVED***

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) ***REMOVED***
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	return sendfile(outfd, infd, offset, count)
***REMOVED***

/*
 * Exposed directly
 */
//sys	Access(path string, mode uint32) (err error)
//sys	Adjtime(delta *Timeval, olddelta *Timeval) (err error)
//sys	Chdir(path string) (err error)
//sys	Chmod(path string, mode uint32) (err error)
//sys	Chown(path string, uid int, gid int) (err error)
//sys	Chroot(path string) (err error)
//sys	Close(fd int) (err error)
//sys	Creat(path string, mode uint32) (fd int, err error)
//sys	Dup(fd int) (nfd int, err error)
//sys	Dup2(oldfd int, newfd int) (err error)
//sys	Exit(code int)
//sys	Faccessat(dirfd int, path string, mode uint32, flags int) (err error)
//sys	Fchdir(fd int) (err error)
//sys	Fchmod(fd int, mode uint32) (err error)
//sys	Fchmodat(dirfd int, path string, mode uint32, flags int) (err error)
//sys	Fchown(fd int, uid int, gid int) (err error)
//sys	Fchownat(dirfd int, path string, uid int, gid int, flags int) (err error)
//sys	Fdatasync(fd int) (err error)
//sys	Flock(fd int, how int) (err error)
//sys	Fpathconf(fd int, name int) (val int, err error)
//sys	Fstat(fd int, stat *Stat_t) (err error)
//sys	Fstatat(fd int, path string, stat *Stat_t, flags int) (err error)
//sys	Fstatvfs(fd int, vfsstat *Statvfs_t) (err error)
//sys	Getdents(fd int, buf []byte, basep *uintptr) (n int, err error)
//sysnb	Getgid() (gid int)
//sysnb	Getpid() (pid int)
//sysnb	Getpgid(pid int) (pgid int, err error)
//sysnb	Getpgrp() (pgid int, err error)
//sys	Geteuid() (euid int)
//sys	Getegid() (egid int)
//sys	Getppid() (ppid int)
//sys	Getpriority(which int, who int) (n int, err error)
//sysnb	Getrlimit(which int, lim *Rlimit) (err error)
//sysnb	Getrusage(who int, rusage *Rusage) (err error)
//sysnb	Gettimeofday(tv *Timeval) (err error)
//sysnb	Getuid() (uid int)
//sys	Kill(pid int, signum syscall.Signal) (err error)
//sys	Lchown(path string, uid int, gid int) (err error)
//sys	Link(path string, link string) (err error)
//sys	Listen(s int, backlog int) (err error) = libsocket.__xnet_llisten
//sys	Lstat(path string, stat *Stat_t) (err error)
//sys	Madvise(b []byte, advice int) (err error)
//sys	Mkdir(path string, mode uint32) (err error)
//sys	Mkdirat(dirfd int, path string, mode uint32) (err error)
//sys	Mkfifo(path string, mode uint32) (err error)
//sys	Mkfifoat(dirfd int, path string, mode uint32) (err error)
//sys	Mknod(path string, mode uint32, dev int) (err error)
//sys	Mknodat(dirfd int, path string, mode uint32, dev int) (err error)
//sys	Mlock(b []byte) (err error)
//sys	Mlockall(flags int) (err error)
//sys	Mprotect(b []byte, prot int) (err error)
//sys	Msync(b []byte, flags int) (err error)
//sys	Munlock(b []byte) (err error)
//sys	Munlockall() (err error)
//sys	Nanosleep(time *Timespec, leftover *Timespec) (err error)
//sys	Open(path string, mode int, perm uint32) (fd int, err error)
//sys	Openat(dirfd int, path string, flags int, mode uint32) (fd int, err error)
//sys	Pathconf(path string, name int) (val int, err error)
//sys	Pause() (err error)
//sys	Pread(fd int, p []byte, offset int64) (n int, err error)
//sys	Pwrite(fd int, p []byte, offset int64) (n int, err error)
//sys	read(fd int, p []byte) (n int, err error)
//sys	Readlink(path string, buf []byte) (n int, err error)
//sys	Rename(from string, to string) (err error)
//sys	Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error)
//sys	Rmdir(path string) (err error)
//sys	Seek(fd int, offset int64, whence int) (newoffset int64, err error) = lseek
//sys	Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error)
//sysnb	Setegid(egid int) (err error)
//sysnb	Seteuid(euid int) (err error)
//sysnb	Setgid(gid int) (err error)
//sys	Sethostname(p []byte) (err error)
//sysnb	Setpgid(pid int, pgid int) (err error)
//sys	Setpriority(which int, who int, prio int) (err error)
//sysnb	Setregid(rgid int, egid int) (err error)
//sysnb	Setreuid(ruid int, euid int) (err error)
//sysnb	Setrlimit(which int, lim *Rlimit) (err error)
//sysnb	Setsid() (pid int, err error)
//sysnb	Setuid(uid int) (err error)
//sys	Shutdown(s int, how int) (err error) = libsocket.shutdown
//sys	Stat(path string, stat *Stat_t) (err error)
//sys	Statvfs(path string, vfsstat *Statvfs_t) (err error)
//sys	Symlink(path string, link string) (err error)
//sys	Sync() (err error)
//sysnb	Times(tms *Tms) (ticks uintptr, err error)
//sys	Truncate(path string, length int64) (err error)
//sys	Fsync(fd int) (err error)
//sys	Ftruncate(fd int, length int64) (err error)
//sys	Umask(mask int) (oldmask int)
//sysnb	Uname(buf *Utsname) (err error)
//sys	Unmount(target string, flags int) (err error) = libc.umount
//sys	Unlink(path string) (err error)
//sys	Unlinkat(dirfd int, path string, flags int) (err error)
//sys	Ustat(dev int, ubuf *Ustat_t) (err error)
//sys	Utime(path string, buf *Utimbuf) (err error)
//sys	bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) = libsocket.__xnet_bind
//sys	connect(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) = libsocket.__xnet_connect
//sys	mmap(addr uintptr, length uintptr, prot int, flag int, fd int, pos int64) (ret uintptr, err error)
//sys	munmap(addr uintptr, length uintptr) (err error)
//sys	sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) = libsendfile.sendfile
//sys	sendto(s int, buf []byte, flags int, to unsafe.Pointer, addrlen _Socklen) (err error) = libsocket.__xnet_sendto
//sys	socket(domain int, typ int, proto int) (fd int, err error) = libsocket.__xnet_socket
//sysnb	socketpair(domain int, typ int, proto int, fd *[2]int32) (err error) = libsocket.__xnet_socketpair
//sys	write(fd int, p []byte) (n int, err error)
//sys	getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error) = libsocket.__xnet_getsockopt
//sysnb	getpeername(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) = libsocket.getpeername
//sys	setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error) = libsocket.setsockopt
//sys	recvfrom(fd int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, err error) = libsocket.recvfrom

func readlen(fd int, buf *byte, nbuf int) (n int, err error) ***REMOVED***
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procread)), 3, uintptr(fd), uintptr(unsafe.Pointer(buf)), uintptr(nbuf), 0, 0, 0)
	n = int(r0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***

func writelen(fd int, buf *byte, nbuf int) (n int, err error) ***REMOVED***
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procwrite)), 3, uintptr(fd), uintptr(unsafe.Pointer(buf)), uintptr(nbuf), 0, 0, 0)
	n = int(r0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***

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
