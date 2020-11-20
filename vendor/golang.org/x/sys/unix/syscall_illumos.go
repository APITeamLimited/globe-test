// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// illumos system calls not present on Solaris.

// +build amd64,illumos

package unix

import "unsafe"

func bytes2iovec(bs [][]byte) []Iovec ***REMOVED***
	iovecs := make([]Iovec, len(bs))
	for i, b := range bs ***REMOVED***
		iovecs[i].SetLen(len(b))
		if len(b) > 0 ***REMOVED***
			// somehow Iovec.Base on illumos is (*int8), not (*byte)
			iovecs[i].Base = (*int8)(unsafe.Pointer(&b[0]))
		***REMOVED*** else ***REMOVED***
			iovecs[i].Base = (*int8)(unsafe.Pointer(&_zero))
		***REMOVED***
	***REMOVED***
	return iovecs
***REMOVED***

//sys	readv(fd int, iovs []Iovec) (n int, err error)

func Readv(fd int, iovs [][]byte) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = readv(fd, iovecs)
	return n, err
***REMOVED***

//sys	preadv(fd int, iovs []Iovec, off int64) (n int, err error)

func Preadv(fd int, iovs [][]byte, off int64) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = preadv(fd, iovecs, off)
	return n, err
***REMOVED***

//sys	writev(fd int, iovs []Iovec) (n int, err error)

func Writev(fd int, iovs [][]byte) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = writev(fd, iovecs)
	return n, err
***REMOVED***

//sys	pwritev(fd int, iovs []Iovec, off int64) (n int, err error)

func Pwritev(fd int, iovs [][]byte, off int64) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = pwritev(fd, iovecs, off)
	return n, err
***REMOVED***

//sys	accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error) = libsocket.accept4

func Accept4(fd int, flags int) (nfd int, sa Sockaddr, err error) ***REMOVED***
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
