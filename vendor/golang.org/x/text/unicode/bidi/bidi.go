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

// TODO:
// The following functionality would not be hard to implement, but hinges on
// the definition of a Segmenter interface. For now this is up to the user.
// - Iterate over paragraphs
// - Segmenter to iterate over runs directly from a given text.
// Also:
// - Transformer for reordering?
// - Transformer (validator, really) for Bidi Rule.

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

type options struct***REMOVED******REMOVED***

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
	panic("unimplemented")
***REMOVED***

// A Paragraph holds a single Paragraph for Bidi processing.
type Paragraph struct ***REMOVED***
	// buffers
***REMOVED***

// SetBytes configures p for the given paragraph text. It replaces text
// previously set by SetBytes or SetString. If b contains a paragraph separator
// it will only process the first paragraph and report the number of bytes
// consumed from b including this separator. Error may be non-nil if options are
// given.
func (p *Paragraph) SetBytes(b []byte, opts ...Option) (n int, err error) ***REMOVED***
	panic("unimplemented")
***REMOVED***

// SetString configures p for the given paragraph text. It replaces text
// previously set by SetBytes or SetString. If b contains a paragraph separator
// it will only process the first paragraph and report the number of bytes
// consumed from b including this separator. Error may be non-nil if options are
// given.
func (p *Paragraph) SetString(s string, opts ...Option) (n int, err error) ***REMOVED***
	panic("unimplemented")
***REMOVED***

// IsLeftToRight reports whether the principle direction of rendering for this
// paragraphs is left-to-right. If this returns false, the principle direction
// of rendering is right-to-left.
func (p *Paragraph) IsLeftToRight() bool ***REMOVED***
	panic("unimplemented")
***REMOVED***

// Direction returns the direction of the text of this paragraph.
//
// The direction may be LeftToRight, RightToLeft, Mixed, or Neutral.
func (p *Paragraph) Direction() Direction ***REMOVED***
	panic("unimplemented")
***REMOVED***

// RunAt reports the Run at the given position of the input text.
//
// This method can be used for computing line breaks on paragraphs.
func (p *Paragraph) RunAt(pos int) Run ***REMOVED***
	panic("unimplemented")
***REMOVED***

// Order computes the visual ordering of all the runs in a Paragraph.
func (p *Paragraph) Order() (Ordering, error) ***REMOVED***
	panic("unimplemented")
***REMOVED***

// Line computes the visual ordering of runs for a single line starting and
// ending at the given positions in the original text.
func (p *Paragraph) Line(start, end int) (Ordering, error) ***REMOVED***
	panic("unimplemented")
***REMOVED***

// An Ordering holds the computed visual order of runs of a Paragraph. Calling
// SetBytes or SetString on the originating Paragraph invalidates an Ordering.
// The methods of an Ordering should only be called by one goroutine at a time.
type Ordering struct***REMOVED******REMOVED***

// Direction reports the directionality of the runs.
//
// The direction may be LeftToRight, RightToLeft, Mixed, or Neutral.
func (o *Ordering) Direction() Direction ***REMOVED***
	panic("unimplemented")
***REMOVED***

// NumRuns returns the number of runs.
func (o *Ordering) NumRuns() int ***REMOVED***
	panic("unimplemented")
***REMOVED***

// Run returns the ith run within the ordering.
func (o *Ordering) Run(i int) Run ***REMOVED***
	panic("unimplemented")
***REMOVED***

// TODO: perhaps with options.
// // Reorder creates a reader that reads the runes in visual order per character.
// // Modifiers remain after the runes they modify.
// func (l *Runs) Reorder() io.Reader ***REMOVED***
// 	panic("unimplemented")
// ***REMOVED***

// A Run is a continuous sequence of characters of a single direction.
type Run struct ***REMOVED***
***REMOVED***

// String returns the text of the run in its original order.
func (r *Run) String() string ***REMOVED***
	panic("unimplemented")
***REMOVED***

// Bytes returns the text of the run in its original order.
func (r *Run) Bytes() []byte ***REMOVED***
	panic("unimplemented")
***REMOVED***

// TODO: methods for
// - Display order
// - headers and footers
// - bracket replacement.

// Direction reports the direction of the run.
func (r *Run) Direction() Direction ***REMOVED***
	panic("unimplemented")
***REMOVED***

// Position of the Run within the text passed to SetBytes or SetString of the
// originating Paragraph value.
func (r *Run) Pos() (start, end int) ***REMOVED***
	panic("unimplemented")
***REMOVED***

// AppendReverse reverses the order of characters of in, appends them to out,
// and returns the result. Modifiers will still follow the runes they modify.
// Brackets are replaced with their counterparts.
func AppendReverse(out, in []byte) []byte ***REMOVED***
	panic("unimplemented")
***REMOVED***

// ReverseString reverses the order of characters in s and returns a new string.
// Modifiers will still follow the runes they modify. Brackets are replaced with
// their counterparts.
func ReverseString(s string) string ***REMOVED***
	panic("unimplemented")
***REMOVED***
