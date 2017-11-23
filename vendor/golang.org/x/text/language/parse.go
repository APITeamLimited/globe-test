// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/internal/tag"
)

// isAlpha returns true if the byte is not a digit.
// b must be an ASCII letter or digit.
func isAlpha(b byte) bool ***REMOVED***
	return b > '9'
***REMOVED***

// isAlphaNum returns true if the string contains only ASCII letters or digits.
func isAlphaNum(s []byte) bool ***REMOVED***
	for _, c := range s ***REMOVED***
		if !('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9') ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// errSyntax is returned by any of the parsing functions when the
// input is not well-formed, according to BCP 47.
// TODO: return the position at which the syntax error occurred?
var errSyntax = errors.New("language: tag is not well-formed")

// ValueError is returned by any of the parsing functions when the
// input is well-formed but the respective subtag is not recognized
// as a valid value.
type ValueError struct ***REMOVED***
	v [8]byte
***REMOVED***

func mkErrInvalid(s []byte) error ***REMOVED***
	var e ValueError
	copy(e.v[:], s)
	return e
***REMOVED***

func (e ValueError) tag() []byte ***REMOVED***
	n := bytes.IndexByte(e.v[:], 0)
	if n == -1 ***REMOVED***
		n = 8
	***REMOVED***
	return e.v[:n]
***REMOVED***

// Error implements the error interface.
func (e ValueError) Error() string ***REMOVED***
	return fmt.Sprintf("language: subtag %q is well-formed but unknown", e.tag())
***REMOVED***

// Subtag returns the subtag for which the error occurred.
func (e ValueError) Subtag() string ***REMOVED***
	return string(e.tag())
***REMOVED***

// scanner is used to scan BCP 47 tokens, which are separated by _ or -.
type scanner struct ***REMOVED***
	b     []byte
	bytes [max99thPercentileSize]byte
	token []byte
	start int // start position of the current token
	end   int // end position of the current token
	next  int // next point for scan
	err   error
	done  bool
***REMOVED***

func makeScannerString(s string) scanner ***REMOVED***
	scan := scanner***REMOVED******REMOVED***
	if len(s) <= len(scan.bytes) ***REMOVED***
		scan.b = scan.bytes[:copy(scan.bytes[:], s)]
	***REMOVED*** else ***REMOVED***
		scan.b = []byte(s)
	***REMOVED***
	scan.init()
	return scan
***REMOVED***

// makeScanner returns a scanner using b as the input buffer.
// b is not copied and may be modified by the scanner routines.
func makeScanner(b []byte) scanner ***REMOVED***
	scan := scanner***REMOVED***b: b***REMOVED***
	scan.init()
	return scan
***REMOVED***

func (s *scanner) init() ***REMOVED***
	for i, c := range s.b ***REMOVED***
		if c == '_' ***REMOVED***
			s.b[i] = '-'
		***REMOVED***
	***REMOVED***
	s.scan()
***REMOVED***

// restToLower converts the string between start and end to lower case.
func (s *scanner) toLower(start, end int) ***REMOVED***
	for i := start; i < end; i++ ***REMOVED***
		c := s.b[i]
		if 'A' <= c && c <= 'Z' ***REMOVED***
			s.b[i] += 'a' - 'A'
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *scanner) setError(e error) ***REMOVED***
	if s.err == nil || (e == errSyntax && s.err != errSyntax) ***REMOVED***
		s.err = e
	***REMOVED***
***REMOVED***

// resizeRange shrinks or grows the array at position oldStart such that
// a new string of size newSize can fit between oldStart and oldEnd.
// Sets the scan point to after the resized range.
func (s *scanner) resizeRange(oldStart, oldEnd, newSize int) ***REMOVED***
	s.start = oldStart
	if end := oldStart + newSize; end != oldEnd ***REMOVED***
		diff := end - oldEnd
		if end < cap(s.b) ***REMOVED***
			b := make([]byte, len(s.b)+diff)
			copy(b, s.b[:oldStart])
			copy(b[end:], s.b[oldEnd:])
			s.b = b
		***REMOVED*** else ***REMOVED***
			s.b = append(s.b[end:], s.b[oldEnd:]...)
		***REMOVED***
		s.next = end + (s.next - s.end)
		s.end = end
	***REMOVED***
***REMOVED***

// replace replaces the current token with repl.
func (s *scanner) replace(repl string) ***REMOVED***
	s.resizeRange(s.start, s.end, len(repl))
	copy(s.b[s.start:], repl)
***REMOVED***

// gobble removes the current token from the input.
// Caller must call scan after calling gobble.
func (s *scanner) gobble(e error) ***REMOVED***
	s.setError(e)
	if s.start == 0 ***REMOVED***
		s.b = s.b[:+copy(s.b, s.b[s.next:])]
		s.end = 0
	***REMOVED*** else ***REMOVED***
		s.b = s.b[:s.start-1+copy(s.b[s.start-1:], s.b[s.end:])]
		s.end = s.start - 1
	***REMOVED***
	s.next = s.start
***REMOVED***

// deleteRange removes the given range from s.b before the current token.
func (s *scanner) deleteRange(start, end int) ***REMOVED***
	s.setError(errSyntax)
	s.b = s.b[:start+copy(s.b[start:], s.b[end:])]
	diff := end - start
	s.next -= diff
	s.start -= diff
	s.end -= diff
***REMOVED***

// scan parses the next token of a BCP 47 string.  Tokens that are larger
// than 8 characters or include non-alphanumeric characters result in an error
// and are gobbled and removed from the output.
// It returns the end position of the last token consumed.
func (s *scanner) scan() (end int) ***REMOVED***
	end = s.end
	s.token = nil
	for s.start = s.next; s.next < len(s.b); ***REMOVED***
		i := bytes.IndexByte(s.b[s.next:], '-')
		if i == -1 ***REMOVED***
			s.end = len(s.b)
			s.next = len(s.b)
			i = s.end - s.start
		***REMOVED*** else ***REMOVED***
			s.end = s.next + i
			s.next = s.end + 1
		***REMOVED***
		token := s.b[s.start:s.end]
		if i < 1 || i > 8 || !isAlphaNum(token) ***REMOVED***
			s.gobble(errSyntax)
			continue
		***REMOVED***
		s.token = token
		return end
	***REMOVED***
	if n := len(s.b); n > 0 && s.b[n-1] == '-' ***REMOVED***
		s.setError(errSyntax)
		s.b = s.b[:len(s.b)-1]
	***REMOVED***
	s.done = true
	return end
***REMOVED***

// acceptMinSize parses multiple tokens of the given size or greater.
// It returns the end position of the last token consumed.
func (s *scanner) acceptMinSize(min int) (end int) ***REMOVED***
	end = s.end
	s.scan()
	for ; len(s.token) >= min; s.scan() ***REMOVED***
		end = s.end
	***REMOVED***
	return end
***REMOVED***

// Parse parses the given BCP 47 string and returns a valid Tag. If parsing
// failed it returns an error and any part of the tag that could be parsed.
// If parsing succeeded but an unknown value was found, it returns
// ValueError. The Tag returned in this case is just stripped of the unknown
// value. All other values are preserved. It accepts tags in the BCP 47 format
// and extensions to this standard defined in
// http://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
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
// http://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
// The resulting tag is canonicalized using the the canonicalization type c.
func (c CanonType) Parse(s string) (t Tag, err error) ***REMOVED***
	// TODO: consider supporting old-style locale key-value pairs.
	if s == "" ***REMOVED***
		return und, errSyntax
	***REMOVED***
	if len(s) <= maxAltTaglen ***REMOVED***
		b := [maxAltTaglen]byte***REMOVED******REMOVED***
		for i, c := range s ***REMOVED***
			// Generating invalid UTF-8 is okay as it won't match.
			if 'A' <= c && c <= 'Z' ***REMOVED***
				c += 'a' - 'A'
			***REMOVED*** else if c == '_' ***REMOVED***
				c = '-'
			***REMOVED***
			b[i] = byte(c)
		***REMOVED***
		if t, ok := grandfathered(b); ok ***REMOVED***
			return t, nil
		***REMOVED***
	***REMOVED***
	scan := makeScannerString(s)
	t, err = parse(&scan, s)
	t, changed := t.canonicalize(c)
	if changed ***REMOVED***
		t.remakeString()
	***REMOVED***
	return t, err
***REMOVED***

func parse(scan *scanner, s string) (t Tag, err error) ***REMOVED***
	t = und
	var end int
	if n := len(scan.token); n <= 1 ***REMOVED***
		scan.toLower(0, len(scan.b))
		if n == 0 || scan.token[0] != 'x' ***REMOVED***
			return t, errSyntax
		***REMOVED***
		end = parseExtensions(scan)
	***REMOVED*** else if n >= 4 ***REMOVED***
		return und, errSyntax
	***REMOVED*** else ***REMOVED*** // the usual case
		t, end = parseTag(scan)
		if n := len(scan.token); n == 1 ***REMOVED***
			t.pExt = uint16(end)
			end = parseExtensions(scan)
		***REMOVED*** else if end < len(scan.b) ***REMOVED***
			scan.setError(errSyntax)
			scan.b = scan.b[:end]
		***REMOVED***
	***REMOVED***
	if int(t.pVariant) < len(scan.b) ***REMOVED***
		if end < len(s) ***REMOVED***
			s = s[:end]
		***REMOVED***
		if len(s) > 0 && tag.Compare(s, scan.b) == 0 ***REMOVED***
			t.str = s
		***REMOVED*** else ***REMOVED***
			t.str = string(scan.b)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		t.pVariant, t.pExt = 0, 0
	***REMOVED***
	return t, scan.err
***REMOVED***

// parseTag parses language, script, region and variants.
// It returns a Tag and the end position in the input that was parsed.
func parseTag(scan *scanner) (t Tag, end int) ***REMOVED***
	var e error
	// TODO: set an error if an unknown lang, script or region is encountered.
	t.lang, e = getLangID(scan.token)
	scan.setError(e)
	scan.replace(t.lang.String())
	langStart := scan.start
	end = scan.scan()
	for len(scan.token) == 3 && isAlpha(scan.token[0]) ***REMOVED***
		// From http://tools.ietf.org/html/bcp47, <lang>-<extlang> tags are equivalent
		// to a tag of the form <extlang>.
		lang, e := getLangID(scan.token)
		if lang != 0 ***REMOVED***
			t.lang = lang
			copy(scan.b[langStart:], lang.String())
			scan.b[langStart+3] = '-'
			scan.start = langStart + 4
		***REMOVED***
		scan.gobble(e)
		end = scan.scan()
	***REMOVED***
	if len(scan.token) == 4 && isAlpha(scan.token[0]) ***REMOVED***
		t.script, e = getScriptID(script, scan.token)
		if t.script == 0 ***REMOVED***
			scan.gobble(e)
		***REMOVED***
		end = scan.scan()
	***REMOVED***
	if n := len(scan.token); n >= 2 && n <= 3 ***REMOVED***
		t.region, e = getRegionID(scan.token)
		if t.region == 0 ***REMOVED***
			scan.gobble(e)
		***REMOVED*** else ***REMOVED***
			scan.replace(t.region.String())
		***REMOVED***
		end = scan.scan()
	***REMOVED***
	scan.toLower(scan.start, len(scan.b))
	t.pVariant = byte(end)
	end = parseVariants(scan, end, t)
	t.pExt = uint16(end)
	return t, end
***REMOVED***

var separator = []byte***REMOVED***'-'***REMOVED***

// parseVariants scans tokens as long as each token is a valid variant string.
// Duplicate variants are removed.
func parseVariants(scan *scanner, end int, t Tag) int ***REMOVED***
	start := scan.start
	varIDBuf := [4]uint8***REMOVED******REMOVED***
	variantBuf := [4][]byte***REMOVED******REMOVED***
	varID := varIDBuf[:0]
	variant := variantBuf[:0]
	last := -1
	needSort := false
	for ; len(scan.token) >= 4; scan.scan() ***REMOVED***
		// TODO: measure the impact of needing this conversion and redesign
		// the data structure if there is an issue.
		v, ok := variantIndex[string(scan.token)]
		if !ok ***REMOVED***
			// unknown variant
			// TODO: allow user-defined variants?
			scan.gobble(mkErrInvalid(scan.token))
			continue
		***REMOVED***
		varID = append(varID, v)
		variant = append(variant, scan.token)
		if !needSort ***REMOVED***
			if last < int(v) ***REMOVED***
				last = int(v)
			***REMOVED*** else ***REMOVED***
				needSort = true
				// There is no legal combinations of more than 7 variants
				// (and this is by no means a useful sequence).
				const maxVariants = 8
				if len(varID) > maxVariants ***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
		end = scan.end
	***REMOVED***
	if needSort ***REMOVED***
		sort.Sort(variantsSort***REMOVED***varID, variant***REMOVED***)
		k, l := 0, -1
		for i, v := range varID ***REMOVED***
			w := int(v)
			if l == w ***REMOVED***
				// Remove duplicates.
				continue
			***REMOVED***
			varID[k] = varID[i]
			variant[k] = variant[i]
			k++
			l = w
		***REMOVED***
		if str := bytes.Join(variant[:k], separator); len(str) == 0 ***REMOVED***
			end = start - 1
		***REMOVED*** else ***REMOVED***
			scan.resizeRange(start, end, len(str))
			copy(scan.b[scan.start:], str)
			end = scan.end
		***REMOVED***
	***REMOVED***
	return end
***REMOVED***

type variantsSort struct ***REMOVED***
	i []uint8
	v [][]byte
***REMOVED***

func (s variantsSort) Len() int ***REMOVED***
	return len(s.i)
***REMOVED***

func (s variantsSort) Swap(i, j int) ***REMOVED***
	s.i[i], s.i[j] = s.i[j], s.i[i]
	s.v[i], s.v[j] = s.v[j], s.v[i]
***REMOVED***

func (s variantsSort) Less(i, j int) bool ***REMOVED***
	return s.i[i] < s.i[j]
***REMOVED***

type bytesSort [][]byte

func (b bytesSort) Len() int ***REMOVED***
	return len(b)
***REMOVED***

func (b bytesSort) Swap(i, j int) ***REMOVED***
	b[i], b[j] = b[j], b[i]
***REMOVED***

func (b bytesSort) Less(i, j int) bool ***REMOVED***
	return bytes.Compare(b[i], b[j]) == -1
***REMOVED***

// parseExtensions parses and normalizes the extensions in the buffer.
// It returns the last position of scan.b that is part of any extension.
// It also trims scan.b to remove excess parts accordingly.
func parseExtensions(scan *scanner) int ***REMOVED***
	start := scan.start
	exts := [][]byte***REMOVED******REMOVED***
	private := []byte***REMOVED******REMOVED***
	end := scan.end
	for len(scan.token) == 1 ***REMOVED***
		extStart := scan.start
		ext := scan.token[0]
		end = parseExtension(scan)
		extension := scan.b[extStart:end]
		if len(extension) < 3 || (ext != 'x' && len(extension) < 4) ***REMOVED***
			scan.setError(errSyntax)
			end = extStart
			continue
		***REMOVED*** else if start == extStart && (ext == 'x' || scan.start == len(scan.b)) ***REMOVED***
			scan.b = scan.b[:end]
			return end
		***REMOVED*** else if ext == 'x' ***REMOVED***
			private = extension
			break
		***REMOVED***
		exts = append(exts, extension)
	***REMOVED***
	sort.Sort(bytesSort(exts))
	if len(private) > 0 ***REMOVED***
		exts = append(exts, private)
	***REMOVED***
	scan.b = scan.b[:start]
	if len(exts) > 0 ***REMOVED***
		scan.b = append(scan.b, bytes.Join(exts, separator)...)
	***REMOVED*** else if start > 0 ***REMOVED***
		// Strip trailing '-'.
		scan.b = scan.b[:start-1]
	***REMOVED***
	return end
***REMOVED***

// parseExtension parses a single extension and returns the position of
// the extension end.
func parseExtension(scan *scanner) int ***REMOVED***
	start, end := scan.start, scan.end
	switch scan.token[0] ***REMOVED***
	case 'u':
		attrStart := end
		scan.scan()
		for last := []byte***REMOVED******REMOVED***; len(scan.token) > 2; scan.scan() ***REMOVED***
			if bytes.Compare(scan.token, last) != -1 ***REMOVED***
				// Attributes are unsorted. Start over from scratch.
				p := attrStart + 1
				scan.next = p
				attrs := [][]byte***REMOVED******REMOVED***
				for scan.scan(); len(scan.token) > 2; scan.scan() ***REMOVED***
					attrs = append(attrs, scan.token)
					end = scan.end
				***REMOVED***
				sort.Sort(bytesSort(attrs))
				copy(scan.b[p:], bytes.Join(attrs, separator))
				break
			***REMOVED***
			last = scan.token
			end = scan.end
		***REMOVED***
		var last, key []byte
		for attrEnd := end; len(scan.token) == 2; last = key ***REMOVED***
			key = scan.token
			keyEnd := scan.end
			end = scan.acceptMinSize(3)
			// TODO: check key value validity
			if keyEnd == end || bytes.Compare(key, last) != 1 ***REMOVED***
				// We have an invalid key or the keys are not sorted.
				// Start scanning keys from scratch and reorder.
				p := attrEnd + 1
				scan.next = p
				keys := [][]byte***REMOVED******REMOVED***
				for scan.scan(); len(scan.token) == 2; ***REMOVED***
					keyStart, keyEnd := scan.start, scan.end
					end = scan.acceptMinSize(3)
					if keyEnd != end ***REMOVED***
						keys = append(keys, scan.b[keyStart:end])
					***REMOVED*** else ***REMOVED***
						scan.setError(errSyntax)
						end = keyStart
					***REMOVED***
				***REMOVED***
				sort.Sort(bytesSort(keys))
				reordered := bytes.Join(keys, separator)
				if e := p + len(reordered); e < end ***REMOVED***
					scan.deleteRange(e, end)
					end = e
				***REMOVED***
				copy(scan.b[p:], bytes.Join(keys, separator))
				break
			***REMOVED***
		***REMOVED***
	case 't':
		scan.scan()
		if n := len(scan.token); n >= 2 && n <= 3 && isAlpha(scan.token[1]) ***REMOVED***
			_, end = parseTag(scan)
			scan.toLower(start, end)
		***REMOVED***
		for len(scan.token) == 2 && !isAlpha(scan.token[1]) ***REMOVED***
			end = scan.acceptMinSize(3)
		***REMOVED***
	case 'x':
		end = scan.acceptMinSize(1)
	default:
		end = scan.acceptMinSize(2)
	***REMOVED***
	return end
***REMOVED***

// Compose creates a Tag from individual parts, which may be of type Tag, Base,
// Script, Region, Variant, []Variant, Extension, []Extension or error. If a
// Base, Script or Region or slice of type Variant or Extension is passed more
// than once, the latter will overwrite the former. Variants and Extensions are
// accumulated, but if two extensions of the same type are passed, the latter
// will replace the former. A Tag overwrites all former values and typically
// only makes sense as the first argument. The resulting tag is returned after
// canonicalizing using the Default CanonType. If one or more errors are
// encountered, one of the errors is returned.
func Compose(part ...interface***REMOVED******REMOVED***) (t Tag, err error) ***REMOVED***
	return Default.Compose(part...)
***REMOVED***

// Compose creates a Tag from individual parts, which may be of type Tag, Base,
// Script, Region, Variant, []Variant, Extension, []Extension or error. If a
// Base, Script or Region or slice of type Variant or Extension is passed more
// than once, the latter will overwrite the former. Variants and Extensions are
// accumulated, but if two extensions of the same type are passed, the latter
// will replace the former. A Tag overwrites all former values and typically
// only makes sense as the first argument. The resulting tag is returned after
// canonicalizing using CanonType c. If one or more errors are encountered,
// one of the errors is returned.
func (c CanonType) Compose(part ...interface***REMOVED******REMOVED***) (t Tag, err error) ***REMOVED***
	var b builder
	if err = b.update(part...); err != nil ***REMOVED***
		return und, err
	***REMOVED***
	t, _ = b.tag.canonicalize(c)

	if len(b.ext) > 0 || len(b.variant) > 0 ***REMOVED***
		sort.Sort(sortVariant(b.variant))
		sort.Strings(b.ext)
		if b.private != "" ***REMOVED***
			b.ext = append(b.ext, b.private)
		***REMOVED***
		n := maxCoreSize + tokenLen(b.variant...) + tokenLen(b.ext...)
		buf := make([]byte, n)
		p := t.genCoreBytes(buf)
		t.pVariant = byte(p)
		p += appendTokens(buf[p:], b.variant...)
		t.pExt = uint16(p)
		p += appendTokens(buf[p:], b.ext...)
		t.str = string(buf[:p])
	***REMOVED*** else if b.private != "" ***REMOVED***
		t.str = b.private
		t.remakeString()
	***REMOVED***
	return
***REMOVED***

type builder struct ***REMOVED***
	tag Tag

	private string // the x extension
	ext     []string
	variant []string

	err error
***REMOVED***

func (b *builder) addExt(e string) ***REMOVED***
	if e == "" ***REMOVED***
	***REMOVED*** else if e[0] == 'x' ***REMOVED***
		b.private = e
	***REMOVED*** else ***REMOVED***
		b.ext = append(b.ext, e)
	***REMOVED***
***REMOVED***

var errInvalidArgument = errors.New("invalid Extension or Variant")

func (b *builder) update(part ...interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	replace := func(l *[]string, s string, eq func(a, b string) bool) bool ***REMOVED***
		if s == "" ***REMOVED***
			b.err = errInvalidArgument
			return true
		***REMOVED***
		for i, v := range *l ***REMOVED***
			if eq(v, s) ***REMOVED***
				(*l)[i] = s
				return true
			***REMOVED***
		***REMOVED***
		return false
	***REMOVED***
	for _, x := range part ***REMOVED***
		switch v := x.(type) ***REMOVED***
		case Tag:
			b.tag.lang = v.lang
			b.tag.region = v.region
			b.tag.script = v.script
			if v.str != "" ***REMOVED***
				b.variant = nil
				for x, s := "", v.str[v.pVariant:v.pExt]; s != ""; ***REMOVED***
					x, s = nextToken(s)
					b.variant = append(b.variant, x)
				***REMOVED***
				b.ext, b.private = nil, ""
				for i, e := int(v.pExt), ""; i < len(v.str); ***REMOVED***
					i, e = getExtension(v.str, i)
					b.addExt(e)
				***REMOVED***
			***REMOVED***
		case Base:
			b.tag.lang = v.langID
		case Script:
			b.tag.script = v.scriptID
		case Region:
			b.tag.region = v.regionID
		case Variant:
			if !replace(&b.variant, v.variant, func(a, b string) bool ***REMOVED*** return a == b ***REMOVED***) ***REMOVED***
				b.variant = append(b.variant, v.variant)
			***REMOVED***
		case Extension:
			if !replace(&b.ext, v.s, func(a, b string) bool ***REMOVED*** return a[0] == b[0] ***REMOVED***) ***REMOVED***
				b.addExt(v.s)
			***REMOVED***
		case []Variant:
			b.variant = nil
			for _, x := range v ***REMOVED***
				b.update(x)
			***REMOVED***
		case []Extension:
			b.ext, b.private = nil, ""
			for _, e := range v ***REMOVED***
				b.update(e)
			***REMOVED***
		// TODO: support parsing of raw strings based on morphology or just extensions?
		case error:
			err = v
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func tokenLen(token ...string) (n int) ***REMOVED***
	for _, t := range token ***REMOVED***
		n += len(t) + 1
	***REMOVED***
	return
***REMOVED***

func appendTokens(b []byte, token ...string) int ***REMOVED***
	p := 0
	for _, t := range token ***REMOVED***
		b[p] = '-'
		copy(b[p+1:], t)
		p += 1 + len(t)
	***REMOVED***
	return p
***REMOVED***

type sortVariant []string

func (s sortVariant) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s sortVariant) Swap(i, j int) ***REMOVED***
	s[j], s[i] = s[i], s[j]
***REMOVED***

func (s sortVariant) Less(i, j int) bool ***REMOVED***
	return variantIndex[s[i]] < variantIndex[s[j]]
***REMOVED***

func findExt(list []string, x byte) int ***REMOVED***
	for i, e := range list ***REMOVED***
		if e[0] == x ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// getExtension returns the name, body and end position of the extension.
func getExtension(s string, p int) (end int, ext string) ***REMOVED***
	if s[p] == '-' ***REMOVED***
		p++
	***REMOVED***
	if s[p] == 'x' ***REMOVED***
		return len(s), s[p:]
	***REMOVED***
	end = nextExtension(s, p)
	return end, s[p:end]
***REMOVED***

// nextExtension finds the next extension within the string, searching
// for the -<char>- pattern from position p.
// In the fast majority of cases, language tags will have at most
// one extension and extensions tend to be small.
func nextExtension(s string, p int) int ***REMOVED***
	for n := len(s) - 3; p < n; ***REMOVED***
		if s[p] == '-' ***REMOVED***
			if s[p+2] == '-' ***REMOVED***
				return p
			***REMOVED***
			p += 3
		***REMOVED*** else ***REMOVED***
			p++
		***REMOVED***
	***REMOVED***
	return len(s)
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
			t = Tag***REMOVED***lang: id***REMOVED***
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

// Add hack mapping to deal with a small number of cases that that occur
// in Accept-Language (with reasonable frequency).
var acceptFallback = map[string]langID***REMOVED***
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
