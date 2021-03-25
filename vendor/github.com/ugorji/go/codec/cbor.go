// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"math"
	"time"
)

// major
const (
	cborMajorUint byte = iota
	cborMajorNegInt
	cborMajorBytes
	cborMajorString
	cborMajorArray
	cborMajorMap
	cborMajorTag
	cborMajorSimpleOrFloat
)

// simple
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

// indefinite
const (
	cborBdIndefiniteBytes  byte = 0x5f
	cborBdIndefiniteString byte = 0x7f
	cborBdIndefiniteArray  byte = 0x9f
	cborBdIndefiniteMap    byte = 0xbf
	cborBdBreak            byte = 0xff
)

// These define some in-stream descriptors for
// manual encoding e.g. when doing explicit indefinite-length
const (
	CborStreamBytes  byte = 0x5f
	CborStreamString byte = 0x7f
	CborStreamArray  byte = 0x9f
	CborStreamMap    byte = 0xbf
	CborStreamBreak  byte = 0xff
)

// base values
const (
	cborBaseUint   byte = 0x00
	cborBaseNegInt byte = 0x20
	cborBaseBytes  byte = 0x40
	cborBaseString byte = 0x60
	cborBaseArray  byte = 0x80
	cborBaseMap    byte = 0xa0
	cborBaseTag    byte = 0xc0
	cborBaseSimple byte = 0xe0
)

// const (
// 	cborSelfDesrTag  byte = 0xd9
// 	cborSelfDesrTag2 byte = 0xd9
// 	cborSelfDesrTag3 byte = 0xf7
// )

func cbordesc(bd byte) string ***REMOVED***
	switch bd >> 5 ***REMOVED***
	case cborMajorUint:
		return "(u)int"
	case cborMajorNegInt:
		return "int"
	case cborMajorBytes:
		return "bytes"
	case cborMajorString:
		return "string"
	case cborMajorArray:
		return "array"
	case cborMajorMap:
		return "map"
	case cborMajorTag:
		return "tag"
	case cborMajorSimpleOrFloat: // default
		switch bd ***REMOVED***
		case cborBdNil:
			return "nil"
		case cborBdFalse:
			return "false"
		case cborBdTrue:
			return "true"
		case cborBdFloat16, cborBdFloat32, cborBdFloat64:
			return "float"
		case cborBdIndefiniteBytes:
			return "bytes*"
		case cborBdIndefiniteString:
			return "string*"
		case cborBdIndefiniteArray:
			return "array*"
		case cborBdIndefiniteMap:
			return "map*"
		default:
			return "unknown(simple)"
		***REMOVED***
	***REMOVED***
	return "unknown"
***REMOVED***

// -------------------

type cborEncDriver struct ***REMOVED***
	noBuiltInTypes
	encDriverNoopContainerWriter
	h *CborHandle
	x [8]byte
	_ [6]uint64 // padding
	e Encoder
***REMOVED***

func (e *cborEncDriver) encoder() *Encoder ***REMOVED***
	return &e.e
***REMOVED***

func (e *cborEncDriver) EncodeNil() ***REMOVED***
	e.e.encWr.writen1(cborBdNil)
***REMOVED***

func (e *cborEncDriver) EncodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.e.encWr.writen1(cborBdTrue)
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(cborBdFalse)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeFloat32(f float32) ***REMOVED***
	e.e.encWr.writen1(cborBdFloat32)
	bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *cborEncDriver) EncodeFloat64(f float64) ***REMOVED***
	e.e.encWr.writen1(cborBdFloat64)
	bigenHelper***REMOVED***e.x[:8], e.e.w()***REMOVED***.writeUint64(math.Float64bits(f))
***REMOVED***

func (e *cborEncDriver) encUint(v uint64, bd byte) ***REMOVED***
	if v <= 0x17 ***REMOVED***
		e.e.encWr.writen1(byte(v) + bd)
	***REMOVED*** else if v <= math.MaxUint8 ***REMOVED***
		e.e.encWr.writen2(bd+0x18, uint8(v))
	***REMOVED*** else if v <= math.MaxUint16 ***REMOVED***
		e.e.encWr.writen1(bd + 0x19)
		bigenHelper***REMOVED***e.x[:2], e.e.w()***REMOVED***.writeUint16(uint16(v))
	***REMOVED*** else if v <= math.MaxUint32 ***REMOVED***
		e.e.encWr.writen1(bd + 0x1a)
		bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(uint32(v))
	***REMOVED*** else ***REMOVED*** // if v <= math.MaxUint64 ***REMOVED***
		e.e.encWr.writen1(bd + 0x1b)
		bigenHelper***REMOVED***e.x[:8], e.e.w()***REMOVED***.writeUint64(v)
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
		e.encStringBytesS(cborBaseString, t.Format(time.RFC3339Nano))
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

func (e *cborEncDriver) EncodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) ***REMOVED***
	e.encUint(uint64(xtag), cborBaseTag)
	if ext == SelfExt ***REMOVED***
		rv2 := baseRV(rv)
		e.e.encodeValue(rv2, e.h.fnNoExt(rv2.Type()))
	***REMOVED*** else if v := ext.ConvertExt(rv); v == nil ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		e.e.encode(v)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeRawExt(re *RawExt) ***REMOVED***
	e.encUint(uint64(re.Tag), cborBaseTag)
	// only encodes re.Value (never re.Data)
	if re.Value != nil ***REMOVED***
		e.e.encode(re.Value)
	***REMOVED*** else ***REMOVED***
		e.EncodeNil()
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) WriteArrayStart(length int) ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		e.e.encWr.writen1(cborBdIndefiniteArray)
	***REMOVED*** else ***REMOVED***
		e.encLen(cborBaseArray, length)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) WriteMapStart(length int) ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		e.e.encWr.writen1(cborBdIndefiniteMap)
	***REMOVED*** else ***REMOVED***
		e.encLen(cborBaseMap, length)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) WriteMapEnd() ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		e.e.encWr.writen1(cborBdBreak)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) WriteArrayEnd() ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		e.e.encWr.writen1(cborBdBreak)
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) EncodeString(v string) ***REMOVED***
	if e.h.StringToRaw ***REMOVED***
		e.EncodeStringBytesRaw(bytesView(v))
		return
	***REMOVED***
	e.encStringBytesS(cborBaseString, v)
***REMOVED***

func (e *cborEncDriver) EncodeStringBytesRaw(v []byte) ***REMOVED***
	if v == nil ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		e.encStringBytesS(cborBaseBytes, stringView(v))
	***REMOVED***
***REMOVED***

func (e *cborEncDriver) encStringBytesS(bb byte, v string) ***REMOVED***
	if e.h.IndefiniteLength ***REMOVED***
		if bb == cborBaseBytes ***REMOVED***
			e.e.encWr.writen1(cborBdIndefiniteBytes)
		***REMOVED*** else ***REMOVED***
			e.e.encWr.writen1(cborBdIndefiniteString)
		***REMOVED***
		var vlen uint = uint(len(v))
		blen := vlen / 4
		if blen == 0 ***REMOVED***
			blen = 64
		***REMOVED*** else if blen > 1024 ***REMOVED***
			blen = 1024
		***REMOVED***
		for i := uint(0); i < vlen; ***REMOVED***
			var v2 string
			i2 := i + blen
			if i2 >= i && i2 < vlen ***REMOVED***
				v2 = v[i:i2]
			***REMOVED*** else ***REMOVED***
				v2 = v[i:]
			***REMOVED***
			e.encLen(bb, len(v2))
			e.e.encWr.writestr(v2)
			i = i2
		***REMOVED***
		e.e.encWr.writen1(cborBdBreak)
	***REMOVED*** else ***REMOVED***
		e.encLen(bb, len(v))
		e.e.encWr.writestr(v)
	***REMOVED***
***REMOVED***

// ----------------------

type cborDecDriver struct ***REMOVED***
	decDriverNoopContainerReader
	h      *CborHandle
	bdRead bool
	bd     byte
	st     bool // skip tags
	fnil   bool // found nil
	noBuiltInTypes
	_ [6]uint64 // padding cache-aligned
	d Decoder
***REMOVED***

func (d *cborDecDriver) decoder() *Decoder ***REMOVED***
	return &d.d
***REMOVED***

func (d *cborDecDriver) readNextBd() ***REMOVED***
	d.bd = d.d.decRd.readn1()
	d.bdRead = true
***REMOVED***

func (d *cborDecDriver) advanceNil() (null bool) ***REMOVED***
	d.fnil = false
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == cborBdNil || d.bd == cborBdUndefined ***REMOVED***
		d.bdRead = false
		d.fnil = true
		null = true
	***REMOVED***
	return
***REMOVED***

// skipTags is called to skip any tags in the stream.
//
// Since any value can be tagged, then we should call skipTags
// before any value is decoded.
//
// By definition, skipTags should not be called before
// checking for break, or nil or undefined.
func (d *cborDecDriver) skipTags() ***REMOVED***
	for d.bd>>5 == cborMajorTag ***REMOVED***
		d.decUint()
		d.bd = d.d.decRd.readn1()
	***REMOVED***
***REMOVED***

func (d *cborDecDriver) uncacheRead() ***REMOVED***
	if d.bdRead ***REMOVED***
		d.d.decRd.unreadn1()
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *cborDecDriver) ContainerType() (vt valueType) ***REMOVED***
	d.fnil = false
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.st ***REMOVED***
		d.skipTags()
	***REMOVED***
	if d.bd == cborBdNil ***REMOVED***
		d.bdRead = false // always consume nil after seeing it in container type
		d.fnil = true
		return valueTypeNil
	***REMOVED*** else if d.bd == cborBdIndefiniteBytes || (d.bd>>5 == cborMajorBytes) ***REMOVED***
		return valueTypeBytes
	***REMOVED*** else if d.bd == cborBdIndefiniteString || (d.bd>>5 == cborMajorString) ***REMOVED***
		return valueTypeString
	***REMOVED*** else if d.bd == cborBdIndefiniteArray || (d.bd>>5 == cborMajorArray) ***REMOVED***
		return valueTypeArray
	***REMOVED*** else if d.bd == cborBdIndefiniteMap || (d.bd>>5 == cborMajorMap) ***REMOVED***
		return valueTypeMap
	***REMOVED***
	return valueTypeUnset
***REMOVED***

func (d *cborDecDriver) Nil() bool ***REMOVED***
	return d.fnil
***REMOVED***

func (d *cborDecDriver) TryNil() bool ***REMOVED***
	return d.advanceNil()
***REMOVED***

func (d *cborDecDriver) CheckBreak() (v bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == cborBdBreak ***REMOVED***
		d.bdRead = false
		v = true
	***REMOVED***
	return
***REMOVED***

func (d *cborDecDriver) decUint() (ui uint64) ***REMOVED***
	v := d.bd & 0x1f
	if v <= 0x17 ***REMOVED***
		ui = uint64(v)
	***REMOVED*** else ***REMOVED***
		if v == 0x18 ***REMOVED***
			ui = uint64(d.d.decRd.readn1())
		***REMOVED*** else if v == 0x19 ***REMOVED***
			ui = uint64(bigen.Uint16(d.d.decRd.readx(2)))
		***REMOVED*** else if v == 0x1a ***REMOVED***
			ui = uint64(bigen.Uint32(d.d.decRd.readx(4)))
		***REMOVED*** else if v == 0x1b ***REMOVED***
			ui = uint64(bigen.Uint64(d.d.decRd.readx(8)))
		***REMOVED*** else ***REMOVED***
			d.d.errorf("invalid descriptor decoding uint: %x/%s", d.bd, cbordesc(d.bd))
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (d *cborDecDriver) decCheckInteger() (neg bool) ***REMOVED***
	if d.st ***REMOVED***
		d.skipTags()
	***REMOVED***
	major := d.bd >> 5
	if major == cborMajorUint ***REMOVED***
	***REMOVED*** else if major == cborMajorNegInt ***REMOVED***
		neg = true
	***REMOVED*** else ***REMOVED***
		d.d.errorf("invalid integer; got major %v from descriptor %x/%s, expected %v or %v",
			major, d.bd, cbordesc(d.bd), cborMajorUint, cborMajorNegInt)
	***REMOVED***
	return
***REMOVED***

func cborDecInt64(ui uint64, neg bool) (i int64) ***REMOVED***
	// check if this number can be converted to an int without overflow
	if neg ***REMOVED***
		i = -(chkOvf.SignedIntV(ui + 1))
	***REMOVED*** else ***REMOVED***
		i = chkOvf.SignedIntV(ui)
	***REMOVED***
	return
***REMOVED***

func (d *cborDecDriver) decLen() int ***REMOVED***
	return int(d.decUint())
***REMOVED***

func (d *cborDecDriver) decAppendIndefiniteBytes(bs []byte) []byte ***REMOVED***
	d.bdRead = false
	for !d.CheckBreak() ***REMOVED***
		if major := d.bd >> 5; major != cborMajorBytes && major != cborMajorString ***REMOVED***
			d.d.errorf("invalid indefinite string/bytes; got major %v, expected %x/%s",
				major, d.bd, cbordesc(d.bd))
		***REMOVED***
		n := uint(d.decLen())
		oldLen := uint(len(bs))
		newLen := oldLen + n
		if newLen > uint(cap(bs)) ***REMOVED***
			bs2 := make([]byte, newLen, 2*uint(cap(bs))+n)
			copy(bs2, bs)
			bs = bs2
		***REMOVED*** else ***REMOVED***
			bs = bs[:newLen]
		***REMOVED***
		d.d.decRd.readb(bs[oldLen:newLen])
		// bs = append(bs, d.d.decRd.readn()...)
		d.bdRead = false
	***REMOVED***
	d.bdRead = false
	return bs
***REMOVED***

func (d *cborDecDriver) DecodeInt64() (i int64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	neg := d.decCheckInteger()
	ui := d.decUint()
	d.bdRead = false
	return cborDecInt64(ui, neg)
***REMOVED***

func (d *cborDecDriver) DecodeUint64() (ui uint64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.decCheckInteger() ***REMOVED***
		d.d.errorf("cannot assign negative signed value to unsigned type")
	***REMOVED***
	ui = d.decUint()
	d.bdRead = false
	return
***REMOVED***

func (d *cborDecDriver) DecodeFloat64() (f float64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.st ***REMOVED***
		d.skipTags()
	***REMOVED***
	switch d.bd ***REMOVED***
	case cborBdFloat16:
		f = float64(math.Float32frombits(halfFloatToFloatBits(bigen.Uint16(d.d.decRd.readx(2)))))
	case cborBdFloat32:
		f = float64(math.Float32frombits(bigen.Uint32(d.d.decRd.readx(4))))
	case cborBdFloat64:
		f = math.Float64frombits(bigen.Uint64(d.d.decRd.readx(8)))
	default:
		major := d.bd >> 5
		if major == cborMajorUint ***REMOVED***
			f = float64(cborDecInt64(d.decUint(), false))
		***REMOVED*** else if major == cborMajorNegInt ***REMOVED***
			f = float64(cborDecInt64(d.decUint(), true))
		***REMOVED*** else ***REMOVED***
			d.d.errorf("invalid float descriptor; got %d/%s, expected float16/32/64 or (-)int",
				d.bd, cbordesc(d.bd))
		***REMOVED***
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool only (single byte).
func (d *cborDecDriver) DecodeBool() (b bool) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.st ***REMOVED***
		d.skipTags()
	***REMOVED***
	if d.bd == cborBdTrue ***REMOVED***
		b = true
	***REMOVED*** else if d.bd == cborBdFalse ***REMOVED***
	***REMOVED*** else ***REMOVED***
		d.d.errorf("not bool - %s %x/%s", msgBadDesc, d.bd, cbordesc(d.bd))
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *cborDecDriver) ReadMapStart() (length int) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return decContainerLenNil
	***REMOVED***
	if d.st ***REMOVED***
		d.skipTags()
	***REMOVED***
	d.bdRead = false
	if d.bd == cborBdIndefiniteMap ***REMOVED***
		return decContainerLenUnknown
	***REMOVED***
	if d.bd>>5 != cborMajorMap ***REMOVED***
		d.d.errorf("error reading map; got major type: %x, expected %x/%s",
			d.bd>>5, cborMajorMap, cbordesc(d.bd))
	***REMOVED***
	return d.decLen()
***REMOVED***

func (d *cborDecDriver) ReadArrayStart() (length int) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return decContainerLenNil
	***REMOVED***
	if d.st ***REMOVED***
		d.skipTags()
	***REMOVED***
	d.bdRead = false
	if d.bd == cborBdIndefiniteArray ***REMOVED***
		return decContainerLenUnknown
	***REMOVED***
	if d.bd>>5 != cborMajorArray ***REMOVED***
		d.d.errorf("invalid array; got major type: %x, expect: %x/%s",
			d.bd>>5, cborMajorArray, cbordesc(d.bd))
	***REMOVED***
	return d.decLen()
***REMOVED***

func (d *cborDecDriver) DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.st ***REMOVED***
		d.skipTags()
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
	if d.bd == cborBdIndefiniteArray ***REMOVED***
		d.bdRead = false
		if zerocopy && len(bs) == 0 ***REMOVED***
			bs = d.d.b[:]
		***REMOVED***
		if bs == nil ***REMOVED***
			bs = []byte***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			bs = bs[:0]
		***REMOVED***
		for !d.CheckBreak() ***REMOVED***
			bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		***REMOVED***
		return bs
	***REMOVED***
	if d.bd>>5 == cborMajorArray ***REMOVED***
		d.bdRead = false
		if zerocopy && len(bs) == 0 ***REMOVED***
			bs = d.d.b[:]
		***REMOVED***
		slen := d.decLen()
		bs = usableByteSlice(bs, slen)
		for i := 0; i < len(bs); i++ ***REMOVED***
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		***REMOVED***
		return bs
	***REMOVED***
	clen := d.decLen()
	d.bdRead = false
	if zerocopy ***REMOVED***
		if d.d.bytes ***REMOVED***
			return d.d.decRd.readx(uint(clen))
		***REMOVED*** else if len(bs) == 0 ***REMOVED***
			bs = d.d.b[:]
		***REMOVED***
	***REMOVED***
	return decByteSlice(d.d.r(), clen, d.h.MaxInitLen, bs)
***REMOVED***

func (d *cborDecDriver) DecodeStringAsBytes() (s []byte) ***REMOVED***
	return d.DecodeBytes(d.d.b[:], true)
***REMOVED***

func (d *cborDecDriver) DecodeTime() (t time.Time) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.bd>>5 != cborMajorTag ***REMOVED***
		d.d.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	***REMOVED***
	xtag := d.decUint()
	d.bdRead = false
	return d.decodeTime(xtag)
***REMOVED***

func (d *cborDecDriver) decodeTime(xtag uint64) (t time.Time) ***REMOVED***
	switch xtag ***REMOVED***
	case 0:
		var err error
		if t, err = time.Parse(time.RFC3339, stringView(d.DecodeStringAsBytes())); err != nil ***REMOVED***
			d.d.errorv(err)
		***REMOVED***
	case 1:
		f1, f2 := math.Modf(d.DecodeFloat64())
		t = time.Unix(int64(f1), int64(f2*1e9))
	default:
		d.d.errorf("invalid tag for time.Time - expecting 0 or 1, got 0x%x", xtag)
	***REMOVED***
	t = t.UTC().Round(time.Microsecond)
	return
***REMOVED***

func (d *cborDecDriver) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.bd>>5 != cborMajorTag ***REMOVED***
		d.d.errorf("error reading tag; expected major type: %x, got: %x", cborMajorTag, d.bd>>5)
	***REMOVED***
	realxtag := d.decUint()
	d.bdRead = false
	if ext == nil ***REMOVED***
		re := rv.(*RawExt)
		re.Tag = realxtag
		d.d.decode(&re.Value)
	***REMOVED*** else if xtag != realxtag ***REMOVED***
		d.d.errorf("Wrong extension tag. Got %b. Expecting: %v", realxtag, xtag)
		return
	***REMOVED*** else if ext == SelfExt ***REMOVED***
		rv2 := baseRV(rv)
		d.d.decodeValue(rv2, d.h.fnNoExt(rv2.Type()))
	***REMOVED*** else ***REMOVED***
		d.d.interfaceExtConvertAndDecode(rv, ext)
	***REMOVED***
	d.bdRead = false
***REMOVED***

func (d *cborDecDriver) DecodeNaked() ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***

	d.fnil = false
	n := d.d.naked()
	var decodeFurther bool

	switch d.bd >> 5 ***REMOVED***
	case cborMajorUint:
		if d.h.SignedInteger ***REMOVED***
			n.v = valueTypeInt
			n.i = d.DecodeInt64()
		***REMOVED*** else ***REMOVED***
			n.v = valueTypeUint
			n.u = d.DecodeUint64()
		***REMOVED***
	case cborMajorNegInt:
		n.v = valueTypeInt
		n.i = d.DecodeInt64()
	case cborMajorBytes:
		decNakedReadRawBytes(d, &d.d, n, d.h.RawToString)
	case cborMajorString:
		n.v = valueTypeString
		n.s = string(d.DecodeStringAsBytes())
	case cborMajorArray:
		n.v = valueTypeArray
		decodeFurther = true
	case cborMajorMap:
		n.v = valueTypeMap
		decodeFurther = true
	case cborMajorTag:
		n.v = valueTypeExt
		n.u = d.decUint()
		n.l = nil
		if n.u == 0 || n.u == 1 ***REMOVED***
			d.bdRead = false
			n.v = valueTypeTime
			n.t = d.decodeTime(n.u)
		***REMOVED*** else if d.st && d.h.getExtForTag(n.u) == nil ***REMOVED***
			// d.skipTags() // no need to call this - tags already skipped
			d.bdRead = false
			d.DecodeNaked()
			return // return when done (as true recursive function)
		***REMOVED***
	case cborMajorSimpleOrFloat:
		switch d.bd ***REMOVED***
		case cborBdNil, cborBdUndefined:
			n.v = valueTypeNil
			d.fnil = true
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
			decNakedReadRawBytes(d, &d.d, n, d.h.RawToString)
		case cborBdIndefiniteString:
			n.v = valueTypeString
			n.s = string(d.DecodeStringAsBytes())
		case cborBdIndefiniteArray:
			n.v = valueTypeArray
			decodeFurther = true
		case cborBdIndefiniteMap:
			n.v = valueTypeMap
			decodeFurther = true
		default:
			d.d.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
		***REMOVED***
	default: // should never happen
		d.d.errorf("decodeNaked: Unrecognized d.bd: 0x%x", d.bd)
	***REMOVED***
	if !decodeFurther ***REMOVED***
		d.bdRead = false
	***REMOVED***
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
	// noElemSeparators
	BasicHandle

	// IndefiniteLength=true, means that we encode using indefinitelength
	IndefiniteLength bool

	// TimeRFC3339 says to encode time.Time using RFC3339 format.
	// If unset, we encode time.Time using seconds past epoch.
	TimeRFC3339 bool

	// SkipUnexpectedTags says to skip over any tags for which extensions are
	// not defined. This is in keeping with the cbor spec on "Optional Tagging of Items".
	//
	// Furthermore, this allows the skipping over of the Self Describing Tag 0xd9d9f7.
	SkipUnexpectedTags bool

	_ [7]uint64 // padding (cache-aligned)
***REMOVED***

// Name returns the name of the handle: cbor
func (h *CborHandle) Name() string ***REMOVED*** return "cbor" ***REMOVED***

func (h *CborHandle) newEncDriver() encDriver ***REMOVED***
	var e = &cborEncDriver***REMOVED***h: h***REMOVED***
	e.e.e = e
	e.e.init(h)
	e.reset()
	return e
***REMOVED***

func (h *CborHandle) newDecDriver() decDriver ***REMOVED***
	d := &cborDecDriver***REMOVED***h: h, st: h.SkipUnexpectedTags***REMOVED***
	d.d.d = d
	d.d.cbor = true
	d.d.init(h)
	d.reset()
	return d
***REMOVED***

func (e *cborEncDriver) reset() ***REMOVED***
***REMOVED***

func (d *cborDecDriver) reset() ***REMOVED***
	d.bd = 0
	d.bdRead = false
	d.fnil = false
	d.st = d.h.SkipUnexpectedTags
***REMOVED***

var _ decDriver = (*cborDecDriver)(nil)
var _ encDriver = (*cborEncDriver)(nil)
