// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// illumos system calls not present on Solaris.

//go:build amd64 && illumos
// +build amd64,illumos

package unix

import (
	"fmt"
	"runtime"
	"unsafe"
)

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

//sys	putmsg(fd int, clptr *strbuf, dataptr *strbuf, flags int) (err error)

func Putmsg(fd int, cl []byte, data []byte, flags int) (err error) ***REMOVED***
	var clp, datap *strbuf
	if len(cl) > 0 ***REMOVED***
		clp = &strbuf***REMOVED***
			Len: int32(len(cl)),
			Buf: (*int8)(unsafe.Pointer(&cl[0])),
		***REMOVED***
	***REMOVED***
	if len(data) > 0 ***REMOVED***
		datap = &strbuf***REMOVED***
			Len: int32(len(data)),
			Buf: (*int8)(unsafe.Pointer(&data[0])),
		***REMOVED***
	***REMOVED***
	return putmsg(fd, clp, datap, flags)
***REMOVED***

//sys	getmsg(fd int, clptr *strbuf, dataptr *strbuf, flags *int) (err error)

func Getmsg(fd int, cl []byte, data []byte) (retCl []byte, retData []byte, flags int, err error) ***REMOVED***
	var clp, datap *strbuf
	if len(cl) > 0 ***REMOVED***
		clp = &strbuf***REMOVED***
			Maxlen: int32(len(cl)),
			Buf:    (*int8)(unsafe.Pointer(&cl[0])),
		***REMOVED***
	***REMOVED***
	if len(data) > 0 ***REMOVED***
		datap = &strbuf***REMOVED***
			Maxlen: int32(len(data)),
			Buf:    (*int8)(unsafe.Pointer(&data[0])),
		***REMOVED***
	***REMOVED***

	if err = getmsg(fd, clp, datap, &flags); err != nil ***REMOVED***
		return nil, nil, 0, err
	***REMOVED***

	if len(cl) > 0 ***REMOVED***
		retCl = cl[:clp.Len]
	***REMOVED***
	if len(data) > 0 ***REMOVED***
		retData = data[:datap.Len]
	***REMOVED***
	return retCl, retData, flags, nil
***REMOVED***

func IoctlSetIntRetInt(fd int, req uint, arg int) (int, error) ***REMOVED***
	return ioctlRet(fd, req, uintptr(arg))
***REMOVED***

func IoctlSetString(fd int, req uint, val string) error ***REMOVED***
	bs := make([]byte, len(val)+1)
	copy(bs[:len(bs)-1], val)
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&bs[0])))
	runtime.KeepAlive(&bs[0])
	return err
***REMOVED***

// Lifreq Helpers

func (l *Lifreq) SetName(name string) error ***REMOVED***
	if len(name) >= len(l.Name) ***REMOVED***
		return fmt.Errorf("name cannot be more than %d characters", len(l.Name)-1)
	***REMOVED***
	for i := range name ***REMOVED***
		l.Name[i] = int8(name[i])
	***REMOVED***
	return nil
***REMOVED***

func (l *Lifreq) SetLifruInt(d int) ***REMOVED***
	*(*int)(unsafe.Pointer(&l.Lifru[0])) = d
***REMOVED***

func (l *Lifreq) GetLifruInt() int ***REMOVED***
	return *(*int)(unsafe.Pointer(&l.Lifru[0]))
***REMOVED***

func IoctlLifreq(fd int, req uint, l *Lifreq) error ***REMOVED***
	return ioctl(fd, req, uintptr(unsafe.Pointer(l)))
***REMOVED***

// Strioctl Helpers

func (s *Strioctl) SetInt(i int) ***REMOVED***
	s.Len = int32(unsafe.Sizeof(i))
	s.Dp = (*int8)(unsafe.Pointer(&i))
***REMOVED***

func IoctlSetStrioctlRetInt(fd int, req uint, s *Strioctl) (int, error) ***REMOVED***
	return ioctlRet(fd, req, uintptr(unsafe.Pointer(s)))
***REMOVED***
