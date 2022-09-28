// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncore

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentBuilder builds a bson document
type DocumentBuilder struct ***REMOVED***
	doc     []byte
	indexes []int32
***REMOVED***

// startDocument reserves the document's length and set the index to where the length begins
func (db *DocumentBuilder) startDocument() *DocumentBuilder ***REMOVED***
	var index int32
	index, db.doc = AppendDocumentStart(db.doc)
	db.indexes = append(db.indexes, index)
	return db
***REMOVED***

// NewDocumentBuilder creates a new DocumentBuilder
func NewDocumentBuilder() *DocumentBuilder ***REMOVED***
	return (&DocumentBuilder***REMOVED******REMOVED***).startDocument()
***REMOVED***

// Build updates the length of the document and index to the beginning of the documents length
// bytes, then returns the document (bson bytes)
func (db *DocumentBuilder) Build() Document ***REMOVED***
	last := len(db.indexes) - 1
	db.doc, _ = AppendDocumentEnd(db.doc, db.indexes[last])
	db.indexes = db.indexes[:last]
	return db.doc
***REMOVED***

// AppendInt32 will append an int32 element using key and i32 to DocumentBuilder.doc
func (db *DocumentBuilder) AppendInt32(key string, i32 int32) *DocumentBuilder ***REMOVED***
	db.doc = AppendInt32Element(db.doc, key, i32)
	return db
***REMOVED***

// AppendDocument will append a bson embedded document element using key
// and doc to DocumentBuilder.doc
func (db *DocumentBuilder) AppendDocument(key string, doc []byte) *DocumentBuilder ***REMOVED***
	db.doc = AppendDocumentElement(db.doc, key, doc)
	return db
***REMOVED***

// AppendArray will append a bson array using key and arr to DocumentBuilder.doc
func (db *DocumentBuilder) AppendArray(key string, arr []byte) *DocumentBuilder ***REMOVED***
	db.doc = AppendHeader(db.doc, bsontype.Array, key)
	db.doc = AppendArray(db.doc, arr)
	return db
***REMOVED***

// AppendDouble will append a double element using key and f to DocumentBuilder.doc
func (db *DocumentBuilder) AppendDouble(key string, f float64) *DocumentBuilder ***REMOVED***
	db.doc = AppendDoubleElement(db.doc, key, f)
	return db
***REMOVED***

// AppendString will append str to DocumentBuilder.doc with the given key
func (db *DocumentBuilder) AppendString(key string, str string) *DocumentBuilder ***REMOVED***
	db.doc = AppendStringElement(db.doc, key, str)
	return db
***REMOVED***

// AppendObjectID will append oid to DocumentBuilder.doc with the given key
func (db *DocumentBuilder) AppendObjectID(key string, oid primitive.ObjectID) *DocumentBuilder ***REMOVED***
	db.doc = AppendObjectIDElement(db.doc, key, oid)
	return db
***REMOVED***

// AppendBinary will append a BSON binary element using key, subtype, and
// b to db.doc
func (db *DocumentBuilder) AppendBinary(key string, subtype byte, b []byte) *DocumentBuilder ***REMOVED***
	db.doc = AppendBinaryElement(db.doc, key, subtype, b)
	return db
***REMOVED***

// AppendUndefined will append a BSON undefined element using key to db.doc
func (db *DocumentBuilder) AppendUndefined(key string) *DocumentBuilder ***REMOVED***
	db.doc = AppendUndefinedElement(db.doc, key)
	return db
***REMOVED***

// AppendBoolean will append a boolean element using key and b to db.doc
func (db *DocumentBuilder) AppendBoolean(key string, b bool) *DocumentBuilder ***REMOVED***
	db.doc = AppendBooleanElement(db.doc, key, b)
	return db
***REMOVED***

// AppendDateTime will append a datetime element using key and dt to db.doc
func (db *DocumentBuilder) AppendDateTime(key string, dt int64) *DocumentBuilder ***REMOVED***
	db.doc = AppendDateTimeElement(db.doc, key, dt)
	return db
***REMOVED***

// AppendNull will append a null element using key to db.doc
func (db *DocumentBuilder) AppendNull(key string) *DocumentBuilder ***REMOVED***
	db.doc = AppendNullElement(db.doc, key)
	return db
***REMOVED***

// AppendRegex will append pattern and options using key to db.doc
func (db *DocumentBuilder) AppendRegex(key, pattern, options string) *DocumentBuilder ***REMOVED***
	db.doc = AppendRegexElement(db.doc, key, pattern, options)
	return db
***REMOVED***

// AppendDBPointer will append ns and oid to using key to db.doc
func (db *DocumentBuilder) AppendDBPointer(key string, ns string, oid primitive.ObjectID) *DocumentBuilder ***REMOVED***
	db.doc = AppendDBPointerElement(db.doc, key, ns, oid)
	return db
***REMOVED***

// AppendJavaScript will append js using the provided key to db.doc
func (db *DocumentBuilder) AppendJavaScript(key, js string) *DocumentBuilder ***REMOVED***
	db.doc = AppendJavaScriptElement(db.doc, key, js)
	return db
***REMOVED***

// AppendSymbol will append a BSON symbol element using key and symbol db.doc
func (db *DocumentBuilder) AppendSymbol(key, symbol string) *DocumentBuilder ***REMOVED***
	db.doc = AppendSymbolElement(db.doc, key, symbol)
	return db
***REMOVED***

// AppendCodeWithScope will append code and scope using key to db.doc
func (db *DocumentBuilder) AppendCodeWithScope(key string, code string, scope Document) *DocumentBuilder ***REMOVED***
	db.doc = AppendCodeWithScopeElement(db.doc, key, code, scope)
	return db
***REMOVED***

// AppendTimestamp will append t and i to db.doc using provided key
func (db *DocumentBuilder) AppendTimestamp(key string, t, i uint32) *DocumentBuilder ***REMOVED***
	db.doc = AppendTimestampElement(db.doc, key, t, i)
	return db
***REMOVED***

// AppendInt64 will append i64 to dst using key to db.doc
func (db *DocumentBuilder) AppendInt64(key string, i64 int64) *DocumentBuilder ***REMOVED***
	db.doc = AppendInt64Element(db.doc, key, i64)
	return db
***REMOVED***

// AppendDecimal128 will append d128 to db.doc using provided key
func (db *DocumentBuilder) AppendDecimal128(key string, d128 primitive.Decimal128) *DocumentBuilder ***REMOVED***
	db.doc = AppendDecimal128Element(db.doc, key, d128)
	return db
***REMOVED***

// AppendMaxKey will append a max key element using key to db.doc
func (db *DocumentBuilder) AppendMaxKey(key string) *DocumentBuilder ***REMOVED***
	db.doc = AppendMaxKeyElement(db.doc, key)
	return db
***REMOVED***

// AppendMinKey will append a min key element using key to db.doc
func (db *DocumentBuilder) AppendMinKey(key string) *DocumentBuilder ***REMOVED***
	db.doc = AppendMinKeyElement(db.doc, key)
	return db
***REMOVED***

// AppendValue will append a BSON element with the provided key and value to the document.
func (db *DocumentBuilder) AppendValue(key string, val Value) *DocumentBuilder ***REMOVED***
	db.doc = AppendValueElement(db.doc, key, val)
	return db
***REMOVED***

// StartDocument starts building an inline document element with the provided key
// After this document is completed, the user must call finishDocument
func (db *DocumentBuilder) StartDocument(key string) *DocumentBuilder ***REMOVED***
	db.doc = AppendHeader(db.doc, bsontype.EmbeddedDocument, key)
	db = db.startDocument()
	return db
***REMOVED***

// FinishDocument builds the most recent document created
func (db *DocumentBuilder) FinishDocument() *DocumentBuilder ***REMOVED***
	db.doc = db.Build()
	return db
***REMOVED***
