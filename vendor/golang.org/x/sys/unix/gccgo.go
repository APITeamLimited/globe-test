// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build gccgo
// +build !aix

package unix

import "syscall"

// We can't use the gc-syntax .s files for gccgo. On the plus side
// much of the functionality can be written directly in Go.

//extern gccgoRealSyscallNoError
func realSyscallNoError(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r uintptr)

//extern gccgoRealSyscall
func realSyscall(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r, errno uintptr)

func SyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr) ***REMOVED***
	syscall.Entersyscall()
	r := realSyscallNoError(trap, a1, a2, a3, 0, 0, 0, 0, 0, 0)
	syscall.Exitsyscall()
	return r, 0
***REMOVED***

func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) ***REMOVED***
	syscall.Entersyscall()
	r, errno := realSyscall(trap, a1, a2, a3, 0, 0, 0, 0, 0, 0)
	syscall.Exitsyscall()
	return r, 0, syscall.Errno(errno)
***REMOVED***

func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno) ***REMOVED***
	syscall.Entersyscall()
	r, errno := realSyscall(trap, a1, a2, a3, a4, a5, a6, 0, 0, 0)
	syscall.Exitsyscall()
	return r, 0, syscall.Errno(errno)
***REMOVED***

func Syscall9(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r1, r2 uintptr, err syscall.Errno) ***REMOVED***
	syscall.Entersyscall()
	r, errno := realSyscall(trap, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	syscall.Exitsyscall()
	return r, 0, syscall.Errno(errno)
***REMOVED***

func RawSyscallNoError(trap, a1, a2, a3 uintptr) (r1, r2 uintptr) ***REMOVED***
	r := realSyscallNoError(trap, a1, a2, a3, 0, 0, 0, 0, 0, 0)
	return r, 0
***REMOVED***

func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) ***REMOVED***
	r, errno := realSyscall(trap, a1, a2, a3, 0, 0, 0, 0, 0, 0)
	return r, 0, syscall.Errno(errno)
***REMOVED***

func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno) ***REMOVED***
	r, errno := realSyscall(trap, a1, a2, a3, a4, a5, a6, 0, 0, 0)
	return r, 0, syscall.Errno(errno)
***REMOVED***
