// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_common.go -output tables.go
//go:generate go run gen_index.go

package language

// TODO: Remove above NOTE after:
// - verifying that tables are dropped correctly (most notably matcher tables).

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// maxCoreSize is the maximum size of a BCP 47 tag without variants and
	// extensions. Equals max lang (3) + script (4) + max reg (3) + 2 dashes.
	maxCoreSize = 12

	// max99thPercentileSize is a somewhat arbitrary buffer size that presumably
	// is large enough to hold at least 99% of the BCP 47 tags.
	max99thPercentileSize = 32

	// maxSimpleUExtensionSize is the maximum size of a -u extension with one
	// key-type pair. Equals len("-u-") + key (2) + dash + max value (8).
	maxSimpleUExtensionSize = 14
)

// Tag represents a BCP 47 language tag. It is used to specify an instance of a
// specific language or locale. All language tag values are guaranteed to be
// well-formed.
type Tag struct ***REMOVED***
	lang   langID
	region regionID
	// TODO: we will soon run out of positions for script. Idea: instead of
	// storing lang, region, and script codes, store only the compact index and
	// have a lookup table from this code to its expansion. This greatly speeds
	// up table lookup, speed up common variant cases.
	// This will also immediately free up 3 extra bytes. Also, the pVariant
	// field can now be moved to the lookup table, as the compact index uniquely
	// determines the offset of a possible variant.
	script   scriptID
	pVariant byte   // offset in str, includes preceding '-'
	pExt     uint16 // offset of first extension, includes preceding '-'

	// str is the string representation of the Tag. It will only be used if the
	// tag has variants or extensions.
	str string
***REMOVED***

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
	return Base***REMOVED***t.lang***REMOVED***, Script***REMOVED***t.script***REMOVED***, Region***REMOVED***t.region***REMOVED***
***REMOVED***

// equalTags compares language, script and region subtags only.
func (t Tag) equalTags(a Tag) bool ***REMOVED***
	return t.lang == a.lang && t.script == a.script && t.region == a.region
***REMOVED***

// IsRoot returns true if t is equal to language "und".
func (t Tag) IsRoot() bool ***REMOVED***
	if int(t.pVariant) < len(t.str) ***REMOVED***
		return false
	***REMOVED***
	return t.equalTags(und)
***REMOVED***

// private reports whether the Tag consists solely of a private use tag.
func (t Tag) private() bool ***REMOVED***
	return t.str != "" && t.pVariant == 0
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
func (t Tag) canonicalize(c CanonType) (Tag, bool) ***REMOVED***
	if c == Raw ***REMOVED***
		return t, false
	***REMOVED***
	changed := false
	if c&SuppressScript != 0 ***REMOVED***
		if t.lang < langNoIndexOffset && uint8(t.script) == suppressScript[t.lang] ***REMOVED***
			t.script = 0
			changed = true
		***REMOVED***
	***REMOVED***
	if c&canonLang != 0 ***REMOVED***
		for ***REMOVED***
			if l, aliasType := normLang(t.lang); l != t.lang ***REMOVED***
				switch aliasType ***REMOVED***
				case langLegacy:
					if c&Legacy != 0 ***REMOVED***
						if t.lang == _sh && t.script == 0 ***REMOVED***
							t.script = _Latn
						***REMOVED***
						t.lang = l
						changed = true
					***REMOVED***
				case langMacro:
					if c&Macro != 0 ***REMOVED***
						// We deviate here from CLDR. The mapping "nb" -> "no"
						// qualifies as a typical Macro language mapping.  However,
						// for legacy reasons, CLDR maps "no", the macro language
						// code for Norwegian, to the dominant variant "nb". This
						// change is currently under consideration for CLDR as well.
						// See http://unicode.org/cldr/trac/ticket/2698 and also
						// http://unicode.org/cldr/trac/ticket/1790 for some of the
						// practical implications. TODO: this check could be removed
						// if CLDR adopts this change.
						if c&CLDR == 0 || t.lang != _nb ***REMOVED***
							changed = true
							t.lang = l
						***REMOVED***
					***REMOVED***
				case langDeprecated:
					if c&DeprecatedBase != 0 ***REMOVED***
						if t.lang == _mo && t.region == 0 ***REMOVED***
							t.region = _MD
						***REMOVED***
						t.lang = l
						changed = true
						// Other canonicalization types may still apply.
						continue
					***REMOVED***
				***REMOVED***
			***REMOVED*** else if c&Legacy != 0 && t.lang == _no && c&CLDR != 0 ***REMOVED***
				t.lang = _nb
				changed = true
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	if c&DeprecatedScript != 0 ***REMOVED***
		if t.script == _Qaai ***REMOVED***
			changed = true
			t.script = _Zinh
		***REMOVED***
	***REMOVED***
	if c&DeprecatedRegion != 0 ***REMOVED***
		if r := normRegion(t.region); r != 0 ***REMOVED***
			changed = true
			t.region = r
		***REMOVED***
	***REMOVED***
	return t, changed
***REMOVED***

// Canonicalize returns the canonicalized equivalent of the tag.
func (c CanonType) Canonicalize(t Tag) (Tag, error) ***REMOVED***
	t, changed := t.canonicalize(c)
	if changed ***REMOVED***
		t.remakeString()
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

// remakeString is used to update t.str in case lang, script or region changed.
// It is assumed that pExt and pVariant still point to the start of the
// respective parts.
func (t *Tag) remakeString() ***REMOVED***
	if t.str == "" ***REMOVED***
		return
	***REMOVED***
	extra := t.str[t.pVariant:]
	if t.pVariant > 0 ***REMOVED***
		extra = extra[1:]
	***REMOVED***
	if t.equalTags(und) && strings.HasPrefix(extra, "x-") ***REMOVED***
		t.str = extra
		t.pVariant = 0
		t.pExt = 0
		return
	***REMOVED***
	var buf [max99thPercentileSize]byte // avoid extra memory allocation in most cases.
	b := buf[:t.genCoreBytes(buf[:])]
	if extra != "" ***REMOVED***
		diff := len(b) - int(t.pVariant)
		b = append(b, '-')
		b = append(b, extra...)
		t.pVariant = uint8(int(t.pVariant) + diff)
		t.pExt = uint16(int(t.pExt) + diff)
	***REMOVED*** else ***REMOVED***
		t.pVariant = uint8(len(b))
		t.pExt = uint16(len(b))
	***REMOVED***
	t.str = string(b)
***REMOVED***

// genCoreBytes writes a string for the base languages, script and region tags
// to the given buffer and returns the number of bytes written. It will never
// write more than maxCoreSize bytes.
func (t *Tag) genCoreBytes(buf []byte) int ***REMOVED***
	n := t.lang.stringToBuf(buf[:])
	if t.script != 0 ***REMOVED***
		n += copy(buf[n:], "-")
		n += copy(buf[n:], t.script.String())
	***REMOVED***
	if t.region != 0 ***REMOVED***
		n += copy(buf[n:], "-")
		n += copy(buf[n:], t.region.String())
	***REMOVED***
	return n
***REMOVED***

// String returns the canonical string representation of the language tag.
func (t Tag) String() string ***REMOVED***
	if t.str != "" ***REMOVED***
		return t.str
	***REMOVED***
	if t.script == 0 && t.region == 0 ***REMOVED***
		return t.lang.String()
	***REMOVED***
	buf := [maxCoreSize]byte***REMOVED******REMOVED***
	return string(buf[:t.genCoreBytes(buf[:])])
***REMOVED***

// MarshalText implements encoding.TextMarshaler.
func (t Tag) MarshalText() (text []byte, err error) ***REMOVED***
	if t.str != "" ***REMOVED***
		text = append(text, t.str...)
	***REMOVED*** else if t.script == 0 && t.region == 0 ***REMOVED***
		text = append(text, t.lang.String()...)
	***REMOVED*** else ***REMOVED***
		buf := [maxCoreSize]byte***REMOVED******REMOVED***
		text = buf[:t.genCoreBytes(buf[:])]
	***REMOVED***
	return text, nil
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Tag) UnmarshalText(text []byte) error ***REMOVED***
	tag, err := Raw.Parse(string(text))
	*t = tag
	return err
***REMOVED***

// Base returns the base language of the language tag. If the base language is
// unspecified, an attempt will be made to infer it from the context.
// It uses a variant of CLDR's Add Likely Subtags algorithm. This is subject to change.
func (t Tag) Base() (Base, Confidence) ***REMOVED***
	if t.lang != 0 ***REMOVED***
		return Base***REMOVED***t.lang***REMOVED***, Exact
	***REMOVED***
	c := High
	if t.script == 0 && !(Region***REMOVED***t.region***REMOVED***).IsCountry() ***REMOVED***
		c = Low
	***REMOVED***
	if tag, err := addTags(t); err == nil && tag.lang != 0 ***REMOVED***
		return Base***REMOVED***tag.lang***REMOVED***, c
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
// See http://www.unicode.org/reports/tr24/#Values for more details. Zzzz is also used for
// unknown value in CLDR.  (Zzzz, Exact) is returned if Zzzz was explicitly specified.
// Note that an inferred script is never guaranteed to be the correct one. Latin is
// almost exclusively used for Afrikaans, but Arabic has been used for some texts
// in the past.  Also, the script that is commonly used may change over time.
// It uses a variant of CLDR's Add Likely Subtags algorithm. This is subject to change.
func (t Tag) Script() (Script, Confidence) ***REMOVED***
	if t.script != 0 ***REMOVED***
		return Script***REMOVED***t.script***REMOVED***, Exact
	***REMOVED***
	sc, c := scriptID(_Zzzz), No
	if t.lang < langNoIndexOffset ***REMOVED***
		if scr := scriptID(suppressScript[t.lang]); scr != 0 ***REMOVED***
			// Note: it is not always the case that a language with a suppress
			// script value is only written in one script (e.g. kk, ms, pa).
			if t.region == 0 ***REMOVED***
				return Script***REMOVED***scriptID(scr)***REMOVED***, High
			***REMOVED***
			sc, c = scr, High
		***REMOVED***
	***REMOVED***
	if tag, err := addTags(t); err == nil ***REMOVED***
		if tag.script != sc ***REMOVED***
			sc, c = tag.script, Low
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		t, _ = (Deprecated | Macro).Canonicalize(t)
		if tag, err := addTags(t); err == nil && tag.script != sc ***REMOVED***
			sc, c = tag.script, Low
		***REMOVED***
	***REMOVED***
	return Script***REMOVED***sc***REMOVED***, c
***REMOVED***

// Region returns the region for the language tag. If it was not explicitly given, it will
// infer a most likely candidate from the context.
// It uses a variant of CLDR's Add Likely Subtags algorithm. This is subject to change.
func (t Tag) Region() (Region, Confidence) ***REMOVED***
	if t.region != 0 ***REMOVED***
		return Region***REMOVED***t.region***REMOVED***, Exact
	***REMOVED***
	if t, err := addTags(t); err == nil ***REMOVED***
		return Region***REMOVED***t.region***REMOVED***, Low // TODO: differentiate between high and low.
	***REMOVED***
	t, _ = (Deprecated | Macro).Canonicalize(t)
	if tag, err := addTags(t); err == nil ***REMOVED***
		return Region***REMOVED***tag.region***REMOVED***, Low
	***REMOVED***
	return Region***REMOVED***_ZZ***REMOVED***, No // TODO: return world instead of undetermined?
***REMOVED***

// Variant returns the variants specified explicitly for this language tag.
// or nil if no variant was specified.
func (t Tag) Variants() []Variant ***REMOVED***
	v := []Variant***REMOVED******REMOVED***
	if int(t.pVariant) < int(t.pExt) ***REMOVED***
		for x, str := "", t.str[t.pVariant:t.pExt]; str != ""; ***REMOVED***
			x, str = nextToken(str)
			v = append(v, Variant***REMOVED***x***REMOVED***)
		***REMOVED***
	***REMOVED***
	return v
***REMOVED***

// Parent returns the CLDR parent of t. In CLDR, missing fields in data for a
// specific language are substituted with fields from the parent language.
// The parent for a language may change for newer versions of CLDR.
func (t Tag) Parent() Tag ***REMOVED***
	if t.str != "" ***REMOVED***
		// Strip the variants and extensions.
		t, _ = Raw.Compose(t.Raw())
		if t.region == 0 && t.script != 0 && t.lang != 0 ***REMOVED***
			base, _ := addTags(Tag***REMOVED***lang: t.lang***REMOVED***)
			if base.script == t.script ***REMOVED***
				return Tag***REMOVED***lang: t.lang***REMOVED***
			***REMOVED***
		***REMOVED***
		return t
	***REMOVED***
	if t.lang != 0 ***REMOVED***
		if t.region != 0 ***REMOVED***
			maxScript := t.script
			if maxScript == 0 ***REMOVED***
				max, _ := addTags(t)
				maxScript = max.script
			***REMOVED***

			for i := range parents ***REMOVED***
				if langID(parents[i].lang) == t.lang && scriptID(parents[i].maxScript) == maxScript ***REMOVED***
					for _, r := range parents[i].fromRegion ***REMOVED***
						if regionID(r) == t.region ***REMOVED***
							return Tag***REMOVED***
								lang:   t.lang,
								script: scriptID(parents[i].script),
								region: regionID(parents[i].toRegion),
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Strip the script if it is the default one.
			base, _ := addTags(Tag***REMOVED***lang: t.lang***REMOVED***)
			if base.script != maxScript ***REMOVED***
				return Tag***REMOVED***lang: t.lang, script: maxScript***REMOVED***
			***REMOVED***
			return Tag***REMOVED***lang: t.lang***REMOVED***
		***REMOVED*** else if t.script != 0 ***REMOVED***
			// The parent for an base-script pair with a non-default script is
			// "und" instead of the base language.
			base, _ := addTags(Tag***REMOVED***lang: t.lang***REMOVED***)
			if base.script != t.script ***REMOVED***
				return und
			***REMOVED***
			return Tag***REMOVED***lang: t.lang***REMOVED***
		***REMOVED***
	***REMOVED***
	return und
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
	scan := makeScannerString(s)
	var end int
	if n := len(scan.token); n != 1 ***REMOVED***
		return Extension***REMOVED******REMOVED***, errSyntax
	***REMOVED***
	scan.toLower(0, len(scan.b))
	end = parseExtension(&scan)
	if end != len(s) ***REMOVED***
		return Extension***REMOVED******REMOVED***, errSyntax
	***REMOVED***
	return Extension***REMOVED***string(scan.b)***REMOVED***, nil
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
	for i := int(t.pExt); i < len(t.str)-1; ***REMOVED***
		var ext string
		i, ext = getExtension(t.str, i)
		if ext[0] == x ***REMOVED***
			return Extension***REMOVED***ext***REMOVED***, true
		***REMOVED***
	***REMOVED***
	return Extension***REMOVED******REMOVED***, false
***REMOVED***

// Extensions returns all extensions of t.
func (t Tag) Extensions() []Extension ***REMOVED***
	e := []Extension***REMOVED******REMOVED***
	for i := int(t.pExt); i < len(t.str)-1; ***REMOVED***
		var ext string
		i, ext = getExtension(t.str, i)
		e = append(e, Extension***REMOVED***ext***REMOVED***)
	***REMOVED***
	return e
***REMOVED***

// TypeForKey returns the type associated with the given key, where key and type
// are of the allowed values defined for the Unicode locale extension ('u') in
// http://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
// TypeForKey will traverse the inheritance chain to get the correct value.
func (t Tag) TypeForKey(key string) string ***REMOVED***
	if start, end, _ := t.findTypeForKey(key); end != start ***REMOVED***
		return t.str[start:end]
	***REMOVED***
	return ""
***REMOVED***

var (
	errPrivateUse       = errors.New("cannot set a key on a private use tag")
	errInvalidArguments = errors.New("invalid key or type")
)

// SetTypeForKey returns a new Tag with the key set to type, where key and type
// are of the allowed values defined for the Unicode locale extension ('u') in
// http://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
// An empty value removes an existing pair with the same key.
func (t Tag) SetTypeForKey(key, value string) (Tag, error) ***REMOVED***
	if t.private() ***REMOVED***
		return t, errPrivateUse
	***REMOVED***
	if len(key) != 2 ***REMOVED***
		return t, errInvalidArguments
	***REMOVED***

	// Remove the setting if value is "".
	if value == "" ***REMOVED***
		start, end, _ := t.findTypeForKey(key)
		if start != end ***REMOVED***
			// Remove key tag and leading '-'.
			start -= 4

			// Remove a possible empty extension.
			if (end == len(t.str) || t.str[end+2] == '-') && t.str[start-2] == '-' ***REMOVED***
				start -= 2
			***REMOVED***
			if start == int(t.pVariant) && end == len(t.str) ***REMOVED***
				t.str = ""
				t.pVariant, t.pExt = 0, 0
			***REMOVED*** else ***REMOVED***
				t.str = fmt.Sprintf("%s%s", t.str[:start], t.str[end:])
			***REMOVED***
		***REMOVED***
		return t, nil
	***REMOVED***

	if len(value) < 3 || len(value) > 8 ***REMOVED***
		return t, errInvalidArguments
	***REMOVED***

	var (
		buf    [maxCoreSize + maxSimpleUExtensionSize]byte
		uStart int // start of the -u extension.
	)

	// Generate the tag string if needed.
	if t.str == "" ***REMOVED***
		uStart = t.genCoreBytes(buf[:])
		buf[uStart] = '-'
		uStart++
	***REMOVED***

	// Create new key-type pair and parse it to verify.
	b := buf[uStart:]
	copy(b, "u-")
	copy(b[2:], key)
	b[4] = '-'
	b = b[:5+copy(b[5:], value)]
	scan := makeScanner(b)
	if parseExtensions(&scan); scan.err != nil ***REMOVED***
		return t, scan.err
	***REMOVED***

	// Assemble the replacement string.
	if t.str == "" ***REMOVED***
		t.pVariant, t.pExt = byte(uStart-1), uint16(uStart-1)
		t.str = string(buf[:uStart+len(b)])
	***REMOVED*** else ***REMOVED***
		s := t.str
		start, end, hasExt := t.findTypeForKey(key)
		if start == end ***REMOVED***
			if hasExt ***REMOVED***
				b = b[2:]
			***REMOVED***
			t.str = fmt.Sprintf("%s-%s%s", s[:start], b, s[end:])
		***REMOVED*** else ***REMOVED***
			t.str = fmt.Sprintf("%s%s%s", s[:start], value, s[end:])
		***REMOVED***
	***REMOVED***
	return t, nil
***REMOVED***

// findKeyAndType returns the start and end position for the type corresponding
// to key or the point at which to insert the key-value pair if the type
// wasn't found. The hasExt return value reports whether an -u extension was present.
// Note: the extensions are typically very small and are likely to contain
// only one key-type pair.
func (t Tag) findTypeForKey(key string) (start, end int, hasExt bool) ***REMOVED***
	p := int(t.pExt)
	if len(key) != 2 || p == len(t.str) || p == 0 ***REMOVED***
		return p, p, false
	***REMOVED***
	s := t.str

	// Find the correct extension.
	for p++; s[p] != 'u'; p++ ***REMOVED***
		if s[p] > 'u' ***REMOVED***
			p--
			return p, p, false
		***REMOVED***
		if p = nextExtension(s, p); p == len(s) ***REMOVED***
			return len(s), len(s), false
		***REMOVED***
	***REMOVED***
	// Proceed to the hyphen following the extension name.
	p++

	// curKey is the key currently being processed.
	curKey := ""

	// Iterate over keys until we get the end of a section.
	for ***REMOVED***
		// p points to the hyphen preceding the current token.
		if p3 := p + 3; s[p3] == '-' ***REMOVED***
			// Found a key.
			// Check whether we just processed the key that was requested.
			if curKey == key ***REMOVED***
				return start, p, true
			***REMOVED***
			// Set to the next key and continue scanning type tokens.
			curKey = s[p+1 : p3]
			if curKey > key ***REMOVED***
				return p, p, true
			***REMOVED***
			// Start of the type token sequence.
			start = p + 4
			// A type is at least 3 characters long.
			p += 7 // 4 + 3
		***REMOVED*** else ***REMOVED***
			// Attribute or type, which is at least 3 characters long.
			p += 4
		***REMOVED***
		// p points past the third character of a type or attribute.
		max := p + 5 // maximum length of token plus hyphen.
		if len(s) < max ***REMOVED***
			max = len(s)
		***REMOVED***
		for ; p < max && s[p] != '-'; p++ ***REMOVED***
		***REMOVED***
		// Bail if we have exhausted all tokens or if the next token starts
		// a new extension.
		if p == len(s) || s[p+2] == '-' ***REMOVED***
			if curKey == key ***REMOVED***
				return start, p, true
			***REMOVED***
			return p, p, true
		***REMOVED***
	***REMOVED***
***REMOVED***

// CompactIndex returns an index, where 0 <= index < NumCompactTags, for tags
// for which data exists in the text repository. The index will change over time
// and should not be stored in persistent storage. Extensions, except for the
// 'va' type of the 'u' extension, are ignored. It will return 0, false if no
// compact tag exists, where 0 is the index for the root language (Und).
func CompactIndex(t Tag) (index int, ok bool) ***REMOVED***
	// TODO: perhaps give more frequent tags a lower index.
	// TODO: we could make the indexes stable. This will excluded some
	//       possibilities for optimization, so don't do this quite yet.
	b, s, r := t.Raw()
	if len(t.str) > 0 ***REMOVED***
		if strings.HasPrefix(t.str, "x-") ***REMOVED***
			// We have no entries for user-defined tags.
			return 0, false
		***REMOVED***
		if uint16(t.pVariant) != t.pExt ***REMOVED***
			// There are no tags with variants and an u-va type.
			if t.TypeForKey("va") != "" ***REMOVED***
				return 0, false
			***REMOVED***
			t, _ = Raw.Compose(b, s, r, t.Variants())
		***REMOVED*** else if _, ok := t.Extension('u'); ok ***REMOVED***
			// Strip all but the 'va' entry.
			variant := t.TypeForKey("va")
			t, _ = Raw.Compose(b, s, r)
			t, _ = t.SetTypeForKey("va", variant)
		***REMOVED***
		if len(t.str) > 0 ***REMOVED***
			// We have some variants.
			for i, s := range specialTags ***REMOVED***
				if s == t ***REMOVED***
					return i + 1, true
				***REMOVED***
			***REMOVED***
			return 0, false
		***REMOVED***
	***REMOVED***
	// No variants specified: just compare core components.
	// The key has the form lllssrrr, where l, s, and r are nibbles for
	// respectively the langID, scriptID, and regionID.
	key := uint32(b.langID) << (8 + 12)
	key |= uint32(s.scriptID) << 12
	key |= uint32(r.regionID)
	x, ok := coreTags[key]
	return int(x), ok
***REMOVED***

// Base is an ISO 639 language code, used for encoding the base language
// of a language tag.
type Base struct ***REMOVED***
	langID
***REMOVED***

// ParseBase parses a 2- or 3-letter ISO 639 code.
// It returns a ValueError if s is a well-formed but unknown language identifier
// or another error if another error occurred.
func ParseBase(s string) (Base, error) ***REMOVED***
	if n := len(s); n < 2 || 3 < n ***REMOVED***
		return Base***REMOVED******REMOVED***, errSyntax
	***REMOVED***
	var buf [3]byte
	l, err := getLangID(buf[:copy(buf[:], s)])
	return Base***REMOVED***l***REMOVED***, err
***REMOVED***

// Script is a 4-letter ISO 15924 code for representing scripts.
// It is idiomatically represented in title case.
type Script struct ***REMOVED***
	scriptID
***REMOVED***

// ParseScript parses a 4-letter ISO 15924 code.
// It returns a ValueError if s is a well-formed but unknown script identifier
// or another error if another error occurred.
func ParseScript(s string) (Script, error) ***REMOVED***
	if len(s) != 4 ***REMOVED***
		return Script***REMOVED******REMOVED***, errSyntax
	***REMOVED***
	var buf [4]byte
	sc, err := getScriptID(script, buf[:copy(buf[:], s)])
	return Script***REMOVED***sc***REMOVED***, err
***REMOVED***

// Region is an ISO 3166-1 or UN M.49 code for representing countries and regions.
type Region struct ***REMOVED***
	regionID
***REMOVED***

// EncodeM49 returns the Region for the given UN M.49 code.
// It returns an error if r is not a valid code.
func EncodeM49(r int) (Region, error) ***REMOVED***
	rid, err := getRegionM49(r)
	return Region***REMOVED***rid***REMOVED***, err
***REMOVED***

// ParseRegion parses a 2- or 3-letter ISO 3166-1 or a UN M.49 code.
// It returns a ValueError if s is a well-formed but unknown region identifier
// or another error if another error occurred.
func ParseRegion(s string) (Region, error) ***REMOVED***
	if n := len(s); n < 2 || 3 < n ***REMOVED***
		return Region***REMOVED******REMOVED***, errSyntax
	***REMOVED***
	var buf [3]byte
	r, err := getRegionID(buf[:copy(buf[:], s)])
	return Region***REMOVED***r***REMOVED***, err
***REMOVED***

// IsCountry returns whether this region is a country or autonomous area. This
// includes non-standard definitions from CLDR.
func (r Region) IsCountry() bool ***REMOVED***
	if r.regionID == 0 || r.IsGroup() || r.IsPrivateUse() && r.regionID != _XK ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// IsGroup returns whether this region defines a collection of regions. This
// includes non-standard definitions from CLDR.
func (r Region) IsGroup() bool ***REMOVED***
	if r.regionID == 0 ***REMOVED***
		return false
	***REMOVED***
	return int(regionInclusion[r.regionID]) < len(regionContainment)
***REMOVED***

// Contains returns whether Region c is contained by Region r. It returns true
// if c == r.
func (r Region) Contains(c Region) bool ***REMOVED***
	return r.regionID.contains(c.regionID)
***REMOVED***

func (r regionID) contains(c regionID) bool ***REMOVED***
	if r == c ***REMOVED***
		return true
	***REMOVED***
	g := regionInclusion[r]
	if g >= nRegionGroups ***REMOVED***
		return false
	***REMOVED***
	m := regionContainment[g]

	d := regionInclusion[c]
	b := regionInclusionBits[d]

	// A contained country may belong to multiple disjoint groups. Matching any
	// of these indicates containment. If the contained region is a group, it
	// must strictly be a subset.
	if d >= nRegionGroups ***REMOVED***
		return b&m != 0
	***REMOVED***
	return b&^m == 0
***REMOVED***

var errNoTLD = errors.New("language: region is not a valid ccTLD")

// TLD returns the country code top-level domain (ccTLD). UK is returned for GB.
// In all other cases it returns either the region itself or an error.
//
// This method may return an error for a region for which there exists a
// canonical form with a ccTLD. To get that ccTLD canonicalize r first. The
// region will already be canonicalized it was obtained from a Tag that was
// obtained using any of the default methods.
func (r Region) TLD() (Region, error) ***REMOVED***
	// See http://en.wikipedia.org/wiki/Country_code_top-level_domain for the
	// difference between ISO 3166-1 and IANA ccTLD.
	if r.regionID == _GB ***REMOVED***
		r = Region***REMOVED***_UK***REMOVED***
	***REMOVED***
	if (r.typ() & ccTLD) == 0 ***REMOVED***
		return Region***REMOVED******REMOVED***, errNoTLD
	***REMOVED***
	return r, nil
***REMOVED***

// Canonicalize returns the region or a possible replacement if the region is
// deprecated. It will not return a replacement for deprecated regions that
// are split into multiple regions.
func (r Region) Canonicalize() Region ***REMOVED***
	if cr := normRegion(r.regionID); cr != 0 ***REMOVED***
		return Region***REMOVED***cr***REMOVED***
	***REMOVED***
	return r
***REMOVED***

// Variant represents a registered variant of a language as defined by BCP 47.
type Variant struct ***REMOVED***
	variant string
***REMOVED***

// ParseVariant parses and returns a Variant. An error is returned if s is not
// a valid variant.
func ParseVariant(s string) (Variant, error) ***REMOVED***
	s = strings.ToLower(s)
	if _, ok := variantIndex[s]; ok ***REMOVED***
		return Variant***REMOVED***s***REMOVED***, nil
	***REMOVED***
	return Variant***REMOVED******REMOVED***, mkErrInvalid([]byte(s))
***REMOVED***

// String returns the string representation of the variant.
func (v Variant) String() string ***REMOVED***
	return v.variant
***REMOVED***
