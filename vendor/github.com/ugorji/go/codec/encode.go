// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bufio"
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"
)

const defEncByteBufSize = 1 << 6 // 4:16, 6:64, 8:256, 10:1024

var errEncoderNotInitialized = errors.New("Encoder not initialized")

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
	EncodeNil()
	EncodeInt(i int64)
	EncodeUint(i uint64)
	EncodeBool(b bool)
	EncodeFloat32(f float32)
	EncodeFloat64(f float64)
	// encodeExtPreamble(xtag byte, length int)
	EncodeRawExt(re *RawExt, e *Encoder)
	EncodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext, e *Encoder)
	EncodeString(c charEncoding, v string)
	// EncodeSymbol(v string)
	EncodeStringBytes(c charEncoding, v []byte)
	EncodeTime(time.Time)
	//encBignum(f *big.Int)
	//encStringRunes(c charEncoding, v []rune)
	WriteArrayStart(length int)
	WriteArrayElem()
	WriteArrayEnd()
	WriteMapStart(length int)
	WriteMapElemKey()
	WriteMapElemValue()
	WriteMapEnd()

	reset()
	atEndOfEncode()
***REMOVED***

type ioEncStringWriter interface ***REMOVED***
	WriteString(s string) (n int, err error)
***REMOVED***

type encDriverAsis interface ***REMOVED***
	EncodeAsis(v []byte)
***REMOVED***

type encDriverNoopContainerWriter struct***REMOVED******REMOVED***

func (encDriverNoopContainerWriter) WriteArrayStart(length int) ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteArrayElem()            ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteArrayEnd()             ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteMapStart(length int)   ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteMapElemKey()           ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteMapElemValue()         ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteMapEnd()               ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) atEndOfEncode()             ***REMOVED******REMOVED***

type encDriverTrackContainerWriter struct ***REMOVED***
	c containerState
***REMOVED***

func (e *encDriverTrackContainerWriter) WriteArrayStart(length int) ***REMOVED*** e.c = containerArrayStart ***REMOVED***
func (e *encDriverTrackContainerWriter) WriteArrayElem()            ***REMOVED*** e.c = containerArrayElem ***REMOVED***
func (e *encDriverTrackContainerWriter) WriteArrayEnd()             ***REMOVED*** e.c = containerArrayEnd ***REMOVED***
func (e *encDriverTrackContainerWriter) WriteMapStart(length int)   ***REMOVED*** e.c = containerMapStart ***REMOVED***
func (e *encDriverTrackContainerWriter) WriteMapElemKey()           ***REMOVED*** e.c = containerMapKey ***REMOVED***
func (e *encDriverTrackContainerWriter) WriteMapElemValue()         ***REMOVED*** e.c = containerMapValue ***REMOVED***
func (e *encDriverTrackContainerWriter) WriteMapEnd()               ***REMOVED*** e.c = containerMapEnd ***REMOVED***
func (e *encDriverTrackContainerWriter) atEndOfEncode()             ***REMOVED******REMOVED***

// type ioEncWriterWriter interface ***REMOVED***
// 	WriteByte(c byte) error
// 	WriteString(s string) (n int, err error)
// 	Write(p []byte) (n int, err error)
// ***REMOVED***

// EncodeOptions captures configuration options during encode.
type EncodeOptions struct ***REMOVED***
	// WriterBufferSize is the size of the buffer used when writing.
	//
	// if > 0, we use a smart buffer internally for performance purposes.
	WriterBufferSize int

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

	// // AsSymbols defines what should be encoded as symbols.
	// //
	// // Encoding as symbols can reduce the encoded size significantly.
	// //
	// // However, during decoding, each string to be encoded as a symbol must
	// // be checked to see if it has been seen before. Consequently, encoding time
	// // will increase if using symbols, because string comparisons has a clear cost.
	// //
	// // Sample values:
	// //   AsSymbolNone
	// //   AsSymbolAll
	// //   AsSymbolMapStringKeys
	// //   AsSymbolMapStringKeysFlag | AsSymbolStructFieldNameFlag
	// AsSymbols AsSymbolFlag
***REMOVED***

// ---------------------------------------------

// ioEncWriter implements encWriter and can write to an io.Writer implementation
type ioEncWriter struct ***REMOVED***
	w  io.Writer
	ww io.Writer
	bw io.ByteWriter
	sw ioEncStringWriter
	fw ioFlusher
	b  [8]byte
***REMOVED***

func (z *ioEncWriter) WriteByte(b byte) (err error) ***REMOVED***
	z.b[0] = b
	_, err = z.w.Write(z.b[:1])
	return
***REMOVED***

func (z *ioEncWriter) WriteString(s string) (n int, err error) ***REMOVED***
	return z.w.Write(bytesView(s))
***REMOVED***

func (z *ioEncWriter) writeb(bs []byte) ***REMOVED***
	if _, err := z.ww.Write(bs); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writestr(s string) ***REMOVED***
	if _, err := z.sw.WriteString(s); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writen1(b byte) ***REMOVED***
	if err := z.bw.WriteByte(b); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (z *ioEncWriter) writen2(b1, b2 byte) ***REMOVED***
	var err error
	if err = z.bw.WriteByte(b1); err == nil ***REMOVED***
		if err = z.bw.WriteByte(b2); err == nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	panic(err)
***REMOVED***

// func (z *ioEncWriter) writen5(b1, b2, b3, b4, b5 byte) ***REMOVED***
// 	z.b[0], z.b[1], z.b[2], z.b[3], z.b[4] = b1, b2, b3, b4, b5
// 	if _, err := z.ww.Write(z.b[:5]); err != nil ***REMOVED***
// 		panic(err)
// 	***REMOVED***
// ***REMOVED***

func (z *ioEncWriter) atEndOfEncode() ***REMOVED***
	if z.fw != nil ***REMOVED***
		z.fw.Flush()
	***REMOVED***
***REMOVED***

// ---------------------------------------------

// bytesEncAppender implements encWriter and can write to an byte slice.
type bytesEncAppender struct ***REMOVED***
	b   []byte
	out *[]byte
***REMOVED***

func (z *bytesEncAppender) writeb(s []byte) ***REMOVED***
	z.b = append(z.b, s...)
***REMOVED***
func (z *bytesEncAppender) writestr(s string) ***REMOVED***
	z.b = append(z.b, s...)
***REMOVED***
func (z *bytesEncAppender) writen1(b1 byte) ***REMOVED***
	z.b = append(z.b, b1)
***REMOVED***
func (z *bytesEncAppender) writen2(b1, b2 byte) ***REMOVED***
	z.b = append(z.b, b1, b2)
***REMOVED***
func (z *bytesEncAppender) atEndOfEncode() ***REMOVED***
	*(z.out) = z.b
***REMOVED***
func (z *bytesEncAppender) reset(in []byte, out *[]byte) ***REMOVED***
	z.b = in[:0]
	z.out = out
***REMOVED***

// ---------------------------------------------

func (e *Encoder) rawExt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeRawExt(rv2i(rv).(*RawExt), e)
***REMOVED***

func (e *Encoder) ext(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeExt(rv2i(rv), f.xfTag, f.xfFn, e)
***REMOVED***

func (e *Encoder) selferMarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv2i(rv).(Selfer).CodecEncodeSelf(e)
***REMOVED***

func (e *Encoder) binaryMarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshal(bs, fnerr, false, cRAW)
***REMOVED***

func (e *Encoder) textMarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshal(bs, fnerr, false, cUTF8)
***REMOVED***

func (e *Encoder) jsonMarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshal(bs, fnerr, true, cUTF8)
***REMOVED***

func (e *Encoder) raw(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.rawBytes(rv2i(rv).(Raw))
***REMOVED***

func (e *Encoder) kInvalid(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeNil()
***REMOVED***

func (e *Encoder) kErr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.errorf("unsupported kind %s, for %#v", rv.Kind(), rv)
***REMOVED***

func (e *Encoder) kSlice(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	ti := f.ti
	ee := e.e
	// array may be non-addressable, so we have to manage with care
	//   (don't call rv.Bytes, rv.Slice, etc).
	// E.g. type struct S***REMOVED***B [2]byte***REMOVED***;
	//   Encode(S***REMOVED******REMOVED***) will bomb on "panic: slice of unaddressable array".
	if f.seq != seqTypeArray ***REMOVED***
		if rv.IsNil() ***REMOVED***
			ee.EncodeNil()
			return
		***REMOVED***
		// If in this method, then there was no extension function defined.
		// So it's okay to treat as []byte.
		if ti.rtid == uint8SliceTypId ***REMOVED***
			ee.EncodeStringBytes(cRAW, rv.Bytes())
			return
		***REMOVED***
	***REMOVED***
	if f.seq == seqTypeChan && ti.chandir&uint8(reflect.RecvDir) == 0 ***REMOVED***
		e.errorf("send-only channel cannot be used for receiving byte(s)")
	***REMOVED***
	elemsep := e.esep
	l := rv.Len()
	rtelem := ti.elem
	rtelemIsByte := uint8TypId == rt2id(rtelem) // NOT rtelem.Kind() == reflect.Uint8
	// if a slice, array or chan of bytes, treat specially
	if rtelemIsByte ***REMOVED***
		switch f.seq ***REMOVED***
		case seqTypeSlice:
			ee.EncodeStringBytes(cRAW, rv.Bytes())
		case seqTypeArray:
			if rv.CanAddr() ***REMOVED***
				ee.EncodeStringBytes(cRAW, rv.Slice(0, l).Bytes())
			***REMOVED*** else ***REMOVED***
				var bs []byte
				if l <= cap(e.b) ***REMOVED***
					bs = e.b[:l]
				***REMOVED*** else ***REMOVED***
					bs = make([]byte, l)
				***REMOVED***
				reflect.Copy(reflect.ValueOf(bs), rv)
				ee.EncodeStringBytes(cRAW, bs)
			***REMOVED***
		case seqTypeChan:
			bs := e.b[:0]
			// do not use range, so that the number of elements encoded
			// does not change, and encoding does not hang waiting on someone to close chan.
			// for b := range rv2i(rv).(<-chan byte) ***REMOVED*** bs = append(bs, b) ***REMOVED***
			// ch := rv2i(rv).(<-chan byte) // fix error - that this is a chan byte, not a <-chan byte.
			irv := rv2i(rv)
			ch, ok := irv.(<-chan byte)
			if !ok ***REMOVED***
				ch = irv.(chan byte)
			***REMOVED***
			for i := 0; i < l; i++ ***REMOVED***
				bs = append(bs, <-ch)
			***REMOVED***
			ee.EncodeStringBytes(cRAW, bs)
		***REMOVED***
		return
	***REMOVED***

	if ti.mbs ***REMOVED***
		if l%2 == 1 ***REMOVED***
			e.errorf("mapBySlice requires even slice length, but got %v", l)
			return
		***REMOVED***
		ee.WriteMapStart(l / 2)
	***REMOVED*** else ***REMOVED***
		ee.WriteArrayStart(l)
	***REMOVED***

	if l > 0 ***REMOVED***
		var fn *codecFn
		for rtelem.Kind() == reflect.Ptr ***REMOVED***
			rtelem = rtelem.Elem()
		***REMOVED***
		// if kind is reflect.Interface, do not pre-determine the
		// encoding type, because preEncodeValue may break it down to
		// a concrete type and kInterface will bomb.
		if rtelem.Kind() != reflect.Interface ***REMOVED***
			fn = e.cfer().get(rtelem, true, true)
		***REMOVED***
		for j := 0; j < l; j++ ***REMOVED***
			if elemsep ***REMOVED***
				if ti.mbs ***REMOVED***
					if j%2 == 0 ***REMOVED***
						ee.WriteMapElemKey()
					***REMOVED*** else ***REMOVED***
						ee.WriteMapElemValue()
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					ee.WriteArrayElem()
				***REMOVED***
			***REMOVED***
			if f.seq == seqTypeChan ***REMOVED***
				if rv2, ok2 := rv.Recv(); ok2 ***REMOVED***
					e.encodeValue(rv2, fn, true)
				***REMOVED*** else ***REMOVED***
					ee.EncodeNil() // WE HAVE TO DO SOMETHING, so nil if nothing received.
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				e.encodeValue(rv.Index(j), fn, true)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if ti.mbs ***REMOVED***
		ee.WriteMapEnd()
	***REMOVED*** else ***REMOVED***
		ee.WriteArrayEnd()
	***REMOVED***
***REMOVED***

func (e *Encoder) kStructNoOmitempty(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	fti := f.ti
	elemsep := e.esep
	tisfi := fti.sfiSrc
	toMap := !(fti.toArray || e.h.StructToArray)
	if toMap ***REMOVED***
		tisfi = fti.sfiSort
	***REMOVED***
	ee := e.e

	sfn := structFieldNode***REMOVED***v: rv, update: false***REMOVED***
	if toMap ***REMOVED***
		ee.WriteMapStart(len(tisfi))
		if elemsep ***REMOVED***
			for _, si := range tisfi ***REMOVED***
				ee.WriteMapElemKey()
				// ee.EncodeString(cUTF8, si.encName)
				encStructFieldKey(ee, fti.keyType, si.encName)
				ee.WriteMapElemValue()
				e.encodeValue(sfn.field(si), nil, true)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for _, si := range tisfi ***REMOVED***
				// ee.EncodeString(cUTF8, si.encName)
				encStructFieldKey(ee, fti.keyType, si.encName)
				e.encodeValue(sfn.field(si), nil, true)
			***REMOVED***
		***REMOVED***
		ee.WriteMapEnd()
	***REMOVED*** else ***REMOVED***
		ee.WriteArrayStart(len(tisfi))
		if elemsep ***REMOVED***
			for _, si := range tisfi ***REMOVED***
				ee.WriteArrayElem()
				e.encodeValue(sfn.field(si), nil, true)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for _, si := range tisfi ***REMOVED***
				e.encodeValue(sfn.field(si), nil, true)
			***REMOVED***
		***REMOVED***
		ee.WriteArrayEnd()
	***REMOVED***
***REMOVED***

func encStructFieldKey(ee encDriver, keyType valueType, s string) ***REMOVED***
	var m must

	// use if-else-if, not switch (which compiles to binary-search)
	// since keyType is typically valueTypeString, branch prediction is pretty good.

	if keyType == valueTypeString ***REMOVED***
		ee.EncodeString(cUTF8, s)
	***REMOVED*** else if keyType == valueTypeInt ***REMOVED***
		ee.EncodeInt(m.Int(strconv.ParseInt(s, 10, 64)))
	***REMOVED*** else if keyType == valueTypeUint ***REMOVED***
		ee.EncodeUint(m.Uint(strconv.ParseUint(s, 10, 64)))
	***REMOVED*** else if keyType == valueTypeFloat ***REMOVED***
		ee.EncodeFloat64(m.Float(strconv.ParseFloat(s, 64)))
	***REMOVED*** else ***REMOVED***
		ee.EncodeString(cUTF8, s)
	***REMOVED***
***REMOVED***

func (e *Encoder) kStruct(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	fti := f.ti
	elemsep := e.esep
	tisfi := fti.sfiSrc
	toMap := !(fti.toArray || e.h.StructToArray)
	// if toMap, use the sorted array. If toArray, use unsorted array (to match sequence in struct)
	if toMap ***REMOVED***
		tisfi = fti.sfiSort
	***REMOVED***
	newlen := len(fti.sfiSort)
	ee := e.e

	// Use sync.Pool to reduce allocating slices unnecessarily.
	// The cost of sync.Pool is less than the cost of new allocation.
	//
	// Each element of the array pools one of encStructPool(8|16|32|64).
	// It allows the re-use of slices up to 64 in length.
	// A performance cost of encoding structs was collecting
	// which values were empty and should be omitted.
	// We needed slices of reflect.Value and string to collect them.
	// This shared pool reduces the amount of unnecessary creation we do.
	// The cost is that of locking sometimes, but sync.Pool is efficient
	// enough to reduce thread contention.

	var spool *sync.Pool
	var poolv interface***REMOVED******REMOVED***
	var fkvs []stringRv
	// fmt.Printf(">>>>>>>>>>>>>> encode.kStruct: newlen: %d\n", newlen)
	if newlen <= 8 ***REMOVED***
		spool, poolv = pool.stringRv8()
		fkvs = poolv.(*[8]stringRv)[:newlen]
	***REMOVED*** else if newlen <= 16 ***REMOVED***
		spool, poolv = pool.stringRv16()
		fkvs = poolv.(*[16]stringRv)[:newlen]
	***REMOVED*** else if newlen <= 32 ***REMOVED***
		spool, poolv = pool.stringRv32()
		fkvs = poolv.(*[32]stringRv)[:newlen]
	***REMOVED*** else if newlen <= 64 ***REMOVED***
		spool, poolv = pool.stringRv64()
		fkvs = poolv.(*[64]stringRv)[:newlen]
	***REMOVED*** else if newlen <= 128 ***REMOVED***
		spool, poolv = pool.stringRv128()
		fkvs = poolv.(*[128]stringRv)[:newlen]
	***REMOVED*** else ***REMOVED***
		fkvs = make([]stringRv, newlen)
	***REMOVED***

	newlen = 0
	var kv stringRv
	recur := e.h.RecursiveEmptyCheck
	sfn := structFieldNode***REMOVED***v: rv, update: false***REMOVED***
	for _, si := range tisfi ***REMOVED***
		// kv.r = si.field(rv, false)
		kv.r = sfn.field(si)
		if toMap ***REMOVED***
			if si.omitEmpty() && isEmptyValue(kv.r, e.h.TypeInfos, recur, recur) ***REMOVED***
				continue
			***REMOVED***
			kv.v = si.encName
		***REMOVED*** else ***REMOVED***
			// use the zero value.
			// if a reference or struct, set to nil (so you do not output too much)
			if si.omitEmpty() && isEmptyValue(kv.r, e.h.TypeInfos, recur, recur) ***REMOVED***
				switch kv.r.Kind() ***REMOVED***
				case reflect.Struct, reflect.Interface, reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice:
					kv.r = reflect.Value***REMOVED******REMOVED*** //encode as nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
		fkvs[newlen] = kv
		newlen++
	***REMOVED***

	if toMap ***REMOVED***
		ee.WriteMapStart(newlen)
		if elemsep ***REMOVED***
			for j := 0; j < newlen; j++ ***REMOVED***
				kv = fkvs[j]
				ee.WriteMapElemKey()
				// ee.EncodeString(cUTF8, kv.v)
				encStructFieldKey(ee, fti.keyType, kv.v)
				ee.WriteMapElemValue()
				e.encodeValue(kv.r, nil, true)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for j := 0; j < newlen; j++ ***REMOVED***
				kv = fkvs[j]
				// ee.EncodeString(cUTF8, kv.v)
				encStructFieldKey(ee, fti.keyType, kv.v)
				e.encodeValue(kv.r, nil, true)
			***REMOVED***
		***REMOVED***
		ee.WriteMapEnd()
	***REMOVED*** else ***REMOVED***
		ee.WriteArrayStart(newlen)
		if elemsep ***REMOVED***
			for j := 0; j < newlen; j++ ***REMOVED***
				ee.WriteArrayElem()
				e.encodeValue(fkvs[j].r, nil, true)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for j := 0; j < newlen; j++ ***REMOVED***
				e.encodeValue(fkvs[j].r, nil, true)
			***REMOVED***
		***REMOVED***
		ee.WriteArrayEnd()
	***REMOVED***

	// do not use defer. Instead, use explicit pool return at end of function.
	// defer has a cost we are trying to avoid.
	// If there is a panic and these slices are not returned, it is ok.
	if spool != nil ***REMOVED***
		spool.Put(poolv)
	***REMOVED***
***REMOVED***

func (e *Encoder) kMap(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	ee := e.e
	if rv.IsNil() ***REMOVED***
		ee.EncodeNil()
		return
	***REMOVED***

	l := rv.Len()
	ee.WriteMapStart(l)
	elemsep := e.esep
	if l == 0 ***REMOVED***
		ee.WriteMapEnd()
		return
	***REMOVED***
	// var asSymbols bool
	// determine the underlying key and val encFn's for the map.
	// This eliminates some work which is done for each loop iteration i.e.
	// rv.Type(), ref.ValueOf(rt).Pointer(), then check map/list for fn.
	//
	// However, if kind is reflect.Interface, do not pre-determine the
	// encoding type, because preEncodeValue may break it down to
	// a concrete type and kInterface will bomb.
	var keyFn, valFn *codecFn
	ti := f.ti
	rtkey0 := ti.key
	rtkey := rtkey0
	rtval0 := ti.elem
	rtval := rtval0
	// rtkeyid := rt2id(rtkey0)
	for rtval.Kind() == reflect.Ptr ***REMOVED***
		rtval = rtval.Elem()
	***REMOVED***
	if rtval.Kind() != reflect.Interface ***REMOVED***
		valFn = e.cfer().get(rtval, true, true)
	***REMOVED***
	mks := rv.MapKeys()

	if e.h.Canonical ***REMOVED***
		e.kMapCanonical(rtkey, rv, mks, valFn)
		ee.WriteMapEnd()
		return
	***REMOVED***

	var keyTypeIsString = stringTypId == rt2id(rtkey0) // rtkeyid
	if !keyTypeIsString ***REMOVED***
		for rtkey.Kind() == reflect.Ptr ***REMOVED***
			rtkey = rtkey.Elem()
		***REMOVED***
		if rtkey.Kind() != reflect.Interface ***REMOVED***
			// rtkeyid = rt2id(rtkey)
			keyFn = e.cfer().get(rtkey, true, true)
		***REMOVED***
	***REMOVED***

	// for j, lmks := 0, len(mks); j < lmks; j++ ***REMOVED***
	for j := range mks ***REMOVED***
		if elemsep ***REMOVED***
			ee.WriteMapElemKey()
		***REMOVED***
		if keyTypeIsString ***REMOVED***
			ee.EncodeString(cUTF8, mks[j].String())
		***REMOVED*** else ***REMOVED***
			e.encodeValue(mks[j], keyFn, true)
		***REMOVED***
		if elemsep ***REMOVED***
			ee.WriteMapElemValue()
		***REMOVED***
		e.encodeValue(rv.MapIndex(mks[j]), valFn, true)

	***REMOVED***
	ee.WriteMapEnd()
***REMOVED***

func (e *Encoder) kMapCanonical(rtkey reflect.Type, rv reflect.Value, mks []reflect.Value, valFn *codecFn) ***REMOVED***
	ee := e.e
	elemsep := e.esep
	// we previously did out-of-band if an extension was registered.
	// This is not necessary, as the natural kind is sufficient for ordering.

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
			if elemsep ***REMOVED***
				ee.WriteMapElemKey()
			***REMOVED***
			ee.EncodeBool(mksv[i].v)
			if elemsep ***REMOVED***
				ee.WriteMapElemValue()
			***REMOVED***
			e.encodeValue(rv.MapIndex(mksv[i].r), valFn, true)
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
			if elemsep ***REMOVED***
				ee.WriteMapElemKey()
			***REMOVED***
			ee.EncodeString(cUTF8, mksv[i].v)
			if elemsep ***REMOVED***
				ee.WriteMapElemValue()
			***REMOVED***
			e.encodeValue(rv.MapIndex(mksv[i].r), valFn, true)
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
			if elemsep ***REMOVED***
				ee.WriteMapElemKey()
			***REMOVED***
			ee.EncodeUint(mksv[i].v)
			if elemsep ***REMOVED***
				ee.WriteMapElemValue()
			***REMOVED***
			e.encodeValue(rv.MapIndex(mksv[i].r), valFn, true)
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
			if elemsep ***REMOVED***
				ee.WriteMapElemKey()
			***REMOVED***
			ee.EncodeInt(mksv[i].v)
			if elemsep ***REMOVED***
				ee.WriteMapElemValue()
			***REMOVED***
			e.encodeValue(rv.MapIndex(mksv[i].r), valFn, true)
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
			if elemsep ***REMOVED***
				ee.WriteMapElemKey()
			***REMOVED***
			ee.EncodeFloat32(float32(mksv[i].v))
			if elemsep ***REMOVED***
				ee.WriteMapElemValue()
			***REMOVED***
			e.encodeValue(rv.MapIndex(mksv[i].r), valFn, true)
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
			if elemsep ***REMOVED***
				ee.WriteMapElemKey()
			***REMOVED***
			ee.EncodeFloat64(mksv[i].v)
			if elemsep ***REMOVED***
				ee.WriteMapElemValue()
			***REMOVED***
			e.encodeValue(rv.MapIndex(mksv[i].r), valFn, true)
		***REMOVED***
	case reflect.Struct:
		if rv.Type() == timeTyp ***REMOVED***
			mksv := make([]timeRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksv[i]
				v.r = k
				v.v = rv2i(k).(time.Time)
			***REMOVED***
			sort.Sort(timeRvSlice(mksv))
			for i := range mksv ***REMOVED***
				if elemsep ***REMOVED***
					ee.WriteMapElemKey()
				***REMOVED***
				ee.EncodeTime(mksv[i].v)
				if elemsep ***REMOVED***
					ee.WriteMapElemValue()
				***REMOVED***
				e.encodeValue(rv.MapIndex(mksv[i].r), valFn, true)
			***REMOVED***
			break
		***REMOVED***
		fallthrough
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
		***REMOVED***
		sort.Sort(bytesRvSlice(mksbv))
		for j := range mksbv ***REMOVED***
			if elemsep ***REMOVED***
				ee.WriteMapElemKey()
			***REMOVED***
			e.asis(mksbv[j].v)
			if elemsep ***REMOVED***
				ee.WriteMapElemValue()
			***REMOVED***
			e.encodeValue(rv.MapIndex(mksbv[j].r), valFn, true)
		***REMOVED***
	***REMOVED***
***REMOVED***

// // --------------------------------------------------

type encWriterSwitch struct ***REMOVED***
	wi *ioEncWriter
	// wb bytesEncWriter
	wb   bytesEncAppender
	wx   bool // if bytes, wx=true
	esep bool // whether it has elem separators
	isas bool // whether e.as != nil
***REMOVED***

// // TODO: Uncomment after mid-stack inlining enabled in go 1.10

// func (z *encWriterSwitch) writeb(s []byte) ***REMOVED***
// 	if z.wx ***REMOVED***
// 		z.wb.writeb(s)
// 	***REMOVED*** else ***REMOVED***
// 		z.wi.writeb(s)
// 	***REMOVED***
// ***REMOVED***
// func (z *encWriterSwitch) writestr(s string) ***REMOVED***
// 	if z.wx ***REMOVED***
// 		z.wb.writestr(s)
// 	***REMOVED*** else ***REMOVED***
// 		z.wi.writestr(s)
// 	***REMOVED***
// ***REMOVED***
// func (z *encWriterSwitch) writen1(b1 byte) ***REMOVED***
// 	if z.wx ***REMOVED***
// 		z.wb.writen1(b1)
// 	***REMOVED*** else ***REMOVED***
// 		z.wi.writen1(b1)
// 	***REMOVED***
// ***REMOVED***
// func (z *encWriterSwitch) writen2(b1, b2 byte) ***REMOVED***
// 	if z.wx ***REMOVED***
// 		z.wb.writen2(b1, b2)
// 	***REMOVED*** else ***REMOVED***
// 		z.wi.writen2(b1, b2)
// 	***REMOVED***
// ***REMOVED***

// An Encoder writes an object to an output stream in the codec format.
type Encoder struct ***REMOVED***
	panicHdl
	// hopefully, reduce derefencing cost by laying the encWriter inside the Encoder
	e encDriver
	// NOTE: Encoder shouldn't call it's write methods,
	// as the handler MAY need to do some coordination.
	w encWriter

	h  *BasicHandle
	bw *bufio.Writer
	as encDriverAsis

	// ---- cpu cache line boundary?

	// ---- cpu cache line boundary?
	encWriterSwitch
	err error

	// ---- cpu cache line boundary?
	codecFnPooler
	ci set
	js bool    // here, so that no need to piggy back on *codecFner for this
	be bool    // here, so that no need to piggy back on *codecFner for this
	_  [6]byte // padding

	// ---- writable fields during execution --- *try* to keep in sep cache line

	// ---- cpu cache line boundary?
	// b [scratchByteArrayLen]byte
	// _ [cacheLineSize - scratchByteArrayLen]byte // padding
	b [cacheLineSize - 0]byte // used for encoding a chan or (non-addressable) array of bytes
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
	e := &Encoder***REMOVED***h: h.getBasicHandle(), err: errEncoderNotInitialized***REMOVED***
	e.hh = h
	e.esep = h.hasElemSeparators()
	return e
***REMOVED***

func (e *Encoder) resetCommon() ***REMOVED***
	if e.e == nil || e.hh.recreateEncDriver(e.e) ***REMOVED***
		e.e = e.hh.newEncDriver(e)
		e.as, e.isas = e.e.(encDriverAsis)
		// e.cr, _ = e.e.(containerStateRecv)
	***REMOVED***
	e.be = e.hh.isBinary()
	_, e.js = e.hh.(*JsonHandle)
	e.e.reset()
	e.err = nil
***REMOVED***

// Reset resets the Encoder with a new output stream.
//
// This accommodates using the state of the Encoder,
// where it has "cached" information about sub-engines.
func (e *Encoder) Reset(w io.Writer) ***REMOVED***
	if w == nil ***REMOVED***
		return
	***REMOVED***
	if e.wi == nil ***REMOVED***
		e.wi = new(ioEncWriter)
	***REMOVED***
	var ok bool
	e.wx = false
	e.wi.w = w
	if e.h.WriterBufferSize > 0 ***REMOVED***
		e.bw = bufio.NewWriterSize(w, e.h.WriterBufferSize)
		e.wi.bw = e.bw
		e.wi.sw = e.bw
		e.wi.fw = e.bw
		e.wi.ww = e.bw
	***REMOVED*** else ***REMOVED***
		if e.wi.bw, ok = w.(io.ByteWriter); !ok ***REMOVED***
			e.wi.bw = e.wi
		***REMOVED***
		if e.wi.sw, ok = w.(ioEncStringWriter); !ok ***REMOVED***
			e.wi.sw = e.wi
		***REMOVED***
		e.wi.fw, _ = w.(ioFlusher)
		e.wi.ww = w
	***REMOVED***
	e.w = e.wi
	e.resetCommon()
***REMOVED***

// ResetBytes resets the Encoder with a new destination output []byte.
func (e *Encoder) ResetBytes(out *[]byte) ***REMOVED***
	if out == nil ***REMOVED***
		return
	***REMOVED***
	var in []byte
	if out != nil ***REMOVED***
		in = *out
	***REMOVED***
	if in == nil ***REMOVED***
		in = make([]byte, defEncByteBufSize)
	***REMOVED***
	e.wx = true
	e.wb.reset(in, out)
	e.w = &e.wb
	e.resetCommon()
***REMOVED***

// Encode writes an object into a stream.
//
// Encoding can be configured via the struct tag for the fields.
// The "codec" key in struct field's tag value is the key name,
// followed by an optional comma and options.
// Note that the "json" key is used in the absence of the "codec" key.
//
// To set an option on all fields (e.g. omitempty on all fields), you
// can create a field called _struct, and set flags on it. The options
// which can be set on _struct are:
//    - omitempty: so all fields are omitted if empty
//    - toarray: so struct is encoded as an array
//    - int: so struct key names are encoded as signed integers (instead of strings)
//    - uint: so struct key names are encoded as unsigned integers (instead of strings)
//    - float: so struct key names are encoded as floats (instead of strings)
// More details on these below.
//
// Struct values "usually" encode as maps. Each exported struct field is encoded unless:
//    - the field's tag is "-", OR
//    - the field is empty (empty or the zero value) and its tag specifies the "omitempty" option.
//
// When encoding as a map, the first string in the tag (before the comma)
// is the map key string to use when encoding.
// ...
// This key is typically encoded as a string.
// However, there are instances where the encoded stream has mapping keys encoded as numbers.
// For example, some cbor streams have keys as integer codes in the stream, but they should map
// to fields in a structured object. Consequently, a struct is the natural representation in code.
// For these, configure the struct to encode/decode the keys as numbers (instead of string).
// This is done with the int,uint or float option on the _struct field (see above).
//
// However, struct values may encode as arrays. This happens when:
//    - StructToArray Encode option is set, OR
//    - the tag on the _struct field sets the "toarray" option
// Note that omitempty is ignored when encoding struct values as arrays,
// as an entry must be encoded for each field, to maintain its position.
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
//          _struct bool    `codec:",toarray"`     //encode struct as an array
//      ***REMOVED***
//
//      type MyStruct struct ***REMOVED***
//          _struct bool    `codec:",uint"`        //encode struct with "unsigned integer" keys
//          Field1 string   `codec:"1"`            //encode Field1 key using: EncodeInt(1)
//          Field2 string   `codec:"2"`            //encode Field2 key using: EncodeInt(2)
//      ***REMOVED***
//
// The mode of encoding is based on the type of the value. When a value is seen:
//   - If a Selfer, call its CodecEncodeSelf method
//   - If an extension is registered for it, call that extension function
//   - If implements encoding.(Binary|Text|JSON)Marshaler, call Marshal(Binary|Text|JSON) method
//   - Else encode it based on its reflect.Kind
//
// Note that struct field names and keys in map[string]XXX will be treated as symbols.
// Some formats support symbols (e.g. binc) and will properly encode the string
// only once in the stream, and use a tag to refer to it thereafter.
func (e *Encoder) Encode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	defer panicToErrs2(e, &e.err, &err)
	defer e.alwaysAtEnd()
	e.MustEncode(v)
	return
***REMOVED***

// MustEncode is like Encode, but panics if unable to Encode.
// This provides insight to the code location that triggered the error.
func (e *Encoder) MustEncode(v interface***REMOVED******REMOVED***) ***REMOVED***
	if e.err != nil ***REMOVED***
		panic(e.err)
	***REMOVED***
	e.encode(v)
	e.e.atEndOfEncode()
	e.w.atEndOfEncode()
	e.alwaysAtEnd()
***REMOVED***

// func (e *Encoder) alwaysAtEnd() ***REMOVED***
// 	e.codecFnPooler.alwaysAtEnd()
// ***REMOVED***

func (e *Encoder) encode(iv interface***REMOVED******REMOVED***) ***REMOVED***
	if iv == nil || definitelyNil(iv) ***REMOVED***
		e.e.EncodeNil()
		return
	***REMOVED***
	if v, ok := iv.(Selfer); ok ***REMOVED***
		v.CodecEncodeSelf(e)
		return
	***REMOVED***

	// a switch with only concrete types can be optimized.
	// consequently, we deal with nil and interfaces outside.

	switch v := iv.(type) ***REMOVED***
	case Raw:
		e.rawBytes(v)
	case reflect.Value:
		e.encodeValue(v, nil, true)

	case string:
		e.e.EncodeString(cUTF8, v)
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
	case uintptr:
		e.e.EncodeUint(uint64(v))
	case float32:
		e.e.EncodeFloat32(v)
	case float64:
		e.e.EncodeFloat64(v)
	case time.Time:
		e.e.EncodeTime(v)
	case []uint8:
		e.e.EncodeStringBytes(cRAW, v)

	case *Raw:
		e.rawBytes(*v)

	case *string:
		e.e.EncodeString(cUTF8, *v)
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
	case *uintptr:
		e.e.EncodeUint(uint64(*v))
	case *float32:
		e.e.EncodeFloat32(*v)
	case *float64:
		e.e.EncodeFloat64(*v)
	case *time.Time:
		e.e.EncodeTime(*v)

	case *[]uint8:
		e.e.EncodeStringBytes(cRAW, *v)

	default:
		if !fastpathEncodeTypeSwitch(iv, e) ***REMOVED***
			// checkfastpath=true (not false), as underlying slice/map type may be fast-path
			e.encodeValue(reflect.ValueOf(iv), nil, true)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Encoder) encodeValue(rv reflect.Value, fn *codecFn, checkFastpath bool) ***REMOVED***
	// if a valid fn is passed, it MUST BE for the dereferenced type of rv
	var sptr uintptr
	var rvp reflect.Value
	var rvpValid bool
TOP:
	switch rv.Kind() ***REMOVED***
	case reflect.Ptr:
		if rv.IsNil() ***REMOVED***
			e.e.EncodeNil()
			return
		***REMOVED***
		rvpValid = true
		rvp = rv
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

	if sptr != 0 && (&e.ci).add(sptr) ***REMOVED***
		e.errorf("circular reference found: # %d", sptr)
	***REMOVED***

	if fn == nil ***REMOVED***
		rt := rv.Type()
		// always pass checkCodecSelfer=true, in case T or ****T is passed, where *T is a Selfer
		fn = e.cfer().get(rt, checkFastpath, true)
	***REMOVED***
	if fn.i.addrE ***REMOVED***
		if rvpValid ***REMOVED***
			fn.fe(e, &fn.i, rvp)
		***REMOVED*** else if rv.CanAddr() ***REMOVED***
			fn.fe(e, &fn.i, rv.Addr())
		***REMOVED*** else ***REMOVED***
			rv2 := reflect.New(rv.Type())
			rv2.Elem().Set(rv)
			fn.fe(e, &fn.i, rv2)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fn.fe(e, &fn.i, rv)
	***REMOVED***
	if sptr != 0 ***REMOVED***
		(&e.ci).remove(sptr)
	***REMOVED***
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
	if e.isas ***REMOVED***
		e.as.EncodeAsis(v)
	***REMOVED*** else ***REMOVED***
		e.w.writeb(v)
	***REMOVED***
***REMOVED***

func (e *Encoder) rawBytes(vv Raw) ***REMOVED***
	v := []byte(vv)
	if !e.h.Raw ***REMOVED***
		e.errorf("Raw values cannot be encoded: %v", v)
	***REMOVED***
	e.asis(v)
***REMOVED***

func (e *Encoder) wrapErrstr(v interface***REMOVED******REMOVED***, err *error) ***REMOVED***
	*err = fmt.Errorf("%s encode error: %v", e.hh.Name(), v)
***REMOVED***
