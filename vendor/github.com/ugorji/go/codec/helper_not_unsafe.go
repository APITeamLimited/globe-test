// +build !go1.7 safe appengine

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"reflect"
	"sync/atomic"
	"time"
)

const safeMode = true

// stringView returns a view of the []byte as a string.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
//
// Usage: Always maintain a reference to v while result of this call is in use,
//        and call keepAlive4BytesView(v) at point where done with view.
func stringView(v []byte) string ***REMOVED***
	return string(v)
***REMOVED***

// bytesView returns a view of the string as a []byte.
// In unsafe mode, it doesn't incur allocation and copying caused by conversion.
// In regular safe mode, it is an allocation and copy.
//
// Usage: Always maintain a reference to v while result of this call is in use,
//        and call keepAlive4BytesView(v) at point where done with view.
func bytesView(v string) []byte ***REMOVED***
	return []byte(v)
***REMOVED***

func definitelyNil(v interface***REMOVED******REMOVED***) bool ***REMOVED***
	// this is a best-effort option.
	// We just return false, so we don't unnecessarily incur the cost of reflection this early.
	return false
***REMOVED***

func rv2i(rv reflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
	return rv.Interface()
***REMOVED***

func rt2id(rt reflect.Type) uintptr ***REMOVED***
	return reflect.ValueOf(rt).Pointer()
***REMOVED***

func rv2rtid(rv reflect.Value) uintptr ***REMOVED***
	return reflect.ValueOf(rv.Type()).Pointer()
***REMOVED***

func i2rtid(i interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return reflect.ValueOf(reflect.TypeOf(i)).Pointer()
***REMOVED***

// --------------------------

func isEmptyValue(v reflect.Value, tinfos *TypeInfos, deref, checkStruct bool) bool ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Invalid:
		return true
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		if deref ***REMOVED***
			if v.IsNil() ***REMOVED***
				return true
			***REMOVED***
			return isEmptyValue(v.Elem(), tinfos, deref, checkStruct)
		***REMOVED***
		return v.IsNil()
	case reflect.Struct:
		return isEmptyStruct(v, tinfos, deref, checkStruct)
	***REMOVED***
	return false
***REMOVED***

// --------------------------
// type ptrToRvMap struct***REMOVED******REMOVED***

// func (*ptrToRvMap) init() ***REMOVED******REMOVED***
// func (*ptrToRvMap) get(i interface***REMOVED******REMOVED***) reflect.Value ***REMOVED***
// 	return reflect.ValueOf(i).Elem()
// ***REMOVED***

// --------------------------
type atomicTypeInfoSlice struct ***REMOVED*** // expected to be 2 words
	v atomic.Value
***REMOVED***

func (x *atomicTypeInfoSlice) load() []rtid2ti ***REMOVED***
	i := x.v.Load()
	if i == nil ***REMOVED***
		return nil
	***REMOVED***
	return i.([]rtid2ti)
***REMOVED***

func (x *atomicTypeInfoSlice) store(p []rtid2ti) ***REMOVED***
	x.v.Store(p)
***REMOVED***

// --------------------------
func (d *Decoder) raw(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetBytes(d.rawBytes())
***REMOVED***

func (d *Decoder) kString(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetString(d.d.DecodeString())
***REMOVED***

func (d *Decoder) kBool(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetBool(d.d.DecodeBool())
***REMOVED***

func (d *Decoder) kTime(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.Set(reflect.ValueOf(d.d.DecodeTime()))
***REMOVED***

func (d *Decoder) kFloat32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	fv := d.d.DecodeFloat64()
	if chkOvf.Float32(fv) ***REMOVED***
		d.errorf("float32 overflow: %v", fv)
	***REMOVED***
	rv.SetFloat(fv)
***REMOVED***

func (d *Decoder) kFloat64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetFloat(d.d.DecodeFloat64())
***REMOVED***

func (d *Decoder) kInt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetInt(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
***REMOVED***

func (d *Decoder) kInt8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetInt(chkOvf.IntV(d.d.DecodeInt64(), 8))
***REMOVED***

func (d *Decoder) kInt16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetInt(chkOvf.IntV(d.d.DecodeInt64(), 16))
***REMOVED***

func (d *Decoder) kInt32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetInt(chkOvf.IntV(d.d.DecodeInt64(), 32))
***REMOVED***

func (d *Decoder) kInt64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetInt(d.d.DecodeInt64())
***REMOVED***

func (d *Decoder) kUint(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
***REMOVED***

func (d *Decoder) kUintptr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
***REMOVED***

func (d *Decoder) kUint8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), 8))
***REMOVED***

func (d *Decoder) kUint16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), 16))
***REMOVED***

func (d *Decoder) kUint32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetUint(chkOvf.UintV(d.d.DecodeUint64(), 32))
***REMOVED***

func (d *Decoder) kUint64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	rv.SetUint(d.d.DecodeUint64())
***REMOVED***

// ----------------

func (e *Encoder) kBool(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeBool(rv.Bool())
***REMOVED***

func (e *Encoder) kTime(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeTime(rv2i(rv).(time.Time))
***REMOVED***

func (e *Encoder) kString(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeString(cUTF8, rv.String())
***REMOVED***

func (e *Encoder) kFloat64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeFloat64(rv.Float())
***REMOVED***

func (e *Encoder) kFloat32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeFloat32(float32(rv.Float()))
***REMOVED***

func (e *Encoder) kInt(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(rv.Int())
***REMOVED***

func (e *Encoder) kInt8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(rv.Int())
***REMOVED***

func (e *Encoder) kInt16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(rv.Int())
***REMOVED***

func (e *Encoder) kInt32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(rv.Int())
***REMOVED***

func (e *Encoder) kInt64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeInt(rv.Int())
***REMOVED***

func (e *Encoder) kUint(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(rv.Uint())
***REMOVED***

func (e *Encoder) kUint8(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(rv.Uint())
***REMOVED***

func (e *Encoder) kUint16(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(rv.Uint())
***REMOVED***

func (e *Encoder) kUint32(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(rv.Uint())
***REMOVED***

func (e *Encoder) kUint64(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(rv.Uint())
***REMOVED***

func (e *Encoder) kUintptr(f *codecFnInfo, rv reflect.Value) ***REMOVED***
	e.e.EncodeUint(rv.Uint())
***REMOVED***

// // keepAlive4BytesView maintains a reference to the input parameter for bytesView.
// //
// // Usage: call this at point where done with the bytes view.
// func keepAlive4BytesView(v string) ***REMOVED******REMOVED***

// // keepAlive4BytesView maintains a reference to the input parameter for stringView.
// //
// // Usage: call this at point where done with the string view.
// func keepAlive4StringView(v []byte) ***REMOVED******REMOVED***

// func definitelyNil(v interface***REMOVED******REMOVED***) bool ***REMOVED***
// 	rv := reflect.ValueOf(v)
// 	switch rv.Kind() ***REMOVED***
// 	case reflect.Invalid:
// 		return true
// 	case reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Slice, reflect.Map, reflect.Func:
// 		return rv.IsNil()
// 	default:
// 		return false
// 	***REMOVED***
// ***REMOVED***
