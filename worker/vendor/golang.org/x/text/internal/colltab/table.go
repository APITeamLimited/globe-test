// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

import (
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// Table holds all collation data for a given collation ordering.
type Table struct ***REMOVED***
	Index Trie // main trie

	// expansion info
	ExpandElem []uint32

	// contraction info
	ContractTries  ContractTrieSet
	ContractElem   []uint32
	MaxContractLen int
	VariableTop    uint32
***REMOVED***

func (t *Table) AppendNext(w []Elem, b []byte) (res []Elem, n int) ***REMOVED***
	return t.appendNext(w, source***REMOVED***bytes: b***REMOVED***)
***REMOVED***

func (t *Table) AppendNextString(w []Elem, s string) (res []Elem, n int) ***REMOVED***
	return t.appendNext(w, source***REMOVED***str: s***REMOVED***)
***REMOVED***

func (t *Table) Start(p int, b []byte) int ***REMOVED***
	// TODO: implement
	panic("not implemented")
***REMOVED***

func (t *Table) StartString(p int, s string) int ***REMOVED***
	// TODO: implement
	panic("not implemented")
***REMOVED***

func (t *Table) Domain() []string ***REMOVED***
	// TODO: implement
	panic("not implemented")
***REMOVED***

func (t *Table) Top() uint32 ***REMOVED***
	return t.VariableTop
***REMOVED***

type source struct ***REMOVED***
	str   string
	bytes []byte
***REMOVED***

func (src *source) lookup(t *Table) (ce Elem, sz int) ***REMOVED***
	if src.bytes == nil ***REMOVED***
		return t.Index.lookupString(src.str)
	***REMOVED***
	return t.Index.lookup(src.bytes)
***REMOVED***

func (src *source) tail(sz int) ***REMOVED***
	if src.bytes == nil ***REMOVED***
		src.str = src.str[sz:]
	***REMOVED*** else ***REMOVED***
		src.bytes = src.bytes[sz:]
	***REMOVED***
***REMOVED***

func (src *source) nfd(buf []byte, end int) []byte ***REMOVED***
	if src.bytes == nil ***REMOVED***
		return norm.NFD.AppendString(buf[:0], src.str[:end])
	***REMOVED***
	return norm.NFD.Append(buf[:0], src.bytes[:end]...)
***REMOVED***

func (src *source) rune() (r rune, sz int) ***REMOVED***
	if src.bytes == nil ***REMOVED***
		return utf8.DecodeRuneInString(src.str)
	***REMOVED***
	return utf8.DecodeRune(src.bytes)
***REMOVED***

func (src *source) properties(f norm.Form) norm.Properties ***REMOVED***
	if src.bytes == nil ***REMOVED***
		return f.PropertiesString(src.str)
	***REMOVED***
	return f.Properties(src.bytes)
***REMOVED***

// appendNext appends the weights corresponding to the next rune or
// contraction in s.  If a contraction is matched to a discontinuous
// sequence of runes, the weights for the interstitial runes are
// appended as well.  It returns a new slice that includes the appended
// weights and the number of bytes consumed from s.
func (t *Table) appendNext(w []Elem, src source) (res []Elem, n int) ***REMOVED***
	ce, sz := src.lookup(t)
	tp := ce.ctype()
	if tp == ceNormal ***REMOVED***
		if ce == 0 ***REMOVED***
			r, _ := src.rune()
			const (
				hangulSize  = 3
				firstHangul = 0xAC00
				lastHangul  = 0xD7A3
			)
			if r >= firstHangul && r <= lastHangul ***REMOVED***
				// TODO: performance can be considerably improved here.
				n = sz
				var buf [16]byte // Used for decomposing Hangul.
				for b := src.nfd(buf[:0], hangulSize); len(b) > 0; b = b[sz:] ***REMOVED***
					ce, sz = t.Index.lookup(b)
					w = append(w, ce)
				***REMOVED***
				return w, n
			***REMOVED***
			ce = makeImplicitCE(implicitPrimary(r))
		***REMOVED***
		w = append(w, ce)
	***REMOVED*** else if tp == ceExpansionIndex ***REMOVED***
		w = t.appendExpansion(w, ce)
	***REMOVED*** else if tp == ceContractionIndex ***REMOVED***
		n := 0
		src.tail(sz)
		if src.bytes == nil ***REMOVED***
			w, n = t.matchContractionString(w, ce, src.str)
		***REMOVED*** else ***REMOVED***
			w, n = t.matchContraction(w, ce, src.bytes)
		***REMOVED***
		sz += n
	***REMOVED*** else if tp == ceDecompose ***REMOVED***
		// Decompose using NFKD and replace tertiary weights.
		t1, t2 := splitDecompose(ce)
		i := len(w)
		nfkd := src.properties(norm.NFKD).Decomposition()
		for p := 0; len(nfkd) > 0; nfkd = nfkd[p:] ***REMOVED***
			w, p = t.appendNext(w, source***REMOVED***bytes: nfkd***REMOVED***)
		***REMOVED***
		w[i] = w[i].updateTertiary(t1)
		if i++; i < len(w) ***REMOVED***
			w[i] = w[i].updateTertiary(t2)
			for i++; i < len(w); i++ ***REMOVED***
				w[i] = w[i].updateTertiary(maxTertiary)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return w, sz
***REMOVED***

func (t *Table) appendExpansion(w []Elem, ce Elem) []Elem ***REMOVED***
	i := splitExpandIndex(ce)
	n := int(t.ExpandElem[i])
	i++
	for _, ce := range t.ExpandElem[i : i+n] ***REMOVED***
		w = append(w, Elem(ce))
	***REMOVED***
	return w
***REMOVED***

func (t *Table) matchContraction(w []Elem, ce Elem, suffix []byte) ([]Elem, int) ***REMOVED***
	index, n, offset := splitContractIndex(ce)

	scan := t.ContractTries.scanner(index, n, suffix)
	buf := [norm.MaxSegmentSize]byte***REMOVED******REMOVED***
	bufp := 0
	p := scan.scan(0)

	if !scan.done && p < len(suffix) && suffix[p] >= utf8.RuneSelf ***REMOVED***
		// By now we should have filtered most cases.
		p0 := p
		bufn := 0
		rune := norm.NFD.Properties(suffix[p:])
		p += rune.Size()
		if rune.LeadCCC() != 0 ***REMOVED***
			prevCC := rune.TrailCCC()
			// A gap may only occur in the last normalization segment.
			// This also ensures that len(scan.s) < norm.MaxSegmentSize.
			if end := norm.NFD.FirstBoundary(suffix[p:]); end != -1 ***REMOVED***
				scan.s = suffix[:p+end]
			***REMOVED***
			for p < len(suffix) && !scan.done && suffix[p] >= utf8.RuneSelf ***REMOVED***
				rune = norm.NFD.Properties(suffix[p:])
				if ccc := rune.LeadCCC(); ccc == 0 || prevCC >= ccc ***REMOVED***
					break
				***REMOVED***
				prevCC = rune.TrailCCC()
				if pp := scan.scan(p); pp != p ***REMOVED***
					// Copy the interstitial runes for later processing.
					bufn += copy(buf[bufn:], suffix[p0:p])
					if scan.pindex == pp ***REMOVED***
						bufp = bufn
					***REMOVED***
					p, p0 = pp, pp
				***REMOVED*** else ***REMOVED***
					p += rune.Size()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Append weights for the matched contraction, which may be an expansion.
	i, n := scan.result()
	ce = Elem(t.ContractElem[i+offset])
	if ce.ctype() == ceNormal ***REMOVED***
		w = append(w, ce)
	***REMOVED*** else ***REMOVED***
		w = t.appendExpansion(w, ce)
	***REMOVED***
	// Append weights for the runes in the segment not part of the contraction.
	for b, p := buf[:bufp], 0; len(b) > 0; b = b[p:] ***REMOVED***
		w, p = t.appendNext(w, source***REMOVED***bytes: b***REMOVED***)
	***REMOVED***
	return w, n
***REMOVED***

// TODO: unify the two implementations. This is best done after first simplifying
// the algorithm taking into account the inclusion of both NFC and NFD forms
// in the table.
func (t *Table) matchContractionString(w []Elem, ce Elem, suffix string) ([]Elem, int) ***REMOVED***
	index, n, offset := splitContractIndex(ce)

	scan := t.ContractTries.scannerString(index, n, suffix)
	buf := [norm.MaxSegmentSize]byte***REMOVED******REMOVED***
	bufp := 0
	p := scan.scan(0)

	if !scan.done && p < len(suffix) && suffix[p] >= utf8.RuneSelf ***REMOVED***
		// By now we should have filtered most cases.
		p0 := p
		bufn := 0
		rune := norm.NFD.PropertiesString(suffix[p:])
		p += rune.Size()
		if rune.LeadCCC() != 0 ***REMOVED***
			prevCC := rune.TrailCCC()
			// A gap may only occur in the last normalization segment.
			// This also ensures that len(scan.s) < norm.MaxSegmentSize.
			if end := norm.NFD.FirstBoundaryInString(suffix[p:]); end != -1 ***REMOVED***
				scan.s = suffix[:p+end]
			***REMOVED***
			for p < len(suffix) && !scan.done && suffix[p] >= utf8.RuneSelf ***REMOVED***
				rune = norm.NFD.PropertiesString(suffix[p:])
				if ccc := rune.LeadCCC(); ccc == 0 || prevCC >= ccc ***REMOVED***
					break
				***REMOVED***
				prevCC = rune.TrailCCC()
				if pp := scan.scan(p); pp != p ***REMOVED***
					// Copy the interstitial runes for later processing.
					bufn += copy(buf[bufn:], suffix[p0:p])
					if scan.pindex == pp ***REMOVED***
						bufp = bufn
					***REMOVED***
					p, p0 = pp, pp
				***REMOVED*** else ***REMOVED***
					p += rune.Size()
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Append weights for the matched contraction, which may be an expansion.
	i, n := scan.result()
	ce = Elem(t.ContractElem[i+offset])
	if ce.ctype() == ceNormal ***REMOVED***
		w = append(w, ce)
	***REMOVED*** else ***REMOVED***
		w = t.appendExpansion(w, ce)
	***REMOVED***
	// Append weights for the runes in the segment not part of the contraction.
	for b, p := buf[:bufp], 0; len(b) > 0; b = b[p:] ***REMOVED***
		w, p = t.appendNext(w, source***REMOVED***bytes: b***REMOVED***)
	***REMOVED***
	return w, n
***REMOVED***
