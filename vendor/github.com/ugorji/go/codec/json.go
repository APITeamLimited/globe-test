// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// By default, this json support uses base64 encoding for bytes, because you cannot
// store and read any arbitrary string in json (only unicode).
// However, the user can configre how to encode/decode bytes.
//
// This library specifically supports UTF-8 for encoding and decoding only.
//
// Note that the library will happily encode/decode things which are not valid
// json e.g. a map[int64]string. We do it for consistency. With valid json,
// we will encode and decode appropriately.
// Users can specify their map type if necessary to force it.
//
// Note:
//   - we cannot use strconv.Quote and strconv.Unquote because json quotes/unquotes differently.
//     We implement it here.

// Top-level methods of json(End|Dec)Driver (which are implementations of (en|de)cDriver
// MUST not call one-another.

import (
	"bytes"
	"encoding/base64"
	"math"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

//--------------------------------

var jsonLiterals = [...]byte***REMOVED***
	'"', 't', 'r', 'u', 'e', '"',
	'"', 'f', 'a', 'l', 's', 'e', '"',
	'"', 'n', 'u', 'l', 'l', '"',
***REMOVED***

const (
	jsonLitTrueQ  = 0
	jsonLitTrue   = 1
	jsonLitFalseQ = 6
	jsonLitFalse  = 7
	jsonLitNullQ  = 13
	jsonLitNull   = 14
)

const (
	jsonU4Chk2 = '0'
	jsonU4Chk1 = 'a' - 10
	jsonU4Chk0 = 'A' - 10

	jsonScratchArrayLen = 64
)

const (
	// If !jsonValidateSymbols, decoding will be faster, by skipping some checks:
	//   - If we see first character of null, false or true,
	//     do not validate subsequent characters.
	//   - e.g. if we see a n, assume null and skip next 3 characters,
	//     and do not validate they are ull.
	// P.S. Do not expect a significant decoding boost from this.
	jsonValidateSymbols = true

	jsonSpacesOrTabsLen = 128

	jsonAlwaysReturnInternString = false
)

var (
	// jsonTabs and jsonSpaces are used as caches for indents
	jsonTabs, jsonSpaces [jsonSpacesOrTabsLen]byte

	jsonCharHtmlSafeSet   bitset128
	jsonCharSafeSet       bitset128
	jsonCharWhitespaceSet bitset256
	jsonNumSet            bitset256
)

func init() ***REMOVED***
	for i := 0; i < jsonSpacesOrTabsLen; i++ ***REMOVED***
		jsonSpaces[i] = ' '
		jsonTabs[i] = '\t'
	***REMOVED***

	// populate the safe values as true: note: ASCII control characters are (0-31)
	// jsonCharSafeSet:     all true except (0-31) " \
	// jsonCharHtmlSafeSet: all true except (0-31) " \ < > &
	var i byte
	for i = 32; i < utf8.RuneSelf; i++ ***REMOVED***
		switch i ***REMOVED***
		case '"', '\\':
		case '<', '>', '&':
			jsonCharSafeSet.set(i) // = true
		default:
			jsonCharSafeSet.set(i)
			jsonCharHtmlSafeSet.set(i)
		***REMOVED***
	***REMOVED***
	for i = 0; i <= utf8.RuneSelf; i++ ***REMOVED***
		switch i ***REMOVED***
		case ' ', '\t', '\r', '\n':
			jsonCharWhitespaceSet.set(i)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'e', 'E', '.', '+', '-':
			jsonNumSet.set(i)
		***REMOVED***
	***REMOVED***
***REMOVED***

// ----------------

type jsonEncDriverTypical struct ***REMOVED***
	w encWriter
	// w  *encWriterSwitch
	b  *[jsonScratchArrayLen]byte
	tw bool // term white space
	c  containerState
***REMOVED***

func (e *jsonEncDriverTypical) typical() ***REMOVED******REMOVED***

func (e *jsonEncDriverTypical) reset(ee *jsonEncDriver) ***REMOVED***
	e.w = ee.ew
	// e.w = &ee.e.encWriterSwitch
	e.b = &ee.b
	e.tw = ee.h.TermWhitespace
	e.c = 0
***REMOVED***

func (e *jsonEncDriverTypical) WriteArrayStart(length int) ***REMOVED***
	e.w.writen1('[')
	e.c = containerArrayStart
***REMOVED***

func (e *jsonEncDriverTypical) WriteArrayElem() ***REMOVED***
	if e.c != containerArrayStart ***REMOVED***
		e.w.writen1(',')
	***REMOVED***
	e.c = containerArrayElem
***REMOVED***

func (e *jsonEncDriverTypical) WriteArrayEnd() ***REMOVED***
	e.w.writen1(']')
	e.c = containerArrayEnd
***REMOVED***

func (e *jsonEncDriverTypical) WriteMapStart(length int) ***REMOVED***
	e.w.writen1('***REMOVED***')
	e.c = containerMapStart
***REMOVED***

func (e *jsonEncDriverTypical) WriteMapElemKey() ***REMOVED***
	if e.c != containerMapStart ***REMOVED***
		e.w.writen1(',')
	***REMOVED***
	e.c = containerMapKey
***REMOVED***

func (e *jsonEncDriverTypical) WriteMapElemValue() ***REMOVED***
	e.w.writen1(':')
	e.c = containerMapValue
***REMOVED***

func (e *jsonEncDriverTypical) WriteMapEnd() ***REMOVED***
	e.w.writen1('***REMOVED***')
	e.c = containerMapEnd
***REMOVED***

func (e *jsonEncDriverTypical) EncodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.w.writeb(jsonLiterals[jsonLitTrue : jsonLitTrue+4])
	***REMOVED*** else ***REMOVED***
		e.w.writeb(jsonLiterals[jsonLitFalse : jsonLitFalse+5])
	***REMOVED***
***REMOVED***

func (e *jsonEncDriverTypical) EncodeFloat64(f float64) ***REMOVED***
	fmt, prec := jsonFloatStrconvFmtPrec(f)
	e.w.writeb(strconv.AppendFloat(e.b[:0], f, fmt, prec, 64))
***REMOVED***

func (e *jsonEncDriverTypical) EncodeInt(v int64) ***REMOVED***
	e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
***REMOVED***

func (e *jsonEncDriverTypical) EncodeUint(v uint64) ***REMOVED***
	e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
***REMOVED***

func (e *jsonEncDriverTypical) EncodeFloat32(f float32) ***REMOVED***
	e.EncodeFloat64(float64(f))
***REMOVED***

func (e *jsonEncDriverTypical) atEndOfEncode() ***REMOVED***
	if e.tw ***REMOVED***
		e.w.writen1(' ')
	***REMOVED***
***REMOVED***

// ----------------

type jsonEncDriverGeneric struct ***REMOVED***
	w encWriter // encWriter // *encWriterSwitch
	b *[jsonScratchArrayLen]byte
	c containerState
	// ds string // indent string
	di int8    // indent per
	d  bool    // indenting?
	dt bool    // indent using tabs
	dl uint16  // indent level
	ks bool    // map key as string
	is byte    // integer as string
	tw bool    // term white space
	_  [7]byte // padding
***REMOVED***

// indent is done as below:
//   - newline and indent are added before each mapKey or arrayElem
//   - newline and indent are added before each ending,
//     except there was no entry (so we can have ***REMOVED******REMOVED*** or [])

func (e *jsonEncDriverGeneric) reset(ee *jsonEncDriver) ***REMOVED***
	e.w = ee.ew
	e.b = &ee.b
	e.tw = ee.h.TermWhitespace
	e.c = 0
	e.d, e.dt, e.dl, e.di = false, false, 0, 0
	h := ee.h
	if h.Indent > 0 ***REMOVED***
		e.d = true
		e.di = int8(h.Indent)
	***REMOVED*** else if h.Indent < 0 ***REMOVED***
		e.d = true
		e.dt = true
		e.di = int8(-h.Indent)
	***REMOVED***
	e.ks = h.MapKeyAsString
	e.is = h.IntegerAsString
***REMOVED***

func (e *jsonEncDriverGeneric) WriteArrayStart(length int) ***REMOVED***
	if e.d ***REMOVED***
		e.dl++
	***REMOVED***
	e.w.writen1('[')
	e.c = containerArrayStart
***REMOVED***

func (e *jsonEncDriverGeneric) WriteArrayElem() ***REMOVED***
	if e.c != containerArrayStart ***REMOVED***
		e.w.writen1(',')
	***REMOVED***
	if e.d ***REMOVED***
		e.writeIndent()
	***REMOVED***
	e.c = containerArrayElem
***REMOVED***

func (e *jsonEncDriverGeneric) WriteArrayEnd() ***REMOVED***
	if e.d ***REMOVED***
		e.dl--
		if e.c != containerArrayStart ***REMOVED***
			e.writeIndent()
		***REMOVED***
	***REMOVED***
	e.w.writen1(']')
	e.c = containerArrayEnd
***REMOVED***

func (e *jsonEncDriverGeneric) WriteMapStart(length int) ***REMOVED***
	if e.d ***REMOVED***
		e.dl++
	***REMOVED***
	e.w.writen1('***REMOVED***')
	e.c = containerMapStart
***REMOVED***

func (e *jsonEncDriverGeneric) WriteMapElemKey() ***REMOVED***
	if e.c != containerMapStart ***REMOVED***
		e.w.writen1(',')
	***REMOVED***
	if e.d ***REMOVED***
		e.writeIndent()
	***REMOVED***
	e.c = containerMapKey
***REMOVED***

func (e *jsonEncDriverGeneric) WriteMapElemValue() ***REMOVED***
	if e.d ***REMOVED***
		e.w.writen2(':', ' ')
	***REMOVED*** else ***REMOVED***
		e.w.writen1(':')
	***REMOVED***
	e.c = containerMapValue
***REMOVED***

func (e *jsonEncDriverGeneric) WriteMapEnd() ***REMOVED***
	if e.d ***REMOVED***
		e.dl--
		if e.c != containerMapStart ***REMOVED***
			e.writeIndent()
		***REMOVED***
	***REMOVED***
	e.w.writen1('***REMOVED***')
	e.c = containerMapEnd
***REMOVED***

func (e *jsonEncDriverGeneric) writeIndent() ***REMOVED***
	e.w.writen1('\n')
	x := int(e.di) * int(e.dl)
	if e.dt ***REMOVED***
		for x > jsonSpacesOrTabsLen ***REMOVED***
			e.w.writeb(jsonTabs[:])
			x -= jsonSpacesOrTabsLen
		***REMOVED***
		e.w.writeb(jsonTabs[:x])
	***REMOVED*** else ***REMOVED***
		for x > jsonSpacesOrTabsLen ***REMOVED***
			e.w.writeb(jsonSpaces[:])
			x -= jsonSpacesOrTabsLen
		***REMOVED***
		e.w.writeb(jsonSpaces[:x])
	***REMOVED***
***REMOVED***

func (e *jsonEncDriverGeneric) EncodeBool(b bool) ***REMOVED***
	if e.ks && e.c == containerMapKey ***REMOVED***
		if b ***REMOVED***
			e.w.writeb(jsonLiterals[jsonLitTrueQ : jsonLitTrueQ+6])
		***REMOVED*** else ***REMOVED***
			e.w.writeb(jsonLiterals[jsonLitFalseQ : jsonLitFalseQ+7])
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if b ***REMOVED***
			e.w.writeb(jsonLiterals[jsonLitTrue : jsonLitTrue+4])
		***REMOVED*** else ***REMOVED***
			e.w.writeb(jsonLiterals[jsonLitFalse : jsonLitFalse+5])
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *jsonEncDriverGeneric) EncodeFloat64(f float64) ***REMOVED***
	// instead of using 'g', specify whether to use 'e' or 'f'
	fmt, prec := jsonFloatStrconvFmtPrec(f)

	var blen int
	if e.ks && e.c == containerMapKey ***REMOVED***
		blen = 2 + len(strconv.AppendFloat(e.b[1:1], f, fmt, prec, 64))
		e.b[0] = '"'
		e.b[blen-1] = '"'
	***REMOVED*** else ***REMOVED***
		blen = len(strconv.AppendFloat(e.b[:0], f, fmt, prec, 64))
	***REMOVED***
	e.w.writeb(e.b[:blen])
***REMOVED***

func (e *jsonEncDriverGeneric) EncodeInt(v int64) ***REMOVED***
	x := e.is
	if x == 'A' || x == 'L' && (v > 1<<53 || v < -(1<<53)) || (e.ks && e.c == containerMapKey) ***REMOVED***
		blen := 2 + len(strconv.AppendInt(e.b[1:1], v, 10))
		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.w.writeb(e.b[:blen])
		return
	***REMOVED***
	e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
***REMOVED***

func (e *jsonEncDriverGeneric) EncodeUint(v uint64) ***REMOVED***
	x := e.is
	if x == 'A' || x == 'L' && v > 1<<53 || (e.ks && e.c == containerMapKey) ***REMOVED***
		blen := 2 + len(strconv.AppendUint(e.b[1:1], v, 10))
		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.w.writeb(e.b[:blen])
		return
	***REMOVED***
	e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
***REMOVED***

func (e *jsonEncDriverGeneric) EncodeFloat32(f float32) ***REMOVED***
	// e.encodeFloat(float64(f), 32)
	// always encode all floats as IEEE 64-bit floating point.
	// It also ensures that we can decode in full precision even if into a float32,
	// as what is written is always to float64 precision.
	e.EncodeFloat64(float64(f))
***REMOVED***

func (e *jsonEncDriverGeneric) atEndOfEncode() ***REMOVED***
	if e.tw ***REMOVED***
		if e.d ***REMOVED***
			e.w.writen1('\n')
		***REMOVED*** else ***REMOVED***
			e.w.writen1(' ')
		***REMOVED***
	***REMOVED***
***REMOVED***

// --------------------

type jsonEncDriver struct ***REMOVED***
	noBuiltInTypes
	e  *Encoder
	h  *JsonHandle
	ew encWriter // encWriter // *encWriterSwitch
	se extWrapper
	// ---- cpu cache line boundary?
	bs []byte // scratch
	// ---- cpu cache line boundary?
	b [jsonScratchArrayLen]byte // scratch (encode time,
***REMOVED***

func (e *jsonEncDriver) EncodeNil() ***REMOVED***
	// We always encode nil as just null (never in quotes)
	// This allows us to easily decode if a nil in the json stream
	// ie if initial token is n.
	e.ew.writeb(jsonLiterals[jsonLitNull : jsonLitNull+4])

	// if e.h.MapKeyAsString && e.c == containerMapKey ***REMOVED***
	// 	e.ew.writeb(jsonLiterals[jsonLitNullQ : jsonLitNullQ+6])
	// ***REMOVED*** else ***REMOVED***
	// 	e.ew.writeb(jsonLiterals[jsonLitNull : jsonLitNull+4])
	// ***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeTime(t time.Time) ***REMOVED***
	// Do NOT use MarshalJSON, as it allocates internally.
	// instead, we call AppendFormat directly, using our scratch buffer (e.b)
	if t.IsZero() ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		e.b[0] = '"'
		b := t.AppendFormat(e.b[1:1], time.RFC3339Nano)
		e.b[len(b)+1] = '"'
		e.ew.writeb(e.b[:len(b)+2])
	***REMOVED***
	// v, err := t.MarshalJSON(); if err != nil ***REMOVED*** e.e.error(err) ***REMOVED*** e.ew.writeb(v)
***REMOVED***

func (e *jsonEncDriver) EncodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext, en *Encoder) ***REMOVED***
	if v := ext.ConvertExt(rv); v == nil ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		en.encode(v)
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeRawExt(re *RawExt, en *Encoder) ***REMOVED***
	// only encodes re.Value (never re.Data)
	if re.Value == nil ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		en.encode(re.Value)
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeString(c charEncoding, v string) ***REMOVED***
	e.quoteStr(v)
***REMOVED***

func (e *jsonEncDriver) EncodeStringBytes(c charEncoding, v []byte) ***REMOVED***
	// if encoding raw bytes and RawBytesExt is configured, use it to encode
	if v == nil ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	if c == cRAW ***REMOVED***
		if e.se.InterfaceExt != nil ***REMOVED***
			e.EncodeExt(v, 0, &e.se, e.e)
			return
		***REMOVED***

		slen := base64.StdEncoding.EncodedLen(len(v))
		if cap(e.bs) >= slen+2 ***REMOVED***
			e.bs = e.bs[:slen+2]
		***REMOVED*** else ***REMOVED***
			e.bs = make([]byte, slen+2)
		***REMOVED***
		e.bs[0] = '"'
		base64.StdEncoding.Encode(e.bs[1:], v)
		e.bs[slen+1] = '"'
		e.ew.writeb(e.bs)
	***REMOVED*** else ***REMOVED***
		e.quoteStr(stringView(v))
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeAsis(v []byte) ***REMOVED***
	e.ew.writeb(v)
***REMOVED***

func (e *jsonEncDriver) quoteStr(s string) ***REMOVED***
	// adapted from std pkg encoding/json
	const hex = "0123456789abcdef"
	w := e.ew
	htmlasis := e.h.HTMLCharsAsIs
	w.writen1('"')
	var start int
	for i, slen := 0, len(s); i < slen; ***REMOVED***
		// encode all bytes < 0x20 (except \r, \n).
		// also encode < > & to prevent security holes when served to some browsers.
		if b := s[i]; b < utf8.RuneSelf ***REMOVED***
			// if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' ***REMOVED***
			// if (htmlasis && jsonCharSafeSet.isset(b)) || jsonCharHtmlSafeSet.isset(b) ***REMOVED***
			if jsonCharHtmlSafeSet.isset(b) || (htmlasis && jsonCharSafeSet.isset(b)) ***REMOVED***
				i++
				continue
			***REMOVED***
			if start < i ***REMOVED***
				w.writestr(s[start:i])
			***REMOVED***
			switch b ***REMOVED***
			case '\\', '"':
				w.writen2('\\', b)
			case '\n':
				w.writen2('\\', 'n')
			case '\r':
				w.writen2('\\', 'r')
			case '\b':
				w.writen2('\\', 'b')
			case '\f':
				w.writen2('\\', 'f')
			case '\t':
				w.writen2('\\', 't')
			default:
				w.writestr(`\u00`)
				w.writen2(hex[b>>4], hex[b&0xF])
			***REMOVED***
			i++
			start = i
			continue
		***REMOVED***
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 ***REMOVED***
			if start < i ***REMOVED***
				w.writestr(s[start:i])
			***REMOVED***
			w.writestr(`\ufffd`)
			i += size
			start = i
			continue
		***REMOVED***
		// U+2028 is LINE SEPARATOR. U+2029 is PARAGRAPH SEPARATOR.
		// Both technically valid JSON, but bomb on JSONP, so fix here unconditionally.
		if c == '\u2028' || c == '\u2029' ***REMOVED***
			if start < i ***REMOVED***
				w.writestr(s[start:i])
			***REMOVED***
			w.writestr(`\u202`)
			w.writen1(hex[c&0xF])
			i += size
			start = i
			continue
		***REMOVED***
		i += size
	***REMOVED***
	if start < len(s) ***REMOVED***
		w.writestr(s[start:])
	***REMOVED***
	w.writen1('"')
***REMOVED***

type jsonDecDriver struct ***REMOVED***
	noBuiltInTypes
	d  *Decoder
	h  *JsonHandle
	r  decReader // *decReaderSwitch // decReader
	se extWrapper

	// ---- writable fields during execution --- *try* to keep in sep cache line

	c containerState
	// tok is used to store the token read right after skipWhiteSpace.
	tok   uint8
	fnull bool    // found null from appendStringAsBytes
	bs    []byte  // scratch. Initialized from b. Used for parsing strings or numbers.
	bstr  [8]byte // scratch used for string \UXXX parsing
	// ---- cpu cache line boundary?
	b  [jsonScratchArrayLen]byte // scratch 1, used for parsing strings or numbers or time.Time
	b2 [jsonScratchArrayLen]byte // scratch 2, used only for readUntil, decNumBytes

	_ [3]uint64 // padding
	// n jsonNum
***REMOVED***

// func jsonIsWS(b byte) bool ***REMOVED***
// 	// return b == ' ' || b == '\t' || b == '\r' || b == '\n'
// 	return jsonCharWhitespaceSet.isset(b)
// ***REMOVED***

func (d *jsonDecDriver) uncacheRead() ***REMOVED***
	if d.tok != 0 ***REMOVED***
		d.r.unreadn1()
		d.tok = 0
	***REMOVED***
***REMOVED***

func (d *jsonDecDriver) ReadMapStart() int ***REMOVED***
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	const xc uint8 = '***REMOVED***'
	if d.tok != xc ***REMOVED***
		d.d.errorf("expect char '%c' but got char '%c'", xc, d.tok)
	***REMOVED***
	d.tok = 0
	d.c = containerMapStart
	return -1
***REMOVED***

func (d *jsonDecDriver) ReadArrayStart() int ***REMOVED***
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	const xc uint8 = '['
	if d.tok != xc ***REMOVED***
		d.d.errorf("expect char '%c' but got char '%c'", xc, d.tok)
	***REMOVED***
	d.tok = 0
	d.c = containerArrayStart
	return -1
***REMOVED***

func (d *jsonDecDriver) CheckBreak() bool ***REMOVED***
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	return d.tok == '***REMOVED***' || d.tok == ']'
***REMOVED***

// For the ReadXXX methods below, we could just delegate to helper functions
// readContainerState(c containerState, xc uint8, check bool)
// - ReadArrayElem would become:
//   readContainerState(containerArrayElem, ',', d.c != containerArrayStart)
//
// However, until mid-stack inlining (go 1.10?) comes, supporting inlining of
// oneliners, we explicitly write them all 5 out to elide the extra func call.
// TODO: For Go 1.10, if inlined, consider consolidating these.

func (d *jsonDecDriver) ReadArrayElem() ***REMOVED***
	const xc uint8 = ','
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	if d.c != containerArrayStart ***REMOVED***
		if d.tok != xc ***REMOVED***
			d.d.errorf("expect char '%c' but got char '%c'", xc, d.tok)
		***REMOVED***
		d.tok = 0
	***REMOVED***
	d.c = containerArrayElem
***REMOVED***

func (d *jsonDecDriver) ReadArrayEnd() ***REMOVED***
	const xc uint8 = ']'
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	if d.tok != xc ***REMOVED***
		d.d.errorf("expect char '%c' but got char '%c'", xc, d.tok)
	***REMOVED***
	d.tok = 0
	d.c = containerArrayEnd
***REMOVED***

func (d *jsonDecDriver) ReadMapElemKey() ***REMOVED***
	const xc uint8 = ','
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	if d.c != containerMapStart ***REMOVED***
		if d.tok != xc ***REMOVED***
			d.d.errorf("expect char '%c' but got char '%c'", xc, d.tok)
		***REMOVED***
		d.tok = 0
	***REMOVED***
	d.c = containerMapKey
***REMOVED***

func (d *jsonDecDriver) ReadMapElemValue() ***REMOVED***
	const xc uint8 = ':'
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	if d.tok != xc ***REMOVED***
		d.d.errorf("expect char '%c' but got char '%c'", xc, d.tok)
	***REMOVED***
	d.tok = 0
	d.c = containerMapValue
***REMOVED***

func (d *jsonDecDriver) ReadMapEnd() ***REMOVED***
	const xc uint8 = '***REMOVED***'
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	if d.tok != xc ***REMOVED***
		d.d.errorf("expect char '%c' but got char '%c'", xc, d.tok)
	***REMOVED***
	d.tok = 0
	d.c = containerMapEnd
***REMOVED***

func (d *jsonDecDriver) readLit(length, fromIdx uint8) ***REMOVED***
	bs := d.r.readx(int(length))
	d.tok = 0
	if jsonValidateSymbols && !bytes.Equal(bs, jsonLiterals[fromIdx:fromIdx+length]) ***REMOVED***
		d.d.errorf("expecting %s: got %s", jsonLiterals[fromIdx:fromIdx+length], bs)
		return
	***REMOVED***
***REMOVED***

func (d *jsonDecDriver) TryDecodeAsNil() bool ***REMOVED***
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	// we shouldn't try to see if "null" was here, right?
	// only the plain string: `null` denotes a nil (ie not quotes)
	if d.tok == 'n' ***REMOVED***
		d.readLit(3, jsonLitNull+1) // (n)ull
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *jsonDecDriver) DecodeBool() (v bool) ***REMOVED***
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	fquot := d.c == containerMapKey && d.tok == '"'
	if fquot ***REMOVED***
		d.tok = d.r.readn1()
	***REMOVED***
	switch d.tok ***REMOVED***
	case 'f':
		d.readLit(4, jsonLitFalse+1) // (f)alse
		// v = false
	case 't':
		d.readLit(3, jsonLitTrue+1) // (t)rue
		v = true
	default:
		d.d.errorf("decode bool: got first char %c", d.tok)
		// v = false // "unreachable"
	***REMOVED***
	if fquot ***REMOVED***
		d.r.readn1()
	***REMOVED***
	return
***REMOVED***

func (d *jsonDecDriver) DecodeTime() (t time.Time) ***REMOVED***
	// read string, and pass the string into json.unmarshal
	d.appendStringAsBytes()
	if d.fnull ***REMOVED***
		return
	***REMOVED***
	t, err := time.Parse(time.RFC3339, stringView(d.bs))
	if err != nil ***REMOVED***
		d.d.errorv(err)
	***REMOVED***
	return
***REMOVED***

func (d *jsonDecDriver) ContainerType() (vt valueType) ***REMOVED***
	// check container type by checking the first char
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***

	// optimize this, so we don't do 4 checks but do one computation.
	// return jsonContainerSet[d.tok]

	// ContainerType is mostly called for Map and Array,
	// so this conditional is good enough (max 2 checks typically)
	if b := d.tok; b == '***REMOVED***' ***REMOVED***
		return valueTypeMap
	***REMOVED*** else if b == '[' ***REMOVED***
		return valueTypeArray
	***REMOVED*** else if b == 'n' ***REMOVED***
		return valueTypeNil
	***REMOVED*** else if b == '"' ***REMOVED***
		return valueTypeString
	***REMOVED***
	return valueTypeUnset
***REMOVED***

func (d *jsonDecDriver) decNumBytes() (bs []byte) ***REMOVED***
	// stores num bytes in d.bs
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	if d.tok == '"' ***REMOVED***
		bs = d.r.readUntil(d.b2[:0], '"')
		bs = bs[:len(bs)-1]
	***REMOVED*** else ***REMOVED***
		d.r.unreadn1()
		bs = d.r.readTo(d.bs[:0], &jsonNumSet)
	***REMOVED***
	d.tok = 0
	return bs
***REMOVED***

func (d *jsonDecDriver) DecodeUint64() (u uint64) ***REMOVED***
	bs := d.decNumBytes()
	n, neg, badsyntax, overflow := jsonParseInteger(bs)
	if overflow ***REMOVED***
		d.d.errorf("overflow parsing unsigned integer: %s", bs)
	***REMOVED*** else if neg ***REMOVED***
		d.d.errorf("minus found parsing unsigned integer: %s", bs)
	***REMOVED*** else if badsyntax ***REMOVED***
		// fallback: try to decode as float, and cast
		n = d.decUint64ViaFloat(stringView(bs))
	***REMOVED***
	return n
***REMOVED***

func (d *jsonDecDriver) DecodeInt64() (i int64) ***REMOVED***
	const cutoff = uint64(1 << uint(64-1))
	bs := d.decNumBytes()
	n, neg, badsyntax, overflow := jsonParseInteger(bs)
	if overflow ***REMOVED***
		d.d.errorf("overflow parsing integer: %s", bs)
	***REMOVED*** else if badsyntax ***REMOVED***
		// d.d.errorf("invalid syntax for integer: %s", bs)
		// fallback: try to decode as float, and cast
		if neg ***REMOVED***
			n = d.decUint64ViaFloat(stringView(bs[1:]))
		***REMOVED*** else ***REMOVED***
			n = d.decUint64ViaFloat(stringView(bs))
		***REMOVED***
	***REMOVED***
	if neg ***REMOVED***
		if n > cutoff ***REMOVED***
			d.d.errorf("overflow parsing integer: %s", bs)
		***REMOVED***
		i = -(int64(n))
	***REMOVED*** else ***REMOVED***
		if n >= cutoff ***REMOVED***
			d.d.errorf("overflow parsing integer: %s", bs)
		***REMOVED***
		i = int64(n)
	***REMOVED***
	return
***REMOVED***

func (d *jsonDecDriver) decUint64ViaFloat(s string) (u uint64) ***REMOVED***
	f, err := strconv.ParseFloat(s, 64)
	if err != nil ***REMOVED***
		d.d.errorf("invalid syntax for integer: %s", s)
		// d.d.errorv(err)
	***REMOVED***
	fi, ff := math.Modf(f)
	if ff > 0 ***REMOVED***
		d.d.errorf("fractional part found parsing integer: %s", s)
	***REMOVED*** else if fi > float64(math.MaxUint64) ***REMOVED***
		d.d.errorf("overflow parsing integer: %s", s)
	***REMOVED***
	return uint64(fi)
***REMOVED***

func (d *jsonDecDriver) DecodeFloat64() (f float64) ***REMOVED***
	bs := d.decNumBytes()
	f, err := strconv.ParseFloat(stringView(bs), 64)
	if err != nil ***REMOVED***
		d.d.errorv(err)
	***REMOVED***
	return
***REMOVED***

func (d *jsonDecDriver) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) (realxtag uint64) ***REMOVED***
	if ext == nil ***REMOVED***
		re := rv.(*RawExt)
		re.Tag = xtag
		d.d.decode(&re.Value)
	***REMOVED*** else ***REMOVED***
		var v interface***REMOVED******REMOVED***
		d.d.decode(&v)
		ext.UpdateExt(rv, v)
	***REMOVED***
	return
***REMOVED***

func (d *jsonDecDriver) DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte) ***REMOVED***
	// if decoding into raw bytes, and the RawBytesExt is configured, use it to decode.
	if d.se.InterfaceExt != nil ***REMOVED***
		bsOut = bs
		d.DecodeExt(&bsOut, 0, &d.se)
		return
	***REMOVED***
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
	if d.tok == '[' ***REMOVED***
		bsOut, _ = fastpathTV.DecSliceUint8V(bs, true, d.d)
		return
	***REMOVED***
	d.appendStringAsBytes()
	// base64 encodes []byte***REMOVED******REMOVED*** as "", and we encode nil []byte as null.
	// Consequently, base64 should decode null as a nil []byte, and "" as an empty []byte***REMOVED******REMOVED***.
	// appendStringAsBytes returns a zero-len slice for both, so as not to reset d.bs.
	// However, it sets a fnull field to true, so we can check if a null was found.
	if len(d.bs) == 0 ***REMOVED***
		if d.fnull ***REMOVED***
			return nil
		***REMOVED***
		return []byte***REMOVED******REMOVED***
	***REMOVED***
	bs0 := d.bs
	slen := base64.StdEncoding.DecodedLen(len(bs0))
	if slen <= cap(bs) ***REMOVED***
		bsOut = bs[:slen]
	***REMOVED*** else if zerocopy && slen <= cap(d.b2) ***REMOVED***
		bsOut = d.b2[:slen]
	***REMOVED*** else ***REMOVED***
		bsOut = make([]byte, slen)
	***REMOVED***
	slen2, err := base64.StdEncoding.Decode(bsOut, bs0)
	if err != nil ***REMOVED***
		d.d.errorf("error decoding base64 binary '%s': %v", bs0, err)
		return nil
	***REMOVED***
	if slen != slen2 ***REMOVED***
		bsOut = bsOut[:slen2]
	***REMOVED***
	return
***REMOVED***

func (d *jsonDecDriver) DecodeString() (s string) ***REMOVED***
	d.appendStringAsBytes()
	return d.bsToString()
***REMOVED***

func (d *jsonDecDriver) DecodeStringAsBytes() (s []byte) ***REMOVED***
	d.appendStringAsBytes()
	return d.bs
***REMOVED***

func (d *jsonDecDriver) appendStringAsBytes() ***REMOVED***
	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***

	d.fnull = false
	if d.tok != '"' ***REMOVED***
		// d.d.errorf("expect char '%c' but got char '%c'", '"', d.tok)
		// handle non-string scalar: null, true, false or a number
		switch d.tok ***REMOVED***
		case 'n':
			d.readLit(3, jsonLitNull+1) // (n)ull
			d.bs = d.bs[:0]
			d.fnull = true
		case 'f':
			d.readLit(4, jsonLitFalse+1) // (f)alse
			d.bs = d.bs[:5]
			copy(d.bs, "false")
		case 't':
			d.readLit(3, jsonLitTrue+1) // (t)rue
			d.bs = d.bs[:4]
			copy(d.bs, "true")
		default:
			// try to parse a valid number
			bs := d.decNumBytes()
			if len(bs) <= cap(d.bs) ***REMOVED***
				d.bs = d.bs[:len(bs)]
			***REMOVED*** else ***REMOVED***
				d.bs = make([]byte, len(bs))
			***REMOVED***
			copy(d.bs, bs)
		***REMOVED***
		return
	***REMOVED***

	d.tok = 0
	r := d.r
	var cs = r.readUntil(d.b2[:0], '"')
	var cslen = len(cs)
	var c uint8
	v := d.bs[:0]
	// append on each byte seen can be expensive, so we just
	// keep track of where we last read a contiguous set of
	// non-special bytes (using cursor variable),
	// and when we see a special byte
	// e.g. end-of-slice, " or \,
	// we will append the full range into the v slice before proceeding
	for i, cursor := 0, 0; ; ***REMOVED***
		if i == cslen ***REMOVED***
			v = append(v, cs[cursor:]...)
			cs = r.readUntil(d.b2[:0], '"')
			cslen = len(cs)
			i, cursor = 0, 0
		***REMOVED***
		c = cs[i]
		if c == '"' ***REMOVED***
			v = append(v, cs[cursor:i]...)
			break
		***REMOVED***
		if c != '\\' ***REMOVED***
			i++
			continue
		***REMOVED***
		v = append(v, cs[cursor:i]...)
		i++
		c = cs[i]
		switch c ***REMOVED***
		case '"', '\\', '/', '\'':
			v = append(v, c)
		case 'b':
			v = append(v, '\b')
		case 'f':
			v = append(v, '\f')
		case 'n':
			v = append(v, '\n')
		case 'r':
			v = append(v, '\r')
		case 't':
			v = append(v, '\t')
		case 'u':
			var r rune
			var rr uint32
			if len(cs) < i+4 ***REMOVED*** // may help reduce bounds-checking
				d.d.errorf("need at least 4 more bytes for unicode sequence")
			***REMOVED***
			// c = cs[i+4] // may help reduce bounds-checking
			for j := 1; j < 5; j++ ***REMOVED***
				// best to use explicit if-else
				// - not a table, etc which involve memory loads, array lookup with bounds checks, etc
				c = cs[i+j]
				if c >= '0' && c <= '9' ***REMOVED***
					rr = rr*16 + uint32(c-jsonU4Chk2)
				***REMOVED*** else if c >= 'a' && c <= 'f' ***REMOVED***
					rr = rr*16 + uint32(c-jsonU4Chk1)
				***REMOVED*** else if c >= 'A' && c <= 'F' ***REMOVED***
					rr = rr*16 + uint32(c-jsonU4Chk0)
				***REMOVED*** else ***REMOVED***
					r = unicode.ReplacementChar
					i += 4
					goto encode_rune
				***REMOVED***
			***REMOVED***
			r = rune(rr)
			i += 4
			if utf16.IsSurrogate(r) ***REMOVED***
				if len(cs) >= i+6 && cs[i+2] == 'u' && cs[i+1] == '\\' ***REMOVED***
					i += 2
					// c = cs[i+4] // may help reduce bounds-checking
					var rr1 uint32
					for j := 1; j < 5; j++ ***REMOVED***
						c = cs[i+j]
						if c >= '0' && c <= '9' ***REMOVED***
							rr = rr*16 + uint32(c-jsonU4Chk2)
						***REMOVED*** else if c >= 'a' && c <= 'f' ***REMOVED***
							rr = rr*16 + uint32(c-jsonU4Chk1)
						***REMOVED*** else if c >= 'A' && c <= 'F' ***REMOVED***
							rr = rr*16 + uint32(c-jsonU4Chk0)
						***REMOVED*** else ***REMOVED***
							r = unicode.ReplacementChar
							i += 4
							goto encode_rune
						***REMOVED***
					***REMOVED***
					r = utf16.DecodeRune(r, rune(rr1))
					i += 4
				***REMOVED*** else ***REMOVED***
					r = unicode.ReplacementChar
					goto encode_rune
				***REMOVED***
			***REMOVED***
		encode_rune:
			w2 := utf8.EncodeRune(d.bstr[:], r)
			v = append(v, d.bstr[:w2]...)
		default:
			d.d.errorf("unsupported escaped value: %c", c)
		***REMOVED***
		i++
		cursor = i
	***REMOVED***
	d.bs = v
***REMOVED***

func (d *jsonDecDriver) nakedNum(z *decNaked, bs []byte) (err error) ***REMOVED***
	const cutoff = uint64(1 << uint(64-1))
	var n uint64
	var neg, badsyntax, overflow bool

	if d.h.PreferFloat ***REMOVED***
		goto F
	***REMOVED***
	n, neg, badsyntax, overflow = jsonParseInteger(bs)
	if badsyntax || overflow ***REMOVED***
		goto F
	***REMOVED***
	if neg ***REMOVED***
		if n > cutoff ***REMOVED***
			goto F
		***REMOVED***
		z.v = valueTypeInt
		z.i = -(int64(n))
	***REMOVED*** else if d.h.SignedInteger ***REMOVED***
		if n >= cutoff ***REMOVED***
			goto F
		***REMOVED***
		z.v = valueTypeInt
		z.i = int64(n)
	***REMOVED*** else ***REMOVED***
		z.v = valueTypeUint
		z.u = n
	***REMOVED***
	return
F:
	z.v = valueTypeFloat
	z.f, err = strconv.ParseFloat(stringView(bs), 64)
	return
***REMOVED***

func (d *jsonDecDriver) bsToString() string ***REMOVED***
	// if x := d.s.sc; x != nil && x.so && x.st == '***REMOVED***' ***REMOVED*** // map key
	if jsonAlwaysReturnInternString || d.c == containerMapKey ***REMOVED***
		return d.d.string(d.bs)
	***REMOVED***
	return string(d.bs)
***REMOVED***

func (d *jsonDecDriver) DecodeNaked() ***REMOVED***
	z := d.d.n
	// var decodeFurther bool

	if d.tok == 0 ***REMOVED***
		d.tok = d.r.skip(&jsonCharWhitespaceSet)
	***REMOVED***
	switch d.tok ***REMOVED***
	case 'n':
		d.readLit(3, jsonLitNull+1) // (n)ull
		z.v = valueTypeNil
	case 'f':
		d.readLit(4, jsonLitFalse+1) // (f)alse
		z.v = valueTypeBool
		z.b = false
	case 't':
		d.readLit(3, jsonLitTrue+1) // (t)rue
		z.v = valueTypeBool
		z.b = true
	case '***REMOVED***':
		z.v = valueTypeMap // don't consume. kInterfaceNaked will call ReadMapStart
	case '[':
		z.v = valueTypeArray // don't consume. kInterfaceNaked will call ReadArrayStart
	case '"':
		// if a string, and MapKeyAsString, then try to decode it as a nil, bool or number first
		d.appendStringAsBytes()
		if len(d.bs) > 0 && d.c == containerMapKey && d.h.MapKeyAsString ***REMOVED***
			switch stringView(d.bs) ***REMOVED***
			case "null":
				z.v = valueTypeNil
			case "true":
				z.v = valueTypeBool
				z.b = true
			case "false":
				z.v = valueTypeBool
				z.b = false
			default:
				// check if a number: float, int or uint
				if err := d.nakedNum(z, d.bs); err != nil ***REMOVED***
					z.v = valueTypeString
					z.s = d.bsToString()
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			z.v = valueTypeString
			z.s = d.bsToString()
		***REMOVED***
	default: // number
		bs := d.decNumBytes()
		if len(bs) == 0 ***REMOVED***
			d.d.errorf("decode number from empty string")
			return
		***REMOVED***
		if err := d.nakedNum(z, bs); err != nil ***REMOVED***
			d.d.errorf("decode number from %s: %v", bs, err)
			return
		***REMOVED***
	***REMOVED***
	// if decodeFurther ***REMOVED***
	// 	d.s.sc.retryRead()
	// ***REMOVED***
	return
***REMOVED***

//----------------------

// JsonHandle is a handle for JSON encoding format.
//
// Json is comprehensively supported:
//    - decodes numbers into interface***REMOVED******REMOVED*** as int, uint or float64
//      based on how the number looks and some config parameters e.g. PreferFloat, SignedInt, etc.
//    - decode integers from float formatted numbers e.g. 1.27e+8
//    - decode any json value (numbers, bool, etc) from quoted strings
//    - configurable way to encode/decode []byte .
//      by default, encodes and decodes []byte using base64 Std Encoding
//    - UTF-8 support for encoding and decoding
//
// It has better performance than the json library in the standard library,
// by leveraging the performance improvements of the codec library.
//
// In addition, it doesn't read more bytes than necessary during a decode, which allows
// reading multiple values from a stream containing json and non-json content.
// For example, a user can read a json value, then a cbor value, then a msgpack value,
// all from the same stream in sequence.
//
// Note that, when decoding quoted strings, invalid UTF-8 or invalid UTF-16 surrogate pairs are
// not treated as an error. Instead, they are replaced by the Unicode replacement character U+FFFD.
type JsonHandle struct ***REMOVED***
	textEncodingType
	BasicHandle

	// Indent indicates how a value is encoded.
	//   - If positive, indent by that number of spaces.
	//   - If negative, indent by that number of tabs.
	Indent int8

	// IntegerAsString controls how integers (signed and unsigned) are encoded.
	//
	// Per the JSON Spec, JSON numbers are 64-bit floating point numbers.
	// Consequently, integers > 2^53 cannot be represented as a JSON number without losing precision.
	// This can be mitigated by configuring how to encode integers.
	//
	// IntegerAsString interpretes the following values:
	//   - if 'L', then encode integers > 2^53 as a json string.
	//   - if 'A', then encode all integers as a json string
	//             containing the exact integer representation as a decimal.
	//   - else    encode all integers as a json number (default)
	IntegerAsString byte

	// HTMLCharsAsIs controls how to encode some special characters to html: < > &
	//
	// By default, we encode them as \uXXX
	// to prevent security holes when served from some browsers.
	HTMLCharsAsIs bool

	// PreferFloat says that we will default to decoding a number as a float.
	// If not set, we will examine the characters of the number and decode as an
	// integer type if it doesn't have any of the characters [.eE].
	PreferFloat bool

	// TermWhitespace says that we add a whitespace character
	// at the end of an encoding.
	//
	// The whitespace is important, especially if using numbers in a context
	// where multiple items are written to a stream.
	TermWhitespace bool

	// MapKeyAsString says to encode all map keys as strings.
	//
	// Use this to enforce strict json output.
	// The only caveat is that nil value is ALWAYS written as null (never as "null")
	MapKeyAsString bool

	// _ [2]byte // padding

	// Note: below, we store hardly-used items e.g. RawBytesExt is cached in the (en|de)cDriver.

	// RawBytesExt, if configured, is used to encode and decode raw bytes in a custom way.
	// If not configured, raw bytes are encoded to/from base64 text.
	RawBytesExt InterfaceExt

	_ [3]uint64 // padding
***REMOVED***

// Name returns the name of the handle: json
func (h *JsonHandle) Name() string            ***REMOVED*** return "json" ***REMOVED***
func (h *JsonHandle) hasElemSeparators() bool ***REMOVED*** return true ***REMOVED***
func (h *JsonHandle) typical() bool ***REMOVED***
	return h.Indent == 0 && !h.MapKeyAsString && h.IntegerAsString != 'A' && h.IntegerAsString != 'L'
***REMOVED***

type jsonTypical interface ***REMOVED***
	typical()
***REMOVED***

func (h *JsonHandle) recreateEncDriver(ed encDriver) (v bool) ***REMOVED***
	_, v = ed.(jsonTypical)
	return v != h.typical()
***REMOVED***

// SetInterfaceExt sets an extension
func (h *JsonHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, &extWrapper***REMOVED***bytesExtFailer***REMOVED******REMOVED***, ext***REMOVED***)
***REMOVED***

type jsonEncDriverTypicalImpl struct ***REMOVED***
	jsonEncDriver
	jsonEncDriverTypical
	_ [1]uint64 // padding
***REMOVED***

func (x *jsonEncDriverTypicalImpl) reset() ***REMOVED***
	x.jsonEncDriver.reset()
	x.jsonEncDriverTypical.reset(&x.jsonEncDriver)
***REMOVED***

type jsonEncDriverGenericImpl struct ***REMOVED***
	jsonEncDriver
	jsonEncDriverGeneric
***REMOVED***

func (x *jsonEncDriverGenericImpl) reset() ***REMOVED***
	x.jsonEncDriver.reset()
	x.jsonEncDriverGeneric.reset(&x.jsonEncDriver)
***REMOVED***

func (h *JsonHandle) newEncDriver(e *Encoder) (ee encDriver) ***REMOVED***
	var hd *jsonEncDriver
	if h.typical() ***REMOVED***
		var v jsonEncDriverTypicalImpl
		ee = &v
		hd = &v.jsonEncDriver
	***REMOVED*** else ***REMOVED***
		var v jsonEncDriverGenericImpl
		ee = &v
		hd = &v.jsonEncDriver
	***REMOVED***
	hd.e, hd.h, hd.bs = e, h, hd.b[:0]
	hd.se.BytesExt = bytesExtFailer***REMOVED******REMOVED***
	ee.reset()
	return
***REMOVED***

func (h *JsonHandle) newDecDriver(d *Decoder) decDriver ***REMOVED***
	// d := jsonDecDriver***REMOVED***r: r.(*bytesDecReader), h: h***REMOVED***
	hd := jsonDecDriver***REMOVED***d: d, h: h***REMOVED***
	hd.se.BytesExt = bytesExtFailer***REMOVED******REMOVED***
	hd.bs = hd.b[:0]
	hd.reset()
	return &hd
***REMOVED***

func (e *jsonEncDriver) reset() ***REMOVED***
	e.ew = e.e.w // e.e.w // &e.e.encWriterSwitch
	e.se.InterfaceExt = e.h.RawBytesExt
	if e.bs != nil ***REMOVED***
		e.bs = e.bs[:0]
	***REMOVED***
***REMOVED***

func (d *jsonDecDriver) reset() ***REMOVED***
	d.r = d.d.r // &d.d.decReaderSwitch // d.d.r
	d.se.InterfaceExt = d.h.RawBytesExt
	if d.bs != nil ***REMOVED***
		d.bs = d.bs[:0]
	***REMOVED***
	d.c, d.tok = 0, 0
	// d.n.reset()
***REMOVED***

func jsonFloatStrconvFmtPrec(f float64) (fmt byte, prec int) ***REMOVED***
	prec = -1
	var abs = math.Abs(f)
	if abs != 0 && (abs < 1e-6 || abs >= 1e21) ***REMOVED***
		fmt = 'e'
	***REMOVED*** else ***REMOVED***
		fmt = 'f'
		// set prec to 1 iff mod is 0.
		//     better than using jsonIsFloatBytesB2 to check if a . or E in the float bytes.
		// this ensures that every float has an e or .0 in it.
		if abs <= 1 ***REMOVED***
			if abs == 0 || abs == 1 ***REMOVED***
				prec = 1
			***REMOVED***
		***REMOVED*** else if _, mod := math.Modf(abs); mod == 0 ***REMOVED***
			prec = 1
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// custom-fitted version of strconv.Parse(Ui|I)nt.
// Also ensures we don't have to search for .eE to determine if a float or not.
func jsonParseInteger(s []byte) (n uint64, neg, badSyntax, overflow bool) ***REMOVED***
	const maxUint64 = (1<<64 - 1)
	const cutoff = maxUint64/10 + 1

	if len(s) == 0 ***REMOVED***
		badSyntax = true
		return
	***REMOVED***
	switch s[0] ***REMOVED***
	case '+':
		s = s[1:]
	case '-':
		s = s[1:]
		neg = true
	***REMOVED***
	for _, c := range s ***REMOVED***
		if c < '0' || c > '9' ***REMOVED***
			badSyntax = true
			return
		***REMOVED***
		// unsigned integers don't overflow well on multiplication, so check cutoff here
		// e.g. (maxUint64-5)*10 doesn't overflow well ...
		if n >= cutoff ***REMOVED***
			overflow = true
			return
		***REMOVED***
		n *= 10
		n1 := n + uint64(c-'0')
		if n1 < n || n1 > maxUint64 ***REMOVED***
			overflow = true
			return
		***REMOVED***
		n = n1
	***REMOVED***
	return
***REMOVED***

var _ decDriver = (*jsonDecDriver)(nil)
var _ encDriver = (*jsonEncDriverGenericImpl)(nil)
var _ encDriver = (*jsonEncDriverTypicalImpl)(nil)
var _ jsonTypical = (*jsonEncDriverTypical)(nil)
