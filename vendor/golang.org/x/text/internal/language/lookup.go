// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"golang.org/x/text/internal/tag"
)

// findIndex tries to find the given tag in idx and returns a standardized error
// if it could not be found.
func findIndex(idx tag.Index, key []byte, form string) (index int, err error) ***REMOVED***
	if !tag.FixCase(form, key) ***REMOVED***
		return 0, ErrSyntax
	***REMOVED***
	i := idx.Index(key)
	if i == -1 ***REMOVED***
		return 0, NewValueError(key)
	***REMOVED***
	return i, nil
***REMOVED***

func searchUint(imap []uint16, key uint16) int ***REMOVED***
	return sort.Search(len(imap), func(i int) bool ***REMOVED***
		return imap[i] >= key
	***REMOVED***)
***REMOVED***

type Language uint16

// getLangID returns the langID of s if s is a canonical subtag
// or langUnknown if s is not a canonical subtag.
func getLangID(s []byte) (Language, error) ***REMOVED***
	if len(s) == 2 ***REMOVED***
		return getLangISO2(s)
	***REMOVED***
	return getLangISO3(s)
***REMOVED***

// TODO language normalization as well as the AliasMaps could be moved to the
// higher level package, but it is a bit tricky to separate the generation.

func (id Language) Canonicalize() (Language, AliasType) ***REMOVED***
	return normLang(id)
***REMOVED***

// mapLang returns the mapped langID of id according to mapping m.
func normLang(id Language) (Language, AliasType) ***REMOVED***
	k := sort.Search(len(AliasMap), func(i int) bool ***REMOVED***
		return AliasMap[i].From >= uint16(id)
	***REMOVED***)
	if k < len(AliasMap) && AliasMap[k].From == uint16(id) ***REMOVED***
		return Language(AliasMap[k].To), AliasTypes[k]
	***REMOVED***
	return id, AliasTypeUnknown
***REMOVED***

// getLangISO2 returns the langID for the given 2-letter ISO language code
// or unknownLang if this does not exist.
func getLangISO2(s []byte) (Language, error) ***REMOVED***
	if !tag.FixCase("zz", s) ***REMOVED***
		return 0, ErrSyntax
	***REMOVED***
	if i := lang.Index(s); i != -1 && lang.Elem(i)[3] != 0 ***REMOVED***
		return Language(i), nil
	***REMOVED***
	return 0, NewValueError(s)
***REMOVED***

const base = 'z' - 'a' + 1

func strToInt(s []byte) uint ***REMOVED***
	v := uint(0)
	for i := 0; i < len(s); i++ ***REMOVED***
		v *= base
		v += uint(s[i] - 'a')
	***REMOVED***
	return v
***REMOVED***

// converts the given integer to the original ASCII string passed to strToInt.
// len(s) must match the number of characters obtained.
func intToStr(v uint, s []byte) ***REMOVED***
	for i := len(s) - 1; i >= 0; i-- ***REMOVED***
		s[i] = byte(v%base) + 'a'
		v /= base
	***REMOVED***
***REMOVED***

// getLangISO3 returns the langID for the given 3-letter ISO language code
// or unknownLang if this does not exist.
func getLangISO3(s []byte) (Language, error) ***REMOVED***
	if tag.FixCase("und", s) ***REMOVED***
		// first try to match canonical 3-letter entries
		for i := lang.Index(s[:2]); i != -1; i = lang.Next(s[:2], i) ***REMOVED***
			if e := lang.Elem(i); e[3] == 0 && e[2] == s[2] ***REMOVED***
				// We treat "und" as special and always translate it to "unspecified".
				// Note that ZZ and Zzzz are private use and are not treated as
				// unspecified by default.
				id := Language(i)
				if id == nonCanonicalUnd ***REMOVED***
					return 0, nil
				***REMOVED***
				return id, nil
			***REMOVED***
		***REMOVED***
		if i := altLangISO3.Index(s); i != -1 ***REMOVED***
			return Language(altLangIndex[altLangISO3.Elem(i)[3]]), nil
		***REMOVED***
		n := strToInt(s)
		if langNoIndex[n/8]&(1<<(n%8)) != 0 ***REMOVED***
			return Language(n) + langNoIndexOffset, nil
		***REMOVED***
		// Check for non-canonical uses of ISO3.
		for i := lang.Index(s[:1]); i != -1; i = lang.Next(s[:1], i) ***REMOVED***
			if e := lang.Elem(i); e[2] == s[1] && e[3] == s[2] ***REMOVED***
				return Language(i), nil
			***REMOVED***
		***REMOVED***
		return 0, NewValueError(s)
	***REMOVED***
	return 0, ErrSyntax
***REMOVED***

// StringToBuf writes the string to b and returns the number of bytes
// written.  cap(b) must be >= 3.
func (id Language) StringToBuf(b []byte) int ***REMOVED***
	if id >= langNoIndexOffset ***REMOVED***
		intToStr(uint(id)-langNoIndexOffset, b[:3])
		return 3
	***REMOVED*** else if id == 0 ***REMOVED***
		return copy(b, "und")
	***REMOVED***
	l := lang[id<<2:]
	if l[3] == 0 ***REMOVED***
		return copy(b, l[:3])
	***REMOVED***
	return copy(b, l[:2])
***REMOVED***

// String returns the BCP 47 representation of the langID.
// Use b as variable name, instead of id, to ensure the variable
// used is consistent with that of Base in which this type is embedded.
func (b Language) String() string ***REMOVED***
	if b == 0 ***REMOVED***
		return "und"
	***REMOVED*** else if b >= langNoIndexOffset ***REMOVED***
		b -= langNoIndexOffset
		buf := [3]byte***REMOVED******REMOVED***
		intToStr(uint(b), buf[:])
		return string(buf[:])
	***REMOVED***
	l := lang.Elem(int(b))
	if l[3] == 0 ***REMOVED***
		return l[:3]
	***REMOVED***
	return l[:2]
***REMOVED***

// ISO3 returns the ISO 639-3 language code.
func (b Language) ISO3() string ***REMOVED***
	if b == 0 || b >= langNoIndexOffset ***REMOVED***
		return b.String()
	***REMOVED***
	l := lang.Elem(int(b))
	if l[3] == 0 ***REMOVED***
		return l[:3]
	***REMOVED*** else if l[2] == 0 ***REMOVED***
		return altLangISO3.Elem(int(l[3]))[:3]
	***REMOVED***
	// This allocation will only happen for 3-letter ISO codes
	// that are non-canonical BCP 47 language identifiers.
	return l[0:1] + l[2:4]
***REMOVED***

// IsPrivateUse reports whether this language code is reserved for private use.
func (b Language) IsPrivateUse() bool ***REMOVED***
	return langPrivateStart <= b && b <= langPrivateEnd
***REMOVED***

// SuppressScript returns the script marked as SuppressScript in the IANA
// language tag repository, or 0 if there is no such script.
func (b Language) SuppressScript() Script ***REMOVED***
	if b < langNoIndexOffset ***REMOVED***
		return Script(suppressScript[b])
	***REMOVED***
	return 0
***REMOVED***

type Region uint16

// getRegionID returns the region id for s if s is a valid 2-letter region code
// or unknownRegion.
func getRegionID(s []byte) (Region, error) ***REMOVED***
	if len(s) == 3 ***REMOVED***
		if isAlpha(s[0]) ***REMOVED***
			return getRegionISO3(s)
		***REMOVED***
		if i, err := strconv.ParseUint(string(s), 10, 10); err == nil ***REMOVED***
			return getRegionM49(int(i))
		***REMOVED***
	***REMOVED***
	return getRegionISO2(s)
***REMOVED***

// getRegionISO2 returns the regionID for the given 2-letter ISO country code
// or unknownRegion if this does not exist.
func getRegionISO2(s []byte) (Region, error) ***REMOVED***
	i, err := findIndex(regionISO, s, "ZZ")
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return Region(i) + isoRegionOffset, nil
***REMOVED***

// getRegionISO3 returns the regionID for the given 3-letter ISO country code
// or unknownRegion if this does not exist.
func getRegionISO3(s []byte) (Region, error) ***REMOVED***
	if tag.FixCase("ZZZ", s) ***REMOVED***
		for i := regionISO.Index(s[:1]); i != -1; i = regionISO.Next(s[:1], i) ***REMOVED***
			if e := regionISO.Elem(i); e[2] == s[1] && e[3] == s[2] ***REMOVED***
				return Region(i) + isoRegionOffset, nil
			***REMOVED***
		***REMOVED***
		for i := 0; i < len(altRegionISO3); i += 3 ***REMOVED***
			if tag.Compare(altRegionISO3[i:i+3], s) == 0 ***REMOVED***
				return Region(altRegionIDs[i/3]), nil
			***REMOVED***
		***REMOVED***
		return 0, NewValueError(s)
	***REMOVED***
	return 0, ErrSyntax
***REMOVED***

func getRegionM49(n int) (Region, error) ***REMOVED***
	if 0 < n && n <= 999 ***REMOVED***
		const (
			searchBits = 7
			regionBits = 9
			regionMask = 1<<regionBits - 1
		)
		idx := n >> searchBits
		buf := fromM49[m49Index[idx]:m49Index[idx+1]]
		val := uint16(n) << regionBits // we rely on bits shifting out
		i := sort.Search(len(buf), func(i int) bool ***REMOVED***
			return buf[i] >= val
		***REMOVED***)
		if r := fromM49[int(m49Index[idx])+i]; r&^regionMask == val ***REMOVED***
			return Region(r & regionMask), nil
		***REMOVED***
	***REMOVED***
	var e ValueError
	fmt.Fprint(bytes.NewBuffer([]byte(e.v[:])), n)
	return 0, e
***REMOVED***

// normRegion returns a region if r is deprecated or 0 otherwise.
// TODO: consider supporting BYS (-> BLR), CSK (-> 200 or CZ), PHI (-> PHL) and AFI (-> DJ).
// TODO: consider mapping split up regions to new most populous one (like CLDR).
func normRegion(r Region) Region ***REMOVED***
	m := regionOldMap
	k := sort.Search(len(m), func(i int) bool ***REMOVED***
		return m[i].From >= uint16(r)
	***REMOVED***)
	if k < len(m) && m[k].From == uint16(r) ***REMOVED***
		return Region(m[k].To)
	***REMOVED***
	return 0
***REMOVED***

const (
	iso3166UserAssigned = 1 << iota
	ccTLD
	bcp47Region
)

func (r Region) typ() byte ***REMOVED***
	return regionTypes[r]
***REMOVED***

// String returns the BCP 47 representation for the region.
// It returns "ZZ" for an unspecified region.
func (r Region) String() string ***REMOVED***
	if r < isoRegionOffset ***REMOVED***
		if r == 0 ***REMOVED***
			return "ZZ"
		***REMOVED***
		return fmt.Sprintf("%03d", r.M49())
	***REMOVED***
	r -= isoRegionOffset
	return regionISO.Elem(int(r))[:2]
***REMOVED***

// ISO3 returns the 3-letter ISO code of r.
// Note that not all regions have a 3-letter ISO code.
// In such cases this method returns "ZZZ".
func (r Region) ISO3() string ***REMOVED***
	if r < isoRegionOffset ***REMOVED***
		return "ZZZ"
	***REMOVED***
	r -= isoRegionOffset
	reg := regionISO.Elem(int(r))
	switch reg[2] ***REMOVED***
	case 0:
		return altRegionISO3[reg[3]:][:3]
	case ' ':
		return "ZZZ"
	***REMOVED***
	return reg[0:1] + reg[2:4]
***REMOVED***

// M49 returns the UN M.49 encoding of r, or 0 if this encoding
// is not defined for r.
func (r Region) M49() int ***REMOVED***
	return int(m49[r])
***REMOVED***

// IsPrivateUse reports whether r has the ISO 3166 User-assigned status. This
// may include private-use tags that are assigned by CLDR and used in this
// implementation. So IsPrivateUse and IsCountry can be simultaneously true.
func (r Region) IsPrivateUse() bool ***REMOVED***
	return r.typ()&iso3166UserAssigned != 0
***REMOVED***

type Script uint8

// getScriptID returns the script id for string s. It assumes that s
// is of the format [A-Z][a-z]***REMOVED***3***REMOVED***.
func getScriptID(idx tag.Index, s []byte) (Script, error) ***REMOVED***
	i, err := findIndex(idx, s, "Zzzz")
	return Script(i), err
***REMOVED***

// String returns the script code in title case.
// It returns "Zzzz" for an unspecified script.
func (s Script) String() string ***REMOVED***
	if s == 0 ***REMOVED***
		return "Zzzz"
	***REMOVED***
	return script.Elem(int(s))
***REMOVED***

// IsPrivateUse reports whether this script code is reserved for private use.
func (s Script) IsPrivateUse() bool ***REMOVED***
	return _Qaaa <= s && s <= _Qabx
***REMOVED***

const (
	maxAltTaglen = len("en-US-POSIX")
	maxLen       = maxAltTaglen
)

var (
	// grandfatheredMap holds a mapping from legacy and grandfathered tags to
	// their base language or index to more elaborate tag.
	grandfatheredMap = map[[maxLen]byte]int16***REMOVED***
		[maxLen]byte***REMOVED***'a', 'r', 't', '-', 'l', 'o', 'j', 'b', 'a', 'n'***REMOVED***: _jbo, // art-lojban
		[maxLen]byte***REMOVED***'i', '-', 'a', 'm', 'i'***REMOVED***:                          _ami, // i-ami
		[maxLen]byte***REMOVED***'i', '-', 'b', 'n', 'n'***REMOVED***:                          _bnn, // i-bnn
		[maxLen]byte***REMOVED***'i', '-', 'h', 'a', 'k'***REMOVED***:                          _hak, // i-hak
		[maxLen]byte***REMOVED***'i', '-', 'k', 'l', 'i', 'n', 'g', 'o', 'n'***REMOVED***:      _tlh, // i-klingon
		[maxLen]byte***REMOVED***'i', '-', 'l', 'u', 'x'***REMOVED***:                          _lb,  // i-lux
		[maxLen]byte***REMOVED***'i', '-', 'n', 'a', 'v', 'a', 'j', 'o'***REMOVED***:           _nv,  // i-navajo
		[maxLen]byte***REMOVED***'i', '-', 'p', 'w', 'n'***REMOVED***:                          _pwn, // i-pwn
		[maxLen]byte***REMOVED***'i', '-', 't', 'a', 'o'***REMOVED***:                          _tao, // i-tao
		[maxLen]byte***REMOVED***'i', '-', 't', 'a', 'y'***REMOVED***:                          _tay, // i-tay
		[maxLen]byte***REMOVED***'i', '-', 't', 's', 'u'***REMOVED***:                          _tsu, // i-tsu
		[maxLen]byte***REMOVED***'n', 'o', '-', 'b', 'o', 'k'***REMOVED***:                     _nb,  // no-bok
		[maxLen]byte***REMOVED***'n', 'o', '-', 'n', 'y', 'n'***REMOVED***:                     _nn,  // no-nyn
		[maxLen]byte***REMOVED***'s', 'g', 'n', '-', 'b', 'e', '-', 'f', 'r'***REMOVED***:      _sfb, // sgn-BE-FR
		[maxLen]byte***REMOVED***'s', 'g', 'n', '-', 'b', 'e', '-', 'n', 'l'***REMOVED***:      _vgt, // sgn-BE-NL
		[maxLen]byte***REMOVED***'s', 'g', 'n', '-', 'c', 'h', '-', 'd', 'e'***REMOVED***:      _sgg, // sgn-CH-DE
		[maxLen]byte***REMOVED***'z', 'h', '-', 'g', 'u', 'o', 'y', 'u'***REMOVED***:           _cmn, // zh-guoyu
		[maxLen]byte***REMOVED***'z', 'h', '-', 'h', 'a', 'k', 'k', 'a'***REMOVED***:           _hak, // zh-hakka
		[maxLen]byte***REMOVED***'z', 'h', '-', 'm', 'i', 'n', '-', 'n', 'a', 'n'***REMOVED***: _nan, // zh-min-nan
		[maxLen]byte***REMOVED***'z', 'h', '-', 'x', 'i', 'a', 'n', 'g'***REMOVED***:           _hsn, // zh-xiang

		// Grandfathered tags with no modern replacement will be converted as
		// follows:
		[maxLen]byte***REMOVED***'c', 'e', 'l', '-', 'g', 'a', 'u', 'l', 'i', 's', 'h'***REMOVED***: -1, // cel-gaulish
		[maxLen]byte***REMOVED***'e', 'n', '-', 'g', 'b', '-', 'o', 'e', 'd'***REMOVED***:           -2, // en-GB-oed
		[maxLen]byte***REMOVED***'i', '-', 'd', 'e', 'f', 'a', 'u', 'l', 't'***REMOVED***:           -3, // i-default
		[maxLen]byte***REMOVED***'i', '-', 'e', 'n', 'o', 'c', 'h', 'i', 'a', 'n'***REMOVED***:      -4, // i-enochian
		[maxLen]byte***REMOVED***'i', '-', 'm', 'i', 'n', 'g', 'o'***REMOVED***:                     -5, // i-mingo
		[maxLen]byte***REMOVED***'z', 'h', '-', 'm', 'i', 'n'***REMOVED***:                          -6, // zh-min

		// CLDR-specific tag.
		[maxLen]byte***REMOVED***'r', 'o', 'o', 't'***REMOVED***:                                    0,  // root
		[maxLen]byte***REMOVED***'e', 'n', '-', 'u', 's', '-', 'p', 'o', 's', 'i', 'x'***REMOVED***: -7, // en_US_POSIX"
	***REMOVED***

	altTagIndex = [...]uint8***REMOVED***0, 17, 31, 45, 61, 74, 86, 102***REMOVED***

	altTags = "xtg-x-cel-gaulishen-GB-oxendicten-x-i-defaultund-x-i-enochiansee-x-i-mingonan-x-zh-minen-US-u-va-posix"
)

func grandfathered(s [maxAltTaglen]byte) (t Tag, ok bool) ***REMOVED***
	if v, ok := grandfatheredMap[s]; ok ***REMOVED***
		if v < 0 ***REMOVED***
			return Make(altTags[altTagIndex[-v-1]:altTagIndex[-v]]), true
		***REMOVED***
		t.LangID = Language(v)
		return t, true
	***REMOVED***
	return t, false
***REMOVED***
