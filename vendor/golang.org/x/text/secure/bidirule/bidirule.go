// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bidirule implements the Bidi Rule defined by RFC 5893.
//
// This package is under development. The API may change without notice and
// without preserving backward compatibility.
package bidirule

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/bidi"
)

// This file contains an implementation of RFC 5893: Right-to-Left Scripts for
// Internationalized Domain Names for Applications (IDNA)
//
// A label is an individual component of a domain name.  Labels are usually
// shown separated by dots; for example, the domain name "www.example.com" is
// composed of three labels: "www", "example", and "com".
//
// An RTL label is a label that contains at least one character of class R, AL,
// or AN. An LTR label is any label that is not an RTL label.
//
// A "Bidi domain name" is a domain name that contains at least one RTL label.
//
//  The following guarantees can be made based on the above:
//
//  o  In a domain name consisting of only labels that satisfy the rule,
//     the requirements of Section 3 are satisfied.  Note that even LTR
//     labels and pure ASCII labels have to be tested.
//
//  o  In a domain name consisting of only LDH labels (as defined in the
//     Definitions document [RFC5890]) and labels that satisfy the rule,
//     the requirements of Section 3 are satisfied as long as a label
//     that starts with an ASCII digit does not come after a
//     right-to-left label.
//
//  No guarantee is given for other combinations.

// ErrInvalid indicates a label is invalid according to the Bidi Rule.
var ErrInvalid = errors.New("bidirule: failed Bidi Rule")

type ruleState uint8

const (
	ruleInitial ruleState = iota
	ruleLTR
	ruleLTRFinal
	ruleRTL
	ruleRTLFinal
	ruleInvalid
)

type ruleTransition struct ***REMOVED***
	next ruleState
	mask uint16
***REMOVED***

var transitions = [...][2]ruleTransition***REMOVED***
	// [2.1] The first character must be a character with Bidi property L, R, or
	// AL. If it has the R or AL property, it is an RTL label; if it has the L
	// property, it is an LTR label.
	ruleInitial: ***REMOVED***
		***REMOVED***ruleLTRFinal, 1 << bidi.L***REMOVED***,
		***REMOVED***ruleRTLFinal, 1<<bidi.R | 1<<bidi.AL***REMOVED***,
	***REMOVED***,
	ruleRTL: ***REMOVED***
		// [2.3] In an RTL label, the end of the label must be a character with
		// Bidi property R, AL, EN, or AN, followed by zero or more characters
		// with Bidi property NSM.
		***REMOVED***ruleRTLFinal, 1<<bidi.R | 1<<bidi.AL | 1<<bidi.EN | 1<<bidi.AN***REMOVED***,

		// [2.2] In an RTL label, only characters with the Bidi properties R,
		// AL, AN, EN, ES, CS, ET, ON, BN, or NSM are allowed.
		// We exclude the entries from [2.3]
		***REMOVED***ruleRTL, 1<<bidi.ES | 1<<bidi.CS | 1<<bidi.ET | 1<<bidi.ON | 1<<bidi.BN | 1<<bidi.NSM***REMOVED***,
	***REMOVED***,
	ruleRTLFinal: ***REMOVED***
		// [2.3] In an RTL label, the end of the label must be a character with
		// Bidi property R, AL, EN, or AN, followed by zero or more characters
		// with Bidi property NSM.
		***REMOVED***ruleRTLFinal, 1<<bidi.R | 1<<bidi.AL | 1<<bidi.EN | 1<<bidi.AN | 1<<bidi.NSM***REMOVED***,

		// [2.2] In an RTL label, only characters with the Bidi properties R,
		// AL, AN, EN, ES, CS, ET, ON, BN, or NSM are allowed.
		// We exclude the entries from [2.3] and NSM.
		***REMOVED***ruleRTL, 1<<bidi.ES | 1<<bidi.CS | 1<<bidi.ET | 1<<bidi.ON | 1<<bidi.BN***REMOVED***,
	***REMOVED***,
	ruleLTR: ***REMOVED***
		// [2.6] In an LTR label, the end of the label must be a character with
		// Bidi property L or EN, followed by zero or more characters with Bidi
		// property NSM.
		***REMOVED***ruleLTRFinal, 1<<bidi.L | 1<<bidi.EN***REMOVED***,

		// [2.5] In an LTR label, only characters with the Bidi properties L,
		// EN, ES, CS, ET, ON, BN, or NSM are allowed.
		// We exclude the entries from [2.6].
		***REMOVED***ruleLTR, 1<<bidi.ES | 1<<bidi.CS | 1<<bidi.ET | 1<<bidi.ON | 1<<bidi.BN | 1<<bidi.NSM***REMOVED***,
	***REMOVED***,
	ruleLTRFinal: ***REMOVED***
		// [2.6] In an LTR label, the end of the label must be a character with
		// Bidi property L or EN, followed by zero or more characters with Bidi
		// property NSM.
		***REMOVED***ruleLTRFinal, 1<<bidi.L | 1<<bidi.EN | 1<<bidi.NSM***REMOVED***,

		// [2.5] In an LTR label, only characters with the Bidi properties L,
		// EN, ES, CS, ET, ON, BN, or NSM are allowed.
		// We exclude the entries from [2.6].
		***REMOVED***ruleLTR, 1<<bidi.ES | 1<<bidi.CS | 1<<bidi.ET | 1<<bidi.ON | 1<<bidi.BN***REMOVED***,
	***REMOVED***,
	ruleInvalid: ***REMOVED***
		***REMOVED***ruleInvalid, 0***REMOVED***,
		***REMOVED***ruleInvalid, 0***REMOVED***,
	***REMOVED***,
***REMOVED***

// [2.4] In an RTL label, if an EN is present, no AN may be present, and
// vice versa.
const exclusiveRTL = uint16(1<<bidi.EN | 1<<bidi.AN)

// From RFC 5893
// An RTL label is a label that contains at least one character of type
// R, AL, or AN.
//
// An LTR label is any label that is not an RTL label.

// Direction reports the direction of the given label as defined by RFC 5893.
// The Bidi Rule does not have to be applied to labels of the category
// LeftToRight.
func Direction(b []byte) bidi.Direction ***REMOVED***
	for i := 0; i < len(b); ***REMOVED***
		e, sz := bidi.Lookup(b[i:])
		if sz == 0 ***REMOVED***
			i++
		***REMOVED***
		c := e.Class()
		if c == bidi.R || c == bidi.AL || c == bidi.AN ***REMOVED***
			return bidi.RightToLeft
		***REMOVED***
		i += sz
	***REMOVED***
	return bidi.LeftToRight
***REMOVED***

// DirectionString reports the direction of the given label as defined by RFC
// 5893. The Bidi Rule does not have to be applied to labels of the category
// LeftToRight.
func DirectionString(s string) bidi.Direction ***REMOVED***
	for i := 0; i < len(s); ***REMOVED***
		e, sz := bidi.LookupString(s[i:])
		if sz == 0 ***REMOVED***
			i++
			continue
		***REMOVED***
		c := e.Class()
		if c == bidi.R || c == bidi.AL || c == bidi.AN ***REMOVED***
			return bidi.RightToLeft
		***REMOVED***
		i += sz
	***REMOVED***
	return bidi.LeftToRight
***REMOVED***

// Valid reports whether b conforms to the BiDi rule.
func Valid(b []byte) bool ***REMOVED***
	var t Transformer
	if n, ok := t.advance(b); !ok || n < len(b) ***REMOVED***
		return false
	***REMOVED***
	return t.isFinal()
***REMOVED***

// ValidString reports whether s conforms to the BiDi rule.
func ValidString(s string) bool ***REMOVED***
	var t Transformer
	if n, ok := t.advanceString(s); !ok || n < len(s) ***REMOVED***
		return false
	***REMOVED***
	return t.isFinal()
***REMOVED***

// New returns a Transformer that verifies that input adheres to the Bidi Rule.
func New() *Transformer ***REMOVED***
	return &Transformer***REMOVED******REMOVED***
***REMOVED***

// Transformer implements transform.Transform.
type Transformer struct ***REMOVED***
	state  ruleState
	hasRTL bool
	seen   uint16
***REMOVED***

// A rule can only be violated for "Bidi Domain names", meaning if one of the
// following categories has been observed.
func (t *Transformer) isRTL() bool ***REMOVED***
	const isRTL = 1<<bidi.R | 1<<bidi.AL | 1<<bidi.AN
	return t.seen&isRTL != 0
***REMOVED***

// Reset implements transform.Transformer.
func (t *Transformer) Reset() ***REMOVED*** *t = Transformer***REMOVED******REMOVED*** ***REMOVED***

// Transform implements transform.Transformer. This Transformer has state and
// needs to be reset between uses.
func (t *Transformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	if len(dst) < len(src) ***REMOVED***
		src = src[:len(dst)]
		atEOF = false
		err = transform.ErrShortDst
	***REMOVED***
	n, err1 := t.Span(src, atEOF)
	copy(dst, src[:n])
	if err == nil || err1 != nil && err1 != transform.ErrShortSrc ***REMOVED***
		err = err1
	***REMOVED***
	return n, n, err
***REMOVED***

// Span returns the first n bytes of src that conform to the Bidi rule.
func (t *Transformer) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	if t.state == ruleInvalid && t.isRTL() ***REMOVED***
		return 0, ErrInvalid
	***REMOVED***
	n, ok := t.advance(src)
	switch ***REMOVED***
	case !ok:
		err = ErrInvalid
	case n < len(src):
		if !atEOF ***REMOVED***
			err = transform.ErrShortSrc
			break
		***REMOVED***
		err = ErrInvalid
	case !t.isFinal():
		err = ErrInvalid
	***REMOVED***
	return n, err
***REMOVED***

// Precomputing the ASCII values decreases running time for the ASCII fast path
// by about 30%.
var asciiTable [128]bidi.Properties

func init() ***REMOVED***
	for i := range asciiTable ***REMOVED***
		p, _ := bidi.LookupRune(rune(i))
		asciiTable[i] = p
	***REMOVED***
***REMOVED***

func (t *Transformer) advance(s []byte) (n int, ok bool) ***REMOVED***
	var e bidi.Properties
	var sz int
	for n < len(s) ***REMOVED***
		if s[n] < utf8.RuneSelf ***REMOVED***
			e, sz = asciiTable[s[n]], 1
		***REMOVED*** else ***REMOVED***
			e, sz = bidi.Lookup(s[n:])
			if sz <= 1 ***REMOVED***
				if sz == 1 ***REMOVED***
					// We always consider invalid UTF-8 to be invalid, even if
					// the string has not yet been determined to be RTL.
					// TODO: is this correct?
					return n, false
				***REMOVED***
				return n, true // incomplete UTF-8 encoding
			***REMOVED***
		***REMOVED***
		// TODO: using CompactClass would result in noticeable speedup.
		// See unicode/bidi/prop.go:Properties.CompactClass.
		c := uint16(1 << e.Class())
		t.seen |= c
		if t.seen&exclusiveRTL == exclusiveRTL ***REMOVED***
			t.state = ruleInvalid
			return n, false
		***REMOVED***
		switch tr := transitions[t.state]; ***REMOVED***
		case tr[0].mask&c != 0:
			t.state = tr[0].next
		case tr[1].mask&c != 0:
			t.state = tr[1].next
		default:
			t.state = ruleInvalid
			if t.isRTL() ***REMOVED***
				return n, false
			***REMOVED***
		***REMOVED***
		n += sz
	***REMOVED***
	return n, true
***REMOVED***

func (t *Transformer) advanceString(s string) (n int, ok bool) ***REMOVED***
	var e bidi.Properties
	var sz int
	for n < len(s) ***REMOVED***
		if s[n] < utf8.RuneSelf ***REMOVED***
			e, sz = asciiTable[s[n]], 1
		***REMOVED*** else ***REMOVED***
			e, sz = bidi.LookupString(s[n:])
			if sz <= 1 ***REMOVED***
				if sz == 1 ***REMOVED***
					return n, false // invalid UTF-8
				***REMOVED***
				return n, true // incomplete UTF-8 encoding
			***REMOVED***
		***REMOVED***
		// TODO: using CompactClass results in noticeable speedup.
		// See unicode/bidi/prop.go:Properties.CompactClass.
		c := uint16(1 << e.Class())
		t.seen |= c
		if t.seen&exclusiveRTL == exclusiveRTL ***REMOVED***
			t.state = ruleInvalid
			return n, false
		***REMOVED***
		switch tr := transitions[t.state]; ***REMOVED***
		case tr[0].mask&c != 0:
			t.state = tr[0].next
		case tr[1].mask&c != 0:
			t.state = tr[1].next
		default:
			t.state = ruleInvalid
			if t.isRTL() ***REMOVED***
				return n, false
			***REMOVED***
		***REMOVED***
		n += sz
	***REMOVED***
	return n, true
***REMOVED***
