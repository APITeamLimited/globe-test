// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

var tRawValue = reflect.TypeOf(RawValue***REMOVED******REMOVED***)
var tRaw = reflect.TypeOf(Raw(nil))

var primitiveCodecs PrimitiveCodecs

// PrimitiveCodecs is a namespace for all of the default bsoncodec.Codecs for the primitive types
// defined in this package.
type PrimitiveCodecs struct***REMOVED******REMOVED***

// RegisterPrimitiveCodecs will register the encode and decode methods attached to PrimitiveCodecs
// with the provided RegistryBuilder. if rb is nil, a new empty RegistryBuilder will be created.
func (pc PrimitiveCodecs) RegisterPrimitiveCodecs(rb *bsoncodec.RegistryBuilder) ***REMOVED***
	if rb == nil ***REMOVED***
		panic(errors.New("argument to RegisterPrimitiveCodecs must not be nil"))
	***REMOVED***

	rb.
		RegisterTypeEncoder(tRawValue, bsoncodec.ValueEncoderFunc(pc.RawValueEncodeValue)).
		RegisterTypeEncoder(tRaw, bsoncodec.ValueEncoderFunc(pc.RawEncodeValue)).
		RegisterTypeDecoder(tRawValue, bsoncodec.ValueDecoderFunc(pc.RawValueDecodeValue)).
		RegisterTypeDecoder(tRaw, bsoncodec.ValueDecoderFunc(pc.RawDecodeValue))
***REMOVED***

// RawValueEncodeValue is the ValueEncoderFunc for RawValue.
func (PrimitiveCodecs) RawValueEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tRawValue ***REMOVED***
		return bsoncodec.ValueEncoderError***REMOVED***Name: "RawValueEncodeValue", Types: []reflect.Type***REMOVED***tRawValue***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	rawvalue := val.Interface().(RawValue)

	return bsonrw.Copier***REMOVED******REMOVED***.CopyValueFromBytes(vw, rawvalue.Type, rawvalue.Value)
***REMOVED***

// RawValueDecodeValue is the ValueDecoderFunc for RawValue.
func (PrimitiveCodecs) RawValueDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tRawValue ***REMOVED***
		return bsoncodec.ValueDecoderError***REMOVED***Name: "RawValueDecodeValue", Types: []reflect.Type***REMOVED***tRawValue***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	t, value, err := bsonrw.Copier***REMOVED******REMOVED***.CopyValueToBytes(vr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(reflect.ValueOf(RawValue***REMOVED***Type: t, Value: value***REMOVED***))
	return nil
***REMOVED***

// RawEncodeValue is the ValueEncoderFunc for Reader.
func (PrimitiveCodecs) RawEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tRaw ***REMOVED***
		return bsoncodec.ValueEncoderError***REMOVED***Name: "RawEncodeValue", Types: []reflect.Type***REMOVED***tRaw***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	rdr := val.Interface().(Raw)

	return bsonrw.Copier***REMOVED******REMOVED***.CopyDocumentFromBytes(vw, rdr)
***REMOVED***

// RawDecodeValue is the ValueDecoderFunc for Reader.
func (PrimitiveCodecs) RawDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tRaw ***REMOVED***
		return bsoncodec.ValueDecoderError***REMOVED***Name: "RawDecodeValue", Types: []reflect.Type***REMOVED***tRaw***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		val.Set(reflect.MakeSlice(val.Type(), 0, 0))
	***REMOVED***

	val.SetLen(0)

	rdr, err := bsonrw.Copier***REMOVED******REMOVED***.AppendDocumentBytes(val.Interface().(Raw), vr)
	val.Set(reflect.ValueOf(rdr))
	return err
***REMOVED***
