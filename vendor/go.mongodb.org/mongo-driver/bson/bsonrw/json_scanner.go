// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"unicode"
	"unicode/utf16"
)

type jsonTokenType byte

const (
	jttBeginObject jsonTokenType = iota
	jttEndObject
	jttBeginArray
	jttEndArray
	jttColon
	jttComma
	jttInt32
	jttInt64
	jttDouble
	jttString
	jttBool
	jttNull
	jttEOF
)

type jsonToken struct ***REMOVED***
	t jsonTokenType
	v interface***REMOVED******REMOVED***
	p int
***REMOVED***

type jsonScanner struct ***REMOVED***
	r           io.Reader
	buf         []byte
	pos         int
	lastReadErr error
***REMOVED***

// nextToken returns the next JSON token if one exists. A token is a character
// of the JSON grammar, a number, a string, or a literal.
func (js *jsonScanner) nextToken() (*jsonToken, error) ***REMOVED***
	c, err := js.readNextByte()

	// keep reading until a non-space is encountered (break on read error or EOF)
	for isWhiteSpace(c) && err == nil ***REMOVED***
		c, err = js.readNextByte()
	***REMOVED***

	if err == io.EOF ***REMOVED***
		return &jsonToken***REMOVED***t: jttEOF***REMOVED***, nil
	***REMOVED*** else if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// switch on the character
	switch c ***REMOVED***
	case '***REMOVED***':
		return &jsonToken***REMOVED***t: jttBeginObject, v: byte('***REMOVED***'), p: js.pos - 1***REMOVED***, nil
	case '***REMOVED***':
		return &jsonToken***REMOVED***t: jttEndObject, v: byte('***REMOVED***'), p: js.pos - 1***REMOVED***, nil
	case '[':
		return &jsonToken***REMOVED***t: jttBeginArray, v: byte('['), p: js.pos - 1***REMOVED***, nil
	case ']':
		return &jsonToken***REMOVED***t: jttEndArray, v: byte(']'), p: js.pos - 1***REMOVED***, nil
	case ':':
		return &jsonToken***REMOVED***t: jttColon, v: byte(':'), p: js.pos - 1***REMOVED***, nil
	case ',':
		return &jsonToken***REMOVED***t: jttComma, v: byte(','), p: js.pos - 1***REMOVED***, nil
	case '"': // RFC-8259 only allows for double quotes (") not single (')
		return js.scanString()
	default:
		// check if it's a number
		if c == '-' || isDigit(c) ***REMOVED***
			return js.scanNumber(c)
		***REMOVED*** else if c == 't' || c == 'f' || c == 'n' ***REMOVED***
			// maybe a literal
			return js.scanLiteral(c)
		***REMOVED*** else ***REMOVED***
			return nil, fmt.Errorf("invalid JSON input. Position: %d. Character: %c", js.pos-1, c)
		***REMOVED***
	***REMOVED***
***REMOVED***

// readNextByte attempts to read the next byte from the buffer. If the buffer
// has been exhausted, this function calls readIntoBuf, thus refilling the
// buffer and resetting the read position to 0
func (js *jsonScanner) readNextByte() (byte, error) ***REMOVED***
	if js.pos >= len(js.buf) ***REMOVED***
		err := js.readIntoBuf()

		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***

	b := js.buf[js.pos]
	js.pos++

	return b, nil
***REMOVED***

// readNNextBytes reads n bytes into dst, starting at offset
func (js *jsonScanner) readNNextBytes(dst []byte, n, offset int) error ***REMOVED***
	var err error

	for i := 0; i < n; i++ ***REMOVED***
		dst[i+offset], err = js.readNextByte()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// readIntoBuf reads up to 512 bytes from the scanner's io.Reader into the buffer
func (js *jsonScanner) readIntoBuf() error ***REMOVED***
	if js.lastReadErr != nil ***REMOVED***
		js.buf = js.buf[:0]
		js.pos = 0
		return js.lastReadErr
	***REMOVED***

	if cap(js.buf) == 0 ***REMOVED***
		js.buf = make([]byte, 0, 512)
	***REMOVED***

	n, err := js.r.Read(js.buf[:cap(js.buf)])
	if err != nil ***REMOVED***
		js.lastReadErr = err
		if n > 0 ***REMOVED***
			err = nil
		***REMOVED***
	***REMOVED***
	js.buf = js.buf[:n]
	js.pos = 0

	return err
***REMOVED***

func isWhiteSpace(c byte) bool ***REMOVED***
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
***REMOVED***

func isDigit(c byte) bool ***REMOVED***
	return unicode.IsDigit(rune(c))
***REMOVED***

func isValueTerminator(c byte) bool ***REMOVED***
	return c == ',' || c == '***REMOVED***' || c == ']' || isWhiteSpace(c)
***REMOVED***

// getu4 decodes the 4-byte hex sequence from the beginning of s, returning the hex value as a rune,
// or it returns -1. Note that the "\u" from the unicode escape sequence should not be present.
// It is copied and lightly modified from the Go JSON decode function at
// https://github.com/golang/go/blob/1b0a0316802b8048d69da49dc23c5a5ab08e8ae8/src/encoding/json/decode.go#L1169-L1188
func getu4(s []byte) rune ***REMOVED***
	if len(s) < 4 ***REMOVED***
		return -1
	***REMOVED***
	var r rune
	for _, c := range s[:4] ***REMOVED***
		switch ***REMOVED***
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		***REMOVED***
		r = r*16 + rune(c)
	***REMOVED***
	return r
***REMOVED***

// scanString reads from an opening '"' to a closing '"' and handles escaped characters
func (js *jsonScanner) scanString() (*jsonToken, error) ***REMOVED***
	var b bytes.Buffer
	var c byte
	var err error

	p := js.pos - 1

	for ***REMOVED***
		c, err = js.readNextByte()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				return nil, errors.New("end of input in JSON string")
			***REMOVED***
			return nil, err
		***REMOVED***

	evalNextChar:
		switch c ***REMOVED***
		case '\\':
			c, err = js.readNextByte()
			if err != nil ***REMOVED***
				if err == io.EOF ***REMOVED***
					return nil, errors.New("end of input in JSON string")
				***REMOVED***
				return nil, err
			***REMOVED***

		evalNextEscapeChar:
			switch c ***REMOVED***
			case '"', '\\', '/':
				b.WriteByte(c)
			case 'b':
				b.WriteByte('\b')
			case 'f':
				b.WriteByte('\f')
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case 'u':
				us := make([]byte, 4)
				err = js.readNNextBytes(us, 4, 0)
				if err != nil ***REMOVED***
					return nil, fmt.Errorf("invalid unicode sequence in JSON string: %s", us)
				***REMOVED***

				rn := getu4(us)

				// If the rune we just decoded is the high or low value of a possible surrogate pair,
				// try to decode the next sequence as the low value of a surrogate pair. We're
				// expecting the next sequence to be another Unicode escape sequence (e.g. "\uDD1E"),
				// but need to handle cases where the input is not a valid surrogate pair.
				// For more context on unicode surrogate pairs, see:
				// https://www.christianfscott.com/rust-chars-vs-go-runes/
				// https://www.unicode.org/glossary/#high_surrogate_code_point
				if utf16.IsSurrogate(rn) ***REMOVED***
					c, err = js.readNextByte()
					if err != nil ***REMOVED***
						if err == io.EOF ***REMOVED***
							return nil, errors.New("end of input in JSON string")
						***REMOVED***
						return nil, err
					***REMOVED***

					// If the next value isn't the beginning of a backslash escape sequence, write
					// the Unicode replacement character for the surrogate value and goto the
					// beginning of the next char eval block.
					if c != '\\' ***REMOVED***
						b.WriteRune(unicode.ReplacementChar)
						goto evalNextChar
					***REMOVED***

					c, err = js.readNextByte()
					if err != nil ***REMOVED***
						if err == io.EOF ***REMOVED***
							return nil, errors.New("end of input in JSON string")
						***REMOVED***
						return nil, err
					***REMOVED***

					// If the next value isn't the beginning of a unicode escape sequence, write the
					// Unicode replacement character for the surrogate value and goto the beginning
					// of the next escape char eval block.
					if c != 'u' ***REMOVED***
						b.WriteRune(unicode.ReplacementChar)
						goto evalNextEscapeChar
					***REMOVED***

					err = js.readNNextBytes(us, 4, 0)
					if err != nil ***REMOVED***
						return nil, fmt.Errorf("invalid unicode sequence in JSON string: %s", us)
					***REMOVED***

					rn2 := getu4(us)

					// Try to decode the pair of runes as a utf16 surrogate pair. If that fails, write
					// the Unicode replacement character for the surrogate value and the 2nd decoded rune.
					if rnPair := utf16.DecodeRune(rn, rn2); rnPair != unicode.ReplacementChar ***REMOVED***
						b.WriteRune(rnPair)
					***REMOVED*** else ***REMOVED***
						b.WriteRune(unicode.ReplacementChar)
						b.WriteRune(rn2)
					***REMOVED***

					break
				***REMOVED***

				b.WriteRune(rn)
			default:
				return nil, fmt.Errorf("invalid escape sequence in JSON string '\\%c'", c)
			***REMOVED***
		case '"':
			return &jsonToken***REMOVED***t: jttString, v: b.String(), p: p***REMOVED***, nil
		default:
			b.WriteByte(c)
		***REMOVED***
	***REMOVED***
***REMOVED***

// scanLiteral reads an unquoted sequence of characters and determines if it is one of
// three valid JSON literals (true, false, null); if so, it returns the appropriate
// jsonToken; otherwise, it returns an error
func (js *jsonScanner) scanLiteral(first byte) (*jsonToken, error) ***REMOVED***
	p := js.pos - 1

	lit := make([]byte, 4)
	lit[0] = first

	err := js.readNNextBytes(lit, 3, 1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	c5, err := js.readNextByte()

	if bytes.Equal([]byte("true"), lit) && (isValueTerminator(c5) || err == io.EOF) ***REMOVED***
		js.pos = int(math.Max(0, float64(js.pos-1)))
		return &jsonToken***REMOVED***t: jttBool, v: true, p: p***REMOVED***, nil
	***REMOVED*** else if bytes.Equal([]byte("null"), lit) && (isValueTerminator(c5) || err == io.EOF) ***REMOVED***
		js.pos = int(math.Max(0, float64(js.pos-1)))
		return &jsonToken***REMOVED***t: jttNull, v: nil, p: p***REMOVED***, nil
	***REMOVED*** else if bytes.Equal([]byte("fals"), lit) ***REMOVED***
		if c5 == 'e' ***REMOVED***
			c5, err = js.readNextByte()

			if isValueTerminator(c5) || err == io.EOF ***REMOVED***
				js.pos = int(math.Max(0, float64(js.pos-1)))
				return &jsonToken***REMOVED***t: jttBool, v: false, p: p***REMOVED***, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil, fmt.Errorf("invalid JSON literal. Position: %d, literal: %s", p, lit)
***REMOVED***

type numberScanState byte

const (
	nssSawLeadingMinus numberScanState = iota
	nssSawLeadingZero
	nssSawIntegerDigits
	nssSawDecimalPoint
	nssSawFractionDigits
	nssSawExponentLetter
	nssSawExponentSign
	nssSawExponentDigits
	nssDone
	nssInvalid
)

// scanNumber reads a JSON number (according to RFC-8259)
func (js *jsonScanner) scanNumber(first byte) (*jsonToken, error) ***REMOVED***
	var b bytes.Buffer
	var s numberScanState
	var c byte
	var err error

	t := jttInt64 // assume it's an int64 until the type can be determined
	start := js.pos - 1

	b.WriteByte(first)

	switch first ***REMOVED***
	case '-':
		s = nssSawLeadingMinus
	case '0':
		s = nssSawLeadingZero
	default:
		s = nssSawIntegerDigits
	***REMOVED***

	for ***REMOVED***
		c, err = js.readNextByte()

		if err != nil && err != io.EOF ***REMOVED***
			return nil, err
		***REMOVED***

		switch s ***REMOVED***
		case nssSawLeadingMinus:
			switch c ***REMOVED***
			case '0':
				s = nssSawLeadingZero
				b.WriteByte(c)
			default:
				if isDigit(c) ***REMOVED***
					s = nssSawIntegerDigits
					b.WriteByte(c)
				***REMOVED*** else ***REMOVED***
					s = nssInvalid
				***REMOVED***
			***REMOVED***
		case nssSawLeadingZero:
			switch c ***REMOVED***
			case '.':
				s = nssSawDecimalPoint
				b.WriteByte(c)
			case 'e', 'E':
				s = nssSawExponentLetter
				b.WriteByte(c)
			case '***REMOVED***', ']', ',':
				s = nssDone
			default:
				if isWhiteSpace(c) || err == io.EOF ***REMOVED***
					s = nssDone
				***REMOVED*** else ***REMOVED***
					s = nssInvalid
				***REMOVED***
			***REMOVED***
		case nssSawIntegerDigits:
			switch c ***REMOVED***
			case '.':
				s = nssSawDecimalPoint
				b.WriteByte(c)
			case 'e', 'E':
				s = nssSawExponentLetter
				b.WriteByte(c)
			case '***REMOVED***', ']', ',':
				s = nssDone
			default:
				if isWhiteSpace(c) || err == io.EOF ***REMOVED***
					s = nssDone
				***REMOVED*** else if isDigit(c) ***REMOVED***
					s = nssSawIntegerDigits
					b.WriteByte(c)
				***REMOVED*** else ***REMOVED***
					s = nssInvalid
				***REMOVED***
			***REMOVED***
		case nssSawDecimalPoint:
			t = jttDouble
			if isDigit(c) ***REMOVED***
				s = nssSawFractionDigits
				b.WriteByte(c)
			***REMOVED*** else ***REMOVED***
				s = nssInvalid
			***REMOVED***
		case nssSawFractionDigits:
			switch c ***REMOVED***
			case 'e', 'E':
				s = nssSawExponentLetter
				b.WriteByte(c)
			case '***REMOVED***', ']', ',':
				s = nssDone
			default:
				if isWhiteSpace(c) || err == io.EOF ***REMOVED***
					s = nssDone
				***REMOVED*** else if isDigit(c) ***REMOVED***
					s = nssSawFractionDigits
					b.WriteByte(c)
				***REMOVED*** else ***REMOVED***
					s = nssInvalid
				***REMOVED***
			***REMOVED***
		case nssSawExponentLetter:
			t = jttDouble
			switch c ***REMOVED***
			case '+', '-':
				s = nssSawExponentSign
				b.WriteByte(c)
			default:
				if isDigit(c) ***REMOVED***
					s = nssSawExponentDigits
					b.WriteByte(c)
				***REMOVED*** else ***REMOVED***
					s = nssInvalid
				***REMOVED***
			***REMOVED***
		case nssSawExponentSign:
			if isDigit(c) ***REMOVED***
				s = nssSawExponentDigits
				b.WriteByte(c)
			***REMOVED*** else ***REMOVED***
				s = nssInvalid
			***REMOVED***
		case nssSawExponentDigits:
			switch c ***REMOVED***
			case '***REMOVED***', ']', ',':
				s = nssDone
			default:
				if isWhiteSpace(c) || err == io.EOF ***REMOVED***
					s = nssDone
				***REMOVED*** else if isDigit(c) ***REMOVED***
					s = nssSawExponentDigits
					b.WriteByte(c)
				***REMOVED*** else ***REMOVED***
					s = nssInvalid
				***REMOVED***
			***REMOVED***
		***REMOVED***

		switch s ***REMOVED***
		case nssInvalid:
			return nil, fmt.Errorf("invalid JSON number. Position: %d", start)
		case nssDone:
			js.pos = int(math.Max(0, float64(js.pos-1)))
			if t != jttDouble ***REMOVED***
				v, err := strconv.ParseInt(b.String(), 10, 64)
				if err == nil ***REMOVED***
					if v < math.MinInt32 || v > math.MaxInt32 ***REMOVED***
						return &jsonToken***REMOVED***t: jttInt64, v: v, p: start***REMOVED***, nil
					***REMOVED***

					return &jsonToken***REMOVED***t: jttInt32, v: int32(v), p: start***REMOVED***, nil
				***REMOVED***
			***REMOVED***

			v, err := strconv.ParseFloat(b.String(), 64)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			return &jsonToken***REMOVED***t: jttDouble, v: v, p: start***REMOVED***, nil
		***REMOVED***
	***REMOVED***
***REMOVED***
