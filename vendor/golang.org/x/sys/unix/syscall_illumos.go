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

//sys   readv(fd int, iovs []Iovec) (n int, err error)

func Readv(fd int, iovs [][]byte) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = readv(fd, iovecs)
	return n, err
***REMOVED***

//sys   preadv(fd int, iovs []Iovec, off int64) (n int, err error)

func Preadv(fd int, iovs [][]byte, off int64) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = preadv(fd, iovecs, off)
	return n, err
***REMOVED***

//sys   writev(fd int, iovs []Iovec) (n int, err error)

func Writev(fd int, iovs [][]byte) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = writev(fd, iovecs)
	return n, err
***REMOVED***

//sys   pwritev(fd int, iovs []Iovec, off int64) (n int, err error)

func Pwritev(fd int, iovs [][]byte, off int64) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = pwritev(fd, iovecs, off)
	return n, err
***REMOVED***
