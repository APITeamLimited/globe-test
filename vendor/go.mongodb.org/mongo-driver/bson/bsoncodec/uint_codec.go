// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"fmt"
	"math"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// UIntCodec is the Codec used for uint values.
type UIntCodec struct ***REMOVED***
	EncodeToMinSize bool
***REMOVED***

var (
	defaultUIntCodec = NewUIntCodec()

	_ ValueCodec  = defaultUIntCodec
	_ typeDecoder = defaultUIntCodec
)

// NewUIntCodec returns a UIntCodec with options opts.
func NewUIntCodec(opts ...*bsonoptions.UIntCodecOptions) *UIntCodec ***REMOVED***
	uintOpt := bsonoptions.MergeUIntCodecOptions(opts...)

	codec := UIntCodec***REMOVED******REMOVED***
	if uintOpt.EncodeToMinSize != nil ***REMOVED***
		codec.EncodeToMinSize = *uintOpt.EncodeToMinSize
	***REMOVED***
	return &codec
***REMOVED***

// EncodeValue is the ValueEncoder for uint types.
func (uic *UIntCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	switch val.Kind() ***REMOVED***
	case reflect.Uint8, reflect.Uint16:
		return vw.WriteInt32(int32(val.Uint()))
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		u64 := val.Uint()

		// If ec.MinSize or if encodeToMinSize is true for a non-uint64 value we should write val as an int32
		useMinSize := ec.MinSize || (uic.EncodeToMinSize && val.Kind() != reflect.Uint64)

		if u64 <= math.MaxInt32 && useMinSize ***REMOVED***
			return vw.WriteInt32(int32(u64))
		***REMOVED***
		if u64 > math.MaxInt64 ***REMOVED***
			return fmt.Errorf("%d overflows int64", u64)
		***REMOVED***
		return vw.WriteInt64(int64(u64))
	***REMOVED***

	return ValueEncoderError***REMOVED***
		Name:     "UintEncodeValue",
		Kinds:    []reflect.Kind***REMOVED***reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint***REMOVED***,
		Received: val,
	***REMOVED***
***REMOVED***

func (uic *UIntCodec) decodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	var i64 int64
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		i64 = int64(i32)
	case bsontype.Int64:
		i64, err = vr.ReadInt64()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Double:
		f64, err := vr.ReadDouble()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if !dc.Truncate && math.Floor(f64) != f64 ***REMOVED***
			return emptyValue, errCannotTruncate
		***REMOVED***
		if f64 > float64(math.MaxInt64) ***REMOVED***
			return emptyValue, fmt.Errorf("%g overflows int64", f64)
		***REMOVED***
		i64 = int64(f64)
	case bsontype.Boolean:
		b, err := vr.ReadBoolean()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if b ***REMOVED***
			i64 = 1
		***REMOVED***
	case bsontype.Null:
		if err = vr.ReadNull(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Undefined:
		if err = vr.ReadUndefined(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into an integer type", vrType)
	***REMOVED***

	switch t.Kind() ***REMOVED***
	case reflect.Uint8:
		if i64 < 0 || i64 > math.MaxUint8 ***REMOVED***
			return emptyValue, fmt.Errorf("%d overflows uint8", i64)
		***REMOVED***

		return reflect.ValueOf(uint8(i64)), nil
	case reflect.Uint16:
		if i64 < 0 || i64 > math.MaxUint16 ***REMOVED***
			return emptyValue, fmt.Errorf("%d overflows uint16", i64)
		***REMOVED***

		return reflect.ValueOf(uint16(i64)), nil
	case reflect.Uint32:
		if i64 < 0 || i64 > math.MaxUint32 ***REMOVED***
			return emptyValue, fmt.Errorf("%d overflows uint32", i64)
		***REMOVED***

		return reflect.ValueOf(uint32(i64)), nil
	case reflect.Uint64:
		if i64 < 0 ***REMOVED***
			return emptyValue, fmt.Errorf("%d overflows uint64", i64)
		***REMOVED***

		return reflect.ValueOf(uint64(i64)), nil
	case reflect.Uint:
		if i64 < 0 || int64(uint(i64)) != i64 ***REMOVED*** // Can we fit this inside of an uint
			return emptyValue, fmt.Errorf("%d overflows uint", i64)
		***REMOVED***

		return reflect.ValueOf(uint(i64)), nil
	default:
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "UintDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***
***REMOVED***

// DecodeValue is the ValueDecoder for uint types.
func (uic *UIntCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() ***REMOVED***
		return ValueDecoderError***REMOVED***
			Name:     "UintDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint***REMOVED***,
			Received: val,
		***REMOVED***
	***REMOVED***

	elem, err := uic.decodeType(dc, vr, val.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.SetUint(elem.Uint())
	return nil
***REMOVED***
