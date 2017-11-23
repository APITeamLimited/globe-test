// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"math"
	"reflect"
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
	encNoSeparator
	e *Encoder
	w encWriter
	h *CborHandle
	x [8]byte
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
	if re.Data != nil ***REMOVED***
		en.encode(re.Data)
	***REMOVED*** else if re.Value == nil ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		en.encode(re.Value)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeArrayStart(length int) ***REMOVED***
	e.encLen(cborBaseArray, length)
***REMOVED***

func (e *cborEncDriver) EncodeMapStart(length int) ***REMOVED***
	e.encLen(cborBaseMap, length)
***REMOVED***

func (e *cborEncDriver) EncodeString(c charEncoding, v string) ***REMOVED***
	e.encLen(cborBaseString, len(v))
	e.w.writestr(v)
***REMOVED***

func (e *cborEncDriver) EncodeSymbol(v string) ***REMOVED***
	e.EncodeString(c_UTF8, v)
***REMOVED***

func (e *cborEncDriver) EncodeStringBytes(c charEncoding, v []byte) ***REMOVED***
	if c == c_RAW ***REMOVED***
		e.encLen(cborBaseBytes, len(v))
	***REMOVED*** else ***REMOVED***
		e.encLen(cborBaseString, len(v))
	***REMOVED***
	e.w.writeb(v)
***REMOVED***

// ----------------------

type cborDecDriver struct ***REMOVED***
	d      *Decoder
	h      *CborHandle
	r      decReader
	b      [scratchByteArrayLen]byte
	br     bool // bytes reader
	bdRead bool
	bd     byte
	noBuiltInTypes
	decNoSeparator
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
	***REMOVED*** else ***REMOVED***
		// d.d.errorf("isContainerType: unsupported parameter: %v", vt)
	***REMOVED***
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

func (d *cborDecDriver) DecodeInt(bitsize uint8) (i int64) ***REMOVED***
	neg := d.decCheckInteger()
	ui := d.decUint()
	// check if this number can be converted to an int without overflow
	var overflow bool
	if neg ***REMOVED***
		if i, overflow = chkOvf.SignedInt(ui + 1); overflow ***REMOVED***
			d.d.errorf("cbor: overflow converting %v to signed integer", ui+1)
			return
		***REMOVED***
		i = -i
	***REMOVED*** else ***REMOVED***
		if i, overflow = chkOvf.SignedInt(ui); overflow ***REMOVED***
			d.d.errorf("cbor: overflow converting %v to signed integer", ui)
			return
		***REMOVED***
	***REMOVED***
	if chkOvf.Int(i, bitsize) ***REMOVED***
		d.d.errorf("cbor: overflow integer: %v", i)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *cborDecDriver) DecodeUint(bitsize uint8) (ui uint64) ***REMOVED***
	if d.decCheckInteger() ***REMOVED***
		d.d.errorf("Assigning negative signed value to unsigned type")
		return
	***REMOVED***
	ui = d.decUint()
	if chkOvf.Uint(ui, bitsize) ***REMOVED***
		d.d.errorf("cbor: overflow integer: %v", ui)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *cborDecDriver) DecodeFloat(chkOverflow32 bool) (f float64) ***REMOVED***
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
		f = float64(d.DecodeInt(64))
	***REMOVED*** else ***REMOVED***
		d.d.errorf("Float only valid from float16/32/64: Invalid descriptor: %v", bd)
		return
	***REMOVED***
	if chkOverflow32 && chkOvf.Float32(f) ***REMOVED***
		d.d.errorf("cbor: float32 overflow: %v", f)
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
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap ***REMOVED***
		return -1
	***REMOVED***
	return d.decLen()
***REMOVED***

func (d *cborDecDriver) ReadArrayStart() (length int) ***REMOVED***
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
			d.d.errorf("cbor: expect bytes or string major type in indefinite string/bytes; got: %v, byte: %v", major, d.bd)
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

func (d *cborDecDriver) DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == cborBdNil || d.bd == cborBdUndefined ***REMOVED***
		d.bdRead = false
		return nil
	***REMOVED***
	if d.bd == cborBdIndefiniteBytes || d.bd == cborBdIndefiniteString ***REMOVED***
		if bs == nil ***REMOVED***
			return d.decAppendIndefiniteBytes(nil)
		***REMOVED***
		return d.decAppendIndefiniteBytes(bs[:0])
	***REMOVED***
	clen := d.decLen()
	d.bdRead = false
	if zerocopy ***REMOVED***
		if d.br ***REMOVED***
			return d.r.readx(clen)
		***REMOVED*** else if len(bs) == 0 ***REMOVED***
			bs = d.b[:]
		***REMOVED***
	***REMOVED***
	return decByteSlice(d.r, clen, d.d.h.MaxInitLen, bs)
***REMOVED***

func (d *cborDecDriver) DecodeString() (s string) ***REMOVED***
	return string(d.DecodeBytes(d.b[:], true, true))
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

	n := &d.d.n
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
	case cborBdFloat16, cborBdFloat32:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat(true)
	case cborBdFloat64:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat(false)
	case cborBdIndefiniteBytes:
		n.v = valueTypeBytes
		n.l = d.DecodeBytes(nil, false, false)
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
				n.i = d.DecodeInt(64)
			***REMOVED*** else ***REMOVED***
				n.v = valueTypeUint
				n.u = d.DecodeUint(64)
			***REMOVED***
		case d.bd >= cborBaseNegInt && d.bd < cborBaseBytes:
			n.v = valueTypeInt
			n.i = d.DecodeInt(64)
		case d.bd >= cborBaseBytes && d.bd < cborBaseString:
			n.v = valueTypeBytes
			n.l = d.DecodeBytes(nil, false, false)
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
//   - timestamp, BigNum, BigFloat, Decimals, Encoded Text (e.g. URL, regexp, base64, MIME Message), etc.
//
// To encode with indefinite lengths (streaming), users will use
// (Must)Encode methods of *Encoder, along with writing CborStreamXXX constants.
//
// For example, to encode "one-byte" as an indefinite length string:
//     var buf bytes.Buffer
//     e := NewEncoder(&buf, new(CborHandle))
//     buf.WriteByte(CborStreamString)
//     e.MustEncode("one-")
//     e.MustEncode("byte")
//     buf.WriteByte(CborStreamBreak)
//     encodedBytes := buf.Bytes()
//     var vv interface***REMOVED******REMOVED***
//     NewDecoderBytes(buf.Bytes(), new(CborHandle)).MustDecode(&vv)
//     // Now, vv contains the same string "one-byte"
//
type CborHandle struct ***REMOVED***
	binaryEncodingType
	BasicHandle
***REMOVED***

func (h *CborHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, &setExtWrapper***REMOVED***i: ext***REMOVED***)
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
