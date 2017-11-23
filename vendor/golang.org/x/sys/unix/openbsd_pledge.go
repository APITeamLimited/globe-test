// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build openbsd
// +build 386 amd64 arm

package unix

import (
	"syscall"
	"unsafe"
)

const (
	SYS_PLEDGE = 108
)

// Pledge implements the pledge syscall. For more information see pledge(2).
func Pledge(promises string, paths []string) error ***REMOVED***
	promisesPtr, err := syscall.BytePtrFromString(promises)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	promisesUnsafe, pathsUnsafe := unsafe.Pointer(promisesPtr), unsafe.Pointer(nil)
	if paths != nil ***REMOVED***
		var pathsPtr []*byte
		if pathsPtr, err = syscall.SlicePtrFromStrings(paths); err != nil ***REMOVED***
			return err
		***REMOVED***
		pathsUnsafe = unsafe.Pointer(&pathsPtr[0])
	***REMOVED***
	_, _, e := syscall.Syscall(SYS_PLEDGE, uintptr(promisesUnsafe), uintptr(pathsUnsafe), 0)
	if e != 0 ***REMOVED***
		return e
	***REMOVED***
	return nil
***REMOVED***
