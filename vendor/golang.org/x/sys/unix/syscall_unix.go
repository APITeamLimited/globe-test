// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package unix

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

var (
	Stdin  = 0
	Stdout = 1
	Stderr = 2
)

const (
	darwin64Bit    = runtime.GOOS == "darwin" && sizeofPtr == 8
	dragonfly64Bit = runtime.GOOS == "dragonfly" && sizeofPtr == 8
	netbsd32Bit    = runtime.GOOS == "netbsd" && sizeofPtr == 4
	solaris64Bit   = runtime.GOOS == "solaris" && sizeofPtr == 8
)

// Do the interface allocations only once for common
// Errno values.
var (
	errEAGAIN error = syscall.EAGAIN
	errEINVAL error = syscall.EINVAL
	errENOENT error = syscall.ENOENT
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

	// Slice memory layout
	var sl = struct ***REMOVED***
		addr uintptr
		len  int
		cap  int
	***REMOVED******REMOVED***addr, length, length***REMOVED***

	// Use unsafe to turn sl into a []byte.
	b := *(*[]byte)(unsafe.Pointer(&sl))

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

type Sockaddr interface ***REMOVED***
	sockaddr() (ptr unsafe.Pointer, len _Socklen, err error) // lowercase; only we can define Sockaddrs
***REMOVED***

type SockaddrInet4 struct ***REMOVED***
	Port int
	Addr [4]byte
	raw  RawSockaddrInet4
***REMOVED***

type SockaddrInet6 struct ***REMOVED***
	Port   int
	ZoneId uint32
	Addr   [16]byte
	raw    RawSockaddrInet6
***REMOVED***

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
	return anyToSockaddr(&rsa)
***REMOVED***

func GetsockoptInt(fd, level, opt int) (value int, err error) ***REMOVED***
	var n int32
	vallen := _Socklen(4)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return int(n), err
***REMOVED***

func Recvfrom(fd int, p []byte, flags int) (n int, from Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if n, err = recvfrom(fd, p, flags, &rsa, &len); err != nil ***REMOVED***
		return
	***REMOVED***
	if rsa.Addr.Family != AF_UNSPEC ***REMOVED***
		from, err = anyToSockaddr(&rsa)
	***REMOVED***
	return
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
	return setsockopt(fd, level, opt, unsafe.Pointer(&[]byte(s)[0]), uintptr(len(s)))
***REMOVED***

func SetsockoptTimeval(fd, level, opt int, tv *Timeval) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(tv), unsafe.Sizeof(*tv))
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

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) ***REMOVED***
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	return sendfile(outfd, infd, offset, count)
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
