// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"math"
	"reflect"
	"time"
)

const (
	_               uint8 = iota
	simpleVdNil           = 1
	simpleVdFalse         = 2
	simpleVdTrue          = 3
	simpleVdFloat32       = 4
	simpleVdFloat64       = 5

	// each lasts for 4 (ie n, n+1, n+2, n+3)
	simpleVdPosInt = 8
	simpleVdNegInt = 12

	simpleVdTime = 24

	// containers: each lasts for 4 (ie n, n+1, n+2, ... n+7)
	simpleVdString    = 216
	simpleVdByteArray = 224
	simpleVdArray     = 232
	simpleVdMap       = 240
	simpleVdExt       = 248
)

type simpleEncDriver struct ***REMOVED***
	noBuiltInTypes
	// encNoSeparator
	e *Encoder
	h *SimpleHandle
	w encWriter
	b [8]byte
	// c containerState
	encDriverTrackContainerWriter
	// encDriverNoopContainerWriter
	_ [2]uint64 // padding
***REMOVED***

func (e *simpleEncDriver) EncodeNil() ***REMOVED***
	e.w.writen1(simpleVdNil)
***REMOVED***

func (e *simpleEncDriver) EncodeBool(b bool) ***REMOVED***
	if e.h.EncZeroValuesAsNil && e.c != containerMapKey && !b ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	if b ***REMOVED***
		e.w.writen1(simpleVdTrue)
	***REMOVED*** else ***REMOVED***
		e.w.writen1(simpleVdFalse)
	***REMOVED***
***REMOVED***

func (e *simpleEncDriver) EncodeFloat32(f float32) ***REMOVED***
	if e.h.EncZeroValuesAsNil && e.c != containerMapKey && f == 0.0 ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	e.w.writen1(simpleVdFloat32)
	bigenHelper***REMOVED***e.b[:4], e.w***REMOVED***.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *simpleEncDriver) EncodeFloat64(f float64) ***REMOVED***
	if e.h.EncZeroValuesAsNil && e.c != containerMapKey && f == 0.0 ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	e.w.writen1(simpleVdFloat64)
	bigenHelper***REMOVED***e.b[:8], e.w***REMOVED***.writeUint64(math.Float64bits(f))
***REMOVED***

func (e *simpleEncDriver) EncodeInt(v int64) ***REMOVED***
	if v < 0 ***REMOVED***
		e.encUint(uint64(-v), simpleVdNegInt)
	***REMOVED*** else ***REMOVED***
		e.encUint(uint64(v), simpleVdPosInt)
	***REMOVED***
***REMOVED***

func (e *simpleEncDriver) EncodeUint(v uint64) ***REMOVED***
	e.encUint(v, simpleVdPosInt)
***REMOVED***

func (e *simpleEncDriver) encUint(v uint64, bd uint8) ***REMOVED***
	if e.h.EncZeroValuesAsNil && e.c != containerMapKey && v == 0 ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	if v <= math.MaxUint8 ***REMOVED***
		e.w.writen2(bd, uint8(v))
	***REMOVED*** else if v <= math.MaxUint16 ***REMOVED***
		e.w.writen1(bd + 1)
		bigenHelper***REMOVED***e.b[:2], e.w***REMOVED***.writeUint16(uint16(v))
	***REMOVED*** else if v <= math.MaxUint32 ***REMOVED***
		e.w.writen1(bd + 2)
		bigenHelper***REMOVED***e.b[:4], e.w***REMOVED***.writeUint32(uint32(v))
	***REMOVED*** else ***REMOVED*** // if v <= math.MaxUint64 ***REMOVED***
		e.w.writen1(bd + 3)
		bigenHelper***REMOVED***e.b[:8], e.w***REMOVED***.writeUint64(v)
	***REMOVED***
***REMOVED***

func (e *simpleEncDriver) encLen(bd byte, length int) ***REMOVED***
	if length == 0 ***REMOVED***
		e.w.writen1(bd)
	***REMOVED*** else if length <= math.MaxUint8 ***REMOVED***
		e.w.writen1(bd + 1)
		e.w.writen1(uint8(length))
	***REMOVED*** else if length <= math.MaxUint16 ***REMOVED***
		e.w.writen1(bd + 2)
		bigenHelper***REMOVED***e.b[:2], e.w***REMOVED***.writeUint16(uint16(length))
	***REMOVED*** else if int64(length) <= math.MaxUint32 ***REMOVED***
		e.w.writen1(bd + 3)
		bigenHelper***REMOVED***e.b[:4], e.w***REMOVED***.writeUint32(uint32(length))
	***REMOVED*** else ***REMOVED***
		e.w.writen1(bd + 4)
		bigenHelper***REMOVED***e.b[:8], e.w***REMOVED***.writeUint64(uint64(length))
	***REMOVED***
***REMOVED***

func (e *simpleEncDriver) EncodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext, _ *Encoder) ***REMOVED***
	bs := ext.WriteExt(rv)
	if bs == nil ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	e.encodeExtPreamble(uint8(xtag), len(bs))
	e.w.writeb(bs)
***REMOVED***

func (e *simpleEncDriver) EncodeRawExt(re *RawExt, _ *Encoder) ***REMOVED***
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
***REMOVED***

func (e *simpleEncDriver) encodeExtPreamble(xtag byte, length int) ***REMOVED***
	e.encLen(simpleVdExt, length)
	e.w.writen1(xtag)
***REMOVED***

func (e *simpleEncDriver) WriteArrayStart(length int) ***REMOVED***
	e.c = containerArrayStart
	e.encLen(simpleVdArray, length)
***REMOVED***

func (e *simpleEncDriver) WriteMapStart(length int) ***REMOVED***
	e.c = containerMapStart
	e.encLen(simpleVdMap, length)
***REMOVED***

func (e *simpleEncDriver) EncodeString(c charEncoding, v string) ***REMOVED***
	if false && e.h.EncZeroValuesAsNil && e.c != containerMapKey && v == "" ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	e.encLen(simpleVdString, len(v))
	e.w.writestr(v)
***REMOVED***

// func (e *simpleEncDriver) EncodeSymbol(v string) ***REMOVED***
// 	e.EncodeString(cUTF8, v)
// ***REMOVED***

func (e *simpleEncDriver) EncodeStringBytes(c charEncoding, v []byte) ***REMOVED***
	// if e.h.EncZeroValuesAsNil && e.c != containerMapKey && v == nil ***REMOVED***
	if v == nil ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	e.encLen(simpleVdByteArray, len(v))
	e.w.writeb(v)
***REMOVED***

func (e *simpleEncDriver) EncodeTime(t time.Time) ***REMOVED***
	// if e.h.EncZeroValuesAsNil && e.c != containerMapKey && t.IsZero() ***REMOVED***
	if t.IsZero() ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	v, err := t.MarshalBinary()
	if err != nil ***REMOVED***
		e.e.errorv(err)
		return
	***REMOVED***
	// time.Time marshalbinary takes about 14 bytes.
	e.w.writen2(simpleVdTime, uint8(len(v)))
	e.w.writeb(v)
***REMOVED***

//------------------------------------

type simpleDecDriver struct ***REMOVED***
	d      *Decoder
	h      *SimpleHandle
	r      decReader
	bdRead bool
	bd     byte
	br     bool // a bytes reader?
	c      containerState
	// b      [scratchByteArrayLen]byte
	noBuiltInTypes
	// noStreamingCodec
	decDriverNoopContainerReader
	_ [3]uint64 // padding
***REMOVED***

func (d *simpleDecDriver) readNextBd() ***REMOVED***
	d.bd = d.r.readn1()
	d.bdRead = true
***REMOVED***

func (d *simpleDecDriver) uncacheRead() ***REMOVED***
	if d.bdRead ***REMOVED***
		d.r.unreadn1()
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *simpleDecDriver) ContainerType() (vt valueType) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	switch d.bd ***REMOVED***
	case simpleVdNil:
		return valueTypeNil
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		return valueTypeBytes
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		return valueTypeString
	case simpleVdArray, simpleVdArray + 1,
		simpleVdArray + 2, simpleVdArray + 3, simpleVdArray + 4:
		return valueTypeArray
	case simpleVdMap, simpleVdMap + 1,
		simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		return valueTypeMap
		// case simpleVdTime:
		// 	return valueTypeTime
	***REMOVED***
	// else ***REMOVED***
	// d.d.errorf("isContainerType: unsupported parameter: %v", vt)
	// ***REMOVED***
	return valueTypeUnset
***REMOVED***

func (d *simpleDecDriver) TryDecodeAsNil() bool ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == simpleVdNil ***REMOVED***
		d.bdRead = false
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *simpleDecDriver) decCheckInteger() (ui uint64, neg bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	switch d.bd ***REMOVED***
	case simpleVdPosInt:
		ui = uint64(d.r.readn1())
	case simpleVdPosInt + 1:
		ui = uint64(bigen.Uint16(d.r.readx(2)))
	case simpleVdPosInt + 2:
		ui = uint64(bigen.Uint32(d.r.readx(4)))
	case simpleVdPosInt + 3:
		ui = uint64(bigen.Uint64(d.r.readx(8)))
	case simpleVdNegInt:
		ui = uint64(d.r.readn1())
		neg = true
	case simpleVdNegInt + 1:
		ui = uint64(bigen.Uint16(d.r.readx(2)))
		neg = true
	case simpleVdNegInt + 2:
		ui = uint64(bigen.Uint32(d.r.readx(4)))
		neg = true
	case simpleVdNegInt + 3:
		ui = uint64(bigen.Uint64(d.r.readx(8)))
		neg = true
	default:
		d.d.errorf("Integer only valid from pos/neg integer1..8. Invalid descriptor: %v", d.bd)
		return
	***REMOVED***
	// don't do this check, because callers may only want the unsigned value.
	// if ui > math.MaxInt64 ***REMOVED***
	// 	d.d.errorf("decIntAny: Integer out of range for signed int64: %v", ui)
	//		return
	// ***REMOVED***
	return
***REMOVED***

func (d *simpleDecDriver) DecodeInt64() (i int64) ***REMOVED***
	ui, neg := d.decCheckInteger()
	i = chkOvf.SignedIntV(ui)
	if neg ***REMOVED***
		i = -i
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) DecodeUint64() (ui uint64) ***REMOVED***
	ui, neg := d.decCheckInteger()
	if neg ***REMOVED***
		d.d.errorf("Assigning negative signed value to unsigned type")
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) DecodeFloat64() (f float64) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == simpleVdFloat32 ***REMOVED***
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readx(4))))
	***REMOVED*** else if d.bd == simpleVdFloat64 ***REMOVED***
		f = math.Float64frombits(bigen.Uint64(d.r.readx(8)))
	***REMOVED*** else ***REMOVED***
		if d.bd >= simpleVdPosInt && d.bd <= simpleVdNegInt+3 ***REMOVED***
			f = float64(d.DecodeInt64())
		***REMOVED*** else ***REMOVED***
			d.d.errorf("Float only valid from float32/64: Invalid descriptor: %v", d.bd)
			return
		***REMOVED***
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool only (single byte).
func (d *simpleDecDriver) DecodeBool() (b bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == simpleVdTrue ***REMOVED***
		b = true
	***REMOVED*** else if d.bd == simpleVdFalse ***REMOVED***
	***REMOVED*** else ***REMOVED***
		d.d.errorf("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) ReadMapStart() (length int) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	d.bdRead = false
	d.c = containerMapStart
	return d.decLen()
***REMOVED***

func (d *simpleDecDriver) ReadArrayStart() (length int) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	d.bdRead = false
	d.c = containerArrayStart
	return d.decLen()
***REMOVED***

func (d *simpleDecDriver) ReadArrayElem() ***REMOVED***
	d.c = containerArrayElem
***REMOVED***

func (d *simpleDecDriver) ReadArrayEnd() ***REMOVED***
	d.c = containerArrayEnd
***REMOVED***

func (d *simpleDecDriver) ReadMapElemKey() ***REMOVED***
	d.c = containerMapKey
***REMOVED***

func (d *simpleDecDriver) ReadMapElemValue() ***REMOVED***
	d.c = containerMapValue
***REMOVED***

func (d *simpleDecDriver) ReadMapEnd() ***REMOVED***
	d.c = containerMapEnd
***REMOVED***

func (d *simpleDecDriver) decLen() int ***REMOVED***
	switch d.bd % 8 ***REMOVED***
	case 0:
		return 0
	case 1:
		return int(d.r.readn1())
	case 2:
		return int(bigen.Uint16(d.r.readx(2)))
	case 3:
		ui := uint64(bigen.Uint32(d.r.readx(4)))
		if chkOvf.Uint(ui, intBitsize) ***REMOVED***
			d.d.errorf("overflow integer: %v", ui)
			return 0
		***REMOVED***
		return int(ui)
	case 4:
		ui := bigen.Uint64(d.r.readx(8))
		if chkOvf.Uint(ui, intBitsize) ***REMOVED***
			d.d.errorf("overflow integer: %v", ui)
			return 0
		***REMOVED***
		return int(ui)
	***REMOVED***
	d.d.errorf("decLen: Cannot read length: bd%%8 must be in range 0..4. Got: %d", d.bd%8)
	return -1
***REMOVED***

func (d *simpleDecDriver) DecodeString() (s string) ***REMOVED***
	return string(d.DecodeBytes(d.d.b[:], true))
***REMOVED***

func (d *simpleDecDriver) DecodeStringAsBytes() (s []byte) ***REMOVED***
	return d.DecodeBytes(d.d.b[:], true)
***REMOVED***

func (d *simpleDecDriver) DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == simpleVdNil ***REMOVED***
		d.bdRead = false
		return
	***REMOVED***
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
	if d.bd >= simpleVdArray && d.bd <= simpleVdMap+4 ***REMOVED***
		if len(bs) == 0 && zerocopy ***REMOVED***
			bs = d.d.b[:]
		***REMOVED***
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
	return decByteSlice(d.r, clen, d.d.h.MaxInitLen, bs)
***REMOVED***

func (d *simpleDecDriver) DecodeTime() (t time.Time) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == simpleVdNil ***REMOVED***
		d.bdRead = false
		return
	***REMOVED***
	if d.bd != simpleVdTime ***REMOVED***
		d.d.errorf("invalid descriptor for time.Time - expect 0x%x, received 0x%x", simpleVdTime, d.bd)
		return
	***REMOVED***
	d.bdRead = false
	clen := int(d.r.readn1())
	b := d.r.readx(clen)
	if err := (&t).UnmarshalBinary(b); err != nil ***REMOVED***
		d.d.errorv(err)
	***REMOVED***
	return
***REMOVED***

func (d *simpleDecDriver) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) (realxtag uint64) ***REMOVED***
	if xtag > 0xff ***REMOVED***
		d.d.errorf("decodeExt: tag must be <= 0xff; got: %v", xtag)
		return
	***REMOVED***
	realxtag1, xbs := d.decodeExtV(ext != nil, uint8(xtag))
	realxtag = uint64(realxtag1)
	if ext == nil ***REMOVED***
		re := rv.(*RawExt)
		re.Tag = realxtag
		re.Data = detachZeroCopyBytes(d.br, re.Data, xbs)
	***REMOVED*** else ***REMOVED***
		ext.ReadExt(rv, xbs)
	***REMOVED***
	return
***REMOVED***

func (d *simpleDecDriver) decodeExtV(verifyTag bool, tag byte) (xtag byte, xbs []byte) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	switch d.bd ***REMOVED***
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag ***REMOVED***
			d.d.errorf("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
			return
		***REMOVED***
		xbs = d.r.readx(l)
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		xbs = d.DecodeBytes(nil, true)
	default:
		d.d.errorf("Invalid descriptor - expecting extensions/bytearray, got: 0x%x", d.bd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *simpleDecDriver) DecodeNaked() ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***

	n := d.d.n
	var decodeFurther bool

	switch d.bd ***REMOVED***
	case simpleVdNil:
		n.v = valueTypeNil
	case simpleVdFalse:
		n.v = valueTypeBool
		n.b = false
	case simpleVdTrue:
		n.v = valueTypeBool
		n.b = true
	case simpleVdPosInt, simpleVdPosInt + 1, simpleVdPosInt + 2, simpleVdPosInt + 3:
		if d.h.SignedInteger ***REMOVED***
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		***REMOVED*** else ***REMOVED***
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		***REMOVED***
	case simpleVdNegInt, simpleVdNegInt + 1, simpleVdNegInt + 2, simpleVdNegInt + 3:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case simpleVdFloat32:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdFloat64:
		n.v = valueTypeFloat
		n.f = d.DecodeFloat64()
	case simpleVdTime:
		n.v = valueTypeTime
		n.t = d.DecodeTime()
	case simpleVdString, simpleVdString + 1,
		simpleVdString + 2, simpleVdString + 3, simpleVdString + 4:
		n.v = valueTypeString
		n.s = d.DecodeString()
	case simpleVdByteArray, simpleVdByteArray + 1,
		simpleVdByteArray + 2, simpleVdByteArray + 3, simpleVdByteArray + 4:
		n.v = valueTypeBytes
		n.l = d.DecodeBytes(nil, false)
	case simpleVdExt, simpleVdExt + 1, simpleVdExt + 2, simpleVdExt + 3, simpleVdExt + 4:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		n.l = d.r.readx(l)
	case simpleVdArray, simpleVdArray + 1, simpleVdArray + 2,
		simpleVdArray + 3, simpleVdArray + 4:
		n.v = valueTypeArray
		decodeFurther = true
	case simpleVdMap, simpleVdMap + 1, simpleVdMap + 2, simpleVdMap + 3, simpleVdMap + 4:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		d.d.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
	***REMOVED***

	if !decodeFurther ***REMOVED***
		d.bdRead = false
	***REMOVED***
	return
***REMOVED***

//------------------------------------

// SimpleHandle is a Handle for a very simple encoding format.
//
// simple is a simplistic codec similar to binc, but not as compact.
//   - Encoding of a value is always preceded by the descriptor byte (bd)
//   - True, false, nil are encoded fully in 1 byte (the descriptor)
//   - Integers (intXXX, uintXXX) are encoded in 1, 2, 4 or 8 bytes (plus a descriptor byte).
//     There are positive (uintXXX and intXXX >= 0) and negative (intXXX < 0) integers.
//   - Floats are encoded in 4 or 8 bytes (plus a descriptor byte)
//   - Length of containers (strings, bytes, array, map, extensions)
//     are encoded in 0, 1, 2, 4 or 8 bytes.
//     Zero-length containers have no length encoded.
//     For others, the number of bytes is given by pow(2, bd%3)
//   - maps are encoded as [bd] [length] [[key][value]]...
//   - arrays are encoded as [bd] [length] [value]...
//   - extensions are encoded as [bd] [length] [tag] [byte]...
//   - strings/bytearrays are encoded as [bd] [length] [byte]...
//   - time.Time are encoded as [bd] [length] [byte]...
//
// The full spec will be published soon.
type SimpleHandle struct ***REMOVED***
	BasicHandle
	binaryEncodingType
	noElemSeparators
	// EncZeroValuesAsNil says to encode zero values for numbers, bool, string, etc as nil
	EncZeroValuesAsNil bool

	_ [1]uint64 // padding
***REMOVED***

// Name returns the name of the handle: simple
func (h *SimpleHandle) Name() string ***REMOVED*** return "simple" ***REMOVED***

// SetBytesExt sets an extension
func (h *SimpleHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, &extWrapper***REMOVED***ext, interfaceExtFailer***REMOVED******REMOVED******REMOVED***)
***REMOVED***

func (h *SimpleHandle) hasElemSeparators() bool ***REMOVED*** return true ***REMOVED*** // as it implements Write(Map|Array)XXX

func (h *SimpleHandle) newEncDriver(e *Encoder) encDriver ***REMOVED***
	return &simpleEncDriver***REMOVED***e: e, w: e.w, h: h***REMOVED***
***REMOVED***

func (h *SimpleHandle) newDecDriver(d *Decoder) decDriver ***REMOVED***
	return &simpleDecDriver***REMOVED***d: d, h: h, r: d.r, br: d.bytes***REMOVED***
***REMOVED***

func (e *simpleEncDriver) reset() ***REMOVED***
	e.c = 0
	e.w = e.e.w
***REMOVED***

func (d *simpleDecDriver) reset() ***REMOVED***
	d.c = 0
	d.r, d.br = d.d.r, d.d.bytes
	d.bd, d.bdRead = 0, false
***REMOVED***

var _ decDriver = (*simpleDecDriver)(nil)
var _ encDriver = (*simpleEncDriver)(nil)
