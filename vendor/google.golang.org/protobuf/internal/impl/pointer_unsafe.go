// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !purego,!appengine

package impl

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

const UnsafeEnabled = true

// Pointer is an opaque pointer type.
type Pointer unsafe.Pointer

// offset represents the offset to a struct field, accessible from a pointer.
// The offset is the byte offset to the field from the start of the struct.
type offset uintptr

// offsetOf returns a field offset for the struct field.
func offsetOf(f reflect.StructField, x exporter) offset ***REMOVED***
	return offset(f.Offset)
***REMOVED***

// IsValid reports whether the offset is valid.
func (f offset) IsValid() bool ***REMOVED*** return f != invalidOffset ***REMOVED***

// invalidOffset is an invalid field offset.
var invalidOffset = ^offset(0)

// zeroOffset is a noop when calling pointer.Apply.
var zeroOffset = offset(0)

// pointer is a pointer to a message struct or field.
type pointer struct***REMOVED*** p unsafe.Pointer ***REMOVED***

// pointerOf returns p as a pointer.
func pointerOf(p Pointer) pointer ***REMOVED***
	return pointer***REMOVED***p: unsafe.Pointer(p)***REMOVED***
***REMOVED***

// pointerOfValue returns v as a pointer.
func pointerOfValue(v reflect.Value) pointer ***REMOVED***
	return pointer***REMOVED***p: unsafe.Pointer(v.Pointer())***REMOVED***
***REMOVED***

// pointerOfIface returns the pointer portion of an interface.
func pointerOfIface(v interface***REMOVED******REMOVED***) pointer ***REMOVED***
	type ifaceHeader struct ***REMOVED***
		Type unsafe.Pointer
		Data unsafe.Pointer
	***REMOVED***
	return pointer***REMOVED***p: (*ifaceHeader)(unsafe.Pointer(&v)).Data***REMOVED***
***REMOVED***

// IsNil reports whether the pointer is nil.
func (p pointer) IsNil() bool ***REMOVED***
	return p.p == nil
***REMOVED***

// Apply adds an offset to the pointer to derive a new pointer
// to a specified field. The pointer must be valid and pointing at a struct.
func (p pointer) Apply(f offset) pointer ***REMOVED***
	if p.IsNil() ***REMOVED***
		panic("invalid nil pointer")
	***REMOVED***
	return pointer***REMOVED***p: unsafe.Pointer(uintptr(p.p) + uintptr(f))***REMOVED***
***REMOVED***

// AsValueOf treats p as a pointer to an object of type t and returns the value.
// It is equivalent to reflect.ValueOf(p.AsIfaceOf(t))
func (p pointer) AsValueOf(t reflect.Type) reflect.Value ***REMOVED***
	return reflect.NewAt(t, p.p)
***REMOVED***

// AsIfaceOf treats p as a pointer to an object of type t and returns the value.
// It is equivalent to p.AsValueOf(t).Interface()
func (p pointer) AsIfaceOf(t reflect.Type) interface***REMOVED******REMOVED*** ***REMOVED***
	// TODO: Use tricky unsafe magic to directly create ifaceHeader.
	return p.AsValueOf(t).Interface()
***REMOVED***

func (p pointer) Bool() *bool                           ***REMOVED*** return (*bool)(p.p) ***REMOVED***
func (p pointer) BoolPtr() **bool                       ***REMOVED*** return (**bool)(p.p) ***REMOVED***
func (p pointer) BoolSlice() *[]bool                    ***REMOVED*** return (*[]bool)(p.p) ***REMOVED***
func (p pointer) Int32() *int32                         ***REMOVED*** return (*int32)(p.p) ***REMOVED***
func (p pointer) Int32Ptr() **int32                     ***REMOVED*** return (**int32)(p.p) ***REMOVED***
func (p pointer) Int32Slice() *[]int32                  ***REMOVED*** return (*[]int32)(p.p) ***REMOVED***
func (p pointer) Int64() *int64                         ***REMOVED*** return (*int64)(p.p) ***REMOVED***
func (p pointer) Int64Ptr() **int64                     ***REMOVED*** return (**int64)(p.p) ***REMOVED***
func (p pointer) Int64Slice() *[]int64                  ***REMOVED*** return (*[]int64)(p.p) ***REMOVED***
func (p pointer) Uint32() *uint32                       ***REMOVED*** return (*uint32)(p.p) ***REMOVED***
func (p pointer) Uint32Ptr() **uint32                   ***REMOVED*** return (**uint32)(p.p) ***REMOVED***
func (p pointer) Uint32Slice() *[]uint32                ***REMOVED*** return (*[]uint32)(p.p) ***REMOVED***
func (p pointer) Uint64() *uint64                       ***REMOVED*** return (*uint64)(p.p) ***REMOVED***
func (p pointer) Uint64Ptr() **uint64                   ***REMOVED*** return (**uint64)(p.p) ***REMOVED***
func (p pointer) Uint64Slice() *[]uint64                ***REMOVED*** return (*[]uint64)(p.p) ***REMOVED***
func (p pointer) Float32() *float32                     ***REMOVED*** return (*float32)(p.p) ***REMOVED***
func (p pointer) Float32Ptr() **float32                 ***REMOVED*** return (**float32)(p.p) ***REMOVED***
func (p pointer) Float32Slice() *[]float32              ***REMOVED*** return (*[]float32)(p.p) ***REMOVED***
func (p pointer) Float64() *float64                     ***REMOVED*** return (*float64)(p.p) ***REMOVED***
func (p pointer) Float64Ptr() **float64                 ***REMOVED*** return (**float64)(p.p) ***REMOVED***
func (p pointer) Float64Slice() *[]float64              ***REMOVED*** return (*[]float64)(p.p) ***REMOVED***
func (p pointer) String() *string                       ***REMOVED*** return (*string)(p.p) ***REMOVED***
func (p pointer) StringPtr() **string                   ***REMOVED*** return (**string)(p.p) ***REMOVED***
func (p pointer) StringSlice() *[]string                ***REMOVED*** return (*[]string)(p.p) ***REMOVED***
func (p pointer) Bytes() *[]byte                        ***REMOVED*** return (*[]byte)(p.p) ***REMOVED***
func (p pointer) BytesSlice() *[][]byte                 ***REMOVED*** return (*[][]byte)(p.p) ***REMOVED***
func (p pointer) WeakFields() *weakFields               ***REMOVED*** return (*weakFields)(p.p) ***REMOVED***
func (p pointer) Extensions() *map[int32]ExtensionField ***REMOVED*** return (*map[int32]ExtensionField)(p.p) ***REMOVED***

func (p pointer) Elem() pointer ***REMOVED***
	return pointer***REMOVED***p: *(*unsafe.Pointer)(p.p)***REMOVED***
***REMOVED***

// PointerSlice loads []*T from p as a []pointer.
// The value returned is aliased with the original slice.
// This behavior differs from the implementation in pointer_reflect.go.
func (p pointer) PointerSlice() []pointer ***REMOVED***
	// Super-tricky - p should point to a []*T where T is a
	// message type. We load it as []pointer.
	return *(*[]pointer)(p.p)
***REMOVED***

// AppendPointerSlice appends v to p, which must be a []*T.
func (p pointer) AppendPointerSlice(v pointer) ***REMOVED***
	*(*[]pointer)(p.p) = append(*(*[]pointer)(p.p), v)
***REMOVED***

// SetPointer sets *p to v.
func (p pointer) SetPointer(v pointer) ***REMOVED***
	*(*unsafe.Pointer)(p.p) = (unsafe.Pointer)(v.p)
***REMOVED***

// Static check that MessageState does not exceed the size of a pointer.
const _ = uint(unsafe.Sizeof(unsafe.Pointer(nil)) - unsafe.Sizeof(MessageState***REMOVED******REMOVED***))

func (Export) MessageStateOf(p Pointer) *messageState ***REMOVED***
	// Super-tricky - see documentation on MessageState.
	return (*messageState)(unsafe.Pointer(p))
***REMOVED***
func (ms *messageState) pointer() pointer ***REMOVED***
	// Super-tricky - see documentation on MessageState.
	return pointer***REMOVED***p: unsafe.Pointer(ms)***REMOVED***
***REMOVED***
func (ms *messageState) messageInfo() *MessageInfo ***REMOVED***
	mi := ms.LoadMessageInfo()
	if mi == nil ***REMOVED***
		panic("invalid nil message info; this suggests memory corruption due to a race or shallow copy on the message struct")
	***REMOVED***
	return mi
***REMOVED***
func (ms *messageState) LoadMessageInfo() *MessageInfo ***REMOVED***
	return (*MessageInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ms.atomicMessageInfo))))
***REMOVED***
func (ms *messageState) StoreMessageInfo(mi *MessageInfo) ***REMOVED***
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&ms.atomicMessageInfo)), unsafe.Pointer(mi))
***REMOVED***

type atomicNilMessage struct***REMOVED*** p unsafe.Pointer ***REMOVED*** // p is a *messageReflectWrapper

func (m *atomicNilMessage) Init(mi *MessageInfo) *messageReflectWrapper ***REMOVED***
	if p := atomic.LoadPointer(&m.p); p != nil ***REMOVED***
		return (*messageReflectWrapper)(p)
	***REMOVED***
	w := &messageReflectWrapper***REMOVED***mi: mi***REMOVED***
	atomic.CompareAndSwapPointer(&m.p, nil, (unsafe.Pointer)(w))
	return (*messageReflectWrapper)(atomic.LoadPointer(&m.p))
***REMOVED***
