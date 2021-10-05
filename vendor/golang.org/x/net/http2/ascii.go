// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import "strings"

// asciiEqualFold is strings.EqualFold, ASCII only. It reports whether s and t
// are equal, ASCII-case-insensitively.
func asciiEqualFold(s, t string) bool ***REMOVED***
	if len(s) != len(t) ***REMOVED***
		return false
	***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		if lower(s[i]) != lower(t[i]) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// lower returns the ASCII lowercase version of b.
func lower(b byte) byte ***REMOVED***
	if 'A' <= b && b <= 'Z' ***REMOVED***
		return b + ('a' - 'A')
	***REMOVED***
	return b
***REMOVED***

// isASCIIPrint returns whether s is ASCII and printable according to
// https://tools.ietf.org/html/rfc20#section-4.2.
func isASCIIPrint(s string) bool ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] < ' ' || s[i] > '~' ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// asciiToLower returns the lowercase version of s if s is ASCII and printable,
// and whether or not it was.
func asciiToLower(s string) (lower string, ok bool) ***REMOVED***
	if !isASCIIPrint(s) ***REMOVED***
		return "", false
	***REMOVED***
	return strings.ToLower(s), true
***REMOVED***
