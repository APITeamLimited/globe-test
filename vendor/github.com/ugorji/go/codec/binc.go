// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"math"
	"time"
)

const bincDoPrune = true // No longer needed. Needed before as C lib did not support pruning.

// vd as low 4 bits (there are 16 slots)
const (
	bincVdSpecial byte = iota
	bincVdPosInt
	bincVdNegInt
	bincVdFloat

	bincVdString
	bincVdByteArray
	bincVdArray
	bincVdMap

	bincVdTimestamp
	bincVdSmallInt
	bincVdUnicodeOther
	bincVdSymbol

	bincVdDecimal
	_               // open slot
	_               // open slot
	bincVdCustomExt = 0x0f
)

const (
	bincSpNil byte = iota
	bincSpFalse
	bincSpTrue
	bincSpNan
	bincSpPosInf
	bincSpNegInf
	bincSpZeroFloat
	bincSpZero
	bincSpNegOne
)

const (
	bincFlBin16 byte = iota
	bincFlBin32
	_ // bincFlBin32e
	bincFlBin64
	_ // bincFlBin64e
	// others not currently supported
)

func bincdesc(vd, vs byte) string ***REMOVED***
	switch vd ***REMOVED***
	case bincVdSpecial:
		switch vs ***REMOVED***
		case bincSpNil:
			return "nil"
		case bincSpFalse:
			return "false"
		case bincSpTrue:
			return "true"
		case bincSpNan, bincSpPosInf, bincSpNegInf, bincSpZeroFloat:
			return "float"
		case bincSpZero:
			return "uint"
		case bincSpNegOne:
			return "int"
		default:
			return "unknown"
		***REMOVED***
	case bincVdSmallInt, bincVdPosInt:
		return "uint"
	case bincVdNegInt:
		return "int"
	case bincVdFloat:
		return "float"
	case bincVdSymbol:
		return "string"
	case bincVdString:
		return "string"
	case bincVdByteArray:
		return "bytes"
	case bincVdTimestamp:
		return "time"
	case bincVdCustomExt:
		return "ext"
	case bincVdArray:
		return "array"
	case bincVdMap:
		return "map"
	default:
		return "unknown"
	***REMOVED***
***REMOVED***

type bincEncDriver struct ***REMOVED***
	noBuiltInTypes
	encDriverNoopContainerWriter
	h *BincHandle
	m map[string]uint16 // symbols
	b [8]byte           // scratch, used for encoding numbers - bigendian style
	s uint16            // symbols sequencer
	_ [4]uint64         // padding
	e Encoder
***REMOVED***

func (e *bincEncDriver) encoder() *Encoder ***REMOVED***
	return &e.e
***REMOVED***

func (e *bincEncDriver) EncodeNil() ***REMOVED***
	e.e.encWr.writen1(bincVdSpecial<<4 | bincSpNil)
***REMOVED***

func (e *bincEncDriver) EncodeTime(t time.Time) ***REMOVED***
	if t.IsZero() ***REMOVED***
		e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		bs := bincEncodeTime(t)
		e.e.encWr.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.e.encWr.writeb(bs)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.e.encWr.writen1(bincVdSpecial<<4 | bincSpTrue)
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(bincVdSpecial<<4 | bincSpFalse)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeFloat32(f float32) ***REMOVED***
	if f == 0 ***REMOVED***
		e.e.encWr.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
		return
	***REMOVED***
	e.e.encWr.writen1(bincVdFloat<<4 | bincFlBin32)
	bigenHelper***REMOVED***e.b[:4], e.e.w()***REMOVED***.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *bincEncDriver) EncodeFloat64(f float64) ***REMOVED***
	if f == 0 ***REMOVED***
		e.e.encWr.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
		return
	***REMOVED***
	bigen.PutUint64(e.b[:8], math.Float64bits(f))
	if bincDoPrune ***REMOVED***
		i := 7
		for ; i >= 0 && (e.b[i] == 0); i-- ***REMOVED***
		***REMOVED***
		i++
		if i <= 6 ***REMOVED***
			e.e.encWr.writen1(bincVdFloat<<4 | 0x8 | bincFlBin64)
			e.e.encWr.writen1(byte(i))
			e.e.encWr.writeb(e.b[:i])
			return
		***REMOVED***
	***REMOVED***
	e.e.encWr.writen1(bincVdFloat<<4 | bincFlBin64)
	e.e.encWr.writeb(e.b[:8])
***REMOVED***

func (e *bincEncDriver) encIntegerPrune(bd byte, pos bool, v uint64, lim uint8) ***REMOVED***
	if lim == 4 ***REMOVED***
		bigen.PutUint32(e.b[:lim], uint32(v))
	***REMOVED*** else ***REMOVED***
		bigen.PutUint64(e.b[:lim], v)
	***REMOVED***
	if bincDoPrune ***REMOVED***
		i := pruneSignExt(e.b[:lim], pos)
		e.e.encWr.writen1(bd | lim - 1 - byte(i))
		e.e.encWr.writeb(e.b[i:lim])
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(bd | lim - 1)
		e.e.encWr.writeb(e.b[:lim])
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeInt(v int64) ***REMOVED***
	// const nbd byte = bincVdNegInt << 4
	if v >= 0 ***REMOVED***
		e.encUint(bincVdPosInt<<4, true, uint64(v))
	***REMOVED*** else if v == -1 ***REMOVED***
		e.e.encWr.writen1(bincVdSpecial<<4 | bincSpNegOne)
	***REMOVED*** else ***REMOVED***
		e.encUint(bincVdNegInt<<4, false, uint64(-v))
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeUint(v uint64) ***REMOVED***
	e.encUint(bincVdPosInt<<4, true, v)
***REMOVED***

func (e *bincEncDriver) encUint(bd byte, pos bool, v uint64) ***REMOVED***
	if v == 0 ***REMOVED***
		e.e.encWr.writen1(bincVdSpecial<<4 | bincSpZero)
	***REMOVED*** else if pos && v >= 1 && v <= 16 ***REMOVED***
		e.e.encWr.writen1(bincVdSmallInt<<4 | byte(v-1))
	***REMOVED*** else if v <= math.MaxUint8 ***REMOVED***
		e.e.encWr.writen2(bd|0x0, byte(v))
	***REMOVED*** else if v <= math.MaxUint16 ***REMOVED***
		e.e.encWr.writen1(bd | 0x01)
		bigenHelper***REMOVED***e.b[:2], e.e.w()***REMOVED***.writeUint16(uint16(v))
	***REMOVED*** else if v <= math.MaxUint32 ***REMOVED***
		e.encIntegerPrune(bd, pos, v, 4)
	***REMOVED*** else ***REMOVED***
		e.encIntegerPrune(bd, pos, v, 8)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext) ***REMOVED***
	var bs []byte
	if ext == SelfExt ***REMOVED***
		bs = e.e.blist.get(1024)[:0]
		e.e.sideEncode(v, &bs)
	***REMOVED*** else ***REMOVED***
		bs = ext.WriteExt(v)
	***REMOVED***
	if bs == nil ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	e.encodeExtPreamble(uint8(xtag), len(bs))
	e.e.encWr.writeb(bs)
	if ext == SelfExt ***REMOVED***
		e.e.blist.put(bs)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeRawExt(re *RawExt) ***REMOVED***
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.e.encWr.writeb(re.Data)
***REMOVED***

func (e *bincEncDriver) encodeExtPreamble(xtag byte, length int) ***REMOVED***
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.e.encWr.writen1(xtag)
***REMOVED***

func (e *bincEncDriver) WriteArrayStart(length int) ***REMOVED***
	e.encLen(bincVdArray<<4, uint64(length))
***REMOVED***

func (e *bincEncDriver) WriteMapStart(length int) ***REMOVED***
	e.encLen(bincVdMap<<4, uint64(length))
***REMOVED***

func (e *bincEncDriver) EncodeSymbol(v string) ***REMOVED***
	// if WriteSymbolsNoRefs ***REMOVED***
	// 	e.encodeString(cUTF8, v)
	// 	return
	// ***REMOVED***

	//symbols only offer benefit when string length > 1.
	//This is because strings with length 1 take only 2 bytes to store
	//(bd with embedded length, and single byte for string val).

	l := len(v)
	if l == 0 ***REMOVED***
		e.encBytesLen(cUTF8, 0)
		return
	***REMOVED*** else if l == 1 ***REMOVED***
		e.encBytesLen(cUTF8, 1)
		e.e.encWr.writen1(v[0])
		return
	***REMOVED***
	if e.m == nil ***REMOVED***
		e.m = make(map[string]uint16, 16)
	***REMOVED***
	ui, ok := e.m[v]
	if ok ***REMOVED***
		if ui <= math.MaxUint8 ***REMOVED***
			e.e.encWr.writen2(bincVdSymbol<<4, byte(ui))
		***REMOVED*** else ***REMOVED***
			e.e.encWr.writen1(bincVdSymbol<<4 | 0x8)
			bigenHelper***REMOVED***e.b[:2], e.e.w()***REMOVED***.writeUint16(ui)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		e.s++
		ui = e.s
		//ui = uint16(atomic.AddUint32(&e.s, 1))
		e.m[v] = ui
		var lenprec uint8
		if l <= math.MaxUint8 ***REMOVED***
			// lenprec = 0
		***REMOVED*** else if l <= math.MaxUint16 ***REMOVED***
			lenprec = 1
		***REMOVED*** else if int64(l) <= math.MaxUint32 ***REMOVED***
			lenprec = 2
		***REMOVED*** else ***REMOVED***
			lenprec = 3
		***REMOVED***
		if ui <= math.MaxUint8 ***REMOVED***
			e.e.encWr.writen2(bincVdSymbol<<4|0x0|0x4|lenprec, byte(ui))
		***REMOVED*** else ***REMOVED***
			e.e.encWr.writen1(bincVdSymbol<<4 | 0x8 | 0x4 | lenprec)
			bigenHelper***REMOVED***e.b[:2], e.e.w()***REMOVED***.writeUint16(ui)
		***REMOVED***
		if lenprec == 0 ***REMOVED***
			e.e.encWr.writen1(byte(l))
		***REMOVED*** else if lenprec == 1 ***REMOVED***
			bigenHelper***REMOVED***e.b[:2], e.e.w()***REMOVED***.writeUint16(uint16(l))
		***REMOVED*** else if lenprec == 2 ***REMOVED***
			bigenHelper***REMOVED***e.b[:4], e.e.w()***REMOVED***.writeUint32(uint32(l))
		***REMOVED*** else ***REMOVED***
			bigenHelper***REMOVED***e.b[:8], e.e.w()***REMOVED***.writeUint64(uint64(l))
		***REMOVED***
		e.e.encWr.writestr(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeString(v string) ***REMOVED***
	if e.h.StringToRaw ***REMOVED***
		e.encLen(bincVdByteArray<<4, uint64(len(v))) // e.encBytesLen(c, l)
		if len(v) > 0 ***REMOVED***
			e.e.encWr.writestr(v)
		***REMOVED***
		return
	***REMOVED***
	e.EncodeStringEnc(cUTF8, v)
***REMOVED***

func (e *bincEncDriver) EncodeStringEnc(c charEncoding, v string) ***REMOVED***
	if e.e.c == containerMapKey && c == cUTF8 && (e.h.AsSymbols == 1) ***REMOVED***
		e.EncodeSymbol(v)
		return
	***REMOVED***
	e.encLen(bincVdString<<4, uint64(len(v))) // e.encBytesLen(c, l)
	if len(v) > 0 ***REMOVED***
		e.e.encWr.writestr(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeStringBytesRaw(v []byte) ***REMOVED***
	if v == nil ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	e.encLen(bincVdByteArray<<4, uint64(len(v))) // e.encBytesLen(c, l)
	if len(v) > 0 ***REMOVED***
		e.e.encWr.writeb(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encBytesLen(c charEncoding, length uint64) ***REMOVED***
	// NOTE: we currently only support UTF-8 (string) and RAW (bytearray).
	// We should consider supporting bincUnicodeOther.

	if c == cRAW ***REMOVED***
		e.encLen(bincVdByteArray<<4, length)
	***REMOVED*** else ***REMOVED***
		e.encLen(bincVdString<<4, length)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encLen(bd byte, l uint64) ***REMOVED***
	if l < 12 ***REMOVED***
		e.e.encWr.writen1(bd | uint8(l+4))
	***REMOVED*** else ***REMOVED***
		e.encLenNumber(bd, l)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encLenNumber(bd byte, v uint64) ***REMOVED***
	if v <= math.MaxUint8 ***REMOVED***
		e.e.encWr.writen2(bd, byte(v))
	***REMOVED*** else if v <= math.MaxUint16 ***REMOVED***
		e.e.encWr.writen1(bd | 0x01)
		bigenHelper***REMOVED***e.b[:2], e.e.w()***REMOVED***.writeUint16(uint16(v))
	***REMOVED*** else if v <= math.MaxUint32 ***REMOVED***
		e.e.encWr.writen1(bd | 0x02)
		bigenHelper***REMOVED***e.b[:4], e.e.w()***REMOVED***.writeUint32(uint32(v))
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(bd | 0x03)
		bigenHelper***REMOVED***e.b[:8], e.e.w()***REMOVED***.writeUint64(uint64(v))
	***REMOVED***
***REMOVED***

//------------------------------------

type bincDecDriver struct ***REMOVED***
	decDriverNoopContainerReader
	noBuiltInTypes

	h      *BincHandle
	bdRead bool
	bd     byte
	vd     byte
	vs     byte

	fnil bool
	// _      [3]byte // padding
	// linear searching on this slice is ok,
	// because we typically expect < 32 symbols in each stream.
	s map[uint16][]byte // []bincDecSymbol

	b [8]byte   // scratch for decoding numbers - big endian style
	_ [4]uint64 // padding cache-aligned

	d Decoder
***REMOVED***

func (d *bincDecDriver) decoder() *Decoder ***REMOVED***
	return &d.d
***REMOVED***

func (d *bincDecDriver) readNextBd() ***REMOVED***
	d.bd = d.d.decRd.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
***REMOVED***

func (d *bincDecDriver) uncacheRead() ***REMOVED***
	if d.bdRead ***REMOVED***
		d.d.decRd.unreadn1()
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) advanceNil() (null bool) ***REMOVED***
	d.fnil = false
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == bincVdSpecial<<4|bincSpNil ***REMOVED***
		d.bdRead = false
		d.fnil = true
		null = true
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) Nil() bool ***REMOVED***
	return d.fnil
***REMOVED***

func (d *bincDecDriver) TryNil() bool ***REMOVED***
	return d.advanceNil()
***REMOVED***

func (d *bincDecDriver) ContainerType() (vt valueType) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	d.fnil = false
	// if d.vd == bincVdSpecial && d.vs == bincSpNil ***REMOVED***
	if d.bd == bincVdSpecial<<4|bincSpNil ***REMOVED***
		d.bdRead = false
		d.fnil = true
		return valueTypeNil
	***REMOVED*** else if d.vd == bincVdByteArray ***REMOVED***
		return valueTypeBytes
	***REMOVED*** else if d.vd == bincVdString ***REMOVED***
		return valueTypeString
	***REMOVED*** else if d.vd == bincVdArray ***REMOVED***
		return valueTypeArray
	***REMOVED*** else if d.vd == bincVdMap ***REMOVED***
		return valueTypeMap
	***REMOVED***
	return valueTypeUnset
***REMOVED***

func (d *bincDecDriver) DecodeTime() (t time.Time) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.vd != bincVdTimestamp ***REMOVED***
		d.d.errorf("cannot decode time - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		return
	***REMOVED***
	t, err := bincDecodeTime(d.d.decRd.readx(uint(d.vs)))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) decFloatPre(vs, defaultLen byte) ***REMOVED***
	if vs&0x8 == 0 ***REMOVED***
		d.d.decRd.readb(d.b[0:defaultLen])
	***REMOVED*** else ***REMOVED***
		l := d.d.decRd.readn1()
		if l > 8 ***REMOVED***
			d.d.errorf("cannot read float - at most 8 bytes used to represent float - received %v bytes", l)
			return
		***REMOVED***
		for i := l; i < 8; i++ ***REMOVED***
			d.b[i] = 0
		***REMOVED***
		d.d.decRd.readb(d.b[0:l])
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) decFloat() (f float64) ***REMOVED***
	//if true ***REMOVED*** f = math.Float64frombits(bigen.Uint64(d.d.decRd.readx(8))); break; ***REMOVED***
	if x := d.vs & 0x7; x == bincFlBin32 ***REMOVED***
		d.decFloatPre(d.vs, 4)
		f = float64(math.Float32frombits(bigen.Uint32(d.b[0:4])))
	***REMOVED*** else if x == bincFlBin64 ***REMOVED***
		d.decFloatPre(d.vs, 8)
		f = math.Float64frombits(bigen.Uint64(d.b[0:8]))
	***REMOVED*** else ***REMOVED***
		d.d.errorf("read float - only float32 and float64 are supported - %s %x-%x/%s",
			msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		return
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decUint() (v uint64) ***REMOVED***
	// need to inline the code (interface conversion and type assertion expensive)
	switch d.vs ***REMOVED***
	case 0:
		v = uint64(d.d.decRd.readn1())
	case 1:
		d.d.decRd.readb(d.b[6:8])
		v = uint64(bigen.Uint16(d.b[6:8]))
	case 2:
		d.b[4] = 0
		d.d.decRd.readb(d.b[5:8])
		v = uint64(bigen.Uint32(d.b[4:8]))
	case 3:
		d.d.decRd.readb(d.b[4:8])
		v = uint64(bigen.Uint32(d.b[4:8]))
	case 4, 5, 6:
		lim := 7 - d.vs
		d.d.decRd.readb(d.b[lim:8])
		for i := uint8(0); i < lim; i++ ***REMOVED***
			d.b[i] = 0
		***REMOVED***
		v = uint64(bigen.Uint64(d.b[:8]))
	case 7:
		d.d.decRd.readb(d.b[:8])
		v = uint64(bigen.Uint64(d.b[:8]))
	default:
		d.d.errorf("unsigned integers with greater than 64 bits of precision not supported")
		return
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decCheckInteger() (ui uint64, neg bool) ***REMOVED***
	vd, vs := d.vd, d.vs
	if vd == bincVdPosInt ***REMOVED***
		ui = d.decUint()
	***REMOVED*** else if vd == bincVdNegInt ***REMOVED***
		ui = d.decUint()
		neg = true
	***REMOVED*** else if vd == bincVdSmallInt ***REMOVED***
		ui = uint64(d.vs) + 1
	***REMOVED*** else if vd == bincVdSpecial ***REMOVED***
		if vs == bincSpZero ***REMOVED***
			//i = 0
		***REMOVED*** else if vs == bincSpNegOne ***REMOVED***
			neg = true
			ui = 1
		***REMOVED*** else ***REMOVED***
			d.d.errorf("integer decode fails - invalid special value from descriptor %x-%x/%s",
				d.vd, d.vs, bincdesc(d.vd, d.vs))
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		d.d.errorf("integer can only be decoded from int/uint. d.bd: 0x%x, d.vd: 0x%x", d.bd, d.vd)
		return
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) DecodeInt64() (i int64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	ui, neg := d.decCheckInteger()
	i = chkOvf.SignedIntV(ui)
	if neg ***REMOVED***
		i = -i
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) DecodeUint64() (ui uint64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	ui, neg := d.decCheckInteger()
	if neg ***REMOVED***
		d.d.errorf("assigning negative signed value to unsigned integer type")
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) DecodeFloat64() (f float64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	vd, vs := d.vd, d.vs
	if vd == bincVdSpecial ***REMOVED***
		d.bdRead = false
		if vs == bincSpNan ***REMOVED***
			return math.NaN()
		***REMOVED*** else if vs == bincSpPosInf ***REMOVED***
			return math.Inf(1)
		***REMOVED*** else if vs == bincSpZeroFloat || vs == bincSpZero ***REMOVED***
			return
		***REMOVED*** else if vs == bincSpNegInf ***REMOVED***
			return math.Inf(-1)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("float - invalid special value from descriptor %x-%x/%s",
				d.vd, d.vs, bincdesc(d.vd, d.vs))
			return
		***REMOVED***
	***REMOVED*** else if vd == bincVdFloat ***REMOVED***
		f = d.decFloat()
	***REMOVED*** else ***REMOVED***
		f = float64(d.DecodeInt64())
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool only (single byte).
func (d *bincDecDriver) DecodeBool() (b bool) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.bd == (bincVdSpecial | bincSpFalse) ***REMOVED***
		// b = false
	***REMOVED*** else if d.bd == (bincVdSpecial | bincSpTrue) ***REMOVED***
		b = true
	***REMOVED*** else ***REMOVED***
		d.d.errorf("bool - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) ReadMapStart() (length int) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return decContainerLenNil
	***REMOVED***
	if d.vd != bincVdMap ***REMOVED***
		d.d.errorf("map - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		return
	***REMOVED***
	length = d.decLen()
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) ReadArrayStart() (length int) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return decContainerLenNil
	***REMOVED***
	if d.vd != bincVdArray ***REMOVED***
		d.d.errorf("array - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		return
	***REMOVED***
	length = d.decLen()
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) decLen() int ***REMOVED***
	if d.vs > 3 ***REMOVED***
		return int(d.vs - 4)
	***REMOVED***
	return int(d.decLenNumber())
***REMOVED***

func (d *bincDecDriver) decLenNumber() (v uint64) ***REMOVED***
	if x := d.vs; x == 0 ***REMOVED***
		v = uint64(d.d.decRd.readn1())
	***REMOVED*** else if x == 1 ***REMOVED***
		d.d.decRd.readb(d.b[6:8])
		v = uint64(bigen.Uint16(d.b[6:8]))
	***REMOVED*** else if x == 2 ***REMOVED***
		d.d.decRd.readb(d.b[4:8])
		v = uint64(bigen.Uint32(d.b[4:8]))
	***REMOVED*** else ***REMOVED***
		d.d.decRd.readb(d.b[:8])
		v = bigen.Uint64(d.b[:8])
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decStringBytes(bs []byte, zerocopy bool) (bs2 []byte) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	var slen = -1
	// var ok bool
	switch d.vd ***REMOVED***
	case bincVdString, bincVdByteArray:
		slen = d.decLen()
		if zerocopy ***REMOVED***
			if d.d.bytes ***REMOVED***
				bs2 = d.d.decRd.readx(uint(slen))
			***REMOVED*** else if len(bs) == 0 ***REMOVED***
				bs2 = decByteSlice(d.d.r(), slen, d.d.h.MaxInitLen, d.d.b[:])
			***REMOVED*** else ***REMOVED***
				bs2 = decByteSlice(d.d.r(), slen, d.d.h.MaxInitLen, bs)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			bs2 = decByteSlice(d.d.r(), slen, d.d.h.MaxInitLen, bs)
		***REMOVED***
	case bincVdSymbol:
		// zerocopy doesn't apply for symbols,
		// as the values must be stored in a table for later use.
		//
		//from vs: extract numSymbolBytes, containsStringVal, strLenPrecision,
		//extract symbol
		//if containsStringVal, read it and put in map
		//else look in map for string value
		var symbol uint16
		vs := d.vs
		if vs&0x8 == 0 ***REMOVED***
			symbol = uint16(d.d.decRd.readn1())
		***REMOVED*** else ***REMOVED***
			symbol = uint16(bigen.Uint16(d.d.decRd.readx(2)))
		***REMOVED***
		if d.s == nil ***REMOVED***
			// d.s = pool4mapU16Bytes.Get().(map[uint16][]byte) // make([]bincDecSymbol, 0, 16)
			d.s = make(map[uint16][]byte, 16)
		***REMOVED***

		if vs&0x4 == 0 ***REMOVED***
			bs2 = d.s[symbol]
		***REMOVED*** else ***REMOVED***
			switch vs & 0x3 ***REMOVED***
			case 0:
				slen = int(d.d.decRd.readn1())
			case 1:
				slen = int(bigen.Uint16(d.d.decRd.readx(2)))
			case 2:
				slen = int(bigen.Uint32(d.d.decRd.readx(4)))
			case 3:
				slen = int(bigen.Uint64(d.d.decRd.readx(8)))
			***REMOVED***
			// since using symbols, do not store any part of
			// the parameter bs in the map, as it might be a shared buffer.
			// bs2 = decByteSlice(d.d.r(), slen, bs)
			bs2 = decByteSlice(d.d.r(), slen, d.d.h.MaxInitLen, nil)
			d.s[symbol] = bs2
			// d.s = append(d.s, bincDecSymbol***REMOVED***i: symbol, s: s, b: bs2***REMOVED***)
		***REMOVED***
	default:
		d.d.errorf("string/bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) DecodeStringAsBytes() (s []byte) ***REMOVED***
	return d.decStringBytes(d.d.b[:], true)
***REMOVED***

func (d *bincDecDriver) DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	// check if an "array" of uint8's (see ContainerType for how to infer if an array)
	if d.vd == bincVdArray ***REMOVED***
		if zerocopy && len(bs) == 0 ***REMOVED***
			bs = d.d.b[:]
		***REMOVED***
		// bsOut, _ = fastpathTV.DecSliceUint8V(bs, true, d.d)
		slen := d.ReadArrayStart()
		bs = usableByteSlice(bs, slen)
		for i := 0; i < slen; i++ ***REMOVED***
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		***REMOVED***
		return bs
	***REMOVED***
	var clen int
	if d.vd == bincVdString || d.vd == bincVdByteArray ***REMOVED***
		clen = d.decLen()
	***REMOVED*** else ***REMOVED***
		d.d.errorf("bytes - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		return
	***REMOVED***
	d.bdRead = false
	if zerocopy ***REMOVED***
		if d.d.bytes ***REMOVED***
			return d.d.decRd.readx(uint(clen))
		***REMOVED*** else if len(bs) == 0 ***REMOVED***
			bs = d.d.b[:]
		***REMOVED***
	***REMOVED***
	return decByteSlice(d.d.r(), clen, d.d.h.MaxInitLen, bs)
***REMOVED***

func (d *bincDecDriver) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) ***REMOVED***
	if xtag > 0xff ***REMOVED***
		d.d.errorf("ext: tag must be <= 0xff; got: %v", xtag)
		return
	***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	realxtag1, xbs := d.decodeExtV(ext != nil, uint8(xtag))
	realxtag := uint64(realxtag1)
	if ext == nil ***REMOVED***
		re := rv.(*RawExt)
		re.Tag = realxtag
		re.Data = detachZeroCopyBytes(d.d.bytes, re.Data, xbs)
	***REMOVED*** else if ext == SelfExt ***REMOVED***
		d.d.sideDecode(rv, xbs)
	***REMOVED*** else ***REMOVED***
		ext.ReadExt(rv, xbs)
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) decodeExtV(verifyTag bool, tag byte) (xtag byte, xbs []byte) ***REMOVED***
	if d.vd == bincVdCustomExt ***REMOVED***
		l := d.decLen()
		xtag = d.d.decRd.readn1()
		if verifyTag && xtag != tag ***REMOVED***
			d.d.errorf("wrong extension tag - got %b, expecting: %v", xtag, tag)
			return
		***REMOVED***
		if d.d.bytes ***REMOVED***
			xbs = d.d.decRd.readx(uint(l))
		***REMOVED*** else ***REMOVED***
			xbs = decByteSlice(d.d.r(), l, d.d.h.MaxInitLen, d.d.b[:])
		***REMOVED***
	***REMOVED*** else if d.vd == bincVdByteArray ***REMOVED***
		xbs = d.DecodeBytes(nil, true)
	***REMOVED*** else ***REMOVED***
		d.d.errorf("ext - expecting extensions or byte array - %s %x-%x/%s",
			msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) DecodeNaked() ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***

	d.fnil = false
	n := d.d.naked()
	var decodeFurther bool

	switch d.vd ***REMOVED***
	case bincVdSpecial:
		switch d.vs ***REMOVED***
		case bincSpNil:
			n.v = valueTypeNil
			d.fnil = true
		case bincSpFalse:
			n.v = valueTypeBool
			n.b = false
		case bincSpTrue:
			n.v = valueTypeBool
			n.b = true
		case bincSpNan:
			n.v = valueTypeFloat
			n.f = math.NaN()
		case bincSpPosInf:
			n.v = valueTypeFloat
			n.f = math.Inf(1)
		case bincSpNegInf:
			n.v = valueTypeFloat
			n.f = math.Inf(-1)
		case bincSpZeroFloat:
			n.v = valueTypeFloat
			n.f = float64(0)
		case bincSpZero:
			n.v = valueTypeUint
			n.u = uint64(0) // int8(0)
		case bincSpNegOne:
			n.v = valueTypeInt
			n.i = int64(-1) // int8(-1)
		default:
			d.d.errorf("cannot infer value - unrecognized special value from descriptor %x-%x/%s",
				d.vd, d.vs, bincdesc(d.vd, d.vs))
		***REMOVED***
	case bincVdSmallInt:
		n.v = valueTypeUint
		n.u = uint64(int8(d.vs)) + 1 // int8(d.vs) + 1
	case bincVdPosInt:
		n.v = valueTypeUint
		n.u = d.decUint()
	case bincVdNegInt:
		n.v = valueTypeInt
		n.i = -(int64(d.decUint()))
	case bincVdFloat:
		n.v = valueTypeFloat
		n.f = d.decFloat()
	case bincVdSymbol:
		n.v = valueTypeSymbol
		n.s = string(d.DecodeStringAsBytes())
	case bincVdString:
		n.v = valueTypeString
		n.s = string(d.DecodeStringAsBytes())
	case bincVdByteArray:
		decNakedReadRawBytes(d, &d.d, n, d.h.RawToString)
	case bincVdTimestamp:
		n.v = valueTypeTime
		tt, err := bincDecodeTime(d.d.decRd.readx(uint(d.vs)))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		n.t = tt
	case bincVdCustomExt:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.d.decRd.readn1())
		if d.d.bytes ***REMOVED***
			n.l = d.d.decRd.readx(uint(l))
		***REMOVED*** else ***REMOVED***
			n.l = decByteSlice(d.d.r(), l, d.d.h.MaxInitLen, d.d.b[:])
		***REMOVED***
	case bincVdArray:
		n.v = valueTypeArray
		decodeFurther = true
	case bincVdMap:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		d.d.errorf("cannot infer value - %s %x-%x/%s", msgBadDesc, d.vd, d.vs, bincdesc(d.vd, d.vs))
	***REMOVED***

	if !decodeFurther ***REMOVED***
		d.bdRead = false
	***REMOVED***
	if n.v == valueTypeUint && d.h.SignedInteger ***REMOVED***
		n.v = valueTypeInt
		n.i = int64(n.u)
	***REMOVED***
***REMOVED***

//------------------------------------

//BincHandle is a Handle for the Binc Schema-Free Encoding Format
//defined at https://github.com/ugorji/binc .
//
//BincHandle currently supports all Binc features with the following EXCEPTIONS:
//  - only integers up to 64 bits of precision are supported.
//    big integers are unsupported.
//  - Only IEEE 754 binary32 and binary64 floats are supported (ie Go float32 and float64 types).
//    extended precision and decimal IEEE 754 floats are unsupported.
//  - Only UTF-8 strings supported.
//    Unicode_Other Binc types (UTF16, UTF32) are currently unsupported.
//
//Note that these EXCEPTIONS are temporary and full support is possible and may happen soon.
type BincHandle struct ***REMOVED***
	BasicHandle
	binaryEncodingType
	// noElemSeparators

	// AsSymbols defines what should be encoded as symbols.
	//
	// Encoding as symbols can reduce the encoded size significantly.
	//
	// However, during decoding, each string to be encoded as a symbol must
	// be checked to see if it has been seen before. Consequently, encoding time
	// will increase if using symbols, because string comparisons has a clear cost.
	//
	// Values:
	// - 0: default: library uses best judgement
	// - 1: use symbols
	// - 2: do not use symbols
	AsSymbols uint8

	// AsSymbols: may later on introduce more options ...
	// - m: map keys
	// - s: struct fields
	// - n: none
	// - a: all: same as m, s, ...

	_ [7]uint64 // padding (cache-aligned)
***REMOVED***

// Name returns the name of the handle: binc
func (h *BincHandle) Name() string ***REMOVED*** return "binc" ***REMOVED***

func (h *BincHandle) newEncDriver() encDriver ***REMOVED***
	var e = &bincEncDriver***REMOVED***h: h***REMOVED***
	e.e.e = e
	e.e.init(h)
	e.reset()
	return e
***REMOVED***

func (h *BincHandle) newDecDriver() decDriver ***REMOVED***
	d := &bincDecDriver***REMOVED***h: h***REMOVED***
	d.d.d = d
	d.d.init(h)
	d.reset()
	return d
***REMOVED***

func (e *bincEncDriver) reset() ***REMOVED***
	e.s = 0
	e.m = nil
***REMOVED***

func (e *bincEncDriver) atEndOfEncode() ***REMOVED***
	if e.m != nil ***REMOVED***
		for k := range e.m ***REMOVED***
			delete(e.m, k)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) reset() ***REMOVED***
	d.s = nil
	d.bd, d.bdRead, d.vd, d.vs = 0, false, 0, 0
	d.fnil = false
***REMOVED***

func (d *bincDecDriver) atEndOfDecode() ***REMOVED***
	if d.s != nil ***REMOVED***
		for k := range d.s ***REMOVED***
			delete(d.s, k)
		***REMOVED***
	***REMOVED***
***REMOVED***

// var timeDigits = [...]byte***REMOVED***'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'***REMOVED***

// EncodeTime encodes a time.Time as a []byte, including
// information on the instant in time and UTC offset.
//
// Format Description
//
//   A timestamp is composed of 3 components:
//
//   - secs: signed integer representing seconds since unix epoch
//   - nsces: unsigned integer representing fractional seconds as a
//     nanosecond offset within secs, in the range 0 <= nsecs < 1e9
//   - tz: signed integer representing timezone offset in minutes east of UTC,
//     and a dst (daylight savings time) flag
//
//   When encoding a timestamp, the first byte is the descriptor, which
//   defines which components are encoded and how many bytes are used to
//   encode secs and nsecs components. *If secs/nsecs is 0 or tz is UTC, it
//   is not encoded in the byte array explicitly*.
//
//       Descriptor 8 bits are of the form `A B C DDD EE`:
//           A:   Is secs component encoded? 1 = true
//           B:   Is nsecs component encoded? 1 = true
//           C:   Is tz component encoded? 1 = true
//           DDD: Number of extra bytes for secs (range 0-7).
//                If A = 1, secs encoded in DDD+1 bytes.
//                    If A = 0, secs is not encoded, and is assumed to be 0.
//                    If A = 1, then we need at least 1 byte to encode secs.
//                    DDD says the number of extra bytes beyond that 1.
//                    E.g. if DDD=0, then secs is represented in 1 byte.
//                         if DDD=2, then secs is represented in 3 bytes.
//           EE:  Number of extra bytes for nsecs (range 0-3).
//                If B = 1, nsecs encoded in EE+1 bytes (similar to secs/DDD above)
//
//   Following the descriptor bytes, subsequent bytes are:
//
//       secs component encoded in `DDD + 1` bytes (if A == 1)
//       nsecs component encoded in `EE + 1` bytes (if B == 1)
//       tz component encoded in 2 bytes (if C == 1)
//
//   secs and nsecs components are integers encoded in a BigEndian
//   2-complement encoding format.
//
//   tz component is encoded as 2 bytes (16 bits). Most significant bit 15 to
//   Least significant bit 0 are described below:
//
//       Timezone offset has a range of -12:00 to +14:00 (ie -720 to +840 minutes).
//       Bit 15 = have\_dst: set to 1 if we set the dst flag.
//       Bit 14 = dst\_on: set to 1 if dst is in effect at the time, or 0 if not.
//       Bits 13..0 = timezone offset in minutes. It is a signed integer in Big Endian format.
//
func bincEncodeTime(t time.Time) []byte ***REMOVED***
	// t := rv2i(rv).(time.Time)
	tsecs, tnsecs := t.Unix(), t.Nanosecond()
	var (
		bd   byte
		btmp [8]byte
		bs   [16]byte
		i    int = 1
	)
	l := t.Location()
	if l == time.UTC ***REMOVED***
		l = nil
	***REMOVED***
	if tsecs != 0 ***REMOVED***
		bd = bd | 0x80
		bigen.PutUint64(btmp[:], uint64(tsecs))
		f := pruneSignExt(btmp[:], tsecs >= 0)
		bd = bd | (byte(7-f) << 2)
		copy(bs[i:], btmp[f:])
		i = i + (8 - f)
	***REMOVED***
	if tnsecs != 0 ***REMOVED***
		bd = bd | 0x40
		bigen.PutUint32(btmp[:4], uint32(tnsecs))
		f := pruneSignExt(btmp[:4], true)
		bd = bd | byte(3-f)
		copy(bs[i:], btmp[f:4])
		i = i + (4 - f)
	***REMOVED***
	if l != nil ***REMOVED***
		bd = bd | 0x20
		// Note that Go Libs do not give access to dst flag.
		_, zoneOffset := t.Zone()
		// zoneName, zoneOffset := t.Zone()
		zoneOffset /= 60
		z := uint16(zoneOffset)
		bigen.PutUint16(btmp[:2], z)
		// clear dst flags
		bs[i] = btmp[0] & 0x3f
		bs[i+1] = btmp[1]
		i = i + 2
	***REMOVED***
	bs[0] = bd
	return bs[0:i]
***REMOVED***

// bincDecodeTime decodes a []byte into a time.Time.
func bincDecodeTime(bs []byte) (tt time.Time, err error) ***REMOVED***
	bd := bs[0]
	var (
		tsec  int64
		tnsec uint32
		tz    uint16
		i     byte = 1
		i2    byte
		n     byte
	)
	if bd&(1<<7) != 0 ***REMOVED***
		var btmp [8]byte
		n = ((bd >> 2) & 0x7) + 1
		i2 = i + n
		copy(btmp[8-n:], bs[i:i2])
		// if first bit of bs[i] is set, then fill btmp[0..8-n] with 0xff (ie sign extend it)
		if bs[i]&(1<<7) != 0 ***REMOVED***
			copy(btmp[0:8-n], bsAll0xff)
			// for j,k := byte(0), 8-n; j < k; j++ ***REMOVED***	btmp[j] = 0xff ***REMOVED***
		***REMOVED***
		i = i2
		tsec = int64(bigen.Uint64(btmp[:]))
	***REMOVED***
	if bd&(1<<6) != 0 ***REMOVED***
		var btmp [4]byte
		n = (bd & 0x3) + 1
		i2 = i + n
		copy(btmp[4-n:], bs[i:i2])
		i = i2
		tnsec = bigen.Uint32(btmp[:])
	***REMOVED***
	if bd&(1<<5) == 0 ***REMOVED***
		tt = time.Unix(tsec, int64(tnsec)).UTC()
		return
	***REMOVED***
	// In stdlib time.Parse, when a date is parsed without a zone name, it uses "" as zone name.
	// However, we need name here, so it can be shown when time is printf.d.
	// Zone name is in form: UTC-08:00.
	// Note that Go Libs do not give access to dst flag, so we ignore dst bits

	i2 = i + 2
	tz = bigen.Uint16(bs[i:i2])
	// i = i2
	// sign extend sign bit into top 2 MSB (which were dst bits):
	if tz&(1<<13) == 0 ***REMOVED*** // positive
		tz = tz & 0x3fff //clear 2 MSBs: dst bits
	***REMOVED*** else ***REMOVED*** // negative
		tz = tz | 0xc000 //set 2 MSBs: dst bits
	***REMOVED***
	tzint := int16(tz)
	if tzint == 0 ***REMOVED***
		tt = time.Unix(tsec, int64(tnsec)).UTC()
	***REMOVED*** else ***REMOVED***
		// For Go Time, do not use a descriptive timezone.
		// It's unnecessary, and makes it harder to do a reflect.DeepEqual.
		// The Offset already tells what the offset should be, if not on UTC and unknown zone name.
		// var zoneName = timeLocUTCName(tzint)
		tt = time.Unix(tsec, int64(tnsec)).In(time.FixedZone("", int(tzint)*60))
	***REMOVED***
	return
***REMOVED***

// func timeLocUTCName(tzint int16) string ***REMOVED***
// 	if tzint == 0 ***REMOVED***
// 		return "UTC"
// 	***REMOVED***
// 	var tzname = []byte("UTC+00:00")
// 	//tzname := fmt.Sprintf("UTC%s%02d:%02d", tzsign, tz/60, tz%60) //perf issue using Sprintf.. inline below.
// 	//tzhr, tzmin := tz/60, tz%60 //faster if u convert to int first
// 	var tzhr, tzmin int16
// 	if tzint < 0 ***REMOVED***
// 		tzname[3] = '-'
// 		tzhr, tzmin = -tzint/60, (-tzint)%60
// 	***REMOVED*** else ***REMOVED***
// 		tzhr, tzmin = tzint/60, tzint%60
// 	***REMOVED***
// 	tzname[4] = timeDigits[tzhr/10]
// 	tzname[5] = timeDigits[tzhr%10]
// 	tzname[7] = timeDigits[tzmin/10]
// 	tzname[8] = timeDigits[tzmin%10]
// 	return string(tzname)
// 	//return time.FixedZone(string(tzname), int(tzint)*60)
// ***REMOVED***

var _ decDriver = (*bincDecDriver)(nil)
var _ encDriver = (*bincEncDriver)(nil)
