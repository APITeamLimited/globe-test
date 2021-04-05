// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go -output tables.go

package language

// TODO: Remove above NOTE after:
// - verifying that tables are dropped correctly (most notably matcher tables).

import (
	"strings"

	"golang.org/x/text/internal/language"
	"golang.org/x/text/internal/language/compact"
)

// Tag represents a BCP 47 language tag. It is used to specify an instance of a
// specific language or locale. All language tag values are guaranteed to be
// well-formed.
type Tag compact.Tag

func makeTag(t language.Tag) (tag Tag) ***REMOVED***
	return Tag(compact.Make(t))
***REMOVED***

func (t *Tag) tag() language.Tag ***REMOVED***
	return (*compact.Tag)(t).Tag()
***REMOVED***

func (t *Tag) isCompact() bool ***REMOVED***
	return (*compact.Tag)(t).IsCompact()
***REMOVED***

// TODO: improve performance.
func (t *Tag) lang() language.Language ***REMOVED*** return t.tag().LangID ***REMOVED***
func (t *Tag) region() language.Region ***REMOVED*** return t.tag().RegionID ***REMOVED***
func (t *Tag) script() language.Script ***REMOVED*** return t.tag().ScriptID ***REMOVED***

// Make is a convenience wrapper for Parse that omits the error.
// In case of an error, a sensible default is returned.
func Make(s string) Tag ***REMOVED***
	return Default.Make(s)
***REMOVED***

// Make is a convenience wrapper for c.Parse that omits the error.
// In case of an error, a sensible default is returned.
func (c CanonType) Make(s string) Tag ***REMOVED***
	t, _ := c.Parse(s)
	return t
***REMOVED***

// Raw returns the raw base language, script and region, without making an
// attempt to infer their values.
func (t Tag) Raw() (b Base, s Script, r Region) ***REMOVED***
	tt := t.tag()
	return Base***REMOVED***tt.LangID***REMOVED***, Script***REMOVED***tt.ScriptID***REMOVED***, Region***REMOVED***tt.RegionID***REMOVED***
***REMOVED***

// IsRoot returns true if t is equal to language "und".
func (t Tag) IsRoot() bool ***REMOVED***
	return compact.Tag(t).IsRoot()
***REMOVED***

// CanonType can be used to enable or disable various types of canonicalization.
type CanonType int

const (
	// Replace deprecated base languages with their preferred replacements.
	DeprecatedBase CanonType = 1 << iota
	// Replace deprecated scripts with their preferred replacements.
	DeprecatedScript
	// Replace deprecated regions with their preferred replacements.
	DeprecatedRegion
	// Remove redundant scripts.
	SuppressScript
	// Normalize legacy encodings. This includes legacy languages defined in
	// CLDR as well as bibliographic codes defined in ISO-639.
	Legacy
	// Map the dominant language of a macro language group to the macro language
	// subtag. For example cmn -> zh.
	Macro
	// The CLDR flag should be used if full compatibility with CLDR is required.
	// There are a few cases where language.Tag may differ from CLDR. To follow all
	// of CLDR's suggestions, use All|CLDR.
	CLDR

	// Raw can be used to Compose or Parse without Canonicalization.
	Raw CanonType = 0

	// Replace all deprecated tags with their preferred replacements.
	Deprecated = DeprecatedBase | DeprecatedScript | DeprecatedRegion

	// All canonicalizations recommended by BCP 47.
	BCP47 = Deprecated | SuppressScript

	// All canonicalizations.
	All = BCP47 | Legacy | Macro

	// Default is the canonicalization used by Parse, Make and Compose. To
	// preserve as much information as possible, canonicalizations that remove
	// potentially valuable information are not included. The Matcher is
	// designed to recognize similar tags that would be the same if
	// they were canonicalized using All.
	Default = Deprecated | Legacy

	canonLang = DeprecatedBase | Legacy | Macro

	// TODO: LikelyScript, LikelyRegion: suppress similar to ICU.
)

// canonicalize returns the canonicalized equivalent of the tag and
// whether there was any change.
func canonicalize(c CanonType, t language.Tag) (language.Tag, bool) ***REMOVED***
	if c == Raw ***REMOVED***
		return t, false
	***REMOVED***
	changed := false
	if c&SuppressScript != 0 ***REMOVED***
		if t.LangID.SuppressScript() == t.ScriptID ***REMOVED***
			t.ScriptID = 0
			changed = true
		***REMOVED***
	***REMOVED***
	if c&canonLang != 0 ***REMOVED***
		for ***REMOVED***
			if l, aliasType := t.LangID.Canonicalize(); l != t.LangID ***REMOVED***
				switch aliasType ***REMOVED***
				case language.Legacy:
					if c&Legacy != 0 ***REMOVED***
						if t.LangID == _sh && t.ScriptID == 0 ***REMOVED***
							t.ScriptID = _Latn
						***REMOVED***
						t.LangID = l
						changed = true
					***REMOVED***
				case language.Macro:
					if c&Macro != 0 ***REMOVED***
						// We deviate here from CLDR. The mapping "nb" -> "no"
						// qualifies as a typical Macro language mapping.  However,
						// for legacy reasons, CLDR maps "no", the macro language
						// code for Norwegian, to the dominant variant "nb". This
						// change is currently under consideration for CLDR as well.
						// See https://unicode.org/cldr/trac/ticket/2698 and also
						// https://unicode.org/cldr/trac/ticket/1790 for some of the
						// practical implications. TODO: this check could be removed
						// if CLDR adopts this change.
						if c&CLDR == 0 || t.LangID != _nb ***REMOVED***
							changed = true
							t.LangID = l
						***REMOVED***
					***REMOVED***
				case language.Deprecated:
					if c&DeprecatedBase != 0 ***REMOVED***
						if t.LangID == _mo && t.RegionID == 0 ***REMOVED***
							t.RegionID = _MD
						***REMOVED***
						t.LangID = l
						changed = true
						// Other canonicalization types may still apply.
						continue
					***REMOVED***
				***REMOVED***
			***REMOVED*** else if c&Legacy != 0 && t.LangID == _no && c&CLDR != 0 ***REMOVED***
				t.LangID = _nb
				changed = true
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if c&DeprecatedScript != 0 ***REMOVED***
		if t.ScriptID == _Qaai ***REMOVED***
			changed = true
			t.ScriptID = _Zinh
		***REMOVED***
	***REMOVED***
	if c&DeprecatedRegion != 0 ***REMOVED***
		if r := t.RegionID.Canonicalize(); r != t.RegionID ***REMOVED***
			changed = true
			t.RegionID = r
		***REMOVED***
	***REMOVED***
	return t, changed
***REMOVED***

// Canonicalize returns the canonicalized equivalent of the tag.
func (c CanonType) Canonicalize(t Tag) (Tag, error) ***REMOVED***
	// First try fast path.
	if t.isCompact() ***REMOVED***
		if _, changed := canonicalize(c, compact.Tag(t).Tag()); !changed ***REMOVED***
			return t, nil
		***REMOVED***
	***REMOVED***
	// It is unlikely that one will canonicalize a tag after matching. So do
	// a slow but simple approach here.
	if tag, changed := canonicalize(c, t.tag()); changed ***REMOVED***
		tag.RemakeString()
		return makeTag(tag), nil
	***REMOVED***
	return t, nil

***REMOVED***

// Confidence indicates the level of certainty for a given return value.
// For example, Serbian may be written in Cyrillic or Latin script.
// The confidence level indicates whether a value was explicitly specified,
// whether it is typically the only possible value, or whether there is
// an ambiguity.
type Confidence int

const (
	No    Confidence = iota // full confidence that there was no match
	Low                     // most likely value picked out of a set of alternatives
	High                    // value is generally assumed to be the correct match
	Exact                   // exact match or explicitly specified value
)

var confName = []string***REMOVED***"No", "Low", "High", "Exact"***REMOVED***

func (c Confidence) String() string ***REMOVED***
	return confName[c]
***REMOVED***

// String returns the canonical string representation of the language tag.
func (t Tag) String() string ***REMOVED***
	return t.tag().String()
***REMOVED***

// MarshalText implements encoding.TextMarshaler.
func (t Tag) MarshalText() (text []byte, err error) ***REMOVED***
	return t.tag().MarshalText()
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Tag) UnmarshalText(text []byte) error ***REMOVED***
	var tag language.Tag
	err := tag.UnmarshalText(text)
	*t = makeTag(tag)
	return err
***REMOVED***

// Base returns the base language of the language tag. If the base language is
// unspecified, an attempt will be made to infer it from the context.
// It uses a variant of CLDR's Add Likely Subtags algorithm. This is subject to change.
func (t Tag) Base() (Base, Confidence) ***REMOVED***
	if b := t.lang(); b != 0 ***REMOVED***
		return Base***REMOVED***b***REMOVED***, Exact
	***REMOVED***
	tt := t.tag()
	c := High
	if tt.ScriptID == 0 && !tt.RegionID.IsCountry() ***REMOVED***
		c = Low
	***REMOVED***
	if tag, err := tt.Maximize(); err == nil && tag.LangID != 0 ***REMOVED***
		return Base***REMOVED***tag.LangID***REMOVED***, c
	***REMOVED***
	return Base***REMOVED***0***REMOVED***, No
***REMOVED***

// Script infers the script for the language tag. If it was not explicitly given, it will infer
// a most likely candidate.
// If more than one script is commonly used for a language, the most likely one
// is returned with a low confidence indication. For example, it returns (Cyrl, Low)
// for Serbian.
// If a script cannot be inferred (Zzzz, No) is returned. We do not use Zyyy (undetermined)
// as one would suspect from the IANA registry for BCP 47. In a Unicode context Zyyy marks
// common characters (like 1, 2, 3, '.', etc.) and is therefore more like multiple scripts.
// See https://www.unicode.org/reports/tr24/#Values for more details. Zzzz is also used for
// unknown value in CLDR.  (Zzzz, Exact) is returned if Zzzz was explicitly specified.
// Note that an inferred script is never guaranteed to be the correct one. Latin is
// almost exclusively used for Afrikaans, but Arabic has been used for some texts
// in the past.  Also, the script that is commonly used may change over time.
// It uses a variant of CLDR's Add Likely Subtags algorithm. This is subject to change.
func (t Tag) Script() (Script, Confidence) ***REMOVED***
	if scr := t.script(); scr != 0 ***REMOVED***
		return Script***REMOVED***scr***REMOVED***, Exact
	***REMOVED***
	tt := t.tag()
	sc, c := language.Script(_Zzzz), No
	if scr := tt.LangID.SuppressScript(); scr != 0 ***REMOVED***
		// Note: it is not always the case that a language with a suppress
		// script value is only written in one script (e.g. kk, ms, pa).
		if tt.RegionID == 0 ***REMOVED***
			return Script***REMOVED***scr***REMOVED***, High
		***REMOVED***
		sc, c = scr, High
	***REMOVED***
	if tag, err := tt.Maximize(); err == nil ***REMOVED***
		if tag.ScriptID != sc ***REMOVED***
			sc, c = tag.ScriptID, Low
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		tt, _ = canonicalize(Deprecated|Macro, tt)
		if tag, err := tt.Maximize(); err == nil && tag.ScriptID != sc ***REMOVED***
			sc, c = tag.ScriptID, Low
		***REMOVED***
	***REMOVED***
	return Script***REMOVED***sc***REMOVED***, c
***REMOVED***

// Region returns the region for the language tag. If it was not explicitly given, it will
// infer a most likely candidate from the context.
// It uses a variant of CLDR's Add Likely Subtags algorithm. This is subject to change.
func (t Tag) Region() (Region, Confidence) ***REMOVED***
	if r := t.region(); r != 0 ***REMOVED***
		return Region***REMOVED***r***REMOVED***, Exact
	***REMOVED***
	tt := t.tag()
	if tt, err := tt.Maximize(); err == nil ***REMOVED***
		return Region***REMOVED***tt.RegionID***REMOVED***, Low // TODO: differentiate between high and low.
	***REMOVED***
	tt, _ = canonicalize(Deprecated|Macro, tt)
	if tag, err := tt.Maximize(); err == nil ***REMOVED***
		return Region***REMOVED***tag.RegionID***REMOVED***, Low
	***REMOVED***
	return Region***REMOVED***_ZZ***REMOVED***, No // TODO: return world instead of undetermined?
***REMOVED***

// Variants returns the variants specified explicitly for this language tag.
// or nil if no variant was specified.
func (t Tag) Variants() []Variant ***REMOVED***
	if !compact.Tag(t).MayHaveVariants() ***REMOVED***
		return nil
	***REMOVED***
	v := []Variant***REMOVED******REMOVED***
	x, str := "", t.tag().Variants()
	for str != "" ***REMOVED***
		x, str = nextToken(str)
		v = append(v, Variant***REMOVED***x***REMOVED***)
	***REMOVED***
	return v
***REMOVED***

// Parent returns the CLDR parent of t. In CLDR, missing fields in data for a
// specific language are substituted with fields from the parent language.
// The parent for a language may change for newer versions of CLDR.
//
// Parent returns a tag for a less specific language that is mutually
// intelligible or Und if there is no such language. This may not be the same as
// simply stripping the last BCP 47 subtag. For instance, the parent of "zh-TW"
// is "zh-Hant", and the parent of "zh-Hant" is "und".
func (t Tag) Parent() Tag ***REMOVED***
	return Tag(compact.Tag(t).Parent())
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

// Extension is a single BCP 47 extension.
type Extension struct ***REMOVED***
	s string
***REMOVED***

// String returns the string representation of the extension, including the
// type tag.
func (e Extension) String() string ***REMOVED***
	return e.s
***REMOVED***

// ParseExtension parses s as an extension and returns it on success.
func ParseExtension(s string) (e Extension, err error) ***REMOVED***
	ext, err := language.ParseExtension(s)
	return Extension***REMOVED***ext***REMOVED***, err
***REMOVED***

// Type returns the one-byte extension type of e. It returns 0 for the zero
// exception.
func (e Extension) Type() byte ***REMOVED***
	if e.s == "" ***REMOVED***
		return 0
	***REMOVED***
	return e.s[0]
***REMOVED***

// Tokens returns the list of tokens of e.
func (e Extension) Tokens() []string ***REMOVED***
	return strings.Split(e.s, "-")
***REMOVED***

// Extension returns the extension of type x for tag t. It will return
// false for ok if t does not have the requested extension. The returned
// extension will be invalid in this case.
func (t Tag) Extension(x byte) (ext Extension, ok bool) ***REMOVED***
	if !compact.Tag(t).MayHaveExtensions() ***REMOVED***
		return Extension***REMOVED******REMOVED***, false
	***REMOVED***
	e, ok := t.tag().Extension(x)
	return Extension***REMOVED***e***REMOVED***, ok
***REMOVED***

// Extensions returns all extensions of t.
func (t Tag) Extensions() []Extension ***REMOVED***
	if !compact.Tag(t).MayHaveExtensions() ***REMOVED***
		return nil
	***REMOVED***
	e := []Extension***REMOVED******REMOVED***
	for _, ext := range t.tag().Extensions() ***REMOVED***
		e = append(e, Extension***REMOVED***ext***REMOVED***)
	***REMOVED***
	return e
***REMOVED***

// TypeForKey returns the type associated with the given key, where key and type
// are of the allowed values defined for the Unicode locale extension ('u') in
// https://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
// TypeForKey will traverse the inheritance chain to get the correct value.
//
// If there are multiple types associated with a key, only the first will be
// returned. If there is no type associated with a key, it returns the empty
// string.
func (t Tag) TypeForKey(key string) string ***REMOVED***
	if !compact.Tag(t).MayHaveExtensions() ***REMOVED***
		if key != "rg" && key != "va" ***REMOVED***
			return ""
		***REMOVED***
	***REMOVED***
	return t.tag().TypeForKey(key)
***REMOVED***

// SetTypeForKey returns a new Tag with the key set to type, where key and type
// are of the allowed values defined for the Unicode locale extension ('u') in
// https://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
// An empty value removes an existing pair with the same key.
func (t Tag) SetTypeForKey(key, value string) (Tag, error) ***REMOVED***
	tt, err := t.tag().SetTypeForKey(key, value)
	return makeTag(tt), err
***REMOVED***

// NumCompactTags is the number of compact tags. The maximum tag is
// NumCompactTags-1.
const NumCompactTags = compact.NumCompactTags

// CompactIndex returns an index, where 0 <= index < NumCompactTags, for tags
// for which data exists in the text repository.The index will change over time
// and should not be stored in persistent storage. If t does not match a compact
// index, exact will be false and the compact index will be returned for the
// first match after repeatedly taking the Parent of t.
func CompactIndex(t Tag) (index int, exact bool) ***REMOVED***
	id, exact := compact.LanguageID(compact.Tag(t))
	return int(id), exact
***REMOVED***

var root = language.Tag***REMOVED******REMOVED***

// Base is an ISO 639 language code, used for encoding the base language
// of a language tag.
type Base struct ***REMOVED***
	langID language.Language
***REMOVED***

// ParseBase parses a 2- or 3-letter ISO 639 code.
// It returns a ValueError if s is a well-formed but unknown language identifier
// or another error if another error occurred.
func ParseBase(s string) (Base, error) ***REMOVED***
	l, err := language.ParseBase(s)
	return Base***REMOVED***l***REMOVED***, err
***REMOVED***

// String returns the BCP 47 representation of the base language.
func (b Base) String() string ***REMOVED***
	return b.langID.String()
***REMOVED***

// ISO3 returns the ISO 639-3 language code.
func (b Base) ISO3() string ***REMOVED***
	return b.langID.ISO3()
***REMOVED***

// IsPrivateUse reports whether this language code is reserved for private use.
func (b Base) IsPrivateUse() bool ***REMOVED***
	return b.langID.IsPrivateUse()
***REMOVED***

// Script is a 4-letter ISO 15924 code for representing scripts.
// It is idiomatically represented in title case.
type Script struct ***REMOVED***
	scriptID language.Script
***REMOVED***

// ParseScript parses a 4-letter ISO 15924 code.
// It returns a ValueError if s is a well-formed but unknown script identifier
// or another error if another error occurred.
func ParseScript(s string) (Script, error) ***REMOVED***
	sc, err := language.ParseScript(s)
	return Script***REMOVED***sc***REMOVED***, err
***REMOVED***

// String returns the script code in title case.
// It returns "Zzzz" for an unspecified script.
func (s Script) String() string ***REMOVED***
	return s.scriptID.String()
***REMOVED***

// IsPrivateUse reports whether this script code is reserved for private use.
func (s Script) IsPrivateUse() bool ***REMOVED***
	return s.scriptID.IsPrivateUse()
***REMOVED***

// Region is an ISO 3166-1 or UN M.49 code for representing countries and regions.
type Region struct ***REMOVED***
	regionID language.Region
***REMOVED***

// EncodeM49 returns the Region for the given UN M.49 code.
// It returns an error if r is not a valid code.
func EncodeM49(r int) (Region, error) ***REMOVED***
	rid, err := language.EncodeM49(r)
	return Region***REMOVED***rid***REMOVED***, err
***REMOVED***

// ParseRegion parses a 2- or 3-letter ISO 3166-1 or a UN M.49 code.
// It returns a ValueError if s is a well-formed but unknown region identifier
// or another error if another error occurred.
func ParseRegion(s string) (Region, error) ***REMOVED***
	r, err := language.ParseRegion(s)
	return Region***REMOVED***r***REMOVED***, err
***REMOVED***

// String returns the BCP 47 representation for the region.
// It returns "ZZ" for an unspecified region.
func (r Region) String() string ***REMOVED***
	return r.regionID.String()
***REMOVED***

// ISO3 returns the 3-letter ISO code of r.
// Note that not all regions have a 3-letter ISO code.
// In such cases this method returns "ZZZ".
func (r Region) ISO3() string ***REMOVED***
	return r.regionID.ISO3()
***REMOVED***

// M49 returns the UN M.49 encoding of r, or 0 if this encoding
// is not defined for r.
func (r Region) M49() int ***REMOVED***
	return r.regionID.M49()
***REMOVED***

// IsPrivateUse reports whether r has the ISO 3166 User-assigned status. This
// may include private-use tags that are assigned by CLDR and used in this
// implementation. So IsPrivateUse and IsCountry can be simultaneously true.
func (r Region) IsPrivateUse() bool ***REMOVED***
	return r.regionID.IsPrivateUse()
***REMOVED***

// IsCountry returns whether this region is a country or autonomous area. This
// includes non-standard definitions from CLDR.
func (r Region) IsCountry() bool ***REMOVED***
	return r.regionID.IsCountry()
***REMOVED***

// IsGroup returns whether this region defines a collection of regions. This
// includes non-standard definitions from CLDR.
func (r Region) IsGroup() bool ***REMOVED***
	return r.regionID.IsGroup()
***REMOVED***

// Contains returns whether Region c is contained by Region r. It returns true
// if c == r.
func (r Region) Contains(c Region) bool ***REMOVED***
	return r.regionID.Contains(c.regionID)
***REMOVED***

// TLD returns the country code top-level domain (ccTLD). UK is returned for GB.
// In all other cases it returns either the region itself or an error.
//
// This method may return an error for a region for which there exists a
// canonical form with a ccTLD. To get that ccTLD canonicalize r first. The
// region will already be canonicalized it was obtained from a Tag that was
// obtained using any of the default methods.
func (r Region) TLD() (Region, error) ***REMOVED***
	tld, err := r.regionID.TLD()
	return Region***REMOVED***tld***REMOVED***, err
***REMOVED***

// Canonicalize returns the region or a possible replacement if the region is
// deprecated. It will not return a replacement for deprecated regions that
// are split into multiple regions.
func (r Region) Canonicalize() Region ***REMOVED***
	return Region***REMOVED***r.regionID.Canonicalize()***REMOVED***
***REMOVED***

// Variant represents a registered variant of a language as defined by BCP 47.
type Variant struct ***REMOVED***
	variant string
***REMOVED***

// ParseVariant parses and returns a Variant. An error is returned if s is not
// a valid variant.
func ParseVariant(s string) (Variant, error) ***REMOVED***
	v, err := language.ParseVariant(s)
	return Variant***REMOVED***v.String()***REMOVED***, err
***REMOVED***

// String returns the string representation of the variant.
func (v Variant) String() string ***REMOVED***
	return v.variant
***REMOVED***
