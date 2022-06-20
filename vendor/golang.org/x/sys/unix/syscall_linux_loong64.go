// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build loong64 && linux
// +build loong64,linux

package unix

import "unsafe"

//sys	EpollWait(epfd int, events []EpollEvent, msec int) (n int, err error) = SYS_EPOLL_PWAIT
//sys	Fadvise(fd int, offset int64, length int64, advice int) (err error) = SYS_FADVISE64
//sys	Fchown(fd int, uid int, gid int) (err error)
//sys	Fstat(fd int, stat *Stat_t) (err error)
//sys	Fstatat(fd int, path string, stat *Stat_t, flags int) (err error)
//sys	Fstatfs(fd int, buf *Statfs_t) (err error)
//sys	Ftruncate(fd int, length int64) (err error)
//sysnb	Getegid() (egid int)
//sysnb	Geteuid() (euid int)
//sysnb	Getgid() (gid int)
//sysnb	Getuid() (uid int)
//sys	Listen(s int, n int) (err error)
//sys	pread(fd int, p []byte, offset int64) (n int, err error) = SYS_PREAD64
//sys	pwrite(fd int, p []byte, offset int64) (n int, err error) = SYS_PWRITE64
//sys	Seek(fd int, offset int64, whence int) (off int64, err error) = SYS_LSEEK

func Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error) ***REMOVED***
	var ts *Timespec
	if timeout != nil ***REMOVED***
		ts = &Timespec***REMOVED***Sec: timeout.Sec, Nsec: timeout.Usec * 1000***REMOVED***
	***REMOVED***
	return Pselect(nfd, r, w, e, ts, nil)
***REMOVED***

//sys	sendfile(outfd int, infd int, offset *int64, count int) (written int, err error)
//sys	setfsgid(gid int) (prev int, err error)
//sys	setfsuid(uid int) (prev int, err error)
//sysnb	Setregid(rgid int, egid int) (err error)
//sysnb	Setresgid(rgid int, egid int, sgid int) (err error)
//sysnb	Setresuid(ruid int, euid int, suid int) (err error)
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

func Ustat(dev int, ubuf *Ustat_t) (err error) ***REMOVED***
	return ENOSYS
***REMOVED***

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

func Getrlimit(resource int, rlim *Rlimit) (err error) ***REMOVED***
	err = Prlimit(0, resource, nil, rlim)
	return
***REMOVED***

func Setrlimit(resource int, rlim *Rlimit) (err error) ***REMOVED***
	err = Prlimit(0, resource, rlim, nil)
	return
***REMOVED***

func futimesat(dirfd int, path string, tv *[2]Timeval) (err error) ***REMOVED***
	if tv == nil ***REMOVED***
		return utimensat(dirfd, path, nil, 0)
	***REMOVED***

	ts := []Timespec***REMOVED***
		NsecToTimespec(TimevalToNsec(tv[0])),
		NsecToTimespec(TimevalToNsec(tv[1])),
	***REMOVED***
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
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

func utimes(path string, tv *[2]Timeval) (err error) ***REMOVED***
	if tv == nil ***REMOVED***
		return utimensat(AT_FDCWD, path, nil, 0)
	***REMOVED***

	ts := []Timespec***REMOVED***
		NsecToTimespec(TimevalToNsec(tv[0])),
		NsecToTimespec(TimevalToNsec(tv[1])),
	***REMOVED***
	return utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
***REMOVED***

func (r *PtraceRegs) PC() uint64 ***REMOVED*** return r.Era ***REMOVED***

func (r *PtraceRegs) SetPC(era uint64) ***REMOVED*** r.Era = era ***REMOVED***

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

func (rsa *RawSockaddrNFCLLCP) SetServiceNameLen(length int) ***REMOVED***
	rsa.Service_name_len = uint64(length)
***REMOVED***

func Pause() error ***REMOVED***
	_, err := ppoll(nil, 0, nil, nil)
	return err
***REMOVED***

func Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error) ***REMOVED***
	return Renameat2(olddirfd, oldpath, newdirfd, newpath, 0)
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