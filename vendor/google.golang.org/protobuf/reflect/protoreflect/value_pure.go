// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build purego appengine

package protoreflect

import "google.golang.org/protobuf/internal/pragma"

type valueType int

const (
	nilType valueType = iota
	boolType
	int32Type
	int64Type
	uint32Type
	uint64Type
	float32Type
	float64Type
	stringType
	bytesType
	enumType
	ifaceType
)

// value is a union where only one type can be represented at a time.
// This uses a distinct field for each type. This is type safe in Go, but
// occupies more memory than necessary (72B).
type value struct ***REMOVED***
	pragma.DoNotCompare // 0B

	typ   valueType   // 8B
	num   uint64      // 8B
	str   string      // 16B
	bin   []byte      // 24B
	iface interface***REMOVED******REMOVED*** // 16B
***REMOVED***

func valueOfString(v string) Value ***REMOVED***
	return Value***REMOVED***typ: stringType, str: v***REMOVED***
***REMOVED***
func valueOfBytes(v []byte) Value ***REMOVED***
	return Value***REMOVED***typ: bytesType, bin: v***REMOVED***
***REMOVED***
func valueOfIface(v interface***REMOVED******REMOVED***) Value ***REMOVED***
	return Value***REMOVED***typ: ifaceType, iface: v***REMOVED***
***REMOVED***

func (v Value) getString() string ***REMOVED***
	return v.str
***REMOVED***
func (v Value) getBytes() []byte ***REMOVED***
	return v.bin
***REMOVED***
func (v Value) getIface() interface***REMOVED******REMOVED*** ***REMOVED***
	return v.iface
***REMOVED***
