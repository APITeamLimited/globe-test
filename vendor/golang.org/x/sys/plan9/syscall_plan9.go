// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Plan 9 system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and
// wrap it in our own nicer implementation.

package plan9

import (
	"bytes"
	"syscall"
	"unsafe"
)

// A Note is a string describing a process note.
// It implements the os.Signal interface.
type Note string

func (n Note) Signal() ***REMOVED******REMOVED***

func (n Note) String() string ***REMOVED***
	return string(n)
***REMOVED***

var (
	Stdin  = 0
	Stdout = 1
	Stderr = 2
)

// For testing: clients can set this flag to force
// creation of IPv6 sockets to return EAFNOSUPPORT.
var SocketDisableIPv6 bool

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.ErrorString)
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.ErrorString)
func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr)
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr)

func atoi(b []byte) (n uint) ***REMOVED***
	n = 0
	for i := 0; i < len(b); i++ ***REMOVED***
		n = n*10 + uint(b[i]-'0')
	***REMOVED***
	return
***REMOVED***

func cstring(s []byte) string ***REMOVED***
	i := bytes.IndexByte(s, 0)
	if i == -1 ***REMOVED***
		i = len(s)
	***REMOVED***
	return string(s[:i])
***REMOVED***

func errstr() string ***REMOVED***
	var buf [ERRMAX]byte

	RawSyscall(SYS_ERRSTR, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0)

	buf[len(buf)-1] = 0
	return cstring(buf[:])
***REMOVED***

// Implemented in assembly to import from runtime.
func exit(code int)

func Exit(code int) ***REMOVED*** exit(code) ***REMOVED***

func readnum(path string) (uint, error) ***REMOVED***
	var b [12]byte

	fd, e := Open(path, O_RDONLY)
	if e != nil ***REMOVED***
		return 0, e
	***REMOVED***
	defer Close(fd)

	n, e := Pread(fd, b[:], 0)

	if e != nil ***REMOVED***
		return 0, e
	***REMOVED***

	m := 0
	for ; m < n && b[m] == ' '; m++ ***REMOVED***
	***REMOVED***

	return atoi(b[m : n-1]), nil
***REMOVED***

func Getpid() (pid int) ***REMOVED***
	n, _ := readnum("#c/pid")
	return int(n)
***REMOVED***

func Getppid() (ppid int) ***REMOVED***
	n, _ := readnum("#c/ppid")
	return int(n)
***REMOVED***

func Read(fd int, p []byte) (n int, err error) ***REMOVED***
	return Pread(fd, p, -1)
***REMOVED***

func Write(fd int, p []byte) (n int, err error) ***REMOVED***
	return Pwrite(fd, p, -1)
***REMOVED***

var ioSync int64

//sys	fd2path(fd int, buf []byte) (err error)
func Fd2path(fd int) (path string, err error) ***REMOVED***
	var buf [512]byte

	e := fd2path(fd, buf[:])
	if e != nil ***REMOVED***
		return "", e
	***REMOVED***
	return cstring(buf[:]), nil
***REMOVED***

//sys	pipe(p *[2]int32) (err error)
func Pipe(p []int) (err error) ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return syscall.ErrorString("bad arg in system call")
	***REMOVED***
	var pp [2]int32
	err = pipe(&pp)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
***REMOVED***

// Underlying system call writes to newoffset via pointer.
// Implemented in assembly to avoid allocation.
func seek(placeholder uintptr, fd int, offset int64, whence int) (newoffset int64, err string)

func Seek(fd int, offset int64, whence int) (newoffset int64, err error) ***REMOVED***
	newoffset, e := seek(0, fd, offset, whence)

	if newoffset == -1 ***REMOVED***
		err = syscall.ErrorString(e)
	***REMOVED***
	return
***REMOVED***

func Mkdir(path string, mode uint32) (err error) ***REMOVED***
	fd, err := Create(path, O_RDONLY, DMDIR|mode)

	if fd != -1 ***REMOVED***
		Close(fd)
	***REMOVED***

	return
***REMOVED***

type Waitmsg struct ***REMOVED***
	Pid  int
	Time [3]uint32
	Msg  string
***REMOVED***

func (w Waitmsg) Exited() bool   ***REMOVED*** return true ***REMOVED***
func (w Waitmsg) Signaled() bool ***REMOVED*** return false ***REMOVED***

func (w Waitmsg) ExitStatus() int ***REMOVED***
	if len(w.Msg) == 0 ***REMOVED***
		// a normal exit returns no message
		return 0
	***REMOVED***
	return 1
***REMOVED***

//sys	await(s []byte) (n int, err error)
func Await(w *Waitmsg) (err error) ***REMOVED***
	var buf [512]byte
	var f [5][]byte

	n, err := await(buf[:])

	if err != nil || w == nil ***REMOVED***
		return
	***REMOVED***

	nf := 0
	p := 0
	for i := 0; i < n && nf < len(f)-1; i++ ***REMOVED***
		if buf[i] == ' ' ***REMOVED***
			f[nf] = buf[p:i]
			p = i + 1
			nf++
		***REMOVED***
	***REMOVED***
	f[nf] = buf[p:]
	nf++

	if nf != len(f) ***REMOVED***
		return syscall.ErrorString("invalid wait message")
	***REMOVED***
	w.Pid = int(atoi(f[0]))
	w.Time[0] = uint32(atoi(f[1]))
	w.Time[1] = uint32(atoi(f[2]))
	w.Time[2] = uint32(atoi(f[3]))
	w.Msg = cstring(f[4])
	if w.Msg == "''" ***REMOVED***
		// await() returns '' for no error
		w.Msg = ""
	***REMOVED***
	return
***REMOVED***

func Unmount(name, old string) (err error) ***REMOVED***
	fixwd()
	oldp, err := BytePtrFromString(old)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	oldptr := uintptr(unsafe.Pointer(oldp))

	var r0 uintptr
	var e syscall.ErrorString

	// bind(2) man page: If name is zero, everything bound or mounted upon old is unbound or unmounted.
	if name == "" ***REMOVED***
		r0, _, e = Syscall(SYS_UNMOUNT, _zero, oldptr, 0)
	***REMOVED*** else ***REMOVED***
		namep, err := BytePtrFromString(name)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		r0, _, e = Syscall(SYS_UNMOUNT, uintptr(unsafe.Pointer(namep)), oldptr, 0)
	***REMOVED***

	if int32(r0) == -1 ***REMOVED***
		err = e
	***REMOVED***
	return
***REMOVED***

func Fchdir(fd int) (err error) ***REMOVED***
	path, err := Fd2path(fd)

	if err != nil ***REMOVED***
		return
	***REMOVED***

	return Chdir(path)
***REMOVED***

type Timespec struct ***REMOVED***
	Sec  int32
	Nsec int32
***REMOVED***

type Timeval struct ***REMOVED***
	Sec  int32
	Usec int32
***REMOVED***

func NsecToTimeval(nsec int64) (tv Timeval) ***REMOVED***
	nsec += 999 // round up to microsecond
	tv.Usec = int32(nsec % 1e9 / 1e3)
	tv.Sec = int32(nsec / 1e9)
	return
***REMOVED***

func nsec() int64 ***REMOVED***
	var scratch int64

	r0, _, _ := Syscall(SYS_NSEC, uintptr(unsafe.Pointer(&scratch)), 0, 0)
	// TODO(aram): remove hack after I fix _nsec in the pc64 kernel.
	if r0 == 0 ***REMOVED***
		return scratch
	***REMOVED***
	return int64(r0)
***REMOVED***

func Gettimeofday(tv *Timeval) error ***REMOVED***
	nsec := nsec()
	*tv = NsecToTimeval(nsec)
	return nil
***REMOVED***

func Getpagesize() int ***REMOVED*** return 0x1000 ***REMOVED***

func Getegid() (egid int) ***REMOVED*** return -1 ***REMOVED***
func Geteuid() (euid int) ***REMOVED*** return -1 ***REMOVED***
func Getgid() (gid int)   ***REMOVED*** return -1 ***REMOVED***
func Getuid() (uid int)   ***REMOVED*** return -1 ***REMOVED***

func Getgroups() (gids []int, err error) ***REMOVED***
	return make([]int, 0), nil
***REMOVED***

//sys	open(path string, mode int) (fd int, err error)
func Open(path string, mode int) (fd int, err error) ***REMOVED***
	fixwd()
	return open(path, mode)
***REMOVED***

//sys	create(path string, mode int, perm uint32) (fd int, err error)
func Create(path string, mode int, perm uint32) (fd int, err error) ***REMOVED***
	fixwd()
	return create(path, mode, perm)
***REMOVED***

//sys	remove(path string) (err error)
func Remove(path string) error ***REMOVED***
	fixwd()
	return remove(path)
***REMOVED***

//sys	stat(path string, edir []byte) (n int, err error)
func Stat(path string, edir []byte) (n int, err error) ***REMOVED***
	fixwd()
	return stat(path, edir)
***REMOVED***

//sys	bind(name string, old string, flag int) (err error)
func Bind(name string, old string, flag int) (err error) ***REMOVED***
	fixwd()
	return bind(name, old, flag)
***REMOVED***

//sys	mount(fd int, afd int, old string, flag int, aname string) (err error)
func Mount(fd int, afd int, old string, flag int, aname string) (err error) ***REMOVED***
	fixwd()
	return mount(fd, afd, old, flag, aname)
***REMOVED***

//sys	wstat(path string, edir []byte) (err error)
func Wstat(path string, edir []byte) (err error) ***REMOVED***
	fixwd()
	return wstat(path, edir)
***REMOVED***

//sys	chdir(path string) (err error)
//sys	Dup(oldfd int, newfd int) (fd int, err error)
//sys	Pread(fd int, p []byte, offset int64) (n int, err error)
//sys	Pwrite(fd int, p []byte, offset int64) (n int, err error)
//sys	Close(fd int) (err error)
//sys	Fstat(fd int, edir []byte) (n int, err error)
//sys	Fwstat(fd int, edir []byte) (err error)
