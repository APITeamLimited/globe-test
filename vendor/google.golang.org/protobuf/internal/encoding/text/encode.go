// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"math"
	"math/bits"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/detrand"
	"google.golang.org/protobuf/internal/errors"
)

// encType represents an encoding type.
type encType uint8

const (
	_ encType = (1 << iota) / 2
	name
	scalar
	messageOpen
	messageClose
)

// Encoder provides methods to write out textproto constructs and values. The user is
// responsible for producing valid sequences of constructs and values.
type Encoder struct ***REMOVED***
	encoderState

	indent      string
	delims      [2]byte
	outputASCII bool
***REMOVED***

type encoderState struct ***REMOVED***
	lastType encType
	indents  []byte
	out      []byte
***REMOVED***

// NewEncoder returns an Encoder.
//
// If indent is a non-empty string, it causes every entry in a List or Message
// to be preceded by the indent and trailed by a newline.
//
// If delims is not the zero value, it controls the delimiter characters used
// for messages (e.g., "***REMOVED******REMOVED***" vs "<>").
//
// If outputASCII is true, strings will be serialized in such a way that
// multi-byte UTF-8 sequences are escaped. This property ensures that the
// overall output is ASCII (as opposed to UTF-8).
func NewEncoder(indent string, delims [2]byte, outputASCII bool) (*Encoder, error) ***REMOVED***
	e := &Encoder***REMOVED******REMOVED***
	if len(indent) > 0 ***REMOVED***
		if strings.Trim(indent, " \t") != "" ***REMOVED***
			return nil, errors.New("indent may only be composed of space and tab characters")
		***REMOVED***
		e.indent = indent
	***REMOVED***
	switch delims ***REMOVED***
	case [2]byte***REMOVED***0, 0***REMOVED***:
		e.delims = [2]byte***REMOVED***'***REMOVED***', '***REMOVED***'***REMOVED***
	case [2]byte***REMOVED***'***REMOVED***', '***REMOVED***'***REMOVED***, [2]byte***REMOVED***'<', '>'***REMOVED***:
		e.delims = delims
	default:
		return nil, errors.New("delimiters may only be \"***REMOVED******REMOVED***\" or \"<>\"")
	***REMOVED***
	e.outputASCII = outputASCII

	return e, nil
***REMOVED***

// Bytes returns the content of the written bytes.
func (e *Encoder) Bytes() []byte ***REMOVED***
	return e.out
***REMOVED***

// StartMessage writes out the '***REMOVED***' or '<' symbol.
func (e *Encoder) StartMessage() ***REMOVED***
	e.prepareNext(messageOpen)
	e.out = append(e.out, e.delims[0])
***REMOVED***

// EndMessage writes out the '***REMOVED***' or '>' symbol.
func (e *Encoder) EndMessage() ***REMOVED***
	e.prepareNext(messageClose)
	e.out = append(e.out, e.delims[1])
***REMOVED***

// WriteName writes out the field name and the separator ':'.
func (e *Encoder) WriteName(s string) ***REMOVED***
	e.prepareNext(name)
	e.out = append(e.out, s...)
	e.out = append(e.out, ':')
***REMOVED***

// WriteBool writes out the given boolean value.
func (e *Encoder) WriteBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.WriteLiteral("true")
	***REMOVED*** else ***REMOVED***
		e.WriteLiteral("false")
	***REMOVED***
***REMOVED***

// WriteString writes out the given string value.
func (e *Encoder) WriteString(s string) ***REMOVED***
	e.prepareNext(scalar)
	e.out = appendString(e.out, s, e.outputASCII)
***REMOVED***

func appendString(out []byte, in string, outputASCII bool) []byte ***REMOVED***
	out = append(out, '"')
	i := indexNeedEscapeInString(in)
	in, out = in[i:], append(out, in[:i]...)
	for len(in) > 0 ***REMOVED***
		switch r, n := utf8.DecodeRuneInString(in); ***REMOVED***
		case r == utf8.RuneError && n == 1:
			// We do not report invalid UTF-8 because strings in the text format
			// are used to represent both the proto string and bytes type.
			r = rune(in[0])
			fallthrough
		case r < ' ' || r == '"' || r == '\\' || r == 0x7f:
			out = append(out, '\\')
			switch r ***REMOVED***
			case '"', '\\':
				out = append(out, byte(r))
			case '\n':
				out = append(out, 'n')
			case '\r':
				out = append(out, 'r')
			case '\t':
				out = append(out, 't')
			default:
				out = append(out, 'x')
				out = append(out, "00"[1+(bits.Len32(uint32(r))-1)/4:]...)
				out = strconv.AppendUint(out, uint64(r), 16)
			***REMOVED***
			in = in[n:]
		case r >= utf8.RuneSelf && (outputASCII || r <= 0x009f):
			out = append(out, '\\')
			if r <= math.MaxUint16 ***REMOVED***
				out = append(out, 'u')
				out = append(out, "0000"[1+(bits.Len32(uint32(r))-1)/4:]...)
				out = strconv.AppendUint(out, uint64(r), 16)
			***REMOVED*** else ***REMOVED***
				out = append(out, 'U')
				out = append(out, "00000000"[1+(bits.Len32(uint32(r))-1)/4:]...)
				out = strconv.AppendUint(out, uint64(r), 16)
			***REMOVED***
			in = in[n:]
		default:
			i := indexNeedEscapeInString(in[n:])
			in, out = in[n+i:], append(out, in[:n+i]...)
		***REMOVED***
	***REMOVED***
	out = append(out, '"')
	return out
***REMOVED***

// indexNeedEscapeInString returns the index of the character that needs
// escaping. If no characters need escaping, this returns the input length.
func indexNeedEscapeInString(s string) int ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		if c := s[i]; c < ' ' || c == '"' || c == '\'' || c == '\\' || c >= 0x7f ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return len(s)
***REMOVED***

// WriteFloat writes out the given float value for given bitSize.
func (e *Encoder) WriteFloat(n float64, bitSize int) ***REMOVED***
	e.prepareNext(scalar)
	e.out = appendFloat(e.out, n, bitSize)
***REMOVED***

func appendFloat(out []byte, n float64, bitSize int) []byte ***REMOVED***
	switch ***REMOVED***
	case math.IsNaN(n):
		return append(out, "nan"...)
	case math.IsInf(n, +1):
		return append(out, "inf"...)
	case math.IsInf(n, -1):
		return append(out, "-inf"...)
	default:
		return strconv.AppendFloat(out, n, 'g', -1, bitSize)
	***REMOVED***
***REMOVED***

// WriteInt writes out the given signed integer value.
func (e *Encoder) WriteInt(n int64) ***REMOVED***
	e.prepareNext(scalar)
	e.out = append(e.out, strconv.FormatInt(n, 10)...)
***REMOVED***

// WriteUint writes out the given unsigned integer value.
func (e *Encoder) WriteUint(n uint64) ***REMOVED***
	e.prepareNext(scalar)
	e.out = append(e.out, strconv.FormatUint(n, 10)...)
***REMOVED***

// WriteLiteral writes out the given string as a literal value without quotes.
// This is used for writing enum literal strings.
func (e *Encoder) WriteLiteral(s string) ***REMOVED***
	e.prepareNext(scalar)
	e.out = append(e.out, s...)
***REMOVED***

// prepareNext adds possible space and indentation for the next value based
// on last encType and indent option. It also updates e.lastType to next.
func (e *Encoder) prepareNext(next encType) ***REMOVED***
	defer func() ***REMOVED***
		e.lastType = next
	***REMOVED***()

	// Single line.
	if len(e.indent) == 0 ***REMOVED***
		// Add space after each field before the next one.
		if e.lastType&(scalar|messageClose) != 0 && next == name ***REMOVED***
			e.out = append(e.out, ' ')
			// Add a random extra space to make output unstable.
			if detrand.Bool() ***REMOVED***
				e.out = append(e.out, ' ')
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	// Multi-line.
	switch ***REMOVED***
	case e.lastType == name:
		e.out = append(e.out, ' ')
		// Add a random extra space after name: to make output unstable.
		if detrand.Bool() ***REMOVED***
			e.out = append(e.out, ' ')
		***REMOVED***

	case e.lastType == messageOpen && next != messageClose:
		e.indents = append(e.indents, e.indent...)
		e.out = append(e.out, '\n')
		e.out = append(e.out, e.indents...)

	case e.lastType&(scalar|messageClose) != 0:
		if next == messageClose ***REMOVED***
			e.indents = e.indents[:len(e.indents)-len(e.indent)]
		***REMOVED***
		e.out = append(e.out, '\n')
		e.out = append(e.out, e.indents...)
	***REMOVED***
***REMOVED***

// Snapshot returns the current snapshot for use in Reset.
func (e *Encoder) Snapshot() encoderState ***REMOVED***
	return e.encoderState
***REMOVED***

// Reset resets the Encoder to the given encoderState from a Snapshot.
func (e *Encoder) Reset(es encoderState) ***REMOVED***
	e.encoderState = es
***REMOVED***
