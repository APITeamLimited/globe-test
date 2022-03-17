// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (darwin && !ios) || linux
// +build darwin,!ios linux

package unix

import (
	"unsafe"

	"golang.org/x/sys/internal/unsafeheader"
)

// SysvShmAttach attaches the Sysv shared memory segment associated with the
// shared memory identifier id.
func SysvShmAttach(id int, addr uintptr, flag int) ([]byte, error) ***REMOVED***
	addr, errno := shmat(id, addr, flag)
	if errno != nil ***REMOVED***
		return nil, errno
	***REMOVED***

	// Retrieve the size of the shared memory to enable slice creation
	var info SysvShmDesc

	_, err := SysvShmCtl(id, IPC_STAT, &info)
	if err != nil ***REMOVED***
		// release the shared memory if we can't find the size

		// ignoring error from shmdt as there's nothing sensible to return here
		shmdt(addr)
		return nil, err
	***REMOVED***

	// Use unsafe to convert addr into a []byte.
	// TODO: convert to unsafe.Slice once we can assume Go 1.17
	var b []byte
	hdr := (*unsafeheader.Slice)(unsafe.Pointer(&b))
	hdr.Data = unsafe.Pointer(addr)
	hdr.Cap = int(info.Segsz)
	hdr.Len = int(info.Segsz)
	return b, nil
***REMOVED***

// SysvShmDetach unmaps the shared memory slice returned from SysvShmAttach.
//
// It is not safe to use the slice after calling this function.
func SysvShmDetach(data []byte) error ***REMOVED***
	if len(data) == 0 ***REMOVED***
		return EINVAL
	***REMOVED***

	return shmdt(uintptr(unsafe.Pointer(&data[0])))
***REMOVED***

// SysvShmGet returns the Sysv shared memory identifier associated with key.
// If the IPC_CREAT flag is specified a new segment is created.
func SysvShmGet(key, size, flag int) (id int, err error) ***REMOVED***
	return shmget(key, size, flag)
***REMOVED***
