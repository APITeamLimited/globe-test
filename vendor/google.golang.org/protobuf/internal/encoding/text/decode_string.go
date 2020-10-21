// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/strs"
)

// parseStringValue parses string field token.
// This differs from parseString since the text format allows
// multiple back-to-back string literals where they are semantically treated
// as a single large string with all values concatenated.
//
// E.g., `"foo" "bar" "baz"` => "foobarbaz"
func (d *Decoder) parseStringValue() (Token, error) ***REMOVED***
	// Note that the ending quote is sufficient to unambiguously mark the end
	// of a string. Thus, the text grammar does not require intervening
	// whitespace or control characters in-between strings.
	// Thus, the following is valid:
	//	`"foo"'bar'"baz"` => "foobarbaz"
	in0 := d.in
	var ss []string
	for len(d.in) > 0 && (d.in[0] == '"' || d.in[0] == '\'') ***REMOVED***
		s, err := d.parseString()
		if err != nil ***REMOVED***
			return Token***REMOVED******REMOVED***, err
		***REMOVED***
		ss = append(ss, s)
	***REMOVED***
	// d.in already points to the end of the value at this point.
	return Token***REMOVED***
		kind:  Scalar,
		attrs: stringValue,
		pos:   len(d.orig) - len(in0),
		raw:   in0[:len(in0)-len(d.in)],
		str:   strings.Join(ss, ""),
	***REMOVED***, nil
***REMOVED***

// parseString parses a string value enclosed in " or '.
func (d *Decoder) parseString() (string, error) ***REMOVED***
	in := d.in
	if len(in) == 0 ***REMOVED***
		return "", ErrUnexpectedEOF
	***REMOVED***
	quote := in[0]
	in = in[1:]
	i := indexNeedEscapeInBytes(in)
	in, out := in[i:], in[:i:i] // set cap to prevent mutations
	for len(in) > 0 ***REMOVED***
		switch r, n := utf8.DecodeRune(in); ***REMOVED***
		case r == utf8.RuneError && n == 1:
			return "", d.newSyntaxError("invalid UTF-8 detected")
		case r == 0 || r == '\n':
			return "", d.newSyntaxError("invalid character %q in string", r)
		case r == rune(quote):
			in = in[1:]
			d.consume(len(d.in) - len(in))
			return string(out), nil
		case r == '\\':
			if len(in) < 2 ***REMOVED***
				return "", ErrUnexpectedEOF
			***REMOVED***
			switch r := in[1]; r ***REMOVED***
			case '"', '\'', '\\', '?':
				in, out = in[2:], append(out, r)
			case 'a':
				in, out = in[2:], append(out, '\a')
			case 'b':
				in, out = in[2:], append(out, '\b')
			case 'n':
				in, out = in[2:], append(out, '\n')
			case 'r':
				in, out = in[2:], append(out, '\r')
			case 't':
				in, out = in[2:], append(out, '\t')
			case 'v':
				in, out = in[2:], append(out, '\v')
			case 'f':
				in, out = in[2:], append(out, '\f')
			case '0', '1', '2', '3', '4', '5', '6', '7':
				// One, two, or three octal characters.
				n := len(in[1:]) - len(bytes.TrimLeft(in[1:], "01234567"))
				if n > 3 ***REMOVED***
					n = 3
				***REMOVED***
				v, err := strconv.ParseUint(string(in[1:1+n]), 8, 8)
				if err != nil ***REMOVED***
					return "", d.newSyntaxError("invalid octal escape code %q in string", in[:1+n])
				***REMOVED***
				in, out = in[1+n:], append(out, byte(v))
			case 'x':
				// One or two hexadecimal characters.
				n := len(in[2:]) - len(bytes.TrimLeft(in[2:], "0123456789abcdefABCDEF"))
				if n > 2 ***REMOVED***
					n = 2
				***REMOVED***
				v, err := strconv.ParseUint(string(in[2:2+n]), 16, 8)
				if err != nil ***REMOVED***
					return "", d.newSyntaxError("invalid hex escape code %q in string", in[:2+n])
				***REMOVED***
				in, out = in[2+n:], append(out, byte(v))
			case 'u', 'U':
				// Four or eight hexadecimal characters
				n := 6
				if r == 'U' ***REMOVED***
					n = 10
				***REMOVED***
				if len(in) < n ***REMOVED***
					return "", ErrUnexpectedEOF
				***REMOVED***
				v, err := strconv.ParseUint(string(in[2:n]), 16, 32)
				if utf8.MaxRune < v || err != nil ***REMOVED***
					return "", d.newSyntaxError("invalid Unicode escape code %q in string", in[:n])
				***REMOVED***
				in = in[n:]

				r := rune(v)
				if utf16.IsSurrogate(r) ***REMOVED***
					if len(in) < 6 ***REMOVED***
						return "", ErrUnexpectedEOF
					***REMOVED***
					v, err := strconv.ParseUint(string(in[2:6]), 16, 16)
					r = utf16.DecodeRune(r, rune(v))
					if in[0] != '\\' || in[1] != 'u' || r == unicode.ReplacementChar || err != nil ***REMOVED***
						return "", d.newSyntaxError("invalid Unicode escape code %q in string", in[:6])
					***REMOVED***
					in = in[6:]
				***REMOVED***
				out = append(out, string(r)...)
			default:
				return "", d.newSyntaxError("invalid escape code %q in string", in[:2])
			***REMOVED***
		default:
			i := indexNeedEscapeInBytes(in[n:])
			in, out = in[n+i:], append(out, in[:n+i]...)
		***REMOVED***
	***REMOVED***
	return "", ErrUnexpectedEOF
***REMOVED***

// indexNeedEscapeInString returns the index of the character that needs
// escaping. If no characters need escaping, this returns the input length.
func indexNeedEscapeInBytes(b []byte) int ***REMOVED*** return indexNeedEscapeInString(strs.UnsafeString(b)) ***REMOVED***

// UnmarshalString returns an unescaped string given a textproto string value.
// String value needs to contain single or double quotes. This is only used by
// internal/encoding/defval package for unmarshaling bytes.
func UnmarshalString(s string) (string, error) ***REMOVED***
	d := NewDecoder([]byte(s))
	return d.parseString()
***REMOVED***
