// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prototext

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/encoding/text"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/order"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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

// Marshal writes the given proto.Message in textproto format using default
// options. Do not depend on the output being stable. It may change over time
// across different versions of the program.
func Marshal(m proto.Message) ([]byte, error) ***REMOVED***
	return MarshalOptions***REMOVED******REMOVED***.Marshal(m)
***REMOVED***

// MarshalOptions is a configurable text format marshaler.
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

	// EmitASCII specifies whether to format strings and bytes as ASCII only
	// as opposed to using UTF-8 encoding when possible.
	EmitASCII bool

	// allowInvalidUTF8 specifies whether to permit the encoding of strings
	// with invalid UTF-8. This is unexported as it is intended to only
	// be specified by the Format method.
	allowInvalidUTF8 bool

	// AllowPartial allows messages that have missing required fields to marshal
	// without returning an error. If AllowPartial is false (the default),
	// Marshal will return error if there are any missing required fields.
	AllowPartial bool

	// EmitUnknown specifies whether to emit unknown fields in the output.
	// If specified, the unmarshaler may be unable to parse the output.
	// The default is to exclude unknown fields.
	EmitUnknown bool

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
	o.allowInvalidUTF8 = true
	o.AllowPartial = true
	o.EmitUnknown = true
	b, _ := o.Marshal(m)
	return string(b)
***REMOVED***

// Marshal writes the given proto.Message in textproto format using options in
// MarshalOptions object. Do not depend on the output being stable. It may
// change over time across different versions of the program.
func (o MarshalOptions) Marshal(m proto.Message) ([]byte, error) ***REMOVED***
	return o.marshal(m)
***REMOVED***

// marshal is a centralized function that all marshal operations go through.
// For profiling purposes, avoid changing the name of this function or
// introducing other code paths for marshal that do not go through this.
func (o MarshalOptions) marshal(m proto.Message) ([]byte, error) ***REMOVED***
	var delims = [2]byte***REMOVED***'***REMOVED***', '***REMOVED***'***REMOVED***

	if o.Multiline && o.Indent == "" ***REMOVED***
		o.Indent = defaultIndent
	***REMOVED***
	if o.Resolver == nil ***REMOVED***
		o.Resolver = protoregistry.GlobalTypes
	***REMOVED***

	internalEnc, err := text.NewEncoder(o.Indent, delims, o.EmitASCII)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Treat nil message interface as an empty message,
	// in which case there is nothing to output.
	if m == nil ***REMOVED***
		return []byte***REMOVED******REMOVED***, nil
	***REMOVED***

	enc := encoder***REMOVED***internalEnc, o***REMOVED***
	err = enc.marshalMessage(m.ProtoReflect(), false)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	out := enc.Bytes()
	if len(o.Indent) > 0 && len(out) > 0 ***REMOVED***
		out = append(out, '\n')
	***REMOVED***
	if o.AllowPartial ***REMOVED***
		return out, nil
	***REMOVED***
	return out, proto.CheckInitialized(m)
***REMOVED***

type encoder struct ***REMOVED***
	*text.Encoder
	opts MarshalOptions
***REMOVED***

// marshalMessage marshals the given protoreflect.Message.
func (e encoder) marshalMessage(m protoreflect.Message, inclDelims bool) error ***REMOVED***
	messageDesc := m.Descriptor()
	if !flags.ProtoLegacy && messageset.IsMessageSet(messageDesc) ***REMOVED***
		return errors.New("no support for proto1 MessageSets")
	***REMOVED***

	if inclDelims ***REMOVED***
		e.StartMessage()
		defer e.EndMessage()
	***REMOVED***

	// Handle Any expansion.
	if messageDesc.FullName() == genid.Any_message_fullname ***REMOVED***
		if e.marshalAny(m) ***REMOVED***
			return nil
		***REMOVED***
		// If unable to expand, continue on to marshal Any as a regular message.
	***REMOVED***

	// Marshal fields.
	var err error
	order.RangeFields(m, order.IndexNameFieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		if err = e.marshalField(fd.TextName(), v, fd); err != nil ***REMOVED***
			return false
		***REMOVED***
		return true
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Marshal unknown fields.
	if e.opts.EmitUnknown ***REMOVED***
		e.marshalUnknown(m.GetUnknown())
	***REMOVED***

	return nil
***REMOVED***

// marshalField marshals the given field with protoreflect.Value.
func (e encoder) marshalField(name string, val protoreflect.Value, fd protoreflect.FieldDescriptor) error ***REMOVED***
	switch ***REMOVED***
	case fd.IsList():
		return e.marshalList(name, val.List(), fd)
	case fd.IsMap():
		return e.marshalMap(name, val.Map(), fd)
	default:
		e.WriteName(name)
		return e.marshalSingular(val, fd)
	***REMOVED***
***REMOVED***

// marshalSingular marshals the given non-repeated field value. This includes
// all scalar types, enums, messages, and groups.
func (e encoder) marshalSingular(val protoreflect.Value, fd protoreflect.FieldDescriptor) error ***REMOVED***
	kind := fd.Kind()
	switch kind ***REMOVED***
	case protoreflect.BoolKind:
		e.WriteBool(val.Bool())

	case protoreflect.StringKind:
		s := val.String()
		if !e.opts.allowInvalidUTF8 && strs.EnforceUTF8(fd) && !utf8.ValidString(s) ***REMOVED***
			return errors.InvalidUTF8(string(fd.FullName()))
		***REMOVED***
		e.WriteString(s)

	case protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		e.WriteInt(val.Int())

	case protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		e.WriteUint(val.Uint())

	case protoreflect.FloatKind:
		// Encoder.WriteFloat handles the special numbers NaN and infinites.
		e.WriteFloat(val.Float(), 32)

	case protoreflect.DoubleKind:
		// Encoder.WriteFloat handles the special numbers NaN and infinites.
		e.WriteFloat(val.Float(), 64)

	case protoreflect.BytesKind:
		e.WriteString(string(val.Bytes()))

	case protoreflect.EnumKind:
		num := val.Enum()
		if desc := fd.Enum().Values().ByNumber(num); desc != nil ***REMOVED***
			e.WriteLiteral(string(desc.Name()))
		***REMOVED*** else ***REMOVED***
			// Use numeric value if there is no enum description.
			e.WriteInt(int64(num))
		***REMOVED***

	case protoreflect.MessageKind, protoreflect.GroupKind:
		return e.marshalMessage(val.Message(), true)

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	***REMOVED***
	return nil
***REMOVED***

// marshalList marshals the given protoreflect.List as multiple name-value fields.
func (e encoder) marshalList(name string, list protoreflect.List, fd protoreflect.FieldDescriptor) error ***REMOVED***
	size := list.Len()
	for i := 0; i < size; i++ ***REMOVED***
		e.WriteName(name)
		if err := e.marshalSingular(list.Get(i), fd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// marshalMap marshals the given protoreflect.Map as multiple name-value fields.
func (e encoder) marshalMap(name string, mmap protoreflect.Map, fd protoreflect.FieldDescriptor) error ***REMOVED***
	var err error
	order.RangeEntries(mmap, order.GenericKeyOrder, func(key protoreflect.MapKey, val protoreflect.Value) bool ***REMOVED***
		e.WriteName(name)
		e.StartMessage()
		defer e.EndMessage()

		e.WriteName(string(genid.MapEntry_Key_field_name))
		err = e.marshalSingular(key.Value(), fd.MapKey())
		if err != nil ***REMOVED***
			return false
		***REMOVED***

		e.WriteName(string(genid.MapEntry_Value_field_name))
		err = e.marshalSingular(val, fd.MapValue())
		if err != nil ***REMOVED***
			return false
		***REMOVED***
		return true
	***REMOVED***)
	return err
***REMOVED***

// marshalUnknown parses the given []byte and marshals fields out.
// This function assumes proper encoding in the given []byte.
func (e encoder) marshalUnknown(b []byte) ***REMOVED***
	const dec = 10
	const hex = 16
	for len(b) > 0 ***REMOVED***
		num, wtype, n := protowire.ConsumeTag(b)
		b = b[n:]
		e.WriteName(strconv.FormatInt(int64(num), dec))

		switch wtype ***REMOVED***
		case protowire.VarintType:
			var v uint64
			v, n = protowire.ConsumeVarint(b)
			e.WriteUint(v)
		case protowire.Fixed32Type:
			var v uint32
			v, n = protowire.ConsumeFixed32(b)
			e.WriteLiteral("0x" + strconv.FormatUint(uint64(v), hex))
		case protowire.Fixed64Type:
			var v uint64
			v, n = protowire.ConsumeFixed64(b)
			e.WriteLiteral("0x" + strconv.FormatUint(v, hex))
		case protowire.BytesType:
			var v []byte
			v, n = protowire.ConsumeBytes(b)
			e.WriteString(string(v))
		case protowire.StartGroupType:
			e.StartMessage()
			var v []byte
			v, n = protowire.ConsumeGroup(num, b)
			e.marshalUnknown(v)
			e.EndMessage()
		default:
			panic(fmt.Sprintf("prototext: error parsing unknown field wire type: %v", wtype))
		***REMOVED***

		b = b[n:]
	***REMOVED***
***REMOVED***

// marshalAny marshals the given google.protobuf.Any message in expanded form.
// It returns true if it was able to marshal, else false.
func (e encoder) marshalAny(any protoreflect.Message) bool ***REMOVED***
	// Construct the embedded message.
	fds := any.Descriptor().Fields()
	fdType := fds.ByNumber(genid.Any_TypeUrl_field_number)
	typeURL := any.Get(fdType).String()
	mt, err := e.opts.Resolver.FindMessageByURL(typeURL)
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	m := mt.New().Interface()

	// Unmarshal bytes into embedded message.
	fdValue := fds.ByNumber(genid.Any_Value_field_number)
	value := any.Get(fdValue)
	err = proto.UnmarshalOptions***REMOVED***
		AllowPartial: true,
		Resolver:     e.opts.Resolver,
	***REMOVED***.Unmarshal(value.Bytes(), m)
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	// Get current encoder position. If marshaling fails, reset encoder output
	// back to this position.
	pos := e.Snapshot()

	// Field name is the proto field name enclosed in [].
	e.WriteName("[" + typeURL + "]")
	err = e.marshalMessage(m.ProtoReflect(), true)
	if err != nil ***REMOVED***
		e.Reset(pos)
		return false
	***REMOVED***
	return true
***REMOVED***
