// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

// condAddrEncoder is the encoder used when a pointer to the encoding value has an encoder.
type condAddrEncoder struct ***REMOVED***
	canAddrEnc ValueEncoder
	elseEnc    ValueEncoder
***REMOVED***

var _ ValueEncoder = (*condAddrEncoder)(nil)

// newCondAddrEncoder returns an condAddrEncoder.
func newCondAddrEncoder(canAddrEnc, elseEnc ValueEncoder) *condAddrEncoder ***REMOVED***
	encoder := condAddrEncoder***REMOVED***canAddrEnc: canAddrEnc, elseEnc: elseEnc***REMOVED***
	return &encoder
***REMOVED***

// EncodeValue is the ValueEncoderFunc for a value that may be addressable.
func (cae *condAddrEncoder) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if val.CanAddr() ***REMOVED***
		return cae.canAddrEnc.EncodeValue(ec, vw, val)
	***REMOVED***
	if cae.elseEnc != nil ***REMOVED***
		return cae.elseEnc.EncodeValue(ec, vw, val)
	***REMOVED***
	return ErrNoEncoder***REMOVED***Type: val.Type()***REMOVED***
***REMOVED***

// condAddrDecoder is the decoder used when a pointer to the value has a decoder.
type condAddrDecoder struct ***REMOVED***
	canAddrDec ValueDecoder
	elseDec    ValueDecoder
***REMOVED***

var _ ValueDecoder = (*condAddrDecoder)(nil)

// newCondAddrDecoder returns an CondAddrDecoder.
func newCondAddrDecoder(canAddrDec, elseDec ValueDecoder) *condAddrDecoder ***REMOVED***
	decoder := condAddrDecoder***REMOVED***canAddrDec: canAddrDec, elseDec: elseDec***REMOVED***
	return &decoder
***REMOVED***

// DecodeValue is the ValueDecoderFunc for a value that may be addressable.
func (cad *condAddrDecoder) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if val.CanAddr() ***REMOVED***
		return cad.canAddrDec.DecodeValue(dc, vr, val)
	***REMOVED***
	if cad.elseDec != nil ***REMOVED***
		return cad.elseDec.DecodeValue(dc, vr, val)
	***REMOVED***
	return ErrNoDecoder***REMOVED***Type: val.Type()***REMOVED***
***REMOVED***
