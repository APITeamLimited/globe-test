// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

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

// ErrSyntax is returned by any of the parsing functions when the
// input is not well-formed, according to BCP 47.
// TODO: return the position at which the syntax error occurred?
var ErrSyntax = errors.New("language: tag is not well-formed")

// ErrDuplicateKey is returned when a tag contains the same key twice with
// different values in the -u section.
var ErrDuplicateKey = errors.New("language: different values for same key in -u extension")

// ValueError is returned by any of the parsing functions when the
// input is well-formed but the respective subtag is not recognized
// as a valid value.
type ValueError struct ***REMOVED***
	v [8]byte
***REMOVED***

// NewValueError creates a new ValueError.
func NewValueError(tag []byte) ValueError ***REMOVED***
	var e ValueError
	copy(e.v[:], tag)
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
	if s.err == nil || (e == ErrSyntax && s.err != ErrSyntax) ***REMOVED***
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
		var b []byte
		if n := len(s.b) + diff; n > cap(s.b) ***REMOVED***
			b = make([]byte, n)
			copy(b, s.b[:oldStart])
		***REMOVED*** else ***REMOVED***
			b = s.b[:n]
		***REMOVED***
		copy(b[end:], s.b[oldEnd:])
		s.b = b
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
			s.gobble(ErrSyntax)
			continue
		***REMOVED***
		s.token = token
		return end
	***REMOVED***
	if n := len(s.b); n > 0 && s.b[n-1] == '-' ***REMOVED***
		s.setError(ErrSyntax)
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
// https://www.unicode.org/reports/tr35/#Unicode_Language_and_Locale_Identifiers.
func Parse(s string) (t Tag, err error) ***REMOVED***
	// TODO: consider supporting old-style locale key-value pairs.
	if s == "" ***REMOVED***
		return Und, ErrSyntax
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
	return parse(&scan, s)
***REMOVED***

func parse(scan *scanner, s string) (t Tag, err error) ***REMOVED***
	t = Und
	var end int
	if n := len(scan.token); n <= 1 ***REMOVED***
		scan.toLower(0, len(scan.b))
		if n == 0 || scan.token[0] != 'x' ***REMOVED***
			return t, ErrSyntax
		***REMOVED***
		end = parseExtensions(scan)
	***REMOVED*** else if n >= 4 ***REMOVED***
		return Und, ErrSyntax
	***REMOVED*** else ***REMOVED*** // the usual case
		t, end = parseTag(scan)
		if n := len(scan.token); n == 1 ***REMOVED***
			t.pExt = uint16(end)
			end = parseExtensions(scan)
		***REMOVED*** else if end < len(scan.b) ***REMOVED***
			scan.setError(ErrSyntax)
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
	t.LangID, e = getLangID(scan.token)
	scan.setError(e)
	scan.replace(t.LangID.String())
	langStart := scan.start
	end = scan.scan()
	for len(scan.token) == 3 && isAlpha(scan.token[0]) ***REMOVED***
		// From http://tools.ietf.org/html/bcp47, <lang>-<extlang> tags are equivalent
		// to a tag of the form <extlang>.
		lang, e := getLangID(scan.token)
		if lang != 0 ***REMOVED***
			t.LangID = lang
			copy(scan.b[langStart:], lang.String())
			scan.b[langStart+3] = '-'
			scan.start = langStart + 4
		***REMOVED***
		scan.gobble(e)
		end = scan.scan()
	***REMOVED***
	if len(scan.token) == 4 && isAlpha(scan.token[0]) ***REMOVED***
		t.ScriptID, e = getScriptID(script, scan.token)
		if t.ScriptID == 0 ***REMOVED***
			scan.gobble(e)
		***REMOVED***
		end = scan.scan()
	***REMOVED***
	if n := len(scan.token); n >= 2 && n <= 3 ***REMOVED***
		t.RegionID, e = getRegionID(scan.token)
		if t.RegionID == 0 ***REMOVED***
			scan.gobble(e)
		***REMOVED*** else ***REMOVED***
			scan.replace(t.RegionID.String())
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
			scan.gobble(NewValueError(scan.token))
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

type bytesSort struct ***REMOVED***
	b [][]byte
	n int // first n bytes to compare
***REMOVED***

func (b bytesSort) Len() int ***REMOVED***
	return len(b.b)
***REMOVED***

func (b bytesSort) Swap(i, j int) ***REMOVED***
	b.b[i], b.b[j] = b.b[j], b.b[i]
***REMOVED***

func (b bytesSort) Less(i, j int) bool ***REMOVED***
	for k := 0; k < b.n; k++ ***REMOVED***
		if b.b[i][k] == b.b[j][k] ***REMOVED***
			continue
		***REMOVED***
		return b.b[i][k] < b.b[j][k]
	***REMOVED***
	return false
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
			scan.setError(ErrSyntax)
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
	sort.Sort(bytesSort***REMOVED***exts, 1***REMOVED***)
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
	case 'u': // https://www.ietf.org/rfc/rfc6067.txt
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
				sort.Sort(bytesSort***REMOVED***attrs, 3***REMOVED***)
				copy(scan.b[p:], bytes.Join(attrs, separator))
				break
			***REMOVED***
			last = scan.token
			end = scan.end
		***REMOVED***
		// Scan key-type sequences. A key is of length 2 and may be followed
		// by 0 or more "type" subtags from 3 to the maximum of 8 letters.
		var last, key []byte
		for attrEnd := end; len(scan.token) == 2; last = key ***REMOVED***
			key = scan.token
			end = scan.end
			for scan.scan(); end < scan.end && len(scan.token) > 2; scan.scan() ***REMOVED***
				end = scan.end
			***REMOVED***
			// TODO: check key value validity
			if bytes.Compare(key, last) != 1 || scan.err != nil ***REMOVED***
				// We have an invalid key or the keys are not sorted.
				// Start scanning keys from scratch and reorder.
				p := attrEnd + 1
				scan.next = p
				keys := [][]byte***REMOVED******REMOVED***
				for scan.scan(); len(scan.token) == 2; ***REMOVED***
					keyStart := scan.start
					end = scan.end
					for scan.scan(); end < scan.end && len(scan.token) > 2; scan.scan() ***REMOVED***
						end = scan.end
					***REMOVED***
					keys = append(keys, scan.b[keyStart:end])
				***REMOVED***
				sort.Stable(bytesSort***REMOVED***keys, 2***REMOVED***)
				if n := len(keys); n > 0 ***REMOVED***
					k := 0
					for i := 1; i < n; i++ ***REMOVED***
						if !bytes.Equal(keys[k][:2], keys[i][:2]) ***REMOVED***
							k++
							keys[k] = keys[i]
						***REMOVED*** else if !bytes.Equal(keys[k], keys[i]) ***REMOVED***
							scan.setError(ErrDuplicateKey)
						***REMOVED***
					***REMOVED***
					keys = keys[:k+1]
				***REMOVED***
				reordered := bytes.Join(keys, separator)
				if e := p + len(reordered); e < end ***REMOVED***
					scan.deleteRange(e, end)
					end = e
				***REMOVED***
				copy(scan.b[p:], reordered)
				break
			***REMOVED***
		***REMOVED***
	case 't': // https://www.ietf.org/rfc/rfc6497.txt
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
