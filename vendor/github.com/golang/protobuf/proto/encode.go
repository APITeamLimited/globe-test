// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package proto

/*
 * Routines for encoding data into the wire format for protocol buffers.
 */

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

// RequiredNotSetError is the error returned if Marshal is called with
// a protocol buffer struct whose required fields have not
// all been initialized. It is also the error returned if Unmarshal is
// called with an encoded protocol buffer that does not include all the
// required fields.
//
// When printed, RequiredNotSetError reports the first unset required field in a
// message. If the field cannot be precisely determined, it is reported as
// "***REMOVED***Unknown***REMOVED***".
type RequiredNotSetError struct ***REMOVED***
	field string
***REMOVED***

func (e *RequiredNotSetError) Error() string ***REMOVED***
	return fmt.Sprintf("proto: required field %q not set", e.field)
***REMOVED***

var (
	// errRepeatedHasNil is the error returned if Marshal is called with
	// a struct with a repeated field containing a nil element.
	errRepeatedHasNil = errors.New("proto: repeated field has nil element")

	// errOneofHasNil is the error returned if Marshal is called with
	// a struct with a oneof field containing a nil element.
	errOneofHasNil = errors.New("proto: oneof field has nil value")

	// ErrNil is the error returned if Marshal is called with nil.
	ErrNil = errors.New("proto: Marshal called with nil")

	// ErrTooLarge is the error returned if Marshal is called with a
	// message that encodes to >2GB.
	ErrTooLarge = errors.New("proto: message encodes to over 2 GB")
)

// The fundamental encoders that put bytes on the wire.
// Those that take integer types all accept uint64 and are
// therefore of type valueEncoder.

const maxVarintBytes = 10 // maximum length of a varint

// maxMarshalSize is the largest allowed size of an encoded protobuf,
// since C++ and Java use signed int32s for the size.
const maxMarshalSize = 1<<31 - 1

// EncodeVarint returns the varint encoding of x.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
// Not used by the package itself, but helpful to clients
// wishing to use the same encoding.
func EncodeVarint(x uint64) []byte ***REMOVED***
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ ***REMOVED***
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	***REMOVED***
	buf[n] = uint8(x)
	n++
	return buf[0:n]
***REMOVED***

// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (p *Buffer) EncodeVarint(x uint64) error ***REMOVED***
	for x >= 1<<7 ***REMOVED***
		p.buf = append(p.buf, uint8(x&0x7f|0x80))
		x >>= 7
	***REMOVED***
	p.buf = append(p.buf, uint8(x))
	return nil
***REMOVED***

// SizeVarint returns the varint encoding size of an integer.
func SizeVarint(x uint64) int ***REMOVED***
	return sizeVarint(x)
***REMOVED***

func sizeVarint(x uint64) (n int) ***REMOVED***
	for ***REMOVED***
		n++
		x >>= 7
		if x == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return n
***REMOVED***

// EncodeFixed64 writes a 64-bit integer to the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (p *Buffer) EncodeFixed64(x uint64) error ***REMOVED***
	p.buf = append(p.buf,
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24),
		uint8(x>>32),
		uint8(x>>40),
		uint8(x>>48),
		uint8(x>>56))
	return nil
***REMOVED***

func sizeFixed64(x uint64) int ***REMOVED***
	return 8
***REMOVED***

// EncodeFixed32 writes a 32-bit integer to the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (p *Buffer) EncodeFixed32(x uint64) error ***REMOVED***
	p.buf = append(p.buf,
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24))
	return nil
***REMOVED***

func sizeFixed32(x uint64) int ***REMOVED***
	return 4
***REMOVED***

// EncodeZigzag64 writes a zigzag-encoded 64-bit integer
// to the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (p *Buffer) EncodeZigzag64(x uint64) error ***REMOVED***
	// use signed number to get arithmetic right shift.
	return p.EncodeVarint((x << 1) ^ uint64((int64(x) >> 63)))
***REMOVED***

func sizeZigzag64(x uint64) int ***REMOVED***
	return sizeVarint((x << 1) ^ uint64((int64(x) >> 63)))
***REMOVED***

// EncodeZigzag32 writes a zigzag-encoded 32-bit integer
// to the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (p *Buffer) EncodeZigzag32(x uint64) error ***REMOVED***
	// use signed number to get arithmetic right shift.
	return p.EncodeVarint(uint64((uint32(x) << 1) ^ uint32((int32(x) >> 31))))
***REMOVED***

func sizeZigzag32(x uint64) int ***REMOVED***
	return sizeVarint(uint64((uint32(x) << 1) ^ uint32((int32(x) >> 31))))
***REMOVED***

// EncodeRawBytes writes a count-delimited byte buffer to the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (p *Buffer) EncodeRawBytes(b []byte) error ***REMOVED***
	p.EncodeVarint(uint64(len(b)))
	p.buf = append(p.buf, b...)
	return nil
***REMOVED***

func sizeRawBytes(b []byte) int ***REMOVED***
	return sizeVarint(uint64(len(b))) +
		len(b)
***REMOVED***

// EncodeStringBytes writes an encoded string to the Buffer.
// This is the format used for the proto2 string type.
func (p *Buffer) EncodeStringBytes(s string) error ***REMOVED***
	p.EncodeVarint(uint64(len(s)))
	p.buf = append(p.buf, s...)
	return nil
***REMOVED***

func sizeStringBytes(s string) int ***REMOVED***
	return sizeVarint(uint64(len(s))) +
		len(s)
***REMOVED***

// Marshaler is the interface representing objects that can marshal themselves.
type Marshaler interface ***REMOVED***
	Marshal() ([]byte, error)
***REMOVED***

// Marshal takes the protocol buffer
// and encodes it into the wire format, returning the data.
func Marshal(pb Message) ([]byte, error) ***REMOVED***
	// Can the object marshal itself?
	if m, ok := pb.(Marshaler); ok ***REMOVED***
		return m.Marshal()
	***REMOVED***
	p := NewBuffer(nil)
	err := p.Marshal(pb)
	if p.buf == nil && err == nil ***REMOVED***
		// Return a non-nil slice on success.
		return []byte***REMOVED******REMOVED***, nil
	***REMOVED***
	return p.buf, err
***REMOVED***

// EncodeMessage writes the protocol buffer to the Buffer,
// prefixed by a varint-encoded length.
func (p *Buffer) EncodeMessage(pb Message) error ***REMOVED***
	t, base, err := getbase(pb)
	if structPointer_IsNil(base) ***REMOVED***
		return ErrNil
	***REMOVED***
	if err == nil ***REMOVED***
		var state errorState
		err = p.enc_len_struct(GetProperties(t.Elem()), base, &state)
	***REMOVED***
	return err
***REMOVED***

// Marshal takes the protocol buffer
// and encodes it into the wire format, writing the result to the
// Buffer.
func (p *Buffer) Marshal(pb Message) error ***REMOVED***
	// Can the object marshal itself?
	if m, ok := pb.(Marshaler); ok ***REMOVED***
		data, err := m.Marshal()
		p.buf = append(p.buf, data...)
		return err
	***REMOVED***

	t, base, err := getbase(pb)
	if structPointer_IsNil(base) ***REMOVED***
		return ErrNil
	***REMOVED***
	if err == nil ***REMOVED***
		err = p.enc_struct(GetProperties(t.Elem()), base)
	***REMOVED***

	if collectStats ***REMOVED***
		(stats).Encode++ // Parens are to work around a goimports bug.
	***REMOVED***

	if len(p.buf) > maxMarshalSize ***REMOVED***
		return ErrTooLarge
	***REMOVED***
	return err
***REMOVED***

// Size returns the encoded size of a protocol buffer.
func Size(pb Message) (n int) ***REMOVED***
	// Can the object marshal itself?  If so, Size is slow.
	// TODO: add Size to Marshaler, or add a Sizer interface.
	if m, ok := pb.(Marshaler); ok ***REMOVED***
		b, _ := m.Marshal()
		return len(b)
	***REMOVED***

	t, base, err := getbase(pb)
	if structPointer_IsNil(base) ***REMOVED***
		return 0
	***REMOVED***
	if err == nil ***REMOVED***
		n = size_struct(GetProperties(t.Elem()), base)
	***REMOVED***

	if collectStats ***REMOVED***
		(stats).Size++ // Parens are to work around a goimports bug.
	***REMOVED***

	return
***REMOVED***

// Individual type encoders.

// Encode a bool.
func (o *Buffer) enc_bool(p *Properties, base structPointer) error ***REMOVED***
	v := *structPointer_Bool(base, p.field)
	if v == nil ***REMOVED***
		return ErrNil
	***REMOVED***
	x := 0
	if *v ***REMOVED***
		x = 1
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
***REMOVED***

func (o *Buffer) enc_proto3_bool(p *Properties, base structPointer) error ***REMOVED***
	v := *structPointer_BoolVal(base, p.field)
	if !v ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, 1)
	return nil
***REMOVED***

func size_bool(p *Properties, base structPointer) int ***REMOVED***
	v := *structPointer_Bool(base, p.field)
	if v == nil ***REMOVED***
		return 0
	***REMOVED***
	return len(p.tagcode) + 1 // each bool takes exactly one byte
***REMOVED***

func size_proto3_bool(p *Properties, base structPointer) int ***REMOVED***
	v := *structPointer_BoolVal(base, p.field)
	if !v && !p.oneof ***REMOVED***
		return 0
	***REMOVED***
	return len(p.tagcode) + 1 // each bool takes exactly one byte
***REMOVED***

// Encode an int32.
func (o *Buffer) enc_int32(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word32(base, p.field)
	if word32_IsNil(v) ***REMOVED***
		return ErrNil
	***REMOVED***
	x := int32(word32_Get(v)) // permit sign extension to use full 64-bit range
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
***REMOVED***

func (o *Buffer) enc_proto3_int32(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word32Val(base, p.field)
	x := int32(word32Val_Get(v)) // permit sign extension to use full 64-bit range
	if x == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
***REMOVED***

func size_int32(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word32(base, p.field)
	if word32_IsNil(v) ***REMOVED***
		return 0
	***REMOVED***
	x := int32(word32_Get(v)) // permit sign extension to use full 64-bit range
	n += len(p.tagcode)
	n += p.valSize(uint64(x))
	return
***REMOVED***

func size_proto3_int32(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word32Val(base, p.field)
	x := int32(word32Val_Get(v)) // permit sign extension to use full 64-bit range
	if x == 0 && !p.oneof ***REMOVED***
		return 0
	***REMOVED***
	n += len(p.tagcode)
	n += p.valSize(uint64(x))
	return
***REMOVED***

// Encode a uint32.
// Exactly the same as int32, except for no sign extension.
func (o *Buffer) enc_uint32(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word32(base, p.field)
	if word32_IsNil(v) ***REMOVED***
		return ErrNil
	***REMOVED***
	x := word32_Get(v)
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
***REMOVED***

func (o *Buffer) enc_proto3_uint32(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word32Val(base, p.field)
	x := word32Val_Get(v)
	if x == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, uint64(x))
	return nil
***REMOVED***

func size_uint32(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word32(base, p.field)
	if word32_IsNil(v) ***REMOVED***
		return 0
	***REMOVED***
	x := word32_Get(v)
	n += len(p.tagcode)
	n += p.valSize(uint64(x))
	return
***REMOVED***

func size_proto3_uint32(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word32Val(base, p.field)
	x := word32Val_Get(v)
	if x == 0 && !p.oneof ***REMOVED***
		return 0
	***REMOVED***
	n += len(p.tagcode)
	n += p.valSize(uint64(x))
	return
***REMOVED***

// Encode an int64.
func (o *Buffer) enc_int64(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word64(base, p.field)
	if word64_IsNil(v) ***REMOVED***
		return ErrNil
	***REMOVED***
	x := word64_Get(v)
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, x)
	return nil
***REMOVED***

func (o *Buffer) enc_proto3_int64(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word64Val(base, p.field)
	x := word64Val_Get(v)
	if x == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	p.valEnc(o, x)
	return nil
***REMOVED***

func size_int64(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word64(base, p.field)
	if word64_IsNil(v) ***REMOVED***
		return 0
	***REMOVED***
	x := word64_Get(v)
	n += len(p.tagcode)
	n += p.valSize(x)
	return
***REMOVED***

func size_proto3_int64(p *Properties, base structPointer) (n int) ***REMOVED***
	v := structPointer_Word64Val(base, p.field)
	x := word64Val_Get(v)
	if x == 0 && !p.oneof ***REMOVED***
		return 0
	***REMOVED***
	n += len(p.tagcode)
	n += p.valSize(x)
	return
***REMOVED***

// Encode a string.
func (o *Buffer) enc_string(p *Properties, base structPointer) error ***REMOVED***
	v := *structPointer_String(base, p.field)
	if v == nil ***REMOVED***
		return ErrNil
	***REMOVED***
	x := *v
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeStringBytes(x)
	return nil
***REMOVED***

func (o *Buffer) enc_proto3_string(p *Properties, base structPointer) error ***REMOVED***
	v := *structPointer_StringVal(base, p.field)
	if v == "" ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeStringBytes(v)
	return nil
***REMOVED***

func size_string(p *Properties, base structPointer) (n int) ***REMOVED***
	v := *structPointer_String(base, p.field)
	if v == nil ***REMOVED***
		return 0
	***REMOVED***
	x := *v
	n += len(p.tagcode)
	n += sizeStringBytes(x)
	return
***REMOVED***

func size_proto3_string(p *Properties, base structPointer) (n int) ***REMOVED***
	v := *structPointer_StringVal(base, p.field)
	if v == "" && !p.oneof ***REMOVED***
		return 0
	***REMOVED***
	n += len(p.tagcode)
	n += sizeStringBytes(v)
	return
***REMOVED***

// All protocol buffer fields are nillable, but be careful.
func isNil(v reflect.Value) bool ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	***REMOVED***
	return false
***REMOVED***

// Encode a message struct.
func (o *Buffer) enc_struct_message(p *Properties, base structPointer) error ***REMOVED***
	var state errorState
	structp := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(structp) ***REMOVED***
		return ErrNil
	***REMOVED***

	// Can the object marshal itself?
	if p.isMarshaler ***REMOVED***
		m := structPointer_Interface(structp, p.stype).(Marshaler)
		data, err := m.Marshal()
		if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
			return err
		***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(data)
		return state.err
	***REMOVED***

	o.buf = append(o.buf, p.tagcode...)
	return o.enc_len_struct(p.sprop, structp, &state)
***REMOVED***

func size_struct_message(p *Properties, base structPointer) int ***REMOVED***
	structp := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(structp) ***REMOVED***
		return 0
	***REMOVED***

	// Can the object marshal itself?
	if p.isMarshaler ***REMOVED***
		m := structPointer_Interface(structp, p.stype).(Marshaler)
		data, _ := m.Marshal()
		n0 := len(p.tagcode)
		n1 := sizeRawBytes(data)
		return n0 + n1
	***REMOVED***

	n0 := len(p.tagcode)
	n1 := size_struct(p.sprop, structp)
	n2 := sizeVarint(uint64(n1)) // size of encoded length
	return n0 + n1 + n2
***REMOVED***

// Encode a group struct.
func (o *Buffer) enc_struct_group(p *Properties, base structPointer) error ***REMOVED***
	var state errorState
	b := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(b) ***REMOVED***
		return ErrNil
	***REMOVED***

	o.EncodeVarint(uint64((p.Tag << 3) | WireStartGroup))
	err := o.enc_struct(p.sprop, b)
	if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
		return err
	***REMOVED***
	o.EncodeVarint(uint64((p.Tag << 3) | WireEndGroup))
	return state.err
***REMOVED***

func size_struct_group(p *Properties, base structPointer) (n int) ***REMOVED***
	b := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(b) ***REMOVED***
		return 0
	***REMOVED***

	n += sizeVarint(uint64((p.Tag << 3) | WireStartGroup))
	n += size_struct(p.sprop, b)
	n += sizeVarint(uint64((p.Tag << 3) | WireEndGroup))
	return
***REMOVED***

// Encode a slice of bools ([]bool).
func (o *Buffer) enc_slice_bool(p *Properties, base structPointer) error ***REMOVED***
	s := *structPointer_BoolSlice(base, p.field)
	l := len(s)
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	for _, x := range s ***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		v := uint64(0)
		if x ***REMOVED***
			v = 1
		***REMOVED***
		p.valEnc(o, v)
	***REMOVED***
	return nil
***REMOVED***

func size_slice_bool(p *Properties, base structPointer) int ***REMOVED***
	s := *structPointer_BoolSlice(base, p.field)
	l := len(s)
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	return l * (len(p.tagcode) + 1) // each bool takes exactly one byte
***REMOVED***

// Encode a slice of bools ([]bool) in packed format.
func (o *Buffer) enc_slice_packed_bool(p *Properties, base structPointer) error ***REMOVED***
	s := *structPointer_BoolSlice(base, p.field)
	l := len(s)
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeVarint(uint64(l)) // each bool takes exactly one byte
	for _, x := range s ***REMOVED***
		v := uint64(0)
		if x ***REMOVED***
			v = 1
		***REMOVED***
		p.valEnc(o, v)
	***REMOVED***
	return nil
***REMOVED***

func size_slice_packed_bool(p *Properties, base structPointer) (n int) ***REMOVED***
	s := *structPointer_BoolSlice(base, p.field)
	l := len(s)
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	n += len(p.tagcode)
	n += sizeVarint(uint64(l))
	n += l // each bool takes exactly one byte
	return
***REMOVED***

// Encode a slice of bytes ([]byte).
func (o *Buffer) enc_slice_byte(p *Properties, base structPointer) error ***REMOVED***
	s := *structPointer_Bytes(base, p.field)
	if s == nil ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(s)
	return nil
***REMOVED***

func (o *Buffer) enc_proto3_slice_byte(p *Properties, base structPointer) error ***REMOVED***
	s := *structPointer_Bytes(base, p.field)
	if len(s) == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	o.buf = append(o.buf, p.tagcode...)
	o.EncodeRawBytes(s)
	return nil
***REMOVED***

func size_slice_byte(p *Properties, base structPointer) (n int) ***REMOVED***
	s := *structPointer_Bytes(base, p.field)
	if s == nil && !p.oneof ***REMOVED***
		return 0
	***REMOVED***
	n += len(p.tagcode)
	n += sizeRawBytes(s)
	return
***REMOVED***

func size_proto3_slice_byte(p *Properties, base structPointer) (n int) ***REMOVED***
	s := *structPointer_Bytes(base, p.field)
	if len(s) == 0 && !p.oneof ***REMOVED***
		return 0
	***REMOVED***
	n += len(p.tagcode)
	n += sizeRawBytes(s)
	return
***REMOVED***

// Encode a slice of int32s ([]int32).
func (o *Buffer) enc_slice_int32(p *Properties, base structPointer) error ***REMOVED***
	s := structPointer_Word32Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	for i := 0; i < l; i++ ***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		x := int32(s.Index(i)) // permit sign extension to use full 64-bit range
		p.valEnc(o, uint64(x))
	***REMOVED***
	return nil
***REMOVED***

func size_slice_int32(p *Properties, base structPointer) (n int) ***REMOVED***
	s := structPointer_Word32Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	for i := 0; i < l; i++ ***REMOVED***
		n += len(p.tagcode)
		x := int32(s.Index(i)) // permit sign extension to use full 64-bit range
		n += p.valSize(uint64(x))
	***REMOVED***
	return
***REMOVED***

// Encode a slice of int32s ([]int32) in packed format.
func (o *Buffer) enc_slice_packed_int32(p *Properties, base structPointer) error ***REMOVED***
	s := structPointer_Word32Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	// TODO: Reuse a Buffer.
	buf := NewBuffer(nil)
	for i := 0; i < l; i++ ***REMOVED***
		x := int32(s.Index(i)) // permit sign extension to use full 64-bit range
		p.valEnc(buf, uint64(x))
	***REMOVED***

	o.buf = append(o.buf, p.tagcode...)
	o.EncodeVarint(uint64(len(buf.buf)))
	o.buf = append(o.buf, buf.buf...)
	return nil
***REMOVED***

func size_slice_packed_int32(p *Properties, base structPointer) (n int) ***REMOVED***
	s := structPointer_Word32Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	var bufSize int
	for i := 0; i < l; i++ ***REMOVED***
		x := int32(s.Index(i)) // permit sign extension to use full 64-bit range
		bufSize += p.valSize(uint64(x))
	***REMOVED***

	n += len(p.tagcode)
	n += sizeVarint(uint64(bufSize))
	n += bufSize
	return
***REMOVED***

// Encode a slice of uint32s ([]uint32).
// Exactly the same as int32, except for no sign extension.
func (o *Buffer) enc_slice_uint32(p *Properties, base structPointer) error ***REMOVED***
	s := structPointer_Word32Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	for i := 0; i < l; i++ ***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		x := s.Index(i)
		p.valEnc(o, uint64(x))
	***REMOVED***
	return nil
***REMOVED***

func size_slice_uint32(p *Properties, base structPointer) (n int) ***REMOVED***
	s := structPointer_Word32Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	for i := 0; i < l; i++ ***REMOVED***
		n += len(p.tagcode)
		x := s.Index(i)
		n += p.valSize(uint64(x))
	***REMOVED***
	return
***REMOVED***

// Encode a slice of uint32s ([]uint32) in packed format.
// Exactly the same as int32, except for no sign extension.
func (o *Buffer) enc_slice_packed_uint32(p *Properties, base structPointer) error ***REMOVED***
	s := structPointer_Word32Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	// TODO: Reuse a Buffer.
	buf := NewBuffer(nil)
	for i := 0; i < l; i++ ***REMOVED***
		p.valEnc(buf, uint64(s.Index(i)))
	***REMOVED***

	o.buf = append(o.buf, p.tagcode...)
	o.EncodeVarint(uint64(len(buf.buf)))
	o.buf = append(o.buf, buf.buf...)
	return nil
***REMOVED***

func size_slice_packed_uint32(p *Properties, base structPointer) (n int) ***REMOVED***
	s := structPointer_Word32Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	var bufSize int
	for i := 0; i < l; i++ ***REMOVED***
		bufSize += p.valSize(uint64(s.Index(i)))
	***REMOVED***

	n += len(p.tagcode)
	n += sizeVarint(uint64(bufSize))
	n += bufSize
	return
***REMOVED***

// Encode a slice of int64s ([]int64).
func (o *Buffer) enc_slice_int64(p *Properties, base structPointer) error ***REMOVED***
	s := structPointer_Word64Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	for i := 0; i < l; i++ ***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		p.valEnc(o, s.Index(i))
	***REMOVED***
	return nil
***REMOVED***

func size_slice_int64(p *Properties, base structPointer) (n int) ***REMOVED***
	s := structPointer_Word64Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	for i := 0; i < l; i++ ***REMOVED***
		n += len(p.tagcode)
		n += p.valSize(s.Index(i))
	***REMOVED***
	return
***REMOVED***

// Encode a slice of int64s ([]int64) in packed format.
func (o *Buffer) enc_slice_packed_int64(p *Properties, base structPointer) error ***REMOVED***
	s := structPointer_Word64Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	// TODO: Reuse a Buffer.
	buf := NewBuffer(nil)
	for i := 0; i < l; i++ ***REMOVED***
		p.valEnc(buf, s.Index(i))
	***REMOVED***

	o.buf = append(o.buf, p.tagcode...)
	o.EncodeVarint(uint64(len(buf.buf)))
	o.buf = append(o.buf, buf.buf...)
	return nil
***REMOVED***

func size_slice_packed_int64(p *Properties, base structPointer) (n int) ***REMOVED***
	s := structPointer_Word64Slice(base, p.field)
	l := s.Len()
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	var bufSize int
	for i := 0; i < l; i++ ***REMOVED***
		bufSize += p.valSize(s.Index(i))
	***REMOVED***

	n += len(p.tagcode)
	n += sizeVarint(uint64(bufSize))
	n += bufSize
	return
***REMOVED***

// Encode a slice of slice of bytes ([][]byte).
func (o *Buffer) enc_slice_slice_byte(p *Properties, base structPointer) error ***REMOVED***
	ss := *structPointer_BytesSlice(base, p.field)
	l := len(ss)
	if l == 0 ***REMOVED***
		return ErrNil
	***REMOVED***
	for i := 0; i < l; i++ ***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeRawBytes(ss[i])
	***REMOVED***
	return nil
***REMOVED***

func size_slice_slice_byte(p *Properties, base structPointer) (n int) ***REMOVED***
	ss := *structPointer_BytesSlice(base, p.field)
	l := len(ss)
	if l == 0 ***REMOVED***
		return 0
	***REMOVED***
	n += l * len(p.tagcode)
	for i := 0; i < l; i++ ***REMOVED***
		n += sizeRawBytes(ss[i])
	***REMOVED***
	return
***REMOVED***

// Encode a slice of strings ([]string).
func (o *Buffer) enc_slice_string(p *Properties, base structPointer) error ***REMOVED***
	ss := *structPointer_StringSlice(base, p.field)
	l := len(ss)
	for i := 0; i < l; i++ ***REMOVED***
		o.buf = append(o.buf, p.tagcode...)
		o.EncodeStringBytes(ss[i])
	***REMOVED***
	return nil
***REMOVED***

func size_slice_string(p *Properties, base structPointer) (n int) ***REMOVED***
	ss := *structPointer_StringSlice(base, p.field)
	l := len(ss)
	n += l * len(p.tagcode)
	for i := 0; i < l; i++ ***REMOVED***
		n += sizeStringBytes(ss[i])
	***REMOVED***
	return
***REMOVED***

// Encode a slice of message structs ([]*struct).
func (o *Buffer) enc_slice_struct_message(p *Properties, base structPointer) error ***REMOVED***
	var state errorState
	s := structPointer_StructPointerSlice(base, p.field)
	l := s.Len()

	for i := 0; i < l; i++ ***REMOVED***
		structp := s.Index(i)
		if structPointer_IsNil(structp) ***REMOVED***
			return errRepeatedHasNil
		***REMOVED***

		// Can the object marshal itself?
		if p.isMarshaler ***REMOVED***
			m := structPointer_Interface(structp, p.stype).(Marshaler)
			data, err := m.Marshal()
			if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
				return err
			***REMOVED***
			o.buf = append(o.buf, p.tagcode...)
			o.EncodeRawBytes(data)
			continue
		***REMOVED***

		o.buf = append(o.buf, p.tagcode...)
		err := o.enc_len_struct(p.sprop, structp, &state)
		if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
			if err == ErrNil ***REMOVED***
				return errRepeatedHasNil
			***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return state.err
***REMOVED***

func size_slice_struct_message(p *Properties, base structPointer) (n int) ***REMOVED***
	s := structPointer_StructPointerSlice(base, p.field)
	l := s.Len()
	n += l * len(p.tagcode)
	for i := 0; i < l; i++ ***REMOVED***
		structp := s.Index(i)
		if structPointer_IsNil(structp) ***REMOVED***
			return // return the size up to this point
		***REMOVED***

		// Can the object marshal itself?
		if p.isMarshaler ***REMOVED***
			m := structPointer_Interface(structp, p.stype).(Marshaler)
			data, _ := m.Marshal()
			n += sizeRawBytes(data)
			continue
		***REMOVED***

		n0 := size_struct(p.sprop, structp)
		n1 := sizeVarint(uint64(n0)) // size of encoded length
		n += n0 + n1
	***REMOVED***
	return
***REMOVED***

// Encode a slice of group structs ([]*struct).
func (o *Buffer) enc_slice_struct_group(p *Properties, base structPointer) error ***REMOVED***
	var state errorState
	s := structPointer_StructPointerSlice(base, p.field)
	l := s.Len()

	for i := 0; i < l; i++ ***REMOVED***
		b := s.Index(i)
		if structPointer_IsNil(b) ***REMOVED***
			return errRepeatedHasNil
		***REMOVED***

		o.EncodeVarint(uint64((p.Tag << 3) | WireStartGroup))

		err := o.enc_struct(p.sprop, b)

		if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
			if err == ErrNil ***REMOVED***
				return errRepeatedHasNil
			***REMOVED***
			return err
		***REMOVED***

		o.EncodeVarint(uint64((p.Tag << 3) | WireEndGroup))
	***REMOVED***
	return state.err
***REMOVED***

func size_slice_struct_group(p *Properties, base structPointer) (n int) ***REMOVED***
	s := structPointer_StructPointerSlice(base, p.field)
	l := s.Len()

	n += l * sizeVarint(uint64((p.Tag<<3)|WireStartGroup))
	n += l * sizeVarint(uint64((p.Tag<<3)|WireEndGroup))
	for i := 0; i < l; i++ ***REMOVED***
		b := s.Index(i)
		if structPointer_IsNil(b) ***REMOVED***
			return // return size up to this point
		***REMOVED***

		n += size_struct(p.sprop, b)
	***REMOVED***
	return
***REMOVED***

// Encode an extension map.
func (o *Buffer) enc_map(p *Properties, base structPointer) error ***REMOVED***
	exts := structPointer_ExtMap(base, p.field)
	if err := encodeExtensionsMap(*exts); err != nil ***REMOVED***
		return err
	***REMOVED***

	return o.enc_map_body(*exts)
***REMOVED***

func (o *Buffer) enc_exts(p *Properties, base structPointer) error ***REMOVED***
	exts := structPointer_Extensions(base, p.field)

	v, mu := exts.extensionsRead()
	if v == nil ***REMOVED***
		return nil
	***REMOVED***

	mu.Lock()
	defer mu.Unlock()
	if err := encodeExtensionsMap(v); err != nil ***REMOVED***
		return err
	***REMOVED***

	return o.enc_map_body(v)
***REMOVED***

func (o *Buffer) enc_map_body(v map[int32]Extension) error ***REMOVED***
	// Fast-path for common cases: zero or one extensions.
	if len(v) <= 1 ***REMOVED***
		for _, e := range v ***REMOVED***
			o.buf = append(o.buf, e.enc...)
		***REMOVED***
		return nil
	***REMOVED***

	// Sort keys to provide a deterministic encoding.
	keys := make([]int, 0, len(v))
	for k := range v ***REMOVED***
		keys = append(keys, int(k))
	***REMOVED***
	sort.Ints(keys)

	for _, k := range keys ***REMOVED***
		o.buf = append(o.buf, v[int32(k)].enc...)
	***REMOVED***
	return nil
***REMOVED***

func size_map(p *Properties, base structPointer) int ***REMOVED***
	v := structPointer_ExtMap(base, p.field)
	return extensionsMapSize(*v)
***REMOVED***

func size_exts(p *Properties, base structPointer) int ***REMOVED***
	v := structPointer_Extensions(base, p.field)
	return extensionsSize(v)
***REMOVED***

// Encode a map field.
func (o *Buffer) enc_new_map(p *Properties, base structPointer) error ***REMOVED***
	var state errorState // XXX: or do we need to plumb this through?

	/*
		A map defined as
			map<key_type, value_type> map_field = N;
		is encoded in the same way as
			message MapFieldEntry ***REMOVED***
				key_type key = 1;
				value_type value = 2;
			***REMOVED***
			repeated MapFieldEntry map_field = N;
	*/

	v := structPointer_NewAt(base, p.field, p.mtype).Elem() // map[K]V
	if v.Len() == 0 ***REMOVED***
		return nil
	***REMOVED***

	keycopy, valcopy, keybase, valbase := mapEncodeScratch(p.mtype)

	enc := func() error ***REMOVED***
		if err := p.mkeyprop.enc(o, p.mkeyprop, keybase); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := p.mvalprop.enc(o, p.mvalprop, valbase); err != nil && err != ErrNil ***REMOVED***
			return err
		***REMOVED***
		return nil
	***REMOVED***

	// Don't sort map keys. It is not required by the spec, and C++ doesn't do it.
	for _, key := range v.MapKeys() ***REMOVED***
		val := v.MapIndex(key)

		keycopy.Set(key)
		valcopy.Set(val)

		o.buf = append(o.buf, p.tagcode...)
		if err := o.enc_len_thing(enc, &state); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func size_new_map(p *Properties, base structPointer) int ***REMOVED***
	v := structPointer_NewAt(base, p.field, p.mtype).Elem() // map[K]V

	keycopy, valcopy, keybase, valbase := mapEncodeScratch(p.mtype)

	n := 0
	for _, key := range v.MapKeys() ***REMOVED***
		val := v.MapIndex(key)
		keycopy.Set(key)
		valcopy.Set(val)

		// Tag codes for key and val are the responsibility of the sub-sizer.
		keysize := p.mkeyprop.size(p.mkeyprop, keybase)
		valsize := p.mvalprop.size(p.mvalprop, valbase)
		entry := keysize + valsize
		// Add on tag code and length of map entry itself.
		n += len(p.tagcode) + sizeVarint(uint64(entry)) + entry
	***REMOVED***
	return n
***REMOVED***

// mapEncodeScratch returns a new reflect.Value matching the map's value type,
// and a structPointer suitable for passing to an encoder or sizer.
func mapEncodeScratch(mapType reflect.Type) (keycopy, valcopy reflect.Value, keybase, valbase structPointer) ***REMOVED***
	// Prepare addressable doubly-indirect placeholders for the key and value types.
	// This is needed because the element-type encoders expect **T, but the map iteration produces T.

	keycopy = reflect.New(mapType.Key()).Elem()                 // addressable K
	keyptr := reflect.New(reflect.PtrTo(keycopy.Type())).Elem() // addressable *K
	keyptr.Set(keycopy.Addr())                                  //
	keybase = toStructPointer(keyptr.Addr())                    // **K

	// Value types are more varied and require special handling.
	switch mapType.Elem().Kind() ***REMOVED***
	case reflect.Slice:
		// []byte
		var dummy []byte
		valcopy = reflect.ValueOf(&dummy).Elem() // addressable []byte
		valbase = toStructPointer(valcopy.Addr())
	case reflect.Ptr:
		// message; the generated field type is map[K]*Msg (so V is *Msg),
		// so we only need one level of indirection.
		valcopy = reflect.New(mapType.Elem()).Elem() // addressable V
		valbase = toStructPointer(valcopy.Addr())
	default:
		// everything else
		valcopy = reflect.New(mapType.Elem()).Elem()                // addressable V
		valptr := reflect.New(reflect.PtrTo(valcopy.Type())).Elem() // addressable *V
		valptr.Set(valcopy.Addr())                                  //
		valbase = toStructPointer(valptr.Addr())                    // **V
	***REMOVED***
	return
***REMOVED***

// Encode a struct.
func (o *Buffer) enc_struct(prop *StructProperties, base structPointer) error ***REMOVED***
	var state errorState
	// Encode fields in tag order so that decoders may use optimizations
	// that depend on the ordering.
	// https://developers.google.com/protocol-buffers/docs/encoding#order
	for _, i := range prop.order ***REMOVED***
		p := prop.Prop[i]
		if p.enc != nil ***REMOVED***
			err := p.enc(o, p, base)
			if err != nil ***REMOVED***
				if err == ErrNil ***REMOVED***
					if p.Required && state.err == nil ***REMOVED***
						state.err = &RequiredNotSetError***REMOVED***p.Name***REMOVED***
					***REMOVED***
				***REMOVED*** else if err == errRepeatedHasNil ***REMOVED***
					// Give more context to nil values in repeated fields.
					return errors.New("repeated field " + p.OrigName + " has nil element")
				***REMOVED*** else if !state.shouldContinue(err, p) ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			if len(o.buf) > maxMarshalSize ***REMOVED***
				return ErrTooLarge
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Do oneof fields.
	if prop.oneofMarshaler != nil ***REMOVED***
		m := structPointer_Interface(base, prop.stype).(Message)
		if err := prop.oneofMarshaler(m, o); err == ErrNil ***REMOVED***
			return errOneofHasNil
		***REMOVED*** else if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Add unrecognized fields at the end.
	if prop.unrecField.IsValid() ***REMOVED***
		v := *structPointer_Bytes(base, prop.unrecField)
		if len(o.buf)+len(v) > maxMarshalSize ***REMOVED***
			return ErrTooLarge
		***REMOVED***
		if len(v) > 0 ***REMOVED***
			o.buf = append(o.buf, v...)
		***REMOVED***
	***REMOVED***

	return state.err
***REMOVED***

func size_struct(prop *StructProperties, base structPointer) (n int) ***REMOVED***
	for _, i := range prop.order ***REMOVED***
		p := prop.Prop[i]
		if p.size != nil ***REMOVED***
			n += p.size(p, base)
		***REMOVED***
	***REMOVED***

	// Add unrecognized fields at the end.
	if prop.unrecField.IsValid() ***REMOVED***
		v := *structPointer_Bytes(base, prop.unrecField)
		n += len(v)
	***REMOVED***

	// Factor in any oneof fields.
	if prop.oneofSizer != nil ***REMOVED***
		m := structPointer_Interface(base, prop.stype).(Message)
		n += prop.oneofSizer(m)
	***REMOVED***

	return
***REMOVED***

var zeroes [20]byte // longer than any conceivable sizeVarint

// Encode a struct, preceded by its encoded length (as a varint).
func (o *Buffer) enc_len_struct(prop *StructProperties, base structPointer, state *errorState) error ***REMOVED***
	return o.enc_len_thing(func() error ***REMOVED*** return o.enc_struct(prop, base) ***REMOVED***, state)
***REMOVED***

// Encode something, preceded by its encoded length (as a varint).
func (o *Buffer) enc_len_thing(enc func() error, state *errorState) error ***REMOVED***
	iLen := len(o.buf)
	o.buf = append(o.buf, 0, 0, 0, 0) // reserve four bytes for length
	iMsg := len(o.buf)
	err := enc()
	if err != nil && !state.shouldContinue(err, nil) ***REMOVED***
		return err
	***REMOVED***
	lMsg := len(o.buf) - iMsg
	lLen := sizeVarint(uint64(lMsg))
	switch x := lLen - (iMsg - iLen); ***REMOVED***
	case x > 0: // actual length is x bytes larger than the space we reserved
		// Move msg x bytes right.
		o.buf = append(o.buf, zeroes[:x]...)
		copy(o.buf[iMsg+x:], o.buf[iMsg:iMsg+lMsg])
	case x < 0: // actual length is x bytes smaller than the space we reserved
		// Move msg x bytes left.
		copy(o.buf[iMsg+x:], o.buf[iMsg:iMsg+lMsg])
		o.buf = o.buf[:len(o.buf)+x] // x is negative
	***REMOVED***
	// Encode the length in the reserved space.
	o.buf = o.buf[:iLen]
	o.EncodeVarint(uint64(lMsg))
	o.buf = o.buf[:len(o.buf)+lMsg]
	return state.err
***REMOVED***

// errorState maintains the first error that occurs and updates that error
// with additional context.
type errorState struct ***REMOVED***
	err error
***REMOVED***

// shouldContinue reports whether encoding should continue upon encountering the
// given error. If the error is RequiredNotSetError, shouldContinue returns true
// and, if this is the first appearance of that error, remembers it for future
// reporting.
//
// If prop is not nil, it may update any error with additional context about the
// field with the error.
func (s *errorState) shouldContinue(err error, prop *Properties) bool ***REMOVED***
	// Ignore unset required fields.
	reqNotSet, ok := err.(*RequiredNotSetError)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	if s.err == nil ***REMOVED***
		if prop != nil ***REMOVED***
			err = &RequiredNotSetError***REMOVED***prop.Name + "." + reqNotSet.field***REMOVED***
		***REMOVED***
		s.err = err
	***REMOVED***
	return true
***REMOVED***
