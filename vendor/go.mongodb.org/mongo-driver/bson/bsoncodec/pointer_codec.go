// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var _ ValueEncoder = &PointerCodec***REMOVED******REMOVED***
var _ ValueDecoder = &PointerCodec***REMOVED******REMOVED***

// PointerCodec is the Codec used for pointers.
type PointerCodec struct ***REMOVED***
	ecache map[reflect.Type]ValueEncoder
	dcache map[reflect.Type]ValueDecoder
	l      sync.RWMutex
***REMOVED***

// NewPointerCodec returns a PointerCodec that has been initialized.
func NewPointerCodec() *PointerCodec ***REMOVED***
	return &PointerCodec***REMOVED***
		ecache: make(map[reflect.Type]ValueEncoder),
		dcache: make(map[reflect.Type]ValueDecoder),
	***REMOVED***
***REMOVED***

// EncodeValue handles encoding a pointer by either encoding it to BSON Null if the pointer is nil
// or looking up an encoder for the type of value the pointer points to.
func (pc *PointerCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if val.Kind() != reflect.Ptr ***REMOVED***
		if !val.IsValid() ***REMOVED***
			return vw.WriteNull()
		***REMOVED***
		return ValueEncoderError***REMOVED***Name: "PointerCodec.EncodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Ptr***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***

	pc.l.RLock()
	enc, ok := pc.ecache[val.Type()]
	pc.l.RUnlock()
	if ok ***REMOVED***
		if enc == nil ***REMOVED***
			return ErrNoEncoder***REMOVED***Type: val.Type()***REMOVED***
		***REMOVED***
		return enc.EncodeValue(ec, vw, val.Elem())
	***REMOVED***

	enc, err := ec.LookupEncoder(val.Type().Elem())
	pc.l.Lock()
	pc.ecache[val.Type()] = enc
	pc.l.Unlock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return enc.EncodeValue(ec, vw, val.Elem())
***REMOVED***

// DecodeValue handles decoding a pointer by looking up a decoder for the type it points to and
// using that to decode. If the BSON value is Null, this method will set the pointer to nil.
func (pc *PointerCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Kind() != reflect.Ptr ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "PointerCodec.DecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Ptr***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if vr.Type() == bsontype.Null ***REMOVED***
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	***REMOVED***
	if vr.Type() == bsontype.Undefined ***REMOVED***
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadUndefined()
	***REMOVED***

	if val.IsNil() ***REMOVED***
		val.Set(reflect.New(val.Type().Elem()))
	***REMOVED***

	pc.l.RLock()
	dec, ok := pc.dcache[val.Type()]
	pc.l.RUnlock()
	if ok ***REMOVED***
		if dec == nil ***REMOVED***
			return ErrNoDecoder***REMOVED***Type: val.Type()***REMOVED***
		***REMOVED***
		return dec.DecodeValue(dc, vr, val.Elem())
	***REMOVED***

	dec, err := dc.LookupDecoder(val.Type().Elem())
	pc.l.Lock()
	pc.dcache[val.Type()] = dec
	pc.l.Unlock()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return dec.DecodeValue(dc, vr, val.Elem())
***REMOVED***
