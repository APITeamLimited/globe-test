// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonx

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var primitiveCodecs PrimitiveCodecs

var tDocument = reflect.TypeOf((Doc)(nil))
var tArray = reflect.TypeOf((Arr)(nil))
var tValue = reflect.TypeOf(Val***REMOVED******REMOVED***)
var tElementSlice = reflect.TypeOf(([]Elem)(nil))

// PrimitiveCodecs is a namespace for all of the default bsoncodec.Codecs for the primitive types
// defined in this package.
type PrimitiveCodecs struct***REMOVED******REMOVED***

// RegisterPrimitiveCodecs will register the encode and decode methods attached to PrimitiveCodecs
// with the provided RegistryBuilder. if rb is nil, a new empty RegistryBuilder will be created.
func (pc PrimitiveCodecs) RegisterPrimitiveCodecs(rb *bsoncodec.RegistryBuilder) ***REMOVED***
	if rb == nil ***REMOVED***
		panic(errors.New("argument to RegisterPrimitiveCodecs must not be nil"))
	***REMOVED***

	rb.
		RegisterTypeEncoder(tDocument, bsoncodec.ValueEncoderFunc(pc.DocumentEncodeValue)).
		RegisterTypeEncoder(tArray, bsoncodec.ValueEncoderFunc(pc.ArrayEncodeValue)).
		RegisterTypeEncoder(tValue, bsoncodec.ValueEncoderFunc(pc.ValueEncodeValue)).
		RegisterTypeEncoder(tElementSlice, bsoncodec.ValueEncoderFunc(pc.ElementSliceEncodeValue)).
		RegisterTypeDecoder(tDocument, bsoncodec.ValueDecoderFunc(pc.DocumentDecodeValue)).
		RegisterTypeDecoder(tArray, bsoncodec.ValueDecoderFunc(pc.ArrayDecodeValue)).
		RegisterTypeDecoder(tValue, bsoncodec.ValueDecoderFunc(pc.ValueDecodeValue)).
		RegisterTypeDecoder(tElementSlice, bsoncodec.ValueDecoderFunc(pc.ElementSliceDecodeValue))
***REMOVED***

// DocumentEncodeValue is the ValueEncoderFunc for *Document.
func (pc PrimitiveCodecs) DocumentEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tDocument ***REMOVED***
		return bsoncodec.ValueEncoderError***REMOVED***Name: "DocumentEncodeValue", Types: []reflect.Type***REMOVED***tDocument***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***

	doc := val.Interface().(Doc)

	dw, err := vw.WriteDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return pc.encodeDocument(ec, dw, doc)
***REMOVED***

// DocumentDecodeValue is the ValueDecoderFunc for *Document.
func (pc PrimitiveCodecs) DocumentDecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tDocument ***REMOVED***
		return bsoncodec.ValueDecoderError***REMOVED***Name: "DocumentDecodeValue", Types: []reflect.Type***REMOVED***tDocument***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return pc.documentDecodeValue(dctx, vr, val.Addr().Interface().(*Doc))
***REMOVED***

func (pc PrimitiveCodecs) documentDecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, doc *Doc) error ***REMOVED***

	dr, err := vr.ReadDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return pc.decodeDocument(dctx, dr, doc)
***REMOVED***

// ArrayEncodeValue is the ValueEncoderFunc for *Array.
func (pc PrimitiveCodecs) ArrayEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tArray ***REMOVED***
		return bsoncodec.ValueEncoderError***REMOVED***Name: "ArrayEncodeValue", Types: []reflect.Type***REMOVED***tArray***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***

	arr := val.Interface().(Arr)

	aw, err := vw.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, val := range arr ***REMOVED***
		dvw, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = pc.encodeValue(ec, dvw, val)

		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

// ArrayDecodeValue is the ValueDecoderFunc for *Array.
func (pc PrimitiveCodecs) ArrayDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tArray ***REMOVED***
		return bsoncodec.ValueDecoderError***REMOVED***Name: "ArrayDecodeValue", Types: []reflect.Type***REMOVED***tArray***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	ar, err := vr.ReadArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if val.IsNil() ***REMOVED***
		val.Set(reflect.MakeSlice(tArray, 0, 0))
	***REMOVED***
	val.SetLen(0)

	for ***REMOVED***
		vr, err := ar.ReadValue()
		if err == bsonrw.ErrEOA ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var elem Val
		err = pc.valueDecodeValue(dc, vr, &elem)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		val.Set(reflect.Append(val, reflect.ValueOf(elem)))
	***REMOVED***

	return nil
***REMOVED***

// ElementSliceEncodeValue is the ValueEncoderFunc for []*Element.
func (pc PrimitiveCodecs) ElementSliceEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tElementSlice ***REMOVED***
		return bsoncodec.ValueEncoderError***REMOVED***Name: "ElementSliceEncodeValue", Types: []reflect.Type***REMOVED***tElementSlice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		return vw.WriteNull()
	***REMOVED***

	return pc.DocumentEncodeValue(ec, vw, val.Convert(tDocument))
***REMOVED***

// ElementSliceDecodeValue is the ValueDecoderFunc for []*Element.
func (pc PrimitiveCodecs) ElementSliceDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tElementSlice ***REMOVED***
		return bsoncodec.ValueDecoderError***REMOVED***Name: "ElementSliceDecodeValue", Types: []reflect.Type***REMOVED***tElementSlice***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	if val.IsNil() ***REMOVED***
		val.Set(reflect.MakeSlice(val.Type(), 0, 0))
	***REMOVED***

	val.SetLen(0)

	dr, err := vr.ReadDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	elems := make([]reflect.Value, 0)
	for ***REMOVED***
		key, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var elem Elem
		err = pc.elementDecodeValue(dc, vr, key, &elem)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		elems = append(elems, reflect.ValueOf(elem))
	***REMOVED***

	val.Set(reflect.Append(val, elems...))
	return nil
***REMOVED***

// ValueEncodeValue is the ValueEncoderFunc for *Value.
func (pc PrimitiveCodecs) ValueEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error ***REMOVED***
	if !val.IsValid() || val.Type() != tValue ***REMOVED***
		return bsoncodec.ValueEncoderError***REMOVED***Name: "ValueEncodeValue", Types: []reflect.Type***REMOVED***tValue***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	v := val.Interface().(Val)

	return pc.encodeValue(ec, vw, v)
***REMOVED***

// ValueDecodeValue is the ValueDecoderFunc for *Value.
func (pc PrimitiveCodecs) ValueDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error ***REMOVED***
	if !val.CanSet() || val.Type() != tValue ***REMOVED***
		return bsoncodec.ValueDecoderError***REMOVED***Name: "ValueDecodeValue", Types: []reflect.Type***REMOVED***tValue***REMOVED***, Received: val***REMOVED***
	***REMOVED***

	return pc.valueDecodeValue(dc, vr, val.Addr().Interface().(*Val))
***REMOVED***

// encodeDocument is a separate function that we use because CodeWithScope
// returns us a DocumentWriter and we need to do the same logic that we would do
// for a document but cannot use a Codec.
func (pc PrimitiveCodecs) encodeDocument(ec bsoncodec.EncodeContext, dw bsonrw.DocumentWriter, doc Doc) error ***REMOVED***
	for _, elem := range doc ***REMOVED***
		dvw, err := dw.WriteDocumentElement(elem.Key)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = pc.encodeValue(ec, dvw, elem.Value)

		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return dw.WriteDocumentEnd()
***REMOVED***

// DecodeDocument haves decoding into a Doc from a bsonrw.DocumentReader.
func (pc PrimitiveCodecs) DecodeDocument(dctx bsoncodec.DecodeContext, dr bsonrw.DocumentReader, pdoc *Doc) error ***REMOVED***
	return pc.decodeDocument(dctx, dr, pdoc)
***REMOVED***

func (pc PrimitiveCodecs) decodeDocument(dctx bsoncodec.DecodeContext, dr bsonrw.DocumentReader, pdoc *Doc) error ***REMOVED***
	if *pdoc == nil ***REMOVED***
		*pdoc = make(Doc, 0)
	***REMOVED***
	*pdoc = (*pdoc)[:0]
	for ***REMOVED***
		key, vr, err := dr.ReadElement()
		if err == bsonrw.ErrEOD ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		var elem Elem
		err = pc.elementDecodeValue(dctx, vr, key, &elem)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		*pdoc = append(*pdoc, elem)
	***REMOVED***
	return nil
***REMOVED***

func (pc PrimitiveCodecs) elementDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, key string, elem *Elem) error ***REMOVED***
	var val Val
	switch vr.Type() ***REMOVED***
	case bsontype.Double:
		f64, err := vr.ReadDouble()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Double(f64)
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = String(str)
	case bsontype.EmbeddedDocument:
		var embeddedDoc Doc
		err := pc.documentDecodeValue(dc, vr, &embeddedDoc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Document(embeddedDoc)
	case bsontype.Array:
		arr := reflect.New(tArray).Elem()
		err := pc.ArrayDecodeValue(dc, vr, arr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Array(arr.Interface().(Arr))
	case bsontype.Binary:
		data, subtype, err := vr.ReadBinary()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Binary(subtype, data)
	case bsontype.Undefined:
		err := vr.ReadUndefined()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Undefined()
	case bsontype.ObjectID:
		oid, err := vr.ReadObjectID()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = ObjectID(oid)
	case bsontype.Boolean:
		b, err := vr.ReadBoolean()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Boolean(b)
	case bsontype.DateTime:
		dt, err := vr.ReadDateTime()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = DateTime(dt)
	case bsontype.Null:
		err := vr.ReadNull()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Null()
	case bsontype.Regex:
		pattern, options, err := vr.ReadRegex()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Regex(pattern, options)
	case bsontype.DBPointer:
		ns, pointer, err := vr.ReadDBPointer()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = DBPointer(ns, pointer)
	case bsontype.JavaScript:
		js, err := vr.ReadJavascript()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = JavaScript(js)
	case bsontype.Symbol:
		symbol, err := vr.ReadSymbol()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Symbol(symbol)
	case bsontype.CodeWithScope:
		code, scope, err := vr.ReadCodeWithScope()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		var doc Doc
		err = pc.decodeDocument(dc, scope, &doc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = CodeWithScope(code, doc)
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Int32(i32)
	case bsontype.Timestamp:
		t, i, err := vr.ReadTimestamp()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Timestamp(t, i)
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Int64(i64)
	case bsontype.Decimal128:
		d128, err := vr.ReadDecimal128()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = Decimal128(d128)
	case bsontype.MinKey:
		err := vr.ReadMinKey()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = MinKey()
	case bsontype.MaxKey:
		err := vr.ReadMaxKey()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		val = MaxKey()
	default:
		return fmt.Errorf("Cannot read unknown BSON type %s", vr.Type())
	***REMOVED***

	*elem = Elem***REMOVED***Key: key, Value: val***REMOVED***
	return nil
***REMOVED***

// encodeValue does not validation, and the callers must perform validation on val before calling
// this method.
func (pc PrimitiveCodecs) encodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val Val) error ***REMOVED***
	var err error
	switch val.Type() ***REMOVED***
	case bsontype.Double:
		err = vw.WriteDouble(val.Double())
	case bsontype.String:
		err = vw.WriteString(val.StringValue())
	case bsontype.EmbeddedDocument:
		var encoder bsoncodec.ValueEncoder
		encoder, err = ec.LookupEncoder(tDocument)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = encoder.EncodeValue(ec, vw, reflect.ValueOf(val.Document()))
	case bsontype.Array:
		var encoder bsoncodec.ValueEncoder
		encoder, err = ec.LookupEncoder(tArray)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = encoder.EncodeValue(ec, vw, reflect.ValueOf(val.Array()))
	case bsontype.Binary:
		// TODO: FIX THIS (╯°□°）╯︵ ┻━┻
		subtype, data := val.Binary()
		err = vw.WriteBinaryWithSubtype(data, subtype)
	case bsontype.Undefined:
		err = vw.WriteUndefined()
	case bsontype.ObjectID:
		err = vw.WriteObjectID(val.ObjectID())
	case bsontype.Boolean:
		err = vw.WriteBoolean(val.Boolean())
	case bsontype.DateTime:
		err = vw.WriteDateTime(val.DateTime())
	case bsontype.Null:
		err = vw.WriteNull()
	case bsontype.Regex:
		err = vw.WriteRegex(val.Regex())
	case bsontype.DBPointer:
		err = vw.WriteDBPointer(val.DBPointer())
	case bsontype.JavaScript:
		err = vw.WriteJavascript(val.JavaScript())
	case bsontype.Symbol:
		err = vw.WriteSymbol(val.Symbol())
	case bsontype.CodeWithScope:
		code, scope := val.CodeWithScope()

		var cwsw bsonrw.DocumentWriter
		cwsw, err = vw.WriteCodeWithScope(code)
		if err != nil ***REMOVED***
			break
		***REMOVED***

		err = pc.encodeDocument(ec, cwsw, scope)
	case bsontype.Int32:
		err = vw.WriteInt32(val.Int32())
	case bsontype.Timestamp:
		err = vw.WriteTimestamp(val.Timestamp())
	case bsontype.Int64:
		err = vw.WriteInt64(val.Int64())
	case bsontype.Decimal128:
		err = vw.WriteDecimal128(val.Decimal128())
	case bsontype.MinKey:
		err = vw.WriteMinKey()
	case bsontype.MaxKey:
		err = vw.WriteMaxKey()
	default:
		err = fmt.Errorf("%T is not a valid BSON type to encode", val.Type())
	***REMOVED***

	return err
***REMOVED***

func (pc PrimitiveCodecs) valueDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val *Val) error ***REMOVED***
	switch vr.Type() ***REMOVED***
	case bsontype.Double:
		f64, err := vr.ReadDouble()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Double(f64)
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = String(str)
	case bsontype.EmbeddedDocument:
		var embeddedDoc Doc
		err := pc.documentDecodeValue(dc, vr, &embeddedDoc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Document(embeddedDoc)
	case bsontype.Array:
		arr := reflect.New(tArray).Elem()
		err := pc.ArrayDecodeValue(dc, vr, arr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Array(arr.Interface().(Arr))
	case bsontype.Binary:
		data, subtype, err := vr.ReadBinary()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Binary(subtype, data)
	case bsontype.Undefined:
		err := vr.ReadUndefined()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Undefined()
	case bsontype.ObjectID:
		oid, err := vr.ReadObjectID()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = ObjectID(oid)
	case bsontype.Boolean:
		b, err := vr.ReadBoolean()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Boolean(b)
	case bsontype.DateTime:
		dt, err := vr.ReadDateTime()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = DateTime(dt)
	case bsontype.Null:
		err := vr.ReadNull()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Null()
	case bsontype.Regex:
		pattern, options, err := vr.ReadRegex()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Regex(pattern, options)
	case bsontype.DBPointer:
		ns, pointer, err := vr.ReadDBPointer()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = DBPointer(ns, pointer)
	case bsontype.JavaScript:
		js, err := vr.ReadJavascript()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = JavaScript(js)
	case bsontype.Symbol:
		symbol, err := vr.ReadSymbol()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Symbol(symbol)
	case bsontype.CodeWithScope:
		code, scope, err := vr.ReadCodeWithScope()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		var scopeDoc Doc
		err = pc.decodeDocument(dc, scope, &scopeDoc)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = CodeWithScope(code, scopeDoc)
	case bsontype.Int32:
		i32, err := vr.ReadInt32()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Int32(i32)
	case bsontype.Timestamp:
		t, i, err := vr.ReadTimestamp()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Timestamp(t, i)
	case bsontype.Int64:
		i64, err := vr.ReadInt64()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Int64(i64)
	case bsontype.Decimal128:
		d128, err := vr.ReadDecimal128()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = Decimal128(d128)
	case bsontype.MinKey:
		err := vr.ReadMinKey()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = MinKey()
	case bsontype.MaxKey:
		err := vr.ReadMaxKey()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		*val = MaxKey()
	default:
		return fmt.Errorf("Cannot read unknown BSON type %s", vr.Type())
	***REMOVED***

	return nil
***REMOVED***
