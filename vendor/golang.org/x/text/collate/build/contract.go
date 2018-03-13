// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/text/internal/colltab"
)

// This file contains code for detecting contractions and generating
// the necessary tables.
// Any Unicode Collation Algorithm (UCA) table entry that has more than
// one rune one the left-hand side is called a contraction.
// See http://www.unicode.org/reports/tr10/#Contractions for more details.
//
// We define the following terms:
//   initial:     a rune that appears as the first rune in a contraction.
//   suffix:      a sequence of runes succeeding the initial rune
//                in a given contraction.
//   non-initial: a rune that appears in a suffix.
//
// A rune may be both an initial and a non-initial and may be so in
// many contractions.  An initial may typically also appear by itself.
// In case of ambiguities, the UCA requires we match the longest
// contraction.
//
// Many contraction rules share the same set of possible suffixes.
// We store sets of suffixes in a trie that associates an index with
// each suffix in the set.  This index can be used to look up a
// collation element associated with the (starter rune, suffix) pair.
//
// The trie is defined on a UTF-8 byte sequence.
// The overall trie is represented as an array of ctEntries.  Each node of the trie
// is represented as a subsequence of ctEntries, where each entry corresponds to
// a possible match of a next character in the search string.  An entry
// also includes the length and offset to the next sequence of entries
// to check in case of a match.

const (
	final   = 0
	noIndex = 0xFF
)

// ctEntry associates to a matching byte an offset and/or next sequence of
// bytes to check. A ctEntry c is called final if a match means that the
// longest suffix has been found.  An entry c is final if c.N == 0.
// A single final entry can match a range of characters to an offset.
// A non-final entry always matches a single byte. Note that a non-final
// entry might still resemble a completed suffix.
// Examples:
// The suffix strings "ab" and "ac" can be represented as:
// []ctEntry***REMOVED***
//     ***REMOVED***'a', 1, 1, noIndex***REMOVED***,  // 'a' by itself does not match, so i is 0xFF.
//     ***REMOVED***'b', 'c', 0, 1***REMOVED***,   // "ab" -> 1, "ac" -> 2
// ***REMOVED***
//
// The suffix strings "ab", "abc", "abd", and "abcd" can be represented as:
// []ctEntry***REMOVED***
//     ***REMOVED***'a', 1, 1, noIndex***REMOVED***, // 'a' must be followed by 'b'.
//     ***REMOVED***'b', 1, 2, 1***REMOVED***,    // "ab" -> 1, may be followed by 'c' or 'd'.
//     ***REMOVED***'d', 'd', final, 3***REMOVED***,  // "abd" -> 3
//     ***REMOVED***'c', 4, 1, 2***REMOVED***,    // "abc" -> 2, may be followed by 'd'.
//     ***REMOVED***'d', 'd', final, 4***REMOVED***,  // "abcd" -> 4
// ***REMOVED***
// See genStateTests in contract_test.go for more examples.
type ctEntry struct ***REMOVED***
	L uint8 // non-final: byte value to match; final: lowest match in range.
	H uint8 // non-final: relative index to next block; final: highest match in range.
	N uint8 // non-final: length of next block; final: final
	I uint8 // result offset. Will be noIndex if more bytes are needed to complete.
***REMOVED***

// contractTrieSet holds a set of contraction tries. The tries are stored
// consecutively in the entry field.
type contractTrieSet []struct***REMOVED*** l, h, n, i uint8 ***REMOVED***

// ctHandle is used to identify a trie in the trie set, consisting in an offset
// in the array and the size of the first node.
type ctHandle struct ***REMOVED***
	index, n int
***REMOVED***

// appendTrie adds a new trie for the given suffixes to the trie set and returns
// a handle to it.  The handle will be invalid on error.
func appendTrie(ct *colltab.ContractTrieSet, suffixes []string) (ctHandle, error) ***REMOVED***
	es := make([]stridx, len(suffixes))
	for i, s := range suffixes ***REMOVED***
		es[i].str = s
	***REMOVED***
	sort.Sort(offsetSort(es))
	for i := range es ***REMOVED***
		es[i].index = i + 1
	***REMOVED***
	sort.Sort(genidxSort(es))
	i := len(*ct)
	n, err := genStates(ct, es)
	if err != nil ***REMOVED***
		*ct = (*ct)[:i]
		return ctHandle***REMOVED******REMOVED***, err
	***REMOVED***
	return ctHandle***REMOVED***i, n***REMOVED***, nil
***REMOVED***

// genStates generates ctEntries for a given suffix set and returns
// the number of entries for the first node.
func genStates(ct *colltab.ContractTrieSet, sis []stridx) (int, error) ***REMOVED***
	if len(sis) == 0 ***REMOVED***
		return 0, fmt.Errorf("genStates: list of suffices must be non-empty")
	***REMOVED***
	start := len(*ct)
	// create entries for differing first bytes.
	for _, si := range sis ***REMOVED***
		s := si.str
		if len(s) == 0 ***REMOVED***
			continue
		***REMOVED***
		added := false
		c := s[0]
		if len(s) > 1 ***REMOVED***
			for j := len(*ct) - 1; j >= start; j-- ***REMOVED***
				if (*ct)[j].L == c ***REMOVED***
					added = true
					break
				***REMOVED***
			***REMOVED***
			if !added ***REMOVED***
				*ct = append(*ct, ctEntry***REMOVED***L: c, I: noIndex***REMOVED***)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for j := len(*ct) - 1; j >= start; j-- ***REMOVED***
				// Update the offset for longer suffixes with the same byte.
				if (*ct)[j].L == c ***REMOVED***
					(*ct)[j].I = uint8(si.index)
					added = true
				***REMOVED***
				// Extend range of final ctEntry, if possible.
				if (*ct)[j].H+1 == c ***REMOVED***
					(*ct)[j].H = c
					added = true
				***REMOVED***
			***REMOVED***
			if !added ***REMOVED***
				*ct = append(*ct, ctEntry***REMOVED***L: c, H: c, N: final, I: uint8(si.index)***REMOVED***)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	n := len(*ct) - start
	// Append nodes for the remainder of the suffixes for each ctEntry.
	sp := 0
	for i, end := start, len(*ct); i < end; i++ ***REMOVED***
		fe := (*ct)[i]
		if fe.H == 0 ***REMOVED*** // uninitialized non-final
			ln := len(*ct) - start - n
			if ln > 0xFF ***REMOVED***
				return 0, fmt.Errorf("genStates: relative block offset too large: %d > 255", ln)
			***REMOVED***
			fe.H = uint8(ln)
			// Find first non-final strings with same byte as current entry.
			for ; sis[sp].str[0] != fe.L; sp++ ***REMOVED***
			***REMOVED***
			se := sp + 1
			for ; se < len(sis) && len(sis[se].str) > 1 && sis[se].str[0] == fe.L; se++ ***REMOVED***
			***REMOVED***
			sl := sis[sp:se]
			sp = se
			for i, si := range sl ***REMOVED***
				sl[i].str = si.str[1:]
			***REMOVED***
			nn, err := genStates(ct, sl)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			fe.N = uint8(nn)
			(*ct)[i] = fe
		***REMOVED***
	***REMOVED***
	sort.Sort(entrySort((*ct)[start : start+n]))
	return n, nil
***REMOVED***

// There may be both a final and non-final entry for a byte if the byte
// is implied in a range of matches in the final entry.
// We need to ensure that the non-final entry comes first in that case.
type entrySort colltab.ContractTrieSet

func (fe entrySort) Len() int      ***REMOVED*** return len(fe) ***REMOVED***
func (fe entrySort) Swap(i, j int) ***REMOVED*** fe[i], fe[j] = fe[j], fe[i] ***REMOVED***
func (fe entrySort) Less(i, j int) bool ***REMOVED***
	return fe[i].L > fe[j].L
***REMOVED***

// stridx is used for sorting suffixes and their associated offsets.
type stridx struct ***REMOVED***
	str   string
	index int
***REMOVED***

// For computing the offsets, we first sort by size, and then by string.
// This ensures that strings that only differ in the last byte by 1
// are sorted consecutively in increasing order such that they can
// be packed as a range in a final ctEntry.
type offsetSort []stridx

func (si offsetSort) Len() int      ***REMOVED*** return len(si) ***REMOVED***
func (si offsetSort) Swap(i, j int) ***REMOVED*** si[i], si[j] = si[j], si[i] ***REMOVED***
func (si offsetSort) Less(i, j int) bool ***REMOVED***
	if len(si[i].str) != len(si[j].str) ***REMOVED***
		return len(si[i].str) > len(si[j].str)
	***REMOVED***
	return si[i].str < si[j].str
***REMOVED***

// For indexing, we want to ensure that strings are sorted in string order, where
// for strings with the same prefix, we put longer strings before shorter ones.
type genidxSort []stridx

func (si genidxSort) Len() int      ***REMOVED*** return len(si) ***REMOVED***
func (si genidxSort) Swap(i, j int) ***REMOVED*** si[i], si[j] = si[j], si[i] ***REMOVED***
func (si genidxSort) Less(i, j int) bool ***REMOVED***
	if strings.HasPrefix(si[j].str, si[i].str) ***REMOVED***
		return false
	***REMOVED***
	if strings.HasPrefix(si[i].str, si[j].str) ***REMOVED***
		return true
	***REMOVED***
	return si[i].str < si[j].str
***REMOVED***

// lookup matches the longest suffix in str and returns the associated offset
// and the number of bytes consumed.
func lookup(ct *colltab.ContractTrieSet, h ctHandle, str []byte) (index, ns int) ***REMOVED***
	states := (*ct)[h.index:]
	p := 0
	n := h.n
	for i := 0; i < n && p < len(str); ***REMOVED***
		e := states[i]
		c := str[p]
		if c >= e.L ***REMOVED***
			if e.L == c ***REMOVED***
				p++
				if e.I != noIndex ***REMOVED***
					index, ns = int(e.I), p
				***REMOVED***
				if e.N != final ***REMOVED***
					// set to new state
					i, states, n = 0, states[int(e.H)+n:], int(e.N)
				***REMOVED*** else ***REMOVED***
					return
				***REMOVED***
				continue
			***REMOVED*** else if e.N == final && c <= e.H ***REMOVED***
				p++
				return int(c-e.L) + int(e.I), p
			***REMOVED***
		***REMOVED***
		i++
	***REMOVED***
	return
***REMOVED***

// print writes the contractTrieSet t as compilable Go code to w. It returns
// the total number of bytes written and the size of the resulting data structure in bytes.
func print(t *colltab.ContractTrieSet, w io.Writer, name string) (n, size int, err error) ***REMOVED***
	update3 := func(nn, sz int, e error) ***REMOVED***
		n += nn
		if err == nil ***REMOVED***
			err = e
		***REMOVED***
		size += sz
	***REMOVED***
	update2 := func(nn int, e error) ***REMOVED*** update3(nn, 0, e) ***REMOVED***

	update3(printArray(*t, w, name))
	update2(fmt.Fprintf(w, "var %sContractTrieSet = ", name))
	update3(printStruct(*t, w, name))
	update2(fmt.Fprintln(w))
	return
***REMOVED***

func printArray(ct colltab.ContractTrieSet, w io.Writer, name string) (n, size int, err error) ***REMOVED***
	p := func(f string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		nn, e := fmt.Fprintf(w, f, a...)
		n += nn
		if err == nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***
	size = len(ct) * 4
	p("// %sCTEntries: %d entries, %d bytes\n", name, len(ct), size)
	p("var %sCTEntries = [%d]struct***REMOVED***L,H,N,I uint8***REMOVED******REMOVED***\n", name, len(ct))
	for _, fe := range ct ***REMOVED***
		p("\t***REMOVED***0x%X, 0x%X, %d, %d***REMOVED***,\n", fe.L, fe.H, fe.N, fe.I)
	***REMOVED***
	p("***REMOVED***\n")
	return
***REMOVED***

func printStruct(ct colltab.ContractTrieSet, w io.Writer, name string) (n, size int, err error) ***REMOVED***
	n, err = fmt.Fprintf(w, "colltab.ContractTrieSet( %sCTEntries[:] )", name)
	size = int(reflect.TypeOf(ct).Size())
	return
***REMOVED***
