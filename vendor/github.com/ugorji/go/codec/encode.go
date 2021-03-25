// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"time"
)

// defEncByteBufSize is the default size of []byte used
// for bufio buffer or []byte (when nil passed)
const defEncByteBufSize = 1 << 10 // 4:16, 6:64, 8:256, 10:1024

var errEncoderNotInitialized = errors.New("Encoder not initialized")

// encDriver abstracts the actual codec (binc vs msgpack, etc)
type encDriver interface ***REMOVED***
	EncodeNil()
	EncodeInt(i int64)
	EncodeUint(i uint64)
	EncodeBool(b bool)
	EncodeFloat32(f float32)
	EncodeFloat64(f float64)
	EncodeRawExt(re *RawExt)
	EncodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext)
	// EncodeString using cUTF8, honor'ing StringToRaw flag
	EncodeString(v string)
	EncodeStringBytesRaw(v []byte)
	EncodeTime(time.Time)
	WriteArrayStart(length int)
	WriteArrayEnd()
	WriteMapStart(length int)
	WriteMapEnd()

	reset()
	atEndOfEncode()
	encoder() *Encoder
***REMOVED***

type encDriverContainerTracker interface ***REMOVED***
	WriteArrayElem()
	WriteMapElemKey()
	WriteMapElemValue()
***REMOVED***

type encodeError struct ***REMOVED***
	codecError
***REMOVED***

func (e encodeError) Error() string ***REMOVED***
	return fmt.Sprintf("%s encode error: %v", e.name, e.err)
***REMOVED***

type encDriverNoopContainerWriter struct***REMOVED******REMOVED***

func (encDriverNoopContainerWriter) WriteArrayStart(length int) ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteArrayEnd()             ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteMapStart(length int)   ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) WriteMapEnd()               ***REMOVED******REMOVED***
func (encDriverNoopContainerWriter) atEndOfEncode()             ***REMOVED******REMOVED***

// EncodeOptions captures configuration options during encode.
type EncodeOptions struct ***REMOVED***
	// WriterBufferSize is the size of the buffer used when writing.
	//
	// if > 0, we use a smart buffer internally for performance purposes.
	WriterBufferSize int

	// ChanRecvTimeout is the timeout used when selecting from a chan.
	//
	// Configuring this controls how we receive from a chan during the encoding process.
	//   - If ==0, we only consume the elements currently available in the chan.
	//   - if  <0, we consume until the chan is closed.
	//   - If  >0, we consume until this timeout.
	ChanRecvTimeout time.Duration

	// StructToArray specifies to encode a struct as an array, and not as a map
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

	// StringToRaw controls how strings are encoded.
	//
	// As a go string is just an (immutable) sequence of bytes,
	// it can be encoded either as raw bytes or as a UTF string.
	//
	// By default, strings are encoded as UTF-8.
	// but can be treated as []byte during an encode.
	//
	// Note that things which we know (by definition) to be UTF-8
	// are ALWAYS encoded as UTF-8 strings.
	// These include encoding.TextMarshaler, time.Format calls, struct field names, etc.
	StringToRaw bool

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

func (e *Encoder) rawExt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeRawExt(rv2i(rv).(*RawExt))
***REMOVED***

func (e *Encoder) ext(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeExt(rv2i(rv), f.xfTag, f.xfFn)
***REMOVED***

func (e *Encoder) selferMarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv2i(rv).(Selfer).CodecEncodeSelf(e)
***REMOVED***

func (e *Encoder) binaryMarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
***REMOVED***

func (e *Encoder) textMarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
***REMOVED***

func (e *Encoder) jsonMarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
***REMOVED***

func (e *Encoder) raw(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.rawBytes(rv2i(rv).(Raw))
***REMOVED***

func (e *Encoder) kBool(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeBool(rvGetBool(rv))
***REMOVED***

func (e *Encoder) kTime(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeTime(rvGetTime(rv))
***REMOVED***

func (e *Encoder) kString(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeString(rvGetString(rv))
***REMOVED***

func (e *Encoder) kFloat64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeFloat64(rvGetFloat64(rv))
***REMOVED***

func (e *Encoder) kFloat32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeFloat32(rvGetFloat32(rv))
***REMOVED***

func (e *Encoder) kInt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(int64(rvGetInt(rv)))
***REMOVED***

func (e *Encoder) kInt8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(int64(rvGetInt8(rv)))
***REMOVED***

func (e *Encoder) kInt16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(int64(rvGetInt16(rv)))
***REMOVED***

func (e *Encoder) kInt32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(int64(rvGetInt32(rv)))
***REMOVED***

func (e *Encoder) kInt64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(int64(rvGetInt64(rv)))
***REMOVED***

func (e *Encoder) kUint(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(uint64(rvGetUint(rv)))
***REMOVED***

func (e *Encoder) kUint8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
***REMOVED***

func (e *Encoder) kUint16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
***REMOVED***

func (e *Encoder) kUint32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
***REMOVED***

func (e *Encoder) kUint64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
***REMOVED***

func (e *Encoder) kUintptr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
***REMOVED***

func (e *Encoder) kInvalid(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeNil()
***REMOVED***

func (e *Encoder) kErr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.errorf("unsupported kind %s, for %#v", rv.Kind(), rv)
***REMOVED***

func chanToSlice(rv reflect.Value, rtslice reflect.Type, timeout time.Duration) (rvcs reflect.Value) ***REMOVED***
	rvcs = reflect.Zero(rtslice)
	if timeout < 0 ***REMOVED*** // consume until close
		for ***REMOVED***
			recv, recvOk := rv.Recv()
			if !recvOk ***REMOVED***
				break
			***REMOVED***
			rvcs = reflect.Append(rvcs, recv)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		cases := make([]reflect.SelectCase, 2)
		cases[0] = reflect.SelectCase***REMOVED***Dir: reflect.SelectRecv, Chan: rv***REMOVED***
		if timeout == 0 ***REMOVED***
			cases[1] = reflect.SelectCase***REMOVED***Dir: reflect.SelectDefault***REMOVED***
		***REMOVED*** else ***REMOVED***
			tt := time.NewTimer(timeout)
			cases[1] = reflect.SelectCase***REMOVED***Dir: reflect.SelectRecv, Chan: rv4i(tt.C)***REMOVED***
		***REMOVED***
		for ***REMOVED***
			chosen, recv, recvOk := reflect.Select(cases)
			if chosen == 1 || !recvOk ***REMOVED***
				break
			***REMOVED***
			rvcs = reflect.Append(rvcs, recv)
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (e *Encoder) kSeqFn(rtelem reflect.Type) (fn *codecFn) ***REMOVED***
	for rtelem.Kind() == reflect.Ptr ***REMOVED***
		rtelem = rtelem.Elem()
	***REMOVED***
	// if kind is reflect.Interface, do not pre-determine the
	// encoding type, because preEncodeValue may break it down to
	// a concrete type and kInterface will bomb.
	if rtelem.Kind() != reflect.Interface ***REMOVED***
		fn = e.h.fn(rtelem)
	***REMOVED***
	return
***REMOVED***

func (e *Encoder) kSliceWMbs(rv reflect.Value, ti *typeInfo) ***REMOVED***
	var l = rvGetSliceLen(rv)
	if l == 0 ***REMOVED***
		e.mapStart(0)
	***REMOVED*** else ***REMOVED***
		if l%2 == 1 ***REMOVED***
			e.errorf("mapBySlice requires even slice length, but got %v", l)
			return
		***REMOVED***
		e.mapStart(l / 2)
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ ***REMOVED***
			if j%2 == 0 ***REMOVED***
				e.mapElemKey()
			***REMOVED*** else ***REMOVED***
				e.mapElemValue()
			***REMOVED***
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		***REMOVED***
	***REMOVED***
	e.mapEnd()
***REMOVED***

func (e *Encoder) kSliceW(rv reflect.Value, ti *typeInfo) ***REMOVED***
	var l = rvGetSliceLen(rv)
	e.arrayStart(l)
	if l > 0 ***REMOVED***
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ ***REMOVED***
			e.arrayElem()
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		***REMOVED***
	***REMOVED***
	e.arrayEnd()
***REMOVED***

func (e *Encoder) kSeqWMbs(rv reflect.Value, ti *typeInfo) ***REMOVED***
	var l = rv.Len()
	if l == 0 ***REMOVED***
		e.mapStart(0)
	***REMOVED*** else ***REMOVED***
		if l%2 == 1 ***REMOVED***
			e.errorf("mapBySlice requires even slice length, but got %v", l)
			return
		***REMOVED***
		e.mapStart(l / 2)
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ ***REMOVED***
			if j%2 == 0 ***REMOVED***
				e.mapElemKey()
			***REMOVED*** else ***REMOVED***
				e.mapElemValue()
			***REMOVED***
			e.encodeValue(rv.Index(j), fn)
		***REMOVED***
	***REMOVED***
	e.mapEnd()
***REMOVED***

func (e *Encoder) kSeqW(rv reflect.Value, ti *typeInfo) ***REMOVED***
	var l = rv.Len()
	e.arrayStart(l)
	if l > 0 ***REMOVED***
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ ***REMOVED***
			e.arrayElem()
			e.encodeValue(rv.Index(j), fn)
		***REMOVED***
	***REMOVED***
	e.arrayEnd()
***REMOVED***

func (e *Encoder) kChan(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	if rvIsNil(rv) ***REMOVED***
		e.e.EncodeNil()
		return
	***REMOVED***
	if f.ti.chandir&uint8(reflect.RecvDir) == 0 ***REMOVED***
		e.errorf("send-only channel cannot be encoded")
		return
	***REMOVED***
	if !f.ti.mbs && uint8TypId == rt2id(f.ti.elem) ***REMOVED***
		e.kSliceBytesChan(rv)
		return
	***REMOVED***
	rtslice := reflect.SliceOf(f.ti.elem)
	rv = chanToSlice(rv, rtslice, e.h.ChanRecvTimeout)
	ti := e.h.getTypeInfo(rt2id(rtslice), rtslice)
	if f.ti.mbs ***REMOVED***
		e.kSliceWMbs(rv, ti)
	***REMOVED*** else ***REMOVED***
		e.kSliceW(rv, ti)
	***REMOVED***
***REMOVED***

func (e *Encoder) kSlice(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	if rvIsNil(rv) ***REMOVED***
		e.e.EncodeNil()
		return
	***REMOVED***
	if f.ti.mbs ***REMOVED***
		e.kSliceWMbs(rv, f.ti)
	***REMOVED*** else ***REMOVED***
		if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) ***REMOVED***
			e.e.EncodeStringBytesRaw(rvGetBytes(rv))
		***REMOVED*** else ***REMOVED***
			e.kSliceW(rv, f.ti)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Encoder) kArray(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	if f.ti.mbs ***REMOVED***
		e.kSeqWMbs(rv, f.ti)
	***REMOVED*** else ***REMOVED***
		if uint8TypId == rt2id(f.ti.elem) ***REMOVED***
			e.e.EncodeStringBytesRaw(rvGetArrayBytesRO(rv, e.b[:]))
		***REMOVED*** else ***REMOVED***
			e.kSeqW(rv, f.ti)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Encoder) kSliceBytesChan(rv reflect.Value) ***REMOVED***
	// do not use range, so that the number of elements encoded
	// does not change, and encoding does not hang waiting on someone to close chan.

	// for b := range rv2i(rv).(<-chan byte) ***REMOVED*** bs = append(bs, b) ***REMOVED***
	// ch := rv2i(rv).(<-chan byte) // fix error - that this is a chan byte, not a <-chan byte.

	bs := e.b[:0]
	irv := rv2i(rv)
	ch, ok := irv.(<-chan byte)
	if !ok ***REMOVED***
		ch = irv.(chan byte)
	***REMOVED***

L1:
	switch timeout := e.h.ChanRecvTimeout; ***REMOVED***
	case timeout == 0: // only consume available
		for ***REMOVED***
			select ***REMOVED***
			case b := <-ch:
				bs = append(bs, b)
			default:
				break L1
			***REMOVED***
		***REMOVED***
	case timeout > 0: // consume until timeout
		tt := time.NewTimer(timeout)
		for ***REMOVED***
			select ***REMOVED***
			case b := <-ch:
				bs = append(bs, b)
			case <-tt.C:
				// close(tt.C)
				break L1
			***REMOVED***
		***REMOVED***
	default: // consume until close
		for b := range ch ***REMOVED***
			bs = append(bs, b)
		***REMOVED***
	***REMOVED***

	e.e.EncodeStringBytesRaw(bs)
***REMOVED***

func (e *Encoder) kStructNoOmitempty(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	sfn := structFieldNode***REMOVED***v: rv, update: false***REMOVED***
	if f.ti.toArray || e.h.StructToArray ***REMOVED*** // toArray
		e.arrayStart(len(f.ti.sfiSrc))
		for _, si := range f.ti.sfiSrc ***REMOVED***
			e.arrayElem()
			e.encodeValue(sfn.field(si), nil)
		***REMOVED***
		e.arrayEnd()
	***REMOVED*** else ***REMOVED***
		e.mapStart(len(f.ti.sfiSort))
		for _, si := range f.ti.sfiSort ***REMOVED***
			e.mapElemKey()
			e.kStructFieldKey(f.ti.keyType, si.encNameAsciiAlphaNum, si.encName)
			e.mapElemValue()
			e.encodeValue(sfn.field(si), nil)
		***REMOVED***
		e.mapEnd()
	***REMOVED***
***REMOVED***

func (e *Encoder) kStructFieldKey(keyType valueType, encNameAsciiAlphaNum bool, encName string) ***REMOVED***
	encStructFieldKey(encName, e.e, e.w(), keyType, encNameAsciiAlphaNum, e.js)
***REMOVED***

func (e *Encoder) kStruct(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	var newlen int
	toMap := !(f.ti.toArray || e.h.StructToArray)
	var mf map[string]interface***REMOVED******REMOVED***
	if f.ti.isFlag(tiflagMissingFielder) ***REMOVED***
		mf = rv2i(rv).(MissingFielder).CodecMissingFields()
		toMap = true
		newlen += len(mf)
	***REMOVED*** else if f.ti.isFlag(tiflagMissingFielderPtr) ***REMOVED***
		if rv.CanAddr() ***REMOVED***
			mf = rv2i(rv.Addr()).(MissingFielder).CodecMissingFields()
		***REMOVED*** else ***REMOVED***
			// make a new addressable value of same one, and use it
			rv2 := reflect.New(rv.Type())
			rvSetDirect(rv2.Elem(), rv)
			mf = rv2i(rv2).(MissingFielder).CodecMissingFields()
		***REMOVED***
		toMap = true
		newlen += len(mf)
	***REMOVED***
	newlen += len(f.ti.sfiSrc)

	var fkvs = e.slist.get(newlen)

	recur := e.h.RecursiveEmptyCheck
	sfn := structFieldNode***REMOVED***v: rv, update: false***REMOVED***

	var kv sfiRv
	var j int
	if toMap ***REMOVED***
		newlen = 0
		for _, si := range f.ti.sfiSort ***REMOVED*** // use sorted array
			kv.r = sfn.field(si)
			if si.omitEmpty() && isEmptyValue(kv.r, e.h.TypeInfos, recur, recur) ***REMOVED***
				continue
			***REMOVED***
			kv.v = si // si.encName
			fkvs[newlen] = kv
			newlen++
		***REMOVED***
		var mflen int
		for k, v := range mf ***REMOVED***
			if k == "" ***REMOVED***
				delete(mf, k)
				continue
			***REMOVED***
			if f.ti.infoFieldOmitempty && isEmptyValue(rv4i(v), e.h.TypeInfos, recur, recur) ***REMOVED***
				delete(mf, k)
				continue
			***REMOVED***
			mflen++
		***REMOVED***
		// encode it all
		e.mapStart(newlen + mflen)
		for j = 0; j < newlen; j++ ***REMOVED***
			kv = fkvs[j]
			e.mapElemKey()
			e.kStructFieldKey(f.ti.keyType, kv.v.encNameAsciiAlphaNum, kv.v.encName)
			e.mapElemValue()
			e.encodeValue(kv.r, nil)
		***REMOVED***
		// now, add the others
		for k, v := range mf ***REMOVED***
			e.mapElemKey()
			e.kStructFieldKey(f.ti.keyType, false, k)
			e.mapElemValue()
			e.encode(v)
		***REMOVED***
		e.mapEnd()
	***REMOVED*** else ***REMOVED***
		newlen = len(f.ti.sfiSrc)
		for i, si := range f.ti.sfiSrc ***REMOVED*** // use unsorted array (to match sequence in struct)
			kv.r = sfn.field(si)
			// use the zero value.
			// if a reference or struct, set to nil (so you do not output too much)
			if si.omitEmpty() && isEmptyValue(kv.r, e.h.TypeInfos, recur, recur) ***REMOVED***
				switch kv.r.Kind() ***REMOVED***
				case reflect.Struct, reflect.Interface, reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice:
					kv.r = reflect.Value***REMOVED******REMOVED*** //encode as nil
				***REMOVED***
			***REMOVED***
			fkvs[i] = kv
		***REMOVED***
		// encode it all
		e.arrayStart(newlen)
		for j = 0; j < newlen; j++ ***REMOVED***
			e.arrayElem()
			e.encodeValue(fkvs[j].r, nil)
		***REMOVED***
		e.arrayEnd()
	***REMOVED***

	// do not use defer. Instead, use explicit pool return at end of function.
	// defer has a cost we are trying to avoid.
	// If there is a panic and these slices are not returned, it is ok.
	// spool.end()
	e.slist.put(fkvs)
***REMOVED***

func (e *Encoder) kMap(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	if rvIsNil(rv) ***REMOVED***
		e.e.EncodeNil()
		return
	***REMOVED***

	l := rv.Len()
	e.mapStart(l)
	if l == 0 ***REMOVED***
		e.mapEnd()
		return
	***REMOVED***

	// determine the underlying key and val encFn's for the map.
	// This eliminates some work which is done for each loop iteration i.e.
	// rv.Type(), ref.ValueOf(rt).Pointer(), then check map/list for fn.
	//
	// However, if kind is reflect.Interface, do not pre-determine the
	// encoding type, because preEncodeValue may break it down to
	// a concrete type and kInterface will bomb.

	var keyFn, valFn *codecFn

	ktypeKind := f.ti.key.Kind()
	vtypeKind := f.ti.elem.Kind()

	rtval := f.ti.elem
	rtvalkind := vtypeKind
	for rtvalkind == reflect.Ptr ***REMOVED***
		rtval = rtval.Elem()
		rtvalkind = rtval.Kind()
	***REMOVED***
	if rtvalkind != reflect.Interface ***REMOVED***
		valFn = e.h.fn(rtval)
	***REMOVED***

	var rvv = mapAddressableRV(f.ti.elem, vtypeKind)

	if e.h.Canonical ***REMOVED***
		e.kMapCanonical(f.ti.key, f.ti.elem, rv, rvv, valFn)
		e.mapEnd()
		return
	***REMOVED***

	rtkey := f.ti.key
	var keyTypeIsString = stringTypId == rt2id(rtkey) // rtkeyid
	if !keyTypeIsString ***REMOVED***
		for rtkey.Kind() == reflect.Ptr ***REMOVED***
			rtkey = rtkey.Elem()
		***REMOVED***
		if rtkey.Kind() != reflect.Interface ***REMOVED***
			keyFn = e.h.fn(rtkey)
		***REMOVED***
	***REMOVED***

	var rvk = mapAddressableRV(f.ti.key, ktypeKind)

	var it mapIter
	mapRange(&it, rv, rvk, rvv, true)
	validKV := it.ValidKV()
	var vx reflect.Value
	for it.Next() ***REMOVED***
		e.mapElemKey()
		if validKV ***REMOVED***
			vx = it.Key()
		***REMOVED*** else ***REMOVED***
			vx = rvk
		***REMOVED***
		if keyTypeIsString ***REMOVED***
			e.e.EncodeString(vx.String())
		***REMOVED*** else ***REMOVED***
			e.encodeValue(vx, keyFn)
		***REMOVED***
		e.mapElemValue()
		if validKV ***REMOVED***
			vx = it.Value()
		***REMOVED*** else ***REMOVED***
			vx = rvv
		***REMOVED***
		e.encodeValue(vx, valFn)
	***REMOVED***
	it.Done()

	e.mapEnd()
***REMOVED***

func (e *Encoder) kMapCanonical(rtkey, rtval reflect.Type, rv, rvv reflect.Value, valFn *codecFn) ***REMOVED***
	// we previously did out-of-band if an extension was registered.
	// This is not necessary, as the natural kind is sufficient for ordering.

	mks := rv.MapKeys()
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
			e.mapElemKey()
			e.e.EncodeBool(mksv[i].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv), valFn)
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
			e.mapElemKey()
			e.e.EncodeString(mksv[i].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv), valFn)
		***REMOVED***
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
		mksv := make([]uint64Rv, len(mks))
		for i, k := range mks ***REMOVED***
			v := &mksv[i]
			v.r = k
			v.v = k.Uint()
		***REMOVED***
		sort.Sort(uint64RvSlice(mksv))
		for i := range mksv ***REMOVED***
			e.mapElemKey()
			e.e.EncodeUint(mksv[i].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv), valFn)
		***REMOVED***
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		mksv := make([]int64Rv, len(mks))
		for i, k := range mks ***REMOVED***
			v := &mksv[i]
			v.r = k
			v.v = k.Int()
		***REMOVED***
		sort.Sort(int64RvSlice(mksv))
		for i := range mksv ***REMOVED***
			e.mapElemKey()
			e.e.EncodeInt(mksv[i].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv), valFn)
		***REMOVED***
	case reflect.Float32:
		mksv := make([]float64Rv, len(mks))
		for i, k := range mks ***REMOVED***
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		***REMOVED***
		sort.Sort(float64RvSlice(mksv))
		for i := range mksv ***REMOVED***
			e.mapElemKey()
			e.e.EncodeFloat32(float32(mksv[i].v))
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv), valFn)
		***REMOVED***
	case reflect.Float64:
		mksv := make([]float64Rv, len(mks))
		for i, k := range mks ***REMOVED***
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		***REMOVED***
		sort.Sort(float64RvSlice(mksv))
		for i := range mksv ***REMOVED***
			e.mapElemKey()
			e.e.EncodeFloat64(mksv[i].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv), valFn)
		***REMOVED***
	case reflect.Struct:
		if rtkey == timeTyp ***REMOVED***
			mksv := make([]timeRv, len(mks))
			for i, k := range mks ***REMOVED***
				v := &mksv[i]
				v.r = k
				v.v = rv2i(k).(time.Time)
			***REMOVED***
			sort.Sort(timeRvSlice(mksv))
			for i := range mksv ***REMOVED***
				e.mapElemKey()
				e.e.EncodeTime(mksv[i].v)
				e.mapElemValue()
				e.encodeValue(mapGet(rv, mksv[i].r, rvv), valFn)
			***REMOVED***
			break
		***REMOVED***
		fallthrough
	default:
		// out-of-band
		// first encode each key to a []byte first, then sort them, then record
		var mksv []byte = e.blist.get(len(mks) * 16)[:0]
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
			e.mapElemKey()
			e.encWr.writeb(mksbv[j].v) // e.asis(mksbv[j].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksbv[j].r, rvv), valFn)
		***REMOVED***
		e.blist.put(mksv)
	***REMOVED***
***REMOVED***

// Encoder writes an object to an output stream in a supported format.
//
// Encoder is NOT safe for concurrent use i.e. a Encoder cannot be used
// concurrently in multiple goroutines.
//
// However, as Encoder could be allocation heavy to initialize, a Reset method is provided
// so its state can be reused to decode new input streams repeatedly.
// This is the idiomatic way to use.
type Encoder struct ***REMOVED***
	panicHdl

	e encDriver

	h *BasicHandle

	// hopefully, reduce derefencing cost by laying the encWriter inside the Encoder
	encWr

	// ---- cpu cache line boundary
	hh Handle

	blist bytesFreelist
	err   error

	// ---- cpu cache line boundary

	// ---- writable fields during execution --- *try* to keep in sep cache line
	ci set // holds set of addresses found during an encoding (if CheckCircularRef=true)

	slist sfiRvFreelist

	b [(2 * 8)]byte // for encoding chan byte, (non-addressable) [N]byte, etc

	// ---- cpu cache line boundary?
***REMOVED***

// NewEncoder returns an Encoder for encoding into an io.Writer.
//
// For efficiency, Users are encouraged to configure WriterBufferSize on the handle
// OR pass in a memory buffered writer (eg bufio.Writer, bytes.Buffer).
func NewEncoder(w io.Writer, h Handle) *Encoder ***REMOVED***
	e := h.newEncDriver().encoder()
	e.Reset(w)
	return e
***REMOVED***

// NewEncoderBytes returns an encoder for encoding directly and efficiently
// into a byte slice, using zero-copying to temporary slices.
//
// It will potentially replace the output byte slice pointed to.
// After encoding, the out parameter contains the encoded contents.
func NewEncoderBytes(out *[]byte, h Handle) *Encoder ***REMOVED***
	e := h.newEncDriver().encoder()
	e.ResetBytes(out)
	return e
***REMOVED***

func (e *Encoder) init(h Handle) ***REMOVED***
	e.err = errEncoderNotInitialized
	e.bytes = true
	e.hh = h
	e.h = basicHandle(h)
	e.be = e.hh.isBinary()
***REMOVED***

func (e *Encoder) w() *encWr ***REMOVED***
	return &e.encWr
***REMOVED***

func (e *Encoder) resetCommon() ***REMOVED***
	e.e.reset()
	if e.ci == nil ***REMOVED***
		// e.ci = (set)(e.cidef[:0])
	***REMOVED*** else ***REMOVED***
		e.ci = e.ci[:0]
	***REMOVED***
	e.c = 0
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
	e.bytes = false
	if e.wf == nil ***REMOVED***
		e.wf = new(bufioEncWriter)
	***REMOVED***
	e.wf.reset(w, e.h.WriterBufferSize, &e.blist)
	e.resetCommon()
***REMOVED***

// ResetBytes resets the Encoder with a new destination output []byte.
func (e *Encoder) ResetBytes(out *[]byte) ***REMOVED***
	if out == nil ***REMOVED***
		return
	***REMOVED***
	var in []byte = *out
	if in == nil ***REMOVED***
		in = make([]byte, defEncByteBufSize)
	***REMOVED***
	e.bytes = true
	e.wb.reset(in, out)
	e.resetCommon()
***REMOVED***

// Encode writes an object into a stream.
//
// Encoding can be configured via the struct tag for the fields.
// The key (in the struct tags) that we look at is configurable.
//
// By default, we look up the "codec" key in the struct field's tags,
// and fall bak to the "json" key if "codec" is absent.
// That key in struct field's tag value is the key name,
// followed by an optional comma and options.
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
	// tried to use closure, as runtime optimizes defer with no params.
	// This seemed to be causing weird issues (like circular reference found, unexpected panic, etc).
	// Also, see https://github.com/golang/go/issues/14939#issuecomment-417836139
	// defer func() ***REMOVED*** e.deferred(&err) ***REMOVED***() ***REMOVED***
	// ***REMOVED*** x, y := e, &err; defer func() ***REMOVED*** x.deferred(y) ***REMOVED***() ***REMOVED***

	if e.err != nil ***REMOVED***
		return e.err
	***REMOVED***
	if recoverPanicToErr ***REMOVED***
		defer func() ***REMOVED***
			// if error occurred during encoding, return that error;
			// else if error occurred on end'ing (i.e. during flush), return that error.
			err = e.w().endErr()
			x := recover()
			if x == nil ***REMOVED***
				if e.err != err ***REMOVED***
					e.err = err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				panicValToErr(e, x, &e.err)
				if e.err != err ***REMOVED***
					err = e.err
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// defer e.deferred(&err)
	e.mustEncode(v)
	return
***REMOVED***

// MustEncode is like Encode, but panics if unable to Encode.
// This provides insight to the code location that triggered the error.
func (e *Encoder) MustEncode(v interface***REMOVED******REMOVED***) ***REMOVED***
	if e.err != nil ***REMOVED***
		panic(e.err)
	***REMOVED***
	e.mustEncode(v)
***REMOVED***

func (e *Encoder) mustEncode(v interface***REMOVED******REMOVED***) ***REMOVED***
	e.calls++
	e.encode(v)
	e.calls--
	if e.calls == 0 ***REMOVED***
		e.e.atEndOfEncode()
		e.w().end()
	***REMOVED***
***REMOVED***

// Release releases shared (pooled) resources.
//
// It is important to call Release() when done with an Encoder, so those resources
// are released instantly for use by subsequently created Encoders.
//
// Deprecated: Release is a no-op as pooled resources are not used with an Encoder.
// This method is kept for compatibility reasons only.
func (e *Encoder) Release() ***REMOVED***
***REMOVED***

func (e *Encoder) encode(iv interface***REMOVED******REMOVED***) ***REMOVED***
	// a switch with only concrete types can be optimized.
	// consequently, we deal with nil and interfaces outside the switch.

	if iv == nil ***REMOVED***
		e.e.EncodeNil()
		return
	***REMOVED***

	rv, ok := isNil(iv)
	if ok ***REMOVED***
		e.e.EncodeNil()
		return
	***REMOVED***

	var vself Selfer

	switch v := iv.(type) ***REMOVED***
	// case nil:
	// case Selfer:
	case Raw:
		e.rawBytes(v)
	case reflect.Value:
		e.encodeValue(v, nil)

	case string:
		e.e.EncodeString(v)
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
		e.e.EncodeStringBytesRaw(v)
	case *Raw:
		e.rawBytes(*v)
	case *string:
		e.e.EncodeString(*v)
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
		if *v == nil ***REMOVED***
			e.e.EncodeNil()
		***REMOVED*** else ***REMOVED***
			e.e.EncodeStringBytesRaw(*v)
		***REMOVED***
	default:
		if vself, ok = iv.(Selfer); ok ***REMOVED***
			vself.CodecEncodeSelf(e)
		***REMOVED*** else if !fastpathEncodeTypeSwitch(iv, e) ***REMOVED***
			if !rv.IsValid() ***REMOVED***
				rv = rv4i(iv)
			***REMOVED***
			e.encodeValue(rv, nil)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (e *Encoder) encodeValue(rv reflect.Value, fn *codecFn) ***REMOVED***
	// if a valid fn is passed, it MUST BE for the dereferenced type of rv

	// We considered using a uintptr (a pointer) retrievable via rv.UnsafeAddr.
	// However, it is possible for the same pointer to point to 2 different types e.g.
	//    type T struct ***REMOVED*** tHelper ***REMOVED***
	//    Here, for var v T; &v and &v.tHelper are the same pointer.
	// Consequently, we need a tuple of type and pointer, which interface***REMOVED******REMOVED*** natively provides.
	var sptr interface***REMOVED******REMOVED*** // uintptr
	var rvp reflect.Value
	var rvpValid bool
TOP:
	switch rv.Kind() ***REMOVED***
	case reflect.Ptr:
		if rvIsNil(rv) ***REMOVED***
			e.e.EncodeNil()
			return
		***REMOVED***
		rvpValid = true
		rvp = rv
		rv = rv.Elem()
		if e.h.CheckCircularRef && rv.Kind() == reflect.Struct ***REMOVED***
			sptr = rv2i(rvp) // rv.UnsafeAddr()
			break TOP
		***REMOVED***
		goto TOP
	case reflect.Interface:
		if rvIsNil(rv) ***REMOVED***
			e.e.EncodeNil()
			return
		***REMOVED***
		rv = rv.Elem()
		goto TOP
	case reflect.Slice, reflect.Map:
		if rvIsNil(rv) ***REMOVED***
			e.e.EncodeNil()
			return
		***REMOVED***
	case reflect.Invalid, reflect.Func:
		e.e.EncodeNil()
		return
	***REMOVED***

	if sptr != nil && (&e.ci).add(sptr) ***REMOVED***
		e.errorf("circular reference found: # %p, %T", sptr, sptr)
	***REMOVED***

	var rt reflect.Type
	if fn == nil ***REMOVED***
		rt = rv.Type()
		fn = e.h.fn(rt)
	***REMOVED***
	if fn.i.addrE ***REMOVED***
		if rvpValid ***REMOVED***
			fn.fe(e, &fn.i, rvp)
		***REMOVED*** else if rv.CanAddr() ***REMOVED***
			fn.fe(e, &fn.i, rv.Addr())
		***REMOVED*** else ***REMOVED***
			if rt == nil ***REMOVED***
				rt = rv.Type()
			***REMOVED***
			rv2 := reflect.New(rt)
			rvSetDirect(rv2.Elem(), rv)
			fn.fe(e, &fn.i, rv2)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fn.fe(e, &fn.i, rv)
	***REMOVED***
	if sptr != 0 ***REMOVED***
		(&e.ci).remove(sptr)
	***REMOVED***
***REMOVED***

func (e *Encoder) marshalUtf8(bs []byte, fnerr error) ***REMOVED***
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
	if bs == nil ***REMOVED***
		e.e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		e.e.EncodeString(stringView(bs))
		// e.e.EncodeStringEnc(cUTF8, stringView(bs))
	***REMOVED***
***REMOVED***

func (e *Encoder) marshalAsis(bs []byte, fnerr error) ***REMOVED***
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
	if bs == nil ***REMOVED***
		e.e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		e.encWr.writeb(bs) // e.asis(bs)
	***REMOVED***
***REMOVED***

func (e *Encoder) marshalRaw(bs []byte, fnerr error) ***REMOVED***
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
	if bs == nil ***REMOVED***
		e.e.EncodeNil()
	***REMOVED*** else ***REMOVED***
		e.e.EncodeStringBytesRaw(bs)
	***REMOVED***
***REMOVED***

func (e *Encoder) rawBytes(vv Raw) ***REMOVED***
	v := []byte(vv)
	if !e.h.Raw ***REMOVED***
		e.errorf("Raw values cannot be encoded: %v", v)
	***REMOVED***
	e.encWr.writeb(v) // e.asis(v)
***REMOVED***

func (e *Encoder) wrapErr(v interface***REMOVED******REMOVED***, err *error) ***REMOVED***
	*err = encodeError***REMOVED***codecError***REMOVED***name: e.hh.Name(), err: v***REMOVED******REMOVED***
***REMOVED***

// ---- container tracker methods
// Note: We update the .c after calling the callback.
// This way, the callback can know what the last status was.

func (e *Encoder) mapStart(length int) ***REMOVED***
	e.e.WriteMapStart(length)
	e.c = containerMapStart
***REMOVED***

func (e *Encoder) mapElemKey() ***REMOVED***
	if e.js ***REMOVED***
		e.jsondriver().WriteMapElemKey()
	***REMOVED***
	e.c = containerMapKey
***REMOVED***

func (e *Encoder) mapElemValue() ***REMOVED***
	if e.js ***REMOVED***
		e.jsondriver().WriteMapElemValue()
	***REMOVED***
	e.c = containerMapValue
***REMOVED***

func (e *Encoder) mapEnd() ***REMOVED***
	e.e.WriteMapEnd()
	// e.c = containerMapEnd
	e.c = 0
***REMOVED***

func (e *Encoder) arrayStart(length int) ***REMOVED***
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
***REMOVED***

func (e *Encoder) arrayElem() ***REMOVED***
	if e.js ***REMOVED***
		e.jsondriver().WriteArrayElem()
	***REMOVED***
	e.c = containerArrayElem
***REMOVED***

func (e *Encoder) arrayEnd() ***REMOVED***
	e.e.WriteArrayEnd()
	e.c = 0
	// e.c = containerArrayEnd
***REMOVED***

// ----------

func (e *Encoder) sideEncode(v interface***REMOVED******REMOVED***, bs *[]byte) ***REMOVED***
	rv := baseRV(v)
	e2 := NewEncoderBytes(bs, e.hh)
	e2.encodeValue(rv, e.h.fnNoExt(rv.Type()))
	e2.e.atEndOfEncode()
	e2.w().end()
***REMOVED***

func encStructFieldKey(encName string, ee encDriver, w *encWr,
	keyType valueType, encNameAsciiAlphaNum bool, js bool) ***REMOVED***
	var m must
	// use if-else-if, not switch (which compiles to binary-search)
	// since keyType is typically valueTypeString, branch prediction is pretty good.
	if keyType == valueTypeString ***REMOVED***
		if js && encNameAsciiAlphaNum ***REMOVED*** // keyType == valueTypeString
			w.writeqstr(encName)
		***REMOVED*** else ***REMOVED*** // keyType == valueTypeString
			ee.EncodeString(encName)
		***REMOVED***
	***REMOVED*** else if keyType == valueTypeInt ***REMOVED***
		ee.EncodeInt(m.Int(strconv.ParseInt(encName, 10, 64)))
	***REMOVED*** else if keyType == valueTypeUint ***REMOVED***
		ee.EncodeUint(m.Uint(strconv.ParseUint(encName, 10, 64)))
	***REMOVED*** else if keyType == valueTypeFloat ***REMOVED***
		ee.EncodeFloat64(m.Float(strconv.ParseFloat(encName, 64)))
	***REMOVED***
***REMOVED***
