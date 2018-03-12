// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build arm64,linux

package unix

//sys	EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error) = SYS_EPOLL_PWAIT
//sys	Fchown(fd int, uid int, gid int) (err error)
//sys	Fstat(fd int, stat *Stat_t) (err error)
//sys	Fstatat(fd int, path string, stat *Stat_t, flags int) (err error)
//sys	Fstatfs(fd int, buf *Statfs_t) (err error)
//sys	Ftruncate(fd int, length int64) (err error)
//sysnb	Getegid() (egid int)
//sysnb	Geteuid() (euid int)
//sysnb	Getgid() (gid int)
//sysnb	Getrlimit(resource int, rlim *Rlimit) (err error)
//sysnb	Getuid() (uid int)
//sys	Listen(s int, n int) (err error)
//sys	Pread(fd int, p []byte, offset int64) (n int, err error) = SYS_PREAD64
//sys	Pwrite(fd int, p []byte, offset int64) (n int, err error) = SYS_PWRITE64
//sys	Seek(fd int, offset int64, whence int) (off int64, err error) = SYS_LSEEK

func Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error) ***REMOVED***
	var ts *Timespec
	if timeout != nil ***REMOVED***
		ts = &Timespec***REMOVED***Sec: timeout.Sec, Nsec: timeout.Usec * 1000***REMOVED***
	***REMOVED***
	return Pselect(nfd, r, w, e, ts, nil)
***REMOVED***

//sys	sendfile(outfd int, infd int, offset *int64, count int) (written int, err error)
//sys	Setfsgid(gid int) (err error)
//sys	Setfsuid(uid int) (err error)
//sysnb	Setregid(rgid int, egid int) (err error)
//sysnb	Setresgid(rgid int, egid int, sgid int) (err error)
//sysnb	Setresuid(ruid int, euid int, suid int) (err error)
//sysnb	Setrlimit(resource int, rlim *Rlimit) (err error)
//sysnb	Setreuid(ruid int, euid int) (err error)
//sys	Shutdown(fd int, how int) (err error)
//sys	Splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int64, err error)

func Stat(path string, stat *Stat_t) (err error) ***REMOVED***
	return Fstatat(AT_FDCWD, path, stat, 0)
***REMOVED***

func Lchown(path string, uid int, gid int) (err error) ***REMOVED***
	return Fchownat(AT_FDCWD, path, uid, gid, AT_SYMLINK_NOFOLLOW)
***REMOVED***

func Lstat(path string, stat *Stat_t) (err error) ***REMOVED***
	return Fstatat(AT_FDCWD, path, stat, AT_SYMLINK_NOFOLLOW)
***REMOVED***

//sys	Statfs(path string, buf *Statfs_t) (err error)
//sys	SyncFileRange(fd int, off int64, n int64, flags int) (err error)
//sys	Truncate(path string, length int64) (err error)
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
//sys	mmap(addr uintptr, length uintptr, prot int, flags int, fd int, offset int64) (xaddr uintptr, err error)

//sysnb	Gettimeofday(tv *Timeval) (err error)

func setTimespec(sec, nsec int64) Timespec ***REMOVED***
	return Timespec***REMOVED***Sec: sec, Nsec: nsec***REMOVED***
***REMOVED***

func setTimeval(sec, usec int64) Timeval ***REMOVED***
	return Timeval***REMOVED***Sec: sec, Usec: usec***REMOVED***
***REMOVED***

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

func Pipe(p []int) (err error) ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	var pp [2]_C_int
	err = pipe2(&pp, 0)
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

func (r *PtraceRegs) PC() uint64 ***REMOVED*** return r.Pc ***REMOVED***

func (r *PtraceRegs) SetPC(pc uint64) ***REMOVED*** r.Pc = pc ***REMOVED***

func (iov *Iovec) SetLen(length int) ***REMOVED***
	iov.Len = uint64(length)
***REMOVED***

func (msghdr *Msghdr) SetControllen(length int) ***REMOVED***
	msghdr.Controllen = uint64(length)
***REMOVED***

func (cmsg *Cmsghdr) SetLen(length int) ***REMOVED***
	cmsg.Len = uint64(length)
***REMOVED***

func InotifyInit() (fd int, err error) ***REMOVED***
	return InotifyInit1(0)
***REMOVED***

func Dup2(oldfd int, newfd int) (err error) ***REMOVED***
	return Dup3(oldfd, newfd, 0)
***REMOVED***

func Pause() (err error) ***REMOVED***
	_, _, e1 := Syscall6(SYS_PPOLL, 0, 0, 0, 0, 0, 0)
	if e1 != 0 ***REMOVED***
		err = errnoErr(e1)
	***REMOVED***
	return
***REMOVED***

// TODO(dfc): constants that should be in zsysnum_linux_arm64.go, remove
// these when the deprecated syscalls that the syscall package relies on
// are removed.
const (
	SYS_GETPGRP      = 1060
	SYS_UTIMES       = 1037
	SYS_FUTIMESAT    = 1066
	SYS_PAUSE        = 1061
	SYS_USTAT        = 1070
	SYS_UTIME        = 1063
	SYS_LCHOWN       = 1032
	SYS_TIME         = 1062
	SYS_EPOLL_CREATE = 1042
	SYS_EPOLL_WAIT   = 1069
)

func Poll(fds []PollFd, timeout int) (n int, err error) ***REMOVED***
	var ts *Timespec
	if timeout >= 0 ***REMOVED***
		ts = new(Timespec)
		*ts = NsecToTimespec(int64(timeout) * 1e6)
	***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return ppoll(nil, 0, ts, nil)
	***REMOVED***
	return ppoll(&fds[0], len(fds), ts, nil)
***REMOVED***
