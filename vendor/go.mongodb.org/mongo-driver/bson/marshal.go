// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"bytes"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const defaultDstCap = 256

var bvwPool = bsonrw.NewBSONValueWriterPool()
var extjPool = bsonrw.NewExtJSONValueWriterPool()

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

// Marshal returns the BSON encoding of val as a BSON document. If val is not a type that can be transformed into a
// document, MarshalValue should be used instead.
//
// Marshal will use the default registry created by NewRegistry to recursively
// marshal val into a []byte. Marshal will inspect struct tags and alter the
// marshaling process accordingly.
func Marshal(val interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return MarshalWithRegistry(DefaultRegistry, val)
***REMOVED***

// MarshalAppend will encode val as a BSON document and append the bytes to dst. If dst is not large enough to hold the
// bytes, it will be grown. If val is not a type that can be transformed into a document, MarshalValueAppend should be
// used instead.
func MarshalAppend(dst []byte, val interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return MarshalAppendWithRegistry(DefaultRegistry, dst, val)
***REMOVED***

// MarshalWithRegistry returns the BSON encoding of val as a BSON document. If val is not a type that can be transformed
// into a document, MarshalValueWithRegistry should be used instead.
func MarshalWithRegistry(r *bsoncodec.Registry, val interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	dst := make([]byte, 0)
	return MarshalAppendWithRegistry(r, dst, val)
***REMOVED***

// MarshalWithContext returns the BSON encoding of val as a BSON document using EncodeContext ec. If val is not a type
// that can be transformed into a document, MarshalValueWithContext should be used instead.
func MarshalWithContext(ec bsoncodec.EncodeContext, val interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	dst := make([]byte, 0)
	return MarshalAppendWithContext(ec, dst, val)
***REMOVED***

// MarshalAppendWithRegistry will encode val as a BSON document using Registry r and append the bytes to dst. If dst is
// not large enough to hold the bytes, it will be grown. If val is not a type that can be transformed into a document,
// MarshalValueAppendWithRegistry should be used instead.
func MarshalAppendWithRegistry(r *bsoncodec.Registry, dst []byte, val interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	return MarshalAppendWithContext(bsoncodec.EncodeContext***REMOVED***Registry: r***REMOVED***, dst, val)
***REMOVED***

// MarshalAppendWithContext will encode val as a BSON document using Registry r and EncodeContext ec and append the
// bytes to dst. If dst is not large enough to hold the bytes, it will be grown. If val is not a type that can be
// transformed into a document, MarshalValueAppendWithContext should be used instead.
func MarshalAppendWithContext(ec bsoncodec.EncodeContext, dst []byte, val interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	sw := new(bsonrw.SliceWriter)
	*sw = dst
	vw := bvwPool.Get(sw)
	defer bvwPool.Put(vw)

	enc := encPool.Get().(*Encoder)
	defer encPool.Put(enc)

	err := enc.Reset(vw)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = enc.SetContext(ec)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = enc.Encode(val)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return *sw, nil
***REMOVED***

// MarshalValue returns the BSON encoding of val.
//
// MarshalValue will use bson.DefaultRegistry to transform val into a BSON value. If val is a struct, this function will
// inspect struct tags and alter the marshalling process accordingly.
func MarshalValue(val interface***REMOVED******REMOVED***) (bsontype.Type, []byte, error) ***REMOVED***
	return MarshalValueWithRegistry(DefaultRegistry, val)
***REMOVED***

// MarshalValueAppend will append the BSON encoding of val to dst. If dst is not large enough to hold the BSON encoding
// of val, dst will be grown.
func MarshalValueAppend(dst []byte, val interface***REMOVED******REMOVED***) (bsontype.Type, []byte, error) ***REMOVED***
	return MarshalValueAppendWithRegistry(DefaultRegistry, dst, val)
***REMOVED***

// MarshalValueWithRegistry returns the BSON encoding of val using Registry r.
func MarshalValueWithRegistry(r *bsoncodec.Registry, val interface***REMOVED******REMOVED***) (bsontype.Type, []byte, error) ***REMOVED***
	dst := make([]byte, 0)
	return MarshalValueAppendWithRegistry(r, dst, val)
***REMOVED***

// MarshalValueWithContext returns the BSON encoding of val using EncodeContext ec.
func MarshalValueWithContext(ec bsoncodec.EncodeContext, val interface***REMOVED******REMOVED***) (bsontype.Type, []byte, error) ***REMOVED***
	dst := make([]byte, 0)
	return MarshalValueAppendWithContext(ec, dst, val)
***REMOVED***

// MarshalValueAppendWithRegistry will append the BSON encoding of val to dst using Registry r. If dst is not large
// enough to hold the BSON encoding of val, dst will be grown.
func MarshalValueAppendWithRegistry(r *bsoncodec.Registry, dst []byte, val interface***REMOVED******REMOVED***) (bsontype.Type, []byte, error) ***REMOVED***
	return MarshalValueAppendWithContext(bsoncodec.EncodeContext***REMOVED***Registry: r***REMOVED***, dst, val)
***REMOVED***

// MarshalValueAppendWithContext will append the BSON encoding of val to dst using EncodeContext ec. If dst is not large
// enough to hold the BSON encoding of val, dst will be grown.
func MarshalValueAppendWithContext(ec bsoncodec.EncodeContext, dst []byte, val interface***REMOVED******REMOVED***) (bsontype.Type, []byte, error) ***REMOVED***
	// get a ValueWriter configured to write to dst
	sw := new(bsonrw.SliceWriter)
	*sw = dst
	vwFlusher := bvwPool.GetAtModeElement(sw)

	// get an Encoder and encode the value
	enc := encPool.Get().(*Encoder)
	defer encPool.Put(enc)
	if err := enc.Reset(vwFlusher); err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***
	if err := enc.SetContext(ec); err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***
	if err := enc.Encode(val); err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***

	// flush the bytes written because we cannot guarantee that a full document has been written
	// after the flush, *sw will be in the format
	// [value type, 0 (null byte to indicate end of empty element name), value bytes..]
	if err := vwFlusher.Flush(); err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***
	buffer := *sw
	return bsontype.Type(buffer[0]), buffer[2:], nil
***REMOVED***

// MarshalExtJSON returns the extended JSON encoding of val.
func MarshalExtJSON(val interface***REMOVED******REMOVED***, canonical, escapeHTML bool) ([]byte, error) ***REMOVED***
	return MarshalExtJSONWithRegistry(DefaultRegistry, val, canonical, escapeHTML)
***REMOVED***

// MarshalExtJSONAppend will append the extended JSON encoding of val to dst.
// If dst is not large enough to hold the extended JSON encoding of val, dst
// will be grown.
func MarshalExtJSONAppend(dst []byte, val interface***REMOVED******REMOVED***, canonical, escapeHTML bool) ([]byte, error) ***REMOVED***
	return MarshalExtJSONAppendWithRegistry(DefaultRegistry, dst, val, canonical, escapeHTML)
***REMOVED***

// MarshalExtJSONWithRegistry returns the extended JSON encoding of val using Registry r.
func MarshalExtJSONWithRegistry(r *bsoncodec.Registry, val interface***REMOVED******REMOVED***, canonical, escapeHTML bool) ([]byte, error) ***REMOVED***
	dst := make([]byte, 0, defaultDstCap)
	return MarshalExtJSONAppendWithContext(bsoncodec.EncodeContext***REMOVED***Registry: r***REMOVED***, dst, val, canonical, escapeHTML)
***REMOVED***

// MarshalExtJSONWithContext returns the extended JSON encoding of val using Registry r.
func MarshalExtJSONWithContext(ec bsoncodec.EncodeContext, val interface***REMOVED******REMOVED***, canonical, escapeHTML bool) ([]byte, error) ***REMOVED***
	dst := make([]byte, 0, defaultDstCap)
	return MarshalExtJSONAppendWithContext(ec, dst, val, canonical, escapeHTML)
***REMOVED***

// MarshalExtJSONAppendWithRegistry will append the extended JSON encoding of
// val to dst using Registry r. If dst is not large enough to hold the BSON
// encoding of val, dst will be grown.
func MarshalExtJSONAppendWithRegistry(r *bsoncodec.Registry, dst []byte, val interface***REMOVED******REMOVED***, canonical, escapeHTML bool) ([]byte, error) ***REMOVED***
	return MarshalExtJSONAppendWithContext(bsoncodec.EncodeContext***REMOVED***Registry: r***REMOVED***, dst, val, canonical, escapeHTML)
***REMOVED***

// MarshalExtJSONAppendWithContext will append the extended JSON encoding of
// val to dst using Registry r. If dst is not large enough to hold the BSON
// encoding of val, dst will be grown.
func MarshalExtJSONAppendWithContext(ec bsoncodec.EncodeContext, dst []byte, val interface***REMOVED******REMOVED***, canonical, escapeHTML bool) ([]byte, error) ***REMOVED***
	sw := new(bsonrw.SliceWriter)
	*sw = dst
	ejvw := extjPool.Get(sw, canonical, escapeHTML)
	defer extjPool.Put(ejvw)

	enc := encPool.Get().(*Encoder)
	defer encPool.Put(enc)

	err := enc.Reset(ejvw)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = enc.SetContext(ec)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err = enc.Encode(val)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return *sw, nil
***REMOVED***

// IndentExtJSON will prefix and indent the provided extended JSON src and append it to dst.
func IndentExtJSON(dst *bytes.Buffer, src []byte, prefix, indent string) error ***REMOVED***
	return json.Indent(dst, src, prefix, indent)
***REMOVED***

// MarshalExtJSONIndent returns the extended JSON encoding of val with each line with prefixed
// and indented.
func MarshalExtJSONIndent(val interface***REMOVED******REMOVED***, canonical, escapeHTML bool, prefix, indent string) ([]byte, error) ***REMOVED***
	marshaled, err := MarshalExtJSON(val, canonical, escapeHTML)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var buf bytes.Buffer
	err = IndentExtJSON(&buf, marshaled, prefix, indent)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return buf.Bytes(), nil
***REMOVED***
