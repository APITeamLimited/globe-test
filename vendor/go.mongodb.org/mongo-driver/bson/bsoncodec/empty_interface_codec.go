// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmptyInterfaceCodec is the Codec used for interface***REMOVED******REMOVED*** values.
type EmptyInterfaceCodec struct ***REMOVED***
	DecodeBinaryAsSlice bool
***REMOVED***

var (
	defaultEmptyInterfaceCodec = NewEmptyInterfaceCodec()

	_ ValueCodec  = defaultEmptyInterfaceCodec
	_ typeDecoder = defaultEmptyInterfaceCodec
)

// NewEmptyInterfaceCodec returns a EmptyInterfaceCodec with options opts.
func NewEmptyInterfaceCodec(opts ...*bsonoptions.EmptyInterfaceCodecOptions) *EmptyInterfaceCodec ***REMOVED***
	interfaceOpt := bsonoptions.MergeEmptyInterfaceCodecOptions(opts...)

	codec := EmptyInterfaceCodec***REMOVED******REMOVED***
	if interfaceOpt.DecodeBinaryAsSlice != nil ***REMOVED***
		codec.DecodeBinaryAsSlice = *interfaceOpt.DecodeBinaryAsSlice
	***REMOVED***
	return &codec
***REMOVED***

// EncodeValue is the ValueEncoderFunc for interface***REMOVED******REMOVED***.
func (eic EmptyInterfaceCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tEmpty ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "EmptyInterfaceEncodeValue", Types: []reflect.Type***REMOVED***tEmpty***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***
	encoder, err := ec.LookupEncoder(val.Elem().Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return encoder.EncodeValue(ec, vw, val.Elem())
***REMOVED***

func (eic EmptyInterfaceCodec) getEmptyInterfaceDecodeType(dc DecodeContext, valueType bsontype.Type) (reflect.Type, error) ***REMOVED***
	isDocument := valueType == bsontype.Type(0) || valueType == bsontype.EmbeddedDocument
	if isDocument ***REMOVED***
		if dc.defaultDocumentType != nil ***REMOVED***
			// If the bsontype is an embedded document and the DocumentType is set on the DecodeContext, then return
			// that type.
			return dc.defaultDocumentType, nil
		***REMOVED***
		if dc.Ancestor != nil ***REMOVED***
			// Using ancestor information rather than looking up the type map entry forces consistent decoding.
			// If we're decoding into a bson.D, subdocuments should also be decoded as bson.D, even if a type map entry
			// has been registered.
			return dc.Ancestor, nil
		***REMOVED***
	***REMOVED***

	rtype, err := dc.LookupTypeMapEntry(valueType)
	if err == nil ***REMOVED***
		return rtype, nil
	***REMOVED***

	if isDocument ***REMOVED***
		// For documents, fallback to looking up a type map entry for bsontype.Type(0) or bsontype.EmbeddedDocument,
		// depending on the original valueType.
		var lookupType bsontype.Type
		switch valueType ***REMOVED***
		case bsontype.Type(0):
			lookupType = bsontype.EmbeddedDocument
		case bsontype.EmbeddedDocument:
			lookupType = bsontype.Type(0)
		***REMOVED***

		rtype, err = dc.LookupTypeMapEntry(lookupType)
		if err == nil ***REMOVED***
			return rtype, nil
		***REMOVED***
	***REMOVED***

	return nil, err
***REMOVED***

func (eic EmptyInterfaceCodec) decodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tEmpty ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***Name: "EmptyInterfaceDecodeValue", Types: []reflect.Type***REMOVED***tEmpty***REMOVED***, Received: reflect.Zero(t)***REMOVED***
	***REMOVED***

	rtype, err := eic.getEmptyInterfaceDecodeType(dc, vr.Type())
	if err != nil ***REMOVED***
		switch vr.Type() ***REMOVED***
		case bsontype.Null:
			return reflect.Zero(t), vr.ReadNull()
		default:
			return emptyValue, err
		***REMOVED***
	***REMOVED***

	decoder, err := dc.LookupDecoder(rtype)
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	elem, err := decodeTypeOrValue(decoder, dc, vr, rtype)
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	if eic.DecodeBinaryAsSlice && rtype == tBinary ***REMOVED***
		binElem := elem.Interface().(primitive.Binary)
		if binElem.Subtype == bsontype.BinaryGeneric || binElem.Subtype == bsontype.BinaryBinaryOld ***REMOVED***
			elem = reflect.ValueOf(binElem.Data)
		***REMOVED***
	***REMOVED***

	return elem, nil
***REMOVED***

// DecodeValue is the ValueDecoderFunc for interface***REMOVED******REMOVED***.
func (eic EmptyInterfaceCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tEmpty ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "EmptyInterfaceDecodeValue", Types: []reflect.Type***REMOVED***tEmpty***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := eic.decodeType(dc, vr, val.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***
