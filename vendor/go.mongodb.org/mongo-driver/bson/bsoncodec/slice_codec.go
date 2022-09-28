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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var defaultSliceCodec = NewSliceCodec()

// SliceCodec is the Codec used for slice values.
type SliceCodec struct ***REMOVED***
	EncodeNilAsEmpty bool
***REMOVED***

var _ ValueCodec = &MapCodec***REMOVED******REMOVED***

// NewSliceCodec returns a MapCodec with options opts.
func NewSliceCodec(opts ...*bsonoptions.SliceCodecOptions) *SliceCodec ***REMOVED***
	sliceOpt := bsonoptions.MergeSliceCodecOptions(opts...)

	codec := SliceCodec***REMOVED******REMOVED***
	if sliceOpt.EncodeNilAsEmpty != nil ***REMOVED***
		codec.EncodeNilAsEmpty = *sliceOpt.EncodeNilAsEmpty
	***REMOVED***
	return &codec
***REMOVED***

// EncodeValue is the ValueEncoder for slice types.
func (sc SliceCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Kind() != reflect.Slice ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "SliceEncodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Slice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() && !sc.EncodeNilAsEmpty ***REMOVED***
		return vw.WriteNull()
	***REMOVED***

	// If we have a []byte we want to treat it as a binary instead of as an array.
	if val.Type().Elem() == tByte ***REMOVED***
		var byteSlice []byte
		for idx := 0; idx < val.Len(); idx++ ***REMOVED***
			byteSlice = append(byteSlice, val.Index(idx).Interface().(byte))
		***REMOVED***
		return vw.WriteBinary(byteSlice)
	***REMOVED***

	// If we have a []primitive.E we want to treat it as a document instead of as an array.
	if val.Type().ConvertibleTo(tD) ***REMOVED***
		d := val.Convert(tD).Interface().(primitive.D)

		dw, err := vw.WriteDocument()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		for _, e := range d ***REMOVED***
			err = encodeElement(ec, dw, e)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		return dw.WriteDocumentEnd()
	***REMOVED***

	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	elemType := val.Type().Elem()
	encoder, err := ec.LookupEncoder(elemType)
	if err != nil && elemType.Kind() != reflect.Interface ***REMOVED***
		return err
	***REMOVED***

	for idx := 0; idx < val.Len(); idx++ ***REMOVED***
		currEncoder, currVal, lookupErr := defaultValueEncoders.lookupElementEncoder(ec, encoder, val.Index(idx))
		if lookupErr != nil && lookupErr != errInvalidValue ***REMOVED***
			return lookupErr
		***REMOVED***

		vw, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if lookupErr == errInvalidValue ***REMOVED***
			err = vw.WriteNull()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		err = currEncoder.EncodeValue(ec, vw, currVal)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return aw.WriteArrayEnd()
***REMOVED***

// DecodeValue is the ValueDecoder for slice types.
func (sc *SliceCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Kind() != reflect.Slice ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "SliceDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Slice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Array:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	case bsontype.Undefined:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadUndefined()
	case bsontype.Type(0), bsontype.EmbeddedDocument:
		if val.Type().Elem() != tE ***REMOVED***
			return fmt.Errorf("cannot decode document into %s", val.Type())
		***REMOVED***
	case bsontype.Binary:
		if val.Type().Elem() != tByte ***REMOVED***
			return fmt.Errorf("SliceDecodeValue can only decode a binary into a byte array, got %v", vrType)
		***REMOVED***
		data, subtype, err := vr.ReadBinary()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if subtype != bsontype.BinaryGeneric && subtype != bsontype.BinaryBinaryOld ***REMOVED***
			return fmt.Errorf("SliceDecodeValue can only be used to decode subtype 0x00 or 0x02 for %s, got %v", bsontype.Binary, subtype)
		***REMOVED***

		if val.IsNil() ***REMOVED***
			val.Set(reflect.MakeSlice(val.Type(), 0, len(data)))
		***REMOVED***

		val.SetLen(0)
		for _, elem := range data ***REMOVED***
			val.Set(reflect.Append(val, reflect.ValueOf(elem)))
		***REMOVED***
		return nil
	case bsontype.String:
		if sliceType := val.Type().Elem(); sliceType != tByte ***REMOVED***
			return fmt.Errorf("SliceDecodeValue can only decode a string into a byte array, got %v", sliceType)
		***REMOVED***
		str, err := vr.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		byteStr := []byte(str)

		if val.IsNil() ***REMOVED***
			val.Set(reflect.MakeSlice(val.Type(), 0, len(byteStr)))
		***REMOVED***

		val.SetLen(0)
		for _, elem := range byteStr ***REMOVED***
			val.Set(reflect.Append(val, reflect.ValueOf(elem)))
		***REMOVED***
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a slice", vrType)
	***REMOVED***

	var elemsFunc func(DecodeContext, bsonrw.ValueReader, reflect.Value) ([]reflect.Value, error)
	switch val.Type().Elem() ***REMOVED***
	case tE:
		dc.Ancestor = val.Type()
		elemsFunc = defaultValueDecoders.decodeD
	default:
		elemsFunc = defaultValueDecoders.decodeDefault
	***REMOVED***

	elems, err := elemsFunc(dc, vr, val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if val.IsNil() ***REMOVED***
		val.Set(reflect.MakeSlice(val.Type(), 0, len(elems)))
	***REMOVED***

	val.SetLen(0)
	val.Set(reflect.Append(val, elems...))

	return nil
***REMOVED***
