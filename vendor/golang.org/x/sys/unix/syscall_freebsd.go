// Copyright 2009,2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// FreeBSD system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and wrap
// it in our own nicer implementation, either here or in
// syscall_bsd.go or syscall_unix.go.

package unix

import (
	"sync"
	"unsafe"
)

const (
	SYS_FSTAT_FREEBSD12         = 551 // ***REMOVED*** int fstat(int fd, _Out_ struct stat *sb); ***REMOVED***
	SYS_FSTATAT_FREEBSD12       = 552 // ***REMOVED*** int fstatat(int fd, _In_z_ char *path, \
	SYS_GETDIRENTRIES_FREEBSD12 = 554 // ***REMOVED*** ssize_t getdirentries(int fd, \
	SYS_STATFS_FREEBSD12        = 555 // ***REMOVED*** int statfs(_In_z_ char *path, \
	SYS_FSTATFS_FREEBSD12       = 556 // ***REMOVED*** int fstatfs(int fd, \
	SYS_GETFSSTAT_FREEBSD12     = 557 // ***REMOVED*** int getfsstat( \
	SYS_MKNODAT_FREEBSD12       = 559 // ***REMOVED*** int mknodat(int fd, _In_z_ char *path, \
)

// See https://www.freebsd.org/doc/en_US.ISO8859-1/books/porters-handbook/versions.html.
var (
	osreldateOnce sync.Once
	osreldate     uint32
)

// INO64_FIRST from /usr/src/lib/libc/sys/compat-ino64.h
const _ino64First = 1200031

func supportsABI(ver uint32) bool ***REMOVED***
	osreldateOnce.Do(func() ***REMOVED*** osreldate, _ = SysctlUint32("kern.osreldate") ***REMOVED***)
	return osreldate >= ver
***REMOVED***

// SockaddrDatalink implements the Sockaddr interface for AF_LINK type sockets.
type SockaddrDatalink struct ***REMOVED***
	Len    uint8
	Family uint8
	Index  uint16
	Type   uint8
	Nlen   uint8
	Alen   uint8
	Slen   uint8
	Data   [46]int8
	raw    RawSockaddrDatalink
***REMOVED***

// Translate "kern.hostname" to []_C_int***REMOVED***0,1,2,3***REMOVED***.
func nametomib(name string) (mib []_C_int, err error) ***REMOVED***
	const siz = unsafe.Sizeof(mib[0])

	// NOTE(rsc): It seems strange to set the buffer to have
	// size CTL_MAXNAME+2 but use only CTL_MAXNAME
	// as the size. I don't know why the +2 is here, but the
	// kernel uses +2 for its own implementation of this function.
	// I am scared that if we don't include the +2 here, the kernel
	// will silently write 2 words farther than we specify
	// and we'll get memory corruption.
	var buf [CTL_MAXNAME + 2]_C_int
	n := uintptr(CTL_MAXNAME) * siz

	p := (*byte)(unsafe.Pointer(&buf[0]))
	bytes, err := ByteSliceFromString(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Magic sysctl: "setting" 0.3 to a string name
	// lets you read back the array of integers form.
	if err = sysctl([]_C_int***REMOVED***0, 3***REMOVED***, p, &n, &bytes[0], uintptr(len(name))); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return buf[0 : n/siz], nil
***REMOVED***

func direntIno(buf []byte) (uint64, bool) ***REMOVED***
	return readInt(buf, unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Fileno), unsafe.Sizeof(Dirent***REMOVED******REMOVED***.Fileno))
***REMOVED***

func direntReclen(buf []byte) (uint64, bool) ***REMOVED***
	return readInt(buf, unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Reclen), unsafe.Sizeof(Dirent***REMOVED******REMOVED***.Reclen))
***REMOVED***

func direntNamlen(buf []byte) (uint64, bool) ***REMOVED***
	return readInt(buf, unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Namlen), unsafe.Sizeof(Dirent***REMOVED******REMOVED***.Namlen))
***REMOVED***

func Pipe(p []int) (err error) ***REMOVED***
	return Pipe2(p, 0)
***REMOVED***

//sysnb	pipe2(p *[2]_C_int, flags int) (err error)

func Pipe2(p []int, flags int) error ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	var pp [2]_C_int
	err := pipe2(&pp, flags)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return err
***REMOVED***

func GetsockoptIPMreqn(fd, level, opt int) (*IPMreqn, error) ***REMOVED***
	var value IPMreqn
	vallen := _Socklen(SizeofIPMreqn)
	errno := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, errno
***REMOVED***

func SetsockoptIPMreqn(fd, level, opt int, mreq *IPMreqn) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), unsafe.Sizeof(*mreq))
***REMOVED***

func Accept4(fd, flags int) (nfd int, sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept4(fd, &rsa, &len, flags)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if len > SizeofSockaddrAny ***REMOVED***
		panic("RawSockaddrAny too small")
	***REMOVED***
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil ***REMOVED***
		Close(nfd)
		nfd = 0
	***REMOVED***
	return
***REMOVED***

const ImplementsGetwd = true

//sys	Getcwd(buf []byte) (n int, err error) = SYS___GETCWD

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

func Getfsstat(buf []Statfs_t, flags int) (n int, err error) ***REMOVED***
	var (
		_p0          unsafe.Pointer
		bufsize      uintptr
		oldBuf       []statfs_freebsd11_t
		needsConvert bool
	)

	if len(buf) > 0 ***REMOVED***
		if supportsABI(_ino64First) ***REMOVED***
			_p0 = unsafe.Pointer(&buf[0])
			bufsize = unsafe.Sizeof(Statfs_t***REMOVED******REMOVED***) * uintptr(len(buf))
		***REMOVED*** else ***REMOVED***
			n := len(buf)
			oldBuf = make([]statfs_freebsd11_t, n)
			_p0 = unsafe.Pointer(&oldBuf[0])
			bufsize = unsafe.Sizeof(statfs_freebsd11_t***REMOVED******REMOVED***) * uintptr(n)
			needsConvert = true
		***REMOVED***
	***REMOVED***
	var sysno uintptr = SYS_GETFSSTAT
	if supportsABI(_ino64First) ***REMOVED***
		sysno = SYS_GETFSSTAT_FREEBSD12
	***REMOVED***
	r0, _, e1 := Syscall(sysno, uintptr(_p0), bufsize, uintptr(flags))
	n = int(r0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	if e1 == 0 && needsConvert ***REMOVED***
		for i := range oldBuf ***REMOVED***
			buf[i].convertFrom(&oldBuf[i])
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func setattrlistTimes(path string, times []Timespec, flags int) error ***REMOVED***
	// used on Darwin for UtimesNano
	return ENOSYS
***REMOVED***

//sys   ioctl(fd int, req uint, arg uintptr) (err error)

//sys   sysctl(mib []_C_int, old *byte, oldlen *uintptr, new *byte, newlen uintptr) (err error) = SYS___SYSCTL

func Uname(uname *Utsname) error ***REMOVED***
	mib := []_C_int***REMOVED***CTL_KERN, KERN_OSTYPE***REMOVED***
	n := unsafe.Sizeof(uname.Sysname)
	if err := sysctl(mib, &uname.Sysname[0], &n, nil, 0); err != nil ***REMOVED***
		return err
	***REMOVED***

	mib = []_C_int***REMOVED***CTL_KERN, KERN_HOSTNAME***REMOVED***
	n = unsafe.Sizeof(uname.Nodename)
	if err := sysctl(mib, &uname.Nodename[0], &n, nil, 0); err != nil ***REMOVED***
		return err
	***REMOVED***

	mib = []_C_int***REMOVED***CTL_KERN, KERN_OSRELEASE***REMOVED***
	n = unsafe.Sizeof(uname.Release)
	if err := sysctl(mib, &uname.Release[0], &n, nil, 0); err != nil ***REMOVED***
		return err
	***REMOVED***

	mib = []_C_int***REMOVED***CTL_KERN, KERN_VERSION***REMOVED***
	n = unsafe.Sizeof(uname.Version)
	if err := sysctl(mib, &uname.Version[0], &n, nil, 0); err != nil ***REMOVED***
		return err
	***REMOVED***

	// The version might have newlines or tabs in it, convert them to
	// spaces.
	for i, b := range uname.Version ***REMOVED***
		if b == '\n' || b == '\t' ***REMOVED***
			if i == len(uname.Version)-1 ***REMOVED***
				uname.Version[i] = 0
			***REMOVED*** else ***REMOVED***
				uname.Version[i] = ' '
			***REMOVED***
		***REMOVED***
	***REMOVED***

	mib = []_C_int***REMOVED***CTL_HW, HW_MACHINE***REMOVED***
	n = unsafe.Sizeof(uname.Machine)
	if err := sysctl(mib, &uname.Machine[0], &n, nil, 0); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func Stat(path string, st *Stat_t) (err error) ***REMOVED***
	var oldStat stat_freebsd11_t
	if supportsABI(_ino64First) ***REMOVED***
		return fstatat_freebsd12(AT_FDCWD, path, st, 0)
	***REMOVED***
	err = stat(path, &oldStat)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	st.convertFrom(&oldStat)
	return nil
***REMOVED***

func Lstat(path string, st *Stat_t) (err error) ***REMOVED***
	var oldStat stat_freebsd11_t
	if supportsABI(_ino64First) ***REMOVED***
		return fstatat_freebsd12(AT_FDCWD, path, st, AT_SYMLINK_NOFOLLOW)
	***REMOVED***
	err = lstat(path, &oldStat)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	st.convertFrom(&oldStat)
	return nil
***REMOVED***

func Fstat(fd int, st *Stat_t) (err error) ***REMOVED***
	var oldStat stat_freebsd11_t
	if supportsABI(_ino64First) ***REMOVED***
		return fstat_freebsd12(fd, st)
	***REMOVED***
	err = fstat(fd, &oldStat)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	st.convertFrom(&oldStat)
	return nil
***REMOVED***

func Fstatat(fd int, path string, st *Stat_t, flags int) (err error) ***REMOVED***
	var oldStat stat_freebsd11_t
	if supportsABI(_ino64First) ***REMOVED***
		return fstatat_freebsd12(fd, path, st, flags)
	***REMOVED***
	err = fstatat(fd, path, &oldStat, flags)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	st.convertFrom(&oldStat)
	return nil
***REMOVED***

func Statfs(path string, st *Statfs_t) (err error) ***REMOVED***
	var oldStatfs statfs_freebsd11_t
	if supportsABI(_ino64First) ***REMOVED***
		return statfs_freebsd12(path, st)
	***REMOVED***
	err = statfs(path, &oldStatfs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	st.convertFrom(&oldStatfs)
	return nil
***REMOVED***

func Fstatfs(fd int, st *Statfs_t) (err error) ***REMOVED***
	var oldStatfs statfs_freebsd11_t
	if supportsABI(_ino64First) ***REMOVED***
		return fstatfs_freebsd12(fd, st)
	***REMOVED***
	err = fstatfs(fd, &oldStatfs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	st.convertFrom(&oldStatfs)
	return nil
***REMOVED***

func Getdents(fd int, buf []byte) (n int, err error) ***REMOVED***
	return Getdirentries(fd, buf, nil)
***REMOVED***

func Getdirentries(fd int, buf []byte, basep *uintptr) (n int, err error) ***REMOVED***
	if supportsABI(_ino64First) ***REMOVED***
		if basep == nil || unsafe.Sizeof(*basep) == 8 ***REMOVED***
			return getdirentries_freebsd12(fd, buf, (*uint64)(unsafe.Pointer(basep)))
		***REMOVED***
		// The freebsd12 syscall needs a 64-bit base. On 32-bit machines
		// we can't just use the basep passed in. See #32498.
		var base uint64 = uint64(*basep)
		n, err = getdirentries_freebsd12(fd, buf, &base)
		*basep = uintptr(base)
		if base>>32 != 0 ***REMOVED***
			// We can't stuff the base back into a uintptr, so any
			// future calls would be suspect. Generate an error.
			// EIO is allowed by getdirentries.
			err = EIO
		***REMOVED***
		return
	***REMOVED***

	// The old syscall entries are smaller than the new. Use 1/4 of the original
	// buffer size rounded up to DIRBLKSIZ (see /usr/src/lib/libc/sys/getdirentries.c).
	oldBufLen := roundup(len(buf)/4, _dirblksiz)
	oldBuf := make([]byte, oldBufLen)
	n, err = getdirentries(fd, oldBuf, basep)
	if err == nil && n > 0 ***REMOVED***
		n = convertFromDirents11(buf, oldBuf[:n])
	***REMOVED***
	return
***REMOVED***

func Mknod(path string, mode uint32, dev uint64) (err error) ***REMOVED***
	var oldDev int
	if supportsABI(_ino64First) ***REMOVED***
		return mknodat_freebsd12(AT_FDCWD, path, mode, dev)
	***REMOVED***
	oldDev = int(dev)
	return mknod(path, mode, oldDev)
***REMOVED***

func Mknodat(fd int, path string, mode uint32, dev uint64) (err error) ***REMOVED***
	var oldDev int
	if supportsABI(_ino64First) ***REMOVED***
		return mknodat_freebsd12(fd, path, mode, dev)
	***REMOVED***
	oldDev = int(dev)
	return mknodat(fd, path, mode, oldDev)
***REMOVED***

// round x to the nearest multiple of y, larger or equal to x.
//
// from /usr/include/sys/param.h Macros for counting and rounding.
// #define roundup(x, y)   ((((x)+((y)-1))/(y))*(y))
func roundup(x, y int) int ***REMOVED***
	return ((x + y - 1) / y) * y
***REMOVED***

func (s *Stat_t) convertFrom(old *stat_freebsd11_t) ***REMOVED***
	*s = Stat_t***REMOVED***
		Dev:     uint64(old.Dev),
		Ino:     uint64(old.Ino),
		Nlink:   uint64(old.Nlink),
		Mode:    old.Mode,
		Uid:     old.Uid,
		Gid:     old.Gid,
		Rdev:    uint64(old.Rdev),
		Atim:    old.Atim,
		Mtim:    old.Mtim,
		Ctim:    old.Ctim,
		Btim:    old.Btim,
		Size:    old.Size,
		Blocks:  old.Blocks,
		Blksize: old.Blksize,
		Flags:   old.Flags,
		Gen:     uint64(old.Gen),
	***REMOVED***
***REMOVED***

func (s *Statfs_t) convertFrom(old *statfs_freebsd11_t) ***REMOVED***
	*s = Statfs_t***REMOVED***
		Version:     _statfsVersion,
		Type:        old.Type,
		Flags:       old.Flags,
		Bsize:       old.Bsize,
		Iosize:      old.Iosize,
		Blocks:      old.Blocks,
		Bfree:       old.Bfree,
		Bavail:      old.Bavail,
		Files:       old.Files,
		Ffree:       old.Ffree,
		Syncwrites:  old.Syncwrites,
		Asyncwrites: old.Asyncwrites,
		Syncreads:   old.Syncreads,
		Asyncreads:  old.Asyncreads,
		// Spare
		Namemax: old.Namemax,
		Owner:   old.Owner,
		Fsid:    old.Fsid,
		// Charspare
		// Fstypename
		// Mntfromname
		// Mntonname
	***REMOVED***

	sl := old.Fstypename[:]
	n := clen(*(*[]byte)(unsafe.Pointer(&sl)))
	copy(s.Fstypename[:], old.Fstypename[:n])

	sl = old.Mntfromname[:]
	n = clen(*(*[]byte)(unsafe.Pointer(&sl)))
	copy(s.Mntfromname[:], old.Mntfromname[:n])

	sl = old.Mntonname[:]
	n = clen(*(*[]byte)(unsafe.Pointer(&sl)))
	copy(s.Mntonname[:], old.Mntonname[:n])
***REMOVED***

func convertFromDirents11(buf []byte, old []byte) int ***REMOVED***
	const (
		fixedSize    = int(unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Name))
		oldFixedSize = int(unsafe.Offsetof(dirent_freebsd11***REMOVED******REMOVED***.Name))
	)

	dstPos := 0
	srcPos := 0
	for dstPos+fixedSize < len(buf) && srcPos+oldFixedSize < len(old) ***REMOVED***
		var dstDirent Dirent
		var srcDirent dirent_freebsd11

		// If multiple direntries are written, sometimes when we reach the final one,
		// we may have cap of old less than size of dirent_freebsd11.
		copy((*[unsafe.Sizeof(srcDirent)]byte)(unsafe.Pointer(&srcDirent))[:], old[srcPos:])

		reclen := roundup(fixedSize+int(srcDirent.Namlen)+1, 8)
		if dstPos+reclen > len(buf) ***REMOVED***
			break
		***REMOVED***

		dstDirent.Fileno = uint64(srcDirent.Fileno)
		dstDirent.Off = 0
		dstDirent.Reclen = uint16(reclen)
		dstDirent.Type = srcDirent.Type
		dstDirent.Pad0 = 0
		dstDirent.Namlen = uint16(srcDirent.Namlen)
		dstDirent.Pad1 = 0

		copy(dstDirent.Name[:], srcDirent.Name[:srcDirent.Namlen])
		copy(buf[dstPos:], (*[unsafe.Sizeof(dstDirent)]byte)(unsafe.Pointer(&dstDirent))[:])
		padding := buf[dstPos+fixedSize+int(dstDirent.Namlen) : dstPos+reclen]
		for i := range padding ***REMOVED***
			padding[i] = 0
		***REMOVED***

		dstPos += int(dstDirent.Reclen)
		srcPos += int(srcDirent.Reclen)
	***REMOVED***

	return dstPos
***REMOVED***

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) ***REMOVED***
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	return sendfile(outfd, infd, offset, count)
***REMOVED***

//sys	ptrace(request int, pid int, addr uintptr, data int) (err error)

func PtraceAttach(pid int) (err error) ***REMOVED***
	return ptrace(PTRACE_ATTACH, pid, 0, 0)
***REMOVED***

func PtraceCont(pid int, signal int) (err error) ***REMOVED***
	return ptrace(PTRACE_CONT, pid, 1, signal)
***REMOVED***

func PtraceDetach(pid int) (err error) ***REMOVED***
	return ptrace(PTRACE_DETACH, pid, 1, 0)
***REMOVED***

func PtraceGetFpRegs(pid int, fpregsout *FpReg) (err error) ***REMOVED***
	return ptrace(PTRACE_GETFPREGS, pid, uintptr(unsafe.Pointer(fpregsout)), 0)
***REMOVED***

func PtraceGetRegs(pid int, regsout *Reg) (err error) ***REMOVED***
	return ptrace(PTRACE_GETREGS, pid, uintptr(unsafe.Pointer(regsout)), 0)
***REMOVED***

func PtraceLwpEvents(pid int, enable int) (err error) ***REMOVED***
	return ptrace(PTRACE_LWPEVENTS, pid, 0, enable)
***REMOVED***

func PtraceLwpInfo(pid int, info uintptr) (err error) ***REMOVED***
	return ptrace(PTRACE_LWPINFO, pid, info, int(unsafe.Sizeof(PtraceLwpInfoStruct***REMOVED******REMOVED***)))
***REMOVED***

func PtracePeekData(pid int, addr uintptr, out []byte) (count int, err error) ***REMOVED***
	return PtraceIO(PIOD_READ_D, pid, addr, out, SizeofLong)
***REMOVED***

func PtracePeekText(pid int, addr uintptr, out []byte) (count int, err error) ***REMOVED***
	return PtraceIO(PIOD_READ_I, pid, addr, out, SizeofLong)
***REMOVED***

func PtracePokeData(pid int, addr uintptr, data []byte) (count int, err error) ***REMOVED***
	return PtraceIO(PIOD_WRITE_D, pid, addr, data, SizeofLong)
***REMOVED***

func PtracePokeText(pid int, addr uintptr, data []byte) (count int, err error) ***REMOVED***
	return PtraceIO(PIOD_WRITE_I, pid, addr, data, SizeofLong)
***REMOVED***

func PtraceSetRegs(pid int, regs *Reg) (err error) ***REMOVED***
	return ptrace(PTRACE_SETREGS, pid, uintptr(unsafe.Pointer(regs)), 0)
***REMOVED***

func PtraceSingleStep(pid int) (err error) ***REMOVED***
	return ptrace(PTRACE_SINGLESTEP, pid, 1, 0)
***REMOVED***

/*
 * Exposed directly
 */
//sys	Access(path string, mode uint32) (err error)
//sys	Adjtime(delta *Timeval, olddelta *Timeval) (err error)
//sys	CapEnter() (err error)
//sys	capRightsGet(version int, fd int, rightsp *CapRights) (err error) = SYS___CAP_RIGHTS_GET
//sys	capRightsLimit(fd int, rightsp *CapRights) (err error)
//sys	Chdir(path string) (err error)
//sys	Chflags(path string, flags int) (err error)
//sys	Chmod(path string, mode uint32) (err error)
//sys	Chown(path string, uid int, gid int) (err error)
//sys	Chroot(path string) (err error)
//sys	Close(fd int) (err error)
//sys	Dup(fd int) (nfd int, err error)
//sys	Dup2(from int, to int) (err error)
//sys	Exit(code int)
//sys	ExtattrGetFd(fd int, attrnamespace int, attrname string, data uintptr, nbytes int) (ret int, err error)
//sys	ExtattrSetFd(fd int, attrnamespace int, attrname string, data uintptr, nbytes int) (ret int, err error)
//sys	ExtattrDeleteFd(fd int, attrnamespace int, attrname string) (err error)
//sys	ExtattrListFd(fd int, attrnamespace int, data uintptr, nbytes int) (ret int, err error)
//sys	ExtattrGetFile(file string, attrnamespace int, attrname string, data uintptr, nbytes int) (ret int, err error)
//sys	ExtattrSetFile(file string, attrnamespace int, attrname string, data uintptr, nbytes int) (ret int, err error)
//sys	ExtattrDeleteFile(file string, attrnamespace int, attrname string) (err error)
//sys	ExtattrListFile(file string, attrnamespace int, data uintptr, nbytes int) (ret int, err error)
//sys	ExtattrGetLink(link string, attrnamespace int, attrname string, data uintptr, nbytes int) (ret int, err error)
//sys	ExtattrSetLink(link string, attrnamespace int, attrname string, data uintptr, nbytes int) (ret int, err error)
//sys	ExtattrDeleteLink(link string, attrnamespace int, attrname string) (err error)
//sys	ExtattrListLink(link string, attrnamespace int, data uintptr, nbytes int) (ret int, err error)
//sys	Fadvise(fd int, offset int64, length int64, advice int) (err error) = SYS_POSIX_FADVISE
//sys	Faccessat(dirfd int, path string, mode uint32, flags int) (err error)
//sys	Fchdir(fd int) (err error)
//sys	Fchflags(fd int, flags int) (err error)
//sys	Fchmod(fd int, mode uint32) (err error)
//sys	Fchmodat(dirfd int, path string, mode uint32, flags int) (err error)
//sys	Fchown(fd int, uid int, gid int) (err error)
//sys	Fchownat(dirfd int, path string, uid int, gid int, flags int) (err error)
//sys	Flock(fd int, how int) (err error)
//sys	Fpathconf(fd int, name int) (val int, err error)
//sys	fstat(fd int, stat *stat_freebsd11_t) (err error)
//sys	fstat_freebsd12(fd int, stat *Stat_t) (err error)
//sys	fstatat(fd int, path string, stat *stat_freebsd11_t, flags int) (err error)
//sys	fstatat_freebsd12(fd int, path string, stat *Stat_t, flags int) (err error)
//sys	fstatfs(fd int, stat *statfs_freebsd11_t) (err error)
//sys	fstatfs_freebsd12(fd int, stat *Statfs_t) (err error)
//sys	Fsync(fd int) (err error)
//sys	Ftruncate(fd int, length int64) (err error)
//sys	getdirentries(fd int, buf []byte, basep *uintptr) (n int, err error)
//sys	getdirentries_freebsd12(fd int, buf []byte, basep *uint64) (n int, err error)
//sys	Getdtablesize() (size int)
//sysnb	Getegid() (egid int)
//sysnb	Geteuid() (uid int)
//sysnb	Getgid() (gid int)
//sysnb	Getpgid(pid int) (pgid int, err error)
//sysnb	Getpgrp() (pgrp int)
//sysnb	Getpid() (pid int)
//sysnb	Getppid() (ppid int)
//sys	Getpriority(which int, who int) (prio int, err error)
//sysnb	Getrlimit(which int, lim *Rlimit) (err error)
//sysnb	Getrusage(who int, rusage *Rusage) (err error)
//sysnb	Getsid(pid int) (sid int, err error)
//sysnb	Gettimeofday(tv *Timeval) (err error)
//sysnb	Getuid() (uid int)
//sys	Issetugid() (tainted bool)
//sys	Kill(pid int, signum syscall.Signal) (err error)
//sys	Kqueue() (fd int, err error)
//sys	Lchown(path string, uid int, gid int) (err error)
//sys	Link(path string, link string) (err error)
//sys	Linkat(pathfd int, path string, linkfd int, link string, flags int) (err error)
//sys	Listen(s int, backlog int) (err error)
//sys	lstat(path string, stat *stat_freebsd11_t) (err error)
//sys	Mkdir(path string, mode uint32) (err error)
//sys	Mkdirat(dirfd int, path string, mode uint32) (err error)
//sys	Mkfifo(path string, mode uint32) (err error)
//sys	mknod(path string, mode uint32, dev int) (err error)
//sys	mknodat(fd int, path string, mode uint32, dev int) (err error)
//sys	mknodat_freebsd12(fd int, path string, mode uint32, dev uint64) (err error)
//sys	Nanosleep(time *Timespec, leftover *Timespec) (err error)
//sys	Open(path string, mode int, perm uint32) (fd int, err error)
//sys	Openat(fdat int, path string, mode int, perm uint32) (fd int, err error)
//sys	Pathconf(path string, name int) (val int, err error)
//sys	Pread(fd int, p []byte, offset int64) (n int, err error)
//sys	Pwrite(fd int, p []byte, offset int64) (n int, err error)
//sys	read(fd int, p []byte) (n int, err error)
//sys	Readlink(path string, buf []byte) (n int, err error)
//sys	Readlinkat(dirfd int, path string, buf []byte) (n int, err error)
//sys	Rename(from string, to string) (err error)
//sys	Renameat(fromfd int, from string, tofd int, to string) (err error)
//sys	Revoke(path string) (err error)
//sys	Rmdir(path string) (err error)
//sys	Seek(fd int, offset int64, whence int) (newoffset int64, err error) = SYS_LSEEK
//sys	Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, err error)
//sysnb	Setegid(egid int) (err error)
//sysnb	Seteuid(euid int) (err error)
//sysnb	Setgid(gid int) (err error)
//sys	Setlogin(name string) (err error)
//sysnb	Setpgid(pid int, pgid int) (err error)
//sys	Setpriority(which int, who int, prio int) (err error)
//sysnb	Setregid(rgid int, egid int) (err error)
//sysnb	Setreuid(ruid int, euid int) (err error)
//sysnb	Setresgid(rgid int, egid int, sgid int) (err error)
//sysnb	Setresuid(ruid int, euid int, suid int) (err error)
//sysnb	Setrlimit(which int, lim *Rlimit) (err error)
//sysnb	Setsid() (pid int, err error)
//sysnb	Settimeofday(tp *Timeval) (err error)
//sysnb	Setuid(uid int) (err error)
//sys	stat(path string, stat *stat_freebsd11_t) (err error)
//sys	statfs(path string, stat *statfs_freebsd11_t) (err error)
//sys	statfs_freebsd12(path string, stat *Statfs_t) (err error)
//sys	Symlink(path string, link string) (err error)
//sys	Symlinkat(oldpath string, newdirfd int, newpath string) (err error)
//sys	Sync() (err error)
//sys	Truncate(path string, length int64) (err error)
//sys	Umask(newmask int) (oldmask int)
//sys	Undelete(path string) (err error)
//sys	Unlink(path string) (err error)
//sys	Unlinkat(dirfd int, path string, flags int) (err error)
//sys	Unmount(path string, flags int) (err error)
//sys	write(fd int, p []byte) (n int, err error)
//sys   mmap(addr uintptr, length uintptr, prot int, flag int, fd int, pos int64) (ret uintptr, err error)
//sys   munmap(addr uintptr, length uintptr) (err error)
//sys	readlen(fd int, buf *byte, nbuf int) (n int, err error) = SYS_READ
//sys	writelen(fd int, buf *byte, nbuf int) (n int, err error) = SYS_WRITE
//sys	accept4(fd int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (nfd int, err error)
//sys	utimensat(dirfd int, path string, times *[2]Timespec, flags int) (err error)

/*
 * Unimplemented
 */
// Profil
// Sigaction
// Sigprocmask
// Getlogin
// Sigpending
// Sigaltstack
// Ioctl
// Reboot
// Execve
// Vfork
// Sbrk
// Sstk
// Ovadvise
// Mincore
// Setitimer
// Swapon
// Select
// Sigsuspend
// Readv
// Writev
// Nfssvc
// Getfh
// Quotactl
// Mount
// Csops
// Waitid
// Add_profil
// Kdebug_trace
// Sigreturn
// Atsocket
// Kqueue_from_portset_np
// Kqueue_portset
// Getattrlist
// Setattrlist
// Getdents
// Getdirentriesattr
// Searchfs
// Delete
// Copyfile
// Watchevent
// Waitevent
// Modwatch
// Fsctl
// Initgroups
// Posix_spawn
// Nfsclnt
// Fhopen
// Minherit
// Semsys
// Msgsys
// Shmsys
// Semctl
// Semget
// Semop
// Msgctl
// Msgget
// Msgsnd
// Msgrcv
// Shmat
// Shmctl
// Shmdt
// Shmget
// Shm_open
// Shm_unlink
// Sem_open
// Sem_close
// Sem_unlink
// Sem_wait
// Sem_trywait
// Sem_post
// Sem_getvalue
// Sem_init
// Sem_destroy
// Open_extended
// Umask_extended
// Stat_extended
// Lstat_extended
// Fstat_extended
// Chmod_extended
// Fchmod_extended
// Access_extended
// Settid
// Gettid
// Setsgroups
// Getsgroups
// Setwgroups
// Getwgroups
// Mkfifo_extended
// Mkdir_extended
// Identitysvc
// Shared_region_check_np
// Shared_region_map_np
// __pthread_mutex_destroy
// __pthread_mutex_init
// __pthread_mutex_lock
// __pthread_mutex_trylock
// __pthread_mutex_unlock
// __pthread_cond_init
// __pthread_cond_destroy
// __pthread_cond_broadcast
// __pthread_cond_signal
// Setsid_with_pid
// __pthread_cond_timedwait
// Aio_fsync
// Aio_return
// Aio_suspend
// Aio_cancel
// Aio_error
// Aio_read
// Aio_write
// Lio_listio
// __pthread_cond_wait
// Iopolicysys
// __pthread_kill
// __pthread_sigmask
// __sigwait
// __disable_threadsignal
// __pthread_markcancel
// __pthread_canceled
// __semwait_signal
// Proc_info
// Stat64_extended
// Lstat64_extended
// Fstat64_extended
// __pthread_chdir
// __pthread_fchdir
// Audit
// Auditon
// Getauid
// Setauid
// Getaudit
// Setaudit
// Getaudit_addr
// Setaudit_addr
// Auditctl
// Bsdthread_create
// Bsdthread_terminate
// Stack_snapshot
// Bsdthread_register
// Workq_open
// Workq_ops
// __mac_execve
// __mac_syscall
// __mac_get_file
// __mac_set_file
// __mac_get_link
// __mac_set_link
// __mac_get_proc
// __mac_set_proc
// __mac_get_fd
// __mac_set_fd
// __mac_get_pid
// __mac_get_lcid
// __mac_get_lctx
// __mac_set_lctx
// Setlcid
// Read_nocancel
// Write_nocancel
// Open_nocancel
// Close_nocancel
// Wait4_nocancel
// Recvmsg_nocancel
// Sendmsg_nocancel
// Recvfrom_nocancel
// Accept_nocancel
// Fcntl_nocancel
// Select_nocancel
// Fsync_nocancel
// Connect_nocancel
// Sigsuspend_nocancel
// Readv_nocancel
// Writev_nocancel
// Sendto_nocancel
// Pread_nocancel
// Pwrite_nocancel
// Waitid_nocancel
// Poll_nocancel
// Msgsnd_nocancel
// Msgrcv_nocancel
// Sem_wait_nocancel
// Aio_suspend_nocancel
// __sigwait_nocancel
// __semwait_signal_nocancel
// __mac_mount
// __mac_get_mount
// __mac_getfsstat
