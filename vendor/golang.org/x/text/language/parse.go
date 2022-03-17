// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"errors"
	"strconv"
	"strings"

	"golang.org/x/text/internal/language"
)

// ValueError is returned by any of the parsing functions when the
// input is well-formed but the respective subtag is not recognized
// as a valid value.
type ValueError interface ***REMOVED***
	error

	// Subtag returns the subtag for which the error occurred.
	Subtag() string
***REMOVED***

// Parse parses the given BCP 47 string and returns a valid Tag. If parsing
// failed it returns an error and any part of the tag that could be parsed.
// If parsing succeeded but an unknown value was found, it returns
// ValueError. The Tag returned in this case is just stripped of the unknown
// value. All other values are preserved. It accepts tags in the BCP 47 format
// and extensions to this standard defined in
// https://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
// The resulting tag is canonicalized using the default canonicalization type.
func Parse(s string) (t Tag, err error) ***REMOVED***
	return Default.Parse(s)
***REMOVED***

// Parse parses the given BCP 47 string and returns a valid Tag. If parsing
// failed it returns an error and any part of the tag that could be parsed.
// If parsing succeeded but an unknown value was found, it returns
// ValueError. The Tag returned in this case is just stripped of the unknown
// value. All other values are preserved. It accepts tags in the BCP 47 format
// and extensions to this standard defined in
// https://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
// The resulting tag is canonicalized using the canonicalization type c.
func (c CanonType) Parse(s string) (t Tag, err error) ***REMOVED***
	defer func() ***REMOVED***
		if recover() != nil ***REMOVED***
			t = Tag***REMOVED******REMOVED***
			err = language.ErrSyntax
		***REMOVED***
	***REMOVED***()

	tt, err := language.Parse(s)
	if err != nil ***REMOVED***
		return makeTag(tt), err
	***REMOVED***
	tt, changed := canonicalize(c, tt)
	if changed ***REMOVED***
		tt.RemakeString()
	***REMOVED***
	return makeTag(tt), err
***REMOVED***

// Compose creates a Tag from individual parts, which may be of type Tag, Base,
// Script, Region, Variant, []Variant, Extension, []Extension or error. If a
// Base, Script or Region or slice of type Variant or Extension is passed more
// than once, the latter will overwrite the former. Variants and Extensions are
// accumulated, but if two extensions of the same type are passed, the latter
// will replace the former. For -u extensions, though, the key-type pairs are
// added, where later values overwrite older ones. A Tag overwrites all former
// values and typically only makes sense as the first argument. The resulting
// tag is returned after canonicalizing using the Default CanonType. If one or
// more errors are encountered, one of the errors is returned.
func Compose(part ...interface***REMOVED******REMOVED***) (t Tag, err error) ***REMOVED***
	return Default.Compose(part...)
***REMOVED***

// Compose creates a Tag from individual parts, which may be of type Tag, Base,
// Script, Region, Variant, []Variant, Extension, []Extension or error. If a
// Base, Script or Region or slice of type Variant or Extension is passed more
// than once, the latter will overwrite the former. Variants and Extensions are
// accumulated, but if two extensions of the same type are passed, the latter
// will replace the former. For -u extensions, though, the key-type pairs are
// added, where later values overwrite older ones. A Tag overwrites all former
// values and typically only makes sense as the first argument. The resulting
// tag is returned after canonicalizing using CanonType c. If one or more errors
// are encountered, one of the errors is returned.
func (c CanonType) Compose(part ...interface***REMOVED******REMOVED***) (t Tag, err error) ***REMOVED***
	defer func() ***REMOVED***
		if recover() != nil ***REMOVED***
			t = Tag***REMOVED******REMOVED***
			err = language.ErrSyntax
		***REMOVED***
	***REMOVED***()

	var b language.Builder
	if err = update(&b, part...); err != nil ***REMOVED***
		return und, err
	***REMOVED***
	b.Tag, _ = canonicalize(c, b.Tag)
	return makeTag(b.Make()), err
***REMOVED***

var errInvalidArgument = errors.New("invalid Extension or Variant")

func update(b *language.Builder, part ...interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	for _, x := range part ***REMOVED***
		switch v := x.(type) ***REMOVED***
		case Tag:
			b.SetTag(v.tag())
		case Base:
			b.Tag.LangID = v.langID
		case Script:
			b.Tag.ScriptID = v.scriptID
		case Region:
			b.Tag.RegionID = v.regionID
		case Variant:
			if v.variant == "" ***REMOVED***
				err = errInvalidArgument
				break
			***REMOVED***
			b.AddVariant(v.variant)
		case Extension:
			if v.s == "" ***REMOVED***
				err = errInvalidArgument
				break
			***REMOVED***
			b.SetExt(v.s)
		case []Variant:
			b.ClearVariants()
			for _, v := range v ***REMOVED***
				b.AddVariant(v.variant)
			***REMOVED***
		case []Extension:
			b.ClearExtensions()
			for _, e := range v ***REMOVED***
				b.SetExt(e.s)
			***REMOVED***
		// TODO: support parsing of raw strings based on morphology or just extensions?
		case error:
			if v != nil ***REMOVED***
				err = v
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

var errInvalidWeight = errors.New("ParseAcceptLanguage: invalid weight")

// ParseAcceptLanguage parses the contents of an Accept-Language header as
// defined in http://www.ietf.org/rfc/rfc2616.txt and returns a list of Tags and
// a list of corresponding quality weights. It is more permissive than RFC 2616
// and may return non-nil slices even if the input is not valid.
// The Tags will be sorted by highest weight first and then by first occurrence.
// Tags with a weight of zero will be dropped. An error will be returned if the
// input could not be parsed.
func ParseAcceptLanguage(s string) (tag []Tag, q []float32, err error) ***REMOVED***
	defer func() ***REMOVED***
		if recover() != nil ***REMOVED***
			tag = nil
			q = nil
			err = language.ErrSyntax
		***REMOVED***
	***REMOVED***()

	var entry string
	for s != "" ***REMOVED***
		if entry, s = split(s, ','); entry == "" ***REMOVED***
			continue
		***REMOVED***

		entry, weight := split(entry, ';')

		// Scan the language.
		t, err := Parse(entry)
		if err != nil ***REMOVED***
			id, ok := acceptFallback[entry]
			if !ok ***REMOVED***
				return nil, nil, err
			***REMOVED***
			t = makeTag(language.Tag***REMOVED***LangID: id***REMOVED***)
		***REMOVED***

		// Scan the optional weight.
		w := 1.0
		if weight != "" ***REMOVED***
			weight = consume(weight, 'q')
			weight = consume(weight, '=')
			// consume returns the empty string when a token could not be
			// consumed, resulting in an error for ParseFloat.
			if w, err = strconv.ParseFloat(weight, 32); err != nil ***REMOVED***
				return nil, nil, errInvalidWeight
			***REMOVED***
			// Drop tags with a quality weight of 0.
			if w <= 0 ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		tag = append(tag, t)
		q = append(q, float32(w))
	***REMOVED***
	sortStable(&tagSort***REMOVED***tag, q***REMOVED***)
	return tag, q, nil
***REMOVED***

// consume removes a leading token c from s and returns the result or the empty
// string if there is no such token.
func consume(s string, c byte) string ***REMOVED***
	if s == "" || s[0] != c ***REMOVED***
		return ""
	***REMOVED***
	return strings.TrimSpace(s[1:])
***REMOVED***

func split(s string, c byte) (head, tail string) ***REMOVED***
	if i := strings.IndexByte(s, c); i >= 0 ***REMOVED***
		return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:])
	***REMOVED***
	return strings.TrimSpace(s), ""
***REMOVED***

// Add hack mapping to deal with a small number of cases that occur
// in Accept-Language (with reasonable frequency).
var acceptFallback = map[string]language.Language***REMOVED***
	"english": _en,
	"deutsch": _de,
	"italian": _it,
	"french":  _fr,
	"*":       _mul, // defined in the spec to match all languages.
***REMOVED***

type tagSort struct ***REMOVED***
	tag []Tag
	q   []float32
***REMOVED***

func (s *tagSort) Len() int ***REMOVED***
	return len(s.q)
***REMOVED***

func (s *tagSort) Less(i, j int) bool ***REMOVED***
	return s.q[i] > s.q[j]
***REMOVED***

func (s *tagSort) Swap(i, j int) ***REMOVED***
	s.tag[i], s.tag[j] = s.tag[j], s.tag[i]
	s.q[i], s.q[j] = s.q[j], s.q[i]
***REMOVED***
