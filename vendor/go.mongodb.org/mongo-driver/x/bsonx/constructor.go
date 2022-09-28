// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"encoding/binary"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDoc is the interface implemented by Doc and MDoc. It allows either of these types to be provided
// to the Document function to create a Value.
type IDoc interface ***REMOVED***
	idoc()
***REMOVED***

// Double constructs a BSON double Value.
func Double(f64 float64) Val ***REMOVED***
	v := Val***REMOVED***t: bsontype.Double***REMOVED***
	binary.LittleEndian.PutUint64(v.bootstrap[0:8], math.Float64bits(f64))
	return v
***REMOVED***

// String constructs a BSON string Value.
func String(str string) Val ***REMOVED*** return Val***REMOVED***t: bsontype.String***REMOVED***.writestring(str) ***REMOVED***

// Document constructs a Value from the given IDoc. If nil is provided, a BSON Null value will be
// returned.
func Document(doc IDoc) Val ***REMOVED***
	var v Val
	switch tt := doc.(type) ***REMOVED***
	case Doc:
		if tt == nil ***REMOVED***
			v.t = bsontype.Null
			break
		***REMOVED***
		v.t = bsontype.EmbeddedDocument
		v.primitive = tt
	case MDoc:
		if tt == nil ***REMOVED***
			v.t = bsontype.Null
			break
		***REMOVED***
		v.t = bsontype.EmbeddedDocument
		v.primitive = tt
	default:
		v.t = bsontype.Null
	***REMOVED***
	return v
***REMOVED***

// Array constructs a Value from arr. If arr is nil, a BSON Null value is returned.
func Array(arr Arr) Val ***REMOVED***
	if arr == nil ***REMOVED***
		return Val***REMOVED***t: bsontype.Null***REMOVED***
	***REMOVED***
	return Val***REMOVED***t: bsontype.Array, primitive: arr***REMOVED***
***REMOVED***

// Binary constructs a BSON binary Value.
func Binary(subtype byte, data []byte) Val ***REMOVED***
	return Val***REMOVED***t: bsontype.Binary, primitive: primitive.Binary***REMOVED***Subtype: subtype, Data: data***REMOVED******REMOVED***
***REMOVED***

// Undefined constructs a BSON binary Value.
func Undefined() Val ***REMOVED*** return Val***REMOVED***t: bsontype.Undefined***REMOVED*** ***REMOVED***

// ObjectID constructs a BSON objectid Value.
func ObjectID(oid primitive.ObjectID) Val ***REMOVED***
	v := Val***REMOVED***t: bsontype.ObjectID***REMOVED***
	copy(v.bootstrap[0:12], oid[:])
	return v
***REMOVED***

// Boolean constructs a BSON boolean Value.
func Boolean(b bool) Val ***REMOVED***
	v := Val***REMOVED***t: bsontype.Boolean***REMOVED***
	if b ***REMOVED***
		v.bootstrap[0] = 0x01
	***REMOVED***
	return v
***REMOVED***

// DateTime constructs a BSON datetime Value.
func DateTime(dt int64) Val ***REMOVED*** return Val***REMOVED***t: bsontype.DateTime***REMOVED***.writei64(dt) ***REMOVED***

// Time constructs a BSON datetime Value.
func Time(t time.Time) Val ***REMOVED***
	return Val***REMOVED***t: bsontype.DateTime***REMOVED***.writei64(t.Unix()*1e3 + int64(t.Nanosecond()/1e6))
***REMOVED***

// Null constructs a BSON binary Value.
func Null() Val ***REMOVED*** return Val***REMOVED***t: bsontype.Null***REMOVED*** ***REMOVED***

// Regex constructs a BSON regex Value.
func Regex(pattern, options string) Val ***REMOVED***
	regex := primitive.Regex***REMOVED***Pattern: pattern, Options: options***REMOVED***
	return Val***REMOVED***t: bsontype.Regex, primitive: regex***REMOVED***
***REMOVED***

// DBPointer constructs a BSON dbpointer Value.
func DBPointer(ns string, ptr primitive.ObjectID) Val ***REMOVED***
	dbptr := primitive.DBPointer***REMOVED***DB: ns, Pointer: ptr***REMOVED***
	return Val***REMOVED***t: bsontype.DBPointer, primitive: dbptr***REMOVED***
***REMOVED***

// JavaScript constructs a BSON javascript Value.
func JavaScript(js string) Val ***REMOVED***
	return Val***REMOVED***t: bsontype.JavaScript***REMOVED***.writestring(js)
***REMOVED***

// Symbol constructs a BSON symbol Value.
func Symbol(symbol string) Val ***REMOVED***
	return Val***REMOVED***t: bsontype.Symbol***REMOVED***.writestring(symbol)
***REMOVED***

// CodeWithScope constructs a BSON code with scope Value.
func CodeWithScope(code string, scope IDoc) Val ***REMOVED***
	cws := primitive.CodeWithScope***REMOVED***Code: primitive.JavaScript(code), Scope: scope***REMOVED***
	return Val***REMOVED***t: bsontype.CodeWithScope, primitive: cws***REMOVED***
***REMOVED***

// Int32 constructs a BSON int32 Value.
func Int32(i32 int32) Val ***REMOVED***
	v := Val***REMOVED***t: bsontype.Int32***REMOVED***
	v.bootstrap[0] = byte(i32)
	v.bootstrap[1] = byte(i32 >> 8)
	v.bootstrap[2] = byte(i32 >> 16)
	v.bootstrap[3] = byte(i32 >> 24)
	return v
***REMOVED***

// Timestamp constructs a BSON timestamp Value.
func Timestamp(t, i uint32) Val ***REMOVED***
	v := Val***REMOVED***t: bsontype.Timestamp***REMOVED***
	v.bootstrap[0] = byte(i)
	v.bootstrap[1] = byte(i >> 8)
	v.bootstrap[2] = byte(i >> 16)
	v.bootstrap[3] = byte(i >> 24)
	v.bootstrap[4] = byte(t)
	v.bootstrap[5] = byte(t >> 8)
	v.bootstrap[6] = byte(t >> 16)
	v.bootstrap[7] = byte(t >> 24)
	return v
***REMOVED***

// Int64 constructs a BSON int64 Value.
func Int64(i64 int64) Val ***REMOVED*** return Val***REMOVED***t: bsontype.Int64***REMOVED***.writei64(i64) ***REMOVED***

// Decimal128 constructs a BSON decimal128 Value.
func Decimal128(d128 primitive.Decimal128) Val ***REMOVED***
	return Val***REMOVED***t: bsontype.Decimal128, primitive: d128***REMOVED***
***REMOVED***

// MinKey constructs a BSON minkey Value.
func MinKey() Val ***REMOVED*** return Val***REMOVED***t: bsontype.MinKey***REMOVED*** ***REMOVED***

// MaxKey constructs a BSON maxkey Value.
func MaxKey() Val ***REMOVED*** return Val***REMOVED***t: bsontype.MaxKey***REMOVED*** ***REMOVED***
