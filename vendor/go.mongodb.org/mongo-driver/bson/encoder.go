// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"errors"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

// This pool is used to keep the allocations of Encoders down. This is only used for the Marshal*
// methods and is not consumable from outside of this package. The Encoders retrieved from this pool
// must have both Reset and SetRegistry called on them.
var encPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return new(Encoder)
	***REMOVED***,
***REMOVED***

// An Encoder writes a serialization format to an output stream. It writes to a bsonrw.ValueWriter
// as the destination of BSON data.
type Encoder struct ***REMOVED***
	ec bsoncodec.EncodeContext
	vw bsonrw.ValueWriter
***REMOVED***

// NewEncoder returns a new encoder that uses the DefaultRegistry to write to vw.
func NewEncoder(vw bsonrw.ValueWriter) (*Encoder, error) ***REMOVED***
	if vw == nil ***REMOVED***
		return nil, errors.New("cannot create a new Encoder with a nil ValueWriter")
	***REMOVED***

	return &Encoder***REMOVED***
		ec: bsoncodec.EncodeContext***REMOVED***Registry: DefaultRegistry***REMOVED***,
		vw: vw,
	***REMOVED***, nil
***REMOVED***

// NewEncoderWithContext returns a new encoder that uses EncodeContext ec to write to vw.
func NewEncoderWithContext(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter) (*Encoder, error) ***REMOVED***
	if ec.Registry == nil ***REMOVED***
		ec = bsoncodec.EncodeContext***REMOVED***Registry: DefaultRegistry***REMOVED***
	***REMOVED***
	if vw == nil ***REMOVED***
		return nil, errors.New("cannot create a new Encoder with a nil ValueWriter")
	***REMOVED***

	return &Encoder***REMOVED***
		ec: ec,
		vw: vw,
	***REMOVED***, nil
***REMOVED***

// Encode writes the BSON encoding of val to the stream.
//
// The documentation for Marshal contains details about the conversion of Go
// values to BSON.
func (e *Encoder) Encode(val interface***REMOVED******REMOVED***) error ***REMOVED***
	if marshaler, ok := val.(Marshaler); ok ***REMOVED***
		// TODO(skriptble): Should we have a MarshalAppender interface so that we can have []byte reuse?
		buf, err := marshaler.MarshalBSON()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return bsonrw.Copier***REMOVED******REMOVED***.CopyDocumentFromBytes(e.vw, buf)
	***REMOVED***

	encoder, err := e.ec.LookupEncoder(reflect.TypeOf(val))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return encoder.EncodeValue(e.ec, e.vw, reflect.ValueOf(val))
***REMOVED***

// Reset will reset the state of the encoder, using the same *EncodeContext used in
// the original construction but using vw.
func (e *Encoder) Reset(vw bsonrw.ValueWriter) error ***REMOVED***
	e.vw = vw
	return nil
***REMOVED***

// SetRegistry replaces the current registry of the encoder with r.
func (e *Encoder) SetRegistry(r *bsoncodec.Registry) error ***REMOVED***
	e.ec.Registry = r
	return nil
***REMOVED***

// SetContext replaces the current EncodeContext of the encoder with er.
func (e *Encoder) SetContext(ec bsoncodec.EncodeContext) error ***REMOVED***
	e.ec = ec
	return nil
***REMOVED***
