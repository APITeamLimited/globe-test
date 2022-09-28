// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

// MustParse is like Parse, but panics if the given BCP 47 tag cannot be parsed.
// It simplifies safe initialization of Tag values.
func MustParse(s string) Tag ***REMOVED***
	t, err := Parse(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return t
***REMOVED***

// MustParseBase is like ParseBase, but panics if the given base cannot be parsed.
// It simplifies safe initialization of Base values.
func MustParseBase(s string) Language ***REMOVED***
	b, err := ParseBase(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return b
***REMOVED***

// MustParseScript is like ParseScript, but panics if the given script cannot be
// parsed. It simplifies safe initialization of Script values.
func MustParseScript(s string) Script ***REMOVED***
	scr, err := ParseScript(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return scr
***REMOVED***

// MustParseRegion is like ParseRegion, but panics if the given region cannot be
// parsed. It simplifies safe initialization of Region values.
func MustParseRegion(s string) Region ***REMOVED***
	r, err := ParseRegion(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return r
***REMOVED***

// Und is the root language.
var Und Tag
