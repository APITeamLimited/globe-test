// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package unix

// Set adds fd to the set fds.
func (fds *FdSet) Set(fd int) ***REMOVED***
	fds.Bits[fd/NFDBITS] |= (1 << (uintptr(fd) % NFDBITS))
***REMOVED***

// Clear removes fd from the set fds.
func (fds *FdSet) Clear(fd int) ***REMOVED***
	fds.Bits[fd/NFDBITS] &^= (1 << (uintptr(fd) % NFDBITS))
***REMOVED***

// IsSet returns whether fd is in the set fds.
func (fds *FdSet) IsSet(fd int) bool ***REMOVED***
	return fds.Bits[fd/NFDBITS]&(1<<(uintptr(fd)%NFDBITS)) != 0
***REMOVED***

// Zero clears the set fds.
func (fds *FdSet) Zero() ***REMOVED***
	for i := range fds.Bits ***REMOVED***
		fds.Bits[i] = 0
	***REMOVED***
***REMOVED***
