// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Extensions to the standard "os" package.
package osext // import "github.com/kardianos/osext"

import "path/filepath"

var cx, ce = executableClean()

func executableClean() (string, error) ***REMOVED***
	p, err := executable()
	return filepath.Clean(p), err
***REMOVED***

// Executable returns an absolute path that can be used to
// re-invoke the current program.
// It may not be valid after the current program exits.
func Executable() (string, error) ***REMOVED***
	return cx, ce
***REMOVED***

// Returns same path as Executable, returns just the folder
// path. Excludes the executable name and any trailing slash.
func ExecutableFolder() (string, error) ***REMOVED***
	p, err := Executable()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return filepath.Dir(p), nil
***REMOVED***
