// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrNilContext is returned when the provided DecodeContext is nil.
var ErrNilContext = errors.New("DecodeContext cannot be nil")

// ErrNilRegistry is returned when the provided registry is nil.
var ErrNilRegistry = errors.New("Registry cannot be nil")

// RawValue represents a BSON value in byte form. It can be used to hold unprocessed BSON or to
// defer processing of BSON. Type is the BSON type of the value and Value are the raw bytes that
// represent the element.
//
// This type wraps bsoncore.Value for most of it's functionality.
type RawValue struct ***REMOVED***
	Type  bsontype.Type
	Value []byte

	r *bsoncodec.Registry
***REMOVED***

// Unmarshal deserializes BSON into the provided val. If RawValue cannot be unmarshaled into val, an
// error is returned. This method will use the registry used to create the RawValue, if the RawValue
// was created from partial BSON processing, or it will use the default registry. Users wishing to
// specify the registry to use should use UnmarshalWithRegistry.
func (rv RawValue) Unmarshal(val interface***REMOVED******REMOVED***) error ***REMOVED***
	reg := rv.r
	if reg == nil ***REMOVED***
		reg = DefaultRegistry
	***REMOVED***
	return rv.UnmarshalWithRegistry(reg, val)
***REMOVED***

// Equal compares rv and rv2 and returns true if they are equal.
func (rv RawValue) Equal(rv2 RawValue) bool ***REMOVED***
	if rv.Type != rv2.Type ***REMOVED***
		return false
	***REMOVED***

	if !bytes.Equal(rv.Value, rv2.Value) ***REMOVED***
		return false
	***REMOVED***

	return true
***REMOVED***

// UnmarshalWithRegistry performs the same unmarshalling as Unmarshal but uses the provided registry
// instead of the one attached or the default registry.
func (rv RawValue) UnmarshalWithRegistry(r *bsoncodec.Registry, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if r == nil ***REMOVED***
		return ErrNilRegistry
	***REMOVED***

	vr := bsonrw.NewBSONValueReader(rv.Type, rv.Value)
	rval := reflect.ValueOf(val)
	if rval.Kind() != reflect.Ptr ***REMOVED***
		return fmt.Errorf("argument to Unmarshal* must be a pointer to a type, but got %v", rval)
	***REMOVED***
	rval = rval.Elem()
	dec, err := r.LookupDecoder(rval.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return dec.DecodeValue(bsoncodec.DecodeContext***REMOVED***Registry: r***REMOVED***, vr, rval)
***REMOVED***

// UnmarshalWithContext performs the same unmarshalling as Unmarshal but uses the provided DecodeContext
// instead of the one attached or the default registry.
func (rv RawValue) UnmarshalWithContext(dc *bsoncodec.DecodeContext, val interface***REMOVED******REMOVED***) error ***REMOVED***
	if dc == nil ***REMOVED***
		return ErrNilContext
	***REMOVED***

	vr := bsonrw.NewBSONValueReader(rv.Type, rv.Value)
	rval := reflect.ValueOf(val)
	if rval.Kind() != reflect.Ptr ***REMOVED***
		return fmt.Errorf("argument to Unmarshal* must be a pointer to a type, but got %v", rval)
	***REMOVED***
	rval = rval.Elem()
	dec, err := dc.LookupDecoder(rval.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return dec.DecodeValue(*dc, vr, rval)
***REMOVED***

func convertFromCoreValue(v bsoncore.Value) RawValue ***REMOVED*** return RawValue***REMOVED***Type: v.Type, Value: v.Data***REMOVED*** ***REMOVED***
func convertToCoreValue(v RawValue) bsoncore.Value ***REMOVED***
	return bsoncore.Value***REMOVED***Type: v.Type, Data: v.Value***REMOVED***
***REMOVED***

// Validate ensures the value is a valid BSON value.
func (rv RawValue) Validate() error ***REMOVED*** return convertToCoreValue(rv).Validate() ***REMOVED***

// IsNumber returns true if the type of v is a numeric BSON type.
func (rv RawValue) IsNumber() bool ***REMOVED*** return convertToCoreValue(rv).IsNumber() ***REMOVED***

// String implements the fmt.String interface. This method will return values in extended JSON
// format. If the value is not valid, this returns an empty string
func (rv RawValue) String() string ***REMOVED*** return convertToCoreValue(rv).String() ***REMOVED***

// DebugString outputs a human readable version of Document. It will attempt to stringify the
// valid components of the document even if the entire document is not valid.
func (rv RawValue) DebugString() string ***REMOVED*** return convertToCoreValue(rv).DebugString() ***REMOVED***

// Double returns the float64 value for this element.
// It panics if e's BSON type is not bsontype.Double.
func (rv RawValue) Double() float64 ***REMOVED*** return convertToCoreValue(rv).Double() ***REMOVED***

// DoubleOK is the same as Double, but returns a boolean instead of panicking.
func (rv RawValue) DoubleOK() (float64, bool) ***REMOVED*** return convertToCoreValue(rv).DoubleOK() ***REMOVED***

// StringValue returns the string value for this element.
// It panics if e's BSON type is not bsontype.String.
//
// NOTE: This method is called StringValue to avoid a collision with the String method which
// implements the fmt.Stringer interface.
func (rv RawValue) StringValue() string ***REMOVED*** return convertToCoreValue(rv).StringValue() ***REMOVED***

// StringValueOK is the same as StringValue, but returns a boolean instead of
// panicking.
func (rv RawValue) StringValueOK() (string, bool) ***REMOVED*** return convertToCoreValue(rv).StringValueOK() ***REMOVED***

// Document returns the BSON document the Value represents as a Document. It panics if the
// value is a BSON type other than document.
func (rv RawValue) Document() Raw ***REMOVED*** return Raw(convertToCoreValue(rv).Document()) ***REMOVED***

// DocumentOK is the same as Document, except it returns a boolean
// instead of panicking.
func (rv RawValue) DocumentOK() (Raw, bool) ***REMOVED***
	doc, ok := convertToCoreValue(rv).DocumentOK()
	return Raw(doc), ok
***REMOVED***

// Array returns the BSON array the Value represents as an Array. It panics if the
// value is a BSON type other than array.
func (rv RawValue) Array() Raw ***REMOVED*** return Raw(convertToCoreValue(rv).Array()) ***REMOVED***

// ArrayOK is the same as Array, except it returns a boolean instead
// of panicking.
func (rv RawValue) ArrayOK() (Raw, bool) ***REMOVED***
	doc, ok := convertToCoreValue(rv).ArrayOK()
	return Raw(doc), ok
***REMOVED***

// Binary returns the BSON binary value the Value represents. It panics if the value is a BSON type
// other than binary.
func (rv RawValue) Binary() (subtype byte, data []byte) ***REMOVED*** return convertToCoreValue(rv).Binary() ***REMOVED***

// BinaryOK is the same as Binary, except it returns a boolean instead of
// panicking.
func (rv RawValue) BinaryOK() (subtype byte, data []byte, ok bool) ***REMOVED***
	return convertToCoreValue(rv).BinaryOK()
***REMOVED***

// ObjectID returns the BSON objectid value the Value represents. It panics if the value is a BSON
// type other than objectid.
func (rv RawValue) ObjectID() primitive.ObjectID ***REMOVED*** return convertToCoreValue(rv).ObjectID() ***REMOVED***

// ObjectIDOK is the same as ObjectID, except it returns a boolean instead of
// panicking.
func (rv RawValue) ObjectIDOK() (primitive.ObjectID, bool) ***REMOVED***
	return convertToCoreValue(rv).ObjectIDOK()
***REMOVED***

// Boolean returns the boolean value the Value represents. It panics if the
// value is a BSON type other than boolean.
func (rv RawValue) Boolean() bool ***REMOVED*** return convertToCoreValue(rv).Boolean() ***REMOVED***

// BooleanOK is the same as Boolean, except it returns a boolean instead of
// panicking.
func (rv RawValue) BooleanOK() (bool, bool) ***REMOVED*** return convertToCoreValue(rv).BooleanOK() ***REMOVED***

// DateTime returns the BSON datetime value the Value represents as a
// unix timestamp. It panics if the value is a BSON type other than datetime.
func (rv RawValue) DateTime() int64 ***REMOVED*** return convertToCoreValue(rv).DateTime() ***REMOVED***

// DateTimeOK is the same as DateTime, except it returns a boolean instead of
// panicking.
func (rv RawValue) DateTimeOK() (int64, bool) ***REMOVED*** return convertToCoreValue(rv).DateTimeOK() ***REMOVED***

// Time returns the BSON datetime value the Value represents. It panics if the value is a BSON
// type other than datetime.
func (rv RawValue) Time() time.Time ***REMOVED*** return convertToCoreValue(rv).Time() ***REMOVED***

// TimeOK is the same as Time, except it returns a boolean instead of
// panicking.
func (rv RawValue) TimeOK() (time.Time, bool) ***REMOVED*** return convertToCoreValue(rv).TimeOK() ***REMOVED***

// Regex returns the BSON regex value the Value represents. It panics if the value is a BSON
// type other than regex.
func (rv RawValue) Regex() (pattern, options string) ***REMOVED*** return convertToCoreValue(rv).Regex() ***REMOVED***

// RegexOK is the same as Regex, except it returns a boolean instead of
// panicking.
func (rv RawValue) RegexOK() (pattern, options string, ok bool) ***REMOVED***
	return convertToCoreValue(rv).RegexOK()
***REMOVED***

// DBPointer returns the BSON dbpointer value the Value represents. It panics if the value is a BSON
// type other than DBPointer.
func (rv RawValue) DBPointer() (string, primitive.ObjectID) ***REMOVED***
	return convertToCoreValue(rv).DBPointer()
***REMOVED***

// DBPointerOK is the same as DBPoitner, except that it returns a boolean
// instead of panicking.
func (rv RawValue) DBPointerOK() (string, primitive.ObjectID, bool) ***REMOVED***
	return convertToCoreValue(rv).DBPointerOK()
***REMOVED***

// JavaScript returns the BSON JavaScript code value the Value represents. It panics if the value is
// a BSON type other than JavaScript code.
func (rv RawValue) JavaScript() string ***REMOVED*** return convertToCoreValue(rv).JavaScript() ***REMOVED***

// JavaScriptOK is the same as Javascript, excepti that it returns a boolean
// instead of panicking.
func (rv RawValue) JavaScriptOK() (string, bool) ***REMOVED*** return convertToCoreValue(rv).JavaScriptOK() ***REMOVED***

// Symbol returns the BSON symbol value the Value represents. It panics if the value is a BSON
// type other than symbol.
func (rv RawValue) Symbol() string ***REMOVED*** return convertToCoreValue(rv).Symbol() ***REMOVED***

// SymbolOK is the same as Symbol, excepti that it returns a boolean
// instead of panicking.
func (rv RawValue) SymbolOK() (string, bool) ***REMOVED*** return convertToCoreValue(rv).SymbolOK() ***REMOVED***

// CodeWithScope returns the BSON JavaScript code with scope the Value represents.
// It panics if the value is a BSON type other than JavaScript code with scope.
func (rv RawValue) CodeWithScope() (string, Raw) ***REMOVED***
	code, scope := convertToCoreValue(rv).CodeWithScope()
	return code, Raw(scope)
***REMOVED***

// CodeWithScopeOK is the same as CodeWithScope, except that it returns a boolean instead of
// panicking.
func (rv RawValue) CodeWithScopeOK() (string, Raw, bool) ***REMOVED***
	code, scope, ok := convertToCoreValue(rv).CodeWithScopeOK()
	return code, Raw(scope), ok
***REMOVED***

// Int32 returns the int32 the Value represents. It panics if the value is a BSON type other than
// int32.
func (rv RawValue) Int32() int32 ***REMOVED*** return convertToCoreValue(rv).Int32() ***REMOVED***

// Int32OK is the same as Int32, except that it returns a boolean instead of
// panicking.
func (rv RawValue) Int32OK() (int32, bool) ***REMOVED*** return convertToCoreValue(rv).Int32OK() ***REMOVED***

// AsInt32 returns a BSON number as an int32. If the BSON type is not a numeric one, this method
// will panic.
func (rv RawValue) AsInt32() int32 ***REMOVED*** return convertToCoreValue(rv).AsInt32() ***REMOVED***

// AsInt32OK is the same as AsInt32, except that it returns a boolean instead of
// panicking.
func (rv RawValue) AsInt32OK() (int32, bool) ***REMOVED*** return convertToCoreValue(rv).AsInt32OK() ***REMOVED***

// Timestamp returns the BSON timestamp value the Value represents. It panics if the value is a
// BSON type other than timestamp.
func (rv RawValue) Timestamp() (t, i uint32) ***REMOVED*** return convertToCoreValue(rv).Timestamp() ***REMOVED***

// TimestampOK is the same as Timestamp, except that it returns a boolean
// instead of panicking.
func (rv RawValue) TimestampOK() (t, i uint32, ok bool) ***REMOVED*** return convertToCoreValue(rv).TimestampOK() ***REMOVED***

// Int64 returns the int64 the Value represents. It panics if the value is a BSON type other than
// int64.
func (rv RawValue) Int64() int64 ***REMOVED*** return convertToCoreValue(rv).Int64() ***REMOVED***

// Int64OK is the same as Int64, except that it returns a boolean instead of
// panicking.
func (rv RawValue) Int64OK() (int64, bool) ***REMOVED*** return convertToCoreValue(rv).Int64OK() ***REMOVED***

// AsInt64 returns a BSON number as an int64. If the BSON type is not a numeric one, this method
// will panic.
func (rv RawValue) AsInt64() int64 ***REMOVED*** return convertToCoreValue(rv).AsInt64() ***REMOVED***

// AsInt64OK is the same as AsInt64, except that it returns a boolean instead of
// panicking.
func (rv RawValue) AsInt64OK() (int64, bool) ***REMOVED*** return convertToCoreValue(rv).AsInt64OK() ***REMOVED***

// Decimal128 returns the decimal the Value represents. It panics if the value is a BSON type other than
// decimal.
func (rv RawValue) Decimal128() primitive.Decimal128 ***REMOVED*** return convertToCoreValue(rv).Decimal128() ***REMOVED***

// Decimal128OK is the same as Decimal128, except that it returns a boolean
// instead of panicking.
func (rv RawValue) Decimal128OK() (primitive.Decimal128, bool) ***REMOVED***
	return convertToCoreValue(rv).Decimal128OK()
***REMOVED***
