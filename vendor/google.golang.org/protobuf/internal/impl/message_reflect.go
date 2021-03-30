// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/internal/detrand"
	"google.golang.org/protobuf/internal/pragma"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

type reflectMessageInfo struct ***REMOVED***
	fields map[pref.FieldNumber]*fieldInfo
	oneofs map[pref.Name]*oneofInfo

	// fieldTypes contains the zero value of an enum or message field.
	// For lists, it contains the element type.
	// For maps, it contains the entry value type.
	fieldTypes map[pref.FieldNumber]interface***REMOVED******REMOVED***

	// denseFields is a subset of fields where:
	//	0 < fieldDesc.Number() < len(denseFields)
	// It provides faster access to the fieldInfo, but may be incomplete.
	denseFields []*fieldInfo

	// rangeInfos is a list of all fields (not belonging to a oneof) and oneofs.
	rangeInfos []interface***REMOVED******REMOVED*** // either *fieldInfo or *oneofInfo

	getUnknown   func(pointer) pref.RawFields
	setUnknown   func(pointer, pref.RawFields)
	extensionMap func(pointer) *extensionMap

	nilMessage atomicNilMessage
***REMOVED***

// makeReflectFuncs generates the set of functions to support reflection.
func (mi *MessageInfo) makeReflectFuncs(t reflect.Type, si structInfo) ***REMOVED***
	mi.makeKnownFieldsFunc(si)
	mi.makeUnknownFieldsFunc(t, si)
	mi.makeExtensionFieldsFunc(t, si)
	mi.makeFieldTypes(si)
***REMOVED***

// makeKnownFieldsFunc generates functions for operations that can be performed
// on each protobuf message field. It takes in a reflect.Type representing the
// Go struct and matches message fields with struct fields.
//
// This code assumes that the struct is well-formed and panics if there are
// any discrepancies.
func (mi *MessageInfo) makeKnownFieldsFunc(si structInfo) ***REMOVED***
	mi.fields = map[pref.FieldNumber]*fieldInfo***REMOVED******REMOVED***
	md := mi.Desc
	fds := md.Fields()
	for i := 0; i < fds.Len(); i++ ***REMOVED***
		fd := fds.Get(i)
		fs := si.fieldsByNumber[fd.Number()]
		var fi fieldInfo
		switch ***REMOVED***
		case fd.ContainingOneof() != nil && !fd.ContainingOneof().IsSynthetic():
			fi = fieldInfoForOneof(fd, si.oneofsByName[fd.ContainingOneof().Name()], mi.Exporter, si.oneofWrappersByNumber[fd.Number()])
		case fd.IsMap():
			fi = fieldInfoForMap(fd, fs, mi.Exporter)
		case fd.IsList():
			fi = fieldInfoForList(fd, fs, mi.Exporter)
		case fd.IsWeak():
			fi = fieldInfoForWeakMessage(fd, si.weakOffset)
		case fd.Message() != nil:
			fi = fieldInfoForMessage(fd, fs, mi.Exporter)
		default:
			fi = fieldInfoForScalar(fd, fs, mi.Exporter)
		***REMOVED***
		mi.fields[fd.Number()] = &fi
	***REMOVED***

	mi.oneofs = map[pref.Name]*oneofInfo***REMOVED******REMOVED***
	for i := 0; i < md.Oneofs().Len(); i++ ***REMOVED***
		od := md.Oneofs().Get(i)
		mi.oneofs[od.Name()] = makeOneofInfo(od, si, mi.Exporter)
	***REMOVED***

	mi.denseFields = make([]*fieldInfo, fds.Len()*2)
	for i := 0; i < fds.Len(); i++ ***REMOVED***
		if fd := fds.Get(i); int(fd.Number()) < len(mi.denseFields) ***REMOVED***
			mi.denseFields[fd.Number()] = mi.fields[fd.Number()]
		***REMOVED***
	***REMOVED***

	for i := 0; i < fds.Len(); ***REMOVED***
		fd := fds.Get(i)
		if od := fd.ContainingOneof(); od != nil && !od.IsSynthetic() ***REMOVED***
			mi.rangeInfos = append(mi.rangeInfos, mi.oneofs[od.Name()])
			i += od.Fields().Len()
		***REMOVED*** else ***REMOVED***
			mi.rangeInfos = append(mi.rangeInfos, mi.fields[fd.Number()])
			i++
		***REMOVED***
	***REMOVED***

	// Introduce instability to iteration order, but keep it deterministic.
	if len(mi.rangeInfos) > 1 && detrand.Bool() ***REMOVED***
		i := detrand.Intn(len(mi.rangeInfos) - 1)
		mi.rangeInfos[i], mi.rangeInfos[i+1] = mi.rangeInfos[i+1], mi.rangeInfos[i]
	***REMOVED***
***REMOVED***

func (mi *MessageInfo) makeUnknownFieldsFunc(t reflect.Type, si structInfo) ***REMOVED***
	switch ***REMOVED***
	case si.unknownOffset.IsValid() && si.unknownType == unknownFieldsAType:
		// Handle as []byte.
		mi.getUnknown = func(p pointer) pref.RawFields ***REMOVED***
			if p.IsNil() ***REMOVED***
				return nil
			***REMOVED***
			return *p.Apply(mi.unknownOffset).Bytes()
		***REMOVED***
		mi.setUnknown = func(p pointer, b pref.RawFields) ***REMOVED***
			if p.IsNil() ***REMOVED***
				panic("invalid SetUnknown on nil Message")
			***REMOVED***
			*p.Apply(mi.unknownOffset).Bytes() = b
		***REMOVED***
	case si.unknownOffset.IsValid() && si.unknownType == unknownFieldsBType:
		// Handle as *[]byte.
		mi.getUnknown = func(p pointer) pref.RawFields ***REMOVED***
			if p.IsNil() ***REMOVED***
				return nil
			***REMOVED***
			bp := p.Apply(mi.unknownOffset).BytesPtr()
			if *bp == nil ***REMOVED***
				return nil
			***REMOVED***
			return **bp
		***REMOVED***
		mi.setUnknown = func(p pointer, b pref.RawFields) ***REMOVED***
			if p.IsNil() ***REMOVED***
				panic("invalid SetUnknown on nil Message")
			***REMOVED***
			bp := p.Apply(mi.unknownOffset).BytesPtr()
			if *bp == nil ***REMOVED***
				*bp = new([]byte)
			***REMOVED***
			**bp = b
		***REMOVED***
	default:
		mi.getUnknown = func(pointer) pref.RawFields ***REMOVED***
			return nil
		***REMOVED***
		mi.setUnknown = func(p pointer, _ pref.RawFields) ***REMOVED***
			if p.IsNil() ***REMOVED***
				panic("invalid SetUnknown on nil Message")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (mi *MessageInfo) makeExtensionFieldsFunc(t reflect.Type, si structInfo) ***REMOVED***
	if si.extensionOffset.IsValid() ***REMOVED***
		mi.extensionMap = func(p pointer) *extensionMap ***REMOVED***
			if p.IsNil() ***REMOVED***
				return (*extensionMap)(nil)
			***REMOVED***
			v := p.Apply(si.extensionOffset).AsValueOf(extensionFieldsType)
			return (*extensionMap)(v.Interface().(*map[int32]ExtensionField))
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		mi.extensionMap = func(pointer) *extensionMap ***REMOVED***
			return (*extensionMap)(nil)
		***REMOVED***
	***REMOVED***
***REMOVED***
func (mi *MessageInfo) makeFieldTypes(si structInfo) ***REMOVED***
	md := mi.Desc
	fds := md.Fields()
	for i := 0; i < fds.Len(); i++ ***REMOVED***
		var ft reflect.Type
		fd := fds.Get(i)
		fs := si.fieldsByNumber[fd.Number()]
		switch ***REMOVED***
		case fd.ContainingOneof() != nil && !fd.ContainingOneof().IsSynthetic():
			if fd.Enum() != nil || fd.Message() != nil ***REMOVED***
				ft = si.oneofWrappersByNumber[fd.Number()].Field(0).Type
			***REMOVED***
		case fd.IsMap():
			if fd.MapValue().Enum() != nil || fd.MapValue().Message() != nil ***REMOVED***
				ft = fs.Type.Elem()
			***REMOVED***
		case fd.IsList():
			if fd.Enum() != nil || fd.Message() != nil ***REMOVED***
				ft = fs.Type.Elem()
			***REMOVED***
		case fd.Enum() != nil:
			ft = fs.Type
			if fd.HasPresence() ***REMOVED***
				ft = ft.Elem()
			***REMOVED***
		case fd.Message() != nil:
			ft = fs.Type
			if fd.IsWeak() ***REMOVED***
				ft = nil
			***REMOVED***
		***REMOVED***
		if ft != nil ***REMOVED***
			if mi.fieldTypes == nil ***REMOVED***
				mi.fieldTypes = make(map[pref.FieldNumber]interface***REMOVED******REMOVED***)
			***REMOVED***
			mi.fieldTypes[fd.Number()] = reflect.Zero(ft).Interface()
		***REMOVED***
	***REMOVED***
***REMOVED***

type extensionMap map[int32]ExtensionField

func (m *extensionMap) Range(f func(pref.FieldDescriptor, pref.Value) bool) ***REMOVED***
	if m != nil ***REMOVED***
		for _, x := range *m ***REMOVED***
			xd := x.Type().TypeDescriptor()
			v := x.Value()
			if xd.IsList() && v.List().Len() == 0 ***REMOVED***
				continue
			***REMOVED***
			if !f(xd, v) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
func (m *extensionMap) Has(xt pref.ExtensionType) (ok bool) ***REMOVED***
	if m == nil ***REMOVED***
		return false
	***REMOVED***
	xd := xt.TypeDescriptor()
	x, ok := (*m)[int32(xd.Number())]
	if !ok ***REMOVED***
		return false
	***REMOVED***
	switch ***REMOVED***
	case xd.IsList():
		return x.Value().List().Len() > 0
	case xd.IsMap():
		return x.Value().Map().Len() > 0
	case xd.Message() != nil:
		return x.Value().Message().IsValid()
	***REMOVED***
	return true
***REMOVED***
func (m *extensionMap) Clear(xt pref.ExtensionType) ***REMOVED***
	delete(*m, int32(xt.TypeDescriptor().Number()))
***REMOVED***
func (m *extensionMap) Get(xt pref.ExtensionType) pref.Value ***REMOVED***
	xd := xt.TypeDescriptor()
	if m != nil ***REMOVED***
		if x, ok := (*m)[int32(xd.Number())]; ok ***REMOVED***
			return x.Value()
		***REMOVED***
	***REMOVED***
	return xt.Zero()
***REMOVED***
func (m *extensionMap) Set(xt pref.ExtensionType, v pref.Value) ***REMOVED***
	xd := xt.TypeDescriptor()
	isValid := true
	switch ***REMOVED***
	case !xt.IsValidValue(v):
		isValid = false
	case xd.IsList():
		isValid = v.List().IsValid()
	case xd.IsMap():
		isValid = v.Map().IsValid()
	case xd.Message() != nil:
		isValid = v.Message().IsValid()
	***REMOVED***
	if !isValid ***REMOVED***
		panic(fmt.Sprintf("%v: assigning invalid value", xt.TypeDescriptor().FullName()))
	***REMOVED***

	if *m == nil ***REMOVED***
		*m = make(map[int32]ExtensionField)
	***REMOVED***
	var x ExtensionField
	x.Set(xt, v)
	(*m)[int32(xd.Number())] = x
***REMOVED***
func (m *extensionMap) Mutable(xt pref.ExtensionType) pref.Value ***REMOVED***
	xd := xt.TypeDescriptor()
	if xd.Kind() != pref.MessageKind && xd.Kind() != pref.GroupKind && !xd.IsList() && !xd.IsMap() ***REMOVED***
		panic("invalid Mutable on field with non-composite type")
	***REMOVED***
	if x, ok := (*m)[int32(xd.Number())]; ok ***REMOVED***
		return x.Value()
	***REMOVED***
	v := xt.New()
	m.Set(xt, v)
	return v
***REMOVED***

// MessageState is a data structure that is nested as the first field in a
// concrete message. It provides a way to implement the ProtoReflect method
// in an allocation-free way without needing to have a shadow Go type generated
// for every message type. This technique only works using unsafe.
//
//
// Example generated code:
//
//	type M struct ***REMOVED***
//		state protoimpl.MessageState
//
//		Field1 int32
//		Field2 string
//		Field3 *BarMessage
//		...
//	***REMOVED***
//
//	func (m *M) ProtoReflect() protoreflect.Message ***REMOVED***
//		mi := &file_fizz_buzz_proto_msgInfos[5]
//		if protoimpl.UnsafeEnabled && m != nil ***REMOVED***
//			ms := protoimpl.X.MessageStateOf(Pointer(m))
//			if ms.LoadMessageInfo() == nil ***REMOVED***
//				ms.StoreMessageInfo(mi)
//			***REMOVED***
//			return ms
//		***REMOVED***
//		return mi.MessageOf(m)
//	***REMOVED***
//
// The MessageState type holds a *MessageInfo, which must be atomically set to
// the message info associated with a given message instance.
// By unsafely converting a *M into a *MessageState, the MessageState object
// has access to all the information needed to implement protobuf reflection.
// It has access to the message info as its first field, and a pointer to the
// MessageState is identical to a pointer to the concrete message value.
//
//
// Requirements:
//	• The type M must implement protoreflect.ProtoMessage.
//	• The address of m must not be nil.
//	• The address of m and the address of m.state must be equal,
//	even though they are different Go types.
type MessageState struct ***REMOVED***
	pragma.NoUnkeyedLiterals
	pragma.DoNotCompare
	pragma.DoNotCopy

	atomicMessageInfo *MessageInfo
***REMOVED***

type messageState MessageState

var (
	_ pref.Message = (*messageState)(nil)
	_ unwrapper    = (*messageState)(nil)
)

// messageDataType is a tuple of a pointer to the message data and
// a pointer to the message type. It is a generalized way of providing a
// reflective view over a message instance. The disadvantage of this approach
// is the need to allocate this tuple of 16B.
type messageDataType struct ***REMOVED***
	p  pointer
	mi *MessageInfo
***REMOVED***

type (
	messageReflectWrapper messageDataType
	messageIfaceWrapper   messageDataType
)

var (
	_ pref.Message      = (*messageReflectWrapper)(nil)
	_ unwrapper         = (*messageReflectWrapper)(nil)
	_ pref.ProtoMessage = (*messageIfaceWrapper)(nil)
	_ unwrapper         = (*messageIfaceWrapper)(nil)
)

// MessageOf returns a reflective view over a message. The input must be a
// pointer to a named Go struct. If the provided type has a ProtoReflect method,
// it must be implemented by calling this method.
func (mi *MessageInfo) MessageOf(m interface***REMOVED******REMOVED***) pref.Message ***REMOVED***
	if reflect.TypeOf(m) != mi.GoReflectType ***REMOVED***
		panic(fmt.Sprintf("type mismatch: got %T, want %v", m, mi.GoReflectType))
	***REMOVED***
	p := pointerOfIface(m)
	if p.IsNil() ***REMOVED***
		return mi.nilMessage.Init(mi)
	***REMOVED***
	return &messageReflectWrapper***REMOVED***p, mi***REMOVED***
***REMOVED***

func (m *messageReflectWrapper) pointer() pointer          ***REMOVED*** return m.p ***REMOVED***
func (m *messageReflectWrapper) messageInfo() *MessageInfo ***REMOVED*** return m.mi ***REMOVED***

func (m *messageIfaceWrapper) ProtoReflect() pref.Message ***REMOVED***
	return (*messageReflectWrapper)(m)
***REMOVED***
func (m *messageIfaceWrapper) protoUnwrap() interface***REMOVED******REMOVED*** ***REMOVED***
	return m.p.AsIfaceOf(m.mi.GoReflectType.Elem())
***REMOVED***

// checkField verifies that the provided field descriptor is valid.
// Exactly one of the returned values is populated.
func (mi *MessageInfo) checkField(fd pref.FieldDescriptor) (*fieldInfo, pref.ExtensionType) ***REMOVED***
	var fi *fieldInfo
	if n := fd.Number(); 0 < n && int(n) < len(mi.denseFields) ***REMOVED***
		fi = mi.denseFields[n]
	***REMOVED*** else ***REMOVED***
		fi = mi.fields[n]
	***REMOVED***
	if fi != nil ***REMOVED***
		if fi.fieldDesc != fd ***REMOVED***
			if got, want := fd.FullName(), fi.fieldDesc.FullName(); got != want ***REMOVED***
				panic(fmt.Sprintf("mismatching field: got %v, want %v", got, want))
			***REMOVED***
			panic(fmt.Sprintf("mismatching field: %v", fd.FullName()))
		***REMOVED***
		return fi, nil
	***REMOVED***

	if fd.IsExtension() ***REMOVED***
		if got, want := fd.ContainingMessage().FullName(), mi.Desc.FullName(); got != want ***REMOVED***
			// TODO: Should this be exact containing message descriptor match?
			panic(fmt.Sprintf("extension %v has mismatching containing message: got %v, want %v", fd.FullName(), got, want))
		***REMOVED***
		if !mi.Desc.ExtensionRanges().Has(fd.Number()) ***REMOVED***
			panic(fmt.Sprintf("extension %v extends %v outside the extension range", fd.FullName(), mi.Desc.FullName()))
		***REMOVED***
		xtd, ok := fd.(pref.ExtensionTypeDescriptor)
		if !ok ***REMOVED***
			panic(fmt.Sprintf("extension %v does not implement protoreflect.ExtensionTypeDescriptor", fd.FullName()))
		***REMOVED***
		return nil, xtd.Type()
	***REMOVED***
	panic(fmt.Sprintf("field %v is invalid", fd.FullName()))
***REMOVED***
