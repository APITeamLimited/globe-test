// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

import (
	"fmt"
	"unicode"
)

// Level identifies the collation comparison level.
// The primary level corresponds to the basic sorting of text.
// The secondary level corresponds to accents and related linguistic elements.
// The tertiary level corresponds to casing and related concepts.
// The quaternary level is derived from the other levels by the
// various algorithms for handling variable elements.
type Level int

const (
	Primary Level = iota
	Secondary
	Tertiary
	Quaternary
	Identity

	NumLevels
)

const (
	defaultSecondary = 0x20
	defaultTertiary  = 0x2
	maxTertiary      = 0x1F
	MaxQuaternary    = 0x1FFFFF // 21 bits.
)

// Elem is a representation of a collation element. This API provides ways to encode
// and decode Elems. Implementations of collation tables may use values greater
// or equal to PrivateUse for their own purposes.  However, these should never be
// returned by AppendNext.
type Elem uint32

const (
	maxCE       Elem = 0xAFFFFFFF
	PrivateUse       = minContract
	minContract      = 0xC0000000
	maxContract      = 0xDFFFFFFF
	minExpand        = 0xE0000000
	maxExpand        = 0xEFFFFFFF
	minDecomp        = 0xF0000000
)

type ceType int

const (
	ceNormal           ceType = iota // ceNormal includes implicits (ce == 0)
	ceContractionIndex               // rune can be a start of a contraction
	ceExpansionIndex                 // rune expands into a sequence of collation elements
	ceDecompose                      // rune expands using NFKC decomposition
)

func (ce Elem) ctype() ceType ***REMOVED***
	if ce <= maxCE ***REMOVED***
		return ceNormal
	***REMOVED***
	if ce <= maxContract ***REMOVED***
		return ceContractionIndex
	***REMOVED*** else ***REMOVED***
		if ce <= maxExpand ***REMOVED***
			return ceExpansionIndex
		***REMOVED***
		return ceDecompose
	***REMOVED***
	panic("should not reach here")
	return ceType(-1)
***REMOVED***

// For normal collation elements, we assume that a collation element either has
// a primary or non-default secondary value, not both.
// Collation elements with a primary value are of the form
// 01pppppp pppppppp ppppppp0 ssssssss
//   - p* is primary collation value
//   - s* is the secondary collation value
// 00pppppp pppppppp ppppppps sssttttt, where
//   - p* is primary collation value
//   - s* offset of secondary from default value.
//   - t* is the tertiary collation value
// 100ttttt cccccccc pppppppp pppppppp
//   - t* is the tertiar collation value
//   - c* is the canonical combining class
//   - p* is the primary collation value
// Collation elements with a secondary value are of the form
// 1010cccc ccccssss ssssssss tttttttt, where
//   - c* is the canonical combining class
//   - s* is the secondary collation value
//   - t* is the tertiary collation value
// 11qqqqqq qqqqqqqq qqqqqqq0 00000000
//   - q* quaternary value
const (
	ceTypeMask              = 0xC0000000
	ceTypeMaskExt           = 0xE0000000
	ceIgnoreMask            = 0xF00FFFFF
	ceType1                 = 0x40000000
	ceType2                 = 0x00000000
	ceType3or4              = 0x80000000
	ceType4                 = 0xA0000000
	ceTypeQ                 = 0xC0000000
	Ignore                  = ceType4
	firstNonPrimary         = 0x80000000
	lastSpecialPrimary      = 0xA0000000
	secondaryMask           = 0x80000000
	hasTertiaryMask         = 0x40000000
	primaryValueMask        = 0x3FFFFE00
	maxPrimaryBits          = 21
	compactPrimaryBits      = 16
	maxSecondaryBits        = 12
	maxTertiaryBits         = 8
	maxCCCBits              = 8
	maxSecondaryCompactBits = 8
	maxSecondaryDiffBits    = 4
	maxTertiaryCompactBits  = 5
	primaryShift            = 9
	compactSecondaryShift   = 5
	minCompactSecondary     = defaultSecondary - 4
)

func makeImplicitCE(primary int) Elem ***REMOVED***
	return ceType1 | Elem(primary<<primaryShift) | defaultSecondary
***REMOVED***

// MakeElem returns an Elem for the given values.  It will return an error
// if the given combination of values is invalid.
func MakeElem(primary, secondary, tertiary int, ccc uint8) (Elem, error) ***REMOVED***
	if w := primary; w >= 1<<maxPrimaryBits || w < 0 ***REMOVED***
		return 0, fmt.Errorf("makeCE: primary weight out of bounds: %x >= %x", w, 1<<maxPrimaryBits)
	***REMOVED***
	if w := secondary; w >= 1<<maxSecondaryBits || w < 0 ***REMOVED***
		return 0, fmt.Errorf("makeCE: secondary weight out of bounds: %x >= %x", w, 1<<maxSecondaryBits)
	***REMOVED***
	if w := tertiary; w >= 1<<maxTertiaryBits || w < 0 ***REMOVED***
		return 0, fmt.Errorf("makeCE: tertiary weight out of bounds: %x >= %x", w, 1<<maxTertiaryBits)
	***REMOVED***
	ce := Elem(0)
	if primary != 0 ***REMOVED***
		if ccc != 0 ***REMOVED***
			if primary >= 1<<compactPrimaryBits ***REMOVED***
				return 0, fmt.Errorf("makeCE: primary weight with non-zero CCC out of bounds: %x >= %x", primary, 1<<compactPrimaryBits)
			***REMOVED***
			if secondary != defaultSecondary ***REMOVED***
				return 0, fmt.Errorf("makeCE: cannot combine non-default secondary value (%x) with non-zero CCC (%x)", secondary, ccc)
			***REMOVED***
			ce = Elem(tertiary << (compactPrimaryBits + maxCCCBits))
			ce |= Elem(ccc) << compactPrimaryBits
			ce |= Elem(primary)
			ce |= ceType3or4
		***REMOVED*** else if tertiary == defaultTertiary ***REMOVED***
			if secondary >= 1<<maxSecondaryCompactBits ***REMOVED***
				return 0, fmt.Errorf("makeCE: secondary weight with non-zero primary out of bounds: %x >= %x", secondary, 1<<maxSecondaryCompactBits)
			***REMOVED***
			ce = Elem(primary<<(maxSecondaryCompactBits+1) + secondary)
			ce |= ceType1
		***REMOVED*** else ***REMOVED***
			d := secondary - defaultSecondary + maxSecondaryDiffBits
			if d >= 1<<maxSecondaryDiffBits || d < 0 ***REMOVED***
				return 0, fmt.Errorf("makeCE: secondary weight diff out of bounds: %x < 0 || %x > %x", d, d, 1<<maxSecondaryDiffBits)
			***REMOVED***
			if tertiary >= 1<<maxTertiaryCompactBits ***REMOVED***
				return 0, fmt.Errorf("makeCE: tertiary weight with non-zero primary out of bounds: %x > %x", tertiary, 1<<maxTertiaryCompactBits)
			***REMOVED***
			ce = Elem(primary<<maxSecondaryDiffBits + d)
			ce = ce<<maxTertiaryCompactBits + Elem(tertiary)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ce = Elem(secondary<<maxTertiaryBits + tertiary)
		ce += Elem(ccc) << (maxSecondaryBits + maxTertiaryBits)
		ce |= ceType4
	***REMOVED***
	return ce, nil
***REMOVED***

// MakeQuaternary returns an Elem with the given quaternary value.
func MakeQuaternary(v int) Elem ***REMOVED***
	return ceTypeQ | Elem(v<<primaryShift)
***REMOVED***

// Mask sets weights for any level smaller than l to 0.
// The resulting Elem can be used to test for equality with
// other Elems to which the same mask has been applied.
func (ce Elem) Mask(l Level) uint32 ***REMOVED***
	return 0
***REMOVED***

// CCC returns the canonical combining class associated with the underlying character,
// if applicable, or 0 otherwise.
func (ce Elem) CCC() uint8 ***REMOVED***
	if ce&ceType3or4 != 0 ***REMOVED***
		if ce&ceType4 == ceType3or4 ***REMOVED***
			return uint8(ce >> 16)
		***REMOVED***
		return uint8(ce >> 20)
	***REMOVED***
	return 0
***REMOVED***

// Primary returns the primary collation weight for ce.
func (ce Elem) Primary() int ***REMOVED***
	if ce >= firstNonPrimary ***REMOVED***
		if ce > lastSpecialPrimary ***REMOVED***
			return 0
		***REMOVED***
		return int(uint16(ce))
	***REMOVED***
	return int(ce&primaryValueMask) >> primaryShift
***REMOVED***

// Secondary returns the secondary collation weight for ce.
func (ce Elem) Secondary() int ***REMOVED***
	switch ce & ceTypeMask ***REMOVED***
	case ceType1:
		return int(uint8(ce))
	case ceType2:
		return minCompactSecondary + int((ce>>compactSecondaryShift)&0xF)
	case ceType3or4:
		if ce < ceType4 ***REMOVED***
			return defaultSecondary
		***REMOVED***
		return int(ce>>8) & 0xFFF
	case ceTypeQ:
		return 0
	***REMOVED***
	panic("should not reach here")
***REMOVED***

// Tertiary returns the tertiary collation weight for ce.
func (ce Elem) Tertiary() uint8 ***REMOVED***
	if ce&hasTertiaryMask == 0 ***REMOVED***
		if ce&ceType3or4 == 0 ***REMOVED***
			return uint8(ce & 0x1F)
		***REMOVED***
		if ce&ceType4 == ceType4 ***REMOVED***
			return uint8(ce)
		***REMOVED***
		return uint8(ce>>24) & 0x1F // type 2
	***REMOVED*** else if ce&ceTypeMask == ceType1 ***REMOVED***
		return defaultTertiary
	***REMOVED***
	// ce is a quaternary value.
	return 0
***REMOVED***

func (ce Elem) updateTertiary(t uint8) Elem ***REMOVED***
	if ce&ceTypeMask == ceType1 ***REMOVED***
		// convert to type 4
		nce := ce & primaryValueMask
		nce |= Elem(uint8(ce)-minCompactSecondary) << compactSecondaryShift
		ce = nce
	***REMOVED*** else if ce&ceTypeMaskExt == ceType3or4 ***REMOVED***
		ce &= ^Elem(maxTertiary << 24)
		return ce | (Elem(t) << 24)
	***REMOVED*** else ***REMOVED***
		// type 2 or 4
		ce &= ^Elem(maxTertiary)
	***REMOVED***
	return ce | Elem(t)
***REMOVED***

// Quaternary returns the quaternary value if explicitly specified,
// 0 if ce == Ignore, or MaxQuaternary otherwise.
// Quaternary values are used only for shifted variants.
func (ce Elem) Quaternary() int ***REMOVED***
	if ce&ceTypeMask == ceTypeQ ***REMOVED***
		return int(ce&primaryValueMask) >> primaryShift
	***REMOVED*** else if ce&ceIgnoreMask == Ignore ***REMOVED***
		return 0
	***REMOVED***
	return MaxQuaternary
***REMOVED***

// Weight returns the collation weight for the given level.
func (ce Elem) Weight(l Level) int ***REMOVED***
	switch l ***REMOVED***
	case Primary:
		return ce.Primary()
	case Secondary:
		return ce.Secondary()
	case Tertiary:
		return int(ce.Tertiary())
	case Quaternary:
		return ce.Quaternary()
	***REMOVED***
	return 0 // return 0 (ignore) for undefined levels.
***REMOVED***

// For contractions, collation elements are of the form
// 110bbbbb bbbbbbbb iiiiiiii iiiinnnn, where
//   - n* is the size of the first node in the contraction trie.
//   - i* is the index of the first node in the contraction trie.
//   - b* is the offset into the contraction collation element table.
// See contract.go for details on the contraction trie.
const (
	maxNBits              = 4
	maxTrieIndexBits      = 12
	maxContractOffsetBits = 13
)

func splitContractIndex(ce Elem) (index, n, offset int) ***REMOVED***
	n = int(ce & (1<<maxNBits - 1))
	ce >>= maxNBits
	index = int(ce & (1<<maxTrieIndexBits - 1))
	ce >>= maxTrieIndexBits
	offset = int(ce & (1<<maxContractOffsetBits - 1))
	return
***REMOVED***

// For expansions, Elems are of the form 11100000 00000000 bbbbbbbb bbbbbbbb,
// where b* is the index into the expansion sequence table.
const maxExpandIndexBits = 16

func splitExpandIndex(ce Elem) (index int) ***REMOVED***
	return int(uint16(ce))
***REMOVED***

// Some runes can be expanded using NFKD decomposition. Instead of storing the full
// sequence of collation elements, we decompose the rune and lookup the collation
// elements for each rune in the decomposition and modify the tertiary weights.
// The Elem, in this case, is of the form 11110000 00000000 wwwwwwww vvvvvvvv, where
//   - v* is the replacement tertiary weight for the first rune,
//   - w* is the replacement tertiary weight for the second rune,
// Tertiary weights of subsequent runes should be replaced with maxTertiary.
// See https://www.unicode.org/reports/tr10/#Compatibility_Decompositions for more details.
func splitDecompose(ce Elem) (t1, t2 uint8) ***REMOVED***
	return uint8(ce), uint8(ce >> 8)
***REMOVED***

const (
	// These constants were taken from https://www.unicode.org/versions/Unicode6.0.0/ch12.pdf.
	minUnified       rune = 0x4E00
	maxUnified            = 0x9FFF
	minCompatibility      = 0xF900
	maxCompatibility      = 0xFAFF
	minRare               = 0x3400
	maxRare               = 0x4DBF
)
const (
	commonUnifiedOffset = 0x10000
	rareUnifiedOffset   = 0x20000 // largest rune in common is U+FAFF
	otherOffset         = 0x50000 // largest rune in rare is U+2FA1D
	illegalOffset       = otherOffset + int(unicode.MaxRune)
	maxPrimary          = illegalOffset + 1
)

// implicitPrimary returns the primary weight for the a rune
// for which there is no entry for the rune in the collation table.
// We take a different approach from the one specified in
// https://unicode.org/reports/tr10/#Implicit_Weights,
// but preserve the resulting relative ordering of the runes.
func implicitPrimary(r rune) int ***REMOVED***
	if unicode.Is(unicode.Ideographic, r) ***REMOVED***
		if r >= minUnified && r <= maxUnified ***REMOVED***
			// The most common case for CJK.
			return int(r) + commonUnifiedOffset
		***REMOVED***
		if r >= minCompatibility && r <= maxCompatibility ***REMOVED***
			// This will typically not hit. The DUCET explicitly specifies mappings
			// for all characters that do not decompose.
			return int(r) + commonUnifiedOffset
		***REMOVED***
		return int(r) + rareUnifiedOffset
	***REMOVED***
	return int(r) + otherOffset
***REMOVED***
