// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen.go gen_common.go -output tables.go

package language // import "golang.org/x/text/internal/language"

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
// well-formed. The zero value of Tag is Und.
type Tag struct ***REMOVED***
	// TODO: the following fields have the form TagTypeID. This name is chosen
	// to allow refactoring the public package without conflicting with its
	// Base, Script, and Region methods. Once the transition is fully completed
	// the ID can be stripped from the name.

	LangID   Language
	RegionID Region
	// TODO: we will soon run out of positions for ScriptID. Idea: instead of
	// storing lang, region, and ScriptID codes, store only the compact index and
	// have a lookup table from this code to its expansion. This greatly speeds
	// up table lookup, speed up common variant cases.
	// This will also immediately free up 3 extra bytes. Also, the pVariant
	// field can now be moved to the lookup table, as the compact index uniquely
	// determines the offset of a possible variant.
	ScriptID Script
	pVariant byte   // offset in str, includes preceding '-'
	pExt     uint16 // offset of first extension, includes preceding '-'

	// str is the string representation of the Tag. It will only be used if the
	// tag has variants or extensions.
	str string
***REMOVED***

// Make is a convenience wrapper for Parse that omits the error.
// In case of an error, a sensible default is returned.
func Make(s string) Tag ***REMOVED***
	t, _ := Parse(s)
	return t
***REMOVED***

// Raw returns the raw base language, script and region, without making an
// attempt to infer their values.
// TODO: consider removing
func (t Tag) Raw() (b Language, s Script, r Region) ***REMOVED***
	return t.LangID, t.ScriptID, t.RegionID
***REMOVED***

// equalTags compares language, script and region subtags only.
func (t Tag) equalTags(a Tag) bool ***REMOVED***
	return t.LangID == a.LangID && t.ScriptID == a.ScriptID && t.RegionID == a.RegionID
***REMOVED***

// IsRoot returns true if t is equal to language "und".
func (t Tag) IsRoot() bool ***REMOVED***
	if int(t.pVariant) < len(t.str) ***REMOVED***
		return false
	***REMOVED***
	return t.equalTags(Und)
***REMOVED***

// IsPrivateUse reports whether the Tag consists solely of an IsPrivateUse use
// tag.
func (t Tag) IsPrivateUse() bool ***REMOVED***
	return t.str != "" && t.pVariant == 0
***REMOVED***

// RemakeString is used to update t.str in case lang, script or region changed.
// It is assumed that pExt and pVariant still point to the start of the
// respective parts.
func (t *Tag) RemakeString() ***REMOVED***
	if t.str == "" ***REMOVED***
		return
	***REMOVED***
	extra := t.str[t.pVariant:]
	if t.pVariant > 0 ***REMOVED***
		extra = extra[1:]
	***REMOVED***
	if t.equalTags(Und) && strings.HasPrefix(extra, "x-") ***REMOVED***
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
	n := t.LangID.StringToBuf(buf[:])
	if t.ScriptID != 0 ***REMOVED***
		n += copy(buf[n:], "-")
		n += copy(buf[n:], t.ScriptID.String())
	***REMOVED***
	if t.RegionID != 0 ***REMOVED***
		n += copy(buf[n:], "-")
		n += copy(buf[n:], t.RegionID.String())
	***REMOVED***
	return n
***REMOVED***

// String returns the canonical string representation of the language tag.
func (t Tag) String() string ***REMOVED***
	if t.str != "" ***REMOVED***
		return t.str
	***REMOVED***
	if t.ScriptID == 0 && t.RegionID == 0 ***REMOVED***
		return t.LangID.String()
	***REMOVED***
	buf := [maxCoreSize]byte***REMOVED******REMOVED***
	return string(buf[:t.genCoreBytes(buf[:])])
***REMOVED***

// MarshalText implements encoding.TextMarshaler.
func (t Tag) MarshalText() (text []byte, err error) ***REMOVED***
	if t.str != "" ***REMOVED***
		text = append(text, t.str...)
	***REMOVED*** else if t.ScriptID == 0 && t.RegionID == 0 ***REMOVED***
		text = append(text, t.LangID.String()...)
	***REMOVED*** else ***REMOVED***
		buf := [maxCoreSize]byte***REMOVED******REMOVED***
		text = buf[:t.genCoreBytes(buf[:])]
	***REMOVED***
	return text, nil
***REMOVED***

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Tag) UnmarshalText(text []byte) error ***REMOVED***
	tag, err := Parse(string(text))
	*t = tag
	return err
***REMOVED***

// Variants returns the part of the tag holding all variants or the empty string
// if there are no variants defined.
func (t Tag) Variants() string ***REMOVED***
	if t.pVariant == 0 ***REMOVED***
		return ""
	***REMOVED***
	return t.str[t.pVariant:t.pExt]
***REMOVED***

// VariantOrPrivateUseTags returns variants or private use tags.
func (t Tag) VariantOrPrivateUseTags() string ***REMOVED***
	if t.pExt > 0 ***REMOVED***
		return t.str[t.pVariant:t.pExt]
	***REMOVED***
	return t.str[t.pVariant:]
***REMOVED***

// HasString reports whether this tag defines more than just the raw
// components.
func (t Tag) HasString() bool ***REMOVED***
	return t.str != ""
***REMOVED***

// Parent returns the CLDR parent of t. In CLDR, missing fields in data for a
// specific language are substituted with fields from the parent language.
// The parent for a language may change for newer versions of CLDR.
func (t Tag) Parent() Tag ***REMOVED***
	if t.str != "" ***REMOVED***
		// Strip the variants and extensions.
		b, s, r := t.Raw()
		t = Tag***REMOVED***LangID: b, ScriptID: s, RegionID: r***REMOVED***
		if t.RegionID == 0 && t.ScriptID != 0 && t.LangID != 0 ***REMOVED***
			base, _ := addTags(Tag***REMOVED***LangID: t.LangID***REMOVED***)
			if base.ScriptID == t.ScriptID ***REMOVED***
				return Tag***REMOVED***LangID: t.LangID***REMOVED***
			***REMOVED***
		***REMOVED***
		return t
	***REMOVED***
	if t.LangID != 0 ***REMOVED***
		if t.RegionID != 0 ***REMOVED***
			maxScript := t.ScriptID
			if maxScript == 0 ***REMOVED***
				max, _ := addTags(t)
				maxScript = max.ScriptID
			***REMOVED***

			for i := range parents ***REMOVED***
				if Language(parents[i].lang) == t.LangID && Script(parents[i].maxScript) == maxScript ***REMOVED***
					for _, r := range parents[i].fromRegion ***REMOVED***
						if Region(r) == t.RegionID ***REMOVED***
							return Tag***REMOVED***
								LangID:   t.LangID,
								ScriptID: Script(parents[i].script),
								RegionID: Region(parents[i].toRegion),
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Strip the script if it is the default one.
			base, _ := addTags(Tag***REMOVED***LangID: t.LangID***REMOVED***)
			if base.ScriptID != maxScript ***REMOVED***
				return Tag***REMOVED***LangID: t.LangID, ScriptID: maxScript***REMOVED***
			***REMOVED***
			return Tag***REMOVED***LangID: t.LangID***REMOVED***
		***REMOVED*** else if t.ScriptID != 0 ***REMOVED***
			// The parent for an base-script pair with a non-default script is
			// "und" instead of the base language.
			base, _ := addTags(Tag***REMOVED***LangID: t.LangID***REMOVED***)
			if base.ScriptID != t.ScriptID ***REMOVED***
				return Und
			***REMOVED***
			return Tag***REMOVED***LangID: t.LangID***REMOVED***
		***REMOVED***
	***REMOVED***
	return Und
***REMOVED***

// ParseExtension parses s as an extension and returns it on success.
func ParseExtension(s string) (ext string, err error) ***REMOVED***
	scan := makeScannerString(s)
	var end int
	if n := len(scan.token); n != 1 ***REMOVED***
		return "", ErrSyntax
	***REMOVED***
	scan.toLower(0, len(scan.b))
	end = parseExtension(&scan)
	if end != len(s) ***REMOVED***
		return "", ErrSyntax
	***REMOVED***
	return string(scan.b), nil
***REMOVED***

// HasVariants reports whether t has variants.
func (t Tag) HasVariants() bool ***REMOVED***
	return uint16(t.pVariant) < t.pExt
***REMOVED***

// HasExtensions reports whether t has extensions.
func (t Tag) HasExtensions() bool ***REMOVED***
	return int(t.pExt) < len(t.str)
***REMOVED***

// Extension returns the extension of type x for tag t. It will return
// false for ok if t does not have the requested extension. The returned
// extension will be invalid in this case.
func (t Tag) Extension(x byte) (ext string, ok bool) ***REMOVED***
	for i := int(t.pExt); i < len(t.str)-1; ***REMOVED***
		var ext string
		i, ext = getExtension(t.str, i)
		if ext[0] == x ***REMOVED***
			return ext, true
		***REMOVED***
	***REMOVED***
	return "", false
***REMOVED***

// Extensions returns all extensions of t.
func (t Tag) Extensions() []string ***REMOVED***
	e := []string***REMOVED******REMOVED***
	for i := int(t.pExt); i < len(t.str)-1; ***REMOVED***
		var ext string
		i, ext = getExtension(t.str, i)
		e = append(e, ext)
	***REMOVED***
	return e
***REMOVED***

// TypeForKey returns the type associated with the given key, where key and type
// are of the allowed values defined for the Unicode locale extension ('u') in
// https://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
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
// https://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
// An empty value removes an existing pair with the same key.
func (t Tag) SetTypeForKey(key, value string) (Tag, error) ***REMOVED***
	if t.IsPrivateUse() ***REMOVED***
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

// ParseBase parses a 2- or 3-letter ISO 639 code.
// It returns a ValueError if s is a well-formed but unknown language identifier
// or another error if another error occurred.
func ParseBase(s string) (Language, error) ***REMOVED***
	if n := len(s); n < 2 || 3 < n ***REMOVED***
		return 0, ErrSyntax
	***REMOVED***
	var buf [3]byte
	return getLangID(buf[:copy(buf[:], s)])
***REMOVED***

// ParseScript parses a 4-letter ISO 15924 code.
// It returns a ValueError if s is a well-formed but unknown script identifier
// or another error if another error occurred.
func ParseScript(s string) (Script, error) ***REMOVED***
	if len(s) != 4 ***REMOVED***
		return 0, ErrSyntax
	***REMOVED***
	var buf [4]byte
	return getScriptID(script, buf[:copy(buf[:], s)])
***REMOVED***

// EncodeM49 returns the Region for the given UN M.49 code.
// It returns an error if r is not a valid code.
func EncodeM49(r int) (Region, error) ***REMOVED***
	return getRegionM49(r)
***REMOVED***

// ParseRegion parses a 2- or 3-letter ISO 3166-1 or a UN M.49 code.
// It returns a ValueError if s is a well-formed but unknown region identifier
// or another error if another error occurred.
func ParseRegion(s string) (Region, error) ***REMOVED***
	if n := len(s); n < 2 || 3 < n ***REMOVED***
		return 0, ErrSyntax
	***REMOVED***
	var buf [3]byte
	return getRegionID(buf[:copy(buf[:], s)])
***REMOVED***

// IsCountry returns whether this region is a country or autonomous area. This
// includes non-standard definitions from CLDR.
func (r Region) IsCountry() bool ***REMOVED***
	if r == 0 || r.IsGroup() || r.IsPrivateUse() && r != _XK ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// IsGroup returns whether this region defines a collection of regions. This
// includes non-standard definitions from CLDR.
func (r Region) IsGroup() bool ***REMOVED***
	if r == 0 ***REMOVED***
		return false
	***REMOVED***
	return int(regionInclusion[r]) < len(regionContainment)
***REMOVED***

// Contains returns whether Region c is contained by Region r. It returns true
// if c == r.
func (r Region) Contains(c Region) bool ***REMOVED***
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
	if r == _GB ***REMOVED***
		r = _UK
	***REMOVED***
	if (r.typ() & ccTLD) == 0 ***REMOVED***
		return 0, errNoTLD
	***REMOVED***
	return r, nil
***REMOVED***

// Canonicalize returns the region or a possible replacement if the region is
// deprecated. It will not return a replacement for deprecated regions that
// are split into multiple regions.
func (r Region) Canonicalize() Region ***REMOVED***
	if cr := normRegion(r); cr != 0 ***REMOVED***
		return cr
	***REMOVED***
	return r
***REMOVED***

// Variant represents a registered variant of a language as defined by BCP 47.
type Variant struct ***REMOVED***
	ID  uint8
	str string
***REMOVED***

// ParseVariant parses and returns a Variant. An error is returned if s is not
// a valid variant.
func ParseVariant(s string) (Variant, error) ***REMOVED***
	s = strings.ToLower(s)
	if id, ok := variantIndex[s]; ok ***REMOVED***
		return Variant***REMOVED***id, s***REMOVED***, nil
	***REMOVED***
	return Variant***REMOVED******REMOVED***, NewValueError([]byte(s))
***REMOVED***

// String returns the string representation of the variant.
func (v Variant) String() string ***REMOVED***
	return v.str
***REMOVED***
