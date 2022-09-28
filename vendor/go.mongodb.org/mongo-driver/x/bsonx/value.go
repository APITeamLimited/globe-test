// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Val represents a BSON value.
type Val struct ***REMOVED***
	// NOTE: The bootstrap is a small amount of space that'll be on the stack. At 15 bytes this
	// doesn't make this type any larger, since there are 7 bytes of padding and we want an int64 to
	// store small values (e.g. boolean, double, int64, etc...). The primitive property is where all
	// of the larger values go. They will use either Go primitives or the primitive.* types.
	t         bsontype.Type
	bootstrap [15]byte
	primitive interface***REMOVED******REMOVED***
***REMOVED***

func (v Val) string() string ***REMOVED***
	if v.primitive != nil ***REMOVED***
		return v.primitive.(string)
	***REMOVED***
	// The string will either end with a null byte or it fills the entire bootstrap space.
	length := v.bootstrap[0]
	return string(v.bootstrap[1 : length+1])
***REMOVED***

func (v Val) writestring(str string) Val ***REMOVED***
	switch ***REMOVED***
	case len(str) < 15:
		v.bootstrap[0] = uint8(len(str))
		copy(v.bootstrap[1:], str)
	default:
		v.primitive = str
	***REMOVED***
	return v
***REMOVED***

func (v Val) i64() int64 ***REMOVED***
	return int64(v.bootstrap[0]) | int64(v.bootstrap[1])<<8 | int64(v.bootstrap[2])<<16 |
		int64(v.bootstrap[3])<<24 | int64(v.bootstrap[4])<<32 | int64(v.bootstrap[5])<<40 |
		int64(v.bootstrap[6])<<48 | int64(v.bootstrap[7])<<56
***REMOVED***

func (v Val) writei64(i64 int64) Val ***REMOVED***
	v.bootstrap[0] = byte(i64)
	v.bootstrap[1] = byte(i64 >> 8)
	v.bootstrap[2] = byte(i64 >> 16)
	v.bootstrap[3] = byte(i64 >> 24)
	v.bootstrap[4] = byte(i64 >> 32)
	v.bootstrap[5] = byte(i64 >> 40)
	v.bootstrap[6] = byte(i64 >> 48)
	v.bootstrap[7] = byte(i64 >> 56)
	return v
***REMOVED***

// IsZero returns true if this value is zero or a BSON null.
func (v Val) IsZero() bool ***REMOVED*** return v.t == bsontype.Type(0) || v.t == bsontype.Null ***REMOVED***

func (v Val) String() string ***REMOVED***
	// TODO(GODRIVER-612): When bsoncore has appenders for extended JSON use that here.
	return fmt.Sprintf("%v", v.Interface())
***REMOVED***

// Interface returns the Go value of this Value as an empty interface.
//
// This method will return nil if it is empty, otherwise it will return a Go primitive or a
// primitive.* instance.
func (v Val) Interface() interface***REMOVED******REMOVED*** ***REMOVED***
	switch v.Type() ***REMOVED***
	case bsontype.Double:
		return v.Double()
	case bsontype.String:
		return v.StringValue()
	case bsontype.EmbeddedDocument:
		switch v.primitive.(type) ***REMOVED***
		case Doc:
			return v.primitive.(Doc)
		case MDoc:
			return v.primitive.(MDoc)
		default:
			return primitive.Null***REMOVED******REMOVED***
		***REMOVED***
	case bsontype.Array:
		return v.Array()
	case bsontype.Binary:
		return v.primitive.(primitive.Binary)
	case bsontype.Undefined:
		return primitive.Undefined***REMOVED******REMOVED***
	case bsontype.ObjectID:
		return v.ObjectID()
	case bsontype.Boolean:
		return v.Boolean()
	case bsontype.DateTime:
		return v.DateTime()
	case bsontype.Null:
		return primitive.Null***REMOVED******REMOVED***
	case bsontype.Regex:
		return v.primitive.(primitive.Regex)
	case bsontype.DBPointer:
		return v.primitive.(primitive.DBPointer)
	case bsontype.JavaScript:
		return v.JavaScript()
	case bsontype.Symbol:
		return v.Symbol()
	case bsontype.CodeWithScope:
		return v.primitive.(primitive.CodeWithScope)
	case bsontype.Int32:
		return v.Int32()
	case bsontype.Timestamp:
		t, i := v.Timestamp()
		return primitive.Timestamp***REMOVED***T: t, I: i***REMOVED***
	case bsontype.Int64:
		return v.Int64()
	case bsontype.Decimal128:
		return v.Decimal128()
	case bsontype.MinKey:
		return primitive.MinKey***REMOVED******REMOVED***
	case bsontype.MaxKey:
		return primitive.MaxKey***REMOVED******REMOVED***
	default:
		return primitive.Null***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// MarshalBSONValue implements the bsoncodec.ValueMarshaler interface.
func (v Val) MarshalBSONValue() (bsontype.Type, []byte, error) ***REMOVED***
	return v.MarshalAppendBSONValue(nil)
***REMOVED***

// MarshalAppendBSONValue is similar to MarshalBSONValue, but allows the caller to specify a slice
// to add the bytes to.
func (v Val) MarshalAppendBSONValue(dst []byte) (bsontype.Type, []byte, error) ***REMOVED***
	t := v.Type()
	switch v.Type() ***REMOVED***
	case bsontype.Double:
		dst = bsoncore.AppendDouble(dst, v.Double())
	case bsontype.String:
		dst = bsoncore.AppendString(dst, v.String())
	case bsontype.EmbeddedDocument:
		switch v.primitive.(type) ***REMOVED***
		case Doc:
			t, dst, _ = v.primitive.(Doc).MarshalBSONValue() // Doc.MarshalBSONValue never returns an error.
		case MDoc:
			t, dst, _ = v.primitive.(MDoc).MarshalBSONValue() // MDoc.MarshalBSONValue never returns an error.
		***REMOVED***
	case bsontype.Array:
		t, dst, _ = v.Array().MarshalBSONValue() // Arr.MarshalBSON never returns an error.
	case bsontype.Binary:
		subtype, bindata := v.Binary()
		dst = bsoncore.AppendBinary(dst, subtype, bindata)
	case bsontype.Undefined:
	case bsontype.ObjectID:
		dst = bsoncore.AppendObjectID(dst, v.ObjectID())
	case bsontype.Boolean:
		dst = bsoncore.AppendBoolean(dst, v.Boolean())
	case bsontype.DateTime:
		dst = bsoncore.AppendDateTime(dst, v.DateTime())
	case bsontype.Null:
	case bsontype.Regex:
		pattern, options := v.Regex()
		dst = bsoncore.AppendRegex(dst, pattern, options)
	case bsontype.DBPointer:
		ns, ptr := v.DBPointer()
		dst = bsoncore.AppendDBPointer(dst, ns, ptr)
	case bsontype.JavaScript:
		dst = bsoncore.AppendJavaScript(dst, v.JavaScript())
	case bsontype.Symbol:
		dst = bsoncore.AppendSymbol(dst, v.Symbol())
	case bsontype.CodeWithScope:
		code, doc := v.CodeWithScope()
		var scope []byte
		scope, _ = doc.MarshalBSON() // Doc.MarshalBSON never returns an error.
		dst = bsoncore.AppendCodeWithScope(dst, code, scope)
	case bsontype.Int32:
		dst = bsoncore.AppendInt32(dst, v.Int32())
	case bsontype.Timestamp:
		t, i := v.Timestamp()
		dst = bsoncore.AppendTimestamp(dst, t, i)
	case bsontype.Int64:
		dst = bsoncore.AppendInt64(dst, v.Int64())
	case bsontype.Decimal128:
		dst = bsoncore.AppendDecimal128(dst, v.Decimal128())
	case bsontype.MinKey:
	case bsontype.MaxKey:
	default:
		panic(fmt.Errorf("invalid BSON type %v", t))
	***REMOVED***

	return t, dst, nil
***REMOVED***

// UnmarshalBSONValue implements the bsoncodec.ValueUnmarshaler interface.
func (v *Val) UnmarshalBSONValue(t bsontype.Type, data []byte) error ***REMOVED***
	if v == nil ***REMOVED***
		return errors.New("cannot unmarshal into nil Value")
	***REMOVED***
	var err error
	var ok = true
	var rem []byte
	switch t ***REMOVED***
	case bsontype.Double:
		var f64 float64
		f64, rem, ok = bsoncore.ReadDouble(data)
		*v = Double(f64)
	case bsontype.String:
		var str string
		str, rem, ok = bsoncore.ReadString(data)
		*v = String(str)
	case bsontype.EmbeddedDocument:
		var raw []byte
		var doc Doc
		raw, rem, ok = bsoncore.ReadDocument(data)
		doc, err = ReadDoc(raw)
		*v = Document(doc)
	case bsontype.Array:
		var raw []byte
		arr := make(Arr, 0)
		raw, rem, ok = bsoncore.ReadArray(data)
		err = arr.UnmarshalBSONValue(t, raw)
		*v = Array(arr)
	case bsontype.Binary:
		var subtype byte
		var bindata []byte
		subtype, bindata, rem, ok = bsoncore.ReadBinary(data)
		*v = Binary(subtype, bindata)
	case bsontype.Undefined:
		*v = Undefined()
	case bsontype.ObjectID:
		var oid primitive.ObjectID
		oid, rem, ok = bsoncore.ReadObjectID(data)
		*v = ObjectID(oid)
	case bsontype.Boolean:
		var b bool
		b, rem, ok = bsoncore.ReadBoolean(data)
		*v = Boolean(b)
	case bsontype.DateTime:
		var dt int64
		dt, rem, ok = bsoncore.ReadDateTime(data)
		*v = DateTime(dt)
	case bsontype.Null:
		*v = Null()
	case bsontype.Regex:
		var pattern, options string
		pattern, options, rem, ok = bsoncore.ReadRegex(data)
		*v = Regex(pattern, options)
	case bsontype.DBPointer:
		var ns string
		var ptr primitive.ObjectID
		ns, ptr, rem, ok = bsoncore.ReadDBPointer(data)
		*v = DBPointer(ns, ptr)
	case bsontype.JavaScript:
		var js string
		js, rem, ok = bsoncore.ReadJavaScript(data)
		*v = JavaScript(js)
	case bsontype.Symbol:
		var symbol string
		symbol, rem, ok = bsoncore.ReadSymbol(data)
		*v = Symbol(symbol)
	case bsontype.CodeWithScope:
		var raw []byte
		var code string
		var scope Doc
		code, raw, rem, ok = bsoncore.ReadCodeWithScope(data)
		scope, err = ReadDoc(raw)
		*v = CodeWithScope(code, scope)
	case bsontype.Int32:
		var i32 int32
		i32, rem, ok = bsoncore.ReadInt32(data)
		*v = Int32(i32)
	case bsontype.Timestamp:
		var i, t uint32
		t, i, rem, ok = bsoncore.ReadTimestamp(data)
		*v = Timestamp(t, i)
	case bsontype.Int64:
		var i64 int64
		i64, rem, ok = bsoncore.ReadInt64(data)
		*v = Int64(i64)
	case bsontype.Decimal128:
		var d128 primitive.Decimal128
		d128, rem, ok = bsoncore.ReadDecimal128(data)
		*v = Decimal128(d128)
	case bsontype.MinKey:
		*v = MinKey()
	case bsontype.MaxKey:
		*v = MaxKey()
	default:
		err = fmt.Errorf("invalid BSON type %v", t)
	***REMOVED***

	if !ok && err == nil ***REMOVED***
		err = bsoncore.NewInsufficientBytesError(data, rem)
	***REMOVED***

	return err
***REMOVED***

// Type returns the BSON type of this value.
func (v Val) Type() bsontype.Type ***REMOVED***
	if v.t == bsontype.Type(0) ***REMOVED***
		return bsontype.Null
	***REMOVED***
	return v.t
***REMOVED***

// IsNumber returns true if the type of v is a numberic BSON type.
func (v Val) IsNumber() bool ***REMOVED***
	switch v.Type() ***REMOVED***
	case bsontype.Double, bsontype.Int32, bsontype.Int64, bsontype.Decimal128:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// Double returns the BSON double value the Value represents. It panics if the value is a BSON type
// other than double.
func (v Val) Double() float64 ***REMOVED***
	if v.t != bsontype.Double ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Double", v.t***REMOVED***)
	***REMOVED***
	return math.Float64frombits(binary.LittleEndian.Uint64(v.bootstrap[0:8]))
***REMOVED***

// DoubleOK is the same as Double, but returns a boolean instead of panicking.
func (v Val) DoubleOK() (float64, bool) ***REMOVED***
	if v.t != bsontype.Double ***REMOVED***
		return 0, false
	***REMOVED***
	return math.Float64frombits(binary.LittleEndian.Uint64(v.bootstrap[0:8])), true
***REMOVED***

// StringValue returns the BSON string the Value represents. It panics if the value is a BSON type
// other than string.
//
// NOTE: This method is called StringValue to avoid it implementing the
// fmt.Stringer interface.
func (v Val) StringValue() string ***REMOVED***
	if v.t != bsontype.String ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.StringValue", v.t***REMOVED***)
	***REMOVED***
	return v.string()
***REMOVED***

// StringValueOK is the same as StringValue, but returns a boolean instead of
// panicking.
func (v Val) StringValueOK() (string, bool) ***REMOVED***
	if v.t != bsontype.String ***REMOVED***
		return "", false
	***REMOVED***
	return v.string(), true
***REMOVED***

func (v Val) asDoc() Doc ***REMOVED***
	doc, ok := v.primitive.(Doc)
	if ok ***REMOVED***
		return doc
	***REMOVED***
	mdoc := v.primitive.(MDoc)
	for k, v := range mdoc ***REMOVED***
		doc = append(doc, Elem***REMOVED***k, v***REMOVED***)
	***REMOVED***
	return doc
***REMOVED***

func (v Val) asMDoc() MDoc ***REMOVED***
	mdoc, ok := v.primitive.(MDoc)
	if ok ***REMOVED***
		return mdoc
	***REMOVED***
	mdoc = make(MDoc)
	doc := v.primitive.(Doc)
	for _, elem := range doc ***REMOVED***
		mdoc[elem.Key] = elem.Value
	***REMOVED***
	return mdoc
***REMOVED***

// Document returns the BSON embedded document value the Value represents. It panics if the value
// is a BSON type other than embedded document.
func (v Val) Document() Doc ***REMOVED***
	if v.t != bsontype.EmbeddedDocument ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Document", v.t***REMOVED***)
	***REMOVED***
	return v.asDoc()
***REMOVED***

// DocumentOK is the same as Document, except it returns a boolean
// instead of panicking.
func (v Val) DocumentOK() (Doc, bool) ***REMOVED***
	if v.t != bsontype.EmbeddedDocument ***REMOVED***
		return nil, false
	***REMOVED***
	return v.asDoc(), true
***REMOVED***

// MDocument returns the BSON embedded document value the Value represents. It panics if the value
// is a BSON type other than embedded document.
func (v Val) MDocument() MDoc ***REMOVED***
	if v.t != bsontype.EmbeddedDocument ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.MDocument", v.t***REMOVED***)
	***REMOVED***
	return v.asMDoc()
***REMOVED***

// MDocumentOK is the same as Document, except it returns a boolean
// instead of panicking.
func (v Val) MDocumentOK() (MDoc, bool) ***REMOVED***
	if v.t != bsontype.EmbeddedDocument ***REMOVED***
		return nil, false
	***REMOVED***
	return v.asMDoc(), true
***REMOVED***

// Array returns the BSON array value the Value represents. It panics if the value is a BSON type
// other than array.
func (v Val) Array() Arr ***REMOVED***
	if v.t != bsontype.Array ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Array", v.t***REMOVED***)
	***REMOVED***
	return v.primitive.(Arr)
***REMOVED***

// ArrayOK is the same as Array, except it returns a boolean
// instead of panicking.
func (v Val) ArrayOK() (Arr, bool) ***REMOVED***
	if v.t != bsontype.Array ***REMOVED***
		return nil, false
	***REMOVED***
	return v.primitive.(Arr), true
***REMOVED***

// Binary returns the BSON binary value the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) Binary() (byte, []byte) ***REMOVED***
	if v.t != bsontype.Binary ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Binary", v.t***REMOVED***)
	***REMOVED***
	bin := v.primitive.(primitive.Binary)
	return bin.Subtype, bin.Data
***REMOVED***

// BinaryOK is the same as Binary, except it returns a boolean instead of
// panicking.
func (v Val) BinaryOK() (byte, []byte, bool) ***REMOVED***
	if v.t != bsontype.Binary ***REMOVED***
		return 0x00, nil, false
	***REMOVED***
	bin := v.primitive.(primitive.Binary)
	return bin.Subtype, bin.Data, true
***REMOVED***

// Undefined returns the BSON undefined the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) Undefined() ***REMOVED***
	if v.t != bsontype.Undefined ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Undefined", v.t***REMOVED***)
	***REMOVED***
***REMOVED***

// UndefinedOK is the same as Undefined, except it returns a boolean instead of
// panicking.
func (v Val) UndefinedOK() bool ***REMOVED***
	return v.t == bsontype.Undefined
***REMOVED***

// ObjectID returns the BSON ObjectID the Value represents. It panics if the value is a BSON type
// other than ObjectID.
func (v Val) ObjectID() primitive.ObjectID ***REMOVED***
	if v.t != bsontype.ObjectID ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.ObjectID", v.t***REMOVED***)
	***REMOVED***
	var oid primitive.ObjectID
	copy(oid[:], v.bootstrap[:12])
	return oid
***REMOVED***

// ObjectIDOK is the same as ObjectID, except it returns a boolean instead of
// panicking.
func (v Val) ObjectIDOK() (primitive.ObjectID, bool) ***REMOVED***
	if v.t != bsontype.ObjectID ***REMOVED***
		return primitive.ObjectID***REMOVED******REMOVED***, false
	***REMOVED***
	var oid primitive.ObjectID
	copy(oid[:], v.bootstrap[:12])
	return oid, true
***REMOVED***

// Boolean returns the BSON boolean the Value represents. It panics if the value is a BSON type
// other than boolean.
func (v Val) Boolean() bool ***REMOVED***
	if v.t != bsontype.Boolean ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Boolean", v.t***REMOVED***)
	***REMOVED***
	return v.bootstrap[0] == 0x01
***REMOVED***

// BooleanOK is the same as Boolean, except it returns a boolean instead of
// panicking.
func (v Val) BooleanOK() (bool, bool) ***REMOVED***
	if v.t != bsontype.Boolean ***REMOVED***
		return false, false
	***REMOVED***
	return v.bootstrap[0] == 0x01, true
***REMOVED***

// DateTime returns the BSON datetime the Value represents. It panics if the value is a BSON type
// other than datetime.
func (v Val) DateTime() int64 ***REMOVED***
	if v.t != bsontype.DateTime ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.DateTime", v.t***REMOVED***)
	***REMOVED***
	return v.i64()
***REMOVED***

// DateTimeOK is the same as DateTime, except it returns a boolean instead of
// panicking.
func (v Val) DateTimeOK() (int64, bool) ***REMOVED***
	if v.t != bsontype.DateTime ***REMOVED***
		return 0, false
	***REMOVED***
	return v.i64(), true
***REMOVED***

// Time returns the BSON datetime the Value represents as time.Time. It panics if the value is a BSON
// type other than datetime.
func (v Val) Time() time.Time ***REMOVED***
	if v.t != bsontype.DateTime ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Time", v.t***REMOVED***)
	***REMOVED***
	i := v.i64()
	return time.Unix(i/1000, i%1000*1000000)
***REMOVED***

// TimeOK is the same as Time, except it returns a boolean instead of
// panicking.
func (v Val) TimeOK() (time.Time, bool) ***REMOVED***
	if v.t != bsontype.DateTime ***REMOVED***
		return time.Time***REMOVED******REMOVED***, false
	***REMOVED***
	i := v.i64()
	return time.Unix(i/1000, i%1000*1000000), true
***REMOVED***

// Null returns the BSON undefined the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) Null() ***REMOVED***
	if v.t != bsontype.Null && v.t != bsontype.Type(0) ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Null", v.t***REMOVED***)
	***REMOVED***
***REMOVED***

// NullOK is the same as Null, except it returns a boolean instead of
// panicking.
func (v Val) NullOK() bool ***REMOVED***
	if v.t != bsontype.Null && v.t != bsontype.Type(0) ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// Regex returns the BSON regex the Value represents. It panics if the value is a BSON type
// other than regex.
func (v Val) Regex() (pattern, options string) ***REMOVED***
	if v.t != bsontype.Regex ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Regex", v.t***REMOVED***)
	***REMOVED***
	regex := v.primitive.(primitive.Regex)
	return regex.Pattern, regex.Options
***REMOVED***

// RegexOK is the same as Regex, except that it returns a boolean
// instead of panicking.
func (v Val) RegexOK() (pattern, options string, ok bool) ***REMOVED***
	if v.t != bsontype.Regex ***REMOVED***
		return "", "", false
	***REMOVED***
	regex := v.primitive.(primitive.Regex)
	return regex.Pattern, regex.Options, true
***REMOVED***

// DBPointer returns the BSON dbpointer the Value represents. It panics if the value is a BSON type
// other than dbpointer.
func (v Val) DBPointer() (string, primitive.ObjectID) ***REMOVED***
	if v.t != bsontype.DBPointer ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.DBPointer", v.t***REMOVED***)
	***REMOVED***
	dbptr := v.primitive.(primitive.DBPointer)
	return dbptr.DB, dbptr.Pointer
***REMOVED***

// DBPointerOK is the same as DBPoitner, except that it returns a boolean
// instead of panicking.
func (v Val) DBPointerOK() (string, primitive.ObjectID, bool) ***REMOVED***
	if v.t != bsontype.DBPointer ***REMOVED***
		return "", primitive.ObjectID***REMOVED******REMOVED***, false
	***REMOVED***
	dbptr := v.primitive.(primitive.DBPointer)
	return dbptr.DB, dbptr.Pointer, true
***REMOVED***

// JavaScript returns the BSON JavaScript the Value represents. It panics if the value is a BSON type
// other than JavaScript.
func (v Val) JavaScript() string ***REMOVED***
	if v.t != bsontype.JavaScript ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.JavaScript", v.t***REMOVED***)
	***REMOVED***
	return v.string()
***REMOVED***

// JavaScriptOK is the same as Javascript, except that it returns a boolean
// instead of panicking.
func (v Val) JavaScriptOK() (string, bool) ***REMOVED***
	if v.t != bsontype.JavaScript ***REMOVED***
		return "", false
	***REMOVED***
	return v.string(), true
***REMOVED***

// Symbol returns the BSON symbol the Value represents. It panics if the value is a BSON type
// other than symbol.
func (v Val) Symbol() string ***REMOVED***
	if v.t != bsontype.Symbol ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Symbol", v.t***REMOVED***)
	***REMOVED***
	return v.string()
***REMOVED***

// SymbolOK is the same as Javascript, except that it returns a boolean
// instead of panicking.
func (v Val) SymbolOK() (string, bool) ***REMOVED***
	if v.t != bsontype.Symbol ***REMOVED***
		return "", false
	***REMOVED***
	return v.string(), true
***REMOVED***

// CodeWithScope returns the BSON code with scope value the Value represents. It panics if the
// value is a BSON type other than code with scope.
func (v Val) CodeWithScope() (string, Doc) ***REMOVED***
	if v.t != bsontype.CodeWithScope ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.CodeWithScope", v.t***REMOVED***)
	***REMOVED***
	cws := v.primitive.(primitive.CodeWithScope)
	return string(cws.Code), cws.Scope.(Doc)
***REMOVED***

// CodeWithScopeOK is the same as JavascriptWithScope,
// except that it returns a boolean instead of panicking.
func (v Val) CodeWithScopeOK() (string, Doc, bool) ***REMOVED***
	if v.t != bsontype.CodeWithScope ***REMOVED***
		return "", nil, false
	***REMOVED***
	cws := v.primitive.(primitive.CodeWithScope)
	return string(cws.Code), cws.Scope.(Doc), true
***REMOVED***

// Int32 returns the BSON int32 the Value represents. It panics if the value is a BSON type
// other than int32.
func (v Val) Int32() int32 ***REMOVED***
	if v.t != bsontype.Int32 ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Int32", v.t***REMOVED***)
	***REMOVED***
	return int32(v.bootstrap[0]) | int32(v.bootstrap[1])<<8 |
		int32(v.bootstrap[2])<<16 | int32(v.bootstrap[3])<<24
***REMOVED***

// Int32OK is the same as Int32, except that it returns a boolean instead of
// panicking.
func (v Val) Int32OK() (int32, bool) ***REMOVED***
	if v.t != bsontype.Int32 ***REMOVED***
		return 0, false
	***REMOVED***
	return int32(v.bootstrap[0]) | int32(v.bootstrap[1])<<8 |
			int32(v.bootstrap[2])<<16 | int32(v.bootstrap[3])<<24,
		true
***REMOVED***

// Timestamp returns the BSON timestamp the Value represents. It panics if the value is a
// BSON type other than timestamp.
func (v Val) Timestamp() (t, i uint32) ***REMOVED***
	if v.t != bsontype.Timestamp ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Timestamp", v.t***REMOVED***)
	***REMOVED***
	return uint32(v.bootstrap[4]) | uint32(v.bootstrap[5])<<8 |
			uint32(v.bootstrap[6])<<16 | uint32(v.bootstrap[7])<<24,
		uint32(v.bootstrap[0]) | uint32(v.bootstrap[1])<<8 |
			uint32(v.bootstrap[2])<<16 | uint32(v.bootstrap[3])<<24
***REMOVED***

// TimestampOK is the same as Timestamp, except that it returns a boolean
// instead of panicking.
func (v Val) TimestampOK() (t uint32, i uint32, ok bool) ***REMOVED***
	if v.t != bsontype.Timestamp ***REMOVED***
		return 0, 0, false
	***REMOVED***
	return uint32(v.bootstrap[4]) | uint32(v.bootstrap[5])<<8 |
			uint32(v.bootstrap[6])<<16 | uint32(v.bootstrap[7])<<24,
		uint32(v.bootstrap[0]) | uint32(v.bootstrap[1])<<8 |
			uint32(v.bootstrap[2])<<16 | uint32(v.bootstrap[3])<<24,
		true
***REMOVED***

// Int64 returns the BSON int64 the Value represents. It panics if the value is a BSON type
// other than int64.
func (v Val) Int64() int64 ***REMOVED***
	if v.t != bsontype.Int64 ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Int64", v.t***REMOVED***)
	***REMOVED***
	return v.i64()
***REMOVED***

// Int64OK is the same as Int64, except that it returns a boolean instead of
// panicking.
func (v Val) Int64OK() (int64, bool) ***REMOVED***
	if v.t != bsontype.Int64 ***REMOVED***
		return 0, false
	***REMOVED***
	return v.i64(), true
***REMOVED***

// Decimal128 returns the BSON decimal128 value the Value represents. It panics if the value is a
// BSON type other than decimal128.
func (v Val) Decimal128() primitive.Decimal128 ***REMOVED***
	if v.t != bsontype.Decimal128 ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.Decimal128", v.t***REMOVED***)
	***REMOVED***
	return v.primitive.(primitive.Decimal128)
***REMOVED***

// Decimal128OK is the same as Decimal128, except that it returns a boolean
// instead of panicking.
func (v Val) Decimal128OK() (primitive.Decimal128, bool) ***REMOVED***
	if v.t != bsontype.Decimal128 ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, false
	***REMOVED***
	return v.primitive.(primitive.Decimal128), true
***REMOVED***

// MinKey returns the BSON minkey the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) MinKey() ***REMOVED***
	if v.t != bsontype.MinKey ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.MinKey", v.t***REMOVED***)
	***REMOVED***
***REMOVED***

// MinKeyOK is the same as MinKey, except it returns a boolean instead of
// panicking.
func (v Val) MinKeyOK() bool ***REMOVED***
	return v.t == bsontype.MinKey
***REMOVED***

// MaxKey returns the BSON maxkey the Value represents. It panics if the value is a BSON type
// other than binary.
func (v Val) MaxKey() ***REMOVED***
	if v.t != bsontype.MaxKey ***REMOVED***
		panic(ElementTypeError***REMOVED***"bson.Value.MaxKey", v.t***REMOVED***)
	***REMOVED***
***REMOVED***

// MaxKeyOK is the same as MaxKey, except it returns a boolean instead of
// panicking.
func (v Val) MaxKeyOK() bool ***REMOVED***
	return v.t == bsontype.MaxKey
***REMOVED***

// Equal compares v to v2 and returns true if they are equal. Unknown BSON types are
// never equal. Two empty values are equal.
func (v Val) Equal(v2 Val) bool ***REMOVED***
	if v.Type() != v2.Type() ***REMOVED***
		return false
	***REMOVED***
	if v.IsZero() && v2.IsZero() ***REMOVED***
		return true
	***REMOVED***

	switch v.Type() ***REMOVED***
	case bsontype.Double, bsontype.DateTime, bsontype.Timestamp, bsontype.Int64:
		return bytes.Equal(v.bootstrap[0:8], v2.bootstrap[0:8])
	case bsontype.String:
		return v.string() == v2.string()
	case bsontype.EmbeddedDocument:
		return v.equalDocs(v2)
	case bsontype.Array:
		return v.Array().Equal(v2.Array())
	case bsontype.Binary:
		return v.primitive.(primitive.Binary).Equal(v2.primitive.(primitive.Binary))
	case bsontype.Undefined:
		return true
	case bsontype.ObjectID:
		return bytes.Equal(v.bootstrap[0:12], v2.bootstrap[0:12])
	case bsontype.Boolean:
		return v.bootstrap[0] == v2.bootstrap[0]
	case bsontype.Null:
		return true
	case bsontype.Regex:
		return v.primitive.(primitive.Regex).Equal(v2.primitive.(primitive.Regex))
	case bsontype.DBPointer:
		return v.primitive.(primitive.DBPointer).Equal(v2.primitive.(primitive.DBPointer))
	case bsontype.JavaScript:
		return v.JavaScript() == v2.JavaScript()
	case bsontype.Symbol:
		return v.Symbol() == v2.Symbol()
	case bsontype.CodeWithScope:
		code1, scope1 := v.primitive.(primitive.CodeWithScope).Code, v.primitive.(primitive.CodeWithScope).Scope
		code2, scope2 := v2.primitive.(primitive.CodeWithScope).Code, v2.primitive.(primitive.CodeWithScope).Scope
		return code1 == code2 && v.equalInterfaceDocs(scope1, scope2)
	case bsontype.Int32:
		return v.Int32() == v2.Int32()
	case bsontype.Decimal128:
		h, l := v.Decimal128().GetBytes()
		h2, l2 := v2.Decimal128().GetBytes()
		return h == h2 && l == l2
	case bsontype.MinKey:
		return true
	case bsontype.MaxKey:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func (v Val) equalDocs(v2 Val) bool ***REMOVED***
	_, ok1 := v.primitive.(MDoc)
	_, ok2 := v2.primitive.(MDoc)
	if ok1 || ok2 ***REMOVED***
		return v.asMDoc().Equal(v2.asMDoc())
	***REMOVED***
	return v.asDoc().Equal(v2.asDoc())
***REMOVED***

func (Val) equalInterfaceDocs(i, i2 interface***REMOVED******REMOVED***) bool ***REMOVED***
	switch d := i.(type) ***REMOVED***
	case MDoc:
		d2, ok := i2.(IDoc)
		if !ok ***REMOVED***
			return false
		***REMOVED***
		return d.Equal(d2)
	case Doc:
		d2, ok := i2.(IDoc)
		if !ok ***REMOVED***
			return false
		***REMOVED***
		return d.Equal(d2)
	case nil:
		return i2 == nil
	default:
		return false
	***REMOVED***
***REMOVED***
