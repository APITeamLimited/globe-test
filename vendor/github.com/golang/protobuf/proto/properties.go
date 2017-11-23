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
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const debug bool = false

// Constants that identify the encoding of a value on the wire.
const (
	WireVarint     = 0
	WireFixed64    = 1
	WireBytes      = 2
	WireStartGroup = 3
	WireEndGroup   = 4
	WireFixed32    = 5
)

const startSize = 10 // initial slice/string sizes

// Encoders are defined in encode.go
// An encoder outputs the full representation of a field, including its
// tag and encoder type.
type encoder func(p *Buffer, prop *Properties, base structPointer) error

// A valueEncoder encodes a single integer in a particular encoding.
type valueEncoder func(o *Buffer, x uint64) error

// Sizers are defined in encode.go
// A sizer returns the encoded size of a field, including its tag and encoder
// type.
type sizer func(prop *Properties, base structPointer) int

// A valueSizer returns the encoded size of a single integer in a particular
// encoding.
type valueSizer func(x uint64) int

// Decoders are defined in decode.go
// A decoder creates a value from its wire representation.
// Unrecognized subelements are saved in unrec.
type decoder func(p *Buffer, prop *Properties, base structPointer) error

// A valueDecoder decodes a single integer in a particular encoding.
type valueDecoder func(o *Buffer) (x uint64, err error)

// A oneofMarshaler does the marshaling for all oneof fields in a message.
type oneofMarshaler func(Message, *Buffer) error

// A oneofUnmarshaler does the unmarshaling for a oneof field in a message.
type oneofUnmarshaler func(Message, int, int, *Buffer) (bool, error)

// A oneofSizer does the sizing for all oneof fields in a message.
type oneofSizer func(Message) int

// tagMap is an optimization over map[int]int for typical protocol buffer
// use-cases. Encoded protocol buffers are often in tag order with small tag
// numbers.
type tagMap struct ***REMOVED***
	fastTags []int
	slowTags map[int]int
***REMOVED***

// tagMapFastLimit is the upper bound on the tag number that will be stored in
// the tagMap slice rather than its map.
const tagMapFastLimit = 1024

func (p *tagMap) get(t int) (int, bool) ***REMOVED***
	if t > 0 && t < tagMapFastLimit ***REMOVED***
		if t >= len(p.fastTags) ***REMOVED***
			return 0, false
		***REMOVED***
		fi := p.fastTags[t]
		return fi, fi >= 0
	***REMOVED***
	fi, ok := p.slowTags[t]
	return fi, ok
***REMOVED***

func (p *tagMap) put(t int, fi int) ***REMOVED***
	if t > 0 && t < tagMapFastLimit ***REMOVED***
		for len(p.fastTags) < t+1 ***REMOVED***
			p.fastTags = append(p.fastTags, -1)
		***REMOVED***
		p.fastTags[t] = fi
		return
	***REMOVED***
	if p.slowTags == nil ***REMOVED***
		p.slowTags = make(map[int]int)
	***REMOVED***
	p.slowTags[t] = fi
***REMOVED***

// StructProperties represents properties for all the fields of a struct.
// decoderTags and decoderOrigNames should only be used by the decoder.
type StructProperties struct ***REMOVED***
	Prop             []*Properties  // properties for each field
	reqCount         int            // required count
	decoderTags      tagMap         // map from proto tag to struct field number
	decoderOrigNames map[string]int // map from original name to struct field number
	order            []int          // list of struct field numbers in tag order
	unrecField       field          // field id of the XXX_unrecognized []byte field
	extendable       bool           // is this an extendable proto

	oneofMarshaler   oneofMarshaler
	oneofUnmarshaler oneofUnmarshaler
	oneofSizer       oneofSizer
	stype            reflect.Type

	// OneofTypes contains information about the oneof fields in this message.
	// It is keyed by the original name of a field.
	OneofTypes map[string]*OneofProperties
***REMOVED***

// OneofProperties represents information about a specific field in a oneof.
type OneofProperties struct ***REMOVED***
	Type  reflect.Type // pointer to generated struct type for this oneof field
	Field int          // struct field number of the containing oneof in the message
	Prop  *Properties
***REMOVED***

// Implement the sorting interface so we can sort the fields in tag order, as recommended by the spec.
// See encode.go, (*Buffer).enc_struct.

func (sp *StructProperties) Len() int ***REMOVED*** return len(sp.order) ***REMOVED***
func (sp *StructProperties) Less(i, j int) bool ***REMOVED***
	return sp.Prop[sp.order[i]].Tag < sp.Prop[sp.order[j]].Tag
***REMOVED***
func (sp *StructProperties) Swap(i, j int) ***REMOVED*** sp.order[i], sp.order[j] = sp.order[j], sp.order[i] ***REMOVED***

// Properties represents the protocol-specific behavior of a single struct field.
type Properties struct ***REMOVED***
	Name     string // name of the field, for error messages
	OrigName string // original name before protocol compiler (always set)
	JSONName string // name to use for JSON; determined by protoc
	Wire     string
	WireType int
	Tag      int
	Required bool
	Optional bool
	Repeated bool
	Packed   bool   // relevant for repeated primitives only
	Enum     string // set for enum types only
	proto3   bool   // whether this is known to be a proto3 field; set for []byte only
	oneof    bool   // whether this is a oneof field

	Default    string // default value
	HasDefault bool   // whether an explicit default was provided
	def_uint64 uint64

	enc           encoder
	valEnc        valueEncoder // set for bool and numeric types only
	field         field
	tagcode       []byte // encoding of EncodeVarint((Tag<<3)|WireType)
	tagbuf        [8]byte
	stype         reflect.Type      // set for struct types only
	sprop         *StructProperties // set for struct types only
	isMarshaler   bool
	isUnmarshaler bool

	mtype    reflect.Type // set for map types only
	mkeyprop *Properties  // set for map types only
	mvalprop *Properties  // set for map types only

	size    sizer
	valSize valueSizer // set for bool and numeric types only

	dec    decoder
	valDec valueDecoder // set for bool and numeric types only

	// If this is a packable field, this will be the decoder for the packed version of the field.
	packedDec decoder
***REMOVED***

// String formats the properties in the protobuf struct field tag style.
func (p *Properties) String() string ***REMOVED***
	s := p.Wire
	s = ","
	s += strconv.Itoa(p.Tag)
	if p.Required ***REMOVED***
		s += ",req"
	***REMOVED***
	if p.Optional ***REMOVED***
		s += ",opt"
	***REMOVED***
	if p.Repeated ***REMOVED***
		s += ",rep"
	***REMOVED***
	if p.Packed ***REMOVED***
		s += ",packed"
	***REMOVED***
	s += ",name=" + p.OrigName
	if p.JSONName != p.OrigName ***REMOVED***
		s += ",json=" + p.JSONName
	***REMOVED***
	if p.proto3 ***REMOVED***
		s += ",proto3"
	***REMOVED***
	if p.oneof ***REMOVED***
		s += ",oneof"
	***REMOVED***
	if len(p.Enum) > 0 ***REMOVED***
		s += ",enum=" + p.Enum
	***REMOVED***
	if p.HasDefault ***REMOVED***
		s += ",def=" + p.Default
	***REMOVED***
	return s
***REMOVED***

// Parse populates p by parsing a string in the protobuf struct field tag style.
func (p *Properties) Parse(s string) ***REMOVED***
	// "bytes,49,opt,name=foo,def=hello!"
	fields := strings.Split(s, ",") // breaks def=, but handled below.
	if len(fields) < 2 ***REMOVED***
		fmt.Fprintf(os.Stderr, "proto: tag has too few fields: %q\n", s)
		return
	***REMOVED***

	p.Wire = fields[0]
	switch p.Wire ***REMOVED***
	case "varint":
		p.WireType = WireVarint
		p.valEnc = (*Buffer).EncodeVarint
		p.valDec = (*Buffer).DecodeVarint
		p.valSize = sizeVarint
	case "fixed32":
		p.WireType = WireFixed32
		p.valEnc = (*Buffer).EncodeFixed32
		p.valDec = (*Buffer).DecodeFixed32
		p.valSize = sizeFixed32
	case "fixed64":
		p.WireType = WireFixed64
		p.valEnc = (*Buffer).EncodeFixed64
		p.valDec = (*Buffer).DecodeFixed64
		p.valSize = sizeFixed64
	case "zigzag32":
		p.WireType = WireVarint
		p.valEnc = (*Buffer).EncodeZigzag32
		p.valDec = (*Buffer).DecodeZigzag32
		p.valSize = sizeZigzag32
	case "zigzag64":
		p.WireType = WireVarint
		p.valEnc = (*Buffer).EncodeZigzag64
		p.valDec = (*Buffer).DecodeZigzag64
		p.valSize = sizeZigzag64
	case "bytes", "group":
		p.WireType = WireBytes
		// no numeric converter for non-numeric types
	default:
		fmt.Fprintf(os.Stderr, "proto: tag has unknown wire type: %q\n", s)
		return
	***REMOVED***

	var err error
	p.Tag, err = strconv.Atoi(fields[1])
	if err != nil ***REMOVED***
		return
	***REMOVED***

	for i := 2; i < len(fields); i++ ***REMOVED***
		f := fields[i]
		switch ***REMOVED***
		case f == "req":
			p.Required = true
		case f == "opt":
			p.Optional = true
		case f == "rep":
			p.Repeated = true
		case f == "packed":
			p.Packed = true
		case strings.HasPrefix(f, "name="):
			p.OrigName = f[5:]
		case strings.HasPrefix(f, "json="):
			p.JSONName = f[5:]
		case strings.HasPrefix(f, "enum="):
			p.Enum = f[5:]
		case f == "proto3":
			p.proto3 = true
		case f == "oneof":
			p.oneof = true
		case strings.HasPrefix(f, "def="):
			p.HasDefault = true
			p.Default = f[4:] // rest of string
			if i+1 < len(fields) ***REMOVED***
				// Commas aren't escaped, and def is always last.
				p.Default += "," + strings.Join(fields[i+1:], ",")
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func logNoSliceEnc(t1, t2 reflect.Type) ***REMOVED***
	fmt.Fprintf(os.Stderr, "proto: no slice oenc for %T = []%T\n", t1, t2)
***REMOVED***

var protoMessageType = reflect.TypeOf((*Message)(nil)).Elem()

// Initialize the fields for encoding and decoding.
func (p *Properties) setEncAndDec(typ reflect.Type, f *reflect.StructField, lockGetProp bool) ***REMOVED***
	p.enc = nil
	p.dec = nil
	p.size = nil

	switch t1 := typ; t1.Kind() ***REMOVED***
	default:
		fmt.Fprintf(os.Stderr, "proto: no coders for %v\n", t1)

	// proto3 scalar types

	case reflect.Bool:
		p.enc = (*Buffer).enc_proto3_bool
		p.dec = (*Buffer).dec_proto3_bool
		p.size = size_proto3_bool
	case reflect.Int32:
		p.enc = (*Buffer).enc_proto3_int32
		p.dec = (*Buffer).dec_proto3_int32
		p.size = size_proto3_int32
	case reflect.Uint32:
		p.enc = (*Buffer).enc_proto3_uint32
		p.dec = (*Buffer).dec_proto3_int32 // can reuse
		p.size = size_proto3_uint32
	case reflect.Int64, reflect.Uint64:
		p.enc = (*Buffer).enc_proto3_int64
		p.dec = (*Buffer).dec_proto3_int64
		p.size = size_proto3_int64
	case reflect.Float32:
		p.enc = (*Buffer).enc_proto3_uint32 // can just treat them as bits
		p.dec = (*Buffer).dec_proto3_int32
		p.size = size_proto3_uint32
	case reflect.Float64:
		p.enc = (*Buffer).enc_proto3_int64 // can just treat them as bits
		p.dec = (*Buffer).dec_proto3_int64
		p.size = size_proto3_int64
	case reflect.String:
		p.enc = (*Buffer).enc_proto3_string
		p.dec = (*Buffer).dec_proto3_string
		p.size = size_proto3_string

	case reflect.Ptr:
		switch t2 := t1.Elem(); t2.Kind() ***REMOVED***
		default:
			fmt.Fprintf(os.Stderr, "proto: no encoder function for %v -> %v\n", t1, t2)
			break
		case reflect.Bool:
			p.enc = (*Buffer).enc_bool
			p.dec = (*Buffer).dec_bool
			p.size = size_bool
		case reflect.Int32:
			p.enc = (*Buffer).enc_int32
			p.dec = (*Buffer).dec_int32
			p.size = size_int32
		case reflect.Uint32:
			p.enc = (*Buffer).enc_uint32
			p.dec = (*Buffer).dec_int32 // can reuse
			p.size = size_uint32
		case reflect.Int64, reflect.Uint64:
			p.enc = (*Buffer).enc_int64
			p.dec = (*Buffer).dec_int64
			p.size = size_int64
		case reflect.Float32:
			p.enc = (*Buffer).enc_uint32 // can just treat them as bits
			p.dec = (*Buffer).dec_int32
			p.size = size_uint32
		case reflect.Float64:
			p.enc = (*Buffer).enc_int64 // can just treat them as bits
			p.dec = (*Buffer).dec_int64
			p.size = size_int64
		case reflect.String:
			p.enc = (*Buffer).enc_string
			p.dec = (*Buffer).dec_string
			p.size = size_string
		case reflect.Struct:
			p.stype = t1.Elem()
			p.isMarshaler = isMarshaler(t1)
			p.isUnmarshaler = isUnmarshaler(t1)
			if p.Wire == "bytes" ***REMOVED***
				p.enc = (*Buffer).enc_struct_message
				p.dec = (*Buffer).dec_struct_message
				p.size = size_struct_message
			***REMOVED*** else ***REMOVED***
				p.enc = (*Buffer).enc_struct_group
				p.dec = (*Buffer).dec_struct_group
				p.size = size_struct_group
			***REMOVED***
		***REMOVED***

	case reflect.Slice:
		switch t2 := t1.Elem(); t2.Kind() ***REMOVED***
		default:
			logNoSliceEnc(t1, t2)
			break
		case reflect.Bool:
			if p.Packed ***REMOVED***
				p.enc = (*Buffer).enc_slice_packed_bool
				p.size = size_slice_packed_bool
			***REMOVED*** else ***REMOVED***
				p.enc = (*Buffer).enc_slice_bool
				p.size = size_slice_bool
			***REMOVED***
			p.dec = (*Buffer).dec_slice_bool
			p.packedDec = (*Buffer).dec_slice_packed_bool
		case reflect.Int32:
			if p.Packed ***REMOVED***
				p.enc = (*Buffer).enc_slice_packed_int32
				p.size = size_slice_packed_int32
			***REMOVED*** else ***REMOVED***
				p.enc = (*Buffer).enc_slice_int32
				p.size = size_slice_int32
			***REMOVED***
			p.dec = (*Buffer).dec_slice_int32
			p.packedDec = (*Buffer).dec_slice_packed_int32
		case reflect.Uint32:
			if p.Packed ***REMOVED***
				p.enc = (*Buffer).enc_slice_packed_uint32
				p.size = size_slice_packed_uint32
			***REMOVED*** else ***REMOVED***
				p.enc = (*Buffer).enc_slice_uint32
				p.size = size_slice_uint32
			***REMOVED***
			p.dec = (*Buffer).dec_slice_int32
			p.packedDec = (*Buffer).dec_slice_packed_int32
		case reflect.Int64, reflect.Uint64:
			if p.Packed ***REMOVED***
				p.enc = (*Buffer).enc_slice_packed_int64
				p.size = size_slice_packed_int64
			***REMOVED*** else ***REMOVED***
				p.enc = (*Buffer).enc_slice_int64
				p.size = size_slice_int64
			***REMOVED***
			p.dec = (*Buffer).dec_slice_int64
			p.packedDec = (*Buffer).dec_slice_packed_int64
		case reflect.Uint8:
			p.dec = (*Buffer).dec_slice_byte
			if p.proto3 ***REMOVED***
				p.enc = (*Buffer).enc_proto3_slice_byte
				p.size = size_proto3_slice_byte
			***REMOVED*** else ***REMOVED***
				p.enc = (*Buffer).enc_slice_byte
				p.size = size_slice_byte
			***REMOVED***
		case reflect.Float32, reflect.Float64:
			switch t2.Bits() ***REMOVED***
			case 32:
				// can just treat them as bits
				if p.Packed ***REMOVED***
					p.enc = (*Buffer).enc_slice_packed_uint32
					p.size = size_slice_packed_uint32
				***REMOVED*** else ***REMOVED***
					p.enc = (*Buffer).enc_slice_uint32
					p.size = size_slice_uint32
				***REMOVED***
				p.dec = (*Buffer).dec_slice_int32
				p.packedDec = (*Buffer).dec_slice_packed_int32
			case 64:
				// can just treat them as bits
				if p.Packed ***REMOVED***
					p.enc = (*Buffer).enc_slice_packed_int64
					p.size = size_slice_packed_int64
				***REMOVED*** else ***REMOVED***
					p.enc = (*Buffer).enc_slice_int64
					p.size = size_slice_int64
				***REMOVED***
				p.dec = (*Buffer).dec_slice_int64
				p.packedDec = (*Buffer).dec_slice_packed_int64
			default:
				logNoSliceEnc(t1, t2)
				break
			***REMOVED***
		case reflect.String:
			p.enc = (*Buffer).enc_slice_string
			p.dec = (*Buffer).dec_slice_string
			p.size = size_slice_string
		case reflect.Ptr:
			switch t3 := t2.Elem(); t3.Kind() ***REMOVED***
			default:
				fmt.Fprintf(os.Stderr, "proto: no ptr oenc for %T -> %T -> %T\n", t1, t2, t3)
				break
			case reflect.Struct:
				p.stype = t2.Elem()
				p.isMarshaler = isMarshaler(t2)
				p.isUnmarshaler = isUnmarshaler(t2)
				if p.Wire == "bytes" ***REMOVED***
					p.enc = (*Buffer).enc_slice_struct_message
					p.dec = (*Buffer).dec_slice_struct_message
					p.size = size_slice_struct_message
				***REMOVED*** else ***REMOVED***
					p.enc = (*Buffer).enc_slice_struct_group
					p.dec = (*Buffer).dec_slice_struct_group
					p.size = size_slice_struct_group
				***REMOVED***
			***REMOVED***
		case reflect.Slice:
			switch t2.Elem().Kind() ***REMOVED***
			default:
				fmt.Fprintf(os.Stderr, "proto: no slice elem oenc for %T -> %T -> %T\n", t1, t2, t2.Elem())
				break
			case reflect.Uint8:
				p.enc = (*Buffer).enc_slice_slice_byte
				p.dec = (*Buffer).dec_slice_slice_byte
				p.size = size_slice_slice_byte
			***REMOVED***
		***REMOVED***

	case reflect.Map:
		p.enc = (*Buffer).enc_new_map
		p.dec = (*Buffer).dec_new_map
		p.size = size_new_map

		p.mtype = t1
		p.mkeyprop = &Properties***REMOVED******REMOVED***
		p.mkeyprop.init(reflect.PtrTo(p.mtype.Key()), "Key", f.Tag.Get("protobuf_key"), nil, lockGetProp)
		p.mvalprop = &Properties***REMOVED******REMOVED***
		vtype := p.mtype.Elem()
		if vtype.Kind() != reflect.Ptr && vtype.Kind() != reflect.Slice ***REMOVED***
			// The value type is not a message (*T) or bytes ([]byte),
			// so we need encoders for the pointer to this type.
			vtype = reflect.PtrTo(vtype)
		***REMOVED***
		p.mvalprop.init(vtype, "Value", f.Tag.Get("protobuf_val"), nil, lockGetProp)
	***REMOVED***

	// precalculate tag code
	wire := p.WireType
	if p.Packed ***REMOVED***
		wire = WireBytes
	***REMOVED***
	x := uint32(p.Tag)<<3 | uint32(wire)
	i := 0
	for i = 0; x > 127; i++ ***REMOVED***
		p.tagbuf[i] = 0x80 | uint8(x&0x7F)
		x >>= 7
	***REMOVED***
	p.tagbuf[i] = uint8(x)
	p.tagcode = p.tagbuf[0 : i+1]

	if p.stype != nil ***REMOVED***
		if lockGetProp ***REMOVED***
			p.sprop = GetProperties(p.stype)
		***REMOVED*** else ***REMOVED***
			p.sprop = getPropertiesLocked(p.stype)
		***REMOVED***
	***REMOVED***
***REMOVED***

var (
	marshalerType   = reflect.TypeOf((*Marshaler)(nil)).Elem()
	unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

// isMarshaler reports whether type t implements Marshaler.
func isMarshaler(t reflect.Type) bool ***REMOVED***
	// We're checking for (likely) pointer-receiver methods
	// so if t is not a pointer, something is very wrong.
	// The calls above only invoke isMarshaler on pointer types.
	if t.Kind() != reflect.Ptr ***REMOVED***
		panic("proto: misuse of isMarshaler")
	***REMOVED***
	return t.Implements(marshalerType)
***REMOVED***

// isUnmarshaler reports whether type t implements Unmarshaler.
func isUnmarshaler(t reflect.Type) bool ***REMOVED***
	// We're checking for (likely) pointer-receiver methods
	// so if t is not a pointer, something is very wrong.
	// The calls above only invoke isUnmarshaler on pointer types.
	if t.Kind() != reflect.Ptr ***REMOVED***
		panic("proto: misuse of isUnmarshaler")
	***REMOVED***
	return t.Implements(unmarshalerType)
***REMOVED***

// Init populates the properties from a protocol buffer struct tag.
func (p *Properties) Init(typ reflect.Type, name, tag string, f *reflect.StructField) ***REMOVED***
	p.init(typ, name, tag, f, true)
***REMOVED***

func (p *Properties) init(typ reflect.Type, name, tag string, f *reflect.StructField, lockGetProp bool) ***REMOVED***
	// "bytes,49,opt,def=hello!"
	p.Name = name
	p.OrigName = name
	if f != nil ***REMOVED***
		p.field = toField(f)
	***REMOVED***
	if tag == "" ***REMOVED***
		return
	***REMOVED***
	p.Parse(tag)
	p.setEncAndDec(typ, f, lockGetProp)
***REMOVED***

var (
	propertiesMu  sync.RWMutex
	propertiesMap = make(map[reflect.Type]*StructProperties)
)

// GetProperties returns the list of properties for the type represented by t.
// t must represent a generated struct type of a protocol message.
func GetProperties(t reflect.Type) *StructProperties ***REMOVED***
	if t.Kind() != reflect.Struct ***REMOVED***
		panic("proto: type must have kind struct")
	***REMOVED***

	// Most calls to GetProperties in a long-running program will be
	// retrieving details for types we have seen before.
	propertiesMu.RLock()
	sprop, ok := propertiesMap[t]
	propertiesMu.RUnlock()
	if ok ***REMOVED***
		if collectStats ***REMOVED***
			stats.Chit++
		***REMOVED***
		return sprop
	***REMOVED***

	propertiesMu.Lock()
	sprop = getPropertiesLocked(t)
	propertiesMu.Unlock()
	return sprop
***REMOVED***

// getPropertiesLocked requires that propertiesMu is held.
func getPropertiesLocked(t reflect.Type) *StructProperties ***REMOVED***
	if prop, ok := propertiesMap[t]; ok ***REMOVED***
		if collectStats ***REMOVED***
			stats.Chit++
		***REMOVED***
		return prop
	***REMOVED***
	if collectStats ***REMOVED***
		stats.Cmiss++
	***REMOVED***

	prop := new(StructProperties)
	// in case of recursive protos, fill this in now.
	propertiesMap[t] = prop

	// build properties
	prop.extendable = reflect.PtrTo(t).Implements(extendableProtoType) ||
		reflect.PtrTo(t).Implements(extendableProtoV1Type)
	prop.unrecField = invalidField
	prop.Prop = make([]*Properties, t.NumField())
	prop.order = make([]int, t.NumField())

	for i := 0; i < t.NumField(); i++ ***REMOVED***
		f := t.Field(i)
		p := new(Properties)
		name := f.Name
		p.init(f.Type, name, f.Tag.Get("protobuf"), &f, false)

		if f.Name == "XXX_InternalExtensions" ***REMOVED*** // special case
			p.enc = (*Buffer).enc_exts
			p.dec = nil // not needed
			p.size = size_exts
		***REMOVED*** else if f.Name == "XXX_extensions" ***REMOVED*** // special case
			p.enc = (*Buffer).enc_map
			p.dec = nil // not needed
			p.size = size_map
		***REMOVED*** else if f.Name == "XXX_unrecognized" ***REMOVED*** // special case
			prop.unrecField = toField(&f)
		***REMOVED***
		oneof := f.Tag.Get("protobuf_oneof") // special case
		if oneof != "" ***REMOVED***
			// Oneof fields don't use the traditional protobuf tag.
			p.OrigName = oneof
		***REMOVED***
		prop.Prop[i] = p
		prop.order[i] = i
		if debug ***REMOVED***
			print(i, " ", f.Name, " ", t.String(), " ")
			if p.Tag > 0 ***REMOVED***
				print(p.String())
			***REMOVED***
			print("\n")
		***REMOVED***
		if p.enc == nil && !strings.HasPrefix(f.Name, "XXX_") && oneof == "" ***REMOVED***
			fmt.Fprintln(os.Stderr, "proto: no encoder for", f.Name, f.Type.String(), "[GetProperties]")
		***REMOVED***
	***REMOVED***

	// Re-order prop.order.
	sort.Sort(prop)

	type oneofMessage interface ***REMOVED***
		XXX_OneofFuncs() (func(Message, *Buffer) error, func(Message, int, int, *Buffer) (bool, error), func(Message) int, []interface***REMOVED******REMOVED***)
	***REMOVED***
	if om, ok := reflect.Zero(reflect.PtrTo(t)).Interface().(oneofMessage); ok ***REMOVED***
		var oots []interface***REMOVED******REMOVED***
		prop.oneofMarshaler, prop.oneofUnmarshaler, prop.oneofSizer, oots = om.XXX_OneofFuncs()
		prop.stype = t

		// Interpret oneof metadata.
		prop.OneofTypes = make(map[string]*OneofProperties)
		for _, oot := range oots ***REMOVED***
			oop := &OneofProperties***REMOVED***
				Type: reflect.ValueOf(oot).Type(), // *T
				Prop: new(Properties),
			***REMOVED***
			sft := oop.Type.Elem().Field(0)
			oop.Prop.Name = sft.Name
			oop.Prop.Parse(sft.Tag.Get("protobuf"))
			// There will be exactly one interface field that
			// this new value is assignable to.
			for i := 0; i < t.NumField(); i++ ***REMOVED***
				f := t.Field(i)
				if f.Type.Kind() != reflect.Interface ***REMOVED***
					continue
				***REMOVED***
				if !oop.Type.AssignableTo(f.Type) ***REMOVED***
					continue
				***REMOVED***
				oop.Field = i
				break
			***REMOVED***
			prop.OneofTypes[oop.Prop.OrigName] = oop
		***REMOVED***
	***REMOVED***

	// build required counts
	// build tags
	reqCount := 0
	prop.decoderOrigNames = make(map[string]int)
	for i, p := range prop.Prop ***REMOVED***
		if strings.HasPrefix(p.Name, "XXX_") ***REMOVED***
			// Internal fields should not appear in tags/origNames maps.
			// They are handled specially when encoding and decoding.
			continue
		***REMOVED***
		if p.Required ***REMOVED***
			reqCount++
		***REMOVED***
		prop.decoderTags.put(p.Tag, i)
		prop.decoderOrigNames[p.OrigName] = i
	***REMOVED***
	prop.reqCount = reqCount

	return prop
***REMOVED***

// Return the Properties object for the x[0]'th field of the structure.
func propByIndex(t reflect.Type, x []int) *Properties ***REMOVED***
	if len(x) != 1 ***REMOVED***
		fmt.Fprintf(os.Stderr, "proto: field index dimension %d (not 1) for type %s\n", len(x), t)
		return nil
	***REMOVED***
	prop := GetProperties(t)
	return prop.Prop[x[0]]
***REMOVED***

// Get the address and type of a pointer to a struct from an interface.
func getbase(pb Message) (t reflect.Type, b structPointer, err error) ***REMOVED***
	if pb == nil ***REMOVED***
		err = ErrNil
		return
	***REMOVED***
	// get the reflect type of the pointer to the struct.
	t = reflect.TypeOf(pb)
	// get the address of the struct.
	value := reflect.ValueOf(pb)
	b = toStructPointer(value)
	return
***REMOVED***

// A global registry of enum types.
// The generated code will register the generated maps by calling RegisterEnum.

var enumValueMaps = make(map[string]map[string]int32)

// RegisterEnum is called from the generated code to install the enum descriptor
// maps into the global table to aid parsing text format protocol buffers.
func RegisterEnum(typeName string, unusedNameMap map[int32]string, valueMap map[string]int32) ***REMOVED***
	if _, ok := enumValueMaps[typeName]; ok ***REMOVED***
		panic("proto: duplicate enum registered: " + typeName)
	***REMOVED***
	enumValueMaps[typeName] = valueMap
***REMOVED***

// EnumValueMap returns the mapping from names to integers of the
// enum type enumType, or a nil if not found.
func EnumValueMap(enumType string) map[string]int32 ***REMOVED***
	return enumValueMaps[enumType]
***REMOVED***

// A registry of all linked message types.
// The string is a fully-qualified proto name ("pkg.Message").
var (
	protoTypes    = make(map[string]reflect.Type)
	revProtoTypes = make(map[reflect.Type]string)
)

// RegisterType is called from generated code and maps from the fully qualified
// proto name to the type (pointer to struct) of the protocol buffer.
func RegisterType(x Message, name string) ***REMOVED***
	if _, ok := protoTypes[name]; ok ***REMOVED***
		// TODO: Some day, make this a panic.
		log.Printf("proto: duplicate proto type registered: %s", name)
		return
	***REMOVED***
	t := reflect.TypeOf(x)
	protoTypes[name] = t
	revProtoTypes[t] = name
***REMOVED***

// MessageName returns the fully-qualified proto name for the given message type.
func MessageName(x Message) string ***REMOVED***
	type xname interface ***REMOVED***
		XXX_MessageName() string
	***REMOVED***
	if m, ok := x.(xname); ok ***REMOVED***
		return m.XXX_MessageName()
	***REMOVED***
	return revProtoTypes[reflect.TypeOf(x)]
***REMOVED***

// MessageType returns the message type (pointer to struct) for a named message.
func MessageType(name string) reflect.Type ***REMOVED*** return protoTypes[name] ***REMOVED***

// A registry of all linked proto files.
var (
	protoFiles = make(map[string][]byte) // file name => fileDescriptor
)

// RegisterFile is called from generated code and maps from the
// full file name of a .proto file to its compressed FileDescriptorProto.
func RegisterFile(filename string, fileDescriptor []byte) ***REMOVED***
	protoFiles[filename] = fileDescriptor
***REMOVED***

// FileDescriptor returns the compressed FileDescriptorProto for a .proto file.
func FileDescriptor(filename string) []byte ***REMOVED*** return protoFiles[filename] ***REMOVED***
