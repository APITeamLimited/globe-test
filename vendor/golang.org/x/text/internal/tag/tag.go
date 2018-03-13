// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tag contains functionality handling tags and related data.
package tag // import "golang.org/x/text/internal/tag"

import "sort"

// An Index converts tags to a compact numeric value.
//
// All elements are of size 4. Tags may be up to 4 bytes long. Excess bytes can
// be used to store additional information about the tag.
type Index string

// Elem returns the element data at the given index.
func (s Index) Elem(x int) string ***REMOVED***
	return string(s[x*4 : x*4+4])
***REMOVED***

// Index reports the index of the given key or -1 if it could not be found.
// Only the first len(key) bytes from the start of the 4-byte entries will be
// considered for the search and the first match in Index will be returned.
func (s Index) Index(key []byte) int ***REMOVED***
	n := len(key)
	// search the index of the first entry with an equal or higher value than
	// key in s.
	index := sort.Search(len(s)/4, func(i int) bool ***REMOVED***
		return cmp(s[i*4:i*4+n], key) != -1
	***REMOVED***)
	i := index * 4
	if cmp(s[i:i+len(key)], key) != 0 ***REMOVED***
		return -1
	***REMOVED***
	return index
***REMOVED***

// Next finds the next occurrence of key after index x, which must have been
// obtained from a call to Index using the same key. It returns x+1 or -1.
func (s Index) Next(key []byte, x int) int ***REMOVED***
	if x++; x*4 < len(s) && cmp(s[x*4:x*4+len(key)], key) == 0 ***REMOVED***
		return x
	***REMOVED***
	return -1
***REMOVED***

// cmp returns an integer comparing a and b lexicographically.
func cmp(a Index, b []byte) int ***REMOVED***
	n := len(a)
	if len(b) < n ***REMOVED***
		n = len(b)
	***REMOVED***
	for i, c := range b[:n] ***REMOVED***
		switch ***REMOVED***
		case a[i] > c:
			return 1
		case a[i] < c:
			return -1
		***REMOVED***
	***REMOVED***
	switch ***REMOVED***
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	***REMOVED***
	return 0
***REMOVED***

// Compare returns an integer comparing a and b lexicographically.
func Compare(a string, b []byte) int ***REMOVED***
	return cmp(Index(a), b)
***REMOVED***

// FixCase reformats b to the same pattern of cases as form.
// If returns false if string b is malformed.
func FixCase(form string, b []byte) bool ***REMOVED***
	if len(form) != len(b) ***REMOVED***
		return false
	***REMOVED***
	for i, c := range b ***REMOVED***
		if form[i] <= 'Z' ***REMOVED***
			if c >= 'a' ***REMOVED***
				c -= 'z' - 'Z'
			***REMOVED***
			if c < 'A' || 'Z' < c ***REMOVED***
				return false
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if c <= 'Z' ***REMOVED***
				c += 'z' - 'Z'
			***REMOVED***
			if c < 'a' || 'z' < c ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		b[i] = c
	***REMOVED***
	return true
***REMOVED***
