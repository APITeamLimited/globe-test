// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/order"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

// MarshalOptions configures the marshaler.
//
// Example usage:
//
//	b, err := MarshalOptions***REMOVED***Deterministic: true***REMOVED***.Marshal(m)
type MarshalOptions struct ***REMOVED***
	pragma.NoUnkeyedLiterals

	// AllowPartial allows messages that have missing required fields to marshal
	// without returning an error. If AllowPartial is false (the default),
	// Marshal will return an error if there are any missing required fields.
	AllowPartial bool

	// Deterministic controls whether the same message will always be
	// serialized to the same bytes within the same binary.
	//
	// Setting this option guarantees that repeated serialization of
	// the same message will return the same bytes, and that different
	// processes of the same binary (which may be executing on different
	// machines) will serialize equal messages to the same bytes.
	// It has no effect on the resulting size of the encoded message compared
	// to a non-deterministic marshal.
	//
	// Note that the deterministic serialization is NOT canonical across
	// languages. It is not guaranteed to remain stable over time. It is
	// unstable across different builds with schema changes due to unknown
	// fields. Users who need canonical serialization (e.g., persistent
	// storage in a canonical form, fingerprinting, etc.) must define
	// their own canonicalization specification and implement their own
	// serializer rather than relying on this API.
	//
	// If deterministic serialization is requested, map entries will be
	// sorted by keys in lexographical order. This is an implementation
	// detail and subject to change.
	Deterministic bool

	// UseCachedSize indicates that the result of a previous Size call
	// may be reused.
	//
	// Setting this option asserts that:
	//
	// 1. Size has previously been called on this message with identical
	// options (except for UseCachedSize itself).
	//
	// 2. The message and all its submessages have not changed in any
	// way since the Size call.
	//
	// If either of these invariants is violated,
	// the results are undefined and may include panics or corrupted output.
	//
	// Implementations MAY take this option into account to provide
	// better performance, but there is no guarantee that they will do so.
	// There is absolutely no guarantee that Size followed by Marshal with
	// UseCachedSize set will perform equivalently to Marshal alone.
	UseCachedSize bool
***REMOVED***

// Marshal returns the wire-format encoding of m.
func Marshal(m Message) ([]byte, error) ***REMOVED***
	// Treat nil message interface as an empty message; nothing to output.
	if m == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	out, err := MarshalOptions***REMOVED******REMOVED***.marshal(nil, m.ProtoReflect())
	if len(out.Buf) == 0 && err == nil ***REMOVED***
		out.Buf = emptyBytesForMessage(m)
	***REMOVED***
	return out.Buf, err
***REMOVED***

// Marshal returns the wire-format encoding of m.
func (o MarshalOptions) Marshal(m Message) ([]byte, error) ***REMOVED***
	// Treat nil message interface as an empty message; nothing to output.
	if m == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	out, err := o.marshal(nil, m.ProtoReflect())
	if len(out.Buf) == 0 && err == nil ***REMOVED***
		out.Buf = emptyBytesForMessage(m)
	***REMOVED***
	return out.Buf, err
***REMOVED***

// emptyBytesForMessage returns a nil buffer if and only if m is invalid,
// otherwise it returns a non-nil empty buffer.
//
// This is to assist the edge-case where user-code does the following:
//
//	m1.OptionalBytes, _ = proto.Marshal(m2)
//
// where they expect the proto2 "optional_bytes" field to be populated
// if any only if m2 is a valid message.
func emptyBytesForMessage(m Message) []byte ***REMOVED***
	if m == nil || !m.ProtoReflect().IsValid() ***REMOVED***
		return nil
	***REMOVED***
	return emptyBuf[:]
***REMOVED***

// MarshalAppend appends the wire-format encoding of m to b,
// returning the result.
func (o MarshalOptions) MarshalAppend(b []byte, m Message) ([]byte, error) ***REMOVED***
	// Treat nil message interface as an empty message; nothing to append.
	if m == nil ***REMOVED***
		return b, nil
	***REMOVED***

	out, err := o.marshal(b, m.ProtoReflect())
	return out.Buf, err
***REMOVED***

// MarshalState returns the wire-format encoding of a message.
//
// This method permits fine-grained control over the marshaler.
// Most users should use Marshal instead.
func (o MarshalOptions) MarshalState(in protoiface.MarshalInput) (protoiface.MarshalOutput, error) ***REMOVED***
	return o.marshal(in.Buf, in.Message)
***REMOVED***

// marshal is a centralized function that all marshal operations go through.
// For profiling purposes, avoid changing the name of this function or
// introducing other code paths for marshal that do not go through this.
func (o MarshalOptions) marshal(b []byte, m protoreflect.Message) (out protoiface.MarshalOutput, err error) ***REMOVED***
	allowPartial := o.AllowPartial
	o.AllowPartial = true
	if methods := protoMethods(m); methods != nil && methods.Marshal != nil &&
		!(o.Deterministic && methods.Flags&protoiface.SupportMarshalDeterministic == 0) ***REMOVED***
		in := protoiface.MarshalInput***REMOVED***
			Message: m,
			Buf:     b,
		***REMOVED***
		if o.Deterministic ***REMOVED***
			in.Flags |= protoiface.MarshalDeterministic
		***REMOVED***
		if o.UseCachedSize ***REMOVED***
			in.Flags |= protoiface.MarshalUseCachedSize
		***REMOVED***
		if methods.Size != nil ***REMOVED***
			sout := methods.Size(protoiface.SizeInput***REMOVED***
				Message: m,
				Flags:   in.Flags,
			***REMOVED***)
			if cap(b) < len(b)+sout.Size ***REMOVED***
				in.Buf = make([]byte, len(b), growcap(cap(b), len(b)+sout.Size))
				copy(in.Buf, b)
			***REMOVED***
			in.Flags |= protoiface.MarshalUseCachedSize
		***REMOVED***
		out, err = methods.Marshal(in)
	***REMOVED*** else ***REMOVED***
		out.Buf, err = o.marshalMessageSlow(b, m)
	***REMOVED***
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	if allowPartial ***REMOVED***
		return out, nil
	***REMOVED***
	return out, checkInitialized(m)
***REMOVED***

func (o MarshalOptions) marshalMessage(b []byte, m protoreflect.Message) ([]byte, error) ***REMOVED***
	out, err := o.marshal(b, m)
	return out.Buf, err
***REMOVED***

// growcap scales up the capacity of a slice.
//
// Given a slice with a current capacity of oldcap and a desired
// capacity of wantcap, growcap returns a new capacity >= wantcap.
//
// The algorithm is mostly identical to the one used by append as of Go 1.14.
func growcap(oldcap, wantcap int) (newcap int) ***REMOVED***
	if wantcap > oldcap*2 ***REMOVED***
		newcap = wantcap
	***REMOVED*** else if oldcap < 1024 ***REMOVED***
		// The Go 1.14 runtime takes this case when len(s) < 1024,
		// not when cap(s) < 1024. The difference doesn't seem
		// significant here.
		newcap = oldcap * 2
	***REMOVED*** else ***REMOVED***
		newcap = oldcap
		for 0 < newcap && newcap < wantcap ***REMOVED***
			newcap += newcap / 4
		***REMOVED***
		if newcap <= 0 ***REMOVED***
			newcap = wantcap
		***REMOVED***
	***REMOVED***
	return newcap
***REMOVED***

func (o MarshalOptions) marshalMessageSlow(b []byte, m protoreflect.Message) ([]byte, error) ***REMOVED***
	if messageset.IsMessageSet(m.Descriptor()) ***REMOVED***
		return o.marshalMessageSet(b, m)
	***REMOVED***
	fieldOrder := order.AnyFieldOrder
	if o.Deterministic ***REMOVED***
		// TODO: This should use a more natural ordering like NumberFieldOrder,
		// but doing so breaks golden tests that make invalid assumption about
		// output stability of this implementation.
		fieldOrder = order.LegacyFieldOrder
	***REMOVED***
	var err error
	order.RangeFields(m, fieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		b, err = o.marshalField(b, fd, v)
		return err == nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return b, err
	***REMOVED***
	b = append(b, m.GetUnknown()...)
	return b, nil
***REMOVED***

func (o MarshalOptions) marshalField(b []byte, fd protoreflect.FieldDescriptor, value protoreflect.Value) ([]byte, error) ***REMOVED***
	switch ***REMOVED***
	case fd.IsList():
		return o.marshalList(b, fd, value.List())
	case fd.IsMap():
		return o.marshalMap(b, fd, value.Map())
	default:
		b = protowire.AppendTag(b, fd.Number(), wireTypes[fd.Kind()])
		return o.marshalSingular(b, fd, value)
	***REMOVED***
***REMOVED***

func (o MarshalOptions) marshalList(b []byte, fd protoreflect.FieldDescriptor, list protoreflect.List) ([]byte, error) ***REMOVED***
	if fd.IsPacked() && list.Len() > 0 ***REMOVED***
		b = protowire.AppendTag(b, fd.Number(), protowire.BytesType)
		b, pos := appendSpeculativeLength(b)
		for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
			var err error
			b, err = o.marshalSingular(b, fd, list.Get(i))
			if err != nil ***REMOVED***
				return b, err
			***REMOVED***
		***REMOVED***
		b = finishSpeculativeLength(b, pos)
		return b, nil
	***REMOVED***

	kind := fd.Kind()
	for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
		var err error
		b = protowire.AppendTag(b, fd.Number(), wireTypes[kind])
		b, err = o.marshalSingular(b, fd, list.Get(i))
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func (o MarshalOptions) marshalMap(b []byte, fd protoreflect.FieldDescriptor, mapv protoreflect.Map) ([]byte, error) ***REMOVED***
	keyf := fd.MapKey()
	valf := fd.MapValue()
	keyOrder := order.AnyKeyOrder
	if o.Deterministic ***REMOVED***
		keyOrder = order.GenericKeyOrder
	***REMOVED***
	var err error
	order.RangeEntries(mapv, keyOrder, func(key protoreflect.MapKey, value protoreflect.Value) bool ***REMOVED***
		b = protowire.AppendTag(b, fd.Number(), protowire.BytesType)
		var pos int
		b, pos = appendSpeculativeLength(b)

		b, err = o.marshalField(b, keyf, key.Value())
		if err != nil ***REMOVED***
			return false
		***REMOVED***
		b, err = o.marshalField(b, valf, value)
		if err != nil ***REMOVED***
			return false
		***REMOVED***
		b = finishSpeculativeLength(b, pos)
		return true
	***REMOVED***)
	return b, err
***REMOVED***

// When encoding length-prefixed fields, we speculatively set aside some number of bytes
// for the length, encode the data, and then encode the length (shifting the data if necessary
// to make room).
const speculativeLength = 1

func appendSpeculativeLength(b []byte) ([]byte, int) ***REMOVED***
	pos := len(b)
	b = append(b, "\x00\x00\x00\x00"[:speculativeLength]...)
	return b, pos
***REMOVED***

func finishSpeculativeLength(b []byte, pos int) []byte ***REMOVED***
	mlen := len(b) - pos - speculativeLength
	msiz := protowire.SizeVarint(uint64(mlen))
	if msiz != speculativeLength ***REMOVED***
		for i := 0; i < msiz-speculativeLength; i++ ***REMOVED***
			b = append(b, 0)
		***REMOVED***
		copy(b[pos+msiz:], b[pos+speculativeLength:])
		b = b[:pos+msiz+mlen]
	***REMOVED***
	protowire.AppendVarint(b[:pos], uint64(mlen))
	return b
***REMOVED***
