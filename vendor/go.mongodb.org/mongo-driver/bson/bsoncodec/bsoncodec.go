// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec // import "go.mongodb.org/mongo-driver/bson/bsoncodec"

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	emptyValue = reflect.Value***REMOVED******REMOVED***
)

// Marshaler is an interface implemented by types that can marshal themselves
// into a BSON document represented as bytes. The bytes returned must be a valid
// BSON document if the error is nil.
type Marshaler interface ***REMOVED***
	MarshalBSON() ([]byte, error)
***REMOVED***

// ValueMarshaler is an interface implemented by types that can marshal
// themselves into a BSON value as bytes. The type must be the valid type for
// the bytes returned. The bytes and byte type together must be valid if the
// error is nil.
type ValueMarshaler interface ***REMOVED***
	MarshalBSONValue() (bsontype.Type, []byte, error)
***REMOVED***

// Unmarshaler is an interface implemented by types that can unmarshal a BSON
// document representation of themselves. The BSON bytes can be assumed to be
// valid. UnmarshalBSON must copy the BSON bytes if it wishes to retain the data
// after returning.
type Unmarshaler interface ***REMOVED***
	UnmarshalBSON([]byte) error
***REMOVED***

// ValueUnmarshaler is an interface implemented by types that can unmarshal a
// BSON value representation of themselves. The BSON bytes and type can be
// assumed to be valid. UnmarshalBSONValue must copy the BSON value bytes if it
// wishes to retain the data after returning.
type ValueUnmarshaler interface ***REMOVED***
	UnmarshalBSONValue(bsontype.Type, []byte) error
***REMOVED***

// ValueEncoderError is an error returned from a ValueEncoder when the provided value can't be
// encoded by the ValueEncoder.
type ValueEncoderError struct ***REMOVED***
	Name     string
	Types    []reflect.Type
	Kinds    []reflect.Kind
	Received reflect.Value
***REMOVED***

func (vee ValueEncoderError) Error() string ***REMOVED***
	typeKinds := make([]string, 0, len(vee.Types)+len(vee.Kinds))
	for _, t := range vee.Types ***REMOVED***
		typeKinds = append(typeKinds, t.String())
	***REMOVED***
	for _, k := range vee.Kinds ***REMOVED***
		if k == reflect.Map ***REMOVED***
			typeKinds = append(typeKinds, "map[string]*")
			continue
		***REMOVED***
		typeKinds = append(typeKinds, k.String())
	***REMOVED***
	received := vee.Received.Kind().String()
	if vee.Received.IsValid() ***REMOVED***
		received = vee.Received.Type().String()
	***REMOVED***
	return fmt.Sprintf("%s can only encode valid %s, but got %s", vee.Name, strings.Join(typeKinds, ", "), received)
***REMOVED***

// ValueDecoderError is an error returned from a ValueDecoder when the provided value can't be
// decoded by the ValueDecoder.
type ValueDecoderError struct ***REMOVED***
	Name     string
	Types    []reflect.Type
	Kinds    []reflect.Kind
	Received reflect.Value
***REMOVED***

func (vde ValueDecoderError) Error() string ***REMOVED***
	typeKinds := make([]string, 0, len(vde.Types)+len(vde.Kinds))
	for _, t := range vde.Types ***REMOVED***
		typeKinds = append(typeKinds, t.String())
	***REMOVED***
	for _, k := range vde.Kinds ***REMOVED***
		if k == reflect.Map ***REMOVED***
			typeKinds = append(typeKinds, "map[string]*")
			continue
		***REMOVED***
		typeKinds = append(typeKinds, k.String())
	***REMOVED***
	received := vde.Received.Kind().String()
	if vde.Received.IsValid() ***REMOVED***
		received = vde.Received.Type().String()
	***REMOVED***
	return fmt.Sprintf("%s can only decode valid and settable %s, but got %s", vde.Name, strings.Join(typeKinds, ", "), received)
***REMOVED***

// EncodeContext is the contextual information required for a Codec to encode a
// value.
type EncodeContext struct ***REMOVED***
	*Registry
	MinSize bool
***REMOVED***

// DecodeContext is the contextual information required for a Codec to decode a
// value.
type DecodeContext struct ***REMOVED***
	*Registry
	Truncate bool

	// Ancestor is the type of a containing document. This is mainly used to determine what type
	// should be used when decoding an embedded document into an empty interface. For example, if
	// Ancestor is a bson.M, BSON embedded document values being decoded into an empty interface
	// will be decoded into a bson.M.
	//
	// Deprecated: Use DefaultDocumentM or DefaultDocumentD instead.
	Ancestor reflect.Type

	// defaultDocumentType specifies the Go type to decode top-level and nested BSON documents into. In particular, the
	// usage for this field is restricted to data typed as "interface***REMOVED******REMOVED***" or "map[string]interface***REMOVED******REMOVED***". If DocumentType is
	// set to a type that a BSON document cannot be unmarshaled into (e.g. "string"), unmarshalling will result in an
	// error. DocumentType overrides the Ancestor field.
	defaultDocumentType reflect.Type
***REMOVED***

// DefaultDocumentM will decode empty documents using the primitive.M type. This behavior is restricted to data typed as
// "interface***REMOVED******REMOVED***" or "map[string]interface***REMOVED******REMOVED***".
func (dc *DecodeContext) DefaultDocumentM() ***REMOVED***
	dc.defaultDocumentType = reflect.TypeOf(primitive.M***REMOVED******REMOVED***)
***REMOVED***

// DefaultDocumentD will decode empty documents using the primitive.D type. This behavior is restricted to data typed as
// "interface***REMOVED******REMOVED***" or "map[string]interface***REMOVED******REMOVED***".
func (dc *DecodeContext) DefaultDocumentD() ***REMOVED***
	dc.defaultDocumentType = reflect.TypeOf(primitive.D***REMOVED******REMOVED***)
***REMOVED***

// ValueCodec is the interface that groups the methods to encode and decode
// values.
type ValueCodec interface ***REMOVED***
	ValueEncoder
	ValueDecoder
***REMOVED***

// ValueEncoder is the interface implemented by types that can handle the encoding of a value.
type ValueEncoder interface ***REMOVED***
	EncodeValue(EncodeContext, bsonrw.ValueWriter, reflect.Value) error
***REMOVED***

// ValueEncoderFunc is an adapter function that allows a function with the correct signature to be
// used as a ValueEncoder.
type ValueEncoderFunc func(EncodeContext, bsonrw.ValueWriter, reflect.Value) error

// EncodeValue implements the ValueEncoder interface.
func (fn ValueEncoderFunc) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	return fn(ec, vw, val)
***REMOVED***

// ValueDecoder is the interface implemented by types that can handle the decoding of a value.
type ValueDecoder interface ***REMOVED***
	DecodeValue(DecodeContext, bsonrw.ValueReader, reflect.Value) error
***REMOVED***

// ValueDecoderFunc is an adapter function that allows a function with the correct signature to be
// used as a ValueDecoder.
type ValueDecoderFunc func(DecodeContext, bsonrw.ValueReader, reflect.Value) error

// DecodeValue implements the ValueDecoder interface.
func (fn ValueDecoderFunc) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	return fn(dc, vr, val)
***REMOVED***

// typeDecoder is the interface implemented by types that can handle the decoding of a value given its type.
type typeDecoder interface ***REMOVED***
	decodeType(DecodeContext, bsonrw.ValueReader, reflect.Type) (reflect.Value, error)
***REMOVED***

// typeDecoderFunc is an adapter function that allows a function with the correct signature to be used as a typeDecoder.
type typeDecoderFunc func(DecodeContext, bsonrw.ValueReader, reflect.Type) (reflect.Value, error)

func (fn typeDecoderFunc) decodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	return fn(dc, vr, t)
***REMOVED***

// decodeAdapter allows two functions with the correct signatures to be used as both a ValueDecoder and typeDecoder.
type decodeAdapter struct ***REMOVED***
	ValueDecoderFunc
	typeDecoderFunc
***REMOVED***

var _ ValueDecoder = decodeAdapter***REMOVED******REMOVED***
var _ typeDecoder = decodeAdapter***REMOVED******REMOVED***

// decodeTypeOrValue calls decoder.decodeType is decoder is a typeDecoder. Otherwise, it allocates a new element of type
// t and calls decoder.DecodeValue on it.
func decodeTypeOrValue(decoder ValueDecoder, dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	td, _ := decoder.(typeDecoder)
	return decodeTypeOrValueWithInfo(decoder, td, dc, vr, t, true)
***REMOVED***

func decodeTypeOrValueWithInfo(vd ValueDecoder, td typeDecoder, dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type, convert bool) (reflect.Value, error) ***REMOVED***
	if td != nil ***REMOVED***
		val, err := td.decodeType(dc, vr, t)
		if err == nil && convert && val.Type() != t ***REMOVED***
			// This conversion step is necessary for slices and maps. If a user declares variables like:
			//
			// type myBool bool
			// var m map[string]myBool
			//
			// and tries to decode BSON bytes into the map, the decoding will fail if this conversion is not present
			// because we'll try to assign a value of type bool to one of type myBool.
			val = val.Convert(t)
		***REMOVED***
		return val, err
	***REMOVED***

	val := reflect.New(t).Elem()
	err := vd.DecodeValue(dc, vr, val)
	return val, err
***REMOVED***

// CodecZeroer is the interface implemented by Codecs that can also determine if
// a value of the type that would be encoded is zero.
type CodecZeroer interface ***REMOVED***
	IsTypeZero(interface***REMOVED******REMOVED***) bool
***REMOVED***
