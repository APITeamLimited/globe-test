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

// +build !appengine,!js

// This file contains the implementation of the proto field accesses using package unsafe.

package proto

import (
	"reflect"
	"unsafe"
)

// NOTE: These type_Foo functions would more idiomatically be methods,
// but Go does not allow methods on pointer types, and we must preserve
// some pointer type for the garbage collector. We use these
// funcs with clunky names as our poor approximation to methods.
//
// An alternative would be
//	type structPointer struct ***REMOVED*** p unsafe.Pointer ***REMOVED***
// but that does not registerize as well.

// A structPointer is a pointer to a struct.
type structPointer unsafe.Pointer

// toStructPointer returns a structPointer equivalent to the given reflect value.
func toStructPointer(v reflect.Value) structPointer ***REMOVED***
	return structPointer(unsafe.Pointer(v.Pointer()))
***REMOVED***

// IsNil reports whether p is nil.
func structPointer_IsNil(p structPointer) bool ***REMOVED***
	return p == nil
***REMOVED***

// Interface returns the struct pointer, assumed to have element type t,
// as an interface value.
func structPointer_Interface(p structPointer, t reflect.Type) interface***REMOVED******REMOVED*** ***REMOVED***
	return reflect.NewAt(t, unsafe.Pointer(p)).Interface()
***REMOVED***

// A field identifies a field in a struct, accessible from a structPointer.
// In this implementation, a field is identified by its byte offset from the start of the struct.
type field uintptr

// toField returns a field equivalent to the given reflect field.
func toField(f *reflect.StructField) field ***REMOVED***
	return field(f.Offset)
***REMOVED***

// invalidField is an invalid field identifier.
const invalidField = ^field(0)

// IsValid reports whether the field identifier is valid.
func (f field) IsValid() bool ***REMOVED***
	return f != ^field(0)
***REMOVED***

// Bytes returns the address of a []byte field in the struct.
func structPointer_Bytes(p structPointer, f field) *[]byte ***REMOVED***
	return (*[]byte)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// BytesSlice returns the address of a [][]byte field in the struct.
func structPointer_BytesSlice(p structPointer, f field) *[][]byte ***REMOVED***
	return (*[][]byte)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// Bool returns the address of a *bool field in the struct.
func structPointer_Bool(p structPointer, f field) **bool ***REMOVED***
	return (**bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// BoolVal returns the address of a bool field in the struct.
func structPointer_BoolVal(p structPointer, f field) *bool ***REMOVED***
	return (*bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// BoolSlice returns the address of a []bool field in the struct.
func structPointer_BoolSlice(p structPointer, f field) *[]bool ***REMOVED***
	return (*[]bool)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// String returns the address of a *string field in the struct.
func structPointer_String(p structPointer, f field) **string ***REMOVED***
	return (**string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// StringVal returns the address of a string field in the struct.
func structPointer_StringVal(p structPointer, f field) *string ***REMOVED***
	return (*string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// StringSlice returns the address of a []string field in the struct.
func structPointer_StringSlice(p structPointer, f field) *[]string ***REMOVED***
	return (*[]string)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// ExtMap returns the address of an extension map field in the struct.
func structPointer_Extensions(p structPointer, f field) *XXX_InternalExtensions ***REMOVED***
	return (*XXX_InternalExtensions)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

func structPointer_ExtMap(p structPointer, f field) *map[int32]Extension ***REMOVED***
	return (*map[int32]Extension)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// NewAt returns the reflect.Value for a pointer to a field in the struct.
func structPointer_NewAt(p structPointer, f field, typ reflect.Type) reflect.Value ***REMOVED***
	return reflect.NewAt(typ, unsafe.Pointer(uintptr(p)+uintptr(f)))
***REMOVED***

// SetStructPointer writes a *struct field in the struct.
func structPointer_SetStructPointer(p structPointer, f field, q structPointer) ***REMOVED***
	*(*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f))) = q
***REMOVED***

// GetStructPointer reads a *struct field in the struct.
func structPointer_GetStructPointer(p structPointer, f field) structPointer ***REMOVED***
	return *(*structPointer)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// StructPointerSlice the address of a []*struct field in the struct.
func structPointer_StructPointerSlice(p structPointer, f field) *structPointerSlice ***REMOVED***
	return (*structPointerSlice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// A structPointerSlice represents a slice of pointers to structs (themselves submessages or groups).
type structPointerSlice []structPointer

func (v *structPointerSlice) Len() int                  ***REMOVED*** return len(*v) ***REMOVED***
func (v *structPointerSlice) Index(i int) structPointer ***REMOVED*** return (*v)[i] ***REMOVED***
func (v *structPointerSlice) Append(p structPointer)    ***REMOVED*** *v = append(*v, p) ***REMOVED***

// A word32 is the address of a "pointer to 32-bit value" field.
type word32 **uint32

// IsNil reports whether *v is nil.
func word32_IsNil(p word32) bool ***REMOVED***
	return *p == nil
***REMOVED***

// Set sets *v to point at a newly allocated word set to x.
func word32_Set(p word32, o *Buffer, x uint32) ***REMOVED***
	if len(o.uint32s) == 0 ***REMOVED***
		o.uint32s = make([]uint32, uint32PoolSize)
	***REMOVED***
	o.uint32s[0] = x
	*p = &o.uint32s[0]
	o.uint32s = o.uint32s[1:]
***REMOVED***

// Get gets the value pointed at by *v.
func word32_Get(p word32) uint32 ***REMOVED***
	return **p
***REMOVED***

// Word32 returns the address of a *int32, *uint32, *float32, or *enum field in the struct.
func structPointer_Word32(p structPointer, f field) word32 ***REMOVED***
	return word32((**uint32)(unsafe.Pointer(uintptr(p) + uintptr(f))))
***REMOVED***

// A word32Val is the address of a 32-bit value field.
type word32Val *uint32

// Set sets *p to x.
func word32Val_Set(p word32Val, x uint32) ***REMOVED***
	*p = x
***REMOVED***

// Get gets the value pointed at by p.
func word32Val_Get(p word32Val) uint32 ***REMOVED***
	return *p
***REMOVED***

// Word32Val returns the address of a *int32, *uint32, *float32, or *enum field in the struct.
func structPointer_Word32Val(p structPointer, f field) word32Val ***REMOVED***
	return word32Val((*uint32)(unsafe.Pointer(uintptr(p) + uintptr(f))))
***REMOVED***

// A word32Slice is a slice of 32-bit values.
type word32Slice []uint32

func (v *word32Slice) Append(x uint32)    ***REMOVED*** *v = append(*v, x) ***REMOVED***
func (v *word32Slice) Len() int           ***REMOVED*** return len(*v) ***REMOVED***
func (v *word32Slice) Index(i int) uint32 ***REMOVED*** return (*v)[i] ***REMOVED***

// Word32Slice returns the address of a []int32, []uint32, []float32, or []enum field in the struct.
func structPointer_Word32Slice(p structPointer, f field) *word32Slice ***REMOVED***
	return (*word32Slice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***

// word64 is like word32 but for 64-bit values.
type word64 **uint64

func word64_Set(p word64, o *Buffer, x uint64) ***REMOVED***
	if len(o.uint64s) == 0 ***REMOVED***
		o.uint64s = make([]uint64, uint64PoolSize)
	***REMOVED***
	o.uint64s[0] = x
	*p = &o.uint64s[0]
	o.uint64s = o.uint64s[1:]
***REMOVED***

func word64_IsNil(p word64) bool ***REMOVED***
	return *p == nil
***REMOVED***

func word64_Get(p word64) uint64 ***REMOVED***
	return **p
***REMOVED***

func structPointer_Word64(p structPointer, f field) word64 ***REMOVED***
	return word64((**uint64)(unsafe.Pointer(uintptr(p) + uintptr(f))))
***REMOVED***

// word64Val is like word32Val but for 64-bit values.
type word64Val *uint64

func word64Val_Set(p word64Val, o *Buffer, x uint64) ***REMOVED***
	*p = x
***REMOVED***

func word64Val_Get(p word64Val) uint64 ***REMOVED***
	return *p
***REMOVED***

func structPointer_Word64Val(p structPointer, f field) word64Val ***REMOVED***
	return word64Val((*uint64)(unsafe.Pointer(uintptr(p) + uintptr(f))))
***REMOVED***

// word64Slice is like word32Slice but for 64-bit values.
type word64Slice []uint64

func (v *word64Slice) Append(x uint64)    ***REMOVED*** *v = append(*v, x) ***REMOVED***
func (v *word64Slice) Len() int           ***REMOVED*** return len(*v) ***REMOVED***
func (v *word64Slice) Index(i int) uint64 ***REMOVED*** return (*v)[i] ***REMOVED***

func structPointer_Word64Slice(p structPointer, f field) *word64Slice ***REMOVED***
	return (*word64Slice)(unsafe.Pointer(uintptr(p) + uintptr(f)))
***REMOVED***
