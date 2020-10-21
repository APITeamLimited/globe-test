// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/strs"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// pointerCoderFuncs is a set of pointer encoding functions.
type pointerCoderFuncs struct ***REMOVED***
	mi        *MessageInfo
	size      func(p pointer, f *coderFieldInfo, opts marshalOptions) int
	marshal   func(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error)
	unmarshal func(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (unmarshalOutput, error)
	isInit    func(p pointer, f *coderFieldInfo) error
	merge     func(dst, src pointer, f *coderFieldInfo, opts mergeOptions)
***REMOVED***

// valueCoderFuncs is a set of protoreflect.Value encoding functions.
type valueCoderFuncs struct ***REMOVED***
	size      func(v pref.Value, tagsize int, opts marshalOptions) int
	marshal   func(b []byte, v pref.Value, wiretag uint64, opts marshalOptions) ([]byte, error)
	unmarshal func(b []byte, v pref.Value, num protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (pref.Value, unmarshalOutput, error)
	isInit    func(v pref.Value) error
	merge     func(dst, src pref.Value, opts mergeOptions) pref.Value
***REMOVED***

// fieldCoder returns pointer functions for a field, used for operating on
// struct fields.
func fieldCoder(fd pref.FieldDescriptor, ft reflect.Type) (*MessageInfo, pointerCoderFuncs) ***REMOVED***
	switch ***REMOVED***
	case fd.IsMap():
		return encoderFuncsForMap(fd, ft)
	case fd.Cardinality() == pref.Repeated && !fd.IsPacked():
		// Repeated fields (not packed).
		if ft.Kind() != reflect.Slice ***REMOVED***
			break
		***REMOVED***
		ft := ft.Elem()
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBoolSlice
			***REMOVED***
		case pref.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnumSlice
			***REMOVED***
		case pref.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32Slice
			***REMOVED***
		case pref.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32Slice
			***REMOVED***
		case pref.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32Slice
			***REMOVED***
		case pref.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64Slice
			***REMOVED***
		case pref.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64Slice
			***REMOVED***
		case pref.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64Slice
			***REMOVED***
		case pref.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32Slice
			***REMOVED***
		case pref.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32Slice
			***REMOVED***
		case pref.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloatSlice
			***REMOVED***
		case pref.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64Slice
			***REMOVED***
		case pref.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64Slice
			***REMOVED***
		case pref.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDoubleSlice
			***REMOVED***
		case pref.StringKind:
			if ft.Kind() == reflect.String && strs.EnforceUTF8(fd) ***REMOVED***
				return nil, coderStringSliceValidateUTF8
			***REMOVED***
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringSlice
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 && strs.EnforceUTF8(fd) ***REMOVED***
				return nil, coderBytesSliceValidateUTF8
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 ***REMOVED***
				return nil, coderBytesSlice
			***REMOVED***
		case pref.BytesKind:
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringSlice
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 ***REMOVED***
				return nil, coderBytesSlice
			***REMOVED***
		case pref.MessageKind:
			return getMessageInfo(ft), makeMessageSliceFieldCoder(fd, ft)
		case pref.GroupKind:
			return getMessageInfo(ft), makeGroupSliceFieldCoder(fd, ft)
		***REMOVED***
	case fd.Cardinality() == pref.Repeated && fd.IsPacked():
		// Packed repeated fields.
		//
		// Only repeated fields of primitive numeric types
		// (Varint, Fixed32, or Fixed64 wire type) can be packed.
		if ft.Kind() != reflect.Slice ***REMOVED***
			break
		***REMOVED***
		ft := ft.Elem()
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBoolPackedSlice
			***REMOVED***
		case pref.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnumPackedSlice
			***REMOVED***
		case pref.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32PackedSlice
			***REMOVED***
		case pref.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32PackedSlice
			***REMOVED***
		case pref.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32PackedSlice
			***REMOVED***
		case pref.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64PackedSlice
			***REMOVED***
		case pref.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64PackedSlice
			***REMOVED***
		case pref.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64PackedSlice
			***REMOVED***
		case pref.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32PackedSlice
			***REMOVED***
		case pref.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32PackedSlice
			***REMOVED***
		case pref.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloatPackedSlice
			***REMOVED***
		case pref.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64PackedSlice
			***REMOVED***
		case pref.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64PackedSlice
			***REMOVED***
		case pref.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDoublePackedSlice
			***REMOVED***
		***REMOVED***
	case fd.Kind() == pref.MessageKind:
		return getMessageInfo(ft), makeMessageFieldCoder(fd, ft)
	case fd.Kind() == pref.GroupKind:
		return getMessageInfo(ft), makeGroupFieldCoder(fd, ft)
	case fd.Syntax() == pref.Proto3 && fd.ContainingOneof() == nil:
		// Populated oneof fields always encode even if set to the zero value,
		// which normally are not encoded in proto3.
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBoolNoZero
			***REMOVED***
		case pref.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnumNoZero
			***REMOVED***
		case pref.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32NoZero
			***REMOVED***
		case pref.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32NoZero
			***REMOVED***
		case pref.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32NoZero
			***REMOVED***
		case pref.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64NoZero
			***REMOVED***
		case pref.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64NoZero
			***REMOVED***
		case pref.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64NoZero
			***REMOVED***
		case pref.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32NoZero
			***REMOVED***
		case pref.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32NoZero
			***REMOVED***
		case pref.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloatNoZero
			***REMOVED***
		case pref.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64NoZero
			***REMOVED***
		case pref.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64NoZero
			***REMOVED***
		case pref.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDoubleNoZero
			***REMOVED***
		case pref.StringKind:
			if ft.Kind() == reflect.String && strs.EnforceUTF8(fd) ***REMOVED***
				return nil, coderStringNoZeroValidateUTF8
			***REMOVED***
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringNoZero
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 && strs.EnforceUTF8(fd) ***REMOVED***
				return nil, coderBytesNoZeroValidateUTF8
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 ***REMOVED***
				return nil, coderBytesNoZero
			***REMOVED***
		case pref.BytesKind:
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringNoZero
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 ***REMOVED***
				return nil, coderBytesNoZero
			***REMOVED***
		***REMOVED***
	case ft.Kind() == reflect.Ptr:
		ft := ft.Elem()
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBoolPtr
			***REMOVED***
		case pref.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnumPtr
			***REMOVED***
		case pref.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32Ptr
			***REMOVED***
		case pref.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32Ptr
			***REMOVED***
		case pref.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32Ptr
			***REMOVED***
		case pref.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64Ptr
			***REMOVED***
		case pref.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64Ptr
			***REMOVED***
		case pref.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64Ptr
			***REMOVED***
		case pref.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32Ptr
			***REMOVED***
		case pref.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32Ptr
			***REMOVED***
		case pref.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloatPtr
			***REMOVED***
		case pref.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64Ptr
			***REMOVED***
		case pref.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64Ptr
			***REMOVED***
		case pref.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDoublePtr
			***REMOVED***
		case pref.StringKind:
			if ft.Kind() == reflect.String && strs.EnforceUTF8(fd) ***REMOVED***
				return nil, coderStringPtrValidateUTF8
			***REMOVED***
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringPtr
			***REMOVED***
		case pref.BytesKind:
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringPtr
			***REMOVED***
		***REMOVED***
	default:
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBool
			***REMOVED***
		case pref.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnum
			***REMOVED***
		case pref.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32
			***REMOVED***
		case pref.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32
			***REMOVED***
		case pref.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32
			***REMOVED***
		case pref.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64
			***REMOVED***
		case pref.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64
			***REMOVED***
		case pref.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64
			***REMOVED***
		case pref.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32
			***REMOVED***
		case pref.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32
			***REMOVED***
		case pref.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloat
			***REMOVED***
		case pref.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64
			***REMOVED***
		case pref.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64
			***REMOVED***
		case pref.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDouble
			***REMOVED***
		case pref.StringKind:
			if ft.Kind() == reflect.String && strs.EnforceUTF8(fd) ***REMOVED***
				return nil, coderStringValidateUTF8
			***REMOVED***
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderString
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 && strs.EnforceUTF8(fd) ***REMOVED***
				return nil, coderBytesValidateUTF8
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 ***REMOVED***
				return nil, coderBytes
			***REMOVED***
		case pref.BytesKind:
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderString
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 ***REMOVED***
				return nil, coderBytes
			***REMOVED***
		***REMOVED***
	***REMOVED***
	panic(fmt.Sprintf("invalid type: no encoder for %v %v %v/%v", fd.FullName(), fd.Cardinality(), fd.Kind(), ft))
***REMOVED***

// encoderFuncsForValue returns value functions for a field, used for
// extension values and map encoding.
func encoderFuncsForValue(fd pref.FieldDescriptor) valueCoderFuncs ***REMOVED***
	switch ***REMOVED***
	case fd.Cardinality() == pref.Repeated && !fd.IsPacked():
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			return coderBoolSliceValue
		case pref.EnumKind:
			return coderEnumSliceValue
		case pref.Int32Kind:
			return coderInt32SliceValue
		case pref.Sint32Kind:
			return coderSint32SliceValue
		case pref.Uint32Kind:
			return coderUint32SliceValue
		case pref.Int64Kind:
			return coderInt64SliceValue
		case pref.Sint64Kind:
			return coderSint64SliceValue
		case pref.Uint64Kind:
			return coderUint64SliceValue
		case pref.Sfixed32Kind:
			return coderSfixed32SliceValue
		case pref.Fixed32Kind:
			return coderFixed32SliceValue
		case pref.FloatKind:
			return coderFloatSliceValue
		case pref.Sfixed64Kind:
			return coderSfixed64SliceValue
		case pref.Fixed64Kind:
			return coderFixed64SliceValue
		case pref.DoubleKind:
			return coderDoubleSliceValue
		case pref.StringKind:
			// We don't have a UTF-8 validating coder for repeated string fields.
			// Value coders are used for extensions and maps.
			// Extensions are never proto3, and maps never contain lists.
			return coderStringSliceValue
		case pref.BytesKind:
			return coderBytesSliceValue
		case pref.MessageKind:
			return coderMessageSliceValue
		case pref.GroupKind:
			return coderGroupSliceValue
		***REMOVED***
	case fd.Cardinality() == pref.Repeated && fd.IsPacked():
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			return coderBoolPackedSliceValue
		case pref.EnumKind:
			return coderEnumPackedSliceValue
		case pref.Int32Kind:
			return coderInt32PackedSliceValue
		case pref.Sint32Kind:
			return coderSint32PackedSliceValue
		case pref.Uint32Kind:
			return coderUint32PackedSliceValue
		case pref.Int64Kind:
			return coderInt64PackedSliceValue
		case pref.Sint64Kind:
			return coderSint64PackedSliceValue
		case pref.Uint64Kind:
			return coderUint64PackedSliceValue
		case pref.Sfixed32Kind:
			return coderSfixed32PackedSliceValue
		case pref.Fixed32Kind:
			return coderFixed32PackedSliceValue
		case pref.FloatKind:
			return coderFloatPackedSliceValue
		case pref.Sfixed64Kind:
			return coderSfixed64PackedSliceValue
		case pref.Fixed64Kind:
			return coderFixed64PackedSliceValue
		case pref.DoubleKind:
			return coderDoublePackedSliceValue
		***REMOVED***
	default:
		switch fd.Kind() ***REMOVED***
		default:
		case pref.BoolKind:
			return coderBoolValue
		case pref.EnumKind:
			return coderEnumValue
		case pref.Int32Kind:
			return coderInt32Value
		case pref.Sint32Kind:
			return coderSint32Value
		case pref.Uint32Kind:
			return coderUint32Value
		case pref.Int64Kind:
			return coderInt64Value
		case pref.Sint64Kind:
			return coderSint64Value
		case pref.Uint64Kind:
			return coderUint64Value
		case pref.Sfixed32Kind:
			return coderSfixed32Value
		case pref.Fixed32Kind:
			return coderFixed32Value
		case pref.FloatKind:
			return coderFloatValue
		case pref.Sfixed64Kind:
			return coderSfixed64Value
		case pref.Fixed64Kind:
			return coderFixed64Value
		case pref.DoubleKind:
			return coderDoubleValue
		case pref.StringKind:
			if strs.EnforceUTF8(fd) ***REMOVED***
				return coderStringValueValidateUTF8
			***REMOVED***
			return coderStringValue
		case pref.BytesKind:
			return coderBytesValue
		case pref.MessageKind:
			return coderMessageValue
		case pref.GroupKind:
			return coderGroupValue
		***REMOVED***
	***REMOVED***
	panic(fmt.Sprintf("invalid field: no encoder for %v %v %v", fd.FullName(), fd.Cardinality(), fd.Kind()))
***REMOVED***
