// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

// TODO: Various sets of commonly use tags and regions.

// MustParse is like Parse, but panics if the given BCP 47 tag cannot be parsed.
// It simplifies safe initialization of Tag values.
func MustParse(s string) Tag ***REMOVED***
	t, err := Parse(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return t
***REMOVED***

// MustParse is like Parse, but panics if the given BCP 47 tag cannot be parsed.
// It simplifies safe initialization of Tag values.
func (c CanonType) MustParse(s string) Tag ***REMOVED***
	t, err := c.Parse(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return t
***REMOVED***

// MustParseBase is like ParseBase, but panics if the given base cannot be parsed.
// It simplifies safe initialization of Base values.
func MustParseBase(s string) Base ***REMOVED***
	b, err := ParseBase(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return b
***REMOVED***

// MustParseScript is like ParseScript, but panics if the given script cannot be
// parsed. It simplifies safe initialization of Script values.
func MustParseScript(s string) Script ***REMOVED***
	scr, err := ParseScript(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return scr
***REMOVED***

// MustParseRegion is like ParseRegion, but panics if the given region cannot be
// parsed. It simplifies safe initialization of Region values.
func MustParseRegion(s string) Region ***REMOVED***
	r, err := ParseRegion(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return r
***REMOVED***

var (
	und = Tag***REMOVED******REMOVED***

	Und Tag = Tag***REMOVED******REMOVED***

	Afrikaans            Tag = Tag***REMOVED***lang: _af***REMOVED***                //  af
	Amharic              Tag = Tag***REMOVED***lang: _am***REMOVED***                //  am
	Arabic               Tag = Tag***REMOVED***lang: _ar***REMOVED***                //  ar
	ModernStandardArabic Tag = Tag***REMOVED***lang: _ar, region: _001***REMOVED***  //  ar-001
	Azerbaijani          Tag = Tag***REMOVED***lang: _az***REMOVED***                //  az
	Bulgarian            Tag = Tag***REMOVED***lang: _bg***REMOVED***                //  bg
	Bengali              Tag = Tag***REMOVED***lang: _bn***REMOVED***                //  bn
	Catalan              Tag = Tag***REMOVED***lang: _ca***REMOVED***                //  ca
	Czech                Tag = Tag***REMOVED***lang: _cs***REMOVED***                //  cs
	Danish               Tag = Tag***REMOVED***lang: _da***REMOVED***                //  da
	German               Tag = Tag***REMOVED***lang: _de***REMOVED***                //  de
	Greek                Tag = Tag***REMOVED***lang: _el***REMOVED***                //  el
	English              Tag = Tag***REMOVED***lang: _en***REMOVED***                //  en
	AmericanEnglish      Tag = Tag***REMOVED***lang: _en, region: _US***REMOVED***   //  en-US
	BritishEnglish       Tag = Tag***REMOVED***lang: _en, region: _GB***REMOVED***   //  en-GB
	Spanish              Tag = Tag***REMOVED***lang: _es***REMOVED***                //  es
	EuropeanSpanish      Tag = Tag***REMOVED***lang: _es, region: _ES***REMOVED***   //  es-ES
	LatinAmericanSpanish Tag = Tag***REMOVED***lang: _es, region: _419***REMOVED***  //  es-419
	Estonian             Tag = Tag***REMOVED***lang: _et***REMOVED***                //  et
	Persian              Tag = Tag***REMOVED***lang: _fa***REMOVED***                //  fa
	Finnish              Tag = Tag***REMOVED***lang: _fi***REMOVED***                //  fi
	Filipino             Tag = Tag***REMOVED***lang: _fil***REMOVED***               //  fil
	French               Tag = Tag***REMOVED***lang: _fr***REMOVED***                //  fr
	CanadianFrench       Tag = Tag***REMOVED***lang: _fr, region: _CA***REMOVED***   //  fr-CA
	Gujarati             Tag = Tag***REMOVED***lang: _gu***REMOVED***                //  gu
	Hebrew               Tag = Tag***REMOVED***lang: _he***REMOVED***                //  he
	Hindi                Tag = Tag***REMOVED***lang: _hi***REMOVED***                //  hi
	Croatian             Tag = Tag***REMOVED***lang: _hr***REMOVED***                //  hr
	Hungarian            Tag = Tag***REMOVED***lang: _hu***REMOVED***                //  hu
	Armenian             Tag = Tag***REMOVED***lang: _hy***REMOVED***                //  hy
	Indonesian           Tag = Tag***REMOVED***lang: _id***REMOVED***                //  id
	Icelandic            Tag = Tag***REMOVED***lang: _is***REMOVED***                //  is
	Italian              Tag = Tag***REMOVED***lang: _it***REMOVED***                //  it
	Japanese             Tag = Tag***REMOVED***lang: _ja***REMOVED***                //  ja
	Georgian             Tag = Tag***REMOVED***lang: _ka***REMOVED***                //  ka
	Kazakh               Tag = Tag***REMOVED***lang: _kk***REMOVED***                //  kk
	Khmer                Tag = Tag***REMOVED***lang: _km***REMOVED***                //  km
	Kannada              Tag = Tag***REMOVED***lang: _kn***REMOVED***                //  kn
	Korean               Tag = Tag***REMOVED***lang: _ko***REMOVED***                //  ko
	Kirghiz              Tag = Tag***REMOVED***lang: _ky***REMOVED***                //  ky
	Lao                  Tag = Tag***REMOVED***lang: _lo***REMOVED***                //  lo
	Lithuanian           Tag = Tag***REMOVED***lang: _lt***REMOVED***                //  lt
	Latvian              Tag = Tag***REMOVED***lang: _lv***REMOVED***                //  lv
	Macedonian           Tag = Tag***REMOVED***lang: _mk***REMOVED***                //  mk
	Malayalam            Tag = Tag***REMOVED***lang: _ml***REMOVED***                //  ml
	Mongolian            Tag = Tag***REMOVED***lang: _mn***REMOVED***                //  mn
	Marathi              Tag = Tag***REMOVED***lang: _mr***REMOVED***                //  mr
	Malay                Tag = Tag***REMOVED***lang: _ms***REMOVED***                //  ms
	Burmese              Tag = Tag***REMOVED***lang: _my***REMOVED***                //  my
	Nepali               Tag = Tag***REMOVED***lang: _ne***REMOVED***                //  ne
	Dutch                Tag = Tag***REMOVED***lang: _nl***REMOVED***                //  nl
	Norwegian            Tag = Tag***REMOVED***lang: _no***REMOVED***                //  no
	Punjabi              Tag = Tag***REMOVED***lang: _pa***REMOVED***                //  pa
	Polish               Tag = Tag***REMOVED***lang: _pl***REMOVED***                //  pl
	Portuguese           Tag = Tag***REMOVED***lang: _pt***REMOVED***                //  pt
	BrazilianPortuguese  Tag = Tag***REMOVED***lang: _pt, region: _BR***REMOVED***   //  pt-BR
	EuropeanPortuguese   Tag = Tag***REMOVED***lang: _pt, region: _PT***REMOVED***   //  pt-PT
	Romanian             Tag = Tag***REMOVED***lang: _ro***REMOVED***                //  ro
	Russian              Tag = Tag***REMOVED***lang: _ru***REMOVED***                //  ru
	Sinhala              Tag = Tag***REMOVED***lang: _si***REMOVED***                //  si
	Slovak               Tag = Tag***REMOVED***lang: _sk***REMOVED***                //  sk
	Slovenian            Tag = Tag***REMOVED***lang: _sl***REMOVED***                //  sl
	Albanian             Tag = Tag***REMOVED***lang: _sq***REMOVED***                //  sq
	Serbian              Tag = Tag***REMOVED***lang: _sr***REMOVED***                //  sr
	SerbianLatin         Tag = Tag***REMOVED***lang: _sr, script: _Latn***REMOVED*** //  sr-Latn
	Swedish              Tag = Tag***REMOVED***lang: _sv***REMOVED***                //  sv
	Swahili              Tag = Tag***REMOVED***lang: _sw***REMOVED***                //  sw
	Tamil                Tag = Tag***REMOVED***lang: _ta***REMOVED***                //  ta
	Telugu               Tag = Tag***REMOVED***lang: _te***REMOVED***                //  te
	Thai                 Tag = Tag***REMOVED***lang: _th***REMOVED***                //  th
	Turkish              Tag = Tag***REMOVED***lang: _tr***REMOVED***                //  tr
	Ukrainian            Tag = Tag***REMOVED***lang: _uk***REMOVED***                //  uk
	Urdu                 Tag = Tag***REMOVED***lang: _ur***REMOVED***                //  ur
	Uzbek                Tag = Tag***REMOVED***lang: _uz***REMOVED***                //  uz
	Vietnamese           Tag = Tag***REMOVED***lang: _vi***REMOVED***                //  vi
	Chinese              Tag = Tag***REMOVED***lang: _zh***REMOVED***                //  zh
	SimplifiedChinese    Tag = Tag***REMOVED***lang: _zh, script: _Hans***REMOVED*** //  zh-Hans
	TraditionalChinese   Tag = Tag***REMOVED***lang: _zh, script: _Hant***REMOVED*** //  zh-Hant
	Zulu                 Tag = Tag***REMOVED***lang: _zu***REMOVED***                //  zu
)
