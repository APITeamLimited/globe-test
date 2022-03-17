// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.5
// +build !go1.5

package plan9

func fixwd() ***REMOVED***
***REMOVED***

func Getwd() (wd string, err error) ***REMOVED***
	fd, err := open(".", O_RDONLY)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	defer Close(fd)
	return Fd2path(fd)
***REMOVED***

func Chdir(path string) error ***REMOVED***
	return chdir(path)
***REMOVED***
