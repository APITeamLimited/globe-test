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
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var defaultValueEncoders DefaultValueEncoders

var bvwPool = bsonrw.NewBSONValueWriterPool()

var errInvalidValue = errors.New("cannot encode invalid element")

var sliceWriterPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		sw := make(bsonrw.SliceWriter, 0)
		return &sw
	***REMOVED***,
***REMOVED***

func encodeElement(ec EncodeContext, dw bsonrw.DocumentWriter, e primitive.E) error ***REMOVED***
	vw, err := dw.WriteDocumentElement(e.Key)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if e.Value == nil ***REMOVED***
		return vw.WriteNull()
	***REMOVED***
	encoder, err := ec.LookupEncoder(reflect.TypeOf(e.Value))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = encoder.EncodeValue(ec, vw, reflect.ValueOf(e.Value))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// DefaultValueEncoders is a namespace type for the default ValueEncoders used
// when creating a registry.
type DefaultValueEncoders struct***REMOVED******REMOVED***

// RegisterDefaultEncoders will register the encoder methods attached to DefaultValueEncoders with
// the provided RegistryBuilder.
func (dve DefaultValueEncoders) RegisterDefaultEncoders(rb *RegistryBuilder) ***REMOVED***
	if rb == nil ***REMOVED***
		panic(errors.New("argument to RegisterDefaultEncoders must not be nil"))
	***REMOVED***
	rb.
		RegisterTypeEncoder(tByteSlice, defaultByteSliceCodec).
		RegisterTypeEncoder(tTime, defaultTimeCodec).
		RegisterTypeEncoder(tEmpty, defaultEmptyInterfaceCodec).
		RegisterTypeEncoder(tCoreArray, defaultArrayCodec).
		RegisterTypeEncoder(tOID, ValueEncoderFunc(dve.ObjectIDEncodeValue)).
		RegisterTypeEncoder(tDecimal, ValueEncoderFunc(dve.Decimal128EncodeValue)).
		RegisterTypeEncoder(tJSONNumber, ValueEncoderFunc(dve.JSONNumberEncodeValue)).
		RegisterTypeEncoder(tURL, ValueEncoderFunc(dve.URLEncodeValue)).
		RegisterTypeEncoder(tJavaScript, ValueEncoderFunc(dve.JavaScriptEncodeValue)).
		RegisterTypeEncoder(tSymbol, ValueEncoderFunc(dve.SymbolEncodeValue)).
		RegisterTypeEncoder(tBinary, ValueEncoderFunc(dve.BinaryEncodeValue)).
		RegisterTypeEncoder(tUndefined, ValueEncoderFunc(dve.UndefinedEncodeValue)).
		RegisterTypeEncoder(tDateTime, ValueEncoderFunc(dve.DateTimeEncodeValue)).
		RegisterTypeEncoder(tNull, ValueEncoderFunc(dve.NullEncodeValue)).
		RegisterTypeEncoder(tRegex, ValueEncoderFunc(dve.RegexEncodeValue)).
		RegisterTypeEncoder(tDBPointer, ValueEncoderFunc(dve.DBPointerEncodeValue)).
		RegisterTypeEncoder(tTimestamp, ValueEncoderFunc(dve.TimestampEncodeValue)).
		RegisterTypeEncoder(tMinKey, ValueEncoderFunc(dve.MinKeyEncodeValue)).
		RegisterTypeEncoder(tMaxKey, ValueEncoderFunc(dve.MaxKeyEncodeValue)).
		RegisterTypeEncoder(tCoreDocument, ValueEncoderFunc(dve.CoreDocumentEncodeValue)).
		RegisterTypeEncoder(tCodeWithScope, ValueEncoderFunc(dve.CodeWithScopeEncodeValue)).
		RegisterDefaultEncoder(reflect.Bool, ValueEncoderFunc(dve.BooleanEncodeValue)).
		RegisterDefaultEncoder(reflect.Int, ValueEncoderFunc(dve.IntEncodeValue)).
		RegisterDefaultEncoder(reflect.Int8, ValueEncoderFunc(dve.IntEncodeValue)).
		RegisterDefaultEncoder(reflect.Int16, ValueEncoderFunc(dve.IntEncodeValue)).
		RegisterDefaultEncoder(reflect.Int32, ValueEncoderFunc(dve.IntEncodeValue)).
		RegisterDefaultEncoder(reflect.Int64, ValueEncoderFunc(dve.IntEncodeValue)).
		RegisterDefaultEncoder(reflect.Uint, defaultUIntCodec).
		RegisterDefaultEncoder(reflect.Uint8, defaultUIntCodec).
		RegisterDefaultEncoder(reflect.Uint16, defaultUIntCodec).
		RegisterDefaultEncoder(reflect.Uint32, defaultUIntCodec).
		RegisterDefaultEncoder(reflect.Uint64, defaultUIntCodec).
		RegisterDefaultEncoder(reflect.Float32, ValueEncoderFunc(dve.FloatEncodeValue)).
		RegisterDefaultEncoder(reflect.Float64, ValueEncoderFunc(dve.FloatEncodeValue)).
		RegisterDefaultEncoder(reflect.Array, ValueEncoderFunc(dve.ArrayEncodeValue)).
		RegisterDefaultEncoder(reflect.Map, defaultMapCodec).
		RegisterDefaultEncoder(reflect.Slice, defaultSliceCodec).
		RegisterDefaultEncoder(reflect.String, defaultStringCodec).
		RegisterDefaultEncoder(reflect.Struct, newDefaultStructCodec()).
		RegisterDefaultEncoder(reflect.Ptr, NewPointerCodec()).
		RegisterHookEncoder(tValueMarshaler, ValueEncoderFunc(dve.ValueMarshalerEncodeValue)).
		RegisterHookEncoder(tMarshaler, ValueEncoderFunc(dve.MarshalerEncodeValue)).
		RegisterHookEncoder(tProxy, ValueEncoderFunc(dve.ProxyEncodeValue))
***REMOVED***

// BooleanEncodeValue is the ValueEncoderFunc for bool types.
func (dve DefaultValueEncoders) BooleanEncodeValue(ectx EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Kind() != reflect.Bool ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "BooleanEncodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Bool***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	return vw.WriteBoolean(val.Bool())
***REMOVED***

func fitsIn32Bits(i int64) bool ***REMOVED***
	return math.MinInt32 <= i && i <= math.MaxInt32
***REMOVED***

// IntEncodeValue is the ValueEncoderFunc for int types.
func (dve DefaultValueEncoders) IntEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	switch val.Kind() ***REMOVED***
	case reflect.Int8, reflect.Int16, reflect.Int32:
		return vw.WriteInt32(int32(val.Int()))
	case reflect.Int:
		i64 := val.Int()
		if fitsIn32Bits(i64) ***REMOVED***
			return vw.WriteInt32(int32(i64))
		***REMOVED***
		return vw.WriteInt64(i64)
	case reflect.Int64:
		i64 := val.Int()
		if ec.MinSize && fitsIn32Bits(i64) ***REMOVED***
			return vw.WriteInt32(int32(i64))
		***REMOVED***
		return vw.WriteInt64(i64)
	***REMOVED***

	return ValueEncoderError***REMOVED***
		Name:     "IntEncodeValue",
		Kinds:    []reflect.Kind***REMOVED***reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int***REMOVED***,
		Received: val,
	***REMOVED***
***REMOVED***

// UintEncodeValue is the ValueEncoderFunc for uint types.
//
// Deprecated: UintEncodeValue is not registered by default. Use UintCodec.EncodeValue instead.
func (dve DefaultValueEncoders) UintEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	switch val.Kind() ***REMOVED***
	case reflect.Uint8, reflect.Uint16:
		return vw.WriteInt32(int32(val.Uint()))
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		u64 := val.Uint()
		if ec.MinSize && u64 <= math.MaxInt32 ***REMOVED***
			return vw.WriteInt32(int32(u64))
		***REMOVED***
		if u64 > math.MaxInt64 ***REMOVED***
			return fmt.Errorf("%d overflows int64", u64)
		***REMOVED***
		return vw.WriteInt64(int64(u64))
	***REMOVED***

	return ValueEncoderError***REMOVED***
		Name:     "UintEncodeValue",
		Kinds:    []reflect.Kind***REMOVED***reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint***REMOVED***,
		Received: val,
	***REMOVED***
***REMOVED***

// FloatEncodeValue is the ValueEncoderFunc for float types.
func (dve DefaultValueEncoders) FloatEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	switch val.Kind() ***REMOVED***
	case reflect.Float32, reflect.Float64:
		return vw.WriteDouble(val.Float())
	***REMOVED***

	return ValueEncoderError***REMOVED***Name: "FloatEncodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Float32, reflect.Float64***REMOVED***, Received: val***REMOVED***
***REMOVED***

// StringEncodeValue is the ValueEncoderFunc for string types.
//
// Deprecated: StringEncodeValue is not registered by default. Use StringCodec.EncodeValue instead.
func (dve DefaultValueEncoders) StringEncodeValue(ectx EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if val.Kind() != reflect.String ***REMOVED***
		return ValueEncoderError***REMOVED***
			Name:     "StringEncodeValue",
			Kinds:    []reflect.Kind***REMOVED***reflect.String***REMOVED***,
			Received: val,
		***REMOVED***
	***REMOVED***

	return vw.WriteString(val.String())
***REMOVED***

// ObjectIDEncodeValue is the ValueEncoderFunc for primitive.ObjectID.
func (dve DefaultValueEncoders) ObjectIDEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tOID ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "ObjectIDEncodeValue", Types: []reflect.Type***REMOVED***tOID***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	return vw.WriteObjectID(val.Interface().(primitive.ObjectID))
***REMOVED***

// Decimal128EncodeValue is the ValueEncoderFunc for primitive.Decimal128.
func (dve DefaultValueEncoders) Decimal128EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tDecimal ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "Decimal128EncodeValue", Types: []reflect.Type***REMOVED***tDecimal***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	return vw.WriteDecimal128(val.Interface().(primitive.Decimal128))
***REMOVED***

// JSONNumberEncodeValue is the ValueEncoderFunc for json.Number.
func (dve DefaultValueEncoders) JSONNumberEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tJSONNumber ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "JSONNumberEncodeValue", Types: []reflect.Type***REMOVED***tJSONNumber***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	jsnum := val.Interface().(json.Number)

	// Attempt int first, then float64
	if i64, err := jsnum.Int64(); err == nil ***REMOVED***
		return dve.IntEncodeValue(ec, vw, reflect.ValueOf(i64))
	***REMOVED***

	f64, err := jsnum.Float64()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return dve.FloatEncodeValue(ec, vw, reflect.ValueOf(f64))
***REMOVED***

// URLEncodeValue is the ValueEncoderFunc for url.URL.
func (dve DefaultValueEncoders) URLEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tURL ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "URLEncodeValue", Types: []reflect.Type***REMOVED***tURL***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	u := val.Interface().(url.URL)
	return vw.WriteString(u.String())
***REMOVED***

// TimeEncodeValue is the ValueEncoderFunc for time.TIme.
//
// Deprecated: TimeEncodeValue is not registered by default. Use TimeCodec.EncodeValue instead.
func (dve DefaultValueEncoders) TimeEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tTime ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "TimeEncodeValue", Types: []reflect.Type***REMOVED***tTime***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	tt := val.Interface().(time.Time)
	dt := primitive.NewDateTimeFromTime(tt)
	return vw.WriteDateTime(int64(dt))
***REMOVED***

// ByteSliceEncodeValue is the ValueEncoderFunc for []byte.
//
// Deprecated: ByteSliceEncodeValue is not registered by default. Use ByteSliceCodec.EncodeValue instead.
func (dve DefaultValueEncoders) ByteSliceEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tByteSlice ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "ByteSliceEncodeValue", Types: []reflect.Type***REMOVED***tByteSlice***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***
	return vw.WriteBinary(val.Interface().([]byte))
***REMOVED***

// MapEncodeValue is the ValueEncoderFunc for map[string]* types.
//
// Deprecated: MapEncodeValue is not registered by default. Use MapCodec.EncodeValue instead.
func (dve DefaultValueEncoders) MapEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Kind() != reflect.Map || val.Type().Key().Kind() != reflect.String ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "MapEncodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Map***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		// If we have a nill map but we can't WriteNull, that means we're probably trying to encode
		// to a TopLevel document. We can't currently tell if this is what actually happened, but if
		// there's a deeper underlying problem, the error will also be returned from WriteDocument,
		// so just continue. The operations on a map reflection value are valid, so we can call
		// MapKeys within mapEncodeValue without a problem.
		err := vw.WriteNull()
		if err == nil ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	dw, err := vw.WriteDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return dve.mapEncodeValue(ec, dw, val, nil)
***REMOVED***

// mapEncodeValue handles encoding of the values of a map. The collisionFn returns
// true if the provided key exists, this is mainly used for inline maps in the
// struct codec.
func (dve DefaultValueEncoders) mapEncodeValue(ec EncodeContext, dw bsonrw.DocumentWriter, val reflect.Value, collisionFn func(string) bool) error ***REMOVED***

	elemType := val.Type().Elem()
	encoder, err := ec.LookupEncoder(elemType)
	if err != nil && elemType.Kind() != reflect.Interface ***REMOVED***
		return err
	***REMOVED***

	keys := val.MapKeys()
	for _, key := range keys ***REMOVED***
		if collisionFn != nil && collisionFn(key.String()) ***REMOVED***
			return fmt.Errorf("Key %s of inlined map conflicts with a struct field name", key)
		***REMOVED***

		currEncoder, currVal, lookupErr := dve.lookupElementEncoder(ec, encoder, val.MapIndex(key))
		if lookupErr != nil && lookupErr != errInvalidValue ***REMOVED***
			return lookupErr
		***REMOVED***

		vw, err := dw.WriteDocumentElement(key.String())
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if lookupErr == errInvalidValue ***REMOVED***
			err = vw.WriteNull()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		err = currEncoder.EncodeValue(ec, vw, currVal)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return dw.WriteDocumentEnd()
***REMOVED***

// ArrayEncodeValue is the ValueEncoderFunc for array types.
func (dve DefaultValueEncoders) ArrayEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Kind() != reflect.Array ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "ArrayEncodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Array***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	// If we have a []primitive.E we want to treat it as a document instead of as an array.
	if val.Type().Elem() == tE ***REMOVED***
		dw, err := vw.WriteDocument()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		for idx := 0; idx < val.Len(); idx++ ***REMOVED***
			e := val.Index(idx).Interface().(primitive.E)
			err = encodeElement(ec, dw, e)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		return dw.WriteDocumentEnd()
	***REMOVED***

	// If we have a []byte we want to treat it as a binary instead of as an array.
	if val.Type().Elem() == tByte ***REMOVED***
		var byteSlice []byte
		for idx := 0; idx < val.Len(); idx++ ***REMOVED***
			byteSlice = append(byteSlice, val.Index(idx).Interface().(byte))
		***REMOVED***
		return vw.WriteBinary(byteSlice)
	***REMOVED***

	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	elemType := val.Type().Elem()
	encoder, err := ec.LookupEncoder(elemType)
	if err != nil && elemType.Kind() != reflect.Interface ***REMOVED***
		return err
	***REMOVED***

	for idx := 0; idx < val.Len(); idx++ ***REMOVED***
		currEncoder, currVal, lookupErr := dve.lookupElementEncoder(ec, encoder, val.Index(idx))
		if lookupErr != nil && lookupErr != errInvalidValue ***REMOVED***
			return lookupErr
		***REMOVED***

		vw, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if lookupErr == errInvalidValue ***REMOVED***
			err = vw.WriteNull()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		err = currEncoder.EncodeValue(ec, vw, currVal)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return aw.WriteArrayEnd()
***REMOVED***

// SliceEncodeValue is the ValueEncoderFunc for slice types.
//
// Deprecated: SliceEncodeValue is not registered by default. Use SliceCodec.EncodeValue instead.
func (dve DefaultValueEncoders) SliceEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Kind() != reflect.Slice ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "SliceEncodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Slice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***

	// If we have a []primitive.E we want to treat it as a document instead of as an array.
	if val.Type().ConvertibleTo(tD) ***REMOVED***
		d := val.Convert(tD).Interface().(primitive.D)

		dw, err := vw.WriteDocument()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		for _, e := range d ***REMOVED***
			err = encodeElement(ec, dw, e)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		return dw.WriteDocumentEnd()
	***REMOVED***

	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	elemType := val.Type().Elem()
	encoder, err := ec.LookupEncoder(elemType)
	if err != nil && elemType.Kind() != reflect.Interface ***REMOVED***
		return err
	***REMOVED***

	for idx := 0; idx < val.Len(); idx++ ***REMOVED***
		currEncoder, currVal, lookupErr := dve.lookupElementEncoder(ec, encoder, val.Index(idx))
		if lookupErr != nil && lookupErr != errInvalidValue ***REMOVED***
			return lookupErr
		***REMOVED***

		vw, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if lookupErr == errInvalidValue ***REMOVED***
			err = vw.WriteNull()
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			continue
		***REMOVED***

		err = currEncoder.EncodeValue(ec, vw, currVal)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return aw.WriteArrayEnd()
***REMOVED***

func (dve DefaultValueEncoders) lookupElementEncoder(ec EncodeContext, origEncoder ValueEncoder, currVal reflect.Value) (ValueEncoder, reflect.Value, error) ***REMOVED***
	if origEncoder != nil || (currVal.Kind() != reflect.Interface) ***REMOVED***
		return origEncoder, currVal, nil
	***REMOVED***
	currVal = currVal.Elem()
	if !currVal.IsValid() ***REMOVED***
		return nil, currVal, errInvalidValue
	***REMOVED***
	currEncoder, err := ec.LookupEncoder(currVal.Type())

	return currEncoder, currVal, err
***REMOVED***

// EmptyInterfaceEncodeValue is the ValueEncoderFunc for interface***REMOVED******REMOVED***.
//
// Deprecated: EmptyInterfaceEncodeValue is not registered by default. Use EmptyInterfaceCodec.EncodeValue instead.
func (dve DefaultValueEncoders) EmptyInterfaceEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tEmpty ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "EmptyInterfaceEncodeValue", Types: []reflect.Type***REMOVED***tEmpty***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***
	encoder, err := ec.LookupEncoder(val.Elem().Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return encoder.EncodeValue(ec, vw, val.Elem())
***REMOVED***

// ValueMarshalerEncodeValue is the ValueEncoderFunc for ValueMarshaler implementations.
func (dve DefaultValueEncoders) ValueMarshalerEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	// Either val or a pointer to val must implement ValueMarshaler
	switch ***REMOVED***
	case !val.IsValid():
		return ValueEncoderError***REMOVED***Name: "ValueMarshalerEncodeValue", Types: []reflect.Type***REMOVED***tValueMarshaler***REMOVED***, Received: val***REMOVED***
	case val.Type().Implements(tValueMarshaler):
		// If ValueMarshaler is implemented on a concrete type, make sure that val isn't a nil pointer
		if isImplementationNil(val, tValueMarshaler) ***REMOVED***
			return vw.WriteNull()
		***REMOVED***
	case reflect.PtrTo(val.Type()).Implements(tValueMarshaler) && val.CanAddr():
		val = val.Addr()
	default:
		return ValueEncoderError***REMOVED***Name: "ValueMarshalerEncodeValue", Types: []reflect.Type***REMOVED***tValueMarshaler***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	fn := val.Convert(tValueMarshaler).MethodByName("MarshalBSONValue")
	returns := fn.Call(nil)
	if !returns[2].IsNil() ***REMOVED***
		return returns[2].Interface().(error)
	***REMOVED***
	t, data := returns[0].Interface().(bsontype.Type), returns[1].Interface().([]byte)
	return bsonrw.Copier***REMOVED******REMOVED***.CopyValueFromBytes(vw, t, data)
***REMOVED***

// MarshalerEncodeValue is the ValueEncoderFunc for Marshaler implementations.
func (dve DefaultValueEncoders) MarshalerEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	// Either val or a pointer to val must implement Marshaler
	switch ***REMOVED***
	case !val.IsValid():
		return ValueEncoderError***REMOVED***Name: "MarshalerEncodeValue", Types: []reflect.Type***REMOVED***tMarshaler***REMOVED***, Received: val***REMOVED***
	case val.Type().Implements(tMarshaler):
		// If Marshaler is implemented on a concrete type, make sure that val isn't a nil pointer
		if isImplementationNil(val, tMarshaler) ***REMOVED***
			return vw.WriteNull()
		***REMOVED***
	case reflect.PtrTo(val.Type()).Implements(tMarshaler) && val.CanAddr():
		val = val.Addr()
	default:
		return ValueEncoderError***REMOVED***Name: "MarshalerEncodeValue", Types: []reflect.Type***REMOVED***tMarshaler***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	fn := val.Convert(tMarshaler).MethodByName("MarshalBSON")
	returns := fn.Call(nil)
	if !returns[1].IsNil() ***REMOVED***
		return returns[1].Interface().(error)
	***REMOVED***
	data := returns[0].Interface().([]byte)
	return bsonrw.Copier***REMOVED******REMOVED***.CopyValueFromBytes(vw, bsontype.EmbeddedDocument, data)
***REMOVED***

// ProxyEncodeValue is the ValueEncoderFunc for Proxy implementations.
func (dve DefaultValueEncoders) ProxyEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	// Either val or a pointer to val must implement Proxy
	switch ***REMOVED***
	case !val.IsValid():
		return ValueEncoderError***REMOVED***Name: "ProxyEncodeValue", Types: []reflect.Type***REMOVED***tProxy***REMOVED***, Received: val***REMOVED***
	case val.Type().Implements(tProxy):
		// If Proxy is implemented on a concrete type, make sure that val isn't a nil pointer
		if isImplementationNil(val, tProxy) ***REMOVED***
			return vw.WriteNull()
		***REMOVED***
	case reflect.PtrTo(val.Type()).Implements(tProxy) && val.CanAddr():
		val = val.Addr()
	default:
		return ValueEncoderError***REMOVED***Name: "ProxyEncodeValue", Types: []reflect.Type***REMOVED***tProxy***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	fn := val.Convert(tProxy).MethodByName("ProxyBSON")
	returns := fn.Call(nil)
	if !returns[1].IsNil() ***REMOVED***
		return returns[1].Interface().(error)
	***REMOVED***
	data := returns[0]
	var encoder ValueEncoder
	var err error
	if data.Elem().IsValid() ***REMOVED***
		encoder, err = ec.LookupEncoder(data.Elem().Type())
	***REMOVED*** else ***REMOVED***
		encoder, err = ec.LookupEncoder(nil)
	***REMOVED***
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return encoder.EncodeValue(ec, vw, data.Elem())
***REMOVED***

// JavaScriptEncodeValue is the ValueEncoderFunc for the primitive.JavaScript type.
func (DefaultValueEncoders) JavaScriptEncodeValue(ectx EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tJavaScript ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "JavaScriptEncodeValue", Types: []reflect.Type***REMOVED***tJavaScript***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return vw.WriteJavascript(val.String())
***REMOVED***

// SymbolEncodeValue is the ValueEncoderFunc for the primitive.Symbol type.
func (DefaultValueEncoders) SymbolEncodeValue(ectx EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tSymbol ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "SymbolEncodeValue", Types: []reflect.Type***REMOVED***tSymbol***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return vw.WriteSymbol(val.String())
***REMOVED***

// BinaryEncodeValue is the ValueEncoderFunc for Binary.
func (DefaultValueEncoders) BinaryEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tBinary ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "BinaryEncodeValue", Types: []reflect.Type***REMOVED***tBinary***REMOVED***, Received: val***REMOVED***
	***REMOVED***
	b := val.Interface().(primitive.Binary)

	return vw.WriteBinaryWithSubtype(b.Data, b.Subtype)
***REMOVED***

// UndefinedEncodeValue is the ValueEncoderFunc for Undefined.
func (DefaultValueEncoders) UndefinedEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tUndefined ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "UndefinedEncodeValue", Types: []reflect.Type***REMOVED***tUndefined***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return vw.WriteUndefined()
***REMOVED***

// DateTimeEncodeValue is the ValueEncoderFunc for DateTime.
func (DefaultValueEncoders) DateTimeEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tDateTime ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "DateTimeEncodeValue", Types: []reflect.Type***REMOVED***tDateTime***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return vw.WriteDateTime(val.Int())
***REMOVED***

// NullEncodeValue is the ValueEncoderFunc for Null.
func (DefaultValueEncoders) NullEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tNull ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "NullEncodeValue", Types: []reflect.Type***REMOVED***tNull***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return vw.WriteNull()
***REMOVED***

// RegexEncodeValue is the ValueEncoderFunc for Regex.
func (DefaultValueEncoders) RegexEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tRegex ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "RegexEncodeValue", Types: []reflect.Type***REMOVED***tRegex***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	regex := val.Interface().(primitive.Regex)

	return vw.WriteRegex(regex.Pattern, regex.Options)
***REMOVED***

// DBPointerEncodeValue is the ValueEncoderFunc for DBPointer.
func (DefaultValueEncoders) DBPointerEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tDBPointer ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "DBPointerEncodeValue", Types: []reflect.Type***REMOVED***tDBPointer***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	dbp := val.Interface().(primitive.DBPointer)

	return vw.WriteDBPointer(dbp.DB, dbp.Pointer)
***REMOVED***

// TimestampEncodeValue is the ValueEncoderFunc for Timestamp.
func (DefaultValueEncoders) TimestampEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tTimestamp ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "TimestampEncodeValue", Types: []reflect.Type***REMOVED***tTimestamp***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	ts := val.Interface().(primitive.Timestamp)

	return vw.WriteTimestamp(ts.T, ts.I)
***REMOVED***

// MinKeyEncodeValue is the ValueEncoderFunc for MinKey.
func (DefaultValueEncoders) MinKeyEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tMinKey ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "MinKeyEncodeValue", Types: []reflect.Type***REMOVED***tMinKey***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return vw.WriteMinKey()
***REMOVED***

// MaxKeyEncodeValue is the ValueEncoderFunc for MaxKey.
func (DefaultValueEncoders) MaxKeyEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tMaxKey ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "MaxKeyEncodeValue", Types: []reflect.Type***REMOVED***tMaxKey***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return vw.WriteMaxKey()
***REMOVED***

// CoreDocumentEncodeValue is the ValueEncoderFunc for bsoncore.Document.
func (DefaultValueEncoders) CoreDocumentEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tCoreDocument ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "CoreDocumentEncodeValue", Types: []reflect.Type***REMOVED***tCoreDocument***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	cdoc := val.Interface().(bsoncore.Document)

	return bsonrw.Copier***REMOVED******REMOVED***.CopyDocumentFromBytes(vw, cdoc)
***REMOVED***

// CodeWithScopeEncodeValue is the ValueEncoderFunc for CodeWithScope.
func (dve DefaultValueEncoders) CodeWithScopeEncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tCodeWithScope ***REMOVED***
		return ValueEncoderError***REMOVED***Name: "CodeWithScopeEncodeValue", Types: []reflect.Type***REMOVED***tCodeWithScope***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	cws := val.Interface().(primitive.CodeWithScope)

	dw, err := vw.WriteCodeWithScope(string(cws.Code))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	sw := sliceWriterPool.Get().(*bsonrw.SliceWriter)
	defer sliceWriterPool.Put(sw)
	*sw = (*sw)[:0]

	scopeVW := bvwPool.Get(sw)
	defer bvwPool.Put(scopeVW)

	encoder, err := ec.LookupEncoder(reflect.TypeOf(cws.Scope))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = encoder.EncodeValue(ec, scopeVW, reflect.ValueOf(cws.Scope))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = bsonrw.Copier***REMOVED******REMOVED***.CopyBytesToDocumentWriter(dw, *sw)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return dw.WriteDocumentEnd()
***REMOVED***

// isImplementationNil returns if val is a nil pointer and inter is implemented on a concrete type
func isImplementationNil(val reflect.Value, inter reflect.Type) bool ***REMOVED***
	vt := val.Type()
	for vt.Kind() == reflect.Ptr ***REMOVED***
		vt = vt.Elem()
	***REMOVED***
	return vt.Implements(inter) && val.Kind() == reflect.Ptr && val.IsNil()
***REMOVED***
