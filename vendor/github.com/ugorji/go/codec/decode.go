// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"
)

// Some tagging information for error messages.
const (
	msgBadDesc            = "Unrecognized descriptor byte"
	msgDecCannotExpandArr = "cannot expand go array from %v to stream length: %v"
)

var (
	onlyMapOrArrayCanDecodeIntoStructErr = errors.New("only encoded map or array can be decoded into a struct")
	cannotDecodeIntoNilErr               = errors.New("cannot decode into nil")
)

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReader interface ***REMOVED***
	unreadn1()

	// readx will use the implementation scratch buffer if possible i.e. n < len(scratchbuf), OR
	// just return a view of the []byte being decoded from.
	// Ensure you call detachZeroCopyBytes later if this needs to be sent outside codec control.
	readx(n int) []byte
	readb([]byte)
	readn1() uint8
	readn1eof() (v uint8, eof bool)
	numread() int // number of bytes read
	track()
	stopTrack() []byte
***REMOVED***

type decReaderByteScanner interface ***REMOVED***
	io.Reader
	io.ByteScanner
***REMOVED***

type decDriver interface ***REMOVED***
	// this will check if the next token is a break.
	CheckBreak() bool
	TryDecodeAsNil() bool
	// vt is one of: Bytes, String, Nil, Slice or Map. Return unSet if not known.
	ContainerType() (vt valueType)
	IsBuiltinType(rt uintptr) bool
	DecodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***)

	// DecodeNaked will decode primitives (number, bool, string, []byte) and RawExt.
	// For maps and arrays, it will not do the decoding in-band, but will signal
	// the decoder, so that is done later, by setting the decNaked.valueType field.
	//
	// Note: Numbers are decoded as int64, uint64, float64 only (no smaller sized number types).
	// for extensions, DecodeNaked must read the tag and the []byte if it exists.
	// if the []byte is not read, then kInterfaceNaked will treat it as a Handle
	// that stores the subsequent value in-band, and complete reading the RawExt.
	//
	// extensions should also use readx to decode them, for efficiency.
	// kInterface will extract the detached byte slice if it has to pass it outside its realm.
	DecodeNaked()
	DecodeInt(bitsize uint8) (i int64)
	DecodeUint(bitsize uint8) (ui uint64)
	DecodeFloat(chkOverflow32 bool) (f float64)
	DecodeBool() (b bool)
	// DecodeString can also decode symbols.
	// It looks redundant as DecodeBytes is available.
	// However, some codecs (e.g. binc) support symbols and can
	// return a pre-stored string value, meaning that it can bypass
	// the cost of []byte->string conversion.
	DecodeString() (s string)

	// DecodeBytes may be called directly, without going through reflection.
	// Consequently, it must be designed to handle possible nil.
	DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte)

	// decodeExt will decode into a *RawExt or into an extension.
	DecodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext) (realxtag uint64)
	// decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte)
	ReadMapStart() int
	ReadArrayStart() int

	reset()
	uncacheRead()
***REMOVED***

type decNoSeparator struct ***REMOVED***
***REMOVED***

func (_ decNoSeparator) ReadEnd() ***REMOVED******REMOVED***

// func (_ decNoSeparator) uncacheRead() ***REMOVED******REMOVED***

type DecodeOptions struct ***REMOVED***
	// MapType specifies type to use during schema-less decoding of a map in the stream.
	// If nil, we use map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
	MapType reflect.Type

	// SliceType specifies type to use during schema-less decoding of an array in the stream.
	// If nil, we use []interface***REMOVED******REMOVED***
	SliceType reflect.Type

	// MaxInitLen defines the maxinum initial length that we "make" a collection (string, slice, map, chan).
	// If 0 or negative, we default to a sensible value based on the size of an element in the collection.
	//
	// For example, when decoding, a stream may say that it has 2^64 elements.
	// We should not auto-matically provision a slice of that length, to prevent Out-Of-Memory crash.
	// Instead, we provision up to MaxInitLen, fill that up, and start appending after that.
	MaxInitLen int

	// If ErrorIfNoField, return an error when decoding a map
	// from a codec stream into a struct, and no matching struct field is found.
	ErrorIfNoField bool

	// If ErrorIfNoArrayExpand, return an error when decoding a slice/array that cannot be expanded.
	// For example, the stream contains an array of 8 items, but you are decoding into a [4]T array,
	// or you are decoding into a slice of length 4 which is non-addressable (and so cannot be set).
	ErrorIfNoArrayExpand bool

	// If SignedInteger, use the int64 during schema-less decoding of unsigned values (not uint64).
	SignedInteger bool

	// MapValueReset controls how we decode into a map value.
	//
	// By default, we MAY retrieve the mapping for a key, and then decode into that.
	// However, especially with big maps, that retrieval may be expensive and unnecessary
	// if the stream already contains all that is necessary to recreate the value.
	//
	// If true, we will never retrieve the previous mapping,
	// but rather decode into a new value and set that in the map.
	//
	// If false, we will retrieve the previous mapping if necessary e.g.
	// the previous mapping is a pointer, or is a struct or array with pre-set state,
	// or is an interface.
	MapValueReset bool

	// InterfaceReset controls how we decode into an interface.
	//
	// By default, when we see a field that is an interface***REMOVED***...***REMOVED***,
	// or a map with interface***REMOVED***...***REMOVED*** value, we will attempt decoding into the
	// "contained" value.
	//
	// However, this prevents us from reading a string into an interface***REMOVED******REMOVED***
	// that formerly contained a number.
	//
	// If true, we will decode into a new "blank" value, and set that in the interface.
	// If false, we will decode into whatever is contained in the interface.
	InterfaceReset bool

	// InternString controls interning of strings during decoding.
	//
	// Some handles, e.g. json, typically will read map keys as strings.
	// If the set of keys are finite, it may help reduce allocation to
	// look them up from a map (than to allocate them afresh).
	//
	// Note: Handles will be smart when using the intern functionality.
	// So everything will not be interned.
	InternString bool

	// PreferArrayOverSlice controls whether to decode to an array or a slice.
	//
	// This only impacts decoding into a nil interface***REMOVED******REMOVED***.
	// Consequently, it has no effect on codecgen.
	//
	// *Note*: This only applies if using go1.5 and above,
	// as it requires reflect.ArrayOf support which was absent before go1.5.
	PreferArrayOverSlice bool
***REMOVED***

// ------------------------------------

// ioDecByteScanner implements Read(), ReadByte(...), UnreadByte(...) methods
// of io.Reader, io.ByteScanner.
type ioDecByteScanner struct ***REMOVED***
	r  io.Reader
	l  byte    // last byte
	ls byte    // last byte status. 0: init-canDoNothing, 1: canRead, 2: canUnread
	b  [1]byte // tiny buffer for reading single bytes
***REMOVED***

func (z *ioDecByteScanner) Read(p []byte) (n int, err error) ***REMOVED***
	var firstByte bool
	if z.ls == 1 ***REMOVED***
		z.ls = 2
		p[0] = z.l
		if len(p) == 1 ***REMOVED***
			n = 1
			return
		***REMOVED***
		firstByte = true
		p = p[1:]
	***REMOVED***
	n, err = z.r.Read(p)
	if n > 0 ***REMOVED***
		if err == io.EOF && n == len(p) ***REMOVED***
			err = nil // read was successful, so postpone EOF (till next time)
		***REMOVED***
		z.l = p[n-1]
		z.ls = 2
	***REMOVED***
	if firstByte ***REMOVED***
		n++
	***REMOVED***
	return
***REMOVED***

func (z *ioDecByteScanner) ReadByte() (c byte, err error) ***REMOVED***
	n, err := z.Read(z.b[:])
	if n == 1 ***REMOVED***
		c = z.b[0]
		if err == io.EOF ***REMOVED***
			err = nil // read was successful, so postpone EOF (till next time)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (z *ioDecByteScanner) UnreadByte() (err error) ***REMOVED***
	x := z.ls
	if x == 0 ***REMOVED***
		err = errors.New("cannot unread - nothing has been read")
	***REMOVED*** else if x == 1 ***REMOVED***
		err = errors.New("cannot unread - last byte has not been read")
	***REMOVED*** else if x == 2 ***REMOVED***
		z.ls = 1
	***REMOVED***
	return
***REMOVED***

// ioDecReader is a decReader that reads off an io.Reader
type ioDecReader struct ***REMOVED***
	br decReaderByteScanner
	// temp byte array re-used internally for efficiency during read.
	// shares buffer with Decoder, so we keep size of struct within 8 words.
	x   *[scratchByteArrayLen]byte
	bs  ioDecByteScanner
	n   int    // num read
	tr  []byte // tracking bytes read
	trb bool
***REMOVED***

func (z *ioDecReader) numread() int ***REMOVED***
	return z.n
***REMOVED***

func (z *ioDecReader) readx(n int) (bs []byte) ***REMOVED***
	if n <= 0 ***REMOVED***
		return
	***REMOVED***
	if n < len(z.x) ***REMOVED***
		bs = z.x[:n]
	***REMOVED*** else ***REMOVED***
		bs = make([]byte, n)
	***REMOVED***
	if _, err := io.ReadAtLeast(z.br, bs, n); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	z.n += len(bs)
	if z.trb ***REMOVED***
		z.tr = append(z.tr, bs...)
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) readb(bs []byte) ***REMOVED***
	if len(bs) == 0 ***REMOVED***
		return
	***REMOVED***
	n, err := io.ReadAtLeast(z.br, bs, len(bs))
	z.n += n
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if z.trb ***REMOVED***
		z.tr = append(z.tr, bs...)
	***REMOVED***
***REMOVED***

func (z *ioDecReader) readn1() (b uint8) ***REMOVED***
	b, err := z.br.ReadByte()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	z.n++
	if z.trb ***REMOVED***
		z.tr = append(z.tr, b)
	***REMOVED***
	return b
***REMOVED***

func (z *ioDecReader) readn1eof() (b uint8, eof bool) ***REMOVED***
	b, err := z.br.ReadByte()
	if err == nil ***REMOVED***
		z.n++
		if z.trb ***REMOVED***
			z.tr = append(z.tr, b)
		***REMOVED***
	***REMOVED*** else if err == io.EOF ***REMOVED***
		eof = true
	***REMOVED*** else ***REMOVED***
		panic(err)
	***REMOVED***
	return
***REMOVED***

func (z *ioDecReader) unreadn1() ***REMOVED***
	err := z.br.UnreadByte()
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	z.n--
	if z.trb ***REMOVED***
		if l := len(z.tr) - 1; l >= 0 ***REMOVED***
			z.tr = z.tr[:l]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (z *ioDecReader) track() ***REMOVED***
	if z.tr != nil ***REMOVED***
		z.tr = z.tr[:0]
	***REMOVED***
	z.trb = true
***REMOVED***

func (z *ioDecReader) stopTrack() (bs []byte) ***REMOVED***
	z.trb = false
	return z.tr
***REMOVED***

// ------------------------------------

var bytesDecReaderCannotUnreadErr = errors.New("cannot unread last byte read")

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct ***REMOVED***
	b []byte // data
	c int    // cursor
	a int    // available
	t int    // track start
***REMOVED***

func (z *bytesDecReader) reset(in []byte) ***REMOVED***
	z.b = in
	z.a = len(in)
	z.c = 0
	z.t = 0
***REMOVED***

func (z *bytesDecReader) numread() int ***REMOVED***
	return z.c
***REMOVED***

func (z *bytesDecReader) unreadn1() ***REMOVED***
	if z.c == 0 || len(z.b) == 0 ***REMOVED***
		panic(bytesDecReaderCannotUnreadErr)
	***REMOVED***
	z.c--
	z.a++
	return
***REMOVED***

func (z *bytesDecReader) readx(n int) (bs []byte) ***REMOVED***
	// slicing from a non-constant start position is more expensive,
	// as more computation is required to decipher the pointer start position.
	// However, we do it only once, and it's better than reslicing both z.b and return value.

	if n <= 0 ***REMOVED***
	***REMOVED*** else if z.a == 0 ***REMOVED***
		panic(io.EOF)
	***REMOVED*** else if n > z.a ***REMOVED***
		panic(io.ErrUnexpectedEOF)
	***REMOVED*** else ***REMOVED***
		c0 := z.c
		z.c = c0 + n
		z.a = z.a - n
		bs = z.b[c0:z.c]
	***REMOVED***
	return
***REMOVED***

func (z *bytesDecReader) readn1() (v uint8) ***REMOVED***
	if z.a == 0 ***REMOVED***
		panic(io.EOF)
	***REMOVED***
	v = z.b[z.c]
	z.c++
	z.a--
	return
***REMOVED***

func (z *bytesDecReader) readn1eof() (v uint8, eof bool) ***REMOVED***
	if z.a == 0 ***REMOVED***
		eof = true
		return
	***REMOVED***
	v = z.b[z.c]
	z.c++
	z.a--
	return
***REMOVED***

func (z *bytesDecReader) readb(bs []byte) ***REMOVED***
	copy(bs, z.readx(len(bs)))
***REMOVED***

func (z *bytesDecReader) track() ***REMOVED***
	z.t = z.c
***REMOVED***

func (z *bytesDecReader) stopTrack() (bs []byte) ***REMOVED***
	return z.b[z.t:z.c]
***REMOVED***

// ------------------------------------

type decFnInfo struct ***REMOVED***
	d     *Decoder
	ti    *typeInfo
	xfFn  Ext
	xfTag uint64
	seq   seqType
***REMOVED***

// ----------------------------------------

type decFn struct ***REMOVED***
	i decFnInfo
	f func(*decFnInfo, reflect.Value)
***REMOVED***

func (f *decFnInfo) builtin(rv reflect.Value) ***REMOVED***
	f.d.d.DecodeBuiltin(f.ti.rtid, rv.Addr().Interface())
***REMOVED***

func (f *decFnInfo) rawExt(rv reflect.Value) ***REMOVED***
	f.d.d.DecodeExt(rv.Addr().Interface(), 0, nil)
***REMOVED***

func (f *decFnInfo) raw(rv reflect.Value) ***REMOVED***
	rv.SetBytes(f.d.raw())
***REMOVED***

func (f *decFnInfo) ext(rv reflect.Value) ***REMOVED***
	f.d.d.DecodeExt(rv.Addr().Interface(), f.xfTag, f.xfFn)
***REMOVED***

func (f *decFnInfo) getValueForUnmarshalInterface(rv reflect.Value, indir int8) (v interface***REMOVED******REMOVED***) ***REMOVED***
	if indir == -1 ***REMOVED***
		v = rv.Addr().Interface()
	***REMOVED*** else if indir == 0 ***REMOVED***
		v = rv.Interface()
	***REMOVED*** else ***REMOVED***
		for j := int8(0); j < indir; j++ ***REMOVED***
			if rv.IsNil() ***REMOVED***
				rv.Set(reflect.New(rv.Type().Elem()))
			***REMOVED***
			rv = rv.Elem()
		***REMOVED***
		v = rv.Interface()
	***REMOVED***
	return
***REMOVED***

func (f *decFnInfo) selferUnmarshal(rv reflect.Value) ***REMOVED***
	f.getValueForUnmarshalInterface(rv, f.ti.csIndir).(Selfer).CodecDecodeSelf(f.d)
***REMOVED***

func (f *decFnInfo) binaryUnmarshal(rv reflect.Value) ***REMOVED***
	bm := f.getValueForUnmarshalInterface(rv, f.ti.bunmIndir).(encoding.BinaryUnmarshaler)
	xbs := f.d.d.DecodeBytes(nil, false, true)
	if fnerr := bm.UnmarshalBinary(xbs); fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (f *decFnInfo) textUnmarshal(rv reflect.Value) ***REMOVED***
	tm := f.getValueForUnmarshalInterface(rv, f.ti.tunmIndir).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(f.d.d.DecodeBytes(f.d.b[:], true, true))
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (f *decFnInfo) jsonUnmarshal(rv reflect.Value) ***REMOVED***
	tm := f.getValueForUnmarshalInterface(rv, f.ti.junmIndir).(jsonUnmarshaler)
	// bs := f.d.d.DecodeBytes(f.d.b[:], true, true)
	// grab the bytes to be read, as UnmarshalJSON needs the full JSON so as to unmarshal it itself.
	fnerr := tm.UnmarshalJSON(f.d.nextValueBytes())
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (f *decFnInfo) kErr(rv reflect.Value) ***REMOVED***
	f.d.errorf("no decoding function defined for kind %v", rv.Kind())
***REMOVED***

func (f *decFnInfo) kString(rv reflect.Value) ***REMOVED***
	rv.SetString(f.d.d.DecodeString())
***REMOVED***

func (f *decFnInfo) kBool(rv reflect.Value) ***REMOVED***
	rv.SetBool(f.d.d.DecodeBool())
***REMOVED***

func (f *decFnInfo) kInt(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.d.d.DecodeInt(intBitsize))
***REMOVED***

func (f *decFnInfo) kInt64(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.d.d.DecodeInt(64))
***REMOVED***

func (f *decFnInfo) kInt32(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.d.d.DecodeInt(32))
***REMOVED***

func (f *decFnInfo) kInt8(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.d.d.DecodeInt(8))
***REMOVED***

func (f *decFnInfo) kInt16(rv reflect.Value) ***REMOVED***
	rv.SetInt(f.d.d.DecodeInt(16))
***REMOVED***

func (f *decFnInfo) kFloat32(rv reflect.Value) ***REMOVED***
	rv.SetFloat(f.d.d.DecodeFloat(true))
***REMOVED***

func (f *decFnInfo) kFloat64(rv reflect.Value) ***REMOVED***
	rv.SetFloat(f.d.d.DecodeFloat(false))
***REMOVED***

func (f *decFnInfo) kUint8(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.d.d.DecodeUint(8))
***REMOVED***

func (f *decFnInfo) kUint64(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.d.d.DecodeUint(64))
***REMOVED***

func (f *decFnInfo) kUint(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.d.d.DecodeUint(uintBitsize))
***REMOVED***

func (f *decFnInfo) kUintptr(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.d.d.DecodeUint(uintBitsize))
***REMOVED***

func (f *decFnInfo) kUint32(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.d.d.DecodeUint(32))
***REMOVED***

func (f *decFnInfo) kUint16(rv reflect.Value) ***REMOVED***
	rv.SetUint(f.d.d.DecodeUint(16))
***REMOVED***

// func (f *decFnInfo) kPtr(rv reflect.Value) ***REMOVED***
// 	debugf(">>>>>>> ??? decode kPtr called - shouldn't get called")
// 	if rv.IsNil() ***REMOVED***
// 		rv.Set(reflect.New(rv.Type().Elem()))
// 	***REMOVED***
// 	f.d.decodeValue(rv.Elem())
// ***REMOVED***

// var kIntfCtr uint64

func (f *decFnInfo) kInterfaceNaked() (rvn reflect.Value) ***REMOVED***
	// nil interface:
	// use some hieristics to decode it appropriately
	// based on the detected next value in the stream.
	d := f.d
	d.d.DecodeNaked()
	n := &d.n
	if n.v == valueTypeNil ***REMOVED***
		return
	***REMOVED***
	// We cannot decode non-nil stream value into nil interface with methods (e.g. io.Reader).
	// if num := f.ti.rt.NumMethod(); num > 0 ***REMOVED***
	if f.ti.numMeth > 0 ***REMOVED***
		d.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
		return
	***REMOVED***
	// var useRvn bool
	switch n.v ***REMOVED***
	case valueTypeMap:
		// if d.h.MapType == nil || d.h.MapType == mapIntfIntfTyp ***REMOVED***
		// ***REMOVED*** else if d.h.MapType == mapStrIntfTyp ***REMOVED*** // for json performance
		// ***REMOVED***
		if d.mtid == 0 || d.mtid == mapIntfIntfTypId ***REMOVED***
			l := len(n.ms)
			n.ms = append(n.ms, nil)
			var v2 interface***REMOVED******REMOVED*** = &n.ms[l]
			d.decode(v2)
			rvn = reflect.ValueOf(v2).Elem()
			n.ms = n.ms[:l]
		***REMOVED*** else if d.mtid == mapStrIntfTypId ***REMOVED*** // for json performance
			l := len(n.ns)
			n.ns = append(n.ns, nil)
			var v2 interface***REMOVED******REMOVED*** = &n.ns[l]
			d.decode(v2)
			rvn = reflect.ValueOf(v2).Elem()
			n.ns = n.ns[:l]
		***REMOVED*** else ***REMOVED***
			rvn = reflect.New(d.h.MapType).Elem()
			d.decodeValue(rvn, nil)
		***REMOVED***
	case valueTypeArray:
		// if d.h.SliceType == nil || d.h.SliceType == intfSliceTyp ***REMOVED***
		if d.stid == 0 || d.stid == intfSliceTypId ***REMOVED***
			l := len(n.ss)
			n.ss = append(n.ss, nil)
			var v2 interface***REMOVED******REMOVED*** = &n.ss[l]
			d.decode(v2)
			n.ss = n.ss[:l]
			rvn = reflect.ValueOf(v2).Elem()
			if reflectArrayOfSupported && d.stid == 0 && d.h.PreferArrayOverSlice ***REMOVED***
				rvn = reflectArrayOf(rvn)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			rvn = reflect.New(d.h.SliceType).Elem()
			d.decodeValue(rvn, nil)
		***REMOVED***
	case valueTypeExt:
		var v interface***REMOVED******REMOVED***
		tag, bytes := n.u, n.l // calling decode below might taint the values
		if bytes == nil ***REMOVED***
			l := len(n.is)
			n.is = append(n.is, nil)
			v2 := &n.is[l]
			d.decode(v2)
			v = *v2
			n.is = n.is[:l]
		***REMOVED***
		bfn := d.h.getExtForTag(tag)
		if bfn == nil ***REMOVED***
			var re RawExt
			re.Tag = tag
			re.Data = detachZeroCopyBytes(d.bytes, nil, bytes)
			rvn = reflect.ValueOf(re)
		***REMOVED*** else ***REMOVED***
			rvnA := reflect.New(bfn.rt)
			rvn = rvnA.Elem()
			if bytes != nil ***REMOVED***
				bfn.ext.ReadExt(rvnA.Interface(), bytes)
			***REMOVED*** else ***REMOVED***
				bfn.ext.UpdateExt(rvnA.Interface(), v)
			***REMOVED***
		***REMOVED***
	case valueTypeNil:
		// no-op
	case valueTypeInt:
		rvn = reflect.ValueOf(&n.i).Elem()
	case valueTypeUint:
		rvn = reflect.ValueOf(&n.u).Elem()
	case valueTypeFloat:
		rvn = reflect.ValueOf(&n.f).Elem()
	case valueTypeBool:
		rvn = reflect.ValueOf(&n.b).Elem()
	case valueTypeString, valueTypeSymbol:
		rvn = reflect.ValueOf(&n.s).Elem()
	case valueTypeBytes:
		rvn = reflect.ValueOf(&n.l).Elem()
	case valueTypeTimestamp:
		rvn = reflect.ValueOf(&n.t).Elem()
	default:
		panic(fmt.Errorf("kInterfaceNaked: unexpected valueType: %d", n.v))
	***REMOVED***
	return
***REMOVED***

func (f *decFnInfo) kInterface(rv reflect.Value) ***REMOVED***
	// debugf("\t===> kInterface")

	// Note:
	// A consequence of how kInterface works, is that
	// if an interface already contains something, we try
	// to decode into what was there before.
	// We do not replace with a generic value (as got from decodeNaked).

	var rvn reflect.Value
	if rv.IsNil() ***REMOVED***
		rvn = f.kInterfaceNaked()
		if rvn.IsValid() ***REMOVED***
			rv.Set(rvn)
		***REMOVED***
	***REMOVED*** else if f.d.h.InterfaceReset ***REMOVED***
		rvn = f.kInterfaceNaked()
		if rvn.IsValid() ***REMOVED***
			rv.Set(rvn)
		***REMOVED*** else ***REMOVED***
			// reset to zero value based on current type in there.
			rv.Set(reflect.Zero(rv.Elem().Type()))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		rvn = rv.Elem()
		// Note: interface***REMOVED******REMOVED*** is settable, but underlying type may not be.
		// Consequently, we have to set the reflect.Value directly.
		// if underlying type is settable (e.g. ptr or interface),
		// we just decode into it.
		// Else we create a settable value, decode into it, and set on the interface.
		if rvn.CanSet() ***REMOVED***
			f.d.decodeValue(rvn, nil)
		***REMOVED*** else ***REMOVED***
			rvn2 := reflect.New(rvn.Type()).Elem()
			rvn2.Set(rvn)
			f.d.decodeValue(rvn2, nil)
			rv.Set(rvn2)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *decFnInfo) kStruct(rv reflect.Value) ***REMOVED***
	fti := f.ti
	d := f.d
	dd := d.d
	cr := d.cr
	ctyp := dd.ContainerType()
	if ctyp == valueTypeMap ***REMOVED***
		containerLen := dd.ReadMapStart()
		if containerLen == 0 ***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapEnd)
			***REMOVED***
			return
		***REMOVED***
		tisfi := fti.sfi
		hasLen := containerLen >= 0
		if hasLen ***REMOVED***
			for j := 0; j < containerLen; j++ ***REMOVED***
				// rvkencname := dd.DecodeString()
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				rvkencname := stringView(dd.DecodeBytes(f.d.b[:], true, true))
				// rvksi := ti.getForEncName(rvkencname)
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				if k := fti.indexForEncName(rvkencname); k > -1 ***REMOVED***
					si := tisfi[k]
					if dd.TryDecodeAsNil() ***REMOVED***
						si.setToZeroValue(rv)
					***REMOVED*** else ***REMOVED***
						d.decodeValue(si.field(rv, true), nil)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					d.structFieldNotFound(-1, rvkencname)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for j := 0; !dd.CheckBreak(); j++ ***REMOVED***
				// rvkencname := dd.DecodeString()
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				rvkencname := stringView(dd.DecodeBytes(f.d.b[:], true, true))
				// rvksi := ti.getForEncName(rvkencname)
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				if k := fti.indexForEncName(rvkencname); k > -1 ***REMOVED***
					si := tisfi[k]
					if dd.TryDecodeAsNil() ***REMOVED***
						si.setToZeroValue(rv)
					***REMOVED*** else ***REMOVED***
						d.decodeValue(si.field(rv, true), nil)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					d.structFieldNotFound(-1, rvkencname)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if cr != nil ***REMOVED***
			cr.sendContainerState(containerMapEnd)
		***REMOVED***
	***REMOVED*** else if ctyp == valueTypeArray ***REMOVED***
		containerLen := dd.ReadArrayStart()
		if containerLen == 0 ***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerArrayEnd)
			***REMOVED***
			return
		***REMOVED***
		// Not much gain from doing it two ways for array.
		// Arrays are not used as much for structs.
		hasLen := containerLen >= 0
		for j, si := range fti.sfip ***REMOVED***
			if hasLen ***REMOVED***
				if j == containerLen ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else if dd.CheckBreak() ***REMOVED***
				break
			***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerArrayElem)
			***REMOVED***
			if dd.TryDecodeAsNil() ***REMOVED***
				si.setToZeroValue(rv)
			***REMOVED*** else ***REMOVED***
				d.decodeValue(si.field(rv, true), nil)
			***REMOVED***
		***REMOVED***
		if containerLen > len(fti.sfip) ***REMOVED***
			// read remaining values and throw away
			for j := len(fti.sfip); j < containerLen; j++ ***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerArrayElem)
				***REMOVED***
				d.structFieldNotFound(j, "")
			***REMOVED***
		***REMOVED***
		if cr != nil ***REMOVED***
			cr.sendContainerState(containerArrayEnd)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		f.d.error(onlyMapOrArrayCanDecodeIntoStructErr)
		return
	***REMOVED***
***REMOVED***

func (f *decFnInfo) kSlice(rv reflect.Value) ***REMOVED***
	// A slice can be set from a map or array in stream.
	// This way, the order can be kept (as order is lost with map).
	ti := f.ti
	d := f.d
	dd := d.d
	rtelem0 := ti.rt.Elem()
	ctyp := dd.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString ***REMOVED***
		// you can only decode bytes or string in the stream into a slice or array of bytes
		if !(ti.rtid == uint8SliceTypId || rtelem0.Kind() == reflect.Uint8) ***REMOVED***
			f.d.errorf("bytes or string in the stream must be decoded into a slice or array of bytes, not %v", ti.rt)
		***REMOVED***
		if f.seq == seqTypeChan ***REMOVED***
			bs2 := dd.DecodeBytes(nil, false, true)
			ch := rv.Interface().(chan<- byte)
			for _, b := range bs2 ***REMOVED***
				ch <- b
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			rvbs := rv.Bytes()
			bs2 := dd.DecodeBytes(rvbs, false, false)
			if rvbs == nil && bs2 != nil || rvbs != nil && bs2 == nil || len(bs2) != len(rvbs) ***REMOVED***
				if rv.CanSet() ***REMOVED***
					rv.SetBytes(bs2)
				***REMOVED*** else ***REMOVED***
					copy(rvbs, bs2)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	// array := f.seq == seqTypeChan

	slh, containerLenS := d.decSliceHelperStart() // only expects valueType(Array|Map)

	// // an array can never return a nil slice. so no need to check f.array here.
	if containerLenS == 0 ***REMOVED***
		if f.seq == seqTypeSlice ***REMOVED***
			if rv.IsNil() ***REMOVED***
				rv.Set(reflect.MakeSlice(ti.rt, 0, 0))
			***REMOVED*** else ***REMOVED***
				rv.SetLen(0)
			***REMOVED***
		***REMOVED*** else if f.seq == seqTypeChan ***REMOVED***
			if rv.IsNil() ***REMOVED***
				rv.Set(reflect.MakeChan(ti.rt, 0))
			***REMOVED***
		***REMOVED***
		slh.End()
		return
	***REMOVED***

	rtelem := rtelem0
	for rtelem.Kind() == reflect.Ptr ***REMOVED***
		rtelem = rtelem.Elem()
	***REMOVED***
	fn := d.getDecFn(rtelem, true, true)

	var rv0, rv9 reflect.Value
	rv0 = rv
	rvChanged := false

	// for j := 0; j < containerLenS; j++ ***REMOVED***
	var rvlen int
	if containerLenS > 0 ***REMOVED*** // hasLen
		if f.seq == seqTypeChan ***REMOVED***
			if rv.IsNil() ***REMOVED***
				rvlen, _ = decInferLen(containerLenS, f.d.h.MaxInitLen, int(rtelem0.Size()))
				rv.Set(reflect.MakeChan(ti.rt, rvlen))
			***REMOVED***
			// handle chan specially:
			for j := 0; j < containerLenS; j++ ***REMOVED***
				rv9 = reflect.New(rtelem0).Elem()
				slh.ElemContainerState(j)
				d.decodeValue(rv9, fn)
				rv.Send(rv9)
			***REMOVED***
		***REMOVED*** else ***REMOVED*** // slice or array
			var truncated bool         // says len of sequence is not same as expected number of elements
			numToRead := containerLenS // if truncated, reset numToRead

			rvcap := rv.Cap()
			rvlen = rv.Len()
			if containerLenS > rvcap ***REMOVED***
				if f.seq == seqTypeArray ***REMOVED***
					d.arrayCannotExpand(rvlen, containerLenS)
				***REMOVED*** else ***REMOVED***
					oldRvlenGtZero := rvlen > 0
					rvlen, truncated = decInferLen(containerLenS, f.d.h.MaxInitLen, int(rtelem0.Size()))
					if truncated ***REMOVED***
						if rvlen <= rvcap ***REMOVED***
							rv.SetLen(rvlen)
						***REMOVED*** else ***REMOVED***
							rv = reflect.MakeSlice(ti.rt, rvlen, rvlen)
							rvChanged = true
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						rv = reflect.MakeSlice(ti.rt, rvlen, rvlen)
						rvChanged = true
					***REMOVED***
					if rvChanged && oldRvlenGtZero && !isImmutableKind(rtelem0.Kind()) ***REMOVED***
						reflect.Copy(rv, rv0) // only copy up to length NOT cap i.e. rv0.Slice(0, rvcap)
					***REMOVED***
					rvcap = rvlen
				***REMOVED***
				numToRead = rvlen
			***REMOVED*** else if containerLenS != rvlen ***REMOVED***
				if f.seq == seqTypeSlice ***REMOVED***
					rv.SetLen(containerLenS)
					rvlen = containerLenS
				***REMOVED***
			***REMOVED***
			j := 0
			// we read up to the numToRead
			for ; j < numToRead; j++ ***REMOVED***
				slh.ElemContainerState(j)
				d.decodeValue(rv.Index(j), fn)
			***REMOVED***

			// if slice, expand and read up to containerLenS (or EOF) iff truncated
			// if array, swallow all the rest.

			if f.seq == seqTypeArray ***REMOVED***
				for ; j < containerLenS; j++ ***REMOVED***
					slh.ElemContainerState(j)
					d.swallow()
				***REMOVED***
			***REMOVED*** else if truncated ***REMOVED*** // slice was truncated, as chan NOT in this block
				for ; j < containerLenS; j++ ***REMOVED***
					rv = expandSliceValue(rv, 1)
					rv9 = rv.Index(j)
					if resetSliceElemToZeroValue ***REMOVED***
						rv9.Set(reflect.Zero(rtelem0))
					***REMOVED***
					slh.ElemContainerState(j)
					d.decodeValue(rv9, fn)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		rvlen = rv.Len()
		j := 0
		for ; !dd.CheckBreak(); j++ ***REMOVED***
			if f.seq == seqTypeChan ***REMOVED***
				slh.ElemContainerState(j)
				rv9 = reflect.New(rtelem0).Elem()
				d.decodeValue(rv9, fn)
				rv.Send(rv9)
			***REMOVED*** else ***REMOVED***
				// if indefinite, etc, then expand the slice if necessary
				var decodeIntoBlank bool
				if j >= rvlen ***REMOVED***
					if f.seq == seqTypeArray ***REMOVED***
						d.arrayCannotExpand(rvlen, j+1)
						decodeIntoBlank = true
					***REMOVED*** else ***REMOVED*** // if f.seq == seqTypeSlice
						// rv = reflect.Append(rv, reflect.Zero(rtelem0)) // uses append logic, plus varargs
						rv = expandSliceValue(rv, 1)
						rv9 = rv.Index(j)
						// rv.Index(rv.Len() - 1).Set(reflect.Zero(rtelem0))
						if resetSliceElemToZeroValue ***REMOVED***
							rv9.Set(reflect.Zero(rtelem0))
						***REMOVED***
						rvlen++
						rvChanged = true
					***REMOVED***
				***REMOVED*** else ***REMOVED*** // slice or array
					rv9 = rv.Index(j)
				***REMOVED***
				slh.ElemContainerState(j)
				if decodeIntoBlank ***REMOVED***
					d.swallow()
				***REMOVED*** else ***REMOVED*** // seqTypeSlice
					d.decodeValue(rv9, fn)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if f.seq == seqTypeSlice ***REMOVED***
			if j < rvlen ***REMOVED***
				rv.SetLen(j)
			***REMOVED*** else if j == 0 && rv.IsNil() ***REMOVED***
				rv = reflect.MakeSlice(ti.rt, 0, 0)
				rvChanged = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	slh.End()

	if rvChanged ***REMOVED***
		rv0.Set(rv)
	***REMOVED***
***REMOVED***

func (f *decFnInfo) kArray(rv reflect.Value) ***REMOVED***
	// f.d.decodeValue(rv.Slice(0, rv.Len()))
	f.kSlice(rv.Slice(0, rv.Len()))
***REMOVED***

func (f *decFnInfo) kMap(rv reflect.Value) ***REMOVED***
	d := f.d
	dd := d.d
	containerLen := dd.ReadMapStart()
	cr := d.cr
	ti := f.ti
	if rv.IsNil() ***REMOVED***
		rv.Set(reflect.MakeMap(ti.rt))
	***REMOVED***

	if containerLen == 0 ***REMOVED***
		if cr != nil ***REMOVED***
			cr.sendContainerState(containerMapEnd)
		***REMOVED***
		return
	***REMOVED***

	ktype, vtype := ti.rt.Key(), ti.rt.Elem()
	ktypeId := reflect.ValueOf(ktype).Pointer()
	vtypeKind := vtype.Kind()
	var keyFn, valFn *decFn
	var xtyp reflect.Type
	for xtyp = ktype; xtyp.Kind() == reflect.Ptr; xtyp = xtyp.Elem() ***REMOVED***
	***REMOVED***
	keyFn = d.getDecFn(xtyp, true, true)
	for xtyp = vtype; xtyp.Kind() == reflect.Ptr; xtyp = xtyp.Elem() ***REMOVED***
	***REMOVED***
	valFn = d.getDecFn(xtyp, true, true)
	var mapGet, mapSet bool
	if !f.d.h.MapValueReset ***REMOVED***
		// if pointer, mapGet = true
		// if interface, mapGet = true if !DecodeNakedAlways (else false)
		// if builtin, mapGet = false
		// else mapGet = true
		if vtypeKind == reflect.Ptr ***REMOVED***
			mapGet = true
		***REMOVED*** else if vtypeKind == reflect.Interface ***REMOVED***
			if !f.d.h.InterfaceReset ***REMOVED***
				mapGet = true
			***REMOVED***
		***REMOVED*** else if !isImmutableKind(vtypeKind) ***REMOVED***
			mapGet = true
		***REMOVED***
	***REMOVED***

	var rvk, rvv, rvz reflect.Value

	// for j := 0; j < containerLen; j++ ***REMOVED***
	if containerLen > 0 ***REMOVED***
		for j := 0; j < containerLen; j++ ***REMOVED***
			rvk = reflect.New(ktype).Elem()
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapKey)
			***REMOVED***
			d.decodeValue(rvk, keyFn)

			// special case if a byte array.
			if ktypeId == intfTypId ***REMOVED***
				rvk = rvk.Elem()
				if rvk.Type() == uint8SliceTyp ***REMOVED***
					rvk = reflect.ValueOf(d.string(rvk.Bytes()))
				***REMOVED***
			***REMOVED***
			mapSet = true // set to false if u do a get, and its a pointer, and exists
			if mapGet ***REMOVED***
				rvv = rv.MapIndex(rvk)
				if rvv.IsValid() ***REMOVED***
					if vtypeKind == reflect.Ptr ***REMOVED***
						mapSet = false
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if rvz.IsValid() ***REMOVED***
						rvz.Set(reflect.Zero(vtype))
					***REMOVED*** else ***REMOVED***
						rvz = reflect.New(vtype).Elem()
					***REMOVED***
					rvv = rvz
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if rvz.IsValid() ***REMOVED***
					rvz.Set(reflect.Zero(vtype))
				***REMOVED*** else ***REMOVED***
					rvz = reflect.New(vtype).Elem()
				***REMOVED***
				rvv = rvz
			***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapValue)
			***REMOVED***
			d.decodeValue(rvv, valFn)
			if mapSet ***REMOVED***
				rv.SetMapIndex(rvk, rvv)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for j := 0; !dd.CheckBreak(); j++ ***REMOVED***
			rvk = reflect.New(ktype).Elem()
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapKey)
			***REMOVED***
			d.decodeValue(rvk, keyFn)

			// special case if a byte array.
			if ktypeId == intfTypId ***REMOVED***
				rvk = rvk.Elem()
				if rvk.Type() == uint8SliceTyp ***REMOVED***
					rvk = reflect.ValueOf(d.string(rvk.Bytes()))
				***REMOVED***
			***REMOVED***
			mapSet = true // set to false if u do a get, and its a pointer, and exists
			if mapGet ***REMOVED***
				rvv = rv.MapIndex(rvk)
				if rvv.IsValid() ***REMOVED***
					if vtypeKind == reflect.Ptr ***REMOVED***
						mapSet = false
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if rvz.IsValid() ***REMOVED***
						rvz.Set(reflect.Zero(vtype))
					***REMOVED*** else ***REMOVED***
						rvz = reflect.New(vtype).Elem()
					***REMOVED***
					rvv = rvz
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if rvz.IsValid() ***REMOVED***
					rvz.Set(reflect.Zero(vtype))
				***REMOVED*** else ***REMOVED***
					rvz = reflect.New(vtype).Elem()
				***REMOVED***
				rvv = rvz
			***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapValue)
			***REMOVED***
			d.decodeValue(rvv, valFn)
			if mapSet ***REMOVED***
				rv.SetMapIndex(rvk, rvv)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if cr != nil ***REMOVED***
		cr.sendContainerState(containerMapEnd)
	***REMOVED***
***REMOVED***

type decRtidFn struct ***REMOVED***
	rtid uintptr
	fn   decFn
***REMOVED***

// decNaked is used to keep track of the primitives decoded.
// Without it, we would have to decode each primitive and wrap it
// in an interface***REMOVED******REMOVED***, causing an allocation.
// In this model, the primitives are decoded in a "pseudo-atomic" fashion,
// so we can rest assured that no other decoding happens while these
// primitives are being decoded.
//
// maps and arrays are not handled by this mechanism.
// However, RawExt is, and we accommodate for extensions that decode
// RawExt from DecodeNaked, but need to decode the value subsequently.
// kInterfaceNaked and swallow, which call DecodeNaked, handle this caveat.
//
// However, decNaked also keeps some arrays of default maps and slices
// used in DecodeNaked. This way, we can get a pointer to it
// without causing a new heap allocation.
//
// kInterfaceNaked will ensure that there is no allocation for the common
// uses.
type decNaked struct ***REMOVED***
	// r RawExt // used for RawExt, uint, []byte.
	u uint64
	i int64
	f float64
	l []byte
	s string
	t time.Time
	b bool
	v valueType

	// stacks for reducing allocation
	is []interface***REMOVED******REMOVED***
	ms []map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
	ns []map[string]interface***REMOVED******REMOVED***
	ss [][]interface***REMOVED******REMOVED***
	// rs []RawExt

	// keep arrays at the bottom? Chance is that they are not used much.
	ia [4]interface***REMOVED******REMOVED***
	ma [4]map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
	na [4]map[string]interface***REMOVED******REMOVED***
	sa [4][]interface***REMOVED******REMOVED***
	// ra [2]RawExt
***REMOVED***

func (n *decNaked) reset() ***REMOVED***
	if n.ss != nil ***REMOVED***
		n.ss = n.ss[:0]
	***REMOVED***
	if n.is != nil ***REMOVED***
		n.is = n.is[:0]
	***REMOVED***
	if n.ms != nil ***REMOVED***
		n.ms = n.ms[:0]
	***REMOVED***
	if n.ns != nil ***REMOVED***
		n.ns = n.ns[:0]
	***REMOVED***
***REMOVED***

// A Decoder reads and decodes an object from an input stream in the codec format.
type Decoder struct ***REMOVED***
	// hopefully, reduce derefencing cost by laying the decReader inside the Decoder.
	// Try to put things that go together to fit within a cache line (8 words).

	d decDriver
	// NOTE: Decoder shouldn't call it's read methods,
	// as the handler MAY need to do some coordination.
	r decReader
	// sa [initCollectionCap]decRtidFn
	h  *BasicHandle
	hh Handle

	be    bool // is binary encoding
	bytes bool // is bytes reader
	js    bool // is json handle

	rb bytesDecReader
	ri ioDecReader
	cr containerStateRecv

	s []decRtidFn
	f map[uintptr]*decFn

	// _  uintptr // for alignment purposes, so next one starts from a cache line

	// cache the mapTypeId and sliceTypeId for faster comparisons
	mtid uintptr
	stid uintptr

	n  decNaked
	b  [scratchByteArrayLen]byte
	is map[string]string // used for interning strings
***REMOVED***

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
//
// For efficiency, Users are encouraged to pass in a memory buffered reader
// (eg bufio.Reader, bytes.Buffer).
func NewDecoder(r io.Reader, h Handle) *Decoder ***REMOVED***
	d := newDecoder(h)
	d.Reset(r)
	return d
***REMOVED***

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) *Decoder ***REMOVED***
	d := newDecoder(h)
	d.ResetBytes(in)
	return d
***REMOVED***

func newDecoder(h Handle) *Decoder ***REMOVED***
	d := &Decoder***REMOVED***hh: h, h: h.getBasicHandle(), be: h.isBinary()***REMOVED***
	n := &d.n
	// n.rs = n.ra[:0]
	n.ms = n.ma[:0]
	n.is = n.ia[:0]
	n.ns = n.na[:0]
	n.ss = n.sa[:0]
	_, d.js = h.(*JsonHandle)
	if d.h.InternString ***REMOVED***
		d.is = make(map[string]string, 32)
	***REMOVED***
	d.d = h.newDecDriver(d)
	d.cr, _ = d.d.(containerStateRecv)
	// d.d = h.newDecDriver(decReaderT***REMOVED***true, &d.rb, &d.ri***REMOVED***)
	return d
***REMOVED***

func (d *Decoder) resetCommon() ***REMOVED***
	d.n.reset()
	d.d.reset()
	// reset all things which were cached from the Handle,
	// but could be changed.
	d.mtid, d.stid = 0, 0
	if d.h.MapType != nil ***REMOVED***
		d.mtid = reflect.ValueOf(d.h.MapType).Pointer()
	***REMOVED***
	if d.h.SliceType != nil ***REMOVED***
		d.stid = reflect.ValueOf(d.h.SliceType).Pointer()
	***REMOVED***
***REMOVED***

func (d *Decoder) Reset(r io.Reader) ***REMOVED***
	d.ri.x = &d.b
	// d.s = d.sa[:0]
	d.ri.bs.r = r
	var ok bool
	d.ri.br, ok = r.(decReaderByteScanner)
	if !ok ***REMOVED***
		d.ri.br = &d.ri.bs
	***REMOVED***
	d.r = &d.ri
	d.resetCommon()
***REMOVED***

func (d *Decoder) ResetBytes(in []byte) ***REMOVED***
	// d.s = d.sa[:0]
	d.bytes = true
	d.rb.reset(in)
	d.r = &d.rb
	d.resetCommon()
***REMOVED***

// func (d *Decoder) sendContainerState(c containerState) ***REMOVED***
// 	if d.cr != nil ***REMOVED***
// 		d.cr.sendContainerState(c)
// 	***REMOVED***
// ***REMOVED***

// Decode decodes the stream from reader and stores the result in the
// value pointed to by v. v cannot be a nil pointer. v can also be
// a reflect.Value of a pointer.
//
// Note that a pointer to a nil interface is not a nil pointer.
// If you do not know what type of stream it is, pass in a pointer to a nil interface.
// We will decode and store a value in that nil interface.
//
// Sample usages:
//   // Decoding into a non-nil typed value
//   var f float32
//   err = codec.NewDecoder(r, handle).Decode(&f)
//
//   // Decoding into nil interface
//   var v interface***REMOVED******REMOVED***
//   dec := codec.NewDecoder(r, handle)
//   err = dec.Decode(&v)
//
// When decoding into a nil interface***REMOVED******REMOVED***, we will decode into an appropriate value based
// on the contents of the stream:
//   - Numbers are decoded as float64, int64 or uint64.
//   - Other values are decoded appropriately depending on the type:
//     bool, string, []byte, time.Time, etc
//   - Extensions are decoded as RawExt (if no ext function registered for the tag)
// Configurations exist on the Handle to override defaults
// (e.g. for MapType, SliceType and how to decode raw bytes).
//
// When decoding into a non-nil interface***REMOVED******REMOVED*** value, the mode of encoding is based on the
// type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryUnmarshaler, call its UnmarshalBinary(data []byte) error
//   - Else decode it based on its reflect.Kind
//
// There are some special rules when decoding into containers (slice/array/map/struct).
// Decode will typically use the stream contents to UPDATE the container.
//   - A map can be decoded from a stream map, by updating matching keys.
//   - A slice can be decoded from a stream array,
//     by updating the first n elements, where n is length of the stream.
//   - A slice can be decoded from a stream map, by decoding as if
//     it contains a sequence of key-value pairs.
//   - A struct can be decoded from a stream map, by updating matching fields.
//   - A struct can be decoded from a stream array,
//     by updating fields as they occur in the struct (by index).
//
// When decoding a stream map or array with length of 0 into a nil map or slice,
// we reset the destination map or slice to a zero-length value.
//
// However, when decoding a stream nil, we reset the destination container
// to its "zero" value (e.g. nil for slice/map, etc).
//
func (d *Decoder) Decode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer panicToErr(&err)
	d.decode(v)
	return
***REMOVED***

// this is not a smart swallow, as it allocates objects and does unnecessary work.
func (d *Decoder) swallowViaHammer() ***REMOVED***
	var blank interface***REMOVED******REMOVED***
	d.decodeValue(reflect.ValueOf(&blank).Elem(), nil)
***REMOVED***

func (d *Decoder) swallow() ***REMOVED***
	// smarter decode that just swallows the content
	dd := d.d
	if dd.TryDecodeAsNil() ***REMOVED***
		return
	***REMOVED***
	cr := d.cr
	switch dd.ContainerType() ***REMOVED***
	case valueTypeMap:
		containerLen := dd.ReadMapStart()
		clenGtEqualZero := containerLen >= 0
		for j := 0; ; j++ ***REMOVED***
			if clenGtEqualZero ***REMOVED***
				if j >= containerLen ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else if dd.CheckBreak() ***REMOVED***
				break
			***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapKey)
			***REMOVED***
			d.swallow()
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapValue)
			***REMOVED***
			d.swallow()
		***REMOVED***
		if cr != nil ***REMOVED***
			cr.sendContainerState(containerMapEnd)
		***REMOVED***
	case valueTypeArray:
		containerLenS := dd.ReadArrayStart()
		clenGtEqualZero := containerLenS >= 0
		for j := 0; ; j++ ***REMOVED***
			if clenGtEqualZero ***REMOVED***
				if j >= containerLenS ***REMOVED***
					break
				***REMOVED***
			***REMOVED*** else if dd.CheckBreak() ***REMOVED***
				break
			***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerArrayElem)
			***REMOVED***
			d.swallow()
		***REMOVED***
		if cr != nil ***REMOVED***
			cr.sendContainerState(containerArrayEnd)
		***REMOVED***
	case valueTypeBytes:
		dd.DecodeBytes(d.b[:], false, true)
	case valueTypeString:
		dd.DecodeBytes(d.b[:], true, true)
		// dd.DecodeStringAsBytes(d.b[:])
	default:
		// these are all primitives, which we can get from decodeNaked
		// if RawExt using Value, complete the processing.
		dd.DecodeNaked()
		if n := &d.n; n.v == valueTypeExt && n.l == nil ***REMOVED***
			l := len(n.is)
			n.is = append(n.is, nil)
			v2 := &n.is[l]
			d.decode(v2)
			n.is = n.is[:l]
		***REMOVED***
	***REMOVED***
***REMOVED***

// MustDecode is like Decode, but panics if unable to Decode.
// This provides insight to the code location that triggered the error.
func (d *Decoder) MustDecode(v interface***REMOVED******REMOVED***) ***REMOVED***
	d.decode(v)
***REMOVED***

func (d *Decoder) decode(iv interface***REMOVED******REMOVED***) ***REMOVED***
	// if ics, ok := iv.(Selfer); ok ***REMOVED***
	// 	ics.CodecDecodeSelf(d)
	// 	return
	// ***REMOVED***

	if d.d.TryDecodeAsNil() ***REMOVED***
		switch v := iv.(type) ***REMOVED***
		case nil:
		case *string:
			*v = ""
		case *bool:
			*v = false
		case *int:
			*v = 0
		case *int8:
			*v = 0
		case *int16:
			*v = 0
		case *int32:
			*v = 0
		case *int64:
			*v = 0
		case *uint:
			*v = 0
		case *uint8:
			*v = 0
		case *uint16:
			*v = 0
		case *uint32:
			*v = 0
		case *uint64:
			*v = 0
		case *float32:
			*v = 0
		case *float64:
			*v = 0
		case *[]uint8:
			*v = nil
		case *Raw:
			*v = nil
		case reflect.Value:
			if v.Kind() != reflect.Ptr || v.IsNil() ***REMOVED***
				d.errNotValidPtrValue(v)
			***REMOVED***
			// d.chkPtrValue(v)
			v = v.Elem()
			if v.IsValid() ***REMOVED***
				v.Set(reflect.Zero(v.Type()))
			***REMOVED***
		default:
			rv := reflect.ValueOf(iv)
			if rv.Kind() != reflect.Ptr || rv.IsNil() ***REMOVED***
				d.errNotValidPtrValue(rv)
			***REMOVED***
			// d.chkPtrValue(rv)
			rv = rv.Elem()
			if rv.IsValid() ***REMOVED***
				rv.Set(reflect.Zero(rv.Type()))
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	switch v := iv.(type) ***REMOVED***
	case nil:
		d.error(cannotDecodeIntoNilErr)
		return

	case Selfer:
		v.CodecDecodeSelf(d)

	case reflect.Value:
		if v.Kind() != reflect.Ptr || v.IsNil() ***REMOVED***
			d.errNotValidPtrValue(v)
		***REMOVED***
		// d.chkPtrValue(v)
		d.decodeValueNotNil(v.Elem(), nil)

	case *string:
		*v = d.d.DecodeString()
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(d.d.DecodeInt(intBitsize))
	case *int8:
		*v = int8(d.d.DecodeInt(8))
	case *int16:
		*v = int16(d.d.DecodeInt(16))
	case *int32:
		*v = int32(d.d.DecodeInt(32))
	case *int64:
		*v = d.d.DecodeInt(64)
	case *uint:
		*v = uint(d.d.DecodeUint(uintBitsize))
	case *uint8:
		*v = uint8(d.d.DecodeUint(8))
	case *uint16:
		*v = uint16(d.d.DecodeUint(16))
	case *uint32:
		*v = uint32(d.d.DecodeUint(32))
	case *uint64:
		*v = d.d.DecodeUint(64)
	case *float32:
		*v = float32(d.d.DecodeFloat(true))
	case *float64:
		*v = d.d.DecodeFloat(false)
	case *[]uint8:
		*v = d.d.DecodeBytes(*v, false, false)

	case *Raw:
		*v = d.raw()

	case *interface***REMOVED******REMOVED***:
		d.decodeValueNotNil(reflect.ValueOf(iv).Elem(), nil)

	default:
		if !fastpathDecodeTypeSwitch(iv, d) ***REMOVED***
			d.decodeI(iv, true, false, false, false)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *Decoder) preDecodeValue(rv reflect.Value, tryNil bool) (rv2 reflect.Value, proceed bool) ***REMOVED***
	if tryNil && d.d.TryDecodeAsNil() ***REMOVED***
		// No need to check if a ptr, recursively, to determine
		// whether to set value to nil.
		// Just always set value to its zero type.
		if rv.IsValid() ***REMOVED*** // rv.CanSet() // always settable, except it's invalid
			rv.Set(reflect.Zero(rv.Type()))
		***REMOVED***
		return
	***REMOVED***

	// If stream is not containing a nil value, then we can deref to the base
	// non-pointer value, and decode into that.
	for rv.Kind() == reflect.Ptr ***REMOVED***
		if rv.IsNil() ***REMOVED***
			rv.Set(reflect.New(rv.Type().Elem()))
		***REMOVED***
		rv = rv.Elem()
	***REMOVED***
	return rv, true
***REMOVED***

func (d *Decoder) decodeI(iv interface***REMOVED******REMOVED***, checkPtr, tryNil, checkFastpath, checkCodecSelfer bool) ***REMOVED***
	rv := reflect.ValueOf(iv)
	if checkPtr ***REMOVED***
		if rv.Kind() != reflect.Ptr || rv.IsNil() ***REMOVED***
			d.errNotValidPtrValue(rv)
		***REMOVED***
		// d.chkPtrValue(rv)
	***REMOVED***
	rv, proceed := d.preDecodeValue(rv, tryNil)
	if proceed ***REMOVED***
		fn := d.getDecFn(rv.Type(), checkFastpath, checkCodecSelfer)
		fn.f(&fn.i, rv)
	***REMOVED***
***REMOVED***

func (d *Decoder) decodeValue(rv reflect.Value, fn *decFn) ***REMOVED***
	if rv, proceed := d.preDecodeValue(rv, true); proceed ***REMOVED***
		if fn == nil ***REMOVED***
			fn = d.getDecFn(rv.Type(), true, true)
		***REMOVED***
		fn.f(&fn.i, rv)
	***REMOVED***
***REMOVED***

func (d *Decoder) decodeValueNotNil(rv reflect.Value, fn *decFn) ***REMOVED***
	if rv, proceed := d.preDecodeValue(rv, false); proceed ***REMOVED***
		if fn == nil ***REMOVED***
			fn = d.getDecFn(rv.Type(), true, true)
		***REMOVED***
		fn.f(&fn.i, rv)
	***REMOVED***
***REMOVED***

func (d *Decoder) getDecFn(rt reflect.Type, checkFastpath, checkCodecSelfer bool) (fn *decFn) ***REMOVED***
	rtid := reflect.ValueOf(rt).Pointer()

	// retrieve or register a focus'ed function for this type
	// to eliminate need to do the retrieval multiple times

	// if d.f == nil && d.s == nil ***REMOVED*** debugf("---->Creating new dec f map for type: %v\n", rt) ***REMOVED***
	var ok bool
	if useMapForCodecCache ***REMOVED***
		fn, ok = d.f[rtid]
	***REMOVED*** else ***REMOVED***
		for i := range d.s ***REMOVED***
			v := &(d.s[i])
			if v.rtid == rtid ***REMOVED***
				fn, ok = &(v.fn), true
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if ok ***REMOVED***
		return
	***REMOVED***

	if useMapForCodecCache ***REMOVED***
		if d.f == nil ***REMOVED***
			d.f = make(map[uintptr]*decFn, initCollectionCap)
		***REMOVED***
		fn = new(decFn)
		d.f[rtid] = fn
	***REMOVED*** else ***REMOVED***
		if d.s == nil ***REMOVED***
			d.s = make([]decRtidFn, 0, initCollectionCap)
		***REMOVED***
		d.s = append(d.s, decRtidFn***REMOVED***rtid: rtid***REMOVED***)
		fn = &(d.s[len(d.s)-1]).fn
	***REMOVED***

	// debugf("\tCreating new dec fn for type: %v\n", rt)
	ti := d.h.getTypeInfo(rtid, rt)
	fi := &(fn.i)
	fi.d = d
	fi.ti = ti

	// An extension can be registered for any type, regardless of the Kind
	// (e.g. type BitSet int64, type MyStruct ***REMOVED*** / * unexported fields * / ***REMOVED***, type X []int, etc.
	//
	// We can't check if it's an extension byte here first, because the user may have
	// registered a pointer or non-pointer type, meaning we may have to recurse first
	// before matching a mapped type, even though the extension byte is already detected.
	//
	// NOTE: if decoding into a nil interface***REMOVED******REMOVED***, we return a non-nil
	// value except even if the container registers a length of 0.
	if checkCodecSelfer && ti.cs ***REMOVED***
		fn.f = (*decFnInfo).selferUnmarshal
	***REMOVED*** else if rtid == rawExtTypId ***REMOVED***
		fn.f = (*decFnInfo).rawExt
	***REMOVED*** else if rtid == rawTypId ***REMOVED***
		fn.f = (*decFnInfo).raw
	***REMOVED*** else if d.d.IsBuiltinType(rtid) ***REMOVED***
		fn.f = (*decFnInfo).builtin
	***REMOVED*** else if xfFn := d.h.getExt(rtid); xfFn != nil ***REMOVED***
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.f = (*decFnInfo).ext
	***REMOVED*** else if supportMarshalInterfaces && d.be && ti.bunm ***REMOVED***
		fn.f = (*decFnInfo).binaryUnmarshal
	***REMOVED*** else if supportMarshalInterfaces && !d.be && d.js && ti.junm ***REMOVED***
		//If JSON, we should check JSONUnmarshal before textUnmarshal
		fn.f = (*decFnInfo).jsonUnmarshal
	***REMOVED*** else if supportMarshalInterfaces && !d.be && ti.tunm ***REMOVED***
		fn.f = (*decFnInfo).textUnmarshal
	***REMOVED*** else ***REMOVED***
		rk := rt.Kind()
		if fastpathEnabled && checkFastpath && (rk == reflect.Map || rk == reflect.Slice) ***REMOVED***
			if rt.PkgPath() == "" ***REMOVED***
				if idx := fastpathAV.index(rtid); idx != -1 ***REMOVED***
					fn.f = fastpathAV[idx].decfn
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// use mapping for underlying type if there
				ok = false
				var rtu reflect.Type
				if rk == reflect.Map ***REMOVED***
					rtu = reflect.MapOf(rt.Key(), rt.Elem())
				***REMOVED*** else ***REMOVED***
					rtu = reflect.SliceOf(rt.Elem())
				***REMOVED***
				rtuid := reflect.ValueOf(rtu).Pointer()
				if idx := fastpathAV.index(rtuid); idx != -1 ***REMOVED***
					xfnf := fastpathAV[idx].decfn
					xrt := fastpathAV[idx].rt
					fn.f = func(xf *decFnInfo, xrv reflect.Value) ***REMOVED***
						// xfnf(xf, xrv.Convert(xrt))
						xfnf(xf, xrv.Addr().Convert(reflect.PtrTo(xrt)).Elem())
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if fn.f == nil ***REMOVED***
			switch rk ***REMOVED***
			case reflect.String:
				fn.f = (*decFnInfo).kString
			case reflect.Bool:
				fn.f = (*decFnInfo).kBool
			case reflect.Int:
				fn.f = (*decFnInfo).kInt
			case reflect.Int64:
				fn.f = (*decFnInfo).kInt64
			case reflect.Int32:
				fn.f = (*decFnInfo).kInt32
			case reflect.Int8:
				fn.f = (*decFnInfo).kInt8
			case reflect.Int16:
				fn.f = (*decFnInfo).kInt16
			case reflect.Float32:
				fn.f = (*decFnInfo).kFloat32
			case reflect.Float64:
				fn.f = (*decFnInfo).kFloat64
			case reflect.Uint8:
				fn.f = (*decFnInfo).kUint8
			case reflect.Uint64:
				fn.f = (*decFnInfo).kUint64
			case reflect.Uint:
				fn.f = (*decFnInfo).kUint
			case reflect.Uint32:
				fn.f = (*decFnInfo).kUint32
			case reflect.Uint16:
				fn.f = (*decFnInfo).kUint16
				// case reflect.Ptr:
				// 	fn.f = (*decFnInfo).kPtr
			case reflect.Uintptr:
				fn.f = (*decFnInfo).kUintptr
			case reflect.Interface:
				fn.f = (*decFnInfo).kInterface
			case reflect.Struct:
				fn.f = (*decFnInfo).kStruct
			case reflect.Chan:
				fi.seq = seqTypeChan
				fn.f = (*decFnInfo).kSlice
			case reflect.Slice:
				fi.seq = seqTypeSlice
				fn.f = (*decFnInfo).kSlice
			case reflect.Array:
				fi.seq = seqTypeArray
				fn.f = (*decFnInfo).kArray
			case reflect.Map:
				fn.f = (*decFnInfo).kMap
			default:
				fn.f = (*decFnInfo).kErr
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func (d *Decoder) structFieldNotFound(index int, rvkencname string) ***REMOVED***
	// NOTE: rvkencname may be a stringView, so don't pass it to another function.
	if d.h.ErrorIfNoField ***REMOVED***
		if index >= 0 ***REMOVED***
			d.errorf("no matching struct field found when decoding stream array at index %v", index)
			return
		***REMOVED*** else if rvkencname != "" ***REMOVED***
			d.errorf("no matching struct field found when decoding stream map with key " + rvkencname)
			return
		***REMOVED***
	***REMOVED***
	d.swallow()
***REMOVED***

func (d *Decoder) arrayCannotExpand(sliceLen, streamLen int) ***REMOVED***
	if d.h.ErrorIfNoArrayExpand ***REMOVED***
		d.errorf("cannot expand array len during decode from %v to %v", sliceLen, streamLen)
	***REMOVED***
***REMOVED***

func (d *Decoder) chkPtrValue(rv reflect.Value) ***REMOVED***
	// We can only decode into a non-nil pointer
	if rv.Kind() == reflect.Ptr && !rv.IsNil() ***REMOVED***
		return
	***REMOVED***
	d.errNotValidPtrValue(rv)
***REMOVED***

func (d *Decoder) errNotValidPtrValue(rv reflect.Value) ***REMOVED***
	if !rv.IsValid() ***REMOVED***
		d.error(cannotDecodeIntoNilErr)
		return
	***REMOVED***
	if !rv.CanInterface() ***REMOVED***
		d.errorf("cannot decode into a value without an interface: %v", rv)
		return
	***REMOVED***
	rvi := rv.Interface()
	d.errorf("cannot decode into non-pointer or nil pointer. Got: %v, %T, %v", rv.Kind(), rvi, rvi)
***REMOVED***

func (d *Decoder) error(err error) ***REMOVED***
	panic(err)
***REMOVED***

func (d *Decoder) errorf(format string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	params2 := make([]interface***REMOVED******REMOVED***, len(params)+1)
	params2[0] = d.r.numread()
	copy(params2[1:], params)
	err := fmt.Errorf("[pos %d]: "+format, params2...)
	panic(err)
***REMOVED***

func (d *Decoder) string(v []byte) (s string) ***REMOVED***
	if d.is != nil ***REMOVED***
		s, ok := d.is[string(v)] // no allocation here.
		if !ok ***REMOVED***
			s = string(v)
			d.is[s] = s
		***REMOVED***
		return s
	***REMOVED***
	return string(v) // don't return stringView, as we need a real string here.
***REMOVED***

func (d *Decoder) intern(s string) ***REMOVED***
	if d.is != nil ***REMOVED***
		d.is[s] = s
	***REMOVED***
***REMOVED***

// nextValueBytes returns the next value in the stream as a set of bytes.
func (d *Decoder) nextValueBytes() []byte ***REMOVED***
	d.d.uncacheRead()
	d.r.track()
	d.swallow()
	return d.r.stopTrack()
***REMOVED***

func (d *Decoder) raw() []byte ***REMOVED***
	// ensure that this is not a view into the bytes
	// i.e. make new copy always.
	bs := d.nextValueBytes()
	bs2 := make([]byte, len(bs))
	copy(bs2, bs)
	return bs2
***REMOVED***

// --------------------------------------------------

// decSliceHelper assists when decoding into a slice, from a map or an array in the stream.
// A slice can be set from a map or array in stream. This supports the MapBySlice interface.
type decSliceHelper struct ***REMOVED***
	d *Decoder
	// ct valueType
	array bool
***REMOVED***

func (d *Decoder) decSliceHelperStart() (x decSliceHelper, clen int) ***REMOVED***
	dd := d.d
	ctyp := dd.ContainerType()
	if ctyp == valueTypeArray ***REMOVED***
		x.array = true
		clen = dd.ReadArrayStart()
	***REMOVED*** else if ctyp == valueTypeMap ***REMOVED***
		clen = dd.ReadMapStart() * 2
	***REMOVED*** else ***REMOVED***
		d.errorf("only encoded map or array can be decoded into a slice (%d)", ctyp)
	***REMOVED***
	// x.ct = ctyp
	x.d = d
	return
***REMOVED***

func (x decSliceHelper) End() ***REMOVED***
	cr := x.d.cr
	if cr == nil ***REMOVED***
		return
	***REMOVED***
	if x.array ***REMOVED***
		cr.sendContainerState(containerArrayEnd)
	***REMOVED*** else ***REMOVED***
		cr.sendContainerState(containerMapEnd)
	***REMOVED***
***REMOVED***

func (x decSliceHelper) ElemContainerState(index int) ***REMOVED***
	cr := x.d.cr
	if cr == nil ***REMOVED***
		return
	***REMOVED***
	if x.array ***REMOVED***
		cr.sendContainerState(containerArrayElem)
	***REMOVED*** else ***REMOVED***
		if index%2 == 0 ***REMOVED***
			cr.sendContainerState(containerMapKey)
		***REMOVED*** else ***REMOVED***
			cr.sendContainerState(containerMapValue)
		***REMOVED***
	***REMOVED***
***REMOVED***

func decByteSlice(r decReader, clen, maxInitLen int, bs []byte) (bsOut []byte) ***REMOVED***
	if clen == 0 ***REMOVED***
		return zeroByteSlice
	***REMOVED***
	if len(bs) == clen ***REMOVED***
		bsOut = bs
		r.readb(bsOut)
	***REMOVED*** else if cap(bs) >= clen ***REMOVED***
		bsOut = bs[:clen]
		r.readb(bsOut)
	***REMOVED*** else ***REMOVED***
		// bsOut = make([]byte, clen)
		len2, _ := decInferLen(clen, maxInitLen, 1)
		bsOut = make([]byte, len2)
		r.readb(bsOut)
		for len2 < clen ***REMOVED***
			len3, _ := decInferLen(clen-len2, maxInitLen, 1)
			// fmt.Printf(">>>>> TESTING: in loop: clen: %v, maxInitLen: %v, len2: %v, len3: %v\n", clen, maxInitLen, len2, len3)
			bs3 := bsOut
			bsOut = make([]byte, len2+len3)
			copy(bsOut, bs3)
			r.readb(bsOut[len2:])
			len2 += len3
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func detachZeroCopyBytes(isBytesReader bool, dest []byte, in []byte) (out []byte) ***REMOVED***
	if xlen := len(in); xlen > 0 ***REMOVED***
		if isBytesReader || xlen <= scratchByteArrayLen ***REMOVED***
			if cap(dest) >= xlen ***REMOVED***
				out = dest[:xlen]
			***REMOVED*** else ***REMOVED***
				out = make([]byte, xlen)
			***REMOVED***
			copy(out, in)
			return
		***REMOVED***
	***REMOVED***
	return in
***REMOVED***

// decInferLen will infer a sensible length, given the following:
//    - clen: length wanted.
//    - maxlen: max length to be returned.
//      if <= 0, it is unset, and we infer it based on the unit size
//    - unit: number of bytes for each element of the collection
func decInferLen(clen, maxlen, unit int) (rvlen int, truncated bool) ***REMOVED***
	// handle when maxlen is not set i.e. <= 0
	if clen <= 0 ***REMOVED***
		return
	***REMOVED***
	if maxlen <= 0 ***REMOVED***
		// no maxlen defined. Use maximum of 256K memory, with a floor of 4K items.
		// maxlen = 256 * 1024 / unit
		// if maxlen < (4 * 1024) ***REMOVED***
		// 	maxlen = 4 * 1024
		// ***REMOVED***
		if unit < (256 / 4) ***REMOVED***
			maxlen = 256 * 1024 / unit
		***REMOVED*** else ***REMOVED***
			maxlen = 4 * 1024
		***REMOVED***
	***REMOVED***
	if clen > maxlen ***REMOVED***
		rvlen = maxlen
		truncated = true
	***REMOVED*** else ***REMOVED***
		rvlen = clen
	***REMOVED***
	return
	// if clen <= 0 ***REMOVED***
	// 	rvlen = 0
	// ***REMOVED*** else if maxlen > 0 && clen > maxlen ***REMOVED***
	// 	rvlen = maxlen
	// 	truncated = true
	// ***REMOVED*** else ***REMOVED***
	// 	rvlen = clen
	// ***REMOVED***
	// return
***REMOVED***

// // implement overall decReader wrapping both, for possible use inline:
// type decReaderT struct ***REMOVED***
// 	bytes bool
// 	rb    *bytesDecReader
// 	ri    *ioDecReader
// ***REMOVED***
//
// // implement *Decoder as a decReader.
// // Using decReaderT (defined just above) caused performance degradation
// // possibly because of constant copying the value,
// // and some value->interface conversion causing allocation.
// func (d *Decoder) unreadn1() ***REMOVED***
// 	if d.bytes ***REMOVED***
// 		d.rb.unreadn1()
// 	***REMOVED*** else ***REMOVED***
// 		d.ri.unreadn1()
// 	***REMOVED***
// ***REMOVED***
// ... for other methods of decReader.
// Testing showed that performance improvement was negligible.
