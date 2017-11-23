// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

/*
MSGPACK

Msgpack-c implementation powers the c, c++, python, ruby, etc libraries.
We need to maintain compatibility with it and how it encodes integer values
without caring about the type.

For compatibility with behaviour of msgpack-c reference implementation:
  - Go intX (>0) and uintX
       IS ENCODED AS
    msgpack +ve fixnum, unsigned
  - Go intX (<0)
       IS ENCODED AS
    msgpack -ve fixnum, signed

*/
package codec

import (
	"fmt"
	"io"
	"math"
	"net/rpc"
	"reflect"
)

const (
	mpPosFixNumMin byte = 0x00
	mpPosFixNumMax      = 0x7f
	mpFixMapMin         = 0x80
	mpFixMapMax         = 0x8f
	mpFixArrayMin       = 0x90
	mpFixArrayMax       = 0x9f
	mpFixStrMin         = 0xa0
	mpFixStrMax         = 0xbf
	mpNil               = 0xc0
	_                   = 0xc1
	mpFalse             = 0xc2
	mpTrue              = 0xc3
	mpFloat             = 0xca
	mpDouble            = 0xcb
	mpUint8             = 0xcc
	mpUint16            = 0xcd
	mpUint32            = 0xce
	mpUint64            = 0xcf
	mpInt8              = 0xd0
	mpInt16             = 0xd1
	mpInt32             = 0xd2
	mpInt64             = 0xd3

	// extensions below
	mpBin8     = 0xc4
	mpBin16    = 0xc5
	mpBin32    = 0xc6
	mpExt8     = 0xc7
	mpExt16    = 0xc8
	mpExt32    = 0xc9
	mpFixExt1  = 0xd4
	mpFixExt2  = 0xd5
	mpFixExt4  = 0xd6
	mpFixExt8  = 0xd7
	mpFixExt16 = 0xd8

	mpStr8  = 0xd9 // new
	mpStr16 = 0xda
	mpStr32 = 0xdb

	mpArray16 = 0xdc
	mpArray32 = 0xdd

	mpMap16 = 0xde
	mpMap32 = 0xdf

	mpNegFixNumMin = 0xe0
	mpNegFixNumMax = 0xff
)

// MsgpackSpecRpcMultiArgs is a special type which signifies to the MsgpackSpecRpcCodec
// that the backend RPC service takes multiple arguments, which have been arranged
// in sequence in the slice.
//
// The Codec then passes it AS-IS to the rpc service (without wrapping it in an
// array of 1 element).
type MsgpackSpecRpcMultiArgs []interface***REMOVED******REMOVED***

// A MsgpackContainer type specifies the different types of msgpackContainers.
type msgpackContainerType struct ***REMOVED***
	fixCutoff                   int
	bFixMin, b8, b16, b32       byte
	hasFixMin, has8, has8Always bool
***REMOVED***

var (
	msgpackContainerStr  = msgpackContainerType***REMOVED***32, mpFixStrMin, mpStr8, mpStr16, mpStr32, true, true, false***REMOVED***
	msgpackContainerBin  = msgpackContainerType***REMOVED***0, 0, mpBin8, mpBin16, mpBin32, false, true, true***REMOVED***
	msgpackContainerList = msgpackContainerType***REMOVED***16, mpFixArrayMin, 0, mpArray16, mpArray32, true, false, false***REMOVED***
	msgpackContainerMap  = msgpackContainerType***REMOVED***16, mpFixMapMin, 0, mpMap16, mpMap32, true, false, false***REMOVED***
)

//---------------------------------------------

type msgpackEncDriver struct ***REMOVED***
	noBuiltInTypes
	encNoSeparator
	e *Encoder
	w encWriter
	h *MsgpackHandle
	x [8]byte
***REMOVED***

func (e *msgpackEncDriver) EncodeNil() ***REMOVED***
	e.w.writen1(mpNil)
***REMOVED***

func (e *msgpackEncDriver) EncodeInt(i int64) ***REMOVED***
	if i >= 0 ***REMOVED***
		e.EncodeUint(uint64(i))
	***REMOVED*** else if i >= -32 ***REMOVED***
		e.w.writen1(byte(i))
	***REMOVED*** else if i >= math.MinInt8 ***REMOVED***
		e.w.writen2(mpInt8, byte(i))
	***REMOVED*** else if i >= math.MinInt16 ***REMOVED***
		e.w.writen1(mpInt16)
		bigenHelper***REMOVED***e.x[:2], e.w***REMOVED***.writeUint16(uint16(i))
	***REMOVED*** else if i >= math.MinInt32 ***REMOVED***
		e.w.writen1(mpInt32)
		bigenHelper***REMOVED***e.x[:4], e.w***REMOVED***.writeUint32(uint32(i))
	***REMOVED*** else ***REMOVED***
		e.w.writen1(mpInt64)
		bigenHelper***REMOVED***e.x[:8], e.w***REMOVED***.writeUint64(uint64(i))
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeUint(i uint64) ***REMOVED***
	if i <= math.MaxInt8 ***REMOVED***
		e.w.writen1(byte(i))
	***REMOVED*** else if i <= math.MaxUint8 ***REMOVED***
		e.w.writen2(mpUint8, byte(i))
	***REMOVED*** else if i <= math.MaxUint16 ***REMOVED***
		e.w.writen1(mpUint16)
		bigenHelper***REMOVED***e.x[:2], e.w***REMOVED***.writeUint16(uint16(i))
	***REMOVED*** else if i <= math.MaxUint32 ***REMOVED***
		e.w.writen1(mpUint32)
		bigenHelper***REMOVED***e.x[:4], e.w***REMOVED***.writeUint32(uint32(i))
	***REMOVED*** else ***REMOVED***
		e.w.writen1(mpUint64)
		bigenHelper***REMOVED***e.x[:8], e.w***REMOVED***.writeUint64(uint64(i))
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.w.writen1(mpTrue)
	***REMOVED*** else ***REMOVED***
		e.w.writen1(mpFalse)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeFloat32(f float32) ***REMOVED***
	e.w.writen1(mpFloat)
	bigenHelper***REMOVED***e.x[:4], e.w***REMOVED***.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *msgpackEncDriver) EncodeFloat64(f float64) ***REMOVED***
	e.w.writen1(mpDouble)
	bigenHelper***REMOVED***e.x[:8], e.w***REMOVED***.writeUint64(math.Float64bits(f))
***REMOVED***

func (e *msgpackEncDriver) EncodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext, _ *Encoder) ***REMOVED***
	bs := ext.WriteExt(v)
	if bs == nil ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	if e.h.WriteExt ***REMOVED***
		e.encodeExtPreamble(uint8(xtag), len(bs))
		e.w.writeb(bs)
	***REMOVED*** else ***REMOVED***
		e.EncodeStringBytes(c_RAW, bs)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeRawExt(re *RawExt, _ *Encoder) ***REMOVED***
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.w.writeb(re.Data)
***REMOVED***

func (e *msgpackEncDriver) encodeExtPreamble(xtag byte, l int) ***REMOVED***
	if l == 1 ***REMOVED***
		e.w.writen2(mpFixExt1, xtag)
	***REMOVED*** else if l == 2 ***REMOVED***
		e.w.writen2(mpFixExt2, xtag)
	***REMOVED*** else if l == 4 ***REMOVED***
		e.w.writen2(mpFixExt4, xtag)
	***REMOVED*** else if l == 8 ***REMOVED***
		e.w.writen2(mpFixExt8, xtag)
	***REMOVED*** else if l == 16 ***REMOVED***
		e.w.writen2(mpFixExt16, xtag)
	***REMOVED*** else if l < 256 ***REMOVED***
		e.w.writen2(mpExt8, byte(l))
		e.w.writen1(xtag)
	***REMOVED*** else if l < 65536 ***REMOVED***
		e.w.writen1(mpExt16)
		bigenHelper***REMOVED***e.x[:2], e.w***REMOVED***.writeUint16(uint16(l))
		e.w.writen1(xtag)
	***REMOVED*** else ***REMOVED***
		e.w.writen1(mpExt32)
		bigenHelper***REMOVED***e.x[:4], e.w***REMOVED***.writeUint32(uint32(l))
		e.w.writen1(xtag)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeArrayStart(length int) ***REMOVED***
	e.writeContainerLen(msgpackContainerList, length)
***REMOVED***

func (e *msgpackEncDriver) EncodeMapStart(length int) ***REMOVED***
	e.writeContainerLen(msgpackContainerMap, length)
***REMOVED***

func (e *msgpackEncDriver) EncodeString(c charEncoding, s string) ***REMOVED***
	if c == c_RAW && e.h.WriteExt ***REMOVED***
		e.writeContainerLen(msgpackContainerBin, len(s))
	***REMOVED*** else ***REMOVED***
		e.writeContainerLen(msgpackContainerStr, len(s))
	***REMOVED***
	if len(s) > 0 ***REMOVED***
		e.w.writestr(s)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeSymbol(v string) ***REMOVED***
	e.EncodeString(c_UTF8, v)
***REMOVED***

func (e *msgpackEncDriver) EncodeStringBytes(c charEncoding, bs []byte) ***REMOVED***
	if c == c_RAW && e.h.WriteExt ***REMOVED***
		e.writeContainerLen(msgpackContainerBin, len(bs))
	***REMOVED*** else ***REMOVED***
		e.writeContainerLen(msgpackContainerStr, len(bs))
	***REMOVED***
	if len(bs) > 0 ***REMOVED***
		e.w.writeb(bs)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) writeContainerLen(ct msgpackContainerType, l int) ***REMOVED***
	if ct.hasFixMin && l < ct.fixCutoff ***REMOVED***
		e.w.writen1(ct.bFixMin | byte(l))
	***REMOVED*** else if ct.has8 && l < 256 && (ct.has8Always || e.h.WriteExt) ***REMOVED***
		e.w.writen2(ct.b8, uint8(l))
	***REMOVED*** else if l < 65536 ***REMOVED***
		e.w.writen1(ct.b16)
		bigenHelper***REMOVED***e.x[:2], e.w***REMOVED***.writeUint16(uint16(l))
	***REMOVED*** else ***REMOVED***
		e.w.writen1(ct.b32)
		bigenHelper***REMOVED***e.x[:4], e.w***REMOVED***.writeUint32(uint32(l))
	***REMOVED***
***REMOVED***

//---------------------------------------------

type msgpackDecDriver struct ***REMOVED***
	d      *Decoder
	r      decReader // *Decoder decReader decReaderT
	h      *MsgpackHandle
	b      [scratchByteArrayLen]byte
	bd     byte
	bdRead bool
	br     bool // bytes reader
	noBuiltInTypes
	noStreamingCodec
	decNoSeparator
***REMOVED***

// Note: This returns either a primitive (int, bool, etc) for non-containers,
// or a containerType, or a specific type denoting nil or extension.
// It is called when a nil interface***REMOVED******REMOVED*** is passed, leaving it up to the DecDriver
// to introspect the stream and decide how best to decode.
// It deciphers the value by looking at the stream first.
func (d *msgpackDecDriver) DecodeNaked() ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	bd := d.bd
	n := &d.d.n
	var decodeFurther bool

	switch bd ***REMOVED***
	case mpNil:
		n.v = valueTypeNil
		d.bdRead = false
	case mpFalse:
		n.v = valueTypeBool
		n.b = false
	case mpTrue:
		n.v = valueTypeBool
		n.b = true

	case mpFloat:
		n.v = valueTypeFloat
		n.f = float64(math.Float32frombits(bigen.Uint32(d.r.readx(4))))
	case mpDouble:
		n.v = valueTypeFloat
		n.f = math.Float64frombits(bigen.Uint64(d.r.readx(8)))

	case mpUint8:
		n.v = valueTypeUint
		n.u = uint64(d.r.readn1())
	case mpUint16:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint16(d.r.readx(2)))
	case mpUint32:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint32(d.r.readx(4)))
	case mpUint64:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint64(d.r.readx(8)))

	case mpInt8:
		n.v = valueTypeInt
		n.i = int64(int8(d.r.readn1()))
	case mpInt16:
		n.v = valueTypeInt
		n.i = int64(int16(bigen.Uint16(d.r.readx(2))))
	case mpInt32:
		n.v = valueTypeInt
		n.i = int64(int32(bigen.Uint32(d.r.readx(4))))
	case mpInt64:
		n.v = valueTypeInt
		n.i = int64(int64(bigen.Uint64(d.r.readx(8))))

	default:
		switch ***REMOVED***
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax:
			// positive fixnum (always signed)
			n.v = valueTypeInt
			n.i = int64(int8(bd))
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:
			// negative fixnum
			n.v = valueTypeInt
			n.i = int64(int8(bd))
		case bd == mpStr8, bd == mpStr16, bd == mpStr32, bd >= mpFixStrMin && bd <= mpFixStrMax:
			if d.h.RawToString ***REMOVED***
				n.v = valueTypeString
				n.s = d.DecodeString()
			***REMOVED*** else ***REMOVED***
				n.v = valueTypeBytes
				n.l = d.DecodeBytes(nil, false, false)
			***REMOVED***
		case bd == mpBin8, bd == mpBin16, bd == mpBin32:
			n.v = valueTypeBytes
			n.l = d.DecodeBytes(nil, false, false)
		case bd == mpArray16, bd == mpArray32, bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			n.v = valueTypeArray
			decodeFurther = true
		case bd == mpMap16, bd == mpMap32, bd >= mpFixMapMin && bd <= mpFixMapMax:
			n.v = valueTypeMap
			decodeFurther = true
		case bd >= mpFixExt1 && bd <= mpFixExt16, bd >= mpExt8 && bd <= mpExt32:
			n.v = valueTypeExt
			clen := d.readExtLen()
			n.u = uint64(d.r.readn1())
			n.l = d.r.readx(clen)
		default:
			d.d.errorf("Nil-Deciphered DecodeValue: %s: hex: %x, dec: %d", msgBadDesc, bd, bd)
		***REMOVED***
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

// int can be decoded from msgpack type: intXXX or uintXXX
func (d *msgpackDecDriver) DecodeInt(bitsize uint8) (i int64) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	switch d.bd ***REMOVED***
	case mpUint8:
		i = int64(uint64(d.r.readn1()))
	case mpUint16:
		i = int64(uint64(bigen.Uint16(d.r.readx(2))))
	case mpUint32:
		i = int64(uint64(bigen.Uint32(d.r.readx(4))))
	case mpUint64:
		i = int64(bigen.Uint64(d.r.readx(8)))
	case mpInt8:
		i = int64(int8(d.r.readn1()))
	case mpInt16:
		i = int64(int16(bigen.Uint16(d.r.readx(2))))
	case mpInt32:
		i = int64(int32(bigen.Uint32(d.r.readx(4))))
	case mpInt64:
		i = int64(bigen.Uint64(d.r.readx(8)))
	default:
		switch ***REMOVED***
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			i = int64(int8(d.bd))
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			i = int64(int8(d.bd))
		default:
			d.d.errorf("Unhandled single-byte unsigned integer value: %s: %x", msgBadDesc, d.bd)
			return
		***REMOVED***
	***REMOVED***
	// check overflow (logic adapted from std pkg reflect/value.go OverflowUint()
	if bitsize > 0 ***REMOVED***
		if trunc := (i << (64 - bitsize)) >> (64 - bitsize); i != trunc ***REMOVED***
			d.d.errorf("Overflow int value: %v", i)
			return
		***REMOVED***
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// uint can be decoded from msgpack type: intXXX or uintXXX
func (d *msgpackDecDriver) DecodeUint(bitsize uint8) (ui uint64) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	switch d.bd ***REMOVED***
	case mpUint8:
		ui = uint64(d.r.readn1())
	case mpUint16:
		ui = uint64(bigen.Uint16(d.r.readx(2)))
	case mpUint32:
		ui = uint64(bigen.Uint32(d.r.readx(4)))
	case mpUint64:
		ui = bigen.Uint64(d.r.readx(8))
	case mpInt8:
		if i := int64(int8(d.r.readn1())); i >= 0 ***REMOVED***
			ui = uint64(i)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("Assigning negative signed value: %v, to unsigned type", i)
			return
		***REMOVED***
	case mpInt16:
		if i := int64(int16(bigen.Uint16(d.r.readx(2)))); i >= 0 ***REMOVED***
			ui = uint64(i)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("Assigning negative signed value: %v, to unsigned type", i)
			return
		***REMOVED***
	case mpInt32:
		if i := int64(int32(bigen.Uint32(d.r.readx(4)))); i >= 0 ***REMOVED***
			ui = uint64(i)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("Assigning negative signed value: %v, to unsigned type", i)
			return
		***REMOVED***
	case mpInt64:
		if i := int64(bigen.Uint64(d.r.readx(8))); i >= 0 ***REMOVED***
			ui = uint64(i)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("Assigning negative signed value: %v, to unsigned type", i)
			return
		***REMOVED***
	default:
		switch ***REMOVED***
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			ui = uint64(d.bd)
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			d.d.errorf("Assigning negative signed value: %v, to unsigned type", int(d.bd))
			return
		default:
			d.d.errorf("Unhandled single-byte unsigned integer value: %s: %x", msgBadDesc, d.bd)
			return
		***REMOVED***
	***REMOVED***
	// check overflow (logic adapted from std pkg reflect/value.go OverflowUint()
	if bitsize > 0 ***REMOVED***
		if trunc := (ui << (64 - bitsize)) >> (64 - bitsize); ui != trunc ***REMOVED***
			d.d.errorf("Overflow uint value: %v", ui)
			return
		***REMOVED***
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// float can either be decoded from msgpack type: float, double or intX
func (d *msgpackDecDriver) DecodeFloat(chkOverflow32 bool) (f float64) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == mpFloat ***REMOVED***
		f = float64(math.Float32frombits(bigen.Uint32(d.r.readx(4))))
	***REMOVED*** else if d.bd == mpDouble ***REMOVED***
		f = math.Float64frombits(bigen.Uint64(d.r.readx(8)))
	***REMOVED*** else ***REMOVED***
		f = float64(d.DecodeInt(0))
	***REMOVED***
	if chkOverflow32 && chkOvf.Float32(f) ***REMOVED***
		d.d.errorf("msgpack: float32 overflow: %v", f)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool, fixnum 0 or 1.
func (d *msgpackDecDriver) DecodeBool() (b bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == mpFalse || d.bd == 0 ***REMOVED***
		// b = false
	***REMOVED*** else if d.bd == mpTrue || d.bd == 1 ***REMOVED***
		b = true
	***REMOVED*** else ***REMOVED***
		d.d.errorf("Invalid single-byte value for bool: %s: %x", msgBadDesc, d.bd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *msgpackDecDriver) DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	var clen int
	// ignore isstring. Expect that the bytes may be found from msgpackContainerStr or msgpackContainerBin
	if bd := d.bd; bd == mpBin8 || bd == mpBin16 || bd == mpBin32 ***REMOVED***
		clen = d.readContainerLen(msgpackContainerBin)
	***REMOVED*** else ***REMOVED***
		clen = d.readContainerLen(msgpackContainerStr)
	***REMOVED***
	// println("DecodeBytes: clen: ", clen)
	d.bdRead = false
	// bytes may be nil, so handle it. if nil, clen=-1.
	if clen < 0 ***REMOVED***
		return nil
	***REMOVED***
	if zerocopy ***REMOVED***
		if d.br ***REMOVED***
			return d.r.readx(clen)
		***REMOVED*** else if len(bs) == 0 ***REMOVED***
			bs = d.b[:]
		***REMOVED***
	***REMOVED***
	return decByteSlice(d.r, clen, d.d.h.MaxInitLen, bs)
***REMOVED***

func (d *msgpackDecDriver) DecodeString() (s string) ***REMOVED***
	return string(d.DecodeBytes(d.b[:], true, true))
***REMOVED***

func (d *msgpackDecDriver) readNextBd() ***REMOVED***
	d.bd = d.r.readn1()
	d.bdRead = true
***REMOVED***

func (d *msgpackDecDriver) uncacheRead() ***REMOVED***
	if d.bdRead ***REMOVED***
		d.r.unreadn1()
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *msgpackDecDriver) ContainerType() (vt valueType) ***REMOVED***
	bd := d.bd
	if bd == mpNil ***REMOVED***
		return valueTypeNil
	***REMOVED*** else if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 ||
		(!d.h.RawToString &&
			(bd == mpStr8 || bd == mpStr16 || bd == mpStr32 || (bd >= mpFixStrMin && bd <= mpFixStrMax))) ***REMOVED***
		return valueTypeBytes
	***REMOVED*** else if d.h.RawToString &&
		(bd == mpStr8 || bd == mpStr16 || bd == mpStr32 || (bd >= mpFixStrMin && bd <= mpFixStrMax)) ***REMOVED***
		return valueTypeString
	***REMOVED*** else if bd == mpArray16 || bd == mpArray32 || (bd >= mpFixArrayMin && bd <= mpFixArrayMax) ***REMOVED***
		return valueTypeArray
	***REMOVED*** else if bd == mpMap16 || bd == mpMap32 || (bd >= mpFixMapMin && bd <= mpFixMapMax) ***REMOVED***
		return valueTypeMap
	***REMOVED*** else ***REMOVED***
		// d.d.errorf("isContainerType: unsupported parameter: %v", vt)
	***REMOVED***
	return valueTypeUnset
***REMOVED***

func (d *msgpackDecDriver) TryDecodeAsNil() (v bool) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == mpNil ***REMOVED***
		d.bdRead = false
		v = true
	***REMOVED***
	return
***REMOVED***

func (d *msgpackDecDriver) readContainerLen(ct msgpackContainerType) (clen int) ***REMOVED***
	bd := d.bd
	if bd == mpNil ***REMOVED***
		clen = -1 // to represent nil
	***REMOVED*** else if bd == ct.b8 ***REMOVED***
		clen = int(d.r.readn1())
	***REMOVED*** else if bd == ct.b16 ***REMOVED***
		clen = int(bigen.Uint16(d.r.readx(2)))
	***REMOVED*** else if bd == ct.b32 ***REMOVED***
		clen = int(bigen.Uint32(d.r.readx(4)))
	***REMOVED*** else if (ct.bFixMin & bd) == ct.bFixMin ***REMOVED***
		clen = int(ct.bFixMin ^ bd)
	***REMOVED*** else ***REMOVED***
		d.d.errorf("readContainerLen: %s: hex: %x, decimal: %d", msgBadDesc, bd, bd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *msgpackDecDriver) ReadMapStart() int ***REMOVED***
	return d.readContainerLen(msgpackContainerMap)
***REMOVED***

func (d *msgpackDecDriver) ReadArrayStart() int ***REMOVED***
	return d.readContainerLen(msgpackContainerList)
***REMOVED***

func (d *msgpackDecDriver) readExtLen() (clen int) ***REMOVED***
	switch d.bd ***REMOVED***
	case mpNil:
		clen = -1 // to represent nil
	case mpFixExt1:
		clen = 1
	case mpFixExt2:
		clen = 2
	case mpFixExt4:
		clen = 4
	case mpFixExt8:
		clen = 8
	case mpFixExt16:
		clen = 16
	case mpExt8:
		clen = int(d.r.readn1())
	case mpExt16:
		clen = int(bigen.Uint16(d.r.readx(2)))
	case mpExt32:
		clen = int(bigen.Uint32(d.r.readx(4)))
	default:
		d.d.errorf("decoding ext bytes: found unexpected byte: %x", d.bd)
		return
	***REMOVED***
	return
***REMOVED***

func (d *msgpackDecDriver) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) (realxtag uint64) ***REMOVED***
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

func (d *msgpackDecDriver) decodeExtV(verifyTag bool, tag byte) (xtag byte, xbs []byte) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	xbd := d.bd
	if xbd == mpBin8 || xbd == mpBin16 || xbd == mpBin32 ***REMOVED***
		xbs = d.DecodeBytes(nil, false, true)
	***REMOVED*** else if xbd == mpStr8 || xbd == mpStr16 || xbd == mpStr32 ||
		(xbd >= mpFixStrMin && xbd <= mpFixStrMax) ***REMOVED***
		xbs = d.DecodeBytes(nil, true, true)
	***REMOVED*** else ***REMOVED***
		clen := d.readExtLen()
		xtag = d.r.readn1()
		if verifyTag && xtag != tag ***REMOVED***
			d.d.errorf("Wrong extension tag. Got %b. Expecting: %v", xtag, tag)
			return
		***REMOVED***
		xbs = d.r.readx(clen)
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

//--------------------------------------------------

//MsgpackHandle is a Handle for the Msgpack Schema-Free Encoding Format.
type MsgpackHandle struct ***REMOVED***
	BasicHandle

	// RawToString controls how raw bytes are decoded into a nil interface***REMOVED******REMOVED***.
	RawToString bool

	// WriteExt flag supports encoding configured extensions with extension tags.
	// It also controls whether other elements of the new spec are encoded (ie Str8).
	//
	// With WriteExt=false, configured extensions are serialized as raw bytes
	// and Str8 is not encoded.
	//
	// A stream can still be decoded into a typed value, provided an appropriate value
	// is provided, but the type cannot be inferred from the stream. If no appropriate
	// type is provided (e.g. decoding into a nil interface***REMOVED******REMOVED***), you get back
	// a []byte or string based on the setting of RawToString.
	WriteExt bool
	binaryEncodingType
***REMOVED***

func (h *MsgpackHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, &setExtWrapper***REMOVED***b: ext***REMOVED***)
***REMOVED***

func (h *MsgpackHandle) newEncDriver(e *Encoder) encDriver ***REMOVED***
	return &msgpackEncDriver***REMOVED***e: e, w: e.w, h: h***REMOVED***
***REMOVED***

func (h *MsgpackHandle) newDecDriver(d *Decoder) decDriver ***REMOVED***
	return &msgpackDecDriver***REMOVED***d: d, h: h, r: d.r, br: d.bytes***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) reset() ***REMOVED***
	e.w = e.e.w
***REMOVED***

func (d *msgpackDecDriver) reset() ***REMOVED***
	d.r, d.br = d.d.r, d.d.bytes
	d.bd, d.bdRead = 0, false
***REMOVED***

//--------------------------------------------------

type msgpackSpecRpcCodec struct ***REMOVED***
	rpcCodec
***REMOVED***

// /////////////// Spec RPC Codec ///////////////////
func (c *msgpackSpecRpcCodec) WriteRequest(r *rpc.Request, body interface***REMOVED******REMOVED***) error ***REMOVED***
	// WriteRequest can write to both a Go service, and other services that do
	// not abide by the 1 argument rule of a Go service.
	// We discriminate based on if the body is a MsgpackSpecRpcMultiArgs
	var bodyArr []interface***REMOVED******REMOVED***
	if m, ok := body.(MsgpackSpecRpcMultiArgs); ok ***REMOVED***
		bodyArr = ([]interface***REMOVED******REMOVED***)(m)
	***REMOVED*** else ***REMOVED***
		bodyArr = []interface***REMOVED******REMOVED******REMOVED***body***REMOVED***
	***REMOVED***
	r2 := []interface***REMOVED******REMOVED******REMOVED***0, uint32(r.Seq), r.ServiceMethod, bodyArr***REMOVED***
	return c.write(r2, nil, false, true)
***REMOVED***

func (c *msgpackSpecRpcCodec) WriteResponse(r *rpc.Response, body interface***REMOVED******REMOVED***) error ***REMOVED***
	var moe interface***REMOVED******REMOVED***
	if r.Error != "" ***REMOVED***
		moe = r.Error
	***REMOVED***
	if moe != nil && body != nil ***REMOVED***
		body = nil
	***REMOVED***
	r2 := []interface***REMOVED******REMOVED******REMOVED***1, uint32(r.Seq), moe, body***REMOVED***
	return c.write(r2, nil, false, true)
***REMOVED***

func (c *msgpackSpecRpcCodec) ReadResponseHeader(r *rpc.Response) error ***REMOVED***
	return c.parseCustomHeader(1, &r.Seq, &r.Error)
***REMOVED***

func (c *msgpackSpecRpcCodec) ReadRequestHeader(r *rpc.Request) error ***REMOVED***
	return c.parseCustomHeader(0, &r.Seq, &r.ServiceMethod)
***REMOVED***

func (c *msgpackSpecRpcCodec) ReadRequestBody(body interface***REMOVED******REMOVED***) error ***REMOVED***
	if body == nil ***REMOVED*** // read and discard
		return c.read(nil)
	***REMOVED***
	bodyArr := []interface***REMOVED******REMOVED******REMOVED***body***REMOVED***
	return c.read(&bodyArr)
***REMOVED***

func (c *msgpackSpecRpcCodec) parseCustomHeader(expectTypeByte byte, msgid *uint64, methodOrError *string) (err error) ***REMOVED***

	if c.isClosed() ***REMOVED***
		return io.EOF
	***REMOVED***

	// We read the response header by hand
	// so that the body can be decoded on its own from the stream at a later time.

	const fia byte = 0x94 //four item array descriptor value
	// Not sure why the panic of EOF is swallowed above.
	// if bs1 := c.dec.r.readn1(); bs1 != fia ***REMOVED***
	// 	err = fmt.Errorf("Unexpected value for array descriptor: Expecting %v. Received %v", fia, bs1)
	// 	return
	// ***REMOVED***
	var b byte
	b, err = c.br.ReadByte()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if b != fia ***REMOVED***
		err = fmt.Errorf("Unexpected value for array descriptor: Expecting %v. Received %v", fia, b)
		return
	***REMOVED***

	if err = c.read(&b); err != nil ***REMOVED***
		return
	***REMOVED***
	if b != expectTypeByte ***REMOVED***
		err = fmt.Errorf("Unexpected byte descriptor in header. Expecting %v. Received %v", expectTypeByte, b)
		return
	***REMOVED***
	if err = c.read(msgid); err != nil ***REMOVED***
		return
	***REMOVED***
	if err = c.read(methodOrError); err != nil ***REMOVED***
		return
	***REMOVED***
	return
***REMOVED***

//--------------------------------------------------

// msgpackSpecRpc is the implementation of Rpc that uses custom communication protocol
// as defined in the msgpack spec at https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md
type msgpackSpecRpc struct***REMOVED******REMOVED***

// MsgpackSpecRpc implements Rpc using the communication protocol defined in
// the msgpack spec at https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md .
// Its methods (ServerCodec and ClientCodec) return values that implement RpcCodecBuffered.
var MsgpackSpecRpc msgpackSpecRpc

func (x msgpackSpecRpc) ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec ***REMOVED***
	return &msgpackSpecRpcCodec***REMOVED***newRPCCodec(conn, h)***REMOVED***
***REMOVED***

func (x msgpackSpecRpc) ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec ***REMOVED***
	return &msgpackSpecRpcCodec***REMOVED***newRPCCodec(conn, h)***REMOVED***
***REMOVED***

var _ decDriver = (*msgpackDecDriver)(nil)
var _ encDriver = (*msgpackEncDriver)(nil)
