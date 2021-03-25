// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate

import (
	"sort"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// newCollator creates a new collator with default options configured.
func newCollator(t colltab.Weighter) *Collator ***REMOVED***
	// Initialize a collator with default options.
	c := &Collator***REMOVED***
		options: options***REMOVED***
			ignore: [colltab.NumLevels]bool***REMOVED***
				colltab.Quaternary: true,
				colltab.Identity:   true,
			***REMOVED***,
			f: norm.NFD,
			t: t,
		***REMOVED***,
	***REMOVED***

	// TODO: store vt in tags or remove.
	c.variableTop = t.Top()

	return c
***REMOVED***

// An Option is used to change the behavior of a Collator. Options override the
// settings passed through the locale identifier.
type Option struct ***REMOVED***
	priority int
	f        func(o *options)
***REMOVED***

type prioritizedOptions []Option

func (p prioritizedOptions) Len() int ***REMOVED***
	return len(p)
***REMOVED***

func (p prioritizedOptions) Swap(i, j int) ***REMOVED***
	p[i], p[j] = p[j], p[i]
***REMOVED***

func (p prioritizedOptions) Less(i, j int) bool ***REMOVED***
	return p[i].priority < p[j].priority
***REMOVED***

type options struct ***REMOVED***
	// ignore specifies which levels to ignore.
	ignore [colltab.NumLevels]bool

	// caseLevel is true if there is an additional level of case matching
	// between the secondary and tertiary levels.
	caseLevel bool

	// backwards specifies the order of sorting at the secondary level.
	// This option exists predominantly to support reverse sorting of accents in French.
	backwards bool

	// numeric specifies whether any sequence of decimal digits (category is Nd)
	// is sorted at a primary level with its numeric value.
	// For example, "A-21" < "A-123".
	// This option is set by wrapping the main Weighter with NewNumericWeighter.
	numeric bool

	// alternate specifies an alternative handling of variables.
	alternate alternateHandling

	// variableTop is the largest primary value that is considered to be
	// variable.
	variableTop uint32

	t colltab.Weighter

	f norm.Form
***REMOVED***

func (o *options) setOptions(opts []Option) ***REMOVED***
	sort.Sort(prioritizedOptions(opts))
	for _, x := range opts ***REMOVED***
		x.f(o)
	***REMOVED***
***REMOVED***

// OptionsFromTag extracts the BCP47 collation options from the tag and
// configures a collator accordingly. These options are set before any other
// option.
func OptionsFromTag(t language.Tag) Option ***REMOVED***
	return Option***REMOVED***0, func(o *options) ***REMOVED***
		o.setFromTag(t)
	***REMOVED******REMOVED***
***REMOVED***

func (o *options) setFromTag(t language.Tag) ***REMOVED***
	o.caseLevel = ldmlBool(t, o.caseLevel, "kc")
	o.backwards = ldmlBool(t, o.backwards, "kb")
	o.numeric = ldmlBool(t, o.numeric, "kn")

	// Extract settings from the BCP47 u extension.
	switch t.TypeForKey("ks") ***REMOVED*** // strength
	case "level1":
		o.ignore[colltab.Secondary] = true
		o.ignore[colltab.Tertiary] = true
	case "level2":
		o.ignore[colltab.Tertiary] = true
	case "level3", "":
		// The default.
	case "level4":
		o.ignore[colltab.Quaternary] = false
	case "identic":
		o.ignore[colltab.Quaternary] = false
		o.ignore[colltab.Identity] = false
	***REMOVED***

	switch t.TypeForKey("ka") ***REMOVED***
	case "shifted":
		o.alternate = altShifted
	// The following two types are not official BCP47, but we support them to
	// give access to this otherwise hidden functionality. The name blanked is
	// derived from the LDML name blanked and posix reflects the main use of
	// the shift-trimmed option.
	case "blanked":
		o.alternate = altBlanked
	case "posix":
		o.alternate = altShiftTrimmed
	***REMOVED***

	// TODO: caseFirst ("kf"), reorder ("kr"), and maybe variableTop ("vt").

	// Not used:
	// - normalization ("kk", not necessary for this implementation)
	// - hiraganaQuatenary ("kh", obsolete)
***REMOVED***

func ldmlBool(t language.Tag, old bool, key string) bool ***REMOVED***
	switch t.TypeForKey(key) ***REMOVED***
	case "true":
		return true
	case "false":
		return false
	default:
		return old
	***REMOVED***
***REMOVED***

var (
	// IgnoreCase sets case-insensitive comparison.
	IgnoreCase Option = ignoreCase
	ignoreCase        = Option***REMOVED***3, ignoreCaseF***REMOVED***

	// IgnoreDiacritics causes diacritical marks to be ignored. ("o" == "ö").
	IgnoreDiacritics Option = ignoreDiacritics
	ignoreDiacritics        = Option***REMOVED***3, ignoreDiacriticsF***REMOVED***

	// IgnoreWidth causes full-width characters to match their half-width
	// equivalents.
	IgnoreWidth Option = ignoreWidth
	ignoreWidth        = Option***REMOVED***2, ignoreWidthF***REMOVED***

	// Loose sets the collator to ignore diacritics, case and width.
	Loose Option = loose
	loose        = Option***REMOVED***4, looseF***REMOVED***

	// Force ordering if strings are equivalent but not equal.
	Force Option = force
	force        = Option***REMOVED***5, forceF***REMOVED***

	// Numeric specifies that numbers should sort numerically ("2" < "12").
	Numeric Option = numeric
	numeric        = Option***REMOVED***5, numericF***REMOVED***
)

func ignoreWidthF(o *options) ***REMOVED***
	o.ignore[colltab.Tertiary] = true
	o.caseLevel = true
***REMOVED***

func ignoreDiacriticsF(o *options) ***REMOVED***
	o.ignore[colltab.Secondary] = true
***REMOVED***

func ignoreCaseF(o *options) ***REMOVED***
	o.ignore[colltab.Tertiary] = true
	o.caseLevel = false
***REMOVED***

func looseF(o *options) ***REMOVED***
	ignoreWidthF(o)
	ignoreDiacriticsF(o)
	ignoreCaseF(o)
***REMOVED***

func forceF(o *options) ***REMOVED***
	o.ignore[colltab.Identity] = false
***REMOVED***

func numericF(o *options) ***REMOVED*** o.numeric = true ***REMOVED***

// Reorder overrides the pre-defined ordering of scripts and character sets.
func Reorder(s ...string) Option ***REMOVED***
	// TODO: need fractional weights to implement this.
	panic("TODO: implement")
***REMOVED***

// TODO: consider making these public again. These options cannot be fully
// specified in BCP47, so an API interface seems warranted. Still a higher-level
// interface would be nice (e.g. a POSIX option for enabling altShiftTrimmed)

// alternateHandling identifies the various ways in which variables are handled.
// A rune with a primary weight lower than the variable top is considered a
// variable.
// See https://www.unicode.org/reports/tr10/#Variable_Weighting for details.
type alternateHandling int

const (
	// altNonIgnorable turns off special handling of variables.
	altNonIgnorable alternateHandling = iota

	// altBlanked sets variables and all subsequent primary ignorables to be
	// ignorable at all levels. This is identical to removing all variables
	// and subsequent primary ignorables from the input.
	altBlanked

	// altShifted sets variables to be ignorable for levels one through three and
	// adds a fourth level based on the values of the ignored levels.
	altShifted

	// altShiftTrimmed is a slight variant of altShifted that is used to
	// emulate POSIX.
	altShiftTrimmed
)
