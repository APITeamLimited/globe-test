// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dynamicpb creates protocol buffer messages using runtime type information.
package dynamicpb

import (
	"math"

	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// enum is a dynamic protoreflect.Enum.
type enum struct ***REMOVED***
	num protoreflect.EnumNumber
	typ protoreflect.EnumType
***REMOVED***

func (e enum) Descriptor() protoreflect.EnumDescriptor ***REMOVED*** return e.typ.Descriptor() ***REMOVED***
func (e enum) Type() protoreflect.EnumType             ***REMOVED*** return e.typ ***REMOVED***
func (e enum) Number() protoreflect.EnumNumber         ***REMOVED*** return e.num ***REMOVED***

// enumType is a dynamic protoreflect.EnumType.
type enumType struct ***REMOVED***
	desc protoreflect.EnumDescriptor
***REMOVED***

// NewEnumType creates a new EnumType with the provided descriptor.
//
// EnumTypes created by this package are equal if their descriptors are equal.
// That is, if ed1 == ed2, then NewEnumType(ed1) == NewEnumType(ed2).
//
// Enum values created by the EnumType are equal if their numbers are equal.
func NewEnumType(desc protoreflect.EnumDescriptor) protoreflect.EnumType ***REMOVED***
	return enumType***REMOVED***desc***REMOVED***
***REMOVED***

func (et enumType) New(n protoreflect.EnumNumber) protoreflect.Enum ***REMOVED*** return enum***REMOVED***n, et***REMOVED*** ***REMOVED***
func (et enumType) Descriptor() protoreflect.EnumDescriptor         ***REMOVED*** return et.desc ***REMOVED***

// extensionType is a dynamic protoreflect.ExtensionType.
type extensionType struct ***REMOVED***
	desc extensionTypeDescriptor
***REMOVED***

// A Message is a dynamically constructed protocol buffer message.
//
// Message implements the proto.Message interface, and may be used with all
// standard proto package functions such as Marshal, Unmarshal, and so forth.
//
// Message also implements the protoreflect.Message interface. See the protoreflect
// package documentation for that interface for how to get and set fields and
// otherwise interact with the contents of a Message.
//
// Reflection API functions which construct messages, such as NewField,
// return new dynamic messages of the appropriate type. Functions which take
// messages, such as Set for a message-value field, will accept any message
// with a compatible type.
//
// Operations which modify a Message are not safe for concurrent use.
type Message struct ***REMOVED***
	typ     messageType
	known   map[protoreflect.FieldNumber]protoreflect.Value
	ext     map[protoreflect.FieldNumber]protoreflect.FieldDescriptor
	unknown protoreflect.RawFields
***REMOVED***

var (
	_ protoreflect.Message      = (*Message)(nil)
	_ protoreflect.ProtoMessage = (*Message)(nil)
	_ protoiface.MessageV1      = (*Message)(nil)
)

// NewMessage creates a new message with the provided descriptor.
func NewMessage(desc protoreflect.MessageDescriptor) *Message ***REMOVED***
	return &Message***REMOVED***
		typ:   messageType***REMOVED***desc***REMOVED***,
		known: make(map[protoreflect.FieldNumber]protoreflect.Value),
		ext:   make(map[protoreflect.FieldNumber]protoreflect.FieldDescriptor),
	***REMOVED***
***REMOVED***

// ProtoMessage implements the legacy message interface.
func (m *Message) ProtoMessage() ***REMOVED******REMOVED***

// ProtoReflect implements the protoreflect.ProtoMessage interface.
func (m *Message) ProtoReflect() protoreflect.Message ***REMOVED***
	return m
***REMOVED***

// String returns a string representation of a message.
func (m *Message) String() string ***REMOVED***
	return protoimpl.X.MessageStringOf(m)
***REMOVED***

// Reset clears the message to be empty, but preserves the dynamic message type.
func (m *Message) Reset() ***REMOVED***
	m.known = make(map[protoreflect.FieldNumber]protoreflect.Value)
	m.ext = make(map[protoreflect.FieldNumber]protoreflect.FieldDescriptor)
	m.unknown = nil
***REMOVED***

// Descriptor returns the message descriptor.
func (m *Message) Descriptor() protoreflect.MessageDescriptor ***REMOVED***
	return m.typ.desc
***REMOVED***

// Type returns the message type.
func (m *Message) Type() protoreflect.MessageType ***REMOVED***
	return m.typ
***REMOVED***

// New returns a newly allocated empty message with the same descriptor.
// See protoreflect.Message for details.
func (m *Message) New() protoreflect.Message ***REMOVED***
	return m.Type().New()
***REMOVED***

// Interface returns the message.
// See protoreflect.Message for details.
func (m *Message) Interface() protoreflect.ProtoMessage ***REMOVED***
	return m
***REMOVED***

// ProtoMethods is an internal detail of the protoreflect.Message interface.
// Users should never call this directly.
func (m *Message) ProtoMethods() *protoiface.Methods ***REMOVED***
	return nil
***REMOVED***

// Range visits every populated field in undefined order.
// See protoreflect.Message for details.
func (m *Message) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) ***REMOVED***
	for num, v := range m.known ***REMOVED***
		fd := m.ext[num]
		if fd == nil ***REMOVED***
			fd = m.Descriptor().Fields().ByNumber(num)
		***REMOVED***
		if !isSet(fd, v) ***REMOVED***
			continue
		***REMOVED***
		if !f(fd, v) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Has reports whether a field is populated.
// See protoreflect.Message for details.
func (m *Message) Has(fd protoreflect.FieldDescriptor) bool ***REMOVED***
	m.checkField(fd)
	if fd.IsExtension() && m.ext[fd.Number()] != fd ***REMOVED***
		return false
	***REMOVED***
	v, ok := m.known[fd.Number()]
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return isSet(fd, v)
***REMOVED***

// Clear clears a field.
// See protoreflect.Message for details.
func (m *Message) Clear(fd protoreflect.FieldDescriptor) ***REMOVED***
	m.checkField(fd)
	num := fd.Number()
	delete(m.known, num)
	delete(m.ext, num)
***REMOVED***

// Get returns the value of a field.
// See protoreflect.Message for details.
func (m *Message) Get(fd protoreflect.FieldDescriptor) protoreflect.Value ***REMOVED***
	m.checkField(fd)
	num := fd.Number()
	if fd.IsExtension() ***REMOVED***
		if fd != m.ext[num] ***REMOVED***
			return fd.(protoreflect.ExtensionTypeDescriptor).Type().Zero()
		***REMOVED***
		return m.known[num]
	***REMOVED***
	if v, ok := m.known[num]; ok ***REMOVED***
		switch ***REMOVED***
		case fd.IsMap():
			if v.Map().Len() > 0 ***REMOVED***
				return v
			***REMOVED***
		case fd.IsList():
			if v.List().Len() > 0 ***REMOVED***
				return v
			***REMOVED***
		default:
			return v
		***REMOVED***
	***REMOVED***
	switch ***REMOVED***
	case fd.IsMap():
		return protoreflect.ValueOfMap(&dynamicMap***REMOVED***desc: fd***REMOVED***)
	case fd.IsList():
		return protoreflect.ValueOfList(emptyList***REMOVED***desc: fd***REMOVED***)
	case fd.Message() != nil:
		return protoreflect.ValueOfMessage(&Message***REMOVED***typ: messageType***REMOVED***fd.Message()***REMOVED******REMOVED***)
	case fd.Kind() == protoreflect.BytesKind:
		return protoreflect.ValueOfBytes(append([]byte(nil), fd.Default().Bytes()...))
	default:
		return fd.Default()
	***REMOVED***
***REMOVED***

// Mutable returns a mutable reference to a repeated, map, or message field.
// See protoreflect.Message for details.
func (m *Message) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value ***REMOVED***
	m.checkField(fd)
	if !fd.IsMap() && !fd.IsList() && fd.Message() == nil ***REMOVED***
		panic(errors.New("%v: getting mutable reference to non-composite type", fd.FullName()))
	***REMOVED***
	if m.known == nil ***REMOVED***
		panic(errors.New("%v: modification of read-only message", fd.FullName()))
	***REMOVED***
	num := fd.Number()
	if fd.IsExtension() ***REMOVED***
		if fd != m.ext[num] ***REMOVED***
			m.ext[num] = fd
			m.known[num] = fd.(protoreflect.ExtensionTypeDescriptor).Type().New()
		***REMOVED***
		return m.known[num]
	***REMOVED***
	if v, ok := m.known[num]; ok ***REMOVED***
		return v
	***REMOVED***
	m.clearOtherOneofFields(fd)
	m.known[num] = m.NewField(fd)
	if fd.IsExtension() ***REMOVED***
		m.ext[num] = fd
	***REMOVED***
	return m.known[num]
***REMOVED***

// Set stores a value in a field.
// See protoreflect.Message for details.
func (m *Message) Set(fd protoreflect.FieldDescriptor, v protoreflect.Value) ***REMOVED***
	m.checkField(fd)
	if m.known == nil ***REMOVED***
		panic(errors.New("%v: modification of read-only message", fd.FullName()))
	***REMOVED***
	if fd.IsExtension() ***REMOVED***
		isValid := true
		switch ***REMOVED***
		case !fd.(protoreflect.ExtensionTypeDescriptor).Type().IsValidValue(v):
			isValid = false
		case fd.IsList():
			isValid = v.List().IsValid()
		case fd.IsMap():
			isValid = v.Map().IsValid()
		case fd.Message() != nil:
			isValid = v.Message().IsValid()
		***REMOVED***
		if !isValid ***REMOVED***
			panic(errors.New("%v: assigning invalid type %T", fd.FullName(), v.Interface()))
		***REMOVED***
		m.ext[fd.Number()] = fd
	***REMOVED*** else ***REMOVED***
		typecheck(fd, v)
	***REMOVED***
	m.clearOtherOneofFields(fd)
	m.known[fd.Number()] = v
***REMOVED***

func (m *Message) clearOtherOneofFields(fd protoreflect.FieldDescriptor) ***REMOVED***
	od := fd.ContainingOneof()
	if od == nil ***REMOVED***
		return
	***REMOVED***
	num := fd.Number()
	for i := 0; i < od.Fields().Len(); i++ ***REMOVED***
		if n := od.Fields().Get(i).Number(); n != num ***REMOVED***
			delete(m.known, n)
		***REMOVED***
	***REMOVED***
***REMOVED***

// NewField returns a new value for assignable to the field of a given descriptor.
// See protoreflect.Message for details.
func (m *Message) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value ***REMOVED***
	m.checkField(fd)
	switch ***REMOVED***
	case fd.IsExtension():
		return fd.(protoreflect.ExtensionTypeDescriptor).Type().New()
	case fd.IsMap():
		return protoreflect.ValueOfMap(&dynamicMap***REMOVED***
			desc: fd,
			mapv: make(map[interface***REMOVED******REMOVED***]protoreflect.Value),
		***REMOVED***)
	case fd.IsList():
		return protoreflect.ValueOfList(&dynamicList***REMOVED***desc: fd***REMOVED***)
	case fd.Message() != nil:
		return protoreflect.ValueOfMessage(NewMessage(fd.Message()).ProtoReflect())
	default:
		return fd.Default()
	***REMOVED***
***REMOVED***

// WhichOneof reports which field in a oneof is populated, returning nil if none are populated.
// See protoreflect.Message for details.
func (m *Message) WhichOneof(od protoreflect.OneofDescriptor) protoreflect.FieldDescriptor ***REMOVED***
	for i := 0; i < od.Fields().Len(); i++ ***REMOVED***
		fd := od.Fields().Get(i)
		if m.Has(fd) ***REMOVED***
			return fd
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetUnknown returns the raw unknown fields.
// See protoreflect.Message for details.
func (m *Message) GetUnknown() protoreflect.RawFields ***REMOVED***
	return m.unknown
***REMOVED***

// SetUnknown sets the raw unknown fields.
// See protoreflect.Message for details.
func (m *Message) SetUnknown(r protoreflect.RawFields) ***REMOVED***
	if m.known == nil ***REMOVED***
		panic(errors.New("%v: modification of read-only message", m.typ.desc.FullName()))
	***REMOVED***
	m.unknown = r
***REMOVED***

// IsValid reports whether the message is valid.
// See protoreflect.Message for details.
func (m *Message) IsValid() bool ***REMOVED***
	return m.known != nil
***REMOVED***

func (m *Message) checkField(fd protoreflect.FieldDescriptor) ***REMOVED***
	if fd.IsExtension() && fd.ContainingMessage().FullName() == m.Descriptor().FullName() ***REMOVED***
		if _, ok := fd.(protoreflect.ExtensionTypeDescriptor); !ok ***REMOVED***
			panic(errors.New("%v: extension field descriptor does not implement ExtensionTypeDescriptor", fd.FullName()))
		***REMOVED***
		return
	***REMOVED***
	if fd.Parent() == m.Descriptor() ***REMOVED***
		return
	***REMOVED***
	fields := m.Descriptor().Fields()
	index := fd.Index()
	if index >= fields.Len() || fields.Get(index) != fd ***REMOVED***
		panic(errors.New("%v: field descriptor does not belong to this message", fd.FullName()))
	***REMOVED***
***REMOVED***

type messageType struct ***REMOVED***
	desc protoreflect.MessageDescriptor
***REMOVED***

// NewMessageType creates a new MessageType with the provided descriptor.
//
// MessageTypes created by this package are equal if their descriptors are equal.
// That is, if md1 == md2, then NewMessageType(md1) == NewMessageType(md2).
func NewMessageType(desc protoreflect.MessageDescriptor) protoreflect.MessageType ***REMOVED***
	return messageType***REMOVED***desc***REMOVED***
***REMOVED***

func (mt messageType) New() protoreflect.Message                  ***REMOVED*** return NewMessage(mt.desc) ***REMOVED***
func (mt messageType) Zero() protoreflect.Message                 ***REMOVED*** return &Message***REMOVED***typ: messageType***REMOVED***mt.desc***REMOVED******REMOVED*** ***REMOVED***
func (mt messageType) Descriptor() protoreflect.MessageDescriptor ***REMOVED*** return mt.desc ***REMOVED***
func (mt messageType) Enum(i int) protoreflect.EnumType ***REMOVED***
	if ed := mt.desc.Fields().Get(i).Enum(); ed != nil ***REMOVED***
		return NewEnumType(ed)
	***REMOVED***
	return nil
***REMOVED***
func (mt messageType) Message(i int) protoreflect.MessageType ***REMOVED***
	if md := mt.desc.Fields().Get(i).Message(); md != nil ***REMOVED***
		return NewMessageType(md)
	***REMOVED***
	return nil
***REMOVED***

type emptyList struct ***REMOVED***
	desc protoreflect.FieldDescriptor
***REMOVED***

func (x emptyList) Len() int                     ***REMOVED*** return 0 ***REMOVED***
func (x emptyList) Get(n int) protoreflect.Value ***REMOVED*** panic(errors.New("out of range")) ***REMOVED***
func (x emptyList) Set(n int, v protoreflect.Value) ***REMOVED***
	panic(errors.New("modification of immutable list"))
***REMOVED***
func (x emptyList) Append(v protoreflect.Value) ***REMOVED*** panic(errors.New("modification of immutable list")) ***REMOVED***
func (x emptyList) AppendMutable() protoreflect.Value ***REMOVED***
	panic(errors.New("modification of immutable list"))
***REMOVED***
func (x emptyList) Truncate(n int)                 ***REMOVED*** panic(errors.New("modification of immutable list")) ***REMOVED***
func (x emptyList) NewElement() protoreflect.Value ***REMOVED*** return newListEntry(x.desc) ***REMOVED***
func (x emptyList) IsValid() bool                  ***REMOVED*** return false ***REMOVED***

type dynamicList struct ***REMOVED***
	desc protoreflect.FieldDescriptor
	list []protoreflect.Value
***REMOVED***

func (x *dynamicList) Len() int ***REMOVED***
	return len(x.list)
***REMOVED***

func (x *dynamicList) Get(n int) protoreflect.Value ***REMOVED***
	return x.list[n]
***REMOVED***

func (x *dynamicList) Set(n int, v protoreflect.Value) ***REMOVED***
	typecheckSingular(x.desc, v)
	x.list[n] = v
***REMOVED***

func (x *dynamicList) Append(v protoreflect.Value) ***REMOVED***
	typecheckSingular(x.desc, v)
	x.list = append(x.list, v)
***REMOVED***

func (x *dynamicList) AppendMutable() protoreflect.Value ***REMOVED***
	if x.desc.Message() == nil ***REMOVED***
		panic(errors.New("%v: invalid AppendMutable on list with non-message type", x.desc.FullName()))
	***REMOVED***
	v := x.NewElement()
	x.Append(v)
	return v
***REMOVED***

func (x *dynamicList) Truncate(n int) ***REMOVED***
	// Zero truncated elements to avoid keeping data live.
	for i := n; i < len(x.list); i++ ***REMOVED***
		x.list[i] = protoreflect.Value***REMOVED******REMOVED***
	***REMOVED***
	x.list = x.list[:n]
***REMOVED***

func (x *dynamicList) NewElement() protoreflect.Value ***REMOVED***
	return newListEntry(x.desc)
***REMOVED***

func (x *dynamicList) IsValid() bool ***REMOVED***
	return true
***REMOVED***

type dynamicMap struct ***REMOVED***
	desc protoreflect.FieldDescriptor
	mapv map[interface***REMOVED******REMOVED***]protoreflect.Value
***REMOVED***

func (x *dynamicMap) Get(k protoreflect.MapKey) protoreflect.Value ***REMOVED*** return x.mapv[k.Interface()] ***REMOVED***
func (x *dynamicMap) Set(k protoreflect.MapKey, v protoreflect.Value) ***REMOVED***
	typecheckSingular(x.desc.MapKey(), k.Value())
	typecheckSingular(x.desc.MapValue(), v)
	x.mapv[k.Interface()] = v
***REMOVED***
func (x *dynamicMap) Has(k protoreflect.MapKey) bool ***REMOVED*** return x.Get(k).IsValid() ***REMOVED***
func (x *dynamicMap) Clear(k protoreflect.MapKey)    ***REMOVED*** delete(x.mapv, k.Interface()) ***REMOVED***
func (x *dynamicMap) Mutable(k protoreflect.MapKey) protoreflect.Value ***REMOVED***
	if x.desc.MapValue().Message() == nil ***REMOVED***
		panic(errors.New("%v: invalid Mutable on map with non-message value type", x.desc.FullName()))
	***REMOVED***
	v := x.Get(k)
	if !v.IsValid() ***REMOVED***
		v = x.NewValue()
		x.Set(k, v)
	***REMOVED***
	return v
***REMOVED***
func (x *dynamicMap) Len() int ***REMOVED*** return len(x.mapv) ***REMOVED***
func (x *dynamicMap) NewValue() protoreflect.Value ***REMOVED***
	if md := x.desc.MapValue().Message(); md != nil ***REMOVED***
		return protoreflect.ValueOfMessage(NewMessage(md).ProtoReflect())
	***REMOVED***
	return x.desc.MapValue().Default()
***REMOVED***
func (x *dynamicMap) IsValid() bool ***REMOVED***
	return x.mapv != nil
***REMOVED***

func (x *dynamicMap) Range(f func(protoreflect.MapKey, protoreflect.Value) bool) ***REMOVED***
	for k, v := range x.mapv ***REMOVED***
		if !f(protoreflect.ValueOf(k).MapKey(), v) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func isSet(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
	switch ***REMOVED***
	case fd.IsMap():
		return v.Map().Len() > 0
	case fd.IsList():
		return v.List().Len() > 0
	case fd.ContainingOneof() != nil:
		return true
	case fd.Syntax() == protoreflect.Proto3 && !fd.IsExtension():
		switch fd.Kind() ***REMOVED***
		case protoreflect.BoolKind:
			return v.Bool()
		case protoreflect.EnumKind:
			return v.Enum() != 0
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
			return v.Int() != 0
		case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
			return v.Uint() != 0
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			return v.Float() != 0 || math.Signbit(v.Float())
		case protoreflect.StringKind:
			return v.String() != ""
		case protoreflect.BytesKind:
			return len(v.Bytes()) > 0
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func typecheck(fd protoreflect.FieldDescriptor, v protoreflect.Value) ***REMOVED***
	if err := typeIsValid(fd, v); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func typeIsValid(fd protoreflect.FieldDescriptor, v protoreflect.Value) error ***REMOVED***
	switch ***REMOVED***
	case !v.IsValid():
		return errors.New("%v: assigning invalid value", fd.FullName())
	case fd.IsMap():
		if mapv, ok := v.Interface().(*dynamicMap); !ok || mapv.desc != fd || !mapv.IsValid() ***REMOVED***
			return errors.New("%v: assigning invalid type %T", fd.FullName(), v.Interface())
		***REMOVED***
		return nil
	case fd.IsList():
		switch list := v.Interface().(type) ***REMOVED***
		case *dynamicList:
			if list.desc == fd && list.IsValid() ***REMOVED***
				return nil
			***REMOVED***
		case emptyList:
			if list.desc == fd && list.IsValid() ***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
		return errors.New("%v: assigning invalid type %T", fd.FullName(), v.Interface())
	default:
		return singularTypeIsValid(fd, v)
	***REMOVED***
***REMOVED***

func typecheckSingular(fd protoreflect.FieldDescriptor, v protoreflect.Value) ***REMOVED***
	if err := singularTypeIsValid(fd, v); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func singularTypeIsValid(fd protoreflect.FieldDescriptor, v protoreflect.Value) error ***REMOVED***
	vi := v.Interface()
	var ok bool
	switch fd.Kind() ***REMOVED***
	case protoreflect.BoolKind:
		_, ok = vi.(bool)
	case protoreflect.EnumKind:
		// We could check against the valid set of enum values, but do not.
		_, ok = vi.(protoreflect.EnumNumber)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		_, ok = vi.(int32)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		_, ok = vi.(uint32)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		_, ok = vi.(int64)
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		_, ok = vi.(uint64)
	case protoreflect.FloatKind:
		_, ok = vi.(float32)
	case protoreflect.DoubleKind:
		_, ok = vi.(float64)
	case protoreflect.StringKind:
		_, ok = vi.(string)
	case protoreflect.BytesKind:
		_, ok = vi.([]byte)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		var m protoreflect.Message
		m, ok = vi.(protoreflect.Message)
		if ok && m.Descriptor().FullName() != fd.Message().FullName() ***REMOVED***
			return errors.New("%v: assigning invalid message type %v", fd.FullName(), m.Descriptor().FullName())
		***REMOVED***
		if dm, ok := vi.(*Message); ok && dm.known == nil ***REMOVED***
			return errors.New("%v: assigning invalid zero-value message", fd.FullName())
		***REMOVED***
	***REMOVED***
	if !ok ***REMOVED***
		return errors.New("%v: assigning invalid type %T", fd.FullName(), v.Interface())
	***REMOVED***
	return nil
***REMOVED***

func newListEntry(fd protoreflect.FieldDescriptor) protoreflect.Value ***REMOVED***
	switch fd.Kind() ***REMOVED***
	case protoreflect.BoolKind:
		return protoreflect.ValueOfBool(false)
	case protoreflect.EnumKind:
		return protoreflect.ValueOfEnum(fd.Enum().Values().Get(0).Number())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return protoreflect.ValueOfInt32(0)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return protoreflect.ValueOfUint32(0)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return protoreflect.ValueOfInt64(0)
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return protoreflect.ValueOfUint64(0)
	case protoreflect.FloatKind:
		return protoreflect.ValueOfFloat32(0)
	case protoreflect.DoubleKind:
		return protoreflect.ValueOfFloat64(0)
	case protoreflect.StringKind:
		return protoreflect.ValueOfString("")
	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes(nil)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return protoreflect.ValueOfMessage(NewMessage(fd.Message()).ProtoReflect())
	***REMOVED***
	panic(errors.New("%v: unknown kind %v", fd.FullName(), fd.Kind()))
***REMOVED***

// NewExtensionType creates a new ExtensionType with the provided descriptor.
//
// Dynamic ExtensionTypes with the same descriptor compare as equal. That is,
// if xd1 == xd2, then NewExtensionType(xd1) == NewExtensionType(xd2).
//
// The InterfaceOf and ValueOf methods of the extension type are defined as:
//
//	func (xt extensionType) ValueOf(iv interface***REMOVED******REMOVED***) protoreflect.Value ***REMOVED***
//		return protoreflect.ValueOf(iv)
//	***REMOVED***
//
//	func (xt extensionType) InterfaceOf(v protoreflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
//		return v.Interface()
//	***REMOVED***
//
// The Go type used by the proto.GetExtension and proto.SetExtension functions
// is determined by these methods, and is therefore equivalent to the Go type
// used to represent a protoreflect.Value. See the protoreflect.Value
// documentation for more details.
func NewExtensionType(desc protoreflect.ExtensionDescriptor) protoreflect.ExtensionType ***REMOVED***
	if xt, ok := desc.(protoreflect.ExtensionTypeDescriptor); ok ***REMOVED***
		desc = xt.Descriptor()
	***REMOVED***
	return extensionType***REMOVED***extensionTypeDescriptor***REMOVED***desc***REMOVED******REMOVED***
***REMOVED***

func (xt extensionType) New() protoreflect.Value ***REMOVED***
	switch ***REMOVED***
	case xt.desc.IsMap():
		return protoreflect.ValueOfMap(&dynamicMap***REMOVED***
			desc: xt.desc,
			mapv: make(map[interface***REMOVED******REMOVED***]protoreflect.Value),
		***REMOVED***)
	case xt.desc.IsList():
		return protoreflect.ValueOfList(&dynamicList***REMOVED***desc: xt.desc***REMOVED***)
	case xt.desc.Message() != nil:
		return protoreflect.ValueOfMessage(NewMessage(xt.desc.Message()))
	default:
		return xt.desc.Default()
	***REMOVED***
***REMOVED***

func (xt extensionType) Zero() protoreflect.Value ***REMOVED***
	switch ***REMOVED***
	case xt.desc.IsMap():
		return protoreflect.ValueOfMap(&dynamicMap***REMOVED***desc: xt.desc***REMOVED***)
	case xt.desc.Cardinality() == protoreflect.Repeated:
		return protoreflect.ValueOfList(emptyList***REMOVED***desc: xt.desc***REMOVED***)
	case xt.desc.Message() != nil:
		return protoreflect.ValueOfMessage(&Message***REMOVED***typ: messageType***REMOVED***xt.desc.Message()***REMOVED******REMOVED***)
	default:
		return xt.desc.Default()
	***REMOVED***
***REMOVED***

func (xt extensionType) TypeDescriptor() protoreflect.ExtensionTypeDescriptor ***REMOVED***
	return xt.desc
***REMOVED***

func (xt extensionType) ValueOf(iv interface***REMOVED******REMOVED***) protoreflect.Value ***REMOVED***
	v := protoreflect.ValueOf(iv)
	typecheck(xt.desc, v)
	return v
***REMOVED***

func (xt extensionType) InterfaceOf(v protoreflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
	typecheck(xt.desc, v)
	return v.Interface()
***REMOVED***

func (xt extensionType) IsValidInterface(iv interface***REMOVED******REMOVED***) bool ***REMOVED***
	return typeIsValid(xt.desc, protoreflect.ValueOf(iv)) == nil
***REMOVED***

func (xt extensionType) IsValidValue(v protoreflect.Value) bool ***REMOVED***
	return typeIsValid(xt.desc, v) == nil
***REMOVED***

type extensionTypeDescriptor struct ***REMOVED***
	protoreflect.ExtensionDescriptor
***REMOVED***

func (xt extensionTypeDescriptor) Type() protoreflect.ExtensionType ***REMOVED***
	return extensionType***REMOVED***xt***REMOVED***
***REMOVED***

func (xt extensionTypeDescriptor) Descriptor() protoreflect.ExtensionDescriptor ***REMOVED***
	return xt.ExtensionDescriptor
***REMOVED***
