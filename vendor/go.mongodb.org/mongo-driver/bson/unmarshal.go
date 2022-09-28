// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"bytes"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

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

// Unmarshal parses the BSON-encoded data and stores the result in the value
// pointed to by val. If val is nil or not a pointer, Unmarshal returns
// InvalidUnmarshalError.
func Unmarshal(data []byte, val interface***REMOVED******REMOVED***) error ***REMOVED***
	return UnmarshalWithRegistry(DefaultRegistry, data, val)
***REMOVED***

// UnmarshalWithRegistry parses the BSON-encoded data using Registry r and
// stores the result in the value pointed to by val. If val is nil or not
// a pointer, UnmarshalWithRegistry returns InvalidUnmarshalError.
func UnmarshalWithRegistry(r *bsoncodec.Registry, data []byte, val interface***REMOVED******REMOVED***) error ***REMOVED***
	vr := bsonrw.NewBSONDocumentReader(data)
	return unmarshalFromReader(bsoncodec.DecodeContext***REMOVED***Registry: r***REMOVED***, vr, val)
***REMOVED***

// UnmarshalWithContext parses the BSON-encoded data using DecodeContext dc and
// stores the result in the value pointed to by val. If val is nil or not
// a pointer, UnmarshalWithRegistry returns InvalidUnmarshalError.
func UnmarshalWithContext(dc bsoncodec.DecodeContext, data []byte, val interface***REMOVED******REMOVED***) error ***REMOVED***
	vr := bsonrw.NewBSONDocumentReader(data)
	return unmarshalFromReader(dc, vr, val)
***REMOVED***

// UnmarshalExtJSON parses the extended JSON-encoded data and stores the result
// in the value pointed to by val. If val is nil or not a pointer, Unmarshal
// returns InvalidUnmarshalError.
func UnmarshalExtJSON(data []byte, canonical bool, val interface***REMOVED******REMOVED***) error ***REMOVED***
	return UnmarshalExtJSONWithRegistry(DefaultRegistry, data, canonical, val)
***REMOVED***

// UnmarshalExtJSONWithRegistry parses the extended JSON-encoded data using
// Registry r and stores the result in the value pointed to by val. If val is
// nil or not a pointer, UnmarshalWithRegistry returns InvalidUnmarshalError.
func UnmarshalExtJSONWithRegistry(r *bsoncodec.Registry, data []byte, canonical bool, val interface***REMOVED******REMOVED***) error ***REMOVED***
	ejvr, err := bsonrw.NewExtJSONValueReader(bytes.NewReader(data), canonical)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return unmarshalFromReader(bsoncodec.DecodeContext***REMOVED***Registry: r***REMOVED***, ejvr, val)
***REMOVED***

// UnmarshalExtJSONWithContext parses the extended JSON-encoded data using
// DecodeContext dc and stores the result in the value pointed to by val. If val is
// nil or not a pointer, UnmarshalWithRegistry returns InvalidUnmarshalError.
func UnmarshalExtJSONWithContext(dc bsoncodec.DecodeContext, data []byte, canonical bool, val interface***REMOVED******REMOVED***) error ***REMOVED***
	ejvr, err := bsonrw.NewExtJSONValueReader(bytes.NewReader(data), canonical)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return unmarshalFromReader(dc, ejvr, val)
***REMOVED***

func unmarshalFromReader(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val interface***REMOVED******REMOVED***) error ***REMOVED***
	dec := decPool.Get().(*Decoder)
	defer decPool.Put(dec)

	err := dec.Reset(vr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = dec.SetContext(dc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return dec.Decode(val)
***REMOVED***
