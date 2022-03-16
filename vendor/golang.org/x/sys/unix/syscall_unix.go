// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package unix

import (
	"bytes"
	"sort"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/internal/unsafeheader"
)

var (
	Stdin  = 0
	Stdout = 1
	Stderr = 2
)

// Do the interface allocations only once for common
// Errno values.
var (
	errEAGAIN error = syscall.EAGAIN
	errEINVAL error = syscall.EINVAL
	errENOENT error = syscall.ENOENT
)

var (
	signalNameMapOnce sync.Once
	signalNameMap     map[string]syscall.Signal
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error ***REMOVED***
	switch e ***REMOVED***
	case 0:
		return nil
	case EAGAIN:
		return errEAGAIN
	case EINVAL:
		return errEINVAL
	case ENOENT:
		return errENOENT
	***REMOVED***
	return e
***REMOVED***

// ErrnoName returns the error name for error number e.
func ErrnoName(e syscall.Errno) string ***REMOVED***
	i := sort.Search(len(errorList), func(i int) bool ***REMOVED***
		return errorList[i].num >= e
	***REMOVED***)
	if i < len(errorList) && errorList[i].num == e ***REMOVED***
		return errorList[i].name
	***REMOVED***
	return ""
***REMOVED***

// SignalName returns the signal name for signal number s.
func SignalName(s syscall.Signal) string ***REMOVED***
	i := sort.Search(len(signalList), func(i int) bool ***REMOVED***
		return signalList[i].num >= s
	***REMOVED***)
	if i < len(signalList) && signalList[i].num == s ***REMOVED***
		return signalList[i].name
	***REMOVED***
	return ""
***REMOVED***

// SignalNum returns the syscall.Signal for signal named s,
// or 0 if a signal with such name is not found.
// The signal name should start with "SIG".
func SignalNum(s string) syscall.Signal ***REMOVED***
	signalNameMapOnce.Do(func() ***REMOVED***
		signalNameMap = make(map[string]syscall.Signal, len(signalList))
		for _, signal := range signalList ***REMOVED***
			signalNameMap[signal.name] = signal.num
		***REMOVED***
	***REMOVED***)
	return signalNameMap[s]
***REMOVED***

// clen returns the index of the first NULL byte in n or len(n) if n contains no NULL byte.
func clen(n []byte) int ***REMOVED***
	i := bytes.IndexByte(n, 0)
	if i == -1 ***REMOVED***
		i = len(n)
	***REMOVED***
	return i
***REMOVED***

// Mmap manager, for use by operating system-specific implementations.

type mmapper struct ***REMOVED***
	sync.Mutex
	active map[*byte][]byte // active mappings; key is last byte in mapping
	mmap   func(addr, length uintptr, prot, flags, fd int, offset int64) (uintptr, error)
	munmap func(addr uintptr, length uintptr) error
***REMOVED***

func (m *mmapper) Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) ***REMOVED***
	if length <= 0 ***REMOVED***
		return nil, EINVAL
	***REMOVED***

	// Map the requested memory.
	addr, errno := m.mmap(0, uintptr(length), prot, flags, fd, offset)
	if errno != nil ***REMOVED***
		return nil, errno
	***REMOVED***

	// Use unsafe to convert addr into a []byte.
	var b []byte
	hdr := (*unsafeheader.Slice)(unsafe.Pointer(&b))
	hdr.Data = unsafe.Pointer(addr)
	hdr.Cap = length
	hdr.Len = length

	// Register mapping in m and return it.
	p := &b[cap(b)-1]
	m.Lock()
	defer m.Unlock()
	m.active[p] = b
	return b, nil
***REMOVED***

func (m *mmapper) Munmap(data []byte) (err error) ***REMOVED***
	if len(data) == 0 || len(data) != cap(data) ***REMOVED***
		return EINVAL
	***REMOVED***

	// Find the base of the mapping.
	p := &data[cap(data)-1]
	m.Lock()
	defer m.Unlock()
	b := m.active[p]
	if b == nil || &b[0] != &data[0] ***REMOVED***
		return EINVAL
	***REMOVED***

	// Unmap the memory and update m.
	if errno := m.munmap(uintptr(unsafe.Pointer(&b[0])), uintptr(len(b))); errno != nil ***REMOVED***
		return errno
	***REMOVED***
	delete(m.active, p)
	return nil
***REMOVED***

func Read(fd int, p []byte) (n int, err error) ***REMOVED***
	n, err = read(fd, p)
	if raceenabled ***REMOVED***
		if n > 0 ***REMOVED***
			raceWriteRange(unsafe.Pointer(&p[0]), n)
		***REMOVED***
		if err == nil ***REMOVED***
			raceAcquire(unsafe.Pointer(&ioSync))
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func Write(fd int, p []byte) (n int, err error) ***REMOVED***
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	n, err = write(fd, p)
	if raceenabled && n > 0 ***REMOVED***
		raceReadRange(unsafe.Pointer(&p[0]), n)
	***REMOVED***
	return
***REMOVED***

// For testing: clients can set this flag to force
// creation of IPv6 sockets to return EAFNOSUPPORT.
var SocketDisableIPv6 bool

// Sockaddr represents a socket address.
type Sockaddr interface ***REMOVED***
	sockaddr() (ptr unsafe.Pointer, len _Socklen, err error) // lowercase; only we can define Sockaddrs
***REMOVED***

// SockaddrInet4 implements the Sockaddr interface for AF_INET type sockets.
type SockaddrInet4 struct ***REMOVED***
	Port int
	Addr [4]byte
	raw  RawSockaddrInet4
***REMOVED***

// SockaddrInet6 implements the Sockaddr interface for AF_INET6 type sockets.
type SockaddrInet6 struct ***REMOVED***
	Port   int
	ZoneId uint32
	Addr   [16]byte
	raw    RawSockaddrInet6
***REMOVED***

// SockaddrUnix implements the Sockaddr interface for AF_UNIX type sockets.
type SockaddrUnix struct ***REMOVED***
	Name string
	raw  RawSockaddrUnix
***REMOVED***

func Bind(fd int, sa Sockaddr) (err error) ***REMOVED***
	ptr, n, err := sa.sockaddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return bind(fd, ptr, n)
***REMOVED***

func Connect(fd int, sa Sockaddr) (err error) ***REMOVED***
	ptr, n, err := sa.sockaddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return connect(fd, ptr, n)
***REMOVED***

func Getpeername(fd int) (sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getpeername(fd, &rsa, &len); err != nil ***REMOVED***
		return
	***REMOVED***
	return anyToSockaddr(fd, &rsa)
***REMOVED***

func GetsockoptByte(fd, level, opt int) (value byte, err error) ***REMOVED***
	var n byte
	vallen := _Socklen(1)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return n, err
***REMOVED***

func GetsockoptInt(fd, level, opt int) (value int, err error) ***REMOVED***
	var n int32
	vallen := _Socklen(4)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return int(n), err
***REMOVED***

func GetsockoptInet4Addr(fd, level, opt int) (value [4]byte, err error) ***REMOVED***
	vallen := _Socklen(4)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&value[0]), &vallen)
	return value, err
***REMOVED***

func GetsockoptIPMreq(fd, level, opt int) (*IPMreq, error) ***REMOVED***
	var value IPMreq
	vallen := _Socklen(SizeofIPMreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptIPv6Mreq(fd, level, opt int) (*IPv6Mreq, error) ***REMOVED***
	var value IPv6Mreq
	vallen := _Socklen(SizeofIPv6Mreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptIPv6MTUInfo(fd, level, opt int) (*IPv6MTUInfo, error) ***REMOVED***
	var value IPv6MTUInfo
	vallen := _Socklen(SizeofIPv6MTUInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptICMPv6Filter(fd, level, opt int) (*ICMPv6Filter, error) ***REMOVED***
	var value ICMPv6Filter
	vallen := _Socklen(SizeofICMPv6Filter)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptLinger(fd, level, opt int) (*Linger, error) ***REMOVED***
	var linger Linger
	vallen := _Socklen(SizeofLinger)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&linger), &vallen)
	return &linger, err
***REMOVED***

func GetsockoptTimeval(fd, level, opt int) (*Timeval, error) ***REMOVED***
	var tv Timeval
	vallen := _Socklen(unsafe.Sizeof(tv))
	err := getsockopt(fd, level, opt, unsafe.Pointer(&tv), &vallen)
	return &tv, err
***REMOVED***

func GetsockoptUint64(fd, level, opt int) (value uint64, err error) ***REMOVED***
	var n uint64
	vallen := _Socklen(8)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return n, err
***REMOVED***

func Recvfrom(fd int, p []byte, flags int) (n int, from Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if n, err = recvfrom(fd, p, flags, &rsa, &len); err != nil ***REMOVED***
		return
	***REMOVED***
	if rsa.Addr.Family != AF_UNSPEC ***REMOVED***
		from, err = anyToSockaddr(fd, &rsa)
	***REMOVED***
	return
***REMOVED***

func Send(s int, buf []byte, flags int) (err error) ***REMOVED***
	return sendto(s, buf, flags, nil, 0)
***REMOVED***

func Sendto(fd int, p []byte, flags int, to Sockaddr) (err error) ***REMOVED***
	ptr, n, err := to.sockaddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return sendto(fd, p, flags, ptr, n)
***REMOVED***

func SetsockoptByte(fd, level, opt int, value byte) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(&value), 1)
***REMOVED***

func SetsockoptInt(fd, level, opt int, value int) (err error) ***REMOVED***
	var n = int32(value)
	return setsockopt(fd, level, opt, unsafe.Pointer(&n), 4)
***REMOVED***

func SetsockoptInet4Addr(fd, level, opt int, value [4]byte) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(&value[0]), 4)
***REMOVED***

func SetsockoptIPMreq(fd, level, opt int, mreq *IPMreq) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), SizeofIPMreq)
***REMOVED***

func SetsockoptIPv6Mreq(fd, level, opt int, mreq *IPv6Mreq) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), SizeofIPv6Mreq)
***REMOVED***

func SetsockoptICMPv6Filter(fd, level, opt int, filter *ICMPv6Filter) error ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(filter), SizeofICMPv6Filter)
***REMOVED***

func SetsockoptLinger(fd, level, opt int, l *Linger) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(l), SizeofLinger)
***REMOVED***

func SetsockoptString(fd, level, opt int, s string) (err error) ***REMOVED***
	var p unsafe.Pointer
	if len(s) > 0 ***REMOVED***
		p = unsafe.Pointer(&[]byte(s)[0])
	***REMOVED***
	return setsockopt(fd, level, opt, p, uintptr(len(s)))
***REMOVED***

func SetsockoptTimeval(fd, level, opt int, tv *Timeval) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(tv), unsafe.Sizeof(*tv))
***REMOVED***

func SetsockoptUint64(fd, level, opt int, value uint64) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(&value), 8)
***REMOVED***

func Socket(domain, typ, proto int) (fd int, err error) ***REMOVED***
	if domain == AF_INET6 && SocketDisableIPv6 ***REMOVED***
		return -1, EAFNOSUPPORT
	***REMOVED***
	fd, err = socket(domain, typ, proto)
	return
***REMOVED***

func Socketpair(domain, typ, proto int) (fd [2]int, err error) ***REMOVED***
	var fdx [2]int32
	err = socketpair(domain, typ, proto, &fdx)
	if err == nil ***REMOVED***
		fd[0] = int(fdx[0])
		fd[1] = int(fdx[1])
	***REMOVED***
	return
***REMOVED***

var ioSync int64

func CloseOnExec(fd int) ***REMOVED*** fcntl(fd, F_SETFD, FD_CLOEXEC) ***REMOVED***

func SetNonblock(fd int, nonblocking bool) (err error) ***REMOVED***
	flag, err := fcntl(fd, F_GETFL, 0)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if nonblocking ***REMOVED***
		flag |= O_NONBLOCK
	***REMOVED*** else ***REMOVED***
		flag &= ^O_NONBLOCK
	***REMOVED***
	_, err = fcntl(fd, F_SETFL, flag)
	return err
***REMOVED***

// Exec calls execve(2), which replaces the calling executable in the process
// tree. argv0 should be the full path to an executable ("/bin/ls") and the
// executable name should also be the first argument in argv (["ls", "-l"]).
// envv are the environment variables that should be passed to the new
// process (["USER=go", "PWD=/tmp"]).
func Exec(argv0 string, argv []string, envv []string) error ***REMOVED***
	return syscall.Exec(argv0, argv, envv)
***REMOVED***

// Lutimes sets the access and modification times tv on path. If path refers to
// a symlink, it is not dereferenced and the timestamps are set on the symlink.
// If tv is nil, the access and modification times are set to the current time.
// Otherwise tv must contain exactly 2 elements, with access time as the first
// element and modification time as the second element.
func Lutimes(path string, tv []Timeval) error ***REMOVED***
	if tv == nil ***REMOVED***
		return UtimesNanoAt(AT_FDCWD, path, nil, AT_SYMLINK_NOFOLLOW)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	ts := []Timespec***REMOVED***
		NsecToTimespec(TimevalToNsec(tv[0])),
		NsecToTimespec(TimevalToNsec(tv[1])),
	***REMOVED***
	return UtimesNanoAt(AT_FDCWD, path, ts, AT_SYMLINK_NOFOLLOW)
***REMOVED***
