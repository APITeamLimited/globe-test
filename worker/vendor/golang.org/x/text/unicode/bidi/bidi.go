// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_trieval.go gen_ranges.go

// Package bidi contains functionality for bidirectional text support.
//
// See https://www.unicode.org/reports/tr9.
//
// NOTE: UNDER CONSTRUCTION. This API may change in backwards incompatible ways
// and without notice.
package bidi // import "golang.org/x/text/unicode/bidi"

// TODO
// - Transformer for reordering?
// - Transformer (validator, really) for Bidi Rule.

import (
	"bytes"
)

// This API tries to avoid dealing with embedding levels for now. Under the hood
// these will be computed, but the question is to which extent the user should
// know they exist. We should at some point allow the user to specify an
// embedding hierarchy, though.

// A Direction indicates the overall flow of text.
type Direction int

const (
	// LeftToRight indicates the text contains no right-to-left characters and
	// that either there are some left-to-right characters or the option
	// DefaultDirection(LeftToRight) was passed.
	LeftToRight Direction = iota

	// RightToLeft indicates the text contains no left-to-right characters and
	// that either there are some right-to-left characters or the option
	// DefaultDirection(RightToLeft) was passed.
	RightToLeft

	// Mixed indicates text contains both left-to-right and right-to-left
	// characters.
	Mixed

	// Neutral means that text contains no left-to-right and right-to-left
	// characters and that no default direction has been set.
	Neutral
)

type options struct ***REMOVED***
	defaultDirection Direction
***REMOVED***

// An Option is an option for Bidi processing.
type Option func(*options)

// ICU allows the user to define embedding levels. This may be used, for example,
// to use hierarchical structure of markup languages to define embeddings.
// The following option may be a way to expose this functionality in this API.
// // LevelFunc sets a function that associates nesting levels with the given text.
// // The levels function will be called with monotonically increasing values for p.
// func LevelFunc(levels func(p int) int) Option ***REMOVED***
// 	panic("unimplemented")
// ***REMOVED***

// DefaultDirection sets the default direction for a Paragraph. The direction is
// overridden if the text contains directional characters.
func DefaultDirection(d Direction) Option ***REMOVED***
	return func(opts *options) ***REMOVED***
		opts.defaultDirection = d
	***REMOVED***
***REMOVED***

// A Paragraph holds a single Paragraph for Bidi processing.
type Paragraph struct ***REMOVED***
	p          []byte
	o          Ordering
	opts       []Option
	types      []Class
	pairTypes  []bracketType
	pairValues []rune
	runes      []rune
	options    options
***REMOVED***

// Initialize the p.pairTypes, p.pairValues and p.types from the input previously
// set by p.SetBytes() or p.SetString(). Also limit the input up to (and including) a paragraph
// separator (bidi class B).
//
// The function p.Order() needs these values to be set, so this preparation could be postponed.
// But since the SetBytes and SetStrings functions return the length of the input up to the paragraph
// separator, the whole input needs to be processed anyway and should not be done twice.
//
// The function has the same return values as SetBytes() / SetString()
func (p *Paragraph) prepareInput() (n int, err error) ***REMOVED***
	p.runes = bytes.Runes(p.p)
	bytecount := 0
	// clear slices from previous SetString or SetBytes
	p.pairTypes = nil
	p.pairValues = nil
	p.types = nil

	for _, r := range p.runes ***REMOVED***
		props, i := LookupRune(r)
		bytecount += i
		cls := props.Class()
		if cls == B ***REMOVED***
			return bytecount, nil
		***REMOVED***
		p.types = append(p.types, cls)
		if props.IsOpeningBracket() ***REMOVED***
			p.pairTypes = append(p.pairTypes, bpOpen)
			p.pairValues = append(p.pairValues, r)
		***REMOVED*** else if props.IsBracket() ***REMOVED***
			// this must be a closing bracket,
			// since IsOpeningBracket is not true
			p.pairTypes = append(p.pairTypes, bpClose)
			p.pairValues = append(p.pairValues, r)
		***REMOVED*** else ***REMOVED***
			p.pairTypes = append(p.pairTypes, bpNone)
			p.pairValues = append(p.pairValues, 0)
		***REMOVED***
	***REMOVED***
	return bytecount, nil
***REMOVED***

// SetBytes configures p for the given paragraph text. It replaces text
// previously set by SetBytes or SetString. If b contains a paragraph separator
// it will only process the first paragraph and report the number of bytes
// consumed from b including this separator. Error may be non-nil if options are
// given.
func (p *Paragraph) SetBytes(b []byte, opts ...Option) (n int, err error) ***REMOVED***
	p.p = b
	p.opts = opts
	return p.prepareInput()
***REMOVED***

// SetString configures s for the given paragraph text. It replaces text
// previously set by SetBytes or SetString. If s contains a paragraph separator
// it will only process the first paragraph and report the number of bytes
// consumed from s including this separator. Error may be non-nil if options are
// given.
func (p *Paragraph) SetString(s string, opts ...Option) (n int, err error) ***REMOVED***
	p.p = []byte(s)
	p.opts = opts
	return p.prepareInput()
***REMOVED***

// IsLeftToRight reports whether the principle direction of rendering for this
// paragraphs is left-to-right. If this returns false, the principle direction
// of rendering is right-to-left.
func (p *Paragraph) IsLeftToRight() bool ***REMOVED***
	return p.Direction() == LeftToRight
***REMOVED***

// Direction returns the direction of the text of this paragraph.
//
// The direction may be LeftToRight, RightToLeft, Mixed, or Neutral.
func (p *Paragraph) Direction() Direction ***REMOVED***
	return p.o.Direction()
***REMOVED***

// TODO: what happens if the position is > len(input)? This should return an error.

// RunAt reports the Run at the given position of the input text.
//
// This method can be used for computing line breaks on paragraphs.
func (p *Paragraph) RunAt(pos int) Run ***REMOVED***
	c := 0
	runNumber := 0
	for i, r := range p.o.runes ***REMOVED***
		c += len(r)
		if pos < c ***REMOVED***
			runNumber = i
		***REMOVED***
	***REMOVED***
	return p.o.Run(runNumber)
***REMOVED***

func calculateOrdering(levels []level, runes []rune) Ordering ***REMOVED***
	var curDir Direction

	prevDir := Neutral
	prevI := 0

	o := Ordering***REMOVED******REMOVED***
	// lvl = 0,2,4,...: left to right
	// lvl = 1,3,5,...: right to left
	for i, lvl := range levels ***REMOVED***
		if lvl%2 == 0 ***REMOVED***
			curDir = LeftToRight
		***REMOVED*** else ***REMOVED***
			curDir = RightToLeft
		***REMOVED***
		if curDir != prevDir ***REMOVED***
			if i > 0 ***REMOVED***
				o.runes = append(o.runes, runes[prevI:i])
				o.directions = append(o.directions, prevDir)
				o.startpos = append(o.startpos, prevI)
			***REMOVED***
			prevI = i
			prevDir = curDir
		***REMOVED***
	***REMOVED***
	o.runes = append(o.runes, runes[prevI:])
	o.directions = append(o.directions, prevDir)
	o.startpos = append(o.startpos, prevI)
	return o
***REMOVED***

// Order computes the visual ordering of all the runs in a Paragraph.
func (p *Paragraph) Order() (Ordering, error) ***REMOVED***
	if len(p.types) == 0 ***REMOVED***
		return Ordering***REMOVED******REMOVED***, nil
	***REMOVED***

	for _, fn := range p.opts ***REMOVED***
		fn(&p.options)
	***REMOVED***
	lvl := level(-1)
	if p.options.defaultDirection == RightToLeft ***REMOVED***
		lvl = 1
	***REMOVED***
	para, err := newParagraph(p.types, p.pairTypes, p.pairValues, lvl)
	if err != nil ***REMOVED***
		return Ordering***REMOVED******REMOVED***, err
	***REMOVED***

	levels := para.getLevels([]int***REMOVED***len(p.types)***REMOVED***)

	p.o = calculateOrdering(levels, p.runes)
	return p.o, nil
***REMOVED***

// Line computes the visual ordering of runs for a single line starting and
// ending at the given positions in the original text.
func (p *Paragraph) Line(start, end int) (Ordering, error) ***REMOVED***
	lineTypes := p.types[start:end]
	para, err := newParagraph(lineTypes, p.pairTypes[start:end], p.pairValues[start:end], -1)
	if err != nil ***REMOVED***
		return Ordering***REMOVED******REMOVED***, err
	***REMOVED***
	levels := para.getLevels([]int***REMOVED***len(lineTypes)***REMOVED***)
	o := calculateOrdering(levels, p.runes[start:end])
	return o, nil
***REMOVED***

// An Ordering holds the computed visual order of runs of a Paragraph. Calling
// SetBytes or SetString on the originating Paragraph invalidates an Ordering.
// The methods of an Ordering should only be called by one goroutine at a time.
type Ordering struct ***REMOVED***
	runes      [][]rune
	directions []Direction
	startpos   []int
***REMOVED***

// Direction reports the directionality of the runs.
//
// The direction may be LeftToRight, RightToLeft, Mixed, or Neutral.
func (o *Ordering) Direction() Direction ***REMOVED***
	return o.directions[0]
***REMOVED***

// NumRuns returns the number of runs.
func (o *Ordering) NumRuns() int ***REMOVED***
	return len(o.runes)
***REMOVED***

// Run returns the ith run within the ordering.
func (o *Ordering) Run(i int) Run ***REMOVED***
	r := Run***REMOVED***
		runes:     o.runes[i],
		direction: o.directions[i],
		startpos:  o.startpos[i],
	***REMOVED***
	return r
***REMOVED***

// TODO: perhaps with options.
// // Reorder creates a reader that reads the runes in visual order per character.
// // Modifiers remain after the runes they modify.
// func (l *Runs) Reorder() io.Reader ***REMOVED***
// 	panic("unimplemented")
// ***REMOVED***

// A Run is a continuous sequence of characters of a single direction.
type Run struct ***REMOVED***
	runes     []rune
	direction Direction
	startpos  int
***REMOVED***

// String returns the text of the run in its original order.
func (r *Run) String() string ***REMOVED***
	return string(r.runes)
***REMOVED***

// Bytes returns the text of the run in its original order.
func (r *Run) Bytes() []byte ***REMOVED***
	return []byte(r.String())
***REMOVED***

// TODO: methods for
// - Display order
// - headers and footers
// - bracket replacement.

// Direction reports the direction of the run.
func (r *Run) Direction() Direction ***REMOVED***
	return r.direction
***REMOVED***

// Pos returns the position of the Run within the text passed to SetBytes or SetString of the
// originating Paragraph value.
func (r *Run) Pos() (start, end int) ***REMOVED***
	return r.startpos, r.startpos + len(r.runes) - 1
***REMOVED***

// AppendReverse reverses the order of characters of in, appends them to out,
// and returns the result. Modifiers will still follow the runes they modify.
// Brackets are replaced with their counterparts.
func AppendReverse(out, in []byte) []byte ***REMOVED***
	ret := make([]byte, len(in)+len(out))
	copy(ret, out)
	inRunes := bytes.Runes(in)

	for i, r := range inRunes ***REMOVED***
		prop, _ := LookupRune(r)
		if prop.IsBracket() ***REMOVED***
			inRunes[i] = prop.reverseBracket(r)
		***REMOVED***
	***REMOVED***

	for i, j := 0, len(inRunes)-1; i < j; i, j = i+1, j-1 ***REMOVED***
		inRunes[i], inRunes[j] = inRunes[j], inRunes[i]
	***REMOVED***
	copy(ret[len(out):], string(inRunes))

	return ret
***REMOVED***

// ReverseString reverses the order of characters in s and returns a new string.
// Modifiers will still follow the runes they modify. Brackets are replaced with
// their counterparts.
func ReverseString(s string) string ***REMOVED***
	input := []rune(s)
	li := len(input)
	ret := make([]rune, li)
	for i, r := range input ***REMOVED***
		prop, _ := LookupRune(r)
		if prop.IsBracket() ***REMOVED***
			ret[li-i-1] = prop.reverseBracket(r)
		***REMOVED*** else ***REMOVED***
			ret[li-i-1] = r
		***REMOVED***
	***REMOVED***
	return string(ret)
***REMOVED***
