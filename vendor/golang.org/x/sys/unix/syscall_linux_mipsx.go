// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && (mips || mipsle)
// +build linux
// +build mips mipsle

package unix

import (
	"syscall"
	"unsafe"
)

func Syscall9(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r1, r2 uintptr, err syscall.Errno)

//sys	dup2(oldfd int, newfd int) (err error)
//sysnb	EpollCreate(size int) (fd int, err error)
//sys	EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error)
//sys	Fadvise(fd int, offset int64, length int64, advice int) (err error) = SYS_FADVISE64
//sys	Fchown(fd int, uid int, gid int) (err error)
//sys	Ftruncate(fd int, length int64) (err error) = SYS_FTRUNCATE64
//sysnb	Getegid() (egid int)
//sysnb	Geteuid() (euid int)
//sysnb	Getgid() (gid int)
//sysnb	Getuid() (uid int)
//sys	Lchown(path string, uid int, gid int) (err error)
//sys	Listen(s int, n int) (err error)
//sys	Pread(fd int, p []byte, offset int64) (n int, err error) = SYS_PREAD64
//sys	Pwrite(fd int, p []byte, offset int64) (n int, err error) = SYS_PWRITE64
//sys	Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error)
//sys	Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error) = SYS__NEWSELECT
//sys	sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) = SYS_SENDFILE64
//sys	setfsgid(gid int) (prev int, err error)
//sys	setfsuid(uid int) (prev int, err error)
//sysnb	Setregid(rgid int, egid int) (err error)
//sysnb	Setresgid(rgid int, egid int, sgid int) (err error)
//sysnb	Setresuid(ruid int, euid int, suid int) (err error)
//sysnb	Setreuid(ruid int, euid int) (err error)
//sys	Shutdown(fd int, how int) (err error)
//sys	Splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int, err error)
//sys	SyncFileRange(fd int, off int64, n int64, flags int) (err error)
//sys	Truncate(path string, length int64) (err error) = SYS_TRUNCATE64
//sys	Ustat(dev int, ubuf *Ustat_t) (err error)
//sys	accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error)
//sys	accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error)
//sys	bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error)
//sys	connect(s int, addr unsafe.Pointer, addrlen _Socklen) (err error)
//sysnb	getgroups(n int, list *_Gid_t) (nn int, err error)
//sysnb	setgroups(n int, list *_Gid_t) (err error)
//sys	getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error)
//sys	setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error)
//sysnb	socket(domain int, typ int, proto int) (fd int, err error)
//sysnb	socketpair(domain int, typ int, proto int, fd *[2]int32) (err error)
//sysnb	getpeername(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error)
//sysnb	getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error)
//sys	recvfrom(fd int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, err error)
//sys	sendto(s int, buf []byte, flags int, to unsafe.Pointer, addrlen _Socklen) (err error)
//sys	recvmsg(s int, msg *Msghdr, flags int) (n int, err error)
//sys	sendmsg(s int, msg *Msghdr, flags int) (n int, err error)

//sysnb	InotifyInit() (fd int, err error)
//sys	Ioperm(from int, num int, on int) (err error)
//sys	Iopl(level int) (err error)

//sys	futimesat(dirfd int, path string, times *[2]Timeval) (err error)
//sysnb	Gettimeofday(tv *Timeval) (err error)
//sysnb	Time(t *Time_t) (tt Time_t, err error)
//sys	Utime(path string, buf *Utimbuf) (err error)
//sys	utimes(path string, times *[2]Timeval) (err error)

//sys	Lstat(path string, stat *Stat_t) (err error) = SYS_LSTAT64
//sys	Fstat(fd int, stat *Stat_t) (err error) = SYS_FSTAT64
//sys	Fstatat(dirfd int, path string, stat *Stat_t, flags int) (err error) = SYS_FSTATAT64
//sys	Stat(path string, stat *Stat_t) (err error) = SYS_STAT64

//sys	Pause() (err error)

func Fstatfs(fd int, buf *Statfs_t) (err error) ***REMOVED***
	_, _, e := Syscall(SYS_FSTATFS64, uintptr(fd), unsafe.Sizeof(*buf), uintptr(unsafe.Pointer(buf)))
	if e != 0 ***REMOVED***
		err = errnoErr(e)
	***REMOVED***
	return
***REMOVED***

func Statfs(path string, buf *Statfs_t) (err error) ***REMOVED***
	p, err := BytePtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	_, _, e := Syscall(SYS_STATFS64, uintptr(unsafe.Pointer(p)), unsafe.Sizeof(*buf), uintptr(unsafe.Pointer(buf)))
	if e != 0 ***REMOVED***
		err = errnoErr(e)
	***REMOVED***
	return
***REMOVED***

func Seek(fd int, offset int64, whence int) (off int64, err error) ***REMOVED***
	_, _, e := Syscall6(SYS__LLSEEK, uintptr(fd), uintptr(offset>>32), uintptr(offset), uintptr(unsafe.Pointer(&off)), uintptr(whence), 0)
	if e != 0 ***REMOVED***
		err = errnoErr(e)
	***REMOVED***
	return
***REMOVED***

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: int32(sec), Nsec: int32(nsec)***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: int32(sec), Usec: int32(usec)***REMOVED***
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

//sysnb	pipe() (p1 int, p2 int, err error)

func Pipe(p []int) (err error) ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	p[0], p[1], err = pipe()
	return
***REMOVED***

//sys	mmap2(addr uintptr, length uintptr, prot int, flags int, fd int, pageOffset uintptr) (xaddr uintptr, err error)

func mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) ***REMOVED***
	page := uintptr(offset / 4096)
	if offset != int64(page)*4096 ***REMOVED***
		return 0, EINVAL
	***REMOVED***
	return mmap2(addr, length, prot, flags, fd, page)
***REMOVED***

const rlimInf32 = ^uint32(0)
const rlimInf64 = ^uint64(0)

type rlimit32 struct ***REMOVED***
	Cur uint32
	Max uint32
***REMOVED***

//sysnb	getrlimit(resource int, rlim *rlimit32) (err error) = SYS_GETRLIMIT

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

func (r *PtraceRegs) PC() uint64 ***REMOVED*** return r.Epc ***REMOVED***

func (r *PtraceRegs) SetPC(pc uint64) ***REMOVED*** r.Epc = pc ***REMOVED***

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
