// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || (darwin && !race) || (linux && !race) || (freebsd && !race) || netbsd || openbsd || solaris || dragonfly || zos
// +build aix darwin,!race linux,!race freebsd,!race netbsd openbsd solaris dragonfly zos

package unix

import (
	"unsafe"
)

const raceenabled = false

func raceAcquire(addr unsafe.Pointer) ***REMOVED***
***REMOVED***

func raceReleaseMerge(addr unsafe.Pointer) ***REMOVED***
***REMOVED***

func raceReadRange(addr unsafe.Pointer, len int) ***REMOVED***
***REMOVED***

func raceWriteRange(addr unsafe.Pointer, len int) ***REMOVED***
***REMOVED***
