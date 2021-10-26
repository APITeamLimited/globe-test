// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package term

import (
	"fmt"
	"runtime"

	"golang.org/x/sys/plan9"
)

type state struct***REMOVED******REMOVED***

func isTerminal(fd int) bool ***REMOVED***
	path, err := plan9.Fd2path(fd)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return path == "/dev/cons" || path == "/mnt/term/dev/cons"
***REMOVED***

func makeRaw(fd int) (*State, error) ***REMOVED***
	return nil, fmt.Errorf("terminal: MakeRaw not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
***REMOVED***

func getState(fd int) (*State, error) ***REMOVED***
	return nil, fmt.Errorf("terminal: GetState not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
***REMOVED***

func restore(fd int, state *State) error ***REMOVED***
	return fmt.Errorf("terminal: Restore not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
***REMOVED***

func getSize(fd int) (width, height int, err error) ***REMOVED***
	return 0, 0, fmt.Errorf("terminal: GetSize not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
***REMOVED***

func readPassword(fd int) ([]byte, error) ***REMOVED***
	return nil, fmt.Errorf("terminal: ReadPassword not implemented on %s/%s", runtime.GOOS, runtime.GOARCH)
***REMOVED***