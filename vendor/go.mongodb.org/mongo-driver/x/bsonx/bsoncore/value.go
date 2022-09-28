// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncore

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ElementTypeError specifies that a method to obtain a BSON value an incorrect type was called on a bson.Value.
type ElementTypeError struct ***REMOVED***
	Method string
	Type   bsontype.Type
***REMOVED***

// Error implements the error interface.
func (ete ElementTypeError) Error() string ***REMOVED***
	return "Call of " + ete.Method + " on " + ete.Type.String() + " type"
***REMOVED***

// Value represents a BSON value with a type and raw bytes.
type Value struct ***REMOVED***
	Type bsontype.Type
	Data []byte
***REMOVED***

// Validate ensures the value is a valid BSON value.
func (v Value) Validate() error ***REMOVED***
	_, _, valid := readValue(v.Data, v.Type)
	if !valid ***REMOVED***
		return NewInsufficientBytesError(v.Data, v.Data)
	***REMOVED***
	return nil
***REMOVED***

// IsNumber returns true if the type of v is a numeric BSON type.
func (v Value) IsNumber() bool ***REMOVED***
	switch v.Type ***REMOVED***
	case bsontype.Double, bsontype.Int32, bsontype.Int64, bsontype.Decimal128:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// AsInt32 returns a BSON number as an int32. If the BSON type is not a numeric one, this method
// will panic.
//
// TODO(skriptble): Add support for Decimal128.
func (v Value) AsInt32() int32 ***REMOVED***
	if !v.IsNumber() ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.AsInt32", v.Type***REMOVED***)
	***REMOVED***
	var i32 int32
	switch v.Type ***REMOVED***
	case bsontype.Double:
		f64, _, ok := ReadDouble(v.Data)
		if !ok ***REMOVED***
			panic(NewInsufficientBytesError(v.Data, v.Data))
		***REMOVED***
		i32 = int32(f64)
	case bsontype.Int32:
		var ok bool
		i32, _, ok = ReadInt32(v.Data)
		if !ok ***REMOVED***
			panic(NewInsufficientBytesError(v.Data, v.Data))
		***REMOVED***
	case bsontype.Int64:
		i64, _, ok := ReadInt64(v.Data)
		if !ok ***REMOVED***
			panic(NewInsufficientBytesError(v.Data, v.Data))
		***REMOVED***
		i32 = int32(i64)
	case bsontype.Decimal128:
		panic(ElementTypeError***REMOVED***"bsoncore.Value.AsInt32", v.Type***REMOVED***)
	***REMOVED***
	return i32
***REMOVED***

// AsInt32OK functions the same as AsInt32 but returns a boolean instead of panicking. False
// indicates an error.
//
// TODO(skriptble): Add support for Decimal128.
func (v Value) AsInt32OK() (int32, bool) ***REMOVED***
	if !v.IsNumber() ***REMOVED***
		return 0, false
	***REMOVED***
	var i32 int32
	switch v.Type ***REMOVED***
	case bsontype.Double:
		f64, _, ok := ReadDouble(v.Data)
		if !ok ***REMOVED***
			return 0, false
		***REMOVED***
		i32 = int32(f64)
	case bsontype.Int32:
		var ok bool
		i32, _, ok = ReadInt32(v.Data)
		if !ok ***REMOVED***
			return 0, false
		***REMOVED***
	case bsontype.Int64:
		i64, _, ok := ReadInt64(v.Data)
		if !ok ***REMOVED***
			return 0, false
		***REMOVED***
		i32 = int32(i64)
	case bsontype.Decimal128:
		return 0, false
	***REMOVED***
	return i32, true
***REMOVED***

// AsInt64 returns a BSON number as an int64. If the BSON type is not a numeric one, this method
// will panic.
//
// TODO(skriptble): Add support for Decimal128.
func (v Value) AsInt64() int64 ***REMOVED***
	if !v.IsNumber() ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.AsInt64", v.Type***REMOVED***)
	***REMOVED***
	var i64 int64
	switch v.Type ***REMOVED***
	case bsontype.Double:
		f64, _, ok := ReadDouble(v.Data)
		if !ok ***REMOVED***
			panic(NewInsufficientBytesError(v.Data, v.Data))
		***REMOVED***
		i64 = int64(f64)
	case bsontype.Int32:
		var ok bool
		i32, _, ok := ReadInt32(v.Data)
		if !ok ***REMOVED***
			panic(NewInsufficientBytesError(v.Data, v.Data))
		***REMOVED***
		i64 = int64(i32)
	case bsontype.Int64:
		var ok bool
		i64, _, ok = ReadInt64(v.Data)
		if !ok ***REMOVED***
			panic(NewInsufficientBytesError(v.Data, v.Data))
		***REMOVED***
	case bsontype.Decimal128:
		panic(ElementTypeError***REMOVED***"bsoncore.Value.AsInt64", v.Type***REMOVED***)
	***REMOVED***
	return i64
***REMOVED***

// AsInt64OK functions the same as AsInt64 but returns a boolean instead of panicking. False
// indicates an error.
//
// TODO(skriptble): Add support for Decimal128.
func (v Value) AsInt64OK() (int64, bool) ***REMOVED***
	if !v.IsNumber() ***REMOVED***
		return 0, false
	***REMOVED***
	var i64 int64
	switch v.Type ***REMOVED***
	case bsontype.Double:
		f64, _, ok := ReadDouble(v.Data)
		if !ok ***REMOVED***
			return 0, false
		***REMOVED***
		i64 = int64(f64)
	case bsontype.Int32:
		var ok bool
		i32, _, ok := ReadInt32(v.Data)
		if !ok ***REMOVED***
			return 0, false
		***REMOVED***
		i64 = int64(i32)
	case bsontype.Int64:
		var ok bool
		i64, _, ok = ReadInt64(v.Data)
		if !ok ***REMOVED***
			return 0, false
		***REMOVED***
	case bsontype.Decimal128:
		return 0, false
	***REMOVED***
	return i64, true
***REMOVED***

// AsFloat64 returns a BSON number as an float64. If the BSON type is not a numeric one, this method
// will panic.
//
// TODO(skriptble): Add support for Decimal128.
func (v Value) AsFloat64() float64 ***REMOVED*** return 0 ***REMOVED***

// AsFloat64OK functions the same as AsFloat64 but returns a boolean instead of panicking. False
// indicates an error.
//
// TODO(skriptble): Add support for Decimal128.
func (v Value) AsFloat64OK() (float64, bool) ***REMOVED*** return 0, false ***REMOVED***

// Add will add this value to another. This is currently only implemented for strings and numbers.
// If either value is a string, the other type is coerced into a string and added to the other.
//
// This method will alter v and will attempt to reuse the []byte of v. If the []byte is too small,
// it will be expanded.
func (v *Value) Add(v2 Value) error ***REMOVED*** return nil ***REMOVED***

// Equal compaes v to v2 and returns true if they are equal.
func (v Value) Equal(v2 Value) bool ***REMOVED***
	if v.Type != v2.Type ***REMOVED***
		return false
	***REMOVED***

	return bytes.Equal(v.Data, v2.Data)
***REMOVED***

// String implements the fmt.String interface. This method will return values in extended JSON
// format. If the value is not valid, this returns an empty string
func (v Value) String() string ***REMOVED***
	switch v.Type ***REMOVED***
	case bsontype.Double:
		f64, ok := v.DoubleOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$numberDouble":"%s"***REMOVED***`, formatDouble(f64))
	case bsontype.String:
		str, ok := v.StringValueOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return escapeString(str)
	case bsontype.EmbeddedDocument:
		doc, ok := v.DocumentOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return doc.String()
	case bsontype.Array:
		arr, ok := v.ArrayOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return arr.String()
	case bsontype.Binary:
		subtype, data, ok := v.BinaryOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$binary":***REMOVED***"base64":"%s","subType":"%02x"***REMOVED******REMOVED***`, base64.StdEncoding.EncodeToString(data), subtype)
	case bsontype.Undefined:
		return `***REMOVED***"$undefined":true***REMOVED***`
	case bsontype.ObjectID:
		oid, ok := v.ObjectIDOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$oid":"%s"***REMOVED***`, oid.Hex())
	case bsontype.Boolean:
		b, ok := v.BooleanOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return strconv.FormatBool(b)
	case bsontype.DateTime:
		dt, ok := v.DateTimeOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$date":***REMOVED***"$numberLong":"%d"***REMOVED******REMOVED***`, dt)
	case bsontype.Null:
		return "null"
	case bsontype.Regex:
		pattern, options, ok := v.RegexOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(
			`***REMOVED***"$regularExpression":***REMOVED***"pattern":%s,"options":"%s"***REMOVED******REMOVED***`,
			escapeString(pattern), sortStringAlphebeticAscending(options),
		)
	case bsontype.DBPointer:
		ns, pointer, ok := v.DBPointerOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$dbPointer":***REMOVED***"$ref":%s,"$id":***REMOVED***"$oid":"%s"***REMOVED******REMOVED******REMOVED***`, escapeString(ns), pointer.Hex())
	case bsontype.JavaScript:
		js, ok := v.JavaScriptOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$code":%s***REMOVED***`, escapeString(js))
	case bsontype.Symbol:
		symbol, ok := v.SymbolOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$symbol":%s***REMOVED***`, escapeString(symbol))
	case bsontype.CodeWithScope:
		code, scope, ok := v.CodeWithScopeOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$code":%s,"$scope":%s***REMOVED***`, code, scope)
	case bsontype.Int32:
		i32, ok := v.Int32OK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$numberInt":"%d"***REMOVED***`, i32)
	case bsontype.Timestamp:
		t, i, ok := v.TimestampOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$timestamp":***REMOVED***"t":"%s","i":"%s"***REMOVED******REMOVED***`, strconv.FormatUint(uint64(t), 10), strconv.FormatUint(uint64(i), 10))
	case bsontype.Int64:
		i64, ok := v.Int64OK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$numberLong":"%d"***REMOVED***`, i64)
	case bsontype.Decimal128:
		d128, ok := v.Decimal128OK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$numberDecimal":"%s"***REMOVED***`, d128.String())
	case bsontype.MinKey:
		return `***REMOVED***"$minKey":1***REMOVED***`
	case bsontype.MaxKey:
		return `***REMOVED***"$maxKey":1***REMOVED***`
	default:
		return ""
	***REMOVED***
***REMOVED***

// DebugString outputs a human readable version of Document. It will attempt to stringify the
// valid components of the document even if the entire document is not valid.
func (v Value) DebugString() string ***REMOVED***
	switch v.Type ***REMOVED***
	case bsontype.String:
		str, ok := v.StringValueOK()
		if !ok ***REMOVED***
			return "<malformed>"
		***REMOVED***
		return escapeString(str)
	case bsontype.EmbeddedDocument:
		doc, ok := v.DocumentOK()
		if !ok ***REMOVED***
			return "<malformed>"
		***REMOVED***
		return doc.DebugString()
	case bsontype.Array:
		arr, ok := v.ArrayOK()
		if !ok ***REMOVED***
			return "<malformed>"
		***REMOVED***
		return arr.DebugString()
	case bsontype.CodeWithScope:
		code, scope, ok := v.CodeWithScopeOK()
		if !ok ***REMOVED***
			return ""
		***REMOVED***
		return fmt.Sprintf(`***REMOVED***"$code":%s,"$scope":%s***REMOVED***`, code, scope.DebugString())
	default:
		str := v.String()
		if str == "" ***REMOVED***
			return "<malformed>"
		***REMOVED***
		return str
	***REMOVED***
***REMOVED***

// Double returns the float64 value for this element.
// It panics if e's BSON type is not bsontype.Double.
func (v Value) Double() float64 ***REMOVED***
	if v.Type != bsontype.Double ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Double", v.Type***REMOVED***)
	***REMOVED***
	f64, _, ok := ReadDouble(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return f64
***REMOVED***

// DoubleOK is the same as Double, but returns a boolean instead of panicking.
func (v Value) DoubleOK() (float64, bool) ***REMOVED***
	if v.Type != bsontype.Double ***REMOVED***
		return 0, false
	***REMOVED***
	f64, _, ok := ReadDouble(v.Data)
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	return f64, true
***REMOVED***

// StringValue returns the string balue for this element.
// It panics if e's BSON type is not bsontype.String.
//
// NOTE: This method is called StringValue to avoid a collision with the String method which
// implements the fmt.Stringer interface.
func (v Value) StringValue() string ***REMOVED***
	if v.Type != bsontype.String ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.StringValue", v.Type***REMOVED***)
	***REMOVED***
	str, _, ok := ReadString(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return str
***REMOVED***

// StringValueOK is the same as StringValue, but returns a boolean instead of
// panicking.
func (v Value) StringValueOK() (string, bool) ***REMOVED***
	if v.Type != bsontype.String ***REMOVED***
		return "", false
	***REMOVED***
	str, _, ok := ReadString(v.Data)
	if !ok ***REMOVED***
		return "", false
	***REMOVED***
	return str, true
***REMOVED***

// Document returns the BSON document the Value represents as a Document. It panics if the
// value is a BSON type other than document.
func (v Value) Document() Document ***REMOVED***
	if v.Type != bsontype.EmbeddedDocument ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Document", v.Type***REMOVED***)
	***REMOVED***
	doc, _, ok := ReadDocument(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return doc
***REMOVED***

// DocumentOK is the same as Document, except it returns a boolean
// instead of panicking.
func (v Value) DocumentOK() (Document, bool) ***REMOVED***
	if v.Type != bsontype.EmbeddedDocument ***REMOVED***
		return nil, false
	***REMOVED***
	doc, _, ok := ReadDocument(v.Data)
	if !ok ***REMOVED***
		return nil, false
	***REMOVED***
	return doc, true
***REMOVED***

// Array returns the BSON array the Value represents as an Array. It panics if the
// value is a BSON type other than array.
func (v Value) Array() Array ***REMOVED***
	if v.Type != bsontype.Array ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Array", v.Type***REMOVED***)
	***REMOVED***
	arr, _, ok := ReadArray(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return arr
***REMOVED***

// ArrayOK is the same as Array, except it returns a boolean instead
// of panicking.
func (v Value) ArrayOK() (Array, bool) ***REMOVED***
	if v.Type != bsontype.Array ***REMOVED***
		return nil, false
	***REMOVED***
	arr, _, ok := ReadArray(v.Data)
	if !ok ***REMOVED***
		return nil, false
	***REMOVED***
	return arr, true
***REMOVED***

// Binary returns the BSON binary value the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Value) Binary() (subtype byte, data []byte) ***REMOVED***
	if v.Type != bsontype.Binary ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Binary", v.Type***REMOVED***)
	***REMOVED***
	subtype, data, _, ok := ReadBinary(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return subtype, data
***REMOVED***

// BinaryOK is the same as Binary, except it returns a boolean instead of
// panicking.
func (v Value) BinaryOK() (subtype byte, data []byte, ok bool) ***REMOVED***
	if v.Type != bsontype.Binary ***REMOVED***
		return 0x00, nil, false
	***REMOVED***
	subtype, data, _, ok = ReadBinary(v.Data)
	if !ok ***REMOVED***
		return 0x00, nil, false
	***REMOVED***
	return subtype, data, true
***REMOVED***

// ObjectID returns the BSON objectid value the Value represents. It panics if the value is a BSON
// type other than objectid.
func (v Value) ObjectID() primitive.ObjectID ***REMOVED***
	if v.Type != bsontype.ObjectID ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.ObjectID", v.Type***REMOVED***)
	***REMOVED***
	oid, _, ok := ReadObjectID(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return oid
***REMOVED***

// ObjectIDOK is the same as ObjectID, except it returns a boolean instead of
// panicking.
func (v Value) ObjectIDOK() (primitive.ObjectID, bool) ***REMOVED***
	if v.Type != bsontype.ObjectID ***REMOVED***
		return primitive.ObjectID***REMOVED******REMOVED***, false
	***REMOVED***
	oid, _, ok := ReadObjectID(v.Data)
	if !ok ***REMOVED***
		return primitive.ObjectID***REMOVED******REMOVED***, false
	***REMOVED***
	return oid, true
***REMOVED***

// Boolean returns the boolean value the Value represents. It panics if the
// value is a BSON type other than boolean.
func (v Value) Boolean() bool ***REMOVED***
	if v.Type != bsontype.Boolean ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Boolean", v.Type***REMOVED***)
	***REMOVED***
	b, _, ok := ReadBoolean(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return b
***REMOVED***

// BooleanOK is the same as Boolean, except it returns a boolean instead of
// panicking.
func (v Value) BooleanOK() (bool, bool) ***REMOVED***
	if v.Type != bsontype.Boolean ***REMOVED***
		return false, false
	***REMOVED***
	b, _, ok := ReadBoolean(v.Data)
	if !ok ***REMOVED***
		return false, false
	***REMOVED***
	return b, true
***REMOVED***

// DateTime returns the BSON datetime value the Value represents as a
// unix timestamp. It panics if the value is a BSON type other than datetime.
func (v Value) DateTime() int64 ***REMOVED***
	if v.Type != bsontype.DateTime ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.DateTime", v.Type***REMOVED***)
	***REMOVED***
	dt, _, ok := ReadDateTime(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return dt
***REMOVED***

// DateTimeOK is the same as DateTime, except it returns a boolean instead of
// panicking.
func (v Value) DateTimeOK() (int64, bool) ***REMOVED***
	if v.Type != bsontype.DateTime ***REMOVED***
		return 0, false
	***REMOVED***
	dt, _, ok := ReadDateTime(v.Data)
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	return dt, true
***REMOVED***

// Time returns the BSON datetime value the Value represents. It panics if the value is a BSON
// type other than datetime.
func (v Value) Time() time.Time ***REMOVED***
	if v.Type != bsontype.DateTime ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Time", v.Type***REMOVED***)
	***REMOVED***
	dt, _, ok := ReadDateTime(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return time.Unix(dt/1000, dt%1000*1000000)
***REMOVED***

// TimeOK is the same as Time, except it returns a boolean instead of
// panicking.
func (v Value) TimeOK() (time.Time, bool) ***REMOVED***
	if v.Type != bsontype.DateTime ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
	dt, _, ok := ReadDateTime(v.Data)
	if !ok ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
	return time.Unix(dt/1000, dt%1000*1000000), true
***REMOVED***

// Regex returns the BSON regex value the Value represents. It panics if the value is a BSON
// type other than regex.
func (v Value) Regex() (pattern, options string) ***REMOVED***
	if v.Type != bsontype.Regex ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Regex", v.Type***REMOVED***)
	***REMOVED***
	pattern, options, _, ok := ReadRegex(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return pattern, options
***REMOVED***

// RegexOK is the same as Regex, except it returns a boolean instead of
// panicking.
func (v Value) RegexOK() (pattern, options string, ok bool) ***REMOVED***
	if v.Type != bsontype.Regex ***REMOVED***
		return "", "", false
	***REMOVED***
	pattern, options, _, ok = ReadRegex(v.Data)
	if !ok ***REMOVED***
		return "", "", false
	***REMOVED***
	return pattern, options, true
***REMOVED***

// DBPointer returns the BSON dbpointer value the Value represents. It panics if the value is a BSON
// type other than DBPointer.
func (v Value) DBPointer() (string, primitive.ObjectID) ***REMOVED***
	if v.Type != bsontype.DBPointer ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.DBPointer", v.Type***REMOVED***)
	***REMOVED***
	ns, pointer, _, ok := ReadDBPointer(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return ns, pointer
***REMOVED***

// DBPointerOK is the same as DBPoitner, except that it returns a boolean
// instead of panicking.
func (v Value) DBPointerOK() (string, primitive.ObjectID, bool) ***REMOVED***
	if v.Type != bsontype.DBPointer ***REMOVED***
		return "", primitive.ObjectID***REMOVED******REMOVED***, false
	***REMOVED***
	ns, pointer, _, ok := ReadDBPointer(v.Data)
	if !ok ***REMOVED***
		return "", primitive.ObjectID***REMOVED******REMOVED***, false
	***REMOVED***
	return ns, pointer, true
***REMOVED***

// JavaScript returns the BSON JavaScript code value the Value represents. It panics if the value is
// a BSON type other than JavaScript code.
func (v Value) JavaScript() string ***REMOVED***
	if v.Type != bsontype.JavaScript ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.JavaScript", v.Type***REMOVED***)
	***REMOVED***
	js, _, ok := ReadJavaScript(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return js
***REMOVED***

// JavaScriptOK is the same as Javascript, excepti that it returns a boolean
// instead of panicking.
func (v Value) JavaScriptOK() (string, bool) ***REMOVED***
	if v.Type != bsontype.JavaScript ***REMOVED***
		return "", false
	***REMOVED***
	js, _, ok := ReadJavaScript(v.Data)
	if !ok ***REMOVED***
		return "", false
	***REMOVED***
	return js, true
***REMOVED***

// Symbol returns the BSON symbol value the Value represents. It panics if the value is a BSON
// type other than symbol.
func (v Value) Symbol() string ***REMOVED***
	if v.Type != bsontype.Symbol ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Symbol", v.Type***REMOVED***)
	***REMOVED***
	symbol, _, ok := ReadSymbol(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return symbol
***REMOVED***

// SymbolOK is the same as Symbol, excepti that it returns a boolean
// instead of panicking.
func (v Value) SymbolOK() (string, bool) ***REMOVED***
	if v.Type != bsontype.Symbol ***REMOVED***
		return "", false
	***REMOVED***
	symbol, _, ok := ReadSymbol(v.Data)
	if !ok ***REMOVED***
		return "", false
	***REMOVED***
	return symbol, true
***REMOVED***

// CodeWithScope returns the BSON JavaScript code with scope the Value represents.
// It panics if the value is a BSON type other than JavaScript code with scope.
func (v Value) CodeWithScope() (string, Document) ***REMOVED***
	if v.Type != bsontype.CodeWithScope ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.CodeWithScope", v.Type***REMOVED***)
	***REMOVED***
	code, scope, _, ok := ReadCodeWithScope(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return code, scope
***REMOVED***

// CodeWithScopeOK is the same as CodeWithScope, except that it returns a boolean instead of
// panicking.
func (v Value) CodeWithScopeOK() (string, Document, bool) ***REMOVED***
	if v.Type != bsontype.CodeWithScope ***REMOVED***
		return "", nil, false
	***REMOVED***
	code, scope, _, ok := ReadCodeWithScope(v.Data)
	if !ok ***REMOVED***
		return "", nil, false
	***REMOVED***
	return code, scope, true
***REMOVED***

// Int32 returns the int32 the Value represents. It panics if the value is a BSON type other than
// int32.
func (v Value) Int32() int32 ***REMOVED***
	if v.Type != bsontype.Int32 ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Int32", v.Type***REMOVED***)
	***REMOVED***
	i32, _, ok := ReadInt32(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return i32
***REMOVED***

// Int32OK is the same as Int32, except that it returns a boolean instead of
// panicking.
func (v Value) Int32OK() (int32, bool) ***REMOVED***
	if v.Type != bsontype.Int32 ***REMOVED***
		return 0, false
	***REMOVED***
	i32, _, ok := ReadInt32(v.Data)
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	return i32, true
***REMOVED***

// Timestamp returns the BSON timestamp value the Value represents. It panics if the value is a
// BSON type other than timestamp.
func (v Value) Timestamp() (t, i uint32) ***REMOVED***
	if v.Type != bsontype.Timestamp ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Timestamp", v.Type***REMOVED***)
	***REMOVED***
	t, i, _, ok := ReadTimestamp(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return t, i
***REMOVED***

// TimestampOK is the same as Timestamp, except that it returns a boolean
// instead of panicking.
func (v Value) TimestampOK() (t, i uint32, ok bool) ***REMOVED***
	if v.Type != bsontype.Timestamp ***REMOVED***
		return 0, 0, false
	***REMOVED***
	t, i, _, ok = ReadTimestamp(v.Data)
	if !ok ***REMOVED***
		return 0, 0, false
	***REMOVED***
	return t, i, true
***REMOVED***

// Int64 returns the int64 the Value represents. It panics if the value is a BSON type other than
// int64.
func (v Value) Int64() int64 ***REMOVED***
	if v.Type != bsontype.Int64 ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Int64", v.Type***REMOVED***)
	***REMOVED***
	i64, _, ok := ReadInt64(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return i64
***REMOVED***

// Int64OK is the same as Int64, except that it returns a boolean instead of
// panicking.
func (v Value) Int64OK() (int64, bool) ***REMOVED***
	if v.Type != bsontype.Int64 ***REMOVED***
		return 0, false
	***REMOVED***
	i64, _, ok := ReadInt64(v.Data)
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	return i64, true
***REMOVED***

// Decimal128 returns the decimal the Value represents. It panics if the value is a BSON type other than
// decimal.
func (v Value) Decimal128() primitive.Decimal128 ***REMOVED***
	if v.Type != bsontype.Decimal128 ***REMOVED***
		panic(ElementTypeError***REMOVED***"bsoncore.Value.Decimal128", v.Type***REMOVED***)
	***REMOVED***
	d128, _, ok := ReadDecimal128(v.Data)
	if !ok ***REMOVED***
		panic(NewInsufficientBytesError(v.Data, v.Data))
	***REMOVED***
	return d128
***REMOVED***

// Decimal128OK is the same as Decimal128, except that it returns a boolean
// instead of panicking.
func (v Value) Decimal128OK() (primitive.Decimal128, bool) ***REMOVED***
	if v.Type != bsontype.Decimal128 ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, false
	***REMOVED***
	d128, _, ok := ReadDecimal128(v.Data)
	if !ok ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, false
	***REMOVED***
	return d128, true
***REMOVED***

var hexChars = "0123456789abcdef"

func escapeString(s string) string ***REMOVED***
	escapeHTML := true
	var buf bytes.Buffer
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); ***REMOVED***
		if b := s[i]; b < utf8.RuneSelf ***REMOVED***
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) ***REMOVED***
				i++
				continue
			***REMOVED***
			if start < i ***REMOVED***
				buf.WriteString(s[start:i])
			***REMOVED***
			switch b ***REMOVED***
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			case '\b':
				buf.WriteByte('\\')
				buf.WriteByte('b')
			case '\f':
				buf.WriteByte('\\')
				buf.WriteByte('f')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hexChars[b>>4])
				buf.WriteByte(hexChars[b&0xF])
			***REMOVED***
			i++
			start = i
			continue
		***REMOVED***
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 ***REMOVED***
			if start < i ***REMOVED***
				buf.WriteString(s[start:i])
			***REMOVED***
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		***REMOVED***
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' ***REMOVED***
			if start < i ***REMOVED***
				buf.WriteString(s[start:i])
			***REMOVED***
			buf.WriteString(`\u202`)
			buf.WriteByte(hexChars[c&0xF])
			i += size
			start = i
			continue
		***REMOVED***
		i += size
	***REMOVED***
	if start < len(s) ***REMOVED***
		buf.WriteString(s[start:])
	***REMOVED***
	buf.WriteByte('"')
	return buf.String()
***REMOVED***

func formatDouble(f float64) string ***REMOVED***
	var s string
	if math.IsInf(f, 1) ***REMOVED***
		s = "Infinity"
	***REMOVED*** else if math.IsInf(f, -1) ***REMOVED***
		s = "-Infinity"
	***REMOVED*** else if math.IsNaN(f) ***REMOVED***
		s = "NaN"
	***REMOVED*** else ***REMOVED***
		// Print exactly one decimalType place for integers; otherwise, print as many are necessary to
		// perfectly represent it.
		s = strconv.FormatFloat(f, 'G', -1, 64)
		if !strings.ContainsRune(s, '.') ***REMOVED***
			s += ".0"
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***

type sortableString []rune

func (ss sortableString) Len() int ***REMOVED***
	return len(ss)
***REMOVED***

func (ss sortableString) Less(i, j int) bool ***REMOVED***
	return ss[i] < ss[j]
***REMOVED***

func (ss sortableString) Swap(i, j int) ***REMOVED***
	oldI := ss[i]
	ss[i] = ss[j]
	ss[j] = oldI
***REMOVED***

func sortStringAlphebeticAscending(s string) string ***REMOVED***
	ss := sortableString([]rune(s))
	sort.Sort(ss)
	return string([]rune(ss))
***REMOVED***
