// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Plan 9 environment variables.

package plan9

import (
	"syscall"
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

func Unsetenv(key string) error ***REMOVED***
	return syscall.Unsetenv(key)
***REMOVED***
