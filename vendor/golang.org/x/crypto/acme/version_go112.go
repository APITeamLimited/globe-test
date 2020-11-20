// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.12

package acme

import "runtime/debug"

func init() ***REMOVED***
	// Set packageVersion if the binary was built in modules mode and x/crypto
	// was not replaced with a different module.
	info, ok := debug.ReadBuildInfo()
	if !ok ***REMOVED***
		return
	***REMOVED***
	for _, m := range info.Deps ***REMOVED***
		if m.Path != "golang.org/x/crypto" ***REMOVED***
			continue
		***REMOVED***
		if m.Replace == nil ***REMOVED***
			packageVersion = m.Version
		***REMOVED***
		break
	***REMOVED***
***REMOVED***
