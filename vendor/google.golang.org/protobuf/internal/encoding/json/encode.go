// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"math"
	"math/bits"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/detrand"
	"google.golang.org/protobuf/internal/errors"
)

// kind represents an encoding type.
type kind uint8

const (
	_ kind = (1 << iota) / 2
	name
	scalar
	objectOpen
	objectClose
	arrayOpen
	arrayClose
)

// Encoder provides methods to write out JSON constructs and values. The user is
// responsible for producing valid sequences of JSON constructs and values.
type Encoder struct ***REMOVED***
	indent   string
	lastKind kind
	indents  []byte
	out      []byte
***REMOVED***

// NewEncoder returns an Encoder.
//
// If indent is a non-empty string, it causes every entry for an Array or Object
// to be preceded by the indent and trailed by a newline.
func NewEncoder(indent string) (*Encoder, error) ***REMOVED***
	e := &Encoder***REMOVED******REMOVED***
	if len(indent) > 0 ***REMOVED***
		if strings.Trim(indent, " \t") != "" ***REMOVED***
			return nil, errors.New("indent may only be composed of space or tab characters")
		***REMOVED***
		e.indent = indent
	***REMOVED***
	return e, nil
***REMOVED***

// Bytes returns the content of the written bytes.
func (e *Encoder) Bytes() []byte ***REMOVED***
	return e.out
***REMOVED***

// WriteNull writes out the null value.
func (e *Encoder) WriteNull() ***REMOVED***
	e.prepareNext(scalar)
	e.out = append(e.out, "null"...)
***REMOVED***

// WriteBool writes out the given boolean value.
func (e *Encoder) WriteBool(b bool) ***REMOVED***
	e.prepareNext(scalar)
	if b ***REMOVED***
		e.out = append(e.out, "true"...)
	***REMOVED*** else ***REMOVED***
		e.out = append(e.out, "false"...)
	***REMOVED***
***REMOVED***

// WriteString writes out the given string in JSON string value. Returns error
// if input string contains invalid UTF-8.
func (e *Encoder) WriteString(s string) error ***REMOVED***
	e.prepareNext(scalar)
	var err error
	if e.out, err = appendString(e.out, s); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Sentinel error used for indicating invalid UTF-8.
var errInvalidUTF8 = errors.New("invalid UTF-8")

func appendString(out []byte, in string) ([]byte, error) ***REMOVED***
	out = append(out, '"')
	i := indexNeedEscapeInString(in)
	in, out = in[i:], append(out, in[:i]...)
	for len(in) > 0 ***REMOVED***
		switch r, n := utf8.DecodeRuneInString(in); ***REMOVED***
		case r == utf8.RuneError && n == 1:
			return out, errInvalidUTF8
		case r < ' ' || r == '"' || r == '\\':
			out = append(out, '\\')
			switch r ***REMOVED***
			case '"', '\\':
				out = append(out, byte(r))
			case '\b':
				out = append(out, 'b')
			case '\f':
				out = append(out, 'f')
			case '\n':
				out = append(out, 'n')
			case '\r':
				out = append(out, 'r')
			case '\t':
				out = append(out, 't')
			default:
				out = append(out, 'u')
				out = append(out, "0000"[1+(bits.Len32(uint32(r))-1)/4:]...)
				out = strconv.AppendUint(out, uint64(r), 16)
			***REMOVED***
			in = in[n:]
		default:
			i := indexNeedEscapeInString(in[n:])
			in, out = in[n+i:], append(out, in[:n+i]...)
		***REMOVED***
	***REMOVED***
	out = append(out, '"')
	return out, nil
***REMOVED***

// indexNeedEscapeInString returns the index of the character that needs
// escaping. If no characters need escaping, this returns the input length.
func indexNeedEscapeInString(s string) int ***REMOVED***
	for i, r := range s ***REMOVED***
		if r < ' ' || r == '\\' || r == '"' || r == utf8.RuneError ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return len(s)
***REMOVED***

// WriteFloat writes out the given float and bitSize in JSON number value.
func (e *Encoder) WriteFloat(n float64, bitSize int) ***REMOVED***
	e.prepareNext(scalar)
	e.out = appendFloat(e.out, n, bitSize)
***REMOVED***

// appendFloat formats given float in bitSize, and appends to the given []byte.
func appendFloat(out []byte, n float64, bitSize int) []byte ***REMOVED***
	switch ***REMOVED***
	case math.IsNaN(n):
		return append(out, `"NaN"`...)
	case math.IsInf(n, +1):
		return append(out, `"Infinity"`...)
	case math.IsInf(n, -1):
		return append(out, `"-Infinity"`...)
	***REMOVED***

	// JSON number formatting logic based on encoding/json.
	// See floatEncoder.encode for reference.
	fmt := byte('f')
	if abs := math.Abs(n); abs != 0 ***REMOVED***
		if bitSize == 64 && (abs < 1e-6 || abs >= 1e21) ||
			bitSize == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) ***REMOVED***
			fmt = 'e'
		***REMOVED***
	***REMOVED***
	out = strconv.AppendFloat(out, n, fmt, -1, bitSize)
	if fmt == 'e' ***REMOVED***
		n := len(out)
		if n >= 4 && out[n-4] == 'e' && out[n-3] == '-' && out[n-2] == '0' ***REMOVED***
			out[n-2] = out[n-1]
			out = out[:n-1]
		***REMOVED***
	***REMOVED***
	return out
***REMOVED***

// WriteInt writes out the given signed integer in JSON number value.
func (e *Encoder) WriteInt(n int64) ***REMOVED***
	e.prepareNext(scalar)
	e.out = append(e.out, strconv.FormatInt(n, 10)...)
***REMOVED***

// WriteUint writes out the given unsigned integer in JSON number value.
func (e *Encoder) WriteUint(n uint64) ***REMOVED***
	e.prepareNext(scalar)
	e.out = append(e.out, strconv.FormatUint(n, 10)...)
***REMOVED***

// StartObject writes out the '***REMOVED***' symbol.
func (e *Encoder) StartObject() ***REMOVED***
	e.prepareNext(objectOpen)
	e.out = append(e.out, '***REMOVED***')
***REMOVED***

// EndObject writes out the '***REMOVED***' symbol.
func (e *Encoder) EndObject() ***REMOVED***
	e.prepareNext(objectClose)
	e.out = append(e.out, '***REMOVED***')
***REMOVED***

// WriteName writes out the given string in JSON string value and the name
// separator ':'. Returns error if input string contains invalid UTF-8, which
// should not be likely as protobuf field names should be valid.
func (e *Encoder) WriteName(s string) error ***REMOVED***
	e.prepareNext(name)
	var err error
	// Append to output regardless of error.
	e.out, err = appendString(e.out, s)
	e.out = append(e.out, ':')
	return err
***REMOVED***

// StartArray writes out the '[' symbol.
func (e *Encoder) StartArray() ***REMOVED***
	e.prepareNext(arrayOpen)
	e.out = append(e.out, '[')
***REMOVED***

// EndArray writes out the ']' symbol.
func (e *Encoder) EndArray() ***REMOVED***
	e.prepareNext(arrayClose)
	e.out = append(e.out, ']')
***REMOVED***

// prepareNext adds possible comma and indentation for the next value based
// on last type and indent option. It also updates lastKind to next.
func (e *Encoder) prepareNext(next kind) ***REMOVED***
	defer func() ***REMOVED***
		// Set lastKind to next.
		e.lastKind = next
	***REMOVED***()

	if len(e.indent) == 0 ***REMOVED***
		// Need to add comma on the following condition.
		if e.lastKind&(scalar|objectClose|arrayClose) != 0 &&
			next&(name|scalar|objectOpen|arrayOpen) != 0 ***REMOVED***
			e.out = append(e.out, ',')
			// For single-line output, add a random extra space after each
			// comma to make output unstable.
			if detrand.Bool() ***REMOVED***
				e.out = append(e.out, ' ')
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	switch ***REMOVED***
	case e.lastKind&(objectOpen|arrayOpen) != 0:
		// If next type is NOT closing, add indent and newline.
		if next&(objectClose|arrayClose) == 0 ***REMOVED***
			e.indents = append(e.indents, e.indent...)
			e.out = append(e.out, '\n')
			e.out = append(e.out, e.indents...)
		***REMOVED***

	case e.lastKind&(scalar|objectClose|arrayClose) != 0:
		switch ***REMOVED***
		// If next type is either a value or name, add comma and newline.
		case next&(name|scalar|objectOpen|arrayOpen) != 0:
			e.out = append(e.out, ',', '\n')

		// If next type is a closing object or array, adjust indentation.
		case next&(objectClose|arrayClose) != 0:
			e.indents = e.indents[:len(e.indents)-len(e.indent)]
			e.out = append(e.out, '\n')
		***REMOVED***
		e.out = append(e.out, e.indents...)

	case e.lastKind&name != 0:
		e.out = append(e.out, ' ')
		// For multi-line output, add a random extra space after key: to make
		// output unstable.
		if detrand.Bool() ***REMOVED***
			e.out = append(e.out, ' ')
		***REMOVED***
	***REMOVED***
***REMOVED***
