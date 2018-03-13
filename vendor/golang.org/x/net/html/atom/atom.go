// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package atom provides integer codes (also known as atoms) for a fixed set of
// frequently occurring HTML strings: tag names and attribute keys such as "p"
// and "id".
//
// Sharing an atom's name between all elements with the same tag can result in
// fewer string allocations when tokenizing and parsing HTML. Integer
// comparisons are also generally faster than string comparisons.
//
// The value of an atom's particular code is not guaranteed to stay the same
// between versions of this package. Neither is any ordering guaranteed:
// whether atom.H1 < atom.H2 may also change. The codes are not guaranteed to
// be dense. The only guarantees are that e.g. looking up "div" will yield
// atom.Div, calling atom.Div.String will return "div", and atom.Div != 0.
package atom // import "golang.org/x/net/html/atom"

// Atom is an integer code for a string. The zero value maps to "".
type Atom uint32

// String returns the atom's name.
func (a Atom) String() string ***REMOVED***
	start := uint32(a >> 8)
	n := uint32(a & 0xff)
	if start+n > uint32(len(atomText)) ***REMOVED***
		return ""
	***REMOVED***
	return atomText[start : start+n]
***REMOVED***

func (a Atom) string() string ***REMOVED***
	return atomText[a>>8 : a>>8+a&0xff]
***REMOVED***

// fnv computes the FNV hash with an arbitrary starting value h.
func fnv(h uint32, s []byte) uint32 ***REMOVED***
	for i := range s ***REMOVED***
		h ^= uint32(s[i])
		h *= 16777619
	***REMOVED***
	return h
***REMOVED***

func match(s string, t []byte) bool ***REMOVED***
	for i, c := range t ***REMOVED***
		if s[i] != c ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// Lookup returns the atom whose name is s. It returns zero if there is no
// such atom. The lookup is case sensitive.
func Lookup(s []byte) Atom ***REMOVED***
	if len(s) == 0 || len(s) > maxAtomLen ***REMOVED***
		return 0
	***REMOVED***
	h := fnv(hash0, s)
	if a := table[h&uint32(len(table)-1)]; int(a&0xff) == len(s) && match(a.string(), s) ***REMOVED***
		return a
	***REMOVED***
	if a := table[(h>>16)&uint32(len(table)-1)]; int(a&0xff) == len(s) && match(a.string(), s) ***REMOVED***
		return a
	***REMOVED***
	return 0
***REMOVED***

// String returns a string whose contents are equal to s. In that sense, it is
// equivalent to string(s) but may be more efficient.
func String(s []byte) string ***REMOVED***
	if a := Lookup(s); a != 0 ***REMOVED***
		return a.String()
	***REMOVED***
	return string(s)
***REMOVED***
