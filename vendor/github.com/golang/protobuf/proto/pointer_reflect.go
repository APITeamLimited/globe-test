// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2012 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// +build appengine js

// This file contains an implementation of proto field accesses using package reflect.
// It is slower than the code in pointer_unsafe.go but it avoids package unsafe and can
// be used on App Engine.

package proto

import (
	"math"
	"reflect"
)

// A structPointer is a pointer to a struct.
type structPointer struct ***REMOVED***
	v reflect.Value
***REMOVED***

// toStructPointer returns a structPointer equivalent to the given reflect value.
// The reflect value must itself be a pointer to a struct.
func toStructPointer(v reflect.Value) structPointer ***REMOVED***
	return structPointer***REMOVED***v***REMOVED***
***REMOVED***

// IsNil reports whether p is nil.
func structPointer_IsNil(p structPointer) bool ***REMOVED***
	return p.v.IsNil()
***REMOVED***

// Interface returns the struct pointer as an interface value.
func structPointer_Interface(p structPointer, _ reflect.Type) interface***REMOVED******REMOVED*** ***REMOVED***
	return p.v.Interface()
***REMOVED***

// A field identifies a field in a struct, accessible from a structPointer.
// In this implementation, a field is identified by the sequence of field indices
// passed to reflect's FieldByIndex.
type field []int

// toField returns a field equivalent to the given reflect field.
func toField(f *reflect.StructField) field ***REMOVED***
	return f.Index
***REMOVED***

// invalidField is an invalid field identifier.
var invalidField = field(nil)

// IsValid reports whether the field identifier is valid.
func (f field) IsValid() bool ***REMOVED*** return f != nil ***REMOVED***

// field returns the given field in the struct as a reflect value.
func structPointer_field(p structPointer, f field) reflect.Value ***REMOVED***
	// Special case: an extension map entry with a value of type T
	// passes a *T to the struct-handling code with a zero field,
	// expecting that it will be treated as equivalent to *struct***REMOVED*** X T ***REMOVED***,
	// which has the same memory layout. We have to handle that case
	// specially, because reflect will panic if we call FieldByIndex on a
	// non-struct.
	if f == nil ***REMOVED***
		return p.v.Elem()
	***REMOVED***

	return p.v.Elem().FieldByIndex(f)
***REMOVED***

// ifield returns the given field in the struct as an interface value.
func structPointer_ifield(p structPointer, f field) interface***REMOVED******REMOVED*** ***REMOVED***
	return structPointer_field(p, f).Addr().Interface()
***REMOVED***

// Bytes returns the address of a []byte field in the struct.
func structPointer_Bytes(p structPointer, f field) *[]byte ***REMOVED***
	return structPointer_ifield(p, f).(*[]byte)
***REMOVED***

// BytesSlice returns the address of a [][]byte field in the struct.
func structPointer_BytesSlice(p structPointer, f field) *[][]byte ***REMOVED***
	return structPointer_ifield(p, f).(*[][]byte)
***REMOVED***

// Bool returns the address of a *bool field in the struct.
func structPointer_Bool(p structPointer, f field) **bool ***REMOVED***
	return structPointer_ifield(p, f).(**bool)
***REMOVED***

// BoolVal returns the address of a bool field in the struct.
func structPointer_BoolVal(p structPointer, f field) *bool ***REMOVED***
	return structPointer_ifield(p, f).(*bool)
***REMOVED***

// BoolSlice returns the address of a []bool field in the struct.
func structPointer_BoolSlice(p structPointer, f field) *[]bool ***REMOVED***
	return structPointer_ifield(p, f).(*[]bool)
***REMOVED***

// String returns the address of a *string field in the struct.
func structPointer_String(p structPointer, f field) **string ***REMOVED***
	return structPointer_ifield(p, f).(**string)
***REMOVED***

// StringVal returns the address of a string field in the struct.
func structPointer_StringVal(p structPointer, f field) *string ***REMOVED***
	return structPointer_ifield(p, f).(*string)
***REMOVED***

// StringSlice returns the address of a []string field in the struct.
func structPointer_StringSlice(p structPointer, f field) *[]string ***REMOVED***
	return structPointer_ifield(p, f).(*[]string)
***REMOVED***

// Extensions returns the address of an extension map field in the struct.
func structPointer_Extensions(p structPointer, f field) *XXX_InternalExtensions ***REMOVED***
	return structPointer_ifield(p, f).(*XXX_InternalExtensions)
***REMOVED***

// ExtMap returns the address of an extension map field in the struct.
func structPointer_ExtMap(p structPointer, f field) *map[int32]Extension ***REMOVED***
	return structPointer_ifield(p, f).(*map[int32]Extension)
***REMOVED***

// NewAt returns the reflect.Value for a pointer to a field in the struct.
func structPointer_NewAt(p structPointer, f field, typ reflect.Type) reflect.Value ***REMOVED***
	return structPointer_field(p, f).Addr()
***REMOVED***

// SetStructPointer writes a *struct field in the struct.
func structPointer_SetStructPointer(p structPointer, f field, q structPointer) ***REMOVED***
	structPointer_field(p, f).Set(q.v)
***REMOVED***

// GetStructPointer reads a *struct field in the struct.
func structPointer_GetStructPointer(p structPointer, f field) structPointer ***REMOVED***
	return structPointer***REMOVED***structPointer_field(p, f)***REMOVED***
***REMOVED***

// StructPointerSlice the address of a []*struct field in the struct.
func structPointer_StructPointerSlice(p structPointer, f field) structPointerSlice ***REMOVED***
	return structPointerSlice***REMOVED***structPointer_field(p, f)***REMOVED***
***REMOVED***

// A structPointerSlice represents the address of a slice of pointers to structs
// (themselves messages or groups). That is, v.Type() is *[]*struct***REMOVED***...***REMOVED***.
type structPointerSlice struct ***REMOVED***
	v reflect.Value
***REMOVED***

func (p structPointerSlice) Len() int                  ***REMOVED*** return p.v.Len() ***REMOVED***
func (p structPointerSlice) Index(i int) structPointer ***REMOVED*** return structPointer***REMOVED***p.v.Index(i)***REMOVED*** ***REMOVED***
func (p structPointerSlice) Append(q structPointer) ***REMOVED***
	p.v.Set(reflect.Append(p.v, q.v))
***REMOVED***

var (
	int32Type   = reflect.TypeOf(int32(0))
	uint32Type  = reflect.TypeOf(uint32(0))
	float32Type = reflect.TypeOf(float32(0))
	int64Type   = reflect.TypeOf(int64(0))
	uint64Type  = reflect.TypeOf(uint64(0))
	float64Type = reflect.TypeOf(float64(0))
)

// A word32 represents a field of type *int32, *uint32, *float32, or *enum.
// That is, v.Type() is *int32, *uint32, *float32, or *enum and v is assignable.
type word32 struct ***REMOVED***
	v reflect.Value
***REMOVED***

// IsNil reports whether p is nil.
func word32_IsNil(p word32) bool ***REMOVED***
	return p.v.IsNil()
***REMOVED***

// Set sets p to point at a newly allocated word with bits set to x.
func word32_Set(p word32, o *Buffer, x uint32) ***REMOVED***
	t := p.v.Type().Elem()
	switch t ***REMOVED***
	case int32Type:
		if len(o.int32s) == 0 ***REMOVED***
			o.int32s = make([]int32, uint32PoolSize)
		***REMOVED***
		o.int32s[0] = int32(x)
		p.v.Set(reflect.ValueOf(&o.int32s[0]))
		o.int32s = o.int32s[1:]
		return
	case uint32Type:
		if len(o.uint32s) == 0 ***REMOVED***
			o.uint32s = make([]uint32, uint32PoolSize)
		***REMOVED***
		o.uint32s[0] = x
		p.v.Set(reflect.ValueOf(&o.uint32s[0]))
		o.uint32s = o.uint32s[1:]
		return
	case float32Type:
		if len(o.float32s) == 0 ***REMOVED***
			o.float32s = make([]float32, uint32PoolSize)
		***REMOVED***
		o.float32s[0] = math.Float32frombits(x)
		p.v.Set(reflect.ValueOf(&o.float32s[0]))
		o.float32s = o.float32s[1:]
		return
	***REMOVED***

	// must be enum
	p.v.Set(reflect.New(t))
	p.v.Elem().SetInt(int64(int32(x)))
***REMOVED***

// Get gets the bits pointed at by p, as a uint32.
func word32_Get(p word32) uint32 ***REMOVED***
	elem := p.v.Elem()
	switch elem.Kind() ***REMOVED***
	case reflect.Int32:
		return uint32(elem.Int())
	case reflect.Uint32:
		return uint32(elem.Uint())
	case reflect.Float32:
		return math.Float32bits(float32(elem.Float()))
	***REMOVED***
	panic("unreachable")
***REMOVED***

// Word32 returns a reference to a *int32, *uint32, *float32, or *enum field in the struct.
func structPointer_Word32(p structPointer, f field) word32 ***REMOVED***
	return word32***REMOVED***structPointer_field(p, f)***REMOVED***
***REMOVED***

// A word32Val represents a field of type int32, uint32, float32, or enum.
// That is, v.Type() is int32, uint32, float32, or enum and v is assignable.
type word32Val struct ***REMOVED***
	v reflect.Value
***REMOVED***

// Set sets *p to x.
func word32Val_Set(p word32Val, x uint32) ***REMOVED***
	switch p.v.Type() ***REMOVED***
	case int32Type:
		p.v.SetInt(int64(x))
		return
	case uint32Type:
		p.v.SetUint(uint64(x))
		return
	case float32Type:
		p.v.SetFloat(float64(math.Float32frombits(x)))
		return
	***REMOVED***

	// must be enum
	p.v.SetInt(int64(int32(x)))
***REMOVED***

// Get gets the bits pointed at by p, as a uint32.
func word32Val_Get(p word32Val) uint32 ***REMOVED***
	elem := p.v
	switch elem.Kind() ***REMOVED***
	case reflect.Int32:
		return uint32(elem.Int())
	case reflect.Uint32:
		return uint32(elem.Uint())
	case reflect.Float32:
		return math.Float32bits(float32(elem.Float()))
	***REMOVED***
	panic("unreachable")
***REMOVED***

// Word32Val returns a reference to a int32, uint32, float32, or enum field in the struct.
func structPointer_Word32Val(p structPointer, f field) word32Val ***REMOVED***
	return word32Val***REMOVED***structPointer_field(p, f)***REMOVED***
***REMOVED***

// A word32Slice is a slice of 32-bit values.
// That is, v.Type() is []int32, []uint32, []float32, or []enum.
type word32Slice struct ***REMOVED***
	v reflect.Value
***REMOVED***

func (p word32Slice) Append(x uint32) ***REMOVED***
	n, m := p.v.Len(), p.v.Cap()
	if n < m ***REMOVED***
		p.v.SetLen(n + 1)
	***REMOVED*** else ***REMOVED***
		t := p.v.Type().Elem()
		p.v.Set(reflect.Append(p.v, reflect.Zero(t)))
	***REMOVED***
	elem := p.v.Index(n)
	switch elem.Kind() ***REMOVED***
	case reflect.Int32:
		elem.SetInt(int64(int32(x)))
	case reflect.Uint32:
		elem.SetUint(uint64(x))
	case reflect.Float32:
		elem.SetFloat(float64(math.Float32frombits(x)))
	***REMOVED***
***REMOVED***

func (p word32Slice) Len() int ***REMOVED***
	return p.v.Len()
***REMOVED***

func (p word32Slice) Index(i int) uint32 ***REMOVED***
	elem := p.v.Index(i)
	switch elem.Kind() ***REMOVED***
	case reflect.Int32:
		return uint32(elem.Int())
	case reflect.Uint32:
		return uint32(elem.Uint())
	case reflect.Float32:
		return math.Float32bits(float32(elem.Float()))
	***REMOVED***
	panic("unreachable")
***REMOVED***

// Word32Slice returns a reference to a []int32, []uint32, []float32, or []enum field in the struct.
func structPointer_Word32Slice(p structPointer, f field) word32Slice ***REMOVED***
	return word32Slice***REMOVED***structPointer_field(p, f)***REMOVED***
***REMOVED***

// word64 is like word32 but for 64-bit values.
type word64 struct ***REMOVED***
	v reflect.Value
***REMOVED***

func word64_Set(p word64, o *Buffer, x uint64) ***REMOVED***
	t := p.v.Type().Elem()
	switch t ***REMOVED***
	case int64Type:
		if len(o.int64s) == 0 ***REMOVED***
			o.int64s = make([]int64, uint64PoolSize)
		***REMOVED***
		o.int64s[0] = int64(x)
		p.v.Set(reflect.ValueOf(&o.int64s[0]))
		o.int64s = o.int64s[1:]
		return
	case uint64Type:
		if len(o.uint64s) == 0 ***REMOVED***
			o.uint64s = make([]uint64, uint64PoolSize)
		***REMOVED***
		o.uint64s[0] = x
		p.v.Set(reflect.ValueOf(&o.uint64s[0]))
		o.uint64s = o.uint64s[1:]
		return
	case float64Type:
		if len(o.float64s) == 0 ***REMOVED***
			o.float64s = make([]float64, uint64PoolSize)
		***REMOVED***
		o.float64s[0] = math.Float64frombits(x)
		p.v.Set(reflect.ValueOf(&o.float64s[0]))
		o.float64s = o.float64s[1:]
		return
	***REMOVED***
	panic("unreachable")
***REMOVED***

func word64_IsNil(p word64) bool ***REMOVED***
	return p.v.IsNil()
***REMOVED***

func word64_Get(p word64) uint64 ***REMOVED***
	elem := p.v.Elem()
	switch elem.Kind() ***REMOVED***
	case reflect.Int64:
		return uint64(elem.Int())
	case reflect.Uint64:
		return elem.Uint()
	case reflect.Float64:
		return math.Float64bits(elem.Float())
	***REMOVED***
	panic("unreachable")
***REMOVED***

func structPointer_Word64(p structPointer, f field) word64 ***REMOVED***
	return word64***REMOVED***structPointer_field(p, f)***REMOVED***
***REMOVED***

// word64Val is like word32Val but for 64-bit values.
type word64Val struct ***REMOVED***
	v reflect.Value
***REMOVED***

func word64Val_Set(p word64Val, o *Buffer, x uint64) ***REMOVED***
	switch p.v.Type() ***REMOVED***
	case int64Type:
		p.v.SetInt(int64(x))
		return
	case uint64Type:
		p.v.SetUint(x)
		return
	case float64Type:
		p.v.SetFloat(math.Float64frombits(x))
		return
	***REMOVED***
	panic("unreachable")
***REMOVED***

func word64Val_Get(p word64Val) uint64 ***REMOVED***
	elem := p.v
	switch elem.Kind() ***REMOVED***
	case reflect.Int64:
		return uint64(elem.Int())
	case reflect.Uint64:
		return elem.Uint()
	case reflect.Float64:
		return math.Float64bits(elem.Float())
	***REMOVED***
	panic("unreachable")
***REMOVED***

func structPointer_Word64Val(p structPointer, f field) word64Val ***REMOVED***
	return word64Val***REMOVED***structPointer_field(p, f)***REMOVED***
***REMOVED***

type word64Slice struct ***REMOVED***
	v reflect.Value
***REMOVED***

func (p word64Slice) Append(x uint64) ***REMOVED***
	n, m := p.v.Len(), p.v.Cap()
	if n < m ***REMOVED***
		p.v.SetLen(n + 1)
	***REMOVED*** else ***REMOVED***
		t := p.v.Type().Elem()
		p.v.Set(reflect.Append(p.v, reflect.Zero(t)))
	***REMOVED***
	elem := p.v.Index(n)
	switch elem.Kind() ***REMOVED***
	case reflect.Int64:
		elem.SetInt(int64(int64(x)))
	case reflect.Uint64:
		elem.SetUint(uint64(x))
	case reflect.Float64:
		elem.SetFloat(float64(math.Float64frombits(x)))
	***REMOVED***
***REMOVED***

func (p word64Slice) Len() int ***REMOVED***
	return p.v.Len()
***REMOVED***

func (p word64Slice) Index(i int) uint64 ***REMOVED***
	elem := p.v.Index(i)
	switch elem.Kind() ***REMOVED***
	case reflect.Int64:
		return uint64(elem.Int())
	case reflect.Uint64:
		return uint64(elem.Uint())
	case reflect.Float64:
		return math.Float64bits(float64(elem.Float()))
	***REMOVED***
	panic("unreachable")
***REMOVED***

func structPointer_Word64Slice(p structPointer, f field) word64Slice ***REMOVED***
	return word64Slice***REMOVED***structPointer_field(p, f)***REMOVED***
***REMOVED***
