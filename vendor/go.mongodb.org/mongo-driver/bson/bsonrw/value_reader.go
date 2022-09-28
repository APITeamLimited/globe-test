// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ ValueReader = (*valueReader)(nil)

var vrPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED***
		return new(valueReader)
	***REMOVED***,
***REMOVED***

// BSONValueReaderPool is a pool for ValueReaders that read BSON.
type BSONValueReaderPool struct ***REMOVED***
	pool sync.Pool
***REMOVED***

// NewBSONValueReaderPool instantiates a new BSONValueReaderPool.
func NewBSONValueReaderPool() *BSONValueReaderPool ***REMOVED***
	return &BSONValueReaderPool***REMOVED***
		pool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED***
				return new(valueReader)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Get retrieves a ValueReader from the pool and uses src as the underlying BSON.
func (bvrp *BSONValueReaderPool) Get(src []byte) ValueReader ***REMOVED***
	vr := bvrp.pool.Get().(*valueReader)
	vr.reset(src)
	return vr
***REMOVED***

// Put inserts a ValueReader into the pool. If the ValueReader is not a BSON ValueReader nothing
// is inserted into the pool and ok will be false.
func (bvrp *BSONValueReaderPool) Put(vr ValueReader) (ok bool) ***REMOVED***
	bvr, ok := vr.(*valueReader)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	bvr.reset(nil)
	bvrp.pool.Put(bvr)
	return true
***REMOVED***

// ErrEOA is the error returned when the end of a BSON array has been reached.
var ErrEOA = errors.New("end of array")

// ErrEOD is the error returned when the end of a BSON document has been reached.
var ErrEOD = errors.New("end of document")

type vrState struct ***REMOVED***
	mode  mode
	vType bsontype.Type
	end   int64
***REMOVED***

// valueReader is for reading BSON values.
type valueReader struct ***REMOVED***
	offset int64
	d      []byte

	stack []vrState
	frame int64
***REMOVED***

// NewBSONDocumentReader returns a ValueReader using b for the underlying BSON
// representation. Parameter b must be a BSON Document.
func NewBSONDocumentReader(b []byte) ValueReader ***REMOVED***
	// TODO(skriptble): There's a lack of symmetry between the reader and writer, since the reader takes a []byte while the
	// TODO writer takes an io.Writer. We should have two versions of each, one that takes a []byte and one that takes an
	// TODO io.Reader or io.Writer. The []byte version will need to return a thing that can return the finished []byte since
	// TODO it might be reallocated when appended to.
	return newValueReader(b)
***REMOVED***

// NewBSONValueReader returns a ValueReader that starts in the Value mode instead of in top
// level document mode. This enables the creation of a ValueReader for a single BSON value.
func NewBSONValueReader(t bsontype.Type, val []byte) ValueReader ***REMOVED***
	stack := make([]vrState, 1, 5)
	stack[0] = vrState***REMOVED***
		mode:  mValue,
		vType: t,
	***REMOVED***
	return &valueReader***REMOVED***
		d:     val,
		stack: stack,
	***REMOVED***
***REMOVED***

func newValueReader(b []byte) *valueReader ***REMOVED***
	stack := make([]vrState, 1, 5)
	stack[0] = vrState***REMOVED***
		mode: mTopLevel,
	***REMOVED***
	return &valueReader***REMOVED***
		d:     b,
		stack: stack,
	***REMOVED***
***REMOVED***

func (vr *valueReader) reset(b []byte) ***REMOVED***
	if vr.stack == nil ***REMOVED***
		vr.stack = make([]vrState, 1, 5)
	***REMOVED***
	vr.stack = vr.stack[:1]
	vr.stack[0] = vrState***REMOVED***mode: mTopLevel***REMOVED***
	vr.d = b
	vr.offset = 0
	vr.frame = 0
***REMOVED***

func (vr *valueReader) advanceFrame() ***REMOVED***
	if vr.frame+1 >= int64(len(vr.stack)) ***REMOVED*** // We need to grow the stack
		length := len(vr.stack)
		if length+1 >= cap(vr.stack) ***REMOVED***
			// double it
			buf := make([]vrState, 2*cap(vr.stack)+1)
			copy(buf, vr.stack)
			vr.stack = buf
		***REMOVED***
		vr.stack = vr.stack[:length+1]
	***REMOVED***
	vr.frame++

	// Clean the stack
	vr.stack[vr.frame].mode = 0
	vr.stack[vr.frame].vType = 0
	vr.stack[vr.frame].end = 0
***REMOVED***

func (vr *valueReader) pushDocument() error ***REMOVED***
	vr.advanceFrame()

	vr.stack[vr.frame].mode = mDocument

	size, err := vr.readLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	vr.stack[vr.frame].end = int64(size) + vr.offset - 4

	return nil
***REMOVED***

func (vr *valueReader) pushArray() error ***REMOVED***
	vr.advanceFrame()

	vr.stack[vr.frame].mode = mArray

	size, err := vr.readLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	vr.stack[vr.frame].end = int64(size) + vr.offset - 4

	return nil
***REMOVED***

func (vr *valueReader) pushElement(t bsontype.Type) ***REMOVED***
	vr.advanceFrame()

	vr.stack[vr.frame].mode = mElement
	vr.stack[vr.frame].vType = t
***REMOVED***

func (vr *valueReader) pushValue(t bsontype.Type) ***REMOVED***
	vr.advanceFrame()

	vr.stack[vr.frame].mode = mValue
	vr.stack[vr.frame].vType = t
***REMOVED***

func (vr *valueReader) pushCodeWithScope() (int64, error) ***REMOVED***
	vr.advanceFrame()

	vr.stack[vr.frame].mode = mCodeWithScope

	size, err := vr.readLength()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	vr.stack[vr.frame].end = int64(size) + vr.offset - 4

	return int64(size), nil
***REMOVED***

func (vr *valueReader) pop() ***REMOVED***
	switch vr.stack[vr.frame].mode ***REMOVED***
	case mElement, mValue:
		vr.frame--
	case mDocument, mArray, mCodeWithScope:
		vr.frame -= 2 // we pop twice to jump over the vrElement: vrDocument -> vrElement -> vrDocument/TopLevel/etc...
	***REMOVED***
***REMOVED***

func (vr *valueReader) invalidTransitionErr(destination mode, name string, modes []mode) error ***REMOVED***
	te := TransitionError***REMOVED***
		name:        name,
		current:     vr.stack[vr.frame].mode,
		destination: destination,
		modes:       modes,
		action:      "read",
	***REMOVED***
	if vr.frame != 0 ***REMOVED***
		te.parent = vr.stack[vr.frame-1].mode
	***REMOVED***
	return te
***REMOVED***

func (vr *valueReader) typeError(t bsontype.Type) error ***REMOVED***
	return fmt.Errorf("positioned on %s, but attempted to read %s", vr.stack[vr.frame].vType, t)
***REMOVED***

func (vr *valueReader) invalidDocumentLengthError() error ***REMOVED***
	return fmt.Errorf("document is invalid, end byte is at %d, but null byte found at %d", vr.stack[vr.frame].end, vr.offset)
***REMOVED***

func (vr *valueReader) ensureElementValue(t bsontype.Type, destination mode, callerName string) error ***REMOVED***
	switch vr.stack[vr.frame].mode ***REMOVED***
	case mElement, mValue:
		if vr.stack[vr.frame].vType != t ***REMOVED***
			return vr.typeError(t)
		***REMOVED***
	default:
		return vr.invalidTransitionErr(destination, callerName, []mode***REMOVED***mElement, mValue***REMOVED***)
	***REMOVED***

	return nil
***REMOVED***

func (vr *valueReader) Type() bsontype.Type ***REMOVED***
	return vr.stack[vr.frame].vType
***REMOVED***

func (vr *valueReader) nextElementLength() (int32, error) ***REMOVED***
	var length int32
	var err error
	switch vr.stack[vr.frame].vType ***REMOVED***
	case bsontype.Array, bsontype.EmbeddedDocument, bsontype.CodeWithScope:
		length, err = vr.peekLength()
	case bsontype.Binary:
		length, err = vr.peekLength()
		length += 4 + 1 // binary length + subtype byte
	case bsontype.Boolean:
		length = 1
	case bsontype.DBPointer:
		length, err = vr.peekLength()
		length += 4 + 12 // string length + ObjectID length
	case bsontype.DateTime, bsontype.Double, bsontype.Int64, bsontype.Timestamp:
		length = 8
	case bsontype.Decimal128:
		length = 16
	case bsontype.Int32:
		length = 4
	case bsontype.JavaScript, bsontype.String, bsontype.Symbol:
		length, err = vr.peekLength()
		length += 4
	case bsontype.MaxKey, bsontype.MinKey, bsontype.Null, bsontype.Undefined:
		length = 0
	case bsontype.ObjectID:
		length = 12
	case bsontype.Regex:
		regex := bytes.IndexByte(vr.d[vr.offset:], 0x00)
		if regex < 0 ***REMOVED***
			err = io.EOF
			break
		***REMOVED***
		pattern := bytes.IndexByte(vr.d[vr.offset+int64(regex)+1:], 0x00)
		if pattern < 0 ***REMOVED***
			err = io.EOF
			break
		***REMOVED***
		length = int32(int64(regex) + 1 + int64(pattern) + 1)
	default:
		return 0, fmt.Errorf("attempted to read bytes of unknown BSON type %v", vr.stack[vr.frame].vType)
	***REMOVED***

	return length, err
***REMOVED***

func (vr *valueReader) ReadValueBytes(dst []byte) (bsontype.Type, []byte, error) ***REMOVED***
	switch vr.stack[vr.frame].mode ***REMOVED***
	case mTopLevel:
		length, err := vr.peekLength()
		if err != nil ***REMOVED***
			return bsontype.Type(0), nil, err
		***REMOVED***
		dst, err = vr.appendBytes(dst, length)
		if err != nil ***REMOVED***
			return bsontype.Type(0), nil, err
		***REMOVED***
		return bsontype.Type(0), dst, nil
	case mElement, mValue:
		length, err := vr.nextElementLength()
		if err != nil ***REMOVED***
			return bsontype.Type(0), dst, err
		***REMOVED***

		dst, err = vr.appendBytes(dst, length)
		t := vr.stack[vr.frame].vType
		vr.pop()
		return t, dst, err
	default:
		return bsontype.Type(0), nil, vr.invalidTransitionErr(0, "ReadValueBytes", []mode***REMOVED***mElement, mValue***REMOVED***)
	***REMOVED***
***REMOVED***

func (vr *valueReader) Skip() error ***REMOVED***
	switch vr.stack[vr.frame].mode ***REMOVED***
	case mElement, mValue:
	default:
		return vr.invalidTransitionErr(0, "Skip", []mode***REMOVED***mElement, mValue***REMOVED***)
	***REMOVED***

	length, err := vr.nextElementLength()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = vr.skipBytes(length)
	vr.pop()
	return err
***REMOVED***

func (vr *valueReader) ReadArray() (ArrayReader, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Array, mArray, "ReadArray"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	err := vr.pushArray()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return vr, nil
***REMOVED***

func (vr *valueReader) ReadBinary() (b []byte, btype byte, err error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Binary, 0, "ReadBinary"); err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	length, err := vr.readLength()
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	btype, err = vr.readByte()
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	// Check length in case it is an old binary without a length.
	if btype == 0x02 && length > 4 ***REMOVED***
		length, err = vr.readLength()
		if err != nil ***REMOVED***
			return nil, 0, err
		***REMOVED***
	***REMOVED***

	b, err = vr.readBytes(length)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***
	// Make a copy of the returned byte slice because it's just a subslice from the valueReader's
	// buffer and is not safe to return in the unmarshaled value.
	cp := make([]byte, len(b))
	copy(cp, b)

	vr.pop()
	return cp, btype, nil
***REMOVED***

func (vr *valueReader) ReadBoolean() (bool, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Boolean, 0, "ReadBoolean"); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	b, err := vr.readByte()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if b > 1 ***REMOVED***
		return false, fmt.Errorf("invalid byte for boolean, %b", b)
	***REMOVED***

	vr.pop()
	return b == 1, nil
***REMOVED***

func (vr *valueReader) ReadDocument() (DocumentReader, error) ***REMOVED***
	switch vr.stack[vr.frame].mode ***REMOVED***
	case mTopLevel:
		// read size
		size, err := vr.readLength()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if int(size) != len(vr.d) ***REMOVED***
			return nil, fmt.Errorf("invalid document length")
		***REMOVED***
		vr.stack[vr.frame].end = int64(size) + vr.offset - 4
		return vr, nil
	case mElement, mValue:
		if vr.stack[vr.frame].vType != bsontype.EmbeddedDocument ***REMOVED***
			return nil, vr.typeError(bsontype.EmbeddedDocument)
		***REMOVED***
	default:
		return nil, vr.invalidTransitionErr(mDocument, "ReadDocument", []mode***REMOVED***mTopLevel, mElement, mValue***REMOVED***)
	***REMOVED***

	err := vr.pushDocument()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return vr, nil
***REMOVED***

func (vr *valueReader) ReadCodeWithScope() (code string, dr DocumentReader, err error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.CodeWithScope, 0, "ReadCodeWithScope"); err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***

	totalLength, err := vr.readLength()
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	strLength, err := vr.readLength()
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	if strLength <= 0 ***REMOVED***
		return "", nil, fmt.Errorf("invalid string length: %d", strLength)
	***REMOVED***
	strBytes, err := vr.readBytes(strLength)
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***
	code = string(strBytes[:len(strBytes)-1])

	size, err := vr.pushCodeWithScope()
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***

	// The total length should equal:
	// 4 (total length) + strLength + 4 (the length of str itself) + (document length)
	componentsLength := int64(4+strLength+4) + size
	if int64(totalLength) != componentsLength ***REMOVED***
		return "", nil, fmt.Errorf(
			"length of CodeWithScope does not match lengths of components; total: %d; components: %d",
			totalLength, componentsLength,
		)
	***REMOVED***
	return code, vr, nil
***REMOVED***

func (vr *valueReader) ReadDBPointer() (ns string, oid primitive.ObjectID, err error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.DBPointer, 0, "ReadDBPointer"); err != nil ***REMOVED***
		return "", oid, err
	***REMOVED***

	ns, err = vr.readString()
	if err != nil ***REMOVED***
		return "", oid, err
	***REMOVED***

	oidbytes, err := vr.readBytes(12)
	if err != nil ***REMOVED***
		return "", oid, err
	***REMOVED***

	copy(oid[:], oidbytes)

	vr.pop()
	return ns, oid, nil
***REMOVED***

func (vr *valueReader) ReadDateTime() (int64, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.DateTime, 0, "ReadDateTime"); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	i, err := vr.readi64()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	vr.pop()
	return i, nil
***REMOVED***

func (vr *valueReader) ReadDecimal128() (primitive.Decimal128, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Decimal128, 0, "ReadDecimal128"); err != nil ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, err
	***REMOVED***

	b, err := vr.readBytes(16)
	if err != nil ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, err
	***REMOVED***

	l := binary.LittleEndian.Uint64(b[0:8])
	h := binary.LittleEndian.Uint64(b[8:16])

	vr.pop()
	return primitive.NewDecimal128(h, l), nil
***REMOVED***

func (vr *valueReader) ReadDouble() (float64, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Double, 0, "ReadDouble"); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	u, err := vr.readu64()
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	vr.pop()
	return math.Float64frombits(u), nil
***REMOVED***

func (vr *valueReader) ReadInt32() (int32, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Int32, 0, "ReadInt32"); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	vr.pop()
	return vr.readi32()
***REMOVED***

func (vr *valueReader) ReadInt64() (int64, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Int64, 0, "ReadInt64"); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	vr.pop()
	return vr.readi64()
***REMOVED***

func (vr *valueReader) ReadJavascript() (code string, err error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.JavaScript, 0, "ReadJavascript"); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	vr.pop()
	return vr.readString()
***REMOVED***

func (vr *valueReader) ReadMaxKey() error ***REMOVED***
	if err := vr.ensureElementValue(bsontype.MaxKey, 0, "ReadMaxKey"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vr.pop()
	return nil
***REMOVED***

func (vr *valueReader) ReadMinKey() error ***REMOVED***
	if err := vr.ensureElementValue(bsontype.MinKey, 0, "ReadMinKey"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vr.pop()
	return nil
***REMOVED***

func (vr *valueReader) ReadNull() error ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Null, 0, "ReadNull"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vr.pop()
	return nil
***REMOVED***

func (vr *valueReader) ReadObjectID() (primitive.ObjectID, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.ObjectID, 0, "ReadObjectID"); err != nil ***REMOVED***
		return primitive.ObjectID***REMOVED******REMOVED***, err
	***REMOVED***

	oidbytes, err := vr.readBytes(12)
	if err != nil ***REMOVED***
		return primitive.ObjectID***REMOVED******REMOVED***, err
	***REMOVED***

	var oid primitive.ObjectID
	copy(oid[:], oidbytes)

	vr.pop()
	return oid, nil
***REMOVED***

func (vr *valueReader) ReadRegex() (string, string, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Regex, 0, "ReadRegex"); err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	pattern, err := vr.readCString()
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	options, err := vr.readCString()
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	vr.pop()
	return pattern, options, nil
***REMOVED***

func (vr *valueReader) ReadString() (string, error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.String, 0, "ReadString"); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	vr.pop()
	return vr.readString()
***REMOVED***

func (vr *valueReader) ReadSymbol() (symbol string, err error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Symbol, 0, "ReadSymbol"); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	vr.pop()
	return vr.readString()
***REMOVED***

func (vr *valueReader) ReadTimestamp() (t uint32, i uint32, err error) ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Timestamp, 0, "ReadTimestamp"); err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***

	i, err = vr.readu32()
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***

	t, err = vr.readu32()
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***

	vr.pop()
	return t, i, nil
***REMOVED***

func (vr *valueReader) ReadUndefined() error ***REMOVED***
	if err := vr.ensureElementValue(bsontype.Undefined, 0, "ReadUndefined"); err != nil ***REMOVED***
		return err
	***REMOVED***

	vr.pop()
	return nil
***REMOVED***

func (vr *valueReader) ReadElement() (string, ValueReader, error) ***REMOVED***
	switch vr.stack[vr.frame].mode ***REMOVED***
	case mTopLevel, mDocument, mCodeWithScope:
	default:
		return "", nil, vr.invalidTransitionErr(mElement, "ReadElement", []mode***REMOVED***mTopLevel, mDocument, mCodeWithScope***REMOVED***)
	***REMOVED***

	t, err := vr.readByte()
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***

	if t == 0 ***REMOVED***
		if vr.offset != vr.stack[vr.frame].end ***REMOVED***
			return "", nil, vr.invalidDocumentLengthError()
		***REMOVED***

		vr.pop()
		return "", nil, ErrEOD
	***REMOVED***

	name, err := vr.readCString()
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***

	vr.pushElement(bsontype.Type(t))
	return name, vr, nil
***REMOVED***

func (vr *valueReader) ReadValue() (ValueReader, error) ***REMOVED***
	switch vr.stack[vr.frame].mode ***REMOVED***
	case mArray:
	default:
		return nil, vr.invalidTransitionErr(mValue, "ReadValue", []mode***REMOVED***mArray***REMOVED***)
	***REMOVED***

	t, err := vr.readByte()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if t == 0 ***REMOVED***
		if vr.offset != vr.stack[vr.frame].end ***REMOVED***
			return nil, vr.invalidDocumentLengthError()
		***REMOVED***

		vr.pop()
		return nil, ErrEOA
	***REMOVED***

	_, err = vr.readCString()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	vr.pushValue(bsontype.Type(t))
	return vr, nil
***REMOVED***

// readBytes reads length bytes from the valueReader starting at the current offset. Note that the
// returned byte slice is a subslice from the valueReader buffer and must be converted or copied
// before returning in an unmarshaled value.
func (vr *valueReader) readBytes(length int32) ([]byte, error) ***REMOVED***
	if length < 0 ***REMOVED***
		return nil, fmt.Errorf("invalid length: %d", length)
	***REMOVED***

	if vr.offset+int64(length) > int64(len(vr.d)) ***REMOVED***
		return nil, io.EOF
	***REMOVED***

	start := vr.offset
	vr.offset += int64(length)

	return vr.d[start : start+int64(length)], nil
***REMOVED***

func (vr *valueReader) appendBytes(dst []byte, length int32) ([]byte, error) ***REMOVED***
	if vr.offset+int64(length) > int64(len(vr.d)) ***REMOVED***
		return nil, io.EOF
	***REMOVED***

	start := vr.offset
	vr.offset += int64(length)
	return append(dst, vr.d[start:start+int64(length)]...), nil
***REMOVED***

func (vr *valueReader) skipBytes(length int32) error ***REMOVED***
	if vr.offset+int64(length) > int64(len(vr.d)) ***REMOVED***
		return io.EOF
	***REMOVED***

	vr.offset += int64(length)
	return nil
***REMOVED***

func (vr *valueReader) readByte() (byte, error) ***REMOVED***
	if vr.offset+1 > int64(len(vr.d)) ***REMOVED***
		return 0x0, io.EOF
	***REMOVED***

	vr.offset++
	return vr.d[vr.offset-1], nil
***REMOVED***

func (vr *valueReader) readCString() (string, error) ***REMOVED***
	idx := bytes.IndexByte(vr.d[vr.offset:], 0x00)
	if idx < 0 ***REMOVED***
		return "", io.EOF
	***REMOVED***
	start := vr.offset
	// idx does not include the null byte
	vr.offset += int64(idx) + 1
	return string(vr.d[start : start+int64(idx)]), nil
***REMOVED***

func (vr *valueReader) readString() (string, error) ***REMOVED***
	length, err := vr.readLength()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if int64(length)+vr.offset > int64(len(vr.d)) ***REMOVED***
		return "", io.EOF
	***REMOVED***

	if length <= 0 ***REMOVED***
		return "", fmt.Errorf("invalid string length: %d", length)
	***REMOVED***

	if vr.d[vr.offset+int64(length)-1] != 0x00 ***REMOVED***
		return "", fmt.Errorf("string does not end with null byte, but with %v", vr.d[vr.offset+int64(length)-1])
	***REMOVED***

	start := vr.offset
	vr.offset += int64(length)
	return string(vr.d[start : start+int64(length)-1]), nil
***REMOVED***

func (vr *valueReader) peekLength() (int32, error) ***REMOVED***
	if vr.offset+4 > int64(len(vr.d)) ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	idx := vr.offset
	return (int32(vr.d[idx]) | int32(vr.d[idx+1])<<8 | int32(vr.d[idx+2])<<16 | int32(vr.d[idx+3])<<24), nil
***REMOVED***

func (vr *valueReader) readLength() (int32, error) ***REMOVED*** return vr.readi32() ***REMOVED***

func (vr *valueReader) readi32() (int32, error) ***REMOVED***
	if vr.offset+4 > int64(len(vr.d)) ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	idx := vr.offset
	vr.offset += 4
	return (int32(vr.d[idx]) | int32(vr.d[idx+1])<<8 | int32(vr.d[idx+2])<<16 | int32(vr.d[idx+3])<<24), nil
***REMOVED***

func (vr *valueReader) readu32() (uint32, error) ***REMOVED***
	if vr.offset+4 > int64(len(vr.d)) ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	idx := vr.offset
	vr.offset += 4
	return (uint32(vr.d[idx]) | uint32(vr.d[idx+1])<<8 | uint32(vr.d[idx+2])<<16 | uint32(vr.d[idx+3])<<24), nil
***REMOVED***

func (vr *valueReader) readi64() (int64, error) ***REMOVED***
	if vr.offset+8 > int64(len(vr.d)) ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	idx := vr.offset
	vr.offset += 8
	return int64(vr.d[idx]) | int64(vr.d[idx+1])<<8 | int64(vr.d[idx+2])<<16 | int64(vr.d[idx+3])<<24 |
		int64(vr.d[idx+4])<<32 | int64(vr.d[idx+5])<<40 | int64(vr.d[idx+6])<<48 | int64(vr.d[idx+7])<<56, nil
***REMOVED***

func (vr *valueReader) readu64() (uint64, error) ***REMOVED***
	if vr.offset+8 > int64(len(vr.d)) ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	idx := vr.offset
	vr.offset += 8
	return uint64(vr.d[idx]) | uint64(vr.d[idx+1])<<8 | uint64(vr.d[idx+2])<<16 | uint64(vr.d[idx+3])<<24 |
		uint64(vr.d[idx+4])<<32 | uint64(vr.d[idx+5])<<40 | uint64(vr.d[idx+6])<<48 | uint64(vr.d[idx+7])<<56, nil
***REMOVED***
