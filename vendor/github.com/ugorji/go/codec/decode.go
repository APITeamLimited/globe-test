// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
)

// Some tagging information for error messages.
const (
	msgBadDesc = "unrecognized descriptor byte"
	// msgDecCannotExpandArr = "cannot expand go array from %v to stream length: %v"
)

const (
	decDefMaxDepth         = 1024 // maximum depth
	decDefSliceCap         = 8
	decDefChanCap          = 64      // should be large, as cap cannot be expanded
	decScratchByteArrayLen = (6 * 8) // ??? cacheLineSize +

	// decContainerLenUnknown is length returned from Read(Map|Array)Len
	// when a format doesn't know apiori.
	// For example, json doesn't pre-determine the length of a container (sequence/map).
	decContainerLenUnknown = -1

	// decContainerLenNil is length returned from Read(Map|Array)Len
	// when a 'nil' was encountered in the stream.
	decContainerLenNil = math.MinInt32

	// decFailNonEmptyIntf configures whether we error
	// when decoding naked into a non-empty interface.
	//
	// Typically, we cannot decode non-nil stream value into
	// nil interface with methods (e.g. io.Reader).
	// However, in some scenarios, this should be allowed:
	//   - MapType
	//   - SliceType
	//   - Extensions
	//
	// Consequently, we should relax this. Put it behind a const flag for now.
	decFailNonEmptyIntf = false
)

var (
	errstrOnlyMapOrArrayCanDecodeIntoStruct = "only encoded map or array can be decoded into a struct"
	errstrCannotDecodeIntoNil               = "cannot decode into nil"

	// errmsgExpandSliceOverflow     = "expand slice: slice overflow"
	errmsgExpandSliceCannotChange = "expand slice: cannot change"

	errDecoderNotInitialized = errors.New("Decoder not initialized")

	errDecUnreadByteNothingToRead   = errors.New("cannot unread - nothing has been read")
	errDecUnreadByteLastByteNotRead = errors.New("cannot unread - last byte has not been read")
	errDecUnreadByteUnknown         = errors.New("cannot unread - reason unknown")
	errMaxDepthExceeded             = errors.New("maximum decoding depth exceeded")

	errBytesDecReaderCannotUnread = errors.New("cannot unread last byte read")
)

type decDriver interface ***REMOVED***
	// this will check if the next token is a break.
	CheckBreak() bool

	// TryNil tries to decode as nil.
	TryNil() bool

	// ContainerType returns one of: Bytes, String, Nil, Slice or Map.
	//
	// Return unSet if not known.
	//
	// Note: Implementations MUST fully consume sentinel container types, specifically Nil.
	ContainerType() (vt valueType)

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

	DecodeInt64() (i int64)
	DecodeUint64() (ui uint64)

	DecodeFloat64() (f float64)
	DecodeBool() (b bool)

	// DecodeStringAsBytes returns the bytes representing a string.
	// By definition, it will return a view into a scratch buffer.
	//
	// Note: This can also decode symbols, if supported.
	//
	// Users should consume it right away and not store it for later use.
	DecodeStringAsBytes() (v []byte)

	// DecodeBytes may be called directly, without going through reflection.
	// Consequently, it must be designed to handle possible nil.
	DecodeBytes(bs []byte, zerocopy bool) (bsOut []byte)
	// DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte)

	// DecodeExt will decode into a *RawExt or into an extension.
	DecodeExt(v interface***REMOVED******REMOVED***, xtag uint64, ext Ext)
	// decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte)

	DecodeTime() (t time.Time)

	// ReadArrayStart will return the length of the array.
	// If the format doesn't prefix the length, it returns decContainerLenUnknown.
	// If the expected array was a nil in the stream, it returns decContainerLenNil.
	ReadArrayStart() int
	ReadArrayEnd()

	// ReadMapStart will return the length of the array.
	// If the format doesn't prefix the length, it returns decContainerLenUnknown.
	// If the expected array was a nil in the stream, it returns decContainerLenNil.
	ReadMapStart() int
	ReadMapEnd()

	reset()
	atEndOfDecode()
	uncacheRead()

	decoder() *Decoder
***REMOVED***

type decDriverContainerTracker interface ***REMOVED***
	ReadArrayElem()
	ReadMapElemKey()
	ReadMapElemValue()
***REMOVED***

type decodeError struct ***REMOVED***
	codecError
	pos int
***REMOVED***

func (d decodeError) Error() string ***REMOVED***
	return fmt.Sprintf("%s decode error [pos %d]: %v", d.name, d.pos, d.err)
***REMOVED***

type decDriverNoopContainerReader struct***REMOVED******REMOVED***

func (x decDriverNoopContainerReader) ReadArrayStart() (v int) ***REMOVED*** return ***REMOVED***
func (x decDriverNoopContainerReader) ReadArrayEnd()           ***REMOVED******REMOVED***
func (x decDriverNoopContainerReader) ReadMapStart() (v int)   ***REMOVED*** return ***REMOVED***
func (x decDriverNoopContainerReader) ReadMapEnd()             ***REMOVED******REMOVED***
func (x decDriverNoopContainerReader) CheckBreak() (v bool)    ***REMOVED*** return ***REMOVED***
func (x decDriverNoopContainerReader) atEndOfDecode()          ***REMOVED******REMOVED***

// DecodeOptions captures configuration options during decode.
type DecodeOptions struct ***REMOVED***
	// MapType specifies type to use during schema-less decoding of a map in the stream.
	// If nil (unset), we default to map[string]interface***REMOVED******REMOVED*** iff json handle and MapStringAsKey=true,
	// else map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***.
	MapType reflect.Type

	// SliceType specifies type to use during schema-less decoding of an array in the stream.
	// If nil (unset), we default to []interface***REMOVED******REMOVED*** for all formats.
	SliceType reflect.Type

	// MaxInitLen defines the maxinum initial length that we "make" a collection
	// (string, slice, map, chan). If 0 or negative, we default to a sensible value
	// based on the size of an element in the collection.
	//
	// For example, when decoding, a stream may say that it has 2^64 elements.
	// We should not auto-matically provision a slice of that size, to prevent Out-Of-Memory crash.
	// Instead, we provision up to MaxInitLen, fill that up, and start appending after that.
	MaxInitLen int

	// ReaderBufferSize is the size of the buffer used when reading.
	//
	// if > 0, we use a smart buffer internally for performance purposes.
	ReaderBufferSize int

	// MaxDepth defines the maximum depth when decoding nested
	// maps and slices. If 0 or negative, we default to a suitably large number (currently 1024).
	MaxDepth int16

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

	// SliceElementReset: on decoding a slice, reset the element to a zero value first.
	//
	// concern: if the slice already contained some garbage, we will decode into that garbage.
	SliceElementReset bool

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
	// Every string should not be interned.
	// An excellent use-case for interning is struct field names,
	// or map keys where key type is string.
	InternString bool

	// PreferArrayOverSlice controls whether to decode to an array or a slice.
	//
	// This only impacts decoding into a nil interface***REMOVED******REMOVED***.
	//
	// Consequently, it has no effect on codecgen.
	//
	// *Note*: This only applies if using go1.5 and above,
	// as it requires reflect.ArrayOf support which was absent before go1.5.
	PreferArrayOverSlice bool

	// DeleteOnNilMapValue controls how to decode a nil value in the stream.
	//
	// If true, we will delete the mapping of the key.
	// Else, just set the mapping to the zero value of the type.
	//
	// Deprecated: This does NOTHING and is left behind for compiling compatibility.
	// This change is necessitated because 'nil' in a stream now consistently
	// means the zero value (ie reset the value to its zero state).
	DeleteOnNilMapValue bool

	// RawToString controls how raw bytes in a stream are decoded into a nil interface***REMOVED******REMOVED***.
	// By default, they are decoded as []byte, but can be decoded as string (if configured).
	RawToString bool
***REMOVED***

// ----------------------------------------

func (d *Decoder) rawExt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	d.d.DecodeExt(rv2i(rv), 0, nil)
***REMOVED***

func (d *Decoder) ext(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	d.d.DecodeExt(rv2i(rv), f.xfTag, f.xfFn)
***REMOVED***

func (d *Decoder) selferUnmarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv2i(rv).(Selfer).CodecDecodeSelf(d)
***REMOVED***

func (d *Decoder) binaryUnmarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs := d.d.DecodeBytes(nil, true)
	if fnerr := bm.UnmarshalBinary(xbs); fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (d *Decoder) textUnmarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(d.d.DecodeStringAsBytes())
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (d *Decoder) jsonUnmarshal(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	tm := rv2i(rv).(jsonUnmarshaler)
	// bs := d.d.DecodeBytes(d.b[:], true, true)
	// grab the bytes to be read, as UnmarshalJSON needs the full JSON so as to unmarshal it itself.
	fnerr := tm.UnmarshalJSON(d.nextValueBytes())
	if fnerr != nil ***REMOVED***
		panic(fnerr)
	***REMOVED***
***REMOVED***

func (d *Decoder) kErr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	d.errorf("no decoding function defined for kind %v", rv.Kind())
***REMOVED***

func (d *Decoder) raw(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetBytes(rv, d.rawBytes())
***REMOVED***

func (d *Decoder) kString(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetString(rv, string(d.d.DecodeStringAsBytes()))
***REMOVED***

func (d *Decoder) kBool(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetBool(rv, d.d.DecodeBool())
***REMOVED***

func (d *Decoder) kTime(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetTime(rv, d.d.DecodeTime())
***REMOVED***

func (d *Decoder) kFloat32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetFloat32(rv, d.decodeFloat32())
***REMOVED***

func (d *Decoder) kFloat64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetFloat64(rv, d.d.DecodeFloat64())
***REMOVED***

func (d *Decoder) kInt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
***REMOVED***

func (d *Decoder) kInt8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
***REMOVED***

func (d *Decoder) kInt16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
***REMOVED***

func (d *Decoder) kInt32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
***REMOVED***

func (d *Decoder) kInt64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetInt64(rv, d.d.DecodeInt64())
***REMOVED***

func (d *Decoder) kUint(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
***REMOVED***

func (d *Decoder) kUintptr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
***REMOVED***

func (d *Decoder) kUint8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
***REMOVED***

func (d *Decoder) kUint16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
***REMOVED***

func (d *Decoder) kUint32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
***REMOVED***

func (d *Decoder) kUint64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rvSetUint64(rv, d.d.DecodeUint64())
***REMOVED***

func (d *Decoder) kInterfaceNaked(f *codecFnInfo) (rvn reflect.Value) ***REMOVED***
	// nil interface:
	// use some hieristics to decode it appropriately
	// based on the detected next value in the stream.
	n := d.naked()
	d.d.DecodeNaked()

	// We cannot decode non-nil stream value into nil interface with methods (e.g. io.Reader).
	// Howver, it is possible that the user has ways to pass in a type for a given interface
	//   - MapType
	//   - SliceType
	//   - Extensions
	//
	// Consequently, we should relax this. Put it behind a const flag for now.
	if decFailNonEmptyIntf && f.ti.numMeth > 0 ***REMOVED***
		d.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
		return
	***REMOVED***
	switch n.v ***REMOVED***
	case valueTypeMap:
		// if json, default to a map type with string keys
		mtid := d.mtid
		if mtid == 0 ***REMOVED***
			if d.jsms ***REMOVED***
				mtid = mapStrIntfTypId
			***REMOVED*** else ***REMOVED***
				mtid = mapIntfIntfTypId
			***REMOVED***
		***REMOVED***
		if mtid == mapIntfIntfTypId ***REMOVED***
			var v2 map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***
			d.decode(&v2)
			rvn = rv4i(&v2).Elem()
		***REMOVED*** else if mtid == mapStrIntfTypId ***REMOVED*** // for json performance
			var v2 map[string]interface***REMOVED******REMOVED***
			d.decode(&v2)
			rvn = rv4i(&v2).Elem()
		***REMOVED*** else ***REMOVED***
			if d.mtr ***REMOVED***
				rvn = reflect.New(d.h.MapType)
				d.decode(rv2i(rvn))
				rvn = rvn.Elem()
			***REMOVED*** else ***REMOVED***
				rvn = rvZeroAddrK(d.h.MapType, reflect.Map)
				d.decodeValue(rvn, nil)
			***REMOVED***
		***REMOVED***
	case valueTypeArray:
		if d.stid == 0 || d.stid == intfSliceTypId ***REMOVED***
			var v2 []interface***REMOVED******REMOVED***
			d.decode(&v2)
			rvn = rv4i(&v2).Elem()
		***REMOVED*** else ***REMOVED***
			if d.str ***REMOVED***
				rvn = reflect.New(d.h.SliceType)
				d.decode(rv2i(rvn))
				rvn = rvn.Elem()
			***REMOVED*** else ***REMOVED***
				rvn = rvZeroAddrK(d.h.SliceType, reflect.Slice)
				d.decodeValue(rvn, nil)
			***REMOVED***
		***REMOVED***
		if reflectArrayOfSupported && d.h.PreferArrayOverSlice ***REMOVED***
			rvn = rvGetArray4Slice(rvn)
		***REMOVED***
	case valueTypeExt:
		tag, bytes := n.u, n.l // calling decode below might taint the values
		bfn := d.h.getExtForTag(tag)
		var re = RawExt***REMOVED***Tag: tag***REMOVED***
		if bytes == nil ***REMOVED***
			// it is one of the InterfaceExt ones: json and cbor.
			// most likely cbor, as json decoding never reveals valueTypeExt (no tagging support)
			if bfn == nil ***REMOVED***
				d.decode(&re.Value)
				rvn = rv4i(&re).Elem()
			***REMOVED*** else ***REMOVED***
				if bfn.ext == SelfExt ***REMOVED***
					rvn = rvZeroAddrK(bfn.rt, bfn.rt.Kind())
					d.decodeValue(rvn, d.h.fnNoExt(bfn.rt))
				***REMOVED*** else ***REMOVED***
					rvn = reflect.New(bfn.rt)
					d.interfaceExtConvertAndDecode(rv2i(rvn), bfn.ext)
					rvn = rvn.Elem()
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// one of the BytesExt ones: binc, msgpack, simple
			if bfn == nil ***REMOVED***
				re.Data = detachZeroCopyBytes(d.bytes, nil, bytes)
				rvn = rv4i(&re).Elem()
			***REMOVED*** else ***REMOVED***
				rvn = reflect.New(bfn.rt)
				if bfn.ext == SelfExt ***REMOVED***
					d.sideDecode(rv2i(rvn), bytes)
				***REMOVED*** else ***REMOVED***
					bfn.ext.ReadExt(rv2i(rvn), bytes)
				***REMOVED***
				rvn = rvn.Elem()
			***REMOVED***
		***REMOVED***
	case valueTypeNil:
		// rvn = reflect.Zero(f.ti.rt)
		// no-op
	case valueTypeInt:
		rvn = n.ri()
	case valueTypeUint:
		rvn = n.ru()
	case valueTypeFloat:
		rvn = n.rf()
	case valueTypeBool:
		rvn = n.rb()
	case valueTypeString, valueTypeSymbol:
		rvn = n.rs()
	case valueTypeBytes:
		rvn = n.rl()
	case valueTypeTime:
		rvn = n.rt()
	default:
		panicv.errorf("kInterfaceNaked: unexpected valueType: %d", n.v)
	***REMOVED***
	return
***REMOVED***

func (d *Decoder) kInterface(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	// Note:
	// A consequence of how kInterface works, is that
	// if an interface already contains something, we try
	// to decode into what was there before.
	// We do not replace with a generic value (as got from decodeNaked).

	// every interface passed here MUST be settable.
	var rvn reflect.Value
	if rvIsNil(rv) || d.h.InterfaceReset ***REMOVED***
		// check if mapping to a type: if so, initialize it and move on
		rvn = d.h.intf2impl(f.ti.rtid)
		if rvn.IsValid() ***REMOVED***
			rv.Set(rvn)
		***REMOVED*** else ***REMOVED***
			rvn = d.kInterfaceNaked(f)
			// xdebugf("kInterface: %v", rvn)
			if rvn.IsValid() ***REMOVED***
				rv.Set(rvn)
			***REMOVED*** else if d.h.InterfaceReset ***REMOVED***
				// reset to zero value based on current type in there.
				if rvelem := rv.Elem(); rvelem.IsValid() ***REMOVED***
					rv.Set(reflect.Zero(rvelem.Type()))
				***REMOVED***
			***REMOVED***
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// now we have a non-nil interface value, meaning it contains a type
		rvn = rv.Elem()
	***REMOVED***

	// Note: interface***REMOVED******REMOVED*** is settable, but underlying type may not be.
	// Consequently, we MAY have to create a decodable value out of the underlying value,
	// decode into it, and reset the interface itself.
	// fmt.Printf(">>>> kInterface: rvn type: %v, rv type: %v\n", rvn.Type(), rv.Type())

	if isDecodeable(rvn) ***REMOVED***
		d.decodeValue(rvn, nil)
		return
	***REMOVED***

	rvn2 := rvZeroAddrK(rvn.Type(), rvn.Kind())
	rvSetDirect(rvn2, rvn)
	d.decodeValue(rvn2, nil)
	rv.Set(rvn2)
***REMOVED***

func decStructFieldKey(dd decDriver, keyType valueType, b *[decScratchByteArrayLen]byte) (rvkencname []byte) ***REMOVED***
	// use if-else-if, not switch (which compiles to binary-search)
	// since keyType is typically valueTypeString, branch prediction is pretty good.

	if keyType == valueTypeString ***REMOVED***
		rvkencname = dd.DecodeStringAsBytes()
	***REMOVED*** else if keyType == valueTypeInt ***REMOVED***
		rvkencname = strconv.AppendInt(b[:0], dd.DecodeInt64(), 10)
	***REMOVED*** else if keyType == valueTypeUint ***REMOVED***
		rvkencname = strconv.AppendUint(b[:0], dd.DecodeUint64(), 10)
	***REMOVED*** else if keyType == valueTypeFloat ***REMOVED***
		rvkencname = strconv.AppendFloat(b[:0], dd.DecodeFloat64(), 'f', -1, 64)
	***REMOVED*** else ***REMOVED***
		rvkencname = dd.DecodeStringAsBytes()
	***REMOVED***
	return
***REMOVED***

func (d *Decoder) kStruct(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	sfn := structFieldNode***REMOVED***v: rv, update: true***REMOVED***
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil ***REMOVED***
		rvSetDirect(rv, f.ti.rv0)
		return
	***REMOVED***
	var mf MissingFielder
	if f.ti.isFlag(tiflagMissingFielder) ***REMOVED***
		mf = rv2i(rv).(MissingFielder)
	***REMOVED*** else if f.ti.isFlag(tiflagMissingFielderPtr) ***REMOVED***
		mf = rv2i(rv.Addr()).(MissingFielder)
	***REMOVED***
	if ctyp == valueTypeMap ***REMOVED***
		containerLen := d.mapStart()
		if containerLen == 0 ***REMOVED***
			d.mapEnd()
			return
		***REMOVED***
		tisfi := f.ti.sfiSort
		hasLen := containerLen >= 0

		var rvkencname []byte
		for j := 0; (hasLen && j < containerLen) || !(hasLen || d.checkBreak()); j++ ***REMOVED***
			d.mapElemKey()
			rvkencname = decStructFieldKey(d.d, f.ti.keyType, &d.b)
			d.mapElemValue()
			if k := f.ti.indexForEncName(rvkencname); k > -1 ***REMOVED***
				si := tisfi[k]
				d.decodeValue(sfn.field(si), nil)
			***REMOVED*** else if mf != nil ***REMOVED***
				// store rvkencname in new []byte, as it previously shares Decoder.b, which is used in decode
				name2 := rvkencname
				rvkencname = make([]byte, len(rvkencname))
				copy(rvkencname, name2)

				var f interface***REMOVED******REMOVED***
				d.decode(&f)
				if !mf.CodecMissingField(rvkencname, f) && d.h.ErrorIfNoField ***REMOVED***
					d.errorf("no matching struct field found when decoding stream map with key: %s ",
						stringView(rvkencname))
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				d.structFieldNotFound(-1, stringView(rvkencname))
			***REMOVED***
			// keepAlive4StringView(rvkencnameB) // not needed, as reference is outside loop
		***REMOVED***
		d.mapEnd()
	***REMOVED*** else if ctyp == valueTypeArray ***REMOVED***
		containerLen := d.arrayStart()
		if containerLen == 0 ***REMOVED***
			d.arrayEnd()
			return
		***REMOVED***
		// Not much gain from doing it two ways for array.
		// Arrays are not used as much for structs.
		hasLen := containerLen >= 0
		var checkbreak bool
		for j, si := range f.ti.sfiSrc ***REMOVED***
			if hasLen && j == containerLen ***REMOVED***
				break
			***REMOVED***
			if !hasLen && d.checkBreak() ***REMOVED***
				checkbreak = true
				break
			***REMOVED***
			d.arrayElem()
			d.decodeValue(sfn.field(si), nil)
		***REMOVED***
		if (hasLen && containerLen > len(f.ti.sfiSrc)) || (!hasLen && !checkbreak) ***REMOVED***
			// read remaining values and throw away
			for j := len(f.ti.sfiSrc); ; j++ ***REMOVED***
				if (hasLen && j == containerLen) || (!hasLen && d.checkBreak()) ***REMOVED***
					break
				***REMOVED***
				d.arrayElem()
				d.structFieldNotFound(j, "")
			***REMOVED***
		***REMOVED***
		d.arrayEnd()
	***REMOVED*** else ***REMOVED***
		d.errorstr(errstrOnlyMapOrArrayCanDecodeIntoStruct)
		return
	***REMOVED***
***REMOVED***

func (d *Decoder) kSlice(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	// A slice can be set from a map or array in stream.
	// This way, the order can be kept (as order is lost with map).

	// Note: rv is a slice type here - guaranteed

	rtelem0 := f.ti.elem
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil ***REMOVED***
		if rv.CanSet() ***REMOVED***
			rvSetDirect(rv, f.ti.rv0)
		***REMOVED***
		return
	***REMOVED***
	if ctyp == valueTypeBytes || ctyp == valueTypeString ***REMOVED***
		// you can only decode bytes or string in the stream into a slice or array of bytes
		if !(f.ti.rtid == uint8SliceTypId || rtelem0.Kind() == reflect.Uint8) ***REMOVED***
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", f.ti.rt)
		***REMOVED***
		rvbs := rvGetBytes(rv)
		bs2 := d.d.DecodeBytes(rvbs, false)
		// if rvbs == nil && bs2 != nil || rvbs != nil && bs2 == nil || len(bs2) != len(rvbs) ***REMOVED***
		if !(len(bs2) > 0 && len(bs2) == len(rvbs) && &bs2[0] == &rvbs[0]) ***REMOVED***
			if rv.CanSet() ***REMOVED***
				rvSetBytes(rv, bs2)
			***REMOVED*** else if len(rvbs) > 0 && len(bs2) > 0 ***REMOVED***
				copy(rvbs, bs2)
			***REMOVED***
		***REMOVED***
		return
	***REMOVED***

	slh, containerLenS := d.decSliceHelperStart() // only expects valueType(Array|Map) - never Nil

	// an array can never return a nil slice. so no need to check f.array here.
	if containerLenS == 0 ***REMOVED***
		if rv.CanSet() ***REMOVED***
			if rvIsNil(rv) ***REMOVED***
				rvSetDirect(rv, reflect.MakeSlice(f.ti.rt, 0, 0))
			***REMOVED*** else ***REMOVED***
				rvSetSliceLen(rv, 0)
			***REMOVED***
		***REMOVED***
		slh.End()
		return
	***REMOVED***

	rtelem0Size := int(rtelem0.Size())
	rtElem0Kind := rtelem0.Kind()
	rtelem0Mut := !isImmutableKind(rtElem0Kind)
	rtelem := rtelem0
	rtelemkind := rtelem.Kind()
	for rtelemkind == reflect.Ptr ***REMOVED***
		rtelem = rtelem.Elem()
		rtelemkind = rtelem.Kind()
	***REMOVED***

	var fn *codecFn

	var rv0 = rv
	var rvChanged bool
	var rvCanset = rv.CanSet()
	var rv9 reflect.Value

	rvlen := rvGetSliceLen(rv)
	rvcap := rvGetSliceCap(rv)
	hasLen := containerLenS > 0
	if hasLen ***REMOVED***
		if containerLenS > rvcap ***REMOVED***
			oldRvlenGtZero := rvlen > 0
			rvlen = decInferLen(containerLenS, d.h.MaxInitLen, int(rtelem0.Size()))
			if rvlen <= rvcap ***REMOVED***
				if rvCanset ***REMOVED***
					rvSetSliceLen(rv, rvlen)
				***REMOVED***
			***REMOVED*** else if rvCanset ***REMOVED***
				rv = reflect.MakeSlice(f.ti.rt, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = true
			***REMOVED*** else ***REMOVED***
				d.errorf("cannot decode into non-settable slice")
			***REMOVED***
			if rvChanged && oldRvlenGtZero && rtelem0Mut ***REMOVED*** // !isImmutableKind(rtelem0.Kind()) ***REMOVED***
				rvCopySlice(rv, rv0) // only copy up to length NOT cap i.e. rv0.Slice(0, rvcap)
			***REMOVED***
		***REMOVED*** else if containerLenS != rvlen ***REMOVED***
			rvlen = containerLenS
			if rvCanset ***REMOVED***
				rvSetSliceLen(rv, rvlen)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// consider creating new element once, and just decoding into it.
	var rtelem0Zero reflect.Value
	var rtelem0ZeroValid bool
	var j int

	for ; (hasLen && j < containerLenS) || !(hasLen || d.checkBreak()); j++ ***REMOVED***
		if j == 0 && f.seq == seqTypeSlice && rvIsNil(rv) ***REMOVED***
			if hasLen ***REMOVED***
				rvlen = decInferLen(containerLenS, d.h.MaxInitLen, rtelem0Size)
			***REMOVED*** else ***REMOVED***
				rvlen = decDefSliceCap
			***REMOVED***
			if rvCanset ***REMOVED***
				rv = reflect.MakeSlice(f.ti.rt, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = true
			***REMOVED*** else ***REMOVED***
				d.errorf("cannot decode into non-settable slice")
			***REMOVED***
		***REMOVED***
		slh.ElemContainerState(j)
		// if indefinite, etc, then expand the slice if necessary
		if j >= rvlen ***REMOVED***
			if f.seq == seqTypeArray ***REMOVED***
				d.arrayCannotExpand(rvlen, j+1)
				// drain completely and return
				d.swallow()
				j++
				for ; (hasLen && j < containerLenS) || !(hasLen || d.checkBreak()); j++ ***REMOVED***
					slh.ElemContainerState(j)
					d.swallow()
				***REMOVED***
				slh.End()
				return
			***REMOVED***
			// rv = reflect.Append(rv, reflect.Zero(rtelem0)) // append logic + varargs

			// expand the slice up to the cap.
			// Note that we did, so we have to reset it later.

			if rvlen < rvcap ***REMOVED***
				if rv.CanSet() ***REMOVED***
					rvSetSliceLen(rv, rvcap)
				***REMOVED*** else if rvCanset ***REMOVED***
					rv = rvSlice(rv, rvcap)
					rvChanged = true
				***REMOVED*** else ***REMOVED***
					d.errorf(errmsgExpandSliceCannotChange)
					return
				***REMOVED***
				rvlen = rvcap
			***REMOVED*** else ***REMOVED***
				if !rvCanset ***REMOVED***
					d.errorf(errmsgExpandSliceCannotChange)
					return
				***REMOVED***
				rvcap = growCap(rvcap, rtelem0Size, rvcap)
				rv9 = reflect.MakeSlice(f.ti.rt, rvcap, rvcap)
				rvCopySlice(rv9, rv)
				rv = rv9
				rvChanged = true
				rvlen = rvcap
			***REMOVED***
		***REMOVED***
		rv9 = rvSliceIndex(rv, j, f.ti)
		if d.h.SliceElementReset ***REMOVED***
			if !rtelem0ZeroValid ***REMOVED***
				rtelem0ZeroValid = true
				rtelem0Zero = reflect.Zero(rtelem0)
			***REMOVED***
			rv9.Set(rtelem0Zero)
		***REMOVED***

		if fn == nil ***REMOVED***
			fn = d.h.fn(rtelem)
		***REMOVED***
		d.decodeValue(rv9, fn)
	***REMOVED***
	if j < rvlen ***REMOVED***
		if rv.CanSet() ***REMOVED***
			rvSetSliceLen(rv, j)
		***REMOVED*** else if rvCanset ***REMOVED***
			rv = rvSlice(rv, j)
			rvChanged = true
		***REMOVED***
		rvlen = j
	***REMOVED*** else if j == 0 && rvIsNil(rv) ***REMOVED***
		if rvCanset ***REMOVED***
			rv = reflect.MakeSlice(f.ti.rt, 0, 0)
			rvChanged = true
		***REMOVED***
	***REMOVED***
	slh.End()

	if rvChanged ***REMOVED*** // infers rvCanset=true, so it can be reset
		rv0.Set(rv)
	***REMOVED***

***REMOVED***

func (d *Decoder) kSliceForChan(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	// A slice can be set from a map or array in stream.
	// This way, the order can be kept (as order is lost with map).

	if f.ti.chandir&uint8(reflect.SendDir) == 0 ***REMOVED***
		d.errorf("receive-only channel cannot be decoded")
	***REMOVED***
	rtelem0 := f.ti.elem
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeNil ***REMOVED***
		rvSetDirect(rv, f.ti.rv0)
		return
	***REMOVED***
	if ctyp == valueTypeBytes || ctyp == valueTypeString ***REMOVED***
		// you can only decode bytes or string in the stream into a slice or array of bytes
		if !(f.ti.rtid == uint8SliceTypId || rtelem0.Kind() == reflect.Uint8) ***REMOVED***
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", f.ti.rt)
		***REMOVED***
		bs2 := d.d.DecodeBytes(nil, true)
		irv := rv2i(rv)
		ch, ok := irv.(chan<- byte)
		if !ok ***REMOVED***
			ch = irv.(chan byte)
		***REMOVED***
		for _, b := range bs2 ***REMOVED***
			ch <- b
		***REMOVED***
		return
	***REMOVED***

	// only expects valueType(Array|Map - nil handled above)
	slh, containerLenS := d.decSliceHelperStart()

	// an array can never return a nil slice. so no need to check f.array here.
	if containerLenS == 0 ***REMOVED***
		if rv.CanSet() && rvIsNil(rv) ***REMOVED***
			rvSetDirect(rv, reflect.MakeChan(f.ti.rt, 0))
		***REMOVED***
		slh.End()
		return
	***REMOVED***

	rtelem0Size := int(rtelem0.Size())
	rtElem0Kind := rtelem0.Kind()
	rtelem0Mut := !isImmutableKind(rtElem0Kind)
	rtelem := rtelem0
	rtelemkind := rtelem.Kind()
	for rtelemkind == reflect.Ptr ***REMOVED***
		rtelem = rtelem.Elem()
		rtelemkind = rtelem.Kind()
	***REMOVED***

	var fn *codecFn

	var rvCanset = rv.CanSet()
	var rvChanged bool
	var rv0 = rv
	var rv9 reflect.Value

	var rvlen int // := rv.Len()
	hasLen := containerLenS > 0

	var j int

	for ; (hasLen && j < containerLenS) || !(hasLen || d.checkBreak()); j++ ***REMOVED***
		if j == 0 && rvIsNil(rv) ***REMOVED***
			if hasLen ***REMOVED***
				rvlen = decInferLen(containerLenS, d.h.MaxInitLen, rtelem0Size)
			***REMOVED*** else ***REMOVED***
				rvlen = decDefChanCap
			***REMOVED***
			if rvCanset ***REMOVED***
				rv = reflect.MakeChan(f.ti.rt, rvlen)
				rvChanged = true
			***REMOVED*** else ***REMOVED***
				d.errorf("cannot decode into non-settable chan")
			***REMOVED***
		***REMOVED***
		slh.ElemContainerState(j)
		if rtelem0Mut || !rv9.IsValid() ***REMOVED*** // || (rtElem0Kind == reflect.Ptr && rvIsNil(rv9)) ***REMOVED***
			rv9 = rvZeroAddrK(rtelem0, rtElem0Kind)
		***REMOVED***
		if fn == nil ***REMOVED***
			fn = d.h.fn(rtelem)
		***REMOVED***
		d.decodeValue(rv9, fn)
		rv.Send(rv9)
	***REMOVED***
	slh.End()

	if rvChanged ***REMOVED*** // infers rvCanset=true, so it can be reset
		rv0.Set(rv)
	***REMOVED***

***REMOVED***

func (d *Decoder) kMap(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	containerLen := d.mapStart()
	if containerLen == decContainerLenNil ***REMOVED***
		rvSetDirect(rv, f.ti.rv0)
		return
	***REMOVED***
	ti := f.ti
	if rvIsNil(rv) ***REMOVED***
		rvlen := decInferLen(containerLen, d.h.MaxInitLen, int(ti.key.Size()+ti.elem.Size()))
		rvSetDirect(rv, makeMapReflect(ti.rt, rvlen))
	***REMOVED***

	if containerLen == 0 ***REMOVED***
		d.mapEnd()
		return
	***REMOVED***

	ktype, vtype := ti.key, ti.elem
	ktypeId := rt2id(ktype)
	vtypeKind := vtype.Kind()
	ktypeKind := ktype.Kind()

	var vtypeElem reflect.Type

	var keyFn, valFn *codecFn
	var ktypeLo, vtypeLo reflect.Type

	for ktypeLo = ktype; ktypeLo.Kind() == reflect.Ptr; ktypeLo = ktypeLo.Elem() ***REMOVED***
	***REMOVED***

	for vtypeLo = vtype; vtypeLo.Kind() == reflect.Ptr; vtypeLo = vtypeLo.Elem() ***REMOVED***
	***REMOVED***

	rvvMut := !isImmutableKind(vtypeKind)

	// we do a doMapGet if kind is mutable, and InterfaceReset=true if interface
	var doMapGet, doMapSet bool
	if !d.h.MapValueReset ***REMOVED***
		if rvvMut ***REMOVED***
			if vtypeKind == reflect.Interface ***REMOVED***
				if !d.h.InterfaceReset ***REMOVED***
					doMapGet = true
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				doMapGet = true
			***REMOVED***
		***REMOVED***
	***REMOVED***

	var rvk, rvkn, rvv, rvvn, rvva reflect.Value
	var rvvaSet bool
	rvkMut := !isImmutableKind(ktype.Kind()) // if ktype is immutable, then re-use the same rvk.
	ktypeIsString := ktypeId == stringTypId
	ktypeIsIntf := ktypeId == intfTypId
	hasLen := containerLen > 0
	var kstrbs []byte

	for j := 0; (hasLen && j < containerLen) || !(hasLen || d.checkBreak()); j++ ***REMOVED***
		if j == 0 ***REMOVED***
			if !rvkMut ***REMOVED***
				rvkn = rvZeroAddrK(ktype, ktypeKind)
			***REMOVED***
			if !rvvMut ***REMOVED***
				rvvn = rvZeroAddrK(vtype, vtypeKind)
			***REMOVED***
		***REMOVED***

		if rvkMut ***REMOVED***
			rvk = rvZeroAddrK(ktype, ktypeKind)
		***REMOVED*** else ***REMOVED***
			rvk = rvkn
		***REMOVED***

		d.mapElemKey()

		if ktypeIsString ***REMOVED***
			kstrbs = d.d.DecodeStringAsBytes()
			rvk.SetString(stringView(kstrbs)) // NOTE: if doing an insert, use real string (not stringview)
		***REMOVED*** else ***REMOVED***
			if keyFn == nil ***REMOVED***
				keyFn = d.h.fn(ktypeLo)
			***REMOVED***
			d.decodeValue(rvk, keyFn)
		***REMOVED***

		// special case if interface wrapping a byte array.
		if ktypeIsIntf ***REMOVED***
			if rvk2 := rvk.Elem(); rvk2.IsValid() && rvk2.Type() == uint8SliceTyp ***REMOVED***
				rvk.Set(rv4i(d.string(rvGetBytes(rvk2))))
			***REMOVED***
			// NOTE: consider failing early if map/slice/func
		***REMOVED***

		d.mapElemValue()

		doMapSet = true // set to false if u do a get, and its a non-nil pointer
		if doMapGet ***REMOVED***
			if !rvvaSet ***REMOVED***
				rvva = mapAddressableRV(vtype, vtypeKind)
				rvvaSet = true
			***REMOVED***
			rvv = mapGet(rv, rvk, rvva) // reflect.Value***REMOVED******REMOVED***)
			if vtypeKind == reflect.Ptr ***REMOVED***
				if rvv.IsValid() && !rvIsNil(rvv) ***REMOVED***
					doMapSet = false
				***REMOVED*** else ***REMOVED***
					if vtypeElem == nil ***REMOVED***
						vtypeElem = vtype.Elem()
					***REMOVED***
					rvv = reflect.New(vtypeElem)
				***REMOVED***
			***REMOVED*** else if rvv.IsValid() && vtypeKind == reflect.Interface && !rvIsNil(rvv) ***REMOVED***
				rvvn = rvZeroAddrK(vtype, vtypeKind)
				rvvn.Set(rvv)
				rvv = rvvn
			***REMOVED*** else if rvvMut ***REMOVED***
				rvv = rvZeroAddrK(vtype, vtypeKind)
			***REMOVED*** else ***REMOVED***
				rvv = rvvn
			***REMOVED***
		***REMOVED*** else if rvvMut ***REMOVED***
			rvv = rvZeroAddrK(vtype, vtypeKind)
		***REMOVED*** else ***REMOVED***
			rvv = rvvn
		***REMOVED***

		if valFn == nil ***REMOVED***
			valFn = d.h.fn(vtypeLo)
		***REMOVED***

		// We MUST be done with the stringview of the key, BEFORE decoding the value (rvv)
		// so that we don't unknowingly reuse the rvk backing buffer during rvv decode.
		if doMapSet && ktypeIsString ***REMOVED*** // set to a real string (not string view)
			rvk.SetString(d.string(kstrbs))
		***REMOVED***
		d.decodeValue(rvv, valFn)
		if doMapSet ***REMOVED***
			mapSet(rv, rvk, rvv)
		***REMOVED***
	***REMOVED***

	d.mapEnd()

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

	// primitives below
	u uint64
	i int64
	f float64
	l []byte
	s string

	// ---- cpu cache line boundary?
	t time.Time
	b bool

	// state
	v valueType
***REMOVED***

// Decoder reads and decodes an object from an input stream in a supported format.
//
// Decoder is NOT safe for concurrent use i.e. a Decoder cannot be used
// concurrently in multiple goroutines.
//
// However, as Decoder could be allocation heavy to initialize, a Reset method is provided
// so its state can be reused to decode new input streams repeatedly.
// This is the idiomatic way to use.
type Decoder struct ***REMOVED***
	panicHdl
	// hopefully, reduce derefencing cost by laying the decReader inside the Decoder.
	// Try to put things that go together to fit within a cache line (8 words).

	d decDriver

	// cache the mapTypeId and sliceTypeId for faster comparisons
	mtid uintptr
	stid uintptr

	h *BasicHandle

	blist bytesFreelist

	// ---- cpu cache line boundary?
	decRd

	// ---- cpu cache line boundary?
	n decNaked

	hh  Handle
	err error

	// ---- cpu cache line boundary?
	is map[string]string // used for interning strings

	// ---- writable fields during execution --- *try* to keep in sep cache line
	maxdepth int16
	depth    int16

	// Extensions can call Decode() within a current Decode() call.
	// We need to know when the top level Decode() call returns,
	// so we can decide whether to Release() or not.
	calls uint16 // what depth in mustDecode are we in now.

	c containerState
	_ [1]byte // padding

	// ---- cpu cache line boundary?

	// b is an always-available scratch buffer used by Decoder and decDrivers.
	// By being always-available, it can be used for one-off things without
	// having to get from freelist, use, and return back to freelist.
	b [decScratchByteArrayLen]byte
***REMOVED***

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
//
// For efficiency, Users are encouraged to configure ReaderBufferSize on the handle
// OR pass in a memory buffered reader (eg bufio.Reader, bytes.Buffer).
func NewDecoder(r io.Reader, h Handle) *Decoder ***REMOVED***
	d := h.newDecDriver().decoder()
	d.Reset(r)
	return d
***REMOVED***

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) *Decoder ***REMOVED***
	d := h.newDecDriver().decoder()
	d.ResetBytes(in)
	return d
***REMOVED***

func (d *Decoder) r() *decRd ***REMOVED***
	return &d.decRd
***REMOVED***

func (d *Decoder) init(h Handle) ***REMOVED***
	d.bytes = true
	d.err = errDecoderNotInitialized
	d.h = basicHandle(h)
	d.hh = h
	d.be = h.isBinary()
	// NOTE: do not initialize d.n here. It is lazily initialized in d.naked()
	if d.h.InternString ***REMOVED***
		d.is = make(map[string]string, 32)
	***REMOVED***
***REMOVED***

func (d *Decoder) resetCommon() ***REMOVED***
	d.d.reset()
	d.err = nil
	d.depth = 0
	d.calls = 0
	d.maxdepth = d.h.MaxDepth
	if d.maxdepth <= 0 ***REMOVED***
		d.maxdepth = decDefMaxDepth
	***REMOVED***
	// reset all things which were cached from the Handle, but could change
	d.mtid, d.stid = 0, 0
	d.mtr, d.str = false, false
	if d.h.MapType != nil ***REMOVED***
		d.mtid = rt2id(d.h.MapType)
		d.mtr = fastpathAV.index(d.mtid) != -1
	***REMOVED***
	if d.h.SliceType != nil ***REMOVED***
		d.stid = rt2id(d.h.SliceType)
		d.str = fastpathAV.index(d.stid) != -1
	***REMOVED***
***REMOVED***

// Reset the Decoder with a new Reader to decode from,
// clearing all state from last run(s).
func (d *Decoder) Reset(r io.Reader) ***REMOVED***
	if r == nil ***REMOVED***
		return
	***REMOVED***
	d.bytes = false
	if d.h.ReaderBufferSize > 0 ***REMOVED***
		if d.bi == nil ***REMOVED***
			d.bi = new(bufioDecReader)
		***REMOVED***
		d.bi.reset(r, d.h.ReaderBufferSize, &d.blist)
		d.bufio = true
	***REMOVED*** else ***REMOVED***
		if d.ri == nil ***REMOVED***
			d.ri = new(ioDecReader)
		***REMOVED***
		d.ri.reset(r, &d.blist)
		d.bufio = false
	***REMOVED***
	d.resetCommon()
***REMOVED***

// ResetBytes resets the Decoder with a new []byte to decode from,
// clearing all state from last run(s).
func (d *Decoder) ResetBytes(in []byte) ***REMOVED***
	if in == nil ***REMOVED***
		return
	***REMOVED***
	d.bytes = true
	d.bufio = false
	d.rb.reset(in)
	d.resetCommon()
***REMOVED***

func (d *Decoder) naked() *decNaked ***REMOVED***
	return &d.n
***REMOVED***

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
// Decode will typically use the stream contents to UPDATE the container i.e. the values
// in these containers will not be zero'ed before decoding.
//   - A map can be decoded from a stream map, by updating matching keys.
//   - A slice can be decoded from a stream array,
//     by updating the first n elements, where n is length of the stream.
//   - A slice can be decoded from a stream map, by decoding as if
//     it contains a sequence of key-value pairs.
//   - A struct can be decoded from a stream map, by updating matching fields.
//   - A struct can be decoded from a stream array,
//     by updating fields as they occur in the struct (by index).
//
// This in-place update maintains consistency in the decoding philosophy (i.e. we ALWAYS update
// in place by default). However, the consequence of this is that values in slices or maps
// which are not zero'ed before hand, will have part of the prior values in place after decode
// if the stream doesn't contain an update for those parts.
//
// This in-place update can be disabled by configuring the MapValueReset and SliceElementReset
// decode options available on every handle.
//
// Furthermore, when decoding a stream map or array with length of 0 into a nil map or slice,
// we reset the destination map or slice to a zero-length value.
//
// However, when decoding a stream nil, we reset the destination container
// to its "zero" value (e.g. nil for slice/map, etc).
//
// Note: we allow nil values in the stream anywhere except for map keys.
// A nil value in the encoded stream where a map key is expected is treated as an error.
func (d *Decoder) Decode(v interface***REMOVED******REMOVED***) (err error) ***REMOVED***
	// tried to use closure, as runtime optimizes defer with no params.
	// This seemed to be causing weird issues (like circular reference found, unexpected panic, etc).
	// Also, see https://github.com/golang/go/issues/14939#issuecomment-417836139
	// defer func() ***REMOVED*** d.deferred(&err) ***REMOVED***()
	// ***REMOVED*** x, y := d, &err; defer func() ***REMOVED*** x.deferred(y) ***REMOVED***() ***REMOVED***
	if d.err != nil ***REMOVED***
		return d.err
	***REMOVED***
	if recoverPanicToErr ***REMOVED***
		defer func() ***REMOVED***
			if x := recover(); x != nil ***REMOVED***
				panicValToErr(d, x, &d.err)
				if d.err != err ***REMOVED***
					err = d.err
				***REMOVED***
			***REMOVED***
		***REMOVED***()
	***REMOVED***

	// defer d.deferred(&err)
	d.mustDecode(v)
	return
***REMOVED***

// MustDecode is like Decode, but panics if unable to Decode.
// This provides insight to the code location that triggered the error.
func (d *Decoder) MustDecode(v interface***REMOVED******REMOVED***) ***REMOVED***
	if d.err != nil ***REMOVED***
		panic(d.err)
	***REMOVED***
	d.mustDecode(v)
***REMOVED***

// MustDecode is like Decode, but panics if unable to Decode.
// This provides insight to the code location that triggered the error.
func (d *Decoder) mustDecode(v interface***REMOVED******REMOVED***) ***REMOVED***
	// Top-level: v is a pointer and not nil.

	d.calls++
	d.decode(v)
	d.calls--
	if d.calls == 0 ***REMOVED***
		d.d.atEndOfDecode()
	***REMOVED***
***REMOVED***

// Release releases shared (pooled) resources.
//
// It is important to call Release() when done with a Decoder, so those resources
// are released instantly for use by subsequently created Decoders.
//
// By default, Release() is automatically called unless the option ExplicitRelease is set.
//
// Deprecated: Release is a no-op as pooled resources are not used with an Decoder.
// This method is kept for compatibility reasons only.
func (d *Decoder) Release() ***REMOVED***
***REMOVED***

func (d *Decoder) swallow() ***REMOVED***
	switch d.d.ContainerType() ***REMOVED***
	case valueTypeNil:
	case valueTypeMap:
		containerLen := d.mapStart()
		hasLen := containerLen >= 0
		for j := 0; (hasLen && j < containerLen) || !(hasLen || d.checkBreak()); j++ ***REMOVED***
			d.mapElemKey()
			d.swallow()
			d.mapElemValue()
			d.swallow()
		***REMOVED***
		d.mapEnd()
	case valueTypeArray:
		containerLen := d.arrayStart()
		hasLen := containerLen >= 0
		for j := 0; (hasLen && j < containerLen) || !(hasLen || d.checkBreak()); j++ ***REMOVED***
			d.arrayElem()
			d.swallow()
		***REMOVED***
		d.arrayEnd()
	case valueTypeBytes:
		d.d.DecodeBytes(d.b[:], true)
	case valueTypeString:
		d.d.DecodeStringAsBytes()
	default:
		// these are all primitives, which we can get from decodeNaked
		// if RawExt using Value, complete the processing.
		n := d.naked()
		d.d.DecodeNaked()
		if n.v == valueTypeExt && n.l == nil ***REMOVED***
			var v2 interface***REMOVED******REMOVED***
			d.decode(&v2)
		***REMOVED***
	***REMOVED***
***REMOVED***

func setZero(iv interface***REMOVED******REMOVED***) ***REMOVED***
	if iv == nil ***REMOVED***
		return
	***REMOVED***
	if _, ok := isNil(iv); ok ***REMOVED***
		return
	***REMOVED***
	// var canDecode bool
	switch v := iv.(type) ***REMOVED***
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
	case *time.Time:
		*v = time.Time***REMOVED******REMOVED***
	case reflect.Value:
		setZeroRV(v)
	default:
		if !fastpathDecodeSetZeroTypeSwitch(iv) ***REMOVED***
			setZeroRV(rv4i(iv))
		***REMOVED***
	***REMOVED***
***REMOVED***

func setZeroRV(v reflect.Value) ***REMOVED***
	// It not decodeable, we do not touch it.
	// We considered empty'ing it if not decodeable e.g.
	//    - if chan, drain it
	//    - if map, clear it
	//    - if slice or array, zero all elements up to len
	//
	// However, we decided instead that we either will set the
	// whole value to the zero value, or leave AS IS.
	if isDecodeable(v) ***REMOVED***
		if v.Kind() == reflect.Ptr ***REMOVED***
			v = v.Elem()
		***REMOVED***
		if v.CanSet() ***REMOVED***
			v.Set(reflect.Zero(v.Type()))
		***REMOVED***
	***REMOVED***
***REMOVED***

func (d *Decoder) decode(iv interface***REMOVED******REMOVED***) ***REMOVED***
	// a switch with only concrete types can be optimized.
	// consequently, we deal with nil and interfaces outside the switch.

	if iv == nil ***REMOVED***
		d.errorstr(errstrCannotDecodeIntoNil)
		return
	***REMOVED***

	switch v := iv.(type) ***REMOVED***
	// case nil:
	// case Selfer:
	case reflect.Value:
		d.ensureDecodeable(v)
		d.decodeValue(v, nil)

	case *string:
		*v = string(d.d.DecodeStringAsBytes())
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	case *int8:
		*v = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
	case *int16:
		*v = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
	case *int32:
		*v = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	case *int64:
		*v = d.d.DecodeInt64()
	case *uint:
		*v = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *uint8:
		*v = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	case *uint16:
		*v = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
	case *uint32:
		*v = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
	case *uint64:
		*v = d.d.DecodeUint64()
	case *float32:
		*v = float32(d.decodeFloat32())
	case *float64:
		*v = d.d.DecodeFloat64()
	case *[]uint8:
		*v = d.d.DecodeBytes(*v, false)
	case []uint8:
		b := d.d.DecodeBytes(v, false)
		if !(len(b) > 0 && len(b) == len(v) && &b[0] == &v[0]) ***REMOVED***
			copy(v, b)
		***REMOVED***
	case *time.Time:
		*v = d.d.DecodeTime()
	case *Raw:
		*v = d.rawBytes()

	case *interface***REMOVED******REMOVED***:
		d.decodeValue(rv4i(iv), nil)

	default:
		if v, ok := iv.(Selfer); ok ***REMOVED***
			v.CodecDecodeSelf(d)
		***REMOVED*** else if !fastpathDecodeTypeSwitch(iv, d) ***REMOVED***
			v := rv4i(iv)
			d.ensureDecodeable(v)
			d.decodeValue(v, nil)
		***REMOVED***
	***REMOVED***
***REMOVED***

// decodeValue MUST be called by the actual value we want to decode into,
// not its addr or a reference to it.
//
// This way, we know if it is itself a pointer, and can handle nil in
// the stream effectively.
func (d *Decoder) decodeValue(rv reflect.Value, fn *codecFn) ***REMOVED***
	// If stream is not containing a nil value, then we can deref to the base
	// non-pointer value, and decode into that.
	var rvp reflect.Value
	var rvpValid bool
	if rv.Kind() == reflect.Ptr ***REMOVED***
		if d.d.TryNil() ***REMOVED***
			if rvelem := rv.Elem(); rvelem.CanSet() ***REMOVED***
				rvelem.Set(reflect.Zero(rvelem.Type()))
			***REMOVED***
			return
		***REMOVED***
		rvpValid = true
		for rv.Kind() == reflect.Ptr ***REMOVED***
			if rvIsNil(rv) ***REMOVED***
				rvSetDirect(rv, reflect.New(rv.Type().Elem()))
			***REMOVED***
			rvp = rv
			rv = rv.Elem()
		***REMOVED***
	***REMOVED***

	if fn == nil ***REMOVED***
		fn = d.h.fn(rv.Type())
	***REMOVED***
	if fn.i.addrD ***REMOVED***
		if rvpValid ***REMOVED***
			fn.fd(d, &fn.i, rvp)
		***REMOVED*** else if rv.CanAddr() ***REMOVED***
			fn.fd(d, &fn.i, rv.Addr())
		***REMOVED*** else if !fn.i.addrF ***REMOVED***
			fn.fd(d, &fn.i, rv)
		***REMOVED*** else ***REMOVED***
			d.errorf("cannot decode into a non-pointer value")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fn.fd(d, &fn.i, rv)
	***REMOVED***
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

func isDecodeable(rv reflect.Value) (canDecode bool) ***REMOVED***
	switch rv.Kind() ***REMOVED***
	case reflect.Array:
		return rv.CanAddr()
	case reflect.Ptr:
		if !rvIsNil(rv) ***REMOVED***
			return true
		***REMOVED***
	case reflect.Slice, reflect.Chan, reflect.Map:
		if !rvIsNil(rv) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (d *Decoder) ensureDecodeable(rv reflect.Value) ***REMOVED***
	// decode can take any reflect.Value that is a inherently addressable i.e.
	//   - array
	//   - non-nil chan    (we will SEND to it)
	//   - non-nil slice   (we will set its elements)
	//   - non-nil map     (we will put into it)
	//   - non-nil pointer (we can "update" it)
	if isDecodeable(rv) ***REMOVED***
		return
	***REMOVED***
	if !rv.IsValid() ***REMOVED***
		d.errorstr(errstrCannotDecodeIntoNil)
		return
	***REMOVED***
	if !rv.CanInterface() ***REMOVED***
		d.errorf("cannot decode into a value without an interface: %v", rv)
		return
	***REMOVED***
	rvi := rv2i(rv)
	rvk := rv.Kind()
	d.errorf("cannot decode into value of kind: %v, type: %T, %#v", rvk, rvi, rvi)
***REMOVED***

func (d *Decoder) depthIncr() ***REMOVED***
	d.depth++
	if d.depth >= d.maxdepth ***REMOVED***
		panic(errMaxDepthExceeded)
	***REMOVED***
***REMOVED***

func (d *Decoder) depthDecr() ***REMOVED***
	d.depth--
***REMOVED***

// Possibly get an interned version of a string
//
// This should mostly be used for map keys, where the key type is string.
// This is because keys of a map/struct are typically reused across many objects.
func (d *Decoder) string(v []byte) (s string) ***REMOVED***
	if v == nil ***REMOVED***
		return
	***REMOVED***
	if d.is == nil ***REMOVED***
		return string(v) // don't return stringView, as we need a real string here.
	***REMOVED***
	s, ok := d.is[string(v)] // no allocation here, per go implementation
	if !ok ***REMOVED***
		s = string(v) // new allocation here
		d.is[s] = s
	***REMOVED***
	return
***REMOVED***

// nextValueBytes returns the next value in the stream as a set of bytes.
func (d *Decoder) nextValueBytes() (bs []byte) ***REMOVED***
	d.d.uncacheRead()
	d.r().track()
	d.swallow()
	bs = d.r().stopTrack()
	return
***REMOVED***

func (d *Decoder) rawBytes() []byte ***REMOVED***
	// ensure that this is not a view into the bytes
	// i.e. make new copy always.
	bs := d.nextValueBytes()
	bs2 := make([]byte, len(bs))
	copy(bs2, bs)
	return bs2
***REMOVED***

func (d *Decoder) wrapErr(v interface***REMOVED******REMOVED***, err *error) ***REMOVED***
	*err = decodeError***REMOVED***codecError: codecError***REMOVED***name: d.hh.Name(), err: v***REMOVED***, pos: d.NumBytesRead()***REMOVED***
***REMOVED***

// NumBytesRead returns the number of bytes read
func (d *Decoder) NumBytesRead() int ***REMOVED***
	return int(d.r().numread())
***REMOVED***

// decodeFloat32 will delegate to an appropriate DecodeFloat32 implementation (if exists),
// else if will call DecodeFloat64 and ensure the value doesn't overflow.
//
// Note that we return float64 to reduce unnecessary conversions
func (d *Decoder) decodeFloat32() float32 ***REMOVED***
	if d.js ***REMOVED***
		return d.jsondriver().DecodeFloat32() // custom implementation for 32-bit
	***REMOVED***
	return float32(chkOvf.Float32V(d.d.DecodeFloat64()))
***REMOVED***

// ---- container tracking
// Note: We update the .c after calling the callback.
// This way, the callback can know what the last status was.

// Note: if you call mapStart and it returns decContainerLenNil,
// then do NOT call mapEnd.

func (d *Decoder) mapStart() (v int) ***REMOVED***
	v = d.d.ReadMapStart()
	if v != decContainerLenNil ***REMOVED***
		d.depthIncr()
		d.c = containerMapStart
	***REMOVED***
	return
***REMOVED***

func (d *Decoder) mapElemKey() ***REMOVED***
	if d.js ***REMOVED***
		d.jsondriver().ReadMapElemKey()
	***REMOVED***
	d.c = containerMapKey
***REMOVED***

func (d *Decoder) mapElemValue() ***REMOVED***
	if d.js ***REMOVED***
		d.jsondriver().ReadMapElemValue()
	***REMOVED***
	d.c = containerMapValue
***REMOVED***

func (d *Decoder) mapEnd() ***REMOVED***
	d.d.ReadMapEnd()
	d.depthDecr()
	// d.c = containerMapEnd
	d.c = 0
***REMOVED***

func (d *Decoder) arrayStart() (v int) ***REMOVED***
	v = d.d.ReadArrayStart()
	if v != decContainerLenNil ***REMOVED***
		d.depthIncr()
		d.c = containerArrayStart
	***REMOVED***
	return
***REMOVED***

func (d *Decoder) arrayElem() ***REMOVED***
	if d.js ***REMOVED***
		d.jsondriver().ReadArrayElem()
	***REMOVED***
	d.c = containerArrayElem
***REMOVED***

func (d *Decoder) arrayEnd() ***REMOVED***
	d.d.ReadArrayEnd()
	d.depthDecr()
	// d.c = containerArrayEnd
	d.c = 0
***REMOVED***

func (d *Decoder) interfaceExtConvertAndDecode(v interface***REMOVED******REMOVED***, ext Ext) ***REMOVED***
	// var v interface***REMOVED******REMOVED*** = ext.ConvertExt(rv)
	// d.d.decode(&v)
	// ext.UpdateExt(rv, v)

	// assume v is a pointer:
	// - if struct|array, pass as is to ConvertExt
	// - else make it non-addressable and pass to ConvertExt
	// - make return value from ConvertExt addressable
	// - decode into it
	// - return the interface for passing into UpdateExt.
	// - interface should be a pointer if struct|array, else a value

	var s interface***REMOVED******REMOVED***
	rv := rv4i(v)
	rv2 := rv.Elem()
	rvk := rv2.Kind()
	if rvk == reflect.Struct || rvk == reflect.Array ***REMOVED***
		s = ext.ConvertExt(v)
	***REMOVED*** else ***REMOVED***
		s = ext.ConvertExt(rv2i(rv2))
	***REMOVED***
	rv = rv4i(s)
	if !rv.CanAddr() ***REMOVED***
		if rv.Kind() == reflect.Ptr ***REMOVED***
			rv2 = reflect.New(rv.Type().Elem())
		***REMOVED*** else ***REMOVED***
			rv2 = rvZeroAddrK(rv.Type(), rv.Kind())
		***REMOVED***
		rvSetDirect(rv2, rv)
		rv = rv2
	***REMOVED***
	d.decodeValue(rv, nil)
	ext.UpdateExt(v, rv2i(rv))
***REMOVED***

func (d *Decoder) sideDecode(v interface***REMOVED******REMOVED***, bs []byte) ***REMOVED***
	rv := baseRV(v)
	NewDecoderBytes(bs, d.hh).decodeValue(rv, d.h.fnNoExt(rv.Type()))
***REMOVED***

// --------------------------------------------------

// decSliceHelper assists when decoding into a slice, from a map or an array in the stream.
// A slice can be set from a map or array in stream. This supports the MapBySlice interface.
//
// Note: if IsNil, do not call ElemContainerState.
type decSliceHelper struct ***REMOVED***
	d     *Decoder
	ct    valueType
	Array bool
	IsNil bool
***REMOVED***

func (d *Decoder) decSliceHelperStart() (x decSliceHelper, clen int) ***REMOVED***
	x.ct = d.d.ContainerType()
	x.d = d
	switch x.ct ***REMOVED***
	case valueTypeNil:
		x.IsNil = true
	case valueTypeArray:
		x.Array = true
		clen = d.arrayStart()
	case valueTypeMap:
		clen = d.mapStart() * 2
	default:
		d.errorf("only encoded map or array can be decoded into a slice (%d)", x.ct)
	***REMOVED***
	return
***REMOVED***

func (x decSliceHelper) End() ***REMOVED***
	if x.IsNil ***REMOVED***
	***REMOVED*** else if x.Array ***REMOVED***
		x.d.arrayEnd()
	***REMOVED*** else ***REMOVED***
		x.d.mapEnd()
	***REMOVED***
***REMOVED***

func (x decSliceHelper) ElemContainerState(index int) ***REMOVED***
	// Note: if isnil, clen=0, so we never call into ElemContainerState

	if x.Array ***REMOVED***
		x.d.arrayElem()
	***REMOVED*** else ***REMOVED***
		if index%2 == 0 ***REMOVED***
			x.d.mapElemKey()
		***REMOVED*** else ***REMOVED***
			x.d.mapElemValue()
		***REMOVED***
	***REMOVED***
***REMOVED***

func decByteSlice(r *decRd, clen, maxInitLen int, bs []byte) (bsOut []byte) ***REMOVED***
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
		len2 := decInferLen(clen, maxInitLen, 1)
		bsOut = make([]byte, len2)
		r.readb(bsOut)
		for len2 < clen ***REMOVED***
			len3 := decInferLen(clen-len2, maxInitLen, 1)
			bs3 := bsOut
			bsOut = make([]byte, len2+len3)
			copy(bsOut, bs3)
			r.readb(bsOut[len2:])
			len2 += len3
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// detachZeroCopyBytes will copy the in bytes into dest,
// or create a new one if not large enough.
//
// It is used to ensure that the []byte returned is not
// part of the input stream or input stream buffers.
func detachZeroCopyBytes(isBytesReader bool, dest []byte, in []byte) (out []byte) ***REMOVED***
	if len(in) > 0 ***REMOVED***
		// if isBytesReader || len(in) <= scratchByteArrayLen ***REMOVED***
		// 	if cap(dest) >= len(in) ***REMOVED***
		// 		out = dest[:len(in)]
		// 	***REMOVED*** else ***REMOVED***
		// 		out = make([]byte, len(in))
		// 	***REMOVED***
		// 	copy(out, in)
		// 	return
		// ***REMOVED***
		if cap(dest) >= len(in) ***REMOVED***
			out = dest[:len(in)]
		***REMOVED*** else ***REMOVED***
			out = make([]byte, len(in))
		***REMOVED***
		copy(out, in)
		return
	***REMOVED***
	return in
***REMOVED***

// decInferLen will infer a sensible length, given the following:
//    - clen: length wanted.
//    - maxlen: max length to be returned.
//      if <= 0, it is unset, and we infer it based on the unit size
//    - unit: number of bytes for each element of the collection
func decInferLen(clen, maxlen, unit int) (rvlen int) ***REMOVED***
	const maxLenIfUnset = 8 // 64
	// handle when maxlen is not set i.e. <= 0

	// clen==0:           use 0
	// maxlen<=0, clen<0: use default
	// maxlen> 0, clen<0: use default
	// maxlen<=0, clen>0: infer maxlen, and cap on it
	// maxlen> 0, clen>0: cap at maxlen

	if clen == 0 ***REMOVED***
		return
	***REMOVED***
	if clen < 0 ***REMOVED***
		if clen == decContainerLenNil ***REMOVED***
			return 0
		***REMOVED***
		return maxLenIfUnset
	***REMOVED***
	if unit == 0 ***REMOVED***
		return clen
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
		// if maxlen > maxLenIfUnset ***REMOVED***
		// 	maxlen = maxLenIfUnset
		// ***REMOVED***
	***REMOVED***
	if clen > maxlen ***REMOVED***
		rvlen = maxlen
	***REMOVED*** else ***REMOVED***
		rvlen = clen
	***REMOVED***
	return
***REMOVED***

func decReadFull(r io.Reader, bs []byte) (n uint, err error) ***REMOVED***
	var nn int
	for n < uint(len(bs)) && err == nil ***REMOVED***
		nn, err = r.Read(bs[n:])
		if nn > 0 ***REMOVED***
			if err == io.EOF ***REMOVED***
				// leave EOF for next time
				err = nil
			***REMOVED***
			n += uint(nn)
		***REMOVED***
	***REMOVED***
	// do not do this - it serves no purpose
	// if n != len(bs) && err == io.EOF ***REMOVED*** err = io.ErrUnexpectedEOF ***REMOVED***
	return
***REMOVED***

func decNakedReadRawBytes(dr decDriver, d *Decoder, n *decNaked, rawToString bool) ***REMOVED***
	if rawToString ***REMOVED***
		n.v = valueTypeString
		n.s = string(dr.DecodeBytes(d.b[:], true))
	***REMOVED*** else ***REMOVED***
		n.v = valueTypeBytes
		n.l = dr.DecodeBytes(nil, false)
	***REMOVED***
***REMOVED***
