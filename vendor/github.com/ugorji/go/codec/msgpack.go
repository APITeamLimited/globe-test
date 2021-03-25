// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

/*
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
	"time"
)

const (
	mpPosFixNumMin byte = 0x00
	mpPosFixNumMax byte = 0x7f
	mpFixMapMin    byte = 0x80
	mpFixMapMax    byte = 0x8f
	mpFixArrayMin  byte = 0x90
	mpFixArrayMax  byte = 0x9f
	mpFixStrMin    byte = 0xa0
	mpFixStrMax    byte = 0xbf
	mpNil          byte = 0xc0
	_              byte = 0xc1
	mpFalse        byte = 0xc2
	mpTrue         byte = 0xc3
	mpFloat        byte = 0xca
	mpDouble       byte = 0xcb
	mpUint8        byte = 0xcc
	mpUint16       byte = 0xcd
	mpUint32       byte = 0xce
	mpUint64       byte = 0xcf
	mpInt8         byte = 0xd0
	mpInt16        byte = 0xd1
	mpInt32        byte = 0xd2
	mpInt64        byte = 0xd3

	// extensions below
	mpBin8     byte = 0xc4
	mpBin16    byte = 0xc5
	mpBin32    byte = 0xc6
	mpExt8     byte = 0xc7
	mpExt16    byte = 0xc8
	mpExt32    byte = 0xc9
	mpFixExt1  byte = 0xd4
	mpFixExt2  byte = 0xd5
	mpFixExt4  byte = 0xd6
	mpFixExt8  byte = 0xd7
	mpFixExt16 byte = 0xd8

	mpStr8  byte = 0xd9 // new
	mpStr16 byte = 0xda
	mpStr32 byte = 0xdb

	mpArray16 byte = 0xdc
	mpArray32 byte = 0xdd

	mpMap16 byte = 0xde
	mpMap32 byte = 0xdf

	mpNegFixNumMin byte = 0xe0
	mpNegFixNumMax byte = 0xff
)

var mpTimeExtTag int8 = -1
var mpTimeExtTagU = uint8(mpTimeExtTag)

// var mpdesc = map[byte]string***REMOVED***
// 	mpPosFixNumMin: "PosFixNumMin",
// 	mpPosFixNumMax: "PosFixNumMax",
// 	mpFixMapMin:    "FixMapMin",
// 	mpFixMapMax:    "FixMapMax",
// 	mpFixArrayMin:  "FixArrayMin",
// 	mpFixArrayMax:  "FixArrayMax",
// 	mpFixStrMin:    "FixStrMin",
// 	mpFixStrMax:    "FixStrMax",
// 	mpNil:          "Nil",
// 	mpFalse:        "False",
// 	mpTrue:         "True",
// 	mpFloat:        "Float",
// 	mpDouble:       "Double",
// 	mpUint8:        "Uint8",
// 	mpUint16:       "Uint16",
// 	mpUint32:       "Uint32",
// 	mpUint64:       "Uint64",
// 	mpInt8:         "Int8",
// 	mpInt16:        "Int16",
// 	mpInt32:        "Int32",
// 	mpInt64:        "Int64",
// 	mpBin8:         "Bin8",
// 	mpBin16:        "Bin16",
// 	mpBin32:        "Bin32",
// 	mpExt8:         "Ext8",
// 	mpExt16:        "Ext16",
// 	mpExt32:        "Ext32",
// 	mpFixExt1:      "FixExt1",
// 	mpFixExt2:      "FixExt2",
// 	mpFixExt4:      "FixExt4",
// 	mpFixExt8:      "FixExt8",
// 	mpFixExt16:     "FixExt16",
// 	mpStr8:         "Str8",
// 	mpStr16:        "Str16",
// 	mpStr32:        "Str32",
// 	mpArray16:      "Array16",
// 	mpArray32:      "Array32",
// 	mpMap16:        "Map16",
// 	mpMap32:        "Map32",
// 	mpNegFixNumMin: "NegFixNumMin",
// 	mpNegFixNumMax: "NegFixNumMax",
// ***REMOVED***

func mpdesc(bd byte) string ***REMOVED***
	switch bd ***REMOVED***
	case mpNil:
		return "nil"
	case mpFalse:
		return "false"
	case mpTrue:
		return "true"
	case mpFloat, mpDouble:
		return "float"
	case mpUint8, mpUint16, mpUint32, mpUint64:
		return "uint"
	case mpInt8, mpInt16, mpInt32, mpInt64:
		return "int"
	default:
		switch ***REMOVED***
		case bd >= mpPosFixNumMin && bd <= mpPosFixNumMax:
			return "int"
		case bd >= mpNegFixNumMin && bd <= mpNegFixNumMax:
			return "int"
		case bd == mpStr8, bd == mpStr16, bd == mpStr32, bd >= mpFixStrMin && bd <= mpFixStrMax:
			return "string|bytes"
		case bd == mpBin8, bd == mpBin16, bd == mpBin32:
			return "bytes"
		case bd == mpArray16, bd == mpArray32, bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			return "array"
		case bd == mpMap16, bd == mpMap32, bd >= mpFixMapMin && bd <= mpFixMapMax:
			return "map"
		case bd >= mpFixExt1 && bd <= mpFixExt16, bd >= mpExt8 && bd <= mpExt32:
			return "ext"
		default:
			return "unknown"
		***REMOVED***
	***REMOVED***
***REMOVED***

// MsgpackSpecRpcMultiArgs is a special type which signifies to the MsgpackSpecRpcCodec
// that the backend RPC service takes multiple arguments, which have been arranged
// in sequence in the slice.
//
// The Codec then passes it AS-IS to the rpc service (without wrapping it in an
// array of 1 element).
type MsgpackSpecRpcMultiArgs []interface***REMOVED******REMOVED***

// A MsgpackContainer type specifies the different types of msgpackContainers.
type msgpackContainerType struct ***REMOVED***
	fixCutoff             uint8
	bFixMin, b8, b16, b32 byte
	// hasFixMin, has8, has8Always bool
***REMOVED***

var (
	msgpackContainerRawLegacy = msgpackContainerType***REMOVED***
		32, mpFixStrMin, 0, mpStr16, mpStr32,
	***REMOVED***
	msgpackContainerStr = msgpackContainerType***REMOVED***
		32, mpFixStrMin, mpStr8, mpStr16, mpStr32, // true, true, false,
	***REMOVED***
	msgpackContainerBin = msgpackContainerType***REMOVED***
		0, 0, mpBin8, mpBin16, mpBin32, // false, true, true,
	***REMOVED***
	msgpackContainerList = msgpackContainerType***REMOVED***
		16, mpFixArrayMin, 0, mpArray16, mpArray32, // true, false, false,
	***REMOVED***
	msgpackContainerMap = msgpackContainerType***REMOVED***
		16, mpFixMapMin, 0, mpMap16, mpMap32, // true, false, false,
	***REMOVED***
)

//---------------------------------------------

type msgpackEncDriver struct ***REMOVED***
	noBuiltInTypes
	encDriverNoopContainerWriter
	h *MsgpackHandle
	x [8]byte
	_ [6]uint64 // padding
	e Encoder
***REMOVED***

func (e *msgpackEncDriver) encoder() *Encoder ***REMOVED***
	return &e.e
***REMOVED***

func (e *msgpackEncDriver) EncodeNil() ***REMOVED***
	e.e.encWr.writen1(mpNil)
***REMOVED***

func (e *msgpackEncDriver) EncodeInt(i int64) ***REMOVED***
	if e.h.PositiveIntUnsigned && i >= 0 ***REMOVED***
		e.EncodeUint(uint64(i))
	***REMOVED*** else if i > math.MaxInt8 ***REMOVED***
		if i <= math.MaxInt16 ***REMOVED***
			e.e.encWr.writen1(mpInt16)
			bigenHelper***REMOVED***e.x[:2], e.e.w()***REMOVED***.writeUint16(uint16(i))
		***REMOVED*** else if i <= math.MaxInt32 ***REMOVED***
			e.e.encWr.writen1(mpInt32)
			bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(uint32(i))
		***REMOVED*** else ***REMOVED***
			e.e.encWr.writen1(mpInt64)
			bigenHelper***REMOVED***e.x[:8], e.e.w()***REMOVED***.writeUint64(uint64(i))
		***REMOVED***
	***REMOVED*** else if i >= -32 ***REMOVED***
		if e.h.NoFixedNum ***REMOVED***
			e.e.encWr.writen2(mpInt8, byte(i))
		***REMOVED*** else ***REMOVED***
			e.e.encWr.writen1(byte(i))
		***REMOVED***
	***REMOVED*** else if i >= math.MinInt8 ***REMOVED***
		e.e.encWr.writen2(mpInt8, byte(i))
	***REMOVED*** else if i >= math.MinInt16 ***REMOVED***
		e.e.encWr.writen1(mpInt16)
		bigenHelper***REMOVED***e.x[:2], e.e.w()***REMOVED***.writeUint16(uint16(i))
	***REMOVED*** else if i >= math.MinInt32 ***REMOVED***
		e.e.encWr.writen1(mpInt32)
		bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(uint32(i))
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(mpInt64)
		bigenHelper***REMOVED***e.x[:8], e.e.w()***REMOVED***.writeUint64(uint64(i))
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeUint(i uint64) ***REMOVED***
	if i <= math.MaxInt8 ***REMOVED***
		if e.h.NoFixedNum ***REMOVED***
			e.e.encWr.writen2(mpUint8, byte(i))
		***REMOVED*** else ***REMOVED***
			e.e.encWr.writen1(byte(i))
		***REMOVED***
	***REMOVED*** else if i <= math.MaxUint8 ***REMOVED***
		e.e.encWr.writen2(mpUint8, byte(i))
	***REMOVED*** else if i <= math.MaxUint16 ***REMOVED***
		e.e.encWr.writen1(mpUint16)
		bigenHelper***REMOVED***e.x[:2], e.e.w()***REMOVED***.writeUint16(uint16(i))
	***REMOVED*** else if i <= math.MaxUint32 ***REMOVED***
		e.e.encWr.writen1(mpUint32)
		bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(uint32(i))
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(mpUint64)
		bigenHelper***REMOVED***e.x[:8], e.e.w()***REMOVED***.writeUint64(uint64(i))
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeBool(b bool) ***REMOVED***
	if b ***REMOVED***
		e.e.encWr.writen1(mpTrue)
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(mpFalse)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeFloat32(f float32) ***REMOVED***
	e.e.encWr.writen1(mpFloat)
	bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(math.Float32bits(f))
***REMOVED***

func (e *msgpackEncDriver) EncodeFloat64(f float64) ***REMOVED***
	e.e.encWr.writen1(mpDouble)
	bigenHelper***REMOVED***e.x[:8], e.e.w()***REMOVED***.writeUint64(math.Float64bits(f))
***REMOVED***

func (e *msgpackEncDriver) EncodeTime(t time.Time) ***REMOVED***
	if t.IsZero() ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	t = t.UTC()
	sec, nsec := t.Unix(), uint64(t.Nanosecond())
	var data64 uint64
	var l = 4
	if sec >= 0 && sec>>34 == 0 ***REMOVED***
		data64 = (nsec << 34) | uint64(sec)
		if data64&0xffffffff00000000 != 0 ***REMOVED***
			l = 8
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		l = 12
	***REMOVED***
	if e.h.WriteExt ***REMOVED***
		e.encodeExtPreamble(mpTimeExtTagU, l)
	***REMOVED*** else ***REMOVED***
		e.writeContainerLen(msgpackContainerRawLegacy, l)
	***REMOVED***
	switch l ***REMOVED***
	case 4:
		bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(uint32(data64))
	case 8:
		bigenHelper***REMOVED***e.x[:8], e.e.w()***REMOVED***.writeUint64(data64)
	case 12:
		bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(uint32(nsec))
		bigenHelper***REMOVED***e.x[:8], e.e.w()***REMOVED***.writeUint64(uint64(sec))
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext) ***REMOVED***
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
	if e.h.WriteExt ***REMOVED***
		e.encodeExtPreamble(uint8(xtag), len(bs))
		e.e.encWr.writeb(bs)
	***REMOVED*** else ***REMOVED***
		e.EncodeStringBytesRaw(bs)
	***REMOVED***
	if ext == SelfExt ***REMOVED***
		e.e.blist.put(bs)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeRawExt(re *RawExt) ***REMOVED***
	e.encodeExtPreamble(uint8(re.Tag), len(re.Data))
	e.e.encWr.writeb(re.Data)
***REMOVED***

func (e *msgpackEncDriver) encodeExtPreamble(xtag byte, l int) ***REMOVED***
	if l == 1 ***REMOVED***
		e.e.encWr.writen2(mpFixExt1, xtag)
	***REMOVED*** else if l == 2 ***REMOVED***
		e.e.encWr.writen2(mpFixExt2, xtag)
	***REMOVED*** else if l == 4 ***REMOVED***
		e.e.encWr.writen2(mpFixExt4, xtag)
	***REMOVED*** else if l == 8 ***REMOVED***
		e.e.encWr.writen2(mpFixExt8, xtag)
	***REMOVED*** else if l == 16 ***REMOVED***
		e.e.encWr.writen2(mpFixExt16, xtag)
	***REMOVED*** else if l < 256 ***REMOVED***
		e.e.encWr.writen2(mpExt8, byte(l))
		e.e.encWr.writen1(xtag)
	***REMOVED*** else if l < 65536 ***REMOVED***
		e.e.encWr.writen1(mpExt16)
		bigenHelper***REMOVED***e.x[:2], e.e.w()***REMOVED***.writeUint16(uint16(l))
		e.e.encWr.writen1(xtag)
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(mpExt32)
		bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(uint32(l))
		e.e.encWr.writen1(xtag)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) WriteArrayStart(length int) ***REMOVED***
	e.writeContainerLen(msgpackContainerList, length)
***REMOVED***

func (e *msgpackEncDriver) WriteMapStart(length int) ***REMOVED***
	e.writeContainerLen(msgpackContainerMap, length)
***REMOVED***

func (e *msgpackEncDriver) EncodeString(s string) ***REMOVED***
	var ct msgpackContainerType
	if e.h.WriteExt ***REMOVED***
		if e.h.StringToRaw ***REMOVED***
			ct = msgpackContainerBin
		***REMOVED*** else ***REMOVED***
			ct = msgpackContainerStr
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ct = msgpackContainerRawLegacy
	***REMOVED***
	e.writeContainerLen(ct, len(s))
	if len(s) > 0 ***REMOVED***
		e.e.encWr.writestr(s)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) EncodeStringBytesRaw(bs []byte) ***REMOVED***
	if bs == nil ***REMOVED***
		e.EncodeNil()
		return
	***REMOVED***
	if e.h.WriteExt ***REMOVED***
		e.writeContainerLen(msgpackContainerBin, len(bs))
	***REMOVED*** else ***REMOVED***
		e.writeContainerLen(msgpackContainerRawLegacy, len(bs))
	***REMOVED***
	if len(bs) > 0 ***REMOVED***
		e.e.encWr.writeb(bs)
	***REMOVED***
***REMOVED***

func (e *msgpackEncDriver) writeContainerLen(ct msgpackContainerType, l int) ***REMOVED***
	if ct.fixCutoff > 0 && l < int(ct.fixCutoff) ***REMOVED***
		e.e.encWr.writen1(ct.bFixMin | byte(l))
	***REMOVED*** else if ct.b8 > 0 && l < 256 ***REMOVED***
		e.e.encWr.writen2(ct.b8, uint8(l))
	***REMOVED*** else if l < 65536 ***REMOVED***
		e.e.encWr.writen1(ct.b16)
		bigenHelper***REMOVED***e.x[:2], e.e.w()***REMOVED***.writeUint16(uint16(l))
	***REMOVED*** else ***REMOVED***
		e.e.encWr.writen1(ct.b32)
		bigenHelper***REMOVED***e.x[:4], e.e.w()***REMOVED***.writeUint32(uint32(l))
	***REMOVED***
***REMOVED***

//---------------------------------------------

type msgpackDecDriver struct ***REMOVED***
	decDriverNoopContainerReader
	h *MsgpackHandle
	// b      [scratchByteArrayLen]byte
	bd     byte
	bdRead bool
	fnil   bool
	noBuiltInTypes
	_ [6]uint64 // padding
	d Decoder
***REMOVED***

func (d *msgpackDecDriver) decoder() *Decoder ***REMOVED***
	return &d.d
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
	d.fnil = false
	bd := d.bd
	n := d.d.naked()
	var decodeFurther bool

	switch bd ***REMOVED***
	case mpNil:
		n.v = valueTypeNil
		d.bdRead = false
		d.fnil = true
	case mpFalse:
		n.v = valueTypeBool
		n.b = false
	case mpTrue:
		n.v = valueTypeBool
		n.b = true

	case mpFloat:
		n.v = valueTypeFloat
		n.f = float64(math.Float32frombits(bigen.Uint32(d.d.decRd.readx(4))))
	case mpDouble:
		n.v = valueTypeFloat
		n.f = math.Float64frombits(bigen.Uint64(d.d.decRd.readx(8)))

	case mpUint8:
		n.v = valueTypeUint
		n.u = uint64(d.d.decRd.readn1())
	case mpUint16:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint16(d.d.decRd.readx(2)))
	case mpUint32:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint32(d.d.decRd.readx(4)))
	case mpUint64:
		n.v = valueTypeUint
		n.u = uint64(bigen.Uint64(d.d.decRd.readx(8)))

	case mpInt8:
		n.v = valueTypeInt
		n.i = int64(int8(d.d.decRd.readn1()))
	case mpInt16:
		n.v = valueTypeInt
		n.i = int64(int16(bigen.Uint16(d.d.decRd.readx(2))))
	case mpInt32:
		n.v = valueTypeInt
		n.i = int64(int32(bigen.Uint32(d.d.decRd.readx(4))))
	case mpInt64:
		n.v = valueTypeInt
		n.i = int64(int64(bigen.Uint64(d.d.decRd.readx(8))))

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
			if d.h.WriteExt || d.h.RawToString ***REMOVED***
				n.v = valueTypeString
				n.s = string(d.DecodeStringAsBytes())
			***REMOVED*** else ***REMOVED***
				n.v = valueTypeBytes
				n.l = d.DecodeBytes(nil, false)
			***REMOVED***
		case bd == mpBin8, bd == mpBin16, bd == mpBin32:
			decNakedReadRawBytes(d, &d.d, n, d.h.RawToString)
		case bd == mpArray16, bd == mpArray32, bd >= mpFixArrayMin && bd <= mpFixArrayMax:
			n.v = valueTypeArray
			decodeFurther = true
		case bd == mpMap16, bd == mpMap32, bd >= mpFixMapMin && bd <= mpFixMapMax:
			n.v = valueTypeMap
			decodeFurther = true
		case bd >= mpFixExt1 && bd <= mpFixExt16, bd >= mpExt8 && bd <= mpExt32:
			n.v = valueTypeExt
			clen := d.readExtLen()
			n.u = uint64(d.d.decRd.readn1())
			if n.u == uint64(mpTimeExtTagU) ***REMOVED***
				n.v = valueTypeTime
				n.t = d.decodeTime(clen)
			***REMOVED*** else if d.d.bytes ***REMOVED***
				n.l = d.d.decRd.readx(uint(clen))
			***REMOVED*** else ***REMOVED***
				n.l = decByteSlice(d.d.r(), clen, d.d.h.MaxInitLen, d.d.b[:])
			***REMOVED***
		default:
			d.d.errorf("cannot infer value: %s: Ox%x/%d/%s", msgBadDesc, bd, bd, mpdesc(bd))
		***REMOVED***
	***REMOVED***
	if !decodeFurther ***REMOVED***
		d.bdRead = false
	***REMOVED***
	if n.v == valueTypeUint && d.h.SignedInteger ***REMOVED***
		n.v = valueTypeInt
		n.i = int64(n.u)
	***REMOVED***
***REMOVED***

// int can be decoded from msgpack type: intXXX or uintXXX
func (d *msgpackDecDriver) DecodeInt64() (i int64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	switch d.bd ***REMOVED***
	case mpUint8:
		i = int64(uint64(d.d.decRd.readn1()))
	case mpUint16:
		i = int64(uint64(bigen.Uint16(d.d.decRd.readx(2))))
	case mpUint32:
		i = int64(uint64(bigen.Uint32(d.d.decRd.readx(4))))
	case mpUint64:
		i = int64(bigen.Uint64(d.d.decRd.readx(8)))
	case mpInt8:
		i = int64(int8(d.d.decRd.readn1()))
	case mpInt16:
		i = int64(int16(bigen.Uint16(d.d.decRd.readx(2))))
	case mpInt32:
		i = int64(int32(bigen.Uint32(d.d.decRd.readx(4))))
	case mpInt64:
		i = int64(bigen.Uint64(d.d.decRd.readx(8)))
	default:
		switch ***REMOVED***
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			i = int64(int8(d.bd))
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			i = int64(int8(d.bd))
		default:
			d.d.errorf("cannot decode signed integer: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
			return
		***REMOVED***
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// uint can be decoded from msgpack type: intXXX or uintXXX
func (d *msgpackDecDriver) DecodeUint64() (ui uint64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	switch d.bd ***REMOVED***
	case mpUint8:
		ui = uint64(d.d.decRd.readn1())
	case mpUint16:
		ui = uint64(bigen.Uint16(d.d.decRd.readx(2)))
	case mpUint32:
		ui = uint64(bigen.Uint32(d.d.decRd.readx(4)))
	case mpUint64:
		ui = bigen.Uint64(d.d.decRd.readx(8))
	case mpInt8:
		if i := int64(int8(d.d.decRd.readn1())); i >= 0 ***REMOVED***
			ui = uint64(i)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("assigning negative signed value: %v, to unsigned type", i)
			return
		***REMOVED***
	case mpInt16:
		if i := int64(int16(bigen.Uint16(d.d.decRd.readx(2)))); i >= 0 ***REMOVED***
			ui = uint64(i)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("assigning negative signed value: %v, to unsigned type", i)
			return
		***REMOVED***
	case mpInt32:
		if i := int64(int32(bigen.Uint32(d.d.decRd.readx(4)))); i >= 0 ***REMOVED***
			ui = uint64(i)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("assigning negative signed value: %v, to unsigned type", i)
			return
		***REMOVED***
	case mpInt64:
		if i := int64(bigen.Uint64(d.d.decRd.readx(8))); i >= 0 ***REMOVED***
			ui = uint64(i)
		***REMOVED*** else ***REMOVED***
			d.d.errorf("assigning negative signed value: %v, to unsigned type", i)
			return
		***REMOVED***
	default:
		switch ***REMOVED***
		case d.bd >= mpPosFixNumMin && d.bd <= mpPosFixNumMax:
			ui = uint64(d.bd)
		case d.bd >= mpNegFixNumMin && d.bd <= mpNegFixNumMax:
			d.d.errorf("assigning negative signed value: %v, to unsigned type", int(d.bd))
			return
		default:
			d.d.errorf("cannot decode unsigned integer: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
			return
		***REMOVED***
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// float can either be decoded from msgpack type: float, double or intX
func (d *msgpackDecDriver) DecodeFloat64() (f float64) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.bd == mpFloat ***REMOVED***
		f = float64(math.Float32frombits(bigen.Uint32(d.d.decRd.readx(4))))
	***REMOVED*** else if d.bd == mpDouble ***REMOVED***
		f = math.Float64frombits(bigen.Uint64(d.d.decRd.readx(8)))
	***REMOVED*** else ***REMOVED***
		f = float64(d.DecodeInt64())
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

// bool can be decoded from bool, fixnum 0 or 1.
func (d *msgpackDecDriver) DecodeBool() (b bool) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	if d.bd == mpFalse || d.bd == 0 ***REMOVED***
		// b = false
	***REMOVED*** else if d.bd == mpTrue || d.bd == 1 ***REMOVED***
		b = true
	***REMOVED*** else ***REMOVED***
		d.d.errorf("cannot decode bool: %s: %x/%s", msgBadDesc, d.bd, mpdesc(d.bd))
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *msgpackDecDriver) DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte) ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***

	bd := d.bd
	var clen int
	if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 ***REMOVED***
		clen = d.readContainerLen(msgpackContainerBin) // binary
	***REMOVED*** else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) ***REMOVED***
		clen = d.readContainerLen(msgpackContainerStr) // string/raw
	***REMOVED*** else if bd == mpArray16 || bd == mpArray32 ||
		(bd >= mpFixArrayMin && bd <= mpFixArrayMax) ***REMOVED***
		// check if an "array" of uint8's
		if zerocopy && len(bs) == 0 ***REMOVED***
			bs = d.d.b[:]
		***REMOVED***
		// bsOut, _ = fastpathTV.DecSliceUint8V(bs, true, d.d)
		slen := d.ReadArrayStart()
		bs = usableByteSlice(bs, slen)
		for i := 0; i < len(bs); i++ ***REMOVED***
			bs[i] = uint8(chkOvf.UintV(d.DecodeUint64(), 8))
		***REMOVED***
		return bs
	***REMOVED*** else ***REMOVED***
		d.d.errorf("invalid byte descriptor for decoding bytes, got: 0x%x", d.bd)
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
	return decByteSlice(d.d.r(), clen, d.h.MaxInitLen, bs)
***REMOVED***

func (d *msgpackDecDriver) DecodeStringAsBytes() (s []byte) ***REMOVED***
	return d.DecodeBytes(d.d.b[:], true)
***REMOVED***

func (d *msgpackDecDriver) readNextBd() ***REMOVED***
	d.bd = d.d.decRd.readn1()
	d.bdRead = true
***REMOVED***

func (d *msgpackDecDriver) uncacheRead() ***REMOVED***
	if d.bdRead ***REMOVED***
		d.d.decRd.unreadn1()
		d.bdRead = false
	***REMOVED***
***REMOVED***

func (d *msgpackDecDriver) advanceNil() (null bool) ***REMOVED***
	d.fnil = false
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	if d.bd == mpNil ***REMOVED***
		d.bdRead = false
		d.fnil = true
		null = true
	***REMOVED***
	return
***REMOVED***

func (d *msgpackDecDriver) Nil() bool ***REMOVED***
	return d.fnil
***REMOVED***

func (d *msgpackDecDriver) ContainerType() (vt valueType) ***REMOVED***
	if !d.bdRead ***REMOVED***
		d.readNextBd()
	***REMOVED***
	bd := d.bd
	d.fnil = false
	if bd == mpNil ***REMOVED***
		d.bdRead = false
		d.fnil = true
		return valueTypeNil
	***REMOVED*** else if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 ***REMOVED***
		return valueTypeBytes
	***REMOVED*** else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) ***REMOVED***
		if d.h.WriteExt || d.h.RawToString ***REMOVED*** // UTF-8 string (new spec)
			return valueTypeString
		***REMOVED***
		return valueTypeBytes // raw (old spec)
	***REMOVED*** else if bd == mpArray16 || bd == mpArray32 || (bd >= mpFixArrayMin && bd <= mpFixArrayMax) ***REMOVED***
		return valueTypeArray
	***REMOVED*** else if bd == mpMap16 || bd == mpMap32 || (bd >= mpFixMapMin && bd <= mpFixMapMax) ***REMOVED***
		return valueTypeMap
	***REMOVED***
	return valueTypeUnset
***REMOVED***

func (d *msgpackDecDriver) TryNil() (v bool) ***REMOVED***
	return d.advanceNil()
***REMOVED***

func (d *msgpackDecDriver) readContainerLen(ct msgpackContainerType) (clen int) ***REMOVED***
	bd := d.bd
	if bd == ct.b8 ***REMOVED***
		clen = int(d.d.decRd.readn1())
	***REMOVED*** else if bd == ct.b16 ***REMOVED***
		clen = int(bigen.Uint16(d.d.decRd.readx(2)))
	***REMOVED*** else if bd == ct.b32 ***REMOVED***
		clen = int(bigen.Uint32(d.d.decRd.readx(4)))
	***REMOVED*** else if (ct.bFixMin & bd) == ct.bFixMin ***REMOVED***
		clen = int(ct.bFixMin ^ bd)
	***REMOVED*** else ***REMOVED***
		d.d.errorf("cannot read container length: %s: hex: %x, decimal: %d", msgBadDesc, bd, bd)
		return
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

func (d *msgpackDecDriver) ReadMapStart() int ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return decContainerLenNil
	***REMOVED***
	return d.readContainerLen(msgpackContainerMap)
***REMOVED***

func (d *msgpackDecDriver) ReadArrayStart() int ***REMOVED***
	if d.advanceNil() ***REMOVED***
		return decContainerLenNil
	***REMOVED***
	return d.readContainerLen(msgpackContainerList)
***REMOVED***

func (d *msgpackDecDriver) readExtLen() (clen int) ***REMOVED***
	switch d.bd ***REMOVED***
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
		clen = int(d.d.decRd.readn1())
	case mpExt16:
		clen = int(bigen.Uint16(d.d.decRd.readx(2)))
	case mpExt32:
		clen = int(bigen.Uint32(d.d.decRd.readx(4)))
	default:
		d.d.errorf("decoding ext bytes: found unexpected byte: %x", d.bd)
		return
	***REMOVED***
	return
***REMOVED***

func (d *msgpackDecDriver) DecodeTime() (t time.Time) ***REMOVED***
	// decode time from string bytes or ext
	if d.advanceNil() ***REMOVED***
		return
	***REMOVED***
	bd := d.bd
	var clen int
	if bd == mpBin8 || bd == mpBin16 || bd == mpBin32 ***REMOVED***
		clen = d.readContainerLen(msgpackContainerBin) // binary
	***REMOVED*** else if bd == mpStr8 || bd == mpStr16 || bd == mpStr32 ||
		(bd >= mpFixStrMin && bd <= mpFixStrMax) ***REMOVED***
		clen = d.readContainerLen(msgpackContainerStr) // string/raw
	***REMOVED*** else ***REMOVED***
		// expect to see mpFixExt4,-1 OR mpFixExt8,-1 OR mpExt8,12,-1
		d.bdRead = false
		b2 := d.d.decRd.readn1()
		if d.bd == mpFixExt4 && b2 == mpTimeExtTagU ***REMOVED***
			clen = 4
		***REMOVED*** else if d.bd == mpFixExt8 && b2 == mpTimeExtTagU ***REMOVED***
			clen = 8
		***REMOVED*** else if d.bd == mpExt8 && b2 == 12 && d.d.decRd.readn1() == mpTimeExtTagU ***REMOVED***
			clen = 12
		***REMOVED*** else ***REMOVED***
			d.d.errorf("invalid stream for decoding time as extension: got 0x%x, 0x%x", d.bd, b2)
			return
		***REMOVED***
	***REMOVED***
	return d.decodeTime(clen)
***REMOVED***

func (d *msgpackDecDriver) decodeTime(clen int) (t time.Time) ***REMOVED***
	// bs = d.d.decRd.readx(clen)
	d.bdRead = false
	switch clen ***REMOVED***
	case 4:
		t = time.Unix(int64(bigen.Uint32(d.d.decRd.readx(4))), 0).UTC()
	case 8:
		tv := bigen.Uint64(d.d.decRd.readx(8))
		t = time.Unix(int64(tv&0x00000003ffffffff), int64(tv>>34)).UTC()
	case 12:
		nsec := bigen.Uint32(d.d.decRd.readx(4))
		sec := bigen.Uint64(d.d.decRd.readx(8))
		t = time.Unix(int64(sec), int64(nsec)).UTC()
	default:
		d.d.errorf("invalid length of bytes for decoding time - expecting 4 or 8 or 12, got %d", clen)
		return
	***REMOVED***
	return
***REMOVED***

func (d *msgpackDecDriver) DecodeExt(rv interface***REMOVED******REMOVED***, xtag uint64, ext Ext) ***REMOVED***
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

func (d *msgpackDecDriver) decodeExtV(verifyTag bool, tag byte) (xtag byte, xbs []byte) ***REMOVED***
	xbd := d.bd
	if xbd == mpBin8 || xbd == mpBin16 || xbd == mpBin32 ***REMOVED***
		xbs = d.DecodeBytes(nil, true)
	***REMOVED*** else if xbd == mpStr8 || xbd == mpStr16 || xbd == mpStr32 ||
		(xbd >= mpFixStrMin && xbd <= mpFixStrMax) ***REMOVED***
		xbs = d.DecodeStringAsBytes()
	***REMOVED*** else ***REMOVED***
		clen := d.readExtLen()
		xtag = d.d.decRd.readn1()
		if verifyTag && xtag != tag ***REMOVED***
			d.d.errorf("wrong extension tag - got %b, expecting %v", xtag, tag)
			return
		***REMOVED***
		if d.d.bytes ***REMOVED***
			xbs = d.d.decRd.readx(uint(clen))
		***REMOVED*** else ***REMOVED***
			xbs = decByteSlice(d.d.r(), clen, d.d.h.MaxInitLen, d.d.b[:])
		***REMOVED***
	***REMOVED***
	d.bdRead = false
	return
***REMOVED***

//--------------------------------------------------

//MsgpackHandle is a Handle for the Msgpack Schema-Free Encoding Format.
type MsgpackHandle struct ***REMOVED***
	binaryEncodingType
	BasicHandle

	// NoFixedNum says to output all signed integers as 2-bytes, never as 1-byte fixednum.
	NoFixedNum bool

	// WriteExt controls whether the new spec is honored.
	//
	// With WriteExt=true, we can encode configured extensions with extension tags
	// and encode string/[]byte/extensions in a way compatible with the new spec
	// but incompatible with the old spec.
	//
	// For compatibility with the old spec, set WriteExt=false.
	//
	// With WriteExt=false:
	//    configured extensions are serialized as raw bytes (not msgpack extensions).
	//    reserved byte descriptors like Str8 and those enabling the new msgpack Binary type
	//    are not encoded.
	WriteExt bool

	// PositiveIntUnsigned says to encode positive integers as unsigned.
	PositiveIntUnsigned bool

	_ [7]uint64 // padding (cache-aligned)
***REMOVED***

// Name returns the name of the handle: msgpack
func (h *MsgpackHandle) Name() string ***REMOVED*** return "msgpack" ***REMOVED***

func (h *MsgpackHandle) newEncDriver() encDriver ***REMOVED***
	var e = &msgpackEncDriver***REMOVED***h: h***REMOVED***
	e.e.e = e
	e.e.init(h)
	e.reset()
	return e
***REMOVED***

func (h *MsgpackHandle) newDecDriver() decDriver ***REMOVED***
	d := &msgpackDecDriver***REMOVED***h: h***REMOVED***
	d.d.d = d
	d.d.init(h)
	d.reset()
	return d
***REMOVED***

func (e *msgpackEncDriver) reset() ***REMOVED***
***REMOVED***

func (d *msgpackDecDriver) reset() ***REMOVED***
	d.bd, d.bdRead = 0, false
	d.fnil = false
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
	return c.write(r2, nil, false)
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
	return c.write(r2, nil, false)
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
	if cls := c.cls.load(); cls.closed ***REMOVED***
		return io.EOF
	***REMOVED***

	// We read the response header by hand
	// so that the body can be decoded on its own from the stream at a later time.

	const fia byte = 0x94 //four item array descriptor value

	var ba [1]byte
	var n int
	for ***REMOVED***
		n, err = c.r.Read(ba[:])
		if err != nil ***REMOVED***
			return
		***REMOVED***
		if n == 1 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	var b = ba[0]
	if b != fia ***REMOVED***
		err = fmt.Errorf("not array - %s %x/%s", msgBadDesc, b, mpdesc(b))
	***REMOVED*** else ***REMOVED***
		err = c.read(&b)
		if err == nil ***REMOVED***
			if b != expectTypeByte ***REMOVED***
				err = fmt.Errorf("%s - expecting %v but got %x/%s",
					msgBadDesc, expectTypeByte, b, mpdesc(b))
			***REMOVED*** else ***REMOVED***
				err = c.read(msgid)
				if err == nil ***REMOVED***
					err = c.read(methodOrError)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

//--------------------------------------------------

// msgpackSpecRpc is the implementation of Rpc that uses custom communication protocol
// as defined in the msgpack spec at https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md
type msgpackSpecRpc struct***REMOVED******REMOVED***

// MsgpackSpecRpc implements Rpc using the communication protocol defined in
// the msgpack spec at https://github.com/msgpack-rpc/msgpack-rpc/blob/master/spec.md .
//
// See GoRpc documentation, for information on buffering for better performance.
var MsgpackSpecRpc msgpackSpecRpc

func (x msgpackSpecRpc) ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec ***REMOVED***
	return &msgpackSpecRpcCodec***REMOVED***newRPCCodec(conn, h)***REMOVED***
***REMOVED***

func (x msgpackSpecRpc) ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec ***REMOVED***
	return &msgpackSpecRpcCodec***REMOVED***newRPCCodec(conn, h)***REMOVED***
***REMOVED***

var _ decDriver = (*msgpackDecDriver)(nil)
var _ encDriver = (*msgpackEncDriver)(nil)
