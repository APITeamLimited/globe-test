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
// We considered storing the error is En|Decoder.
//   - once it has its err field set, it cannot be used again.
//   - panicing will be optional, controlled by const flag.
//   - code should always check error first and return early.
//
// We eventually decided against it as it makes the code clumsier to always
// check for these error conditions.
//
// ------------------------------------------
// We use sync.Pool only for the aid of long-lived objects shared across multiple goroutines.
// Encoder, Decoder, enc|decDriver, reader|writer, etc do not fall into this bucket.
//
// Also, GC is much better now, eliminating some of the reasons to use a shared pool structure.
// Instead, the short-lived objects use free-lists that live as long as the object exists.
//
// ------------------------------------------
// Performance is affected by the following:
//    - Bounds Checking
//    - Inlining
//    - Pointer chasing
// This package tries hard to manage the performance impact of these.
//
// ------------------------------------------
// To alleviate performance due to pointer-chasing:
//    - Prefer non-pointer values in a struct field
//    - Refer to these directly within helper classes
//      e.g. json.go refers directly to d.d.decRd
//
// We made the changes to embed En/Decoder in en/decDriver,
// but we had to explicitly reference the fields as opposed to using a function
// to get the better performance that we were looking for.
// For example, we explicitly call d.d.decRd.fn() instead of d.d.r().fn().
//
// ------------------------------------------
// Bounds Checking
//    - Allow bytesDecReader to incur "bounds check error", and
//      recover that as an io.EOF.
//      This allows the bounds check branch to always be taken by the branch predictor,
//      giving better performance (in theory), while ensuring that the code is shorter.
//
// ------------------------------------------
// Escape Analysis
//    - Prefer to return non-pointers if the value is used right away.
//      Newly allocated values returned as pointers will be heap-allocated as they escape.
//
// Prefer functions and methods that
//    - take no parameters and
//    - return no results and
//    - do not allocate.
// These are optimized by the runtime.
// For example, in json, we have dedicated functions for ReadMapElemKey, etc
// which do not delegate to readDelim, as readDelim takes a parameter.
// The difference in runtime was as much as 5%.

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
	"sync/atomic"
	"time"
)

const (
	// rvNLen is the length of the array for readn or writen calls
	rwNLen = 7

	// scratchByteArrayLen = 64
	// initCollectionCap   = 16 // 32 is defensive. 16 is preferred.

	// Support encoding.(Binary|Text)(Unm|M)arshaler.
	// This constant flag will enable or disable it.
	supportMarshalInterfaces = true

	// for debugging, set this to false, to catch panic traces.
	// Note that this will always cause rpc tests to fail, since they need io.EOF sent via panic.
	recoverPanicToErr = true

	// arrayCacheLen is the length of the cache used in encoder or decoder for
	// allowing zero-alloc initialization.
	// arrayCacheLen = 8

	// size of the cacheline: defaulting to value for archs: amd64, arm64, 386
	// should use "runtime/internal/sys".CacheLineSize, but that is not exposed.
	cacheLineSize = 64

	wordSizeBits = 32 << (^uint(0) >> 63) // strconv.IntSize
	wordSize     = wordSizeBits / 8

	// so structFieldInfo fits into 8 bytes
	maxLevelsEmbedding = 14

	// xdebug controls whether xdebugf prints any output
	xdebug = true
)

var (
	oneByteArr    [1]byte
	zeroByteSlice = oneByteArr[:0:0]

	codecgen bool

	panicv panicHdl

	refBitset    bitset32
	isnilBitset  bitset32
	scalarBitset bitset32
)

var (
	errMapTypeNotMapKind     = errors.New("MapType MUST be of Map Kind")
	errSliceTypeNotSliceKind = errors.New("SliceType MUST be of Slice Kind")
)

var pool4tiload = sync.Pool***REMOVED***New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(typeInfoLoadArray) ***REMOVED******REMOVED***

func init() ***REMOVED***
	refBitset = refBitset.
		set(byte(reflect.Map)).
		set(byte(reflect.Ptr)).
		set(byte(reflect.Func)).
		set(byte(reflect.Chan)).
		set(byte(reflect.UnsafePointer))

	isnilBitset = isnilBitset.
		set(byte(reflect.Map)).
		set(byte(reflect.Ptr)).
		set(byte(reflect.Func)).
		set(byte(reflect.Chan)).
		set(byte(reflect.UnsafePointer)).
		set(byte(reflect.Interface)).
		set(byte(reflect.Slice))

	scalarBitset = scalarBitset.
		set(byte(reflect.Bool)).
		set(byte(reflect.Int)).
		set(byte(reflect.Int8)).
		set(byte(reflect.Int16)).
		set(byte(reflect.Int32)).
		set(byte(reflect.Int64)).
		set(byte(reflect.Uint)).
		set(byte(reflect.Uint8)).
		set(byte(reflect.Uint16)).
		set(byte(reflect.Uint32)).
		set(byte(reflect.Uint64)).
		set(byte(reflect.Uintptr)).
		set(byte(reflect.Float32)).
		set(byte(reflect.Float64)).
		set(byte(reflect.Complex64)).
		set(byte(reflect.Complex128)).
		set(byte(reflect.String))

***REMOVED***

type handleFlag uint8

const (
	initedHandleFlag handleFlag = 1 << iota
	binaryHandleFlag
	jsonHandleFlag
)

type clsErr struct ***REMOVED***
	closed    bool  // is it closed?
	errClosed error // error on closing
***REMOVED***

type charEncoding uint8

const (
	_ charEncoding = iota // make 0 unset
	cUTF8
	cUTF16LE
	cUTF16BE
	cUTF32LE
	cUTF32BE
	// Deprecated: not a true char encoding value
	cRAW charEncoding = 255
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

	containerMapStart
	containerMapKey
	containerMapValue
	containerMapEnd
	containerArrayStart
	containerArrayElem
	containerArrayEnd
)

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

// typeInfoLoad is a transient object used while loading up a typeInfo.
type typeInfoLoad struct ***REMOVED***
	etypes []uintptr
	sfis   []structFieldInfo
***REMOVED***

// typeInfoLoadArray is a cache object used to efficiently load up a typeInfo without
// much allocation.
type typeInfoLoadArray struct ***REMOVED***
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

type codecError struct ***REMOVED***
	name string
	err  interface***REMOVED******REMOVED***
***REMOVED***

func (e codecError) Cause() error ***REMOVED***
	switch xerr := e.err.(type) ***REMOVED***
	case nil:
		return nil
	case error:
		return xerr
	case string:
		return errors.New(xerr)
	case fmt.Stringer:
		return errors.New(xerr.String())
	default:
		return fmt.Errorf("%v", e.err)
	***REMOVED***
***REMOVED***

func (e codecError) Error() string ***REMOVED***
	return fmt.Sprintf("%s error: %v", e.name, e.err)
***REMOVED***

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

	selferTyp         = reflect.TypeOf((*Selfer)(nil)).Elem()
	missingFielderTyp = reflect.TypeOf((*MissingFielder)(nil)).Elem()
	iszeroTyp         = reflect.TypeOf((*isZeroer)(nil)).Elem()

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

	// bsAll0x00 = []byte***REMOVED***0, 0, 0, 0, 0, 0, 0, 0***REMOVED***
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

// SelfExt is a sentinel extension signifying that types
// registered with it SHOULD be encoded and decoded
// based on the native mode of the format.
//
// This allows users to define a tag for an extension,
// but signify that the types should be encoded/decoded as the native encoding.
// This way, users need not also define how to encode or decode the extension.
var SelfExt = &extFailWrapper***REMOVED******REMOVED***

// Selfer defines methods by which a value can encode or decode itself.
//
// Any type which implements Selfer will be able to encode or decode itself.
// Consequently, during (en|de)code, this takes precedence over
// (text|binary)(M|Unm)arshal or extension support.
//
// By definition, it is not allowed for a Selfer to directly call Encode or Decode on itself.
// If that is done, Encode/Decode will rightfully fail with a Stack Overflow style error.
// For example, the snippet below will cause such an error.
//     type testSelferRecur struct***REMOVED******REMOVED***
//     func (s *testSelferRecur) CodecEncodeSelf(e *Encoder) ***REMOVED*** e.MustEncode(s) ***REMOVED***
//     func (s *testSelferRecur) CodecDecodeSelf(d *Decoder) ***REMOVED*** d.MustDecode(s) ***REMOVED***
//
// Note: *the first set of bytes of any value MUST NOT represent nil in the format*.
// This is because, during each decode, we first check the the next set of bytes
// represent nil, and if so, we just set the value to nil.
type Selfer interface ***REMOVED***
	CodecEncodeSelf(*Encoder)
	CodecDecodeSelf(*Decoder)
***REMOVED***

// MissingFielder defines the interface allowing structs to internally decode or encode
// values which do not map to struct fields.
//
// We expect that this interface is bound to a pointer type (so the mutation function works).
//
// A use-case is if a version of a type unexports a field, but you want compatibility between
// both versions during encoding and decoding.
//
// Note that the interface is completely ignored during codecgen.
type MissingFielder interface ***REMOVED***
	// CodecMissingField is called to set a missing field and value pair.
	//
	// It returns true if the missing field was set on the struct.
	CodecMissingField(field []byte, value interface***REMOVED******REMOVED***) bool

	// CodecMissingFields returns the set of fields which are not struct fields
	CodecMissingFields() map[string]interface***REMOVED******REMOVED***
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

	rtidFns      atomicRtidFnSlice
	rtidFnsNoExt atomicRtidFnSlice

	// ---- cache line

	DecodeOptions

	// ---- cache line

	EncodeOptions

	intf2impls

	mu     sync.Mutex
	inited uint32 // holds if inited, and also handle flags (binary encoding, json handler, etc)

	RPCOptions

	// TimeNotBuiltin configures whether time.Time should be treated as a builtin type.
	//
	// All Handlers should know how to encode/decode time.Time as part of the core
	// format specification, or as a standard extension defined by the format.
	//
	// However, users can elect to handle time.Time as a custom extension, or via the
	// standard library's encoding.Binary(M|Unm)arshaler or Text(M|Unm)arshaler interface.
	// To elect this behavior, users can set TimeNotBuiltin=true.
	//
	// Note: Setting TimeNotBuiltin=true can be used to enable the legacy behavior
	// (for Cbor and Msgpack), where time.Time was not a builtin supported type.
	//
	// Note: DO NOT CHANGE AFTER FIRST USE.
	//
	// Once a Handle has been used, do not modify this option.
	// It will lead to unexpected behaviour during encoding and decoding.
	TimeNotBuiltin bool

	// ExplicitRelease configures whether Release() is implicitly called after an encode or
	// decode call.
	//
	// If you will hold onto an Encoder or Decoder for re-use, by calling Reset(...)
	// on it or calling (Must)Encode repeatedly into a given []byte or io.Writer,
	// then you do not want it to be implicitly closed after each Encode/Decode call.
	// Doing so will unnecessarily return resources to the shared pool, only for you to
	// grab them right after again to do another Encode/Decode call.
	//
	// Instead, you configure ExplicitRelease=true, and you explicitly call Release() when
	// you are truly done.
	//
	// As an alternative, you can explicitly set a finalizer - so its resources
	// are returned to the shared pool before it is garbage-collected. Do it as below:
	//    runtime.SetFinalizer(e, (*Encoder).Release)
	//    runtime.SetFinalizer(d, (*Decoder).Release)
	//
	// Deprecated: This is not longer used as pools are only used for long-lived objects
	// which are shared across goroutines.
	// Setting this value has no effect. It is maintained for backward compatibility.
	ExplicitRelease bool

	// ---- cache line
***REMOVED***

// basicHandle returns an initialized BasicHandle from the Handle.
func basicHandle(hh Handle) (x *BasicHandle) ***REMOVED***
	x = hh.getBasicHandle()
	// ** We need to simulate once.Do, to ensure no data race within the block.
	// ** Consequently, below would not work.
	// if atomic.CompareAndSwapUint32(&x.inited, 0, 1) ***REMOVED***
	// 	x.be = hh.isBinary()
	// 	_, x.js = hh.(*JsonHandle)
	// 	x.n = hh.Name()[0]
	// ***REMOVED***

	// simulate once.Do using our own stored flag and mutex as a CompareAndSwap
	// is not sufficient, since a race condition can occur within init(Handle) function.
	// init is made noinline, so that this function can be inlined by its caller.
	if atomic.LoadUint32(&x.inited) == 0 ***REMOVED***
		x.init(hh)
	***REMOVED***
	return
***REMOVED***

func (x *BasicHandle) isJs() bool ***REMOVED***
	return handleFlag(x.inited)&jsonHandleFlag != 0
***REMOVED***

func (x *BasicHandle) isBe() bool ***REMOVED***
	return handleFlag(x.inited)&binaryHandleFlag != 0
***REMOVED***

//go:noinline
func (x *BasicHandle) init(hh Handle) ***REMOVED***
	// make it uninlineable, as it is called at most once
	x.mu.Lock()
	if x.inited == 0 ***REMOVED***
		var f = initedHandleFlag
		if hh.isBinary() ***REMOVED***
			f |= binaryHandleFlag
		***REMOVED***
		if _, b := hh.(*JsonHandle); b ***REMOVED***
			f |= jsonHandleFlag
		***REMOVED***
		atomic.StoreUint32(&x.inited, uint32(f))
		// ensure MapType and SliceType are of correct type
		if x.MapType != nil && x.MapType.Kind() != reflect.Map ***REMOVED***
			panic(errMapTypeNotMapKind)
		***REMOVED***
		if x.SliceType != nil && x.SliceType.Kind() != reflect.Slice ***REMOVED***
			panic(errSliceTypeNotSliceKind)
		***REMOVED***
	***REMOVED***
	x.mu.Unlock()
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

func findFn(s []codecRtidFn, rtid uintptr) (i uint, fn *codecFn) ***REMOVED***
	// binary search. adapted from sort/search.go.
	// Note: we use goto (instead of for loop) so this can be inlined.

	// h, i, j := 0, 0, len(s)
	var h uint // var h, i uint
	var j = uint(len(s))
LOOP:
	if i < j ***REMOVED***
		h = i + (j-i)/2
		if s[h].rtid < rtid ***REMOVED***
			i = h + 1
		***REMOVED*** else ***REMOVED***
			j = h
		***REMOVED***
		goto LOOP
	***REMOVED***
	if i < uint(len(s)) && s[i].rtid == rtid ***REMOVED***
		fn = s[i].fn
	***REMOVED***
	return
***REMOVED***

func (x *BasicHandle) fn(rt reflect.Type) (fn *codecFn) ***REMOVED***
	return x.fnVia(rt, &x.rtidFns, true)
***REMOVED***

func (x *BasicHandle) fnNoExt(rt reflect.Type) (fn *codecFn) ***REMOVED***
	return x.fnVia(rt, &x.rtidFnsNoExt, false)
***REMOVED***

func (x *BasicHandle) fnVia(rt reflect.Type, fs *atomicRtidFnSlice, checkExt bool) (fn *codecFn) ***REMOVED***
	rtid := rt2id(rt)
	sp := fs.load()
	if sp != nil ***REMOVED***
		if _, fn = findFn(sp, rtid); fn != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	fn = x.fnLoad(rt, rtid, checkExt)
	x.mu.Lock()
	var sp2 []codecRtidFn
	sp = fs.load()
	if sp == nil ***REMOVED***
		sp2 = []codecRtidFn***REMOVED******REMOVED***rtid, fn***REMOVED******REMOVED***
		fs.store(sp2)
	***REMOVED*** else ***REMOVED***
		idx, fn2 := findFn(sp, rtid)
		if fn2 == nil ***REMOVED***
			sp2 = make([]codecRtidFn, len(sp)+1)
			copy(sp2, sp[:idx])
			copy(sp2[idx+1:], sp[idx:])
			sp2[idx] = codecRtidFn***REMOVED***rtid, fn***REMOVED***
			fs.store(sp2)
		***REMOVED***
	***REMOVED***
	x.mu.Unlock()
	return
***REMOVED***

func (x *BasicHandle) fnLoad(rt reflect.Type, rtid uintptr, checkExt bool) (fn *codecFn) ***REMOVED***
	fn = new(codecFn)
	fi := &(fn.i)
	ti := x.getTypeInfo(rtid, rt)
	fi.ti = ti

	rk := reflect.Kind(ti.kind)

	// anything can be an extension except the built-in ones: time, raw and rawext

	if rtid == timeTypId && !x.TimeNotBuiltin ***REMOVED***
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
	***REMOVED*** else if xfFn := x.getExt(rtid, checkExt); xfFn != nil ***REMOVED***
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*Encoder).ext
		fn.fd = (*Decoder).ext
		fi.addrF = true
		fi.addrD = true
		if rk == reflect.Struct || rk == reflect.Array ***REMOVED***
			fi.addrE = true
		***REMOVED***
	***REMOVED*** else if ti.isFlag(tiflagSelfer) || ti.isFlag(tiflagSelferPtr) ***REMOVED***
		fn.fe = (*Encoder).selferMarshal
		fn.fd = (*Decoder).selferUnmarshal
		fi.addrF = true
		fi.addrD = ti.isFlag(tiflagSelferPtr)
		fi.addrE = ti.isFlag(tiflagSelferPtr)
	***REMOVED*** else if supportMarshalInterfaces && x.isBe() &&
		(ti.isFlag(tiflagBinaryMarshaler) || ti.isFlag(tiflagBinaryMarshalerPtr)) &&
		(ti.isFlag(tiflagBinaryUnmarshaler) || ti.isFlag(tiflagBinaryUnmarshalerPtr)) ***REMOVED***
		fn.fe = (*Encoder).binaryMarshal
		fn.fd = (*Decoder).binaryUnmarshal
		fi.addrF = true
		fi.addrD = ti.isFlag(tiflagBinaryUnmarshalerPtr)
		fi.addrE = ti.isFlag(tiflagBinaryMarshalerPtr)
	***REMOVED*** else if supportMarshalInterfaces && !x.isBe() && x.isJs() &&
		(ti.isFlag(tiflagJsonMarshaler) || ti.isFlag(tiflagJsonMarshalerPtr)) &&
		(ti.isFlag(tiflagJsonUnmarshaler) || ti.isFlag(tiflagJsonUnmarshalerPtr)) ***REMOVED***
		//If JSON, we should check JSONMarshal before textMarshal
		fn.fe = (*Encoder).jsonMarshal
		fn.fd = (*Decoder).jsonUnmarshal
		fi.addrF = true
		fi.addrD = ti.isFlag(tiflagJsonUnmarshalerPtr)
		fi.addrE = ti.isFlag(tiflagJsonMarshalerPtr)
	***REMOVED*** else if supportMarshalInterfaces && !x.isBe() &&
		(ti.isFlag(tiflagTextMarshaler) || ti.isFlag(tiflagTextMarshalerPtr)) &&
		(ti.isFlag(tiflagTextUnmarshaler) || ti.isFlag(tiflagTextUnmarshalerPtr)) ***REMOVED***
		fn.fe = (*Encoder).textMarshal
		fn.fd = (*Decoder).textUnmarshal
		fi.addrF = true
		fi.addrD = ti.isFlag(tiflagTextUnmarshalerPtr)
		fi.addrE = ti.isFlag(tiflagTextMarshalerPtr)
	***REMOVED*** else ***REMOVED***
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice) ***REMOVED***
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
						xfnf(e, xf, rvConvert(xrv, xrt))
					***REMOVED***
					fi.addrD = true
					fi.addrF = false // meaning it can be an address(ptr) or a value
					xfnf2 := fastpathAV[idx].decfn
					xptr2rt := reflect.PtrTo(xrt)
					fn.fd = func(d *Decoder, xf *codecFnInfo, xrv reflect.Value) ***REMOVED***
						if xrv.Kind() == reflect.Ptr ***REMOVED***
							xfnf2(d, xf, rvConvert(xrv, xptr2rt))
						***REMOVED*** else ***REMOVED***
							xfnf2(d, xf, rvConvert(xrv, xrt))
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
				// Do not use different functions based on StringToRaw option,
				// as that will statically set the function for a string type,
				// and if the Handle is modified thereafter, behaviour is non-deterministic.
				// i.e. DO NOT DO:
				//   if x.StringToRaw ***REMOVED***
				//   	fn.fe = (*Encoder).kStringToRaw
				//   ***REMOVED*** else ***REMOVED***
				//   	fn.fe = (*Encoder).kStringEnc
				//   ***REMOVED***

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
				fn.fe = (*Encoder).kChan
				fn.fd = (*Decoder).kSliceForChan
			case reflect.Slice:
				fi.seq = seqTypeSlice
				fn.fe = (*Encoder).kSlice
				fn.fd = (*Decoder).kSlice
			case reflect.Array:
				fi.seq = seqTypeArray
				fn.fe = (*Encoder).kArray
				fi.addrF = false
				fi.addrD = false
				rt2 := reflect.SliceOf(ti.elem)
				fn.fd = func(d *Decoder, xf *codecFnInfo, xrv reflect.Value) ***REMOVED***
					// call fnVia directly, so fn(...) is not recursive, and can be inlined
					d.h.fnVia(rt2, &x.rtidFns, true).fd(d, xf, rvGetSlice4Array(xrv, rt2))
				***REMOVED***
			case reflect.Struct:
				if ti.anyOmitEmpty ||
					ti.isFlag(tiflagMissingFielder) ||
					ti.isFlag(tiflagMissingFielderPtr) ***REMOVED***
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

// Handle defines a specific encoding format. It also stores any runtime state
// used during an Encoding or Decoding session e.g. stored state about Types, etc.
//
// Once a handle is configured, it can be shared across multiple Encoders and Decoders.
//
// Note that a Handle is NOT safe for concurrent modification.
//
// A Handle also should not be modified after it is configured and has
// been used at least once. This is because stored state may be out of sync with the
// new configuration, and a data race can occur when multiple goroutines access it.
// i.e. multiple Encoders or Decoders in different goroutines.
//
// Consequently, the typical usage model is that a Handle is pre-configured
// before first time use, and not modified while in use.
// Such a pre-configured Handle is safe for concurrent access.
type Handle interface ***REMOVED***
	Name() string
	// return the basic handle. It may not have been inited.
	// Prefer to use basicHandle() helper function that ensures it has been inited.
	getBasicHandle() *BasicHandle
	newEncDriver() encDriver
	newDecDriver() decDriver
	isBinary() bool
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
	bs, err := x.encFn(rv4i(v))
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return bs
***REMOVED***

func (x addExtWrapper) ReadExt(v interface***REMOVED******REMOVED***, bs []byte) ***REMOVED***
	if err := x.decFn(rv4i(v), bs); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (x addExtWrapper) ConvertExt(v interface***REMOVED******REMOVED***) interface***REMOVED******REMOVED*** ***REMOVED***
	return x.WriteExt(v)
***REMOVED***

func (x addExtWrapper) UpdateExt(dest interface***REMOVED******REMOVED***, v interface***REMOVED******REMOVED***) ***REMOVED***
	x.ReadExt(dest, v.([]byte))
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

type bytesExtWrapper struct ***REMOVED***
	interfaceExtFailer
	BytesExt
***REMOVED***

type interfaceExtWrapper struct ***REMOVED***
	bytesExtFailer
	InterfaceExt
***REMOVED***

type extFailWrapper struct ***REMOVED***
	bytesExtFailer
	interfaceExtFailer
***REMOVED***

type binaryEncodingType struct***REMOVED******REMOVED***

func (binaryEncodingType) isBinary() bool ***REMOVED*** return true ***REMOVED***

type textEncodingType struct***REMOVED******REMOVED***

func (textEncodingType) isBinary() bool ***REMOVED*** return false ***REMOVED***

// noBuiltInTypes is embedded into many types which do not support builtins
// e.g. msgpack, simple, cbor.

type noBuiltInTypes struct***REMOVED******REMOVED***

func (noBuiltInTypes) EncodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***
func (noBuiltInTypes) DecodeBuiltin(rt uintptr, v interface***REMOVED******REMOVED***) ***REMOVED******REMOVED***

// bigenHelper.
// Users must already slice the x completely, because we will not reslice.
type bigenHelper struct ***REMOVED***
	x []byte // must be correctly sliced to appropriate len. slicing is a cost.
	w *encWr
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
	// _       [1]uint64 // padding
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
	// We EXPECT *o is a pointer to a non-nil extHandle.

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
		// all natively supported type, so cannot have an extension.
		// However, we do not return an error for these, as we do not document that.
		// Instead, we silently treat as a no-op, and return.
		return
	***REMOVED***
	o2 := *o
	for i := range o2 ***REMOVED***
		v := &o2[i]
		if v.rtid == rtid ***REMOVED***
			v.tag, v.ext = tag, ext
			return
		***REMOVED***
	***REMOVED***
	rtidptr := rt2id(reflect.PtrTo(rt))
	*o = append(o2, extTypeTagFn***REMOVED***rtid, rtidptr, rt, tag, ext***REMOVED***) // , [1]uint64***REMOVED******REMOVED******REMOVED***)
	return
***REMOVED***

func (o extHandle) getExt(rtid uintptr, check bool) (v *extTypeTagFn) ***REMOVED***
	if !check ***REMOVED***
		return
	***REMOVED***
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
			vkind := v.impl.Kind()
			if vkind == reflect.Ptr ***REMOVED***
				return reflect.New(v.impl.Elem())
			***REMOVED***
			return rvZeroAddrK(v.impl, vkind)
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

	encNameAsciiAlphaNum bool // the encName only contains ascii alphabet and numbers
	structFieldInfoFlag
	// _ [1]byte // padding
***REMOVED***

// func (si *structFieldInfo) setToZeroValue(v reflect.Value) ***REMOVED***
// 	if v, valid := si.field(v, false); valid ***REMOVED***
// 		v.Set(reflect.Zero(v.Type()))
// 	***REMOVED***
// ***REMOVED***

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
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type sfiSortedByEncName []*structFieldInfo

func (p sfiSortedByEncName) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p sfiSortedByEncName) Less(i, j int) bool ***REMOVED*** return p[uint(i)].encName < p[uint(j)].encName ***REMOVED***
func (p sfiSortedByEncName) Swap(i, j int)      ***REMOVED*** p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] ***REMOVED***

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
		if rvIsNil(v) ***REMOVED***
			if !update ***REMOVED***
				return
			***REMOVED***
			rvSetDirect(v, reflect.New(v.Type().Elem()))
		***REMOVED***
		v = v.Elem()
	***REMOVED***
	return v, true
***REMOVED***

type tiflag uint32

const (
	_ tiflag = 1 << iota

	tiflagComparable

	tiflagIsZeroer
	tiflagIsZeroerPtr

	tiflagBinaryMarshaler
	tiflagBinaryMarshalerPtr

	tiflagBinaryUnmarshaler
	tiflagBinaryUnmarshalerPtr

	tiflagTextMarshaler
	tiflagTextMarshalerPtr

	tiflagTextUnmarshaler
	tiflagTextUnmarshalerPtr

	tiflagJsonMarshaler
	tiflagJsonMarshalerPtr

	tiflagJsonUnmarshaler
	tiflagJsonUnmarshalerPtr

	tiflagSelfer
	tiflagSelferPtr

	tiflagMissingFielder
	tiflagMissingFielderPtr
)

// typeInfo keeps static (non-changing readonly)information
// about each (non-ptr) type referenced in the encode/decode sequence.
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

	// rv0 is the zero value for the type.
	// It is mostly beneficial for all non-reference kinds
	// i.e. all but map/chan/func/ptr/unsafe.pointer
	// so beneficial for intXX, bool, slices, structs, etc
	rv0 reflect.Value

	elemsize uintptr

	// other flags, with individual bits representing if set.
	flags tiflag

	infoFieldOmitempty bool

	elemkind uint8
	_        [2]byte // padding
	// _ [1]uint64 // padding
***REMOVED***

func (ti *typeInfo) isFlag(f tiflag) bool ***REMOVED***
	return ti.flags&f != 0
***REMOVED***

func (ti *typeInfo) flag(when bool, f tiflag) *typeInfo ***REMOVED***
	if when ***REMOVED***
		ti.flags |= f
	***REMOVED***
	return ti
***REMOVED***

func (ti *typeInfo) indexForEncName(name []byte) (index int16) ***REMOVED***
	var sn []byte
	if len(name)+2 <= 32 ***REMOVED***
		var buf [32]byte // should not escape to heap
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
	_     uint64 // padding (cache-aligned)
	tags  []string
	_     uint64 // padding (cache-aligned)
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

func findTypeInfo(s []rtid2ti, rtid uintptr) (i uint, ti *typeInfo) ***REMOVED***
	// binary search. adapted from sort/search.go.
	// Note: we use goto (instead of for loop) so this can be inlined.

	// h, i, j := 0, 0, len(s)
	var h uint // var h, i uint
	var j = uint(len(s))
LOOP:
	if i < j ***REMOVED***
		h = i + (j-i)/2
		if s[h].rtid < rtid ***REMOVED***
			i = h + 1
		***REMOVED*** else ***REMOVED***
			j = h
		***REMOVED***
		goto LOOP
	***REMOVED***
	if i < uint(len(s)) && s[i].rtid == rtid ***REMOVED***
		ti = s[i].ti
	***REMOVED***
	return
***REMOVED***

func (x *TypeInfos) get(rtid uintptr, rt reflect.Type) (pti *typeInfo) ***REMOVED***
	sp := x.infos.load()
	if sp != nil ***REMOVED***
		_, pti = findTypeInfo(sp, rtid)
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
	ti := typeInfo***REMOVED***
		rt:      rt,
		rtid:    rtid,
		kind:    uint8(rk),
		pkgpath: rt.PkgPath(),
		keyType: valueTypeString, // default it - so it's never 0
	***REMOVED***
	ti.rv0 = reflect.Zero(rt)

	ti.numMeth = uint16(rt.NumMethod())

	var b1, b2 bool
	b1, b2 = implIntf(rt, binaryMarshalerTyp)
	ti.flag(b1, tiflagBinaryMarshaler).flag(b2, tiflagBinaryMarshalerPtr)
	b1, b2 = implIntf(rt, binaryUnmarshalerTyp)
	ti.flag(b1, tiflagBinaryUnmarshaler).flag(b2, tiflagBinaryUnmarshalerPtr)
	b1, b2 = implIntf(rt, textMarshalerTyp)
	ti.flag(b1, tiflagTextMarshaler).flag(b2, tiflagTextMarshalerPtr)
	b1, b2 = implIntf(rt, textUnmarshalerTyp)
	ti.flag(b1, tiflagTextUnmarshaler).flag(b2, tiflagTextUnmarshalerPtr)
	b1, b2 = implIntf(rt, jsonMarshalerTyp)
	ti.flag(b1, tiflagJsonMarshaler).flag(b2, tiflagJsonMarshalerPtr)
	b1, b2 = implIntf(rt, jsonUnmarshalerTyp)
	ti.flag(b1, tiflagJsonUnmarshaler).flag(b2, tiflagJsonUnmarshalerPtr)
	b1, b2 = implIntf(rt, selferTyp)
	ti.flag(b1, tiflagSelfer).flag(b2, tiflagSelferPtr)
	b1, b2 = implIntf(rt, missingFielderTyp)
	ti.flag(b1, tiflagMissingFielder).flag(b2, tiflagMissingFielderPtr)
	b1, b2 = implIntf(rt, iszeroTyp)
	ti.flag(b1, tiflagIsZeroer).flag(b2, tiflagIsZeroerPtr)
	b1 = rt.Comparable()
	ti.flag(b1, tiflagComparable)

	switch rk ***REMOVED***
	case reflect.Struct:
		var omitEmpty bool
		if f, ok := rt.FieldByName(structInfoFieldName); ok ***REMOVED***
			ti.toArray, omitEmpty, ti.keyType = parseStructInfo(x.structTag(f.Tag))
			ti.infoFieldOmitempty = omitEmpty
		***REMOVED*** else ***REMOVED***
			ti.keyType = valueTypeString
		***REMOVED***
		pp, pi := &pool4tiload, pool4tiload.Get() // pool.tiLoad()
		pv := pi.(*typeInfoLoadArray)
		pv.etypes[0] = ti.rtid
		// vv := typeInfoLoad***REMOVED***pv.fNames[:0], pv.encNames[:0], pv.etypes[:1], pv.sfis[:0]***REMOVED***
		vv := typeInfoLoad***REMOVED***pv.etypes[:1], pv.sfis[:0]***REMOVED***
		x.rget(rt, rtid, omitEmpty, nil, &vv)
		ti.sfiSrc, ti.sfiSort, ti.sfiNamesSort, ti.anyOmitEmpty = rgetResolveSFI(rt, vv.sfis, pv)
		pp.Put(pi)
	case reflect.Map:
		ti.elem = rt.Elem()
		ti.key = rt.Key()
	case reflect.Slice:
		ti.mbs, _ = implIntf(rt, mapBySliceTyp)
		ti.elem = rt.Elem()
		ti.elemsize = ti.elem.Size()
		ti.elemkind = uint8(ti.elem.Kind())
	case reflect.Chan:
		ti.elem = rt.Elem()
		ti.chandir = uint8(rt.ChanDir())
	case reflect.Array:
		ti.elem = rt.Elem()
		ti.elemsize = ti.elem.Size()
		ti.elemkind = uint8(ti.elem.Kind())
	case reflect.Ptr:
		ti.elem = rt.Elem()
	***REMOVED***

	x.mu.Lock()
	sp = x.infos.load()
	var sp2 []rtid2ti
	if sp == nil ***REMOVED***
		pti = &ti
		sp2 = []rtid2ti***REMOVED******REMOVED***rtid, pti***REMOVED******REMOVED***
		x.infos.store(sp2)
	***REMOVED*** else ***REMOVED***
		var idx uint
		idx, pti = findTypeInfo(sp, rtid)
		if pti == nil ***REMOVED***
			pti = &ti
			sp2 = make([]rtid2ti, len(sp)+1)
			copy(sp2, sp[:idx])
			copy(sp2[idx+1:], sp[idx:])
			sp2[idx] = rtid2ti***REMOVED***rtid, pti***REMOVED***
			x.infos.store(sp2)
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
		si.encNameAsciiAlphaNum = true
		for i := len(si.encName) - 1; i >= 0; i-- ***REMOVED*** // bounds-check elimination
			b := si.encName[i]
			if (b >= '0' && b <= '9') || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') ***REMOVED***
				continue
			***REMOVED***
			si.encNameAsciiAlphaNum = false
			break
		***REMOVED***
		si.fieldName = f.Name
		si.flagSet(structFieldInfoFlagReady)

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
	// Tried the following before settling on correct implementation:
	//   return 0xfe - (name[0] & 63)
	//   return 0xfe - (name[0] & 63) - uint8(len(name))
	//   return 0xfe - (name[0] & 63) - uint8(len(name)&63)
	//   return ((0xfe - (name[0] & 63)) & 0xf8) | (uint8(len(name) & 0x07))

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
		if len(xn)+2 > cap(sn) ***REMOVED***
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
			// one of them must be cleared (reset to nil),
			// and the index updated appropriately
			i2clear := ui                // index to be cleared
			if x[i].nis < x[index].nis ***REMOVED*** // this one is shallower
				// update the index to point to this later one.
				sa[j+len(sn)], sa[j+len(sn)+1] = byte(ui>>8), byte(ui)
				// clear the earlier one, as this later one is shallower.
				i2clear = index
			***REMOVED***
			if x[i2clear].ready() ***REMOVED***
				x[i2clear].flagClr(structFieldInfoFlagReady)
				n--
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
	if ti.isFlag(tiflagIsZeroerPtr) && v.CanAddr() ***REMOVED***
		return rv2i(v.Addr()).(isZeroer).IsZero()
	***REMOVED***
	if ti.isFlag(tiflagIsZeroer) ***REMOVED***
		return rv2i(v).(isZeroer).IsZero()
	***REMOVED***
	if ti.isFlag(tiflagComparable) ***REMOVED***
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

func panicToErr(h errDecorator, err *error) ***REMOVED***
	// Note: This method MUST be called directly from defer i.e. defer panicToErr ...
	// else it seems the recover is not fully handled
	if recoverPanicToErr ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			// fmt.Printf("panic'ing with: %v\n", x)
			// debug.PrintStack()
			panicValToErr(h, x, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func isSliceBoundsError(s string) bool ***REMOVED***
	return strings.Contains(s, "index out of range") ||
		strings.Contains(s, "slice bounds out of range")
***REMOVED***

func panicValToErr(h errDecorator, v interface***REMOVED******REMOVED***, err *error) ***REMOVED***
	d, dok := h.(*Decoder)
	switch xerr := v.(type) ***REMOVED***
	case nil:
	case error:
		switch xerr ***REMOVED***
		case nil:
		case io.EOF, io.ErrUnexpectedEOF, errEncoderNotInitialized, errDecoderNotInitialized:
			// treat as special (bubble up)
			*err = xerr
		default:
			if dok && d.bytes && isSliceBoundsError(xerr.Error()) ***REMOVED***
				*err = io.EOF
			***REMOVED*** else ***REMOVED***
				h.wrapErr(xerr, err)
			***REMOVED***
		***REMOVED***
	case string:
		if xerr != "" ***REMOVED***
			if dok && d.bytes && isSliceBoundsError(xerr) ***REMOVED***
				*err = io.EOF
			***REMOVED*** else ***REMOVED***
				h.wrapErr(xerr, err)
			***REMOVED***
		***REMOVED***
	case fmt.Stringer:
		if xerr != nil ***REMOVED***
			h.wrapErr(xerr, err)
		***REMOVED***
	default:
		h.wrapErr(v, err)
	***REMOVED***
***REMOVED***

func isImmutableKind(k reflect.Kind) (v bool) ***REMOVED***
	// return immutableKindsSet[k]
	// since we know reflect.Kind is in range 0..31, then use the k%32 == k constraint
	return immutableKindsSet[k%reflect.Kind(len(immutableKindsSet))] // bounds-check-elimination
***REMOVED***

func usableByteSlice(bs []byte, slen int) []byte ***REMOVED***
	if cap(bs) >= slen ***REMOVED***
		if bs == nil ***REMOVED***
			return []byte***REMOVED******REMOVED***
		***REMOVED***
		return bs[:slen]
	***REMOVED***
	return make([]byte, slen)
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
***REMOVED***

// codecFn encapsulates the captured variables and the encode function.
// This way, we only do some calculations one times, and pass to the
// code block that should be called (encapsulated in a function)
// instead of executing the checks every time.
type codecFn struct ***REMOVED***
	i  codecFnInfo
	fe func(*Encoder, *codecFnInfo, reflect.Value)
	fd func(*Decoder, *codecFnInfo, reflect.Value)
	_  [1]uint64 // padding (cache-aligned)
***REMOVED***

type codecRtidFn struct ***REMOVED***
	rtid uintptr
	fn   *codecFn
***REMOVED***

func makeExt(ext interface***REMOVED******REMOVED***) Ext ***REMOVED***
	if ext == nil ***REMOVED***
		return &extFailWrapper***REMOVED******REMOVED***
	***REMOVED***
	switch t := ext.(type) ***REMOVED***
	case nil:
		return &extFailWrapper***REMOVED******REMOVED***
	case Ext:
		return t
	case BytesExt:
		return &bytesExtWrapper***REMOVED***BytesExt: t***REMOVED***
	case InterfaceExt:
		return &interfaceExtWrapper***REMOVED***InterfaceExt: t***REMOVED***
	***REMOVED***
	return &extFailWrapper***REMOVED******REMOVED***
***REMOVED***

func baseRV(v interface***REMOVED******REMOVED***) (rv reflect.Value) ***REMOVED***
	for rv = rv4i(v); rv.Kind() == reflect.Ptr; rv = rv.Elem() ***REMOVED***
	***REMOVED***
	return
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

// ------------------ FLOATING POINT -----------------

func isNaN64(f float64) bool ***REMOVED*** return f != f ***REMOVED***
func isNaN32(f float32) bool ***REMOVED*** return f != f ***REMOVED***
func abs32(f float32) float32 ***REMOVED***
	return math.Float32frombits(math.Float32bits(f) &^ (1 << 31))
***REMOVED***

// Per go spec, floats are represented in memory as
// IEEE single or double precision floating point values.
//
// We also looked at the source for stdlib math/modf.go,
// reviewed https://github.com/chewxy/math32
// and read wikipedia documents describing the formats.
//
// It became clear that we could easily look at the bits to determine
// whether any fraction exists.
//
// This is all we need for now.

func noFrac64(f float64) (v bool) ***REMOVED***
	x := math.Float64bits(f)
	e := uint64(x>>52)&0x7FF - 1023 // uint(x>>shift)&mask - bias
	// clear top 12+e bits, the integer part; if the rest is 0, then no fraction.
	if e < 52 ***REMOVED***
		// return x&((1<<64-1)>>(12+e)) == 0
		return x<<(12+e) == 0
	***REMOVED***
	return
***REMOVED***

func noFrac32(f float32) (v bool) ***REMOVED***
	x := math.Float32bits(f)
	e := uint32(x>>23)&0xFF - 127 // uint(x>>shift)&mask - bias
	// clear top 9+e bits, the integer part; if the rest is 0, then no fraction.
	if e < 23 ***REMOVED***
		// return x&((1<<32-1)>>(9+e)) == 0
		return x<<(9+e) == 0
	***REMOVED***
	return
***REMOVED***

// func noFrac(f float64) bool ***REMOVED***
// 	_, frac := math.Modf(float64(f))
// 	return frac == 0
// ***REMOVED***

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

type sfiRv struct ***REMOVED***
	v *structFieldInfo
	r reflect.Value
***REMOVED***

// -----------------

type set []interface***REMOVED******REMOVED***

func (s *set) add(v interface***REMOVED******REMOVED***) (exists bool) ***REMOVED***
	// e.ci is always nil, or len >= 1
	x := *s

	if x == nil ***REMOVED***
		x = make([]interface***REMOVED******REMOVED***, 1, 8)
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

func (s *set) remove(v interface***REMOVED******REMOVED***) (exists bool) ***REMOVED***
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
//
// Also, since pos is a byte (0-255), there's no bounds checks on indexing (cheap).
//
// We previously had bitset128 [16]byte, and bitset32 [4]byte, but those introduces
// bounds checking, so we discarded them, and everyone uses bitset256.
//
// given x > 0 and n > 0 and x is exactly 2^n, then pos/x === pos>>n AND pos%x === pos&(x-1).
// consequently, pos/32 === pos>>5, pos/16 === pos>>4, pos/8 === pos>>3, pos%8 == pos&7

type bitset256 [32]byte

func (x *bitset256) check(pos byte) uint8 ***REMOVED***
	return x[pos>>3] & (1 << (pos & 7))
***REMOVED***

func (x *bitset256) isset(pos byte) bool ***REMOVED***
	return x.check(pos) != 0
	// return x[pos>>3]&(1<<(pos&7)) != 0
***REMOVED***

// func (x *bitset256) issetv(pos byte) byte ***REMOVED***
// 	return x[pos>>3] & (1 << (pos & 7))
// ***REMOVED***

func (x *bitset256) set(pos byte) ***REMOVED***
	x[pos>>3] |= (1 << (pos & 7))
***REMOVED***

type bitset32 uint32

func (x bitset32) set(pos byte) bitset32 ***REMOVED***
	return x | (1 << pos)
***REMOVED***

func (x bitset32) check(pos byte) uint32 ***REMOVED***
	return uint32(x) & (1 << pos)
***REMOVED***
func (x bitset32) isset(pos byte) bool ***REMOVED***
	return x.check(pos) != 0
	// return x&(1<<pos) != 0
***REMOVED***

// func (x *bitset256) unset(pos byte) ***REMOVED***
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
	if len(params) != 0 ***REMOVED***
		panic(fmt.Sprintf(format, params...))
	***REMOVED***
	if len(params) == 0 ***REMOVED***
		panic(format)
	***REMOVED***
	panic("undefined error")
***REMOVED***

// ----------------------------------------------------

type errDecorator interface ***REMOVED***
	wrapErr(in interface***REMOVED******REMOVED***, out *error)
***REMOVED***

type errDecoratorDef struct***REMOVED******REMOVED***

func (errDecoratorDef) wrapErr(v interface***REMOVED******REMOVED***, e *error) ***REMOVED*** *e = fmt.Errorf("%v", v) ***REMOVED***

// ----------------------------------------------------

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

// -------------------

func freelistCapacity(length int) (capacity int) ***REMOVED***
	for capacity = 8; capacity < length; capacity *= 2 ***REMOVED***
	***REMOVED***
	return
***REMOVED***

type bytesFreelist [][]byte

func (x *bytesFreelist) get(length int) (out []byte) ***REMOVED***
	var j int = -1
	for i := 0; i < len(*x); i++ ***REMOVED***
		if cap((*x)[i]) >= length && (j == -1 || cap((*x)[j]) > cap((*x)[i])) ***REMOVED***
			j = i
		***REMOVED***
	***REMOVED***
	if j == -1 ***REMOVED***
		return make([]byte, length, freelistCapacity(length))
	***REMOVED***
	out = (*x)[j][:length]
	(*x)[j] = nil
	for i := 0; i < len(out); i++ ***REMOVED***
		out[i] = 0
	***REMOVED***
	return
***REMOVED***

func (x *bytesFreelist) put(v []byte) ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return
	***REMOVED***
	for i := 0; i < len(*x); i++ ***REMOVED***
		if cap((*x)[i]) == 0 ***REMOVED***
			(*x)[i] = v
			return
		***REMOVED***
	***REMOVED***
	*x = append(*x, v)
***REMOVED***

func (x *bytesFreelist) check(v []byte, length int) (out []byte) ***REMOVED***
	if cap(v) < length ***REMOVED***
		x.put(v)
		return x.get(length)
	***REMOVED***
	return v[:length]
***REMOVED***

// -------------------------

type sfiRvFreelist [][]sfiRv

func (x *sfiRvFreelist) get(length int) (out []sfiRv) ***REMOVED***
	var j int = -1
	for i := 0; i < len(*x); i++ ***REMOVED***
		if cap((*x)[i]) >= length && (j == -1 || cap((*x)[j]) > cap((*x)[i])) ***REMOVED***
			j = i
		***REMOVED***
	***REMOVED***
	if j == -1 ***REMOVED***
		return make([]sfiRv, length, freelistCapacity(length))
	***REMOVED***
	out = (*x)[j][:length]
	(*x)[j] = nil
	for i := 0; i < len(out); i++ ***REMOVED***
		out[i] = sfiRv***REMOVED******REMOVED***
	***REMOVED***
	return
***REMOVED***

func (x *sfiRvFreelist) put(v []sfiRv) ***REMOVED***
	for i := 0; i < len(*x); i++ ***REMOVED***
		if cap((*x)[i]) == 0 ***REMOVED***
			(*x)[i] = v
			return
		***REMOVED***
	***REMOVED***
	*x = append(*x, v)
***REMOVED***

// -----------

// xdebugf printf. the message in red on the terminal.
// Use it in place of fmt.Printf (which it calls internally)
func xdebugf(pattern string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	xdebugAnyf("31", pattern, args...)
***REMOVED***

// xdebug2f printf. the message in blue on the terminal.
// Use it in place of fmt.Printf (which it calls internally)
func xdebug2f(pattern string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	xdebugAnyf("34", pattern, args...)
***REMOVED***

func xdebugAnyf(colorcode, pattern string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if !xdebug ***REMOVED***
		return
	***REMOVED***
	var delim string
	if len(pattern) > 0 && pattern[len(pattern)-1] != '\n' ***REMOVED***
		delim = "\n"
	***REMOVED***
	fmt.Printf("\033[1;"+colorcode+"m"+pattern+delim+"\033[0m", args...)
	// os.Stderr.Flush()
***REMOVED***

// register these here, so that staticcheck stops barfing
var _ = xdebug2f
var _ = xdebugf
var _ = isNaN32
