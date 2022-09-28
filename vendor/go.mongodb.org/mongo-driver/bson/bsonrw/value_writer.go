// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var _ ValueWriter = (*valueWriter)(nil)

var vwPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return new(valueWriter)
	***REMOVED***,
***REMOVED***

// BSONValueWriterPool is a pool for BSON ValueWriters.
type BSONValueWriterPool struct ***REMOVED***
	pool sync.Pool
***REMOVED***

// NewBSONValueWriterPool creates a new pool for ValueWriter instances that write to BSON.
func NewBSONValueWriterPool() *BSONValueWriterPool ***REMOVED***
	return &BSONValueWriterPool***REMOVED***
		pool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED***
				return new(valueWriter)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Get retrieves a BSON ValueWriter from the pool and resets it to use w as the destination.
func (bvwp *BSONValueWriterPool) Get(w io.Writer) ValueWriter ***REMOVED***
	vw := bvwp.pool.Get().(*valueWriter)

	// TODO: Having to call reset here with the same buffer doesn't really make sense.
	vw.reset(vw.buf)
	vw.buf = vw.buf[:0]
	vw.w = w
	return vw
***REMOVED***

// GetAtModeElement retrieves a ValueWriterFlusher from the pool and resets it to use w as the destination.
func (bvwp *BSONValueWriterPool) GetAtModeElement(w io.Writer) ValueWriterFlusher ***REMOVED***
	vw := bvwp.Get(w).(*valueWriter)
	vw.push(mElement)
	return vw
***REMOVED***

// Put inserts a ValueWriter into the pool. If the ValueWriter is not a BSON ValueWriter, nothing
// happens and ok will be false.
func (bvwp *BSONValueWriterPool) Put(vw ValueWriter) (ok bool) ***REMOVED***
	bvw, ok := vw.(*valueWriter)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	bvwp.pool.Put(bvw)
	return true
***REMOVED***

// This is here so that during testing we can change it and not require
// allocating a 4GB slice.
var maxSize = math.MaxInt32

var errNilWriter = errors.New("cannot create a ValueWriter from a nil io.Writer")

type errMaxDocumentSizeExceeded struct ***REMOVED***
	size int64
***REMOVED***

func (mdse errMaxDocumentSizeExceeded) Error() string ***REMOVED***
	return fmt.Sprintf("document size (%d) is larger than the max int32", mdse.size)
***REMOVED***

type vwMode int

const (
	_ vwMode = iota
	vwTopLevel
	vwDocument
	vwArray
	vwValue
	vwElement
	vwCodeWithScope
)

func (vm vwMode) String() string ***REMOVED***
	var str string

	switch vm ***REMOVED***
	case vwTopLevel:
		str = "TopLevel"
	case vwDocument:
		str = "DocumentMode"
	case vwArray:
		str = "ArrayMode"
	case vwValue:
		str = "ValueMode"
	case vwElement:
		str = "ElementMode"
	case vwCodeWithScope:
		str = "CodeWithScopeMode"
	default:
		str = "UnknownMode"
	***REMOVED***

	return str
***REMOVED***

type vwState struct ***REMOVED***
	mode   mode
	key    string
	arrkey int
	start  int32
***REMOVED***

type valueWriter struct ***REMOVED***
	w   io.Writer
	buf []byte

	stack []vwState
	frame int64
***REMOVED***

func (vw *valueWriter) advanceFrame() ***REMOVED***
	if vw.frame+1 >= int64(len(vw.stack)) ***REMOVED*** // We need to grow the stack
		length := len(vw.stack)
		if length+1 >= cap(vw.stack) ***REMOVED***
			// double it
			buf := make([]vwState, 2*cap(vw.stack)+1)
			copy(buf, vw.stack)
			vw.stack = buf
		***REMOVED***
		vw.stack = vw.stack[:length+1]
	***REMOVED***
	vw.frame++
***REMOVED***

func (vw *valueWriter) push(m mode) ***REMOVED***
	vw.advanceFrame()

	// Clean the stack
	vw.stack[vw.frame].mode = m
	vw.stack[vw.frame].key = ""
	vw.stack[vw.frame].arrkey = 0
	vw.stack[vw.frame].start = 0

	vw.stack[vw.frame].mode = m
	switch m ***REMOVED***
	case mDocument, mArray, mCodeWithScope:
		vw.reserveLength()
	***REMOVED***
***REMOVED***

func (vw *valueWriter) reserveLength() ***REMOVED***
	vw.stack[vw.frame].start = int32(len(vw.buf))
	vw.buf = append(vw.buf, 0x00, 0x00, 0x00, 0x00)
***REMOVED***

func (vw *valueWriter) pop() ***REMOVED***
	switch vw.stack[vw.frame].mode ***REMOVED***
	case mElement, mValue:
		vw.frame--
	case mDocument, mArray, mCodeWithScope:
		vw.frame -= 2 // we pop twice to jump over the mElement: mDocument -> mElement -> mDocument/mTopLevel/etc...
	***REMOVED***
***REMOVED***

// NewBSONValueWriter creates a ValueWriter that writes BSON to w.
//
// This ValueWriter will only write entire documents to the io.Writer and it
// will buffer the document as it is built.
func NewBSONValueWriter(w io.Writer) (ValueWriter, error) ***REMOVED***
	if w == nil ***REMOVED***
		return nil, errNilWriter
	***REMOVED***
	return newValueWriter(w), nil
***REMOVED***

func newValueWriter(w io.Writer) *valueWriter ***REMOVED***
	vw := new(valueWriter)
	stack := make([]vwState, 1, 5)
	stack[0] = vwState***REMOVED***mode: mTopLevel***REMOVED***
	vw.w = w
	vw.stack = stack

	return vw
***REMOVED***

func newValueWriterFromSlice(buf []byte) *valueWriter ***REMOVED***
	vw := new(valueWriter)
	stack := make([]vwState, 1, 5)
	stack[0] = vwState***REMOVED***mode: mTopLevel***REMOVED***
	vw.stack = stack
	vw.buf = buf

	return vw
***REMOVED***

func (vw *valueWriter) reset(buf []byte) ***REMOVED***
	if vw.stack == nil ***REMOVED***
		vw.stack = make([]vwState, 1, 5)
	***REMOVED***
	vw.stack = vw.stack[:1]
	vw.stack[0] = vwState***REMOVED***mode: mTopLevel***REMOVED***
	vw.buf = buf
	vw.frame = 0
	vw.w = nil
***REMOVED***

func (vw *valueWriter) invalidTransitionError(destination mode, name string, modes []mode) error ***REMOVED***
	te := TransitionError***REMOVED***
		name:        name,
		current:     vw.stack[vw.frame].mode,
		destination: destination,
		modes:       modes,
		action:      "write",
	***REMOVED***
	if vw.frame != 0 ***REMOVED***
		te.parent = vw.stack[vw.frame-1].mode
	***REMOVED***
	return te
***REMOVED***

func (vw *valueWriter) writeElementHeader(t bsontype.Type, destination mode, callerName string, addmodes ...mode) error ***REMOVED***
	switch vw.stack[vw.frame].mode ***REMOVED***
	case mElement:
		key := vw.stack[vw.frame].key
		if !isValidCString(key) ***REMOVED***
			return errors.New("BSON element key cannot contain null bytes")
		***REMOVED***

		vw.buf = bsoncore.AppendHeader(vw.buf, t, key)
	case mValue:
		// TODO: Do this with a cache of the first 1000 or so array keys.
		vw.buf = bsoncore.AppendHeader(vw.buf, t, strconv.Itoa(vw.stack[vw.frame].arrkey))
	default:
		modes := []mode***REMOVED***mElement, mValue***REMOVED***
		if addmodes != nil ***REMOVED***
			modes = append(modes, addmodes...)
		***REMOVED***
		return vw.invalidTransitionError(destination, callerName, modes)
	***REMOVED***

	return nil
***REMOVED***

func (vw *valueWriter) WriteValueBytes(t bsontype.Type, b []byte) error ***REMOVED***
	if err := vw.writeElementHeader(t, mode(0), "WriteValueBytes"); err != nil ***REMOVED***
		return err
	***REMOVED***
	vw.buf = append(vw.buf, b...)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteArray() (ArrayWriter, error) ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Array, mArray, "WriteArray"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vw.push(mArray)

	return vw, nil
***REMOVED***

func (vw *valueWriter) WriteBinary(b []byte) error ***REMOVED***
	return vw.WriteBinaryWithSubtype(b, 0x00)
***REMOVED***

func (vw *valueWriter) WriteBinaryWithSubtype(b []byte, btype byte) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Binary, mode(0), "WriteBinaryWithSubtype"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendBinary(vw.buf, btype, b)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteBoolean(b bool) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Boolean, mode(0), "WriteBoolean"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendBoolean(vw.buf, b)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteCodeWithScope(code string) (DocumentWriter, error) ***REMOVED***
	if err := vw.writeElementHeader(bsontype.CodeWithScope, mCodeWithScope, "WriteCodeWithScope"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// CodeWithScope is a different than other types because we need an extra
	// frame on the stack. In the EndDocument code, we write the document
	// length, pop, write the code with scope length, and pop. To simplify the
	// pop code, we push a spacer frame that we'll always jump over.
	vw.push(mCodeWithScope)
	vw.buf = bsoncore.AppendString(vw.buf, code)
	vw.push(mSpacer)
	vw.push(mDocument)

	return vw, nil
***REMOVED***

func (vw *valueWriter) WriteDBPointer(ns string, oid primitive.ObjectID) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.DBPointer, mode(0), "WriteDBPointer"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendDBPointer(vw.buf, ns, oid)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteDateTime(dt int64) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.DateTime, mode(0), "WriteDateTime"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendDateTime(vw.buf, dt)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteDecimal128(d128 primitive.Decimal128) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Decimal128, mode(0), "WriteDecimal128"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendDecimal128(vw.buf, d128)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteDouble(f float64) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Double, mode(0), "WriteDouble"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendDouble(vw.buf, f)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteInt32(i32 int32) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Int32, mode(0), "WriteInt32"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendInt32(vw.buf, i32)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteInt64(i64 int64) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Int64, mode(0), "WriteInt64"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendInt64(vw.buf, i64)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteJavascript(code string) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.JavaScript, mode(0), "WriteJavascript"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendJavaScript(vw.buf, code)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteMaxKey() error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.MaxKey, mode(0), "WriteMaxKey"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteMinKey() error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.MinKey, mode(0), "WriteMinKey"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteNull() error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Null, mode(0), "WriteNull"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteObjectID(oid primitive.ObjectID) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.ObjectID, mode(0), "WriteObjectID"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendObjectID(vw.buf, oid)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteRegex(pattern string, options string) error ***REMOVED***
	if !isValidCString(pattern) || !isValidCString(options) ***REMOVED***
		return errors.New("BSON regex values cannot contain null bytes")
	***REMOVED***
	if err := vw.writeElementHeader(bsontype.Regex, mode(0), "WriteRegex"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendRegex(vw.buf, pattern, sortStringAlphebeticAscending(options))
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteString(s string) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.String, mode(0), "WriteString"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendString(vw.buf, s)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteDocument() (DocumentWriter, error) ***REMOVED***
	if vw.stack[vw.frame].mode == mTopLevel ***REMOVED***
		vw.reserveLength()
		return vw, nil
	***REMOVED***
	if err := vw.writeElementHeader(bsontype.EmbeddedDocument, mDocument, "WriteDocument", mTopLevel); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vw.push(mDocument)
	return vw, nil
***REMOVED***

func (vw *valueWriter) WriteSymbol(symbol string) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Symbol, mode(0), "WriteSymbol"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendSymbol(vw.buf, symbol)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteTimestamp(t uint32, i uint32) error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Timestamp, mode(0), "WriteTimestamp"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.buf = bsoncore.AppendTimestamp(vw.buf, t, i)
	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteUndefined() error ***REMOVED***
	if err := vw.writeElementHeader(bsontype.Undefined, mode(0), "WriteUndefined"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.pop()
	return nil
***REMOVED***

func (vw *valueWriter) WriteDocumentElement(key string) (ValueWriter, error) ***REMOVED***
	switch vw.stack[vw.frame].mode ***REMOVED***
	case mTopLevel, mDocument:
	default:
		return nil, vw.invalidTransitionError(mElement, "WriteDocumentElement", []mode***REMOVED***mTopLevel, mDocument***REMOVED***)
	***REMOVED***

	vw.push(mElement)
	vw.stack[vw.frame].key = key

	return vw, nil
***REMOVED***

func (vw *valueWriter) WriteDocumentEnd() error ***REMOVED***
	switch vw.stack[vw.frame].mode ***REMOVED***
	case mTopLevel, mDocument:
	default:
		return fmt.Errorf("incorrect mode to end document: %s", vw.stack[vw.frame].mode)
	***REMOVED***

	vw.buf = append(vw.buf, 0x00)

	err := vw.writeLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if vw.stack[vw.frame].mode == mTopLevel ***REMOVED***
		if err = vw.Flush(); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	vw.pop()

	if vw.stack[vw.frame].mode == mCodeWithScope ***REMOVED***
		// We ignore the error here because of the guarantee of writeLength.
		// See the docs for writeLength for more info.
		_ = vw.writeLength()
		vw.pop()
	***REMOVED***
	return nil
***REMOVED***

func (vw *valueWriter) Flush() error ***REMOVED***
	if vw.w == nil ***REMOVED***
		return nil
	***REMOVED***

	if _, err := vw.w.Write(vw.buf); err != nil ***REMOVED***
		return err
	***REMOVED***
	// reset buffer
	vw.buf = vw.buf[:0]
	return nil
***REMOVED***

func (vw *valueWriter) WriteArrayElement() (ValueWriter, error) ***REMOVED***
	if vw.stack[vw.frame].mode != mArray ***REMOVED***
		return nil, vw.invalidTransitionError(mValue, "WriteArrayElement", []mode***REMOVED***mArray***REMOVED***)
	***REMOVED***

	arrkey := vw.stack[vw.frame].arrkey
	vw.stack[vw.frame].arrkey++

	vw.push(mValue)
	vw.stack[vw.frame].arrkey = arrkey

	return vw, nil
***REMOVED***

func (vw *valueWriter) WriteArrayEnd() error ***REMOVED***
	if vw.stack[vw.frame].mode != mArray ***REMOVED***
		return fmt.Errorf("incorrect mode to end array: %s", vw.stack[vw.frame].mode)
	***REMOVED***

	vw.buf = append(vw.buf, 0x00)

	err := vw.writeLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	vw.pop()
	return nil
***REMOVED***

// NOTE: We assume that if we call writeLength more than once the same function
// within the same function without altering the vw.buf that this method will
// not return an error. If this changes ensure that the following methods are
// updated:
//
// - WriteDocumentEnd
func (vw *valueWriter) writeLength() error ***REMOVED***
	length := len(vw.buf)
	if length > maxSize ***REMOVED***
		return errMaxDocumentSizeExceeded***REMOVED***size: int64(len(vw.buf))***REMOVED***
	***REMOVED***
	length = length - int(vw.stack[vw.frame].start)
	start := vw.stack[vw.frame].start

	vw.buf[start+0] = byte(length)
	vw.buf[start+1] = byte(length >> 8)
	vw.buf[start+2] = byte(length >> 16)
	vw.buf[start+3] = byte(length >> 24)
	return nil
***REMOVED***

func isValidCString(cs string) bool ***REMOVED***
	return !strings.ContainsRune(cs, '\x00')
***REMOVED***
