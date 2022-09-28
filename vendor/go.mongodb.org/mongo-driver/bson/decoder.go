// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

// ErrDecodeToNil is the error returned when trying to decode to a nil value
var ErrDecodeToNil = errors.New("cannot Decode to nil value")

// This pool is used to keep the allocations of Decoders down. This is only used for the Marshal*
// methods and is not consumable from outside of this package. The Decoders retrieved from this pool
// must have both Reset and SetRegistry called on them.
var decPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return new(Decoder)
	***REMOVED***,
***REMOVED***

// A Decoder reads and decodes BSON documents from a stream. It reads from a bsonrw.ValueReader as
// the source of BSON data.
type Decoder struct ***REMOVED***
	dc bsoncodec.DecodeContext
	vr bsonrw.ValueReader

	// We persist defaultDocumentM and defaultDocumentD on the Decoder to prevent overwriting from
	// (*Decoder).SetContext.
	defaultDocumentM bool
	defaultDocumentD bool
***REMOVED***

// NewDecoder returns a new decoder that uses the DefaultRegistry to read from vr.
func NewDecoder(vr bsonrw.ValueReader) (*Decoder, error) ***REMOVED***
	if vr == nil ***REMOVED***
		return nil, errors.New("cannot create a new Decoder with a nil ValueReader")
	***REMOVED***

	return &Decoder***REMOVED***
		dc: bsoncodec.DecodeContext***REMOVED***Registry: DefaultRegistry***REMOVED***,
		vr: vr,
	***REMOVED***, nil
***REMOVED***

// NewDecoderWithContext returns a new decoder that uses DecodeContext dc to read from vr.
func NewDecoderWithContext(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader) (*Decoder, error) ***REMOVED***
	if dc.Registry == nil ***REMOVED***
		dc.Registry = DefaultRegistry
	***REMOVED***
	if vr == nil ***REMOVED***
		return nil, errors.New("cannot create a new Decoder with a nil ValueReader")
	***REMOVED***

	return &Decoder***REMOVED***
		dc: dc,
		vr: vr,
	***REMOVED***, nil
***REMOVED***

// Decode reads the next BSON document from the stream and decodes it into the
// value pointed to by val.
//
// The documentation for Unmarshal contains details about of BSON into a Go
// value.
func (d *Decoder) Decode(val interface***REMOVED******REMOVED***) error ***REMOVED***
	if unmarshaler, ok := val.(Unmarshaler); ok ***REMOVED***
		// TODO(skriptble): Reuse a []byte here and use the AppendDocumentBytes method.
		buf, err := bsonrw.Copier***REMOVED******REMOVED***.CopyDocumentToBytes(d.vr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return unmarshaler.UnmarshalBSON(buf)
	***REMOVED***

	rval := reflect.ValueOf(val)
	switch rval.Kind() ***REMOVED***
	case reflect.Ptr:
		if rval.IsNil() ***REMOVED***
			return ErrDecodeToNil
		***REMOVED***
		rval = rval.Elem()
	case reflect.Map:
		if rval.IsNil() ***REMOVED***
			return ErrDecodeToNil
		***REMOVED***
	default:
		return fmt.Errorf("argument to Decode must be a pointer or a map, but got %v", rval)
	***REMOVED***
	decoder, err := d.dc.LookupDecoder(rval.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if d.defaultDocumentM ***REMOVED***
		d.dc.DefaultDocumentM()
	***REMOVED***
	if d.defaultDocumentD ***REMOVED***
		d.dc.DefaultDocumentD()
	***REMOVED***
	return decoder.DecodeValue(d.dc, d.vr, rval)
***REMOVED***

// Reset will reset the state of the decoder, using the same *DecodeContext used in
// the original construction but using vr for reading.
func (d *Decoder) Reset(vr bsonrw.ValueReader) error ***REMOVED***
	d.vr = vr
	return nil
***REMOVED***

// SetRegistry replaces the current registry of the decoder with r.
func (d *Decoder) SetRegistry(r *bsoncodec.Registry) error ***REMOVED***
	d.dc.Registry = r
	return nil
***REMOVED***

// SetContext replaces the current registry of the decoder with dc.
func (d *Decoder) SetContext(dc bsoncodec.DecodeContext) error ***REMOVED***
	d.dc = dc
	return nil
***REMOVED***

// DefaultDocumentM will decode empty documents using the primitive.M type. This behavior is restricted to data typed as
// "interface***REMOVED******REMOVED***" or "map[string]interface***REMOVED******REMOVED***".
func (d *Decoder) DefaultDocumentM() ***REMOVED***
	d.defaultDocumentM = true
***REMOVED***

// DefaultDocumentD will decode empty documents using the primitive.D type. This behavior is restricted to data typed as
// "interface***REMOVED******REMOVED***" or "map[string]interface***REMOVED******REMOVED***".
func (d *Decoder) DefaultDocumentD() ***REMOVED***
	d.defaultDocumentD = true
***REMOVED***
