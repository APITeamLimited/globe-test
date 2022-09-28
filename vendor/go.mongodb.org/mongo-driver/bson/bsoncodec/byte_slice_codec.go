// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// ByteSliceCodec is the Codec used for []byte values.
type ByteSliceCodec struct ***REMOVED***
	EncodeNilAsEmpty bool
***REMOVED***

var (
	defaultByteSliceCodec = NewByteSliceCodec()

	_ ValueCodec  = defaultByteSliceCodec
	_ typeDecoder = defaultByteSliceCodec
)

// NewByteSliceCodec returns a StringCodec with options opts.
func NewByteSliceCodec(opts ...*bsonoptions.ByteSliceCodecOptions) *ByteSliceCodec ***REMOVED***
	byteSliceOpt := bsonoptions.MergeByteSliceCodecOptions(opts...)
	codec := ByteSliceCodec***REMOVED******REMOVED***
	if byteSliceOpt.EncodeNilAsEmpty != nil ***REMOVED***
		codec.EncodeNilAsEmpty = *byteSliceOpt.EncodeNilAsEmpty
	***REMOVED***
	return &codec
***REMOVED***

// EncodeValue is the ValueEncoder for []byte.
func (bsc *ByteSliceCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tByteSlice ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "ByteSliceEncodeValue", Types: []reflect.Type***REMOVED***tByteSlice***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	if val.IsNil() && !bsc.EncodeNilAsEmpty ***REMOVED***
		return vw.WriteNull()
	***REMOVED***
	return vw.WriteBinary(val.Interface().([]byte))
***REMOVED***

func (bsc *ByteSliceCodec) decodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tByteSlice ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "ByteSliceDecodeValue",
			Types:    []reflect.Type***REMOVED***tByteSlice***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var data []byte
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		data = []byte(str)
	case bsontype.Symbol:
		sym, err := vr.ReadSymbol()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		data = []byte(sym)
	case bsontype.Binary:
		var subtype byte
		data, subtype, err = vr.ReadBinary()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if subtype != bsontype.BinaryGeneric && subtype != bsontype.BinaryBinaryOld ***REMOVED***
			return emptyValue, decodeBinaryError***REMOVED***subtype: subtype, typeName: "[]byte"***REMOVED***
		***REMOVED***
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a []byte", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(data), nil
***REMOVED***

// DecodeValue is the ValueDecoder for []byte.
func (bsc *ByteSliceCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tByteSlice ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "ByteSliceDecodeValue", Types: []reflect.Type***REMOVED***tByteSlice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := bsc.decodeType(dc, vr, tByteSlice)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***
