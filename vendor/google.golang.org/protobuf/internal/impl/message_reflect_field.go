// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"math"
	"reflect"
	"sync"

	"google.golang.org/protobuf/internal/flags"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	preg "google.golang.org/protobuf/reflect/protoregistry"
)

type fieldInfo struct ***REMOVED***
	fieldDesc pref.FieldDescriptor

	// These fields are used for protobuf reflection support.
	has        func(pointer) bool
	clear      func(pointer)
	get        func(pointer) pref.Value
	set        func(pointer, pref.Value)
	mutable    func(pointer) pref.Value
	newMessage func() pref.Message
	newField   func() pref.Value
***REMOVED***

func fieldInfoForMissing(fd pref.FieldDescriptor) fieldInfo ***REMOVED***
	// This never occurs for generated message types.
	// It implies that a hand-crafted type has missing Go fields
	// for specific protobuf message fields.
	return fieldInfo***REMOVED***
		fieldDesc: fd,
		has: func(p pointer) bool ***REMOVED***
			return false
		***REMOVED***,
		clear: func(p pointer) ***REMOVED***
			panic("missing Go struct field for " + string(fd.FullName()))
		***REMOVED***,
		get: func(p pointer) pref.Value ***REMOVED***
			return fd.Default()
		***REMOVED***,
		set: func(p pointer, v pref.Value) ***REMOVED***
			panic("missing Go struct field for " + string(fd.FullName()))
		***REMOVED***,
		mutable: func(p pointer) pref.Value ***REMOVED***
			panic("missing Go struct field for " + string(fd.FullName()))
		***REMOVED***,
		newMessage: func() pref.Message ***REMOVED***
			panic("missing Go struct field for " + string(fd.FullName()))
		***REMOVED***,
		newField: func() pref.Value ***REMOVED***
			if v := fd.Default(); v.IsValid() ***REMOVED***
				return v
			***REMOVED***
			panic("missing Go struct field for " + string(fd.FullName()))
		***REMOVED***,
	***REMOVED***
***REMOVED***

func fieldInfoForOneof(fd pref.FieldDescriptor, fs reflect.StructField, x exporter, ot reflect.Type) fieldInfo ***REMOVED***
	ft := fs.Type
	if ft.Kind() != reflect.Interface ***REMOVED***
		panic(fmt.Sprintf("field %v has invalid type: got %v, want interface kind", fd.FullName(), ft))
	***REMOVED***
	if ot.Kind() != reflect.Struct ***REMOVED***
		panic(fmt.Sprintf("field %v has invalid type: got %v, want struct kind", fd.FullName(), ot))
	***REMOVED***
	if !reflect.PtrTo(ot).Implements(ft) ***REMOVED***
		panic(fmt.Sprintf("field %v has invalid type: %v does not implement %v", fd.FullName(), ot, ft))
	***REMOVED***
	conv := NewConverter(ot.Field(0).Type, fd)
	isMessage := fd.Message() != nil

	// TODO: Implement unsafe fast path?
	fieldOffset := offsetOf(fs, x)
	return fieldInfo***REMOVED***
		// NOTE: The logic below intentionally assumes that oneof fields are
		// well-formatted. That is, the oneof interface never contains a
		// typed nil pointer to one of the wrapper structs.

		fieldDesc: fd,
		has: func(p pointer) bool ***REMOVED***
			if p.IsNil() ***REMOVED***
				return false
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if rv.IsNil() || rv.Elem().Type().Elem() != ot || rv.Elem().IsNil() ***REMOVED***
				return false
			***REMOVED***
			return true
		***REMOVED***,
		clear: func(p pointer) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if rv.IsNil() || rv.Elem().Type().Elem() != ot ***REMOVED***
				// NOTE: We intentionally don't check for rv.Elem().IsNil()
				// so that (*OneofWrapperType)(nil) gets cleared to nil.
				return
			***REMOVED***
			rv.Set(reflect.Zero(rv.Type()))
		***REMOVED***,
		get: func(p pointer) pref.Value ***REMOVED***
			if p.IsNil() ***REMOVED***
				return conv.Zero()
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if rv.IsNil() || rv.Elem().Type().Elem() != ot || rv.Elem().IsNil() ***REMOVED***
				return conv.Zero()
			***REMOVED***
			rv = rv.Elem().Elem().Field(0)
			return conv.PBValueOf(rv)
		***REMOVED***,
		set: func(p pointer, v pref.Value) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if rv.IsNil() || rv.Elem().Type().Elem() != ot || rv.Elem().IsNil() ***REMOVED***
				rv.Set(reflect.New(ot))
			***REMOVED***
			rv = rv.Elem().Elem().Field(0)
			rv.Set(conv.GoValueOf(v))
		***REMOVED***,
		mutable: func(p pointer) pref.Value ***REMOVED***
			if !isMessage ***REMOVED***
				panic(fmt.Sprintf("field %v with invalid Mutable call on field with non-composite type", fd.FullName()))
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if rv.IsNil() || rv.Elem().Type().Elem() != ot || rv.Elem().IsNil() ***REMOVED***
				rv.Set(reflect.New(ot))
			***REMOVED***
			rv = rv.Elem().Elem().Field(0)
			if rv.Kind() == reflect.Ptr && rv.IsNil() ***REMOVED***
				rv.Set(conv.GoValueOf(pref.ValueOfMessage(conv.New().Message())))
			***REMOVED***
			return conv.PBValueOf(rv)
		***REMOVED***,
		newMessage: func() pref.Message ***REMOVED***
			return conv.New().Message()
		***REMOVED***,
		newField: func() pref.Value ***REMOVED***
			return conv.New()
		***REMOVED***,
	***REMOVED***
***REMOVED***

func fieldInfoForMap(fd pref.FieldDescriptor, fs reflect.StructField, x exporter) fieldInfo ***REMOVED***
	ft := fs.Type
	if ft.Kind() != reflect.Map ***REMOVED***
		panic(fmt.Sprintf("field %v has invalid type: got %v, want map kind", fd.FullName(), ft))
	***REMOVED***
	conv := NewConverter(ft, fd)

	// TODO: Implement unsafe fast path?
	fieldOffset := offsetOf(fs, x)
	return fieldInfo***REMOVED***
		fieldDesc: fd,
		has: func(p pointer) bool ***REMOVED***
			if p.IsNil() ***REMOVED***
				return false
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			return rv.Len() > 0
		***REMOVED***,
		clear: func(p pointer) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			rv.Set(reflect.Zero(rv.Type()))
		***REMOVED***,
		get: func(p pointer) pref.Value ***REMOVED***
			if p.IsNil() ***REMOVED***
				return conv.Zero()
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if rv.Len() == 0 ***REMOVED***
				return conv.Zero()
			***REMOVED***
			return conv.PBValueOf(rv)
		***REMOVED***,
		set: func(p pointer, v pref.Value) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			pv := conv.GoValueOf(v)
			if pv.IsNil() ***REMOVED***
				panic(fmt.Sprintf("map field %v cannot be set with read-only value", fd.FullName()))
			***REMOVED***
			rv.Set(pv)
		***REMOVED***,
		mutable: func(p pointer) pref.Value ***REMOVED***
			v := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if v.IsNil() ***REMOVED***
				v.Set(reflect.MakeMap(fs.Type))
			***REMOVED***
			return conv.PBValueOf(v)
		***REMOVED***,
		newField: func() pref.Value ***REMOVED***
			return conv.New()
		***REMOVED***,
	***REMOVED***
***REMOVED***

func fieldInfoForList(fd pref.FieldDescriptor, fs reflect.StructField, x exporter) fieldInfo ***REMOVED***
	ft := fs.Type
	if ft.Kind() != reflect.Slice ***REMOVED***
		panic(fmt.Sprintf("field %v has invalid type: got %v, want slice kind", fd.FullName(), ft))
	***REMOVED***
	conv := NewConverter(reflect.PtrTo(ft), fd)

	// TODO: Implement unsafe fast path?
	fieldOffset := offsetOf(fs, x)
	return fieldInfo***REMOVED***
		fieldDesc: fd,
		has: func(p pointer) bool ***REMOVED***
			if p.IsNil() ***REMOVED***
				return false
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			return rv.Len() > 0
		***REMOVED***,
		clear: func(p pointer) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			rv.Set(reflect.Zero(rv.Type()))
		***REMOVED***,
		get: func(p pointer) pref.Value ***REMOVED***
			if p.IsNil() ***REMOVED***
				return conv.Zero()
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type)
			if rv.Elem().Len() == 0 ***REMOVED***
				return conv.Zero()
			***REMOVED***
			return conv.PBValueOf(rv)
		***REMOVED***,
		set: func(p pointer, v pref.Value) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			pv := conv.GoValueOf(v)
			if pv.IsNil() ***REMOVED***
				panic(fmt.Sprintf("list field %v cannot be set with read-only value", fd.FullName()))
			***REMOVED***
			rv.Set(pv.Elem())
		***REMOVED***,
		mutable: func(p pointer) pref.Value ***REMOVED***
			v := p.Apply(fieldOffset).AsValueOf(fs.Type)
			return conv.PBValueOf(v)
		***REMOVED***,
		newField: func() pref.Value ***REMOVED***
			return conv.New()
		***REMOVED***,
	***REMOVED***
***REMOVED***

var (
	nilBytes   = reflect.ValueOf([]byte(nil))
	emptyBytes = reflect.ValueOf([]byte***REMOVED******REMOVED***)
)

func fieldInfoForScalar(fd pref.FieldDescriptor, fs reflect.StructField, x exporter) fieldInfo ***REMOVED***
	ft := fs.Type
	nullable := fd.HasPresence()
	isBytes := ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Uint8
	if nullable ***REMOVED***
		if ft.Kind() != reflect.Ptr && ft.Kind() != reflect.Slice ***REMOVED***
			// This never occurs for generated message types.
			// Despite the protobuf type system specifying presence,
			// the Go field type cannot represent it.
			nullable = false
		***REMOVED***
		if ft.Kind() == reflect.Ptr ***REMOVED***
			ft = ft.Elem()
		***REMOVED***
	***REMOVED***
	conv := NewConverter(ft, fd)

	// TODO: Implement unsafe fast path?
	fieldOffset := offsetOf(fs, x)
	return fieldInfo***REMOVED***
		fieldDesc: fd,
		has: func(p pointer) bool ***REMOVED***
			if p.IsNil() ***REMOVED***
				return false
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if nullable ***REMOVED***
				return !rv.IsNil()
			***REMOVED***
			switch rv.Kind() ***REMOVED***
			case reflect.Bool:
				return rv.Bool()
			case reflect.Int32, reflect.Int64:
				return rv.Int() != 0
			case reflect.Uint32, reflect.Uint64:
				return rv.Uint() != 0
			case reflect.Float32, reflect.Float64:
				return rv.Float() != 0 || math.Signbit(rv.Float())
			case reflect.String, reflect.Slice:
				return rv.Len() > 0
			default:
				panic(fmt.Sprintf("field %v has invalid type: %v", fd.FullName(), rv.Type())) // should never happen
			***REMOVED***
		***REMOVED***,
		clear: func(p pointer) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			rv.Set(reflect.Zero(rv.Type()))
		***REMOVED***,
		get: func(p pointer) pref.Value ***REMOVED***
			if p.IsNil() ***REMOVED***
				return conv.Zero()
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if nullable ***REMOVED***
				if rv.IsNil() ***REMOVED***
					return conv.Zero()
				***REMOVED***
				if rv.Kind() == reflect.Ptr ***REMOVED***
					rv = rv.Elem()
				***REMOVED***
			***REMOVED***
			return conv.PBValueOf(rv)
		***REMOVED***,
		set: func(p pointer, v pref.Value) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if nullable && rv.Kind() == reflect.Ptr ***REMOVED***
				if rv.IsNil() ***REMOVED***
					rv.Set(reflect.New(ft))
				***REMOVED***
				rv = rv.Elem()
			***REMOVED***
			rv.Set(conv.GoValueOf(v))
			if isBytes && rv.Len() == 0 ***REMOVED***
				if nullable ***REMOVED***
					rv.Set(emptyBytes) // preserve presence
				***REMOVED*** else ***REMOVED***
					rv.Set(nilBytes) // do not preserve presence
				***REMOVED***
			***REMOVED***
		***REMOVED***,
		newField: func() pref.Value ***REMOVED***
			return conv.New()
		***REMOVED***,
	***REMOVED***
***REMOVED***

func fieldInfoForWeakMessage(fd pref.FieldDescriptor, weakOffset offset) fieldInfo ***REMOVED***
	if !flags.ProtoLegacy ***REMOVED***
		panic("no support for proto1 weak fields")
	***REMOVED***

	var once sync.Once
	var messageType pref.MessageType
	lazyInit := func() ***REMOVED***
		once.Do(func() ***REMOVED***
			messageName := fd.Message().FullName()
			messageType, _ = preg.GlobalTypes.FindMessageByName(messageName)
			if messageType == nil ***REMOVED***
				panic(fmt.Sprintf("weak message %v for field %v is not linked in", messageName, fd.FullName()))
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	num := fd.Number()
	return fieldInfo***REMOVED***
		fieldDesc: fd,
		has: func(p pointer) bool ***REMOVED***
			if p.IsNil() ***REMOVED***
				return false
			***REMOVED***
			_, ok := p.Apply(weakOffset).WeakFields().get(num)
			return ok
		***REMOVED***,
		clear: func(p pointer) ***REMOVED***
			p.Apply(weakOffset).WeakFields().clear(num)
		***REMOVED***,
		get: func(p pointer) pref.Value ***REMOVED***
			lazyInit()
			if p.IsNil() ***REMOVED***
				return pref.ValueOfMessage(messageType.Zero())
			***REMOVED***
			m, ok := p.Apply(weakOffset).WeakFields().get(num)
			if !ok ***REMOVED***
				return pref.ValueOfMessage(messageType.Zero())
			***REMOVED***
			return pref.ValueOfMessage(m.ProtoReflect())
		***REMOVED***,
		set: func(p pointer, v pref.Value) ***REMOVED***
			lazyInit()
			m := v.Message()
			if m.Descriptor() != messageType.Descriptor() ***REMOVED***
				if got, want := m.Descriptor().FullName(), messageType.Descriptor().FullName(); got != want ***REMOVED***
					panic(fmt.Sprintf("field %v has mismatching message descriptor: got %v, want %v", fd.FullName(), got, want))
				***REMOVED***
				panic(fmt.Sprintf("field %v has mismatching message descriptor: %v", fd.FullName(), m.Descriptor().FullName()))
			***REMOVED***
			p.Apply(weakOffset).WeakFields().set(num, m.Interface())
		***REMOVED***,
		mutable: func(p pointer) pref.Value ***REMOVED***
			lazyInit()
			fs := p.Apply(weakOffset).WeakFields()
			m, ok := fs.get(num)
			if !ok ***REMOVED***
				m = messageType.New().Interface()
				fs.set(num, m)
			***REMOVED***
			return pref.ValueOfMessage(m.ProtoReflect())
		***REMOVED***,
		newMessage: func() pref.Message ***REMOVED***
			lazyInit()
			return messageType.New()
		***REMOVED***,
		newField: func() pref.Value ***REMOVED***
			lazyInit()
			return pref.ValueOfMessage(messageType.New())
		***REMOVED***,
	***REMOVED***
***REMOVED***

func fieldInfoForMessage(fd pref.FieldDescriptor, fs reflect.StructField, x exporter) fieldInfo ***REMOVED***
	ft := fs.Type
	conv := NewConverter(ft, fd)

	// TODO: Implement unsafe fast path?
	fieldOffset := offsetOf(fs, x)
	return fieldInfo***REMOVED***
		fieldDesc: fd,
		has: func(p pointer) bool ***REMOVED***
			if p.IsNil() ***REMOVED***
				return false
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if fs.Type.Kind() != reflect.Ptr ***REMOVED***
				return !isZero(rv)
			***REMOVED***
			return !rv.IsNil()
		***REMOVED***,
		clear: func(p pointer) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			rv.Set(reflect.Zero(rv.Type()))
		***REMOVED***,
		get: func(p pointer) pref.Value ***REMOVED***
			if p.IsNil() ***REMOVED***
				return conv.Zero()
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			return conv.PBValueOf(rv)
		***REMOVED***,
		set: func(p pointer, v pref.Value) ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			rv.Set(conv.GoValueOf(v))
			if fs.Type.Kind() == reflect.Ptr && rv.IsNil() ***REMOVED***
				panic(fmt.Sprintf("field %v has invalid nil pointer", fd.FullName()))
			***REMOVED***
		***REMOVED***,
		mutable: func(p pointer) pref.Value ***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if fs.Type.Kind() == reflect.Ptr && rv.IsNil() ***REMOVED***
				rv.Set(conv.GoValueOf(conv.New()))
			***REMOVED***
			return conv.PBValueOf(rv)
		***REMOVED***,
		newMessage: func() pref.Message ***REMOVED***
			return conv.New().Message()
		***REMOVED***,
		newField: func() pref.Value ***REMOVED***
			return conv.New()
		***REMOVED***,
	***REMOVED***
***REMOVED***

type oneofInfo struct ***REMOVED***
	oneofDesc pref.OneofDescriptor
	which     func(pointer) pref.FieldNumber
***REMOVED***

func makeOneofInfo(od pref.OneofDescriptor, si structInfo, x exporter) *oneofInfo ***REMOVED***
	oi := &oneofInfo***REMOVED***oneofDesc: od***REMOVED***
	if od.IsSynthetic() ***REMOVED***
		fs := si.fieldsByNumber[od.Fields().Get(0).Number()]
		fieldOffset := offsetOf(fs, x)
		oi.which = func(p pointer) pref.FieldNumber ***REMOVED***
			if p.IsNil() ***REMOVED***
				return 0
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if rv.IsNil() ***REMOVED*** // valid on either *T or []byte
				return 0
			***REMOVED***
			return od.Fields().Get(0).Number()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fs := si.oneofsByName[od.Name()]
		fieldOffset := offsetOf(fs, x)
		oi.which = func(p pointer) pref.FieldNumber ***REMOVED***
			if p.IsNil() ***REMOVED***
				return 0
			***REMOVED***
			rv := p.Apply(fieldOffset).AsValueOf(fs.Type).Elem()
			if rv.IsNil() ***REMOVED***
				return 0
			***REMOVED***
			rv = rv.Elem()
			if rv.IsNil() ***REMOVED***
				return 0
			***REMOVED***
			return si.oneofWrappersByType[rv.Type().Elem()]
		***REMOVED***
	***REMOVED***
	return oi
***REMOVED***

// isZero is identical to reflect.Value.IsZero.
// TODO: Remove this when Go1.13 is the minimally supported Go version.
func isZero(v reflect.Value) bool ***REMOVED***
	switch v.Kind() ***REMOVED***
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ ***REMOVED***
			if !isZero(v.Index(i)) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ ***REMOVED***
			if !isZero(v.Field(i)) ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	default:
		panic(&reflect.ValueError***REMOVED***"reflect.Value.IsZero", v.Kind()***REMOVED***)
	***REMOVED***
***REMOVED***
