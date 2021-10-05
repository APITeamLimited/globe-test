// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.5

package plan9

import "syscall"

func fixwd() ***REMOVED***
	syscall.Fixwd()
***REMOVED***

func Getwd() (wd string, err error) ***REMOVED***
	return syscall.Getwd()
***REMOVED***

func Chdir(path string) error ***REMOVED***
	return syscall.Chdir(path)
***REMOVED***
