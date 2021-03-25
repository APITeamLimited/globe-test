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

// isNil says whether the value v is nil.
// This applies to references like map/ptr/unsafepointer/chan/func,
// and non-reference values like interface/slice.
func isNil(v interface***REMOVED******REMOVED***) (rv reflect.Value, isnil bool) ***REMOVED***
	rv = rv4i(v)
	if isnilBitset.isset(byte(rv.Kind())) ***REMOVED***
		isnil = rv.IsNil()
	***REMOVED***
	return
***REMOVED***

func rv4i(i interface***REMOVED******REMOVED***) reflect.Value ***REMOVED***
	return reflect.ValueOf(i)
***REMOVED***

func rv2i(rv reflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
	return rv.Interface()
***REMOVED***

func rvIsNil(rv reflect.Value) bool ***REMOVED***
	return rv.IsNil()
***REMOVED***

func rvSetSliceLen(rv reflect.Value, length int) ***REMOVED***
	rv.SetLen(length)
***REMOVED***

func rvZeroAddrK(t reflect.Type, k reflect.Kind) reflect.Value ***REMOVED***
	return reflect.New(t).Elem()
***REMOVED***

func rvConvert(v reflect.Value, t reflect.Type) (rv reflect.Value) ***REMOVED***
	return v.Convert(t)
***REMOVED***

func rt2id(rt reflect.Type) uintptr ***REMOVED***
	return rv4i(rt).Pointer()
***REMOVED***

func i2rtid(i interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return rv4i(reflect.TypeOf(i)).Pointer()
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
type atomicClsErr struct ***REMOVED***
	v atomic.Value
***REMOVED***

func (x *atomicClsErr) load() (e clsErr) ***REMOVED***
	if i := x.v.Load(); i != nil ***REMOVED***
		e = i.(clsErr)
	***REMOVED***
	return
***REMOVED***

func (x *atomicClsErr) store(p clsErr) ***REMOVED***
	x.v.Store(p)
***REMOVED***

// --------------------------
type atomicTypeInfoSlice struct ***REMOVED*** // expected to be 2 words
	v atomic.Value
***REMOVED***

func (x *atomicTypeInfoSlice) load() (e []rtid2ti) ***REMOVED***
	if i := x.v.Load(); i != nil ***REMOVED***
		e = i.([]rtid2ti)
	***REMOVED***
	return
***REMOVED***

func (x *atomicTypeInfoSlice) store(p []rtid2ti) ***REMOVED***
	x.v.Store(p)
***REMOVED***

// --------------------------
type atomicRtidFnSlice struct ***REMOVED*** // expected to be 2 words
	v atomic.Value
***REMOVED***

func (x *atomicRtidFnSlice) load() (e []codecRtidFn) ***REMOVED***
	if i := x.v.Load(); i != nil ***REMOVED***
		e = i.([]codecRtidFn)
	***REMOVED***
	return
***REMOVED***

func (x *atomicRtidFnSlice) store(p []codecRtidFn) ***REMOVED***
	x.v.Store(p)
***REMOVED***

// --------------------------
func (n *decNaked) ru() reflect.Value ***REMOVED***
	return rv4i(&n.u).Elem()
***REMOVED***
func (n *decNaked) ri() reflect.Value ***REMOVED***
	return rv4i(&n.i).Elem()
***REMOVED***
func (n *decNaked) rf() reflect.Value ***REMOVED***
	return rv4i(&n.f).Elem()
***REMOVED***
func (n *decNaked) rl() reflect.Value ***REMOVED***
	return rv4i(&n.l).Elem()
***REMOVED***
func (n *decNaked) rs() reflect.Value ***REMOVED***
	return rv4i(&n.s).Elem()
***REMOVED***
func (n *decNaked) rt() reflect.Value ***REMOVED***
	return rv4i(&n.t).Elem()
***REMOVED***
func (n *decNaked) rb() reflect.Value ***REMOVED***
	return rv4i(&n.b).Elem()
***REMOVED***

// --------------------------
func rvSetBytes(rv reflect.Value, v []byte) ***REMOVED***
	rv.SetBytes(v)
***REMOVED***

func rvSetString(rv reflect.Value, v string) ***REMOVED***
	rv.SetString(v)
***REMOVED***

func rvSetBool(rv reflect.Value, v bool) ***REMOVED***
	rv.SetBool(v)
***REMOVED***

func rvSetTime(rv reflect.Value, v time.Time) ***REMOVED***
	rv.Set(rv4i(v))
***REMOVED***

func rvSetFloat32(rv reflect.Value, v float32) ***REMOVED***
	rv.SetFloat(float64(v))
***REMOVED***

func rvSetFloat64(rv reflect.Value, v float64) ***REMOVED***
	rv.SetFloat(v)
***REMOVED***

func rvSetInt(rv reflect.Value, v int) ***REMOVED***
	rv.SetInt(int64(v))
***REMOVED***

func rvSetInt8(rv reflect.Value, v int8) ***REMOVED***
	rv.SetInt(int64(v))
***REMOVED***

func rvSetInt16(rv reflect.Value, v int16) ***REMOVED***
	rv.SetInt(int64(v))
***REMOVED***

func rvSetInt32(rv reflect.Value, v int32) ***REMOVED***
	rv.SetInt(int64(v))
***REMOVED***

func rvSetInt64(rv reflect.Value, v int64) ***REMOVED***
	rv.SetInt(v)
***REMOVED***

func rvSetUint(rv reflect.Value, v uint) ***REMOVED***
	rv.SetUint(uint64(v))
***REMOVED***

func rvSetUintptr(rv reflect.Value, v uintptr) ***REMOVED***
	rv.SetUint(uint64(v))
***REMOVED***

func rvSetUint8(rv reflect.Value, v uint8) ***REMOVED***
	rv.SetUint(uint64(v))
***REMOVED***

func rvSetUint16(rv reflect.Value, v uint16) ***REMOVED***
	rv.SetUint(uint64(v))
***REMOVED***

func rvSetUint32(rv reflect.Value, v uint32) ***REMOVED***
	rv.SetUint(uint64(v))
***REMOVED***

func rvSetUint64(rv reflect.Value, v uint64) ***REMOVED***
	rv.SetUint(v)
***REMOVED***

// ----------------

// rvSetDirect is rv.Set for all kinds except reflect.Interface
func rvSetDirect(rv reflect.Value, v reflect.Value) ***REMOVED***
	rv.Set(v)
***REMOVED***

// rvSlice returns a slice of the slice of lenth
func rvSlice(rv reflect.Value, length int) reflect.Value ***REMOVED***
	return rv.Slice(0, length)
***REMOVED***

// ----------------

func rvSliceIndex(rv reflect.Value, i int, ti *typeInfo) reflect.Value ***REMOVED***
	return rv.Index(i)
***REMOVED***

func rvGetSliceLen(rv reflect.Value) int ***REMOVED***
	return rv.Len()
***REMOVED***

func rvGetSliceCap(rv reflect.Value) int ***REMOVED***
	return rv.Cap()
***REMOVED***

func rvGetArrayBytesRO(rv reflect.Value, scratch []byte) (bs []byte) ***REMOVED***
	l := rv.Len()
	if rv.CanAddr() ***REMOVED***
		return rvGetBytes(rv.Slice(0, l))
	***REMOVED***

	if l <= cap(scratch) ***REMOVED***
		bs = scratch[:l]
	***REMOVED*** else ***REMOVED***
		bs = make([]byte, l)
	***REMOVED***
	reflect.Copy(rv4i(bs), rv)
	return
***REMOVED***

func rvGetArray4Slice(rv reflect.Value) (v reflect.Value) ***REMOVED***
	v = rvZeroAddrK(reflectArrayOf(rvGetSliceLen(rv), rv.Type().Elem()), reflect.Array)
	reflect.Copy(v, rv)
	return
***REMOVED***

func rvGetSlice4Array(rv reflect.Value, tslice reflect.Type) (v reflect.Value) ***REMOVED***
	return rv.Slice(0, rv.Len())
***REMOVED***

func rvCopySlice(dest, src reflect.Value) ***REMOVED***
	reflect.Copy(dest, src)
***REMOVED***

// ------------

func rvGetBool(rv reflect.Value) bool ***REMOVED***
	return rv.Bool()
***REMOVED***

func rvGetBytes(rv reflect.Value) []byte ***REMOVED***
	return rv.Bytes()
***REMOVED***

func rvGetTime(rv reflect.Value) time.Time ***REMOVED***
	return rv2i(rv).(time.Time)
***REMOVED***

func rvGetString(rv reflect.Value) string ***REMOVED***
	return rv.String()
***REMOVED***

func rvGetFloat64(rv reflect.Value) float64 ***REMOVED***
	return rv.Float()
***REMOVED***

func rvGetFloat32(rv reflect.Value) float32 ***REMOVED***
	return float32(rv.Float())
***REMOVED***

func rvGetInt(rv reflect.Value) int ***REMOVED***
	return int(rv.Int())
***REMOVED***

func rvGetInt8(rv reflect.Value) int8 ***REMOVED***
	return int8(rv.Int())
***REMOVED***

func rvGetInt16(rv reflect.Value) int16 ***REMOVED***
	return int16(rv.Int())
***REMOVED***

func rvGetInt32(rv reflect.Value) int32 ***REMOVED***
	return int32(rv.Int())
***REMOVED***

func rvGetInt64(rv reflect.Value) int64 ***REMOVED***
	return rv.Int()
***REMOVED***

func rvGetUint(rv reflect.Value) uint ***REMOVED***
	return uint(rv.Uint())
***REMOVED***

func rvGetUint8(rv reflect.Value) uint8 ***REMOVED***
	return uint8(rv.Uint())
***REMOVED***

func rvGetUint16(rv reflect.Value) uint16 ***REMOVED***
	return uint16(rv.Uint())
***REMOVED***

func rvGetUint32(rv reflect.Value) uint32 ***REMOVED***
	return uint32(rv.Uint())
***REMOVED***

func rvGetUint64(rv reflect.Value) uint64 ***REMOVED***
	return rv.Uint()
***REMOVED***

func rvGetUintptr(rv reflect.Value) uintptr ***REMOVED***
	return uintptr(rv.Uint())
***REMOVED***

// ------------ map range and map indexing ----------

func mapGet(m, k, v reflect.Value) (vv reflect.Value) ***REMOVED***
	return m.MapIndex(k)
***REMOVED***

func mapSet(m, k, v reflect.Value) ***REMOVED***
	m.SetMapIndex(k, v)
***REMOVED***

func mapDelete(m, k reflect.Value) ***REMOVED***
	m.SetMapIndex(k, reflect.Value***REMOVED******REMOVED***)
***REMOVED***

// return an addressable reflect value that can be used in mapRange and mapGet operations.
//
// all calls to mapGet or mapRange will call here to get an addressable reflect.Value.
func mapAddressableRV(t reflect.Type, k reflect.Kind) (r reflect.Value) ***REMOVED***
	return // reflect.New(t).Elem()
***REMOVED***

// ---------- ENCODER optimized ---------------

func (e *Encoder) jsondriver() *jsonEncDriver ***REMOVED***
	return e.e.(*jsonEncDriver)
***REMOVED***

// ---------- DECODER optimized ---------------

func (d *Decoder) checkBreak() bool ***REMOVED***
	return d.d.CheckBreak()
***REMOVED***

func (d *Decoder) jsondriver() *jsonDecDriver ***REMOVED***
	return d.d.(*jsonDecDriver)
***REMOVED***
