// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

// BaseLanguages returns the list of all supported base languages. It generates
// the list by traversing the internal structures.
func BaseLanguages() []Language ***REMOVED***
	base := make([]Language, 0, NumLanguages)
	for i := 0; i < langNoIndexOffset; i++ ***REMOVED***
		// We included "und" already for the value 0.
		if i != nonCanonicalUnd ***REMOVED***
			base = append(base, Language(i))
		***REMOVED***
	***REMOVED***
	i := langNoIndexOffset
	for _, v := range langNoIndex ***REMOVED***
		for k := 0; k < 8; k++ ***REMOVED***
			if v&1 == 1 ***REMOVED***
				base = append(base, Language(i))
			***REMOVED***
			v >>= 1
			i++
		***REMOVED***
	***REMOVED***
	return base
***REMOVED***
