// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_trieval.go

// Package cases provides general and language-specific case mappers.
package cases // import "golang.org/x/text/cases"

import (
	"golang.org/x/text/language"
	"golang.org/x/text/transform"
)

// References:
// - Unicode Reference Manual Chapter 3.13, 4.2, and 5.18.
// - https://www.unicode.org/reports/tr29/
// - https://www.unicode.org/Public/6.3.0/ucd/CaseFolding.txt
// - https://www.unicode.org/Public/6.3.0/ucd/SpecialCasing.txt
// - https://www.unicode.org/Public/6.3.0/ucd/DerivedCoreProperties.txt
// - https://www.unicode.org/Public/6.3.0/ucd/auxiliary/WordBreakProperty.txt
// - https://www.unicode.org/Public/6.3.0/ucd/auxiliary/WordBreakTest.txt
// - http://userguide.icu-project.org/transforms/casemappings

// TODO:
// - Case folding
// - Wide and Narrow?
// - Segmenter option for title casing.
// - ASCII fast paths
// - Encode Soft-Dotted property within trie somehow.

// A Caser transforms given input to a certain case. It implements
// transform.Transformer.
//
// A Caser may be stateful and should therefore not be shared between
// goroutines.
type Caser struct ***REMOVED***
	t transform.SpanningTransformer
***REMOVED***

// Bytes returns a new byte slice with the result of converting b to the case
// form implemented by c.
func (c Caser) Bytes(b []byte) []byte ***REMOVED***
	b, _, _ = transform.Bytes(c.t, b)
	return b
***REMOVED***

// String returns a string with the result of transforming s to the case form
// implemented by c.
func (c Caser) String(s string) string ***REMOVED***
	s, _, _ = transform.String(c.t, s)
	return s
***REMOVED***

// Reset resets the Caser to be reused for new input after a previous call to
// Transform.
func (c Caser) Reset() ***REMOVED*** c.t.Reset() ***REMOVED***

// Transform implements the transform.Transformer interface and transforms the
// given input to the case form implemented by c.
func (c Caser) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	return c.t.Transform(dst, src, atEOF)
***REMOVED***

// Span implements the transform.SpanningTransformer interface.
func (c Caser) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	return c.t.Span(src, atEOF)
***REMOVED***

// Upper returns a Caser for language-specific uppercasing.
func Upper(t language.Tag, opts ...Option) Caser ***REMOVED***
	return Caser***REMOVED***makeUpper(t, getOpts(opts...))***REMOVED***
***REMOVED***

// Lower returns a Caser for language-specific lowercasing.
func Lower(t language.Tag, opts ...Option) Caser ***REMOVED***
	return Caser***REMOVED***makeLower(t, getOpts(opts...))***REMOVED***
***REMOVED***

// Title returns a Caser for language-specific title casing. It uses an
// approximation of the default Unicode Word Break algorithm.
func Title(t language.Tag, opts ...Option) Caser ***REMOVED***
	return Caser***REMOVED***makeTitle(t, getOpts(opts...))***REMOVED***
***REMOVED***

// Fold returns a Caser that implements Unicode case folding. The returned Caser
// is stateless and safe to use concurrently by multiple goroutines.
//
// Case folding does not normalize the input and may not preserve a normal form.
// Use the collate or search package for more convenient and linguistically
// sound comparisons. Use golang.org/x/text/secure/precis for string comparisons
// where security aspects are a concern.
func Fold(opts ...Option) Caser ***REMOVED***
	return Caser***REMOVED***makeFold(getOpts(opts...))***REMOVED***
***REMOVED***

// An Option is used to modify the behavior of a Caser.
type Option func(o options) options

// TODO: consider these options to take a boolean as well, like FinalSigma.
// The advantage of using this approach is that other providers of a lower-case
// algorithm could set different defaults by prefixing a user-provided slice
// of options with their own. This is handy, for instance, for the precis
// package which would override the default to not handle the Greek final sigma.

var (
	// NoLower disables the lowercasing of non-leading letters for a title
	// caser.
	NoLower Option = noLower

	// Compact omits mappings in case folding for characters that would grow the
	// input. (Unimplemented.)
	Compact Option = compact
)

// TODO: option to preserve a normal form, if applicable?

type options struct ***REMOVED***
	noLower bool
	simple  bool

	// TODO: segmenter, max ignorable, alternative versions, etc.

	ignoreFinalSigma bool
***REMOVED***

func getOpts(o ...Option) (res options) ***REMOVED***
	for _, f := range o ***REMOVED***
		res = f(res)
	***REMOVED***
	return
***REMOVED***

func noLower(o options) options ***REMOVED***
	o.noLower = true
	return o
***REMOVED***

func compact(o options) options ***REMOVED***
	o.simple = true
	return o
***REMOVED***

// HandleFinalSigma specifies whether the special handling of Greek final sigma
// should be enabled. Unicode prescribes handling the Greek final sigma for all
// locales, but standards like IDNA and PRECIS override this default.
func HandleFinalSigma(enable bool) Option ***REMOVED***
	if enable ***REMOVED***
		return handleFinalSigma
	***REMOVED***
	return ignoreFinalSigma
***REMOVED***

func ignoreFinalSigma(o options) options ***REMOVED***
	o.ignoreFinalSigma = true
	return o
***REMOVED***

func handleFinalSigma(o options) options ***REMOVED***
	o.ignoreFinalSigma = false
	return o
***REMOVED***
