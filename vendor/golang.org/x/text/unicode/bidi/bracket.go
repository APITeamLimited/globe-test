// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bidi

import (
	"container/list"
	"fmt"
	"sort"
)

// This file contains a port of the reference implementation of the
// Bidi Parentheses Algorithm:
// https://www.unicode.org/Public/PROGRAMS/BidiReferenceJava/BidiPBAReference.java
//
// The implementation in this file covers definitions BD14-BD16 and rule N0
// of UAX#9.
//
// Some preprocessing is done for each rune before data is passed to this
// algorithm:
//  - opening and closing brackets are identified
//  - a bracket pair type, like '(' and ')' is assigned a unique identifier that
//    is identical for the opening and closing bracket. It is left to do these
//    mappings.
//  - The BPA algorithm requires that bracket characters that are canonical
//    equivalents of each other be able to be substituted for each other.
//    It is the responsibility of the caller to do this canonicalization.
//
// In implementing BD16, this implementation departs slightly from the "logical"
// algorithm defined in UAX#9. In particular, the stack referenced there
// supports operations that go beyond a "basic" stack. An equivalent
// implementation based on a linked list is used here.

// Bidi_Paired_Bracket_Type
// BD14. An opening paired bracket is a character whose
// Bidi_Paired_Bracket_Type property value is Open.
//
// BD15. A closing paired bracket is a character whose
// Bidi_Paired_Bracket_Type property value is Close.
type bracketType byte

const (
	bpNone bracketType = iota
	bpOpen
	bpClose
)

// bracketPair holds a pair of index values for opening and closing bracket
// location of a bracket pair.
type bracketPair struct ***REMOVED***
	opener int
	closer int
***REMOVED***

func (b *bracketPair) String() string ***REMOVED***
	return fmt.Sprintf("(%v, %v)", b.opener, b.closer)
***REMOVED***

// bracketPairs is a slice of bracketPairs with a sort.Interface implementation.
type bracketPairs []bracketPair

func (b bracketPairs) Len() int           ***REMOVED*** return len(b) ***REMOVED***
func (b bracketPairs) Swap(i, j int)      ***REMOVED*** b[i], b[j] = b[j], b[i] ***REMOVED***
func (b bracketPairs) Less(i, j int) bool ***REMOVED*** return b[i].opener < b[j].opener ***REMOVED***

// resolvePairedBrackets runs the paired bracket part of the UBA algorithm.
//
// For each rune, it takes the indexes into the original string, the class the
// bracket type (in pairTypes) and the bracket identifier (pairValues). It also
// takes the direction type for the start-of-sentence and the embedding level.
//
// The identifiers for bracket types are the rune of the canonicalized opening
// bracket for brackets (open or close) or 0 for runes that are not brackets.
func resolvePairedBrackets(s *isolatingRunSequence) ***REMOVED***
	p := bracketPairer***REMOVED***
		sos:              s.sos,
		openers:          list.New(),
		codesIsolatedRun: s.types,
		indexes:          s.indexes,
	***REMOVED***
	dirEmbed := L
	if s.level&1 != 0 ***REMOVED***
		dirEmbed = R
	***REMOVED***
	p.locateBrackets(s.p.pairTypes, s.p.pairValues)
	p.resolveBrackets(dirEmbed, s.p.initialTypes)
***REMOVED***

type bracketPairer struct ***REMOVED***
	sos Class // direction corresponding to start of sequence

	// The following is a restatement of BD 16 using non-algorithmic language.
	//
	// A bracket pair is a pair of characters consisting of an opening
	// paired bracket and a closing paired bracket such that the
	// Bidi_Paired_Bracket property value of the former equals the latter,
	// subject to the following constraints.
	// - both characters of a pair occur in the same isolating run sequence
	// - the closing character of a pair follows the opening character
	// - any bracket character can belong at most to one pair, the earliest possible one
	// - any bracket character not part of a pair is treated like an ordinary character
	// - pairs may nest properly, but their spans may not overlap otherwise

	// Bracket characters with canonical decompositions are supposed to be
	// treated as if they had been normalized, to allow normalized and non-
	// normalized text to give the same result. In this implementation that step
	// is pushed out to the caller. The caller has to ensure that the pairValue
	// slices contain the rune of the opening bracket after normalization for
	// any opening or closing bracket.

	openers *list.List // list of positions for opening brackets

	// bracket pair positions sorted by location of opening bracket
	pairPositions bracketPairs

	codesIsolatedRun []Class // directional bidi codes for an isolated run
	indexes          []int   // array of index values into the original string

***REMOVED***

// matchOpener reports whether characters at given positions form a matching
// bracket pair.
func (p *bracketPairer) matchOpener(pairValues []rune, opener, closer int) bool ***REMOVED***
	return pairValues[p.indexes[opener]] == pairValues[p.indexes[closer]]
***REMOVED***

const maxPairingDepth = 63

// locateBrackets locates matching bracket pairs according to BD16.
//
// This implementation uses a linked list instead of a stack, because, while
// elements are added at the front (like a push) they are not generally removed
// in atomic 'pop' operations, reducing the benefit of the stack archetype.
func (p *bracketPairer) locateBrackets(pairTypes []bracketType, pairValues []rune) ***REMOVED***
	// traverse the run
	// do that explicitly (not in a for-each) so we can record position
	for i, index := range p.indexes ***REMOVED***

		// look at the bracket type for each character
		if pairTypes[index] == bpNone || p.codesIsolatedRun[i] != ON ***REMOVED***
			// continue scanning
			continue
		***REMOVED***
		switch pairTypes[index] ***REMOVED***
		case bpOpen:
			// check if maximum pairing depth reached
			if p.openers.Len() == maxPairingDepth ***REMOVED***
				p.openers.Init()
				return
			***REMOVED***
			// remember opener location, most recent first
			p.openers.PushFront(i)

		case bpClose:
			// see if there is a match
			count := 0
			for elem := p.openers.Front(); elem != nil; elem = elem.Next() ***REMOVED***
				count++
				opener := elem.Value.(int)
				if p.matchOpener(pairValues, opener, i) ***REMOVED***
					// if the opener matches, add nested pair to the ordered list
					p.pairPositions = append(p.pairPositions, bracketPair***REMOVED***opener, i***REMOVED***)
					// remove up to and including matched opener
					for ; count > 0; count-- ***REMOVED***
						p.openers.Remove(p.openers.Front())
					***REMOVED***
					break
				***REMOVED***
			***REMOVED***
			sort.Sort(p.pairPositions)
			// if we get here, the closing bracket matched no openers
			// and gets ignored
		***REMOVED***
	***REMOVED***
***REMOVED***

// Bracket pairs within an isolating run sequence are processed as units so
// that both the opening and the closing paired bracket in a pair resolve to
// the same direction.
//
// N0. Process bracket pairs in an isolating run sequence sequentially in
// the logical order of the text positions of the opening paired brackets
// using the logic given below. Within this scope, bidirectional types EN
// and AN are treated as R.
//
// Identify the bracket pairs in the current isolating run sequence
// according to BD16. For each bracket-pair element in the list of pairs of
// text positions:
//
// a Inspect the bidirectional types of the characters enclosed within the
// bracket pair.
//
// b If any strong type (either L or R) matching the embedding direction is
// found, set the type for both brackets in the pair to match the embedding
// direction.
//
// o [ e ] o -> o e e e o
//
// o [ o e ] -> o e o e e
//
// o [ NI e ] -> o e NI e e
//
// c Otherwise, if a strong type (opposite the embedding direction) is
// found, test for adjacent strong types as follows: 1 First, check
// backwards before the opening paired bracket until the first strong type
// (L, R, or sos) is found. If that first preceding strong type is opposite
// the embedding direction, then set the type for both brackets in the pair
// to that type. 2 Otherwise, set the type for both brackets in the pair to
// the embedding direction.
//
// o [ o ] e -> o o o o e
//
// o [ o NI ] o -> o o o NI o o
//
// e [ o ] o -> e e o e o
//
// e [ o ] e -> e e o e e
//
// e ( o [ o ] NI ) e -> e e o o o o NI e e
//
// d Otherwise, do not set the type for the current bracket pair. Note that
// if the enclosed text contains no strong types the paired brackets will
// both resolve to the same level when resolved individually using rules N1
// and N2.
//
// e ( NI ) o -> e ( NI ) o

// getStrongTypeN0 maps character's directional code to strong type as required
// by rule N0.
//
// TODO: have separate type for "strong" directionality.
func (p *bracketPairer) getStrongTypeN0(index int) Class ***REMOVED***
	switch p.codesIsolatedRun[index] ***REMOVED***
	// in the scope of N0, number types are treated as R
	case EN, AN, AL, R:
		return R
	case L:
		return L
	default:
		return ON
	***REMOVED***
***REMOVED***

// classifyPairContent reports the strong types contained inside a Bracket Pair,
// assuming the given embedding direction.
//
// It returns ON if no strong type is found. If a single strong type is found,
// it returns this type. Otherwise it returns the embedding direction.
//
// TODO: use separate type for "strong" directionality.
func (p *bracketPairer) classifyPairContent(loc bracketPair, dirEmbed Class) Class ***REMOVED***
	dirOpposite := ON
	for i := loc.opener + 1; i < loc.closer; i++ ***REMOVED***
		dir := p.getStrongTypeN0(i)
		if dir == ON ***REMOVED***
			continue
		***REMOVED***
		if dir == dirEmbed ***REMOVED***
			return dir // type matching embedding direction found
		***REMOVED***
		dirOpposite = dir
	***REMOVED***
	// return ON if no strong type found, or class opposite to dirEmbed
	return dirOpposite
***REMOVED***

// classBeforePair determines which strong types are present before a Bracket
// Pair. Return R or L if strong type found, otherwise ON.
func (p *bracketPairer) classBeforePair(loc bracketPair) Class ***REMOVED***
	for i := loc.opener - 1; i >= 0; i-- ***REMOVED***
		if dir := p.getStrongTypeN0(i); dir != ON ***REMOVED***
			return dir
		***REMOVED***
	***REMOVED***
	// no strong types found, return sos
	return p.sos
***REMOVED***

// assignBracketType implements rule N0 for a single bracket pair.
func (p *bracketPairer) assignBracketType(loc bracketPair, dirEmbed Class, initialTypes []Class) ***REMOVED***
	// rule "N0, a", inspect contents of pair
	dirPair := p.classifyPairContent(loc, dirEmbed)

	// dirPair is now L, R, or N (no strong type found)

	// the following logical tests are performed out of order compared to
	// the statement of the rules but yield the same results
	if dirPair == ON ***REMOVED***
		return // case "d" - nothing to do
	***REMOVED***

	if dirPair != dirEmbed ***REMOVED***
		// case "c": strong type found, opposite - check before (c.1)
		dirPair = p.classBeforePair(loc)
		if dirPair == dirEmbed || dirPair == ON ***REMOVED***
			// no strong opposite type found before - use embedding (c.2)
			dirPair = dirEmbed
		***REMOVED***
	***REMOVED***
	// else: case "b", strong type found matching embedding,
	// no explicit action needed, as dirPair is already set to embedding
	// direction

	// set the bracket types to the type found
	p.setBracketsToType(loc, dirPair, initialTypes)
***REMOVED***

func (p *bracketPairer) setBracketsToType(loc bracketPair, dirPair Class, initialTypes []Class) ***REMOVED***
	p.codesIsolatedRun[loc.opener] = dirPair
	p.codesIsolatedRun[loc.closer] = dirPair

	for i := loc.opener + 1; i < loc.closer; i++ ***REMOVED***
		index := p.indexes[i]
		if initialTypes[index] != NSM ***REMOVED***
			break
		***REMOVED***
		p.codesIsolatedRun[i] = dirPair
	***REMOVED***

	for i := loc.closer + 1; i < len(p.indexes); i++ ***REMOVED***
		index := p.indexes[i]
		if initialTypes[index] != NSM ***REMOVED***
			break
		***REMOVED***
		p.codesIsolatedRun[i] = dirPair
	***REMOVED***
***REMOVED***

// resolveBrackets implements rule N0 for a list of pairs.
func (p *bracketPairer) resolveBrackets(dirEmbed Class, initialTypes []Class) ***REMOVED***
	for _, loc := range p.pairPositions ***REMOVED***
		p.assignBracketType(loc, dirEmbed, initialTypes)
	***REMOVED***
***REMOVED***
