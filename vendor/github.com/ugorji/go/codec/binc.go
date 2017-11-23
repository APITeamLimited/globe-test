// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"math"
	"reflect"
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

type bincEncDriver struct ***REMOVED***
	e *Encoder
	w encWriter
	m map[string]uint16 // symbols
	b [scratchByteArrayLen]byte
	s uint16 // symbols sequencer
	encNoSeparator
***REMOVED***

func (e *bincEncDriver) IsBuiltinType(rt uintptr) bool ***REMOVED***
	return rt == timeTypId
***REMOVED***

func (e *bincEncDriver) EncodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED***
	if rt == timeTypId ***REMOVED***
		var bs []byte
		switch x := v.(type) ***REMOVED***
		case time.Time:
			bs = encodeTime(x)
		case *time.Time:
			bs = encodeTime(*x)
		default:
			e.e.errorf("binc error encoding builtin: expect time.Time, received %T", v)
		***REMOVED***
		e.w.writen1(bincVdTimestamp<<4 | uint8(len(bs)))
		e.w.writeb(bs)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeNil() ***REMOVED***
	e.w.writen1(bincVdSpecial<<4 | bincSpNil)
***REMOVED***

func (e *bincEncDriver) EncodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpTrue)
	***REMOVED*** else ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpFalse)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeFloat32(f float32) ***REMOVED***
	if f == 0 ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
		return
	***REMOVED***
	e.w.writen1(bincVdFloat<<4 | bincFlBin32)
	bigenHelper***REMOVED***e.b[:4], e.w***REMOVED***.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *bincEncDriver) EncodeFloat64(f float64) ***REMOVED***
	if f == 0 ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpZeroFloat)
		return
	***REMOVED***
	bigen.PutUint64(e.b[:8], math.Float64bits(f))
	if bincDoPrune ***REMOVED***
		i := 7
		for ; i >= 0 && (e.b[i] == 0); i-- ***REMOVED***
		***REMOVED***
		i++
		if i <= 6 ***REMOVED***
			e.w.writen1(bincVdFloat<<4 | 0x8 | bincFlBin64)
			e.w.writen1(byte(i))
			e.w.writeb(e.b[:i])
			return
		***REMOVED***
	***REMOVED***
	e.w.writen1(bincVdFloat<<4 | bincFlBin64)
	e.w.writeb(e.b[:8])
***REMOVED***

func (e *bincEncDriver) encIntegerPrune(bd byte, pos bool, v uint64, lim uint8) ***REMOVED***
	if lim == 4 ***REMOVED***
		bigen.PutUint32(e.b[:lim], uint32(v))
	***REMOVED*** else ***REMOVED***
		bigen.PutUint64(e.b[:lim], v)
	***REMOVED***
	if bincDoPrune ***REMOVED***
		i := pruneSignExt(e.b[:lim], pos)
		e.w.writen1(bd | lim - 1 - byte(i))
		e.w.writeb(e.b[i:lim])
	***REMOVED*** else ***REMOVED***
		e.w.writen1(bd | lim - 1)
		e.w.writeb(e.b[:lim])
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeInt(v int64) ***REMOVED***
	const nbd byte = bincVdNegInt << 4
	if v >= 0 ***REMOVED***
		e.encUint(bincVdPosInt<<4, true, uint64(v))
	***REMOVED*** else if v == -1 ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpNegOne)
	***REMOVED*** else ***REMOVED***
		e.encUint(bincVdNegInt<<4, false, uint64(-v))
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeUint(v uint64) ***REMOVED***
	e.encUint(bincVdPosInt<<4, true, v)
***REMOVED***

func (e *bincEncDriver) encUint(bd byte, pos bool, v uint64) ***REMOVED***
	if v == 0 ***REMOVED***
		e.w.writen1(bincVdSpecial<<4 | bincSpZero)
	***REMOVED*** else if pos && v >= 1 && v <= 16 ***REMOVED***
		e.w.writen1(bincVdSmallInt<<4 | byte(v-1))
	***REMOVED*** else if v <= math.MaxUint8 ***REMOVED***
		e.w.writen2(bd|0x0, byte(v))
	***REMOVED*** else if v <= math.MaxUint16 ***REMOVED***
		e.w.writen1(bd | 0x01)
		bigenHelper***REMOVED***e.b[:2], e.w***REMOVED***.writeUint16(uint16(v))
	***REMOVED*** else if v <= math.MaxUint32 ***REMOVED***
		e.encIntegerPrune(bd, pos, v, 4)
	***REMOVED*** else ***REMOVED***
		e.encIntegerPrune(bd, pos, v, 8)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext, _ *Encoder) ***REMOVED***
	bs := ext.WriteExt(rv)
	if bs == nil ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	e.encodeExtPreamble(uint8(xtag), len(bs))
	e.w.writeb(bs)
***REMOVED***

func (e *bincEncDriver) EncodeRawExt(re *RawExt, _ *Encoder) ***REMOVED***
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
***REMOVED***

func (e *bincEncDriver) encodeExtPreamble(xtag byte, length int) ***REMOVED***
	e.encLen(bincVdCustomExt<<4, uint64(length))
	e.w.writen1(xtag)
***REMOVED***

func (e *bincEncDriver) EncodeArrayStart(length int) ***REMOVED***
	e.encLen(bincVdArray<<4, uint64(length))
***REMOVED***

func (e *bincEncDriver) EncodeMapStart(length int) ***REMOVED***
	e.encLen(bincVdMap<<4, uint64(length))
***REMOVED***

func (e *bincEncDriver) EncodeString(c charEncoding, v string) ***REMOVED***
	l := uint64(len(v))
	e.encBytesLen(c, l)
	if l > 0 ***REMOVED***
		e.w.writestr(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeSymbol(v string) ***REMOVED***
	// if WriteSymbolsNoRefs ***REMOVED***
	// 	e.encodeString(c_UTF8, v)
	// 	return
	// ***REMOVED***

	//symbols only offer benefit when string length > 1.
	//This is because strings with length 1 take only 2 bytes to store
	//(bd with embedded length, and single byte for string val).

	l := len(v)
	if l == 0 ***REMOVED***
		e.encBytesLen(c_UTF8, 0)
		return
	***REMOVED*** else if l == 1 ***REMOVED***
		e.encBytesLen(c_UTF8, 1)
		e.w.writen1(v[0])
		return
	***REMOVED***
	if e.m == nil ***REMOVED***
		e.m = make(map[string]uint16, 16)
	***REMOVED***
	ui, ok := e.m[v]
	if ok ***REMOVED***
		if ui <= math.MaxUint8 ***REMOVED***
			e.w.writen2(bincVdSymbol<<4, byte(ui))
		***REMOVED*** else ***REMOVED***
			e.w.writen1(bincVdSymbol<<4 | 0x8)
			bigenHelper***REMOVED***e.b[:2], e.w***REMOVED***.writeUint16(ui)
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
			e.w.writen2(bincVdSymbol<<4|0x0|0x4|lenprec, byte(ui))
		***REMOVED*** else ***REMOVED***
			e.w.writen1(bincVdSymbol<<4 | 0x8 | 0x4 | lenprec)
			bigenHelper***REMOVED***e.b[:2], e.w***REMOVED***.writeUint16(ui)
		***REMOVED***
		if lenprec == 0 ***REMOVED***
			e.w.writen1(byte(l))
		***REMOVED*** else if lenprec == 1 ***REMOVED***
			bigenHelper***REMOVED***e.b[:2], e.w***REMOVED***.writeUint16(uint16(l))
		***REMOVED*** else if lenprec == 2 ***REMOVED***
			bigenHelper***REMOVED***e.b[:4], e.w***REMOVED***.writeUint32(uint32(l))
		***REMOVED*** else ***REMOVED***
			bigenHelper***REMOVED***e.b[:8], e.w***REMOVED***.writeUint64(uint64(l))
		***REMOVED***
		e.w.writestr(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) EncodeStringBytes(c charEncoding, v []byte) ***REMOVED***
	l := uint64(len(v))
	e.encBytesLen(c, l)
	if l > 0 ***REMOVED***
		e.w.writeb(v)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encBytesLen(c charEncoding, length uint64) ***REMOVED***
	//TODO: support bincUnicodeOther (for now, just use string or bytearray)
	if c == c_RAW ***REMOVED***
		e.encLen(bincVdByteArray<<4, length)
	***REMOVED*** else ***REMOVED***
		e.encLen(bincVdString<<4, length)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encLen(bd byte, l uint64) ***REMOVED***
	if l < 12 ***REMOVED***
		e.w.writen1(bd | uint8(l+4))
	***REMOVED*** else ***REMOVED***
		e.encLenNumber(bd, l)
	***REMOVED***
***REMOVED***

func (e *bincEncDriver) encLenNumber(bd byte, v uint64) ***REMOVED***
	if v <= math.MaxUint8 ***REMOVED***
		e.w.writen2(bd, byte(v))
	***REMOVED*** else if v <= math.MaxUint16 ***REMOVED***
		e.w.writen1(bd | 0x01)
		bigenHelper***REMOVED***e.b[:2], e.w***REMOVED***.writeUint16(uint16(v))
	***REMOVED*** else if v <= math.MaxUint32 ***REMOVED***
		e.w.writen1(bd | 0x02)
		bigenHelper***REMOVED***e.b[:4], e.w***REMOVED***.writeUint32(uint32(v))
	***REMOVED*** else ***REMOVED***
		e.w.writen1(bd | 0x03)
		bigenHelper***REMOVED***e.b[:8], e.w***REMOVED***.writeUint64(uint64(v))
	***REMOVED***
***REMOVED***

//------------------------------------

type bincDecSymbol struct ***REMOVED***
	s string
	b []byte
	i uint16
***REMOVED***

type bincDecDriver struct ***REMOVED***
	d      *Decoder
	h      *BincHandle
	r      decReader
	br     bool // bytes reader
	bdRead bool
	bd     byte
	vd     byte
	vs     byte
	noStreamingCodec
	decNoSeparator
	b [scratchByteArrayLen]byte

	// linear searching on this slice is ok,
	// because we typically expect < 32 symbols in each stream.
	s []bincDecSymbol
***REMOVED***

func (d *bincDecDriver) readNextBd() ***REMOVED***
	d.bd = d.r.readn1()
	d.vd = d.bd >> 4
	d.vs = d.bd & 0x0f
	d.bdRead = true
***REMOVED***

func (d *bincDecDriver) uncacheRead() ***REMOVED***
	if d.bdRead ***REMOVED***
		d.r.unreadn1()
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) ContainerType() (vt valueType) ***REMOVED***
	if d.vd == bincVdSpecial && d.vs == bincSpNil ***REMOVED***
		return valueTypeNil
	***REMOVED*** else if d.vd == bincVdByteArray ***REMOVED***
		return valueTypeBytes
	***REMOVED*** else if d.vd == bincVdString ***REMOVED***
		return valueTypeString
	***REMOVED*** else if d.vd == bincVdArray ***REMOVED***
		return valueTypeArray
	***REMOVED*** else if d.vd == bincVdMap ***REMOVED***
		return valueTypeMap
	***REMOVED*** else ***REMOVED***
		// d.d.errorf("isContainerType: unsupported parameter: %v", vt)
	***REMOVED***
	return valueTypeUnset
***REMOVED***

func (d *bincDecDriver) TryDecodeAsNil() bool ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == bincVdSpecial<<4|bincSpNil ***REMOVED***
		d.bdRead = false
		return true
	***REMOVED***
	return false
***REMOVED***

func (d *bincDecDriver) IsBuiltinType(rt uintptr) bool ***REMOVED***
	return rt == timeTypId
***REMOVED***

func (d *bincDecDriver) DecodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if rt == timeTypId ***REMOVED***
		if d.vd != bincVdTimestamp ***REMOVED***
			d.d.errorf("Invalid d.vd. Expecting 0x%x. Received: 0x%x", bincVdTimestamp, d.vd)
			return
		***REMOVED***
		tt, err := decodeTime(d.r.readx(int(d.vs)))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		var vt *time.Time = v.(*time.Time)
		*vt = tt
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) decFloatPre(vs, defaultLen byte) ***REMOVED***
	if vs&0x8 == 0 ***REMOVED***
		d.r.readb(d.b[0:defaultLen])
	***REMOVED*** else ***REMOVED***
		l := d.r.readn1()
		if l > 8 ***REMOVED***
			d.d.errorf("At most 8 bytes used to represent float. Received: %v bytes", l)
			return
		***REMOVED***
		for i := l; i < 8; i++ ***REMOVED***
			d.b[i] = 0
		***REMOVED***
		d.r.readb(d.b[0:l])
	***REMOVED***
***REMOVED***

func (d *bincDecDriver) decFloat() (f float64) ***REMOVED***
	//if true ***REMOVED*** f = math.Float64frombits(bigen.Uint64(d.r.readx(8))); break; ***REMOVED***
	if x := d.vs & 0x7; x == bincFlBin32 ***REMOVED***
		d.decFloatPre(d.vs, 4)
		f = float64(math.Float32frombits(bigen.Uint32(d.b[0:4])))
	***REMOVED*** else if x == bincFlBin64 ***REMOVED***
		d.decFloatPre(d.vs, 8)
		f = math.Float64frombits(bigen.Uint64(d.b[0:8]))
	***REMOVED*** else ***REMOVED***
		d.d.errorf("only float32 and float64 are supported. d.vd: 0x%x, d.vs: 0x%x", d.vd, d.vs)
		return
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decUint() (v uint64) ***REMOVED***
	// need to inline the code (interface conversion and type assertion expensive)
	switch d.vs ***REMOVED***
	case 0:
		v = uint64(d.r.readn1())
	case 1:
		d.r.readb(d.b[6:8])
		v = uint64(bigen.Uint16(d.b[6:8]))
	case 2:
		d.b[4] = 0
		d.r.readb(d.b[5:8])
		v = uint64(bigen.Uint32(d.b[4:8]))
	case 3:
		d.r.readb(d.b[4:8])
		v = uint64(bigen.Uint32(d.b[4:8]))
	case 4, 5, 6:
		lim := int(7 - d.vs)
		d.r.readb(d.b[lim:8])
		for i := 0; i < lim; i++ ***REMOVED***
			d.b[i] = 0
		***REMOVED***
		v = uint64(bigen.Uint64(d.b[:8]))
	case 7:
		d.r.readb(d.b[:8])
		v = uint64(bigen.Uint64(d.b[:8]))
	default:
		d.d.errorf("unsigned integers with greater than 64 bits of precision not supported")
		return
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decCheckInteger() (ui uint64, neg bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
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
			d.d.errorf("numeric decode fails for special value: d.vs: 0x%x", d.vs)
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		d.d.errorf("number can only be decoded from uint or int values. d.bd: 0x%x, d.vd: 0x%x", d.bd, d.vd)
		return
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) DecodeInt(bitsize uint8) (i int64) ***REMOVED***
	ui, neg := d.decCheckInteger()
	i, overflow := chkOvf.SignedInt(ui)
	if overflow ***REMOVED***
		d.d.errorf("simple: overflow converting %v to signed integer", ui)
		return
	***REMOVED***
	if neg ***REMOVED***
		i = -i
	***REMOVED***
	if chkOvf.Int(i, bitsize) ***REMOVED***
		d.d.errorf("binc: overflow integer: %v", i)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) DecodeUint(bitsize uint8) (ui uint64) ***REMOVED***
	ui, neg := d.decCheckInteger()
	if neg ***REMOVED***
		d.d.errorf("Assigning negative signed value to unsigned type")
		return
	***REMOVED***
	if chkOvf.Uint(ui, bitsize) ***REMOVED***
		d.d.errorf("binc: overflow integer: %v", ui)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) DecodeFloat(chkOverflow32 bool) (f float64) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
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
			d.d.errorf("Invalid d.vs decoding float where d.vd=bincVdSpecial: %v", d.vs)
			return
		***REMOVED***
	***REMOVED*** else if vd == bincVdFloat ***REMOVED***
		f = d.decFloat()
	***REMOVED*** else ***REMOVED***
		f = float64(d.DecodeInt(64))
	***REMOVED***
	if chkOverflow32 && chkOvf.Float32(f) ***REMOVED***
		d.d.errorf("binc: float32 overflow: %v", f)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool only (single byte).
func (d *bincDecDriver) DecodeBool() (b bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if bd := d.bd; bd == (bincVdSpecial | bincSpFalse) ***REMOVED***
		// b = false
	***REMOVED*** else if bd == (bincVdSpecial | bincSpTrue) ***REMOVED***
		b = true
	***REMOVED*** else ***REMOVED***
		d.d.errorf("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) ReadMapStart() (length int) ***REMOVED***
	if d.vd != bincVdMap ***REMOVED***
		d.d.errorf("Invalid d.vd for map. Expecting 0x%x. Got: 0x%x", bincVdMap, d.vd)
		return
	***REMOVED***
	length = d.decLen()
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) ReadArrayStart() (length int) ***REMOVED***
	if d.vd != bincVdArray ***REMOVED***
		d.d.errorf("Invalid d.vd for array. Expecting 0x%x. Got: 0x%x", bincVdArray, d.vd)
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
		v = uint64(d.r.readn1())
	***REMOVED*** else if x == 1 ***REMOVED***
		d.r.readb(d.b[6:8])
		v = uint64(bigen.Uint16(d.b[6:8]))
	***REMOVED*** else if x == 2 ***REMOVED***
		d.r.readb(d.b[4:8])
		v = uint64(bigen.Uint32(d.b[4:8]))
	***REMOVED*** else ***REMOVED***
		d.r.readb(d.b[:8])
		v = bigen.Uint64(d.b[:8])
	***REMOVED***
	return
***REMOVED***

func (d *bincDecDriver) decStringAndBytes(bs []byte, withString, zerocopy bool) (bs2 []byte, s string) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == bincVdSpecial<<4|bincSpNil ***REMOVED***
		d.bdRead = false
		return
	***REMOVED***
	var slen int = -1
	// var ok bool
	switch d.vd ***REMOVED***
	case bincVdString, bincVdByteArray:
		slen = d.decLen()
		if zerocopy ***REMOVED***
			if d.br ***REMOVED***
				bs2 = d.r.readx(slen)
			***REMOVED*** else if len(bs) == 0 ***REMOVED***
				bs2 = decByteSlice(d.r, slen, d.d.h.MaxInitLen, d.b[:])
			***REMOVED*** else ***REMOVED***
				bs2 = decByteSlice(d.r, slen, d.d.h.MaxInitLen, bs)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			bs2 = decByteSlice(d.r, slen, d.d.h.MaxInitLen, bs)
		***REMOVED***
		if withString ***REMOVED***
			s = string(bs2)
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
			symbol = uint16(d.r.readn1())
		***REMOVED*** else ***REMOVED***
			symbol = uint16(bigen.Uint16(d.r.readx(2)))
		***REMOVED***
		if d.s == nil ***REMOVED***
			d.s = make([]bincDecSymbol, 0, 16)
		***REMOVED***

		if vs&0x4 == 0 ***REMOVED***
			for i := range d.s ***REMOVED***
				j := &d.s[i]
				if j.i == symbol ***REMOVED***
					bs2 = j.b
					if withString ***REMOVED***
						if j.s == "" && bs2 != nil ***REMOVED***
							j.s = string(bs2)
						***REMOVED***
						s = j.s
					***REMOVED***
					break
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			switch vs & 0x3 ***REMOVED***
			case 0:
				slen = int(d.r.readn1())
			case 1:
				slen = int(bigen.Uint16(d.r.readx(2)))
			case 2:
				slen = int(bigen.Uint32(d.r.readx(4)))
			case 3:
				slen = int(bigen.Uint64(d.r.readx(8)))
			***REMOVED***
			// since using symbols, do not store any part of
			// the parameter bs in the map, as it might be a shared buffer.
			// bs2 = decByteSlice(d.r, slen, bs)
			bs2 = decByteSlice(d.r, slen, d.d.h.MaxInitLen, nil)
			if withString ***REMOVED***
				s = string(bs2)
			***REMOVED***
			d.s = append(d.s, bincDecSymbol***REMOVED***i: symbol, s: s, b: bs2***REMOVED***)
		***REMOVED***
	default:
		d.d.errorf("Invalid d.vd. Expecting string:0x%x, bytearray:0x%x or symbol: 0x%x. Got: 0x%x",
			bincVdString, bincVdByteArray, bincVdSymbol, d.vd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) DecodeString() (s string) ***REMOVED***
	// DecodeBytes does not accommodate symbols, whose impl stores string version in map.
	// Use decStringAndBytes directly.
	// return string(d.DecodeBytes(d.b[:], true, true))
	_, s = d.decStringAndBytes(d.b[:], true, true)
	return
***REMOVED***

func (d *bincDecDriver) DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte) ***REMOVED***
	if isstring ***REMOVED***
		bsOut, _ = d.decStringAndBytes(bs, false, zerocopy)
		return
	***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == bincVdSpecial<<4|bincSpNil ***REMOVED***
		d.bdRead = false
		return nil
	***REMOVED***
	var clen int
	if d.vd == bincVdString || d.vd == bincVdByteArray ***REMOVED***
		clen = d.decLen()
	***REMOVED*** else ***REMOVED***
		d.d.errorf("Invalid d.vd for bytes. Expecting string:0x%x or bytearray:0x%x. Got: 0x%x",
			bincVdString, bincVdByteArray, d.vd)
		return
	***REMOVED***
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

func (d *bincDecDriver) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) (realxtag uint64) ***REMOVED***
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

func (d *bincDecDriver) decodeExtV(verifyTag bool, tag byte) (xtag byte, xbs []byte) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.vd == bincVdCustomExt ***REMOVED***
		l := d.decLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag ***REMOVED***
			d.d.errorf("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
			return
		***REMOVED***
		xbs = d.r.readx(l)
	***REMOVED*** else if d.vd == bincVdByteArray ***REMOVED***
		xbs = d.DecodeBytes(nil, false, true)
	***REMOVED*** else ***REMOVED***
		d.d.errorf("Invalid d.vd for extensions (Expecting extensions or byte array). Got: 0x%x", d.vd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *bincDecDriver) DecodeNaked() ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***

	n := &d.d.n
	var decodeFurther bool

	switch d.vd ***REMOVED***
	case bincVdSpecial:
		switch d.vs ***REMOVED***
		case bincSpNil:
			n.v = valueTypeNil
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
			d.d.errorf("decodeNaked: Unrecognized special value 0x%x", d.vs)
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
		n.s = d.DecodeString()
	case bincVdString:
		n.v = valueTypeString
		n.s = d.DecodeString()
	case bincVdByteArray:
		n.v = valueTypeBytes
		n.l = d.DecodeBytes(nil, false, false)
	case bincVdTimestamp:
		n.v = valueTypeTimestamp
		tt, err := decodeTime(d.r.readx(int(d.vs)))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		n.t = tt
	case bincVdCustomExt:
		n.v = valueTypeExt
		l := d.decLen()
		n.u = uint64(d.r.readn1())
		n.l = d.r.readx(l)
	case bincVdArray:
		n.v = valueTypeArray
		decodeFurther = true
	case bincVdMap:
		n.v = valueTypeMap
		decodeFurther = true
	default:
		d.d.errorf("decodeNaked: Unrecognized d.vd: 0x%x", d.vd)
	***REMOVED***

	if !decodeFurther ***REMOVED***
		d.bdRead = false
	***REMOVED***
	if n.v == valueTypeUint && d.h.SignedInteger ***REMOVED***
		n.v = valueTypeInt
		n.i = int64(n.u)
	***REMOVED***
	return
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
***REMOVED***

func (h *BincHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, &setExtWrapper***REMOVED***b: ext***REMOVED***)
***REMOVED***

func (h *BincHandle) newEncDriver(e *Encoder) encDriver ***REMOVED***
	return &bincEncDriver***REMOVED***e: e, w: e.w***REMOVED***
***REMOVED***

func (h *BincHandle) newDecDriver(d *Decoder) decDriver ***REMOVED***
	return &bincDecDriver***REMOVED***d: d, h: h, r: d.r, br: d.bytes***REMOVED***
***REMOVED***

func (e *bincEncDriver) reset() ***REMOVED***
	e.w = e.e.w
	e.s = 0
	e.m = nil
***REMOVED***

func (d *bincDecDriver) reset() ***REMOVED***
	d.r, d.br = d.d.r, d.d.bytes
	d.s = nil
	d.bd, d.bdRead, d.vd, d.vs = 0, false, 0, 0
***REMOVED***

var _ decDriver = (*bincDecDriver)(nil)
var _ encDriver = (*bincEncDriver)(nil)
