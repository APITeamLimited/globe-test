// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import "encoding/binary"

// This file contains Form-specific logic and wrappers for data in tables.go.

// Rune info is stored in a separate trie per composing form. A composing form
// and its corresponding decomposing form share the same trie.  Each trie maps
// a rune to a uint16. The values take two forms.  For v >= 0x8000:
//   bits
//   15:    1 (inverse of NFD_QC bit of qcInfo)
//   13..7: qcInfo (see below). isYesD is always true (no decompostion).
//    6..0: ccc (compressed CCC value).
// For v < 0x8000, the respective rune has a decomposition and v is an index
// into a byte array of UTF-8 decomposition sequences and additional info and
// has the form:
//    <header> <decomp_byte>* [<tccc> [<lccc>]]
// The header contains the number of bytes in the decomposition (excluding this
// length byte). The two most significant bits of this length byte correspond
// to bit 5 and 4 of qcInfo (see below).  The byte sequence itself starts at v+1.
// The byte sequence is followed by a trailing and leading CCC if the values
// for these are not zero.  The value of v determines which ccc are appended
// to the sequences.  For v < firstCCC, there are none, for v >= firstCCC,
// the sequence is followed by a trailing ccc, and for v >= firstLeadingCC
// there is an additional leading ccc. The value of tccc itself is the
// trailing CCC shifted left 2 bits. The two least-significant bits of tccc
// are the number of trailing non-starters.

const (
	qcInfoMask      = 0x3F // to clear all but the relevant bits in a qcInfo
	headerLenMask   = 0x3F // extract the length value from the header byte
	headerFlagsMask = 0xC0 // extract the qcInfo bits from the header byte
)

// Properties provides access to normalization properties of a rune.
type Properties struct ***REMOVED***
	pos   uint8  // start position in reorderBuffer; used in composition.go
	size  uint8  // length of UTF-8 encoding of this rune
	ccc   uint8  // leading canonical combining class (ccc if not decomposition)
	tccc  uint8  // trailing canonical combining class (ccc if not decomposition)
	nLead uint8  // number of leading non-starters.
	flags qcInfo // quick check flags
	index uint16
***REMOVED***

// functions dispatchable per form
type lookupFunc func(b input, i int) Properties

// formInfo holds Form-specific functions and tables.
type formInfo struct ***REMOVED***
	form                     Form
	composing, compatibility bool // form type
	info                     lookupFunc
	nextMain                 iterFunc
***REMOVED***

var formTable = []*formInfo***REMOVED******REMOVED***
	form:          NFC,
	composing:     true,
	compatibility: false,
	info:          lookupInfoNFC,
	nextMain:      nextComposed,
***REMOVED***, ***REMOVED***
	form:          NFD,
	composing:     false,
	compatibility: false,
	info:          lookupInfoNFC,
	nextMain:      nextDecomposed,
***REMOVED***, ***REMOVED***
	form:          NFKC,
	composing:     true,
	compatibility: true,
	info:          lookupInfoNFKC,
	nextMain:      nextComposed,
***REMOVED***, ***REMOVED***
	form:          NFKD,
	composing:     false,
	compatibility: true,
	info:          lookupInfoNFKC,
	nextMain:      nextDecomposed,
***REMOVED******REMOVED***

// We do not distinguish between boundaries for NFC, NFD, etc. to avoid
// unexpected behavior for the user.  For example, in NFD, there is a boundary
// after 'a'.  However, 'a' might combine with modifiers, so from the application's
// perspective it is not a good boundary. We will therefore always use the
// boundaries for the combining variants.

// BoundaryBefore returns true if this rune starts a new segment and
// cannot combine with any rune on the left.
func (p Properties) BoundaryBefore() bool ***REMOVED***
	if p.ccc == 0 && !p.combinesBackward() ***REMOVED***
		return true
	***REMOVED***
	// We assume that the CCC of the first character in a decomposition
	// is always non-zero if different from info.ccc and that we can return
	// false at this point. This is verified by maketables.
	return false
***REMOVED***

// BoundaryAfter returns true if runes cannot combine with or otherwise
// interact with this or previous runes.
func (p Properties) BoundaryAfter() bool ***REMOVED***
	// TODO: loosen these conditions.
	return p.isInert()
***REMOVED***

// We pack quick check data in 4 bits:
//   5:    Combines forward  (0 == false, 1 == true)
//   4..3: NFC_QC Yes(00), No (10), or Maybe (11)
//   2:    NFD_QC Yes (0) or No (1). No also means there is a decomposition.
//   1..0: Number of trailing non-starters.
//
// When all 4 bits are zero, the character is inert, meaning it is never
// influenced by normalization.
type qcInfo uint8

func (p Properties) isYesC() bool ***REMOVED*** return p.flags&0x10 == 0 ***REMOVED***
func (p Properties) isYesD() bool ***REMOVED*** return p.flags&0x4 == 0 ***REMOVED***

func (p Properties) combinesForward() bool  ***REMOVED*** return p.flags&0x20 != 0 ***REMOVED***
func (p Properties) combinesBackward() bool ***REMOVED*** return p.flags&0x8 != 0 ***REMOVED*** // == isMaybe
func (p Properties) hasDecomposition() bool ***REMOVED*** return p.flags&0x4 != 0 ***REMOVED*** // == isNoD

func (p Properties) isInert() bool ***REMOVED***
	return p.flags&qcInfoMask == 0 && p.ccc == 0
***REMOVED***

func (p Properties) multiSegment() bool ***REMOVED***
	return p.index >= firstMulti && p.index < endMulti
***REMOVED***

func (p Properties) nLeadingNonStarters() uint8 ***REMOVED***
	return p.nLead
***REMOVED***

func (p Properties) nTrailingNonStarters() uint8 ***REMOVED***
	return uint8(p.flags & 0x03)
***REMOVED***

// Decomposition returns the decomposition for the underlying rune
// or nil if there is none.
func (p Properties) Decomposition() []byte ***REMOVED***
	// TODO: create the decomposition for Hangul?
	if p.index == 0 ***REMOVED***
		return nil
	***REMOVED***
	i := p.index
	n := decomps[i] & headerLenMask
	i++
	return decomps[i : i+uint16(n)]
***REMOVED***

// Size returns the length of UTF-8 encoding of the rune.
func (p Properties) Size() int ***REMOVED***
	return int(p.size)
***REMOVED***

// CCC returns the canonical combining class of the underlying rune.
func (p Properties) CCC() uint8 ***REMOVED***
	if p.index >= firstCCCZeroExcept ***REMOVED***
		return 0
	***REMOVED***
	return ccc[p.ccc]
***REMOVED***

// LeadCCC returns the CCC of the first rune in the decomposition.
// If there is no decomposition, LeadCCC equals CCC.
func (p Properties) LeadCCC() uint8 ***REMOVED***
	return ccc[p.ccc]
***REMOVED***

// TrailCCC returns the CCC of the last rune in the decomposition.
// If there is no decomposition, TrailCCC equals CCC.
func (p Properties) TrailCCC() uint8 ***REMOVED***
	return ccc[p.tccc]
***REMOVED***

func buildRecompMap() ***REMOVED***
	recompMap = make(map[uint32]rune, len(recompMapPacked)/8)
	var buf [8]byte
	for i := 0; i < len(recompMapPacked); i += 8 ***REMOVED***
		copy(buf[:], recompMapPacked[i:i+8])
		key := binary.BigEndian.Uint32(buf[:4])
		val := binary.BigEndian.Uint32(buf[4:])
		recompMap[key] = rune(val)
	***REMOVED***
***REMOVED***

// Recomposition
// We use 32-bit keys instead of 64-bit for the two codepoint keys.
// This clips off the bits of three entries, but we know this will not
// result in a collision. In the unlikely event that changes to
// UnicodeData.txt introduce collisions, the compiler will catch it.
// Note that the recomposition map for NFC and NFKC are identical.

// combine returns the combined rune or 0 if it doesn't exist.
//
// The caller is responsible for calling
// recompMapOnce.Do(buildRecompMap) sometime before this is called.
func combine(a, b rune) rune ***REMOVED***
	key := uint32(uint16(a))<<16 + uint32(uint16(b))
	if recompMap == nil ***REMOVED***
		panic("caller error") // see func comment
	***REMOVED***
	return recompMap[key]
***REMOVED***

func lookupInfoNFC(b input, i int) Properties ***REMOVED***
	v, sz := b.charinfoNFC(i)
	return compInfo(v, sz)
***REMOVED***

func lookupInfoNFKC(b input, i int) Properties ***REMOVED***
	v, sz := b.charinfoNFKC(i)
	return compInfo(v, sz)
***REMOVED***

// Properties returns properties for the first rune in s.
func (f Form) Properties(s []byte) Properties ***REMOVED***
	if f == NFC || f == NFD ***REMOVED***
		return compInfo(nfcData.lookup(s))
	***REMOVED***
	return compInfo(nfkcData.lookup(s))
***REMOVED***

// PropertiesString returns properties for the first rune in s.
func (f Form) PropertiesString(s string) Properties ***REMOVED***
	if f == NFC || f == NFD ***REMOVED***
		return compInfo(nfcData.lookupString(s))
	***REMOVED***
	return compInfo(nfkcData.lookupString(s))
***REMOVED***

// compInfo converts the information contained in v and sz
// to a Properties.  See the comment at the top of the file
// for more information on the format.
func compInfo(v uint16, sz int) Properties ***REMOVED***
	if v == 0 ***REMOVED***
		return Properties***REMOVED***size: uint8(sz)***REMOVED***
	***REMOVED*** else if v >= 0x8000 ***REMOVED***
		p := Properties***REMOVED***
			size:  uint8(sz),
			ccc:   uint8(v),
			tccc:  uint8(v),
			flags: qcInfo(v >> 8),
		***REMOVED***
		if p.ccc > 0 || p.combinesBackward() ***REMOVED***
			p.nLead = uint8(p.flags & 0x3)
		***REMOVED***
		return p
	***REMOVED***
	// has decomposition
	h := decomps[v]
	f := (qcInfo(h&headerFlagsMask) >> 2) | 0x4
	p := Properties***REMOVED***size: uint8(sz), flags: f, index: v***REMOVED***
	if v >= firstCCC ***REMOVED***
		v += uint16(h&headerLenMask) + 1
		c := decomps[v]
		p.tccc = c >> 2
		p.flags |= qcInfo(c & 0x3)
		if v >= firstLeadingCCC ***REMOVED***
			p.nLead = c & 0x3
			if v >= firstStarterWithNLead ***REMOVED***
				// We were tricked. Remove the decomposition.
				p.flags &= 0x03
				p.index = 0
				return p
			***REMOVED***
			p.ccc = decomps[v+1]
		***REMOVED***
	***REMOVED***
	return p
***REMOVED***
