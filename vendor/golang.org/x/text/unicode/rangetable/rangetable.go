// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rangetable provides utilities for creating and inspecting
// unicode.RangeTables.
package rangetable

import (
	"sort"
	"unicode"
)

// New creates a RangeTable from the given runes, which may contain duplicates.
func New(r ...rune) *unicode.RangeTable ***REMOVED***
	if len(r) == 0 ***REMOVED***
		return &unicode.RangeTable***REMOVED******REMOVED***
	***REMOVED***

	sort.Sort(byRune(r))

	// Remove duplicates.
	k := 1
	for i := 1; i < len(r); i++ ***REMOVED***
		if r[k-1] != r[i] ***REMOVED***
			r[k] = r[i]
			k++
		***REMOVED***
	***REMOVED***

	var rt unicode.RangeTable
	for _, r := range r[:k] ***REMOVED***
		if r <= 0xFFFF ***REMOVED***
			rt.R16 = append(rt.R16, unicode.Range16***REMOVED***Lo: uint16(r), Hi: uint16(r), Stride: 1***REMOVED***)
		***REMOVED*** else ***REMOVED***
			rt.R32 = append(rt.R32, unicode.Range32***REMOVED***Lo: uint32(r), Hi: uint32(r), Stride: 1***REMOVED***)
		***REMOVED***
	***REMOVED***

	// Optimize RangeTable.
	return Merge(&rt)
***REMOVED***

type byRune []rune

func (r byRune) Len() int           ***REMOVED*** return len(r) ***REMOVED***
func (r byRune) Swap(i, j int)      ***REMOVED*** r[i], r[j] = r[j], r[i] ***REMOVED***
func (r byRune) Less(i, j int) bool ***REMOVED*** return r[i] < r[j] ***REMOVED***

// Visit visits all runes in the given RangeTable in order, calling fn for each.
func Visit(rt *unicode.RangeTable, fn func(rune)) ***REMOVED***
	for _, r16 := range rt.R16 ***REMOVED***
		for r := rune(r16.Lo); r <= rune(r16.Hi); r += rune(r16.Stride) ***REMOVED***
			fn(r)
		***REMOVED***
	***REMOVED***
	for _, r32 := range rt.R32 ***REMOVED***
		for r := rune(r32.Lo); r <= rune(r32.Hi); r += rune(r32.Stride) ***REMOVED***
			fn(r)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Assigned returns a RangeTable with all assigned code points for a given
// Unicode version. This includes graphic, format, control, and private-use
// characters. It returns nil if the data for the given version is not
// available.
func Assigned(version string) *unicode.RangeTable ***REMOVED***
	return assigned[version]
***REMOVED***
