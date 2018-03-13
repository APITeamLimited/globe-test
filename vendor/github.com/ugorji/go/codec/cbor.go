// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"math"
	"reflect"
	"time"
)

const (
	cborMajorUint byte = iota
	cborMajorNegInt
	cborMajorBytes
	cborMajorText
	cborMajorArray
	cborMajorMap
	cborMajorTag
	cborMajorOther
)

const (
	cborBdFalse byte = 0xf4 + iota
	cborBdTrue
	cborBdNil
	cborBdUndefined
	cborBdExt
	cborBdFloat16
	cborBdFloat32
	cborBdFloat64
)

const (
	cborBdIndefiniteBytes  byte = 0x5f
	cborBdIndefiniteString      = 0x7f
	cborBdIndefiniteArray       = 0x9f
	cborBdIndefiniteMap         = 0xbf
	cborBdBreak                 = 0xff
)

// These define some in-stream descriptors for
// manual encoding e.g. when doing explicit indefinite-length
const (
	CborStreamBytes  byte = 0x5f
	CborStreamString      = 0x7f
	CborStreamArray       = 0x9f
	CborStreamMap         = 0xbf
	CborStreamBreak       = 0xff
)

const (
	cborBaseUint   byte = 0x00
	cborBaseNegInt      = 0x20
	cborBaseBytes       = 0x40
	cborBaseString      = 0x60
	cborBaseArray       = 0x80
	cborBaseMap         = 0xa0
	cborBaseTag         = 0xc0
	cborBaseSimple      = 0xe0
)

// -------------------

type cborEncDriver struct ***REMOVED***
	noBuiltInTypes
	encDriverNoopContainerWriter
	// encNoSeparator
	e *Encoder
	w encWriter
	h *CborHandle
	x [8]byte
	_ [3]uint64 // padding
***REMOVED***

func (e *cborEncDriver) EncodeNil() ***REMOVED***
	e.w.writen1(cborBdNil)
***REMOVED***

func (e *cborEncDriver) EncodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.w.writen1(cborBdTrue)
	***REMOVED*** else ***REMOVED***
		e.w.writen1(cborBdFalse)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeFloat32(f float32) ***REMOVED***
	e.w.writen1(cborBdFloat32)
	bigenHelper***REMOVED***e.x[:4], e.w***REMOVED***.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *cborEncDriver) EncodeFloat64(f float64) ***REMOVED***
	e.w.writen1(cborBdFloat64)
	bigenHelper***REMOVED***e.x[:8], e.w***REMOVED***.writeUint64(math.Float64bits(f))
***REMOVED***

func (e *cborEncDriver) encUint(v uint64, bd byte) ***REMOVED***
	if v <= 0x17 ***REMOVED***
		e.w.writen1(byte(v) + bd)
	***REMOVED*** else if v <= math.MaxUint8 ***REMOVED***
		e.w.writen2(bd+0x18, uint8(v))
	***REMOVED*** else if v <= math.MaxUint16 ***REMOVED***
		e.w.writen1(bd + 0x19)
		bigenHelper***REMOVED***e.x[:2], e.w***REMOVED***.writeUint16(uint16(v))
	***REMOVED*** else if v <= math.MaxUint32 ***REMOVED***
		e.w.writen1(bd + 0x1a)
		bigenHelper***REMOVED***e.x[:4], e.w***REMOVED***.writeUint32(uint32(v))
	***REMOVED*** else ***REMOVED*** // if v <= math.MaxUint64 ***REMOVED***
		e.w.writen1(bd + 0x1b)
		bigenHelper***REMOVED***e.x[:8], e.w***REMOVED***.writeUint64(v)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeInt(v int64) ***REMOVED***
	if v < 0 ***REMOVED***
		e.encUint(uint64(-1-v), cborBaseNegInt)
	***REMOVED*** else ***REMOVED***
		e.encUint(uint64(v), cborBaseUint)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeUint(v uint64) ***REMOVED***
	e.encUint(v, cborBaseUint)
***REMOVED***

func (e *cborEncDriver) encLen(bd byte, length int) ***REMOVED***
	e.encUint(uint64(length), bd)
***REMOVED***

func (e *cborEncDriver) EncodeTime(t time.Time) ***REMOVED***
	if t.IsZero() ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else if e.h.TimeRFC3339 ***REMOVED***
		e.encUint(0, cborBaseTag)
		e.EncodeString(cUTF8, t.Format(time.RFC3339Nano))
	***REMOVED*** else ***REMOVED***
		e.encUint(1, cborBaseTag)
		t = t.UTC().Round(time.Microsecond)
		sec, nsec := t.Unix(), uint64(t.Nanosecond())
		if nsec == 0 ***REMOVED***
			e.EncodeInt(sec)
		***REMOVED*** else ***REMOVED***
			e.EncodeFloat64(float64(sec) + float64(nsec)/1e9)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext, en *Encoder) ***REMOVED***
	e.encUint(uint64(xtag), cborBaseTag)
	if v := ext.ConvertExt(rv); v == nil ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		en.encode(v)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeRawExt(re *RawExt, en *Encoder) ***REMOVED***
	e.encUint(uint64(re.Tag), cborBaseTag)
	if false && re.Data != nil ***REMOVED***
		en.encode(re.Data)
	***REMOVED*** else if re.Value != nil ***REMOVED***
		en.encode(re.Value)
	***REMOVED*** else ***REMOVED***
		e.EncodeNil()
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) WriteArrayStart(length int) ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		e.w.writen1(cborBdIndefiniteArray)
	***REMOVED*** else ***REMOVED***
		e.encLen(cborBaseArray, length)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) WriteMapStart(length int) ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		e.w.writen1(cborBdIndefiniteMap)
	***REMOVED*** else ***REMOVED***
		e.encLen(cborBaseMap, length)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) WriteMapEnd() ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		e.w.writen1(cborBdBreak)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) WriteArrayEnd() ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		e.w.writen1(cborBdBreak)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeString(c charEncoding, v string) ***REMOVED***
	e.encStringBytesS(cborBaseString, v)
***REMOVED***

func (e *cborEncDriver) EncodeStringBytes(c charEncoding, v []byte) ***REMOVED***
	if v == nil ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else if c == cRAW ***REMOVED***
		e.encStringBytesS(cborBaseBytes, stringView(v))
	***REMOVED*** else ***REMOVED***
		e.encStringBytesS(cborBaseString, stringView(v))
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) encStringBytesS(bb byte, v string) ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		if bb == cborBaseBytes ***REMOVED***
			e.w.writen1(cborBdIndefiniteBytes)
		***REMOVED*** else ***REMOVED***
			e.w.writen1(cborBdIndefiniteString)
		***REMOVED***
		blen := len(v) / 4
		if blen == 0 ***REMOVED***
			blen = 64
		***REMOVED*** else if blen > 1024 ***REMOVED***
			blen = 1024
		***REMOVED***
		for i := 0; i < len(v); ***REMOVED***
			var v2 string
			i2 := i + blen
			if i2 < len(v) ***REMOVED***
				v2 = v[i:i2]
			***REMOVED*** else ***REMOVED***
				v2 = v[i:]
			***REMOVED***
			e.encLen(bb, len(v2))
			e.w.writestr(v2)
			i = i2
		***REMOVED***
		e.w.writen1(cborBdBreak)
	***REMOVED*** else ***REMOVED***
		e.encLen(bb, len(v))
		e.w.writestr(v)
	***REMOVED***
***REMOVED***

// ----------------------

type cborDecDriver struct ***REMOVED***
	d *Decoder
	h *CborHandle
	r decReader
	// b      [scratchByteArrayLen]byte
	br     bool // bytes reader
	bdRead bool
	bd     byte
	noBuiltInTypes
	// decNoSeparator
	decDriverNoopContainerReader
	_ [3]uint64 // padding
***REMOVED***

func (d *cborDecDriver) readNextBd() ***REMOVED***
	d.bd = d.r.readn1()
	d.bdRead = true
***REMOVED***

func (d *cborDecDriver) uncacheRead() ***REMOVED***
	if d.bdRead ***REMOVED***
		d.r.unreadn1()
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *cborDecDriver) ContainerType() (vt valueType) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == cborBdNil ***REMOVED***
		return valueTypeNil
	***REMOVED*** else if d.bd == cborBdIndefiniteBytes || (d.bd >= cborBaseBytes && d.bd < cborBaseString) ***REMOVED***
		return valueTypeBytes
	***REMOVED*** else if d.bd == cborBdIndefiniteString || (d.bd >= cborBaseString && d.bd < cborBaseArray) ***REMOVED***
		return valueTypeString
	***REMOVED*** else if d.bd == cborBdIndefiniteArray || (d.bd >= cborBaseArray && d.bd < cborBaseMap) ***REMOVED***
		return valueTypeArray
	***REMOVED*** else if d.bd == cborBdIndefiniteMap || (d.bd >= cborBaseMap && d.bd < cborBaseTag) ***REMOVED***
		return valueTypeMap
	***REMOVED***
	// else ***REMOVED***
	// d.d.errorf("isContainerType: unsupported parameter: %v", vt)
	// ***REMOVED***
	return valueTypeUnset
***REMOVED***

func (d *cborDecDriver) TryDecodeAsNil() bool ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	// treat Nil and Undefined as nil values
	if d.bd == cborBdNil || d.bd == cborBdUndefined ***REMOVED***
		d.bdRead = false
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *cborDecDriver) CheckBreak() bool ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == cborBdBreak ***REMOVED***
		d.bdRead = false
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *cborDecDriver) decUint() (ui uint64) ***REMOVED***
	v := d.bd & 0x1f
	if v <= 0x17 ***REMOVED***
		ui = uint64(v)
	***REMOVED*** else ***REMOVED***
		if v == 0x18 ***REMOVED***
			ui = uint64(d.r.readn1())
		***REMOVED*** else if v == 0x19 ***REMOVED***
			ui = uint64(bigen.Uint16(d.r.readx(2)))
		***REMOVED*** else if v == 0x1a ***REMOVED***
			ui = uint64(bigen.Uint32(d.r.readx(4)))
		***REMOVED*** else if v == 0x1b ***REMOVED***
			ui = uint64(bigen.Uint64(d.r.readx(8)))
		***REMOVED*** else ***REMOVED***
			d.d.errorf("decUint: Invalid descriptor: %v", d.bd)
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (d *cborDecDriver) decCheckInteger() (neg bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	major := d.bd >> 5
	if major == cborMajorUint ***REMOVED***
	***REMOVED*** else if major == cborMajorNegInt ***REMOVED***
		neg = true
	***REMOVED*** else ***REMOVED***
		d.d.errorf("invalid major: %v (bd: %v)", major, d.bd)
		return
	***REMOVED***
	return
***REMOVED***

func (d *cborDecDriver) DecodeInt64() (i int64) ***REMOVED***
	neg := d.decCheckInteger()
	ui := d.decUint()
	// check if this number can be converted to an int without overflow
	if neg ***REMOVED***
		i = -(chkOvf.SignedIntV(ui + 1))
	***REMOVED*** else ***REMOVED***
		i = chkOvf.SignedIntV(ui)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *cborDecDriver) DecodeUint64() (ui uint64) ***REMOVED***
	if d.decCheckInteger() ***REMOVED***
		d.d.errorf("Assigning negative signed value to unsigned type")
		return
	***REMOVED***
	ui = d.decUint()
	d.bdRead = false
	return
***REMOVED***

func (d *cborDecDriver) DecodeFloat64() (f float64) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if bd := d.bd; bd == cborBdFloat16 ***REMOVED***
		f = float64(math.Float32frombits(halfFloatToFloatBits(bigen.Uint16(d.r.readx(2)))))
	***REMOVED*** else if bd == cborBdFloat32 ***REMOVED***
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readx(4))))
	***REMOVED*** else if bd == cborBdFloat64 ***REMOVED***
		f = math.Float64frombits(bigen.Uint64(d.r.readx(8)))
	***REMOVED*** else if bd >= cborBaseUint && bd < cborBaseBytes ***REMOVED***
		f = float64(d.DecodeInt64())
	***REMOVED*** else ***REMOVED***
		d.d.errorf("Float only valid from float16/32/64: Invalid descriptor: %v", bd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool only (single byte).
func (d *cborDecDriver) DecodeBool() (b bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if bd := d.bd; bd == cborBdTrue ***REMOVED***
		b = true
	***REMOVED*** else if bd == cborBdFalse ***REMOVED***
	***REMOVED*** else ***REMOVED***
		d.d.errorf("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *cborDecDriver) ReadMapStart() (length int) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap ***REMOVED***
		return -1
	***REMOVED***
	return d.decLen()
***REMOVED***

func (d *cborDecDriver) ReadArrayStart() (length int) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	d.bdRead = false
	if d.bd == cborBdIndefiniteArray ***REMOVED***
		return -1
	***REMOVED***
	return d.decLen()
***REMOVED***

func (d *cborDecDriver) decLen() int ***REMOVED***
	return int(d.decUint())
***REMOVED***

func (d *cborDecDriver) decAppendIndefiniteBytes(bs []byte) []byte ***REMOVED***
	d.bdRead = false
	for ***REMOVED***
		if d.CheckBreak() ***REMOVED***
			break
		***REMOVED***
		if major := d.bd >> 5; major != cborMajorBytes && major != cborMajorText ***REMOVED***
			d.d.errorf("expect bytes/string major type in indefinite string/bytes;"+
				" got: %v, byte: %v", major, d.bd)
			return nil
		***REMOVED***
		n := d.decLen()
		oldLen := len(bs)
		newLen := oldLen + n
		if newLen > cap(bs) ***REMOVED***
			bs2 := make([]byte, newLen, 2*cap(bs)+n)
			copy(bs2, bs)
			bs = bs2
		***REMOVED*** else ***REMOVED***
			bs = bs[:newLen]
		***REMOVED***
		d.r.readb(bs[oldLen:newLen])
		// bs = append(bs, d.r.readn()...)
		d.bdRead = false
	***REMOVED***
	d.bdRead = false
	return bs
***REMOVED***

func (d *cborDecDriver) DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == cborBdNil || d.bd == cborBdUndefined ***REMOVED***
		d.bdRead = false
		return nil
	***REMOVED***
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString ***REMOVED***
		d.bdRead = false
		if bs == nil ***REMOVED***
			if zerocopy ***REMOVED***
				return d.decAppendIndefiniteBytes(d.d.b[:0])
			***REMOVED***
			return d.decAppendIndefiniteBytes(zeroByteSlice)
		***REMOVED***
		return d.decAppendIndefiniteBytes(bs[:0])
	***REMOVED***
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
	if d.bd == cborBdIndefiniteArray || (d.bd >= cborBaseArray && d.bd < cborBaseMap) ***REMOVED***
		bsOut, _ = fastpathTV.DecSliceUint8V(bs, true, d.d)
		return
	***REMOVED***
	clen := d.decLen()
	d.bdRead = false
	if zerocopy ***REMOVED***
		if d.br ***REMOVED***
			return d.r.readx(clen)
		***REMOVED*** else if len(bs) == 0 ***REMOVED***
			bs = d.d.b[:]
		***REMOVED***
	***REMOVED***
	return decByteSlice(d.r, clen, d.h.MaxInitLen, bs)
***REMOVED***

func (d *cborDecDriver) DecodeString() (s string) ***REMOVED***
	return string(d.DecodeBytes(d.d.b[:], true))
***REMOVED***

func (d *cborDecDriver) DecodeStringAsBytes() (s []byte) ***REMOVED***
	return d.DecodeBytes(d.d.b[:], true)
***REMOVED***

func (d *cborDecDriver) DecodeTime() (t time.Time) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == cborBdNil || d.bd == cborBdUndefined ***REMOVED***
		d.bdRead = false
		return
	***REMOVED***
	xtag := d.decUint()
	d.bdRead = false
	return d.decodeTime(xtag)
***REMOVED***

func (d *cborDecDriver) decodeTime(xtag uint64) (t time.Time) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	switch xtag ***REMOVED***
	case 0:
		var err error
		if t, err = time.Parse(time.RFC3339, stringView(d.DecodeStringAsBytes())); err != nil ***REMOVED***
			d.d.errorv(err)
		***REMOVED***
	case 1:
		// decode an int64 or a float, and infer time.Time from there.
		// for floats, round to microseconds, as that is what is guaranteed to fit well.
		switch ***REMOVED***
		case d.bd == cborBdFloat16, d.bd == cborBdFloat32:
			f1, f2 := math.Modf(d.DecodeFloat64())
			t = time.Unix(int64(f1), int64(f2*1e9))
		case d.bd == cborBdFloat64:
			f1, f2 := math.Modf(d.DecodeFloat64())
			t = time.Unix(int64(f1), int64(f2*1e9))
		case d.bd >= cborBaseUint && d.bd < cborBaseNegInt,
			d.bd >= cborBaseNegInt && d.bd < cborBaseBytes:
			t = time.Unix(d.DecodeInt64(), 0)
		default:
			d.d.errorf("time.Time can only be decoded from a number (or RFC3339 string)")
		***REMOVED***
	default:
		d.d.errorf("invalid tag for time.Time - expecting 0 or 1, got 0x%x", xtag)
	***REMOVED***
	t = t.UTC().Round(time.Microsecond)
	return
***REMOVED***

func (d *cborDecDriver) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) (realxtag uint64) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	u := d.decUint()
	d.bdRead = false
	realxtag = u
	if ext == nil ***REMOVED***
		re := rv.(*RawExt)
		re.Tag = realxtag
		d.d.decode(&re.Value)
	***REMOVED*** else if xtag != realxtag ***REMOVED***
		d.d.errorf("Wrong extension tag. Got %b. Expecting: %v", realxtag, xtag)
		return
	***REMOVED*** else ***REMOVED***
		var v interface***REMOVED******REMOVED***
		d.d.decode(&v)
		ext.UpdateExt(rv, v)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *cborDecDriver) DecodeNaked() ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***

	n := d.d.n
	var decodeFurther bool

	switch d.bd ***REMOVED***
	case cborBdNil:
		n.v = valueTypeNil
	case cborBdFalse:
		n.v = valueTypeBool
		n.b = false
	case cborBdTrue:
		n.v = valueTypeBool
		n.b = true
	case cborBdFloat16, cborBdFloat32, cborBdFloat64:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case cborBdIndefiniteBytes:
		n.v = valueTypeBytes
		n.l = d.DecodeBytes(nil, false)
	case cborBdIndefiniteString:
		n.v = valueTypeString
		n.s = d.DecodeString()
	case cborBdIndefiniteArray:
		n.v = valueTypeArray
		decodeFurther = true
	case cborBdIndefiniteMap:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		switch ***REMOVED***
		case d.bd >= cborBaseUint && d.bd < cborBaseNegInt:
			if d.h.SignedInteger ***REMOVED***
				n.v = valueTypeInt
				n.i = d.DecodeInt64()
			***REMOVED*** else ***REMOVED***
				n.v = valueTypeUint
				n.u = d.DecodeUint64()
			***REMOVED***
		case d.bd >= cborBaseNegInt && d.bd < cborBaseBytes:
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		case d.bd >= cborBaseBytes && d.bd < cborBaseString:
			n.v = valueTypeBytes
			n.l = d.DecodeBytes(nil, false)
		case d.bd >= cborBaseString && d.bd < cborBaseArray:
			n.v = valueTypeString
			n.s = d.DecodeString()
		case d.bd >= cborBaseArray && d.bd < cborBaseMap:
			n.v = valueTypeArray
			decodeFurther = true
		case d.bd >= cborBaseMap && d.bd < cborBaseTag:
			n.v = valueTypeMap
			decodeFurther = true
		case d.bd >= cborBaseTag && d.bd < cborBaseSimple:
			n.v = valueTypeExt
			n.u = d.decUint()
			n.l = nil
			if n.u == 0 || n.u == 1 ***REMOVED***
				d.bdRead = false
				n.v = valueTypeTime
				n.t = d.decodeTime(n.u)
			***REMOVED***
			// d.bdRead = false
			// d.d.decode(&re.Value) // handled by decode itself.
			// decodeFurther = true
		default:
			d.d.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
			return
		***REMOVED***
	***REMOVED***

	if !decodeFurther ***REMOVED***
		d.bdRead = false
	***REMOVED***
	return
***REMOVED***

// -------------------------

// CborHandle is a Handle for the CBOR encoding format,
// defined at http://tools.ietf.org/html/rfc7049 and documented further at http://cbor.io .
//
// CBOR is comprehensively supported, including support for:
//   - indefinite-length arrays/maps/bytes/strings
//   - (extension) tags in range 0..0xffff (0 .. 65535)
//   - half, single and double-precision floats
//   - all numbers (1, 2, 4 and 8-byte signed and unsigned integers)
//   - nil, true, false, ...
//   - arrays and maps, bytes and text strings
//
// None of the optional extensions (with tags) defined in the spec are supported out-of-the-box.
// Users can implement them as needed (using SetExt), including spec-documented ones:
//   - timestamp, BigNum, BigFloat, Decimals,
//   - Encoded Text (e.g. URL, regexp, base64, MIME Message), etc.
type CborHandle struct ***REMOVED***
	binaryEncodingType
	noElemSeparators
	BasicHandle

	// IndefiniteLength=true, means that we encode using indefinitelength
	IndefiniteLength bool

	// TimeRFC3339 says to encode time.Time using RFC3339 format.
	// If unset, we encode time.Time using seconds past epoch.
	TimeRFC3339 bool

	_ [1]uint64 // padding
***REMOVED***

// Name returns the name of the handle: cbor
func (h *CborHandle) Name() string ***REMOVED*** return "cbor" ***REMOVED***

// SetInterfaceExt sets an extension
func (h *CborHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, &extWrapper***REMOVED***bytesExtFailer***REMOVED******REMOVED***, ext***REMOVED***)
***REMOVED***

func (h *CborHandle) newEncDriver(e *Encoder) encDriver ***REMOVED***
	return &cborEncDriver***REMOVED***e: e, w: e.w, h: h***REMOVED***
***REMOVED***

func (h *CborHandle) newDecDriver(d *Decoder) decDriver ***REMOVED***
	return &cborDecDriver***REMOVED***d: d, h: h, r: d.r, br: d.bytes***REMOVED***
***REMOVED***

func (e *cborEncDriver) reset() ***REMOVED***
	e.w = e.e.w
***REMOVED***

func (d *cborDecDriver) reset() ***REMOVED***
	d.r, d.br = d.d.r, d.d.bytes
	d.bd, d.bdRead = 0, false
***REMOVED***

var _ decDriver = (*cborDecDriver)(nil)
var _ encDriver = (*cborEncDriver)(nil)
