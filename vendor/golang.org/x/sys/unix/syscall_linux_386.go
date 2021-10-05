// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build 386 && linux
// +build 386,linux

package unix

import (
	"unsafe"
)

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: int32(sec), Nsec: int32(nsec)***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: int32(sec), Usec: int32(usec)***REMOVED***
***REMOVED***

//sysnb	pipe(p *[2]_C_int) (err error)

func Pipe(p []int) (err error) ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	var pp [2]_C_int
	err = pipe(&pp)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
***REMOVED***

//sysnb	pipe2(p *[2]_C_int, flags int) (err error)

func Pipe2(p []int, flags int) (err error) ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	var pp [2]_C_int
	err = pipe2(&pp, flags)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
***REMOVED***

// 64-bit file system and 32-bit uid calls
// (386 default is 32-bit file system and 16-bit uid).
//sys	dup2(oldfd int, newfd int) (err error)
//sysnb	EpollCreate(size int) (fd int, err error)
//sys	EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error)
//sys	Fadvise(fd int, offset int64, length int64, advice int) (err error) = SYS_FADVISE64_64
//sys	Fchown(fd int, uid int, gid int) (err error) = SYS_FCHOWN32
//sys	Fstat(fd int, stat *Stat_t) (err error) = SYS_FSTAT64
//sys	Fstatat(dirfd int, path string, stat *Stat_t, flags int) (err error) = SYS_FSTATAT64
//sys	Ftruncate(fd int, length int64) (err error) = SYS_FTRUNCATE64
//sysnb	Getegid() (egid int) = SYS_GETEGID32
//sysnb	Geteuid() (euid int) = SYS_GETEUID32
//sysnb	Getgid() (gid int) = SYS_GETGID32
//sysnb	Getuid() (uid int) = SYS_GETUID32
//sysnb	InotifyInit() (fd int, err error)
//sys	Ioperm(from int, num int, on int) (err error)
//sys	Iopl(level int) (err error)
//sys	Lchown(path string, uid int, gid int) (err error) = SYS_LCHOWN32
//sys	Lstat(path string, stat *Stat_t) (err error) = SYS_LSTAT64
//sys	Pread(fd int, p []byte, offset int64) (n int, err error) = SYS_PREAD64
//sys	Pwrite(fd int, p []byte, offset int64) (n int, err error) = SYS_PWRITE64
//sys	Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error)
//sys	sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) = SYS_SENDFILE64
//sys	setfsgid(gid int) (prev int, err error) = SYS_SETFSGID32
//sys	setfsuid(uid int) (prev int, err error) = SYS_SETFSUID32
//sysnb	Setregid(rgid int, egid int) (err error) = SYS_SETREGID32
//sysnb	Setresgid(rgid int, egid int, sgid int) (err error) = SYS_SETRESGID32
//sysnb	Setresuid(ruid int, euid int, suid int) (err error) = SYS_SETRESUID32
//sysnb	Setreuid(ruid int, euid int) (err error) = SYS_SETREUID32
//sys	Splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int, err error)
//sys	Stat(path string, stat *Stat_t) (err error) = SYS_STAT64
//sys	SyncFileRange(fd int, off int64, n int64, flags int) (err error)
//sys	Truncate(path string, length int64) (err error) = SYS_TRUNCATE64
//sys	Ustat(dev int, ubuf *Ustat_t) (err error)
//sysnb	getgroups(n int, list *_Gid_t) (nn int, err error) = SYS_GETGROUPS32
//sysnb	setgroups(n int, list *_Gid_t) (err error) = SYS_SETGROUPS32
//sys	Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error) = SYS__NEWSELECT

//sys	mmap2(addr uintptr, length uintptr, prot int, flags int, fd int, pageOffset uintptr) (xaddr uintptr, err error)
//sys	Pause() (err error)

func mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) ***REMOVED***
	page := uintptr(offset / 4096)
	if offset != int64(page)*4096 ***REMOVED***
		return 0, EINVAL
	***REMOVED***
	return mmap2(addr, length, prot, flags, fd, page)
***REMOVED***

type rlimit32 struct ***REMOVED***
	Cur uint32
	Max uint32
***REMOVED***

//sysnb	getrlimit(resource int, rlim *rlimit32) (err error) = SYS_GETRLIMIT

const rlimInf32 = ^uint32(0)
const rlimInf64 = ^uint64(0)

func Getrlimit(resource int, rlim *Rlimit) (err error) ***REMOVED***
	err = prlimit(0, resource, nil, rlim)
	if err != ENOSYS ***REMOVED***
		return err
	***REMOVED***

	rl := rlimit32***REMOVED******REMOVED***
	err = getrlimit(resource, &rl)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if rl.Cur == rlimInf32 ***REMOVED***
		rlim.Cur = rlimInf64
	***REMOVED*** else ***REMOVED***
		rlim.Cur = uint64(rl.Cur)
	***REMOVED***

	if rl.Max == rlimInf32 ***REMOVED***
		rlim.Max = rlimInf64
	***REMOVED*** else ***REMOVED***
		rlim.Max = uint64(rl.Max)
	***REMOVED***
	return
***REMOVED***

//sysnb	setrlimit(resource int, rlim *rlimit32) (err error) = SYS_SETRLIMIT

func Setrlimit(resource int, rlim *Rlimit) (err error) ***REMOVED***
	err = prlimit(0, resource, rlim, nil)
	if err != ENOSYS ***REMOVED***
		return err
	***REMOVED***

	rl := rlimit32***REMOVED******REMOVED***
	if rlim.Cur == rlimInf64 ***REMOVED***
		rl.Cur = rlimInf32
	***REMOVED*** else if rlim.Cur < uint64(rlimInf32) ***REMOVED***
		rl.Cur = uint32(rlim.Cur)
	***REMOVED*** else ***REMOVED***
		return EINVAL
	***REMOVED***
	if rlim.Max == rlimInf64 ***REMOVED***
		rl.Max = rlimInf32
	***REMOVED*** else if rlim.Max < uint64(rlimInf32) ***REMOVED***
		rl.Max = uint32(rlim.Max)
	***REMOVED*** else ***REMOVED***
		return EINVAL
	***REMOVED***

	return setrlimit(resource, &rl)
***REMOVED***

func Seek(fd int, offset int64, whence int) (newoffset int64, err error) ***REMOVED***
	newoffset, errno := seek(fd, offset, whence)
	if errno != 0 ***REMOVED***
		return 0, errno
	***REMOVED***
	return newoffset, nil
***REMOVED***

//sys	futimesat(dirfd int, path string, times *[2]Timeval) (err error)
//sysnb	Gettimeofday(tv *Timeval) (err error)
//sysnb	Time(t *Time_t) (tt Time_t, err error)
//sys	Utime(path string, buf *Utimbuf) (err error)
//sys	utimes(path string, times *[2]Timeval) (err error)

// On x86 Linux, all the socket calls go through an extra indirection,
// I think because the 5-register system call interface can't handle
// the 6-argument calls like sendto and recvfrom. Instead the
// arguments to the underlying system call are the number below
// and a pointer to an array of uintptr. We hide the pointer in the
// socketcall assembly to avoid allocation on every system call.

const (
	// see linux/net.h
	_SOCKET      = 1
	_BIND        = 2
	_CONNECT     = 3
	_LISTEN      = 4
	_ACCEPT      = 5
	_GETSOCKNAME = 6
	_GETPEERNAME = 7
	_SOCKETPAIR  = 8
	_SEND        = 9
	_RECV        = 10
	_SENDTO      = 11
	_RECVFROM    = 12
	_SHUTDOWN    = 13
	_SETSOCKOPT  = 14
	_GETSOCKOPT  = 15
	_SENDMSG     = 16
	_RECVMSG     = 17
	_ACCEPT4     = 18
	_RECVMMSG    = 19
	_SENDMMSG    = 20
)

func accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error) ***REMOVED***
	fd, e := socketcall(_ACCEPT, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)), 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error) ***REMOVED***
	fd, e := socketcall(_ACCEPT4, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)), uintptr(flags), 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func getsockname(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) ***REMOVED***
	_, e := rawsocketcall(_GETSOCKNAME, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)), 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func getpeername(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error) ***REMOVED***
	_, e := rawsocketcall(_GETPEERNAME, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)), 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func socketpair(domain int, typ int, flags int, fd *[2]int32) (err error) ***REMOVED***
	_, e := rawsocketcall(_SOCKETPAIR, uintptr(domain), uintptr(typ), uintptr(flags), uintptr(unsafe.Pointer(fd)), 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) ***REMOVED***
	_, e := socketcall(_BIND, uintptr(s), uintptr(addr), uintptr(addrlen), 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func connect(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) ***REMOVED***
	_, e := socketcall(_CONNECT, uintptr(s), uintptr(addr), uintptr(addrlen), 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func socket(domain int, typ int, proto int) (fd int, err error) ***REMOVED***
	fd, e := rawsocketcall(_SOCKET, uintptr(domain), uintptr(typ), uintptr(proto), 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error) ***REMOVED***
	_, e := socketcall(_GETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(unsafe.Pointer(vallen)), 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error) ***REMOVED***
	_, e := socketcall(_SETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), vallen, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func recvfrom(s int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, err error) ***REMOVED***
	var base uintptr
	if len(p) > 0 ***REMOVED***
		base = uintptr(unsafe.Pointer(&p[0]))
	***REMOVED***
	n, e := socketcall(_RECVFROM, uintptr(s), base, uintptr(len(p)), uintptr(flags), uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(fromlen)))
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func sendto(s int, p []byte, flags int, to unsafe.Pointer, addrlen _Socklen) (err error) ***REMOVED***
	var base uintptr
	if len(p) > 0 ***REMOVED***
		base = uintptr(unsafe.Pointer(&p[0]))
	***REMOVED***
	_, e := socketcall(_SENDTO, uintptr(s), base, uintptr(len(p)), uintptr(flags), uintptr(to), uintptr(addrlen))
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func recvmsg(s int, msg *Msghdr, flags int) (n int, err error) ***REMOVED***
	n, e := socketcall(_RECVMSG, uintptr(s), uintptr(unsafe.Pointer(msg)), uintptr(flags), 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func sendmsg(s int, msg *Msghdr, flags int) (n int, err error) ***REMOVED***
	n, e := socketcall(_SENDMSG, uintptr(s), uintptr(unsafe.Pointer(msg)), uintptr(flags), 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func Listen(s int, n int) (err error) ***REMOVED***
	_, e := socketcall(_LISTEN, uintptr(s), uintptr(n), 0, 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func Shutdown(s, how int) (err error) ***REMOVED***
	_, e := socketcall(_SHUTDOWN, uintptr(s), uintptr(how), 0, 0, 0, 0)
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func Fstatfs(fd int, buf *Statfs_t) (err error) ***REMOVED***
	_, _, e := Syscall(SYS_FSTATFS64, uintptr(fd), unsafe.Sizeof(*buf), uintptr(unsafe.Pointer(buf)))
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func Statfs(path string, buf *Statfs_t) (err error) ***REMOVED***
	pathp, err := BytePtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, _, e := Syscall(SYS_STATFS64, uintptr(unsafe.Pointer(pathp)), unsafe.Sizeof(*buf), uintptr(unsafe.Pointer(buf)))
	if e != 0 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func (r *PtraceRegs) PC() uint64 ***REMOVED*** return uint64(uint32(r.Eip)) ***REMOVED***

func (r *PtraceRegs) SetPC(pc uint64) ***REMOVED*** r.Eip = int32(pc) ***REMOVED***

func (iov *Iovec) SetLen(length int) ***REMOVED***
	iov.Len = uint32(length)
***REMOVED***

func (msghdr *Msghdr) SetControllen(length int) ***REMOVED***
	msghdr.Controllen = uint32(length)
***REMOVED***

func (msghdr *Msghdr) SetIovlen(length int) ***REMOVED***
	msghdr.Iovlen = uint32(length)
***REMOVED***

func (cmsg *Cmsghdr) SetLen(length int) ***REMOVED***
	cmsg.Len = uint32(length)
***REMOVED***

//sys	poll(fds *PollFd, nfds int, timeout int) (n int, err error)

func Poll(fds []PollFd, timeout int) (n int, err error) ***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return poll(nil, 0, timeout)
	***REMOVED***
	return poll(&fds[0], len(fds), timeout)
***REMOVED***
