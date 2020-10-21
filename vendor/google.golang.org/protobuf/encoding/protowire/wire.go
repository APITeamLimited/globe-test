// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protowire parses and formats the raw wire encoding.
// See https://developers.google.com/protocol-buffers/docs/encoding.
//
// For marshaling and unmarshaling entire protobuf messages,
// use the "google.golang.org/protobuf/proto" package instead.
package protowire

import (
	"io"
	"math"
	"math/bits"

	"google.golang.org/protobuf/internal/errors"
)

// Number represents the field number.
type Number int32

const (
	MinValidNumber      Number = 1
	FirstReservedNumber Number = 19000
	LastReservedNumber  Number = 19999
	MaxValidNumber      Number = 1<<29 - 1
)

// IsValid reports whether the field number is semantically valid.
//
// Note that while numbers within the reserved range are semantically invalid,
// they are syntactically valid in the wire format.
// Implementations may treat records with reserved field numbers as unknown.
func (n Number) IsValid() bool ***REMOVED***
	return MinValidNumber <= n && n < FirstReservedNumber || LastReservedNumber < n && n <= MaxValidNumber
***REMOVED***

// Type represents the wire type.
type Type int8

const (
	VarintType     Type = 0
	Fixed32Type    Type = 5
	Fixed64Type    Type = 1
	BytesType      Type = 2
	StartGroupType Type = 3
	EndGroupType   Type = 4
)

const (
	_ = -iota
	errCodeTruncated
	errCodeFieldNumber
	errCodeOverflow
	errCodeReserved
	errCodeEndGroup
)

var (
	errFieldNumber = errors.New("invalid field number")
	errOverflow    = errors.New("variable length integer overflow")
	errReserved    = errors.New("cannot parse reserved wire type")
	errEndGroup    = errors.New("mismatching end group marker")
	errParse       = errors.New("parse error")
)

// ParseError converts an error code into an error value.
// This returns nil if n is a non-negative number.
func ParseError(n int) error ***REMOVED***
	if n >= 0 ***REMOVED***
		return nil
	***REMOVED***
	switch n ***REMOVED***
	case errCodeTruncated:
		return io.ErrUnexpectedEOF
	case errCodeFieldNumber:
		return errFieldNumber
	case errCodeOverflow:
		return errOverflow
	case errCodeReserved:
		return errReserved
	case errCodeEndGroup:
		return errEndGroup
	default:
		return errParse
	***REMOVED***
***REMOVED***

// ConsumeField parses an entire field record (both tag and value) and returns
// the field number, the wire type, and the total length.
// This returns a negative length upon an error (see ParseError).
//
// The total length includes the tag header and the end group marker (if the
// field is a group).
func ConsumeField(b []byte) (Number, Type, int) ***REMOVED***
	num, typ, n := ConsumeTag(b)
	if n < 0 ***REMOVED***
		return 0, 0, n // forward error code
	***REMOVED***
	m := ConsumeFieldValue(num, typ, b[n:])
	if m < 0 ***REMOVED***
		return 0, 0, m // forward error code
	***REMOVED***
	return num, typ, n + m
***REMOVED***

// ConsumeFieldValue parses a field value and returns its length.
// This assumes that the field Number and wire Type have already been parsed.
// This returns a negative length upon an error (see ParseError).
//
// When parsing a group, the length includes the end group marker and
// the end group is verified to match the starting field number.
func ConsumeFieldValue(num Number, typ Type, b []byte) (n int) ***REMOVED***
	switch typ ***REMOVED***
	case VarintType:
		_, n = ConsumeVarint(b)
		return n
	case Fixed32Type:
		_, n = ConsumeFixed32(b)
		return n
	case Fixed64Type:
		_, n = ConsumeFixed64(b)
		return n
	case BytesType:
		_, n = ConsumeBytes(b)
		return n
	case StartGroupType:
		n0 := len(b)
		for ***REMOVED***
			num2, typ2, n := ConsumeTag(b)
			if n < 0 ***REMOVED***
				return n // forward error code
			***REMOVED***
			b = b[n:]
			if typ2 == EndGroupType ***REMOVED***
				if num != num2 ***REMOVED***
					return errCodeEndGroup
				***REMOVED***
				return n0 - len(b)
			***REMOVED***

			n = ConsumeFieldValue(num2, typ2, b)
			if n < 0 ***REMOVED***
				return n // forward error code
			***REMOVED***
			b = b[n:]
		***REMOVED***
	case EndGroupType:
		return errCodeEndGroup
	default:
		return errCodeReserved
	***REMOVED***
***REMOVED***

// AppendTag encodes num and typ as a varint-encoded tag and appends it to b.
func AppendTag(b []byte, num Number, typ Type) []byte ***REMOVED***
	return AppendVarint(b, EncodeTag(num, typ))
***REMOVED***

// ConsumeTag parses b as a varint-encoded tag, reporting its length.
// This returns a negative length upon an error (see ParseError).
func ConsumeTag(b []byte) (Number, Type, int) ***REMOVED***
	v, n := ConsumeVarint(b)
	if n < 0 ***REMOVED***
		return 0, 0, n // forward error code
	***REMOVED***
	num, typ := DecodeTag(v)
	if num < MinValidNumber ***REMOVED***
		return 0, 0, errCodeFieldNumber
	***REMOVED***
	return num, typ, n
***REMOVED***

func SizeTag(num Number) int ***REMOVED***
	return SizeVarint(EncodeTag(num, 0)) // wire type has no effect on size
***REMOVED***

// AppendVarint appends v to b as a varint-encoded uint64.
func AppendVarint(b []byte, v uint64) []byte ***REMOVED***
	switch ***REMOVED***
	case v < 1<<7:
		b = append(b, byte(v))
	case v < 1<<14:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte(v>>7))
	case v < 1<<21:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte(v>>14))
	case v < 1<<28:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte(v>>21))
	case v < 1<<35:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte(v>>28))
	case v < 1<<42:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte(v>>35))
	case v < 1<<49:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte(v>>42))
	case v < 1<<56:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte(v>>49))
	case v < 1<<63:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte(v>>56))
	default:
		b = append(b,
			byte((v>>0)&0x7f|0x80),
			byte((v>>7)&0x7f|0x80),
			byte((v>>14)&0x7f|0x80),
			byte((v>>21)&0x7f|0x80),
			byte((v>>28)&0x7f|0x80),
			byte((v>>35)&0x7f|0x80),
			byte((v>>42)&0x7f|0x80),
			byte((v>>49)&0x7f|0x80),
			byte((v>>56)&0x7f|0x80),
			1)
	***REMOVED***
	return b
***REMOVED***

// ConsumeVarint parses b as a varint-encoded uint64, reporting its length.
// This returns a negative length upon an error (see ParseError).
func ConsumeVarint(b []byte) (v uint64, n int) ***REMOVED***
	var y uint64
	if len(b) <= 0 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	v = uint64(b[0])
	if v < 0x80 ***REMOVED***
		return v, 1
	***REMOVED***
	v -= 0x80

	if len(b) <= 1 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[1])
	v += y << 7
	if y < 0x80 ***REMOVED***
		return v, 2
	***REMOVED***
	v -= 0x80 << 7

	if len(b) <= 2 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[2])
	v += y << 14
	if y < 0x80 ***REMOVED***
		return v, 3
	***REMOVED***
	v -= 0x80 << 14

	if len(b) <= 3 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[3])
	v += y << 21
	if y < 0x80 ***REMOVED***
		return v, 4
	***REMOVED***
	v -= 0x80 << 21

	if len(b) <= 4 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[4])
	v += y << 28
	if y < 0x80 ***REMOVED***
		return v, 5
	***REMOVED***
	v -= 0x80 << 28

	if len(b) <= 5 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[5])
	v += y << 35
	if y < 0x80 ***REMOVED***
		return v, 6
	***REMOVED***
	v -= 0x80 << 35

	if len(b) <= 6 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[6])
	v += y << 42
	if y < 0x80 ***REMOVED***
		return v, 7
	***REMOVED***
	v -= 0x80 << 42

	if len(b) <= 7 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[7])
	v += y << 49
	if y < 0x80 ***REMOVED***
		return v, 8
	***REMOVED***
	v -= 0x80 << 49

	if len(b) <= 8 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[8])
	v += y << 56
	if y < 0x80 ***REMOVED***
		return v, 9
	***REMOVED***
	v -= 0x80 << 56

	if len(b) <= 9 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	y = uint64(b[9])
	v += y << 63
	if y < 2 ***REMOVED***
		return v, 10
	***REMOVED***
	return 0, errCodeOverflow
***REMOVED***

// SizeVarint returns the encoded size of a varint.
// The size is guaranteed to be within 1 and 10, inclusive.
func SizeVarint(v uint64) int ***REMOVED***
	// This computes 1 + (bits.Len64(v)-1)/7.
	// 9/64 is a good enough approximation of 1/7
	return int(9*uint32(bits.Len64(v))+64) / 64
***REMOVED***

// AppendFixed32 appends v to b as a little-endian uint32.
func AppendFixed32(b []byte, v uint32) []byte ***REMOVED***
	return append(b,
		byte(v>>0),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24))
***REMOVED***

// ConsumeFixed32 parses b as a little-endian uint32, reporting its length.
// This returns a negative length upon an error (see ParseError).
func ConsumeFixed32(b []byte) (v uint32, n int) ***REMOVED***
	if len(b) < 4 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	v = uint32(b[0])<<0 | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	return v, 4
***REMOVED***

// SizeFixed32 returns the encoded size of a fixed32; which is always 4.
func SizeFixed32() int ***REMOVED***
	return 4
***REMOVED***

// AppendFixed64 appends v to b as a little-endian uint64.
func AppendFixed64(b []byte, v uint64) []byte ***REMOVED***
	return append(b,
		byte(v>>0),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56))
***REMOVED***

// ConsumeFixed64 parses b as a little-endian uint64, reporting its length.
// This returns a negative length upon an error (see ParseError).
func ConsumeFixed64(b []byte) (v uint64, n int) ***REMOVED***
	if len(b) < 8 ***REMOVED***
		return 0, errCodeTruncated
	***REMOVED***
	v = uint64(b[0])<<0 | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	return v, 8
***REMOVED***

// SizeFixed64 returns the encoded size of a fixed64; which is always 8.
func SizeFixed64() int ***REMOVED***
	return 8
***REMOVED***

// AppendBytes appends v to b as a length-prefixed bytes value.
func AppendBytes(b []byte, v []byte) []byte ***REMOVED***
	return append(AppendVarint(b, uint64(len(v))), v...)
***REMOVED***

// ConsumeBytes parses b as a length-prefixed bytes value, reporting its length.
// This returns a negative length upon an error (see ParseError).
func ConsumeBytes(b []byte) (v []byte, n int) ***REMOVED***
	m, n := ConsumeVarint(b)
	if n < 0 ***REMOVED***
		return nil, n // forward error code
	***REMOVED***
	if m > uint64(len(b[n:])) ***REMOVED***
		return nil, errCodeTruncated
	***REMOVED***
	return b[n:][:m], n + int(m)
***REMOVED***

// SizeBytes returns the encoded size of a length-prefixed bytes value,
// given only the length.
func SizeBytes(n int) int ***REMOVED***
	return SizeVarint(uint64(n)) + n
***REMOVED***

// AppendString appends v to b as a length-prefixed bytes value.
func AppendString(b []byte, v string) []byte ***REMOVED***
	return append(AppendVarint(b, uint64(len(v))), v...)
***REMOVED***

// ConsumeString parses b as a length-prefixed bytes value, reporting its length.
// This returns a negative length upon an error (see ParseError).
func ConsumeString(b []byte) (v string, n int) ***REMOVED***
	bb, n := ConsumeBytes(b)
	return string(bb), n
***REMOVED***

// AppendGroup appends v to b as group value, with a trailing end group marker.
// The value v must not contain the end marker.
func AppendGroup(b []byte, num Number, v []byte) []byte ***REMOVED***
	return AppendVarint(append(b, v...), EncodeTag(num, EndGroupType))
***REMOVED***

// ConsumeGroup parses b as a group value until the trailing end group marker,
// and verifies that the end marker matches the provided num. The value v
// does not contain the end marker, while the length does contain the end marker.
// This returns a negative length upon an error (see ParseError).
func ConsumeGroup(num Number, b []byte) (v []byte, n int) ***REMOVED***
	n = ConsumeFieldValue(num, StartGroupType, b)
	if n < 0 ***REMOVED***
		return nil, n // forward error code
	***REMOVED***
	b = b[:n]

	// Truncate off end group marker, but need to handle denormalized varints.
	// Assuming end marker is never 0 (which is always the case since
	// EndGroupType is non-zero), we can truncate all trailing bytes where the
	// lower 7 bits are all zero (implying that the varint is denormalized).
	for len(b) > 0 && b[len(b)-1]&0x7f == 0 ***REMOVED***
		b = b[:len(b)-1]
	***REMOVED***
	b = b[:len(b)-SizeTag(num)]
	return b, n
***REMOVED***

// SizeGroup returns the encoded size of a group, given only the length.
func SizeGroup(num Number, n int) int ***REMOVED***
	return n + SizeTag(num)
***REMOVED***

// DecodeTag decodes the field Number and wire Type from its unified form.
// The Number is -1 if the decoded field number overflows int32.
// Other than overflow, this does not check for field number validity.
func DecodeTag(x uint64) (Number, Type) ***REMOVED***
	// NOTE: MessageSet allows for larger field numbers than normal.
	if x>>3 > uint64(math.MaxInt32) ***REMOVED***
		return -1, 0
	***REMOVED***
	return Number(x >> 3), Type(x & 7)
***REMOVED***

// EncodeTag encodes the field Number and wire Type into its unified form.
func EncodeTag(num Number, typ Type) uint64 ***REMOVED***
	return uint64(num)<<3 | uint64(typ&7)
***REMOVED***

// DecodeZigZag decodes a zig-zag-encoded uint64 as an int64.
//	Input:  ***REMOVED***…,  5,  3,  1,  0,  2,  4,  6, …***REMOVED***
//	Output: ***REMOVED***…, -3, -2, -1,  0, +1, +2, +3, …***REMOVED***
func DecodeZigZag(x uint64) int64 ***REMOVED***
	return int64(x>>1) ^ int64(x)<<63>>63
***REMOVED***

// EncodeZigZag encodes an int64 as a zig-zag-encoded uint64.
//	Input:  ***REMOVED***…, -3, -2, -1,  0, +1, +2, +3, …***REMOVED***
//	Output: ***REMOVED***…,  5,  3,  1,  0,  2,  4,  6, …***REMOVED***
func EncodeZigZag(x int64) uint64 ***REMOVED***
	return uint64(x<<1) ^ uint64(x>>63)
***REMOVED***

// DecodeBool decodes a uint64 as a bool.
//	Input:  ***REMOVED***    0,    1,    2, …***REMOVED***
//	Output: ***REMOVED***false, true, true, …***REMOVED***
func DecodeBool(x uint64) bool ***REMOVED***
	return x != 0
***REMOVED***

// EncodeBool encodes a bool as a uint64.
//	Input:  ***REMOVED***false, true***REMOVED***
//	Output: ***REMOVED***    0,    1***REMOVED***
func EncodeBool(x bool) uint64 ***REMOVED***
	if x ***REMOVED***
		return 1
	***REMOVED***
	return 0
***REMOVED***
