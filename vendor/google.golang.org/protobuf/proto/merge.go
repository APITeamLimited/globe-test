// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

// Merge merges src into dst, which must be a message with the same descriptor.
//
// Populated scalar fields in src are copied to dst, while populated
// singular messages in src are merged into dst by recursively calling Merge.
// The elements of every list field in src is appended to the corresponded
// list fields in dst. The entries of every map field in src is copied into
// the corresponding map field in dst, possibly replacing existing entries.
// The unknown fields of src are appended to the unknown fields of dst.
//
// It is semantically equivalent to unmarshaling the encoded form of src
// into dst with the UnmarshalOptions.Merge option specified.
func Merge(dst, src Message) ***REMOVED***
	// TODO: Should nil src be treated as semantically equivalent to a
	// untyped, read-only, empty message? What about a nil dst?

	dstMsg, srcMsg := dst.ProtoReflect(), src.ProtoReflect()
	if dstMsg.Descriptor() != srcMsg.Descriptor() ***REMOVED***
		if got, want := dstMsg.Descriptor().FullName(), srcMsg.Descriptor().FullName(); got != want ***REMOVED***
			panic(fmt.Sprintf("descriptor mismatch: %v != %v", got, want))
		***REMOVED***
		panic("descriptor mismatch")
	***REMOVED***
	mergeOptions***REMOVED******REMOVED***.mergeMessage(dstMsg, srcMsg)
***REMOVED***

// Clone returns a deep copy of m.
// If the top-level message is invalid, it returns an invalid message as well.
func Clone(m Message) Message ***REMOVED***
	// NOTE: Most usages of Clone assume the following properties:
	//	t := reflect.TypeOf(m)
	//	t == reflect.TypeOf(m.ProtoReflect().New().Interface())
	//	t == reflect.TypeOf(m.ProtoReflect().Type().Zero().Interface())
	//
	// Embedding protobuf messages breaks this since the parent type will have
	// a forwarded ProtoReflect method, but the Interface method will return
	// the underlying embedded message type.
	if m == nil ***REMOVED***
		return nil
	***REMOVED***
	src := m.ProtoReflect()
	if !src.IsValid() ***REMOVED***
		return src.Type().Zero().Interface()
	***REMOVED***
	dst := src.New()
	mergeOptions***REMOVED******REMOVED***.mergeMessage(dst, src)
	return dst.Interface()
***REMOVED***

// mergeOptions provides a namespace for merge functions, and can be
// exported in the future if we add user-visible merge options.
type mergeOptions struct***REMOVED******REMOVED***

func (o mergeOptions) mergeMessage(dst, src protoreflect.Message) ***REMOVED***
	methods := protoMethods(dst)
	if methods != nil && methods.Merge != nil ***REMOVED***
		in := protoiface.MergeInput***REMOVED***
			Destination: dst,
			Source:      src,
		***REMOVED***
		out := methods.Merge(in)
		if out.Flags&protoiface.MergeComplete != 0 ***REMOVED***
			return
		***REMOVED***
	***REMOVED***

	if !dst.IsValid() ***REMOVED***
		panic(fmt.Sprintf("cannot merge into invalid %v message", dst.Descriptor().FullName()))
	***REMOVED***

	src.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		switch ***REMOVED***
		case fd.IsList():
			o.mergeList(dst.Mutable(fd).List(), v.List(), fd)
		case fd.IsMap():
			o.mergeMap(dst.Mutable(fd).Map(), v.Map(), fd.MapValue())
		case fd.Message() != nil:
			o.mergeMessage(dst.Mutable(fd).Message(), v.Message())
		case fd.Kind() == protoreflect.BytesKind:
			dst.Set(fd, o.cloneBytes(v))
		default:
			dst.Set(fd, v)
		***REMOVED***
		return true
	***REMOVED***)

	if len(src.GetUnknown()) > 0 ***REMOVED***
		dst.SetUnknown(append(dst.GetUnknown(), src.GetUnknown()...))
	***REMOVED***
***REMOVED***

func (o mergeOptions) mergeList(dst, src protoreflect.List, fd protoreflect.FieldDescriptor) ***REMOVED***
	// Merge semantics appends to the end of the existing list.
	for i, n := 0, src.Len(); i < n; i++ ***REMOVED***
		switch v := src.Get(i); ***REMOVED***
		case fd.Message() != nil:
			dstv := dst.NewElement()
			o.mergeMessage(dstv.Message(), v.Message())
			dst.Append(dstv)
		case fd.Kind() == protoreflect.BytesKind:
			dst.Append(o.cloneBytes(v))
		default:
			dst.Append(v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (o mergeOptions) mergeMap(dst, src protoreflect.Map, fd protoreflect.FieldDescriptor) ***REMOVED***
	// Merge semantics replaces, rather than merges into existing entries.
	src.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool ***REMOVED***
		switch ***REMOVED***
		case fd.Message() != nil:
			dstv := dst.NewValue()
			o.mergeMessage(dstv.Message(), v.Message())
			dst.Set(k, dstv)
		case fd.Kind() == protoreflect.BytesKind:
			dst.Set(k, o.cloneBytes(v))
		default:
			dst.Set(k, v)
		***REMOVED***
		return true
	***REMOVED***)
***REMOVED***

func (o mergeOptions) cloneBytes(v protoreflect.Value) protoreflect.Value ***REMOVED***
	return protoreflect.ValueOfBytes(append([]byte***REMOVED******REMOVED***, v.Bytes()...))
***REMOVED***
