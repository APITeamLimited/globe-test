// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (darwin && race) || (linux && race) || (freebsd && race)
// +build darwin,race linux,race freebsd,race

package unix

import (
	"runtime"
	"unsafe"
)

const raceenabled = true

func raceAcquire(addr unsafe.Pointer) ***REMOVED***
	runtime.RaceAcquire(addr)
***REMOVED***

func raceReleaseMerge(addr unsafe.Pointer) ***REMOVED***
	runtime.RaceReleaseMerge(addr)
***REMOVED***

func raceReadRange(addr unsafe.Pointer, len int) ***REMOVED***
	runtime.RaceReadRange(addr, len)
***REMOVED***

func raceWriteRange(addr unsafe.Pointer, len int) ***REMOVED***
	runtime.RaceWriteRange(addr, len)
***REMOVED***
