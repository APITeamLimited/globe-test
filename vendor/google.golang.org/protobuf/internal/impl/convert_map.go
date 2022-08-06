// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type mapConverter struct ***REMOVED***
	goType           reflect.Type // map[K]V
	keyConv, valConv Converter
***REMOVED***

func newMapConverter(t reflect.Type, fd protoreflect.FieldDescriptor) *mapConverter ***REMOVED***
	if t.Kind() != reflect.Map ***REMOVED***
		panic(fmt.Sprintf("invalid Go type %v for field %v", t, fd.FullName()))
	***REMOVED***
	return &mapConverter***REMOVED***
		goType:  t,
		keyConv: newSingularConverter(t.Key(), fd.MapKey()),
		valConv: newSingularConverter(t.Elem(), fd.MapValue()),
	***REMOVED***
***REMOVED***

func (c *mapConverter) PBValueOf(v reflect.Value) protoreflect.Value ***REMOVED***
	if v.Type() != c.goType ***REMOVED***
		panic(fmt.Sprintf("invalid type: got %v, want %v", v.Type(), c.goType))
	***REMOVED***
	return protoreflect.ValueOfMap(&mapReflect***REMOVED***v, c.keyConv, c.valConv***REMOVED***)
***REMOVED***

func (c *mapConverter) GoValueOf(v protoreflect.Value) reflect.Value ***REMOVED***
	return v.Map().(*mapReflect).v
***REMOVED***

func (c *mapConverter) IsValidPB(v protoreflect.Value) bool ***REMOVED***
	mapv, ok := v.Interface().(*mapReflect)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return mapv.v.Type() == c.goType
***REMOVED***

func (c *mapConverter) IsValidGo(v reflect.Value) bool ***REMOVED***
	return v.IsValid() && v.Type() == c.goType
***REMOVED***

func (c *mapConverter) New() protoreflect.Value ***REMOVED***
	return c.PBValueOf(reflect.MakeMap(c.goType))
***REMOVED***

func (c *mapConverter) Zero() protoreflect.Value ***REMOVED***
	return c.PBValueOf(reflect.Zero(c.goType))
***REMOVED***

type mapReflect struct ***REMOVED***
	v       reflect.Value // map[K]V
	keyConv Converter
	valConv Converter
***REMOVED***

func (ms *mapReflect) Len() int ***REMOVED***
	return ms.v.Len()
***REMOVED***
func (ms *mapReflect) Has(k protoreflect.MapKey) bool ***REMOVED***
	rk := ms.keyConv.GoValueOf(k.Value())
	rv := ms.v.MapIndex(rk)
	return rv.IsValid()
***REMOVED***
func (ms *mapReflect) Get(k protoreflect.MapKey) protoreflect.Value ***REMOVED***
	rk := ms.keyConv.GoValueOf(k.Value())
	rv := ms.v.MapIndex(rk)
	if !rv.IsValid() ***REMOVED***
		return protoreflect.Value***REMOVED******REMOVED***
	***REMOVED***
	return ms.valConv.PBValueOf(rv)
***REMOVED***
func (ms *mapReflect) Set(k protoreflect.MapKey, v protoreflect.Value) ***REMOVED***
	rk := ms.keyConv.GoValueOf(k.Value())
	rv := ms.valConv.GoValueOf(v)
	ms.v.SetMapIndex(rk, rv)
***REMOVED***
func (ms *mapReflect) Clear(k protoreflect.MapKey) ***REMOVED***
	rk := ms.keyConv.GoValueOf(k.Value())
	ms.v.SetMapIndex(rk, reflect.Value***REMOVED******REMOVED***)
***REMOVED***
func (ms *mapReflect) Mutable(k protoreflect.MapKey) protoreflect.Value ***REMOVED***
	if _, ok := ms.valConv.(*messageConverter); !ok ***REMOVED***
		panic("invalid Mutable on map with non-message value type")
	***REMOVED***
	v := ms.Get(k)
	if !v.IsValid() ***REMOVED***
		v = ms.NewValue()
		ms.Set(k, v)
	***REMOVED***
	return v
***REMOVED***
func (ms *mapReflect) Range(f func(protoreflect.MapKey, protoreflect.Value) bool) ***REMOVED***
	iter := mapRange(ms.v)
	for iter.Next() ***REMOVED***
		k := ms.keyConv.PBValueOf(iter.Key()).MapKey()
		v := ms.valConv.PBValueOf(iter.Value())
		if !f(k, v) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***
func (ms *mapReflect) NewValue() protoreflect.Value ***REMOVED***
	return ms.valConv.New()
***REMOVED***
func (ms *mapReflect) IsValid() bool ***REMOVED***
	return !ms.v.IsNil()
***REMOVED***
func (ms *mapReflect) protoUnwrap() interface***REMOVED******REMOVED*** ***REMOVED***
	return ms.v.Interface()
***REMOVED***
