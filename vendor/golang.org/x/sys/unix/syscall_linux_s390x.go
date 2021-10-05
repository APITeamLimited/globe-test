// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build s390x && linux
// +build s390x,linux

package unix

import (
	"unsafe"
)

//sys	dup2(oldfd int, newfd int) (err error)
//sysnb	EpollCreate(size int) (fd int, err error)
//sys	EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error)
//sys	Fadvise(fd int, offset int64, length int64, advice int) (err error) = SYS_FADVISE64
//sys	Fchown(fd int, uid int, gid int) (err error)
//sys	Fstat(fd int, stat *Stat_t) (err error)
//sys	Fstatat(dirfd int, path string, stat *Stat_t, flags int) (err error) = SYS_NEWFSTATAT
//sys	Fstatfs(fd int, buf *Statfs_t) (err error)
//sys	Ftruncate(fd int, length int64) (err error)
//sysnb	Getegid() (egid int)
//sysnb	Geteuid() (euid int)
//sysnb	Getgid() (gid int)
//sysnb	Getrlimit(resource int, rlim *Rlimit) (err error)
//sysnb	Getuid() (uid int)
//sysnb	InotifyInit() (fd int, err error)
//sys	Lchown(path string, uid int, gid int) (err error)
//sys	Lstat(path string, stat *Stat_t) (err error)
//sys	Pause() (err error)
//sys	Pread(fd int, p []byte, offset int64) (n int, err error) = SYS_PREAD64
//sys	Pwrite(fd int, p []byte, offset int64) (n int, err error) = SYS_PWRITE64
//sys	Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error)
//sys	Seek(fd int, offset int64, whence int) (off int64, err error) = SYS_LSEEK
//sys	Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error)
//sys	sendfile(outfd int, infd int, offset *int64, count int) (written int, err error)
//sys	setfsgid(gid int) (prev int, err error)
//sys	setfsuid(uid int) (prev int, err error)
//sysnb	Setregid(rgid int, egid int) (err error)
//sysnb	Setresgid(rgid int, egid int, sgid int) (err error)
//sysnb	Setresuid(ruid int, euid int, suid int) (err error)
//sysnb	Setrlimit(resource int, rlim *Rlimit) (err error)
//sysnb	Setreuid(ruid int, euid int) (err error)
//sys	Splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int64, err error)
//sys	Stat(path string, stat *Stat_t) (err error)
//sys	Statfs(path string, buf *Statfs_t) (err error)
//sys	SyncFileRange(fd int, off int64, n int64, flags int) (err error)
//sys	Truncate(path string, length int64) (err error)
//sys	Ustat(dev int, ubuf *Ustat_t) (err error)
//sysnb	getgroups(n int, list *_Gid_t) (nn int, err error)
//sysnb	setgroups(n int, list *_Gid_t) (err error)

//sys	futimesat(dirfd int, path string, times *[2]Timeval) (err error)
//sysnb	Gettimeofday(tv *Timeval) (err error)

func Time(t *Time_t) (tt Time_t, err error) ***REMOVED***
	var tv Timeval
	err = Gettimeofday(&tv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if t != nil ***REMOVED***
		*t = Time_t(tv.Sec)
	***REMOVED***
	return Time_t(tv.Sec), nil
***REMOVED***

//sys	Utime(path string, buf *Utimbuf) (err error)
//sys	utimes(path string, times *[2]Timeval) (err error)

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: sec, Nsec: nsec***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: sec, Usec: usec***REMOVED***
***REMOVED***

//sysnb	pipe2(p *[2]_C_int, flags int) (err error)

func Pipe(p []int) (err error) ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	var pp [2]_C_int
	err = pipe2(&pp, 0) // pipe2 is the same as pipe when flags are set to 0.
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
***REMOVED***

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

func Ioperm(from int, num int, on int) (err error) ***REMOVED***
	return ENOSYS
***REMOVED***

func Iopl(level int) (err error) ***REMOVED***
	return ENOSYS
***REMOVED***

func (r *PtraceRegs) PC() uint64 ***REMOVED*** return r.Psw.Addr ***REMOVED***

func (r *PtraceRegs) SetPC(pc uint64) ***REMOVED*** r.Psw.Addr = pc ***REMOVED***

func (iov *Iovec) SetLen(length int) ***REMOVED***
	iov.Len = uint64(length)
***REMOVED***

func (msghdr *Msghdr) SetControllen(length int) ***REMOVED***
	msghdr.Controllen = uint64(length)
***REMOVED***

func (msghdr *Msghdr) SetIovlen(length int) ***REMOVED***
	msghdr.Iovlen = uint64(length)
***REMOVED***

func (cmsg *Cmsghdr) SetLen(length int) ***REMOVED***
	cmsg.Len = uint64(length)
***REMOVED***

// Linux on s390x uses the old mmap interface, which requires arguments to be passed in a struct.
// mmap2 also requires arguments to be passed in a struct; it is currently not exposed in <asm/unistd.h>.
func mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error) ***REMOVED***
	mmap_args := [6]uintptr***REMOVED***addr, length, uintptr(prot), uintptr(flags), uintptr(fd), uintptr(offset)***REMOVED***
	r0, _, e1 := Syscall(SYS_MMAP, uintptr(unsafe.Pointer(&mmap_args[0])), 0, 0)
	xaddr = uintptr(r0)
	if e1 != 0 ***REMOVED***
		err = errnoErr(e1)
	***REMOVED***
	return
***REMOVED***

// On s390x Linux, all the socket calls go through an extra indirection.
// The arguments to the underlying system call (SYS_SOCKETCALL) are the
// number below and a pointer to an array of uintptr.
const (
	// see linux/net.h
	netSocket      = 1
	netBind        = 2
	netConnect     = 3
	netListen      = 4
	netAccept      = 5
	netGetSockName = 6
	netGetPeerName = 7
	netSocketPair  = 8
	netSend        = 9
	netRecv        = 10
	netSendTo      = 11
	netRecvFrom    = 12
	netShutdown    = 13
	netSetSockOpt  = 14
	netGetSockOpt  = 15
	netSendMsg     = 16
	netRecvMsg     = 17
	netAccept4     = 18
	netRecvMMsg    = 19
	netSendMMsg    = 20
)

func accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (int, error) ***REMOVED***
	args := [3]uintptr***REMOVED***uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen))***REMOVED***
	fd, _, err := Syscall(SYS_SOCKETCALL, netAccept, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return 0, err
	***REMOVED***
	return int(fd), nil
***REMOVED***

func accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (int, error) ***REMOVED***
	args := [4]uintptr***REMOVED***uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)), uintptr(flags)***REMOVED***
	fd, _, err := Syscall(SYS_SOCKETCALL, netAccept4, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return 0, err
	***REMOVED***
	return int(fd), nil
***REMOVED***

func getsockname(s int, rsa *RawSockaddrAny, addrlen *_Socklen) error ***REMOVED***
	args := [3]uintptr***REMOVED***uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen))***REMOVED***
	_, _, err := RawSyscall(SYS_SOCKETCALL, netGetSockName, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func getpeername(s int, rsa *RawSockaddrAny, addrlen *_Socklen) error ***REMOVED***
	args := [3]uintptr***REMOVED***uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen))***REMOVED***
	_, _, err := RawSyscall(SYS_SOCKETCALL, netGetPeerName, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func socketpair(domain int, typ int, flags int, fd *[2]int32) error ***REMOVED***
	args := [4]uintptr***REMOVED***uintptr(domain), uintptr(typ), uintptr(flags), uintptr(unsafe.Pointer(fd))***REMOVED***
	_, _, err := RawSyscall(SYS_SOCKETCALL, netSocketPair, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func bind(s int, addr unsafe.Pointer, addrlen _Socklen) error ***REMOVED***
	args := [3]uintptr***REMOVED***uintptr(s), uintptr(addr), uintptr(addrlen)***REMOVED***
	_, _, err := Syscall(SYS_SOCKETCALL, netBind, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func connect(s int, addr unsafe.Pointer, addrlen _Socklen) error ***REMOVED***
	args := [3]uintptr***REMOVED***uintptr(s), uintptr(addr), uintptr(addrlen)***REMOVED***
	_, _, err := Syscall(SYS_SOCKETCALL, netConnect, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func socket(domain int, typ int, proto int) (int, error) ***REMOVED***
	args := [3]uintptr***REMOVED***uintptr(domain), uintptr(typ), uintptr(proto)***REMOVED***
	fd, _, err := RawSyscall(SYS_SOCKETCALL, netSocket, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return 0, err
	***REMOVED***
	return int(fd), nil
***REMOVED***

func getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) error ***REMOVED***
	args := [5]uintptr***REMOVED***uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(unsafe.Pointer(vallen))***REMOVED***
	_, _, err := Syscall(SYS_SOCKETCALL, netGetSockOpt, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) error ***REMOVED***
	args := [5]uintptr***REMOVED***uintptr(s), uintptr(level), uintptr(name), uintptr(val), vallen***REMOVED***
	_, _, err := Syscall(SYS_SOCKETCALL, netSetSockOpt, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func recvfrom(s int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (int, error) ***REMOVED***
	var base uintptr
	if len(p) > 0 ***REMOVED***
		base = uintptr(unsafe.Pointer(&p[0]))
	***REMOVED***
	args := [6]uintptr***REMOVED***uintptr(s), base, uintptr(len(p)), uintptr(flags), uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(fromlen))***REMOVED***
	n, _, err := Syscall(SYS_SOCKETCALL, netRecvFrom, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return 0, err
	***REMOVED***
	return int(n), nil
***REMOVED***

func sendto(s int, p []byte, flags int, to unsafe.Pointer, addrlen _Socklen) error ***REMOVED***
	var base uintptr
	if len(p) > 0 ***REMOVED***
		base = uintptr(unsafe.Pointer(&p[0]))
	***REMOVED***
	args := [6]uintptr***REMOVED***uintptr(s), base, uintptr(len(p)), uintptr(flags), uintptr(to), uintptr(addrlen)***REMOVED***
	_, _, err := Syscall(SYS_SOCKETCALL, netSendTo, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func recvmsg(s int, msg *Msghdr, flags int) (int, error) ***REMOVED***
	args := [3]uintptr***REMOVED***uintptr(s), uintptr(unsafe.Pointer(msg)), uintptr(flags)***REMOVED***
	n, _, err := Syscall(SYS_SOCKETCALL, netRecvMsg, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return 0, err
	***REMOVED***
	return int(n), nil
***REMOVED***

func sendmsg(s int, msg *Msghdr, flags int) (int, error) ***REMOVED***
	args := [3]uintptr***REMOVED***uintptr(s), uintptr(unsafe.Pointer(msg)), uintptr(flags)***REMOVED***
	n, _, err := Syscall(SYS_SOCKETCALL, netSendMsg, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return 0, err
	***REMOVED***
	return int(n), nil
***REMOVED***

func Listen(s int, n int) error ***REMOVED***
	args := [2]uintptr***REMOVED***uintptr(s), uintptr(n)***REMOVED***
	_, _, err := Syscall(SYS_SOCKETCALL, netListen, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func Shutdown(s, how int) error ***REMOVED***
	args := [2]uintptr***REMOVED***uintptr(s), uintptr(how)***REMOVED***
	_, _, err := Syscall(SYS_SOCKETCALL, netShutdown, uintptr(unsafe.Pointer(&args)), 0)
	if err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

//sys	poll(fds *PollFd, nfds int, timeout int) (n int, err error)

func Poll(fds []PollFd, timeout int) (n int, err error) ***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return poll(nil, 0, timeout)
	***REMOVED***
	return poll(&fds[0], len(fds), timeout)
***REMOVED***

//sys	kexecFileLoad(kernelFd int, initrdFd int, cmdlineLen int, cmdline string, flags int) (err error)

func KexecFileLoad(kernelFd int, initrdFd int, cmdline string, flags int) error ***REMOVED***
	cmdlineLen := len(cmdline)
	if cmdlineLen > 0 ***REMOVED***
		// Account for the additional NULL byte added by
		// BytePtrFromString in kexecFileLoad. The kexec_file_load
		// syscall expects a NULL-terminated string.
		cmdlineLen++
	***REMOVED***
	return kexecFileLoad(kernelFd, initrdFd, cmdlineLen, cmdline, flags)
***REMOVED***
