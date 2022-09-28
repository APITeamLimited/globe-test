// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Copier is a type that allows copying between ValueReaders, ValueWriters, and
// []byte values.
type Copier struct***REMOVED******REMOVED***

// NewCopier creates a new copier with the given registry. If a nil registry is provided
// a default registry is used.
func NewCopier() Copier ***REMOVED***
	return Copier***REMOVED******REMOVED***
***REMOVED***

// CopyDocument handles copying a document from src to dst.
func CopyDocument(dst ValueWriter, src ValueReader) error ***REMOVED***
	return Copier***REMOVED******REMOVED***.CopyDocument(dst, src)
***REMOVED***

// CopyDocument handles copying one document from the src to the dst.
func (c Copier) CopyDocument(dst ValueWriter, src ValueReader) error ***REMOVED***
	dr, err := src.ReadDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	dw, err := dst.WriteDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return c.copyDocumentCore(dw, dr)
***REMOVED***

// CopyArrayFromBytes copies the values from a BSON array represented as a
// []byte to a ValueWriter.
func (c Copier) CopyArrayFromBytes(dst ValueWriter, src []byte) error ***REMOVED***
	aw, err := dst.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = c.CopyBytesToArrayWriter(aw, src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

// CopyDocumentFromBytes copies the values from a BSON document represented as a
// []byte to a ValueWriter.
func (c Copier) CopyDocumentFromBytes(dst ValueWriter, src []byte) error ***REMOVED***
	dw, err := dst.WriteDocument()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = c.CopyBytesToDocumentWriter(dw, src)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return dw.WriteDocumentEnd()
***REMOVED***

type writeElementFn func(key string) (ValueWriter, error)

// CopyBytesToArrayWriter copies the values from a BSON Array represented as a []byte to an
// ArrayWriter.
func (c Copier) CopyBytesToArrayWriter(dst ArrayWriter, src []byte) error ***REMOVED***
	wef := func(_ string) (ValueWriter, error) ***REMOVED***
		return dst.WriteArrayElement()
	***REMOVED***

	return c.copyBytesToValueWriter(src, wef)
***REMOVED***

// CopyBytesToDocumentWriter copies the values from a BSON document represented as a []byte to a
// DocumentWriter.
func (c Copier) CopyBytesToDocumentWriter(dst DocumentWriter, src []byte) error ***REMOVED***
	wef := func(key string) (ValueWriter, error) ***REMOVED***
		return dst.WriteDocumentElement(key)
	***REMOVED***

	return c.copyBytesToValueWriter(src, wef)
***REMOVED***

func (c Copier) copyBytesToValueWriter(src []byte, wef writeElementFn) error ***REMOVED***
	// TODO(skriptble): Create errors types here. Anything thats a tag should be a property.
	length, rem, ok := bsoncore.ReadLength(src)
	if !ok ***REMOVED***
		return fmt.Errorf("couldn't read length from src, not enough bytes. length=%d", len(src))
	***REMOVED***
	if len(src) < int(length) ***REMOVED***
		return fmt.Errorf("length read exceeds number of bytes available. length=%d bytes=%d", len(src), length)
	***REMOVED***
	rem = rem[:length-4]

	var t bsontype.Type
	var key string
	var val bsoncore.Value
	for ***REMOVED***
		t, rem, ok = bsoncore.ReadType(rem)
		if !ok ***REMOVED***
			return io.EOF
		***REMOVED***
		if t == bsontype.Type(0) ***REMOVED***
			if len(rem) != 0 ***REMOVED***
				return fmt.Errorf("document end byte found before end of document. remaining bytes=%v", rem)
			***REMOVED***
			break
		***REMOVED***

		key, rem, ok = bsoncore.ReadKey(rem)
		if !ok ***REMOVED***
			return fmt.Errorf("invalid key found. remaining bytes=%v", rem)
		***REMOVED***

		// write as either array element or document element using writeElementFn
		vw, err := wef(key)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		val, rem, ok = bsoncore.ReadValue(rem, t)
		if !ok ***REMOVED***
			return fmt.Errorf("not enough bytes available to read type. bytes=%d type=%s", len(rem), t)
		***REMOVED***
		err = c.CopyValueFromBytes(vw, t, val.Data)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// CopyDocumentToBytes copies an entire document from the ValueReader and
// returns it as bytes.
func (c Copier) CopyDocumentToBytes(src ValueReader) ([]byte, error) ***REMOVED***
	return c.AppendDocumentBytes(nil, src)
***REMOVED***

// AppendDocumentBytes functions the same as CopyDocumentToBytes, but will
// append the result to dst.
func (c Copier) AppendDocumentBytes(dst []byte, src ValueReader) ([]byte, error) ***REMOVED***
	if br, ok := src.(BytesReader); ok ***REMOVED***
		_, dst, err := br.ReadValueBytes(dst)
		return dst, err
	***REMOVED***

	vw := vwPool.Get().(*valueWriter)
	defer vwPool.Put(vw)

	vw.reset(dst)

	err := c.CopyDocument(vw, src)
	dst = vw.buf
	return dst, err
***REMOVED***

// AppendArrayBytes copies an array from the ValueReader to dst.
func (c Copier) AppendArrayBytes(dst []byte, src ValueReader) ([]byte, error) ***REMOVED***
	if br, ok := src.(BytesReader); ok ***REMOVED***
		_, dst, err := br.ReadValueBytes(dst)
		return dst, err
	***REMOVED***

	vw := vwPool.Get().(*valueWriter)
	defer vwPool.Put(vw)

	vw.reset(dst)

	err := c.copyArray(vw, src)
	dst = vw.buf
	return dst, err
***REMOVED***

// CopyValueFromBytes will write the value represtend by t and src to dst.
func (c Copier) CopyValueFromBytes(dst ValueWriter, t bsontype.Type, src []byte) error ***REMOVED***
	if wvb, ok := dst.(BytesWriter); ok ***REMOVED***
		return wvb.WriteValueBytes(t, src)
	***REMOVED***

	vr := vrPool.Get().(*valueReader)
	defer vrPool.Put(vr)

	vr.reset(src)
	vr.pushElement(t)

	return c.CopyValue(dst, vr)
***REMOVED***

// CopyValueToBytes copies a value from src and returns it as a bsontype.Type and a
// []byte.
func (c Copier) CopyValueToBytes(src ValueReader) (bsontype.Type, []byte, error) ***REMOVED***
	return c.AppendValueBytes(nil, src)
***REMOVED***

// AppendValueBytes functions the same as CopyValueToBytes, but will append the
// result to dst.
func (c Copier) AppendValueBytes(dst []byte, src ValueReader) (bsontype.Type, []byte, error) ***REMOVED***
	if br, ok := src.(BytesReader); ok ***REMOVED***
		return br.ReadValueBytes(dst)
	***REMOVED***

	vw := vwPool.Get().(*valueWriter)
	defer vwPool.Put(vw)

	start := len(dst)

	vw.reset(dst)
	vw.push(mElement)

	err := c.CopyValue(vw, src)
	if err != nil ***REMOVED***
		return 0, dst, err
	***REMOVED***

	return bsontype.Type(vw.buf[start]), vw.buf[start+2:], nil
***REMOVED***

// CopyValue will copy a single value from src to dst.
func (c Copier) CopyValue(dst ValueWriter, src ValueReader) error ***REMOVED***
	var err error
	switch src.Type() ***REMOVED***
	case bsontype.Double:
		var f64 float64
		f64, err = src.ReadDouble()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteDouble(f64)
	case bsontype.String:
		var str string
		str, err = src.ReadString()
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		err = dst.WriteString(str)
	case bsontype.EmbeddedDocument:
		err = c.CopyDocument(dst, src)
	case bsontype.Array:
		err = c.copyArray(dst, src)
	case bsontype.Binary:
		var data []byte
		var subtype byte
		data, subtype, err = src.ReadBinary()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteBinaryWithSubtype(data, subtype)
	case bsontype.Undefined:
		err = src.ReadUndefined()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteUndefined()
	case bsontype.ObjectID:
		var oid primitive.ObjectID
		oid, err = src.ReadObjectID()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteObjectID(oid)
	case bsontype.Boolean:
		var b bool
		b, err = src.ReadBoolean()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteBoolean(b)
	case bsontype.DateTime:
		var dt int64
		dt, err = src.ReadDateTime()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteDateTime(dt)
	case bsontype.Null:
		err = src.ReadNull()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteNull()
	case bsontype.Regex:
		var pattern, options string
		pattern, options, err = src.ReadRegex()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteRegex(pattern, options)
	case bsontype.DBPointer:
		var ns string
		var pointer primitive.ObjectID
		ns, pointer, err = src.ReadDBPointer()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteDBPointer(ns, pointer)
	case bsontype.JavaScript:
		var js string
		js, err = src.ReadJavascript()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteJavascript(js)
	case bsontype.Symbol:
		var symbol string
		symbol, err = src.ReadSymbol()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteSymbol(symbol)
	case bsontype.CodeWithScope:
		var code string
		var srcScope DocumentReader
		code, srcScope, err = src.ReadCodeWithScope()
		if err != nil ***REMOVED***
			break
		***REMOVED***

		var dstScope DocumentWriter
		dstScope, err = dst.WriteCodeWithScope(code)
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = c.copyDocumentCore(dstScope, srcScope)
	case bsontype.Int32:
		var i32 int32
		i32, err = src.ReadInt32()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteInt32(i32)
	case bsontype.Timestamp:
		var t, i uint32
		t, i, err = src.ReadTimestamp()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteTimestamp(t, i)
	case bsontype.Int64:
		var i64 int64
		i64, err = src.ReadInt64()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteInt64(i64)
	case bsontype.Decimal128:
		var d128 primitive.Decimal128
		d128, err = src.ReadDecimal128()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteDecimal128(d128)
	case bsontype.MinKey:
		err = src.ReadMinKey()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteMinKey()
	case bsontype.MaxKey:
		err = src.ReadMaxKey()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		err = dst.WriteMaxKey()
	default:
		err = fmt.Errorf("Cannot copy unknown BSON type %s", src.Type())
	***REMOVED***

	return err
***REMOVED***

func (c Copier) copyArray(dst ValueWriter, src ValueReader) error ***REMOVED***
	ar, err := src.ReadArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	aw, err := dst.WriteArray()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for ***REMOVED***
		vr, err := ar.ReadValue()
		if err == ErrEOA ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		vw, err := aw.WriteArrayElement()
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = c.CopyValue(vw, vr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return aw.WriteArrayEnd()
***REMOVED***

func (c Copier) copyDocumentCore(dw DocumentWriter, dr DocumentReader) error ***REMOVED***
	for ***REMOVED***
		key, vr, err := dr.ReadElement()
		if err == ErrEOD ***REMOVED***
			break
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		vw, err := dw.WriteDocumentElement(key)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		err = c.CopyValue(vw, vr)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return dw.WriteDocumentEnd()
***REMOVED***
