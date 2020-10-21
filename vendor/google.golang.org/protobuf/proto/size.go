// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

// Size returns the size in bytes of the wire-format encoding of m.
func Size(m Message) int ***REMOVED***
	return MarshalOptions***REMOVED******REMOVED***.Size(m)
***REMOVED***

// Size returns the size in bytes of the wire-format encoding of m.
func (o MarshalOptions) Size(m Message) int ***REMOVED***
	// Treat a nil message interface as an empty message; nothing to output.
	if m == nil ***REMOVED***
		return 0
	***REMOVED***

	return o.size(m.ProtoReflect())
***REMOVED***

// size is a centralized function that all size operations go through.
// For profiling purposes, avoid changing the name of this function or
// introducing other code paths for size that do not go through this.
func (o MarshalOptions) size(m protoreflect.Message) (size int) ***REMOVED***
	methods := protoMethods(m)
	if methods != nil && methods.Size != nil ***REMOVED***
		out := methods.Size(protoiface.SizeInput***REMOVED***
			Message: m,
		***REMOVED***)
		return out.Size
	***REMOVED***
	if methods != nil && methods.Marshal != nil ***REMOVED***
		// This is not efficient, but we don't have any choice.
		// This case is mainly used for legacy types with a Marshal method.
		out, _ := methods.Marshal(protoiface.MarshalInput***REMOVED***
			Message: m,
		***REMOVED***)
		return len(out.Buf)
	***REMOVED***
	return o.sizeMessageSlow(m)
***REMOVED***

func (o MarshalOptions) sizeMessageSlow(m protoreflect.Message) (size int) ***REMOVED***
	if messageset.IsMessageSet(m.Descriptor()) ***REMOVED***
		return o.sizeMessageSet(m)
	***REMOVED***
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool ***REMOVED***
		size += o.sizeField(fd, v)
		return true
	***REMOVED***)
	size += len(m.GetUnknown())
	return size
***REMOVED***

func (o MarshalOptions) sizeField(fd protoreflect.FieldDescriptor, value protoreflect.Value) (size int) ***REMOVED***
	num := fd.Number()
	switch ***REMOVED***
	case fd.IsList():
		return o.sizeList(num, fd, value.List())
	case fd.IsMap():
		return o.sizeMap(num, fd, value.Map())
	default:
		return protowire.SizeTag(num) + o.sizeSingular(num, fd.Kind(), value)
	***REMOVED***
***REMOVED***

func (o MarshalOptions) sizeList(num protowire.Number, fd protoreflect.FieldDescriptor, list protoreflect.List) (size int) ***REMOVED***
	if fd.IsPacked() && list.Len() > 0 ***REMOVED***
		content := 0
		for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
			content += o.sizeSingular(num, fd.Kind(), list.Get(i))
		***REMOVED***
		return protowire.SizeTag(num) + protowire.SizeBytes(content)
	***REMOVED***

	for i, llen := 0, list.Len(); i < llen; i++ ***REMOVED***
		size += protowire.SizeTag(num) + o.sizeSingular(num, fd.Kind(), list.Get(i))
	***REMOVED***
	return size
***REMOVED***

func (o MarshalOptions) sizeMap(num protowire.Number, fd protoreflect.FieldDescriptor, mapv protoreflect.Map) (size int) ***REMOVED***
	mapv.Range(func(key protoreflect.MapKey, value protoreflect.Value) bool ***REMOVED***
		size += protowire.SizeTag(num)
		size += protowire.SizeBytes(o.sizeField(fd.MapKey(), key.Value()) + o.sizeField(fd.MapValue(), value))
		return true
	***REMOVED***)
	return size
***REMOVED***
