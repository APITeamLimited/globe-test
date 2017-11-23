// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// Contains code shared by both encode and decode.

// Some shared ideas around encoding/decoding
// ------------------------------------------
//
// If an interface***REMOVED******REMOVED*** is passed, we first do a type assertion to see if it is
// a primitive type or a map/slice of primitive types, and use a fastpath to handle it.
//
// If we start with a reflect.Value, we are already in reflect.Value land and
// will try to grab the function for the underlying Type and directly call that function.
// This is more performant than calling reflect.Value.Interface().
//
// This still helps us bypass many layers of reflection, and give best performance.
//
// Containers
// ------------
// Containers in the stream are either associative arrays (key-value pairs) or
// regular arrays (indexed by incrementing integers).
//
// Some streams support indefinite-length containers, and use a breaking
// byte-sequence to denote that the container has come to an end.
//
// Some streams also are text-based, and use explicit separators to denote the
// end/beginning of different values.
//
// During encode, we use a high-level condition to determine how to iterate through
// the container. That decision is based on whether the container is text-based (with
// separators) or binary (without separators). If binary, we do not even call the
// encoding of separators.
//
// During decode, we use a different high-level condition to determine how to iterate
// through the containers. That decision is based on whether the stream contained
// a length prefix, or if it used explicit breaks. If length-prefixed, we assume that
// it has to be binary, and we do not even try to read separators.
//
// Philosophy
// ------------
// On decode, this codec will update containers appropriately:
//    - If struct, update fields from stream into fields of struct.
//      If field in stream not found in struct, handle appropriately (based on option).
//      If a struct field has no corresponding value in the stream, leave it AS IS.
//      If nil in stream, set value to nil/zero value.
//    - If map, update map from stream.
//      If the stream value is NIL, set the map to nil.
//    - if slice, try to update up to length of array in stream.
//      if container len is less than stream array length,
//      and container cannot be expanded, handled (based on option).
//      This means you can decode 4-element stream array into 1-element array.
//
// ------------------------------------
// On encode, user can specify omitEmpty. This means that the value will be omitted
// if the zero value. The problem may occur during decode, where omitted values do not affect
// the value being decoded into. This means that if decoding into a struct with an
// int field with current value=5, and the field is omitted in the stream, then after
// decoding, the value will still be 5 (not 0).
// omitEmpty only works if you guarantee that you always decode into zero-values.
//
// ------------------------------------
// We could have truncated a map to remove keys not available in the stream,
// or set values in the struct which are not in the stream to their zero values.
// We decided against it because there is no efficient way to do it.
// We may introduce it as an option later.
// However, that will require enabling it for both runtime and code generation modes.
//
// To support truncate, we need to do 2 passes over the container:
//   map
//   - first collect all keys (e.g. in k1)
//   - for each key in stream, mark k1 that the key should not be removed
//   - after updating map, do second pass and call delete for all keys in k1 which are not marked
//   struct:
//   - for each field, track the *typeInfo s1
//   - iterate through all s1, and for each one not marked, set value to zero
//   - this involves checking the possible anonymous fields which are nil ptrs.
//     too much work.
//
// ------------------------------------------
// Error Handling is done within the library using panic.
//
// This way, the code doesn't have to keep checking if an error has happened,
// and we don't have to keep sending the error value along with each call
// or storing it in the En|Decoder and checking it constantly along the way.
//
// The disadvantage is that small functions which use panics cannot be inlined.
// The code accounts for that by only using panics behind an interface;
// since interface calls cannot be inlined, this is irrelevant.
//
// We considered storing the error is En|Decoder.
//   - once it has its err field set, it cannot be used again.
//   - panicing will be optional, controlled by const flag.
//   - code should always check error first and return early.
// We eventually decided against it as it makes the code clumsier to always
// check for these error conditions.

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	scratchByteArrayLen = 32
	initCollectionCap   = 32 // 32 is defensive. 16 is preferred.

	// Support encoding.(Binary|Text)(Unm|M)arshaler.
	// This constant flag will enable or disable it.
	supportMarshalInterfaces = true

	// Each Encoder or Decoder uses a cache of functions based on conditionals,
	// so that the conditionals are not run every time.
	//
	// Either a map or a slice is used to keep track of the functions.
	// The map is more natural, but has a higher cost than a slice/array.
	// This flag (useMapForCodecCache) controls which is used.
	//
	// From benchmarks, slices with linear search perform better with < 32 entries.
	// We have typically seen a high threshold of about 24 entries.
	useMapForCodecCache = false

	// for debugging, set this to false, to catch panic traces.
	// Note that this will always cause rpc tests to fail, since they need io.EOF sent via panic.
	recoverPanicToErr = true

	// if resetSliceElemToZeroValue, then on decoding a slice, reset the element to a zero value first.
	// Only concern is that, if the slice already contained some garbage, we will decode into that garbage.
	// The chances of this are slim, so leave this "optimization".
	// TODO: should this be true, to ensure that we always decode into a "zero" "empty" value?
	resetSliceElemToZeroValue bool = false
)

var (
	oneByteArr    = [1]byte***REMOVED***0***REMOVED***
	zeroByteSlice = oneByteArr[:0:0]
)

type charEncoding uint8

const (
	c_RAW charEncoding = iota
	c_UTF8
	c_UTF16LE
	c_UTF16BE
	c_UTF32LE
	c_UTF32BE
)

// valueType is the stream type
type valueType uint8

const (
	valueTypeUnset valueType = iota
	valueTypeNil
	valueTypeInt
	valueTypeUint
	valueTypeFloat
	valueTypeBool
	valueTypeString
	valueTypeSymbol
	valueTypeBytes
	valueTypeMap
	valueTypeArray
	valueTypeTimestamp
	valueTypeExt

	// valueTypeInvalid = 0xff
)

type seqType uint8

const (
	_ seqType = iota
	seqTypeArray
	seqTypeSlice
	seqTypeChan
)

// note that containerMapStart and containerArraySend are not sent.
// This is because the ReadXXXStart and EncodeXXXStart already does these.
type containerState uint8

const (
	_ containerState = iota

	containerMapStart // slot left open, since Driver method already covers it
	containerMapKey
	containerMapValue
	containerMapEnd
	containerArrayStart // slot left open, since Driver methods already cover it
	containerArrayElem
	containerArrayEnd
)

// sfiIdx used for tracking where a (field/enc)Name is seen in a []*structFieldInfo
type sfiIdx struct ***REMOVED***
	name  string
	index int
***REMOVED***

// do not recurse if a containing type refers to an embedded type
// which refers back to its containing type (via a pointer).
// The second time this back-reference happens, break out,
// so as not to cause an infinite loop.
const rgetMaxRecursion = 2

// Anecdotally, we believe most types have <= 12 fields.
// Java's PMD rules set TooManyFields threshold to 15.
const rgetPoolTArrayLen = 12

type rgetT struct ***REMOVED***
	fNames   []string
	encNames []string
	etypes   []uintptr
	sfis     []*structFieldInfo
***REMOVED***

type rgetPoolT struct ***REMOVED***
	fNames   [rgetPoolTArrayLen]string
	encNames [rgetPoolTArrayLen]string
	etypes   [rgetPoolTArrayLen]uintptr
	sfis     [rgetPoolTArrayLen]*structFieldInfo
	sfiidx   [rgetPoolTArrayLen]sfiIdx
***REMOVED***

var rgetPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(rgetPoolT) ***REMOVED***,
***REMOVED***

type containerStateRecv interface ***REMOVED***
	sendContainerState(containerState)
***REMOVED***

// mirror json.Marshaler and json.Unmarshaler here,
// so we don't import the encoding/json package
type jsonMarshaler interface ***REMOVED***
	MarshalJSON() ([]byte, error)
***REMOVED***
type jsonUnmarshaler interface ***REMOVED***
	UnmarshalJSON([]byte) error
***REMOVED***

var (
	bigen               = binary.BigEndian
	structInfoFieldName = "_struct"

	mapStrIntfTyp  = reflect.TypeOf(map[string]interface***REMOVED******REMOVED***(nil))
	mapIntfIntfTyp = reflect.TypeOf(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***(nil))
	intfSliceTyp   = reflect.TypeOf([]interface***REMOVED******REMOVED***(nil))
	intfTyp        = intfSliceTyp.Elem()

	stringTyp     = reflect.TypeOf("")
	timeTyp       = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	rawExtTyp     = reflect.TypeOf(RawExt***REMOVED******REMOVED***)
	rawTyp        = reflect.TypeOf(Raw***REMOVED******REMOVED***)
	uint8SliceTyp = reflect.TypeOf([]uint8(nil))

	mapBySliceTyp = reflect.TypeOf((*MapBySlice)(nil)).Elem()

	binaryMarshalerTyp   = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	binaryUnmarshalerTyp = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()

	textMarshalerTyp   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerTyp = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

	jsonMarshalerTyp   = reflect.TypeOf((*jsonMarshaler)(nil)).Elem()
	jsonUnmarshalerTyp = reflect.TypeOf((*jsonUnmarshaler)(nil)).Elem()

	selferTyp = reflect.TypeOf((*Selfer)(nil)).Elem()

	uint8SliceTypId = reflect.ValueOf(uint8SliceTyp).Pointer()
	rawExtTypId     = reflect.ValueOf(rawExtTyp).Pointer()
	rawTypId        = reflect.ValueOf(rawTyp).Pointer()
	intfTypId       = reflect.ValueOf(intfTyp).Pointer()
	timeTypId       = reflect.ValueOf(timeTyp).Pointer()
	stringTypId     = reflect.ValueOf(stringTyp).Pointer()

	mapStrIntfTypId  = reflect.ValueOf(mapStrIntfTyp).Pointer()
	mapIntfIntfTypId = reflect.ValueOf(mapIntfIntfTyp).Pointer()
	intfSliceTypId   = reflect.ValueOf(intfSliceTyp).Pointer()
	// mapBySliceTypId  = reflect.ValueOf(mapBySliceTyp).Pointer()

	intBitsize  uint8 = uint8(reflect.TypeOf(int(0)).Bits())
	uintBitsize uint8 = uint8(reflect.TypeOf(uint(0)).Bits())

	bsAll0x00 = []byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0***REMOVED***
	bsAll0xff = []byte***REMOVED***0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff***REMOVED***

	chkOvf checkOverflow

	noFieldNameToStructFieldInfoErr = errors.New("no field name passed to parseStructFieldInfo")
)

var defTypeInfos = NewTypeInfos([]string***REMOVED***"codec", "json"***REMOVED***)

// Selfer defines methods by which a value can encode or decode itself.
//
// Any type which implements Selfer will be able to encode or decode itself.
// Consequently, during (en|de)code, this takes precedence over
// (text|binary)(M|Unm)arshal or extension support.
type Selfer interface ***REMOVED***
	CodecEncodeSelf(*Encoder)
	CodecDecodeSelf(*Decoder)
***REMOVED***

// MapBySlice represents a slice which should be encoded as a map in the stream.
// The slice contains a sequence of key-value pairs.
// This affords storing a map in a specific sequence in the stream.
//
// The support of MapBySlice affords the following:
//   - A slice type which implements MapBySlice will be encoded as a map
//   - A slice can be decoded from a map in the stream
type MapBySlice interface ***REMOVED***
	MapBySlice()
***REMOVED***

// WARNING: DO NOT USE DIRECTLY. EXPORTED FOR GODOC BENEFIT. WILL BE REMOVED.
//
// BasicHandle encapsulates the common options and extension functions.
type BasicHandle struct ***REMOVED***
	// TypeInfos is used to get the type info for any type.
	//
	// If not configured, the default TypeInfos is used, which uses struct tag keys: codec, json
	TypeInfos *TypeInfos

	extHandle
	EncodeOptions
	DecodeOptions
***REMOVED***

func (x *BasicHandle) getBasicHandle() *BasicHandle ***REMOVED***
	return x
***REMOVED***

func (x *BasicHandle) getTypeInfo(rtid uintptr, rt reflect.Type) (pti *typeInfo) ***REMOVED***
	if x.TypeInfos != nil ***REMOVED***
		return x.TypeInfos.get(rtid, rt)
	***REMOVED***
	return defTypeInfos.get(rtid, rt)
***REMOVED***

// Handle is the interface for a specific encoding format.
//
// Typically, a Handle is pre-configured before first time use,
// and not modified while in use. Such a pre-configured Handle
// is safe for concurrent access.
type Handle interface ***REMOVED***
	getBasicHandle() *BasicHandle
	newEncDriver(w *Encoder) encDriver
	newDecDriver(r *Decoder) decDriver
	isBinary() bool
***REMOVED***

// Raw represents raw formatted bytes.
// We "blindly" store it during encode and store the raw bytes during decode.
// Note: it is dangerous during encode, so we may gate the behaviour behind an Encode flag which must be explicitly set.
type Raw []byte

// RawExt represents raw unprocessed extension data.
// Some codecs will decode extension data as a *RawExt if there is no registered extension for the tag.
//
// Only one of Data or Value is nil. If Data is nil, then the content of the RawExt is in the Value.
type RawExt struct ***REMOVED***
	Tag uint64
	// Data is the []byte which represents the raw ext. If Data is nil, ext is exposed in Value.
	// Data is used by codecs (e.g. binc, msgpack, simple) which do custom serialization of the types
	Data []byte
	// Value represents the extension, if Data is nil.
	// Value is used by codecs (e.g. cbor, json) which use the format to do custom serialization of the types.
	Value interface***REMOVED******REMOVED***
***REMOVED***

// BytesExt handles custom (de)serialization of types to/from []byte.
// It is used by codecs (e.g. binc, msgpack, simple) which do custom serialization of the types.
type BytesExt interface ***REMOVED***
	// WriteExt converts a value to a []byte.
	//
	// Note: v *may* be a pointer to the extension type, if the extension type was a struct or array.
	WriteExt(v interface***REMOVED******REMOVED***) []byte

	// ReadExt updates a value from a []byte.
	ReadExt(dst interface***REMOVED******REMOVED***, src []byte)
***REMOVED***

// InterfaceExt handles custom (de)serialization of types to/from another interface***REMOVED******REMOVED*** value.
// The Encoder or Decoder will then handle the further (de)serialization of that known type.
//
// It is used by codecs (e.g. cbor, json) which use the format to do custom serialization of the types.
type InterfaceExt interface ***REMOVED***
	// ConvertExt converts a value into a simpler interface for easy encoding e.g. convert time.Time to int64.
	//
	// Note: v *may* be a pointer to the extension type, if the extension type was a struct or array.
	ConvertExt(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED***

	// UpdateExt updates a value from a simpler interface for easy decoding e.g. convert int64 to time.Time.
	UpdateExt(dst interface***REMOVED******REMOVED***, src interface***REMOVED******REMOVED***)
***REMOVED***

// Ext handles custom (de)serialization of custom types / extensions.
type Ext interface ***REMOVED***
	BytesExt
	InterfaceExt
***REMOVED***

// addExtWrapper is a wrapper implementation to support former AddExt exported method.
type addExtWrapper struct ***REMOVED***
	encFn func(reflect.Value) ([]byte, error)
	decFn func(reflect.Value, []byte) error
***REMOVED***

func (x addExtWrapper) WriteExt(v interface***REMOVED******REMOVED***) []byte ***REMOVED***
	bs, err := x.encFn(reflect.ValueOf(v))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return bs
***REMOVED***

func (x addExtWrapper) ReadExt(v interface***REMOVED******REMOVED***, bs []byte) ***REMOVED***
	if err := x.decFn(reflect.ValueOf(v), bs); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (x addExtWrapper) ConvertExt(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return x.WriteExt(v)
***REMOVED***

func (x addExtWrapper) UpdateExt(dest interface***REMOVED******REMOVED***, v interface***REMOVED******REMOVED***) ***REMOVED***
	x.ReadExt(dest, v.([]byte))
***REMOVED***

type setExtWrapper struct ***REMOVED***
	b BytesExt
	i InterfaceExt
***REMOVED***

func (x *setExtWrapper) WriteExt(v interface***REMOVED******REMOVED***) []byte ***REMOVED***
	if x.b == nil ***REMOVED***
		panic("BytesExt.WriteExt is not supported")
	***REMOVED***
	return x.b.WriteExt(v)
***REMOVED***

func (x *setExtWrapper) ReadExt(v interface***REMOVED******REMOVED***, bs []byte) ***REMOVED***
	if x.b == nil ***REMOVED***
		panic("BytesExt.WriteExt is not supported")

	***REMOVED***
	x.b.ReadExt(v, bs)
***REMOVED***

func (x *setExtWrapper) ConvertExt(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	if x.i == nil ***REMOVED***
		panic("InterfaceExt.ConvertExt is not supported")

	***REMOVED***
	return x.i.ConvertExt(v)
***REMOVED***

func (x *setExtWrapper) UpdateExt(dest interface***REMOVED******REMOVED***, v interface***REMOVED******REMOVED***) ***REMOVED***
	if x.i == nil ***REMOVED***
		panic("InterfaceExxt.UpdateExt is not supported")

	***REMOVED***
	x.i.UpdateExt(dest, v)
***REMOVED***

// type errorString string
// func (x errorString) Error() string ***REMOVED*** return string(x) ***REMOVED***

type binaryEncodingType struct***REMOVED******REMOVED***

func (_ binaryEncodingType) isBinary() bool ***REMOVED*** return true ***REMOVED***

type textEncodingType struct***REMOVED******REMOVED***

func (_ textEncodingType) isBinary() bool ***REMOVED*** return false ***REMOVED***

// noBuiltInTypes is embedded into many types which do not support builtins
// e.g. msgpack, simple, cbor.
type noBuiltInTypes struct***REMOVED******REMOVED***

func (_ noBuiltInTypes) IsBuiltinType(rt uintptr) bool           ***REMOVED*** return false ***REMOVED***
func (_ noBuiltInTypes) EncodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***
func (_ noBuiltInTypes) DecodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***

type noStreamingCodec struct***REMOVED******REMOVED***

func (_ noStreamingCodec) CheckBreak() bool ***REMOVED*** return false ***REMOVED***

// bigenHelper.
// Users must already slice the x completely, because we will not reslice.
type bigenHelper struct ***REMOVED***
	x []byte // must be correctly sliced to appropriate len. slicing is a cost.
	w encWriter
***REMOVED***

func (z bigenHelper) writeUint16(v uint16) ***REMOVED***
	bigen.PutUint16(z.x, v)
	z.w.writeb(z.x)
***REMOVED***

func (z bigenHelper) writeUint32(v uint32) ***REMOVED***
	bigen.PutUint32(z.x, v)
	z.w.writeb(z.x)
***REMOVED***

func (z bigenHelper) writeUint64(v uint64) ***REMOVED***
	bigen.PutUint64(z.x, v)
	z.w.writeb(z.x)
***REMOVED***

type extTypeTagFn struct ***REMOVED***
	rtid uintptr
	rt   reflect.Type
	tag  uint64
	ext  Ext
***REMOVED***

type extHandle []extTypeTagFn

// DEPRECATED: Use SetBytesExt or SetInterfaceExt on the Handle instead.
//
// AddExt registes an encode and decode function for a reflect.Type.
// AddExt internally calls SetExt.
// To deregister an Ext, call AddExt with nil encfn and/or nil decfn.
func (o *extHandle) AddExt(
	rt reflect.Type, tag byte,
	encfn func(reflect.Value) ([]byte, error), decfn func(reflect.Value, []byte) error,
) (err error) ***REMOVED***
	if encfn == nil || decfn == nil ***REMOVED***
		return o.SetExt(rt, uint64(tag), nil)
	***REMOVED***
	return o.SetExt(rt, uint64(tag), addExtWrapper***REMOVED***encfn, decfn***REMOVED***)
***REMOVED***

// DEPRECATED: Use SetBytesExt or SetInterfaceExt on the Handle instead.
//
// Note that the type must be a named type, and specifically not
// a pointer or Interface. An error is returned if that is not honored.
//
// To Deregister an ext, call SetExt with nil Ext
func (o *extHandle) SetExt(rt reflect.Type, tag uint64, ext Ext) (err error) ***REMOVED***
	// o is a pointer, because we may need to initialize it
	if rt.PkgPath() == "" || rt.Kind() == reflect.Interface ***REMOVED***
		err = fmt.Errorf("codec.Handle.AddExt: Takes named type, not a pointer or interface: %T",
			reflect.Zero(rt).Interface())
		return
	***REMOVED***

	rtid := reflect.ValueOf(rt).Pointer()
	for _, v := range *o ***REMOVED***
		if v.rtid == rtid ***REMOVED***
			v.tag, v.ext = tag, ext
			return
		***REMOVED***
	***REMOVED***

	if *o == nil ***REMOVED***
		*o = make([]extTypeTagFn, 0, 4)
	***REMOVED***
	*o = append(*o, extTypeTagFn***REMOVED***rtid, rt, tag, ext***REMOVED***)
	return
***REMOVED***

func (o extHandle) getExt(rtid uintptr) *extTypeTagFn ***REMOVED***
	var v *extTypeTagFn
	for i := range o ***REMOVED***
		v = &o[i]
		if v.rtid == rtid ***REMOVED***
			return v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (o extHandle) getExtForTag(tag uint64) *extTypeTagFn ***REMOVED***
	var v *extTypeTagFn
	for i := range o ***REMOVED***
		v = &o[i]
		if v.tag == tag ***REMOVED***
			return v
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type structFieldInfo struct ***REMOVED***
	encName   string // encode name
	fieldName string // field name

	// only one of 'i' or 'is' can be set. If 'i' is -1, then 'is' has been set.

	is        []int // (recursive/embedded) field index in struct
	i         int16 // field index in struct
	omitEmpty bool
	toArray   bool // if field is _struct, is the toArray set?
***REMOVED***

// func (si *structFieldInfo) isZero() bool ***REMOVED***
// 	return si.encName == "" && len(si.is) == 0 && si.i == 0 && !si.omitEmpty && !si.toArray
// ***REMOVED***

// rv returns the field of the struct.
// If anonymous, it returns an Invalid
func (si *structFieldInfo) field(v reflect.Value, update bool) (rv2 reflect.Value) ***REMOVED***
	if si.i != -1 ***REMOVED***
		v = v.Field(int(si.i))
		return v
	***REMOVED***
	// replicate FieldByIndex
	for _, x := range si.is ***REMOVED***
		for v.Kind() == reflect.Ptr ***REMOVED***
			if v.IsNil() ***REMOVED***
				if !update ***REMOVED***
					return
				***REMOVED***
				v.Set(reflect.New(v.Type().Elem()))
			***REMOVED***
			v = v.Elem()
		***REMOVED***
		v = v.Field(x)
	***REMOVED***
	return v
***REMOVED***

func (si *structFieldInfo) setToZeroValue(v reflect.Value) ***REMOVED***
	if si.i != -1 ***REMOVED***
		v = v.Field(int(si.i))
		v.Set(reflect.Zero(v.Type()))
		// v.Set(reflect.New(v.Type()).Elem())
		// v.Set(reflect.New(v.Type()))
	***REMOVED*** else ***REMOVED***
		// replicate FieldByIndex
		for _, x := range si.is ***REMOVED***
			for v.Kind() == reflect.Ptr ***REMOVED***
				if v.IsNil() ***REMOVED***
					return
				***REMOVED***
				v = v.Elem()
			***REMOVED***
			v = v.Field(x)
		***REMOVED***
		v.Set(reflect.Zero(v.Type()))
	***REMOVED***
***REMOVED***

func parseStructFieldInfo(fname string, stag string) *structFieldInfo ***REMOVED***
	// if fname == "" ***REMOVED***
	// 	panic(noFieldNameToStructFieldInfoErr)
	// ***REMOVED***
	si := structFieldInfo***REMOVED***
		encName: fname,
	***REMOVED***

	if stag != "" ***REMOVED***
		for i, s := range strings.Split(stag, ",") ***REMOVED***
			if i == 0 ***REMOVED***
				if s != "" ***REMOVED***
					si.encName = s
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if s == "omitempty" ***REMOVED***
					si.omitEmpty = true
				***REMOVED*** else if s == "toarray" ***REMOVED***
					si.toArray = true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// si.encNameBs = []byte(si.encName)
	return &si
***REMOVED***

type sfiSortedByEncName []*structFieldInfo

func (p sfiSortedByEncName) Len() int ***REMOVED***
	return len(p)
***REMOVED***

func (p sfiSortedByEncName) Less(i, j int) bool ***REMOVED***
	return p[i].encName < p[j].encName
***REMOVED***

func (p sfiSortedByEncName) Swap(i, j int) ***REMOVED***
	p[i], p[j] = p[j], p[i]
***REMOVED***

// typeInfo keeps information about each type referenced in the encode/decode sequence.
//
// During an encode/decode sequence, we work as below:
//   - If base is a built in type, en/decode base value
//   - If base is registered as an extension, en/decode base value
//   - If type is binary(M/Unm)arshaler, call Binary(M/Unm)arshal method
//   - If type is text(M/Unm)arshaler, call Text(M/Unm)arshal method
//   - Else decode appropriately based on the reflect.Kind
type typeInfo struct ***REMOVED***
	sfi  []*structFieldInfo // sorted. Used when enc/dec struct to map.
	sfip []*structFieldInfo // unsorted. Used when enc/dec struct to array.

	rt   reflect.Type
	rtid uintptr

	numMeth uint16 // number of methods

	// baseId gives pointer to the base reflect.Type, after deferencing
	// the pointers. E.g. base type of ***time.Time is time.Time.
	base      reflect.Type
	baseId    uintptr
	baseIndir int8 // number of indirections to get to base

	mbs bool // base type (T or *T) is a MapBySlice

	bm        bool // base type (T or *T) is a binaryMarshaler
	bunm      bool // base type (T or *T) is a binaryUnmarshaler
	bmIndir   int8 // number of indirections to get to binaryMarshaler type
	bunmIndir int8 // number of indirections to get to binaryUnmarshaler type

	tm        bool // base type (T or *T) is a textMarshaler
	tunm      bool // base type (T or *T) is a textUnmarshaler
	tmIndir   int8 // number of indirections to get to textMarshaler type
	tunmIndir int8 // number of indirections to get to textUnmarshaler type

	jm        bool // base type (T or *T) is a jsonMarshaler
	junm      bool // base type (T or *T) is a jsonUnmarshaler
	jmIndir   int8 // number of indirections to get to jsonMarshaler type
	junmIndir int8 // number of indirections to get to jsonUnmarshaler type

	cs      bool // base type (T or *T) is a Selfer
	csIndir int8 // number of indirections to get to Selfer type

	toArray bool // whether this (struct) type should be encoded as an array
***REMOVED***

func (ti *typeInfo) indexForEncName(name string) int ***REMOVED***
	// NOTE: name may be a stringView, so don't pass it to another function.
	//tisfi := ti.sfi
	const binarySearchThreshold = 16
	if sfilen := len(ti.sfi); sfilen < binarySearchThreshold ***REMOVED***
		// linear search. faster than binary search in my testing up to 16-field structs.
		for i, si := range ti.sfi ***REMOVED***
			if si.encName == name ***REMOVED***
				return i
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// binary search. adapted from sort/search.go.
		h, i, j := 0, 0, sfilen
		for i < j ***REMOVED***
			h = i + (j-i)/2
			if ti.sfi[h].encName < name ***REMOVED***
				i = h + 1
			***REMOVED*** else ***REMOVED***
				j = h
			***REMOVED***
		***REMOVED***
		if i < sfilen && ti.sfi[i].encName == name ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// TypeInfos caches typeInfo for each type on first inspection.
//
// It is configured with a set of tag keys, which are used to get
// configuration for the type.
type TypeInfos struct ***REMOVED***
	infos map[uintptr]*typeInfo
	mu    sync.RWMutex
	tags  []string
***REMOVED***

// NewTypeInfos creates a TypeInfos given a set of struct tags keys.
//
// This allows users customize the struct tag keys which contain configuration
// of their types.
func NewTypeInfos(tags []string) *TypeInfos ***REMOVED***
	return &TypeInfos***REMOVED***tags: tags, infos: make(map[uintptr]*typeInfo, 64)***REMOVED***
***REMOVED***

func (x *TypeInfos) structTag(t reflect.StructTag) (s string) ***REMOVED***
	// check for tags: codec, json, in that order.
	// this allows seamless support for many configured structs.
	for _, x := range x.tags ***REMOVED***
		s = t.Get(x)
		if s != "" ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (x *TypeInfos) get(rtid uintptr, rt reflect.Type) (pti *typeInfo) ***REMOVED***
	var ok bool
	x.mu.RLock()
	pti, ok = x.infos[rtid]
	x.mu.RUnlock()
	if ok ***REMOVED***
		return
	***REMOVED***

	// do not hold lock while computing this.
	// it may lead to duplication, but that's ok.
	ti := typeInfo***REMOVED***rt: rt, rtid: rtid***REMOVED***
	ti.numMeth = uint16(rt.NumMethod())

	var indir int8
	if ok, indir = implementsIntf(rt, binaryMarshalerTyp); ok ***REMOVED***
		ti.bm, ti.bmIndir = true, indir
	***REMOVED***
	if ok, indir = implementsIntf(rt, binaryUnmarshalerTyp); ok ***REMOVED***
		ti.bunm, ti.bunmIndir = true, indir
	***REMOVED***
	if ok, indir = implementsIntf(rt, textMarshalerTyp); ok ***REMOVED***
		ti.tm, ti.tmIndir = true, indir
	***REMOVED***
	if ok, indir = implementsIntf(rt, textUnmarshalerTyp); ok ***REMOVED***
		ti.tunm, ti.tunmIndir = true, indir
	***REMOVED***
	if ok, indir = implementsIntf(rt, jsonMarshalerTyp); ok ***REMOVED***
		ti.jm, ti.jmIndir = true, indir
	***REMOVED***
	if ok, indir = implementsIntf(rt, jsonUnmarshalerTyp); ok ***REMOVED***
		ti.junm, ti.junmIndir = true, indir
	***REMOVED***
	if ok, indir = implementsIntf(rt, selferTyp); ok ***REMOVED***
		ti.cs, ti.csIndir = true, indir
	***REMOVED***
	if ok, _ = implementsIntf(rt, mapBySliceTyp); ok ***REMOVED***
		ti.mbs = true
	***REMOVED***

	pt := rt
	var ptIndir int8
	// for ; pt.Kind() == reflect.Ptr; pt, ptIndir = pt.Elem(), ptIndir+1 ***REMOVED*** ***REMOVED***
	for pt.Kind() == reflect.Ptr ***REMOVED***
		pt = pt.Elem()
		ptIndir++
	***REMOVED***
	if ptIndir == 0 ***REMOVED***
		ti.base = rt
		ti.baseId = rtid
	***REMOVED*** else ***REMOVED***
		ti.base = pt
		ti.baseId = reflect.ValueOf(pt).Pointer()
		ti.baseIndir = ptIndir
	***REMOVED***

	if rt.Kind() == reflect.Struct ***REMOVED***
		var omitEmpty bool
		if f, ok := rt.FieldByName(structInfoFieldName); ok ***REMOVED***
			siInfo := parseStructFieldInfo(structInfoFieldName, x.structTag(f.Tag))
			ti.toArray = siInfo.toArray
			omitEmpty = siInfo.omitEmpty
		***REMOVED***
		pi := rgetPool.Get()
		pv := pi.(*rgetPoolT)
		pv.etypes[0] = ti.baseId
		vv := rgetT***REMOVED***pv.fNames[:0], pv.encNames[:0], pv.etypes[:1], pv.sfis[:0]***REMOVED***
		x.rget(rt, rtid, omitEmpty, nil, &vv)
		ti.sfip, ti.sfi = rgetResolveSFI(vv.sfis, pv.sfiidx[:0])
		rgetPool.Put(pi)
	***REMOVED***
	// sfi = sfip

	x.mu.Lock()
	if pti, ok = x.infos[rtid]; !ok ***REMOVED***
		pti = &ti
		x.infos[rtid] = pti
	***REMOVED***
	x.mu.Unlock()
	return
***REMOVED***

func (x *TypeInfos) rget(rt reflect.Type, rtid uintptr, omitEmpty bool,
	indexstack []int, pv *rgetT,
) ***REMOVED***
	// Read up fields and store how to access the value.
	//
	// It uses go's rules for message selectors,
	// which say that the field with the shallowest depth is selected.
	//
	// Note: we consciously use slices, not a map, to simulate a set.
	//       Typically, types have < 16 fields,
	//       and iteration using equals is faster than maps there

LOOP:
	for j, jlen := 0, rt.NumField(); j < jlen; j++ ***REMOVED***
		f := rt.Field(j)
		fkind := f.Type.Kind()
		// skip if a func type, or is unexported, or structTag value == "-"
		switch fkind ***REMOVED***
		case reflect.Func, reflect.Complex64, reflect.Complex128, reflect.UnsafePointer:
			continue LOOP
		***REMOVED***

		// if r1, _ := utf8.DecodeRuneInString(f.Name);
		// r1 == utf8.RuneError || !unicode.IsUpper(r1) ***REMOVED***
		if f.PkgPath != "" && !f.Anonymous ***REMOVED*** // unexported, not embedded
			continue
		***REMOVED***
		stag := x.structTag(f.Tag)
		if stag == "-" ***REMOVED***
			continue
		***REMOVED***
		var si *structFieldInfo
		// if anonymous and no struct tag (or it's blank),
		// and a struct (or pointer to struct), inline it.
		if f.Anonymous && fkind != reflect.Interface ***REMOVED***
			doInline := stag == ""
			if !doInline ***REMOVED***
				si = parseStructFieldInfo("", stag)
				doInline = si.encName == ""
				// doInline = si.isZero()
			***REMOVED***
			if doInline ***REMOVED***
				ft := f.Type
				for ft.Kind() == reflect.Ptr ***REMOVED***
					ft = ft.Elem()
				***REMOVED***
				if ft.Kind() == reflect.Struct ***REMOVED***
					// if etypes contains this, don't call rget again (as fields are already seen here)
					ftid := reflect.ValueOf(ft).Pointer()
					// We cannot recurse forever, but we need to track other field depths.
					// So - we break if we see a type twice (not the first time).
					// This should be sufficient to handle an embedded type that refers to its
					// owning type, which then refers to its embedded type.
					processIt := true
					numk := 0
					for _, k := range pv.etypes ***REMOVED***
						if k == ftid ***REMOVED***
							numk++
							if numk == rgetMaxRecursion ***REMOVED***
								processIt = false
								break
							***REMOVED***
						***REMOVED***
					***REMOVED***
					if processIt ***REMOVED***
						pv.etypes = append(pv.etypes, ftid)
						indexstack2 := make([]int, len(indexstack)+1)
						copy(indexstack2, indexstack)
						indexstack2[len(indexstack)] = j
						// indexstack2 := append(append(make([]int, 0, len(indexstack)+4), indexstack...), j)
						x.rget(ft, ftid, omitEmpty, indexstack2, pv)
					***REMOVED***
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// after the anonymous dance: if an unexported field, skip
		if f.PkgPath != "" ***REMOVED*** // unexported
			continue
		***REMOVED***

		if f.Name == "" ***REMOVED***
			panic(noFieldNameToStructFieldInfoErr)
		***REMOVED***

		pv.fNames = append(pv.fNames, f.Name)

		if si == nil ***REMOVED***
			si = parseStructFieldInfo(f.Name, stag)
		***REMOVED*** else if si.encName == "" ***REMOVED***
			si.encName = f.Name
		***REMOVED***
		si.fieldName = f.Name

		pv.encNames = append(pv.encNames, si.encName)

		// si.ikind = int(f.Type.Kind())
		if len(indexstack) == 0 ***REMOVED***
			si.i = int16(j)
		***REMOVED*** else ***REMOVED***
			si.i = -1
			si.is = make([]int, len(indexstack)+1)
			copy(si.is, indexstack)
			si.is[len(indexstack)] = j
			// si.is = append(append(make([]int, 0, len(indexstack)+4), indexstack...), j)
		***REMOVED***

		if omitEmpty ***REMOVED***
			si.omitEmpty = true
		***REMOVED***
		pv.sfis = append(pv.sfis, si)
	***REMOVED***
***REMOVED***

// resolves the struct field info got from a call to rget.
// Returns a trimmed, unsorted and sorted []*structFieldInfo.
func rgetResolveSFI(x []*structFieldInfo, pv []sfiIdx) (y, z []*structFieldInfo) ***REMOVED***
	var n int
	for i, v := range x ***REMOVED***
		xn := v.encName //TODO: fieldName or encName? use encName for now.
		var found bool
		for j, k := range pv ***REMOVED***
			if k.name == xn ***REMOVED***
				// one of them must be reset to nil, and the index updated appropriately to the other one
				if len(v.is) == len(x[k.index].is) ***REMOVED***
				***REMOVED*** else if len(v.is) < len(x[k.index].is) ***REMOVED***
					pv[j].index = i
					if x[k.index] != nil ***REMOVED***
						x[k.index] = nil
						n++
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if x[i] != nil ***REMOVED***
						x[i] = nil
						n++
					***REMOVED***
				***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			pv = append(pv, sfiIdx***REMOVED***xn, i***REMOVED***)
		***REMOVED***
	***REMOVED***

	// remove all the nils
	y = make([]*structFieldInfo, len(x)-n)
	n = 0
	for _, v := range x ***REMOVED***
		if v == nil ***REMOVED***
			continue
		***REMOVED***
		y[n] = v
		n++
	***REMOVED***

	z = make([]*structFieldInfo, len(y))
	copy(z, y)
	sort.Sort(sfiSortedByEncName(z))
	return
***REMOVED***

func panicToErr(err *error) ***REMOVED***
	if recoverPanicToErr ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			//debug.PrintStack()
			panicValToErr(x, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// func doPanic(tag string, format string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
// 	params2 := make([]interface***REMOVED******REMOVED***, len(params)+1)
// 	params2[0] = tag
// 	copy(params2[1:], params)
// 	panic(fmt.Errorf("%s: "+format, params2...))
// ***REMOVED***

func isImmutableKind(k reflect.Kind) (v bool) ***REMOVED***
	return false ||
		k == reflect.Int ||
		k == reflect.Int8 ||
		k == reflect.Int16 ||
		k == reflect.Int32 ||
		k == reflect.Int64 ||
		k == reflect.Uint ||
		k == reflect.Uint8 ||
		k == reflect.Uint16 ||
		k == reflect.Uint32 ||
		k == reflect.Uint64 ||
		k == reflect.Uintptr ||
		k == reflect.Float32 ||
		k == reflect.Float64 ||
		k == reflect.Bool ||
		k == reflect.String
***REMOVED***

// these functions must be inlinable, and not call anybody
type checkOverflow struct***REMOVED******REMOVED***

func (_ checkOverflow) Float32(f float64) (overflow bool) ***REMOVED***
	if f < 0 ***REMOVED***
		f = -f
	***REMOVED***
	return math.MaxFloat32 < f && f <= math.MaxFloat64
***REMOVED***

func (_ checkOverflow) Uint(v uint64, bitsize uint8) (overflow bool) ***REMOVED***
	if bitsize == 0 || bitsize >= 64 || v == 0 ***REMOVED***
		return
	***REMOVED***
	if trunc := (v << (64 - bitsize)) >> (64 - bitsize); v != trunc ***REMOVED***
		overflow = true
	***REMOVED***
	return
***REMOVED***

func (_ checkOverflow) Int(v int64, bitsize uint8) (overflow bool) ***REMOVED***
	if bitsize == 0 || bitsize >= 64 || v == 0 ***REMOVED***
		return
	***REMOVED***
	if trunc := (v << (64 - bitsize)) >> (64 - bitsize); v != trunc ***REMOVED***
		overflow = true
	***REMOVED***
	return
***REMOVED***

func (_ checkOverflow) SignedInt(v uint64) (i int64, overflow bool) ***REMOVED***
	//e.g. -127 to 128 for int8
	pos := (v >> 63) == 0
	ui2 := v & 0x7fffffffffffffff
	if pos ***REMOVED***
		if ui2 > math.MaxInt64 ***REMOVED***
			overflow = true
			return
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if ui2 > math.MaxInt64-1 ***REMOVED***
			overflow = true
			return
		***REMOVED***
	***REMOVED***
	i = int64(v)
	return
***REMOVED***

// ------------------ SORT -----------------

func isNaN(f float64) bool ***REMOVED*** return f != f ***REMOVED***

// -----------------------

type intSlice []int64
type uintSlice []uint64
type floatSlice []float64
type boolSlice []bool
type stringSlice []string
type bytesSlice [][]byte

func (p intSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p intSlice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p intSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p uintSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p uintSlice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p uintSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p floatSlice) Len() int ***REMOVED*** return len(p) ***REMOVED***
func (p floatSlice) Less(i, j int) bool ***REMOVED***
	return p[i] < p[j] || isNaN(p[i]) && !isNaN(p[j])
***REMOVED***
func (p floatSlice) Swap(i, j int) ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p stringSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p stringSlice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p stringSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p bytesSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p bytesSlice) Less(i, j int) bool ***REMOVED*** return bytes.Compare(p[i], p[j]) == -1 ***REMOVED***
func (p bytesSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p boolSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p boolSlice) Less(i, j int) bool ***REMOVED*** return !p[i] && p[j] ***REMOVED***
func (p boolSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

// ---------------------

type intRv struct ***REMOVED***
	v int64
	r reflect.Value
***REMOVED***
type intRvSlice []intRv
type uintRv struct ***REMOVED***
	v uint64
	r reflect.Value
***REMOVED***
type uintRvSlice []uintRv
type floatRv struct ***REMOVED***
	v float64
	r reflect.Value
***REMOVED***
type floatRvSlice []floatRv
type boolRv struct ***REMOVED***
	v bool
	r reflect.Value
***REMOVED***
type boolRvSlice []boolRv
type stringRv struct ***REMOVED***
	v string
	r reflect.Value
***REMOVED***
type stringRvSlice []stringRv
type bytesRv struct ***REMOVED***
	v []byte
	r reflect.Value
***REMOVED***
type bytesRvSlice []bytesRv

func (p intRvSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p intRvSlice) Less(i, j int) bool ***REMOVED*** return p[i].v < p[j].v ***REMOVED***
func (p intRvSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p uintRvSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p uintRvSlice) Less(i, j int) bool ***REMOVED*** return p[i].v < p[j].v ***REMOVED***
func (p uintRvSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p floatRvSlice) Len() int ***REMOVED*** return len(p) ***REMOVED***
func (p floatRvSlice) Less(i, j int) bool ***REMOVED***
	return p[i].v < p[j].v || isNaN(p[i].v) && !isNaN(p[j].v)
***REMOVED***
func (p floatRvSlice) Swap(i, j int) ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p stringRvSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p stringRvSlice) Less(i, j int) bool ***REMOVED*** return p[i].v < p[j].v ***REMOVED***
func (p stringRvSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p bytesRvSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p bytesRvSlice) Less(i, j int) bool ***REMOVED*** return bytes.Compare(p[i].v, p[j].v) == -1 ***REMOVED***
func (p bytesRvSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p boolRvSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p boolRvSlice) Less(i, j int) bool ***REMOVED*** return !p[i].v && p[j].v ***REMOVED***
func (p boolRvSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

// -----------------

type bytesI struct ***REMOVED***
	v []byte
	i interface***REMOVED******REMOVED***
***REMOVED***

type bytesISlice []bytesI

func (p bytesISlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p bytesISlice) Less(i, j int) bool ***REMOVED*** return bytes.Compare(p[i].v, p[j].v) == -1 ***REMOVED***
func (p bytesISlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

// -----------------

type set []uintptr

func (s *set) add(v uintptr) (exists bool) ***REMOVED***
	// e.ci is always nil, or len >= 1
	// defer func() ***REMOVED*** fmt.Printf("$$$$$$$$$$$ cirRef Add: %v, exists: %v\n", v, exists) ***REMOVED***()
	x := *s
	if x == nil ***REMOVED***
		x = make([]uintptr, 1, 8)
		x[0] = v
		*s = x
		return
	***REMOVED***
	// typically, length will be 1. make this perform.
	if len(x) == 1 ***REMOVED***
		if j := x[0]; j == 0 ***REMOVED***
			x[0] = v
		***REMOVED*** else if j == v ***REMOVED***
			exists = true
		***REMOVED*** else ***REMOVED***
			x = append(x, v)
			*s = x
		***REMOVED***
		return
	***REMOVED***
	// check if it exists
	for _, j := range x ***REMOVED***
		if j == v ***REMOVED***
			exists = true
			return
		***REMOVED***
	***REMOVED***
	// try to replace a "deleted" slot
	for i, j := range x ***REMOVED***
		if j == 0 ***REMOVED***
			x[i] = v
			return
		***REMOVED***
	***REMOVED***
	// if unable to replace deleted slot, just append it.
	x = append(x, v)
	*s = x
	return
***REMOVED***

func (s *set) remove(v uintptr) (exists bool) ***REMOVED***
	// defer func() ***REMOVED*** fmt.Printf("$$$$$$$$$$$ cirRef Rm: %v, exists: %v\n", v, exists) ***REMOVED***()
	x := *s
	if len(x) == 0 ***REMOVED***
		return
	***REMOVED***
	if len(x) == 1 ***REMOVED***
		if x[0] == v ***REMOVED***
			x[0] = 0
		***REMOVED***
		return
	***REMOVED***
	for i, j := range x ***REMOVED***
		if j == v ***REMOVED***
			exists = true
			x[i] = 0 // set it to 0, as way to delete it.
			// copy(x[i:], x[i+1:])
			// x = x[:len(x)-1]
			return
		***REMOVED***
	***REMOVED***
	return
***REMOVED***
