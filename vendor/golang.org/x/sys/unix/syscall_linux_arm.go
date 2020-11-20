// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build arm,linux

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
	// Try pipe2 first for Android O, then try pipe for kernel 2.6.23.
	err = pipe2(&pp, 0)
	if err == ENOSYS ***REMOVED***
		err = pipe(&pp)
	***REMOVED***
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
***REMOVED***

//sysnb pipe2(p *[2]_C_int, flags int) (err error)

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

func Seek(fd int, offset int64, whence int) (newoffset int64, err error) ***REMOVED***
	newoffset, errno := seek(fd, offset, whence)
	if errno != 0 ***REMOVED***
		return 0, errno
	***REMOVED***
	return newoffset, nil
***REMOVED***

//sys	accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error)
//sys	accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error)
//sys	bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error)
//sys	connect(s int, addr unsafe.Pointer, addrlen _Socklen) (err error)
//sysnb	getgroups(n int, list *_Gid_t) (nn int, err error) = SYS_GETGROUPS32
//sysnb	setgroups(n int, list *_Gid_t) (err error) = SYS_SETGROUPS32
//sys	getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error)
//sys	setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error)
//sysnb	socket(domain int, typ int, proto int) (fd int, err error)
//sysnb	getpeername(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error)
//sysnb	getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error)
//sys	recvfrom(fd int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, err error)
//sys	sendto(s int, buf []byte, flags int, to unsafe.Pointer, addrlen _Socklen) (err error)
//sysnb	socketpair(domain int, typ int, flags int, fd *[2]int32) (err error)
//sys	recvmsg(s int, msg *Msghdr, flags int) (n int, err error)
//sys	sendmsg(s int, msg *Msghdr, flags int) (n int, err error)

// 64-bit file system and 32-bit uid calls
// (16-bit uid calls are not always supported in newer kernels)
//sys	dup2(oldfd int, newfd int) (err error)
//sysnb	EpollCreate(size int) (fd int, err error)
//sys	EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error)
//sys	Fchown(fd int, uid int, gid int) (err error) = SYS_FCHOWN32
//sys	Fstat(fd int, stat *Stat_t) (err error) = SYS_FSTAT64
//sys	Fstatat(dirfd int, path string, stat *Stat_t, flags int) (err error) = SYS_FSTATAT64
//sysnb	Getegid() (egid int) = SYS_GETEGID32
//sysnb	Geteuid() (euid int) = SYS_GETEUID32
//sysnb	Getgid() (gid int) = SYS_GETGID32
//sysnb	Getuid() (uid int) = SYS_GETUID32
//sysnb	InotifyInit() (fd int, err error)
//sys	Lchown(path string, uid int, gid int) (err error) = SYS_LCHOWN32
//sys	Listen(s int, n int) (err error)
//sys	Lstat(path string, stat *Stat_t) (err error) = SYS_LSTAT64
//sys	Pause() (err error)
//sys	Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error)
//sys	sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) = SYS_SENDFILE64
//sys	Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error) = SYS__NEWSELECT
//sys	setfsgid(gid int) (prev int, err error) = SYS_SETFSGID32
//sys	setfsuid(uid int) (prev int, err error) = SYS_SETFSUID32
//sysnb	Setregid(rgid int, egid int) (err error) = SYS_SETREGID32
//sysnb	Setresgid(rgid int, egid int, sgid int) (err error) = SYS_SETRESGID32
//sysnb	Setresuid(ruid int, euid int, suid int) (err error) = SYS_SETRESUID32
//sysnb	Setreuid(ruid int, euid int) (err error) = SYS_SETREUID32
//sys	Shutdown(fd int, how int) (err error)
//sys	Splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int, err error)
//sys	Stat(path string, stat *Stat_t) (err error) = SYS_STAT64
//sys	Ustat(dev int, ubuf *Ustat_t) (err error)

//sys	futimesat(dirfd int, path string, times *[2]Timeval) (err error)
//sysnb	Gettimeofday(tv *Timeval) (err error)

func Time(t *Time_t) (Time_t, error) ***REMOVED***
	var tv Timeval
	err := Gettimeofday(&tv)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if t != nil ***REMOVED***
		*t = Time_t(tv.Sec)
	***REMOVED***
	return Time_t(tv.Sec), nil
***REMOVED***

func Utime(path string, buf *Utimbuf) error ***REMOVED***
	tv := []Timeval***REMOVED***
		***REMOVED***Sec: buf.Actime***REMOVED***,
		***REMOVED***Sec: buf.Modtime***REMOVED***,
	***REMOVED***
	return Utimes(path, tv)
***REMOVED***

//sys	utimes(path string, times *[2]Timeval) (err error)

//sys   Pread(fd int, p []byte, offset int64) (n int, err error) = SYS_PREAD64
//sys   Pwrite(fd int, p []byte, offset int64) (n int, err error) = SYS_PWRITE64
//sys	Truncate(path string, length int64) (err error) = SYS_TRUNCATE64
//sys	Ftruncate(fd int, length int64) (err error) = SYS_FTRUNCATE64

func Fadvise(fd int, offset int64, length int64, advice int) (err error) ***REMOVED***
	_, _, e1 := Syscall6(SYS_ARM_FADVISE64_64, uintptr(fd), uintptr(advice), uintptr(offset), uintptr(offset>>32), uintptr(length), uintptr(length>>32))
	if e1 != 0 ***REMOVED***
		err = errnoErr(e1)
	***REMOVED***
	return
***REMOVED***

//sys	mmap2(addr uintptr, length uintptr, prot int, flags int, fd int, pageOffset uintptr) (xaddr uintptr, err error)

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

//sysnb getrlimit(resource int, rlim *rlimit32) (err error) = SYS_UGETRLIMIT

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

//sysnb setrlimit(resource int, rlim *rlimit32) (err error) = SYS_SETRLIMIT

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

func (r *PtraceRegs) PC() uint64 ***REMOVED*** return uint64(r.Uregs[15]) ***REMOVED***

func (r *PtraceRegs) SetPC(pc uint64) ***REMOVED*** r.Uregs[15] = uint32(pc) ***REMOVED***

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

//sys	armSyncFileRange(fd int, flags int, off int64, n int64) (err error) = SYS_ARM_SYNC_FILE_RANGE

func SyncFileRange(fd int, off int64, n int64, flags int) error ***REMOVED***
	// The sync_file_range and arm_sync_file_range syscalls differ only in the
	// order of their arguments.
	return armSyncFileRange(fd, flags, off, n)
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
