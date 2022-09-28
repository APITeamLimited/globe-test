// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Windows environment variables.

package windows

import (
	"syscall"
	"unsafe"
)

func Getenv(key string) (value string, found bool) ***REMOVED***
	return syscall.Getenv(key)
***REMOVED***

func Setenv(key, value string) error ***REMOVED***
	return syscall.Setenv(key, value)
***REMOVED***

func Clearenv() ***REMOVED***
	syscall.Clearenv()
***REMOVED***

func Environ() []string ***REMOVED***
	return syscall.Environ()
***REMOVED***

// Returns a default environment associated with the token, rather than the current
// process. If inheritExisting is true, then this environment also inherits the
// environment of the current process.
func (token Token) Environ(inheritExisting bool) (env []string, err error) ***REMOVED***
	var block *uint16
	err = CreateEnvironmentBlock(&block, token, inheritExisting)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer DestroyEnvironmentBlock(block)
	blockp := uintptr(unsafe.Pointer(block))
	for ***REMOVED***
		entry := UTF16PtrToString((*uint16)(unsafe.Pointer(blockp)))
		if len(entry) == 0 ***REMOVED***
			break
		***REMOVED***
		env = append(env, entry)
		blockp += 2 * (uintptr(len(entry)) + 1)
	***REMOVED***
	return env, nil
***REMOVED***

func Unsetenv(key string) error ***REMOVED***
	return syscall.Unsetenv(key)
***REMOVED***
