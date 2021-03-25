// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cases

// This file contains the definitions of case mappings for all supported
// languages. The rules for the language-specific tailorings were taken and
// modified from the CLDR transform definitions in common/transforms.

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/internal"
	"golang.org/x/text/language"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// A mapFunc takes a context set to the current rune and writes the mapped
// version to the same context. It may advance the context to the next rune. It
// returns whether a checkpoint is possible: whether the pDst bytes written to
// dst so far won't need changing as we see more source bytes.
type mapFunc func(*context) bool

// A spanFunc takes a context set to the current rune and returns whether this
// rune would be altered when written to the output. It may advance the context
// to the next rune. It returns whether a checkpoint is possible.
type spanFunc func(*context) bool

// maxIgnorable defines the maximum number of ignorables to consider for
// lookahead operations.
const maxIgnorable = 30

// supported lists the language tags for which we have tailorings.
const supported = "und af az el lt nl tr"

func init() ***REMOVED***
	tags := []language.Tag***REMOVED******REMOVED***
	for _, s := range strings.Split(supported, " ") ***REMOVED***
		tags = append(tags, language.MustParse(s))
	***REMOVED***
	matcher = internal.NewInheritanceMatcher(tags)
	Supported = language.NewCoverage(tags)
***REMOVED***

var (
	matcher *internal.InheritanceMatcher

	Supported language.Coverage

	// We keep the following lists separate, instead of having a single per-
	// language struct, to give the compiler a chance to remove unused code.

	// Some uppercase mappers are stateless, so we can precompute the
	// Transformers and save a bit on runtime allocations.
	upperFunc = []struct ***REMOVED***
		upper mapFunc
		span  spanFunc
	***REMOVED******REMOVED***
		***REMOVED***nil, nil***REMOVED***,                  // und
		***REMOVED***nil, nil***REMOVED***,                  // af
		***REMOVED***aztrUpper(upper), isUpper***REMOVED***, // az
		***REMOVED***elUpper, noSpan***REMOVED***,           // el
		***REMOVED***ltUpper(upper), noSpan***REMOVED***,    // lt
		***REMOVED***nil, nil***REMOVED***,                  // nl
		***REMOVED***aztrUpper(upper), isUpper***REMOVED***, // tr
	***REMOVED***

	undUpper            transform.SpanningTransformer = &undUpperCaser***REMOVED******REMOVED***
	undLower            transform.SpanningTransformer = &undLowerCaser***REMOVED******REMOVED***
	undLowerIgnoreSigma transform.SpanningTransformer = &undLowerIgnoreSigmaCaser***REMOVED******REMOVED***

	lowerFunc = []mapFunc***REMOVED***
		nil,       // und
		nil,       // af
		aztrLower, // az
		nil,       // el
		ltLower,   // lt
		nil,       // nl
		aztrLower, // tr
	***REMOVED***

	titleInfos = []struct ***REMOVED***
		title     mapFunc
		lower     mapFunc
		titleSpan spanFunc
		rewrite   func(*context)
	***REMOVED******REMOVED***
		***REMOVED***title, lower, isTitle, nil***REMOVED***,                // und
		***REMOVED***title, lower, isTitle, afnlRewrite***REMOVED***,        // af
		***REMOVED***aztrUpper(title), aztrLower, isTitle, nil***REMOVED***, // az
		***REMOVED***title, lower, isTitle, nil***REMOVED***,                // el
		***REMOVED***ltUpper(title), ltLower, noSpan, nil***REMOVED***,      // lt
		***REMOVED***nlTitle, lower, nlTitleSpan, afnlRewrite***REMOVED***,  // nl
		***REMOVED***aztrUpper(title), aztrLower, isTitle, nil***REMOVED***, // tr
	***REMOVED***
)

func makeUpper(t language.Tag, o options) transform.SpanningTransformer ***REMOVED***
	_, i, _ := matcher.Match(t)
	f := upperFunc[i].upper
	if f == nil ***REMOVED***
		return undUpper
	***REMOVED***
	return &simpleCaser***REMOVED***f: f, span: upperFunc[i].span***REMOVED***
***REMOVED***

func makeLower(t language.Tag, o options) transform.SpanningTransformer ***REMOVED***
	_, i, _ := matcher.Match(t)
	f := lowerFunc[i]
	if f == nil ***REMOVED***
		if o.ignoreFinalSigma ***REMOVED***
			return undLowerIgnoreSigma
		***REMOVED***
		return undLower
	***REMOVED***
	if o.ignoreFinalSigma ***REMOVED***
		return &simpleCaser***REMOVED***f: f, span: isLower***REMOVED***
	***REMOVED***
	return &lowerCaser***REMOVED***
		first:   f,
		midWord: finalSigma(f),
	***REMOVED***
***REMOVED***

func makeTitle(t language.Tag, o options) transform.SpanningTransformer ***REMOVED***
	_, i, _ := matcher.Match(t)
	x := &titleInfos[i]
	lower := x.lower
	if o.noLower ***REMOVED***
		lower = (*context).copy
	***REMOVED*** else if !o.ignoreFinalSigma ***REMOVED***
		lower = finalSigma(lower)
	***REMOVED***
	return &titleCaser***REMOVED***
		title:     x.title,
		lower:     lower,
		titleSpan: x.titleSpan,
		rewrite:   x.rewrite,
	***REMOVED***
***REMOVED***

func noSpan(c *context) bool ***REMOVED***
	c.err = transform.ErrEndOfSpan
	return false
***REMOVED***

// TODO: consider a similar special case for the fast majority lower case. This
// is a bit more involved so will require some more precise benchmarking to
// justify it.

type undUpperCaser struct***REMOVED*** transform.NopResetter ***REMOVED***

// undUpperCaser implements the Transformer interface for doing an upper case
// mapping for the root locale (und). It eliminates the need for an allocation
// as it prevents escaping by not using function pointers.
func (t undUpperCaser) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	c := context***REMOVED***dst: dst, src: src, atEOF: atEOF***REMOVED***
	for c.next() ***REMOVED***
		upper(&c)
		c.checkpoint()
	***REMOVED***
	return c.ret()
***REMOVED***

func (t undUpperCaser) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	c := context***REMOVED***src: src, atEOF: atEOF***REMOVED***
	for c.next() && isUpper(&c) ***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.retSpan()
***REMOVED***

// undLowerIgnoreSigmaCaser implements the Transformer interface for doing
// a lower case mapping for the root locale (und) ignoring final sigma
// handling. This casing algorithm is used in some performance-critical packages
// like secure/precis and x/net/http/idna, which warrants its special-casing.
type undLowerIgnoreSigmaCaser struct***REMOVED*** transform.NopResetter ***REMOVED***

func (t undLowerIgnoreSigmaCaser) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	c := context***REMOVED***dst: dst, src: src, atEOF: atEOF***REMOVED***
	for c.next() && lower(&c) ***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.ret()

***REMOVED***

// Span implements a generic lower-casing. This is possible as isLower works
// for all lowercasing variants. All lowercase variants only vary in how they
// transform a non-lowercase letter. They will never change an already lowercase
// letter. In addition, there is no state.
func (t undLowerIgnoreSigmaCaser) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	c := context***REMOVED***src: src, atEOF: atEOF***REMOVED***
	for c.next() && isLower(&c) ***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.retSpan()
***REMOVED***

type simpleCaser struct ***REMOVED***
	context
	f    mapFunc
	span spanFunc
***REMOVED***

// simpleCaser implements the Transformer interface for doing a case operation
// on a rune-by-rune basis.
func (t *simpleCaser) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	c := context***REMOVED***dst: dst, src: src, atEOF: atEOF***REMOVED***
	for c.next() && t.f(&c) ***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.ret()
***REMOVED***

func (t *simpleCaser) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	c := context***REMOVED***src: src, atEOF: atEOF***REMOVED***
	for c.next() && t.span(&c) ***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.retSpan()
***REMOVED***

// undLowerCaser implements the Transformer interface for doing a lower case
// mapping for the root locale (und) ignoring final sigma handling. This casing
// algorithm is used in some performance-critical packages like secure/precis
// and x/net/http/idna, which warrants its special-casing.
type undLowerCaser struct***REMOVED*** transform.NopResetter ***REMOVED***

func (t undLowerCaser) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	c := context***REMOVED***dst: dst, src: src, atEOF: atEOF***REMOVED***

	for isInterWord := true; c.next(); ***REMOVED***
		if isInterWord ***REMOVED***
			if c.info.isCased() ***REMOVED***
				if !lower(&c) ***REMOVED***
					break
				***REMOVED***
				isInterWord = false
			***REMOVED*** else if !c.copy() ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if c.info.isNotCasedAndNotCaseIgnorable() ***REMOVED***
				if !c.copy() ***REMOVED***
					break
				***REMOVED***
				isInterWord = true
			***REMOVED*** else if !c.hasPrefix("Σ") ***REMOVED***
				if !lower(&c) ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else if !finalSigmaBody(&c) ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.ret()
***REMOVED***

func (t undLowerCaser) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	c := context***REMOVED***src: src, atEOF: atEOF***REMOVED***
	for c.next() && isLower(&c) ***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.retSpan()
***REMOVED***

// lowerCaser implements the Transformer interface. The default Unicode lower
// casing requires different treatment for the first and subsequent characters
// of a word, most notably to handle the Greek final Sigma.
type lowerCaser struct ***REMOVED***
	undLowerIgnoreSigmaCaser

	context

	first, midWord mapFunc
***REMOVED***

func (t *lowerCaser) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	t.context = context***REMOVED***dst: dst, src: src, atEOF: atEOF***REMOVED***
	c := &t.context

	for isInterWord := true; c.next(); ***REMOVED***
		if isInterWord ***REMOVED***
			if c.info.isCased() ***REMOVED***
				if !t.first(c) ***REMOVED***
					break
				***REMOVED***
				isInterWord = false
			***REMOVED*** else if !c.copy() ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if c.info.isNotCasedAndNotCaseIgnorable() ***REMOVED***
				if !c.copy() ***REMOVED***
					break
				***REMOVED***
				isInterWord = true
			***REMOVED*** else if !t.midWord(c) ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.ret()
***REMOVED***

// titleCaser implements the Transformer interface. Title casing algorithms
// distinguish between the first letter of a word and subsequent letters of the
// same word. It uses state to avoid requiring a potentially infinite lookahead.
type titleCaser struct ***REMOVED***
	context

	// rune mappings used by the actual casing algorithms.
	title     mapFunc
	lower     mapFunc
	titleSpan spanFunc

	rewrite func(*context)
***REMOVED***

// Transform implements the standard Unicode title case algorithm as defined in
// Chapter 3 of The Unicode Standard:
// toTitlecase(X): Find the word boundaries in X according to Unicode Standard
// Annex #29, "Unicode Text Segmentation." For each word boundary, find the
// first cased character F following the word boundary. If F exists, map F to
// Titlecase_Mapping(F); then map all characters C between F and the following
// word boundary to Lowercase_Mapping(C).
func (t *titleCaser) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	t.context = context***REMOVED***dst: dst, src: src, atEOF: atEOF, isMidWord: t.isMidWord***REMOVED***
	c := &t.context

	if !c.next() ***REMOVED***
		return c.ret()
	***REMOVED***

	for ***REMOVED***
		p := c.info
		if t.rewrite != nil ***REMOVED***
			t.rewrite(c)
		***REMOVED***

		wasMid := p.isMid()
		// Break out of this loop on failure to ensure we do not modify the
		// state incorrectly.
		if p.isCased() ***REMOVED***
			if !c.isMidWord ***REMOVED***
				if !t.title(c) ***REMOVED***
					break
				***REMOVED***
				c.isMidWord = true
			***REMOVED*** else if !t.lower(c) ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else if !c.copy() ***REMOVED***
			break
		***REMOVED*** else if p.isBreak() ***REMOVED***
			c.isMidWord = false
		***REMOVED***

		// As we save the state of the transformer, it is safe to call
		// checkpoint after any successful write.
		if !(c.isMidWord && wasMid) ***REMOVED***
			c.checkpoint()
		***REMOVED***

		if !c.next() ***REMOVED***
			break
		***REMOVED***
		if wasMid && c.info.isMid() ***REMOVED***
			c.isMidWord = false
		***REMOVED***
	***REMOVED***
	return c.ret()
***REMOVED***

func (t *titleCaser) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	t.context = context***REMOVED***src: src, atEOF: atEOF, isMidWord: t.isMidWord***REMOVED***
	c := &t.context

	if !c.next() ***REMOVED***
		return c.retSpan()
	***REMOVED***

	for ***REMOVED***
		p := c.info
		if t.rewrite != nil ***REMOVED***
			t.rewrite(c)
		***REMOVED***

		wasMid := p.isMid()
		// Break out of this loop on failure to ensure we do not modify the
		// state incorrectly.
		if p.isCased() ***REMOVED***
			if !c.isMidWord ***REMOVED***
				if !t.titleSpan(c) ***REMOVED***
					break
				***REMOVED***
				c.isMidWord = true
			***REMOVED*** else if !isLower(c) ***REMOVED***
				break
			***REMOVED***
		***REMOVED*** else if p.isBreak() ***REMOVED***
			c.isMidWord = false
		***REMOVED***
		// As we save the state of the transformer, it is safe to call
		// checkpoint after any successful write.
		if !(c.isMidWord && wasMid) ***REMOVED***
			c.checkpoint()
		***REMOVED***

		if !c.next() ***REMOVED***
			break
		***REMOVED***
		if wasMid && c.info.isMid() ***REMOVED***
			c.isMidWord = false
		***REMOVED***
	***REMOVED***
	return c.retSpan()
***REMOVED***

// finalSigma adds Greek final Sigma handing to another casing function. It
// determines whether a lowercased sigma should be σ or ς, by looking ahead for
// case-ignorables and a cased letters.
func finalSigma(f mapFunc) mapFunc ***REMOVED***
	return func(c *context) bool ***REMOVED***
		if !c.hasPrefix("Σ") ***REMOVED***
			return f(c)
		***REMOVED***
		return finalSigmaBody(c)
	***REMOVED***
***REMOVED***

func finalSigmaBody(c *context) bool ***REMOVED***
	// Current rune must be ∑.

	// ::NFD();
	// # 03A3; 03C2; 03A3; 03A3; Final_Sigma; # GREEK CAPITAL LETTER SIGMA
	// Σ ***REMOVED*** [:case-ignorable:]* [:cased:] → σ;
	// [:cased:] [:case-ignorable:]* ***REMOVED*** Σ → ς;
	// ::Any-Lower;
	// ::NFC();

	p := c.pDst
	c.writeString("ς")

	// TODO: we should do this here, but right now this will never have an
	// effect as this is called when the prefix is Sigma, whereas Dutch and
	// Afrikaans only test for an apostrophe.
	//
	// if t.rewrite != nil ***REMOVED***
	// 	t.rewrite(c)
	// ***REMOVED***

	// We need to do one more iteration after maxIgnorable, as a cased
	// letter is not an ignorable and may modify the result.
	wasMid := false
	for i := 0; i < maxIgnorable+1; i++ ***REMOVED***
		if !c.next() ***REMOVED***
			return false
		***REMOVED***
		if !c.info.isCaseIgnorable() ***REMOVED***
			// All Midword runes are also case ignorable, so we are
			// guaranteed to have a letter or word break here. As we are
			// unreading the run, there is no need to unset c.isMidWord;
			// the title caser will handle this.
			if c.info.isCased() ***REMOVED***
				// p+1 is guaranteed to be in bounds: if writing ς was
				// successful, p+1 will contain the second byte of ς. If not,
				// this function will have returned after c.next returned false.
				c.dst[p+1]++ // ς → σ
			***REMOVED***
			c.unreadRune()
			return true
		***REMOVED***
		// A case ignorable may also introduce a word break, so we may need
		// to continue searching even after detecting a break.
		isMid := c.info.isMid()
		if (wasMid && isMid) || c.info.isBreak() ***REMOVED***
			c.isMidWord = false
		***REMOVED***
		wasMid = isMid
		c.copy()
	***REMOVED***
	return true
***REMOVED***

// finalSigmaSpan would be the same as isLower.

// elUpper implements Greek upper casing, which entails removing a predefined
// set of non-blocked modifiers. Note that these accents should not be removed
// for title casing!
// Example: "Οδός" -> "ΟΔΟΣ".
func elUpper(c *context) bool ***REMOVED***
	// From CLDR:
	// [:Greek:] [^[:ccc=Not_Reordered:][:ccc=Above:]]*? ***REMOVED*** [\u0313\u0314\u0301\u0300\u0306\u0342\u0308\u0304] → ;
	// [:Greek:] [^[:ccc=Not_Reordered:][:ccc=Iota_Subscript:]]*? ***REMOVED*** \u0345 → ;

	r, _ := utf8.DecodeRune(c.src[c.pSrc:])
	oldPDst := c.pDst
	if !upper(c) ***REMOVED***
		return false
	***REMOVED***
	if !unicode.Is(unicode.Greek, r) ***REMOVED***
		return true
	***REMOVED***
	i := 0
	// Take the properties of the uppercased rune that is already written to the
	// destination. This saves us the trouble of having to uppercase the
	// decomposed rune again.
	if b := norm.NFD.Properties(c.dst[oldPDst:]).Decomposition(); b != nil ***REMOVED***
		// Restore the destination position and process the decomposed rune.
		r, sz := utf8.DecodeRune(b)
		if r <= 0xFF ***REMOVED*** // See A.6.1
			return true
		***REMOVED***
		c.pDst = oldPDst
		// Insert the first rune and ignore the modifiers. See A.6.2.
		c.writeBytes(b[:sz])
		i = len(b[sz:]) / 2 // Greek modifiers are always of length 2.
	***REMOVED***

	for ; i < maxIgnorable && c.next(); i++ ***REMOVED***
		switch r, _ := utf8.DecodeRune(c.src[c.pSrc:]); r ***REMOVED***
		// Above and Iota Subscript
		case 0x0300, // U+0300 COMBINING GRAVE ACCENT
			0x0301, // U+0301 COMBINING ACUTE ACCENT
			0x0304, // U+0304 COMBINING MACRON
			0x0306, // U+0306 COMBINING BREVE
			0x0308, // U+0308 COMBINING DIAERESIS
			0x0313, // U+0313 COMBINING COMMA ABOVE
			0x0314, // U+0314 COMBINING REVERSED COMMA ABOVE
			0x0342, // U+0342 COMBINING GREEK PERISPOMENI
			0x0345: // U+0345 COMBINING GREEK YPOGEGRAMMENI
			// No-op. Gobble the modifier.

		default:
			switch v, _ := trie.lookup(c.src[c.pSrc:]); info(v).cccType() ***REMOVED***
			case cccZero:
				c.unreadRune()
				return true

			// We don't need to test for IotaSubscript as the only rune that
			// qualifies (U+0345) was already excluded in the switch statement
			// above. See A.4.

			case cccAbove:
				return c.copy()
			default:
				// Some other modifier. We're still allowed to gobble Greek
				// modifiers after this.
				c.copy()
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return i == maxIgnorable
***REMOVED***

// TODO: implement elUpperSpan (low-priority: complex and infrequent).

func ltLower(c *context) bool ***REMOVED***
	// From CLDR:
	// # Introduce an explicit dot above when lowercasing capital I's and J's
	// # whenever there are more accents above.
	// # (of the accents used in Lithuanian: grave, acute, tilde above, and ogonek)
	// # 0049; 0069 0307; 0049; 0049; lt More_Above; # LATIN CAPITAL LETTER I
	// # 004A; 006A 0307; 004A; 004A; lt More_Above; # LATIN CAPITAL LETTER J
	// # 012E; 012F 0307; 012E; 012E; lt More_Above; # LATIN CAPITAL LETTER I WITH OGONEK
	// # 00CC; 0069 0307 0300; 00CC; 00CC; lt; # LATIN CAPITAL LETTER I WITH GRAVE
	// # 00CD; 0069 0307 0301; 00CD; 00CD; lt; # LATIN CAPITAL LETTER I WITH ACUTE
	// # 0128; 0069 0307 0303; 0128; 0128; lt; # LATIN CAPITAL LETTER I WITH TILDE
	// ::NFD();
	// I ***REMOVED*** [^[:ccc=Not_Reordered:][:ccc=Above:]]* [:ccc=Above:] → i \u0307;
	// J ***REMOVED*** [^[:ccc=Not_Reordered:][:ccc=Above:]]* [:ccc=Above:] → j \u0307;
	// I \u0328 (Į) ***REMOVED*** [^[:ccc=Not_Reordered:][:ccc=Above:]]* [:ccc=Above:] → i \u0328 \u0307;
	// I \u0300 (Ì) → i \u0307 \u0300;
	// I \u0301 (Í) → i \u0307 \u0301;
	// I \u0303 (Ĩ) → i \u0307 \u0303;
	// ::Any-Lower();
	// ::NFC();

	i := 0
	if r := c.src[c.pSrc]; r < utf8.RuneSelf ***REMOVED***
		lower(c)
		if r != 'I' && r != 'J' ***REMOVED***
			return true
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		p := norm.NFD.Properties(c.src[c.pSrc:])
		if d := p.Decomposition(); len(d) >= 3 && (d[0] == 'I' || d[0] == 'J') ***REMOVED***
			// UTF-8 optimization: the decomposition will only have an above
			// modifier if the last rune of the decomposition is in [U+300-U+311].
			// In all other cases, a decomposition starting with I is always
			// an I followed by modifiers that are not cased themselves. See A.2.
			if d[1] == 0xCC && d[2] <= 0x91 ***REMOVED*** // A.2.4.
				if !c.writeBytes(d[:1]) ***REMOVED***
					return false
				***REMOVED***
				c.dst[c.pDst-1] += 'a' - 'A' // lower

				// Assumption: modifier never changes on lowercase. See A.1.
				// Assumption: all modifiers added have CCC = Above. See A.2.3.
				return c.writeString("\u0307") && c.writeBytes(d[1:])
			***REMOVED***
			// In all other cases the additional modifiers will have a CCC
			// that is less than 230 (Above). We will insert the U+0307, if
			// needed, after these modifiers so that a string in FCD form
			// will remain so. See A.2.2.
			lower(c)
			i = 1
		***REMOVED*** else ***REMOVED***
			return lower(c)
		***REMOVED***
	***REMOVED***

	for ; i < maxIgnorable && c.next(); i++ ***REMOVED***
		switch c.info.cccType() ***REMOVED***
		case cccZero:
			c.unreadRune()
			return true
		case cccAbove:
			return c.writeString("\u0307") && c.copy() // See A.1.
		default:
			c.copy() // See A.1.
		***REMOVED***
	***REMOVED***
	return i == maxIgnorable
***REMOVED***

// ltLowerSpan would be the same as isLower.

func ltUpper(f mapFunc) mapFunc ***REMOVED***
	return func(c *context) bool ***REMOVED***
		// Unicode:
		// 0307; 0307; ; ; lt After_Soft_Dotted; # COMBINING DOT ABOVE
		//
		// From CLDR:
		// # Remove \u0307 following soft-dotteds (i, j, and the like), with possible
		// # intervening non-230 marks.
		// ::NFD();
		// [:Soft_Dotted:] [^[:ccc=Not_Reordered:][:ccc=Above:]]* ***REMOVED*** \u0307 → ;
		// ::Any-Upper();
		// ::NFC();

		// TODO: See A.5. A soft-dotted rune never has an exception. This would
		// allow us to overload the exception bit and encode this property in
		// info. Need to measure performance impact of this.
		r, _ := utf8.DecodeRune(c.src[c.pSrc:])
		oldPDst := c.pDst
		if !f(c) ***REMOVED***
			return false
		***REMOVED***
		if !unicode.Is(unicode.Soft_Dotted, r) ***REMOVED***
			return true
		***REMOVED***

		// We don't need to do an NFD normalization, as a soft-dotted rune never
		// contains U+0307. See A.3.

		i := 0
		for ; i < maxIgnorable && c.next(); i++ ***REMOVED***
			switch c.info.cccType() ***REMOVED***
			case cccZero:
				c.unreadRune()
				return true
			case cccAbove:
				if c.hasPrefix("\u0307") ***REMOVED***
					// We don't do a full NFC, but rather combine runes for
					// some of the common cases. (Returning NFC or
					// preserving normal form is neither a requirement nor
					// a possibility anyway).
					if !c.next() ***REMOVED***
						return false
					***REMOVED***
					if c.dst[oldPDst] == 'I' && c.pDst == oldPDst+1 && c.src[c.pSrc] == 0xcc ***REMOVED***
						s := ""
						switch c.src[c.pSrc+1] ***REMOVED***
						case 0x80: // U+0300 COMBINING GRAVE ACCENT
							s = "\u00cc" // U+00CC LATIN CAPITAL LETTER I WITH GRAVE
						case 0x81: // U+0301 COMBINING ACUTE ACCENT
							s = "\u00cd" // U+00CD LATIN CAPITAL LETTER I WITH ACUTE
						case 0x83: // U+0303 COMBINING TILDE
							s = "\u0128" // U+0128 LATIN CAPITAL LETTER I WITH TILDE
						case 0x88: // U+0308 COMBINING DIAERESIS
							s = "\u00cf" // U+00CF LATIN CAPITAL LETTER I WITH DIAERESIS
						default:
						***REMOVED***
						if s != "" ***REMOVED***
							c.pDst = oldPDst
							return c.writeString(s)
						***REMOVED***
					***REMOVED***
				***REMOVED***
				return c.copy()
			default:
				c.copy()
			***REMOVED***
		***REMOVED***
		return i == maxIgnorable
	***REMOVED***
***REMOVED***

// TODO: implement ltUpperSpan (low priority: complex and infrequent).

func aztrUpper(f mapFunc) mapFunc ***REMOVED***
	return func(c *context) bool ***REMOVED***
		// i→İ;
		if c.src[c.pSrc] == 'i' ***REMOVED***
			return c.writeString("İ")
		***REMOVED***
		return f(c)
	***REMOVED***
***REMOVED***

func aztrLower(c *context) (done bool) ***REMOVED***
	// From CLDR:
	// # I and i-dotless; I-dot and i are case pairs in Turkish and Azeri
	// # 0130; 0069; 0130; 0130; tr; # LATIN CAPITAL LETTER I WITH DOT ABOVE
	// İ→i;
	// # When lowercasing, remove dot_above in the sequence I + dot_above, which will turn into i.
	// # This matches the behavior of the canonically equivalent I-dot_above
	// # 0307; ; 0307; 0307; tr After_I; # COMBINING DOT ABOVE
	// # When lowercasing, unless an I is before a dot_above, it turns into a dotless i.
	// # 0049; 0131; 0049; 0049; tr Not_Before_Dot; # LATIN CAPITAL LETTER I
	// I([^[:ccc=Not_Reordered:][:ccc=Above:]]*)\u0307 → i$1 ;
	// I→ı ;
	// ::Any-Lower();
	if c.hasPrefix("\u0130") ***REMOVED*** // İ
		return c.writeString("i")
	***REMOVED***
	if c.src[c.pSrc] != 'I' ***REMOVED***
		return lower(c)
	***REMOVED***

	// We ignore the lower-case I for now, but insert it later when we know
	// which form we need.
	start := c.pSrc + c.sz

	i := 0
Loop:
	// We check for up to n ignorables before \u0307. As \u0307 is an
	// ignorable as well, n is maxIgnorable-1.
	for ; i < maxIgnorable && c.next(); i++ ***REMOVED***
		switch c.info.cccType() ***REMOVED***
		case cccAbove:
			if c.hasPrefix("\u0307") ***REMOVED***
				return c.writeString("i") && c.writeBytes(c.src[start:c.pSrc]) // ignore U+0307
			***REMOVED***
			done = true
			break Loop
		case cccZero:
			c.unreadRune()
			done = true
			break Loop
		default:
			// We'll write this rune after we know which starter to use.
		***REMOVED***
	***REMOVED***
	if i == maxIgnorable ***REMOVED***
		done = true
	***REMOVED***
	return c.writeString("ı") && c.writeBytes(c.src[start:c.pSrc+c.sz]) && done
***REMOVED***

// aztrLowerSpan would be the same as isLower.

func nlTitle(c *context) bool ***REMOVED***
	// From CLDR:
	// # Special titlecasing for Dutch initial "ij".
	// ::Any-Title();
	// # Fix up Ij at the beginning of a "word" (per Any-Title, notUAX #29)
	// [:^WB=ALetter:] [:WB=Extend:]* [[:WB=MidLetter:][:WB=MidNumLet:]]? ***REMOVED*** Ij ***REMOVED*** → IJ ;
	if c.src[c.pSrc] != 'I' && c.src[c.pSrc] != 'i' ***REMOVED***
		return title(c)
	***REMOVED***

	if !c.writeString("I") || !c.next() ***REMOVED***
		return false
	***REMOVED***
	if c.src[c.pSrc] == 'j' || c.src[c.pSrc] == 'J' ***REMOVED***
		return c.writeString("J")
	***REMOVED***
	c.unreadRune()
	return true
***REMOVED***

func nlTitleSpan(c *context) bool ***REMOVED***
	// From CLDR:
	// # Special titlecasing for Dutch initial "ij".
	// ::Any-Title();
	// # Fix up Ij at the beginning of a "word" (per Any-Title, notUAX #29)
	// [:^WB=ALetter:] [:WB=Extend:]* [[:WB=MidLetter:][:WB=MidNumLet:]]? ***REMOVED*** Ij ***REMOVED*** → IJ ;
	if c.src[c.pSrc] != 'I' ***REMOVED***
		return isTitle(c)
	***REMOVED***
	if !c.next() || c.src[c.pSrc] == 'j' ***REMOVED***
		return false
	***REMOVED***
	if c.src[c.pSrc] != 'J' ***REMOVED***
		c.unreadRune()
	***REMOVED***
	return true
***REMOVED***

// Not part of CLDR, but see https://unicode.org/cldr/trac/ticket/7078.
func afnlRewrite(c *context) ***REMOVED***
	if c.hasPrefix("'") || c.hasPrefix("’") ***REMOVED***
		c.isMidWord = true
	***REMOVED***
***REMOVED***
