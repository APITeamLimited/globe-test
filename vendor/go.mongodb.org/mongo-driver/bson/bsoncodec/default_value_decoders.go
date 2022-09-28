// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var (
	defaultValueDecoders DefaultValueDecoders
	errCannotTruncate    = errors.New("float64 can only be truncated to an integer type when truncation is enabled")
)

type decodeBinaryError struct ***REMOVED***
	subtype  byte
	typeName string
***REMOVED***

func (d decodeBinaryError) Error() string ***REMOVED***
	return fmt.Sprintf("only binary values with subtype 0x00 or 0x02 can be decoded into %s, but got subtype %v", d.typeName, d.subtype)
***REMOVED***

func newDefaultStructCodec() *StructCodec ***REMOVED***
	codec, err := NewStructCodec(DefaultStructTagParser)
	if err != nil ***REMOVED***
		// This function is called from the codec registration path, so errors can't be propagated. If there's an error
		// constructing the StructCodec, we panic to avoid losing it.
		panic(fmt.Errorf("error creating default StructCodec: %v", err))
	***REMOVED***
	return codec
***REMOVED***

// DefaultValueDecoders is a namespace type for the default ValueDecoders used
// when creating a registry.
type DefaultValueDecoders struct***REMOVED******REMOVED***

// RegisterDefaultDecoders will register the decoder methods attached to DefaultValueDecoders with
// the provided RegistryBuilder.
//
// There is no support for decoding map[string]interface***REMOVED******REMOVED*** because there is no decoder for
// interface***REMOVED******REMOVED***, so users must either register this decoder themselves or use the
// EmptyInterfaceDecoder available in the bson package.
func (dvd DefaultValueDecoders) RegisterDefaultDecoders(rb *RegistryBuilder) ***REMOVED***
	if rb == nil ***REMOVED***
		panic(errors.New("argument to RegisterDefaultDecoders must not be nil"))
	***REMOVED***

	intDecoder := decodeAdapter***REMOVED***dvd.IntDecodeValue, dvd.intDecodeType***REMOVED***
	floatDecoder := decodeAdapter***REMOVED***dvd.FloatDecodeValue, dvd.floatDecodeType***REMOVED***

	rb.
		RegisterTypeDecoder(tD, ValueDecoderFunc(dvd.DDecodeValue)).
		RegisterTypeDecoder(tBinary, decodeAdapter***REMOVED***dvd.BinaryDecodeValue, dvd.binaryDecodeType***REMOVED***).
		RegisterTypeDecoder(tUndefined, decodeAdapter***REMOVED***dvd.UndefinedDecodeValue, dvd.undefinedDecodeType***REMOVED***).
		RegisterTypeDecoder(tDateTime, decodeAdapter***REMOVED***dvd.DateTimeDecodeValue, dvd.dateTimeDecodeType***REMOVED***).
		RegisterTypeDecoder(tNull, decodeAdapter***REMOVED***dvd.NullDecodeValue, dvd.nullDecodeType***REMOVED***).
		RegisterTypeDecoder(tRegex, decodeAdapter***REMOVED***dvd.RegexDecodeValue, dvd.regexDecodeType***REMOVED***).
		RegisterTypeDecoder(tDBPointer, decodeAdapter***REMOVED***dvd.DBPointerDecodeValue, dvd.dBPointerDecodeType***REMOVED***).
		RegisterTypeDecoder(tTimestamp, decodeAdapter***REMOVED***dvd.TimestampDecodeValue, dvd.timestampDecodeType***REMOVED***).
		RegisterTypeDecoder(tMinKey, decodeAdapter***REMOVED***dvd.MinKeyDecodeValue, dvd.minKeyDecodeType***REMOVED***).
		RegisterTypeDecoder(tMaxKey, decodeAdapter***REMOVED***dvd.MaxKeyDecodeValue, dvd.maxKeyDecodeType***REMOVED***).
		RegisterTypeDecoder(tJavaScript, decodeAdapter***REMOVED***dvd.JavaScriptDecodeValue, dvd.javaScriptDecodeType***REMOVED***).
		RegisterTypeDecoder(tSymbol, decodeAdapter***REMOVED***dvd.SymbolDecodeValue, dvd.symbolDecodeType***REMOVED***).
		RegisterTypeDecoder(tByteSlice, defaultByteSliceCodec).
		RegisterTypeDecoder(tTime, defaultTimeCodec).
		RegisterTypeDecoder(tEmpty, defaultEmptyInterfaceCodec).
		RegisterTypeDecoder(tCoreArray, defaultArrayCodec).
		RegisterTypeDecoder(tOID, decodeAdapter***REMOVED***dvd.ObjectIDDecodeValue, dvd.objectIDDecodeType***REMOVED***).
		RegisterTypeDecoder(tDecimal, decodeAdapter***REMOVED***dvd.Decimal128DecodeValue, dvd.decimal128DecodeType***REMOVED***).
		RegisterTypeDecoder(tJSONNumber, decodeAdapter***REMOVED***dvd.JSONNumberDecodeValue, dvd.jsonNumberDecodeType***REMOVED***).
		RegisterTypeDecoder(tURL, decodeAdapter***REMOVED***dvd.URLDecodeValue, dvd.urlDecodeType***REMOVED***).
		RegisterTypeDecoder(tCoreDocument, ValueDecoderFunc(dvd.CoreDocumentDecodeValue)).
		RegisterTypeDecoder(tCodeWithScope, decodeAdapter***REMOVED***dvd.CodeWithScopeDecodeValue, dvd.codeWithScopeDecodeType***REMOVED***).
		RegisterDefaultDecoder(reflect.Bool, decodeAdapter***REMOVED***dvd.BooleanDecodeValue, dvd.booleanDecodeType***REMOVED***).
		RegisterDefaultDecoder(reflect.Int, intDecoder).
		RegisterDefaultDecoder(reflect.Int8, intDecoder).
		RegisterDefaultDecoder(reflect.Int16, intDecoder).
		RegisterDefaultDecoder(reflect.Int32, intDecoder).
		RegisterDefaultDecoder(reflect.Int64, intDecoder).
		RegisterDefaultDecoder(reflect.Uint, defaultUIntCodec).
		RegisterDefaultDecoder(reflect.Uint8, defaultUIntCodec).
		RegisterDefaultDecoder(reflect.Uint16, defaultUIntCodec).
		RegisterDefaultDecoder(reflect.Uint32, defaultUIntCodec).
		RegisterDefaultDecoder(reflect.Uint64, defaultUIntCodec).
		RegisterDefaultDecoder(reflect.Float32, floatDecoder).
		RegisterDefaultDecoder(reflect.Float64, floatDecoder).
		RegisterDefaultDecoder(reflect.Array, ValueDecoderFunc(dvd.ArrayDecodeValue)).
		RegisterDefaultDecoder(reflect.Map, defaultMapCodec).
		RegisterDefaultDecoder(reflect.Slice, defaultSliceCodec).
		RegisterDefaultDecoder(reflect.String, defaultStringCodec).
		RegisterDefaultDecoder(reflect.Struct, newDefaultStructCodec()).
		RegisterDefaultDecoder(reflect.Ptr, NewPointerCodec()).
		RegisterTypeMapEntry(bsontype.Double, tFloat64).
		RegisterTypeMapEntry(bsontype.String, tString).
		RegisterTypeMapEntry(bsontype.Array, tA).
		RegisterTypeMapEntry(bsontype.Binary, tBinary).
		RegisterTypeMapEntry(bsontype.Undefined, tUndefined).
		RegisterTypeMapEntry(bsontype.ObjectID, tOID).
		RegisterTypeMapEntry(bsontype.Boolean, tBool).
		RegisterTypeMapEntry(bsontype.DateTime, tDateTime).
		RegisterTypeMapEntry(bsontype.Regex, tRegex).
		RegisterTypeMapEntry(bsontype.DBPointer, tDBPointer).
		RegisterTypeMapEntry(bsontype.JavaScript, tJavaScript).
		RegisterTypeMapEntry(bsontype.Symbol, tSymbol).
		RegisterTypeMapEntry(bsontype.CodeWithScope, tCodeWithScope).
		RegisterTypeMapEntry(bsontype.Int32, tInt32).
		RegisterTypeMapEntry(bsontype.Int64, tInt64).
		RegisterTypeMapEntry(bsontype.Timestamp, tTimestamp).
		RegisterTypeMapEntry(bsontype.Decimal128, tDecimal).
		RegisterTypeMapEntry(bsontype.MinKey, tMinKey).
		RegisterTypeMapEntry(bsontype.MaxKey, tMaxKey).
		RegisterTypeMapEntry(bsontype.Type(0), tD).
		RegisterTypeMapEntry(bsontype.EmbeddedDocument, tD).
		RegisterHookDecoder(tValueUnmarshaler, ValueDecoderFunc(dvd.ValueUnmarshalerDecodeValue)).
		RegisterHookDecoder(tUnmarshaler, ValueDecoderFunc(dvd.UnmarshalerDecodeValue))
***REMOVED***

// DDecodeValue is the ValueDecoderFunc for primitive.D instances.
func (dvd DefaultValueDecoders) DDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || !val.CanSet() || val.Type() != tD ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "DDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Slice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Type(0), bsontype.EmbeddedDocument:
		dc.Ancestor = tD
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	default:
		return fmt.Errorf("cannot decode %v into a primitive.D", vrType)
	***REMOVED***

	dr, err := vr.ReadDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	decoder, err := dc.LookupDecoder(tEmpty)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	tEmptyTypeDecoder, _ := decoder.(typeDecoder)

	// Use the elements in the provided value if it's non nil. Otherwise, allocate a new D instance.
	var elems primitive.D
	if !val.IsNil() ***REMOVED***
		val.SetLen(0)
		elems = val.Interface().(primitive.D)
	***REMOVED*** else ***REMOVED***
		elems = make(primitive.D, 0)
	***REMOVED***

	for ***REMOVED***
		key, elemVr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD ***REMOVED***
			break
		***REMOVED*** else if err != nil ***REMOVED***
			return err
		***REMOVED***

		// Pass false for convert because we don't need to call reflect.Value.Convert for tEmpty.
		elem, err := decodeTypeOrValueWithInfo(decoder, tEmptyTypeDecoder, dc, elemVr, tEmpty, false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		elems = append(elems, primitive.E***REMOVED***Key: key, Value: elem.Interface()***REMOVED***)
	***REMOVED***

	val.Set(reflect.ValueOf(elems))
	return nil
***REMOVED***

func (dvd DefaultValueDecoders) booleanDecodeType(dctx DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t.Kind() != reflect.Bool ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "BooleanDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Bool***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var b bool
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		b = (i32 != 0)
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		b = (i64 != 0)
	case bsontype.Double:
		f64, err := vr.ReadDouble()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		b = (f64 != 0)
	case bsontype.Boolean:
		b, err = vr.ReadBoolean()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a boolean", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(b), nil
***REMOVED***

// BooleanDecodeValue is the ValueDecoderFunc for bool types.
func (dvd DefaultValueDecoders) BooleanDecodeValue(dctx DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || !val.CanSet() || val.Kind() != reflect.Bool ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "BooleanDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Bool***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.booleanDecodeType(dctx, vr, val.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.SetBool(elem.Bool())
	return nil
***REMOVED***

func (DefaultValueDecoders) intDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	var i64 int64
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		i64 = int64(i32)
	case bsontype.Int64:
		i64, err = vr.ReadInt64()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Double:
		f64, err := vr.ReadDouble()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if !dc.Truncate && math.Floor(f64) != f64 ***REMOVED***
			return emptyValue, errCannotTruncate
		***REMOVED***
		if f64 > float64(math.MaxInt64) ***REMOVED***
			return emptyValue, fmt.Errorf("%g overflows int64", f64)
		***REMOVED***
		i64 = int64(f64)
	case bsontype.Boolean:
		b, err := vr.ReadBoolean()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if b ***REMOVED***
			i64 = 1
		***REMOVED***
	case bsontype.Null:
		if err = vr.ReadNull(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Undefined:
		if err = vr.ReadUndefined(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into an integer type", vrType)
	***REMOVED***

	switch t.Kind() ***REMOVED***
	case reflect.Int8:
		if i64 < math.MinInt8 || i64 > math.MaxInt8 ***REMOVED***
			return emptyValue, fmt.Errorf("%d overflows int8", i64)
		***REMOVED***

		return reflect.ValueOf(int8(i64)), nil
	case reflect.Int16:
		if i64 < math.MinInt16 || i64 > math.MaxInt16 ***REMOVED***
			return emptyValue, fmt.Errorf("%d overflows int16", i64)
		***REMOVED***

		return reflect.ValueOf(int16(i64)), nil
	case reflect.Int32:
		if i64 < math.MinInt32 || i64 > math.MaxInt32 ***REMOVED***
			return emptyValue, fmt.Errorf("%d overflows int32", i64)
		***REMOVED***

		return reflect.ValueOf(int32(i64)), nil
	case reflect.Int64:
		return reflect.ValueOf(i64), nil
	case reflect.Int:
		if int64(int(i64)) != i64 ***REMOVED*** // Can we fit this inside of an int
			return emptyValue, fmt.Errorf("%d overflows int", i64)
		***REMOVED***

		return reflect.ValueOf(int(i64)), nil
	default:
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "IntDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***
***REMOVED***

// IntDecodeValue is the ValueDecoderFunc for int types.
func (dvd DefaultValueDecoders) IntDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() ***REMOVED***
		return ValueDecoderError***REMOVED***
			Name:     "IntDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int***REMOVED***,
			Received: val,
		***REMOVED***
	***REMOVED***

	elem, err := dvd.intDecodeType(dc, vr, val.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.SetInt(elem.Int())
	return nil
***REMOVED***

// UintDecodeValue is the ValueDecoderFunc for uint types.
//
// Deprecated: UintDecodeValue is not registered by default. Use UintCodec.DecodeValue instead.
func (dvd DefaultValueDecoders) UintDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	var i64 int64
	var err error
	switch vr.Type() ***REMOVED***
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		i64 = int64(i32)
	case bsontype.Int64:
		i64, err = vr.ReadInt64()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	case bsontype.Double:
		f64, err := vr.ReadDouble()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if !dc.Truncate && math.Floor(f64) != f64 ***REMOVED***
			return errors.New("UintDecodeValue can only truncate float64 to an integer type when truncation is enabled")
		***REMOVED***
		if f64 > float64(math.MaxInt64) ***REMOVED***
			return fmt.Errorf("%g overflows int64", f64)
		***REMOVED***
		i64 = int64(f64)
	case bsontype.Boolean:
		b, err := vr.ReadBoolean()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if b ***REMOVED***
			i64 = 1
		***REMOVED***
	default:
		return fmt.Errorf("cannot decode %v into an integer type", vr.Type())
	***REMOVED***

	if !val.CanSet() ***REMOVED***
		return ValueDecoderError***REMOVED***
			Name:     "UintDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint***REMOVED***,
			Received: val,
		***REMOVED***
	***REMOVED***

	switch val.Kind() ***REMOVED***
	case reflect.Uint8:
		if i64 < 0 || i64 > math.MaxUint8 ***REMOVED***
			return fmt.Errorf("%d overflows uint8", i64)
		***REMOVED***
	case reflect.Uint16:
		if i64 < 0 || i64 > math.MaxUint16 ***REMOVED***
			return fmt.Errorf("%d overflows uint16", i64)
		***REMOVED***
	case reflect.Uint32:
		if i64 < 0 || i64 > math.MaxUint32 ***REMOVED***
			return fmt.Errorf("%d overflows uint32", i64)
		***REMOVED***
	case reflect.Uint64:
		if i64 < 0 ***REMOVED***
			return fmt.Errorf("%d overflows uint64", i64)
		***REMOVED***
	case reflect.Uint:
		if i64 < 0 || int64(uint(i64)) != i64 ***REMOVED*** // Can we fit this inside of an uint
			return fmt.Errorf("%d overflows uint", i64)
		***REMOVED***
	default:
		return ValueDecoderError***REMOVED***
			Name:     "UintDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint***REMOVED***,
			Received: val,
		***REMOVED***
	***REMOVED***

	val.SetUint(uint64(i64))
	return nil
***REMOVED***

func (dvd DefaultValueDecoders) floatDecodeType(ec DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	var f float64
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		f = float64(i32)
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		f = float64(i64)
	case bsontype.Double:
		f, err = vr.ReadDouble()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Boolean:
		b, err := vr.ReadBoolean()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if b ***REMOVED***
			f = 1
		***REMOVED***
	case bsontype.Null:
		if err = vr.ReadNull(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Undefined:
		if err = vr.ReadUndefined(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a float32 or float64 type", vrType)
	***REMOVED***

	switch t.Kind() ***REMOVED***
	case reflect.Float32:
		if !ec.Truncate && float64(float32(f)) != f ***REMOVED***
			return emptyValue, errCannotTruncate
		***REMOVED***

		return reflect.ValueOf(float32(f)), nil
	case reflect.Float64:
		return reflect.ValueOf(f), nil
	default:
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "FloatDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Float32, reflect.Float64***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***
***REMOVED***

// FloatDecodeValue is the ValueDecoderFunc for float types.
func (dvd DefaultValueDecoders) FloatDecodeValue(ec DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() ***REMOVED***
		return ValueDecoderError***REMOVED***
			Name:     "FloatDecodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.Float32, reflect.Float64***REMOVED***,
			Received: val,
		***REMOVED***
	***REMOVED***

	elem, err := dvd.floatDecodeType(ec, vr, val.Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.SetFloat(elem.Float())
	return nil
***REMOVED***

// StringDecodeValue is the ValueDecoderFunc for string types.
//
// Deprecated: StringDecodeValue is not registered by default. Use StringCodec.DecodeValue instead.
func (dvd DefaultValueDecoders) StringDecodeValue(dctx DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	var str string
	var err error
	switch vr.Type() ***REMOVED***
	// TODO(GODRIVER-577): Handle JavaScript and Symbol BSON types when allowed.
	case bsontype.String:
		str, err = vr.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	default:
		return fmt.Errorf("cannot decode %v into a string type", vr.Type())
	***REMOVED***
	if !val.CanSet() || val.Kind() != reflect.String ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "StringDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.String***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	val.SetString(str)
	return nil
***REMOVED***

func (DefaultValueDecoders) javaScriptDecodeType(dctx DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tJavaScript ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "JavaScriptDecodeValue",
			Types:    []reflect.Type***REMOVED***tJavaScript***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var js string
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.JavaScript:
		js, err = vr.ReadJavascript()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a primitive.JavaScript", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.JavaScript(js)), nil
***REMOVED***

// JavaScriptDecodeValue is the ValueDecoderFunc for the primitive.JavaScript type.
func (dvd DefaultValueDecoders) JavaScriptDecodeValue(dctx DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tJavaScript ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "JavaScriptDecodeValue", Types: []reflect.Type***REMOVED***tJavaScript***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.javaScriptDecodeType(dctx, vr, tJavaScript)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.SetString(elem.String())
	return nil
***REMOVED***

func (DefaultValueDecoders) symbolDecodeType(dctx DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tSymbol ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "SymbolDecodeValue",
			Types:    []reflect.Type***REMOVED***tSymbol***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var symbol string
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.String:
		symbol, err = vr.ReadString()
	case bsontype.Symbol:
		symbol, err = vr.ReadSymbol()
	case bsontype.Binary:
		data, subtype, err := vr.ReadBinary()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***

		if subtype != bsontype.BinaryGeneric && subtype != bsontype.BinaryBinaryOld ***REMOVED***
			return emptyValue, decodeBinaryError***REMOVED***subtype: subtype, typeName: "primitive.Symbol"***REMOVED***
		***REMOVED***
		symbol = string(data)
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a primitive.Symbol", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.Symbol(symbol)), nil
***REMOVED***

// SymbolDecodeValue is the ValueDecoderFunc for the primitive.Symbol type.
func (dvd DefaultValueDecoders) SymbolDecodeValue(dctx DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tSymbol ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "SymbolDecodeValue", Types: []reflect.Type***REMOVED***tSymbol***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.symbolDecodeType(dctx, vr, tSymbol)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.SetString(elem.String())
	return nil
***REMOVED***

func (DefaultValueDecoders) binaryDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tBinary ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "BinaryDecodeValue",
			Types:    []reflect.Type***REMOVED***tBinary***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Binary:
		data, subtype, err = vr.ReadBinary()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a Binary", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.Binary***REMOVED***Subtype: subtype, Data: data***REMOVED***), nil
***REMOVED***

// BinaryDecodeValue is the ValueDecoderFunc for Binary.
func (dvd DefaultValueDecoders) BinaryDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tBinary ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "BinaryDecodeValue", Types: []reflect.Type***REMOVED***tBinary***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.binaryDecodeType(dc, vr, tBinary)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (DefaultValueDecoders) undefinedDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tUndefined ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "UndefinedDecodeValue",
			Types:    []reflect.Type***REMOVED***tUndefined***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	case bsontype.Null:
		err = vr.ReadNull()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into an Undefined", vr.Type())
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.Undefined***REMOVED******REMOVED***), nil
***REMOVED***

// UndefinedDecodeValue is the ValueDecoderFunc for Undefined.
func (dvd DefaultValueDecoders) UndefinedDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tUndefined ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "UndefinedDecodeValue", Types: []reflect.Type***REMOVED***tUndefined***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.undefinedDecodeType(dc, vr, tUndefined)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

// Accept both 12-byte string and pretty-printed 24-byte hex string formats.
func (dvd DefaultValueDecoders) objectIDDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tOID ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "ObjectIDDecodeValue",
			Types:    []reflect.Type***REMOVED***tOID***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var oid primitive.ObjectID
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.ObjectID:
		oid, err = vr.ReadObjectID()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		if oid, err = primitive.ObjectIDFromHex(str); err == nil ***REMOVED***
			break
		***REMOVED***
		if len(str) != 12 ***REMOVED***
			return emptyValue, fmt.Errorf("an ObjectID string must be exactly 12 bytes long (got %v)", len(str))
		***REMOVED***
		byteArr := []byte(str)
		copy(oid[:], byteArr)
	case bsontype.Null:
		if err = vr.ReadNull(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	case bsontype.Undefined:
		if err = vr.ReadUndefined(); err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into an ObjectID", vrType)
	***REMOVED***

	return reflect.ValueOf(oid), nil
***REMOVED***

// ObjectIDDecodeValue is the ValueDecoderFunc for primitive.ObjectID.
func (dvd DefaultValueDecoders) ObjectIDDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tOID ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "ObjectIDDecodeValue", Types: []reflect.Type***REMOVED***tOID***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.objectIDDecodeType(dc, vr, tOID)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (DefaultValueDecoders) dateTimeDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tDateTime ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "DateTimeDecodeValue",
			Types:    []reflect.Type***REMOVED***tDateTime***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var dt int64
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.DateTime:
		dt, err = vr.ReadDateTime()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a DateTime", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.DateTime(dt)), nil
***REMOVED***

// DateTimeDecodeValue is the ValueDecoderFunc for DateTime.
func (dvd DefaultValueDecoders) DateTimeDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tDateTime ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "DateTimeDecodeValue", Types: []reflect.Type***REMOVED***tDateTime***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.dateTimeDecodeType(dc, vr, tDateTime)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (DefaultValueDecoders) nullDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tNull ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "NullDecodeValue",
			Types:    []reflect.Type***REMOVED***tNull***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	case bsontype.Null:
		err = vr.ReadNull()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a Null", vr.Type())
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.Null***REMOVED******REMOVED***), nil
***REMOVED***

// NullDecodeValue is the ValueDecoderFunc for Null.
func (dvd DefaultValueDecoders) NullDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tNull ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "NullDecodeValue", Types: []reflect.Type***REMOVED***tNull***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.nullDecodeType(dc, vr, tNull)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (DefaultValueDecoders) regexDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tRegex ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "RegexDecodeValue",
			Types:    []reflect.Type***REMOVED***tRegex***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var pattern, options string
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Regex:
		pattern, options, err = vr.ReadRegex()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a Regex", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.Regex***REMOVED***Pattern: pattern, Options: options***REMOVED***), nil
***REMOVED***

// RegexDecodeValue is the ValueDecoderFunc for Regex.
func (dvd DefaultValueDecoders) RegexDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tRegex ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "RegexDecodeValue", Types: []reflect.Type***REMOVED***tRegex***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.regexDecodeType(dc, vr, tRegex)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (DefaultValueDecoders) dBPointerDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tDBPointer ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "DBPointerDecodeValue",
			Types:    []reflect.Type***REMOVED***tDBPointer***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var ns string
	var pointer primitive.ObjectID
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.DBPointer:
		ns, pointer, err = vr.ReadDBPointer()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a DBPointer", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.DBPointer***REMOVED***DB: ns, Pointer: pointer***REMOVED***), nil
***REMOVED***

// DBPointerDecodeValue is the ValueDecoderFunc for DBPointer.
func (dvd DefaultValueDecoders) DBPointerDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tDBPointer ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "DBPointerDecodeValue", Types: []reflect.Type***REMOVED***tDBPointer***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.dBPointerDecodeType(dc, vr, tDBPointer)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (DefaultValueDecoders) timestampDecodeType(dc DecodeContext, vr bsonrw.ValueReader, reflectType reflect.Type) (reflect.Value, error) ***REMOVED***
	if reflectType != tTimestamp ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "TimestampDecodeValue",
			Types:    []reflect.Type***REMOVED***tTimestamp***REMOVED***,
			Received: reflect.Zero(reflectType),
		***REMOVED***
	***REMOVED***

	var t, incr uint32
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Timestamp:
		t, incr, err = vr.ReadTimestamp()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a Timestamp", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.Timestamp***REMOVED***T: t, I: incr***REMOVED***), nil
***REMOVED***

// TimestampDecodeValue is the ValueDecoderFunc for Timestamp.
func (dvd DefaultValueDecoders) TimestampDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tTimestamp ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "TimestampDecodeValue", Types: []reflect.Type***REMOVED***tTimestamp***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.timestampDecodeType(dc, vr, tTimestamp)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (DefaultValueDecoders) minKeyDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tMinKey ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "MinKeyDecodeValue",
			Types:    []reflect.Type***REMOVED***tMinKey***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.MinKey:
		err = vr.ReadMinKey()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a MinKey", vr.Type())
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.MinKey***REMOVED******REMOVED***), nil
***REMOVED***

// MinKeyDecodeValue is the ValueDecoderFunc for MinKey.
func (dvd DefaultValueDecoders) MinKeyDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tMinKey ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "MinKeyDecodeValue", Types: []reflect.Type***REMOVED***tMinKey***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.minKeyDecodeType(dc, vr, tMinKey)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (DefaultValueDecoders) maxKeyDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tMaxKey ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "MaxKeyDecodeValue",
			Types:    []reflect.Type***REMOVED***tMaxKey***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.MaxKey:
		err = vr.ReadMaxKey()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a MaxKey", vr.Type())
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(primitive.MaxKey***REMOVED******REMOVED***), nil
***REMOVED***

// MaxKeyDecodeValue is the ValueDecoderFunc for MaxKey.
func (dvd DefaultValueDecoders) MaxKeyDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tMaxKey ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "MaxKeyDecodeValue", Types: []reflect.Type***REMOVED***tMaxKey***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.maxKeyDecodeType(dc, vr, tMaxKey)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (dvd DefaultValueDecoders) decimal128DecodeType(dctx DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tDecimal ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "Decimal128DecodeValue",
			Types:    []reflect.Type***REMOVED***tDecimal***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var d128 primitive.Decimal128
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Decimal128:
		d128, err = vr.ReadDecimal128()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a primitive.Decimal128", vr.Type())
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(d128), nil
***REMOVED***

// Decimal128DecodeValue is the ValueDecoderFunc for primitive.Decimal128.
func (dvd DefaultValueDecoders) Decimal128DecodeValue(dctx DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tDecimal ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "Decimal128DecodeValue", Types: []reflect.Type***REMOVED***tDecimal***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.decimal128DecodeType(dctx, vr, tDecimal)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (dvd DefaultValueDecoders) jsonNumberDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tJSONNumber ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "JSONNumberDecodeValue",
			Types:    []reflect.Type***REMOVED***tJSONNumber***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var jsonNum json.Number
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Double:
		f64, err := vr.ReadDouble()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		jsonNum = json.Number(strconv.FormatFloat(f64, 'f', -1, 64))
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		jsonNum = json.Number(strconv.FormatInt(int64(i32), 10))
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***
		jsonNum = json.Number(strconv.FormatInt(i64, 10))
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a json.Number", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(jsonNum), nil
***REMOVED***

// JSONNumberDecodeValue is the ValueDecoderFunc for json.Number.
func (dvd DefaultValueDecoders) JSONNumberDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tJSONNumber ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "JSONNumberDecodeValue", Types: []reflect.Type***REMOVED***tJSONNumber***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.jsonNumberDecodeType(dc, vr, tJSONNumber)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (dvd DefaultValueDecoders) urlDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tURL ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "URLDecodeValue",
			Types:    []reflect.Type***REMOVED***tURL***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	urlPtr := &url.URL***REMOVED******REMOVED***
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.String:
		var str string // Declare str here to avoid shadowing err during the ReadString call.
		str, err = vr.ReadString()
		if err != nil ***REMOVED***
			return emptyValue, err
		***REMOVED***

		urlPtr, err = url.Parse(str)
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a *url.URL", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(urlPtr).Elem(), nil
***REMOVED***

// URLDecodeValue is the ValueDecoderFunc for url.URL.
func (dvd DefaultValueDecoders) URLDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tURL ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "URLDecodeValue", Types: []reflect.Type***REMOVED***tURL***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.urlDecodeType(dc, vr, tURL)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

// TimeDecodeValue is the ValueDecoderFunc for time.Time.
//
// Deprecated: TimeDecodeValue is not registered by default. Use TimeCodec.DecodeValue instead.
func (dvd DefaultValueDecoders) TimeDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if vr.Type() != bsontype.DateTime ***REMOVED***
		return fmt.Errorf("cannot decode %v into a time.Time", vr.Type())
	***REMOVED***

	dt, err := vr.ReadDateTime()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if !val.CanSet() || val.Type() != tTime ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "TimeDecodeValue", Types: []reflect.Type***REMOVED***tTime***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	val.Set(reflect.ValueOf(time.Unix(dt/1000, dt%1000*1000000).UTC()))
	return nil
***REMOVED***

// ByteSliceDecodeValue is the ValueDecoderFunc for []byte.
//
// Deprecated: ByteSliceDecodeValue is not registered by default. Use ByteSliceCodec.DecodeValue instead.
func (dvd DefaultValueDecoders) ByteSliceDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if vr.Type() != bsontype.Binary && vr.Type() != bsontype.Null ***REMOVED***
		return fmt.Errorf("cannot decode %v into a []byte", vr.Type())
	***REMOVED***

	if !val.CanSet() || val.Type() != tByteSlice ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "ByteSliceDecodeValue", Types: []reflect.Type***REMOVED***tByteSlice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if vr.Type() == bsontype.Null ***REMOVED***
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	***REMOVED***

	data, subtype, err := vr.ReadBinary()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if subtype != 0x00 ***REMOVED***
		return fmt.Errorf("ByteSliceDecodeValue can only be used to decode subtype 0x00 for %s, got %v", bsontype.Binary, subtype)
	***REMOVED***

	val.Set(reflect.ValueOf(data))
	return nil
***REMOVED***

// MapDecodeValue is the ValueDecoderFunc for map[string]* types.
//
// Deprecated: MapDecodeValue is not registered by default. Use MapCodec.DecodeValue instead.
func (dvd DefaultValueDecoders) MapDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Kind() != reflect.Map || val.Type().Key().Kind() != reflect.String ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "MapDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Map***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	switch vr.Type() ***REMOVED***
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	default:
		return fmt.Errorf("cannot decode %v into a %s", vr.Type(), val.Type())
	***REMOVED***

	dr, err := vr.ReadDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if val.IsNil() ***REMOVED***
		val.Set(reflect.MakeMap(val.Type()))
	***REMOVED***

	eType := val.Type().Elem()
	decoder, err := dc.LookupDecoder(eType)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if eType == tEmpty ***REMOVED***
		dc.Ancestor = val.Type()
	***REMOVED***

	keyType := val.Type().Key()
	for ***REMOVED***
		key, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		elem := reflect.New(eType).Elem()

		err = decoder.DecodeValue(dc, vr, elem)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		val.SetMapIndex(reflect.ValueOf(key).Convert(keyType), elem)
	***REMOVED***
	return nil
***REMOVED***

// ArrayDecodeValue is the ValueDecoderFunc for array types.
func (dvd DefaultValueDecoders) ArrayDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Kind() != reflect.Array ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "ArrayDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Array***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Array:
	case bsontype.Type(0), bsontype.EmbeddedDocument:
		if val.Type().Elem() != tE ***REMOVED***
			return fmt.Errorf("cannot decode document into %s", val.Type())
		***REMOVED***
	case bsontype.Binary:
		if val.Type().Elem() != tByte ***REMOVED***
			return fmt.Errorf("ArrayDecodeValue can only be used to decode binary into a byte array, got %v", vrType)
		***REMOVED***
		data, subtype, err := vr.ReadBinary()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if subtype != bsontype.BinaryGeneric && subtype != bsontype.BinaryBinaryOld ***REMOVED***
			return fmt.Errorf("ArrayDecodeValue can only be used to decode subtype 0x00 or 0x02 for %s, got %v", bsontype.Binary, subtype)
		***REMOVED***

		if len(data) > val.Len() ***REMOVED***
			return fmt.Errorf("more elements returned in array than can fit inside %s", val.Type())
		***REMOVED***

		for idx, elem := range data ***REMOVED***
			val.Index(idx).Set(reflect.ValueOf(elem))
		***REMOVED***
		return nil
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	case bsontype.Undefined:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into an array", vrType)
	***REMOVED***

	var elemsFunc func(DecodeContext, bsonrw.ValueReader, reflect.Value) ([]reflect.Value, error)
	switch val.Type().Elem() ***REMOVED***
	case tE:
		elemsFunc = dvd.decodeD
	default:
		elemsFunc = dvd.decodeDefault
	***REMOVED***

	elems, err := elemsFunc(dc, vr, val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if len(elems) > val.Len() ***REMOVED***
		return fmt.Errorf("more elements returned in array than can fit inside %s, got %v elements", val.Type(), len(elems))
	***REMOVED***

	for idx, elem := range elems ***REMOVED***
		val.Index(idx).Set(elem)
	***REMOVED***

	return nil
***REMOVED***

// SliceDecodeValue is the ValueDecoderFunc for slice types.
//
// Deprecated: SliceDecodeValue is not registered by default. Use SliceCodec.DecodeValue instead.
func (dvd DefaultValueDecoders) SliceDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Kind() != reflect.Slice ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "SliceDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Slice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	switch vr.Type() ***REMOVED***
	case bsontype.Array:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	case bsontype.Type(0), bsontype.EmbeddedDocument:
		if val.Type().Elem() != tE ***REMOVED***
			return fmt.Errorf("cannot decode document into %s", val.Type())
		***REMOVED***
	default:
		return fmt.Errorf("cannot decode %v into a slice", vr.Type())
	***REMOVED***

	var elemsFunc func(DecodeContext, bsonrw.ValueReader, reflect.Value) ([]reflect.Value, error)
	switch val.Type().Elem() ***REMOVED***
	case tE:
		dc.Ancestor = val.Type()
		elemsFunc = dvd.decodeD
	default:
		elemsFunc = dvd.decodeDefault
	***REMOVED***

	elems, err := elemsFunc(dc, vr, val)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if val.IsNil() ***REMOVED***
		val.Set(reflect.MakeSlice(val.Type(), 0, len(elems)))
	***REMOVED***

	val.SetLen(0)
	val.Set(reflect.Append(val, elems...))

	return nil
***REMOVED***

// ValueUnmarshalerDecodeValue is the ValueDecoderFunc for ValueUnmarshaler implementations.
func (dvd DefaultValueDecoders) ValueUnmarshalerDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || (!val.Type().Implements(tValueUnmarshaler) && !reflect.PtrTo(val.Type()).Implements(tValueUnmarshaler)) ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "ValueUnmarshalerDecodeValue", Types: []reflect.Type***REMOVED***tValueUnmarshaler***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.Kind() == reflect.Ptr && val.IsNil() ***REMOVED***
		if !val.CanSet() ***REMOVED***
			return ValueDecoderError***REMOVED***Name: "ValueUnmarshalerDecodeValue", Types: []reflect.Type***REMOVED***tValueUnmarshaler***REMOVED***, Received: val***REMOVED***
		***REMOVED***
		val.Set(reflect.New(val.Type().Elem()))
	***REMOVED***

	if !val.Type().Implements(tValueUnmarshaler) ***REMOVED***
		if !val.CanAddr() ***REMOVED***
			return ValueDecoderError***REMOVED***Name: "ValueUnmarshalerDecodeValue", Types: []reflect.Type***REMOVED***tValueUnmarshaler***REMOVED***, Received: val***REMOVED***
		***REMOVED***
		val = val.Addr() // If the type doesn't implement the interface, a pointer to it must.
	***REMOVED***

	t, src, err := bsonrw.Copier***REMOVED******REMOVED***.CopyValueToBytes(vr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	fn := val.Convert(tValueUnmarshaler).MethodByName("UnmarshalBSONValue")
	errVal := fn.Call([]reflect.Value***REMOVED***reflect.ValueOf(t), reflect.ValueOf(src)***REMOVED***)[0]
	if !errVal.IsNil() ***REMOVED***
		return errVal.Interface().(error)
	***REMOVED***
	return nil
***REMOVED***

// UnmarshalerDecodeValue is the ValueDecoderFunc for Unmarshaler implementations.
func (dvd DefaultValueDecoders) UnmarshalerDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || (!val.Type().Implements(tUnmarshaler) && !reflect.PtrTo(val.Type()).Implements(tUnmarshaler)) ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "UnmarshalerDecodeValue", Types: []reflect.Type***REMOVED***tUnmarshaler***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.Kind() == reflect.Ptr && val.IsNil() ***REMOVED***
		if !val.CanSet() ***REMOVED***
			return ValueDecoderError***REMOVED***Name: "UnmarshalerDecodeValue", Types: []reflect.Type***REMOVED***tUnmarshaler***REMOVED***, Received: val***REMOVED***
		***REMOVED***
		val.Set(reflect.New(val.Type().Elem()))
	***REMOVED***

	_, src, err := bsonrw.Copier***REMOVED******REMOVED***.CopyValueToBytes(vr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// If the target Go value is a pointer and the BSON field value is empty, set the value to the
	// zero value of the pointer (nil) and don't call UnmarshalBSON. UnmarshalBSON has no way to
	// change the pointer value from within the function (only the value at the pointer address),
	// so it can't set the pointer to "nil" itself. Since the most common Go value for an empty BSON
	// field value is "nil", we set "nil" here and don't call UnmarshalBSON. This behavior matches
	// the behavior of the Go "encoding/json" unmarshaler when the target Go value is a pointer and
	// the JSON field value is "null".
	if val.Kind() == reflect.Ptr && len(src) == 0 ***REMOVED***
		val.Set(reflect.Zero(val.Type()))
		return nil
	***REMOVED***

	if !val.Type().Implements(tUnmarshaler) ***REMOVED***
		if !val.CanAddr() ***REMOVED***
			return ValueDecoderError***REMOVED***Name: "UnmarshalerDecodeValue", Types: []reflect.Type***REMOVED***tUnmarshaler***REMOVED***, Received: val***REMOVED***
		***REMOVED***
		val = val.Addr() // If the type doesn't implement the interface, a pointer to it must.
	***REMOVED***

	fn := val.Convert(tUnmarshaler).MethodByName("UnmarshalBSON")
	errVal := fn.Call([]reflect.Value***REMOVED***reflect.ValueOf(src)***REMOVED***)[0]
	if !errVal.IsNil() ***REMOVED***
		return errVal.Interface().(error)
	***REMOVED***
	return nil
***REMOVED***

// EmptyInterfaceDecodeValue is the ValueDecoderFunc for interface***REMOVED******REMOVED***.
//
// Deprecated: EmptyInterfaceDecodeValue is not registered by default. Use EmptyInterfaceCodec.DecodeValue instead.
func (dvd DefaultValueDecoders) EmptyInterfaceDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tEmpty ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "EmptyInterfaceDecodeValue", Types: []reflect.Type***REMOVED***tEmpty***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	rtype, err := dc.LookupTypeMapEntry(vr.Type())
	if err != nil ***REMOVED***
		switch vr.Type() ***REMOVED***
		case bsontype.EmbeddedDocument:
			if dc.Ancestor != nil ***REMOVED***
				rtype = dc.Ancestor
				break
			***REMOVED***
			rtype = tD
		case bsontype.Null:
			val.Set(reflect.Zero(val.Type()))
			return vr.ReadNull()
		default:
			return err
		***REMOVED***
	***REMOVED***

	decoder, err := dc.LookupDecoder(rtype)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	elem := reflect.New(rtype).Elem()
	err = decoder.DecodeValue(dc, vr, elem)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

// CoreDocumentDecodeValue is the ValueDecoderFunc for bsoncore.Document.
func (DefaultValueDecoders) CoreDocumentDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tCoreDocument ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "CoreDocumentDecodeValue", Types: []reflect.Type***REMOVED***tCoreDocument***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		val.Set(reflect.MakeSlice(val.Type(), 0, 0))
	***REMOVED***

	val.SetLen(0)

	cdoc, err := bsonrw.Copier***REMOVED******REMOVED***.AppendDocumentBytes(val.Interface().(bsoncore.Document), vr)
	val.Set(reflect.ValueOf(cdoc))
	return err
***REMOVED***

func (dvd DefaultValueDecoders) decodeDefault(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) ([]reflect.Value, error) ***REMOVED***
	elems := make([]reflect.Value, 0)

	ar, err := vr.ReadArray()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	eType := val.Type().Elem()

	decoder, err := dc.LookupDecoder(eType)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	eTypeDecoder, _ := decoder.(typeDecoder)

	idx := 0
	for ***REMOVED***
		vr, err := ar.ReadValue()
		if err == bsonrw.ErrEOA ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		elem, err := decodeTypeOrValueWithInfo(decoder, eTypeDecoder, dc, vr, eType, true)
		if err != nil ***REMOVED***
			return nil, newDecodeError(strconv.Itoa(idx), err)
		***REMOVED***
		elems = append(elems, elem)
		idx++
	***REMOVED***

	return elems, nil
***REMOVED***

func (dvd DefaultValueDecoders) readCodeWithScope(dc DecodeContext, vr bsonrw.ValueReader) (primitive.CodeWithScope, error) ***REMOVED***
	var cws primitive.CodeWithScope

	code, dr, err := vr.ReadCodeWithScope()
	if err != nil ***REMOVED***
		return cws, err
	***REMOVED***

	scope := reflect.New(tD).Elem()
	elems, err := dvd.decodeElemsFromDocumentReader(dc, dr)
	if err != nil ***REMOVED***
		return cws, err
	***REMOVED***

	scope.Set(reflect.MakeSlice(tD, 0, len(elems)))
	scope.Set(reflect.Append(scope, elems...))

	cws = primitive.CodeWithScope***REMOVED***
		Code:  primitive.JavaScript(code),
		Scope: scope.Interface().(primitive.D),
	***REMOVED***
	return cws, nil
***REMOVED***

func (dvd DefaultValueDecoders) codeWithScopeDecodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) ***REMOVED***
	if t != tCodeWithScope ***REMOVED***
		return emptyValue, ValueDecoderError***REMOVED***
			Name:     "CodeWithScopeDecodeValue",
			Types:    []reflect.Type***REMOVED***tCodeWithScope***REMOVED***,
			Received: reflect.Zero(t),
		***REMOVED***
	***REMOVED***

	var cws primitive.CodeWithScope
	var err error
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.CodeWithScope:
		cws, err = dvd.readCodeWithScope(dc, vr)
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a primitive.CodeWithScope", vrType)
	***REMOVED***
	if err != nil ***REMOVED***
		return emptyValue, err
	***REMOVED***

	return reflect.ValueOf(cws), nil
***REMOVED***

// CodeWithScopeDecodeValue is the ValueDecoderFunc for CodeWithScope.
func (dvd DefaultValueDecoders) CodeWithScopeDecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tCodeWithScope ***REMOVED***
		return ValueDecoderError***REMOVED***Name: "CodeWithScopeDecodeValue", Types: []reflect.Type***REMOVED***tCodeWithScope***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	elem, err := dvd.codeWithScopeDecodeType(dc, vr, tCodeWithScope)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(elem)
	return nil
***REMOVED***

func (dvd DefaultValueDecoders) decodeD(dc DecodeContext, vr bsonrw.ValueReader, _ reflect.Value) ([]reflect.Value, error) ***REMOVED***
	switch vr.Type() ***REMOVED***
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	default:
		return nil, fmt.Errorf("cannot decode %v into a D", vr.Type())
	***REMOVED***

	dr, err := vr.ReadDocument()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return dvd.decodeElemsFromDocumentReader(dc, dr)
***REMOVED***

func (DefaultValueDecoders) decodeElemsFromDocumentReader(dc DecodeContext, dr bsonrw.DocumentReader) ([]reflect.Value, error) ***REMOVED***
	decoder, err := dc.LookupDecoder(tEmpty)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	elems := make([]reflect.Value, 0)
	for ***REMOVED***
		key, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		val := reflect.New(tEmpty).Elem()
		err = decoder.DecodeValue(dc, vr, val)
		if err != nil ***REMOVED***
			return nil, newDecodeError(key, err)
		***REMOVED***

		elems = append(elems, reflect.ValueOf(primitive.E***REMOVED***Key: key, Value: val.Interface()***REMOVED***))
	***REMOVED***

	return elems, nil
***REMOVED***
