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
 * Routines for decoding protocol buffer data to construct in-memory representations.
 */

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

// errOverflow is returned when an integer is too large to be represented.
var errOverflow = errors.New("proto: integer overflow")

// ErrInternalBadWireType is returned by generated code when an incorrect
// wire type is encountered. It does not get returned to user code.
var ErrInternalBadWireType = errors.New("proto: internal error: bad wiretype for oneof")

// The fundamental decoders that interpret bytes on the wire.
// Those that take integer types all return uint64 and are
// therefore of type valueDecoder.

// DecodeVarint reads a varint-encoded integer from the slice.
// It returns the integer and the number of bytes consumed, or
// zero if there is not enough.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func DecodeVarint(buf []byte) (x uint64, n int) ***REMOVED***
	for shift := uint(0); shift < 64; shift += 7 ***REMOVED***
		if n >= len(buf) ***REMOVED***
			return 0, 0
		***REMOVED***
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 ***REMOVED***
			return x, n
		***REMOVED***
	***REMOVED***

	// The number is too large to represent in a 64-bit value.
	return 0, 0
***REMOVED***

func (p *Buffer) decodeVarintSlow() (x uint64, err error) ***REMOVED***
	i := p.index
	l := len(p.buf)

	for shift := uint(0); shift < 64; shift += 7 ***REMOVED***
		if i >= l ***REMOVED***
			err = io.ErrUnexpectedEOF
			return
		***REMOVED***
		b := p.buf[i]
		i++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 ***REMOVED***
			p.index = i
			return
		***REMOVED***
	***REMOVED***

	// The number is too large to represent in a 64-bit value.
	err = errOverflow
	return
***REMOVED***

// DecodeVarint reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (p *Buffer) DecodeVarint() (x uint64, err error) ***REMOVED***
	i := p.index
	buf := p.buf

	if i >= len(buf) ***REMOVED***
		return 0, io.ErrUnexpectedEOF
	***REMOVED*** else if buf[i] < 0x80 ***REMOVED***
		p.index++
		return uint64(buf[i]), nil
	***REMOVED*** else if len(buf)-i < 10 ***REMOVED***
		return p.decodeVarintSlow()
	***REMOVED***

	var b uint64
	// we already checked the first byte
	x = uint64(buf[i]) - 0x80
	i++

	b = uint64(buf[i])
	i++
	x += b << 7
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	x -= 0x80 << 7

	b = uint64(buf[i])
	i++
	x += b << 14
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	x -= 0x80 << 14

	b = uint64(buf[i])
	i++
	x += b << 21
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	x -= 0x80 << 21

	b = uint64(buf[i])
	i++
	x += b << 28
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	x -= 0x80 << 28

	b = uint64(buf[i])
	i++
	x += b << 35
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	x -= 0x80 << 35

	b = uint64(buf[i])
	i++
	x += b << 42
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	x -= 0x80 << 42

	b = uint64(buf[i])
	i++
	x += b << 49
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	x -= 0x80 << 49

	b = uint64(buf[i])
	i++
	x += b << 56
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	x -= 0x80 << 56

	b = uint64(buf[i])
	i++
	x += b << 63
	if b&0x80 == 0 ***REMOVED***
		goto done
	***REMOVED***
	// x -= 0x80 << 63 // Always zero.

	return 0, errOverflow

done:
	p.index = i
	return x, nil
***REMOVED***

// DecodeFixed64 reads a 64-bit integer from the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (p *Buffer) DecodeFixed64() (x uint64, err error) ***REMOVED***
	// x, err already 0
	i := p.index + 8
	if i < 0 || i > len(p.buf) ***REMOVED***
		err = io.ErrUnexpectedEOF
		return
	***REMOVED***
	p.index = i

	x = uint64(p.buf[i-8])
	x |= uint64(p.buf[i-7]) << 8
	x |= uint64(p.buf[i-6]) << 16
	x |= uint64(p.buf[i-5]) << 24
	x |= uint64(p.buf[i-4]) << 32
	x |= uint64(p.buf[i-3]) << 40
	x |= uint64(p.buf[i-2]) << 48
	x |= uint64(p.buf[i-1]) << 56
	return
***REMOVED***

// DecodeFixed32 reads a 32-bit integer from the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (p *Buffer) DecodeFixed32() (x uint64, err error) ***REMOVED***
	// x, err already 0
	i := p.index + 4
	if i < 0 || i > len(p.buf) ***REMOVED***
		err = io.ErrUnexpectedEOF
		return
	***REMOVED***
	p.index = i

	x = uint64(p.buf[i-4])
	x |= uint64(p.buf[i-3]) << 8
	x |= uint64(p.buf[i-2]) << 16
	x |= uint64(p.buf[i-1]) << 24
	return
***REMOVED***

// DecodeZigzag64 reads a zigzag-encoded 64-bit integer
// from the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (p *Buffer) DecodeZigzag64() (x uint64, err error) ***REMOVED***
	x, err = p.DecodeVarint()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	x = (x >> 1) ^ uint64((int64(x&1)<<63)>>63)
	return
***REMOVED***

// DecodeZigzag32 reads a zigzag-encoded 32-bit integer
// from  the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (p *Buffer) DecodeZigzag32() (x uint64, err error) ***REMOVED***
	x, err = p.DecodeVarint()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	x = uint64((uint32(x) >> 1) ^ uint32((int32(x&1)<<31)>>31))
	return
***REMOVED***

// These are not ValueDecoders: they produce an array of bytes or a string.
// bytes, embedded messages

// DecodeRawBytes reads a count-delimited byte buffer from the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (p *Buffer) DecodeRawBytes(alloc bool) (buf []byte, err error) ***REMOVED***
	n, err := p.DecodeVarint()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	nb := int(n)
	if nb < 0 ***REMOVED***
		return nil, fmt.Errorf("proto: bad byte length %d", nb)
	***REMOVED***
	end := p.index + nb
	if end < p.index || end > len(p.buf) ***REMOVED***
		return nil, io.ErrUnexpectedEOF
	***REMOVED***

	if !alloc ***REMOVED***
		// todo: check if can get more uses of alloc=false
		buf = p.buf[p.index:end]
		p.index += nb
		return
	***REMOVED***

	buf = make([]byte, nb)
	copy(buf, p.buf[p.index:])
	p.index += nb
	return
***REMOVED***

// DecodeStringBytes reads an encoded string from the Buffer.
// This is the format used for the proto2 string type.
func (p *Buffer) DecodeStringBytes() (s string, err error) ***REMOVED***
	buf, err := p.DecodeRawBytes(false)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return string(buf), nil
***REMOVED***

// Skip the next item in the buffer. Its wire type is decoded and presented as an argument.
// If the protocol buffer has extensions, and the field matches, add it as an extension.
// Otherwise, if the XXX_unrecognized field exists, append the skipped data there.
func (o *Buffer) skipAndSave(t reflect.Type, tag, wire int, base structPointer, unrecField field) error ***REMOVED***
	oi := o.index

	err := o.skip(t, tag, wire)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !unrecField.IsValid() ***REMOVED***
		return nil
	***REMOVED***

	ptr := structPointer_Bytes(base, unrecField)

	// Add the skipped field to struct field
	obuf := o.buf

	o.buf = *ptr
	o.EncodeVarint(uint64(tag<<3 | wire))
	*ptr = append(o.buf, obuf[oi:o.index]...)

	o.buf = obuf

	return nil
***REMOVED***

// Skip the next item in the buffer. Its wire type is decoded and presented as an argument.
func (o *Buffer) skip(t reflect.Type, tag, wire int) error ***REMOVED***

	var u uint64
	var err error

	switch wire ***REMOVED***
	case WireVarint:
		_, err = o.DecodeVarint()
	case WireFixed64:
		_, err = o.DecodeFixed64()
	case WireBytes:
		_, err = o.DecodeRawBytes(false)
	case WireFixed32:
		_, err = o.DecodeFixed32()
	case WireStartGroup:
		for ***REMOVED***
			u, err = o.DecodeVarint()
			if err != nil ***REMOVED***
				break
			***REMOVED***
			fwire := int(u & 0x7)
			if fwire == WireEndGroup ***REMOVED***
				break
			***REMOVED***
			ftag := int(u >> 3)
			err = o.skip(t, ftag, fwire)
			if err != nil ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	default:
		err = fmt.Errorf("proto: can't skip unknown wire type %d for %s", wire, t)
	***REMOVED***
	return err
***REMOVED***

// Unmarshaler is the interface representing objects that can
// unmarshal themselves.  The method should reset the receiver before
// decoding starts.  The argument points to data that may be
// overwritten, so implementations should not keep references to the
// buffer.
type Unmarshaler interface ***REMOVED***
	Unmarshal([]byte) error
***REMOVED***

// Unmarshal parses the protocol buffer representation in buf and places the
// decoded result in pb.  If the struct underlying pb does not match
// the data in buf, the results can be unpredictable.
//
// Unmarshal resets pb before starting to unmarshal, so any
// existing data in pb is always removed. Use UnmarshalMerge
// to preserve and append to existing data.
func Unmarshal(buf []byte, pb Message) error ***REMOVED***
	pb.Reset()
	return UnmarshalMerge(buf, pb)
***REMOVED***

// UnmarshalMerge parses the protocol buffer representation in buf and
// writes the decoded result to pb.  If the struct underlying pb does not match
// the data in buf, the results can be unpredictable.
//
// UnmarshalMerge merges into existing data in pb.
// Most code should use Unmarshal instead.
func UnmarshalMerge(buf []byte, pb Message) error ***REMOVED***
	// If the object can unmarshal itself, let it.
	if u, ok := pb.(Unmarshaler); ok ***REMOVED***
		return u.Unmarshal(buf)
	***REMOVED***
	return NewBuffer(buf).Unmarshal(pb)
***REMOVED***

// DecodeMessage reads a count-delimited message from the Buffer.
func (p *Buffer) DecodeMessage(pb Message) error ***REMOVED***
	enc, err := p.DecodeRawBytes(false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return NewBuffer(enc).Unmarshal(pb)
***REMOVED***

// DecodeGroup reads a tag-delimited group from the Buffer.
func (p *Buffer) DecodeGroup(pb Message) error ***REMOVED***
	typ, base, err := getbase(pb)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return p.unmarshalType(typ.Elem(), GetProperties(typ.Elem()), true, base)
***REMOVED***

// Unmarshal parses the protocol buffer representation in the
// Buffer and places the decoded result in pb.  If the struct
// underlying pb does not match the data in the buffer, the results can be
// unpredictable.
//
// Unlike proto.Unmarshal, this does not reset pb before starting to unmarshal.
func (p *Buffer) Unmarshal(pb Message) error ***REMOVED***
	// If the object can unmarshal itself, let it.
	if u, ok := pb.(Unmarshaler); ok ***REMOVED***
		err := u.Unmarshal(p.buf[p.index:])
		p.index = len(p.buf)
		return err
	***REMOVED***

	typ, base, err := getbase(pb)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = p.unmarshalType(typ.Elem(), GetProperties(typ.Elem()), false, base)

	if collectStats ***REMOVED***
		stats.Decode++
	***REMOVED***

	return err
***REMOVED***

// unmarshalType does the work of unmarshaling a structure.
func (o *Buffer) unmarshalType(st reflect.Type, prop *StructProperties, is_group bool, base structPointer) error ***REMOVED***
	var state errorState
	required, reqFields := prop.reqCount, uint64(0)

	var err error
	for err == nil && o.index < len(o.buf) ***REMOVED***
		oi := o.index
		var u uint64
		u, err = o.DecodeVarint()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		wire := int(u & 0x7)
		if wire == WireEndGroup ***REMOVED***
			if is_group ***REMOVED***
				if required > 0 ***REMOVED***
					// Not enough information to determine the exact field.
					// (See below.)
					return &RequiredNotSetError***REMOVED***"***REMOVED***Unknown***REMOVED***"***REMOVED***
				***REMOVED***
				return nil // input is satisfied
			***REMOVED***
			return fmt.Errorf("proto: %s: wiretype end group for non-group", st)
		***REMOVED***
		tag := int(u >> 3)
		if tag <= 0 ***REMOVED***
			return fmt.Errorf("proto: %s: illegal tag %d (wire type %d)", st, tag, wire)
		***REMOVED***
		fieldnum, ok := prop.decoderTags.get(tag)
		if !ok ***REMOVED***
			// Maybe it's an extension?
			if prop.extendable ***REMOVED***
				if e, _ := extendable(structPointer_Interface(base, st)); isExtensionField(e, int32(tag)) ***REMOVED***
					if err = o.skip(st, tag, wire); err == nil ***REMOVED***
						extmap := e.extensionsWrite()
						ext := extmap[int32(tag)] // may be missing
						ext.enc = append(ext.enc, o.buf[oi:o.index]...)
						extmap[int32(tag)] = ext
					***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			// Maybe it's a oneof?
			if prop.oneofUnmarshaler != nil ***REMOVED***
				m := structPointer_Interface(base, st).(Message)
				// First return value indicates whether tag is a oneof field.
				ok, err = prop.oneofUnmarshaler(m, tag, wire, o)
				if err == ErrInternalBadWireType ***REMOVED***
					// Map the error to something more descriptive.
					// Do the formatting here to save generated code space.
					err = fmt.Errorf("bad wiretype for oneof field in %T", m)
				***REMOVED***
				if ok ***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
			err = o.skipAndSave(st, tag, wire, base, prop.unrecField)
			continue
		***REMOVED***
		p := prop.Prop[fieldnum]

		if p.dec == nil ***REMOVED***
			fmt.Fprintf(os.Stderr, "proto: no protobuf decoder for %s.%s\n", st, st.Field(fieldnum).Name)
			continue
		***REMOVED***
		dec := p.dec
		if wire != WireStartGroup && wire != p.WireType ***REMOVED***
			if wire == WireBytes && p.packedDec != nil ***REMOVED***
				// a packable field
				dec = p.packedDec
			***REMOVED*** else ***REMOVED***
				err = fmt.Errorf("proto: bad wiretype for field %s.%s: got wiretype %d, want %d", st, st.Field(fieldnum).Name, wire, p.WireType)
				continue
			***REMOVED***
		***REMOVED***
		decErr := dec(o, p, base)
		if decErr != nil && !state.shouldContinue(decErr, p) ***REMOVED***
			err = decErr
		***REMOVED***
		if err == nil && p.Required ***REMOVED***
			// Successfully decoded a required field.
			if tag <= 64 ***REMOVED***
				// use bitmap for fields 1-64 to catch field reuse.
				var mask uint64 = 1 << uint64(tag-1)
				if reqFields&mask == 0 ***REMOVED***
					// new required field
					reqFields |= mask
					required--
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// This is imprecise. It can be fooled by a required field
				// with a tag > 64 that is encoded twice; that's very rare.
				// A fully correct implementation would require allocating
				// a data structure, which we would like to avoid.
				required--
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if err == nil ***REMOVED***
		if is_group ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		if state.err != nil ***REMOVED***
			return state.err
		***REMOVED***
		if required > 0 ***REMOVED***
			// Not enough information to determine the exact field. If we use extra
			// CPU, we could determine the field only if the missing required field
			// has a tag <= 64 and we check reqFields.
			return &RequiredNotSetError***REMOVED***"***REMOVED***Unknown***REMOVED***"***REMOVED***
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

// Individual type decoders
// For each,
//	u is the decoded value,
//	v is a pointer to the field (pointer) in the struct

// Sizes of the pools to allocate inside the Buffer.
// The goal is modest amortization and allocation
// on at least 16-byte boundaries.
const (
	boolPoolSize   = 16
	uint32PoolSize = 8
	uint64PoolSize = 4
)

// Decode a bool.
func (o *Buffer) dec_bool(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if len(o.bools) == 0 ***REMOVED***
		o.bools = make([]bool, boolPoolSize)
	***REMOVED***
	o.bools[0] = u != 0
	*structPointer_Bool(base, p.field) = &o.bools[0]
	o.bools = o.bools[1:]
	return nil
***REMOVED***

func (o *Buffer) dec_proto3_bool(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*structPointer_BoolVal(base, p.field) = u != 0
	return nil
***REMOVED***

// Decode an int32.
func (o *Buffer) dec_int32(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	word32_Set(structPointer_Word32(base, p.field), o, uint32(u))
	return nil
***REMOVED***

func (o *Buffer) dec_proto3_int32(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	word32Val_Set(structPointer_Word32Val(base, p.field), uint32(u))
	return nil
***REMOVED***

// Decode an int64.
func (o *Buffer) dec_int64(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	word64_Set(structPointer_Word64(base, p.field), o, u)
	return nil
***REMOVED***

func (o *Buffer) dec_proto3_int64(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	word64Val_Set(structPointer_Word64Val(base, p.field), o, u)
	return nil
***REMOVED***

// Decode a string.
func (o *Buffer) dec_string(p *Properties, base structPointer) error ***REMOVED***
	s, err := o.DecodeStringBytes()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*structPointer_String(base, p.field) = &s
	return nil
***REMOVED***

func (o *Buffer) dec_proto3_string(p *Properties, base structPointer) error ***REMOVED***
	s, err := o.DecodeStringBytes()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*structPointer_StringVal(base, p.field) = s
	return nil
***REMOVED***

// Decode a slice of bytes ([]byte).
func (o *Buffer) dec_slice_byte(p *Properties, base structPointer) error ***REMOVED***
	b, err := o.DecodeRawBytes(true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*structPointer_Bytes(base, p.field) = b
	return nil
***REMOVED***

// Decode a slice of bools ([]bool).
func (o *Buffer) dec_slice_bool(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	v := structPointer_BoolSlice(base, p.field)
	*v = append(*v, u != 0)
	return nil
***REMOVED***

// Decode a slice of bools ([]bool) in packed format.
func (o *Buffer) dec_slice_packed_bool(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_BoolSlice(base, p.field)

	nn, err := o.DecodeVarint()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	nb := int(nn) // number of bytes of encoded bools
	fin := o.index + nb
	if fin < o.index ***REMOVED***
		return errOverflow
	***REMOVED***

	y := *v
	for o.index < fin ***REMOVED***
		u, err := p.valDec(o)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		y = append(y, u != 0)
	***REMOVED***

	*v = y
	return nil
***REMOVED***

// Decode a slice of int32s ([]int32).
func (o *Buffer) dec_slice_int32(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	structPointer_Word32Slice(base, p.field).Append(uint32(u))
	return nil
***REMOVED***

// Decode a slice of int32s ([]int32) in packed format.
func (o *Buffer) dec_slice_packed_int32(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word32Slice(base, p.field)

	nn, err := o.DecodeVarint()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	nb := int(nn) // number of bytes of encoded int32s

	fin := o.index + nb
	if fin < o.index ***REMOVED***
		return errOverflow
	***REMOVED***
	for o.index < fin ***REMOVED***
		u, err := p.valDec(o)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v.Append(uint32(u))
	***REMOVED***
	return nil
***REMOVED***

// Decode a slice of int64s ([]int64).
func (o *Buffer) dec_slice_int64(p *Properties, base structPointer) error ***REMOVED***
	u, err := p.valDec(o)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	structPointer_Word64Slice(base, p.field).Append(u)
	return nil
***REMOVED***

// Decode a slice of int64s ([]int64) in packed format.
func (o *Buffer) dec_slice_packed_int64(p *Properties, base structPointer) error ***REMOVED***
	v := structPointer_Word64Slice(base, p.field)

	nn, err := o.DecodeVarint()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	nb := int(nn) // number of bytes of encoded int64s

	fin := o.index + nb
	if fin < o.index ***REMOVED***
		return errOverflow
	***REMOVED***
	for o.index < fin ***REMOVED***
		u, err := p.valDec(o)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		v.Append(u)
	***REMOVED***
	return nil
***REMOVED***

// Decode a slice of strings ([]string).
func (o *Buffer) dec_slice_string(p *Properties, base structPointer) error ***REMOVED***
	s, err := o.DecodeStringBytes()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	v := structPointer_StringSlice(base, p.field)
	*v = append(*v, s)
	return nil
***REMOVED***

// Decode a slice of slice of bytes ([][]byte).
func (o *Buffer) dec_slice_slice_byte(p *Properties, base structPointer) error ***REMOVED***
	b, err := o.DecodeRawBytes(true)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	v := structPointer_BytesSlice(base, p.field)
	*v = append(*v, b)
	return nil
***REMOVED***

// Decode a map field.
func (o *Buffer) dec_new_map(p *Properties, base structPointer) error ***REMOVED***
	raw, err := o.DecodeRawBytes(false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	oi := o.index       // index at the end of this map entry
	o.index -= len(raw) // move buffer back to start of map entry

	mptr := structPointer_NewAt(base, p.field, p.mtype) // *map[K]V
	if mptr.Elem().IsNil() ***REMOVED***
		mptr.Elem().Set(reflect.MakeMap(mptr.Type().Elem()))
	***REMOVED***
	v := mptr.Elem() // map[K]V

	// Prepare addressable doubly-indirect placeholders for the key and value types.
	// See enc_new_map for why.
	keyptr := reflect.New(reflect.PtrTo(p.mtype.Key())).Elem() // addressable *K
	keybase := toStructPointer(keyptr.Addr())                  // **K

	var valbase structPointer
	var valptr reflect.Value
	switch p.mtype.Elem().Kind() ***REMOVED***
	case reflect.Slice:
		// []byte
		var dummy []byte
		valptr = reflect.ValueOf(&dummy)  // *[]byte
		valbase = toStructPointer(valptr) // *[]byte
	case reflect.Ptr:
		// message; valptr is **Msg; need to allocate the intermediate pointer
		valptr = reflect.New(reflect.PtrTo(p.mtype.Elem())).Elem() // addressable *V
		valptr.Set(reflect.New(valptr.Type().Elem()))
		valbase = toStructPointer(valptr)
	default:
		// everything else
		valptr = reflect.New(reflect.PtrTo(p.mtype.Elem())).Elem() // addressable *V
		valbase = toStructPointer(valptr.Addr())                   // **V
	***REMOVED***

	// Decode.
	// This parses a restricted wire format, namely the encoding of a message
	// with two fields. See enc_new_map for the format.
	for o.index < oi ***REMOVED***
		// tagcode for key and value properties are always a single byte
		// because they have tags 1 and 2.
		tagcode := o.buf[o.index]
		o.index++
		switch tagcode ***REMOVED***
		case p.mkeyprop.tagcode[0]:
			if err := p.mkeyprop.dec(o, p.mkeyprop, keybase); err != nil ***REMOVED***
				return err
			***REMOVED***
		case p.mvalprop.tagcode[0]:
			if err := p.mvalprop.dec(o, p.mvalprop, valbase); err != nil ***REMOVED***
				return err
			***REMOVED***
		default:
			// TODO: Should we silently skip this instead?
			return fmt.Errorf("proto: bad map data tag %d", raw[0])
		***REMOVED***
	***REMOVED***
	keyelem, valelem := keyptr.Elem(), valptr.Elem()
	if !keyelem.IsValid() ***REMOVED***
		keyelem = reflect.Zero(p.mtype.Key())
	***REMOVED***
	if !valelem.IsValid() ***REMOVED***
		valelem = reflect.Zero(p.mtype.Elem())
	***REMOVED***

	v.SetMapIndex(keyelem, valelem)
	return nil
***REMOVED***

// Decode a group.
func (o *Buffer) dec_struct_group(p *Properties, base structPointer) error ***REMOVED***
	bas := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(bas) ***REMOVED***
		// allocate new nested message
		bas = toStructPointer(reflect.New(p.stype))
		structPointer_SetStructPointer(base, p.field, bas)
	***REMOVED***
	return o.unmarshalType(p.stype, p.sprop, true, bas)
***REMOVED***

// Decode an embedded message.
func (o *Buffer) dec_struct_message(p *Properties, base structPointer) (err error) ***REMOVED***
	raw, e := o.DecodeRawBytes(false)
	if e != nil ***REMOVED***
		return e
	***REMOVED***

	bas := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(bas) ***REMOVED***
		// allocate new nested message
		bas = toStructPointer(reflect.New(p.stype))
		structPointer_SetStructPointer(base, p.field, bas)
	***REMOVED***

	// If the object can unmarshal itself, let it.
	if p.isUnmarshaler ***REMOVED***
		iv := structPointer_Interface(bas, p.stype)
		return iv.(Unmarshaler).Unmarshal(raw)
	***REMOVED***

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	err = o.unmarshalType(p.stype, p.sprop, false, bas)
	o.buf = obuf
	o.index = oi

	return err
***REMOVED***

// Decode a slice of embedded messages.
func (o *Buffer) dec_slice_struct_message(p *Properties, base structPointer) error ***REMOVED***
	return o.dec_slice_struct(p, false, base)
***REMOVED***

// Decode a slice of embedded groups.
func (o *Buffer) dec_slice_struct_group(p *Properties, base structPointer) error ***REMOVED***
	return o.dec_slice_struct(p, true, base)
***REMOVED***

// Decode a slice of structs ([]*struct).
func (o *Buffer) dec_slice_struct(p *Properties, is_group bool, base structPointer) error ***REMOVED***
	v := reflect.New(p.stype)
	bas := toStructPointer(v)
	structPointer_StructPointerSlice(base, p.field).Append(bas)

	if is_group ***REMOVED***
		err := o.unmarshalType(p.stype, p.sprop, is_group, bas)
		return err
	***REMOVED***

	raw, err := o.DecodeRawBytes(false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the object can unmarshal itself, let it.
	if p.isUnmarshaler ***REMOVED***
		iv := v.Interface()
		return iv.(Unmarshaler).Unmarshal(raw)
	***REMOVED***

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	err = o.unmarshalType(p.stype, p.sprop, is_group, bas)

	o.buf = obuf
	o.index = oi

	return err
***REMOVED***
