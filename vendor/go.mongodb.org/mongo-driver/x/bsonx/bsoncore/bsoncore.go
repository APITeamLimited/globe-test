// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package bsoncore contains functions that can be used to encode and decode BSON
// elements and values to or from a slice of bytes. These functions are aimed at
// allowing low level manipulation of BSON and can be used to build a higher
// level BSON library.
//
// The Read* functions within this package return the values of the element and
// a boolean indicating if the values are valid. A boolean was used instead of
// an error because any error that would be returned would be the same: not
// enough bytes. This library attempts to do no validation, it will only return
// false if there are not enough bytes for an item to be read. For example, the
// ReadDocument function checks the length, if that length is larger than the
// number of bytes available, it will return false, if there are enough bytes, it
// will return those bytes and true. It is the consumers responsibility to
// validate those bytes.
//
// The Append* functions within this package will append the type value to the
// given dst slice. If the slice has enough capacity, it will not grow the
// slice. The Append*Element functions within this package operate in the same
// way, but additionally append the BSON type and the key before the value.
package bsoncore // import "go.mongodb.org/mongo-driver/x/bsonx/bsoncore"

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// EmptyDocumentLength is the length of a document that has been started/ended but has no elements.
	EmptyDocumentLength = 5
	// nullTerminator is a string version of the 0 byte that is appended at the end of cstrings.
	nullTerminator       = string(byte(0))
	invalidKeyPanicMsg   = "BSON element keys cannot contain null bytes"
	invalidRegexPanicMsg = "BSON regex values cannot contain null bytes"
)

// AppendType will append t to dst and return the extended buffer.
func AppendType(dst []byte, t bsontype.Type) []byte ***REMOVED*** return append(dst, byte(t)) ***REMOVED***

// AppendKey will append key to dst and return the extended buffer.
func AppendKey(dst []byte, key string) []byte ***REMOVED*** return append(dst, key+nullTerminator...) ***REMOVED***

// AppendHeader will append Type t and key to dst and return the extended
// buffer.
func AppendHeader(dst []byte, t bsontype.Type, key string) []byte ***REMOVED***
	if !isValidCString(key) ***REMOVED***
		panic(invalidKeyPanicMsg)
	***REMOVED***

	dst = AppendType(dst, t)
	dst = append(dst, key...)
	return append(dst, 0x00)
	// return append(AppendType(dst, t), key+string(0x00)...)
***REMOVED***

// TODO(skriptble): All of the Read* functions should return src resliced to start just after what was read.

// ReadType will return the first byte of the provided []byte as a type. If
// there is no available byte, false is returned.
func ReadType(src []byte) (bsontype.Type, []byte, bool) ***REMOVED***
	if len(src) < 1 ***REMOVED***
		return 0, src, false
	***REMOVED***
	return bsontype.Type(src[0]), src[1:], true
***REMOVED***

// ReadKey will read a key from src. The 0x00 byte will not be present
// in the returned string. If there are not enough bytes available, false is
// returned.
func ReadKey(src []byte) (string, []byte, bool) ***REMOVED*** return readcstring(src) ***REMOVED***

// ReadKeyBytes will read a key from src as bytes. The 0x00 byte will
// not be present in the returned string. If there are not enough bytes
// available, false is returned.
func ReadKeyBytes(src []byte) ([]byte, []byte, bool) ***REMOVED*** return readcstringbytes(src) ***REMOVED***

// ReadHeader will read a type byte and a key from src. If both of these
// values cannot be read, false is returned.
func ReadHeader(src []byte) (t bsontype.Type, key string, rem []byte, ok bool) ***REMOVED***
	t, rem, ok = ReadType(src)
	if !ok ***REMOVED***
		return 0, "", src, false
	***REMOVED***
	key, rem, ok = ReadKey(rem)
	if !ok ***REMOVED***
		return 0, "", src, false
	***REMOVED***

	return t, key, rem, true
***REMOVED***

// ReadHeaderBytes will read a type and a key from src and the remainder of the bytes
// are returned as rem. If either the type or key cannot be red, ok will be false.
func ReadHeaderBytes(src []byte) (header []byte, rem []byte, ok bool) ***REMOVED***
	if len(src) < 1 ***REMOVED***
		return nil, src, false
	***REMOVED***
	idx := bytes.IndexByte(src[1:], 0x00)
	if idx == -1 ***REMOVED***
		return nil, src, false
	***REMOVED***
	return src[:idx], src[idx+1:], true
***REMOVED***

// ReadElement reads the next full element from src. It returns the element, the remaining bytes in
// the slice, and a boolean indicating if the read was successful.
func ReadElement(src []byte) (Element, []byte, bool) ***REMOVED***
	if len(src) < 1 ***REMOVED***
		return nil, src, false
	***REMOVED***
	t := bsontype.Type(src[0])
	idx := bytes.IndexByte(src[1:], 0x00)
	if idx == -1 ***REMOVED***
		return nil, src, false
	***REMOVED***
	length, ok := valueLength(src[idx+2:], t) // We add 2 here because we called IndexByte with src[1:]
	if !ok ***REMOVED***
		return nil, src, false
	***REMOVED***
	elemLength := 1 + idx + 1 + int(length)
	if elemLength > len(src) ***REMOVED***
		return nil, src, false
	***REMOVED***
	if elemLength < 0 ***REMOVED***
		return nil, src, false
	***REMOVED***
	return src[:elemLength], src[elemLength:], true
***REMOVED***

// AppendValueElement appends value to dst as an element using key as the element's key.
func AppendValueElement(dst []byte, key string, value Value) []byte ***REMOVED***
	dst = AppendHeader(dst, value.Type, key)
	dst = append(dst, value.Data...)
	return dst
***REMOVED***

// ReadValue reads the next value as the provided types and returns a Value, the remaining bytes,
// and a boolean indicating if the read was successful.
func ReadValue(src []byte, t bsontype.Type) (Value, []byte, bool) ***REMOVED***
	data, rem, ok := readValue(src, t)
	if !ok ***REMOVED***
		return Value***REMOVED******REMOVED***, src, false
	***REMOVED***
	return Value***REMOVED***Type: t, Data: data***REMOVED***, rem, true
***REMOVED***

// AppendDouble will append f to dst and return the extended buffer.
func AppendDouble(dst []byte, f float64) []byte ***REMOVED***
	return appendu64(dst, math.Float64bits(f))
***REMOVED***

// AppendDoubleElement will append a BSON double element using key and f to dst
// and return the extended buffer.
func AppendDoubleElement(dst []byte, key string, f float64) []byte ***REMOVED***
	return AppendDouble(AppendHeader(dst, bsontype.Double, key), f)
***REMOVED***

// ReadDouble will read a float64 from src. If there are not enough bytes it
// will return false.
func ReadDouble(src []byte) (float64, []byte, bool) ***REMOVED***
	bits, src, ok := readu64(src)
	if !ok ***REMOVED***
		return 0, src, false
	***REMOVED***
	return math.Float64frombits(bits), src, true
***REMOVED***

// AppendString will append s to dst and return the extended buffer.
func AppendString(dst []byte, s string) []byte ***REMOVED***
	return appendstring(dst, s)
***REMOVED***

// AppendStringElement will append a BSON string element using key and val to dst
// and return the extended buffer.
func AppendStringElement(dst []byte, key, val string) []byte ***REMOVED***
	return AppendString(AppendHeader(dst, bsontype.String, key), val)
***REMOVED***

// ReadString will read a string from src. If there are not enough bytes it
// will return false.
func ReadString(src []byte) (string, []byte, bool) ***REMOVED***
	return readstring(src)
***REMOVED***

// AppendDocumentStart reserves a document's length and returns the index where the length begins.
// This index can later be used to write the length of the document.
func AppendDocumentStart(dst []byte) (index int32, b []byte) ***REMOVED***
	// TODO(skriptble): We really need AppendDocumentStart and AppendDocumentEnd.  AppendDocumentStart would handle calling
	// TODO ReserveLength and providing the index of the start of the document. AppendDocumentEnd would handle taking that
	// TODO start index, adding the null byte, calculating the length, and filling in the length at the start of the
	// TODO document.
	return ReserveLength(dst)
***REMOVED***

// AppendDocumentStartInline functions the same as AppendDocumentStart but takes a pointer to the
// index int32 which allows this function to be used inline.
func AppendDocumentStartInline(dst []byte, index *int32) []byte ***REMOVED***
	idx, doc := AppendDocumentStart(dst)
	*index = idx
	return doc
***REMOVED***

// AppendDocumentElementStart writes a document element header and then reserves the length bytes.
func AppendDocumentElementStart(dst []byte, key string) (index int32, b []byte) ***REMOVED***
	return AppendDocumentStart(AppendHeader(dst, bsontype.EmbeddedDocument, key))
***REMOVED***

// AppendDocumentEnd writes the null byte for a document and updates the length of the document.
// The index should be the beginning of the document's length bytes.
func AppendDocumentEnd(dst []byte, index int32) ([]byte, error) ***REMOVED***
	if int(index) > len(dst)-4 ***REMOVED***
		return dst, fmt.Errorf("not enough bytes available after index to write length")
	***REMOVED***
	dst = append(dst, 0x00)
	dst = UpdateLength(dst, index, int32(len(dst[index:])))
	return dst, nil
***REMOVED***

// AppendDocument will append doc to dst and return the extended buffer.
func AppendDocument(dst []byte, doc []byte) []byte ***REMOVED*** return append(dst, doc...) ***REMOVED***

// AppendDocumentElement will append a BSON embedded document element using key
// and doc to dst and return the extended buffer.
func AppendDocumentElement(dst []byte, key string, doc []byte) []byte ***REMOVED***
	return AppendDocument(AppendHeader(dst, bsontype.EmbeddedDocument, key), doc)
***REMOVED***

// BuildDocument will create a document with the given slice of elements and will append
// it to dst and return the extended buffer.
func BuildDocument(dst []byte, elems ...[]byte) []byte ***REMOVED***
	idx, dst := ReserveLength(dst)
	for _, elem := range elems ***REMOVED***
		dst = append(dst, elem...)
	***REMOVED***
	dst = append(dst, 0x00)
	dst = UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst
***REMOVED***

// BuildDocumentValue creates an Embedded Document value from the given elements.
func BuildDocumentValue(elems ...[]byte) Value ***REMOVED***
	return Value***REMOVED***Type: bsontype.EmbeddedDocument, Data: BuildDocument(nil, elems...)***REMOVED***
***REMOVED***

// BuildDocumentElement will append a BSON embedded document elemnt using key and the provided
// elements and return the extended buffer.
func BuildDocumentElement(dst []byte, key string, elems ...[]byte) []byte ***REMOVED***
	return BuildDocument(AppendHeader(dst, bsontype.EmbeddedDocument, key), elems...)
***REMOVED***

// BuildDocumentFromElements is an alaias for the BuildDocument function.
var BuildDocumentFromElements = BuildDocument

// ReadDocument will read a document from src. If there are not enough bytes it
// will return false.
func ReadDocument(src []byte) (doc Document, rem []byte, ok bool) ***REMOVED*** return readLengthBytes(src) ***REMOVED***

// AppendArrayStart appends the length bytes to an array and then returns the index of the start
// of those length bytes.
func AppendArrayStart(dst []byte) (index int32, b []byte) ***REMOVED*** return ReserveLength(dst) ***REMOVED***

// AppendArrayElementStart appends an array element header and then the length bytes for an array,
// returning the index where the length starts.
func AppendArrayElementStart(dst []byte, key string) (index int32, b []byte) ***REMOVED***
	return AppendArrayStart(AppendHeader(dst, bsontype.Array, key))
***REMOVED***

// AppendArrayEnd appends the null byte to an array and calculates the length, inserting that
// calculated length starting at index.
func AppendArrayEnd(dst []byte, index int32) ([]byte, error) ***REMOVED*** return AppendDocumentEnd(dst, index) ***REMOVED***

// AppendArray will append arr to dst and return the extended buffer.
func AppendArray(dst []byte, arr []byte) []byte ***REMOVED*** return append(dst, arr...) ***REMOVED***

// AppendArrayElement will append a BSON array element using key and arr to dst
// and return the extended buffer.
func AppendArrayElement(dst []byte, key string, arr []byte) []byte ***REMOVED***
	return AppendArray(AppendHeader(dst, bsontype.Array, key), arr)
***REMOVED***

// BuildArray will append a BSON array to dst built from values.
func BuildArray(dst []byte, values ...Value) []byte ***REMOVED***
	idx, dst := ReserveLength(dst)
	for pos, val := range values ***REMOVED***
		dst = AppendValueElement(dst, strconv.Itoa(pos), val)
	***REMOVED***
	dst = append(dst, 0x00)
	dst = UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst
***REMOVED***

// BuildArrayElement will create an array element using the provided values.
func BuildArrayElement(dst []byte, key string, values ...Value) []byte ***REMOVED***
	return BuildArray(AppendHeader(dst, bsontype.Array, key), values...)
***REMOVED***

// ReadArray will read an array from src. If there are not enough bytes it
// will return false.
func ReadArray(src []byte) (arr Array, rem []byte, ok bool) ***REMOVED*** return readLengthBytes(src) ***REMOVED***

// AppendBinary will append subtype and b to dst and return the extended buffer.
func AppendBinary(dst []byte, subtype byte, b []byte) []byte ***REMOVED***
	if subtype == 0x02 ***REMOVED***
		return appendBinarySubtype2(dst, subtype, b)
	***REMOVED***
	dst = append(appendLength(dst, int32(len(b))), subtype)
	return append(dst, b...)
***REMOVED***

// AppendBinaryElement will append a BSON binary element using key, subtype, and
// b to dst and return the extended buffer.
func AppendBinaryElement(dst []byte, key string, subtype byte, b []byte) []byte ***REMOVED***
	return AppendBinary(AppendHeader(dst, bsontype.Binary, key), subtype, b)
***REMOVED***

// ReadBinary will read a subtype and bin from src. If there are not enough bytes it
// will return false.
func ReadBinary(src []byte) (subtype byte, bin []byte, rem []byte, ok bool) ***REMOVED***
	length, rem, ok := ReadLength(src)
	if !ok ***REMOVED***
		return 0x00, nil, src, false
	***REMOVED***
	if len(rem) < 1 ***REMOVED*** // subtype
		return 0x00, nil, src, false
	***REMOVED***
	subtype, rem = rem[0], rem[1:]

	if len(rem) < int(length) ***REMOVED***
		return 0x00, nil, src, false
	***REMOVED***

	if subtype == 0x02 ***REMOVED***
		length, rem, ok = ReadLength(rem)
		if !ok || len(rem) < int(length) ***REMOVED***
			return 0x00, nil, src, false
		***REMOVED***
	***REMOVED***

	return subtype, rem[:length], rem[length:], true
***REMOVED***

// AppendUndefinedElement will append a BSON undefined element using key to dst
// and return the extended buffer.
func AppendUndefinedElement(dst []byte, key string) []byte ***REMOVED***
	return AppendHeader(dst, bsontype.Undefined, key)
***REMOVED***

// AppendObjectID will append oid to dst and return the extended buffer.
func AppendObjectID(dst []byte, oid primitive.ObjectID) []byte ***REMOVED*** return append(dst, oid[:]...) ***REMOVED***

// AppendObjectIDElement will append a BSON ObjectID element using key and oid to dst
// and return the extended buffer.
func AppendObjectIDElement(dst []byte, key string, oid primitive.ObjectID) []byte ***REMOVED***
	return AppendObjectID(AppendHeader(dst, bsontype.ObjectID, key), oid)
***REMOVED***

// ReadObjectID will read an ObjectID from src. If there are not enough bytes it
// will return false.
func ReadObjectID(src []byte) (primitive.ObjectID, []byte, bool) ***REMOVED***
	if len(src) < 12 ***REMOVED***
		return primitive.ObjectID***REMOVED******REMOVED***, src, false
	***REMOVED***
	var oid primitive.ObjectID
	copy(oid[:], src[0:12])
	return oid, src[12:], true
***REMOVED***

// AppendBoolean will append b to dst and return the extended buffer.
func AppendBoolean(dst []byte, b bool) []byte ***REMOVED***
	if b ***REMOVED***
		return append(dst, 0x01)
	***REMOVED***
	return append(dst, 0x00)
***REMOVED***

// AppendBooleanElement will append a BSON boolean element using key and b to dst
// and return the extended buffer.
func AppendBooleanElement(dst []byte, key string, b bool) []byte ***REMOVED***
	return AppendBoolean(AppendHeader(dst, bsontype.Boolean, key), b)
***REMOVED***

// ReadBoolean will read a bool from src. If there are not enough bytes it
// will return false.
func ReadBoolean(src []byte) (bool, []byte, bool) ***REMOVED***
	if len(src) < 1 ***REMOVED***
		return false, src, false
	***REMOVED***

	return src[0] == 0x01, src[1:], true
***REMOVED***

// AppendDateTime will append dt to dst and return the extended buffer.
func AppendDateTime(dst []byte, dt int64) []byte ***REMOVED*** return appendi64(dst, dt) ***REMOVED***

// AppendDateTimeElement will append a BSON datetime element using key and dt to dst
// and return the extended buffer.
func AppendDateTimeElement(dst []byte, key string, dt int64) []byte ***REMOVED***
	return AppendDateTime(AppendHeader(dst, bsontype.DateTime, key), dt)
***REMOVED***

// ReadDateTime will read an int64 datetime from src. If there are not enough bytes it
// will return false.
func ReadDateTime(src []byte) (int64, []byte, bool) ***REMOVED*** return readi64(src) ***REMOVED***

// AppendTime will append time as a BSON DateTime to dst and return the extended buffer.
func AppendTime(dst []byte, t time.Time) []byte ***REMOVED***
	return AppendDateTime(dst, t.Unix()*1000+int64(t.Nanosecond()/1e6))
***REMOVED***

// AppendTimeElement will append a BSON datetime element using key and dt to dst
// and return the extended buffer.
func AppendTimeElement(dst []byte, key string, t time.Time) []byte ***REMOVED***
	return AppendTime(AppendHeader(dst, bsontype.DateTime, key), t)
***REMOVED***

// ReadTime will read an time.Time datetime from src. If there are not enough bytes it
// will return false.
func ReadTime(src []byte) (time.Time, []byte, bool) ***REMOVED***
	dt, rem, ok := readi64(src)
	return time.Unix(dt/1e3, dt%1e3*1e6), rem, ok
***REMOVED***

// AppendNullElement will append a BSON null element using key to dst
// and return the extended buffer.
func AppendNullElement(dst []byte, key string) []byte ***REMOVED*** return AppendHeader(dst, bsontype.Null, key) ***REMOVED***

// AppendRegex will append pattern and options to dst and return the extended buffer.
func AppendRegex(dst []byte, pattern, options string) []byte ***REMOVED***
	if !isValidCString(pattern) || !isValidCString(options) ***REMOVED***
		panic(invalidRegexPanicMsg)
	***REMOVED***

	return append(dst, pattern+nullTerminator+options+nullTerminator...)
***REMOVED***

// AppendRegexElement will append a BSON regex element using key, pattern, and
// options to dst and return the extended buffer.
func AppendRegexElement(dst []byte, key, pattern, options string) []byte ***REMOVED***
	return AppendRegex(AppendHeader(dst, bsontype.Regex, key), pattern, options)
***REMOVED***

// ReadRegex will read a pattern and options from src. If there are not enough bytes it
// will return false.
func ReadRegex(src []byte) (pattern, options string, rem []byte, ok bool) ***REMOVED***
	pattern, rem, ok = readcstring(src)
	if !ok ***REMOVED***
		return "", "", src, false
	***REMOVED***
	options, rem, ok = readcstring(rem)
	if !ok ***REMOVED***
		return "", "", src, false
	***REMOVED***
	return pattern, options, rem, true
***REMOVED***

// AppendDBPointer will append ns and oid to dst and return the extended buffer.
func AppendDBPointer(dst []byte, ns string, oid primitive.ObjectID) []byte ***REMOVED***
	return append(appendstring(dst, ns), oid[:]...)
***REMOVED***

// AppendDBPointerElement will append a BSON DBPointer element using key, ns,
// and oid to dst and return the extended buffer.
func AppendDBPointerElement(dst []byte, key, ns string, oid primitive.ObjectID) []byte ***REMOVED***
	return AppendDBPointer(AppendHeader(dst, bsontype.DBPointer, key), ns, oid)
***REMOVED***

// ReadDBPointer will read a ns and oid from src. If there are not enough bytes it
// will return false.
func ReadDBPointer(src []byte) (ns string, oid primitive.ObjectID, rem []byte, ok bool) ***REMOVED***
	ns, rem, ok = readstring(src)
	if !ok ***REMOVED***
		return "", primitive.ObjectID***REMOVED******REMOVED***, src, false
	***REMOVED***
	oid, rem, ok = ReadObjectID(rem)
	if !ok ***REMOVED***
		return "", primitive.ObjectID***REMOVED******REMOVED***, src, false
	***REMOVED***
	return ns, oid, rem, true
***REMOVED***

// AppendJavaScript will append js to dst and return the extended buffer.
func AppendJavaScript(dst []byte, js string) []byte ***REMOVED*** return appendstring(dst, js) ***REMOVED***

// AppendJavaScriptElement will append a BSON JavaScript element using key and
// js to dst and return the extended buffer.
func AppendJavaScriptElement(dst []byte, key, js string) []byte ***REMOVED***
	return AppendJavaScript(AppendHeader(dst, bsontype.JavaScript, key), js)
***REMOVED***

// ReadJavaScript will read a js string from src. If there are not enough bytes it
// will return false.
func ReadJavaScript(src []byte) (js string, rem []byte, ok bool) ***REMOVED*** return readstring(src) ***REMOVED***

// AppendSymbol will append symbol to dst and return the extended buffer.
func AppendSymbol(dst []byte, symbol string) []byte ***REMOVED*** return appendstring(dst, symbol) ***REMOVED***

// AppendSymbolElement will append a BSON symbol element using key and symbol to dst
// and return the extended buffer.
func AppendSymbolElement(dst []byte, key, symbol string) []byte ***REMOVED***
	return AppendSymbol(AppendHeader(dst, bsontype.Symbol, key), symbol)
***REMOVED***

// ReadSymbol will read a symbol string from src. If there are not enough bytes it
// will return false.
func ReadSymbol(src []byte) (symbol string, rem []byte, ok bool) ***REMOVED*** return readstring(src) ***REMOVED***

// AppendCodeWithScope will append code and scope to dst and return the extended buffer.
func AppendCodeWithScope(dst []byte, code string, scope []byte) []byte ***REMOVED***
	length := int32(4 + 4 + len(code) + 1 + len(scope)) // length of cws, length of code, code, 0x00, scope
	dst = appendLength(dst, length)

	return append(appendstring(dst, code), scope...)
***REMOVED***

// AppendCodeWithScopeElement will append a BSON code with scope element using
// key, code, and scope to dst
// and return the extended buffer.
func AppendCodeWithScopeElement(dst []byte, key, code string, scope []byte) []byte ***REMOVED***
	return AppendCodeWithScope(AppendHeader(dst, bsontype.CodeWithScope, key), code, scope)
***REMOVED***

// ReadCodeWithScope will read code and scope from src. If there are not enough bytes it
// will return false.
func ReadCodeWithScope(src []byte) (code string, scope []byte, rem []byte, ok bool) ***REMOVED***
	length, rem, ok := ReadLength(src)
	if !ok || len(src) < int(length) ***REMOVED***
		return "", nil, src, false
	***REMOVED***

	code, rem, ok = readstring(rem)
	if !ok ***REMOVED***
		return "", nil, src, false
	***REMOVED***

	scope, rem, ok = ReadDocument(rem)
	if !ok ***REMOVED***
		return "", nil, src, false
	***REMOVED***
	return code, scope, rem, true
***REMOVED***

// AppendInt32 will append i32 to dst and return the extended buffer.
func AppendInt32(dst []byte, i32 int32) []byte ***REMOVED*** return appendi32(dst, i32) ***REMOVED***

// AppendInt32Element will append a BSON int32 element using key and i32 to dst
// and return the extended buffer.
func AppendInt32Element(dst []byte, key string, i32 int32) []byte ***REMOVED***
	return AppendInt32(AppendHeader(dst, bsontype.Int32, key), i32)
***REMOVED***

// ReadInt32 will read an int32 from src. If there are not enough bytes it
// will return false.
func ReadInt32(src []byte) (int32, []byte, bool) ***REMOVED*** return readi32(src) ***REMOVED***

// AppendTimestamp will append t and i to dst and return the extended buffer.
func AppendTimestamp(dst []byte, t, i uint32) []byte ***REMOVED***
	return appendu32(appendu32(dst, i), t) // i is the lower 4 bytes, t is the higher 4 bytes
***REMOVED***

// AppendTimestampElement will append a BSON timestamp element using key, t, and
// i to dst and return the extended buffer.
func AppendTimestampElement(dst []byte, key string, t, i uint32) []byte ***REMOVED***
	return AppendTimestamp(AppendHeader(dst, bsontype.Timestamp, key), t, i)
***REMOVED***

// ReadTimestamp will read t and i from src. If there are not enough bytes it
// will return false.
func ReadTimestamp(src []byte) (t, i uint32, rem []byte, ok bool) ***REMOVED***
	i, rem, ok = readu32(src)
	if !ok ***REMOVED***
		return 0, 0, src, false
	***REMOVED***
	t, rem, ok = readu32(rem)
	if !ok ***REMOVED***
		return 0, 0, src, false
	***REMOVED***
	return t, i, rem, true
***REMOVED***

// AppendInt64 will append i64 to dst and return the extended buffer.
func AppendInt64(dst []byte, i64 int64) []byte ***REMOVED*** return appendi64(dst, i64) ***REMOVED***

// AppendInt64Element will append a BSON int64 element using key and i64 to dst
// and return the extended buffer.
func AppendInt64Element(dst []byte, key string, i64 int64) []byte ***REMOVED***
	return AppendInt64(AppendHeader(dst, bsontype.Int64, key), i64)
***REMOVED***

// ReadInt64 will read an int64 from src. If there are not enough bytes it
// will return false.
func ReadInt64(src []byte) (int64, []byte, bool) ***REMOVED*** return readi64(src) ***REMOVED***

// AppendDecimal128 will append d128 to dst and return the extended buffer.
func AppendDecimal128(dst []byte, d128 primitive.Decimal128) []byte ***REMOVED***
	high, low := d128.GetBytes()
	return appendu64(appendu64(dst, low), high)
***REMOVED***

// AppendDecimal128Element will append a BSON primitive.28 element using key and
// d128 to dst and return the extended buffer.
func AppendDecimal128Element(dst []byte, key string, d128 primitive.Decimal128) []byte ***REMOVED***
	return AppendDecimal128(AppendHeader(dst, bsontype.Decimal128, key), d128)
***REMOVED***

// ReadDecimal128 will read a primitive.Decimal128 from src. If there are not enough bytes it
// will return false.
func ReadDecimal128(src []byte) (primitive.Decimal128, []byte, bool) ***REMOVED***
	l, rem, ok := readu64(src)
	if !ok ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, src, false
	***REMOVED***

	h, rem, ok := readu64(rem)
	if !ok ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, src, false
	***REMOVED***

	return primitive.NewDecimal128(h, l), rem, true
***REMOVED***

// AppendMaxKeyElement will append a BSON max key element using key to dst
// and return the extended buffer.
func AppendMaxKeyElement(dst []byte, key string) []byte ***REMOVED***
	return AppendHeader(dst, bsontype.MaxKey, key)
***REMOVED***

// AppendMinKeyElement will append a BSON min key element using key to dst
// and return the extended buffer.
func AppendMinKeyElement(dst []byte, key string) []byte ***REMOVED***
	return AppendHeader(dst, bsontype.MinKey, key)
***REMOVED***

// EqualValue will return true if the two values are equal.
func EqualValue(t1, t2 bsontype.Type, v1, v2 []byte) bool ***REMOVED***
	if t1 != t2 ***REMOVED***
		return false
	***REMOVED***
	v1, _, ok := readValue(v1, t1)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	v2, _, ok = readValue(v2, t2)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return bytes.Equal(v1, v2)
***REMOVED***

// valueLength will determine the length of the next value contained in src as if it
// is type t. The returned bool will be false if there are not enough bytes in src for
// a value of type t.
func valueLength(src []byte, t bsontype.Type) (int32, bool) ***REMOVED***
	var length int32
	ok := true
	switch t ***REMOVED***
	case bsontype.Array, bsontype.EmbeddedDocument, bsontype.CodeWithScope:
		length, _, ok = ReadLength(src)
	case bsontype.Binary:
		length, _, ok = ReadLength(src)
		length += 4 + 1 // binary length + subtype byte
	case bsontype.Boolean:
		length = 1
	case bsontype.DBPointer:
		length, _, ok = ReadLength(src)
		length += 4 + 12 // string length + ObjectID length
	case bsontype.DateTime, bsontype.Double, bsontype.Int64, bsontype.Timestamp:
		length = 8
	case bsontype.Decimal128:
		length = 16
	case bsontype.Int32:
		length = 4
	case bsontype.JavaScript, bsontype.String, bsontype.Symbol:
		length, _, ok = ReadLength(src)
		length += 4
	case bsontype.MaxKey, bsontype.MinKey, bsontype.Null, bsontype.Undefined:
		length = 0
	case bsontype.ObjectID:
		length = 12
	case bsontype.Regex:
		regex := bytes.IndexByte(src, 0x00)
		if regex < 0 ***REMOVED***
			ok = false
			break
		***REMOVED***
		pattern := bytes.IndexByte(src[regex+1:], 0x00)
		if pattern < 0 ***REMOVED***
			ok = false
			break
		***REMOVED***
		length = int32(int64(regex) + 1 + int64(pattern) + 1)
	default:
		ok = false
	***REMOVED***

	return length, ok
***REMOVED***

func readValue(src []byte, t bsontype.Type) ([]byte, []byte, bool) ***REMOVED***
	length, ok := valueLength(src, t)
	if !ok || int(length) > len(src) ***REMOVED***
		return nil, src, false
	***REMOVED***

	return src[:length], src[length:], true
***REMOVED***

// ReserveLength reserves the space required for length and returns the index where to write the length
// and the []byte with reserved space.
func ReserveLength(dst []byte) (int32, []byte) ***REMOVED***
	index := len(dst)
	return int32(index), append(dst, 0x00, 0x00, 0x00, 0x00)
***REMOVED***

// UpdateLength updates the length at index with length and returns the []byte.
func UpdateLength(dst []byte, index, length int32) []byte ***REMOVED***
	dst[index] = byte(length)
	dst[index+1] = byte(length >> 8)
	dst[index+2] = byte(length >> 16)
	dst[index+3] = byte(length >> 24)
	return dst
***REMOVED***

func appendLength(dst []byte, l int32) []byte ***REMOVED*** return appendi32(dst, l) ***REMOVED***

func appendi32(dst []byte, i32 int32) []byte ***REMOVED***
	return append(dst, byte(i32), byte(i32>>8), byte(i32>>16), byte(i32>>24))
***REMOVED***

// ReadLength reads an int32 length from src and returns the length and the remaining bytes. If
// there aren't enough bytes to read a valid length, src is returned unomdified and the returned
// bool will be false.
func ReadLength(src []byte) (int32, []byte, bool) ***REMOVED***
	ln, src, ok := readi32(src)
	if ln < 0 ***REMOVED***
		return ln, src, false
	***REMOVED***
	return ln, src, ok
***REMOVED***

func readi32(src []byte) (int32, []byte, bool) ***REMOVED***
	if len(src) < 4 ***REMOVED***
		return 0, src, false
	***REMOVED***
	return (int32(src[0]) | int32(src[1])<<8 | int32(src[2])<<16 | int32(src[3])<<24), src[4:], true
***REMOVED***

func appendi64(dst []byte, i64 int64) []byte ***REMOVED***
	return append(dst,
		byte(i64), byte(i64>>8), byte(i64>>16), byte(i64>>24),
		byte(i64>>32), byte(i64>>40), byte(i64>>48), byte(i64>>56),
	)
***REMOVED***

func readi64(src []byte) (int64, []byte, bool) ***REMOVED***
	if len(src) < 8 ***REMOVED***
		return 0, src, false
	***REMOVED***
	i64 := (int64(src[0]) | int64(src[1])<<8 | int64(src[2])<<16 | int64(src[3])<<24 |
		int64(src[4])<<32 | int64(src[5])<<40 | int64(src[6])<<48 | int64(src[7])<<56)
	return i64, src[8:], true
***REMOVED***

func appendu32(dst []byte, u32 uint32) []byte ***REMOVED***
	return append(dst, byte(u32), byte(u32>>8), byte(u32>>16), byte(u32>>24))
***REMOVED***

func readu32(src []byte) (uint32, []byte, bool) ***REMOVED***
	if len(src) < 4 ***REMOVED***
		return 0, src, false
	***REMOVED***

	return (uint32(src[0]) | uint32(src[1])<<8 | uint32(src[2])<<16 | uint32(src[3])<<24), src[4:], true
***REMOVED***

func appendu64(dst []byte, u64 uint64) []byte ***REMOVED***
	return append(dst,
		byte(u64), byte(u64>>8), byte(u64>>16), byte(u64>>24),
		byte(u64>>32), byte(u64>>40), byte(u64>>48), byte(u64>>56),
	)
***REMOVED***

func readu64(src []byte) (uint64, []byte, bool) ***REMOVED***
	if len(src) < 8 ***REMOVED***
		return 0, src, false
	***REMOVED***
	u64 := (uint64(src[0]) | uint64(src[1])<<8 | uint64(src[2])<<16 | uint64(src[3])<<24 |
		uint64(src[4])<<32 | uint64(src[5])<<40 | uint64(src[6])<<48 | uint64(src[7])<<56)
	return u64, src[8:], true
***REMOVED***

// keep in sync with readcstringbytes
func readcstring(src []byte) (string, []byte, bool) ***REMOVED***
	idx := bytes.IndexByte(src, 0x00)
	if idx < 0 ***REMOVED***
		return "", src, false
	***REMOVED***
	return string(src[:idx]), src[idx+1:], true
***REMOVED***

// keep in sync with readcstring
func readcstringbytes(src []byte) ([]byte, []byte, bool) ***REMOVED***
	idx := bytes.IndexByte(src, 0x00)
	if idx < 0 ***REMOVED***
		return nil, src, false
	***REMOVED***
	return src[:idx], src[idx+1:], true
***REMOVED***

func appendstring(dst []byte, s string) []byte ***REMOVED***
	l := int32(len(s) + 1)
	dst = appendLength(dst, l)
	dst = append(dst, s...)
	return append(dst, 0x00)
***REMOVED***

func readstring(src []byte) (string, []byte, bool) ***REMOVED***
	l, rem, ok := ReadLength(src)
	if !ok ***REMOVED***
		return "", src, false
	***REMOVED***
	if len(src[4:]) < int(l) || l == 0 ***REMOVED***
		return "", src, false
	***REMOVED***

	return string(rem[:l-1]), rem[l:], true
***REMOVED***

// readLengthBytes attempts to read a length and that number of bytes. This
// function requires that the length include the four bytes for itself.
func readLengthBytes(src []byte) ([]byte, []byte, bool) ***REMOVED***
	l, _, ok := ReadLength(src)
	if !ok ***REMOVED***
		return nil, src, false
	***REMOVED***
	if len(src) < int(l) ***REMOVED***
		return nil, src, false
	***REMOVED***
	return src[:l], src[l:], true
***REMOVED***

func appendBinarySubtype2(dst []byte, subtype byte, b []byte) []byte ***REMOVED***
	dst = appendLength(dst, int32(len(b)+4)) // The bytes we'll encode need to be 4 larger for the length bytes
	dst = append(dst, subtype)
	dst = appendLength(dst, int32(len(b)))
	return append(dst, b...)
***REMOVED***

func isValidCString(cs string) bool ***REMOVED***
	return !strings.ContainsRune(cs, '\x00')
***REMOVED***
