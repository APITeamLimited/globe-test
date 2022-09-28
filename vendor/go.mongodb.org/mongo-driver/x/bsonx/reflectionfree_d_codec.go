// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	tPrimitiveD          = reflect.TypeOf(primitive.D***REMOVED******REMOVED***)
	tPrimitiveCWS        = reflect.TypeOf(primitive.CodeWithScope***REMOVED******REMOVED***)
	defaultValueEncoders = bsoncodec.DefaultValueEncoders***REMOVED******REMOVED***
	defaultValueDecoders = bsoncodec.DefaultValueDecoders***REMOVED******REMOVED***
)

type reflectionFreeDCodec struct***REMOVED******REMOVED***

// ReflectionFreeDCodec is a ValueEncoder for the primitive.D type that does not use reflection.
var ReflectionFreeDCodec bsoncodec.ValueCodec = &reflectionFreeDCodec***REMOVED******REMOVED***

func (r *reflectionFreeDCodec) EncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tPrimitiveD ***REMOVED***
		return bsoncodec.ValueEncoderError***REMOVED***Name: "DEncodeValue", Types: []reflect.Type***REMOVED***tPrimitiveD***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***

	doc := val.Interface().(primitive.D)
	return r.encodeDocument(ec, vw, doc)
***REMOVED***

func (r *reflectionFreeDCodec) DecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || !val.CanSet() || val.Type() != tPrimitiveD ***REMOVED***
		return bsoncodec.ValueDecoderError***REMOVED***Name: "DDecodeValue", Kinds: []reflect.Kind***REMOVED***reflect.Slice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	default:
		return fmt.Errorf("cannot decode %v into a primitive.D", vrType)
	***REMOVED***

	doc, err := r.decodeDocument(dc, vr)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	val.Set(reflect.ValueOf(doc))
	return nil
***REMOVED***

func (r *reflectionFreeDCodec) decodeDocument(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader) (primitive.D, error) ***REMOVED***
	dr, err := vr.ReadDocument()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	doc := primitive.D***REMOVED******REMOVED***
	for ***REMOVED***
		key, elemVr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		val, err := r.decodeValue(dc, elemVr)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		doc = append(doc, primitive.E***REMOVED***Key: key, Value: val***REMOVED***)
	***REMOVED***

	return doc, nil
***REMOVED***

func (r *reflectionFreeDCodec) decodeArray(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader) (primitive.A, error) ***REMOVED***
	ar, err := vr.ReadArray()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	array := primitive.A***REMOVED******REMOVED***
	for ***REMOVED***
		arrayValReader, err := ar.ReadValue()
		if err == bsonrw.ErrEOA ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		val, err := r.decodeValue(dc, arrayValReader)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		array = append(array, val)
	***REMOVED***

	return array, nil
***REMOVED***

func (r *reflectionFreeDCodec) decodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	switch vrType := vr.Type(); vrType ***REMOVED***
	case bsontype.Null:
		return nil, vr.ReadNull()
	case bsontype.Double:
		return vr.ReadDouble()
	case bsontype.String:
		return vr.ReadString()
	case bsontype.Binary:
		data, subtype, err := vr.ReadBinary()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return primitive.Binary***REMOVED***
			Data:    data,
			Subtype: subtype,
		***REMOVED***, nil
	case bsontype.Undefined:
		return primitive.Undefined***REMOVED******REMOVED***, vr.ReadUndefined()
	case bsontype.ObjectID:
		return vr.ReadObjectID()
	case bsontype.Boolean:
		return vr.ReadBoolean()
	case bsontype.DateTime:
		dt, err := vr.ReadDateTime()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return primitive.DateTime(dt), nil
	case bsontype.Regex:
		pattern, options, err := vr.ReadRegex()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return primitive.Regex***REMOVED***
			Pattern: pattern,
			Options: options,
		***REMOVED***, nil
	case bsontype.DBPointer:
		ns, oid, err := vr.ReadDBPointer()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return primitive.DBPointer***REMOVED***
			DB:      ns,
			Pointer: oid,
		***REMOVED***, nil
	case bsontype.JavaScript:
		js, err := vr.ReadJavascript()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return primitive.JavaScript(js), nil
	case bsontype.Symbol:
		sym, err := vr.ReadSymbol()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return primitive.Symbol(sym), nil
	case bsontype.CodeWithScope:
		cws := reflect.New(tPrimitiveCWS).Elem()
		err := defaultValueDecoders.CodeWithScopeDecodeValue(dc, vr, cws)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return cws.Interface().(primitive.CodeWithScope), nil
	case bsontype.Int32:
		return vr.ReadInt32()
	case bsontype.Int64:
		return vr.ReadInt64()
	case bsontype.Timestamp:
		t, i, err := vr.ReadTimestamp()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return primitive.Timestamp***REMOVED***
			T: t,
			I: i,
		***REMOVED***, nil
	case bsontype.Decimal128:
		return vr.ReadDecimal128()
	case bsontype.MinKey:
		return primitive.MinKey***REMOVED******REMOVED***, vr.ReadMinKey()
	case bsontype.MaxKey:
		return primitive.MaxKey***REMOVED******REMOVED***, vr.ReadMaxKey()
	case bsontype.Type(0), bsontype.EmbeddedDocument:
		return r.decodeDocument(dc, vr)
	case bsontype.Array:
		return r.decodeArray(dc, vr)
	default:
		return nil, fmt.Errorf("cannot decode invalid BSON type %s", vrType)
	***REMOVED***
***REMOVED***

func (r *reflectionFreeDCodec) encodeDocumentValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, v interface***REMOVED******REMOVED***) error ***REMOVED***
	switch val := v.(type) ***REMOVED***
	case int:
		return r.encodeInt(vw, val)
	case int8:
		return vw.WriteInt32(int32(val))
	case int16:
		return vw.WriteInt32(int32(val))
	case int32:
		return vw.WriteInt32(val)
	case int64:
		return r.encodeInt64(ec, vw, val)
	case uint:
		return r.encodeUint64(ec, vw, uint64(val))
	case uint8:
		return vw.WriteInt32(int32(val))
	case uint16:
		return vw.WriteInt32(int32(val))
	case uint32:
		return r.encodeUint64(ec, vw, uint64(val))
	case uint64:
		return r.encodeUint64(ec, vw, val)
	case float32:
		return vw.WriteDouble(float64(val))
	case float64:
		return vw.WriteDouble(val)
	case []byte:
		return vw.WriteBinary(val)
	case primitive.Binary:
		return vw.WriteBinaryWithSubtype(val.Data, val.Subtype)
	case bool:
		return vw.WriteBoolean(val)
	case primitive.CodeWithScope:
		return defaultValueEncoders.CodeWithScopeEncodeValue(ec, vw, reflect.ValueOf(val))
	case primitive.DBPointer:
		return vw.WriteDBPointer(val.DB, val.Pointer)
	case primitive.DateTime:
		return vw.WriteDateTime(int64(val))
	case time.Time:
		dt := primitive.NewDateTimeFromTime(val)
		return vw.WriteDateTime(int64(dt))
	case primitive.Decimal128:
		return vw.WriteDecimal128(val)
	case primitive.JavaScript:
		return vw.WriteJavascript(string(val))
	case primitive.MinKey:
		return vw.WriteMinKey()
	case primitive.MaxKey:
		return vw.WriteMaxKey()
	case primitive.Null, nil:
		return vw.WriteNull()
	case primitive.ObjectID:
		return vw.WriteObjectID(val)
	case primitive.Regex:
		return vw.WriteRegex(val.Pattern, val.Options)
	case string:
		return vw.WriteString(val)
	case primitive.Symbol:
		return vw.WriteSymbol(string(val))
	case primitive.Timestamp:
		return vw.WriteTimestamp(val.T, val.I)
	case primitive.Undefined:
		return vw.WriteUndefined()
	case primitive.D:
		return r.encodeDocument(ec, vw, val)
	case primitive.A:
		return r.encodePrimitiveA(ec, vw, val)
	case []interface***REMOVED******REMOVED***:
		return r.encodePrimitiveA(ec, vw, val)
	case []primitive.D:
		return r.encodeSliceD(ec, vw, val)
	case []int:
		return r.encodeSliceInt(vw, val)
	case []int8:
		return r.encodeSliceInt8(vw, val)
	case []int16:
		return r.encodeSliceInt16(vw, val)
	case []int32:
		return r.encodeSliceInt32(vw, val)
	case []int64:
		return r.encodeSliceInt64(ec, vw, val)
	case []uint:
		return r.encodeSliceUint(ec, vw, val)
	case []uint16:
		return r.encodeSliceUint16(vw, val)
	case []uint32:
		return r.encodeSliceUint32(ec, vw, val)
	case []uint64:
		return r.encodeSliceUint64(ec, vw, val)
	case [][]byte:
		return r.encodeSliceByteSlice(vw, val)
	case []primitive.Binary:
		return r.encodeSliceBinary(vw, val)
	case []bool:
		return r.encodeSliceBoolean(vw, val)
	case []primitive.CodeWithScope:
		return r.encodeSliceCWS(ec, vw, val)
	case []primitive.DBPointer:
		return r.encodeSliceDBPointer(vw, val)
	case []primitive.DateTime:
		return r.encodeSliceDateTime(vw, val)
	case []time.Time:
		return r.encodeSliceTimeTime(vw, val)
	case []primitive.Decimal128:
		return r.encodeSliceDecimal128(vw, val)
	case []float32:
		return r.encodeSliceFloat32(vw, val)
	case []float64:
		return r.encodeSliceFloat64(vw, val)
	case []primitive.JavaScript:
		return r.encodeSliceJavaScript(vw, val)
	case []primitive.MinKey:
		return r.encodeSliceMinKey(vw, val)
	case []primitive.MaxKey:
		return r.encodeSliceMaxKey(vw, val)
	case []primitive.Null:
		return r.encodeSliceNull(vw, val)
	case []primitive.ObjectID:
		return r.encodeSliceObjectID(vw, val)
	case []primitive.Regex:
		return r.encodeSliceRegex(vw, val)
	case []string:
		return r.encodeSliceString(vw, val)
	case []primitive.Symbol:
		return r.encodeSliceSymbol(vw, val)
	case []primitive.Timestamp:
		return r.encodeSliceTimestamp(vw, val)
	case []primitive.Undefined:
		return r.encodeSliceUndefined(vw, val)
	default:
		return fmt.Errorf("value of type %T not supported", v)
	***REMOVED***
***REMOVED***

func (r *reflectionFreeDCodec) encodeInt(vw bsonrw.ValueWriter, val int) error ***REMOVED***
	if fitsIn32Bits(int64(val)) ***REMOVED***
		return vw.WriteInt32(int32(val))
	***REMOVED***
	return vw.WriteInt64(int64(val))
***REMOVED***

func (r *reflectionFreeDCodec) encodeInt64(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val int64) error ***REMOVED***
	if ec.MinSize && fitsIn32Bits(val) ***REMOVED***
		return vw.WriteInt32(int32(val))
	***REMOVED***
	return vw.WriteInt64(val)
***REMOVED***

func (r *reflectionFreeDCodec) encodeUint64(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val uint64) error ***REMOVED***
	if ec.MinSize && val <= math.MaxInt32 ***REMOVED***
		return vw.WriteInt32(int32(val))
	***REMOVED***
	if val > math.MaxInt64 ***REMOVED***
		return fmt.Errorf("%d overflows int64", val)
	***REMOVED***

	return vw.WriteInt64(int64(val))
***REMOVED***

func (r *reflectionFreeDCodec) encodeDocument(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, doc primitive.D) error ***REMOVED***
	dw, err := vw.WriteDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, elem := range doc ***REMOVED***
		docValWriter, err := dw.WriteDocumentElement(elem.Key)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := r.encodeDocumentValue(ec, docValWriter, elem.Value); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return dw.WriteDocumentEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceByteSlice(vw bsonrw.ValueWriter, arr [][]byte) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteBinary(val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceBinary(vw bsonrw.ValueWriter, arr []primitive.Binary) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteBinaryWithSubtype(val.Data, val.Subtype); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceBoolean(vw bsonrw.ValueWriter, arr []bool) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteBoolean(val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceCWS(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, arr []primitive.CodeWithScope) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := defaultValueEncoders.CodeWithScopeEncodeValue(ec, arrayValWriter, reflect.ValueOf(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceDBPointer(vw bsonrw.ValueWriter, arr []primitive.DBPointer) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteDBPointer(val.DB, val.Pointer); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceDateTime(vw bsonrw.ValueWriter, arr []primitive.DateTime) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteDateTime(int64(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceTimeTime(vw bsonrw.ValueWriter, arr []time.Time) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		dt := primitive.NewDateTimeFromTime(val)
		if err := arrayValWriter.WriteDateTime(int64(dt)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceDecimal128(vw bsonrw.ValueWriter, arr []primitive.Decimal128) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteDecimal128(val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceFloat32(vw bsonrw.ValueWriter, arr []float32) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteDouble(float64(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceFloat64(vw bsonrw.ValueWriter, arr []float64) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteDouble(val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceJavaScript(vw bsonrw.ValueWriter, arr []primitive.JavaScript) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteJavascript(string(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceMinKey(vw bsonrw.ValueWriter, arr []primitive.MinKey) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteMinKey(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceMaxKey(vw bsonrw.ValueWriter, arr []primitive.MaxKey) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteMaxKey(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceNull(vw bsonrw.ValueWriter, arr []primitive.Null) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteNull(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceObjectID(vw bsonrw.ValueWriter, arr []primitive.ObjectID) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteObjectID(val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceRegex(vw bsonrw.ValueWriter, arr []primitive.Regex) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteRegex(val.Pattern, val.Options); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceString(vw bsonrw.ValueWriter, arr []string) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteString(val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceSymbol(vw bsonrw.ValueWriter, arr []primitive.Symbol) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteSymbol(string(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceTimestamp(vw bsonrw.ValueWriter, arr []primitive.Timestamp) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteTimestamp(val.T, val.I); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceUndefined(vw bsonrw.ValueWriter, arr []primitive.Undefined) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteUndefined(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodePrimitiveA(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, arr primitive.A) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := r.encodeDocumentValue(ec, arrayValWriter, val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceD(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, arr []primitive.D) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := r.encodeDocument(ec, arrayValWriter, val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceInt(vw bsonrw.ValueWriter, arr []int) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := r.encodeInt(arrayValWriter, val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceInt8(vw bsonrw.ValueWriter, arr []int8) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteInt32(int32(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceInt16(vw bsonrw.ValueWriter, arr []int16) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteInt32(int32(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceInt32(vw bsonrw.ValueWriter, arr []int32) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteInt32(val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceInt64(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, arr []int64) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := r.encodeInt64(ec, arrayValWriter, val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceUint(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, arr []uint) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := r.encodeUint64(ec, arrayValWriter, uint64(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceUint16(vw bsonrw.ValueWriter, arr []uint16) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := arrayValWriter.WriteInt32(int32(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceUint32(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, arr []uint32) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := r.encodeUint64(ec, arrayValWriter, uint64(val)); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (r *reflectionFreeDCodec) encodeSliceUint64(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, arr []uint64) error ***REMOVED***
	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		arrayValWriter, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		if err := r.encodeUint64(ec, arrayValWriter, val); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func fitsIn32Bits(i int64) bool ***REMOVED***
	return math.MinInt32 <= i && i <= math.MaxInt32
***REMOVED***
