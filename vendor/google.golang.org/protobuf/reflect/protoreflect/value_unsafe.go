// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !purego,!appengine

package protoreflect

import (
	"unsafe"

	"google.golang.org/protobuf/internal/pragma"
)

type (
	stringHeader struct ***REMOVED***
		Data unsafe.Pointer
		Len  int
	***REMOVED***
	sliceHeader struct ***REMOVED***
		Data unsafe.Pointer
		Len  int
		Cap  int
	***REMOVED***
	ifaceHeader struct ***REMOVED***
		Type unsafe.Pointer
		Data unsafe.Pointer
	***REMOVED***
)

var (
	nilType     = typeOf(nil)
	boolType    = typeOf(*new(bool))
	int32Type   = typeOf(*new(int32))
	int64Type   = typeOf(*new(int64))
	uint32Type  = typeOf(*new(uint32))
	uint64Type  = typeOf(*new(uint64))
	float32Type = typeOf(*new(float32))
	float64Type = typeOf(*new(float64))
	stringType  = typeOf(*new(string))
	bytesType   = typeOf(*new([]byte))
	enumType    = typeOf(*new(EnumNumber))
)

// typeOf returns a pointer to the Go type information.
// The pointer is comparable and equal if and only if the types are identical.
func typeOf(t interface***REMOVED******REMOVED***) unsafe.Pointer ***REMOVED***
	return (*ifaceHeader)(unsafe.Pointer(&t)).Type
***REMOVED***

// value is a union where only one type can be represented at a time.
// The struct is 24B large on 64-bit systems and requires the minimum storage
// necessary to represent each possible type.
//
// The Go GC needs to be able to scan variables containing pointers.
// As such, pointers and non-pointers cannot be intermixed.
type value struct ***REMOVED***
	pragma.DoNotCompare // 0B

	// typ stores the type of the value as a pointer to the Go type.
	typ unsafe.Pointer // 8B

	// ptr stores the data pointer for a String, Bytes, or interface value.
	ptr unsafe.Pointer // 8B

	// num stores a Bool, Int32, Int64, Uint32, Uint64, Float32, Float64, or
	// Enum value as a raw uint64.
	//
	// It is also used to store the length of a String or Bytes value;
	// the capacity is ignored.
	num uint64 // 8B
***REMOVED***

func valueOfString(v string) Value ***REMOVED***
	p := (*stringHeader)(unsafe.Pointer(&v))
	return Value***REMOVED***typ: stringType, ptr: p.Data, num: uint64(len(v))***REMOVED***
***REMOVED***
func valueOfBytes(v []byte) Value ***REMOVED***
	p := (*sliceHeader)(unsafe.Pointer(&v))
	return Value***REMOVED***typ: bytesType, ptr: p.Data, num: uint64(len(v))***REMOVED***
***REMOVED***
func valueOfIface(v interface***REMOVED******REMOVED***) Value ***REMOVED***
	p := (*ifaceHeader)(unsafe.Pointer(&v))
	return Value***REMOVED***typ: p.Type, ptr: p.Data***REMOVED***
***REMOVED***

func (v Value) getString() (x string) ***REMOVED***
	*(*stringHeader)(unsafe.Pointer(&x)) = stringHeader***REMOVED***Data: v.ptr, Len: int(v.num)***REMOVED***
	return x
***REMOVED***
func (v Value) getBytes() (x []byte) ***REMOVED***
	*(*sliceHeader)(unsafe.Pointer(&x)) = sliceHeader***REMOVED***Data: v.ptr, Len: int(v.num), Cap: int(v.num)***REMOVED***
	return x
***REMOVED***
func (v Value) getIface() (x interface***REMOVED******REMOVED***) ***REMOVED***
	*(*ifaceHeader)(unsafe.Pointer(&x)) = ifaceHeader***REMOVED***Type: v.typ, Data: v.ptr***REMOVED***
	return x
***REMOVED***
