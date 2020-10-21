// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dynamicpb creates protocol buffer messages using runtime type information.
package dynamicpb

import (
	"math"

	"google.golang.org/protobuf/internal/errors"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/runtime/protoimpl"
)

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
	known   map[pref.FieldNumber]pref.Value
	ext     map[pref.FieldNumber]pref.FieldDescriptor
	unknown pref.RawFields
***REMOVED***

var (
	_ pref.Message         = (*Message)(nil)
	_ pref.ProtoMessage    = (*Message)(nil)
	_ protoiface.MessageV1 = (*Message)(nil)
)

// NewMessage creates a new message with the provided descriptor.
func NewMessage(desc pref.MessageDescriptor) *Message ***REMOVED***
	return &Message***REMOVED***
		typ:   messageType***REMOVED***desc***REMOVED***,
		known: make(map[pref.FieldNumber]pref.Value),
		ext:   make(map[pref.FieldNumber]pref.FieldDescriptor),
	***REMOVED***
***REMOVED***

// ProtoMessage implements the legacy message interface.
func (m *Message) ProtoMessage() ***REMOVED******REMOVED***

// ProtoReflect implements the protoreflect.ProtoMessage interface.
func (m *Message) ProtoReflect() pref.Message ***REMOVED***
	return m
***REMOVED***

// String returns a string representation of a message.
func (m *Message) String() string ***REMOVED***
	return protoimpl.X.MessageStringOf(m)
***REMOVED***

// Reset clears the message to be empty, but preserves the dynamic message type.
func (m *Message) Reset() ***REMOVED***
	m.known = make(map[pref.FieldNumber]pref.Value)
	m.ext = make(map[pref.FieldNumber]pref.FieldDescriptor)
	m.unknown = nil
***REMOVED***

// Descriptor returns the message descriptor.
func (m *Message) Descriptor() pref.MessageDescriptor ***REMOVED***
	return m.typ.desc
***REMOVED***

// Type returns the message type.
func (m *Message) Type() pref.MessageType ***REMOVED***
	return m.typ
***REMOVED***

// New returns a newly allocated empty message with the same descriptor.
// See protoreflect.Message for details.
func (m *Message) New() pref.Message ***REMOVED***
	return m.Type().New()
***REMOVED***

// Interface returns the message.
// See protoreflect.Message for details.
func (m *Message) Interface() pref.ProtoMessage ***REMOVED***
	return m
***REMOVED***

// ProtoMethods is an internal detail of the protoreflect.Message interface.
// Users should never call this directly.
func (m *Message) ProtoMethods() *protoiface.Methods ***REMOVED***
	return nil
***REMOVED***

// Range visits every populated field in undefined order.
// See protoreflect.Message for details.
func (m *Message) Range(f func(pref.FieldDescriptor, pref.Value) bool) ***REMOVED***
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
func (m *Message) Has(fd pref.FieldDescriptor) bool ***REMOVED***
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
func (m *Message) Clear(fd pref.FieldDescriptor) ***REMOVED***
	m.checkField(fd)
	num := fd.Number()
	delete(m.known, num)
	delete(m.ext, num)
***REMOVED***

// Get returns the value of a field.
// See protoreflect.Message for details.
func (m *Message) Get(fd pref.FieldDescriptor) pref.Value ***REMOVED***
	m.checkField(fd)
	num := fd.Number()
	if fd.IsExtension() ***REMOVED***
		if fd != m.ext[num] ***REMOVED***
			return fd.(pref.ExtensionTypeDescriptor).Type().Zero()
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
		return pref.ValueOfMap(&dynamicMap***REMOVED***desc: fd***REMOVED***)
	case fd.IsList():
		return pref.ValueOfList(emptyList***REMOVED***desc: fd***REMOVED***)
	case fd.Message() != nil:
		return pref.ValueOfMessage(&Message***REMOVED***typ: messageType***REMOVED***fd.Message()***REMOVED******REMOVED***)
	case fd.Kind() == pref.BytesKind:
		return pref.ValueOfBytes(append([]byte(nil), fd.Default().Bytes()...))
	default:
		return fd.Default()
	***REMOVED***
***REMOVED***

// Mutable returns a mutable reference to a repeated, map, or message field.
// See protoreflect.Message for details.
func (m *Message) Mutable(fd pref.FieldDescriptor) pref.Value ***REMOVED***
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
			m.known[num] = fd.(pref.ExtensionTypeDescriptor).Type().New()
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
func (m *Message) Set(fd pref.FieldDescriptor, v pref.Value) ***REMOVED***
	m.checkField(fd)
	if m.known == nil ***REMOVED***
		panic(errors.New("%v: modification of read-only message", fd.FullName()))
	***REMOVED***
	if fd.IsExtension() ***REMOVED***
		isValid := true
		switch ***REMOVED***
		case !fd.(pref.ExtensionTypeDescriptor).Type().IsValidValue(v):
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

func (m *Message) clearOtherOneofFields(fd pref.FieldDescriptor) ***REMOVED***
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
func (m *Message) NewField(fd pref.FieldDescriptor) pref.Value ***REMOVED***
	m.checkField(fd)
	switch ***REMOVED***
	case fd.IsExtension():
		return fd.(pref.ExtensionTypeDescriptor).Type().New()
	case fd.IsMap():
		return pref.ValueOfMap(&dynamicMap***REMOVED***
			desc: fd,
			mapv: make(map[interface***REMOVED******REMOVED***]pref.Value),
		***REMOVED***)
	case fd.IsList():
		return pref.ValueOfList(&dynamicList***REMOVED***desc: fd***REMOVED***)
	case fd.Message() != nil:
		return pref.ValueOfMessage(NewMessage(fd.Message()).ProtoReflect())
	default:
		return fd.Default()
	***REMOVED***
***REMOVED***

// WhichOneof reports which field in a oneof is populated, returning nil if none are populated.
// See protoreflect.Message for details.
func (m *Message) WhichOneof(od pref.OneofDescriptor) pref.FieldDescriptor ***REMOVED***
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
func (m *Message) GetUnknown() pref.RawFields ***REMOVED***
	return m.unknown
***REMOVED***

// SetUnknown sets the raw unknown fields.
// See protoreflect.Message for details.
func (m *Message) SetUnknown(r pref.RawFields) ***REMOVED***
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

func (m *Message) checkField(fd pref.FieldDescriptor) ***REMOVED***
	if fd.IsExtension() && fd.ContainingMessage().FullName() == m.Descriptor().FullName() ***REMOVED***
		if _, ok := fd.(pref.ExtensionTypeDescriptor); !ok ***REMOVED***
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
	desc pref.MessageDescriptor
***REMOVED***

// NewMessageType creates a new MessageType with the provided descriptor.
//
// MessageTypes created by this package are equal if their descriptors are equal.
// That is, if md1 == md2, then NewMessageType(md1) == NewMessageType(md2).
func NewMessageType(desc pref.MessageDescriptor) pref.MessageType ***REMOVED***
	return messageType***REMOVED***desc***REMOVED***
***REMOVED***

func (mt messageType) New() pref.Message                  ***REMOVED*** return NewMessage(mt.desc) ***REMOVED***
func (mt messageType) Zero() pref.Message                 ***REMOVED*** return &Message***REMOVED***typ: messageType***REMOVED***mt.desc***REMOVED******REMOVED*** ***REMOVED***
func (mt messageType) Descriptor() pref.MessageDescriptor ***REMOVED*** return mt.desc ***REMOVED***

type emptyList struct ***REMOVED***
	desc pref.FieldDescriptor
***REMOVED***

func (x emptyList) Len() int                  ***REMOVED*** return 0 ***REMOVED***
func (x emptyList) Get(n int) pref.Value      ***REMOVED*** panic(errors.New("out of range")) ***REMOVED***
func (x emptyList) Set(n int, v pref.Value)   ***REMOVED*** panic(errors.New("modification of immutable list")) ***REMOVED***
func (x emptyList) Append(v pref.Value)       ***REMOVED*** panic(errors.New("modification of immutable list")) ***REMOVED***
func (x emptyList) AppendMutable() pref.Value ***REMOVED*** panic(errors.New("modification of immutable list")) ***REMOVED***
func (x emptyList) Truncate(n int)            ***REMOVED*** panic(errors.New("modification of immutable list")) ***REMOVED***
func (x emptyList) NewElement() pref.Value    ***REMOVED*** return newListEntry(x.desc) ***REMOVED***
func (x emptyList) IsValid() bool             ***REMOVED*** return false ***REMOVED***

type dynamicList struct ***REMOVED***
	desc pref.FieldDescriptor
	list []pref.Value
***REMOVED***

func (x *dynamicList) Len() int ***REMOVED***
	return len(x.list)
***REMOVED***

func (x *dynamicList) Get(n int) pref.Value ***REMOVED***
	return x.list[n]
***REMOVED***

func (x *dynamicList) Set(n int, v pref.Value) ***REMOVED***
	typecheckSingular(x.desc, v)
	x.list[n] = v
***REMOVED***

func (x *dynamicList) Append(v pref.Value) ***REMOVED***
	typecheckSingular(x.desc, v)
	x.list = append(x.list, v)
***REMOVED***

func (x *dynamicList) AppendMutable() pref.Value ***REMOVED***
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
		x.list[i] = pref.Value***REMOVED******REMOVED***
	***REMOVED***
	x.list = x.list[:n]
***REMOVED***

func (x *dynamicList) NewElement() pref.Value ***REMOVED***
	return newListEntry(x.desc)
***REMOVED***

func (x *dynamicList) IsValid() bool ***REMOVED***
	return true
***REMOVED***

type dynamicMap struct ***REMOVED***
	desc pref.FieldDescriptor
	mapv map[interface***REMOVED******REMOVED***]pref.Value
***REMOVED***

func (x *dynamicMap) Get(k pref.MapKey) pref.Value ***REMOVED*** return x.mapv[k.Interface()] ***REMOVED***
func (x *dynamicMap) Set(k pref.MapKey, v pref.Value) ***REMOVED***
	typecheckSingular(x.desc.MapKey(), k.Value())
	typecheckSingular(x.desc.MapValue(), v)
	x.mapv[k.Interface()] = v
***REMOVED***
func (x *dynamicMap) Has(k pref.MapKey) bool ***REMOVED*** return x.Get(k).IsValid() ***REMOVED***
func (x *dynamicMap) Clear(k pref.MapKey)    ***REMOVED*** delete(x.mapv, k.Interface()) ***REMOVED***
func (x *dynamicMap) Mutable(k pref.MapKey) pref.Value ***REMOVED***
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
func (x *dynamicMap) NewValue() pref.Value ***REMOVED***
	if md := x.desc.MapValue().Message(); md != nil ***REMOVED***
		return pref.ValueOfMessage(NewMessage(md).ProtoReflect())
	***REMOVED***
	return x.desc.MapValue().Default()
***REMOVED***
func (x *dynamicMap) IsValid() bool ***REMOVED***
	return x.mapv != nil
***REMOVED***

func (x *dynamicMap) Range(f func(pref.MapKey, pref.Value) bool) ***REMOVED***
	for k, v := range x.mapv ***REMOVED***
		if !f(pref.ValueOf(k).MapKey(), v) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func isSet(fd pref.FieldDescriptor, v pref.Value) bool ***REMOVED***
	switch ***REMOVED***
	case fd.IsMap():
		return v.Map().Len() > 0
	case fd.IsList():
		return v.List().Len() > 0
	case fd.ContainingOneof() != nil:
		return true
	case fd.Syntax() == pref.Proto3 && !fd.IsExtension():
		switch fd.Kind() ***REMOVED***
		case pref.BoolKind:
			return v.Bool()
		case pref.EnumKind:
			return v.Enum() != 0
		case pref.Int32Kind, pref.Sint32Kind, pref.Int64Kind, pref.Sint64Kind, pref.Sfixed32Kind, pref.Sfixed64Kind:
			return v.Int() != 0
		case pref.Uint32Kind, pref.Uint64Kind, pref.Fixed32Kind, pref.Fixed64Kind:
			return v.Uint() != 0
		case pref.FloatKind, pref.DoubleKind:
			return v.Float() != 0 || math.Signbit(v.Float())
		case pref.StringKind:
			return v.String() != ""
		case pref.BytesKind:
			return len(v.Bytes()) > 0
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func typecheck(fd pref.FieldDescriptor, v pref.Value) ***REMOVED***
	if err := typeIsValid(fd, v); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func typeIsValid(fd pref.FieldDescriptor, v pref.Value) error ***REMOVED***
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

func typecheckSingular(fd pref.FieldDescriptor, v pref.Value) ***REMOVED***
	if err := singularTypeIsValid(fd, v); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func singularTypeIsValid(fd pref.FieldDescriptor, v pref.Value) error ***REMOVED***
	vi := v.Interface()
	var ok bool
	switch fd.Kind() ***REMOVED***
	case pref.BoolKind:
		_, ok = vi.(bool)
	case pref.EnumKind:
		// We could check against the valid set of enum values, but do not.
		_, ok = vi.(pref.EnumNumber)
	case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
		_, ok = vi.(int32)
	case pref.Uint32Kind, pref.Fixed32Kind:
		_, ok = vi.(uint32)
	case pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
		_, ok = vi.(int64)
	case pref.Uint64Kind, pref.Fixed64Kind:
		_, ok = vi.(uint64)
	case pref.FloatKind:
		_, ok = vi.(float32)
	case pref.DoubleKind:
		_, ok = vi.(float64)
	case pref.StringKind:
		_, ok = vi.(string)
	case pref.BytesKind:
		_, ok = vi.([]byte)
	case pref.MessageKind, pref.GroupKind:
		var m pref.Message
		m, ok = vi.(pref.Message)
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

func newListEntry(fd pref.FieldDescriptor) pref.Value ***REMOVED***
	switch fd.Kind() ***REMOVED***
	case pref.BoolKind:
		return pref.ValueOfBool(false)
	case pref.EnumKind:
		return pref.ValueOfEnum(fd.Enum().Values().Get(0).Number())
	case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
		return pref.ValueOfInt32(0)
	case pref.Uint32Kind, pref.Fixed32Kind:
		return pref.ValueOfUint32(0)
	case pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
		return pref.ValueOfInt64(0)
	case pref.Uint64Kind, pref.Fixed64Kind:
		return pref.ValueOfUint64(0)
	case pref.FloatKind:
		return pref.ValueOfFloat32(0)
	case pref.DoubleKind:
		return pref.ValueOfFloat64(0)
	case pref.StringKind:
		return pref.ValueOfString("")
	case pref.BytesKind:
		return pref.ValueOfBytes(nil)
	case pref.MessageKind, pref.GroupKind:
		return pref.ValueOfMessage(NewMessage(fd.Message()).ProtoReflect())
	***REMOVED***
	panic(errors.New("%v: unknown kind %v", fd.FullName(), fd.Kind()))
***REMOVED***

// extensionType is a dynamic protoreflect.ExtensionType.
type extensionType struct ***REMOVED***
	desc extensionTypeDescriptor
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
func NewExtensionType(desc pref.ExtensionDescriptor) pref.ExtensionType ***REMOVED***
	if xt, ok := desc.(pref.ExtensionTypeDescriptor); ok ***REMOVED***
		desc = xt.Descriptor()
	***REMOVED***
	return extensionType***REMOVED***extensionTypeDescriptor***REMOVED***desc***REMOVED******REMOVED***
***REMOVED***

func (xt extensionType) New() pref.Value ***REMOVED***
	switch ***REMOVED***
	case xt.desc.IsMap():
		return pref.ValueOfMap(&dynamicMap***REMOVED***
			desc: xt.desc,
			mapv: make(map[interface***REMOVED******REMOVED***]pref.Value),
		***REMOVED***)
	case xt.desc.IsList():
		return pref.ValueOfList(&dynamicList***REMOVED***desc: xt.desc***REMOVED***)
	case xt.desc.Message() != nil:
		return pref.ValueOfMessage(NewMessage(xt.desc.Message()))
	default:
		return xt.desc.Default()
	***REMOVED***
***REMOVED***

func (xt extensionType) Zero() pref.Value ***REMOVED***
	switch ***REMOVED***
	case xt.desc.IsMap():
		return pref.ValueOfMap(&dynamicMap***REMOVED***desc: xt.desc***REMOVED***)
	case xt.desc.Cardinality() == pref.Repeated:
		return pref.ValueOfList(emptyList***REMOVED***desc: xt.desc***REMOVED***)
	case xt.desc.Message() != nil:
		return pref.ValueOfMessage(&Message***REMOVED***typ: messageType***REMOVED***xt.desc.Message()***REMOVED******REMOVED***)
	default:
		return xt.desc.Default()
	***REMOVED***
***REMOVED***

func (xt extensionType) TypeDescriptor() pref.ExtensionTypeDescriptor ***REMOVED***
	return xt.desc
***REMOVED***

func (xt extensionType) ValueOf(iv interface***REMOVED******REMOVED***) pref.Value ***REMOVED***
	v := pref.ValueOf(iv)
	typecheck(xt.desc, v)
	return v
***REMOVED***

func (xt extensionType) InterfaceOf(v pref.Value) interface***REMOVED******REMOVED*** ***REMOVED***
	typecheck(xt.desc, v)
	return v.Interface()
***REMOVED***

func (xt extensionType) IsValidInterface(iv interface***REMOVED******REMOVED***) bool ***REMOVED***
	return typeIsValid(xt.desc, pref.ValueOf(iv)) == nil
***REMOVED***

func (xt extensionType) IsValidValue(v pref.Value) bool ***REMOVED***
	return typeIsValid(xt.desc, v) == nil
***REMOVED***

type extensionTypeDescriptor struct ***REMOVED***
	pref.ExtensionDescriptor
***REMOVED***

func (xt extensionTypeDescriptor) Type() pref.ExtensionType ***REMOVED***
	return extensionType***REMOVED***xt***REMOVED***
***REMOVED***

func (xt extensionTypeDescriptor) Descriptor() pref.ExtensionDescriptor ***REMOVED***
	return xt.ExtensionDescriptor
***REMOVED***
