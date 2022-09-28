// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"fmt"
	"strconv"
)

// Kind represents a token kind expressible in the JSON format.
type Kind uint16

const (
	Invalid Kind = (1 << iota) / 2
	EOF
	Null
	Bool
	Number
	String
	Name
	ObjectOpen
	ObjectClose
	ArrayOpen
	ArrayClose

	// comma is only for parsing in between tokens and
	// does not need to be exported.
	comma
)

func (k Kind) String() string ***REMOVED***
	switch k ***REMOVED***
	case EOF:
		return "eof"
	case Null:
		return "null"
	case Bool:
		return "bool"
	case Number:
		return "number"
	case String:
		return "string"
	case ObjectOpen:
		return "***REMOVED***"
	case ObjectClose:
		return "***REMOVED***"
	case Name:
		return "name"
	case ArrayOpen:
		return "["
	case ArrayClose:
		return "]"
	case comma:
		return ","
	***REMOVED***
	return "<invalid>"
***REMOVED***

// Token provides a parsed token kind and value.
//
// Values are provided by the difference accessor methods. The accessor methods
// Name, Bool, and ParsedString will panic if called on the wrong kind. There
// are different accessor methods for the Number kind for converting to the
// appropriate Go numeric type and those methods have the ok return value.
type Token struct ***REMOVED***
	// Token kind.
	kind Kind
	// pos provides the position of the token in the original input.
	pos int
	// raw bytes of the serialized token.
	// This is a subslice into the original input.
	raw []byte
	// boo is parsed boolean value.
	boo bool
	// str is parsed string value.
	str string
***REMOVED***

// Kind returns the token kind.
func (t Token) Kind() Kind ***REMOVED***
	return t.kind
***REMOVED***

// RawString returns the read value in string.
func (t Token) RawString() string ***REMOVED***
	return string(t.raw)
***REMOVED***

// Pos returns the token position from the input.
func (t Token) Pos() int ***REMOVED***
	return t.pos
***REMOVED***

// Name returns the object name if token is Name, else it panics.
func (t Token) Name() string ***REMOVED***
	if t.kind == Name ***REMOVED***
		return t.str
	***REMOVED***
	panic(fmt.Sprintf("Token is not a Name: %v", t.RawString()))
***REMOVED***

// Bool returns the bool value if token kind is Bool, else it panics.
func (t Token) Bool() bool ***REMOVED***
	if t.kind == Bool ***REMOVED***
		return t.boo
	***REMOVED***
	panic(fmt.Sprintf("Token is not a Bool: %v", t.RawString()))
***REMOVED***

// ParsedString returns the string value for a JSON string token or the read
// value in string if token is not a string.
func (t Token) ParsedString() string ***REMOVED***
	if t.kind == String ***REMOVED***
		return t.str
	***REMOVED***
	panic(fmt.Sprintf("Token is not a String: %v", t.RawString()))
***REMOVED***

// Float returns the floating-point number if token kind is Number.
//
// The floating-point precision is specified by the bitSize parameter: 32 for
// float32 or 64 for float64. If bitSize=32, the result still has type float64,
// but it will be convertible to float32 without changing its value. It will
// return false if the number exceeds the floating point limits for given
// bitSize.
func (t Token) Float(bitSize int) (float64, bool) ***REMOVED***
	if t.kind != Number ***REMOVED***
		return 0, false
	***REMOVED***
	f, err := strconv.ParseFloat(t.RawString(), bitSize)
	if err != nil ***REMOVED***
		return 0, false
	***REMOVED***
	return f, true
***REMOVED***

// Int returns the signed integer number if token is Number.
//
// The given bitSize specifies the integer type that the result must fit into.
// It returns false if the number is not an integer value or if the result
// exceeds the limits for given bitSize.
func (t Token) Int(bitSize int) (int64, bool) ***REMOVED***
	s, ok := t.getIntStr()
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	n, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil ***REMOVED***
		return 0, false
	***REMOVED***
	return n, true
***REMOVED***

// Uint returns the signed integer number if token is Number.
//
// The given bitSize specifies the unsigned integer type that the result must
// fit into. It returns false if the number is not an unsigned integer value
// or if the result exceeds the limits for given bitSize.
func (t Token) Uint(bitSize int) (uint64, bool) ***REMOVED***
	s, ok := t.getIntStr()
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	n, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil ***REMOVED***
		return 0, false
	***REMOVED***
	return n, true
***REMOVED***

func (t Token) getIntStr() (string, bool) ***REMOVED***
	if t.kind != Number ***REMOVED***
		return "", false
	***REMOVED***
	parts, ok := parseNumberParts(t.raw)
	if !ok ***REMOVED***
		return "", false
	***REMOVED***
	return normalizeToIntString(parts)
***REMOVED***

// TokenEquals returns true if given Tokens are equal, else false.
func TokenEquals(x, y Token) bool ***REMOVED***
	return x.kind == y.kind &&
		x.pos == y.pos &&
		bytes.Equal(x.raw, y.raw) &&
		x.boo == y.boo &&
		x.str == y.str
***REMOVED***
