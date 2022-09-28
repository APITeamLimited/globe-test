// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonrw

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ExtJSONValueWriterPool is a pool for ExtJSON ValueWriters.
type ExtJSONValueWriterPool struct ***REMOVED***
	pool sync.Pool
***REMOVED***

// NewExtJSONValueWriterPool creates a new pool for ValueWriter instances that write to ExtJSON.
func NewExtJSONValueWriterPool() *ExtJSONValueWriterPool ***REMOVED***
	return &ExtJSONValueWriterPool***REMOVED***
		pool: sync.Pool***REMOVED***
			New: func() interface***REMOVED******REMOVED*** ***REMOVED***
				return new(extJSONValueWriter)
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Get retrieves a ExtJSON ValueWriter from the pool and resets it to use w as the destination.
func (bvwp *ExtJSONValueWriterPool) Get(w io.Writer, canonical, escapeHTML bool) ValueWriter ***REMOVED***
	vw := bvwp.pool.Get().(*extJSONValueWriter)
	if writer, ok := w.(*SliceWriter); ok ***REMOVED***
		vw.reset(*writer, canonical, escapeHTML)
		vw.w = writer
		return vw
	***REMOVED***
	vw.buf = vw.buf[:0]
	vw.w = w
	return vw
***REMOVED***

// Put inserts a ValueWriter into the pool. If the ValueWriter is not a ExtJSON ValueWriter, nothing
// happens and ok will be false.
func (bvwp *ExtJSONValueWriterPool) Put(vw ValueWriter) (ok bool) ***REMOVED***
	bvw, ok := vw.(*extJSONValueWriter)
	if !ok ***REMOVED***
		return false
	***REMOVED***

	if _, ok := bvw.w.(*SliceWriter); ok ***REMOVED***
		bvw.buf = nil
	***REMOVED***
	bvw.w = nil

	bvwp.pool.Put(bvw)
	return true
***REMOVED***

type ejvwState struct ***REMOVED***
	mode mode
***REMOVED***

type extJSONValueWriter struct ***REMOVED***
	w   io.Writer
	buf []byte

	stack      []ejvwState
	frame      int64
	canonical  bool
	escapeHTML bool
***REMOVED***

// NewExtJSONValueWriter creates a ValueWriter that writes Extended JSON to w.
func NewExtJSONValueWriter(w io.Writer, canonical, escapeHTML bool) (ValueWriter, error) ***REMOVED***
	if w == nil ***REMOVED***
		return nil, errNilWriter
	***REMOVED***

	return newExtJSONWriter(w, canonical, escapeHTML), nil
***REMOVED***

func newExtJSONWriter(w io.Writer, canonical, escapeHTML bool) *extJSONValueWriter ***REMOVED***
	stack := make([]ejvwState, 1, 5)
	stack[0] = ejvwState***REMOVED***mode: mTopLevel***REMOVED***

	return &extJSONValueWriter***REMOVED***
		w:          w,
		buf:        []byte***REMOVED******REMOVED***,
		stack:      stack,
		canonical:  canonical,
		escapeHTML: escapeHTML,
	***REMOVED***
***REMOVED***

func newExtJSONWriterFromSlice(buf []byte, canonical, escapeHTML bool) *extJSONValueWriter ***REMOVED***
	stack := make([]ejvwState, 1, 5)
	stack[0] = ejvwState***REMOVED***mode: mTopLevel***REMOVED***

	return &extJSONValueWriter***REMOVED***
		buf:        buf,
		stack:      stack,
		canonical:  canonical,
		escapeHTML: escapeHTML,
	***REMOVED***
***REMOVED***

func (ejvw *extJSONValueWriter) reset(buf []byte, canonical, escapeHTML bool) ***REMOVED***
	if ejvw.stack == nil ***REMOVED***
		ejvw.stack = make([]ejvwState, 1, 5)
	***REMOVED***

	ejvw.stack = ejvw.stack[:1]
	ejvw.stack[0] = ejvwState***REMOVED***mode: mTopLevel***REMOVED***
	ejvw.canonical = canonical
	ejvw.escapeHTML = escapeHTML
	ejvw.frame = 0
	ejvw.buf = buf
	ejvw.w = nil
***REMOVED***

func (ejvw *extJSONValueWriter) advanceFrame() ***REMOVED***
	if ejvw.frame+1 >= int64(len(ejvw.stack)) ***REMOVED*** // We need to grow the stack
		length := len(ejvw.stack)
		if length+1 >= cap(ejvw.stack) ***REMOVED***
			// double it
			buf := make([]ejvwState, 2*cap(ejvw.stack)+1)
			copy(buf, ejvw.stack)
			ejvw.stack = buf
		***REMOVED***
		ejvw.stack = ejvw.stack[:length+1]
	***REMOVED***
	ejvw.frame++
***REMOVED***

func (ejvw *extJSONValueWriter) push(m mode) ***REMOVED***
	ejvw.advanceFrame()

	ejvw.stack[ejvw.frame].mode = m
***REMOVED***

func (ejvw *extJSONValueWriter) pop() ***REMOVED***
	switch ejvw.stack[ejvw.frame].mode ***REMOVED***
	case mElement, mValue:
		ejvw.frame--
	case mDocument, mArray, mCodeWithScope:
		ejvw.frame -= 2 // we pop twice to jump over the mElement: mDocument -> mElement -> mDocument/mTopLevel/etc...
	***REMOVED***
***REMOVED***

func (ejvw *extJSONValueWriter) invalidTransitionErr(destination mode, name string, modes []mode) error ***REMOVED***
	te := TransitionError***REMOVED***
		name:        name,
		current:     ejvw.stack[ejvw.frame].mode,
		destination: destination,
		modes:       modes,
		action:      "write",
	***REMOVED***
	if ejvw.frame != 0 ***REMOVED***
		te.parent = ejvw.stack[ejvw.frame-1].mode
	***REMOVED***
	return te
***REMOVED***

func (ejvw *extJSONValueWriter) ensureElementValue(destination mode, callerName string, addmodes ...mode) error ***REMOVED***
	switch ejvw.stack[ejvw.frame].mode ***REMOVED***
	case mElement, mValue:
	default:
		modes := []mode***REMOVED***mElement, mValue***REMOVED***
		if addmodes != nil ***REMOVED***
			modes = append(modes, addmodes...)
		***REMOVED***
		return ejvw.invalidTransitionErr(destination, callerName, modes)
	***REMOVED***

	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) writeExtendedSingleValue(key string, value string, quotes bool) ***REMOVED***
	var s string
	if quotes ***REMOVED***
		s = fmt.Sprintf(`***REMOVED***"$%s":"%s"***REMOVED***`, key, value)
	***REMOVED*** else ***REMOVED***
		s = fmt.Sprintf(`***REMOVED***"$%s":%s***REMOVED***`, key, value)
	***REMOVED***

	ejvw.buf = append(ejvw.buf, []byte(s)...)
***REMOVED***

func (ejvw *extJSONValueWriter) WriteArray() (ArrayWriter, error) ***REMOVED***
	if err := ejvw.ensureElementValue(mArray, "WriteArray"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ejvw.buf = append(ejvw.buf, '[')

	ejvw.push(mArray)
	return ejvw, nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteBinary(b []byte) error ***REMOVED***
	return ejvw.WriteBinaryWithSubtype(b, 0x00)
***REMOVED***

func (ejvw *extJSONValueWriter) WriteBinaryWithSubtype(b []byte, btype byte) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteBinaryWithSubtype"); err != nil ***REMOVED***
		return err
	***REMOVED***

	var buf bytes.Buffer
	buf.WriteString(`***REMOVED***"$binary":***REMOVED***"base64":"`)
	buf.WriteString(base64.StdEncoding.EncodeToString(b))
	buf.WriteString(fmt.Sprintf(`","subType":"%02x"***REMOVED******REMOVED***,`, btype))

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteBoolean(b bool) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteBoolean"); err != nil ***REMOVED***
		return err
	***REMOVED***

	ejvw.buf = append(ejvw.buf, []byte(strconv.FormatBool(b))...)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteCodeWithScope(code string) (DocumentWriter, error) ***REMOVED***
	if err := ejvw.ensureElementValue(mCodeWithScope, "WriteCodeWithScope"); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var buf bytes.Buffer
	buf.WriteString(`***REMOVED***"$code":`)
	writeStringWithEscapes(code, &buf, ejvw.escapeHTML)
	buf.WriteString(`,"$scope":***REMOVED***`)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.push(mCodeWithScope)
	return ejvw, nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteDBPointer(ns string, oid primitive.ObjectID) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteDBPointer"); err != nil ***REMOVED***
		return err
	***REMOVED***

	var buf bytes.Buffer
	buf.WriteString(`***REMOVED***"$dbPointer":***REMOVED***"$ref":"`)
	buf.WriteString(ns)
	buf.WriteString(`","$id":***REMOVED***"$oid":"`)
	buf.WriteString(oid.Hex())
	buf.WriteString(`"***REMOVED******REMOVED******REMOVED***,`)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteDateTime(dt int64) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteDateTime"); err != nil ***REMOVED***
		return err
	***REMOVED***

	t := time.Unix(dt/1e3, dt%1e3*1e6).UTC()

	if ejvw.canonical || t.Year() < 1970 || t.Year() > 9999 ***REMOVED***
		s := fmt.Sprintf(`***REMOVED***"$numberLong":"%d"***REMOVED***`, dt)
		ejvw.writeExtendedSingleValue("date", s, false)
	***REMOVED*** else ***REMOVED***
		ejvw.writeExtendedSingleValue("date", t.Format(rfc3339Milli), true)
	***REMOVED***

	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteDecimal128(d primitive.Decimal128) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteDecimal128"); err != nil ***REMOVED***
		return err
	***REMOVED***

	ejvw.writeExtendedSingleValue("numberDecimal", d.String(), true)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteDocument() (DocumentWriter, error) ***REMOVED***
	if ejvw.stack[ejvw.frame].mode == mTopLevel ***REMOVED***
		ejvw.buf = append(ejvw.buf, '***REMOVED***')
		return ejvw, nil
	***REMOVED***

	if err := ejvw.ensureElementValue(mDocument, "WriteDocument", mTopLevel); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ejvw.buf = append(ejvw.buf, '***REMOVED***')
	ejvw.push(mDocument)
	return ejvw, nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteDouble(f float64) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteDouble"); err != nil ***REMOVED***
		return err
	***REMOVED***

	s := formatDouble(f)

	if ejvw.canonical ***REMOVED***
		ejvw.writeExtendedSingleValue("numberDouble", s, true)
	***REMOVED*** else ***REMOVED***
		switch s ***REMOVED***
		case "Infinity":
			fallthrough
		case "-Infinity":
			fallthrough
		case "NaN":
			s = fmt.Sprintf(`***REMOVED***"$numberDouble":"%s"***REMOVED***`, s)
		***REMOVED***
		ejvw.buf = append(ejvw.buf, []byte(s)...)
	***REMOVED***

	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteInt32(i int32) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteInt32"); err != nil ***REMOVED***
		return err
	***REMOVED***

	s := strconv.FormatInt(int64(i), 10)

	if ejvw.canonical ***REMOVED***
		ejvw.writeExtendedSingleValue("numberInt", s, true)
	***REMOVED*** else ***REMOVED***
		ejvw.buf = append(ejvw.buf, []byte(s)...)
	***REMOVED***

	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteInt64(i int64) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteInt64"); err != nil ***REMOVED***
		return err
	***REMOVED***

	s := strconv.FormatInt(i, 10)

	if ejvw.canonical ***REMOVED***
		ejvw.writeExtendedSingleValue("numberLong", s, true)
	***REMOVED*** else ***REMOVED***
		ejvw.buf = append(ejvw.buf, []byte(s)...)
	***REMOVED***

	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteJavascript(code string) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteJavascript"); err != nil ***REMOVED***
		return err
	***REMOVED***

	var buf bytes.Buffer
	writeStringWithEscapes(code, &buf, ejvw.escapeHTML)

	ejvw.writeExtendedSingleValue("code", buf.String(), false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteMaxKey() error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteMaxKey"); err != nil ***REMOVED***
		return err
	***REMOVED***

	ejvw.writeExtendedSingleValue("maxKey", "1", false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteMinKey() error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteMinKey"); err != nil ***REMOVED***
		return err
	***REMOVED***

	ejvw.writeExtendedSingleValue("minKey", "1", false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteNull() error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteNull"); err != nil ***REMOVED***
		return err
	***REMOVED***

	ejvw.buf = append(ejvw.buf, []byte("null")...)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteObjectID(oid primitive.ObjectID) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteObjectID"); err != nil ***REMOVED***
		return err
	***REMOVED***

	ejvw.writeExtendedSingleValue("oid", oid.Hex(), true)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteRegex(pattern string, options string) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteRegex"); err != nil ***REMOVED***
		return err
	***REMOVED***

	var buf bytes.Buffer
	buf.WriteString(`***REMOVED***"$regularExpression":***REMOVED***"pattern":`)
	writeStringWithEscapes(pattern, &buf, ejvw.escapeHTML)
	buf.WriteString(`,"options":"`)
	buf.WriteString(sortStringAlphebeticAscending(options))
	buf.WriteString(`"***REMOVED******REMOVED***,`)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteString(s string) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteString"); err != nil ***REMOVED***
		return err
	***REMOVED***

	var buf bytes.Buffer
	writeStringWithEscapes(s, &buf, ejvw.escapeHTML)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteSymbol(symbol string) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteSymbol"); err != nil ***REMOVED***
		return err
	***REMOVED***

	var buf bytes.Buffer
	writeStringWithEscapes(symbol, &buf, ejvw.escapeHTML)

	ejvw.writeExtendedSingleValue("symbol", buf.String(), false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteTimestamp(t uint32, i uint32) error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteTimestamp"); err != nil ***REMOVED***
		return err
	***REMOVED***

	var buf bytes.Buffer
	buf.WriteString(`***REMOVED***"$timestamp":***REMOVED***"t":`)
	buf.WriteString(strconv.FormatUint(uint64(t), 10))
	buf.WriteString(`,"i":`)
	buf.WriteString(strconv.FormatUint(uint64(i), 10))
	buf.WriteString(`***REMOVED******REMOVED***,`)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteUndefined() error ***REMOVED***
	if err := ejvw.ensureElementValue(mode(0), "WriteUndefined"); err != nil ***REMOVED***
		return err
	***REMOVED***

	ejvw.writeExtendedSingleValue("undefined", "true", false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteDocumentElement(key string) (ValueWriter, error) ***REMOVED***
	switch ejvw.stack[ejvw.frame].mode ***REMOVED***
	case mDocument, mTopLevel, mCodeWithScope:
		var buf bytes.Buffer
		writeStringWithEscapes(key, &buf, ejvw.escapeHTML)

		ejvw.buf = append(ejvw.buf, []byte(fmt.Sprintf(`%s:`, buf.String()))...)
		ejvw.push(mElement)
	default:
		return nil, ejvw.invalidTransitionErr(mElement, "WriteDocumentElement", []mode***REMOVED***mDocument, mTopLevel, mCodeWithScope***REMOVED***)
	***REMOVED***

	return ejvw, nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteDocumentEnd() error ***REMOVED***
	switch ejvw.stack[ejvw.frame].mode ***REMOVED***
	case mDocument, mTopLevel, mCodeWithScope:
	default:
		return fmt.Errorf("incorrect mode to end document: %s", ejvw.stack[ejvw.frame].mode)
	***REMOVED***

	// close the document
	if ejvw.buf[len(ejvw.buf)-1] == ',' ***REMOVED***
		ejvw.buf[len(ejvw.buf)-1] = '***REMOVED***'
	***REMOVED*** else ***REMOVED***
		ejvw.buf = append(ejvw.buf, '***REMOVED***')
	***REMOVED***

	switch ejvw.stack[ejvw.frame].mode ***REMOVED***
	case mCodeWithScope:
		ejvw.buf = append(ejvw.buf, '***REMOVED***')
		fallthrough
	case mDocument:
		ejvw.buf = append(ejvw.buf, ',')
	case mTopLevel:
		if ejvw.w != nil ***REMOVED***
			if _, err := ejvw.w.Write(ejvw.buf); err != nil ***REMOVED***
				return err
			***REMOVED***
			ejvw.buf = ejvw.buf[:0]
		***REMOVED***
	***REMOVED***

	ejvw.pop()
	return nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteArrayElement() (ValueWriter, error) ***REMOVED***
	switch ejvw.stack[ejvw.frame].mode ***REMOVED***
	case mArray:
		ejvw.push(mValue)
	default:
		return nil, ejvw.invalidTransitionErr(mValue, "WriteArrayElement", []mode***REMOVED***mArray***REMOVED***)
	***REMOVED***

	return ejvw, nil
***REMOVED***

func (ejvw *extJSONValueWriter) WriteArrayEnd() error ***REMOVED***
	switch ejvw.stack[ejvw.frame].mode ***REMOVED***
	case mArray:
		// close the array
		if ejvw.buf[len(ejvw.buf)-1] == ',' ***REMOVED***
			ejvw.buf[len(ejvw.buf)-1] = ']'
		***REMOVED*** else ***REMOVED***
			ejvw.buf = append(ejvw.buf, ']')
		***REMOVED***

		ejvw.buf = append(ejvw.buf, ',')

		ejvw.pop()
	default:
		return fmt.Errorf("incorrect mode to end array: %s", ejvw.stack[ejvw.frame].mode)
	***REMOVED***

	return nil
***REMOVED***

func formatDouble(f float64) string ***REMOVED***
	var s string
	if math.IsInf(f, 1) ***REMOVED***
		s = "Infinity"
	***REMOVED*** else if math.IsInf(f, -1) ***REMOVED***
		s = "-Infinity"
	***REMOVED*** else if math.IsNaN(f) ***REMOVED***
		s = "NaN"
	***REMOVED*** else ***REMOVED***
		// Print exactly one decimalType place for integers; otherwise, print as many are necessary to
		// perfectly represent it.
		s = strconv.FormatFloat(f, 'G', -1, 64)
		if !strings.ContainsRune(s, 'E') && !strings.ContainsRune(s, '.') ***REMOVED***
			s += ".0"
		***REMOVED***
	***REMOVED***

	return s
***REMOVED***

var hexChars = "0123456789abcdef"

func writeStringWithEscapes(s string, buf *bytes.Buffer, escapeHTML bool) ***REMOVED***
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); ***REMOVED***
		if b := s[i]; b < utf8.RuneSelf ***REMOVED***
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) ***REMOVED***
				i++
				continue
			***REMOVED***
			if start < i ***REMOVED***
				buf.WriteString(s[start:i])
			***REMOVED***
			switch b ***REMOVED***
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			case '\b':
				buf.WriteByte('\\')
				buf.WriteByte('b')
			case '\f':
				buf.WriteByte('\\')
				buf.WriteByte('f')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hexChars[b>>4])
				buf.WriteByte(hexChars[b&0xF])
			***REMOVED***
			i++
			start = i
			continue
		***REMOVED***
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 ***REMOVED***
			if start < i ***REMOVED***
				buf.WriteString(s[start:i])
			***REMOVED***
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		***REMOVED***
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' ***REMOVED***
			if start < i ***REMOVED***
				buf.WriteString(s[start:i])
			***REMOVED***
			buf.WriteString(`\u202`)
			buf.WriteByte(hexChars[c&0xF])
			i += size
			start = i
			continue
		***REMOVED***
		i += size
	***REMOVED***
	if start < len(s) ***REMOVED***
		buf.WriteString(s[start:])
	***REMOVED***
	buf.WriteByte('"')
***REMOVED***

type sortableString []rune

func (ss sortableString) Len() int ***REMOVED***
	return len(ss)
***REMOVED***

func (ss sortableString) Less(i, j int) bool ***REMOVED***
	return ss[i] < ss[j]
***REMOVED***

func (ss sortableString) Swap(i, j int) ***REMOVED***
	oldI := ss[i]
	ss[i] = ss[j]
	ss[j] = oldI
***REMOVED***

func sortStringAlphebeticAscending(s string) string ***REMOVED***
	ss := sortableString([]rune(s))
	sort.Sort(ss)
	return string([]rune(ss))
***REMOVED***
