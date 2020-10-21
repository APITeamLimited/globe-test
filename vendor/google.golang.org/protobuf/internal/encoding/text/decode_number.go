// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

// parseNumberValue parses a number from the input and returns a Token object.
func (d *Decoder) parseNumberValue() (Token, bool) ***REMOVED***
	in := d.in
	num := parseNumber(in)
	if num.size == 0 ***REMOVED***
		return Token***REMOVED******REMOVED***, false
	***REMOVED***
	numAttrs := num.kind
	if num.neg ***REMOVED***
		numAttrs |= isNegative
	***REMOVED***
	strSize := num.size
	last := num.size - 1
	if num.kind == numFloat && (d.in[last] == 'f' || d.in[last] == 'F') ***REMOVED***
		strSize = last
	***REMOVED***
	tok := Token***REMOVED***
		kind:     Scalar,
		attrs:    numberValue,
		pos:      len(d.orig) - len(d.in),
		raw:      d.in[:num.size],
		str:      string(d.in[:strSize]),
		numAttrs: numAttrs,
	***REMOVED***
	d.consume(num.size)
	return tok, true
***REMOVED***

const (
	numDec uint8 = (1 << iota) / 2
	numHex
	numOct
	numFloat
)

// number is the result of parsing out a valid number from parseNumber. It
// contains data for doing float or integer conversion via the strconv package
// in conjunction with the input bytes.
type number struct ***REMOVED***
	kind uint8
	neg  bool
	size int
***REMOVED***

// parseNumber constructs a number object from given input. It allows for the
// following patterns:
//   integer: ^-?([1-9][0-9]*|0[xX][0-9a-fA-F]+|0[0-7]*)
//   float: ^-?((0|[1-9][0-9]*)?([.][0-9]*)?([eE][+-]?[0-9]+)?[fF]?)
// It also returns the number of parsed bytes for the given number, 0 if it is
// not a number.
func parseNumber(input []byte) number ***REMOVED***
	kind := numDec
	var size int
	var neg bool

	s := input
	if len(s) == 0 ***REMOVED***
		return number***REMOVED******REMOVED***
	***REMOVED***

	// Optional -
	if s[0] == '-' ***REMOVED***
		neg = true
		s = s[1:]
		size++
		if len(s) == 0 ***REMOVED***
			return number***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***

	// C++ allows for whitespace and comments in between the negative sign and
	// the rest of the number. This logic currently does not but is consistent
	// with v1.

	switch ***REMOVED***
	case s[0] == '0':
		if len(s) > 1 ***REMOVED***
			switch ***REMOVED***
			case s[1] == 'x' || s[1] == 'X':
				// Parse as hex number.
				kind = numHex
				n := 2
				s = s[2:]
				for len(s) > 0 && (('0' <= s[0] && s[0] <= '9') ||
					('a' <= s[0] && s[0] <= 'f') ||
					('A' <= s[0] && s[0] <= 'F')) ***REMOVED***
					s = s[1:]
					n++
				***REMOVED***
				if n == 2 ***REMOVED***
					return number***REMOVED******REMOVED***
				***REMOVED***
				size += n

			case '0' <= s[1] && s[1] <= '7':
				// Parse as octal number.
				kind = numOct
				n := 2
				s = s[2:]
				for len(s) > 0 && '0' <= s[0] && s[0] <= '7' ***REMOVED***
					s = s[1:]
					n++
				***REMOVED***
				size += n
			***REMOVED***

			if kind&(numHex|numOct) > 0 ***REMOVED***
				if len(s) > 0 && !isDelim(s[0]) ***REMOVED***
					return number***REMOVED******REMOVED***
				***REMOVED***
				return number***REMOVED***kind: kind, neg: neg, size: size***REMOVED***
			***REMOVED***
		***REMOVED***
		s = s[1:]
		size++

	case '1' <= s[0] && s[0] <= '9':
		n := 1
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' ***REMOVED***
			s = s[1:]
			n++
		***REMOVED***
		size += n

	case s[0] == '.':
		// Set kind to numFloat to signify the intent to parse as float. And
		// that it needs to have other digits after '.'.
		kind = numFloat

	default:
		return number***REMOVED******REMOVED***
	***REMOVED***

	// . followed by 0 or more digits.
	if len(s) > 0 && s[0] == '.' ***REMOVED***
		n := 1
		s = s[1:]
		// If decimal point was before any digits, it should be followed by
		// other digits.
		if len(s) == 0 && kind == numFloat ***REMOVED***
			return number***REMOVED******REMOVED***
		***REMOVED***
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' ***REMOVED***
			s = s[1:]
			n++
		***REMOVED***
		size += n
		kind = numFloat
	***REMOVED***

	// e or E followed by an optional - or + and 1 or more digits.
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') ***REMOVED***
		kind = numFloat
		s = s[1:]
		n := 1
		if s[0] == '+' || s[0] == '-' ***REMOVED***
			s = s[1:]
			n++
			if len(s) == 0 ***REMOVED***
				return number***REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' ***REMOVED***
			s = s[1:]
			n++
		***REMOVED***
		size += n
	***REMOVED***

	// Optional suffix f or F for floats.
	if len(s) > 0 && (s[0] == 'f' || s[0] == 'F') ***REMOVED***
		kind = numFloat
		s = s[1:]
		size++
	***REMOVED***

	// Check that next byte is a delimiter or it is at the end.
	if len(s) > 0 && !isDelim(s[0]) ***REMOVED***
		return number***REMOVED******REMOVED***
	***REMOVED***

	return number***REMOVED***kind: kind, neg: neg, size: size***REMOVED***
***REMOVED***
