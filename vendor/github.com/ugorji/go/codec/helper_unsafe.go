// +build !safe
// +build !appengine
// +build go1.7

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"
)

// This file has unsafe variants of some helper methods.
// NOTE: See helper_not_unsafe.go for the usage information.

// var zeroRTv [4]uintptr

const safeMode = false
const unsafeFlagIndir = 1 << 7 // keep in sync with GO_ROOT/src/reflect/value.go

type unsafeString struct ***REMOVED***
	Data unsafe.Pointer
	Len  int
***REMOVED***

type unsafeSlice struct ***REMOVED***
	Data unsafe.Pointer
	Len  int
	Cap  int
***REMOVED***

type unsafeIntf struct ***REMOVED***
	typ  unsafe.Pointer
	word unsafe.Pointer
***REMOVED***

type unsafeReflectValue struct ***REMOVED***
	typ  unsafe.Pointer
	ptr  unsafe.Pointer
	flag uintptr
***REMOVED***

func stringView(v []byte) string ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return ""
	***REMOVED***
	bx := (*unsafeSlice)(unsafe.Pointer(&v))
	return *(*string)(unsafe.Pointer(&unsafeString***REMOVED***bx.Data, bx.Len***REMOVED***))
***REMOVED***

func bytesView(v string) []byte ***REMOVED***
	if len(v) == 0 ***REMOVED***
		return zeroByteSlice
	***REMOVED***
	sx := (*unsafeString)(unsafe.Pointer(&v))
	return *(*[]byte)(unsafe.Pointer(&unsafeSlice***REMOVED***sx.Data, sx.Len, sx.Len***REMOVED***))
***REMOVED***

func definitelyNil(v interface***REMOVED******REMOVED***) bool ***REMOVED***
	// There is no global way of checking if an interface is nil.
	// For true references (map, ptr, func, chan), you can just look
	// at the word of the interface. However, for slices, you have to dereference
	// the word, and get a pointer to the 3-word interface value.
	//
	// However, the following are cheap calls
	// - TypeOf(interface): cheap 2-line call.
	// - ValueOf(interface***REMOVED******REMOVED***): expensive
	// - type.Kind: cheap call through an interface
	// - Value.Type(): cheap call
	//                 except it's a method value (e.g. r.Read, which implies that it is a Func)

	return ((*unsafeIntf)(unsafe.Pointer(&v))).word == nil
***REMOVED***

func rv2i(rv reflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
	// TODO: consider a more generally-known optimization for reflect.Value ==> Interface
	//
	// Currently, we use this fragile method that taps into implememtation details from
	// the source go stdlib reflect/value.go, and trims the implementation.

	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	// true references (map, func, chan, ptr - NOT slice) may be double-referenced as flagIndir
	var ptr unsafe.Pointer
	if refBitset.isset(byte(urv.flag&(1<<5-1))) && urv.flag&unsafeFlagIndir != 0 ***REMOVED***
		ptr = *(*unsafe.Pointer)(urv.ptr)
	***REMOVED*** else ***REMOVED***
		ptr = urv.ptr
	***REMOVED***
	return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***typ: urv.typ, word: ptr***REMOVED***))
***REMOVED***

func rt2id(rt reflect.Type) uintptr ***REMOVED***
	return uintptr(((*unsafeIntf)(unsafe.Pointer(&rt))).word)
***REMOVED***

func rv2rtid(rv reflect.Value) uintptr ***REMOVED***
	return uintptr((*unsafeReflectValue)(unsafe.Pointer(&rv)).typ)
***REMOVED***

func i2rtid(i interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return uintptr(((*unsafeIntf)(unsafe.Pointer(&i))).typ)
***REMOVED***

// --------------------------

func isEmptyValue(v reflect.Value, tinfos *TypeInfos, deref, checkStruct bool) bool ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	if urv.flag == 0 ***REMOVED***
		return true
	***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Invalid:
		return true
	case reflect.String:
		return (*unsafeString)(urv.ptr).Len == 0
	case reflect.Slice:
		return (*unsafeSlice)(urv.ptr).Len == 0
	case reflect.Bool:
		return !*(*bool)(urv.ptr)
	case reflect.Int:
		return *(*int)(urv.ptr) == 0
	case reflect.Int8:
		return *(*int8)(urv.ptr) == 0
	case reflect.Int16:
		return *(*int16)(urv.ptr) == 0
	case reflect.Int32:
		return *(*int32)(urv.ptr) == 0
	case reflect.Int64:
		return *(*int64)(urv.ptr) == 0
	case reflect.Uint:
		return *(*uint)(urv.ptr) == 0
	case reflect.Uint8:
		return *(*uint8)(urv.ptr) == 0
	case reflect.Uint16:
		return *(*uint16)(urv.ptr) == 0
	case reflect.Uint32:
		return *(*uint32)(urv.ptr) == 0
	case reflect.Uint64:
		return *(*uint64)(urv.ptr) == 0
	case reflect.Uintptr:
		return *(*uintptr)(urv.ptr) == 0
	case reflect.Float32:
		return *(*float32)(urv.ptr) == 0
	case reflect.Float64:
		return *(*float64)(urv.ptr) == 0
	case reflect.Interface:
		isnil := urv.ptr == nil || *(*unsafe.Pointer)(urv.ptr) == nil
		if deref ***REMOVED***
			if isnil ***REMOVED***
				return true
			***REMOVED***
			return isEmptyValue(v.Elem(), tinfos, deref, checkStruct)
		***REMOVED***
		return isnil
	case reflect.Ptr:
		isnil := urv.ptr == nil
		if deref ***REMOVED***
			if isnil ***REMOVED***
				return true
			***REMOVED***
			return isEmptyValue(v.Elem(), tinfos, deref, checkStruct)
		***REMOVED***
		return isnil
	case reflect.Struct:
		return isEmptyStruct(v, tinfos, deref, checkStruct)
	case reflect.Map, reflect.Array, reflect.Chan:
		return v.Len() == 0
	***REMOVED***
	return false
***REMOVED***

// --------------------------

type atomicTypeInfoSlice struct ***REMOVED*** // expected to be 2 words
	v unsafe.Pointer // data array - Pointer (not uintptr) to maintain GC reference
	l int64          // length of the data array
***REMOVED***

func (x *atomicTypeInfoSlice) load() []rtid2ti ***REMOVED***
	l := int(atomic.LoadInt64(&x.l))
	if l == 0 ***REMOVED***
		return nil
	***REMOVED***
	return *(*[]rtid2ti)(unsafe.Pointer(&unsafeSlice***REMOVED***Data: atomic.LoadPointer(&x.v), Len: l, Cap: l***REMOVED***))
	// return (*[]rtid2ti)(atomic.LoadPointer(&x.v))
***REMOVED***

func (x *atomicTypeInfoSlice) store(p []rtid2ti) ***REMOVED***
	s := (*unsafeSlice)(unsafe.Pointer(&p))
	atomic.StorePointer(&x.v, s.Data)
	atomic.StoreInt64(&x.l, int64(s.Len))
	// atomic.StorePointer(&x.v, unsafe.Pointer(p))
***REMOVED***

// --------------------------
func (d *Decoder) raw(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*[]byte)(urv.ptr) = d.rawBytes()
***REMOVED***

func (d *Decoder) kString(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*string)(urv.ptr) = d.d.DecodeString()
***REMOVED***

func (d *Decoder) kBool(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*bool)(urv.ptr) = d.d.DecodeBool()
***REMOVED***

func (d *Decoder) kTime(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*time.Time)(urv.ptr) = d.d.DecodeTime()
***REMOVED***

func (d *Decoder) kFloat32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	fv := d.d.DecodeFloat64()
	if chkOvf.Float32(fv) ***REMOVED***
		d.errorf("float32 overflow: %v", fv)
	***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*float32)(urv.ptr) = float32(fv)
***REMOVED***

func (d *Decoder) kFloat64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*float64)(urv.ptr) = d.d.DecodeFloat64()
***REMOVED***

func (d *Decoder) kInt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int)(urv.ptr) = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
***REMOVED***

func (d *Decoder) kInt8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int8)(urv.ptr) = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
***REMOVED***

func (d *Decoder) kInt16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int16)(urv.ptr) = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
***REMOVED***

func (d *Decoder) kInt32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int32)(urv.ptr) = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
***REMOVED***

func (d *Decoder) kInt64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int64)(urv.ptr) = d.d.DecodeInt64()
***REMOVED***

func (d *Decoder) kUint(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint)(urv.ptr) = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
***REMOVED***

func (d *Decoder) kUintptr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uintptr)(urv.ptr) = uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
***REMOVED***

func (d *Decoder) kUint8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint8)(urv.ptr) = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
***REMOVED***

func (d *Decoder) kUint16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint16)(urv.ptr) = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
***REMOVED***

func (d *Decoder) kUint32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint32)(urv.ptr) = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
***REMOVED***

func (d *Decoder) kUint64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint64)(urv.ptr) = d.d.DecodeUint64()
***REMOVED***

// ------------

func (e *Encoder) kBool(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeBool(*(*bool)(v.ptr))
***REMOVED***

func (e *Encoder) kTime(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeTime(*(*time.Time)(v.ptr))
***REMOVED***

func (e *Encoder) kString(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeString(cUTF8, *(*string)(v.ptr))
***REMOVED***

func (e *Encoder) kFloat64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeFloat64(*(*float64)(v.ptr))
***REMOVED***

func (e *Encoder) kFloat32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeFloat32(*(*float32)(v.ptr))
***REMOVED***

func (e *Encoder) kInt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int)(v.ptr)))
***REMOVED***

func (e *Encoder) kInt8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int8)(v.ptr)))
***REMOVED***

func (e *Encoder) kInt16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int16)(v.ptr)))
***REMOVED***

func (e *Encoder) kInt32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int32)(v.ptr)))
***REMOVED***

func (e *Encoder) kInt64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeInt(int64(*(*int64)(v.ptr)))
***REMOVED***

func (e *Encoder) kUint(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint)(v.ptr)))
***REMOVED***

func (e *Encoder) kUint8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint8)(v.ptr)))
***REMOVED***

func (e *Encoder) kUint16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint16)(v.ptr)))
***REMOVED***

func (e *Encoder) kUint32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint32)(v.ptr)))
***REMOVED***

func (e *Encoder) kUint64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uint64)(v.ptr)))
***REMOVED***

func (e *Encoder) kUintptr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	e.e.EncodeUint(uint64(*(*uintptr)(v.ptr)))
***REMOVED***

// ------------

// func (d *Decoder) raw(f *codecFnInfo, rv reflect.Value) ***REMOVED***
// 	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
// 	// if urv.flag&unsafeFlagIndir != 0 ***REMOVED***
// 	// 	urv.ptr = *(*unsafe.Pointer)(urv.ptr)
// 	// ***REMOVED***
// 	*(*[]byte)(urv.ptr) = d.rawBytes()
// ***REMOVED***

// func rv0t(rt reflect.Type) reflect.Value ***REMOVED***
// 	ut := (*unsafeIntf)(unsafe.Pointer(&rt))
// 	// we need to determine whether ifaceIndir, and then whether to just pass 0 as the ptr
// 	uv := unsafeReflectValue***REMOVED***ut.word, &zeroRTv, flag(rt.Kind())***REMOVED***
// 	return *(*reflect.Value)(unsafe.Pointer(&uv***REMOVED***)
// ***REMOVED***

// func rv2i(rv reflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
// 	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
// 	// true references (map, func, chan, ptr - NOT slice) may be double-referenced as flagIndir
// 	var ptr unsafe.Pointer
// 	// kk := reflect.Kind(urv.flag & (1<<5 - 1))
// 	// if (kk == reflect.Map || kk == reflect.Ptr || kk == reflect.Chan || kk == reflect.Func) && urv.flag&unsafeFlagIndir != 0 ***REMOVED***
// 	if refBitset.isset(byte(urv.flag&(1<<5-1))) && urv.flag&unsafeFlagIndir != 0 ***REMOVED***
// 		ptr = *(*unsafe.Pointer)(urv.ptr)
// 	***REMOVED*** else ***REMOVED***
// 		ptr = urv.ptr
// 	***REMOVED***
// 	return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***typ: urv.typ, word: ptr***REMOVED***))
// 	// return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: *(*unsafe.Pointer)(urv.ptr), typ: urv.typ***REMOVED***))
// 	// return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: urv.ptr, typ: urv.typ***REMOVED***))
// ***REMOVED***

// func definitelyNil(v interface***REMOVED******REMOVED***) bool ***REMOVED***
// 	var ui *unsafeIntf = (*unsafeIntf)(unsafe.Pointer(&v))
// 	if ui.word == nil ***REMOVED***
// 		return true
// 	***REMOVED***
// 	var tk = reflect.TypeOf(v).Kind()
// 	return (tk == reflect.Interface || tk == reflect.Slice) && *(*unsafe.Pointer)(ui.word) == nil
// 	fmt.Printf(">>>> definitely nil: isnil: %v, TYPE: \t%T, word: %v, *word: %v, type: %v, nil: %v\n",
// 	v == nil, v, word, *((*unsafe.Pointer)(word)), ui.typ, nil)
// ***REMOVED***

// func keepAlive4BytesView(v string) ***REMOVED***
// 	runtime.KeepAlive(v)
// ***REMOVED***

// func keepAlive4StringView(v []byte) ***REMOVED***
// 	runtime.KeepAlive(v)
// ***REMOVED***

// func rt2id(rt reflect.Type) uintptr ***REMOVED***
// 	return uintptr(((*unsafeIntf)(unsafe.Pointer(&rt))).word)
// 	// var i interface***REMOVED******REMOVED*** = rt
// 	// // ui := (*unsafeIntf)(unsafe.Pointer(&i))
// 	// return ((*unsafeIntf)(unsafe.Pointer(&i))).word
// ***REMOVED***

// func rv2i(rv reflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
// 	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
// 	// non-reference type: already indir
// 	// reference type: depend on flagIndir property ('cos maybe was double-referenced)
// 	// const (unsafeRvFlagKindMask    = 1<<5 - 1 , unsafeRvFlagIndir       = 1 << 7 )
// 	// rvk := reflect.Kind(urv.flag & (1<<5 - 1))
// 	// if (rvk == reflect.Chan ||
// 	// 	rvk == reflect.Func ||
// 	// 	rvk == reflect.Interface ||
// 	// 	rvk == reflect.Map ||
// 	// 	rvk == reflect.Ptr ||
// 	// 	rvk == reflect.UnsafePointer) && urv.flag&(1<<8) != 0 ***REMOVED***
// 	// 	fmt.Printf(">>>>> ---- double indirect reference: %v, %v\n", rvk, rv.Type())
// 	// 	return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: *(*unsafe.Pointer)(urv.ptr), typ: urv.typ***REMOVED***))
// 	// ***REMOVED***
// 	if urv.flag&(1<<5-1) == uintptr(reflect.Map) && urv.flag&(1<<7) != 0 ***REMOVED***
// 		// fmt.Printf(">>>>> ---- double indirect reference: %v, %v\n", rvk, rv.Type())
// 		return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: *(*unsafe.Pointer)(urv.ptr), typ: urv.typ***REMOVED***))
// 	***REMOVED***
// 	// fmt.Printf(">>>>> ++++ direct reference: %v, %v\n", rvk, rv.Type())
// 	return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: urv.ptr, typ: urv.typ***REMOVED***))
// ***REMOVED***

// const (
// 	unsafeRvFlagKindMask    = 1<<5 - 1
// 	unsafeRvKindDirectIface = 1 << 5
// 	unsafeRvFlagIndir       = 1 << 7
// 	unsafeRvFlagAddr        = 1 << 8
// 	unsafeRvFlagMethod      = 1 << 9

// 	_USE_RV_INTERFACE bool = false
// 	_UNSAFE_RV_DEBUG       = true
// )

// type unsafeRtype struct ***REMOVED***
// 	_    [2]uintptr
// 	_    uint32
// 	_    uint8
// 	_    uint8
// 	_    uint8
// 	kind uint8
// 	_    [2]uintptr
// 	_    int32
// ***REMOVED***

// func _rv2i(rv reflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
// 	// Note: From use,
// 	//   - it's never an interface
// 	//   - the only calls here are for ifaceIndir types.
// 	//     (though that conditional is wrong)
// 	//     To know for sure, we need the value of t.kind (which is not exposed).
// 	//
// 	// Need to validate the path: type is indirect ==> only value is indirect ==> default (value is direct)
// 	//    - Type indirect, Value indirect: ==> numbers, boolean, slice, struct, array, string
// 	//    - Type Direct,   Value indirect: ==> map???
// 	//    - Type Direct,   Value direct:   ==> pointers, unsafe.Pointer, func, chan, map
// 	//
// 	// TRANSLATES TO:
// 	//    if typeIndirect ***REMOVED*** ***REMOVED*** else if valueIndirect ***REMOVED*** ***REMOVED*** else ***REMOVED*** ***REMOVED***
// 	//
// 	// Since we don't deal with funcs, then "flagNethod" is unset, and can be ignored.

// 	if _USE_RV_INTERFACE ***REMOVED***
// 		return rv.Interface()
// 	***REMOVED***
// 	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))

// 	// if urv.flag&unsafeRvFlagMethod != 0 || urv.flag&unsafeRvFlagKindMask == uintptr(reflect.Interface) ***REMOVED***
// 	// 	println("***** IS flag method or interface: delegating to rv.Interface()")
// 	// 	return rv.Interface()
// 	// ***REMOVED***

// 	// if urv.flag&unsafeRvFlagKindMask == uintptr(reflect.Interface) ***REMOVED***
// 	// 	println("***** IS Interface: delegate to rv.Interface")
// 	// 	return rv.Interface()
// 	// ***REMOVED***
// 	// if urv.flag&unsafeRvFlagKindMask&unsafeRvKindDirectIface == 0 ***REMOVED***
// 	// 	if urv.flag&unsafeRvFlagAddr == 0 ***REMOVED***
// 	// 		println("***** IS ifaceIndir typ")
// 	// 		// ui := unsafeIntf***REMOVED***word: urv.ptr, typ: urv.typ***REMOVED***
// 	// 		// return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&ui))
// 	// 		// return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: urv.ptr, typ: urv.typ***REMOVED***))
// 	// 	***REMOVED***
// 	// ***REMOVED*** else if urv.flag&unsafeRvFlagIndir != 0 ***REMOVED***
// 	// 	println("***** IS flagindir")
// 	// 	// return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: *(*unsafe.Pointer)(urv.ptr), typ: urv.typ***REMOVED***))
// 	// ***REMOVED*** else ***REMOVED***
// 	// 	println("***** NOT flagindir")
// 	// 	return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: urv.ptr, typ: urv.typ***REMOVED***))
// 	// ***REMOVED***
// 	// println("***** default: delegate to rv.Interface")

// 	urt := (*unsafeRtype)(unsafe.Pointer(urv.typ))
// 	if _UNSAFE_RV_DEBUG ***REMOVED***
// 		fmt.Printf(">>>> start: %v: ", rv.Type())
// 		fmt.Printf("%v - %v\n", *urv, *urt)
// 	***REMOVED***
// 	if urt.kind&unsafeRvKindDirectIface == 0 ***REMOVED***
// 		if _UNSAFE_RV_DEBUG ***REMOVED***
// 			fmt.Printf("**** +ifaceIndir type: %v\n", rv.Type())
// 		***REMOVED***
// 		// println("***** IS ifaceIndir typ")
// 		// if true || urv.flag&unsafeRvFlagAddr == 0 ***REMOVED***
// 		// 	// println("    ***** IS NOT addr")
// 		return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: urv.ptr, typ: urv.typ***REMOVED***))
// 		// ***REMOVED***
// 	***REMOVED*** else if urv.flag&unsafeRvFlagIndir != 0 ***REMOVED***
// 		if _UNSAFE_RV_DEBUG ***REMOVED***
// 			fmt.Printf("**** +flagIndir type: %v\n", rv.Type())
// 		***REMOVED***
// 		// println("***** IS flagindir")
// 		return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: *(*unsafe.Pointer)(urv.ptr), typ: urv.typ***REMOVED***))
// 	***REMOVED*** else ***REMOVED***
// 		if _UNSAFE_RV_DEBUG ***REMOVED***
// 			fmt.Printf("**** -flagIndir type: %v\n", rv.Type())
// 		***REMOVED***
// 		// println("***** NOT flagindir")
// 		return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&unsafeIntf***REMOVED***word: urv.ptr, typ: urv.typ***REMOVED***))
// 	***REMOVED***
// 	// println("***** default: delegating to rv.Interface()")
// 	// return rv.Interface()
// ***REMOVED***

// var staticM0 = make(map[string]uint64)
// var staticI0 = (int32)(-5)

// func staticRv2iTest() ***REMOVED***
// 	i0 := (int32)(-5)
// 	m0 := make(map[string]uint16)
// 	m0["1"] = 1
// 	for _, i := range []interface***REMOVED******REMOVED******REMOVED***
// 		(int)(7),
// 		(uint)(8),
// 		(int16)(-9),
// 		(uint16)(19),
// 		(uintptr)(77),
// 		(bool)(true),
// 		float32(-32.7),
// 		float64(64.9),
// 		complex(float32(19), 5),
// 		complex(float64(-32), 7),
// 		[4]uint64***REMOVED***1, 2, 3, 4***REMOVED***,
// 		(chan<- int)(nil), // chan,
// 		rv2i,              // func
// 		io.Writer(ioutil.Discard),
// 		make(map[string]uint),
// 		(map[string]uint)(nil),
// 		staticM0,
// 		m0,
// 		&m0,
// 		i0,
// 		&i0,
// 		&staticI0,
// 		&staticM0,
// 		[]uint32***REMOVED***6, 7, 8***REMOVED***,
// 		"abc",
// 		Raw***REMOVED******REMOVED***,
// 		RawExt***REMOVED******REMOVED***,
// 		&Raw***REMOVED******REMOVED***,
// 		&RawExt***REMOVED******REMOVED***,
// 		unsafe.Pointer(&i0),
// 	***REMOVED*** ***REMOVED***
// 		i2 := rv2i(reflect.ValueOf(i))
// 		eq := reflect.DeepEqual(i, i2)
// 		fmt.Printf(">>>> %v == %v? %v\n", i, i2, eq)
// 	***REMOVED***
// 	// os.Exit(0)
// ***REMOVED***

// func init() ***REMOVED***
// 	staticRv2iTest()
// ***REMOVED***

// func rv2i(rv reflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
// 	if _USE_RV_INTERFACE || rv.Kind() == reflect.Interface || rv.CanAddr() ***REMOVED***
// 		return rv.Interface()
// 	***REMOVED***
// 	// var i interface***REMOVED******REMOVED***
// 	// ui := (*unsafeIntf)(unsafe.Pointer(&i))
// 	var ui unsafeIntf
// 	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
// 	// fmt.Printf("urv: flag: %b, typ: %b, ptr: %b\n", urv.flag, uintptr(urv.typ), uintptr(urv.ptr))
// 	if (urv.flag&unsafeRvFlagKindMask)&unsafeRvKindDirectIface == 0 ***REMOVED***
// 		if urv.flag&unsafeRvFlagAddr != 0 ***REMOVED***
// 			println("***** indirect and addressable! Needs typed move - delegate to rv.Interface()")
// 			return rv.Interface()
// 		***REMOVED***
// 		println("****** indirect type/kind")
// 		ui.word = urv.ptr
// 	***REMOVED*** else if urv.flag&unsafeRvFlagIndir != 0 ***REMOVED***
// 		println("****** unsafe rv flag indir")
// 		ui.word = *(*unsafe.Pointer)(urv.ptr)
// 	***REMOVED*** else ***REMOVED***
// 		println("****** default: assign prt to word directly")
// 		ui.word = urv.ptr
// 	***REMOVED***
// 	// ui.word = urv.ptr
// 	ui.typ = urv.typ
// 	// fmt.Printf("(pointers) ui.typ: %p, word: %p\n", ui.typ, ui.word)
// 	// fmt.Printf("(binary)   ui.typ: %b, word: %b\n", uintptr(ui.typ), uintptr(ui.word))
// 	return *(*interface***REMOVED******REMOVED***)(unsafe.Pointer(&ui))
// 	// return i
// ***REMOVED***
