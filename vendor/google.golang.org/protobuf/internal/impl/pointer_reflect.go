// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build purego appengine

package impl

import (
	"fmt"
	"reflect"
	"sync"
)

const UnsafeEnabled = false

// Pointer is an opaque pointer type.
type Pointer interface***REMOVED******REMOVED***

// offset represents the offset to a struct field, accessible from a pointer.
// The offset is the field index into a struct.
type offset struct ***REMOVED***
	index  int
	export exporter
***REMOVED***

// offsetOf returns a field offset for the struct field.
func offsetOf(f reflect.StructField, x exporter) offset ***REMOVED***
	if len(f.Index) != 1 ***REMOVED***
		panic("embedded structs are not supported")
	***REMOVED***
	if f.PkgPath == "" ***REMOVED***
		return offset***REMOVED***index: f.Index[0]***REMOVED*** // field is already exported
	***REMOVED***
	if x == nil ***REMOVED***
		panic("exporter must be provided for unexported field")
	***REMOVED***
	return offset***REMOVED***index: f.Index[0], export: x***REMOVED***
***REMOVED***

// IsValid reports whether the offset is valid.
func (f offset) IsValid() bool ***REMOVED*** return f.index >= 0 ***REMOVED***

// invalidOffset is an invalid field offset.
var invalidOffset = offset***REMOVED***index: -1***REMOVED***

// zeroOffset is a noop when calling pointer.Apply.
var zeroOffset = offset***REMOVED***index: 0***REMOVED***

// pointer is an abstract representation of a pointer to a struct or field.
type pointer struct***REMOVED*** v reflect.Value ***REMOVED***

// pointerOf returns p as a pointer.
func pointerOf(p Pointer) pointer ***REMOVED***
	return pointerOfIface(p)
***REMOVED***

// pointerOfValue returns v as a pointer.
func pointerOfValue(v reflect.Value) pointer ***REMOVED***
	return pointer***REMOVED***v: v***REMOVED***
***REMOVED***

// pointerOfIface returns the pointer portion of an interface.
func pointerOfIface(v interface***REMOVED******REMOVED***) pointer ***REMOVED***
	return pointer***REMOVED***v: reflect.ValueOf(v)***REMOVED***
***REMOVED***

// IsNil reports whether the pointer is nil.
func (p pointer) IsNil() bool ***REMOVED***
	return p.v.IsNil()
***REMOVED***

// Apply adds an offset to the pointer to derive a new pointer
// to a specified field. The current pointer must be pointing at a struct.
func (p pointer) Apply(f offset) pointer ***REMOVED***
	if f.export != nil ***REMOVED***
		if v := reflect.ValueOf(f.export(p.v.Interface(), f.index)); v.IsValid() ***REMOVED***
			return pointer***REMOVED***v: v***REMOVED***
		***REMOVED***
	***REMOVED***
	return pointer***REMOVED***v: p.v.Elem().Field(f.index).Addr()***REMOVED***
***REMOVED***

// AsValueOf treats p as a pointer to an object of type t and returns the value.
// It is equivalent to reflect.ValueOf(p.AsIfaceOf(t))
func (p pointer) AsValueOf(t reflect.Type) reflect.Value ***REMOVED***
	if got := p.v.Type().Elem(); got != t ***REMOVED***
		panic(fmt.Sprintf("invalid type: got %v, want %v", got, t))
	***REMOVED***
	return p.v
***REMOVED***

// AsIfaceOf treats p as a pointer to an object of type t and returns the value.
// It is equivalent to p.AsValueOf(t).Interface()
func (p pointer) AsIfaceOf(t reflect.Type) interface***REMOVED******REMOVED*** ***REMOVED***
	return p.AsValueOf(t).Interface()
***REMOVED***

func (p pointer) Bool() *bool              ***REMOVED*** return p.v.Interface().(*bool) ***REMOVED***
func (p pointer) BoolPtr() **bool          ***REMOVED*** return p.v.Interface().(**bool) ***REMOVED***
func (p pointer) BoolSlice() *[]bool       ***REMOVED*** return p.v.Interface().(*[]bool) ***REMOVED***
func (p pointer) Int32() *int32            ***REMOVED*** return p.v.Interface().(*int32) ***REMOVED***
func (p pointer) Int32Ptr() **int32        ***REMOVED*** return p.v.Interface().(**int32) ***REMOVED***
func (p pointer) Int32Slice() *[]int32     ***REMOVED*** return p.v.Interface().(*[]int32) ***REMOVED***
func (p pointer) Int64() *int64            ***REMOVED*** return p.v.Interface().(*int64) ***REMOVED***
func (p pointer) Int64Ptr() **int64        ***REMOVED*** return p.v.Interface().(**int64) ***REMOVED***
func (p pointer) Int64Slice() *[]int64     ***REMOVED*** return p.v.Interface().(*[]int64) ***REMOVED***
func (p pointer) Uint32() *uint32          ***REMOVED*** return p.v.Interface().(*uint32) ***REMOVED***
func (p pointer) Uint32Ptr() **uint32      ***REMOVED*** return p.v.Interface().(**uint32) ***REMOVED***
func (p pointer) Uint32Slice() *[]uint32   ***REMOVED*** return p.v.Interface().(*[]uint32) ***REMOVED***
func (p pointer) Uint64() *uint64          ***REMOVED*** return p.v.Interface().(*uint64) ***REMOVED***
func (p pointer) Uint64Ptr() **uint64      ***REMOVED*** return p.v.Interface().(**uint64) ***REMOVED***
func (p pointer) Uint64Slice() *[]uint64   ***REMOVED*** return p.v.Interface().(*[]uint64) ***REMOVED***
func (p pointer) Float32() *float32        ***REMOVED*** return p.v.Interface().(*float32) ***REMOVED***
func (p pointer) Float32Ptr() **float32    ***REMOVED*** return p.v.Interface().(**float32) ***REMOVED***
func (p pointer) Float32Slice() *[]float32 ***REMOVED*** return p.v.Interface().(*[]float32) ***REMOVED***
func (p pointer) Float64() *float64        ***REMOVED*** return p.v.Interface().(*float64) ***REMOVED***
func (p pointer) Float64Ptr() **float64    ***REMOVED*** return p.v.Interface().(**float64) ***REMOVED***
func (p pointer) Float64Slice() *[]float64 ***REMOVED*** return p.v.Interface().(*[]float64) ***REMOVED***
func (p pointer) String() *string          ***REMOVED*** return p.v.Interface().(*string) ***REMOVED***
func (p pointer) StringPtr() **string      ***REMOVED*** return p.v.Interface().(**string) ***REMOVED***
func (p pointer) StringSlice() *[]string   ***REMOVED*** return p.v.Interface().(*[]string) ***REMOVED***
func (p pointer) Bytes() *[]byte           ***REMOVED*** return p.v.Interface().(*[]byte) ***REMOVED***
func (p pointer) BytesSlice() *[][]byte    ***REMOVED*** return p.v.Interface().(*[][]byte) ***REMOVED***
func (p pointer) WeakFields() *weakFields  ***REMOVED*** return (*weakFields)(p.v.Interface().(*WeakFields)) ***REMOVED***
func (p pointer) Extensions() *map[int32]ExtensionField ***REMOVED***
	return p.v.Interface().(*map[int32]ExtensionField)
***REMOVED***

func (p pointer) Elem() pointer ***REMOVED***
	return pointer***REMOVED***v: p.v.Elem()***REMOVED***
***REMOVED***

// PointerSlice copies []*T from p as a new []pointer.
// This behavior differs from the implementation in pointer_unsafe.go.
func (p pointer) PointerSlice() []pointer ***REMOVED***
	// TODO: reconsider this
	if p.v.IsNil() ***REMOVED***
		return nil
	***REMOVED***
	n := p.v.Elem().Len()
	s := make([]pointer, n)
	for i := 0; i < n; i++ ***REMOVED***
		s[i] = pointer***REMOVED***v: p.v.Elem().Index(i)***REMOVED***
	***REMOVED***
	return s
***REMOVED***

// AppendPointerSlice appends v to p, which must be a []*T.
func (p pointer) AppendPointerSlice(v pointer) ***REMOVED***
	sp := p.v.Elem()
	sp.Set(reflect.Append(sp, v.v))
***REMOVED***

// SetPointer sets *p to v.
func (p pointer) SetPointer(v pointer) ***REMOVED***
	p.v.Elem().Set(v.v)
***REMOVED***

func (Export) MessageStateOf(p Pointer) *messageState     ***REMOVED*** panic("not supported") ***REMOVED***
func (ms *messageState) pointer() pointer                 ***REMOVED*** panic("not supported") ***REMOVED***
func (ms *messageState) messageInfo() *MessageInfo        ***REMOVED*** panic("not supported") ***REMOVED***
func (ms *messageState) LoadMessageInfo() *MessageInfo    ***REMOVED*** panic("not supported") ***REMOVED***
func (ms *messageState) StoreMessageInfo(mi *MessageInfo) ***REMOVED*** panic("not supported") ***REMOVED***

type atomicNilMessage struct ***REMOVED***
	once sync.Once
	m    messageReflectWrapper
***REMOVED***

func (m *atomicNilMessage) Init(mi *MessageInfo) *messageReflectWrapper ***REMOVED***
	m.once.Do(func() ***REMOVED***
		m.m.p = pointerOfIface(reflect.Zero(mi.GoReflectType).Interface())
		m.m.mi = mi
	***REMOVED***)
	return &m.m
***REMOVED***
