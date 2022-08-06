// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	size      func(v protoreflect.Value, tagsize int, opts marshalOptions) int
	marshal   func(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error)
	unmarshal func(b []byte, v protoreflect.Value, num protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (protoreflect.Value, unmarshalOutput, error)
	isInit    func(v protoreflect.Value) error
	merge     func(dst, src protoreflect.Value, opts mergeOptions) protoreflect.Value
***REMOVED***

// fieldCoder returns pointer functions for a field, used for operating on
// struct fields.
func fieldCoder(fd protoreflect.FieldDescriptor, ft reflect.Type) (*MessageInfo, pointerCoderFuncs) ***REMOVED***
	switch ***REMOVED***
	case fd.IsMap():
		return encoderFuncsForMap(fd, ft)
	case fd.Cardinality() == protoreflect.Repeated && !fd.IsPacked():
		// Repeated fields (not packed).
		if ft.Kind() != reflect.Slice ***REMOVED***
			break
		***REMOVED***
		ft := ft.Elem()
		switch fd.Kind() ***REMOVED***
		case protoreflect.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBoolSlice
			***REMOVED***
		case protoreflect.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnumSlice
			***REMOVED***
		case protoreflect.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32Slice
			***REMOVED***
		case protoreflect.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32Slice
			***REMOVED***
		case protoreflect.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32Slice
			***REMOVED***
		case protoreflect.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64Slice
			***REMOVED***
		case protoreflect.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64Slice
			***REMOVED***
		case protoreflect.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64Slice
			***REMOVED***
		case protoreflect.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32Slice
			***REMOVED***
		case protoreflect.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32Slice
			***REMOVED***
		case protoreflect.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloatSlice
			***REMOVED***
		case protoreflect.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64Slice
			***REMOVED***
		case protoreflect.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64Slice
			***REMOVED***
		case protoreflect.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDoubleSlice
			***REMOVED***
		case protoreflect.StringKind:
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
		case protoreflect.BytesKind:
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringSlice
			***REMOVED***
			if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8 ***REMOVED***
				return nil, coderBytesSlice
			***REMOVED***
		case protoreflect.MessageKind:
			return getMessageInfo(ft), makeMessageSliceFieldCoder(fd, ft)
		case protoreflect.GroupKind:
			return getMessageInfo(ft), makeGroupSliceFieldCoder(fd, ft)
		***REMOVED***
	case fd.Cardinality() == protoreflect.Repeated && fd.IsPacked():
		// Packed repeated fields.
		//
		// Only repeated fields of primitive numeric types
		// (Varint, Fixed32, or Fixed64 wire type) can be packed.
		if ft.Kind() != reflect.Slice ***REMOVED***
			break
		***REMOVED***
		ft := ft.Elem()
		switch fd.Kind() ***REMOVED***
		case protoreflect.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBoolPackedSlice
			***REMOVED***
		case protoreflect.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnumPackedSlice
			***REMOVED***
		case protoreflect.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32PackedSlice
			***REMOVED***
		case protoreflect.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32PackedSlice
			***REMOVED***
		case protoreflect.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32PackedSlice
			***REMOVED***
		case protoreflect.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64PackedSlice
			***REMOVED***
		case protoreflect.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64PackedSlice
			***REMOVED***
		case protoreflect.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64PackedSlice
			***REMOVED***
		case protoreflect.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32PackedSlice
			***REMOVED***
		case protoreflect.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32PackedSlice
			***REMOVED***
		case protoreflect.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloatPackedSlice
			***REMOVED***
		case protoreflect.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64PackedSlice
			***REMOVED***
		case protoreflect.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64PackedSlice
			***REMOVED***
		case protoreflect.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDoublePackedSlice
			***REMOVED***
		***REMOVED***
	case fd.Kind() == protoreflect.MessageKind:
		return getMessageInfo(ft), makeMessageFieldCoder(fd, ft)
	case fd.Kind() == protoreflect.GroupKind:
		return getMessageInfo(ft), makeGroupFieldCoder(fd, ft)
	case fd.Syntax() == protoreflect.Proto3 && fd.ContainingOneof() == nil:
		// Populated oneof fields always encode even if set to the zero value,
		// which normally are not encoded in proto3.
		switch fd.Kind() ***REMOVED***
		case protoreflect.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBoolNoZero
			***REMOVED***
		case protoreflect.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnumNoZero
			***REMOVED***
		case protoreflect.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32NoZero
			***REMOVED***
		case protoreflect.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32NoZero
			***REMOVED***
		case protoreflect.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32NoZero
			***REMOVED***
		case protoreflect.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64NoZero
			***REMOVED***
		case protoreflect.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64NoZero
			***REMOVED***
		case protoreflect.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64NoZero
			***REMOVED***
		case protoreflect.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32NoZero
			***REMOVED***
		case protoreflect.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32NoZero
			***REMOVED***
		case protoreflect.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloatNoZero
			***REMOVED***
		case protoreflect.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64NoZero
			***REMOVED***
		case protoreflect.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64NoZero
			***REMOVED***
		case protoreflect.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDoubleNoZero
			***REMOVED***
		case protoreflect.StringKind:
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
		case protoreflect.BytesKind:
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
		case protoreflect.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBoolPtr
			***REMOVED***
		case protoreflect.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnumPtr
			***REMOVED***
		case protoreflect.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32Ptr
			***REMOVED***
		case protoreflect.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32Ptr
			***REMOVED***
		case protoreflect.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32Ptr
			***REMOVED***
		case protoreflect.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64Ptr
			***REMOVED***
		case protoreflect.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64Ptr
			***REMOVED***
		case protoreflect.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64Ptr
			***REMOVED***
		case protoreflect.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32Ptr
			***REMOVED***
		case protoreflect.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32Ptr
			***REMOVED***
		case protoreflect.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloatPtr
			***REMOVED***
		case protoreflect.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64Ptr
			***REMOVED***
		case protoreflect.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64Ptr
			***REMOVED***
		case protoreflect.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDoublePtr
			***REMOVED***
		case protoreflect.StringKind:
			if ft.Kind() == reflect.String && strs.EnforceUTF8(fd) ***REMOVED***
				return nil, coderStringPtrValidateUTF8
			***REMOVED***
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringPtr
			***REMOVED***
		case protoreflect.BytesKind:
			if ft.Kind() == reflect.String ***REMOVED***
				return nil, coderStringPtr
			***REMOVED***
		***REMOVED***
	default:
		switch fd.Kind() ***REMOVED***
		case protoreflect.BoolKind:
			if ft.Kind() == reflect.Bool ***REMOVED***
				return nil, coderBool
			***REMOVED***
		case protoreflect.EnumKind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderEnum
			***REMOVED***
		case protoreflect.Int32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderInt32
			***REMOVED***
		case protoreflect.Sint32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSint32
			***REMOVED***
		case protoreflect.Uint32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderUint32
			***REMOVED***
		case protoreflect.Int64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderInt64
			***REMOVED***
		case protoreflect.Sint64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSint64
			***REMOVED***
		case protoreflect.Uint64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderUint64
			***REMOVED***
		case protoreflect.Sfixed32Kind:
			if ft.Kind() == reflect.Int32 ***REMOVED***
				return nil, coderSfixed32
			***REMOVED***
		case protoreflect.Fixed32Kind:
			if ft.Kind() == reflect.Uint32 ***REMOVED***
				return nil, coderFixed32
			***REMOVED***
		case protoreflect.FloatKind:
			if ft.Kind() == reflect.Float32 ***REMOVED***
				return nil, coderFloat
			***REMOVED***
		case protoreflect.Sfixed64Kind:
			if ft.Kind() == reflect.Int64 ***REMOVED***
				return nil, coderSfixed64
			***REMOVED***
		case protoreflect.Fixed64Kind:
			if ft.Kind() == reflect.Uint64 ***REMOVED***
				return nil, coderFixed64
			***REMOVED***
		case protoreflect.DoubleKind:
			if ft.Kind() == reflect.Float64 ***REMOVED***
				return nil, coderDouble
			***REMOVED***
		case protoreflect.StringKind:
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
		case protoreflect.BytesKind:
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
func encoderFuncsForValue(fd protoreflect.FieldDescriptor) valueCoderFuncs ***REMOVED***
	switch ***REMOVED***
	case fd.Cardinality() == protoreflect.Repeated && !fd.IsPacked():
		switch fd.Kind() ***REMOVED***
		case protoreflect.BoolKind:
			return coderBoolSliceValue
		case protoreflect.EnumKind:
			return coderEnumSliceValue
		case protoreflect.Int32Kind:
			return coderInt32SliceValue
		case protoreflect.Sint32Kind:
			return coderSint32SliceValue
		case protoreflect.Uint32Kind:
			return coderUint32SliceValue
		case protoreflect.Int64Kind:
			return coderInt64SliceValue
		case protoreflect.Sint64Kind:
			return coderSint64SliceValue
		case protoreflect.Uint64Kind:
			return coderUint64SliceValue
		case protoreflect.Sfixed32Kind:
			return coderSfixed32SliceValue
		case protoreflect.Fixed32Kind:
			return coderFixed32SliceValue
		case protoreflect.FloatKind:
			return coderFloatSliceValue
		case protoreflect.Sfixed64Kind:
			return coderSfixed64SliceValue
		case protoreflect.Fixed64Kind:
			return coderFixed64SliceValue
		case protoreflect.DoubleKind:
			return coderDoubleSliceValue
		case protoreflect.StringKind:
			// We don't have a UTF-8 validating coder for repeated string fields.
			// Value coders are used for extensions and maps.
			// Extensions are never proto3, and maps never contain lists.
			return coderStringSliceValue
		case protoreflect.BytesKind:
			return coderBytesSliceValue
		case protoreflect.MessageKind:
			return coderMessageSliceValue
		case protoreflect.GroupKind:
			return coderGroupSliceValue
		***REMOVED***
	case fd.Cardinality() == protoreflect.Repeated && fd.IsPacked():
		switch fd.Kind() ***REMOVED***
		case protoreflect.BoolKind:
			return coderBoolPackedSliceValue
		case protoreflect.EnumKind:
			return coderEnumPackedSliceValue
		case protoreflect.Int32Kind:
			return coderInt32PackedSliceValue
		case protoreflect.Sint32Kind:
			return coderSint32PackedSliceValue
		case protoreflect.Uint32Kind:
			return coderUint32PackedSliceValue
		case protoreflect.Int64Kind:
			return coderInt64PackedSliceValue
		case protoreflect.Sint64Kind:
			return coderSint64PackedSliceValue
		case protoreflect.Uint64Kind:
			return coderUint64PackedSliceValue
		case protoreflect.Sfixed32Kind:
			return coderSfixed32PackedSliceValue
		case protoreflect.Fixed32Kind:
			return coderFixed32PackedSliceValue
		case protoreflect.FloatKind:
			return coderFloatPackedSliceValue
		case protoreflect.Sfixed64Kind:
			return coderSfixed64PackedSliceValue
		case protoreflect.Fixed64Kind:
			return coderFixed64PackedSliceValue
		case protoreflect.DoubleKind:
			return coderDoublePackedSliceValue
		***REMOVED***
	default:
		switch fd.Kind() ***REMOVED***
		default:
		case protoreflect.BoolKind:
			return coderBoolValue
		case protoreflect.EnumKind:
			return coderEnumValue
		case protoreflect.Int32Kind:
			return coderInt32Value
		case protoreflect.Sint32Kind:
			return coderSint32Value
		case protoreflect.Uint32Kind:
			return coderUint32Value
		case protoreflect.Int64Kind:
			return coderInt64Value
		case protoreflect.Sint64Kind:
			return coderSint64Value
		case protoreflect.Uint64Kind:
			return coderUint64Value
		case protoreflect.Sfixed32Kind:
			return coderSfixed32Value
		case protoreflect.Fixed32Kind:
			return coderFixed32Value
		case protoreflect.FloatKind:
			return coderFloatValue
		case protoreflect.Sfixed64Kind:
			return coderSfixed64Value
		case protoreflect.Fixed64Kind:
			return coderFixed64Value
		case protoreflect.DoubleKind:
			return coderDoubleValue
		case protoreflect.StringKind:
			if strs.EnforceUTF8(fd) ***REMOVED***
				return coderStringValueValidateUTF8
			***REMOVED***
			return coderStringValue
		case protoreflect.BytesKind:
			return coderBytesValue
		case protoreflect.MessageKind:
			return coderMessageValue
		case protoreflect.GroupKind:
			return coderGroupValue
		***REMOVED***
	***REMOVED***
	panic(fmt.Sprintf("invalid field: no encoder for %v %v %v", fd.FullName(), fd.Cardinality(), fd.Kind()))
***REMOVED***
