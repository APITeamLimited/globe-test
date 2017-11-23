// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
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
//   - Also, strconv.ParseXXX for floats and integers
//     - only works on strings resulting in unnecessary allocation and []byte-string conversion.
//     - it does a lot of redundant checks, because json numbers are simpler that what it supports.
//   - We parse numbers (floats and integers) directly here.
//     We only delegate parsing floats if it is a hairy float which could cause a loss of precision.
//     In that case, we delegate to strconv.ParseFloat.
//
// Note:
//   - encode does not beautify. There is no whitespace when encoding.
//   - rpc calls which take single integer arguments or write single numeric arguments will need care.

// Top-level methods of json(End|Dec)Driver (which are implementations of (en|de)cDriver
// MUST not call one-another.

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

//--------------------------------

var (
	jsonLiterals = [...]byte***REMOVED***'t', 'r', 'u', 'e', 'f', 'a', 'l', 's', 'e', 'n', 'u', 'l', 'l'***REMOVED***

	jsonFloat64Pow10 = [...]float64***REMOVED***
		1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
		1e20, 1e21, 1e22,
	***REMOVED***

	jsonUint64Pow10 = [...]uint64***REMOVED***
		1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
		1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	***REMOVED***

	// jsonTabs and jsonSpaces are used as caches for indents
	jsonTabs, jsonSpaces string
)

const (
	// jsonUnreadAfterDecNum controls whether we unread after decoding a number.
	//
	// instead of unreading, just update d.tok (iff it's not a whitespace char)
	// However, doing this means that we may HOLD onto some data which belongs to another stream.
	// Thus, it is safest to unread the data when done.
	// keep behind a constant flag for now.
	jsonUnreadAfterDecNum = true

	// If !jsonValidateSymbols, decoding will be faster, by skipping some checks:
	//   - If we see first character of null, false or true,
	//     do not validate subsequent characters.
	//   - e.g. if we see a n, assume null and skip next 3 characters,
	//     and do not validate they are ull.
	// P.S. Do not expect a significant decoding boost from this.
	jsonValidateSymbols = true

	// if jsonTruncateMantissa, truncate mantissa if trailing 0's.
	// This is important because it could allow some floats to be decoded without
	// deferring to strconv.ParseFloat.
	jsonTruncateMantissa = true

	// if mantissa >= jsonNumUintCutoff before multiplying by 10, this is an overflow
	jsonNumUintCutoff = (1<<64-1)/uint64(10) + 1 // cutoff64(base)

	// if mantissa >= jsonNumUintMaxVal, this is an overflow
	jsonNumUintMaxVal = 1<<uint64(64) - 1

	// jsonNumDigitsUint64Largest = 19

	jsonSpacesOrTabsLen = 128
)

func init() ***REMOVED***
	var bs [jsonSpacesOrTabsLen]byte
	for i := 0; i < jsonSpacesOrTabsLen; i++ ***REMOVED***
		bs[i] = ' '
	***REMOVED***
	jsonSpaces = string(bs[:])

	for i := 0; i < jsonSpacesOrTabsLen; i++ ***REMOVED***
		bs[i] = '\t'
	***REMOVED***
	jsonTabs = string(bs[:])
***REMOVED***

type jsonEncDriver struct ***REMOVED***
	e  *Encoder
	w  encWriter
	h  *JsonHandle
	b  [64]byte // scratch
	bs []byte   // scratch
	se setExtWrapper
	ds string // indent string
	dl uint16 // indent level
	dt bool   // indent using tabs
	d  bool   // indent
	c  containerState
	noBuiltInTypes
***REMOVED***

// indent is done as below:
//   - newline and indent are added before each mapKey or arrayElem
//   - newline and indent are added before each ending,
//     except there was no entry (so we can have ***REMOVED******REMOVED*** or [])

func (e *jsonEncDriver) sendContainerState(c containerState) ***REMOVED***
	// determine whether to output separators
	if c == containerMapKey ***REMOVED***
		if e.c != containerMapStart ***REMOVED***
			e.w.writen1(',')
		***REMOVED***
		if e.d ***REMOVED***
			e.writeIndent()
		***REMOVED***
	***REMOVED*** else if c == containerMapValue ***REMOVED***
		if e.d ***REMOVED***
			e.w.writen2(':', ' ')
		***REMOVED*** else ***REMOVED***
			e.w.writen1(':')
		***REMOVED***
	***REMOVED*** else if c == containerMapEnd ***REMOVED***
		if e.d ***REMOVED***
			e.dl--
			if e.c != containerMapStart ***REMOVED***
				e.writeIndent()
			***REMOVED***
		***REMOVED***
		e.w.writen1('***REMOVED***')
	***REMOVED*** else if c == containerArrayElem ***REMOVED***
		if e.c != containerArrayStart ***REMOVED***
			e.w.writen1(',')
		***REMOVED***
		if e.d ***REMOVED***
			e.writeIndent()
		***REMOVED***
	***REMOVED*** else if c == containerArrayEnd ***REMOVED***
		if e.d ***REMOVED***
			e.dl--
			if e.c != containerArrayStart ***REMOVED***
				e.writeIndent()
			***REMOVED***
		***REMOVED***
		e.w.writen1(']')
	***REMOVED***
	e.c = c
***REMOVED***

func (e *jsonEncDriver) writeIndent() ***REMOVED***
	e.w.writen1('\n')
	if x := len(e.ds) * int(e.dl); x <= jsonSpacesOrTabsLen ***REMOVED***
		if e.dt ***REMOVED***
			e.w.writestr(jsonTabs[:x])
		***REMOVED*** else ***REMOVED***
			e.w.writestr(jsonSpaces[:x])
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for i := uint16(0); i < e.dl; i++ ***REMOVED***
			e.w.writestr(e.ds)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeNil() ***REMOVED***
	e.w.writeb(jsonLiterals[9:13]) // null
***REMOVED***

func (e *jsonEncDriver) EncodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.w.writeb(jsonLiterals[0:4]) // true
	***REMOVED*** else ***REMOVED***
		e.w.writeb(jsonLiterals[4:9]) // false
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeFloat32(f float32) ***REMOVED***
	e.encodeFloat(float64(f), 32)
***REMOVED***

func (e *jsonEncDriver) EncodeFloat64(f float64) ***REMOVED***
	// e.w.writestr(strconv.FormatFloat(f, 'E', -1, 64))
	e.encodeFloat(f, 64)
***REMOVED***

func (e *jsonEncDriver) encodeFloat(f float64, numbits int) ***REMOVED***
	x := strconv.AppendFloat(e.b[:0], f, 'G', -1, numbits)
	e.w.writeb(x)
	if bytes.IndexByte(x, 'E') == -1 && bytes.IndexByte(x, '.') == -1 ***REMOVED***
		e.w.writen2('.', '0')
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeInt(v int64) ***REMOVED***
	if x := e.h.IntegerAsString; x == 'A' || x == 'L' && (v > 1<<53 || v < -(1<<53)) ***REMOVED***
		e.w.writen1('"')
		e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
		e.w.writen1('"')
		return
	***REMOVED***
	e.w.writeb(strconv.AppendInt(e.b[:0], v, 10))
***REMOVED***

func (e *jsonEncDriver) EncodeUint(v uint64) ***REMOVED***
	if x := e.h.IntegerAsString; x == 'A' || x == 'L' && v > 1<<53 ***REMOVED***
		e.w.writen1('"')
		e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
		e.w.writen1('"')
		return
	***REMOVED***
	e.w.writeb(strconv.AppendUint(e.b[:0], v, 10))
***REMOVED***

func (e *jsonEncDriver) EncodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext, en *Encoder) ***REMOVED***
	if v := ext.ConvertExt(rv); v == nil ***REMOVED***
		e.w.writeb(jsonLiterals[9:13]) // null // e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		en.encode(v)
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeRawExt(re *RawExt, en *Encoder) ***REMOVED***
	// only encodes re.Value (never re.Data)
	if re.Value == nil ***REMOVED***
		e.w.writeb(jsonLiterals[9:13]) // null // e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		en.encode(re.Value)
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeArrayStart(length int) ***REMOVED***
	if e.d ***REMOVED***
		e.dl++
	***REMOVED***
	e.w.writen1('[')
	e.c = containerArrayStart
***REMOVED***

func (e *jsonEncDriver) EncodeMapStart(length int) ***REMOVED***
	if e.d ***REMOVED***
		e.dl++
	***REMOVED***
	e.w.writen1('***REMOVED***')
	e.c = containerMapStart
***REMOVED***

func (e *jsonEncDriver) EncodeString(c charEncoding, v string) ***REMOVED***
	// e.w.writestr(strconv.Quote(v))
	e.quoteStr(v)
***REMOVED***

func (e *jsonEncDriver) EncodeSymbol(v string) ***REMOVED***
	// e.EncodeString(c_UTF8, v)
	e.quoteStr(v)
***REMOVED***

func (e *jsonEncDriver) EncodeStringBytes(c charEncoding, v []byte) ***REMOVED***
	// if encoding raw bytes and RawBytesExt is configured, use it to encode
	if c == c_RAW && e.se.i != nil ***REMOVED***
		e.EncodeExt(v, 0, &e.se, e.e)
		return
	***REMOVED***
	if c == c_RAW ***REMOVED***
		slen := base64.StdEncoding.EncodedLen(len(v))
		if cap(e.bs) >= slen ***REMOVED***
			e.bs = e.bs[:slen]
		***REMOVED*** else ***REMOVED***
			e.bs = make([]byte, slen)
		***REMOVED***
		base64.StdEncoding.Encode(e.bs, v)
		e.w.writen1('"')
		e.w.writeb(e.bs)
		e.w.writen1('"')
	***REMOVED*** else ***REMOVED***
		// e.EncodeString(c, string(v))
		e.quoteStr(stringView(v))
	***REMOVED***
***REMOVED***

func (e *jsonEncDriver) EncodeAsis(v []byte) ***REMOVED***
	e.w.writeb(v)
***REMOVED***

func (e *jsonEncDriver) quoteStr(s string) ***REMOVED***
	// adapted from std pkg encoding/json
	const hex = "0123456789abcdef"
	w := e.w
	w.writen1('"')
	start := 0
	for i := 0; i < len(s); ***REMOVED***
		// encode all bytes < 0x20 (except \r, \n).
		// also encode < > & to prevent security holes when served to some browsers.
		if b := s[i]; b < utf8.RuneSelf ***REMOVED***
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' ***REMOVED***
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
			case '<', '>', '&':
				if e.h.HTMLCharsAsIs ***REMOVED***
					w.writen1(b)
				***REMOVED*** else ***REMOVED***
					w.writestr(`\u00`)
					w.writen2(hex[b>>4], hex[b&0xF])
				***REMOVED***
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

//--------------------------------

type jsonNum struct ***REMOVED***
	// bytes            []byte // may have [+-.eE0-9]
	mantissa         uint64 // where mantissa ends, and maybe dot begins.
	exponent         int16  // exponent value.
	manOverflow      bool
	neg              bool // started with -. No initial sign in the bytes above.
	dot              bool // has dot
	explicitExponent bool // explicit exponent
***REMOVED***

func (x *jsonNum) reset() ***REMOVED***
	x.manOverflow = false
	x.neg = false
	x.dot = false
	x.explicitExponent = false
	x.mantissa = 0
	x.exponent = 0
***REMOVED***

// uintExp is called only if exponent > 0.
func (x *jsonNum) uintExp() (n uint64, overflow bool) ***REMOVED***
	n = x.mantissa
	e := x.exponent
	if e >= int16(len(jsonUint64Pow10)) ***REMOVED***
		overflow = true
		return
	***REMOVED***
	n *= jsonUint64Pow10[e]
	if n < x.mantissa || n > jsonNumUintMaxVal ***REMOVED***
		overflow = true
		return
	***REMOVED***
	return
	// for i := int16(0); i < e; i++ ***REMOVED***
	// 	if n >= jsonNumUintCutoff ***REMOVED***
	// 		overflow = true
	// 		return
	// 	***REMOVED***
	// 	n *= 10
	// ***REMOVED***
	// return
***REMOVED***

// these constants are only used withn floatVal.
// They are brought out, so that floatVal can be inlined.
const (
	jsonUint64MantissaBits = 52
	jsonMaxExponent        = int16(len(jsonFloat64Pow10)) - 1
)

func (x *jsonNum) floatVal() (f float64, parseUsingStrConv bool) ***REMOVED***
	// We do not want to lose precision.
	// Consequently, we will delegate to strconv.ParseFloat if any of the following happen:
	//    - There are more digits than in math.MaxUint64: 18446744073709551615 (20 digits)
	//      We expect up to 99.... (19 digits)
	//    - The mantissa cannot fit into a 52 bits of uint64
	//    - The exponent is beyond our scope ie beyong 22.
	parseUsingStrConv = x.manOverflow ||
		x.exponent > jsonMaxExponent ||
		(x.exponent < 0 && -(x.exponent) > jsonMaxExponent) ||
		x.mantissa>>jsonUint64MantissaBits != 0

	if parseUsingStrConv ***REMOVED***
		return
	***REMOVED***

	// all good. so handle parse here.
	f = float64(x.mantissa)
	// fmt.Printf(".Float: uint64 value: %v, float: %v\n", m, f)
	if x.neg ***REMOVED***
		f = -f
	***REMOVED***
	if x.exponent > 0 ***REMOVED***
		f *= jsonFloat64Pow10[x.exponent]
	***REMOVED*** else if x.exponent < 0 ***REMOVED***
		f /= jsonFloat64Pow10[-x.exponent]
	***REMOVED***
	return
***REMOVED***

type jsonDecDriver struct ***REMOVED***
	noBuiltInTypes
	d *Decoder
	h *JsonHandle
	r decReader

	c containerState
	// tok is used to store the token read right after skipWhiteSpace.
	tok uint8

	bstr [8]byte  // scratch used for string \UXXX parsing
	b    [64]byte // scratch, used for parsing strings or numbers
	b2   [64]byte // scratch, used only for decodeBytes (after base64)
	bs   []byte   // scratch. Initialized from b. Used for parsing strings or numbers.

	se setExtWrapper

	n jsonNum
***REMOVED***

func jsonIsWS(b byte) bool ***REMOVED***
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
***REMOVED***

// // This will skip whitespace characters and return the next byte to read.
// // The next byte determines what the value will be one of.
// func (d *jsonDecDriver) skipWhitespace() ***REMOVED***
// 	// fast-path: do not enter loop. Just check first (in case no whitespace).
// 	b := d.r.readn1()
// 	if jsonIsWS(b) ***REMOVED***
// 		r := d.r
// 		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
// 		***REMOVED***
// 	***REMOVED***
// 	d.tok = b
// ***REMOVED***

func (d *jsonDecDriver) uncacheRead() ***REMOVED***
	if d.tok != 0 ***REMOVED***
		d.r.unreadn1()
		d.tok = 0
	***REMOVED***
***REMOVED***

func (d *jsonDecDriver) sendContainerState(c containerState) ***REMOVED***
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
	var xc uint8 // char expected
	if c == containerMapKey ***REMOVED***
		if d.c != containerMapStart ***REMOVED***
			xc = ','
		***REMOVED***
	***REMOVED*** else if c == containerMapValue ***REMOVED***
		xc = ':'
	***REMOVED*** else if c == containerMapEnd ***REMOVED***
		xc = '***REMOVED***'
	***REMOVED*** else if c == containerArrayElem ***REMOVED***
		if d.c != containerArrayStart ***REMOVED***
			xc = ','
		***REMOVED***
	***REMOVED*** else if c == containerArrayEnd ***REMOVED***
		xc = ']'
	***REMOVED***
	if xc != 0 ***REMOVED***
		if d.tok != xc ***REMOVED***
			d.d.errorf("json: expect char '%c' but got char '%c'", xc, d.tok)
		***REMOVED***
		d.tok = 0
	***REMOVED***
	d.c = c
***REMOVED***

func (d *jsonDecDriver) CheckBreak() bool ***REMOVED***
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
	if d.tok == '***REMOVED***' || d.tok == ']' ***REMOVED***
		// d.tok = 0 // only checking, not consuming
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *jsonDecDriver) readStrIdx(fromIdx, toIdx uint8) ***REMOVED***
	bs := d.r.readx(int(toIdx - fromIdx))
	d.tok = 0
	if jsonValidateSymbols ***REMOVED***
		if !bytes.Equal(bs, jsonLiterals[fromIdx:toIdx]) ***REMOVED***
			d.d.errorf("json: expecting %s: got %s", jsonLiterals[fromIdx:toIdx], bs)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *jsonDecDriver) TryDecodeAsNil() bool ***REMOVED***
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
	if d.tok == 'n' ***REMOVED***
		d.readStrIdx(10, 13) // ull
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *jsonDecDriver) DecodeBool() bool ***REMOVED***
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
	if d.tok == 'f' ***REMOVED***
		d.readStrIdx(5, 9) // alse
		return false
	***REMOVED***
	if d.tok == 't' ***REMOVED***
		d.readStrIdx(1, 4) // rue
		return true
	***REMOVED***
	d.d.errorf("json: decode bool: got first char %c", d.tok)
	return false // "unreachable"
***REMOVED***

func (d *jsonDecDriver) ReadMapStart() int ***REMOVED***
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
	if d.tok != '***REMOVED***' ***REMOVED***
		d.d.errorf("json: expect char '%c' but got char '%c'", '***REMOVED***', d.tok)
	***REMOVED***
	d.tok = 0
	d.c = containerMapStart
	return -1
***REMOVED***

func (d *jsonDecDriver) ReadArrayStart() int ***REMOVED***
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
	if d.tok != '[' ***REMOVED***
		d.d.errorf("json: expect char '%c' but got char '%c'", '[', d.tok)
	***REMOVED***
	d.tok = 0
	d.c = containerArrayStart
	return -1
***REMOVED***

func (d *jsonDecDriver) ContainerType() (vt valueType) ***REMOVED***
	// check container type by checking the first char
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
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
	// d.d.errorf("isContainerType: unsupported parameter: %v", vt)
	// return false // "unreachable"
***REMOVED***

func (d *jsonDecDriver) decNum(storeBytes bool) ***REMOVED***
	// If it is has a . or an e|E, decode as a float; else decode as an int.
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
	b := d.tok
	var str bool
	if b == '"' ***REMOVED***
		str = true
		b = d.r.readn1()
	***REMOVED***
	if !(b == '+' || b == '-' || b == '.' || (b >= '0' && b <= '9')) ***REMOVED***
		d.d.errorf("json: decNum: got first char '%c'", b)
		return
	***REMOVED***
	d.tok = 0

	const cutoff = (1<<64-1)/uint64(10) + 1 // cutoff64(base)
	const jsonNumUintMaxVal = 1<<uint64(64) - 1

	n := &d.n
	r := d.r
	n.reset()
	d.bs = d.bs[:0]

	if str && storeBytes ***REMOVED***
		d.bs = append(d.bs, '"')
	***REMOVED***

	// The format of a number is as below:
	// parsing:     sign? digit* dot? digit* e?  sign? digit*
	// states:  0   1*    2      3*   4      5*  6     7
	// We honor this state so we can break correctly.
	var state uint8 = 0
	var eNeg bool
	var e int16
	var eof bool
LOOP:
	for !eof ***REMOVED***
		// fmt.Printf("LOOP: b: %q\n", b)
		switch b ***REMOVED***
		case '+':
			switch state ***REMOVED***
			case 0:
				state = 2
				// do not add sign to the slice ...
				b, eof = r.readn1eof()
				continue
			case 6: // typ = jsonNumFloat
				state = 7
			default:
				break LOOP
			***REMOVED***
		case '-':
			switch state ***REMOVED***
			case 0:
				state = 2
				n.neg = true
				// do not add sign to the slice ...
				b, eof = r.readn1eof()
				continue
			case 6: // typ = jsonNumFloat
				eNeg = true
				state = 7
			default:
				break LOOP
			***REMOVED***
		case '.':
			switch state ***REMOVED***
			case 0, 2: // typ = jsonNumFloat
				state = 4
				n.dot = true
			default:
				break LOOP
			***REMOVED***
		case 'e', 'E':
			switch state ***REMOVED***
			case 0, 2, 4: // typ = jsonNumFloat
				state = 6
				// n.mantissaEndIndex = int16(len(n.bytes))
				n.explicitExponent = true
			default:
				break LOOP
			***REMOVED***
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			switch state ***REMOVED***
			case 0:
				state = 2
				fallthrough
			case 2:
				fallthrough
			case 4:
				if n.dot ***REMOVED***
					n.exponent--
				***REMOVED***
				if n.mantissa >= jsonNumUintCutoff ***REMOVED***
					n.manOverflow = true
					break
				***REMOVED***
				v := uint64(b - '0')
				n.mantissa *= 10
				if v != 0 ***REMOVED***
					n1 := n.mantissa + v
					if n1 < n.mantissa || n1 > jsonNumUintMaxVal ***REMOVED***
						n.manOverflow = true // n+v overflows
						break
					***REMOVED***
					n.mantissa = n1
				***REMOVED***
			case 6:
				state = 7
				fallthrough
			case 7:
				if !(b == '0' && e == 0) ***REMOVED***
					e = e*10 + int16(b-'0')
				***REMOVED***
			default:
				break LOOP
			***REMOVED***
		case '"':
			if str ***REMOVED***
				if storeBytes ***REMOVED***
					d.bs = append(d.bs, '"')
				***REMOVED***
				b, eof = r.readn1eof()
			***REMOVED***
			break LOOP
		default:
			break LOOP
		***REMOVED***
		if storeBytes ***REMOVED***
			d.bs = append(d.bs, b)
		***REMOVED***
		b, eof = r.readn1eof()
	***REMOVED***

	if jsonTruncateMantissa && n.mantissa != 0 ***REMOVED***
		for n.mantissa%10 == 0 ***REMOVED***
			n.mantissa /= 10
			n.exponent++
		***REMOVED***
	***REMOVED***

	if e != 0 ***REMOVED***
		if eNeg ***REMOVED***
			n.exponent -= e
		***REMOVED*** else ***REMOVED***
			n.exponent += e
		***REMOVED***
	***REMOVED***

	// d.n = n

	if !eof ***REMOVED***
		if jsonUnreadAfterDecNum ***REMOVED***
			r.unreadn1()
		***REMOVED*** else ***REMOVED***
			if !jsonIsWS(b) ***REMOVED***
				d.tok = b
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// fmt.Printf("1: n: bytes: %s, neg: %v, dot: %v, exponent: %v, mantissaEndIndex: %v\n",
	// 	n.bytes, n.neg, n.dot, n.exponent, n.mantissaEndIndex)
	return
***REMOVED***

func (d *jsonDecDriver) DecodeInt(bitsize uint8) (i int64) ***REMOVED***
	d.decNum(false)
	n := &d.n
	if n.manOverflow ***REMOVED***
		d.d.errorf("json: overflow integer after: %v", n.mantissa)
		return
	***REMOVED***
	var u uint64
	if n.exponent == 0 ***REMOVED***
		u = n.mantissa
	***REMOVED*** else if n.exponent < 0 ***REMOVED***
		d.d.errorf("json: fractional integer")
		return
	***REMOVED*** else if n.exponent > 0 ***REMOVED***
		var overflow bool
		if u, overflow = n.uintExp(); overflow ***REMOVED***
			d.d.errorf("json: overflow integer")
			return
		***REMOVED***
	***REMOVED***
	i = int64(u)
	if n.neg ***REMOVED***
		i = -i
	***REMOVED***
	if chkOvf.Int(i, bitsize) ***REMOVED***
		d.d.errorf("json: overflow %v bits: %s", bitsize, d.bs)
		return
	***REMOVED***
	// fmt.Printf("DecodeInt: %v\n", i)
	return
***REMOVED***

// floatVal MUST only be called after a decNum, as d.bs now contains the bytes of the number
func (d *jsonDecDriver) floatVal() (f float64) ***REMOVED***
	f, useStrConv := d.n.floatVal()
	if useStrConv ***REMOVED***
		var err error
		if f, err = strconv.ParseFloat(stringView(d.bs), 64); err != nil ***REMOVED***
			panic(fmt.Errorf("parse float: %s, %v", d.bs, err))
		***REMOVED***
		if d.n.neg ***REMOVED***
			f = -f
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (d *jsonDecDriver) DecodeUint(bitsize uint8) (u uint64) ***REMOVED***
	d.decNum(false)
	n := &d.n
	if n.neg ***REMOVED***
		d.d.errorf("json: unsigned integer cannot be negative")
		return
	***REMOVED***
	if n.manOverflow ***REMOVED***
		d.d.errorf("json: overflow integer after: %v", n.mantissa)
		return
	***REMOVED***
	if n.exponent == 0 ***REMOVED***
		u = n.mantissa
	***REMOVED*** else if n.exponent < 0 ***REMOVED***
		d.d.errorf("json: fractional integer")
		return
	***REMOVED*** else if n.exponent > 0 ***REMOVED***
		var overflow bool
		if u, overflow = n.uintExp(); overflow ***REMOVED***
			d.d.errorf("json: overflow integer")
			return
		***REMOVED***
	***REMOVED***
	if chkOvf.Uint(u, bitsize) ***REMOVED***
		d.d.errorf("json: overflow %v bits: %s", bitsize, d.bs)
		return
	***REMOVED***
	// fmt.Printf("DecodeUint: %v\n", u)
	return
***REMOVED***

func (d *jsonDecDriver) DecodeFloat(chkOverflow32 bool) (f float64) ***REMOVED***
	d.decNum(true)
	f = d.floatVal()
	if chkOverflow32 && chkOvf.Float32(f) ***REMOVED***
		d.d.errorf("json: overflow float32: %v, %s", f, d.bs)
		return
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

func (d *jsonDecDriver) DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte) ***REMOVED***
	// if decoding into raw bytes, and the RawBytesExt is configured, use it to decode.
	if !isstring && d.se.i != nil ***REMOVED***
		bsOut = bs
		d.DecodeExt(&bsOut, 0, &d.se)
		return
	***REMOVED***
	d.appendStringAsBytes()
	// if isstring, then just return the bytes, even if it is using the scratch buffer.
	// the bytes will be converted to a string as needed.
	if isstring ***REMOVED***
		return d.bs
	***REMOVED***
	// if appendStringAsBytes returned a zero-len slice, then treat as nil.
	// This should only happen for null, and "".
	if len(d.bs) == 0 ***REMOVED***
		return nil
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
		d.d.errorf("json: error decoding base64 binary '%s': %v", bs0, err)
		return nil
	***REMOVED***
	if slen != slen2 ***REMOVED***
		bsOut = bsOut[:slen2]
	***REMOVED***
	return
***REMOVED***

func (d *jsonDecDriver) DecodeString() (s string) ***REMOVED***
	d.appendStringAsBytes()
	// if x := d.s.sc; x != nil && x.so && x.st == '***REMOVED***' ***REMOVED*** // map key
	if d.c == containerMapKey ***REMOVED***
		return d.d.string(d.bs)
	***REMOVED***
	return string(d.bs)
***REMOVED***

func (d *jsonDecDriver) appendStringAsBytes() ***REMOVED***
	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***

	// handle null as a string
	if d.tok == 'n' ***REMOVED***
		d.readStrIdx(10, 13) // ull
		d.bs = d.bs[:0]
		return
	***REMOVED***

	if d.tok != '"' ***REMOVED***
		d.d.errorf("json: expect char '%c' but got char '%c'", '"', d.tok)
	***REMOVED***
	d.tok = 0

	v := d.bs[:0]
	var c uint8
	r := d.r
	for ***REMOVED***
		c = r.readn1()
		if c == '"' ***REMOVED***
			break
		***REMOVED*** else if c == '\\' ***REMOVED***
			c = r.readn1()
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
				rr := d.jsonU4(false)
				// fmt.Printf("$$$$$$$$$: is surrogate: %v\n", utf16.IsSurrogate(rr))
				if utf16.IsSurrogate(rr) ***REMOVED***
					rr = utf16.DecodeRune(rr, d.jsonU4(true))
				***REMOVED***
				w2 := utf8.EncodeRune(d.bstr[:], rr)
				v = append(v, d.bstr[:w2]...)
			default:
				d.d.errorf("json: unsupported escaped value: %c", c)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			v = append(v, c)
		***REMOVED***
	***REMOVED***
	d.bs = v
***REMOVED***

func (d *jsonDecDriver) jsonU4(checkSlashU bool) rune ***REMOVED***
	r := d.r
	if checkSlashU && !(r.readn1() == '\\' && r.readn1() == 'u') ***REMOVED***
		d.d.errorf(`json: unquoteStr: invalid unicode sequence. Expecting \u`)
		return 0
	***REMOVED***
	// u, _ := strconv.ParseUint(string(d.bstr[:4]), 16, 64)
	var u uint32
	for i := 0; i < 4; i++ ***REMOVED***
		v := r.readn1()
		if '0' <= v && v <= '9' ***REMOVED***
			v = v - '0'
		***REMOVED*** else if 'a' <= v && v <= 'z' ***REMOVED***
			v = v - 'a' + 10
		***REMOVED*** else if 'A' <= v && v <= 'Z' ***REMOVED***
			v = v - 'A' + 10
		***REMOVED*** else ***REMOVED***
			d.d.errorf(`json: unquoteStr: invalid hex char in \u unicode sequence: %q`, v)
			return 0
		***REMOVED***
		u = u*16 + uint32(v)
	***REMOVED***
	return rune(u)
***REMOVED***

func (d *jsonDecDriver) DecodeNaked() ***REMOVED***
	z := &d.d.n
	// var decodeFurther bool

	if d.tok == 0 ***REMOVED***
		var b byte
		r := d.r
		for b = r.readn1(); jsonIsWS(b); b = r.readn1() ***REMOVED***
		***REMOVED***
		d.tok = b
	***REMOVED***
	switch d.tok ***REMOVED***
	case 'n':
		d.readStrIdx(10, 13) // ull
		z.v = valueTypeNil
	case 'f':
		d.readStrIdx(5, 9) // alse
		z.v = valueTypeBool
		z.b = false
	case 't':
		d.readStrIdx(1, 4) // rue
		z.v = valueTypeBool
		z.b = true
	case '***REMOVED***':
		z.v = valueTypeMap
		// d.tok = 0 // don't consume. kInterfaceNaked will call ReadMapStart
		// decodeFurther = true
	case '[':
		z.v = valueTypeArray
		// d.tok = 0 // don't consume. kInterfaceNaked will call ReadArrayStart
		// decodeFurther = true
	case '"':
		z.v = valueTypeString
		z.s = d.DecodeString()
	default: // number
		d.decNum(true)
		n := &d.n
		// if the string had a any of [.eE], then decode as float.
		switch ***REMOVED***
		case n.explicitExponent, n.dot, n.exponent < 0, n.manOverflow:
			z.v = valueTypeFloat
			z.f = d.floatVal()
		case n.exponent == 0:
			u := n.mantissa
			switch ***REMOVED***
			case n.neg:
				z.v = valueTypeInt
				z.i = -int64(u)
			case d.h.SignedInteger:
				z.v = valueTypeInt
				z.i = int64(u)
			default:
				z.v = valueTypeUint
				z.u = u
			***REMOVED***
		default:
			u, overflow := n.uintExp()
			switch ***REMOVED***
			case overflow:
				z.v = valueTypeFloat
				z.f = d.floatVal()
			case n.neg:
				z.v = valueTypeInt
				z.i = -int64(u)
			case d.h.SignedInteger:
				z.v = valueTypeInt
				z.i = int64(u)
			default:
				z.v = valueTypeUint
				z.u = u
			***REMOVED***
		***REMOVED***
		// fmt.Printf("DecodeNaked: Number: %T, %v\n", v, v)
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
//    - configurable way to encode/decode []byte .
//      by default, encodes and decodes []byte using base64 Std Encoding
//    - UTF-8 support for encoding and decoding
//
// It has better performance than the json library in the standard library,
// by leveraging the performance improvements of the codec library and
// minimizing allocations.
//
// In addition, it doesn't read more bytes than necessary during a decode, which allows
// reading multiple values from a stream containing json and non-json content.
// For example, a user can read a json value, then a cbor value, then a msgpack value,
// all from the same stream in sequence.
type JsonHandle struct ***REMOVED***
	textEncodingType
	BasicHandle
	// RawBytesExt, if configured, is used to encode and decode raw bytes in a custom way.
	// If not configured, raw bytes are encoded to/from base64 text.
	RawBytesExt InterfaceExt

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
	IntegerAsString uint8

	// HTMLCharsAsIs controls how to encode some special characters to html: < > &
	//
	// By default, we encode them as \uXXX
	// to prevent security holes when served from some browsers.
	HTMLCharsAsIs bool
***REMOVED***

func (h *JsonHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, &setExtWrapper***REMOVED***i: ext***REMOVED***)
***REMOVED***

func (h *JsonHandle) newEncDriver(e *Encoder) encDriver ***REMOVED***
	hd := jsonEncDriver***REMOVED***e: e, h: h***REMOVED***
	hd.bs = hd.b[:0]

	hd.reset()

	return &hd
***REMOVED***

func (h *JsonHandle) newDecDriver(d *Decoder) decDriver ***REMOVED***
	// d := jsonDecDriver***REMOVED***r: r.(*bytesDecReader), h: h***REMOVED***
	hd := jsonDecDriver***REMOVED***d: d, h: h***REMOVED***
	hd.bs = hd.b[:0]
	hd.reset()
	return &hd
***REMOVED***

func (e *jsonEncDriver) reset() ***REMOVED***
	e.w = e.e.w
	e.se.i = e.h.RawBytesExt
	if e.bs != nil ***REMOVED***
		e.bs = e.bs[:0]
	***REMOVED***
	e.d, e.dt, e.dl, e.ds = false, false, 0, ""
	e.c = 0
	if e.h.Indent > 0 ***REMOVED***
		e.d = true
		e.ds = jsonSpaces[:e.h.Indent]
	***REMOVED*** else if e.h.Indent < 0 ***REMOVED***
		e.d = true
		e.dt = true
		e.ds = jsonTabs[:-(e.h.Indent)]
	***REMOVED***
***REMOVED***

func (d *jsonDecDriver) reset() ***REMOVED***
	d.r = d.d.r
	d.se.i = d.h.RawBytesExt
	if d.bs != nil ***REMOVED***
		d.bs = d.bs[:0]
	***REMOVED***
	d.c, d.tok = 0, 0
	d.n.reset()
***REMOVED***

var jsonEncodeTerminate = []byte***REMOVED***' '***REMOVED***

func (h *JsonHandle) rpcEncodeTerminate() []byte ***REMOVED***
	return jsonEncodeTerminate
***REMOVED***

var _ decDriver = (*jsonDecDriver)(nil)
var _ encDriver = (*jsonEncDriver)(nil)
