// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

import "unicode/utf8"

// For a description of ContractTrieSet, see text/collate/build/contract.go.

type ContractTrieSet []struct***REMOVED*** L, H, N, I uint8 ***REMOVED***

// ctScanner is used to match a trie to an input sequence.
// A contraction may match a non-contiguous sequence of bytes in an input string.
// For example, if there is a contraction for <a, combining_ring>, it should match
// the sequence <a, combining_cedilla, combining_ring>, as combining_cedilla does
// not block combining_ring.
// ctScanner does not automatically skip over non-blocking non-starters, but rather
// retains the state of the last match and leaves it up to the user to continue
// the match at the appropriate points.
type ctScanner struct ***REMOVED***
	states ContractTrieSet
	s      []byte
	n      int
	index  int
	pindex int
	done   bool
***REMOVED***

type ctScannerString struct ***REMOVED***
	states ContractTrieSet
	s      string
	n      int
	index  int
	pindex int
	done   bool
***REMOVED***

func (t ContractTrieSet) scanner(index, n int, b []byte) ctScanner ***REMOVED***
	return ctScanner***REMOVED***s: b, states: t[index:], n: n***REMOVED***
***REMOVED***

func (t ContractTrieSet) scannerString(index, n int, str string) ctScannerString ***REMOVED***
	return ctScannerString***REMOVED***s: str, states: t[index:], n: n***REMOVED***
***REMOVED***

// result returns the offset i and bytes consumed p so far.  If no suffix
// matched, i and p will be 0.
func (s *ctScanner) result() (i, p int) ***REMOVED***
	return s.index, s.pindex
***REMOVED***

func (s *ctScannerString) result() (i, p int) ***REMOVED***
	return s.index, s.pindex
***REMOVED***

const (
	final   = 0
	noIndex = 0xFF
)

// scan matches the longest suffix at the current location in the input
// and returns the number of bytes consumed.
func (s *ctScanner) scan(p int) int ***REMOVED***
	pr := p // the p at the rune start
	str := s.s
	states, n := s.states, s.n
	for i := 0; i < n && p < len(str); ***REMOVED***
		e := states[i]
		c := str[p]
		// TODO: a significant number of contractions are of a form that
		// cannot match discontiguous UTF-8 in a normalized string. We could let
		// a negative value of e.n mean that we can set s.done = true and avoid
		// the need for additional matches.
		if c >= e.L ***REMOVED***
			if e.L == c ***REMOVED***
				p++
				if e.I != noIndex ***REMOVED***
					s.index = int(e.I)
					s.pindex = p
				***REMOVED***
				if e.N != final ***REMOVED***
					i, states, n = 0, states[int(e.H)+n:], int(e.N)
					if p >= len(str) || utf8.RuneStart(str[p]) ***REMOVED***
						s.states, s.n, pr = states, n, p
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					s.done = true
					return p
				***REMOVED***
				continue
			***REMOVED*** else if e.N == final && c <= e.H ***REMOVED***
				p++
				s.done = true
				s.index = int(c-e.L) + int(e.I)
				s.pindex = p
				return p
			***REMOVED***
		***REMOVED***
		i++
	***REMOVED***
	return pr
***REMOVED***

// scan is a verbatim copy of ctScanner.scan.
func (s *ctScannerString) scan(p int) int ***REMOVED***
	pr := p // the p at the rune start
	str := s.s
	states, n := s.states, s.n
	for i := 0; i < n && p < len(str); ***REMOVED***
		e := states[i]
		c := str[p]
		// TODO: a significant number of contractions are of a form that
		// cannot match discontiguous UTF-8 in a normalized string. We could let
		// a negative value of e.n mean that we can set s.done = true and avoid
		// the need for additional matches.
		if c >= e.L ***REMOVED***
			if e.L == c ***REMOVED***
				p++
				if e.I != noIndex ***REMOVED***
					s.index = int(e.I)
					s.pindex = p
				***REMOVED***
				if e.N != final ***REMOVED***
					i, states, n = 0, states[int(e.H)+n:], int(e.N)
					if p >= len(str) || utf8.RuneStart(str[p]) ***REMOVED***
						s.states, s.n, pr = states, n, p
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					s.done = true
					return p
				***REMOVED***
				continue
			***REMOVED*** else if e.N == final && c <= e.H ***REMOVED***
				p++
				s.done = true
				s.index = int(c-e.L) + int(e.I)
				s.pindex = p
				return p
			***REMOVED***
		***REMOVED***
		i++
	***REMOVED***
	return pr
***REMOVED***
