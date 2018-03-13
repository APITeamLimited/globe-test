// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
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
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	scratchByteArrayLen = 32
	// initCollectionCap   = 16 // 32 is defensive. 16 is preferred.

	// Support encoding.(Binary|Text)(Unm|M)arshaler.
	// This constant flag will enable or disable it.
	supportMarshalInterfaces = true

	// for debugging, set this to false, to catch panic traces.
	// Note that this will always cause rpc tests to fail, since they need io.EOF sent via panic.
	recoverPanicToErr = true

	// arrayCacheLen is the length of the cache used in encoder or decoder for
	// allowing zero-alloc initialization.
	arrayCacheLen = 8

	// size of the cacheline: defaulting to value for archs: amd64, arm64, 386
	// should use "runtime/internal/sys".CacheLineSize, but that is not exposed.
	cacheLineSize = 64

	wordSizeBits = 32 << (^uint(0) >> 63) // strconv.IntSize
	wordSize     = wordSizeBits / 8

	maxLevelsEmbedding = 15 // use this, so structFieldInfo fits into 8 bytes
)

var (
	oneByteArr    = [1]byte***REMOVED***0***REMOVED***
	zeroByteSlice = oneByteArr[:0:0]
)

var refBitset bitset32
var pool pooler
var panicv panicHdl

func init() ***REMOVED***
	pool.init()

	refBitset.set(byte(reflect.Map))
	refBitset.set(byte(reflect.Ptr))
	refBitset.set(byte(reflect.Func))
	refBitset.set(byte(reflect.Chan))
***REMOVED***

type charEncoding uint8

const (
	cRAW charEncoding = iota
	cUTF8
	cUTF16LE
	cUTF16BE
	cUTF32LE
	cUTF32BE
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
	valueTypeTime
	valueTypeExt

	// valueTypeInvalid = 0xff
)

var valueTypeStrings = [...]string***REMOVED***
	"Unset",
	"Nil",
	"Int",
	"Uint",
	"Float",
	"Bool",
	"String",
	"Symbol",
	"Bytes",
	"Map",
	"Array",
	"Timestamp",
	"Ext",
***REMOVED***

func (x valueType) String() string ***REMOVED***
	if int(x) < len(valueTypeStrings) ***REMOVED***
		return valueTypeStrings[x]
	***REMOVED***
	return strconv.FormatInt(int64(x), 10)
***REMOVED***

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

// // sfiIdx used for tracking where a (field/enc)Name is seen in a []*structFieldInfo
// type sfiIdx struct ***REMOVED***
// 	name  string
// 	index int
// ***REMOVED***

// do not recurse if a containing type refers to an embedded type
// which refers back to its containing type (via a pointer).
// The second time this back-reference happens, break out,
// so as not to cause an infinite loop.
const rgetMaxRecursion = 2

// Anecdotally, we believe most types have <= 12 fields.
// - even Java's PMD rules set TooManyFields threshold to 15.
// However, go has embedded fields, which should be regarded as
// top level, allowing structs to possibly double or triple.
// In addition, we don't want to keep creating transient arrays,
// especially for the sfi index tracking, and the evtypes tracking.
//
// So - try to keep typeInfoLoadArray within 2K bytes
const (
	typeInfoLoadArraySfisLen   = 16
	typeInfoLoadArraySfiidxLen = 8 * 112
	typeInfoLoadArrayEtypesLen = 12
	typeInfoLoadArrayBLen      = 8 * 4
)

type typeInfoLoad struct ***REMOVED***
	// fNames   []string
	// encNames []string
	etypes []uintptr
	sfis   []structFieldInfo
***REMOVED***

type typeInfoLoadArray struct ***REMOVED***
	// fNames   [typeInfoLoadArrayLen]string
	// encNames [typeInfoLoadArrayLen]string
	sfis   [typeInfoLoadArraySfisLen]structFieldInfo
	sfiidx [typeInfoLoadArraySfiidxLen]byte
	etypes [typeInfoLoadArrayEtypesLen]uintptr
	b      [typeInfoLoadArrayBLen]byte // scratch - used for struct field names
***REMOVED***

// mirror json.Marshaler and json.Unmarshaler here,
// so we don't import the encoding/json package

type jsonMarshaler interface ***REMOVED***
	MarshalJSON() ([]byte, error)
***REMOVED***
type jsonUnmarshaler interface ***REMOVED***
	UnmarshalJSON([]byte) error
***REMOVED***

type isZeroer interface ***REMOVED***
	IsZero() bool
***REMOVED***

// type byteAccepter func(byte) bool

var (
	bigen               = binary.BigEndian
	structInfoFieldName = "_struct"

	mapStrIntfTyp  = reflect.TypeOf(map[string]interface***REMOVED******REMOVED***(nil))
	mapIntfIntfTyp = reflect.TypeOf(map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***(nil))
	intfSliceTyp   = reflect.TypeOf([]interface***REMOVED******REMOVED***(nil))
	intfTyp        = intfSliceTyp.Elem()

	reflectValTyp = reflect.TypeOf((*reflect.Value)(nil)).Elem()

	stringTyp     = reflect.TypeOf("")
	timeTyp       = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	rawExtTyp     = reflect.TypeOf(RawExt***REMOVED******REMOVED***)
	rawTyp        = reflect.TypeOf(Raw***REMOVED******REMOVED***)
	uintptrTyp    = reflect.TypeOf(uintptr(0))
	uint8Typ      = reflect.TypeOf(uint8(0))
	uint8SliceTyp = reflect.TypeOf([]uint8(nil))
	uintTyp       = reflect.TypeOf(uint(0))
	intTyp        = reflect.TypeOf(int(0))

	mapBySliceTyp = reflect.TypeOf((*MapBySlice)(nil)).Elem()

	binaryMarshalerTyp   = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	binaryUnmarshalerTyp = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()

	textMarshalerTyp   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerTyp = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

	jsonMarshalerTyp   = reflect.TypeOf((*jsonMarshaler)(nil)).Elem()
	jsonUnmarshalerTyp = reflect.TypeOf((*jsonUnmarshaler)(nil)).Elem()

	selferTyp = reflect.TypeOf((*Selfer)(nil)).Elem()
	iszeroTyp = reflect.TypeOf((*isZeroer)(nil)).Elem()

	uint8TypId      = rt2id(uint8Typ)
	uint8SliceTypId = rt2id(uint8SliceTyp)
	rawExtTypId     = rt2id(rawExtTyp)
	rawTypId        = rt2id(rawTyp)
	intfTypId       = rt2id(intfTyp)
	timeTypId       = rt2id(timeTyp)
	stringTypId     = rt2id(stringTyp)

	mapStrIntfTypId  = rt2id(mapStrIntfTyp)
	mapIntfIntfTypId = rt2id(mapIntfIntfTyp)
	intfSliceTypId   = rt2id(intfSliceTyp)
	// mapBySliceTypId  = rt2id(mapBySliceTyp)

	intBitsize  = uint8(intTyp.Bits())
	uintBitsize = uint8(uintTyp.Bits())

	bsAll0x00 = []byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0***REMOVED***
	bsAll0xff = []byte***REMOVED***0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff***REMOVED***

	chkOvf checkOverflow

	errNoFieldNameToStructFieldInfo = errors.New("no field name passed to parseStructFieldInfo")
)

var defTypeInfos = NewTypeInfos([]string***REMOVED***"codec", "json"***REMOVED***)

var immutableKindsSet = [32]bool***REMOVED***
	// reflect.Invalid:  ,
	reflect.Bool:       true,
	reflect.Int:        true,
	reflect.Int8:       true,
	reflect.Int16:      true,
	reflect.Int32:      true,
	reflect.Int64:      true,
	reflect.Uint:       true,
	reflect.Uint8:      true,
	reflect.Uint16:     true,
	reflect.Uint32:     true,
	reflect.Uint64:     true,
	reflect.Uintptr:    true,
	reflect.Float32:    true,
	reflect.Float64:    true,
	reflect.Complex64:  true,
	reflect.Complex128: true,
	// reflect.Array
	// reflect.Chan
	// reflect.Func: true,
	// reflect.Interface
	// reflect.Map
	// reflect.Ptr
	// reflect.Slice
	reflect.String: true,
	// reflect.Struct
	// reflect.UnsafePointer
***REMOVED***

// Selfer defines methods by which a value can encode or decode itself.
//
// Any type which implements Selfer will be able to encode or decode itself.
// Consequently, during (en|de)code, this takes precedence over
// (text|binary)(M|Unm)arshal or extension support.
type Selfer interface ***REMOVED***
	CodecEncodeSelf(*Encoder)
	CodecDecodeSelf(*Decoder)
***REMOVED***

// MapBySlice is a tag interface that denotes wrapped slice should encode as a map in the stream.
// The slice contains a sequence of key-value pairs.
// This affords storing a map in a specific sequence in the stream.
//
// Example usage:
//    type T1 []string         // or []int or []Point or any other "slice" type
//    func (_ T1) MapBySlice***REMOVED******REMOVED*** // T1 now implements MapBySlice, and will be encoded as a map
//    type T2 struct ***REMOVED*** KeyValues T1 ***REMOVED***
//
//    var kvs = []string***REMOVED***"one", "1", "two", "2", "three", "3"***REMOVED***
//    var v2 = T2***REMOVED*** KeyValues: T1(kvs) ***REMOVED***
//    // v2 will be encoded like the map: ***REMOVED***"KeyValues": ***REMOVED***"one": "1", "two": "2", "three": "3"***REMOVED*** ***REMOVED***
//
// The support of MapBySlice affords the following:
//   - A slice type which implements MapBySlice will be encoded as a map
//   - A slice can be decoded from a map in the stream
//   - It MUST be a slice type (not a pointer receiver) that implements MapBySlice
type MapBySlice interface ***REMOVED***
	MapBySlice()
***REMOVED***

// BasicHandle encapsulates the common options and extension functions.
//
// Deprecated: DO NOT USE DIRECTLY. EXPORTED FOR GODOC BENEFIT. WILL BE REMOVED.
type BasicHandle struct ***REMOVED***
	// BasicHandle is always a part of a different type.
	// It doesn't have to fit into it own cache lines.

	// TypeInfos is used to get the type info for any type.
	//
	// If not configured, the default TypeInfos is used, which uses struct tag keys: codec, json
	TypeInfos *TypeInfos

	// Note: BasicHandle is not comparable, due to these slices here (extHandle, intf2impls).
	// If *[]T is used instead, this becomes comparable, at the cost of extra indirection.
	// Thses slices are used all the time, so keep as slices (not pointers).

	extHandle

	intf2impls

	RPCOptions

	// ---- cache line

	DecodeOptions

	// ---- cache line

	EncodeOptions

	// noBuiltInTypeChecker
***REMOVED***

func (x *BasicHandle) getBasicHandle() *BasicHandle ***REMOVED***
	return x
***REMOVED***

func (x *BasicHandle) getTypeInfo(rtid uintptr, rt reflect.Type) (pti *typeInfo) ***REMOVED***
	if x.TypeInfos == nil ***REMOVED***
		return defTypeInfos.get(rtid, rt)
	***REMOVED***
	return x.TypeInfos.get(rtid, rt)
***REMOVED***

// Handle is the interface for a specific encoding format.
//
// Typically, a Handle is pre-configured before first time use,
// and not modified while in use. Such a pre-configured Handle
// is safe for concurrent access.
type Handle interface ***REMOVED***
	Name() string
	getBasicHandle() *BasicHandle
	recreateEncDriver(encDriver) bool
	newEncDriver(w *Encoder) encDriver
	newDecDriver(r *Decoder) decDriver
	isBinary() bool
	hasElemSeparators() bool
	// IsBuiltinType(rtid uintptr) bool
***REMOVED***

// Raw represents raw formatted bytes.
// We "blindly" store it during encode and retrieve the raw bytes during decode.
// Note: it is dangerous during encode, so we may gate the behaviour
// behind an Encode flag which must be explicitly set.
type Raw []byte

// RawExt represents raw unprocessed extension data.
// Some codecs will decode extension data as a *RawExt
// if there is no registered extension for the tag.
//
// Only one of Data or Value is nil.
// If Data is nil, then the content of the RawExt is in the Value.
type RawExt struct ***REMOVED***
	Tag uint64
	// Data is the []byte which represents the raw ext. If nil, ext is exposed in Value.
	// Data is used by codecs (e.g. binc, msgpack, simple) which do custom serialization of types
	Data []byte
	// Value represents the extension, if Data is nil.
	// Value is used by codecs (e.g. cbor, json) which leverage the format to do
	// custom serialization of the types.
	Value interface***REMOVED******REMOVED***
***REMOVED***

// BytesExt handles custom (de)serialization of types to/from []byte.
// It is used by codecs (e.g. binc, msgpack, simple) which do custom serialization of the types.
type BytesExt interface ***REMOVED***
	// WriteExt converts a value to a []byte.
	//
	// Note: v is a pointer iff the registered extension type is a struct or array kind.
	WriteExt(v interface***REMOVED******REMOVED***) []byte

	// ReadExt updates a value from a []byte.
	//
	// Note: dst is always a pointer kind to the registered extension type.
	ReadExt(dst interface***REMOVED******REMOVED***, src []byte)
***REMOVED***

// InterfaceExt handles custom (de)serialization of types to/from another interface***REMOVED******REMOVED*** value.
// The Encoder or Decoder will then handle the further (de)serialization of that known type.
//
// It is used by codecs (e.g. cbor, json) which use the format to do custom serialization of types.
type InterfaceExt interface ***REMOVED***
	// ConvertExt converts a value into a simpler interface for easy encoding
	// e.g. convert time.Time to int64.
	//
	// Note: v is a pointer iff the registered extension type is a struct or array kind.
	ConvertExt(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED***

	// UpdateExt updates a value from a simpler interface for easy decoding
	// e.g. convert int64 to time.Time.
	//
	// Note: dst is always a pointer kind to the registered extension type.
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

type extWrapper struct ***REMOVED***
	BytesExt
	InterfaceExt
***REMOVED***

type bytesExtFailer struct***REMOVED******REMOVED***

func (bytesExtFailer) WriteExt(v interface***REMOVED******REMOVED***) []byte ***REMOVED***
	panicv.errorstr("BytesExt.WriteExt is not supported")
	return nil
***REMOVED***
func (bytesExtFailer) ReadExt(v interface***REMOVED******REMOVED***, bs []byte) ***REMOVED***
	panicv.errorstr("BytesExt.ReadExt is not supported")
***REMOVED***

type interfaceExtFailer struct***REMOVED******REMOVED***

func (interfaceExtFailer) ConvertExt(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	panicv.errorstr("InterfaceExt.ConvertExt is not supported")
	return nil
***REMOVED***
func (interfaceExtFailer) UpdateExt(dest interface***REMOVED******REMOVED***, v interface***REMOVED******REMOVED***) ***REMOVED***
	panicv.errorstr("InterfaceExt.UpdateExt is not supported")
***REMOVED***

type binaryEncodingType struct***REMOVED******REMOVED***

func (binaryEncodingType) isBinary() bool ***REMOVED*** return true ***REMOVED***

type textEncodingType struct***REMOVED******REMOVED***

func (textEncodingType) isBinary() bool ***REMOVED*** return false ***REMOVED***

// noBuiltInTypes is embedded into many types which do not support builtins
// e.g. msgpack, simple, cbor.

// type noBuiltInTypeChecker struct***REMOVED******REMOVED***
// func (noBuiltInTypeChecker) IsBuiltinType(rt uintptr) bool ***REMOVED*** return false ***REMOVED***
// type noBuiltInTypes struct***REMOVED*** noBuiltInTypeChecker ***REMOVED***

type noBuiltInTypes struct***REMOVED******REMOVED***

func (noBuiltInTypes) EncodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***
func (noBuiltInTypes) DecodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***

// type noStreamingCodec struct***REMOVED******REMOVED***
// func (noStreamingCodec) CheckBreak() bool ***REMOVED*** return false ***REMOVED***
// func (noStreamingCodec) hasElemSeparators() bool ***REMOVED*** return false ***REMOVED***

type noElemSeparators struct***REMOVED******REMOVED***

func (noElemSeparators) hasElemSeparators() (v bool)            ***REMOVED*** return ***REMOVED***
func (noElemSeparators) recreateEncDriver(e encDriver) (v bool) ***REMOVED*** return ***REMOVED***

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
	rtid    uintptr
	rtidptr uintptr
	rt      reflect.Type
	tag     uint64
	ext     Ext
	_       [1]uint64 // padding
***REMOVED***

type extHandle []extTypeTagFn

// AddExt registes an encode and decode function for a reflect.Type.
// To deregister an Ext, call AddExt with nil encfn and/or nil decfn.
//
// Deprecated: Use SetBytesExt or SetInterfaceExt on the Handle instead.
func (o *extHandle) AddExt(rt reflect.Type, tag byte,
	encfn func(reflect.Value) ([]byte, error),
	decfn func(reflect.Value, []byte) error) (err error) ***REMOVED***
	if encfn == nil || decfn == nil ***REMOVED***
		return o.SetExt(rt, uint64(tag), nil)
	***REMOVED***
	return o.SetExt(rt, uint64(tag), addExtWrapper***REMOVED***encfn, decfn***REMOVED***)
***REMOVED***

// SetExt will set the extension for a tag and reflect.Type.
// Note that the type must be a named type, and specifically not a pointer or Interface.
// An error is returned if that is not honored.
// To Deregister an ext, call SetExt with nil Ext.
//
// Deprecated: Use SetBytesExt or SetInterfaceExt on the Handle instead.
func (o *extHandle) SetExt(rt reflect.Type, tag uint64, ext Ext) (err error) ***REMOVED***
	// o is a pointer, because we may need to initialize it
	rk := rt.Kind()
	for rk == reflect.Ptr ***REMOVED***
		rt = rt.Elem()
		rk = rt.Kind()
	***REMOVED***

	if rt.PkgPath() == "" || rk == reflect.Interface ***REMOVED*** // || rk == reflect.Ptr ***REMOVED***
		return fmt.Errorf("codec.Handle.SetExt: Takes named type, not a pointer or interface: %v", rt)
	***REMOVED***

	rtid := rt2id(rt)
	switch rtid ***REMOVED***
	case timeTypId, rawTypId, rawExtTypId:
		// all natively supported type, so cannot have an extension
		return // TODO: should we silently ignore, or return an error???
	***REMOVED***
	// if o == nil ***REMOVED***
	// 	return errors.New("codec.Handle.SetExt: extHandle not initialized")
	// ***REMOVED***
	o2 := *o
	// if o2 == nil ***REMOVED***
	// 	return errors.New("codec.Handle.SetExt: extHandle not initialized")
	// ***REMOVED***
	for i := range o2 ***REMOVED***
		v := &o2[i]
		if v.rtid == rtid ***REMOVED***
			v.tag, v.ext = tag, ext
			return
		***REMOVED***
	***REMOVED***
	rtidptr := rt2id(reflect.PtrTo(rt))
	*o = append(o2, extTypeTagFn***REMOVED***rtid, rtidptr, rt, tag, ext, [1]uint64***REMOVED******REMOVED******REMOVED***)
	return
***REMOVED***

func (o extHandle) getExt(rtid uintptr) (v *extTypeTagFn) ***REMOVED***
	for i := range o ***REMOVED***
		v = &o[i]
		if v.rtid == rtid || v.rtidptr == rtid ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (o extHandle) getExtForTag(tag uint64) (v *extTypeTagFn) ***REMOVED***
	for i := range o ***REMOVED***
		v = &o[i]
		if v.tag == tag ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type intf2impl struct ***REMOVED***
	rtid uintptr // for intf
	impl reflect.Type
	// _    [1]uint64 // padding // not-needed, as *intf2impl is never returned.
***REMOVED***

type intf2impls []intf2impl

// Intf2Impl maps an interface to an implementing type.
// This allows us support infering the concrete type
// and populating it when passed an interface.
// e.g. var v io.Reader can be decoded as a bytes.Buffer, etc.
//
// Passing a nil impl will clear the mapping.
func (o *intf2impls) Intf2Impl(intf, impl reflect.Type) (err error) ***REMOVED***
	if impl != nil && !impl.Implements(intf) ***REMOVED***
		return fmt.Errorf("Intf2Impl: %v does not implement %v", impl, intf)
	***REMOVED***
	rtid := rt2id(intf)
	o2 := *o
	for i := range o2 ***REMOVED***
		v := &o2[i]
		if v.rtid == rtid ***REMOVED***
			v.impl = impl
			return
		***REMOVED***
	***REMOVED***
	*o = append(o2, intf2impl***REMOVED***rtid, impl***REMOVED***)
	return
***REMOVED***

func (o intf2impls) intf2impl(rtid uintptr) (rv reflect.Value) ***REMOVED***
	for i := range o ***REMOVED***
		v := &o[i]
		if v.rtid == rtid ***REMOVED***
			if v.impl == nil ***REMOVED***
				return
			***REMOVED***
			if v.impl.Kind() == reflect.Ptr ***REMOVED***
				return reflect.New(v.impl.Elem())
			***REMOVED***
			return reflect.New(v.impl).Elem()
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

type structFieldInfoFlag uint8

const (
	_ structFieldInfoFlag = 1 << iota
	structFieldInfoFlagReady
	structFieldInfoFlagOmitEmpty
)

func (x *structFieldInfoFlag) flagSet(f structFieldInfoFlag) ***REMOVED***
	*x = *x | f
***REMOVED***

func (x *structFieldInfoFlag) flagClr(f structFieldInfoFlag) ***REMOVED***
	*x = *x &^ f
***REMOVED***

func (x structFieldInfoFlag) flagGet(f structFieldInfoFlag) bool ***REMOVED***
	return x&f != 0
***REMOVED***

func (x structFieldInfoFlag) omitEmpty() bool ***REMOVED***
	return x.flagGet(structFieldInfoFlagOmitEmpty)
***REMOVED***

func (x structFieldInfoFlag) ready() bool ***REMOVED***
	return x.flagGet(structFieldInfoFlagReady)
***REMOVED***

type structFieldInfo struct ***REMOVED***
	encName   string // encode name
	fieldName string // field name

	is  [maxLevelsEmbedding]uint16 // (recursive/embedded) field index in struct
	nis uint8                      // num levels of embedding. if 1, then it's not embedded.
	structFieldInfoFlag
***REMOVED***

func (si *structFieldInfo) setToZeroValue(v reflect.Value) ***REMOVED***
	if v, valid := si.field(v, false); valid ***REMOVED***
		v.Set(reflect.Zero(v.Type()))
	***REMOVED***
***REMOVED***

// rv returns the field of the struct.
// If anonymous, it returns an Invalid
func (si *structFieldInfo) field(v reflect.Value, update bool) (rv2 reflect.Value, valid bool) ***REMOVED***
	// replicate FieldByIndex
	for i, x := range si.is ***REMOVED***
		if uint8(i) == si.nis ***REMOVED***
			break
		***REMOVED***
		if v, valid = baseStructRv(v, update); !valid ***REMOVED***
			return
		***REMOVED***
		v = v.Field(int(x))
	***REMOVED***

	return v, true
***REMOVED***

// func (si *structFieldInfo) fieldval(v reflect.Value, update bool) reflect.Value ***REMOVED***
// 	v, _ = si.field(v, update)
// 	return v
// ***REMOVED***

func parseStructInfo(stag string) (toArray, omitEmpty bool, keytype valueType) ***REMOVED***
	keytype = valueTypeString // default
	if stag == "" ***REMOVED***
		return
	***REMOVED***
	for i, s := range strings.Split(stag, ",") ***REMOVED***
		if i == 0 ***REMOVED***
		***REMOVED*** else ***REMOVED***
			switch s ***REMOVED***
			case "omitempty":
				omitEmpty = true
			case "toarray":
				toArray = true
			case "int":
				keytype = valueTypeInt
			case "uint":
				keytype = valueTypeUint
			case "float":
				keytype = valueTypeFloat
				// case "bool":
				// 	keytype = valueTypeBool
			case "string":
				keytype = valueTypeString
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (si *structFieldInfo) parseTag(stag string) ***REMOVED***
	// if fname == "" ***REMOVED***
	// 	panic(errNoFieldNameToStructFieldInfo)
	// ***REMOVED***

	if stag == "" ***REMOVED***
		return
	***REMOVED***
	for i, s := range strings.Split(stag, ",") ***REMOVED***
		if i == 0 ***REMOVED***
			if s != "" ***REMOVED***
				si.encName = s
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			switch s ***REMOVED***
			case "omitempty":
				si.flagSet(structFieldInfoFlagOmitEmpty)
				// si.omitEmpty = true
				// case "toarray":
				// 	si.toArray = true
			***REMOVED***
		***REMOVED***
	***REMOVED***
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

const structFieldNodeNumToCache = 4

type structFieldNodeCache struct ***REMOVED***
	rv  [structFieldNodeNumToCache]reflect.Value
	idx [structFieldNodeNumToCache]uint32
	num uint8
***REMOVED***

func (x *structFieldNodeCache) get(key uint32) (fv reflect.Value, valid bool) ***REMOVED***
	for i, k := range &x.idx ***REMOVED***
		if uint8(i) == x.num ***REMOVED***
			return // break
		***REMOVED***
		if key == k ***REMOVED***
			return x.rv[i], true
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (x *structFieldNodeCache) tryAdd(fv reflect.Value, key uint32) ***REMOVED***
	if x.num < structFieldNodeNumToCache ***REMOVED***
		x.rv[x.num] = fv
		x.idx[x.num] = key
		x.num++
		return
	***REMOVED***
***REMOVED***

type structFieldNode struct ***REMOVED***
	v      reflect.Value
	cache2 structFieldNodeCache
	cache3 structFieldNodeCache
	update bool
***REMOVED***

func (x *structFieldNode) field(si *structFieldInfo) (fv reflect.Value) ***REMOVED***
	// return si.fieldval(x.v, x.update)
	// Note: we only cache if nis=2 or nis=3 i.e. up to 2 levels of embedding
	// This mostly saves us time on the repeated calls to v.Elem, v.Field, etc.
	var valid bool
	switch si.nis ***REMOVED***
	case 1:
		fv = x.v.Field(int(si.is[0]))
	case 2:
		if fv, valid = x.cache2.get(uint32(si.is[0])); valid ***REMOVED***
			fv = fv.Field(int(si.is[1]))
			return
		***REMOVED***
		fv = x.v.Field(int(si.is[0]))
		if fv, valid = baseStructRv(fv, x.update); !valid ***REMOVED***
			return
		***REMOVED***
		x.cache2.tryAdd(fv, uint32(si.is[0]))
		fv = fv.Field(int(si.is[1]))
	case 3:
		var key uint32 = uint32(si.is[0])<<16 | uint32(si.is[1])
		if fv, valid = x.cache3.get(key); valid ***REMOVED***
			fv = fv.Field(int(si.is[2]))
			return
		***REMOVED***
		fv = x.v.Field(int(si.is[0]))
		if fv, valid = baseStructRv(fv, x.update); !valid ***REMOVED***
			return
		***REMOVED***
		fv = fv.Field(int(si.is[1]))
		if fv, valid = baseStructRv(fv, x.update); !valid ***REMOVED***
			return
		***REMOVED***
		x.cache3.tryAdd(fv, key)
		fv = fv.Field(int(si.is[2]))
	default:
		fv, _ = si.field(x.v, x.update)
	***REMOVED***
	return
***REMOVED***

func baseStructRv(v reflect.Value, update bool) (v2 reflect.Value, valid bool) ***REMOVED***
	for v.Kind() == reflect.Ptr ***REMOVED***
		if v.IsNil() ***REMOVED***
			if !update ***REMOVED***
				return
			***REMOVED***
			v.Set(reflect.New(v.Type().Elem()))
		***REMOVED***
		v = v.Elem()
	***REMOVED***
	return v, true
***REMOVED***

type typeInfoFlag uint8

const (
	typeInfoFlagComparable = 1 << iota
	typeInfoFlagIsZeroer
	typeInfoFlagIsZeroerPtr
)

// typeInfo keeps information about each (non-ptr) type referenced in the encode/decode sequence.
//
// During an encode/decode sequence, we work as below:
//   - If base is a built in type, en/decode base value
//   - If base is registered as an extension, en/decode base value
//   - If type is binary(M/Unm)arshaler, call Binary(M/Unm)arshal method
//   - If type is text(M/Unm)arshaler, call Text(M/Unm)arshal method
//   - Else decode appropriately based on the reflect.Kind
type typeInfo struct ***REMOVED***
	rt      reflect.Type
	elem    reflect.Type
	pkgpath string

	rtid uintptr
	// rv0  reflect.Value // saved zero value, used if immutableKind

	numMeth uint16 // number of methods
	kind    uint8
	chandir uint8

	anyOmitEmpty bool      // true if a struct, and any of the fields are tagged "omitempty"
	toArray      bool      // whether this (struct) type should be encoded as an array
	keyType      valueType // if struct, how is the field name stored in a stream? default is string
	mbs          bool      // base type (T or *T) is a MapBySlice

	// ---- cpu cache line boundary?
	sfiSort []*structFieldInfo // sorted. Used when enc/dec struct to map.
	sfiSrc  []*structFieldInfo // unsorted. Used when enc/dec struct to array.

	key reflect.Type

	// ---- cpu cache line boundary?
	// sfis         []structFieldInfo // all sfi, in src order, as created.
	sfiNamesSort []byte // all names, with indexes into the sfiSort

	// format of marshal type fields below: [btj][mu]p? OR csp?

	bm  bool // T is a binaryMarshaler
	bmp bool // *T is a binaryMarshaler
	bu  bool // T is a binaryUnmarshaler
	bup bool // *T is a binaryUnmarshaler
	tm  bool // T is a textMarshaler
	tmp bool // *T is a textMarshaler
	tu  bool // T is a textUnmarshaler
	tup bool // *T is a textUnmarshaler

	jm  bool // T is a jsonMarshaler
	jmp bool // *T is a jsonMarshaler
	ju  bool // T is a jsonUnmarshaler
	jup bool // *T is a jsonUnmarshaler
	cs  bool // T is a Selfer
	csp bool // *T is a Selfer

	// other flags, with individual bits representing if set.
	flags typeInfoFlag

	// _ [2]byte   // padding
	_ [3]uint64 // padding
***REMOVED***

func (ti *typeInfo) isFlag(f typeInfoFlag) bool ***REMOVED***
	return ti.flags&f != 0
***REMOVED***

func (ti *typeInfo) indexForEncName(name []byte) (index int16) ***REMOVED***
	var sn []byte
	if len(name)+2 <= 32 ***REMOVED***
		var buf [32]byte // should not escape
		sn = buf[:len(name)+2]
	***REMOVED*** else ***REMOVED***
		sn = make([]byte, len(name)+2)
	***REMOVED***
	copy(sn[1:], name)
	sn[0], sn[len(sn)-1] = tiSep2(name), 0xff
	j := bytes.Index(ti.sfiNamesSort, sn)
	if j < 0 ***REMOVED***
		return -1
	***REMOVED***
	index = int16(uint16(ti.sfiNamesSort[j+len(sn)+1]) | uint16(ti.sfiNamesSort[j+len(sn)])<<8)
	return
***REMOVED***

type rtid2ti struct ***REMOVED***
	rtid uintptr
	ti   *typeInfo
***REMOVED***

// TypeInfos caches typeInfo for each type on first inspection.
//
// It is configured with a set of tag keys, which are used to get
// configuration for the type.
type TypeInfos struct ***REMOVED***
	// infos: formerly map[uintptr]*typeInfo, now *[]rtid2ti, 2 words expected
	infos atomicTypeInfoSlice
	mu    sync.Mutex
	tags  []string
	_     [2]uint64 // padding
***REMOVED***

// NewTypeInfos creates a TypeInfos given a set of struct tags keys.
//
// This allows users customize the struct tag keys which contain configuration
// of their types.
func NewTypeInfos(tags []string) *TypeInfos ***REMOVED***
	return &TypeInfos***REMOVED***tags: tags***REMOVED***
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

func (x *TypeInfos) find(s []rtid2ti, rtid uintptr) (idx int, ti *typeInfo) ***REMOVED***
	// binary search. adapted from sort/search.go.
	// if sp == nil ***REMOVED***
	// 	return -1, nil
	// ***REMOVED***
	// s := *sp
	h, i, j := 0, 0, len(s)
	for i < j ***REMOVED***
		h = i + (j-i)/2
		if s[h].rtid < rtid ***REMOVED***
			i = h + 1
		***REMOVED*** else ***REMOVED***
			j = h
		***REMOVED***
	***REMOVED***
	if i < len(s) && s[i].rtid == rtid ***REMOVED***
		return i, s[i].ti
	***REMOVED***
	return i, nil
***REMOVED***

func (x *TypeInfos) get(rtid uintptr, rt reflect.Type) (pti *typeInfo) ***REMOVED***
	sp := x.infos.load()
	var idx int
	if sp != nil ***REMOVED***
		idx, pti = x.find(sp, rtid)
		if pti != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	rk := rt.Kind()

	if rk == reflect.Ptr ***REMOVED*** // || (rk == reflect.Interface && rtid != intfTypId) ***REMOVED***
		panicv.errorf("invalid kind passed to TypeInfos.get: %v - %v", rk, rt)
	***REMOVED***

	// do not hold lock while computing this.
	// it may lead to duplication, but that's ok.
	ti := typeInfo***REMOVED***rt: rt, rtid: rtid, kind: uint8(rk), pkgpath: rt.PkgPath()***REMOVED***
	// ti.rv0 = reflect.Zero(rt)

	// ti.comparable = rt.Comparable()
	ti.numMeth = uint16(rt.NumMethod())

	ti.bm, ti.bmp = implIntf(rt, binaryMarshalerTyp)
	ti.bu, ti.bup = implIntf(rt, binaryUnmarshalerTyp)
	ti.tm, ti.tmp = implIntf(rt, textMarshalerTyp)
	ti.tu, ti.tup = implIntf(rt, textUnmarshalerTyp)
	ti.jm, ti.jmp = implIntf(rt, jsonMarshalerTyp)
	ti.ju, ti.jup = implIntf(rt, jsonUnmarshalerTyp)
	ti.cs, ti.csp = implIntf(rt, selferTyp)

	b1, b2 := implIntf(rt, iszeroTyp)
	if b1 ***REMOVED***
		ti.flags |= typeInfoFlagIsZeroer
	***REMOVED***
	if b2 ***REMOVED***
		ti.flags |= typeInfoFlagIsZeroerPtr
	***REMOVED***
	if rt.Comparable() ***REMOVED***
		ti.flags |= typeInfoFlagComparable
	***REMOVED***

	switch rk ***REMOVED***
	case reflect.Struct:
		var omitEmpty bool
		if f, ok := rt.FieldByName(structInfoFieldName); ok ***REMOVED***
			ti.toArray, omitEmpty, ti.keyType = parseStructInfo(x.structTag(f.Tag))
		***REMOVED*** else ***REMOVED***
			ti.keyType = valueTypeString
		***REMOVED***
		pp, pi := pool.tiLoad()
		pv := pi.(*typeInfoLoadArray)
		pv.etypes[0] = ti.rtid
		// vv := typeInfoLoad***REMOVED***pv.fNames[:0], pv.encNames[:0], pv.etypes[:1], pv.sfis[:0]***REMOVED***
		vv := typeInfoLoad***REMOVED***pv.etypes[:1], pv.sfis[:0]***REMOVED***
		x.rget(rt, rtid, omitEmpty, nil, &vv)
		// ti.sfis = vv.sfis
		ti.sfiSrc, ti.sfiSort, ti.sfiNamesSort, ti.anyOmitEmpty = rgetResolveSFI(rt, vv.sfis, pv)
		pp.Put(pi)
	case reflect.Map:
		ti.elem = rt.Elem()
		ti.key = rt.Key()
	case reflect.Slice:
		ti.mbs, _ = implIntf(rt, mapBySliceTyp)
		ti.elem = rt.Elem()
	case reflect.Chan:
		ti.elem = rt.Elem()
		ti.chandir = uint8(rt.ChanDir())
	case reflect.Array, reflect.Ptr:
		ti.elem = rt.Elem()
	***REMOVED***
	// sfi = sfiSrc

	x.mu.Lock()
	sp = x.infos.load()
	if sp == nil ***REMOVED***
		pti = &ti
		vs := []rtid2ti***REMOVED******REMOVED***rtid, pti***REMOVED******REMOVED***
		x.infos.store(vs)
	***REMOVED*** else ***REMOVED***
		idx, pti = x.find(sp, rtid)
		if pti == nil ***REMOVED***
			pti = &ti
			vs := make([]rtid2ti, len(sp)+1)
			copy(vs, sp[:idx])
			copy(vs[idx+1:], sp[idx:])
			vs[idx] = rtid2ti***REMOVED***rtid, pti***REMOVED***
			x.infos.store(vs)
		***REMOVED***
	***REMOVED***
	x.mu.Unlock()
	return
***REMOVED***

func (x *TypeInfos) rget(rt reflect.Type, rtid uintptr, omitEmpty bool,
	indexstack []uint16, pv *typeInfoLoad) ***REMOVED***
	// Read up fields and store how to access the value.
	//
	// It uses go's rules for message selectors,
	// which say that the field with the shallowest depth is selected.
	//
	// Note: we consciously use slices, not a map, to simulate a set.
	//       Typically, types have < 16 fields,
	//       and iteration using equals is faster than maps there
	flen := rt.NumField()
	if flen > (1<<maxLevelsEmbedding - 1) ***REMOVED***
		panicv.errorf("codec: types with > %v fields are not supported - has %v fields",
			(1<<maxLevelsEmbedding - 1), flen)
	***REMOVED***
	// pv.sfis = make([]structFieldInfo, flen)
LOOP:
	for j, jlen := uint16(0), uint16(flen); j < jlen; j++ ***REMOVED***
		f := rt.Field(int(j))
		fkind := f.Type.Kind()
		// skip if a func type, or is unexported, or structTag value == "-"
		switch fkind ***REMOVED***
		case reflect.Func, reflect.Complex64, reflect.Complex128, reflect.UnsafePointer:
			continue LOOP
		***REMOVED***

		isUnexported := f.PkgPath != ""
		if isUnexported && !f.Anonymous ***REMOVED***
			continue
		***REMOVED***
		stag := x.structTag(f.Tag)
		if stag == "-" ***REMOVED***
			continue
		***REMOVED***
		var si structFieldInfo
		var parsed bool
		// if anonymous and no struct tag (or it's blank),
		// and a struct (or pointer to struct), inline it.
		if f.Anonymous && fkind != reflect.Interface ***REMOVED***
			// ^^ redundant but ok: per go spec, an embedded pointer type cannot be to an interface
			ft := f.Type
			isPtr := ft.Kind() == reflect.Ptr
			for ft.Kind() == reflect.Ptr ***REMOVED***
				ft = ft.Elem()
			***REMOVED***
			isStruct := ft.Kind() == reflect.Struct

			// Ignore embedded fields of unexported non-struct types.
			// Also, from go1.10, ignore pointers to unexported struct types
			// because unmarshal cannot assign a new struct to an unexported field.
			// See https://golang.org/issue/21357
			if (isUnexported && !isStruct) || (!allowSetUnexportedEmbeddedPtr && isUnexported && isPtr) ***REMOVED***
				continue
			***REMOVED***
			doInline := stag == ""
			if !doInline ***REMOVED***
				si.parseTag(stag)
				parsed = true
				doInline = si.encName == ""
				// doInline = si.isZero()
			***REMOVED***
			if doInline && isStruct ***REMOVED***
				// if etypes contains this, don't call rget again (as fields are already seen here)
				ftid := rt2id(ft)
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
					indexstack2 := make([]uint16, len(indexstack)+1)
					copy(indexstack2, indexstack)
					indexstack2[len(indexstack)] = j
					// indexstack2 := append(append(make([]int, 0, len(indexstack)+4), indexstack...), j)
					x.rget(ft, ftid, omitEmpty, indexstack2, pv)
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED***

		// after the anonymous dance: if an unexported field, skip
		if isUnexported ***REMOVED***
			continue
		***REMOVED***

		if f.Name == "" ***REMOVED***
			panic(errNoFieldNameToStructFieldInfo)
		***REMOVED***

		// pv.fNames = append(pv.fNames, f.Name)
		// if si.encName == "" ***REMOVED***

		if !parsed ***REMOVED***
			si.encName = f.Name
			si.parseTag(stag)
			parsed = true
		***REMOVED*** else if si.encName == "" ***REMOVED***
			si.encName = f.Name
		***REMOVED***
		si.fieldName = f.Name
		si.flagSet(structFieldInfoFlagReady)

		// pv.encNames = append(pv.encNames, si.encName)

		// si.ikind = int(f.Type.Kind())
		if len(indexstack) > maxLevelsEmbedding-1 ***REMOVED***
			panicv.errorf("codec: only supports up to %v depth of embedding - type has %v depth",
				maxLevelsEmbedding-1, len(indexstack))
		***REMOVED***
		si.nis = uint8(len(indexstack)) + 1
		copy(si.is[:], indexstack)
		si.is[len(indexstack)] = j

		if omitEmpty ***REMOVED***
			si.flagSet(structFieldInfoFlagOmitEmpty)
		***REMOVED***
		pv.sfis = append(pv.sfis, si)
	***REMOVED***
***REMOVED***

func tiSep(name string) uint8 ***REMOVED***
	// (xn[0]%64) // (between 192-255 - outside ascii BMP)
	// return 0xfe - (name[0] & 63)
	// return 0xfe - (name[0] & 63) - uint8(len(name))
	// return 0xfe - (name[0] & 63) - uint8(len(name)&63)
	// return ((0xfe - (name[0] & 63)) & 0xf8) | (uint8(len(name) & 0x07))
	return 0xfe - (name[0] & 63) - uint8(len(name)&63)
***REMOVED***

func tiSep2(name []byte) uint8 ***REMOVED***
	return 0xfe - (name[0] & 63) - uint8(len(name)&63)
***REMOVED***

// resolves the struct field info got from a call to rget.
// Returns a trimmed, unsorted and sorted []*structFieldInfo.
func rgetResolveSFI(rt reflect.Type, x []structFieldInfo, pv *typeInfoLoadArray) (
	y, z []*structFieldInfo, ss []byte, anyOmitEmpty bool) ***REMOVED***
	sa := pv.sfiidx[:0]
	sn := pv.b[:]
	n := len(x)

	var xn string
	var ui uint16
	var sep byte

	for i := range x ***REMOVED***
		ui = uint16(i)
		xn = x[i].encName // fieldName or encName? use encName for now.
		if len(xn)+2 > cap(pv.b) ***REMOVED***
			sn = make([]byte, len(xn)+2)
		***REMOVED*** else ***REMOVED***
			sn = sn[:len(xn)+2]
		***REMOVED***
		// use a custom sep, so that misses are less frequent,
		// since the sep (first char in search) is as unique as first char in field name.
		sep = tiSep(xn)
		sn[0], sn[len(sn)-1] = sep, 0xff
		copy(sn[1:], xn)
		j := bytes.Index(sa, sn)
		if j == -1 ***REMOVED***
			sa = append(sa, sep)
			sa = append(sa, xn...)
			sa = append(sa, 0xff, byte(ui>>8), byte(ui))
		***REMOVED*** else ***REMOVED***
			index := uint16(sa[j+len(sn)+1]) | uint16(sa[j+len(sn)])<<8
			// one of them must be reset to nil,
			// and the index updated appropriately to the other one
			if x[i].nis == x[index].nis ***REMOVED***
			***REMOVED*** else if x[i].nis < x[index].nis ***REMOVED***
				sa[j+len(sn)], sa[j+len(sn)+1] = byte(ui>>8), byte(ui)
				if x[index].ready() ***REMOVED***
					x[index].flagClr(structFieldInfoFlagReady)
					n--
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if x[i].ready() ***REMOVED***
					x[i].flagClr(structFieldInfoFlagReady)
					n--
				***REMOVED***
			***REMOVED***
		***REMOVED***

	***REMOVED***
	var w []structFieldInfo
	sharingArray := len(x) <= typeInfoLoadArraySfisLen // sharing array with typeInfoLoadArray
	if sharingArray ***REMOVED***
		w = make([]structFieldInfo, n)
	***REMOVED***

	// remove all the nils (non-ready)
	y = make([]*structFieldInfo, n)
	n = 0
	var sslen int
	for i := range x ***REMOVED***
		if !x[i].ready() ***REMOVED***
			continue
		***REMOVED***
		if !anyOmitEmpty && x[i].omitEmpty() ***REMOVED***
			anyOmitEmpty = true
		***REMOVED***
		if sharingArray ***REMOVED***
			w[n] = x[i]
			y[n] = &w[n]
		***REMOVED*** else ***REMOVED***
			y[n] = &x[i]
		***REMOVED***
		sslen = sslen + len(x[i].encName) + 4
		n++
	***REMOVED***
	if n != len(y) ***REMOVED***
		panicv.errorf("failure reading struct %v - expecting %d of %d valid fields, got %d",
			rt, len(y), len(x), n)
	***REMOVED***

	z = make([]*structFieldInfo, len(y))
	copy(z, y)
	sort.Sort(sfiSortedByEncName(z))

	sharingArray = len(sa) <= typeInfoLoadArraySfiidxLen
	if sharingArray ***REMOVED***
		ss = make([]byte, 0, sslen)
	***REMOVED*** else ***REMOVED***
		ss = sa[:0] // reuse the newly made sa array if necessary
	***REMOVED***
	for i := range z ***REMOVED***
		xn = z[i].encName
		sep = tiSep(xn)
		ui = uint16(i)
		ss = append(ss, sep)
		ss = append(ss, xn...)
		ss = append(ss, 0xff, byte(ui>>8), byte(ui))
	***REMOVED***
	return
***REMOVED***

func implIntf(rt, iTyp reflect.Type) (base bool, indir bool) ***REMOVED***
	return rt.Implements(iTyp), reflect.PtrTo(rt).Implements(iTyp)
***REMOVED***

// isEmptyStruct is only called from isEmptyValue, and checks if a struct is empty:
//    - does it implement IsZero() bool
//    - is it comparable, and can i compare directly using ==
//    - if checkStruct, then walk through the encodable fields
//      and check if they are empty or not.
func isEmptyStruct(v reflect.Value, tinfos *TypeInfos, deref, checkStruct bool) bool ***REMOVED***
	// v is a struct kind - no need to check again.
	// We only check isZero on a struct kind, to reduce the amount of times
	// that we lookup the rtid and typeInfo for each type as we walk the tree.

	vt := v.Type()
	rtid := rt2id(vt)
	if tinfos == nil ***REMOVED***
		tinfos = defTypeInfos
	***REMOVED***
	ti := tinfos.get(rtid, vt)
	if ti.rtid == timeTypId ***REMOVED***
		return rv2i(v).(time.Time).IsZero()
	***REMOVED***
	if ti.isFlag(typeInfoFlagIsZeroerPtr) && v.CanAddr() ***REMOVED***
		return rv2i(v.Addr()).(isZeroer).IsZero()
	***REMOVED***
	if ti.isFlag(typeInfoFlagIsZeroer) ***REMOVED***
		return rv2i(v).(isZeroer).IsZero()
	***REMOVED***
	if ti.isFlag(typeInfoFlagComparable) ***REMOVED***
		return rv2i(v) == rv2i(reflect.Zero(vt))
	***REMOVED***
	if !checkStruct ***REMOVED***
		return false
	***REMOVED***
	// We only care about what we can encode/decode,
	// so that is what we use to check omitEmpty.
	for _, si := range ti.sfiSrc ***REMOVED***
		sfv, valid := si.field(v, false)
		if valid && !isEmptyValue(sfv, tinfos, deref, checkStruct) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// func roundFloat(x float64) float64 ***REMOVED***
// 	t := math.Trunc(x)
// 	if math.Abs(x-t) >= 0.5 ***REMOVED***
// 		return t + math.Copysign(1, x)
// 	***REMOVED***
// 	return t
// ***REMOVED***

func panicToErr(h errstrDecorator, err *error) ***REMOVED***
	if recoverPanicToErr ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			// fmt.Printf("panic'ing with: %v\n", x)
			// debug.PrintStack()
			panicValToErr(h, x, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func panicToErrs2(h errstrDecorator, err1, err2 *error) ***REMOVED***
	if recoverPanicToErr ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			panicValToErr(h, x, err1)
			panicValToErr(h, x, err2)
		***REMOVED***
	***REMOVED***
***REMOVED***

func panicValToErr(h errstrDecorator, v interface***REMOVED******REMOVED***, err *error) ***REMOVED***
	switch xerr := v.(type) ***REMOVED***
	case nil:
	case error:
		switch xerr ***REMOVED***
		case nil:
		case io.EOF, io.ErrUnexpectedEOF, errEncoderNotInitialized, errDecoderNotInitialized:
			// treat as special (bubble up)
			*err = xerr
		default:
			h.wrapErrstr(xerr.Error(), err)
		***REMOVED***
	case string:
		if xerr != "" ***REMOVED***
			h.wrapErrstr(xerr, err)
		***REMOVED***
	case fmt.Stringer:
		if xerr != nil ***REMOVED***
			h.wrapErrstr(xerr.String(), err)
		***REMOVED***
	default:
		h.wrapErrstr(v, err)
	***REMOVED***
***REMOVED***

func isImmutableKind(k reflect.Kind) (v bool) ***REMOVED***
	return immutableKindsSet[k]
***REMOVED***

// ----

type codecFnInfo struct ***REMOVED***
	ti    *typeInfo
	xfFn  Ext
	xfTag uint64
	seq   seqType
	addrD bool
	addrF bool // if addrD, this says whether decode function can take a value or a ptr
	addrE bool
	ready bool // ready to use
***REMOVED***

// codecFn encapsulates the captured variables and the encode function.
// This way, we only do some calculations one times, and pass to the
// code block that should be called (encapsulated in a function)
// instead of executing the checks every time.
type codecFn struct ***REMOVED***
	i  codecFnInfo
	fe func(*Encoder, *codecFnInfo, reflect.Value)
	fd func(*Decoder, *codecFnInfo, reflect.Value)
	_  [1]uint64 // padding
***REMOVED***

type codecRtidFn struct ***REMOVED***
	rtid uintptr
	fn   *codecFn
***REMOVED***

type codecFner struct ***REMOVED***
	// hh Handle
	h  *BasicHandle
	s  []codecRtidFn
	be bool
	js bool
	_  [6]byte   // padding
	_  [3]uint64 // padding
***REMOVED***

func (c *codecFner) reset(hh Handle) ***REMOVED***
	bh := hh.getBasicHandle()
	// only reset iff extensions changed or *TypeInfos changed
	var hhSame = true &&
		c.h == bh && c.h.TypeInfos == bh.TypeInfos &&
		len(c.h.extHandle) == len(bh.extHandle) &&
		(len(c.h.extHandle) == 0 || &c.h.extHandle[0] == &bh.extHandle[0])
	if !hhSame ***REMOVED***
		// c.hh = hh
		c.h, bh = bh, c.h // swap both
		_, c.js = hh.(*JsonHandle)
		c.be = hh.isBinary()
		for i := range c.s ***REMOVED***
			c.s[i].fn.i.ready = false
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *codecFner) get(rt reflect.Type, checkFastpath, checkCodecSelfer bool) (fn *codecFn) ***REMOVED***
	rtid := rt2id(rt)

	for _, x := range c.s ***REMOVED***
		if x.rtid == rtid ***REMOVED***
			// if rtid exists, then there's a *codenFn attached (non-nil)
			fn = x.fn
			if fn.i.ready ***REMOVED***
				return
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	var ti *typeInfo
	if fn == nil ***REMOVED***
		fn = new(codecFn)
		if c.s == nil ***REMOVED***
			c.s = make([]codecRtidFn, 0, 8)
		***REMOVED***
		c.s = append(c.s, codecRtidFn***REMOVED***rtid, fn***REMOVED***)
	***REMOVED*** else ***REMOVED***
		ti = fn.i.ti
		*fn = codecFn***REMOVED******REMOVED***
		fn.i.ti = ti
		// fn.fe, fn.fd = nil, nil
	***REMOVED***
	fi := &(fn.i)
	fi.ready = true
	if ti == nil ***REMOVED***
		ti = c.h.getTypeInfo(rtid, rt)
		fi.ti = ti
	***REMOVED***

	rk := reflect.Kind(ti.kind)

	if checkCodecSelfer && (ti.cs || ti.csp) ***REMOVED***
		fn.fe = (*Encoder).selferMarshal
		fn.fd = (*Decoder).selferUnmarshal
		fi.addrF = true
		fi.addrD = ti.csp
		fi.addrE = ti.csp
	***REMOVED*** else if rtid == timeTypId ***REMOVED***
		fn.fe = (*Encoder).kTime
		fn.fd = (*Decoder).kTime
	***REMOVED*** else if rtid == rawTypId ***REMOVED***
		fn.fe = (*Encoder).raw
		fn.fd = (*Decoder).raw
	***REMOVED*** else if rtid == rawExtTypId ***REMOVED***
		fn.fe = (*Encoder).rawExt
		fn.fd = (*Decoder).rawExt
		fi.addrF = true
		fi.addrD = true
		fi.addrE = true
	***REMOVED*** else if xfFn := c.h.getExt(rtid); xfFn != nil ***REMOVED***
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*Encoder).ext
		fn.fd = (*Decoder).ext
		fi.addrF = true
		fi.addrD = true
		if rk == reflect.Struct || rk == reflect.Array ***REMOVED***
			fi.addrE = true
		***REMOVED***
	***REMOVED*** else if supportMarshalInterfaces && c.be && (ti.bm || ti.bmp) && (ti.bu || ti.bup) ***REMOVED***
		fn.fe = (*Encoder).binaryMarshal
		fn.fd = (*Decoder).binaryUnmarshal
		fi.addrF = true
		fi.addrD = ti.bup
		fi.addrE = ti.bmp
	***REMOVED*** else if supportMarshalInterfaces && !c.be && c.js && (ti.jm || ti.jmp) && (ti.ju || ti.jup) ***REMOVED***
		//If JSON, we should check JSONMarshal before textMarshal
		fn.fe = (*Encoder).jsonMarshal
		fn.fd = (*Decoder).jsonUnmarshal
		fi.addrF = true
		fi.addrD = ti.jup
		fi.addrE = ti.jmp
	***REMOVED*** else if supportMarshalInterfaces && !c.be && (ti.tm || ti.tmp) && (ti.tu || ti.tup) ***REMOVED***
		fn.fe = (*Encoder).textMarshal
		fn.fd = (*Decoder).textUnmarshal
		fi.addrF = true
		fi.addrD = ti.tup
		fi.addrE = ti.tmp
	***REMOVED*** else ***REMOVED***
		if fastpathEnabled && checkFastpath && (rk == reflect.Map || rk == reflect.Slice) ***REMOVED***
			if ti.pkgpath == "" ***REMOVED*** // un-named slice or map
				if idx := fastpathAV.index(rtid); idx != -1 ***REMOVED***
					fn.fe = fastpathAV[idx].encfn
					fn.fd = fastpathAV[idx].decfn
					fi.addrD = true
					fi.addrF = false
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// use mapping for underlying type if there
				var rtu reflect.Type
				if rk == reflect.Map ***REMOVED***
					rtu = reflect.MapOf(ti.key, ti.elem)
				***REMOVED*** else ***REMOVED***
					rtu = reflect.SliceOf(ti.elem)
				***REMOVED***
				rtuid := rt2id(rtu)
				if idx := fastpathAV.index(rtuid); idx != -1 ***REMOVED***
					xfnf := fastpathAV[idx].encfn
					xrt := fastpathAV[idx].rt
					fn.fe = func(e *Encoder, xf *codecFnInfo, xrv reflect.Value) ***REMOVED***
						xfnf(e, xf, xrv.Convert(xrt))
					***REMOVED***
					fi.addrD = true
					fi.addrF = false // meaning it can be an address(ptr) or a value
					xfnf2 := fastpathAV[idx].decfn
					fn.fd = func(d *Decoder, xf *codecFnInfo, xrv reflect.Value) ***REMOVED***
						if xrv.Kind() == reflect.Ptr ***REMOVED***
							xfnf2(d, xf, xrv.Convert(reflect.PtrTo(xrt)))
						***REMOVED*** else ***REMOVED***
							xfnf2(d, xf, xrv.Convert(xrt))
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if fn.fe == nil && fn.fd == nil ***REMOVED***
			switch rk ***REMOVED***
			case reflect.Bool:
				fn.fe = (*Encoder).kBool
				fn.fd = (*Decoder).kBool
			case reflect.String:
				fn.fe = (*Encoder).kString
				fn.fd = (*Decoder).kString
			case reflect.Int:
				fn.fd = (*Decoder).kInt
				fn.fe = (*Encoder).kInt
			case reflect.Int8:
				fn.fe = (*Encoder).kInt8
				fn.fd = (*Decoder).kInt8
			case reflect.Int16:
				fn.fe = (*Encoder).kInt16
				fn.fd = (*Decoder).kInt16
			case reflect.Int32:
				fn.fe = (*Encoder).kInt32
				fn.fd = (*Decoder).kInt32
			case reflect.Int64:
				fn.fe = (*Encoder).kInt64
				fn.fd = (*Decoder).kInt64
			case reflect.Uint:
				fn.fd = (*Decoder).kUint
				fn.fe = (*Encoder).kUint
			case reflect.Uint8:
				fn.fe = (*Encoder).kUint8
				fn.fd = (*Decoder).kUint8
			case reflect.Uint16:
				fn.fe = (*Encoder).kUint16
				fn.fd = (*Decoder).kUint16
			case reflect.Uint32:
				fn.fe = (*Encoder).kUint32
				fn.fd = (*Decoder).kUint32
			case reflect.Uint64:
				fn.fe = (*Encoder).kUint64
				fn.fd = (*Decoder).kUint64
			case reflect.Uintptr:
				fn.fe = (*Encoder).kUintptr
				fn.fd = (*Decoder).kUintptr
			case reflect.Float32:
				fn.fe = (*Encoder).kFloat32
				fn.fd = (*Decoder).kFloat32
			case reflect.Float64:
				fn.fe = (*Encoder).kFloat64
				fn.fd = (*Decoder).kFloat64
			case reflect.Invalid:
				fn.fe = (*Encoder).kInvalid
				fn.fd = (*Decoder).kErr
			case reflect.Chan:
				fi.seq = seqTypeChan
				fn.fe = (*Encoder).kSlice
				fn.fd = (*Decoder).kSlice
			case reflect.Slice:
				fi.seq = seqTypeSlice
				fn.fe = (*Encoder).kSlice
				fn.fd = (*Decoder).kSlice
			case reflect.Array:
				fi.seq = seqTypeArray
				fn.fe = (*Encoder).kSlice
				fi.addrF = false
				fi.addrD = false
				rt2 := reflect.SliceOf(ti.elem)
				fn.fd = func(d *Decoder, xf *codecFnInfo, xrv reflect.Value) ***REMOVED***
					d.cfer().get(rt2, true, false).fd(d, xf, xrv.Slice(0, xrv.Len()))
				***REMOVED***
				// fn.fd = (*Decoder).kArray
			case reflect.Struct:
				if ti.anyOmitEmpty ***REMOVED***
					fn.fe = (*Encoder).kStruct
				***REMOVED*** else ***REMOVED***
					fn.fe = (*Encoder).kStructNoOmitempty
				***REMOVED***
				fn.fd = (*Decoder).kStruct
			case reflect.Map:
				fn.fe = (*Encoder).kMap
				fn.fd = (*Decoder).kMap
			case reflect.Interface:
				// encode: reflect.Interface are handled already by preEncodeValue
				fn.fd = (*Decoder).kInterface
				fn.fe = (*Encoder).kErr
			default:
				// reflect.Ptr and reflect.Interface are handled already by preEncodeValue
				fn.fe = (*Encoder).kErr
				fn.fd = (*Decoder).kErr
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

type codecFnPooler struct ***REMOVED***
	cf  *codecFner
	cfp *sync.Pool
	hh  Handle
***REMOVED***

func (d *codecFnPooler) cfer() *codecFner ***REMOVED***
	if d.cf == nil ***REMOVED***
		var v interface***REMOVED******REMOVED***
		d.cfp, v = pool.codecFner()
		d.cf = v.(*codecFner)
		d.cf.reset(d.hh)
	***REMOVED***
	return d.cf
***REMOVED***

func (d *codecFnPooler) alwaysAtEnd() ***REMOVED***
	if d.cf != nil ***REMOVED***
		d.cfp.Put(d.cf)
		d.cf, d.cfp = nil, nil
	***REMOVED***
***REMOVED***

// ----

// these "checkOverflow" functions must be inlinable, and not call anybody.
// Overflow means that the value cannot be represented without wrapping/overflow.
// Overflow=false does not mean that the value can be represented without losing precision
// (especially for floating point).

type checkOverflow struct***REMOVED******REMOVED***

// func (checkOverflow) Float16(f float64) (overflow bool) ***REMOVED***
// 	panicv.errorf("unimplemented")
// 	if f < 0 ***REMOVED***
// 		f = -f
// 	***REMOVED***
// 	return math.MaxFloat32 < f && f <= math.MaxFloat64
// ***REMOVED***

func (checkOverflow) Float32(v float64) (overflow bool) ***REMOVED***
	if v < 0 ***REMOVED***
		v = -v
	***REMOVED***
	return math.MaxFloat32 < v && v <= math.MaxFloat64
***REMOVED***
func (checkOverflow) Uint(v uint64, bitsize uint8) (overflow bool) ***REMOVED***
	if bitsize == 0 || bitsize >= 64 || v == 0 ***REMOVED***
		return
	***REMOVED***
	if trunc := (v << (64 - bitsize)) >> (64 - bitsize); v != trunc ***REMOVED***
		overflow = true
	***REMOVED***
	return
***REMOVED***
func (checkOverflow) Int(v int64, bitsize uint8) (overflow bool) ***REMOVED***
	if bitsize == 0 || bitsize >= 64 || v == 0 ***REMOVED***
		return
	***REMOVED***
	if trunc := (v << (64 - bitsize)) >> (64 - bitsize); v != trunc ***REMOVED***
		overflow = true
	***REMOVED***
	return
***REMOVED***
func (checkOverflow) SignedInt(v uint64) (overflow bool) ***REMOVED***
	//e.g. -127 to 128 for int8
	pos := (v >> 63) == 0
	ui2 := v & 0x7fffffffffffffff
	if pos ***REMOVED***
		if ui2 > math.MaxInt64 ***REMOVED***
			overflow = true
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if ui2 > math.MaxInt64-1 ***REMOVED***
			overflow = true
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func (x checkOverflow) Float32V(v float64) float64 ***REMOVED***
	if x.Float32(v) ***REMOVED***
		panicv.errorf("float32 overflow: %v", v)
	***REMOVED***
	return v
***REMOVED***
func (x checkOverflow) UintV(v uint64, bitsize uint8) uint64 ***REMOVED***
	if x.Uint(v, bitsize) ***REMOVED***
		panicv.errorf("uint64 overflow: %v", v)
	***REMOVED***
	return v
***REMOVED***
func (x checkOverflow) IntV(v int64, bitsize uint8) int64 ***REMOVED***
	if x.Int(v, bitsize) ***REMOVED***
		panicv.errorf("int64 overflow: %v", v)
	***REMOVED***
	return v
***REMOVED***
func (x checkOverflow) SignedIntV(v uint64) int64 ***REMOVED***
	if x.SignedInt(v) ***REMOVED***
		panicv.errorf("uint64 to int64 overflow: %v", v)
	***REMOVED***
	return int64(v)
***REMOVED***

// ------------------ SORT -----------------

func isNaN(f float64) bool ***REMOVED*** return f != f ***REMOVED***

// -----------------------

type ioFlusher interface ***REMOVED***
	Flush() error
***REMOVED***

type ioPeeker interface ***REMOVED***
	Peek(int) ([]byte, error)
***REMOVED***

type ioBuffered interface ***REMOVED***
	Buffered() int
***REMOVED***

// -----------------------

type intSlice []int64
type uintSlice []uint64

// type uintptrSlice []uintptr
type floatSlice []float64
type boolSlice []bool
type stringSlice []string

// type bytesSlice [][]byte

func (p intSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p intSlice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p intSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p uintSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p uintSlice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p uintSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

// func (p uintptrSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
// func (p uintptrSlice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
// func (p uintptrSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p floatSlice) Len() int ***REMOVED*** return len(p) ***REMOVED***
func (p floatSlice) Less(i, j int) bool ***REMOVED***
	return p[i] < p[j] || isNaN(p[i]) && !isNaN(p[j])
***REMOVED***
func (p floatSlice) Swap(i, j int) ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

func (p stringSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p stringSlice) Less(i, j int) bool ***REMOVED*** return p[i] < p[j] ***REMOVED***
func (p stringSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

// func (p bytesSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
// func (p bytesSlice) Less(i, j int) bool ***REMOVED*** return bytes.Compare(p[i], p[j]) == -1 ***REMOVED***
// func (p bytesSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

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
type timeRv struct ***REMOVED***
	v time.Time
	r reflect.Value
***REMOVED***
type timeRvSlice []timeRv

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

func (p timeRvSlice) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p timeRvSlice) Less(i, j int) bool ***REMOVED*** return p[i].v.Before(p[j].v) ***REMOVED***
func (p timeRvSlice) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

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

// ------

// bitset types are better than [256]bool, because they permit the whole
// bitset array being on a single cache line and use less memory.

// given x > 0 and n > 0 and x is exactly 2^n, then pos/x === pos>>n AND pos%x === pos&(x-1).
// consequently, pos/32 === pos>>5, pos/16 === pos>>4, pos/8 === pos>>3, pos%8 == pos&7

type bitset256 [32]byte

func (x *bitset256) isset(pos byte) bool ***REMOVED***
	return x[pos>>3]&(1<<(pos&7)) != 0
***REMOVED***
func (x *bitset256) issetv(pos byte) byte ***REMOVED***
	return x[pos>>3] & (1 << (pos & 7))
***REMOVED***
func (x *bitset256) set(pos byte) ***REMOVED***
	x[pos>>3] |= (1 << (pos & 7))
***REMOVED***

// func (x *bitset256) unset(pos byte) ***REMOVED***
// 	x[pos>>3] &^= (1 << (pos & 7))
// ***REMOVED***

type bitset128 [16]byte

func (x *bitset128) isset(pos byte) bool ***REMOVED***
	return x[pos>>3]&(1<<(pos&7)) != 0
***REMOVED***
func (x *bitset128) set(pos byte) ***REMOVED***
	x[pos>>3] |= (1 << (pos & 7))
***REMOVED***

// func (x *bitset128) unset(pos byte) ***REMOVED***
// 	x[pos>>3] &^= (1 << (pos & 7))
// ***REMOVED***

type bitset32 [4]byte

func (x *bitset32) isset(pos byte) bool ***REMOVED***
	return x[pos>>3]&(1<<(pos&7)) != 0
***REMOVED***
func (x *bitset32) set(pos byte) ***REMOVED***
	x[pos>>3] |= (1 << (pos & 7))
***REMOVED***

// func (x *bitset32) unset(pos byte) ***REMOVED***
// 	x[pos>>3] &^= (1 << (pos & 7))
// ***REMOVED***

// type bit2set256 [64]byte

// func (x *bit2set256) set(pos byte, v1, v2 bool) ***REMOVED***
// 	var pos2 uint8 = (pos & 3) << 1 // returning 0, 2, 4 or 6
// 	if v1 ***REMOVED***
// 		x[pos>>2] |= 1 << (pos2 + 1)
// 	***REMOVED***
// 	if v2 ***REMOVED***
// 		x[pos>>2] |= 1 << pos2
// 	***REMOVED***
// ***REMOVED***
// func (x *bit2set256) get(pos byte) uint8 ***REMOVED***
// 	var pos2 uint8 = (pos & 3) << 1     // returning 0, 2, 4 or 6
// 	return x[pos>>2] << (6 - pos2) >> 6 // 11000000 -> 00000011
// ***REMOVED***

// ------------

type pooler struct ***REMOVED***
	dn                                          sync.Pool // for decNaked
	cfn                                         sync.Pool // for codecFner
	tiload                                      sync.Pool
	strRv8, strRv16, strRv32, strRv64, strRv128 sync.Pool // for stringRV
***REMOVED***

func (p *pooler) init() ***REMOVED***
	p.strRv8.New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([8]stringRv) ***REMOVED***
	p.strRv16.New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([16]stringRv) ***REMOVED***
	p.strRv32.New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([32]stringRv) ***REMOVED***
	p.strRv64.New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([64]stringRv) ***REMOVED***
	p.strRv128.New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new([128]stringRv) ***REMOVED***
	p.dn.New = func() interface***REMOVED******REMOVED*** ***REMOVED*** x := new(decNaked); x.init(); return x ***REMOVED***
	p.tiload.New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(typeInfoLoadArray) ***REMOVED***
	p.cfn.New = func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(codecFner) ***REMOVED***
***REMOVED***

func (p *pooler) stringRv8() (sp *sync.Pool, v interface***REMOVED******REMOVED***) ***REMOVED***
	return &p.strRv8, p.strRv8.Get()
***REMOVED***
func (p *pooler) stringRv16() (sp *sync.Pool, v interface***REMOVED******REMOVED***) ***REMOVED***
	return &p.strRv16, p.strRv16.Get()
***REMOVED***
func (p *pooler) stringRv32() (sp *sync.Pool, v interface***REMOVED******REMOVED***) ***REMOVED***
	return &p.strRv32, p.strRv32.Get()
***REMOVED***
func (p *pooler) stringRv64() (sp *sync.Pool, v interface***REMOVED******REMOVED***) ***REMOVED***
	return &p.strRv64, p.strRv64.Get()
***REMOVED***
func (p *pooler) stringRv128() (sp *sync.Pool, v interface***REMOVED******REMOVED***) ***REMOVED***
	return &p.strRv128, p.strRv128.Get()
***REMOVED***
func (p *pooler) decNaked() (sp *sync.Pool, v interface***REMOVED******REMOVED***) ***REMOVED***
	return &p.dn, p.dn.Get()
***REMOVED***
func (p *pooler) codecFner() (sp *sync.Pool, v interface***REMOVED******REMOVED***) ***REMOVED***
	return &p.cfn, p.cfn.Get()
***REMOVED***
func (p *pooler) tiLoad() (sp *sync.Pool, v interface***REMOVED******REMOVED***) ***REMOVED***
	return &p.tiload, p.tiload.Get()
***REMOVED***

// func (p *pooler) decNaked() (v *decNaked, f func(*decNaked) ) ***REMOVED***
// 	sp := &(p.dn)
// 	vv := sp.Get()
// 	return vv.(*decNaked), func(x *decNaked) ***REMOVED*** sp.Put(vv) ***REMOVED***
// ***REMOVED***
// func (p *pooler) decNakedGet() (v interface***REMOVED******REMOVED***) ***REMOVED***
// 	return p.dn.Get()
// ***REMOVED***
// func (p *pooler) codecFnerGet() (v interface***REMOVED******REMOVED***) ***REMOVED***
// 	return p.cfn.Get()
// ***REMOVED***
// func (p *pooler) tiLoadGet() (v interface***REMOVED******REMOVED***) ***REMOVED***
// 	return p.tiload.Get()
// ***REMOVED***
// func (p *pooler) decNakedPut(v interface***REMOVED******REMOVED***) ***REMOVED***
// 	p.dn.Put(v)
// ***REMOVED***
// func (p *pooler) codecFnerPut(v interface***REMOVED******REMOVED***) ***REMOVED***
// 	p.cfn.Put(v)
// ***REMOVED***
// func (p *pooler) tiLoadPut(v interface***REMOVED******REMOVED***) ***REMOVED***
// 	p.tiload.Put(v)
// ***REMOVED***

type panicHdl struct***REMOVED******REMOVED***

func (panicHdl) errorv(err error) ***REMOVED***
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (panicHdl) errorstr(message string) ***REMOVED***
	if message != "" ***REMOVED***
		panic(message)
	***REMOVED***
***REMOVED***

func (panicHdl) errorf(format string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	if format != "" ***REMOVED***
		if len(params) == 0 ***REMOVED***
			panic(format)
		***REMOVED*** else ***REMOVED***
			panic(fmt.Sprintf(format, params...))
		***REMOVED***
	***REMOVED***
***REMOVED***

type errstrDecorator interface ***REMOVED***
	wrapErrstr(interface***REMOVED******REMOVED***, *error)
***REMOVED***

type errstrDecoratorDef struct***REMOVED******REMOVED***

func (errstrDecoratorDef) wrapErrstr(v interface***REMOVED******REMOVED***, e *error) ***REMOVED*** *e = fmt.Errorf("%v", v) ***REMOVED***

type must struct***REMOVED******REMOVED***

func (must) String(s string, err error) string ***REMOVED***
	if err != nil ***REMOVED***
		panicv.errorv(err)
	***REMOVED***
	return s
***REMOVED***
func (must) Int(s int64, err error) int64 ***REMOVED***
	if err != nil ***REMOVED***
		panicv.errorv(err)
	***REMOVED***
	return s
***REMOVED***
func (must) Uint(s uint64, err error) uint64 ***REMOVED***
	if err != nil ***REMOVED***
		panicv.errorv(err)
	***REMOVED***
	return s
***REMOVED***
func (must) Float(s float64, err error) float64 ***REMOVED***
	if err != nil ***REMOVED***
		panicv.errorv(err)
	***REMOVED***
	return s
***REMOVED***

// xdebugf prints the message in red on the terminal.
// Use it in place of fmt.Printf (which it calls internally)
func xdebugf(pattern string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	var delim string
	if len(pattern) > 0 && pattern[len(pattern)-1] != '\n' ***REMOVED***
		delim = "\n"
	***REMOVED***
	fmt.Printf("\033[1;31m"+pattern+delim+"\033[0m", args...)
***REMOVED***

// func isImmutableKind(k reflect.Kind) (v bool) ***REMOVED***
// 	return false ||
// 		k == reflect.Int ||
// 		k == reflect.Int8 ||
// 		k == reflect.Int16 ||
// 		k == reflect.Int32 ||
// 		k == reflect.Int64 ||
// 		k == reflect.Uint ||
// 		k == reflect.Uint8 ||
// 		k == reflect.Uint16 ||
// 		k == reflect.Uint32 ||
// 		k == reflect.Uint64 ||
// 		k == reflect.Uintptr ||
// 		k == reflect.Float32 ||
// 		k == reflect.Float64 ||
// 		k == reflect.Bool ||
// 		k == reflect.String
// ***REMOVED***

// func timeLocUTCName(tzint int16) string ***REMOVED***
// 	if tzint == 0 ***REMOVED***
// 		return "UTC"
// 	***REMOVED***
// 	var tzname = []byte("UTC+00:00")
// 	//tzname := fmt.Sprintf("UTC%s%02d:%02d", tzsign, tz/60, tz%60) //perf issue using Sprintf. inline below.
// 	//tzhr, tzmin := tz/60, tz%60 //faster if u convert to int first
// 	var tzhr, tzmin int16
// 	if tzint < 0 ***REMOVED***
// 		tzname[3] = '-' // (TODO: verify. this works here)
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
