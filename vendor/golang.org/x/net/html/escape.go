// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

// These replacements permit compatibility with old numeric entities that
// assumed Windows-1252 encoding.
// https://html.spec.whatwg.org/multipage/syntax.html#consume-a-character-reference
var replacementTable = [...]rune***REMOVED***
	'\u20AC', // First entry is what 0x80 should be replaced with.
	'\u0081',
	'\u201A',
	'\u0192',
	'\u201E',
	'\u2026',
	'\u2020',
	'\u2021',
	'\u02C6',
	'\u2030',
	'\u0160',
	'\u2039',
	'\u0152',
	'\u008D',
	'\u017D',
	'\u008F',
	'\u0090',
	'\u2018',
	'\u2019',
	'\u201C',
	'\u201D',
	'\u2022',
	'\u2013',
	'\u2014',
	'\u02DC',
	'\u2122',
	'\u0161',
	'\u203A',
	'\u0153',
	'\u009D',
	'\u017E',
	'\u0178', // Last entry is 0x9F.
	// 0x00->'\uFFFD' is handled programmatically.
	// 0x0D->'\u000D' is a no-op.
***REMOVED***

// unescapeEntity reads an entity like "&lt;" from b[src:] and writes the
// corresponding "<" to b[dst:], returning the incremented dst and src cursors.
// Precondition: b[src] == '&' && dst <= src.
// attribute should be true if parsing an attribute value.
func unescapeEntity(b []byte, dst, src int, attribute bool) (dst1, src1 int) ***REMOVED***
	// https://html.spec.whatwg.org/multipage/syntax.html#consume-a-character-reference

	// i starts at 1 because we already know that s[0] == '&'.
	i, s := 1, b[src:]

	if len(s) <= 1 ***REMOVED***
		b[dst] = b[src]
		return dst + 1, src + 1
	***REMOVED***

	if s[i] == '#' ***REMOVED***
		if len(s) <= 3 ***REMOVED*** // We need to have at least "&#.".
			b[dst] = b[src]
			return dst + 1, src + 1
		***REMOVED***
		i++
		c := s[i]
		hex := false
		if c == 'x' || c == 'X' ***REMOVED***
			hex = true
			i++
		***REMOVED***

		x := '\x00'
		for i < len(s) ***REMOVED***
			c = s[i]
			i++
			if hex ***REMOVED***
				if '0' <= c && c <= '9' ***REMOVED***
					x = 16*x + rune(c) - '0'
					continue
				***REMOVED*** else if 'a' <= c && c <= 'f' ***REMOVED***
					x = 16*x + rune(c) - 'a' + 10
					continue
				***REMOVED*** else if 'A' <= c && c <= 'F' ***REMOVED***
					x = 16*x + rune(c) - 'A' + 10
					continue
				***REMOVED***
			***REMOVED*** else if '0' <= c && c <= '9' ***REMOVED***
				x = 10*x + rune(c) - '0'
				continue
			***REMOVED***
			if c != ';' ***REMOVED***
				i--
			***REMOVED***
			break
		***REMOVED***

		if i <= 3 ***REMOVED*** // No characters matched.
			b[dst] = b[src]
			return dst + 1, src + 1
		***REMOVED***

		if 0x80 <= x && x <= 0x9F ***REMOVED***
			// Replace characters from Windows-1252 with UTF-8 equivalents.
			x = replacementTable[x-0x80]
		***REMOVED*** else if x == 0 || (0xD800 <= x && x <= 0xDFFF) || x > 0x10FFFF ***REMOVED***
			// Replace invalid characters with the replacement character.
			x = '\uFFFD'
		***REMOVED***

		return dst + utf8.EncodeRune(b[dst:], x), src + i
	***REMOVED***

	// Consume the maximum number of characters possible, with the
	// consumed characters matching one of the named references.

	for i < len(s) ***REMOVED***
		c := s[i]
		i++
		// Lower-cased characters are more common in entities, so we check for them first.
		if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' ***REMOVED***
			continue
		***REMOVED***
		if c != ';' ***REMOVED***
			i--
		***REMOVED***
		break
	***REMOVED***

	entityName := string(s[1:i])
	if entityName == "" ***REMOVED***
		// No-op.
	***REMOVED*** else if attribute && entityName[len(entityName)-1] != ';' && len(s) > i && s[i] == '=' ***REMOVED***
		// No-op.
	***REMOVED*** else if x := entity[entityName]; x != 0 ***REMOVED***
		return dst + utf8.EncodeRune(b[dst:], x), src + i
	***REMOVED*** else if x := entity2[entityName]; x[0] != 0 ***REMOVED***
		dst1 := dst + utf8.EncodeRune(b[dst:], x[0])
		return dst1 + utf8.EncodeRune(b[dst1:], x[1]), src + i
	***REMOVED*** else if !attribute ***REMOVED***
		maxLen := len(entityName) - 1
		if maxLen > longestEntityWithoutSemicolon ***REMOVED***
			maxLen = longestEntityWithoutSemicolon
		***REMOVED***
		for j := maxLen; j > 1; j-- ***REMOVED***
			if x := entity[entityName[:j]]; x != 0 ***REMOVED***
				return dst + utf8.EncodeRune(b[dst:], x), src + j + 1
			***REMOVED***
		***REMOVED***
	***REMOVED***

	dst1, src1 = dst+i, src+i
	copy(b[dst:dst1], b[src:src1])
	return dst1, src1
***REMOVED***

// unescape unescapes b's entities in-place, so that "a&lt;b" becomes "a<b".
// attribute should be true if parsing an attribute value.
func unescape(b []byte, attribute bool) []byte ***REMOVED***
	for i, c := range b ***REMOVED***
		if c == '&' ***REMOVED***
			dst, src := unescapeEntity(b, i, i, attribute)
			for src < len(b) ***REMOVED***
				c := b[src]
				if c == '&' ***REMOVED***
					dst, src = unescapeEntity(b, dst, src, attribute)
				***REMOVED*** else ***REMOVED***
					b[dst] = c
					dst, src = dst+1, src+1
				***REMOVED***
			***REMOVED***
			return b[0:dst]
		***REMOVED***
	***REMOVED***
	return b
***REMOVED***

// lower lower-cases the A-Z bytes in b in-place, so that "aBc" becomes "abc".
func lower(b []byte) []byte ***REMOVED***
	for i, c := range b ***REMOVED***
		if 'A' <= c && c <= 'Z' ***REMOVED***
			b[i] = c + 'a' - 'A'
		***REMOVED***
	***REMOVED***
	return b
***REMOVED***

const escapedChars = "&'<>\"\r"

func escape(w writer, s string) error ***REMOVED***
	i := strings.IndexAny(s, escapedChars)
	for i != -1 ***REMOVED***
		if _, err := w.WriteString(s[:i]); err != nil ***REMOVED***
			return err
		***REMOVED***
		var esc string
		switch s[i] ***REMOVED***
		case '&':
			esc = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		***REMOVED***
		s = s[i+1:]
		if _, err := w.WriteString(esc); err != nil ***REMOVED***
			return err
		***REMOVED***
		i = strings.IndexAny(s, escapedChars)
	***REMOVED***
	_, err := w.WriteString(s)
	return err
***REMOVED***

// EscapeString escapes special characters like "<" to become "&lt;". It
// escapes only five such characters: <, >, &, ' and ".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
func EscapeString(s string) string ***REMOVED***
	if strings.IndexAny(s, escapedChars) == -1 ***REMOVED***
		return s
	***REMOVED***
	var buf bytes.Buffer
	escape(&buf, s)
	return buf.String()
***REMOVED***

// UnescapeString unescapes entities like "&lt;" to become "<". It unescapes a
// larger range of entities than EscapeString escapes. For example, "&aacute;"
// unescapes to "á", as does "&#225;" and "&xE1;".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
func UnescapeString(s string) string ***REMOVED***
	for _, c := range s ***REMOVED***
		if c == '&' ***REMOVED***
			return string(unescape([]byte(s), false))
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***