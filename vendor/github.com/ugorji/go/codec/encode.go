// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"
	"fmt"
	"io"
	"reflect"
	"sort"
	"sync"
)

const (
	defEncByteBufSize = 1 << 6 // 4:16, 6:64, 8:256, 10:1024
)

// AsSymbolFlag defines what should be encoded as symbols.
type AsSymbolFlag uint8

const (
	// AsSymbolDefault is default.
	// Currently, this means only encode struct field names as symbols.
	// The default is subject to change.
	AsSymbolDefault AsSymbolFlag = iota

	// AsSymbolAll means encode anything which could be a symbol as a symbol.
	AsSymbolAll = 0xfe

	// AsSymbolNone means do not encode anything as a symbol.
	AsSymbolNone = 1 << iota

	// AsSymbolMapStringKeys means encode keys in map[string]XXX as symbols.
	AsSymbolMapStringKeysFlag

	// AsSymbolStructFieldName means encode struct field names as symbols.
	AsSymbolStructFieldNameFlag
)

// encWriter abstracts writing to a byte array or to an io.Writer.
type encWriter interface ***REMOVED***
	writeb([]byte)
	writestr(string)
	writen1(byte)
	writen2(byte, byte)
	atEndOfEncode()
***REMOVED***

// encDriver abstracts the actual codec (binc vs msgpack, etc)
type encDriver interface ***REMOVED***
	IsBuiltinType(rt uintptr) bool
	EncodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***)
	EncodeNil()
	EncodeInt(i int64)
	EncodeUint(i uint64)
	EncodeBool(b bool)
	EncodeFloat32(f float32)
	EncodeFloat64(f float64)
	// encodeExtPreamble(xtag byte, length int)
	EncodeRawExt(re *RawExt, e *Encoder)
	EncodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext, e *Encoder)
	EncodeArrayStart(length int)
	EncodeMapStart(length int)
	EncodeString(c charEncoding, v string)
	EncodeSymbol(v string)
	EncodeStringBytes(c charEncoding, v []byte)
	//TODO
	//encBignum(f *big.Int)
	//encStringRunes(c charEncoding, v []rune)

	reset()
***REMOVED***

type encDriverAsis interface ***REMOVED***
	EncodeAsis(v []byte)
***REMOVED***

type encNoSeparator struct***REMOVED******REMOVED***

func (_ encNoSeparator) EncodeEnd() ***REMOVED******REMOVED***

type ioEncWriterWriter interface ***REMOVED***
	WriteByte(c byte) error
	WriteString(s string) (n int, err error)
	Write(p []byte) (n int, err error)
***REMOVED***

type ioEncStringWriter interface ***REMOVED***
	WriteString(s string) (n int, err error)
***REMOVED***

type EncodeOptions struct ***REMOVED***
	// Encode a struct as an array, and not as a map
	StructToArray bool

	// Canonical representation means that encoding a value will always result in the same
	// sequence of bytes.
	//
	// This only affects maps, as the iteration order for maps is random.
	//
	// The implementation MAY use the natural sort order for the map keys if possible:
	//
	//     - If there is a natural sort order (ie for number, bool, string or []byte keys),
	//       then the map keys are first sorted in natural order and then written
	//       with corresponding map values to the strema.
	//     - If there is no natural sort order, then the map keys will first be
	//       encoded into []byte, and then sorted,
	//       before writing the sorted keys and the corresponding map values to the stream.
	//
	Canonical bool

	// CheckCircularRef controls whether we check for circular references
	// and error fast during an encode.
	//
	// If enabled, an error is received if a pointer to a struct
	// references itself either directly or through one of its fields (iteratively).
	//
	// This is opt-in, as there may be a performance hit to checking circular references.
	CheckCircularRef bool

	// RecursiveEmptyCheck controls whether we descend into interfaces, structs and pointers
	// when checking if a value is empty.
	//
	// Note that this may make OmitEmpty more expensive, as it incurs a lot more reflect calls.
	RecursiveEmptyCheck bool

	// Raw controls whether we encode Raw values.
	// This is a "dangerous" option and must be explicitly set.
	// If set, we blindly encode Raw values as-is, without checking
	// if they are a correct representation of a value in that format.
	// If unset, we error out.
	Raw bool

	// AsSymbols defines what should be encoded as symbols.
	//
	// Encoding as symbols can reduce the encoded size significantly.
	//
	// However, during decoding, each string to be encoded as a symbol must
	// be checked to see if it has been seen before. Consequently, encoding time
	// will increase if using symbols, because string comparisons has a clear cost.
	//
	// Sample values:
	//   AsSymbolNone
	//   AsSymbolAll
	//   AsSymbolMapStringKeys
	//   AsSymbolMapStringKeysFlag | AsSymbolStructFieldNameFlag
	AsSymbols AsSymbolFlag
***REMOVED***

// ---------------------------------------------

type simpleIoEncWriterWriter struct ***REMOVED***
	w  io.Writer
	bw io.ByteWriter
	sw ioEncStringWriter
	bs [1]byte
***REMOVED***

func (o *simpleIoEncWriterWriter) WriteByte(c byte) (err error) ***REMOVED***
	if o.bw != nil ***REMOVED***
		return o.bw.WriteByte(c)
	***REMOVED***
	// _, err = o.w.Write([]byte***REMOVED***c***REMOVED***)
	o.bs[0] = c
	_, err = o.w.Write(o.bs[:])
	return
***REMOVED***

func (o *simpleIoEncWriterWriter) WriteString(s string) (n int, err error) ***REMOVED***
	if o.sw != nil ***REMOVED***
		return o.sw.WriteString(s)
	***REMOVED***
	// return o.w.Write([]byte(s))
	return o.w.Write(bytesView(s))
***REMOVED***

func (o *simpleIoEncWriterWriter) Write(p []byte) (n int, err error) ***REMOVED***
	return o.w.Write(p)
***REMOVED***

// ----------------------------------------

// ioEncWriter implements encWriter and can write to an io.Writer implementation
type ioEncWriter struct ***REMOVED***
	w ioEncWriterWriter
	s simpleIoEncWriterWriter
	// x [8]byte // temp byte array re-used internally for efficiency
***REMOVED***

func (z *ioEncWriter) writeb(bs []byte) ***REMOVED***
	if len(bs) == 0 ***REMOVED***
		return
	***REMOVED***
	n, err := z.w.Write(bs)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if n != len(bs) ***REMOVED***
		panic(fmt.Errorf("incorrect num bytes written. Expecting: %v, Wrote: %v", len(bs), n))
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writestr(s string) ***REMOVED***
	n, err := z.w.WriteString(s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if n != len(s) ***REMOVED***
		panic(fmt.Errorf("incorrect num bytes written. Expecting: %v, Wrote: %v", len(s), n))
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writen1(b byte) ***REMOVED***
	if err := z.w.WriteByte(b); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writen2(b1 byte, b2 byte) ***REMOVED***
	z.writen1(b1)
	z.writen1(b2)
***REMOVED***

func (z *ioEncWriter) atEndOfEncode() ***REMOVED******REMOVED***

// ----------------------------------------

// bytesEncWriter implements encWriter and can write to an byte slice.
// It is used by Marshal function.
type bytesEncWriter struct ***REMOVED***
	b   []byte
	c   int     // cursor
	out *[]byte // write out on atEndOfEncode
***REMOVED***

func (z *bytesEncWriter) writeb(s []byte) ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return
	***REMOVED***
	oc, a := z.growNoAlloc(len(s))
	if a ***REMOVED***
		z.growAlloc(len(s), oc)
	***REMOVED***
	copy(z.b[oc:], s)
***REMOVED***

func (z *bytesEncWriter) writestr(s string) ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return
	***REMOVED***
	oc, a := z.growNoAlloc(len(s))
	if a ***REMOVED***
		z.growAlloc(len(s), oc)
	***REMOVED***
	copy(z.b[oc:], s)
***REMOVED***

func (z *bytesEncWriter) writen1(b1 byte) ***REMOVED***
	oc, a := z.growNoAlloc(1)
	if a ***REMOVED***
		z.growAlloc(1, oc)
	***REMOVED***
	z.b[oc] = b1
***REMOVED***

func (z *bytesEncWriter) writen2(b1 byte, b2 byte) ***REMOVED***
	oc, a := z.growNoAlloc(2)
	if a ***REMOVED***
		z.growAlloc(2, oc)
	***REMOVED***
	z.b[oc+1] = b2
	z.b[oc] = b1
***REMOVED***

func (z *bytesEncWriter) atEndOfEncode() ***REMOVED***
	*(z.out) = z.b[:z.c]
***REMOVED***

// have a growNoalloc(n int), which can be inlined.
// if allocation is needed, then call growAlloc(n int)

func (z *bytesEncWriter) growNoAlloc(n int) (oldcursor int, allocNeeded bool) ***REMOVED***
	oldcursor = z.c
	z.c = z.c + n
	if z.c > len(z.b) ***REMOVED***
		if z.c > cap(z.b) ***REMOVED***
			allocNeeded = true
		***REMOVED*** else ***REMOVED***
			z.b = z.b[:cap(z.b)]
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (z *bytesEncWriter) growAlloc(n int, oldcursor int) ***REMOVED***
	// appendslice logic (if cap < 1024, *2, else *1.25): more expensive. many copy calls.
	// bytes.Buffer model (2*cap + n): much better
	// bs := make([]byte, 2*cap(z.b)+n)
	bs := make([]byte, growCap(cap(z.b), 1, n))
	copy(bs, z.b[:oldcursor])
	z.b = bs
***REMOVED***

// ---------------------------------------------

type encFnInfo struct ***REMOVED***
	e     *Encoder
	ti    *typeInfo
	xfFn  Ext
	xfTag uint64
	seq   seqType
***REMOVED***

func (f *encFnInfo) builtin(rv reflect.Value) ***REMOVED***
	f.e.e.EncodeBuiltin(f.ti.rtid, rv.Interface())
***REMOVED***

func (f *encFnInfo) raw(rv reflect.Value) ***REMOVED***
	f.e.raw(rv.Interface().(Raw))
***REMOVED***

func (f *encFnInfo) rawExt(rv reflect.Value) ***REMOVED***
	// rev := rv.Interface().(RawExt)
	// f.e.e.EncodeRawExt(&rev, f.e)
	var re *RawExt
	if rv.CanAddr() ***REMOVED***
		re = rv.Addr().Interface().(*RawExt)
	***REMOVED*** else ***REMOVED***
		rev := rv.Interface().(RawExt)
		re = &rev
	***REMOVED***
	f.e.e.EncodeRawExt(re, f.e)
***REMOVED***

func (f *encFnInfo) ext(rv reflect.Value) ***REMOVED***
	// if this is a struct|array and it was addressable, then pass the address directly (not the value)
	if k := rv.Kind(); (k == reflect.Struct || k == reflect.Array) && rv.CanAddr() ***REMOVED***
		rv = rv.Addr()
	***REMOVED***
	f.e.e.EncodeExt(rv.Interface(), f.xfTag, f.xfFn, f.e)
***REMOVED***

func (f *encFnInfo) getValueForMarshalInterface(rv reflect.Value, indir int8) (v interface***REMOVED******REMOVED***, proceed bool) ***REMOVED***
	if indir == 0 ***REMOVED***
		v = rv.Interface()
	***REMOVED*** else if indir == -1 ***REMOVED***
		// If a non-pointer was passed to Encode(), then that value is not addressable.
		// Take addr if addressable, else copy value to an addressable value.
		if rv.CanAddr() ***REMOVED***
			v = rv.Addr().Interface()
		***REMOVED*** else ***REMOVED***
			rv2 := reflect.New(rv.Type())
			rv2.Elem().Set(rv)
			v = rv2.Interface()
			// fmt.Printf("rv.Type: %v, rv2.Type: %v, v: %v\n", rv.Type(), rv2.Type(), v)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for j := int8(0); j < indir; j++ ***REMOVED***
			if rv.IsNil() ***REMOVED***
				f.e.e.EncodeNil()
				return
			***REMOVED***
			rv = rv.Elem()
		***REMOVED***
		v = rv.Interface()
	***REMOVED***
	return v, true
***REMOVED***

func (f *encFnInfo) selferMarshal(rv reflect.Value) ***REMOVED***
	if v, proceed := f.getValueForMarshalInterface(rv, f.ti.csIndir); proceed ***REMOVED***
		v.(Selfer).CodecEncodeSelf(f.e)
	***REMOVED***
***REMOVED***

func (f *encFnInfo) binaryMarshal(rv reflect.Value) ***REMOVED***
	if v, proceed := f.getValueForMarshalInterface(rv, f.ti.bmIndir); proceed ***REMOVED***
		bs, fnerr := v.(encoding.BinaryMarshaler).MarshalBinary()
		f.e.marshal(bs, fnerr, false, c_RAW)
	***REMOVED***
***REMOVED***

func (f *encFnInfo) textMarshal(rv reflect.Value) ***REMOVED***
	if v, proceed := f.getValueForMarshalInterface(rv, f.ti.tmIndir); proceed ***REMOVED***
		// debugf(">>>> encoding.TextMarshaler: %T", rv.Interface())
		bs, fnerr := v.(encoding.TextMarshaler).MarshalText()
		f.e.marshal(bs, fnerr, false, c_UTF8)
	***REMOVED***
***REMOVED***

func (f *encFnInfo) jsonMarshal(rv reflect.Value) ***REMOVED***
	if v, proceed := f.getValueForMarshalInterface(rv, f.ti.jmIndir); proceed ***REMOVED***
		bs, fnerr := v.(jsonMarshaler).MarshalJSON()
		f.e.marshal(bs, fnerr, true, c_UTF8)
	***REMOVED***
***REMOVED***

func (f *encFnInfo) kBool(rv reflect.Value) ***REMOVED***
	f.e.e.EncodeBool(rv.Bool())
***REMOVED***

func (f *encFnInfo) kString(rv reflect.Value) ***REMOVED***
	f.e.e.EncodeString(c_UTF8, rv.String())
***REMOVED***

func (f *encFnInfo) kFloat64(rv reflect.Value) ***REMOVED***
	f.e.e.EncodeFloat64(rv.Float())
***REMOVED***

func (f *encFnInfo) kFloat32(rv reflect.Value) ***REMOVED***
	f.e.e.EncodeFloat32(float32(rv.Float()))
***REMOVED***

func (f *encFnInfo) kInt(rv reflect.Value) ***REMOVED***
	f.e.e.EncodeInt(rv.Int())
***REMOVED***

func (f *encFnInfo) kUint(rv reflect.Value) ***REMOVED***
	f.e.e.EncodeUint(rv.Uint())
***REMOVED***

func (f *encFnInfo) kInvalid(rv reflect.Value) ***REMOVED***
	f.e.e.EncodeNil()
***REMOVED***

func (f *encFnInfo) kErr(rv reflect.Value) ***REMOVED***
	f.e.errorf("unsupported kind %s, for %#v", rv.Kind(), rv)
***REMOVED***

func (f *encFnInfo) kSlice(rv reflect.Value) ***REMOVED***
	ti := f.ti
	// array may be non-addressable, so we have to manage with care
	//   (don't call rv.Bytes, rv.Slice, etc).
	// E.g. type struct S***REMOVED***B [2]byte***REMOVED***;
	//   Encode(S***REMOVED******REMOVED***) will bomb on "panic: slice of unaddressable array".
	e := f.e
	if f.seq != seqTypeArray ***REMOVED***
		if rv.IsNil() ***REMOVED***
			e.e.EncodeNil()
			return
		***REMOVED***
		// If in this method, then there was no extension function defined.
		// So it's okay to treat as []byte.
		if ti.rtid == uint8SliceTypId ***REMOVED***
			e.e.EncodeStringBytes(c_RAW, rv.Bytes())
			return
		***REMOVED***
	***REMOVED***
	cr := e.cr
	rtelem := ti.rt.Elem()
	l := rv.Len()
	if ti.rtid == uint8SliceTypId || rtelem.Kind() == reflect.Uint8 ***REMOVED***
		switch f.seq ***REMOVED***
		case seqTypeArray:
			// if l == 0 ***REMOVED*** e.e.encodeStringBytes(c_RAW, nil) ***REMOVED*** else
			if rv.CanAddr() ***REMOVED***
				e.e.EncodeStringBytes(c_RAW, rv.Slice(0, l).Bytes())
			***REMOVED*** else ***REMOVED***
				var bs []byte
				if l <= cap(e.b) ***REMOVED***
					bs = e.b[:l]
				***REMOVED*** else ***REMOVED***
					bs = make([]byte, l)
				***REMOVED***
				reflect.Copy(reflect.ValueOf(bs), rv)
				// TODO: Test that reflect.Copy works instead of manual one-by-one
				// for i := 0; i < l; i++ ***REMOVED***
				// 	bs[i] = byte(rv.Index(i).Uint())
				// ***REMOVED***
				e.e.EncodeStringBytes(c_RAW, bs)
			***REMOVED***
		case seqTypeSlice:
			e.e.EncodeStringBytes(c_RAW, rv.Bytes())
		case seqTypeChan:
			bs := e.b[:0]
			// do not use range, so that the number of elements encoded
			// does not change, and encoding does not hang waiting on someone to close chan.
			// for b := range rv.Interface().(<-chan byte) ***REMOVED***
			// 	bs = append(bs, b)
			// ***REMOVED***
			ch := rv.Interface().(<-chan byte)
			for i := 0; i < l; i++ ***REMOVED***
				bs = append(bs, <-ch)
			***REMOVED***
			e.e.EncodeStringBytes(c_RAW, bs)
		***REMOVED***
		return
	***REMOVED***

	if ti.mbs ***REMOVED***
		if l%2 == 1 ***REMOVED***
			e.errorf("mapBySlice requires even slice length, but got %v", l)
			return
		***REMOVED***
		e.e.EncodeMapStart(l / 2)
	***REMOVED*** else ***REMOVED***
		e.e.EncodeArrayStart(l)
	***REMOVED***

	if l > 0 ***REMOVED***
		for rtelem.Kind() == reflect.Ptr ***REMOVED***
			rtelem = rtelem.Elem()
		***REMOVED***
		// if kind is reflect.Interface, do not pre-determine the
		// encoding type, because preEncodeValue may break it down to
		// a concrete type and kInterface will bomb.
		var fn *encFn
		if rtelem.Kind() != reflect.Interface ***REMOVED***
			rtelemid := reflect.ValueOf(rtelem).Pointer()
			fn = e.getEncFn(rtelemid, rtelem, true, true)
		***REMOVED***
		// TODO: Consider perf implication of encoding odd index values as symbols if type is string
		for j := 0; j < l; j++ ***REMOVED***
			if cr != nil ***REMOVED***
				if ti.mbs ***REMOVED***
					if j%2 == 0 ***REMOVED***
						cr.sendContainerState(containerMapKey)
					***REMOVED*** else ***REMOVED***
						cr.sendContainerState(containerMapValue)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					cr.sendContainerState(containerArrayElem)
				***REMOVED***
			***REMOVED***
			if f.seq == seqTypeChan ***REMOVED***
				if rv2, ok2 := rv.Recv(); ok2 ***REMOVED***
					e.encodeValue(rv2, fn)
				***REMOVED*** else ***REMOVED***
					e.encode(nil) // WE HAVE TO DO SOMETHING, so nil if nothing received.
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				e.encodeValue(rv.Index(j), fn)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if cr != nil ***REMOVED***
		if ti.mbs ***REMOVED***
			cr.sendContainerState(containerMapEnd)
		***REMOVED*** else ***REMOVED***
			cr.sendContainerState(containerArrayEnd)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *encFnInfo) kStruct(rv reflect.Value) ***REMOVED***
	fti := f.ti
	e := f.e
	cr := e.cr
	tisfi := fti.sfip
	toMap := !(fti.toArray || e.h.StructToArray)
	newlen := len(fti.sfi)

	// Use sync.Pool to reduce allocating slices unnecessarily.
	// The cost of sync.Pool is less than the cost of new allocation.
	pool, poolv, fkvs := encStructPoolGet(newlen)

	// if toMap, use the sorted array. If toArray, use unsorted array (to match sequence in struct)
	if toMap ***REMOVED***
		tisfi = fti.sfi
	***REMOVED***
	newlen = 0
	var kv stringRv
	recur := e.h.RecursiveEmptyCheck
	for _, si := range tisfi ***REMOVED***
		kv.r = si.field(rv, false)
		if toMap ***REMOVED***
			if si.omitEmpty && isEmptyValue(kv.r, recur, recur) ***REMOVED***
				continue
			***REMOVED***
			kv.v = si.encName
		***REMOVED*** else ***REMOVED***
			// use the zero value.
			// if a reference or struct, set to nil (so you do not output too much)
			if si.omitEmpty && isEmptyValue(kv.r, recur, recur) ***REMOVED***
				switch kv.r.Kind() ***REMOVED***
				case reflect.Struct, reflect.Interface, reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice:
					kv.r = reflect.Value***REMOVED******REMOVED*** //encode as nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
		fkvs[newlen] = kv
		newlen++
	***REMOVED***

	// debugf(">>>> kStruct: newlen: %v", newlen)
	// sep := !e.be
	ee := e.e //don't dereference every time

	if toMap ***REMOVED***
		ee.EncodeMapStart(newlen)
		// asSymbols := e.h.AsSymbols&AsSymbolStructFieldNameFlag != 0
		asSymbols := e.h.AsSymbols == AsSymbolDefault || e.h.AsSymbols&AsSymbolStructFieldNameFlag != 0
		for j := 0; j < newlen; j++ ***REMOVED***
			kv = fkvs[j]
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapKey)
			***REMOVED***
			if asSymbols ***REMOVED***
				ee.EncodeSymbol(kv.v)
			***REMOVED*** else ***REMOVED***
				ee.EncodeString(c_UTF8, kv.v)
			***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapValue)
			***REMOVED***
			e.encodeValue(kv.r, nil)
		***REMOVED***
		if cr != nil ***REMOVED***
			cr.sendContainerState(containerMapEnd)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		ee.EncodeArrayStart(newlen)
		for j := 0; j < newlen; j++ ***REMOVED***
			kv = fkvs[j]
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerArrayElem)
			***REMOVED***
			e.encodeValue(kv.r, nil)
		***REMOVED***
		if cr != nil ***REMOVED***
			cr.sendContainerState(containerArrayEnd)
		***REMOVED***
	***REMOVED***

	// do not use defer. Instead, use explicit pool return at end of function.
	// defer has a cost we are trying to avoid.
	// If there is a panic and these slices are not returned, it is ok.
	if pool != nil ***REMOVED***
		pool.Put(poolv)
	***REMOVED***
***REMOVED***

// func (f *encFnInfo) kPtr(rv reflect.Value) ***REMOVED***
// 	debugf(">>>>>>> ??? encode kPtr called - shouldn't get called")
// 	if rv.IsNil() ***REMOVED***
// 		f.e.e.encodeNil()
// 		return
// 	***REMOVED***
// 	f.e.encodeValue(rv.Elem())
// ***REMOVED***

// func (f *encFnInfo) kInterface(rv reflect.Value) ***REMOVED***
// 	println("kInterface called")
// 	debug.PrintStack()
// 	if rv.IsNil() ***REMOVED***
// 		f.e.e.EncodeNil()
// 		return
// 	***REMOVED***
// 	f.e.encodeValue(rv.Elem(), nil)
// ***REMOVED***

func (f *encFnInfo) kMap(rv reflect.Value) ***REMOVED***
	ee := f.e.e
	if rv.IsNil() ***REMOVED***
		ee.EncodeNil()
		return
	***REMOVED***

	l := rv.Len()
	ee.EncodeMapStart(l)
	e := f.e
	cr := e.cr
	if l == 0 ***REMOVED***
		if cr != nil ***REMOVED***
			cr.sendContainerState(containerMapEnd)
		***REMOVED***
		return
	***REMOVED***
	var asSymbols bool
	// determine the underlying key and val encFn's for the map.
	// This eliminates some work which is done for each loop iteration i.e.
	// rv.Type(), ref.ValueOf(rt).Pointer(), then check map/list for fn.
	//
	// However, if kind is reflect.Interface, do not pre-determine the
	// encoding type, because preEncodeValue may break it down to
	// a concrete type and kInterface will bomb.
	var keyFn, valFn *encFn
	ti := f.ti
	rtkey := ti.rt.Key()
	rtval := ti.rt.Elem()
	rtkeyid := reflect.ValueOf(rtkey).Pointer()
	// keyTypeIsString := f.ti.rt.Key().Kind() == reflect.String
	var keyTypeIsString = rtkeyid == stringTypId
	if keyTypeIsString ***REMOVED***
		asSymbols = e.h.AsSymbols&AsSymbolMapStringKeysFlag != 0
	***REMOVED*** else ***REMOVED***
		for rtkey.Kind() == reflect.Ptr ***REMOVED***
			rtkey = rtkey.Elem()
		***REMOVED***
		if rtkey.Kind() != reflect.Interface ***REMOVED***
			rtkeyid = reflect.ValueOf(rtkey).Pointer()
			keyFn = e.getEncFn(rtkeyid, rtkey, true, true)
		***REMOVED***
	***REMOVED***
	for rtval.Kind() == reflect.Ptr ***REMOVED***
		rtval = rtval.Elem()
	***REMOVED***
	if rtval.Kind() != reflect.Interface ***REMOVED***
		rtvalid := reflect.ValueOf(rtval).Pointer()
		valFn = e.getEncFn(rtvalid, rtval, true, true)
	***REMOVED***
	mks := rv.MapKeys()
	// for j, lmks := 0, len(mks); j < lmks; j++ ***REMOVED***

	if e.h.Canonical ***REMOVED***
		e.kMapCanonical(rtkeyid, rtkey, rv, mks, valFn, asSymbols)
	***REMOVED*** else ***REMOVED***
		for j := range mks ***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapKey)
			***REMOVED***
			if keyTypeIsString ***REMOVED***
				if asSymbols ***REMOVED***
					ee.EncodeSymbol(mks[j].String())
				***REMOVED*** else ***REMOVED***
					ee.EncodeString(c_UTF8, mks[j].String())
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				e.encodeValue(mks[j], keyFn)
			***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapValue)
			***REMOVED***
			e.encodeValue(rv.MapIndex(mks[j]), valFn)
		***REMOVED***
	***REMOVED***
	if cr != nil ***REMOVED***
		cr.sendContainerState(containerMapEnd)
	***REMOVED***
***REMOVED***

func (e *Encoder) kMapCanonical(rtkeyid uintptr, rtkey reflect.Type, rv reflect.Value, mks []reflect.Value, valFn *encFn, asSymbols bool) ***REMOVED***
	ee := e.e
	cr := e.cr
	// we previously did out-of-band if an extension was registered.
	// This is not necessary, as the natural kind is sufficient for ordering.

	if rtkeyid == uint8SliceTypId ***REMOVED***
		mksv := make([]bytesRv, len(mks))
		for i, k := range mks ***REMOVED***
			v := &mksv[i]
			v.r = k
			v.v = k.Bytes()
		***REMOVED***
		sort.Sort(bytesRvSlice(mksv))
		for i := range mksv ***REMOVED***
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapKey)
			***REMOVED***
			ee.EncodeStringBytes(c_RAW, mksv[i].v)
			if cr != nil ***REMOVED***
				cr.sendContainerState(containerMapValue)
			***REMOVED***
			e.encodeValue(rv.MapIndex(mksv[i].r), valFn)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		switch rtkey.Kind() ***REMOVED***
		case reflect.Bool:
			mksv := make([]boolRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksv[i]
				v.r = k
				v.v = k.Bool()
			***REMOVED***
			sort.Sort(boolRvSlice(mksv))
			for i := range mksv ***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				ee.EncodeBool(mksv[i].v)
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				e.encodeValue(rv.MapIndex(mksv[i].r), valFn)
			***REMOVED***
		case reflect.String:
			mksv := make([]stringRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksv[i]
				v.r = k
				v.v = k.String()
			***REMOVED***
			sort.Sort(stringRvSlice(mksv))
			for i := range mksv ***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				if asSymbols ***REMOVED***
					ee.EncodeSymbol(mksv[i].v)
				***REMOVED*** else ***REMOVED***
					ee.EncodeString(c_UTF8, mksv[i].v)
				***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				e.encodeValue(rv.MapIndex(mksv[i].r), valFn)
			***REMOVED***
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
			mksv := make([]uintRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksv[i]
				v.r = k
				v.v = k.Uint()
			***REMOVED***
			sort.Sort(uintRvSlice(mksv))
			for i := range mksv ***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				ee.EncodeUint(mksv[i].v)
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				e.encodeValue(rv.MapIndex(mksv[i].r), valFn)
			***REMOVED***
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			mksv := make([]intRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksv[i]
				v.r = k
				v.v = k.Int()
			***REMOVED***
			sort.Sort(intRvSlice(mksv))
			for i := range mksv ***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				ee.EncodeInt(mksv[i].v)
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				e.encodeValue(rv.MapIndex(mksv[i].r), valFn)
			***REMOVED***
		case reflect.Float32:
			mksv := make([]floatRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksv[i]
				v.r = k
				v.v = k.Float()
			***REMOVED***
			sort.Sort(floatRvSlice(mksv))
			for i := range mksv ***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				ee.EncodeFloat32(float32(mksv[i].v))
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				e.encodeValue(rv.MapIndex(mksv[i].r), valFn)
			***REMOVED***
		case reflect.Float64:
			mksv := make([]floatRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksv[i]
				v.r = k
				v.v = k.Float()
			***REMOVED***
			sort.Sort(floatRvSlice(mksv))
			for i := range mksv ***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				ee.EncodeFloat64(mksv[i].v)
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				e.encodeValue(rv.MapIndex(mksv[i].r), valFn)
			***REMOVED***
		default:
			// out-of-band
			// first encode each key to a []byte first, then sort them, then record
			var mksv []byte = make([]byte, 0, len(mks)*16) // temporary byte slice for the encoding
			e2 := NewEncoderBytes(&mksv, e.hh)
			mksbv := make([]bytesRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksbv[i]
				l := len(mksv)
				e2.MustEncode(k)
				v.r = k
				v.v = mksv[l:]
				// fmt.Printf(">>>>> %s\n", mksv[l:])
			***REMOVED***
			sort.Sort(bytesRvSlice(mksbv))
			for j := range mksbv ***REMOVED***
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapKey)
				***REMOVED***
				e.asis(mksbv[j].v)
				if cr != nil ***REMOVED***
					cr.sendContainerState(containerMapValue)
				***REMOVED***
				e.encodeValue(rv.MapIndex(mksbv[j].r), valFn)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// --------------------------------------------------

// encFn encapsulates the captured variables and the encode function.
// This way, we only do some calculations one times, and pass to the
// code block that should be called (encapsulated in a function)
// instead of executing the checks every time.
type encFn struct ***REMOVED***
	i encFnInfo
	f func(*encFnInfo, reflect.Value)
***REMOVED***

// --------------------------------------------------

type encRtidFn struct ***REMOVED***
	rtid uintptr
	fn   encFn
***REMOVED***

// An Encoder writes an object to an output stream in the codec format.
type Encoder struct ***REMOVED***
	// hopefully, reduce derefencing cost by laying the encWriter inside the Encoder
	e encDriver
	// NOTE: Encoder shouldn't call it's write methods,
	// as the handler MAY need to do some coordination.
	w  encWriter
	s  []encRtidFn
	ci set
	be bool // is binary encoding
	js bool // is json handle

	wi ioEncWriter
	wb bytesEncWriter

	h  *BasicHandle
	hh Handle

	cr containerStateRecv
	as encDriverAsis

	f map[uintptr]*encFn
	b [scratchByteArrayLen]byte
***REMOVED***

// NewEncoder returns an Encoder for encoding into an io.Writer.
//
// For efficiency, Users are encouraged to pass in a memory buffered writer
// (eg bufio.Writer, bytes.Buffer).
func NewEncoder(w io.Writer, h Handle) *Encoder ***REMOVED***
	e := newEncoder(h)
	e.Reset(w)
	return e
***REMOVED***

// NewEncoderBytes returns an encoder for encoding directly and efficiently
// into a byte slice, using zero-copying to temporary slices.
//
// It will potentially replace the output byte slice pointed to.
// After encoding, the out parameter contains the encoded contents.
func NewEncoderBytes(out *[]byte, h Handle) *Encoder ***REMOVED***
	e := newEncoder(h)
	e.ResetBytes(out)
	return e
***REMOVED***

func newEncoder(h Handle) *Encoder ***REMOVED***
	e := &Encoder***REMOVED***hh: h, h: h.getBasicHandle(), be: h.isBinary()***REMOVED***
	_, e.js = h.(*JsonHandle)
	e.e = h.newEncDriver(e)
	e.as, _ = e.e.(encDriverAsis)
	e.cr, _ = e.e.(containerStateRecv)
	return e
***REMOVED***

// Reset the Encoder with a new output stream.
//
// This accommodates using the state of the Encoder,
// where it has "cached" information about sub-engines.
func (e *Encoder) Reset(w io.Writer) ***REMOVED***
	ww, ok := w.(ioEncWriterWriter)
	if ok ***REMOVED***
		e.wi.w = ww
	***REMOVED*** else ***REMOVED***
		sww := &e.wi.s
		sww.w = w
		sww.bw, _ = w.(io.ByteWriter)
		sww.sw, _ = w.(ioEncStringWriter)
		e.wi.w = sww
		//ww = bufio.NewWriterSize(w, defEncByteBufSize)
	***REMOVED***
	e.w = &e.wi
	e.e.reset()
***REMOVED***

func (e *Encoder) ResetBytes(out *[]byte) ***REMOVED***
	in := *out
	if in == nil ***REMOVED***
		in = make([]byte, defEncByteBufSize)
	***REMOVED***
	e.wb.b, e.wb.out, e.wb.c = in, out, 0
	e.w = &e.wb
	e.e.reset()
***REMOVED***

// func (e *Encoder) sendContainerState(c containerState) ***REMOVED***
// 	if e.cr != nil ***REMOVED***
// 		e.cr.sendContainerState(c)
// 	***REMOVED***
// ***REMOVED***

// Encode writes an object into a stream.
//
// Encoding can be configured via the struct tag for the fields.
// The "codec" key in struct field's tag value is the key name,
// followed by an optional comma and options.
// Note that the "json" key is used in the absence of the "codec" key.
//
// To set an option on all fields (e.g. omitempty on all fields), you
// can create a field called _struct, and set flags on it.
//
// Struct values "usually" encode as maps. Each exported struct field is encoded unless:
//    - the field's tag is "-", OR
//    - the field is empty (empty or the zero value) and its tag specifies the "omitempty" option.
//
// When encoding as a map, the first string in the tag (before the comma)
// is the map key string to use when encoding.
//
// However, struct values may encode as arrays. This happens when:
//    - StructToArray Encode option is set, OR
//    - the tag on the _struct field sets the "toarray" option
//
// Values with types that implement MapBySlice are encoded as stream maps.
//
// The empty values (for omitempty option) are false, 0, any nil pointer
// or interface value, and any array, slice, map, or string of length zero.
//
// Anonymous fields are encoded inline except:
//    - the struct tag specifies a replacement name (first value)
//    - the field is of an interface type
//
// Examples:
//
//      // NOTE: 'json:' can be used as struct tag key, in place 'codec:' below.
//      type MyStruct struct ***REMOVED***
//          _struct bool    `codec:",omitempty"`   //set omitempty for every field
//          Field1 string   `codec:"-"`            //skip this field
//          Field2 int      `codec:"myName"`       //Use key "myName" in encode stream
//          Field3 int32    `codec:",omitempty"`   //use key "Field3". Omit if empty.
//          Field4 bool     `codec:"f4,omitempty"` //use key "f4". Omit if empty.
//          io.Reader                              //use key "Reader".
//          MyStruct        `codec:"my1"           //use key "my1".
//          MyStruct                               //inline it
//          ...
//      ***REMOVED***
//
//      type MyStruct struct ***REMOVED***
//          _struct bool    `codec:",omitempty,toarray"`   //set omitempty for every field
//                                                         //and encode struct as an array
//      ***REMOVED***
//
// The mode of encoding is based on the type of the value. When a value is seen:
//   - If a Selfer, call its CodecEncodeSelf method
//   - If an extension is registered for it, call that extension function
//   - If it implements encoding.(Binary|Text|JSON)Marshaler, call its Marshal(Binary|Text|JSON) method
//   - Else encode it based on its reflect.Kind
//
// Note that struct field names and keys in map[string]XXX will be treated as symbols.
// Some formats support symbols (e.g. binc) and will properly encode the string
// only once in the stream, and use a tag to refer to it thereafter.
func (e *Encoder) Encode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer panicToErr(&err)
	e.encode(v)
	e.w.atEndOfEncode()
	return
***REMOVED***

// MustEncode is like Encode, but panics if unable to Encode.
// This provides insight to the code location that triggered the error.
func (e *Encoder) MustEncode(v interface***REMOVED******REMOVED***) ***REMOVED***
	e.encode(v)
	e.w.atEndOfEncode()
***REMOVED***

func (e *Encoder) encode(iv interface***REMOVED******REMOVED***) ***REMOVED***
	// if ics, ok := iv.(Selfer); ok ***REMOVED***
	// 	ics.CodecEncodeSelf(e)
	// 	return
	// ***REMOVED***

	switch v := iv.(type) ***REMOVED***
	case nil:
		e.e.EncodeNil()
	case Selfer:
		v.CodecEncodeSelf(e)
	case Raw:
		e.raw(v)
	case reflect.Value:
		e.encodeValue(v, nil)

	case string:
		e.e.EncodeString(c_UTF8, v)
	case bool:
		e.e.EncodeBool(v)
	case int:
		e.e.EncodeInt(int64(v))
	case int8:
		e.e.EncodeInt(int64(v))
	case int16:
		e.e.EncodeInt(int64(v))
	case int32:
		e.e.EncodeInt(int64(v))
	case int64:
		e.e.EncodeInt(v)
	case uint:
		e.e.EncodeUint(uint64(v))
	case uint8:
		e.e.EncodeUint(uint64(v))
	case uint16:
		e.e.EncodeUint(uint64(v))
	case uint32:
		e.e.EncodeUint(uint64(v))
	case uint64:
		e.e.EncodeUint(v)
	case float32:
		e.e.EncodeFloat32(v)
	case float64:
		e.e.EncodeFloat64(v)

	case []uint8:
		e.e.EncodeStringBytes(c_RAW, v)

	case *string:
		e.e.EncodeString(c_UTF8, *v)
	case *bool:
		e.e.EncodeBool(*v)
	case *int:
		e.e.EncodeInt(int64(*v))
	case *int8:
		e.e.EncodeInt(int64(*v))
	case *int16:
		e.e.EncodeInt(int64(*v))
	case *int32:
		e.e.EncodeInt(int64(*v))
	case *int64:
		e.e.EncodeInt(*v)
	case *uint:
		e.e.EncodeUint(uint64(*v))
	case *uint8:
		e.e.EncodeUint(uint64(*v))
	case *uint16:
		e.e.EncodeUint(uint64(*v))
	case *uint32:
		e.e.EncodeUint(uint64(*v))
	case *uint64:
		e.e.EncodeUint(*v)
	case *float32:
		e.e.EncodeFloat32(*v)
	case *float64:
		e.e.EncodeFloat64(*v)

	case *[]uint8:
		e.e.EncodeStringBytes(c_RAW, *v)

	default:
		const checkCodecSelfer1 = true // in case T is passed, where *T is a Selfer, still checkCodecSelfer
		if !fastpathEncodeTypeSwitch(iv, e) ***REMOVED***
			e.encodeI(iv, false, checkCodecSelfer1)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Encoder) preEncodeValue(rv reflect.Value) (rv2 reflect.Value, sptr uintptr, proceed bool) ***REMOVED***
	// use a goto statement instead of a recursive function for ptr/interface.
TOP:
	switch rv.Kind() ***REMOVED***
	case reflect.Ptr:
		if rv.IsNil() ***REMOVED***
			e.e.EncodeNil()
			return
		***REMOVED***
		rv = rv.Elem()
		if e.h.CheckCircularRef && rv.Kind() == reflect.Struct ***REMOVED***
			// TODO: Movable pointers will be an issue here. Future problem.
			sptr = rv.UnsafeAddr()
			break TOP
		***REMOVED***
		goto TOP
	case reflect.Interface:
		if rv.IsNil() ***REMOVED***
			e.e.EncodeNil()
			return
		***REMOVED***
		rv = rv.Elem()
		goto TOP
	case reflect.Slice, reflect.Map:
		if rv.IsNil() ***REMOVED***
			e.e.EncodeNil()
			return
		***REMOVED***
	case reflect.Invalid, reflect.Func:
		e.e.EncodeNil()
		return
	***REMOVED***

	proceed = true
	rv2 = rv
	return
***REMOVED***

func (e *Encoder) doEncodeValue(rv reflect.Value, fn *encFn, sptr uintptr,
	checkFastpath, checkCodecSelfer bool) ***REMOVED***
	if sptr != 0 ***REMOVED***
		if (&e.ci).add(sptr) ***REMOVED***
			e.errorf("circular reference found: # %d", sptr)
		***REMOVED***
	***REMOVED***
	if fn == nil ***REMOVED***
		rt := rv.Type()
		rtid := reflect.ValueOf(rt).Pointer()
		// fn = e.getEncFn(rtid, rt, true, true)
		fn = e.getEncFn(rtid, rt, checkFastpath, checkCodecSelfer)
	***REMOVED***
	fn.f(&fn.i, rv)
	if sptr != 0 ***REMOVED***
		(&e.ci).remove(sptr)
	***REMOVED***
***REMOVED***

func (e *Encoder) encodeI(iv interface***REMOVED******REMOVED***, checkFastpath, checkCodecSelfer bool) ***REMOVED***
	if rv, sptr, proceed := e.preEncodeValue(reflect.ValueOf(iv)); proceed ***REMOVED***
		e.doEncodeValue(rv, nil, sptr, checkFastpath, checkCodecSelfer)
	***REMOVED***
***REMOVED***

func (e *Encoder) encodeValue(rv reflect.Value, fn *encFn) ***REMOVED***
	// if a valid fn is passed, it MUST BE for the dereferenced type of rv
	if rv, sptr, proceed := e.preEncodeValue(rv); proceed ***REMOVED***
		e.doEncodeValue(rv, fn, sptr, true, true)
	***REMOVED***
***REMOVED***

func (e *Encoder) getEncFn(rtid uintptr, rt reflect.Type, checkFastpath, checkCodecSelfer bool) (fn *encFn) ***REMOVED***
	// rtid := reflect.ValueOf(rt).Pointer()
	var ok bool
	if useMapForCodecCache ***REMOVED***
		fn, ok = e.f[rtid]
	***REMOVED*** else ***REMOVED***
		for i := range e.s ***REMOVED***
			v := &(e.s[i])
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
		if e.f == nil ***REMOVED***
			e.f = make(map[uintptr]*encFn, initCollectionCap)
		***REMOVED***
		fn = new(encFn)
		e.f[rtid] = fn
	***REMOVED*** else ***REMOVED***
		if e.s == nil ***REMOVED***
			e.s = make([]encRtidFn, 0, initCollectionCap)
		***REMOVED***
		e.s = append(e.s, encRtidFn***REMOVED***rtid: rtid***REMOVED***)
		fn = &(e.s[len(e.s)-1]).fn
	***REMOVED***

	ti := e.h.getTypeInfo(rtid, rt)
	fi := &(fn.i)
	fi.e = e
	fi.ti = ti

	if checkCodecSelfer && ti.cs ***REMOVED***
		fn.f = (*encFnInfo).selferMarshal
	***REMOVED*** else if rtid == rawTypId ***REMOVED***
		fn.f = (*encFnInfo).raw
	***REMOVED*** else if rtid == rawExtTypId ***REMOVED***
		fn.f = (*encFnInfo).rawExt
	***REMOVED*** else if e.e.IsBuiltinType(rtid) ***REMOVED***
		fn.f = (*encFnInfo).builtin
	***REMOVED*** else if xfFn := e.h.getExt(rtid); xfFn != nil ***REMOVED***
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.f = (*encFnInfo).ext
	***REMOVED*** else if supportMarshalInterfaces && e.be && ti.bm ***REMOVED***
		fn.f = (*encFnInfo).binaryMarshal
	***REMOVED*** else if supportMarshalInterfaces && !e.be && e.js && ti.jm ***REMOVED***
		//If JSON, we should check JSONMarshal before textMarshal
		fn.f = (*encFnInfo).jsonMarshal
	***REMOVED*** else if supportMarshalInterfaces && !e.be && ti.tm ***REMOVED***
		fn.f = (*encFnInfo).textMarshal
	***REMOVED*** else ***REMOVED***
		rk := rt.Kind()
		if fastpathEnabled && checkFastpath && (rk == reflect.Map || rk == reflect.Slice) ***REMOVED***
			if rt.PkgPath() == "" ***REMOVED*** // un-named slice or map
				if idx := fastpathAV.index(rtid); idx != -1 ***REMOVED***
					fn.f = fastpathAV[idx].encfn
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				ok = false
				// use mapping for underlying type if there
				var rtu reflect.Type
				if rk == reflect.Map ***REMOVED***
					rtu = reflect.MapOf(rt.Key(), rt.Elem())
				***REMOVED*** else ***REMOVED***
					rtu = reflect.SliceOf(rt.Elem())
				***REMOVED***
				rtuid := reflect.ValueOf(rtu).Pointer()
				if idx := fastpathAV.index(rtuid); idx != -1 ***REMOVED***
					xfnf := fastpathAV[idx].encfn
					xrt := fastpathAV[idx].rt
					fn.f = func(xf *encFnInfo, xrv reflect.Value) ***REMOVED***
						xfnf(xf, xrv.Convert(xrt))
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if fn.f == nil ***REMOVED***
			switch rk ***REMOVED***
			case reflect.Bool:
				fn.f = (*encFnInfo).kBool
			case reflect.String:
				fn.f = (*encFnInfo).kString
			case reflect.Float64:
				fn.f = (*encFnInfo).kFloat64
			case reflect.Float32:
				fn.f = (*encFnInfo).kFloat32
			case reflect.Int, reflect.Int8, reflect.Int64, reflect.Int32, reflect.Int16:
				fn.f = (*encFnInfo).kInt
			case reflect.Uint8, reflect.Uint64, reflect.Uint, reflect.Uint32, reflect.Uint16, reflect.Uintptr:
				fn.f = (*encFnInfo).kUint
			case reflect.Invalid:
				fn.f = (*encFnInfo).kInvalid
			case reflect.Chan:
				fi.seq = seqTypeChan
				fn.f = (*encFnInfo).kSlice
			case reflect.Slice:
				fi.seq = seqTypeSlice
				fn.f = (*encFnInfo).kSlice
			case reflect.Array:
				fi.seq = seqTypeArray
				fn.f = (*encFnInfo).kSlice
			case reflect.Struct:
				fn.f = (*encFnInfo).kStruct
				// reflect.Ptr and reflect.Interface are handled already by preEncodeValue
				// case reflect.Ptr:
				// 	fn.f = (*encFnInfo).kPtr
				// case reflect.Interface:
				// 	fn.f = (*encFnInfo).kInterface
			case reflect.Map:
				fn.f = (*encFnInfo).kMap
			default:
				fn.f = (*encFnInfo).kErr
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return
***REMOVED***

func (e *Encoder) marshal(bs []byte, fnerr error, asis bool, c charEncoding) ***REMOVED***
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
	if bs == nil ***REMOVED***
		e.e.EncodeNil()
	***REMOVED*** else if asis ***REMOVED***
		e.asis(bs)
	***REMOVED*** else ***REMOVED***
		e.e.EncodeStringBytes(c, bs)
	***REMOVED***
***REMOVED***

func (e *Encoder) asis(v []byte) ***REMOVED***
	if e.as == nil ***REMOVED***
		e.w.writeb(v)
	***REMOVED*** else ***REMOVED***
		e.as.EncodeAsis(v)
	***REMOVED***
***REMOVED***

func (e *Encoder) raw(vv Raw) ***REMOVED***
	v := []byte(vv)
	if !e.h.Raw ***REMOVED***
		e.errorf("Raw values cannot be encoded: %v", v)
	***REMOVED***
	if e.as == nil ***REMOVED***
		e.w.writeb(v)
	***REMOVED*** else ***REMOVED***
		e.as.EncodeAsis(v)
	***REMOVED***
***REMOVED***

func (e *Encoder) errorf(format string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	err := fmt.Errorf(format, params...)
	panic(err)
***REMOVED***

// ----------------------------------------

const encStructPoolLen = 5

// encStructPool is an array of sync.Pool.
// Each element of the array pools one of encStructPool(8|16|32|64).
// It allows the re-use of slices up to 64 in length.
// A performance cost of encoding structs was collecting
// which values were empty and should be omitted.
// We needed slices of reflect.Value and string to collect them.
// This shared pool reduces the amount of unnecessary creation we do.
// The cost is that of locking sometimes, but sync.Pool is efficient
// enough to reduce thread contention.
var encStructPool [encStructPoolLen]sync.Pool

func init() ***REMOVED***
	encStructPool[0].New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([8]stringRv) ***REMOVED***
	encStructPool[1].New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([16]stringRv) ***REMOVED***
	encStructPool[2].New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([32]stringRv) ***REMOVED***
	encStructPool[3].New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([64]stringRv) ***REMOVED***
	encStructPool[4].New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([128]stringRv) ***REMOVED***
***REMOVED***

func encStructPoolGet(newlen int) (p *sync.Pool, v interface***REMOVED******REMOVED***, s []stringRv) ***REMOVED***
	// if encStructPoolLen != 5 ***REMOVED*** // constant chec, so removed at build time.
	// 	panic(errors.New("encStructPoolLen must be equal to 4")) // defensive, in case it is changed
	// ***REMOVED***
	// idxpool := newlen / 8
	if newlen <= 8 ***REMOVED***
		p = &encStructPool[0]
		v = p.Get()
		s = v.(*[8]stringRv)[:newlen]
	***REMOVED*** else if newlen <= 16 ***REMOVED***
		p = &encStructPool[1]
		v = p.Get()
		s = v.(*[16]stringRv)[:newlen]
	***REMOVED*** else if newlen <= 32 ***REMOVED***
		p = &encStructPool[2]
		v = p.Get()
		s = v.(*[32]stringRv)[:newlen]
	***REMOVED*** else if newlen <= 64 ***REMOVED***
		p = &encStructPool[3]
		v = p.Get()
		s = v.(*[64]stringRv)[:newlen]
	***REMOVED*** else if newlen <= 128 ***REMOVED***
		p = &encStructPool[4]
		v = p.Get()
		s = v.(*[128]stringRv)[:newlen]
	***REMOVED*** else ***REMOVED***
		s = make([]stringRv, newlen)
	***REMOVED***
	return
***REMOVED***

// ----------------------------------------

// func encErr(format string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
// 	doPanic(msgTagEnc, format, params...)
// ***REMOVED***
