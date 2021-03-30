// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoiface"
)

// UnmarshalOptions configures the unmarshaler.
//
// Example usage:
//   err := UnmarshalOptions***REMOVED***DiscardUnknown: true***REMOVED***.Unmarshal(b, m)
type UnmarshalOptions struct ***REMOVED***
	pragma.NoUnkeyedLiterals

	// Merge merges the input into the destination message.
	// The default behavior is to always reset the message before unmarshaling,
	// unless Merge is specified.
	Merge bool

	// AllowPartial accepts input for messages that will result in missing
	// required fields. If AllowPartial is false (the default), Unmarshal will
	// return an error if there are any missing required fields.
	AllowPartial bool

	// If DiscardUnknown is set, unknown fields are ignored.
	DiscardUnknown bool

	// Resolver is used for looking up types when unmarshaling extension fields.
	// If nil, this defaults to using protoregistry.GlobalTypes.
	Resolver interface ***REMOVED***
		FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error)
		FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error)
	***REMOVED***
***REMOVED***

// Unmarshal parses the wire-format message in b and places the result in m.
func Unmarshal(b []byte, m Message) error ***REMOVED***
	_, err := UnmarshalOptions***REMOVED******REMOVED***.unmarshal(b, m.ProtoReflect())
	return err
***REMOVED***

// Unmarshal parses the wire-format message in b and places the result in m.
func (o UnmarshalOptions) Unmarshal(b []byte, m Message) error ***REMOVED***
	_, err := o.unmarshal(b, m.ProtoReflect())
	return err
***REMOVED***

// UnmarshalState parses a wire-format message and places the result in m.
//
// This method permits fine-grained control over the unmarshaler.
// Most users should use Unmarshal instead.
func (o UnmarshalOptions) UnmarshalState(in protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) ***REMOVED***
	return o.unmarshal(in.Buf, in.Message)
***REMOVED***

// unmarshal is a centralized function that all unmarshal operations go through.
// For profiling purposes, avoid changing the name of this function or
// introducing other code paths for unmarshal that do not go through this.
func (o UnmarshalOptions) unmarshal(b []byte, m protoreflect.Message) (out protoiface.UnmarshalOutput, err error) ***REMOVED***
	if o.Resolver == nil ***REMOVED***
		o.Resolver = protoregistry.GlobalTypes
	***REMOVED***
	if !o.Merge ***REMOVED***
		Reset(m.Interface())
	***REMOVED***
	allowPartial := o.AllowPartial
	o.Merge = true
	o.AllowPartial = true
	methods := protoMethods(m)
	if methods != nil && methods.Unmarshal != nil &&
		!(o.DiscardUnknown && methods.Flags&protoiface.SupportUnmarshalDiscardUnknown == 0) ***REMOVED***
		in := protoiface.UnmarshalInput***REMOVED***
			Message:  m,
			Buf:      b,
			Resolver: o.Resolver,
		***REMOVED***
		if o.DiscardUnknown ***REMOVED***
			in.Flags |= protoiface.UnmarshalDiscardUnknown
		***REMOVED***
		out, err = methods.Unmarshal(in)
	***REMOVED*** else ***REMOVED***
		err = o.unmarshalMessageSlow(b, m)
	***REMOVED***
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	if allowPartial || (out.Flags&protoiface.UnmarshalInitialized != 0) ***REMOVED***
		return out, nil
	***REMOVED***
	return out, checkInitialized(m)
***REMOVED***

func (o UnmarshalOptions) unmarshalMessage(b []byte, m protoreflect.Message) error ***REMOVED***
	_, err := o.unmarshal(b, m)
	return err
***REMOVED***

func (o UnmarshalOptions) unmarshalMessageSlow(b []byte, m protoreflect.Message) error ***REMOVED***
	md := m.Descriptor()
	if messageset.IsMessageSet(md) ***REMOVED***
		return o.unmarshalMessageSet(b, m)
	***REMOVED***
	fields := md.Fields()
	for len(b) > 0 ***REMOVED***
		// Parse the tag (field number and wire type).
		num, wtyp, tagLen := protowire.ConsumeTag(b)
		if tagLen < 0 ***REMOVED***
			return errDecode
		***REMOVED***
		if num > protowire.MaxValidNumber ***REMOVED***
			return errDecode
		***REMOVED***

		// Find the field descriptor for this field number.
		fd := fields.ByNumber(num)
		if fd == nil && md.ExtensionRanges().Has(num) ***REMOVED***
			extType, err := o.Resolver.FindExtensionByNumber(md.FullName(), num)
			if err != nil && err != protoregistry.NotFound ***REMOVED***
				return errors.New("%v: unable to resolve extension %v: %v", md.FullName(), num, err)
			***REMOVED***
			if extType != nil ***REMOVED***
				fd = extType.TypeDescriptor()
			***REMOVED***
		***REMOVED***
		var err error
		if fd == nil ***REMOVED***
			err = errUnknown
		***REMOVED*** else if flags.ProtoLegacy ***REMOVED***
			if fd.IsWeak() && fd.Message().IsPlaceholder() ***REMOVED***
				err = errUnknown // weak referent is not linked in
			***REMOVED***
		***REMOVED***

		// Parse the field value.
		var valLen int
		switch ***REMOVED***
		case err != nil:
		case fd.IsList():
			valLen, err = o.unmarshalList(b[tagLen:], wtyp, m.Mutable(fd).List(), fd)
		case fd.IsMap():
			valLen, err = o.unmarshalMap(b[tagLen:], wtyp, m.Mutable(fd).Map(), fd)
		default:
			valLen, err = o.unmarshalSingular(b[tagLen:], wtyp, m, fd)
		***REMOVED***
		if err != nil ***REMOVED***
			if err != errUnknown ***REMOVED***
				return err
			***REMOVED***
			valLen = protowire.ConsumeFieldValue(num, wtyp, b[tagLen:])
			if valLen < 0 ***REMOVED***
				return errDecode
			***REMOVED***
			if !o.DiscardUnknown ***REMOVED***
				m.SetUnknown(append(m.GetUnknown(), b[:tagLen+valLen]...))
			***REMOVED***
		***REMOVED***
		b = b[tagLen+valLen:]
	***REMOVED***
	return nil
***REMOVED***

func (o UnmarshalOptions) unmarshalSingular(b []byte, wtyp protowire.Type, m protoreflect.Message, fd protoreflect.FieldDescriptor) (n int, err error) ***REMOVED***
	v, n, err := o.unmarshalScalar(b, wtyp, fd)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	switch fd.Kind() ***REMOVED***
	case protoreflect.GroupKind, protoreflect.MessageKind:
		m2 := m.Mutable(fd).Message()
		if err := o.unmarshalMessage(v.Bytes(), m2); err != nil ***REMOVED***
			return n, err
		***REMOVED***
	default:
		// Non-message scalars replace the previous value.
		m.Set(fd, v)
	***REMOVED***
	return n, nil
***REMOVED***

func (o UnmarshalOptions) unmarshalMap(b []byte, wtyp protowire.Type, mapv protoreflect.Map, fd protoreflect.FieldDescriptor) (n int, err error) ***REMOVED***
	if wtyp != protowire.BytesType ***REMOVED***
		return 0, errUnknown
	***REMOVED***
	b, n = protowire.ConsumeBytes(b)
	if n < 0 ***REMOVED***
		return 0, errDecode
	***REMOVED***
	var (
		keyField = fd.MapKey()
		valField = fd.MapValue()
		key      protoreflect.Value
		val      protoreflect.Value
		haveKey  bool
		haveVal  bool
	)
	switch valField.Kind() ***REMOVED***
	case protoreflect.GroupKind, protoreflect.MessageKind:
		val = mapv.NewValue()
	***REMOVED***
	// Map entries are represented as a two-element message with fields
	// containing the key and value.
	for len(b) > 0 ***REMOVED***
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 ***REMOVED***
			return 0, errDecode
		***REMOVED***
		if num > protowire.MaxValidNumber ***REMOVED***
			return 0, errDecode
		***REMOVED***
		b = b[n:]
		err = errUnknown
		switch num ***REMOVED***
		case genid.MapEntry_Key_field_number:
			key, n, err = o.unmarshalScalar(b, wtyp, keyField)
			if err != nil ***REMOVED***
				break
			***REMOVED***
			haveKey = true
		case genid.MapEntry_Value_field_number:
			var v protoreflect.Value
			v, n, err = o.unmarshalScalar(b, wtyp, valField)
			if err != nil ***REMOVED***
				break
			***REMOVED***
			switch valField.Kind() ***REMOVED***
			case protoreflect.GroupKind, protoreflect.MessageKind:
				if err := o.unmarshalMessage(v.Bytes(), val.Message()); err != nil ***REMOVED***
					return 0, err
				***REMOVED***
			default:
				val = v
			***REMOVED***
			haveVal = true
		***REMOVED***
		if err == errUnknown ***REMOVED***
			n = protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 ***REMOVED***
				return 0, errDecode
			***REMOVED***
		***REMOVED*** else if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		b = b[n:]
	***REMOVED***
	// Every map entry should have entries for key and value, but this is not strictly required.
	if !haveKey ***REMOVED***
		key = keyField.Default()
	***REMOVED***
	if !haveVal ***REMOVED***
		switch valField.Kind() ***REMOVED***
		case protoreflect.GroupKind, protoreflect.MessageKind:
		default:
			val = valField.Default()
		***REMOVED***
	***REMOVED***
	mapv.Set(key.MapKey(), val)
	return n, nil
***REMOVED***

// errUnknown is used internally to indicate fields which should be added
// to the unknown field set of a message. It is never returned from an exported
// function.
var errUnknown = errors.New("BUG: internal error (unknown)")

var errDecode = errors.New("cannot parse invalid wire-format data")
