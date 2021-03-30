// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/order"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func (o MarshalOptions) sizeMessageSet(m protoreflect.Message) (size int) ***REMOVED***
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		size += messageset.SizeField(fd.Number())
		size += protowire.SizeTag(messageset.FieldMessage)
		size += protowire.SizeBytes(o.size(v.Message()))
		return true
	***REMOVED***)
	size += messageset.SizeUnknown(m.GetUnknown())
	return size
***REMOVED***

func (o MarshalOptions) marshalMessageSet(b []byte, m protoreflect.Message) ([]byte, error) ***REMOVED***
	if !flags.ProtoLegacy ***REMOVED***
		return b, errors.New("no support for message_set_wire_format")
	***REMOVED***
	fieldOrder := order.AnyFieldOrder
	if o.Deterministic ***REMOVED***
		fieldOrder = order.NumberFieldOrder
	***REMOVED***
	var err error
	order.RangeFields(m, fieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		b, err = o.marshalMessageSetField(b, fd, v)
		return err == nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return b, err
	***REMOVED***
	return messageset.AppendUnknown(b, m.GetUnknown())
***REMOVED***

func (o MarshalOptions) marshalMessageSetField(b []byte, fd protoreflect.FieldDescriptor, value protoreflect.Value) ([]byte, error) ***REMOVED***
	b = messageset.AppendFieldStart(b, fd.Number())
	b = protowire.AppendTag(b, messageset.FieldMessage, protowire.BytesType)
	b = protowire.AppendVarint(b, uint64(o.Size(value.Message().Interface())))
	b, err := o.marshalMessage(b, value.Message())
	if err != nil ***REMOVED***
		return b, err
	***REMOVED***
	b = messageset.AppendFieldEnd(b)
	return b, nil
***REMOVED***

func (o UnmarshalOptions) unmarshalMessageSet(b []byte, m protoreflect.Message) error ***REMOVED***
	if !flags.ProtoLegacy ***REMOVED***
		return errors.New("no support for message_set_wire_format")
	***REMOVED***
	return messageset.Unmarshal(b, false, func(num protowire.Number, v []byte) error ***REMOVED***
		err := o.unmarshalMessageSetField(m, num, v)
		if err == errUnknown ***REMOVED***
			unknown := m.GetUnknown()
			unknown = protowire.AppendTag(unknown, num, protowire.BytesType)
			unknown = protowire.AppendBytes(unknown, v)
			m.SetUnknown(unknown)
			return nil
		***REMOVED***
		return err
	***REMOVED***)
***REMOVED***

func (o UnmarshalOptions) unmarshalMessageSetField(m protoreflect.Message, num protowire.Number, v []byte) error ***REMOVED***
	md := m.Descriptor()
	if !md.ExtensionRanges().Has(num) ***REMOVED***
		return errUnknown
	***REMOVED***
	xt, err := o.Resolver.FindExtensionByNumber(md.FullName(), num)
	if err == protoregistry.NotFound ***REMOVED***
		return errUnknown
	***REMOVED***
	if err != nil ***REMOVED***
		return errors.New("%v: unable to resolve extension %v: %v", md.FullName(), num, err)
	***REMOVED***
	xd := xt.TypeDescriptor()
	if err := o.unmarshalMessage(v, m.Mutable(xd).Message()); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***
