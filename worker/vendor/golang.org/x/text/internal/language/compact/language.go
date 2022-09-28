// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_index.go -output tables.go
//go:generate go run gen_parents.go

package compact

// TODO: Remove above NOTE after:
// - verifying that tables are dropped correctly (most notably matcher tables).

import (
	"strings"

	"golang.org/x/text/internal/language"
)

// Tag represents a BCP 47 language tag. It is used to specify an instance of a
// specific language or locale. All language tag values are guaranteed to be
// well-formed.
type Tag struct ***REMOVED***
	// NOTE: exported tags will become part of the public API.
	language ID
	locale   ID
	full     fullTag // always a language.Tag for now.
***REMOVED***

const _und = 0

type fullTag interface ***REMOVED***
	IsRoot() bool
	Parent() language.Tag
***REMOVED***

// Make a compact Tag from a fully specified internal language Tag.
func Make(t language.Tag) (tag Tag) ***REMOVED***
	if region := t.TypeForKey("rg"); len(region) == 6 && region[2:] == "zzzz" ***REMOVED***
		if r, err := language.ParseRegion(region[:2]); err == nil ***REMOVED***
			tFull := t
			t, _ = t.SetTypeForKey("rg", "")
			// TODO: should we not consider "va" for the language tag?
			var exact1, exact2 bool
			tag.language, exact1 = FromTag(t)
			t.RegionID = r
			tag.locale, exact2 = FromTag(t)
			if !exact1 || !exact2 ***REMOVED***
				tag.full = tFull
			***REMOVED***
			return tag
		***REMOVED***
	***REMOVED***
	lang, ok := FromTag(t)
	tag.language = lang
	tag.locale = lang
	if !ok ***REMOVED***
		tag.full = t
	***REMOVED***
	return tag
***REMOVED***

// Tag returns an internal language Tag version of this tag.
func (t Tag) Tag() language.Tag ***REMOVED***
	if t.full != nil ***REMOVED***
		return t.full.(language.Tag)
	***REMOVED***
	tag := t.language.Tag()
	if t.language != t.locale ***REMOVED***
		loc := t.locale.Tag()
		tag, _ = tag.SetTypeForKey("rg", strings.ToLower(loc.RegionID.String())+"zzzz")
	***REMOVED***
	return tag
***REMOVED***

// IsCompact reports whether this tag is fully defined in terms of ID.
func (t *Tag) IsCompact() bool ***REMOVED***
	return t.full == nil
***REMOVED***

// MayHaveVariants reports whether a tag may have variants. If it returns false
// it is guaranteed the tag does not have variants.
func (t Tag) MayHaveVariants() bool ***REMOVED***
	return t.full != nil || int(t.language) >= len(coreTags)
***REMOVED***

// MayHaveExtensions reports whether a tag may have extensions. If it returns
// false it is guaranteed the tag does not have them.
func (t Tag) MayHaveExtensions() bool ***REMOVED***
	return t.full != nil ||
		int(t.language) >= len(coreTags) ||
		t.language != t.locale
***REMOVED***

// IsRoot returns true if t is equal to language "und".
func (t Tag) IsRoot() bool ***REMOVED***
	if t.full != nil ***REMOVED***
		return t.full.IsRoot()
	***REMOVED***
	return t.language == _und
***REMOVED***

// Parent returns the CLDR parent of t. In CLDR, missing fields in data for a
// specific language are substituted with fields from the parent language.
// The parent for a language may change for newer versions of CLDR.
func (t Tag) Parent() Tag ***REMOVED***
	if t.full != nil ***REMOVED***
		return Make(t.full.Parent())
	***REMOVED***
	if t.language != t.locale ***REMOVED***
		// Simulate stripping -u-rg-xxxxxx
		return Tag***REMOVED***language: t.language, locale: t.language***REMOVED***
	***REMOVED***
	// TODO: use parent lookup table once cycle from internal package is
	// removed. Probably by internalizing the table and declaring this fast
	// enough.
	// lang := compactID(internal.Parent(uint16(t.language)))
	lang, _ := FromTag(t.language.Tag().Parent())
	return Tag***REMOVED***language: lang, locale: lang***REMOVED***
***REMOVED***

// returns token t and the rest of the string.
func nextToken(s string) (t, tail string) ***REMOVED***
	p := strings.Index(s[1:], "-")
	if p == -1 ***REMOVED***
		return s[1:], ""
	***REMOVED***
	p++
	return s[1:p], s[p:]
***REMOVED***

// LanguageID returns an index, where 0 <= index < NumCompactTags, for tags
// for which data exists in the text repository.The index will change over time
// and should not be stored in persistent storage. If t does not match a compact
// index, exact will be false and the compact index will be returned for the
// first match after repeatedly taking the Parent of t.
func LanguageID(t Tag) (id ID, exact bool) ***REMOVED***
	return t.language, t.full == nil
***REMOVED***

// RegionalID returns the ID for the regional variant of this tag. This index is
// used to indicate region-specific overrides, such as default currency, default
// calendar and week data, default time cycle, and default measurement system
// and unit preferences.
//
// For instance, the tag en-GB-u-rg-uszzzz specifies British English with US
// settings for currency, number formatting, etc. The CompactIndex for this tag
// will be that for en-GB, while the RegionalID will be the one corresponding to
// en-US.
func RegionalID(t Tag) (id ID, exact bool) ***REMOVED***
	return t.locale, t.full == nil
***REMOVED***

// LanguageTag returns t stripped of regional variant indicators.
//
// At the moment this means it is stripped of a regional and variant subtag "rg"
// and "va" in the "u" extension.
func (t Tag) LanguageTag() Tag ***REMOVED***
	if t.full == nil ***REMOVED***
		return Tag***REMOVED***language: t.language, locale: t.language***REMOVED***
	***REMOVED***
	tt := t.Tag()
	tt.SetTypeForKey("rg", "")
	tt.SetTypeForKey("va", "")
	return Make(tt)
***REMOVED***

// RegionalTag returns the regional variant of the tag.
//
// At the moment this means that the region is set from the regional subtag
// "rg" in the "u" extension.
func (t Tag) RegionalTag() Tag ***REMOVED***
	rt := Tag***REMOVED***language: t.locale, locale: t.locale***REMOVED***
	if t.full == nil ***REMOVED***
		return rt
	***REMOVED***
	b := language.Builder***REMOVED******REMOVED***
	tag := t.Tag()
	// tag, _ = tag.SetTypeForKey("rg", "")
	b.SetTag(t.locale.Tag())
	if v := tag.Variants(); v != "" ***REMOVED***
		for _, v := range strings.Split(v, "-") ***REMOVED***
			b.AddVariant(v)
		***REMOVED***
	***REMOVED***
	for _, e := range tag.Extensions() ***REMOVED***
		b.AddExt(e)
	***REMOVED***
	return t
***REMOVED***

// FromTag reports closest matching ID for an internal language Tag.
func FromTag(t language.Tag) (id ID, exact bool) ***REMOVED***
	// TODO: perhaps give more frequent tags a lower index.
	// TODO: we could make the indexes stable. This will excluded some
	//       possibilities for optimization, so don't do this quite yet.
	exact = true

	b, s, r := t.Raw()
	if t.HasString() ***REMOVED***
		if t.IsPrivateUse() ***REMOVED***
			// We have no entries for user-defined tags.
			return 0, false
		***REMOVED***
		hasExtra := false
		if t.HasVariants() ***REMOVED***
			if t.HasExtensions() ***REMOVED***
				build := language.Builder***REMOVED******REMOVED***
				build.SetTag(language.Tag***REMOVED***LangID: b, ScriptID: s, RegionID: r***REMOVED***)
				build.AddVariant(t.Variants())
				exact = false
				t = build.Make()
			***REMOVED***
			hasExtra = true
		***REMOVED*** else if _, ok := t.Extension('u'); ok ***REMOVED***
			// TODO: va may mean something else. Consider not considering it.
			// Strip all but the 'va' entry.
			old := t
			variant := t.TypeForKey("va")
			t = language.Tag***REMOVED***LangID: b, ScriptID: s, RegionID: r***REMOVED***
			if variant != "" ***REMOVED***
				t, _ = t.SetTypeForKey("va", variant)
				hasExtra = true
			***REMOVED***
			exact = old == t
		***REMOVED*** else ***REMOVED***
			exact = false
		***REMOVED***
		if hasExtra ***REMOVED***
			// We have some variants.
			for i, s := range specialTags ***REMOVED***
				if s == t ***REMOVED***
					return ID(i + len(coreTags)), exact
				***REMOVED***
			***REMOVED***
			exact = false
		***REMOVED***
	***REMOVED***
	if x, ok := getCoreIndex(t); ok ***REMOVED***
		return x, exact
	***REMOVED***
	exact = false
	if r != 0 && s == 0 ***REMOVED***
		// Deal with cases where an extra script is inserted for the region.
		t, _ := t.Maximize()
		if x, ok := getCoreIndex(t); ok ***REMOVED***
			return x, exact
		***REMOVED***
	***REMOVED***
	for t = t.Parent(); t != root; t = t.Parent() ***REMOVED***
		// No variants specified: just compare core components.
		// The key has the form lllssrrr, where l, s, and r are nibbles for
		// respectively the langID, scriptID, and regionID.
		if x, ok := getCoreIndex(t); ok ***REMOVED***
			return x, exact
		***REMOVED***
	***REMOVED***
	return 0, exact
***REMOVED***

var root = language.Tag***REMOVED******REMOVED***
