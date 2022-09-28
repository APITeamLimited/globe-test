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

// StringCodec is the Codec used for struct values.
type StringCodec struct ***REMOVED***
	DecodeObjectIDAsHex bool
***REMOVED***

var (
	defaultStringCodec = NewStringCodec()

	_ ValueCodec  = defaultStringCodec
	_ typeDecoder = defaultStringCodec
)

// NewStringCodec returns a StringCodec with options opts.
func NewStringCodec(opts ...*bsonoptions.StringCodecOptions) *StringCodec ***REMOVED***
	stringOpt := bsonoptions.MergeStringCodecOptions(opts...)
	return &StringCodec***REMOVED****stringOpt.DecodeObjectIDAsHex***REMOVED***
***REMOVED***

// EncodeValue is the ValueEncoder for string types.
func (sc *StringCodec) EncodeValue(ectx EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if val.Kind() != reflect.String ***REMOVED***
		return ValueEncoderError***REMOVED***
			Name:     "StringEncodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.String***REMOVED***,
			Received: val,
		***REMOVED***
	***REMOVED***

	return vw.WriteString(val.String())
***REMOVED***

func (sc *StringCodec) decodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t.Kind() != reflect.String ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "StringDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.String***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var str string
	var err error
	switch vr.Type() ***REMOVED***
	case bsontype.String:
		str, err = vr.ReadString()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.ObjectID:
		oid, err := vr.ReadObjectID()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if sc.DecodeObjectIDAsHex ***REMOVED***
			str = oid.Hex()
		***REMOVED*** else ***REMOVED***
			byteArray := [12]byte(oid)
			str = string(byteArray[:])
		***REMOVED***
	case bsontype.Symbol:
		str, err = vr.ReadSymbol()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Binary:
		data, subtype, err := vr.ReadBinary()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if subtype != bsontype.BinaryGeneric && subtype != bsontype.BinaryBinaryOld ***REMOVED***
			return emptyValue, decodeBinaryError***REMOVED***subtype: subtype, typeName: "string"***REMOVED***
		***REMOVED***
		str = string(data)
	case bsontype.Null:
		if err = vr.ReadNull(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Undefined:
		if err = vr.ReadUndefined(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a string type", vr.Type())
	***REMOVED***

	return reflect.ValueOf(str), nil
***REMOVED***

// DecodeValue is the ValueDecoder for string types.
func (sc *StringCodec) DecodeValue(dctx DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Kind() != reflect.String ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "StringDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.String***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := sc.decodeType(dctx, vr, val.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.SetString(elem.String())
	return nil
***REMOVED***
