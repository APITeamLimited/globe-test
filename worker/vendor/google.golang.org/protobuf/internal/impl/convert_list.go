// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func newListConverter(t reflect.Type, fd protoreflect.FieldDescriptor) Converter ***REMOVED***
	switch ***REMOVED***
	case t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice:
		return &listPtrConverter***REMOVED***t, newSingularConverter(t.Elem().Elem(), fd)***REMOVED***
	case t.Kind() == reflect.Slice:
		return &listConverter***REMOVED***t, newSingularConverter(t.Elem(), fd)***REMOVED***
	***REMOVED***
	panic(fmt.Sprintf("invalid Go type %v for field %v", t, fd.FullName()))
***REMOVED***

type listConverter struct ***REMOVED***
	goType reflect.Type // []T
	c      Converter
***REMOVED***

func (c *listConverter) PBValueOf(v reflect.Value) protoreflect.Value ***REMOVED***
	if v.Type() != c.goType ***REMOVED***
		panic(fmt.Sprintf("invalid type: got %v, want %v", v.Type(), c.goType))
	***REMOVED***
	pv := reflect.New(c.goType)
	pv.Elem().Set(v)
	return protoreflect.ValueOfList(&listReflect***REMOVED***pv, c.c***REMOVED***)
***REMOVED***

func (c *listConverter) GoValueOf(v protoreflect.Value) reflect.Value ***REMOVED***
	rv := v.List().(*listReflect).v
	if rv.IsNil() ***REMOVED***
		return reflect.Zero(c.goType)
	***REMOVED***
	return rv.Elem()
***REMOVED***

func (c *listConverter) IsValidPB(v protoreflect.Value) bool ***REMOVED***
	list, ok := v.Interface().(*listReflect)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return list.v.Type().Elem() == c.goType
***REMOVED***

func (c *listConverter) IsValidGo(v reflect.Value) bool ***REMOVED***
	return v.IsValid() && v.Type() == c.goType
***REMOVED***

func (c *listConverter) New() protoreflect.Value ***REMOVED***
	return protoreflect.ValueOfList(&listReflect***REMOVED***reflect.New(c.goType), c.c***REMOVED***)
***REMOVED***

func (c *listConverter) Zero() protoreflect.Value ***REMOVED***
	return protoreflect.ValueOfList(&listReflect***REMOVED***reflect.Zero(reflect.PtrTo(c.goType)), c.c***REMOVED***)
***REMOVED***

type listPtrConverter struct ***REMOVED***
	goType reflect.Type // *[]T
	c      Converter
***REMOVED***

func (c *listPtrConverter) PBValueOf(v reflect.Value) protoreflect.Value ***REMOVED***
	if v.Type() != c.goType ***REMOVED***
		panic(fmt.Sprintf("invalid type: got %v, want %v", v.Type(), c.goType))
	***REMOVED***
	return protoreflect.ValueOfList(&listReflect***REMOVED***v, c.c***REMOVED***)
***REMOVED***

func (c *listPtrConverter) GoValueOf(v protoreflect.Value) reflect.Value ***REMOVED***
	return v.List().(*listReflect).v
***REMOVED***

func (c *listPtrConverter) IsValidPB(v protoreflect.Value) bool ***REMOVED***
	list, ok := v.Interface().(*listReflect)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return list.v.Type() == c.goType
***REMOVED***

func (c *listPtrConverter) IsValidGo(v reflect.Value) bool ***REMOVED***
	return v.IsValid() && v.Type() == c.goType
***REMOVED***

func (c *listPtrConverter) New() protoreflect.Value ***REMOVED***
	return c.PBValueOf(reflect.New(c.goType.Elem()))
***REMOVED***

func (c *listPtrConverter) Zero() protoreflect.Value ***REMOVED***
	return c.PBValueOf(reflect.Zero(c.goType))
***REMOVED***

type listReflect struct ***REMOVED***
	v    reflect.Value // *[]T
	conv Converter
***REMOVED***

func (ls *listReflect) Len() int ***REMOVED***
	if ls.v.IsNil() ***REMOVED***
		return 0
	***REMOVED***
	return ls.v.Elem().Len()
***REMOVED***
func (ls *listReflect) Get(i int) protoreflect.Value ***REMOVED***
	return ls.conv.PBValueOf(ls.v.Elem().Index(i))
***REMOVED***
func (ls *listReflect) Set(i int, v protoreflect.Value) ***REMOVED***
	ls.v.Elem().Index(i).Set(ls.conv.GoValueOf(v))
***REMOVED***
func (ls *listReflect) Append(v protoreflect.Value) ***REMOVED***
	ls.v.Elem().Set(reflect.Append(ls.v.Elem(), ls.conv.GoValueOf(v)))
***REMOVED***
func (ls *listReflect) AppendMutable() protoreflect.Value ***REMOVED***
	if _, ok := ls.conv.(*messageConverter); !ok ***REMOVED***
		panic("invalid AppendMutable on list with non-message type")
	***REMOVED***
	v := ls.NewElement()
	ls.Append(v)
	return v
***REMOVED***
func (ls *listReflect) Truncate(i int) ***REMOVED***
	ls.v.Elem().Set(ls.v.Elem().Slice(0, i))
***REMOVED***
func (ls *listReflect) NewElement() protoreflect.Value ***REMOVED***
	return ls.conv.New()
***REMOVED***
func (ls *listReflect) IsValid() bool ***REMOVED***
	return !ls.v.IsNil()
***REMOVED***
func (ls *listReflect) protoUnwrap() interface***REMOVED******REMOVED*** ***REMOVED***
	return ls.v.Interface()
***REMOVED***
