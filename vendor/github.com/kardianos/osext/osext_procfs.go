// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.8,android !go1.8,linux !go1.8,netbsd !go1.8,solaris !go1.8,dragonfly

package osext

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func executable() (string, error) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "linux", "android":
		const deletedTag = " (deleted)"
		execpath, err := os.Readlink("/proc/self/exe")
		if err != nil ***REMOVED***
			return execpath, err
		***REMOVED***
		execpath = strings.TrimSuffix(execpath, deletedTag)
		execpath = strings.TrimPrefix(execpath, deletedTag)
		return execpath, nil
	case "netbsd":
		return os.Readlink("/proc/curproc/exe")
	case "dragonfly":
		return os.Readlink("/proc/curproc/file")
	case "solaris":
		return os.Readlink(fmt.Sprintf("/proc/%d/path/a.out", os.Getpid()))
	***REMOVED***
	return "", errors.New("ExecPath not implemented for " + runtime.GOOS)
***REMOVED***
