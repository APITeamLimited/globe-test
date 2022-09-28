// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"fmt"
	"io"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ExtJSONValueReaderPool is a pool for ValueReaders that read ExtJSON.
type ExtJSONValueReaderPool struct ***REMOVED***
	pool sync.Pool
***REMOVED***

// NewExtJSONValueReaderPool instantiates a new ExtJSONValueReaderPool.
func NewExtJSONValueReaderPool() *ExtJSONValueReaderPool ***REMOVED***
	return &ExtJSONValueReaderPool***REMOVED***
		pool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED***
				return new(extJSONValueReader)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Get retrieves a ValueReader from the pool and uses src as the underlying ExtJSON.
func (bvrp *ExtJSONValueReaderPool) Get(r io.Reader, canonical bool) (ValueReader, error) ***REMOVED***
	vr := bvrp.pool.Get().(*extJSONValueReader)
	return vr.reset(r, canonical)
***REMOVED***

// Put inserts a ValueReader into the pool. If the ValueReader is not a ExtJSON ValueReader nothing
// is inserted into the pool and ok will be false.
func (bvrp *ExtJSONValueReaderPool) Put(vr ValueReader) (ok bool) ***REMOVED***
	bvr, ok := vr.(*extJSONValueReader)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	bvr, _ = bvr.reset(nil, false)
	bvrp.pool.Put(bvr)
	return true
***REMOVED***

type ejvrState struct ***REMOVED***
	mode  mode
	vType bsontype.Type
	depth int
***REMOVED***

// extJSONValueReader is for reading extended JSON.
type extJSONValueReader struct ***REMOVED***
	p *extJSONParser

	stack []ejvrState
	frame int
***REMOVED***

// NewExtJSONValueReader creates a new ValueReader from a given io.Reader
// It will interpret the JSON of r as canonical or relaxed according to the
// given canonical flag
func NewExtJSONValueReader(r io.Reader, canonical bool) (ValueReader, error) ***REMOVED***
	return newExtJSONValueReader(r, canonical)
***REMOVED***

func newExtJSONValueReader(r io.Reader, canonical bool) (*extJSONValueReader, error) ***REMOVED***
	ejvr := new(extJSONValueReader)
	return ejvr.reset(r, canonical)
***REMOVED***

func (ejvr *extJSONValueReader) reset(r io.Reader, canonical bool) (*extJSONValueReader, error) ***REMOVED***
	p := newExtJSONParser(r, canonical)
	typ, err := p.peekType()

	if err != nil ***REMOVED***
		return nil, ErrInvalidJSON
	***REMOVED***

	var m mode
	switch typ ***REMOVED***
	case bsontype.EmbeddedDocument:
		m = mTopLevel
	case bsontype.Array:
		m = mArray
	default:
		m = mValue
	***REMOVED***

	stack := make([]ejvrState, 1, 5)
	stack[0] = ejvrState***REMOVED***
		mode:  m,
		vType: typ,
	***REMOVED***
	return &extJSONValueReader***REMOVED***
		p:     p,
		stack: stack,
	***REMOVED***, nil
***REMOVED***

func (ejvr *extJSONValueReader) advanceFrame() ***REMOVED***
	if ejvr.frame+1 >= len(ejvr.stack) ***REMOVED*** // We need to grow the stack
		length := len(ejvr.stack)
		if length+1 >= cap(ejvr.stack) ***REMOVED***
			// double it
			buf := make([]ejvrState, 2*cap(ejvr.stack)+1)
			copy(buf, ejvr.stack)
			ejvr.stack = buf
		***REMOVED***
		ejvr.stack = ejvr.stack[:length+1]
	***REMOVED***
	ejvr.frame++

	// Clean the stack
	ejvr.stack[ejvr.frame].mode = 0
	ejvr.stack[ejvr.frame].vType = 0
	ejvr.stack[ejvr.frame].depth = 0
***REMOVED***

func (ejvr *extJSONValueReader) pushDocument() ***REMOVED***
	ejvr.advanceFrame()

	ejvr.stack[ejvr.frame].mode = mDocument
	ejvr.stack[ejvr.frame].depth = ejvr.p.depth
***REMOVED***

func (ejvr *extJSONValueReader) pushCodeWithScope() ***REMOVED***
	ejvr.advanceFrame()

	ejvr.stack[ejvr.frame].mode = mCodeWithScope
***REMOVED***

func (ejvr *extJSONValueReader) pushArray() ***REMOVED***
	ejvr.advanceFrame()

	ejvr.stack[ejvr.frame].mode = mArray
***REMOVED***

func (ejvr *extJSONValueReader) push(m mode, t bsontype.Type) ***REMOVED***
	ejvr.advanceFrame()

	ejvr.stack[ejvr.frame].mode = m
	ejvr.stack[ejvr.frame].vType = t
***REMOVED***

func (ejvr *extJSONValueReader) pop() ***REMOVED***
	switch ejvr.stack[ejvr.frame].mode ***REMOVED***
	case mElement, mValue:
		ejvr.frame--
	case mDocument, mArray, mCodeWithScope:
		ejvr.frame -= 2 // we pop twice to jump over the vrElement: vrDocument -> vrElement -> vrDocument/TopLevel/etc...
	***REMOVED***
***REMOVED***

func (ejvr *extJSONValueReader) skipObject() ***REMOVED***
	// read entire object until depth returns to 0 (last ending ***REMOVED*** or ] seen)
	depth := 1
	for depth > 0 ***REMOVED***
		ejvr.p.advanceState()

		// If object is empty, raise depth and continue. When emptyObject is true, the
		// parser has already read both the opening and closing brackets of an empty
		// object ("***REMOVED******REMOVED***"), so the next valid token will be part of the parent document,
		// not part of the nested document.
		//
		// If there is a comma, there are remaining fields, emptyObject must be set back
		// to false, and comma must be skipped with advanceState().
		if ejvr.p.emptyObject ***REMOVED***
			if ejvr.p.s == jpsSawComma ***REMOVED***
				ejvr.p.emptyObject = false
				ejvr.p.advanceState()
			***REMOVED***
			depth--
			continue
		***REMOVED***

		switch ejvr.p.s ***REMOVED***
		case jpsSawBeginObject, jpsSawBeginArray:
			depth++
		case jpsSawEndObject, jpsSawEndArray:
			depth--
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ejvr *extJSONValueReader) invalidTransitionErr(destination mode, name string, modes []mode) error ***REMOVED***
	te := TransitionError***REMOVED***
		name:        name,
		current:     ejvr.stack[ejvr.frame].mode,
		destination: destination,
		modes:       modes,
		action:      "read",
	***REMOVED***
	if ejvr.frame != 0 ***REMOVED***
		te.parent = ejvr.stack[ejvr.frame-1].mode
	***REMOVED***
	return te
***REMOVED***

func (ejvr *extJSONValueReader) typeError(t bsontype.Type) error ***REMOVED***
	return fmt.Errorf("positioned on %s, but attempted to read %s", ejvr.stack[ejvr.frame].vType, t)
***REMOVED***

func (ejvr *extJSONValueReader) ensureElementValue(t bsontype.Type, destination mode, callerName string, addModes ...mode) error ***REMOVED***
	switch ejvr.stack[ejvr.frame].mode ***REMOVED***
	case mElement, mValue:
		if ejvr.stack[ejvr.frame].vType != t ***REMOVED***
			return ejvr.typeError(t)
		***REMOVED***
	default:
		modes := []mode***REMOVED***mElement, mValue***REMOVED***
		if addModes != nil ***REMOVED***
			modes = append(modes, addModes...)
		***REMOVED***
		return ejvr.invalidTransitionErr(destination, callerName, modes)
	***REMOVED***

	return nil
***REMOVED***

func (ejvr *extJSONValueReader) Type() bsontype.Type ***REMOVED***
	return ejvr.stack[ejvr.frame].vType
***REMOVED***

func (ejvr *extJSONValueReader) Skip() error ***REMOVED***
	switch ejvr.stack[ejvr.frame].mode ***REMOVED***
	case mElement, mValue:
	default:
		return ejvr.invalidTransitionErr(0, "Skip", []mode***REMOVED***mElement, mValue***REMOVED***)
	***REMOVED***

	defer ejvr.pop()

	t := ejvr.stack[ejvr.frame].vType
	switch t ***REMOVED***
	case bsontype.Array, bsontype.EmbeddedDocument, bsontype.CodeWithScope:
		// read entire array, doc or CodeWithScope
		ejvr.skipObject()
	default:
		_, err := ejvr.p.readValue(t)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (ejvr *extJSONValueReader) ReadArray() (ArrayReader, error) ***REMOVED***
	switch ejvr.stack[ejvr.frame].mode ***REMOVED***
	case mTopLevel: // allow reading array from top level
	case mArray:
		return ejvr, nil
	default:
		if err := ejvr.ensureElementValue(bsontype.Array, mArray, "ReadArray", mTopLevel, mArray); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	ejvr.pushArray()

	return ejvr, nil
***REMOVED***

func (ejvr *extJSONValueReader) ReadBinary() (b []byte, btype byte, err error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.Binary, 0, "ReadBinary"); err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Binary)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	b, btype, err = v.parseBinary()

	ejvr.pop()
	return b, btype, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadBoolean() (bool, error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.Boolean, 0, "ReadBoolean"); err != nil ***REMOVED***
		return false, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Boolean)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if v.t != bsontype.Boolean ***REMOVED***
		return false, fmt.Errorf("expected type bool, but got type %s", v.t)
	***REMOVED***

	ejvr.pop()
	return v.v.(bool), nil
***REMOVED***

func (ejvr *extJSONValueReader) ReadDocument() (DocumentReader, error) ***REMOVED***
	switch ejvr.stack[ejvr.frame].mode ***REMOVED***
	case mTopLevel:
		return ejvr, nil
	case mElement, mValue:
		if ejvr.stack[ejvr.frame].vType != bsontype.EmbeddedDocument ***REMOVED***
			return nil, ejvr.typeError(bsontype.EmbeddedDocument)
		***REMOVED***

		ejvr.pushDocument()
		return ejvr, nil
	default:
		return nil, ejvr.invalidTransitionErr(mDocument, "ReadDocument", []mode***REMOVED***mTopLevel, mElement, mValue***REMOVED***)
	***REMOVED***
***REMOVED***

func (ejvr *extJSONValueReader) ReadCodeWithScope() (code string, dr DocumentReader, err error) ***REMOVED***
	if err = ejvr.ensureElementValue(bsontype.CodeWithScope, 0, "ReadCodeWithScope"); err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.CodeWithScope)
	if err != nil ***REMOVED***
		return "", nil, err
	***REMOVED***

	code, err = v.parseJavascript()

	ejvr.pushCodeWithScope()
	return code, ejvr, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadDBPointer() (ns string, oid primitive.ObjectID, err error) ***REMOVED***
	if err = ejvr.ensureElementValue(bsontype.DBPointer, 0, "ReadDBPointer"); err != nil ***REMOVED***
		return "", primitive.NilObjectID, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.DBPointer)
	if err != nil ***REMOVED***
		return "", primitive.NilObjectID, err
	***REMOVED***

	ns, oid, err = v.parseDBPointer()

	ejvr.pop()
	return ns, oid, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadDateTime() (int64, error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.DateTime, 0, "ReadDateTime"); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.DateTime)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	d, err := v.parseDateTime()

	ejvr.pop()
	return d, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadDecimal128() (primitive.Decimal128, error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.Decimal128, 0, "ReadDecimal128"); err != nil ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Decimal128)
	if err != nil ***REMOVED***
		return primitive.Decimal128***REMOVED******REMOVED***, err
	***REMOVED***

	d, err := v.parseDecimal128()

	ejvr.pop()
	return d, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadDouble() (float64, error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.Double, 0, "ReadDouble"); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Double)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	d, err := v.parseDouble()

	ejvr.pop()
	return d, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadInt32() (int32, error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.Int32, 0, "ReadInt32"); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Int32)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	i, err := v.parseInt32()

	ejvr.pop()
	return i, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadInt64() (int64, error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.Int64, 0, "ReadInt64"); err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Int64)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	i, err := v.parseInt64()

	ejvr.pop()
	return i, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadJavascript() (code string, err error) ***REMOVED***
	if err = ejvr.ensureElementValue(bsontype.JavaScript, 0, "ReadJavascript"); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.JavaScript)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	code, err = v.parseJavascript()

	ejvr.pop()
	return code, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadMaxKey() error ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.MaxKey, 0, "ReadMaxKey"); err != nil ***REMOVED***
		return err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.MaxKey)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = v.parseMinMaxKey("max")

	ejvr.pop()
	return err
***REMOVED***

func (ejvr *extJSONValueReader) ReadMinKey() error ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.MinKey, 0, "ReadMinKey"); err != nil ***REMOVED***
		return err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.MinKey)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = v.parseMinMaxKey("min")

	ejvr.pop()
	return err
***REMOVED***

func (ejvr *extJSONValueReader) ReadNull() error ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.Null, 0, "ReadNull"); err != nil ***REMOVED***
		return err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Null)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if v.t != bsontype.Null ***REMOVED***
		return fmt.Errorf("expected type null but got type %s", v.t)
	***REMOVED***

	ejvr.pop()
	return nil
***REMOVED***

func (ejvr *extJSONValueReader) ReadObjectID() (primitive.ObjectID, error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.ObjectID, 0, "ReadObjectID"); err != nil ***REMOVED***
		return primitive.ObjectID***REMOVED******REMOVED***, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.ObjectID)
	if err != nil ***REMOVED***
		return primitive.ObjectID***REMOVED******REMOVED***, err
	***REMOVED***

	oid, err := v.parseObjectID()

	ejvr.pop()
	return oid, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadRegex() (pattern string, options string, err error) ***REMOVED***
	if err = ejvr.ensureElementValue(bsontype.Regex, 0, "ReadRegex"); err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Regex)
	if err != nil ***REMOVED***
		return "", "", err
	***REMOVED***

	pattern, options, err = v.parseRegex()

	ejvr.pop()
	return pattern, options, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadString() (string, error) ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.String, 0, "ReadString"); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.String)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	if v.t != bsontype.String ***REMOVED***
		return "", fmt.Errorf("expected type string but got type %s", v.t)
	***REMOVED***

	ejvr.pop()
	return v.v.(string), nil
***REMOVED***

func (ejvr *extJSONValueReader) ReadSymbol() (symbol string, err error) ***REMOVED***
	if err = ejvr.ensureElementValue(bsontype.Symbol, 0, "ReadSymbol"); err != nil ***REMOVED***
		return "", err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Symbol)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	symbol, err = v.parseSymbol()

	ejvr.pop()
	return symbol, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadTimestamp() (t uint32, i uint32, err error) ***REMOVED***
	if err = ejvr.ensureElementValue(bsontype.Timestamp, 0, "ReadTimestamp"); err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Timestamp)
	if err != nil ***REMOVED***
		return 0, 0, err
	***REMOVED***

	t, i, err = v.parseTimestamp()

	ejvr.pop()
	return t, i, err
***REMOVED***

func (ejvr *extJSONValueReader) ReadUndefined() error ***REMOVED***
	if err := ejvr.ensureElementValue(bsontype.Undefined, 0, "ReadUndefined"); err != nil ***REMOVED***
		return err
	***REMOVED***

	v, err := ejvr.p.readValue(bsontype.Undefined)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = v.parseUndefined()

	ejvr.pop()
	return err
***REMOVED***

func (ejvr *extJSONValueReader) ReadElement() (string, ValueReader, error) ***REMOVED***
	switch ejvr.stack[ejvr.frame].mode ***REMOVED***
	case mTopLevel, mDocument, mCodeWithScope:
	default:
		return "", nil, ejvr.invalidTransitionErr(mElement, "ReadElement", []mode***REMOVED***mTopLevel, mDocument, mCodeWithScope***REMOVED***)
	***REMOVED***

	name, t, err := ejvr.p.readKey()

	if err != nil ***REMOVED***
		if err == ErrEOD ***REMOVED***
			if ejvr.stack[ejvr.frame].mode == mCodeWithScope ***REMOVED***
				_, err := ejvr.p.peekType()
				if err != nil ***REMOVED***
					return "", nil, err
				***REMOVED***
			***REMOVED***

			ejvr.pop()
		***REMOVED***

		return "", nil, err
	***REMOVED***

	ejvr.push(mElement, t)
	return name, ejvr, nil
***REMOVED***

func (ejvr *extJSONValueReader) ReadValue() (ValueReader, error) ***REMOVED***
	switch ejvr.stack[ejvr.frame].mode ***REMOVED***
	case mArray:
	default:
		return nil, ejvr.invalidTransitionErr(mValue, "ReadValue", []mode***REMOVED***mArray***REMOVED***)
	***REMOVED***

	t, err := ejvr.p.peekType()
	if err != nil ***REMOVED***
		if err == ErrEOA ***REMOVED***
			ejvr.pop()
		***REMOVED***

		return nil, err
	***REMOVED***

	ejvr.push(mValue, t)
	return ejvr, nil
***REMOVED***
