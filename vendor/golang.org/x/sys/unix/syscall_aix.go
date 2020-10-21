// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build aix

// Aix system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and
// wrap it in our own nicer implementation.

package unix

import "unsafe"

/*
 * Wrapped
 */

//sys	utimes(path string, times *[2]Timeval) (err error)
func Utimes(path string, tv []Timeval) error ***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

//sys	utimensat(dirfd int, path string, times *[2]Timespec, flag int) (err error)
func UtimesNano(path string, ts []Timespec) error ***REMOVED***
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
	if n > len(sa.raw.Path) ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	if n == len(sa.raw.Path) && name[0] != '@' ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ ***REMOVED***
		sa.raw.Path[i] = uint8(name[i])
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

func Getsockname(fd int) (sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil ***REMOVED***
		return
	***REMOVED***
	return anyToSockaddr(fd, &rsa)
***REMOVED***

//sys	getcwd(buf []byte) (err error)

const ImplementsGetwd = true

func Getwd() (ret string, err error) ***REMOVED***
	for len := uint64(4096); ; len *= 2 ***REMOVED***
		b := make([]byte, len)
		err := getcwd(b)
		if err == nil ***REMOVED***
			i := 0
			for b[i] != 0 ***REMOVED***
				i++
			***REMOVED***
			return string(b[0:i]), nil
		***REMOVED***
		if err != ERANGE ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
***REMOVED***

func Getcwd(buf []byte) (n int, err error) ***REMOVED***
	err = getcwd(buf)
	if err == nil ***REMOVED***
		i := 0
		for buf[i] != 0 ***REMOVED***
			i++
		***REMOVED***
		n = i + 1
	***REMOVED***
	return
***REMOVED***

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

/*
 * Socket
 */

//sys	accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error)

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

func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error) ***REMOVED***
	// Recvmsg not implemented on AIX
	sa := new(SockaddrUnix)
	return -1, -1, -1, sa, ENOSYS
***REMOVED***

func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error) ***REMOVED***
	_, err = SendmsgN(fd, p, oob, to, flags)
	return
***REMOVED***

func SendmsgN(fd int, p, oob []byte, to Sockaddr, flags int) (n int, err error) ***REMOVED***
	// SendmsgN not implemented on AIX
	return -1, ENOSYS
***REMOVED***

func anyToSockaddr(fd int, rsa *RawSockaddrAny) (Sockaddr, error) ***REMOVED***
	switch rsa.Addr.Family ***REMOVED***

	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)

		// Some versions of AIX have a bug in getsockname (see IV78655).
		// We can't rely on sa.Len being set correctly.
		n := SizeofSockaddrUnix - 3 // subtract leading Family, Len, terminating NUL.
		for i := 0; i < n; i++ ***REMOVED***
			if pp.Path[i] == 0 ***REMOVED***
				n = i
				break
			***REMOVED***
		***REMOVED***

		bytes := (*[10000]byte)(unsafe.Pointer(&pp.Path[0]))[0:n]
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

func Gettimeofday(tv *Timeval) (err error) ***REMOVED***
	err = gettimeofday(tv, nil)
	return
***REMOVED***

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) ***REMOVED***
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	return sendfile(outfd, infd, offset, count)
***REMOVED***

// TODO
func sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) ***REMOVED***
	return -1, ENOSYS
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

//sys	getdirent(fd int, buf []byte) (n int, err error)
func Getdents(fd int, buf []byte) (n int, err error) ***REMOVED***
	return getdirent(fd, buf)
***REMOVED***

//sys	wait4(pid Pid_t, status *_C_int, options int, rusage *Rusage) (wpid Pid_t, err error)
func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) ***REMOVED***
	var status _C_int
	var r Pid_t
	err = ERESTART
	// AIX wait4 may return with ERESTART errno, while the processus is still
	// active.
	for err == ERESTART ***REMOVED***
		r, err = wait4(Pid_t(pid), &status, options, rusage)
	***REMOVED***
	wpid = int(r)
	if wstatus != nil ***REMOVED***
		*wstatus = WaitStatus(status)
	***REMOVED***
	return
***REMOVED***

/*
 * Wait
 */

type WaitStatus uint32

func (w WaitStatus) Stopped() bool ***REMOVED*** return w&0x40 != 0 ***REMOVED***
func (w WaitStatus) StopSignal() Signal ***REMOVED***
	if !w.Stopped() ***REMOVED***
		return -1
	***REMOVED***
	return Signal(w>>8) & 0xFF
***REMOVED***

func (w WaitStatus) Exited() bool ***REMOVED*** return w&0xFF == 0 ***REMOVED***
func (w WaitStatus) ExitStatus() int ***REMOVED***
	if !w.Exited() ***REMOVED***
		return -1
	***REMOVED***
	return int((w >> 8) & 0xFF)
***REMOVED***

func (w WaitStatus) Signaled() bool ***REMOVED*** return w&0x40 == 0 && w&0xFF != 0 ***REMOVED***
func (w WaitStatus) Signal() Signal ***REMOVED***
	if !w.Signaled() ***REMOVED***
		return -1
	***REMOVED***
	return Signal(w>>16) & 0xFF
***REMOVED***

func (w WaitStatus) Continued() bool ***REMOVED*** return w&0x01000000 != 0 ***REMOVED***

func (w WaitStatus) CoreDump() bool ***REMOVED*** return w&0x80 == 0x80 ***REMOVED***

func (w WaitStatus) TrapCause() int ***REMOVED*** return -1 ***REMOVED***

//sys	ioctl(fd int, req uint, arg uintptr) (err error)

// fcntl must never be called with cmd=F_DUP2FD because it doesn't work on AIX
// There is no way to create a custom fcntl and to keep //sys fcntl easily,
// Therefore, the programmer must call dup2 instead of fcntl in this case.

// FcntlInt performs a fcntl syscall on fd with the provided command and argument.
//sys	FcntlInt(fd uintptr, cmd int, arg int) (r int,err error) = fcntl

// FcntlFlock performs a fcntl syscall for the F_GETLK, F_SETLK or F_SETLKW command.
//sys	FcntlFlock(fd uintptr, cmd int, lk *Flock_t) (err error) = fcntl

//sys	fcntl(fd int, cmd int, arg int) (val int, err error)

/*
 * Direct access
 */

//sys	Acct(path string) (err error)
//sys	Chdir(path string) (err error)
//sys	Chroot(path string) (err error)
//sys	Close(fd int) (err error)
//sys	Dup(oldfd int) (fd int, err error)
//sys	Exit(code int)
//sys	Faccessat(dirfd int, path string, mode uint32, flags int) (err error)
//sys	Fchdir(fd int) (err error)
//sys	Fchmod(fd int, mode uint32) (err error)
//sys	Fchmodat(dirfd int, path string, mode uint32, flags int) (err error)
//sys	Fchownat(dirfd int, path string, uid int, gid int, flags int) (err error)
//sys	Fdatasync(fd int) (err error)
//sys	Fsync(fd int) (err error)
// readdir_r
//sysnb	Getpgid(pid int) (pgid int, err error)

//sys	Getpgrp() (pid int)

//sysnb	Getpid() (pid int)
//sysnb	Getppid() (ppid int)
//sys	Getpriority(which int, who int) (prio int, err error)
//sysnb	Getrusage(who int, rusage *Rusage) (err error)
//sysnb	Getsid(pid int) (sid int, err error)
//sysnb	Kill(pid int, sig Signal) (err error)
//sys	Klogctl(typ int, buf []byte) (n int, err error) = syslog
//sys	Mkdir(dirfd int, path string, mode uint32) (err error)
//sys	Mkdirat(dirfd int, path string, mode uint32) (err error)
//sys	Mkfifo(path string, mode uint32) (err error)
//sys	Mknod(path string, mode uint32, dev int) (err error)
//sys	Mknodat(dirfd int, path string, mode uint32, dev int) (err error)
//sys	Nanosleep(time *Timespec, leftover *Timespec) (err error)
//sys   Open(path string, mode int, perm uint32) (fd int, err error) = open64
//sys   Openat(dirfd int, path string, flags int, mode uint32) (fd int, err error)
//sys	read(fd int, p []byte) (n int, err error)
//sys	Readlink(path string, buf []byte) (n int, err error)
//sys	Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error)
//sys	Setdomainname(p []byte) (err error)
//sys	Sethostname(p []byte) (err error)
//sysnb	Setpgid(pid int, pgid int) (err error)
//sysnb	Setsid() (pid int, err error)
//sysnb	Settimeofday(tv *Timeval) (err error)

//sys	Setuid(uid int) (err error)
//sys	Setgid(uid int) (err error)

//sys	Setpriority(which int, who int, prio int) (err error)
//sys	Statx(dirfd int, path string, flags int, mask int, stat *Statx_t) (err error)
//sys	Sync()
//sysnb	Times(tms *Tms) (ticks uintptr, err error)
//sysnb	Umask(mask int) (oldmask int)
//sysnb	Uname(buf *Utsname) (err error)
//sys   Unlink(path string) (err error)
//sys   Unlinkat(dirfd int, path string, flags int) (err error)
//sys	Ustat(dev int, ubuf *Ustat_t) (err error)
//sys	write(fd int, p []byte) (n int, err error)
//sys	readlen(fd int, p *byte, np int) (n int, err error) = read
//sys	writelen(fd int, p *byte, np int) (n int, err error) = write

//sys	Dup2(oldfd int, newfd int) (err error)
//sys	Fadvise(fd int, offset int64, length int64, advice int) (err error) = posix_fadvise64
//sys	Fchown(fd int, uid int, gid int) (err error)
//sys	fstat(fd int, stat *Stat_t) (err error)
//sys	fstatat(dirfd int, path string, stat *Stat_t, flags int) (err error) = fstatat
//sys	Fstatfs(fd int, buf *Statfs_t) (err error)
//sys	Ftruncate(fd int, length int64) (err error)
//sysnb	Getegid() (egid int)
//sysnb	Geteuid() (euid int)
//sysnb	Getgid() (gid int)
//sysnb	Getuid() (uid int)
//sys	Lchown(path string, uid int, gid int) (err error)
//sys	Listen(s int, n int) (err error)
//sys	lstat(path string, stat *Stat_t) (err error)
//sys	Pause() (err error)
//sys	Pread(fd int, p []byte, offset int64) (n int, err error) = pread64
//sys	Pwrite(fd int, p []byte, offset int64) (n int, err error) = pwrite64
//sys	Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error)
//sys	Pselect(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timespec, sigmask *Sigset_t) (n int, err error)
//sysnb	Setregid(rgid int, egid int) (err error)
//sysnb	Setreuid(ruid int, euid int) (err error)
//sys	Shutdown(fd int, how int) (err error)
//sys	Splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int64, err error)
//sys	stat(path string, statptr *Stat_t) (err error)
//sys	Statfs(path string, buf *Statfs_t) (err error)
//sys	Truncate(path string, length int64) (err error)

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

// In order to use msghdr structure with Control, Controllen, nrecvmsg and nsendmsg must be used.
//sys	recvmsg(s int, msg *Msghdr, flags int) (n int, err error) = nrecvmsg
//sys	sendmsg(s int, msg *Msghdr, flags int) (n int, err error) = nsendmsg

//sys	munmap(addr uintptr, length uintptr) (err error)

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

//sys	Madvise(b []byte, advice int) (err error)
//sys	Mprotect(b []byte, prot int) (err error)
//sys	Mlock(b []byte) (err error)
//sys	Mlockall(flags int) (err error)
//sys	Msync(b []byte, flags int) (err error)
//sys	Munlock(b []byte) (err error)
//sys	Munlockall() (err error)

//sysnb pipe(p *[2]_C_int) (err error)

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

//sys	poll(fds *PollFd, nfds int, timeout int) (n int, err error)

func Poll(fds []PollFd, timeout int) (n int, err error) ***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return poll(nil, 0, timeout)
	***REMOVED***
	return poll(&fds[0], len(fds), timeout)
***REMOVED***

//sys	gettimeofday(tv *Timeval, tzp *Timezone) (err error)
//sysnb	Time(t *Time_t) (tt Time_t, err error)
//sys	Utime(path string, buf *Utimbuf) (err error)

//sys	Getsystemcfg(label int) (n uint64)

//sys	umount(target string) (err error)
func Unmount(target string, flags int) (err error) ***REMOVED***
	if flags != 0 ***REMOVED***
		// AIX doesn't have any flags for umount.
		return ENOSYS
	***REMOVED***
	return umount(target)
***REMOVED***
