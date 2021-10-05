// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package unix

import "unsafe"

// readInt returns the size-bytes unsigned integer in native byte order at offset off.
func readInt(b []byte, off, size uintptr) (u uint64, ok bool) ***REMOVED***
	if len(b) < int(off+size) ***REMOVED***
		return 0, false
	***REMOVED***
	if isBigEndian ***REMOVED***
		return readIntBE(b[off:], size), true
	***REMOVED***
	return readIntLE(b[off:], size), true
***REMOVED***

func readIntBE(b []byte, size uintptr) uint64 ***REMOVED***
	switch size ***REMOVED***
	case 1:
		return uint64(b[0])
	case 2:
		_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[1]) | uint64(b[0])<<8
	case 4:
		_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[3]) | uint64(b[2])<<8 | uint64(b[1])<<16 | uint64(b[0])<<24
	case 8:
		_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
			uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	default:
		panic("syscall: readInt with unsupported size")
	***REMOVED***
***REMOVED***

func readIntLE(b []byte, size uintptr) uint64 ***REMOVED***
	switch size ***REMOVED***
	case 1:
		return uint64(b[0])
	case 2:
		_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8
	case 4:
		_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24
	case 8:
		_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
			uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	default:
		panic("syscall: readInt with unsupported size")
	***REMOVED***
***REMOVED***

// ParseDirent parses up to max directory entries in buf,
// appending the names to names. It returns the number of
// bytes consumed from buf, the number of entries added
// to names, and the new names slice.
func ParseDirent(buf []byte, max int, names []string) (consumed int, count int, newnames []string) ***REMOVED***
	origlen := len(buf)
	count = 0
	for max != 0 && len(buf) > 0 ***REMOVED***
		reclen, ok := direntReclen(buf)
		if !ok || reclen > uint64(len(buf)) ***REMOVED***
			return origlen, count, names
		***REMOVED***
		rec := buf[:reclen]
		buf = buf[reclen:]
		ino, ok := direntIno(rec)
		if !ok ***REMOVED***
			break
		***REMOVED***
		if ino == 0 ***REMOVED*** // File absent in directory.
			continue
		***REMOVED***
		const namoff = uint64(unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Name))
		namlen, ok := direntNamlen(rec)
		if !ok || namoff+namlen > uint64(len(rec)) ***REMOVED***
			break
		***REMOVED***
		name := rec[namoff : namoff+namlen]
		for i, c := range name ***REMOVED***
			if c == 0 ***REMOVED***
				name = name[:i]
				break
			***REMOVED***
		***REMOVED***
		// Check for useless names before allocating a string.
		if string(name) == "." || string(name) == ".." ***REMOVED***
			continue
		***REMOVED***
		max--
		count++
		names = append(names, string(name))
	***REMOVED***
	return origlen - len(buf), count, names
***REMOVED***
