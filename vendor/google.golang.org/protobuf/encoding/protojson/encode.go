// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protojson

import (
	"encoding/base64"
	"fmt"
	"sort"

	"google.golang.org/protobuf/internal/encoding/json"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const defaultIndent = "  "

// Format formats the message as a multiline string.
// This function is only intended for human consumption and ignores errors.
// Do not depend on the output being stable. It may change over time across
// different versions of the program.
func Format(m proto.Message) string ***REMOVED***
	return MarshalOptions***REMOVED***Multiline: true***REMOVED***.Format(m)
***REMOVED***

// Marshal writes the given proto.Message in JSON format using default options.
// Do not depend on the output being stable. It may change over time across
// different versions of the program.
func Marshal(m proto.Message) ([]byte, error) ***REMOVED***
	return MarshalOptions***REMOVED******REMOVED***.Marshal(m)
***REMOVED***

// MarshalOptions is a configurable JSON format marshaler.
type MarshalOptions struct ***REMOVED***
	pragma.NoUnkeyedLiterals

	// Multiline specifies whether the marshaler should format the output in
	// indented-form with every textual element on a new line.
	// If Indent is an empty string, then an arbitrary indent is chosen.
	Multiline bool

	// Indent specifies the set of indentation characters to use in a multiline
	// formatted output such that every entry is preceded by Indent and
	// terminated by a newline. If non-empty, then Multiline is treated as true.
	// Indent can only be composed of space or tab characters.
	Indent string

	// AllowPartial allows messages that have missing required fields to marshal
	// without returning an error. If AllowPartial is false (the default),
	// Marshal will return error if there are any missing required fields.
	AllowPartial bool

	// UseProtoNames uses proto field name instead of lowerCamelCase name in JSON
	// field names.
	UseProtoNames bool

	// UseEnumNumbers emits enum values as numbers.
	UseEnumNumbers bool

	// EmitUnpopulated specifies whether to emit unpopulated fields. It does not
	// emit unpopulated oneof fields or unpopulated extension fields.
	// The JSON value emitted for unpopulated fields are as follows:
	//  ╔═══════╤════════════════════════════╗
	//  ║ JSON  │ Protobuf field             ║
	//  ╠═══════╪════════════════════════════╣
	//  ║ false │ proto3 boolean fields      ║
	//  ║ 0     │ proto3 numeric fields      ║
	//  ║ ""    │ proto3 string/bytes fields ║
	//  ║ null  │ proto2 scalar fields       ║
	//  ║ null  │ message fields             ║
	//  ║ []    │ list fields                ║
	//  ║ ***REMOVED******REMOVED***    │ map fields                 ║
	//  ╚═══════╧════════════════════════════╝
	EmitUnpopulated bool

	// Resolver is used for looking up types when expanding google.protobuf.Any
	// messages. If nil, this defaults to using protoregistry.GlobalTypes.
	Resolver interface ***REMOVED***
		protoregistry.ExtensionTypeResolver
		protoregistry.MessageTypeResolver
	***REMOVED***
***REMOVED***

// Format formats the message as a string.
// This method is only intended for human consumption and ignores errors.
// Do not depend on the output being stable. It may change over time across
// different versions of the program.
func (o MarshalOptions) Format(m proto.Message) string ***REMOVED***
	if m == nil || !m.ProtoReflect().IsValid() ***REMOVED***
		return "<nil>" // invalid syntax, but okay since this is for debugging
	***REMOVED***
	o.AllowPartial = true
	b, _ := o.Marshal(m)
	return string(b)
***REMOVED***

// Marshal marshals the given proto.Message in the JSON format using options in
// MarshalOptions. Do not depend on the output being stable. It may change over
// time across different versions of the program.
func (o MarshalOptions) Marshal(m proto.Message) ([]byte, error) ***REMOVED***
	return o.marshal(m)
***REMOVED***

// marshal is a centralized function that all marshal operations go through.
// For profiling purposes, avoid changing the name of this function or
// introducing other code paths for marshal that do not go through this.
func (o MarshalOptions) marshal(m proto.Message) ([]byte, error) ***REMOVED***
	if o.Multiline && o.Indent == "" ***REMOVED***
		o.Indent = defaultIndent
	***REMOVED***
	if o.Resolver == nil ***REMOVED***
		o.Resolver = protoregistry.GlobalTypes
	***REMOVED***

	internalEnc, err := json.NewEncoder(o.Indent)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Treat nil message interface as an empty message,
	// in which case the output in an empty JSON object.
	if m == nil ***REMOVED***
		return []byte("***REMOVED******REMOVED***"), nil
	***REMOVED***

	enc := encoder***REMOVED***internalEnc, o***REMOVED***
	if err := enc.marshalMessage(m.ProtoReflect()); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if o.AllowPartial ***REMOVED***
		return enc.Bytes(), nil
	***REMOVED***
	return enc.Bytes(), proto.CheckInitialized(m)
***REMOVED***

type encoder struct ***REMOVED***
	*json.Encoder
	opts MarshalOptions
***REMOVED***

// marshalMessage marshals the given protoreflect.Message.
func (e encoder) marshalMessage(m pref.Message) error ***REMOVED***
	if marshal := wellKnownTypeMarshaler(m.Descriptor().FullName()); marshal != nil ***REMOVED***
		return marshal(e, m)
	***REMOVED***

	e.StartObject()
	defer e.EndObject()
	if err := e.marshalFields(m); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

// marshalFields marshals the fields in the given protoreflect.Message.
func (e encoder) marshalFields(m pref.Message) error ***REMOVED***
	messageDesc := m.Descriptor()
	if !flags.ProtoLegacy && messageset.IsMessageSet(messageDesc) ***REMOVED***
		return errors.New("no support for proto1 MessageSets")
	***REMOVED***

	// Marshal out known fields.
	fieldDescs := messageDesc.Fields()
	for i := 0; i < fieldDescs.Len(); ***REMOVED***
		fd := fieldDescs.Get(i)
		if od := fd.ContainingOneof(); od != nil ***REMOVED***
			fd = m.WhichOneof(od)
			i += od.Fields().Len()
			if fd == nil ***REMOVED***
				continue // unpopulated oneofs are not affected by EmitUnpopulated
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			i++
		***REMOVED***

		val := m.Get(fd)
		if !m.Has(fd) ***REMOVED***
			if !e.opts.EmitUnpopulated ***REMOVED***
				continue
			***REMOVED***
			isProto2Scalar := fd.Syntax() == pref.Proto2 && fd.Default().IsValid()
			isSingularMessage := fd.Cardinality() != pref.Repeated && fd.Message() != nil
			if isProto2Scalar || isSingularMessage ***REMOVED***
				// Use invalid value to emit null.
				val = pref.Value***REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***

		name := fd.JSONName()
		if e.opts.UseProtoNames ***REMOVED***
			name = string(fd.Name())
			// Use type name for group field name.
			if fd.Kind() == pref.GroupKind ***REMOVED***
				name = string(fd.Message().Name())
			***REMOVED***
		***REMOVED***
		if err := e.WriteName(name); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := e.marshalValue(val, fd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Marshal out extensions.
	if err := e.marshalExtensions(m); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// marshalValue marshals the given protoreflect.Value.
func (e encoder) marshalValue(val pref.Value, fd pref.FieldDescriptor) error ***REMOVED***
	switch ***REMOVED***
	case fd.IsList():
		return e.marshalList(val.List(), fd)
	case fd.IsMap():
		return e.marshalMap(val.Map(), fd)
	default:
		return e.marshalSingular(val, fd)
	***REMOVED***
***REMOVED***

// marshalSingular marshals the given non-repeated field value. This includes
// all scalar types, enums, messages, and groups.
func (e encoder) marshalSingular(val pref.Value, fd pref.FieldDescriptor) error ***REMOVED***
	if !val.IsValid() ***REMOVED***
		e.WriteNull()
		return nil
	***REMOVED***

	switch kind := fd.Kind(); kind ***REMOVED***
	case pref.BoolKind:
		e.WriteBool(val.Bool())

	case pref.StringKind:
		if e.WriteString(val.String()) != nil ***REMOVED***
			return errors.InvalidUTF8(string(fd.FullName()))
		***REMOVED***

	case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind:
		e.WriteInt(val.Int())

	case pref.Uint32Kind, pref.Fixed32Kind:
		e.WriteUint(val.Uint())

	case pref.Int64Kind, pref.Sint64Kind, pref.Uint64Kind,
		pref.Sfixed64Kind, pref.Fixed64Kind:
		// 64-bit integers are written out as JSON string.
		e.WriteString(val.String())

	case pref.FloatKind:
		// Encoder.WriteFloat handles the special numbers NaN and infinites.
		e.WriteFloat(val.Float(), 32)

	case pref.DoubleKind:
		// Encoder.WriteFloat handles the special numbers NaN and infinites.
		e.WriteFloat(val.Float(), 64)

	case pref.BytesKind:
		e.WriteString(base64.StdEncoding.EncodeToString(val.Bytes()))

	case pref.EnumKind:
		if fd.Enum().FullName() == genid.NullValue_enum_fullname ***REMOVED***
			e.WriteNull()
		***REMOVED*** else ***REMOVED***
			desc := fd.Enum().Values().ByNumber(val.Enum())
			if e.opts.UseEnumNumbers || desc == nil ***REMOVED***
				e.WriteInt(int64(val.Enum()))
			***REMOVED*** else ***REMOVED***
				e.WriteString(string(desc.Name()))
			***REMOVED***
		***REMOVED***

	case pref.MessageKind, pref.GroupKind:
		if err := e.marshalMessage(val.Message()); err != nil ***REMOVED***
			return err
		***REMOVED***

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	***REMOVED***
	return nil
***REMOVED***

// marshalList marshals the given protoreflect.List.
func (e encoder) marshalList(list pref.List, fd pref.FieldDescriptor) error ***REMOVED***
	e.StartArray()
	defer e.EndArray()

	for i := 0; i < list.Len(); i++ ***REMOVED***
		item := list.Get(i)
		if err := e.marshalSingular(item, fd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type mapEntry struct ***REMOVED***
	key   pref.MapKey
	value pref.Value
***REMOVED***

// marshalMap marshals given protoreflect.Map.
func (e encoder) marshalMap(mmap pref.Map, fd pref.FieldDescriptor) error ***REMOVED***
	e.StartObject()
	defer e.EndObject()

	// Get a sorted list based on keyType first.
	entries := make([]mapEntry, 0, mmap.Len())
	mmap.Range(func(key pref.MapKey, val pref.Value) bool ***REMOVED***
		entries = append(entries, mapEntry***REMOVED***key: key, value: val***REMOVED***)
		return true
	***REMOVED***)
	sortMap(fd.MapKey().Kind(), entries)

	// Write out sorted list.
	for _, entry := range entries ***REMOVED***
		if err := e.WriteName(entry.key.String()); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := e.marshalSingular(entry.value, fd.MapValue()); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// sortMap orders list based on value of key field for deterministic ordering.
func sortMap(keyKind pref.Kind, values []mapEntry) ***REMOVED***
	sort.Slice(values, func(i, j int) bool ***REMOVED***
		switch keyKind ***REMOVED***
		case pref.Int32Kind, pref.Sint32Kind, pref.Sfixed32Kind,
			pref.Int64Kind, pref.Sint64Kind, pref.Sfixed64Kind:
			return values[i].key.Int() < values[j].key.Int()

		case pref.Uint32Kind, pref.Fixed32Kind,
			pref.Uint64Kind, pref.Fixed64Kind:
			return values[i].key.Uint() < values[j].key.Uint()
		***REMOVED***
		return values[i].key.String() < values[j].key.String()
	***REMOVED***)
***REMOVED***

// marshalExtensions marshals extension fields.
func (e encoder) marshalExtensions(m pref.Message) error ***REMOVED***
	type entry struct ***REMOVED***
		key   string
		value pref.Value
		desc  pref.FieldDescriptor
	***REMOVED***

	// Get a sorted list based on field key first.
	var entries []entry
	m.Range(func(fd pref.FieldDescriptor, v pref.Value) bool ***REMOVED***
		if !fd.IsExtension() ***REMOVED***
			return true
		***REMOVED***

		// For MessageSet extensions, the name used is the parent message.
		name := fd.FullName()
		if messageset.IsMessageSetExtension(fd) ***REMOVED***
			name = name.Parent()
		***REMOVED***

		// Use [name] format for JSON field name.
		entries = append(entries, entry***REMOVED***
			key:   string(name),
			value: v,
			desc:  fd,
		***REMOVED***)
		return true
	***REMOVED***)

	// Sort extensions lexicographically.
	sort.Slice(entries, func(i, j int) bool ***REMOVED***
		return entries[i].key < entries[j].key
	***REMOVED***)

	// Write out sorted list.
	for _, entry := range entries ***REMOVED***
		// JSON field name is the proto field name enclosed in [], similar to
		// textproto. This is consistent with Go v1 lib. C++ lib v3.7.0 does not
		// marshal out extension fields.
		if err := e.WriteName("[" + entry.key + "]"); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := e.marshalValue(entry.value, entry.desc); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
